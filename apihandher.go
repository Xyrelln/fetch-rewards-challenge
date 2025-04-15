package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func processHandler(w http.ResponseWriter, r *http.Request) {
	var rcpt receipt
	if err := json.NewDecoder(r.Body).Decode(&rcpt); err != nil {
		http.Error(w, "The receipt is invalid.", http.StatusBadRequest)
		return
	}

	var err error
	var id string
	id, err = rm.calculateAndSavePoints(rcpt)
	if err != nil {
		http.Error(w, "The receipt is invalid.", http.StatusBadRequest)
		return
	}

	// response content
	res := struct {
		ID string `json:"id"`
	}{
		ID: id,
	}

	// write res
	writeJSON(w, res)
}

func getPointsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "No IDs are given", http.StatusBadRequest)
		return
	}

	points, err := rm.getPoint(id)
	if err != nil {
		http.Error(w, "No receipt found for that ID.", http.StatusBadRequest)
		return
	}

	// response content
	res := struct {
		Points int `json:"points"`
	}{
		Points: points,
	}

	// write res
	writeJSON(w, res)
}
