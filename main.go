package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	//Establish a new router instance
	router := mux.NewRouter()

	//Define API endpoints
	router.HandleFunc("/receipts/process", ProcessReceipt).Methods("POST")

	//Start the HTTP server
	log.Fatal(http.ListenAndServe(":8080", router))

}
