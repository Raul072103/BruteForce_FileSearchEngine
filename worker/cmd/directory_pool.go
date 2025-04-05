package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type DirectoryResponse struct {
	Path string `json:"path"`
}

// requestDirectoryPool queries the directory pool endpoint and returns the path if successful
func (w *worker) requestDirectoryPool() (string, error) {
	url := w.config.managerURL + w.config.directoryPoolEndpoint

	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var dirResp DirectoryResponse
	if err := json.NewDecoder(resp.Body).Decode(&dirResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	return dirResp.Path, nil
}
