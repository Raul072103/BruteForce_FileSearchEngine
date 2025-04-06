package main

import (
	"BruteForce_SearchEnginer/common/model"
	"BruteForce_SearchEnginer/common/utils"
	"go.uber.org/zap"
	"net/http"
)

func (app *application) updateResultPoolHandler(w http.ResponseWriter, r *http.Request) {
	var fileSearchResponse model.FileSearchResponse
	err := utils.ReadJSON(w, r, &fileSearchResponse)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	app.resultPool.AddResult(fileSearchResponse)
}

func (app *application) getResultPoolHandler(w http.ResponseWriter, r *http.Request) {
	fileSearchResponses := app.resultPool.GetAllResults()

	err := utils.WriteJSON(w, http.StatusOK, fileSearchResponses)
	if err != nil {
		app.logger.Error("error sending results through get all results endpoint", zap.Error(err))
	}
}
