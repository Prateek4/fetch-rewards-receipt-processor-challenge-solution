package main

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"strings"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/google/uuid"
)

type Receipt struct {
	Retailer     string   `json:"retailer"`
	PurchaseDate string   `json:"purchaseDate"`
	PurchaseTime string   `json:"purchaseTime"`
	Items        []Item   `json:"items"`
	Total        string   `json:"total"`
}

type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

type ReceiptResponse struct {
	ID string `json:"id"`
}

type PointsResponse struct {
	Points int `json:"points"`
}

var receiptData = make(map[string]Receipt) // Store receipt data in-memory
var pointsData = make(map[string]int)      // Store points associated with receipt ID


func calculatePoints(receipt Receipt) int {
	points := 0

	// 1. Count alphanumeric characters in retailer name
	retailer_name := receipt.Retailer
	for _, c := range retailer_name {
		if isAlphanumeric(c) {
			points++
		}
	}

	// 2. Points for total being a round dollar amount
	total, err := strconv.ParseFloat(receipt.Total, 64)
	if err == nil {
		// Check if total is a round dollar
		if total == float64(int(total)) {
			points += 50
		}

		// 3. Check if the total is a multiple of 0.25
		if total*4 == float64(int(total*4)) {
			points += 25
		}
	}

	// 4. Points for every two items
	points += (len(receipt.Items) / 2) * 5

	// 5. For items with description length multiple of 3, apply the price logic
	for _, item := range receipt.Items {
		trimmed_desc := strings.TrimSpace(item.ShortDescription)
		if len(trimmed_desc)%3 == 0 { // Check if length is multiple of 3
			price, err := strconv.ParseFloat(item.Price, 64)
			if err == nil {
				points += int(math.Ceil(price * 0.2)) // Apply the multiplier and round up
			}
		}
	}

	// 6. Check if the purchase day is odd
	dateParts := strings.Split(receipt.PurchaseDate, "-")
	if len(dateParts) == 3 {
		day, err := strconv.Atoi(dateParts[2])
		if err == nil && day%2 != 0 {
			points += 6
		}
	}

	// 7. Points if purchase time is between 2:01 PM and 3:59 PM
	timeParts := strings.Split(receipt.PurchaseTime, ":")
	if len(timeParts) >= 2 {
		hour, err := strconv.Atoi(timeParts[0])
		minute, err2 := strconv.Atoi(timeParts[1])
		if err == nil && err2 == nil {
			if (hour == 14 && minute > 0) || (hour == 15) {
				points += 10
			}
		}
	}

	return points
}

func isAlphanumeric(c rune) bool {
	// Simple check for alphanumeric characters (letters and digits)
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
}

func processReceipt(w http.ResponseWriter, r *http.Request) {
	var receipt Receipt
	err := json.NewDecoder(r.Body).Decode(&receipt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Generate a unique ID for the receipt
	receiptID := uuid.New().String()
	receiptData[receiptID] = receipt
	pointsData[receiptID] = calculatePoints(receipt)

	response := ReceiptResponse{ID: receiptID}
	json.NewEncoder(w).Encode(response)
}

func getPoints(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	points, found := pointsData[id]
	if !found {
		http.Error(w, "Receipt not found", http.StatusNotFound)
		return
	}

	response := PointsResponse{Points: points}
	json.NewEncoder(w).Encode(response)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/receipts/process", processReceipt).Methods("POST")
	router.HandleFunc("/receipts/{id}/points", getPoints).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}
