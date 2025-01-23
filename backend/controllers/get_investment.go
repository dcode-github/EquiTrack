package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Investment struct {
	UserId           int       `json:"user_id"`
	Instrument       string    `json:"instrument"`
	Qty              int       `json:"qty"`
	Avg              float64   `json:"avg"`
	Price            float64   `json:"ltp"`
	CurVal           float64   `json:"currVal"`
	PNL              float64   `json:"pnl"`
	NetChg           float64   `json:"netChng"`
	PercentageChange float64   `json:"dayChng"`
	Date             time.Time `json:"date"`
}

func GetInvestment(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr := r.URL.Query().Get("userId")
		user_id, err := strconv.Atoi(userIDStr)
		if err != nil || user_id <= 0 {
			http.Error(w, "Invalid user_id", http.StatusBadRequest)
			return
		}

		rows, err := db.Query("SELECT id, user_id, instrument, qty, avg FROM portfolio WHERE user_id = ?", user_id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error querying portfolio: %v", err), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var investments []Investment

		for rows.Next() {
			var investment Investment
			var id int
			if err := rows.Scan(&id, &investment.UserId, &investment.Instrument, &investment.Qty, &investment.Avg); err != nil {
				http.Error(w, fmt.Sprintf("Error scanning investment data: %v", err), http.StatusInternalServerError)
				return
			}

			stock, err := fetchLivePrice(investment.Instrument)
			if err != nil {
				log.Println("Error fetching live price:", err)
				investment.Price = 0
				investment.PercentageChange = 0
			} else {
				investment.Price = stock.Price
				investment.PercentageChange = stock.PercentageChange
			}

			investment.CurVal = roundToTwoDecimalPlaces(investment.Price * float64(investment.Qty))
			investment.PNL = roundToTwoDecimalPlaces(float64(investment.Qty) * (investment.Price - investment.Avg))
			investment.NetChg = roundToTwoDecimalPlaces(investment.PNL / (investment.Avg * float64(investment.Qty)) * 100.0)

			investments = append(investments, investment)
		}

		if err := rows.Err(); err != nil {
			http.Error(w, fmt.Sprintf("Error iterating over rows: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if len(investments) == 0 {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"message": "No investments found for the user"})
		} else {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(investments)
		}
	}
}

func fetchLivePrice(instrument string) (*StockData, error) {
	url := fmt.Sprintf("http://localhost:8080/price?instrument=%s", instrument)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: %s", resp.Status)
	}

	var stock StockData
	err = json.NewDecoder(resp.Body).Decode(&stock)
	if err != nil {
		return nil, fmt.Errorf("error decoding price data: %v", err)
	}

	return &stock, nil
}
