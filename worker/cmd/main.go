package main

import (
	"BruteForce_SearchEnginer/common/logger"
	"BruteForce_SearchEnginer/common/model"
	"BruteForce_SearchEnginer/worker/internal/crawler"
	"BruteForce_SearchEnginer/worker/internal/matcher"
	"BruteForce_SearchEnginer/worker/internal/repo"
	"context"
	"go.uber.org/zap"
	"math/rand/v2"
	"strconv"
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
	workerStopEndpoint    string
	directoryPoolEndpoint string
	resultPoolEndpoint    string
}

func main() {
	var appWorker worker
	appWorker.setup()

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

	// TODO() send stop signal to the manager

	// stop
	appWorker.logger.Info("worker finished", zap.Int64("worker_id", appWorker.id))
}

func (w *worker) setup() {
	w.id = rand.Int64()
	w.config = workerConfig{
		managerURL:            "http://127.0.0.1",
		workerStopEndpoint:    "/stop",
		directoryPoolEndpoint: "/directory-pool",
		resultPoolEndpoint:    "/results-pool",
	}

	zapLogger := logger.InitLogger("../logs/worker_" + strconv.FormatInt(w.id, 10) + ".log")
	typeMap, err := model.ParseFileTypesConfig("../common/file_types_config.json")
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
