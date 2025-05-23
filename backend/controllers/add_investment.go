package controllers

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func AddInvestment(db *sql.DB, redisClient *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var investment Investment
		if err := json.NewDecoder(r.Body).Decode(&investment); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		cacheKey := "user:" + strconv.Itoa(investment.UserId)

		ctx := context.Background()
		err := redisClient.Del(ctx, cacheKey).Err()
		if err != nil {
			log.Println("Error invalidating cache: ", err)
		} else {
			log.Println("Cache invalidated for user: ", investment.UserId)
		}

		current_time := time.Now().Format("2006-01-02")
		_, err = db.Exec(
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
