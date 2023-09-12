package main

import (
	"log"
	"net/http"

	"receipt-processor/pkg/api"

	"github.com/gorilla/mux"
)

func main() {
	//Establish a new router instance
	router := mux.NewRouter()

	//Define API endpoints
	router.HandleFunc("/receipts/process", api.ProcessReceipt).Methods("POST")
	router.HandleFunc("/receipts/{id}/points", api.GetPoints).Methods("GET")

	//Start the HTTP server
	log.Fatal(http.ListenAndServe(":8080", router))

}
