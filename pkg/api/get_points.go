package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"receipt-processor/pkg/models"

	"github.com/gorilla/mux"
)

func GetPoints(w http.ResponseWriter, r *http.Request) {
	// Retrieve ID from URL
	vars := mux.Vars(r)
	receiptID := vars["id"]

	// Retrieve points
	points, err := getPoints(receiptID)
	// Error when ID doesn't exist
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Create response struct
	response := models.GetPointsResponse{Points: points}

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

func getPoints(id string) (int64, error) {
	if points, ok := models.PointStore[id]; ok {
		return points, nil
	}

	// Create and Format error response as JSON
	errResponses := models.ErrorResponse{
		Errors: []string{fmt.Sprintf("no receipt found for ID %s", id)},
	}

	errJSON, err := json.Marshal(errResponses)
	if err != nil {
		return 0, err
	}

	return 0, fmt.Errorf(string(errJSON))
}
