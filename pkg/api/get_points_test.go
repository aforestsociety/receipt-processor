package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"receipt-processor/pkg/models"

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

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			if testCase.receiptID == "12345" {
				models.PointStore[testCase.receiptID] = testCase.expectedPoints
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
			var response models.GetPointsResponse
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
