package utils

import (
	"receipt-processor/pkg/models"
	"testing"
)

func TestCalculatePoints(t *testing.T) {
	testCases := []struct {
		description    string
		receipt        models.Receipt
		expectedPoints int64
	}{
		{
			description: "Sample Receipt 1",
			receipt: models.Receipt{
				Retailer:     "Target",
				PurchaseDate: "2022-01-02",
				PurchaseTime: "11:11",
				Items: []models.Item{
					{ShortDescription: "Pepsi - 12-oz", Price: "1.25"},
				},
				Total: "1.25",
			},
			expectedPoints: int64(31),
		},
		{
			description: "Sample Receipt 2",
			receipt: models.Receipt{
				Retailer:     "Walgreens",
				PurchaseDate: "2022-01-21",
				PurchaseTime: "08:13",
				Items: []models.Item{
					{ShortDescription: "Pepsi - 12-oz", Price: "1.25"},
					{ShortDescription: "Dasani", Price: "1.40"},
				},
				Total: "2.65",
			},
			expectedPoints: int64(21),
		},
		{
			description: "Sample Receipt 3",
			receipt: models.Receipt{
				Retailer:     "Record Store",
				PurchaseDate: "2022-01-22",
				PurchaseTime: "15:42",
				Items: []models.Item{
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
			receipt: models.Receipt{
				Retailer:     "T's Candy",
				PurchaseDate: "2022-01-09",
				PurchaseTime: "16:00",
				Items: []models.Item{
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
			points, _ := CalculatePoints(testCase.receipt)

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
		items         []models.Item
		expectedCount int
	}{
		{
			description:   "No Items",
			items:         []models.Item{},
			expectedCount: 0,
		},
		{
			description:   "1 Item",
			items:         []models.Item{{ShortDescription: "Item 1", Price: "1.23"}},
			expectedCount: 1,
		},
		{
			description: "5 Items",
			items: []models.Item{
				{ShortDescription: "Item 1", Price: "1.23"},
				{ShortDescription: "Item 2", Price: "1.23"},
				{ShortDescription: "Item 3", Price: "1.23"},
				{ShortDescription: "Item 4", Price: "1.23"},
				{ShortDescription: "Item 5", Price: "1.23"}},
			expectedCount: 5,
		},
		{
			description: "8 Items",
			items: []models.Item{
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
		{"12345", false}, // Invalid
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
