package main

import (
	"BruteForce_SearchEnginer/common/utils"
	"log"
	"net/http"
)

func (app *application) searchHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"version": version,
	}

	if err := utils.WriteJSON(w, http.StatusOK, data); err != nil {
		err := utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		if err != nil {
			log.Println("Failed writing to JSON error")
		}
	}
}
