package main

import (
	"BruteForce_SearchEnginer/common/utils"
	"go.uber.org/zap"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warn(
		"not found response",
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
		zap.Error(err))

	_ = utils.WriteJSONError(w, http.StatusInternalServerError, "the server encountered a problem")
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warn(
		"not found response",
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
		zap.Error(err))

	_ = utils.WriteJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warn(
		"not found response",
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
		zap.Error(err))

	_ = utils.WriteJSONError(w, http.StatusNotFound, "not found")
}
