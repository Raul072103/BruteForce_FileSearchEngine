package crawler

import (
	"BruteForce_SearchEnginer/common/model"
	"BruteForce_SearchEnginer/worker/internal/matcher"
	"BruteForce_SearchEnginer/worker/internal/repo"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// Crawler is the basis for the component which crawls every file and directory starting from the root path.
// Logs any errors happening throughout the process and jumps over files specified in the configuration.
type Crawler interface {
	Run(ctx context.Context, directoryPath string, request model.SearchRequest)
}

type crawler struct {
	id             int64
	fileRepo       repo.FileRepo
	requestMatcher matcher.RequestMatcher
	config         Config
	logger         *zap.Logger
	typeMap        model.FileTypesConfig
}

type Config struct {
	DirectoryPoolEndpoint string
	ResultsPoolEndpoint   string
}

func New(
	id int64,
	fileRepo repo.FileRepo,
	requestMatcher matcher.RequestMatcher,
	logger *zap.Logger,
	config Config,
	typeMap model.FileTypesConfig) Crawler {
	return &crawler{
		id:             id,
		fileRepo:       fileRepo,
		requestMatcher: requestMatcher,
		config:         config,
		logger:         logger,
		typeMap:        typeMap,
	}
}

func (c *crawler) Run(ctx context.Context, directoryPath string, request model.SearchRequest) {
	c.logger.Info("Starting crawler", zap.String("root", directoryPath), zap.Int64("worker_id", c.id))
	c.crawl(ctx, directoryPath, request)
}

// crawl goes through the current folder in a DFS manner, every folder except the last one is sent back to the
// pool of directories for the current worker or the others to take on when they have resources.
func (c *crawler) crawl(ctx context.Context, path string, request model.SearchRequest) {
	select {
	case <-ctx.Done():
		c.logger.Info("Crawler stopped before going further", zap.String("path", path))
		return
	default:
		var directories = make([]string, 0)

		entries, err := os.ReadDir(path)

		if err != nil {
			c.logger.Error(
				"Error reading directory for further traversing",
				zap.String("path", path),
				zap.Error(err),
				zap.Int64("worker_id", c.id))

			return
		}

		// Look for the results in each file
		for _, entry := range entries {
			entryPath := filepath.Join(path, entry.Name())
			fileMetadata, err := c.fileRepo.Read(entryPath)
			if err != nil {
				c.logger.Info("Error reading file or dir",
					zap.String("path", path),
					zap.Error(err),
					zap.Int64("worker_id", c.id))
				return
			}

			if fileMetadata.Extension == "" {
				directories = append(directories, entryPath)
			} else {
				matchesRequest := c.requestMatcher.MatchFile(fileMetadata, request)
				if matchesRequest {
					go func() {
						var fileSearchResponse model.FileSearchResponse

						if c.typeMap.GetTypeByExtension(fileMetadata.Extension) == ".txt" {
							textContent := string(fileMetadata.Content)
							preview := textContent[:min(len(textContent), 200)]

							fileSearchResponse = model.ConvertToResponse(fileMetadata, preview)
						} else {
							fileSearchResponse = model.ConvertToResponse(fileMetadata, "")
						}

						err := c.sendFileToResultsPool(fileSearchResponse)
						if err != nil {
							c.logger.Error(
								"error sending result to the result pool",
								zap.Error(err),
								zap.Int64("worker_id", c.id))
						}
					}()
				}
			}

		}

		// Send every directory path back to the directories pool, besides the last one
		for i, dir := range directories {
			if i < len(entries)-1 {
				go func() {
					directoryNetworkResponse := model.DirectoryResponse{
						Path:          dir,
						SearchRequest: model.ConvertSearchRequest(request),
					}
					err := c.sendDirectoryToPool(directoryNetworkResponse)
					if err != nil {
						c.logger.Error("Error sending directory to directory pool", zap.Error(err), zap.Int64("worker_id", c.id))
					}
				}()
			} else {
				// go further the last one
				c.crawl(ctx, dir, request)
			}
		}

	}
}

// sendDirectoryToPool sends a directory to the directory pool endpoint
func (c *crawler) sendDirectoryToPool(directoryNetworkResponse model.DirectoryResponse) error {
	jsonData, err := json.Marshal(directoryNetworkResponse)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Post(c.config.DirectoryPoolEndpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send directory to pool, status code: %d", resp.StatusCode)
	}

	return nil
}

// sendFileToResultsPool sends a file search result to the results pool endpoint
func (c *crawler) sendFileToResultsPool(fileSearchResponse model.FileSearchResponse) error {
	jsonData, err := json.Marshal(fileSearchResponse)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Post(c.config.ResultsPoolEndpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send file to results pool, status code: %d", resp.StatusCode)
	}

	return nil
}
