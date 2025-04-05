package main

import (
	"BruteForce_SearchEnginer/common/model"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// requestDirectoryPool queries the directory pool endpoint and returns the path if successful
func (w *worker) requestDirectoryPool() (model.DirectoryResponse, error) {
	var dirResp model.DirectoryResponse
	url := w.config.managerURL + w.config.directoryPoolEndpoint

	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return dirResp, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return dirResp, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&dirResp); err != nil {
		return dirResp, fmt.Errorf("failed to parse response: %w", err)
	}

	return dirResp, nil
}
