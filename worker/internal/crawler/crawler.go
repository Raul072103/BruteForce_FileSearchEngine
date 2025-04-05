package crawler

import (
	"BruteForce_SearchEnginer/common/model"
	"BruteForce_SearchEnginer/worker/internal/matcher"
	"BruteForce_SearchEnginer/worker/internal/repo"
	"context"
	"go.uber.org/zap"
	"os"
	"path/filepath"
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
	logger         *zap.Logger
}

func New(id int64, fileRepo repo.FileRepo, requestMatcher matcher.RequestMatcher, logger *zap.Logger) Crawler {
	return &crawler{
		id:             id,
		fileRepo:       fileRepo,
		requestMatcher: requestMatcher,
		logger:         logger,
	}
}

func (c *crawler) Run(ctx context.Context, directoryPath string, request model.SearchRequest) {
	c.logger.Info("Starting crawler", zap.String("root", directoryPath), zap.Int64("worker_id", c.id))
	c.crawl(ctx, directoryPath, request)
	c.logger.Info("Crawler finished")
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
					// TODO() send the file to the results pool
				}
			}

		}

		// Send every directory path back to the directories pool, besides the last one
		for i, dir := range directories {
			if i < len(entries)-1 {
				// TODO() send every one of the directory to the pool of directories
			} else {
				// go further the last one
				c.crawl(ctx, dir, request)
			}
		}

	}
}
