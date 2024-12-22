package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func DeleteInvestment(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		investmentId := r.URL.Query().Get("id")
		id, err := strconv.Atoi(investmentId)
		if err != nil || id <= 0 {
			http.Error(w, "Invalid investment Id", http.StatusBadRequest)
			return
		}

		var investment Investment
		err = db.QueryRow("SELECT user_id, stock, units, price FROM investments WHERE id = ?", id).Scan(&investment.UserId, &investment.Stock, &investment.Units, &investment.Price)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error querying investments: %v", err), http.StatusInternalServerError)
			return
		}

		_, err = db.Query("DELETE FROM investments WHERE id = ?", id)
		if err != nil {
			http.Error(w, "Error removing the investment", http.StatusInternalServerError)
			return
		}

		var portId int
		var tot_units int
		var avg_price float64

		err = db.QueryRow("SELECT id, units, price FROM portfolio WHERE user_id = ? AND stock = ?", investment.UserId, investment.Stock).Scan(&portId, &tot_units, &avg_price)
		if err != nil {
			http.Error(w, "Error fetching portfolio details", http.StatusInternalServerError)
			return
		}

		avg_price = (avg_price*float64(tot_units) - investment.Price*float64(investment.Units)) / (float64(tot_units) - float64(investment.Units))
		tot_units -= investment.Units
		tot_amt := avg_price * float64(tot_units)

		if tot_units == 0 {
			_, err = db.Query("DELETE FROM portfolio WHERE id = ?", portId)
			if err != nil {
				http.Error(w, "Error delete investment", http.StatusInternalServerError)
				return
			}
		} else {
			_, err = db.Exec(
				"UPDATE portfolio SET units = ?, price = ?, tot_amt = ? WHERE id = ?",
				tot_units,
				avg_price,
				tot_amt,
				portId,
			)
			if err != nil {
				http.Error(w, "Error updating portfolio", http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Investment deleted successfully"})
	}
}
