package main

import (
	"BruteForce_SearchEnginer/common/model"
	"BruteForce_SearchEnginer/common/utils"
	"errors"
	"go.uber.org/zap"
	"net/http"
)

var (
	ErrNoDirectoryInPool = errors.New("no directory is in the pool")
)

func (app *application) updateDirectoryPoolHandler(w http.ResponseWriter, r *http.Request) {
	var dir model.DirectoryResponse
	err := utils.ReadJSON(w, r, &dir)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	app.directoryPool.AddDirectory(dir)
}

func (app *application) getDirectoryPoolHandler(w http.ResponseWriter, r *http.Request) {
	dir, dirExists := app.directoryPool.PopDirectory()
	if !dirExists {
		app.badRequestResponse(w, r, ErrNoDirectoryInPool)
		return
	}

	err := utils.WriteJSON(w, http.StatusOK, dir)
	if err != nil {
		app.logger.Error("error sending directory to worker", zap.Error(err))
		return
	}
}

func (app *application) getAllDirectoryPoolHandler(w http.ResponseWriter, r *http.Request) {
	directories := app.directoryPool.GetAllDirectories()
	err := utils.WriteJSON(w, http.StatusOK, directories)
	if err != nil {
		app.logger.Error("error sending directory to worker", zap.Error(err))
		return
	}
}
