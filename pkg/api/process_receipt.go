package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"receipt-processor/pkg/models"
	"receipt-processor/pkg/utils"

	"github.com/google/uuid"
)

func ProcessReceipt(w http.ResponseWriter, r *http.Request) {
	// Parse JSON request into a Receipt struct
	var receipt models.Receipt
	if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate receipt fields
	validationErrors := validateReceipt(receipt)

	if len(validationErrors.Errors) > 0 {
		errorResponse := models.ErrorResponse{Errors: validationErrors.Errors}
		jsonResponse, err := json.Marshal(errorResponse)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonResponse)
		return
	}

	// Calculate points for Receipt
	points, breakdown := utils.CalculatePoints(receipt)
	fmt.Print(breakdown)

	// Generate ID and save points to data store
	receiptID := generateUniqueID()
	setPoints(receiptID, points)

	fmt.Printf("Successfully saved Receipt with ID: %s and Points: %d\n", receiptID, points)

	// Create response struct
	response := models.PostReceiptResponse{ID: receiptID}

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

func setPoints(id string, points int64) string {
	models.PointStore[id] = points
	return id
}

func validateReceipt(receipt models.Receipt) models.ErrorResponse {
	var validationErrors []error
	// Validate receipt fields
	if receipt.Retailer == "" {
		validationErrors = append(validationErrors, errors.New("field 'retailer' is required"))
	}
	if receipt.PurchaseDate == "" {
		validationErrors = append(validationErrors, errors.New("field 'purchaseDate' is required"))
	}
	if receipt.PurchaseTime == "" {
		validationErrors = append(validationErrors, errors.New("field 'purchaseTime' is required"))
	}
	if len(receipt.Items) == 0 {
		validationErrors = append(validationErrors, errors.New("at least one item is required"))
	} else {
		// Check each item for missing fields
		for _, item := range receipt.Items {
			if item.ShortDescription == "" {
				validationErrors = append(validationErrors, errors.New("field 'shortDescription' is required for items"))
			}
			if item.Price == "" {
				validationErrors = append(validationErrors, errors.New("field 'price' is required for items"))
			}
		}
	}
	if receipt.Total == "" {
		validationErrors = append(validationErrors, errors.New("field 'total' is required"))
	}

	errorStrings := make([]string, len(validationErrors))
	for i, err := range validationErrors {
		errorStrings[i] = err.Error()
	}

	return models.ErrorResponse{Errors: errorStrings}
}

func generateUniqueID() string {
	return uuid.NewString()
}
