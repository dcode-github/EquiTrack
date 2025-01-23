package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type StockData struct {
	Price            float64 `json:"price"`
	PercentageChange float64 `json:"per_change"`
}

func LivePrice(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		instrument := r.URL.Query().Get("instrument")
		if instrument == "" {
			http.Error(w, "Instrument is required", http.StatusBadRequest)
			return
		}
		url := fmt.Sprintf("https://www.screener.in/company/%s/consolidated/", instrument)
		price, percentageChange, err := extractStockData(url)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		price = roundToTwoDecimalPlaces(price)
		percentageChange = roundToTwoDecimalPlaces(percentageChange)
		data := StockData{
			Price:            price,
			PercentageChange: percentageChange,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}
}

func parsePrice(price string) (float64, error) {
	cleanedPrice := strings.Replace(price, "₹", "", -1)
	cleanedPrice = strings.Replace(cleanedPrice, ",", "", -1)
	parsedPrice, err := strconv.ParseFloat(strings.TrimSpace(cleanedPrice), 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing price: %v", err)
	}

	return parsedPrice, nil
}

func parsePercentage(percentage string) (float64, error) {
	cleanedPercentage := strings.TrimSpace(percentage)
	cleanedPercentage = strings.TrimSuffix(cleanedPercentage, "%")
	parsedPercentage, err := strconv.ParseFloat(cleanedPercentage, 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing percentage change: %v", err)
	}
	return parsedPercentage, nil
}

func roundToTwoDecimalPlaces(value float64) float64 {
	return math.Round(value*100) / 100
}

func extractStockData(url string) (float64, float64, error) {
	res, err := http.Get(url)
	if err != nil {
		return 0, 0, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return 0, 0, fmt.Errorf("error: Status Code %d", res.StatusCode)
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return 0, 0, err
	}
	priceText := doc.Find(".font-size-18.strong .flex span").First().Text()
	percentageText := doc.Find(".font-size-12.down.margin-left-4").Text()
	parsedPrice, err := parsePrice(priceText)
	if err != nil {
		return 0, 0, err
	}
	parsedPercentageChange, err := parsePercentage(percentageText)
	if err != nil {
		return 0, 0, err
	}

	return parsedPrice, parsedPercentageChange, nil
}
