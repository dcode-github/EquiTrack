package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

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
		// var totalInvestment, totalCurrentVal, totalPNL float64

		for rows.Next() {
			var investment Investment
			var id int
			if err := rows.Scan(&id, &investment.UserId, &investment.Instrument, &investment.Qty, &investment.Avg); err != nil {
				http.Error(w, fmt.Sprintf("Error scanning investment data: %v", err), http.StatusInternalServerError)
				return
			}
			investment.Avg = roundToTwoDecimalPlaces(investment.Avg)

			// stock, err := fetchLivePrice(investment.Instrument)
			// if err != nil {
			// 	log.Println("Error fetching live price:", err)
			// 	investment.Price = 0
			// 	investment.PercentageChange = 0
			// } else {
			// 	investment.Price = stock.Price
			// 	investment.PercentageChange = stock.PercentageChange
			// }

			investment.TotInvestment = float64(investment.Qty) * investment.Avg
			// investment.CurVal = roundToTwoDecimalPlaces(investment.Price * float64(investment.Qty))
			// investment.PNL = roundToTwoDecimalPlaces(float64(investment.Qty) * (investment.Price - investment.Avg))
			// investment.NetChg = roundToTwoDecimalPlaces(investment.PNL / (investment.Avg * float64(investment.Qty)) * 100.0)

			// totalInvestment += investment.TotInvestment
			// totalCurrentVal += investment.CurVal
			// totalPNL += investment.PNL

			investments = append(investments, investment)
		}

		if err := rows.Err(); err != nil {
			http.Error(w, fmt.Sprintf("Error iterating over rows: %v", err), http.StatusInternalServerError)
			return
		}

		// var totalPNLPercent float64
		// if totalInvestment > 0 {
		// 	totalPNLPercent = roundToTwoDecimalPlaces((totalPNL / totalInvestment) * 100)
		// }

		// totalData := TotalInvestmentData{
		// 	TotalInvestment: 0,
		// 	TotalCurrentVal: 0,
		// 	TotalPNL:        0,
		// 	TotalPNLPercent: 0,
		// }

		response := struct {
			Investments []Investment `json:"investments"`
			// TotalInvestmentData TotalInvestmentData `json:"total_investment_data"`
		}{
			Investments: investments,
			// TotalInvestmentData: totalData,
		}

		w.Header().Set("Content-Type", "application/json")
		if len(investments) == 0 {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"message": "No investments found for the user"})
		} else {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
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
