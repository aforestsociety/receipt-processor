package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func TestGetPointsHandler(t *testing.T) {
	// Define slice of test cases
	testCases := []struct {
		description    string
		requestPath    string
		receiptID      string
		expectedStatus int
		expectedPoints int64
	}{
		{
			description:    "Valid ID",
			requestPath:    "/receipts/12345/points",
			receiptID:      "12345",
			expectedStatus: http.StatusOK,
			expectedPoints: 100,
		},
		{
			description:    "ID does not exist",
			requestPath:    "/receipts/99999/points",
			receiptID:      "99999",
			expectedStatus: http.StatusNotFound,
			expectedPoints: 0,
		},
		{
			description:    "No ID provided",
			requestPath:    "/receipts/points",
			receiptID:      "",
			expectedStatus: http.StatusNotFound,
			expectedPoints: 0,
		},
		{
			description:    "Invalud URL",
			requestPath:    "/invalid",
			receiptID:      "",
			expectedStatus: http.StatusNotFound,
			expectedPoints: 0,
		},
	}

	// Check cases one by one
	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			if testCase.receiptID == "12345" {
				pointStore[testCase.receiptID] = testCase.expectedPoints
			}

			request := httptest.NewRequest("GET", testCase.requestPath, nil)
			recorder := httptest.NewRecorder()

			// Inject Mock Vars
			request = mux.SetURLVars(request, map[string]string{
				"id": testCase.receiptID,
			})

			GetPoints(recorder, request)

			// Check for expected status code
			if recorder.Code != testCase.expectedStatus {
				t.Errorf("Expected status code %d, got %d", testCase.expectedStatus, recorder.Code)
			}

			// Extract points from json response
			var response GetPointsResponse
			if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
				t.Errorf("Error parsing response: %v", err)
			}

			// Compare expected point value to actual
			if response.Points != testCase.expectedPoints {
				t.Errorf("Want %d, got %d", testCase.expectedPoints, response.Points)
			}
		})
	}
}

func TestPostReceiptHandler(t *testing.T) {
	testCases := []struct {
		description    string
		requestBody    string
		expectedStatus int
	}{
		{
			description:    "Valid Receipt",
			requestBody:    `{"retailer": "Test Retailer", "purchaseDate": "2022-01-01", "purchaseTime": "13:01", "items": [{"shortDescription": "Test Item", "price": "9.99"}], "total": "9.99"}`,
			expectedStatus: http.StatusOK,
		},
		{
			description:    "Invalid JSON",
			requestBody:    `invalid`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			description:    "Missing Retailer Field",
			requestBody:    `{"purchaseDate": "2022-01-01", "purchaseTime": "13:01", "items": [{"shortDescription": "Test Item", "price": "9.99"}], "total": "9.99"}`,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			request := httptest.NewRequest("POST", "/receipts/process", strings.NewReader(testCase.requestBody))
			recorder := httptest.NewRecorder()

			ProcessReceipt(recorder, request)

			// Check for expected status code
			if recorder.Code != testCase.expectedStatus {
				t.Errorf("Expected status code %d, got %d", testCase.expectedStatus, recorder.Code)
			}

			// Check for empty response body
			responeBody := recorder.Body.String()
			if len(responeBody) == 0 {
				t.Error("Expected non empty response body, got empty")
			}
		})
	}
}
