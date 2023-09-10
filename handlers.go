package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func ProcessReceipt(w http.ResponseWriter, r *http.Request) {
	// Parse JSON request into a Receipt struct
	var receipt Receipt
	if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Calculate points for Receipt
	points := calculatePoints(receipt)

	// Generate ID and save points to data store
	receiptID := generateUniqueID()
	setPoints(receiptID, points)

	fmt.Printf("Successfully saved Receipt with ID: %s and Points: %d", receiptID, points)

	// Create response struct
	response := PostReceiptResponse{ID: receiptID}

	// Serialize response into JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the header and send the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func calculatePoints(receipt Receipt) int64 {
	return int64(100)
}

func setPoints(id string, points int64) string {
	pointStore[id] = points
	return id
}

func generateUniqueID() string {
	return uuid.NewString()
}