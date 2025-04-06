package main

import (
	"BruteForce_SearchEnginer/common/json"
	"log"
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"version": version,
	}

	if err := json.WriteJSON(w, http.StatusOK, data); err != nil {
		err := json.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		if err != nil {
			log.Println("Failed writing to JSON error")
		}
	}
}
