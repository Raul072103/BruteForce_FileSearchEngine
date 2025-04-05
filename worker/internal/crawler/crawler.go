package crawler

import (
	"BruteForce_SearchEnginer/worker/internal/repo"
	"context"
	"go.uber.org/zap"
	"os"
	"path/filepath"
)

// Crawler is the basis for the component which crawls every file and directory starting from the root path.
// Logs any errors happening throughout the process and jumps over files specified in the configuration.
type Crawler interface {
	Run(ctx context.Context, directoryPath string)
}

type crawler struct {
	id       int64
	fileRepo repo.FileRepo
	logger   *zap.Logger
}

func New(fileRepo repo.FileRepo, logger *zap.Logger) Crawler {
	return &crawler{
		fileRepo: fileRepo,
		logger:   logger,
	}
}

func (c *crawler) Run(ctx context.Context, directoryPath string) {
	c.logger.Info("Starting crawler", zap.String("root", directoryPath), zap.Int64("worker_id", c.id))
	c.crawl(ctx, directoryPath)
	c.logger.Info("Crawler finished")
}

// crawl goes through the current folder in a DFS manner, every folder except the last one is sent back to the
// pool of directories for the current worker or the others to take on when they have resources.
func (c *crawler) crawl(ctx context.Context, path string) {
	select {
	case <-ctx.Done():
		c.logger.Info("Crawler stopped before going further", zap.String("path", path))
		return
	default:
		fileModel, err := c.fileRepo.Read(path)
		if err != nil {
			c.logger.Info("Error reading file or dir", zap.String("path", path), zap.Error(err))
			return
		}

		if fileModel.Extension == "" {
			entries, err := os.ReadDir(path)
			if err != nil {
				c.logger.Error(
					"Error reading directory for further traversing",
					zap.String("path", path),
					zap.Error(err))

				return
			}

			// Recur for each entry
			for i, entry := range entries {
				entryPath := filepath.Join(path, entry.Name())
				// send every one of the directory to the pool of directories
				if i < len(entries)-1 {

				} else {
					// go further the last one
					c.crawl(ctx, entryPath)
				}
			}
		}
	}
}
