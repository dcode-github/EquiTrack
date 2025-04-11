package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func extractStockArray(scriptContent string) ([]float64, error) {
	re := regexp.MustCompile(`"INR",\[([^\]]+)\]`)
	matches := re.FindStringSubmatch(scriptContent)

	if len(matches) < 2 {
		return nil, fmt.Errorf("no INR array found")
	}

	arrayStr := matches[1]
	stringValues := strings.Split(arrayStr, ",")
	stockArray := make([]float64, len(stringValues))

	for i, valStr := range stringValues {
		val, err := strconv.ParseFloat(strings.TrimSpace(valStr), 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing value: %v", err)
		}
		stockArray[i] = val
	}

	return stockArray, nil
}

func LivePrice(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		instrument := r.URL.Query().Get("instrument")
		if instrument == "" {
			log.Println("No instrument in URL")
			http.Error(w, "Instrument is required", http.StatusBadRequest)
			return
		}
		url := fmt.Sprintf("https://www.google.com/finance/quote/%s:NSE", instrument)
		resp, err := http.Get(url)
		if err != nil {
			log.Println("Error accessing external API")
			http.Error(w, "Error fetching data: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error reading body")
			http.Error(w, "Error reading response: "+err.Error(), http.StatusInternalServerError)
			return
		}

		re := regexp.MustCompile(`<script class="ds:2"[^>]*>(.*?)</script>`)
		matches := re.FindStringSubmatch(string(body))

		if len(matches) < 2 {
			http.Error(w, "No script tag with class 'ds:2' found", http.StatusInternalServerError)
			return
		}

		scriptContent := matches[1]
		stockArray, err := extractStockArray(scriptContent)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		price := roundToTwoDecimalPlaces(stockArray[0])
		percentageChange := roundToTwoDecimalPlaces(stockArray[2])

		data := StockData{
			Price:            price,
			PercentageChange: percentageChange,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}
}

func roundToTwoDecimalPlaces(value float64) float64 {
	return math.Round(value*100) / 100
}
