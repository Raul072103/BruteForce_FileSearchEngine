package main

import (
	"BruteForce_SearchEnginer/common/logger"
	"BruteForce_SearchEnginer/common/model"
	"BruteForce_SearchEnginer/worker/internal/crawler"
	"BruteForce_SearchEnginer/worker/internal/matcher"
	"BruteForce_SearchEnginer/worker/internal/repo"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"math/rand/v2"
	"net/http"
	"strconv"
	"time"
)

type worker struct {
	id             int64
	logger         *zap.Logger
	fileRepo       repo.FileRepo
	requestMatcher matcher.RequestMatcher
	typeMap        model.FileTypesConfig
	crawler        crawler.Crawler
	config         workerConfig
}

type workerConfig struct {
	managerURL            string
	workerStartEndpoint   string
	workerStopEndpoint    string
	directoryPoolEndpoint string
	resultPoolEndpoint    string
}

func main() {
	var appWorker worker
	appWorker.setup()

	// send start signal to manager
	err := appWorker.sendStartSignal()
	if err != nil {
		appWorker.logger.Error("error sending start signal to manager",
			zap.Error(err), zap.Int64("worker_id", appWorker.id))
		return
	}

	// get file search request from directory pool
	dirResponse, err := appWorker.requestDirectoryPool()
	if err != nil {
		appWorker.logger.Panic("Failed to request directory pool",
			zap.Error(err),
			zap.Int64("worker_id", appWorker.id))
	}

	// start crawler
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	searchRequest := model.ConvertNetworkSearchRequest(dirResponse.SearchRequest)
	appWorker.crawler.Run(ctx, dirResponse.Path, searchRequest)

	// send stop signal to the manager
	err = appWorker.sendStopSignal()
	if err != nil {
		appWorker.logger.Error("error sending stop signal to manager",
			zap.Error(err), zap.Int64("worker_id", appWorker.id))
		return
	}

	// stop
	appWorker.logger.Info("worker finished", zap.Int64("worker_id", appWorker.id))
}

func (w *worker) setup() {
	w.id = rand.Int64()
	w.config = workerConfig{
		managerURL:            "http://127.0.0.1",
		workerStopEndpoint:    "/stop",
		workerStartEndpoint:   "/start",
		directoryPoolEndpoint: "/directory-pool",
		resultPoolEndpoint:    "/results-pool",
	}

	zapLogger := logger.InitLogger("./../logs/worker_" + strconv.FormatInt(w.id, 10) + ".log")
	typeMap, err := model.ParseFileTypesConfig("./../common/file_types_config.utils")
	if err != nil {
		zapLogger.Panic("Type map panic", zap.Error(err))
		return
	}
	fileRepo := repo.New(typeMap)
	requestMatcher := matcher.New(typeMap)

	crawlerConfig := crawler.Config{
		DirectoryPoolEndpoint: w.config.managerURL + w.config.directoryPoolEndpoint,
		ResultsPoolEndpoint:   w.config.managerURL + w.config.resultPoolEndpoint,
	}
	directoryCrawler := crawler.New(w.id, fileRepo, requestMatcher, zapLogger, crawlerConfig, typeMap)

	w.logger = zapLogger
	w.typeMap = typeMap
	w.fileRepo = fileRepo
	w.requestMatcher = requestMatcher
	w.crawler = directoryCrawler
}

// sendStopSignal sends a stop signal to the manager when the worker finishes its task
func (w *worker) sendStopSignal() error {
	payload := model.StopSignal{WorkerId: w.id}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 5 * time.Second}

	stopURL := w.config.managerURL + w.config.workerStopEndpoint
	resp, err := client.Post(stopURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send stop signal, status code: %d", resp.StatusCode)
	}

	return nil
}

// sendStartSignal notifies the manager that the worker has started
func (w *worker) sendStartSignal() error {
	payload := model.StartSignal{WorkerId: w.id}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 5 * time.Second}

	startURL := w.config.managerURL + w.config.workerStartEndpoint
	resp, err := client.Post(startURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send start signal, status code: %d", resp.StatusCode)
	}

	return nil
}
