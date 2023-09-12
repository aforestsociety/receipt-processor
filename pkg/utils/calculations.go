package utils

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
	"unicode"

	"receipt-processor/pkg/models"
)

func CalculatePoints(receipt models.Receipt) (int64, string) {
	points := int64(0)
	breakdown := ""

	// Rule 1: One point for every alphanumeric character in the retailer name.
	namePoints := countAlphaNumeric(receipt.Retailer)
	points += namePoints
	breakdown += fmt.Sprintf("%d points - retailer name (%s) has %d alphanumeric characters\n", namePoints, receipt.Retailer, namePoints)
	// Rule 2: 50 points if the total is a round dollar amount with no cents.
	if isRoundDollarAmount(receipt.Total) {
		points += 50
		breakdown += "50 points - total is a round dollar amount\n"
	}
	// Rule 3: 25 points if the total is a multiple of 0.25.
	if isMultipleOf25Cents(receipt.Total) {
		points += 25
		breakdown += "25 points - total is a multiple of 0.25\n"
	}
	// Rule 4: 5 points for every two items on the receipt.
	numItems := numItemsOnReceipt(receipt.Items)
	itemPoints := int64(numItems/2) * 5
	points += itemPoints
	breakdown += fmt.Sprintf("%d points - %d items (%d pairs @ 5 points each)\n", itemPoints, numItems, numItems/2)
	// Rule 5: If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer. The result is the number of points earned.
	for _, item := range receipt.Items {
		if isMultipleOf3(item.ShortDescription) {
			itemPrice := stringToFloat(item.Price)
			itemPoints := int64(math.Ceil(itemPrice * 0.2))
			points += itemPoints
			breakdown += fmt.Sprintf("%d points - item description (%s) is a multiple of 3, price: %.2f\n", itemPoints, item.ShortDescription, itemPrice)
		}
	}
	// Rule 6: 6 points if the day in the purchase date is odd.
	if isOddDay(receipt.PurchaseDate) {
		points += 6
		breakdown += "6 points - purchase date day is odd\n"
	}
	// Rule 7: 10 points if the time of purchase is after 2:00pm and before 4:00pm.
	if isBetween2And4PM(receipt.PurchaseTime) {
		points += 10
		breakdown += "10 points - purchase time is between 2:00pm and 4:00pm\n"
	}

	return points, fmt.Sprintf("Total Points: %d\nBreakdown:\n%s+ ---------\n= %d points\n", points, breakdown, points)
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

func numItemsOnReceipt(items []models.Item) int {
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
