package process_manager

import (
	"BruteForce_SearchEnginer/manager/internal/pool"
	"BruteForce_SearchEnginer/manager/internal/worker_manager"
	"context"
	"fmt"
	"go.uber.org/zap"
	"os"
	"os/exec"
)

const MaximumNrWorkers = 30

type ProcessManager struct {
	workerManager *worker_manager.WorkerManager
	directoryPool *pool.DirectoryPool
	resultPool    *pool.ResultPool
	logger        *zap.Logger
}

func New(
	workerManager *worker_manager.WorkerManager,
	directoryPool *pool.DirectoryPool,
	resultPool *pool.ResultPool,
	logger *zap.Logger) *ProcessManager {

	return &ProcessManager{
		workerManager: workerManager,
		directoryPool: directoryPool,
		resultPool:    resultPool,
		logger:        logger}
}

func (processManager *ProcessManager) Run(ctx context.Context) {
	filePath := "./../workers_progress.txt"
	file, err := os.Create(filePath)
	if err != nil {
		processManager.logger.Error("Failed to create file", zap.Error(err))
		return
	}
	defer file.Close()

	for {
		select {
		case <-ctx.Done():
			processManager.logger.Info("process manager has stopped")
			return
		default:
			currNoWorkers := processManager.workerManager.NoOfWorkers()
			currNoSearchRequests := processManager.directoryPool.Length()

			if currNoSearchRequests > 0 && currNoWorkers < MaximumNrWorkers && currNoWorkers < int(currNoSearchRequests) {
				cmd := exec.Command("./../worker/worker.exe")
				err := cmd.Start()
				if err != nil {
					processManager.logger.Error("error starting process", zap.Error(err))
				} else {
					processManager.logger.Info("started process", zap.Int("PID", cmd.Process.Pid))
				}
			}

			updateWorkersFile(filePath, currNoWorkers)
		}
	}
}

// updateWorkersFile writes the ASCII visualization to the file
func updateWorkersFile(filePath string, numWorkers int) {
	asciiArt := generateBarChart(numWorkers, MaximumNrWorkers)

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Failed to open file:", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(asciiArt)
	if err != nil {
		fmt.Println("Failed to write to file:", err)
	}
}

// generateBarChart creates ASCII visualization
func generateBarChart(curr, max int) string {
	bar := "["
	for i := 0; i < max; i++ {
		if i < curr {
			bar += "[::]"
		} else {
			bar += "----"
		}
	}
	bar += "]"

	return fmt.Sprintf("Workers: %d/%d\n%s\n", curr, max, bar)
}
