package pool

import (
	"BruteForce_SearchEnginer/common/model"
	"sync"
)

type DirectoryPool struct {
	mu          sync.Mutex
	directories []model.DirectoryResponse
}

func NewDirectoryPool() *DirectoryPool {
	return &DirectoryPool{
		mu:          sync.Mutex{},
		directories: make([]model.DirectoryResponse, 0),
	}
}

// AddDirectory adds a new directory to the pool
func (pool *DirectoryPool) AddDirectory(directory model.DirectoryResponse) {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	pool.directories = append(pool.directories, directory)
}

// PopDirectory removes and returns a directory from the pool (FIFO)
func (pool *DirectoryPool) PopDirectory() (model.DirectoryResponse, bool) {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	if len(pool.directories) == 0 {
		return model.DirectoryResponse{}, false
	}

	dir := pool.directories[0]
	pool.directories = pool.directories[1:]
	return dir, true
}

// GetAllDirectories returns all the directories in the pool
func (pool *DirectoryPool) GetAllDirectories() []model.DirectoryResponse {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	if len(pool.directories) == 0 {
		return nil
	}

	return pool.directories
}

// Length returns the number of directories in the pool
func (pool *DirectoryPool) Length() int64 {
	pool.mu.Lock()
	defer pool.mu.Unlock()
	return int64(len(pool.directories))
}
