package worker_manager

import (
	"errors"
	"go.uber.org/zap"
	"sync"
	"time"
)

var (
	ErrWorkerNotRegistered     = errors.New("tried accessing an unregistered worker")
	ErrWorkerAlreadyRegistered = errors.New("tried creating an worker that already exists")
)

type WorkerManager struct {
	workerIDs     map[int64]any
	activeWorkers map[int64]WorkerStatus
	mutex         *sync.Mutex
	logger        *zap.Logger
}

type WorkerStatus struct {
	ID             int64
	RegisteredTime time.Time
	Active         bool
}

func New(logger *zap.Logger) *WorkerManager {
	mutex := &sync.Mutex{}

	return &WorkerManager{
		workerIDs:     make(map[int64]any),
		activeWorkers: make(map[int64]WorkerStatus),
		mutex:         mutex,
		logger:        logger,
	}
}

func (manager *WorkerManager) AddWorker(id int64) error {
	worker := WorkerStatus{
		ID:             id,
		RegisteredTime: time.Now(),
		Active:         true,
	}

	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	_, exists := manager.activeWorkers[id]
	if exists {
		return ErrWorkerAlreadyRegistered
	}

	manager.workerIDs[id] = struct{}{}
	manager.activeWorkers[id] = worker

	manager.logger.Info(
		"worker started",
		zap.Int64("worker_id", worker.ID),
		zap.Time("time_registered", worker.RegisteredTime))

	return nil
}

func (manager *WorkerManager) RemoveWorker(id int64) error {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	worker, exists := manager.activeWorkers[id]
	if !exists {
		return ErrWorkerNotRegistered
	}

	delete(manager.activeWorkers, id)

	manager.logger.Info(
		"worker ended",
		zap.Int64("worker_id", worker.ID),
		zap.Time("time_unregistered", time.Now()))

	return nil
}

func (manager *WorkerManager) NoOfWorkers() int {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	return len(manager.activeWorkers)
}
