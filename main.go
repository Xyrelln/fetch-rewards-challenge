package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// setup receipt manager
var rm = newReceiptManager()

func main() {
	// get port flag
	portFlag := flag.Int("port", PORT, "Port to run the server on")
	flag.Parse()
	PORT = *portFlag

	// router and routes
	r := mux.NewRouter()
	r.HandleFunc("/receipts/process", processHandler).Methods("POST")
	r.HandleFunc("/receipts/{id}/points", getPointsHandler).Methods("GET")

	// CORS
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		handlers.ExposedHeaders([]string{"Content-Type"}),
	)

	// run server
	log.Printf("Starting server on :%d\n", PORT) // on terminal
	if err := http.ListenAndServe(fmt.Sprintf(":%d", PORT), corsHandler(r)); err != nil {
		log.Fatal("Server failed to start")
	}
}
