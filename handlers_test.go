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

func TestCalculatePoints(t *testing.T) {
	testCases := []struct {
		description    string
		receipt        Receipt
		expectedPoints int64
	}{
		{
			description: "Sample Receipt 1",
			receipt: Receipt{
				Retailer:     "Target",
				PurchaseDate: "2022-01-02",
				PurchaseTime: "11:11",
				Items: []Item{
					{ShortDescription: "Pepsi - 12-oz", Price: "1.25"},
				},
				Total: "1.25",
			},
			expectedPoints: int64(31),
		},
		{
			description: "Sample Receipt 2",
			receipt: Receipt{
				Retailer:     "Walgreens",
				PurchaseDate: "2022-01-21",
				PurchaseTime: "08:13",
				Items: []Item{
					{ShortDescription: "Pepsi - 12-oz", Price: "1.25"},
					{ShortDescription: "Dasani", Price: "1.40"},
				},
				Total: "2.65",
			},
			expectedPoints: int64(21),
		},
		{
			description: "Sample Receipt 3",
			receipt: Receipt{
				Retailer:     "Record Store",
				PurchaseDate: "2022-01-22",
				PurchaseTime: "15:42",
				Items: []Item{
					{ShortDescription: "Blonde - Frank Ocean", Price: "24.99"},
					{ShortDescription: "Volcano - Jungle", Price: "17.99"},
					{ShortDescription: "Derealised - Jadu Heart", Price: "18.50"},
					{ShortDescription: "Rare Pleasure - Mndsgn", Price: "22.99"},
					{ShortDescription: "Maps - billy woods", Price: "19.99"},
				},
				Total: "104.46",
			},
			expectedPoints: int64(35),
		},
		{
			description: "Sample Receipt 4",
			receipt: Receipt{
				Retailer:     "T's Candy",
				PurchaseDate: "2022-01-09",
				PurchaseTime: "16:00",
				Items: []Item{
					{ShortDescription: "Mambas", Price: "2.00"},
					{ShortDescription: "Skittles", Price: "1.75"},
					{ShortDescription: "Albanese Gummy Bears", Price: "3.25"},
				},
				Total: "7.00",
			},
			expectedPoints: int64(95),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			points := calculatePoints(testCase.receipt)

			if points != testCase.expectedPoints {
				t.Errorf("Expected points: %d, got %d", testCase.expectedPoints, points)
			}
		})
	}
}

func TestCountAlphaNumeric(t *testing.T) {
	testCases := []struct {
		input    string
		expected int64
	}{
		{"Hello World!", 10}, // 10 alphanumeric characters
		{"12345", 5},         // 5 numeric characters
		{"ABCdef", 6},        // 6 letter characters
		{"!@#$%^&*()_+", 0},  // No alphanumeric characters
		{"", 0},              // Empty string
		{"   ", 0},           // Only spaces, no alphanumeric characters
	}

	for _, testCase := range testCases {
		t.Run(testCase.input, func(t *testing.T) {
			count := countAlphaNumeric(testCase.input)
			if count != testCase.expected {
				t.Errorf("Expected %d alphanumeric characters, but got %d", testCase.expected, count)
			}
		})
	}
}

func TestIsRoundDollarMount(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{"0", true},
		{"12.34", false},
		{"5.0", true},
		{"1,000,000.00", true},
		{"0.01", false},
		{"abc", false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.input, func(t *testing.T) {
			actual := isRoundDollarAmount(testCase.input)
			if actual != testCase.expected {
				t.Errorf("Expected isRoundDollarAmount(%s) to be %v, but got %v", testCase.input, testCase.expected, actual)
			}
		})
	}
}

func TestIsMultipleOf25Cents(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{"0.00", true},
		{"0.25", true},
		{"1.23", false},
		{"1,000,000.75", true},
		{"19.4", false},
		{"abc", false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.input, func(t *testing.T) {
			actual := isMultipleOf25Cents(testCase.input)
			if actual != testCase.expected {
				t.Errorf("Expected isRoundDollarAmount(%s) to be %v, but got %v", testCase.input, testCase.expected, actual)
			}
		})
	}
}

func TestNumItemsOnReceipt(t *testing.T) {
	testCases := []struct {
		description   string
		items         []Item
		expectedCount int
	}{
		{
			description:   "No Items",
			items:         []Item{},
			expectedCount: 0,
		},
		{
			description:   "1 Item",
			items:         []Item{{ShortDescription: "Item 1", Price: "1.23"}},
			expectedCount: 1,
		},
		{
			description: "5 Items",
			items: []Item{
				{ShortDescription: "Item 1", Price: "1.23"},
				{ShortDescription: "Item 2", Price: "1.23"},
				{ShortDescription: "Item 3", Price: "1.23"},
				{ShortDescription: "Item 4", Price: "1.23"},
				{ShortDescription: "Item 5", Price: "1.23"}},
			expectedCount: 5,
		},
		{
			description: "8 Items",
			items: []Item{
				{ShortDescription: "Item 1", Price: "1.23"},
				{ShortDescription: "Item 2", Price: "1.23"},
				{ShortDescription: "Item 3", Price: "1.23"},
				{ShortDescription: "Item 4", Price: "1.23"},
				{ShortDescription: "Item 5", Price: "1.23"},
				{ShortDescription: "Item 6", Price: "1.23"},
				{ShortDescription: "Item 7", Price: "1.23"},
				{ShortDescription: "Item 8", Price: "1.23"}},
			expectedCount: 8,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			count := numItemsOnReceipt(testCase.items)
			if count != testCase.expectedCount {
				t.Errorf("Expected %d items, but got %d", testCase.expectedCount, count)
			}
		})
	}
}

func TestIsMultipleOf3(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{"", false},
		{"abc", true},
		{"abc 123", true},
		{"sally sells sea shells", false},
		{"! !  ! !", false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.input, func(t *testing.T) {
			actual := isMultipleOf3(testCase.input)
			if actual != testCase.expected {
				t.Errorf("Expected isMultipleOf3(%s) to be %v, but got %v", testCase.input, testCase.expected, actual)
			}
		})
	}
}

func TestIsOddDay(t *testing.T) {
	testCases := []struct {
		purchaseDate string
		expected     bool
	}{
		{"2022-01-01", true},  // January 1, odd
		{"2023-02-14", false}, // February 14, even
		{"2023-10-31", true},  // October 31, odd
		{"2023-12-25", true},  // December 25, odd
		{"2023-11-22", false}, // November 22 (my birthday)
	}

	for _, testCase := range testCases {
		t.Run(testCase.purchaseDate, func(t *testing.T) {
			actual := isOddDay(testCase.purchaseDate)
			if actual != testCase.expected {
				t.Errorf("Expected isOddDay(%s) to be %v, but got %v", testCase.purchaseDate, testCase.expected, actual)
			}
		})
	}
}

func TestIsBetween2And4PM(t *testing.T) {
	testCases := []struct {
		purchaseTime string
		expected     bool
	}{
		{"11:11", false}, // Before
		{"14:00", false}, // 2 PM
		{"14:01", true},  // Between
		{"15:59", true},  // Between
		{"16:00", false}, // 4 PM
		{"16:01", false}, // After
		{"20:12", false}, // After
	}

	for _, testCase := range testCases {
		t.Run(testCase.purchaseTime, func(t *testing.T) {
			actual := isBetween2And4PM(testCase.purchaseTime)
			if actual != testCase.expected {
				t.Errorf("Expected isBetween2And4PM(%s) to be %v, but got %v", testCase.purchaseTime, testCase.expected, actual)
			}
		})
	}
}
