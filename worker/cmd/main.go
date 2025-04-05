package main

import (
	"BruteForce_SearchEnginer/common/model"
	"BruteForce_SearchEnginer/worker/internal/matcher"
	"BruteForce_SearchEnginer/worker/internal/repo"
	"go.uber.org/zap"
)

type worker struct {
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
	// TODO() setup logger, repo, matcher, fileTypesConfig

	// TODO() lookup pool directories

	// stop
	// TODO() send stop signal to the manager

}
