package main

import (
	"BruteForce_SearchEnginer/common/logger"
	"BruteForce_SearchEnginer/manager/internal/pool"
	"BruteForce_SearchEnginer/manager/internal/process_manager"
	"BruteForce_SearchEnginer/manager/internal/worker_manager"
	"expvar"
	"go.uber.org/zap"
	"runtime"
)

const (
	version = "0.0.0"
)

type application struct {
	directoryPool  *pool.DirectoryPool
	resultPool     *pool.ResultPool
	logger         *zap.Logger
	workerManager  *worker_manager.WorkerManager
	processManager *process_manager.ProcessManager
	config         config
}

type config struct {
	addr   string
	apiURL string
}

func main() {

	// Performance metrics
	expvar.NewString("version").Set(version)
	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))

	// manager entry point setup
	app := setup()

	mux := app.mount()

	err := app.run(mux)
	if err != nil {
		app.logger.Fatal("server error", zap.Error(err))
	}

	// TODO() create the process of creating other workers

}

func setup() *application {
	var app application

	dirPool := pool.NewDirectoryPool()
	resultPool := pool.NewResultPool()

	appLogger := logger.InitLogger("./../manager.log")
	workerManagerLogger := logger.InitLogger("./../worker_manger.log")
	processManagerLogger := logger.InitLogger("./../process_manager.log")

	workerManager := worker_manager.New(workerManagerLogger)
	processManager := process_manager.New(workerManager, dirPool, resultPool, processManagerLogger)

	config := config{
		addr:   ":8080",
		apiURL: "localhost:8080",
	}

	app.directoryPool = dirPool
	app.resultPool = resultPool
	app.logger = appLogger
	app.config = config
	app.workerManager = workerManager
	app.processManager = processManager

	return &app
}
