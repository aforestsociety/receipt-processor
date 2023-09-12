package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPostReceiptHandler(t *testing.T) {
	// Define slice of test cases
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
