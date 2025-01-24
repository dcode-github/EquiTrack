package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func GetIndvInvestment(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr := r.URL.Query().Get("userId")
		instrument := r.URL.Query().Get("instrument")
		user_id, err := strconv.Atoi(userIDStr)
		if err != nil || user_id <= 0 {
			http.Error(w, "Invalid user_id", http.StatusBadRequest)
			return
		}

		rows, err := db.Query("SELECT id, qty, avg, purchase_date FROM investments WHERE user_id = ? AND instrument = ?", user_id, instrument)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error querying investments: %v", err), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var investments []Investment

		for rows.Next() {
			var investment Investment
			var id int
			if err := rows.Scan(&id, &investment.Qty, &investment.Avg, &investment.Date); err != nil {
				http.Error(w, fmt.Sprintf("Error scanning investment data: %v", err), http.StatusInternalServerError)
				return
			}
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
