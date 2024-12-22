package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Investment struct {
	UserId int     `json:"user_id"`
	Stock  string  `json:"stock"`
	Units  int     `json:"units"`
	Price  float64 `json:"price"`
}

func GetInvestment(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr := r.URL.Query().Get("user_id")
		user_id, err := strconv.Atoi(userIDStr)
		if err != nil || user_id <= 0 {
			http.Error(w, "Invalid user_id", http.StatusBadRequest)
			return
		}

		rows, err := db.Query("SELECT id, user_id, stock, units, price FROM portfolio WHERE user_id = ?", user_id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error querying portfolio: %v", err), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var investments []Investment

		for rows.Next() {
			var investment Investment
			var id int
			if err := rows.Scan(&id, &investment.UserId, &investment.Stock, &investment.Units, &investment.Price); err != nil {
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
