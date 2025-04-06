package main

import (
	"BruteForce_SearchEnginer/common/logger"
	"BruteForce_SearchEnginer/manager/internal/pool"
	"expvar"
	"go.uber.org/zap"
	"runtime"
)

const (
	version = "0.0.0"
)

type application struct {
	directoryPool *pool.DirectoryPool
	resultPool    *pool.ResultPool
	config        config
	logger        *zap.Logger
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
	app.logger.Fatal("server error", zap.Error(app.run(mux)))

	err := app.run(mux)
	if err != nil {
		app.logger.Fatal("server error", zap.Error(err))
	}

	// TODO() create the directory pool

	// TODO() create the results pool

	// TODO() create the endpoints for start, stop, directory-pool, results-pool, get-results

	// TODO() create the process of creating other workers
}

func setup() *application {
	var app application

	dirPool := pool.DirectoryPool{}
	resultPool := pool.ResultPool{}
	appLogger := logger.InitLogger("./../../manager.log")

	config := config{
		addr:   ":8080",
		apiURL: "localhost:8080",
	}

	app.directoryPool = &dirPool
	app.resultPool = &resultPool
	app.logger = appLogger
	app.config = config

	return &app
}
