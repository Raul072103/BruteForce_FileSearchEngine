package pool

import (
	"BruteForce_SearchEnginer/common/model"
	"sync"
)

type ResultPool struct {
	mu      sync.Mutex
	results []model.FileSearchResponse
}

func NewResultPool() *ResultPool {
	return &ResultPool{
		mu:      sync.Mutex{},
		results: make([]model.FileSearchResponse, 0),
	}
}

// AddResult adds a new search result to the pool
func (pool *ResultPool) AddResult(result model.FileSearchResponse) {
	pool.mu.Lock()
	defer pool.mu.Unlock()
	pool.results = append(pool.results, result)
}

// PopResult removes and returns a search result from the pool (FIFO)
func (pool *ResultPool) PopResult() (model.FileSearchResponse, bool) {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	if len(pool.results) == 0 {
		return model.FileSearchResponse{}, false
	}

	result := pool.results[0]
	pool.results = pool.results[1:] // Remove first element
	return result, true
}

// Length returns the number of results in the pool
func (pool *ResultPool) Length() int64 {
	pool.mu.Lock()
	defer pool.mu.Unlock()
	return int64(len(pool.results))
}
