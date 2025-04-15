package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func writeJSON(w http.ResponseWriter, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func generateReceiptId() string {
	return uuid.New().String()
}
