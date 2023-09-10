package main

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

var pointStore = make(map[string]int64)
