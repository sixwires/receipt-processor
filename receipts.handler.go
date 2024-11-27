package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Root struct for the JSON data
type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Total        string `json:"total"`
	Items        []Item `json:"items"`
}

// Item represents each item in the receipt
type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

func (s *APIServer) ReceiptsHandler() *mux.Router {
	r := mux.NewRouter().PathPrefix("/receipts").Subrouter()

	r.HandleFunc("/process", makeHTTPHandleFunc(s.ProcessReceipts)).Methods("POST")
	r.HandleFunc("/{id}/points", makeHTTPHandleFunc(s.GetPoints)).Methods("GET")

	return r
}

// Submits a receipt for processing
func (s *APIServer) ProcessReceipts(w http.ResponseWriter, r *http.Request) error {
	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	// Parse body into struct
	var receipt Receipt
	err = json.Unmarshal(body, &receipt)
	if err != nil {
		return err
	}

	// Pseudo-store receipt in struct
	// In practice, this is where we would store in our DB.
	u := uuid.New()
	Receipts[u] = receipt
	return WriteJson(w, http.StatusOK, map[string]uuid.UUID{"id": u})
}

// Returns the points awarded for the receipt
func (s *APIServer) GetPoints(w http.ResponseWriter, r *http.Request) error {
	// Get id from req url
	idStr := mux.Vars(r)["id"]

	// Attempt to parse the string into a UUID
	id, err := uuid.Parse(idStr)
	if err != nil {
		// If there's an error parsing, return a 400 Bad Request with an error message
		return WriteJson(w, http.StatusBadRequest, "Invalid UUID format.")
	}

	// Search for id in stored receipts
	receipt, found := Receipts[id]
	if !found {
		return WriteJson(w, http.StatusBadRequest, "Specified ID could not be found.")
	}

	// Generate and return points
	points := calculatePoints(receipt)
	return WriteJson(w, http.StatusOK, map[string]int{"points": points})
}

func calculatePoints(r Receipt) int {
	points := 0

	// One point for every alphanumeric character in the retailer name.
	points += getAlphaNumericCount(r.Retailer)

	// 50 points if the total is a round dollar amount with no cents.
	// 25 points if the total is a multiple of 0.25.
	points += getPointsFromTotal(r.Total)

	// 5 points for every two items on the receipt.
	// If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer. The result is the number of points earned.
	points += getPointsFromItems(r.Items)

	// 6 points if the day in the purchase date is odd.
	points += getPointsFromPurchaseDate(r.PurchaseDate)

	// 10 points if the time of purchase is after 2:00pm and before 4:00pm.
	points += getPointsFromPurchaseTime(r.PurchaseTime)
	return points
}

// Returns number of alphanumeric characters in string
func getAlphaNumericCount(s string) int {
	points := 0
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			points++
		}
	}
	return points
}

func getPointsFromTotal(total string) int {
	points := 0
	if total[len(total)-2:] == "00" {
		points += 50
	}

	// Parse the string into a float64
	value, err := strconv.ParseFloat(total, 64)
	if err != nil {
		fmt.Println("Error parsing total:", err)
	}

	// Scale the value to work with integers
	scaledValue := int(value * 100)

	// Check if the scaled value is divisible by 25
	if scaledValue%25 == 0 {
		points += 25
	}
	return points
}

func getPointsFromItems(items []Item) int {
	// 5 points for every two items
	points := 5 * int(math.Floor(float64(len(items))/2.0))

	// If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer. The result is the number of points earned.
	for _, item := range items {
		trimmed := strings.TrimSpace(item.ShortDescription)
		if len(trimmed)%3 != 0 {
			continue
		}

		value, err := strconv.ParseFloat(item.Price, 64)
		if err != nil {
			fmt.Println("Error parsing total")
			continue
		}

		fmt.Printf("%v: %v: %d\n", item.ShortDescription, item.Price, int(math.Ceil(value*0.2)))
		points += int(math.Ceil(value * 0.2))
	}

	return points
}

func getPointsFromPurchaseDate(d string) int {
	points := 0
	date, err := strconv.Atoi(d[len(d)-2:])
	if err != nil {
		fmt.Println("Error converting day to int:", err)
		return points
	}

	// 6 points if the day in the purchase date is odd.
	if date%2 == 1 {
		points += 6
	}
	return points
}

func getPointsFromPurchaseTime(t string) int {
	// Parse the input time (24-hour format HH:mm)
	parsedTime, err := time.Parse("15:04", t)
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return 0
	}

	// Start and end times (2 PM and 4 PM)
	startTime, _ := time.Parse("15:04", "14:00")
	endTime, _ := time.Parse("15:04", "16:00")

	// Check if the parsed time is between startTime and endTime
	if parsedTime.After(startTime) && parsedTime.Before(endTime) {
		return 10
	}
	return 0
}
