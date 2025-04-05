package main

import (
	"BruteForce_SearchEnginer/common/logger"
	"BruteForce_SearchEnginer/common/model"
	"BruteForce_SearchEnginer/worker/internal/matcher"
	"BruteForce_SearchEnginer/worker/internal/repo"
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

	// TODO() lookup pool directories

	// stop
	// TODO() send stop signal to the manager

}

func (w *worker) setup() {
	w.id = rand.Int64()

	zapLogger := logger.InitLogger("../logs/worker_" + strconv.FormatInt(w.id, 10) + ".log")
	typeMap, err := model.ParseFileTypesConfig("../common/file_types_config.json")
	if err != nil {
		zapLogger.Panic("Type map panic", zap.Error(err))
		return
	}
	fileRepo := repo.New(typeMap)
	requestMatcher := matcher.New(typeMap)

	workerConfig := workerConfig{
		managerURL:            "http://127.0.0.1",
		workerStopEndpoint:    "/stop",
		directoryPoolEndpoint: "/directory-pool",
		resultPoolEndpoint:    "/results-pool",
	}

	w.logger = zapLogger
	w.typeMap = typeMap
	w.fileRepo = fileRepo
	w.requestMatcher = requestMatcher
	w.config = workerConfig
}
