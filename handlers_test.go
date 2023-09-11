package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestGetPointsHandler(t *testing.T) {
	receiptID := "12345"
	pointStore[receiptID] = 100

	request := httptest.NewRequest("GET", "/receipts/"+receiptID+"/points", nil)
	recorder := httptest.NewRecorder()

	// Inject Mock Vars
	request = mux.SetURLVars(request, map[string]string{
		"id": receiptID,
	})

	GetPoints(recorder, request)

	// Check for OK status code
	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d. %v", http.StatusOK, recorder.Code, recorder)
	}

	// Extract points from json response
	var response GetPointsResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Errorf("Error parsing response: %v", err)
	}

	// Compare expected point value to actual
	want := int64(100)
	if response.Points != want {
		t.Errorf("Want %d, got %d", want, response.Points)
	}
}
