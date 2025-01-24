package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

type Investment struct {
	UserId           int       `json:"user_id"`
	Instrument       string    `json:"instrument"`
	Qty              int       `json:"qty"`
	Avg              float64   `json:"avg"`
	Price            float64   `json:"ltp"`
	TotInvestment    float64   `json:"tot_invest"`
	CurVal           float64   `json:"currVal"`
	PNL              float64   `json:"pnl"`
	NetChg           float64   `json:"netChng"`
	PercentageChange float64   `json:"dayChng"`
	Date             time.Time `json:"date"`
}

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
			"INSERT INTO investments (user_id, instrument, qty, avg, purchase_date) VALUES (?, ?, ?, ?, ?)",
			investment.UserId,
			investment.Instrument,
			investment.Qty,
			investment.Avg,
			current_time,
		)
		if err != nil {
			http.Error(w, "Error adding investment", http.StatusInternalServerError)
			return
		}

		var units int
		var price float64

		err = db.QueryRow("SELECT qty, avg FROM portfolio WHERE user_id = ? AND instrument = ?", investment.UserId, investment.Instrument).Scan(&units, &price)
		if err == nil {

			tot_units := units + investment.Qty
			avg_price := (price*float64(units) + investment.Avg*float64(investment.Qty)) / float64(tot_units)

			_, err = db.Exec(
				"UPDATE portfolio SET qty = ?, avg = ?, tot_amt = ? WHERE user_id = ? AND instrument = ?",
				tot_units,
				avg_price,
				roundToTwoDecimalPlaces(avg_price*float64(tot_units)),
				investment.UserId,
				investment.Instrument,
			)
			if err != nil {
				http.Error(w, "Error updating portfolio", http.StatusInternalServerError)
				return
			}
		} else {
			_, err = db.Exec(
				"INSERT INTO portfolio (user_id, instrument, qty, avg, tot_amt) VALUES (?, ?, ?, ?, ?)",
				investment.UserId,
				investment.Instrument,
				investment.Qty,
				investment.Avg,
				investment.Avg*float64(investment.Qty),
			)
			if err != nil {
				http.Error(w, "Error adding investment", http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Investment added successfully"})
	}
}
