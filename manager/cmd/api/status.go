package main

import (
	"BruteForce_SearchEnginer/common/model"
	"BruteForce_SearchEnginer/common/utils"
	"go.uber.org/zap"
	"net/http"
)

func (app *application) startHandler(w http.ResponseWriter, r *http.Request) {
	var startSignal model.StartSignal
	err := utils.ReadJSON(w, r, &startSignal)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.workerManager.AddWorker(startSignal.WorkerId)
	if err != nil {
		app.logger.Error("error creating worker", zap.Error(err))
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) stopHandler(w http.ResponseWriter, r *http.Request) {
	var stopSignal model.StopSignal
	err := utils.ReadJSON(w, r, &stopSignal)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.workerManager.RemoveWorker(stopSignal.WorkerId)
	if err != nil {
		app.logger.Error("error deleting worker", zap.Error(err))
		app.internalServerError(w, r, err)
		return
	}
}
