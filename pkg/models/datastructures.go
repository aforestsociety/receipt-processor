package models

type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []Item `json:"items"`
	Total        string `json:"total"`
}

type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

type PostReceiptResponse struct {
	ID string `json:"id"`
}

type GetPointsResponse struct {
	Points int64 `json:"points"`
}

type ErrorResponse struct {
	Errors []string `json:"errors"`
}

var PointStore = make(map[string]int64)
