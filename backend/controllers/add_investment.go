package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func AddInvestment(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var investment Investment
		if err := json.NewDecoder(r.Body).Decode(&investment); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		current_time := time.Now().Format("2006-01-02")
		_, err := db.Exec(
			"INSERT INTO investments (user_id, stock, units, price, purchase_date) VALUES (?, ?, ?, ?, ?)",
			investment.UserId,
			investment.Stock,
			investment.Units,
			investment.Price,
			current_time,
		)
		if err != nil {
			http.Error(w, "Error adding investment", http.StatusInternalServerError)
			return
		}

		var units int
		var price float64

		err = db.QueryRow("SELECT units, price FROM portfolio WHERE user_id = ? AND stock = ?", investment.UserId, investment.Stock).Scan(&units, &price)
		if err == nil {
			fmt.Println("Current units:", units, "Current price:", price)

			tot_units := units + investment.Units
			avg_price := (price*float64(units) + investment.Price*float64(investment.Units)) / float64(tot_units)

			fmt.Println("Updated total units:", tot_units, "Updated average price:", avg_price)

			_, err = db.Exec(
				"UPDATE portfolio SET units = ?, price = ?, tot_amt = ? WHERE user_id = ? AND stock = ?",
				tot_units,
				avg_price,
				avg_price*float64(tot_units),
				investment.UserId,
				investment.Stock,
			)
			if err != nil {
				fmt.Println("Error updating portfolio:", err)
				http.Error(w, "Error updating portfolio", http.StatusInternalServerError)
				return
			}
		} else {
			_, err = db.Exec(
				"INSERT INTO portfolio (user_id, stock, units, price, tot_amt) VALUES (?, ?, ?, ?, ?)",
				investment.UserId,
				investment.Stock,
				investment.Units,
				investment.Price,
				investment.Price*float64(investment.Units),
			)
			if err != nil {
				fmt.Println("Portfolio update error")
				http.Error(w, "Error adding investment", http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Investment added successfully"})
	}
}
