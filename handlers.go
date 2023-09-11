package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type ErrorResponse struct {
	Errors []string `json:"errors"`
}

func ProcessReceipt(w http.ResponseWriter, r *http.Request) {
	// Parse JSON request into a Receipt struct
	var receipt Receipt
	if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate receipt fields
	validationErrors := validateReceipt(receipt)

	if len(validationErrors.Errors) > 0 {
		errorResponse := ErrorResponse{Errors: validationErrors.Errors}
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
	points := calculatePoints(receipt)

	// Generate ID and save points to data store
	receiptID := generateUniqueID()
	setPoints(receiptID, points)

	fmt.Printf("Successfully saved Receipt with ID: %s and Points: %d\n", receiptID, points)

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
	response := GetPointsResponse{Points: points}

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
	points := int64(0)

	// Rule 1: One point for every alphanumeric character in the retailer name.
	points += countAlphaNumeric(receipt.Retailer)
	// Rule 2: 50 points if the total is a round dollar amount with no cents.
	if isRoundDollarAmount(receipt.Total) {
		points += 50
	}
	// Rule 3: 25 points if the total is a multiple of 0.25.
	if isMultipleOf25Cents(receipt.Total) {
		points += 25
	}
	// Rule 4: 5 points for every two items on the receipt.
	points += int64(numItemsOnReceipt(receipt.Items)/2) * 5
	// Rule 5: If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer. The result is the number of points earned.
	for _, item := range receipt.Items {
		if isMultipleOf3(item.ShortDescription) {
			points += int64(math.Ceil(stringToFloat(item.Price) * 0.2))
		}
	}
	// Rule 6: 6 points if the day in the purchase date is odd.
	if isOddDay(receipt.PurchaseDate) {
		points += 6
	}
	// Rule 7: 10 points if the time of purchase is after 2:00pm and before 4:00pm.
	if isBetween2And4PM(receipt.PurchaseTime) {
		points += 10
	}

	return points
}

func countAlphaNumeric(s string) int64 {
	count := 0
	for _, char := range s {
		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			count++
		}
	}
	return int64(count)
}

func isRoundDollarAmount(s string) bool {
	// Clean and convert string
	val := stringToFloat(s)
	if val == -1 {
		return false
	}

	// Compare float value to its int counterpart
	return val == math.Trunc(val)
}

func isMultipleOf25Cents(s string) bool {
	// Clean and convert string
	val := stringToFloat(s)
	if val == -1 {
		return false
	}

	// Check if cleanly divisible by 0.25
	return math.Mod(val, 0.25) == 0
}

func numItemsOnReceipt(items []Item) int {
	return len(items)
}

func isMultipleOf3(s string) bool {
	// Remove whitespace
	length := len(strings.ReplaceAll(s, " ", ""))
	if length > 0 {
		return length%3 == 0
	}

	return false
}

func isOddDay(s string) bool {
	day, err := getDayFromDate(s)
	if err != nil {
		return false
	}

	return day%2 != 0
}

func isBetween2And4PM(s string) bool {
	parsedTime, _ := parseTime(s)

	startTime, _ := parseTime("14:00")
	endTime, _ := parseTime("16:00")

	return parsedTime.After(startTime) && parsedTime.Before(endTime)
}

func setPoints(id string, points int64) string {
	pointStore[id] = points
	return id
}

func getPoints(id string) (int64, error) {
	if points, ok := pointStore[id]; ok {
		return points, nil
	}

	// Create and Format error response as JSON
	errResponses := ErrorResponse{
		Errors: []string{fmt.Sprintf("no receipt found for ID %s", id)},
	}

	errJSON, err := json.Marshal(errResponses)
	if err != nil {
		return 0, err
	}

	return 0, fmt.Errorf(string(errJSON))
}

func validateReceipt(receipt Receipt) ErrorResponse {
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

	return ErrorResponse{Errors: errorStrings}
}

func generateUniqueID() string {
	return uuid.NewString()
}

func stringToFloat(s string) float64 {
	// Remove commas if they exist
	s = strings.ReplaceAll(s, ",", "")

	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return -1
	}
	return val
}

func getDayFromDate(s string) (int, error) {
	parsedDate, err := time.Parse(time.DateOnly, s)
	if err != nil {
		return 0, err
	}
	return parsedDate.Day(), nil
}

func parseTime(s string) (time.Time, error) {
	return time.Parse("15:04", s)
}
