package controllers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func DeleteInvestment(db *sql.DB, redisClient *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		investmentId := r.URL.Query().Get("id")
		id, err := strconv.Atoi(investmentId)
		if err != nil || id <= 0 {
			http.Error(w, "Invalid investment Id", http.StatusBadRequest)
			return
		}

		var investment Investment
		err = db.QueryRow("SELECT user_id, instrument, qty, avg FROM investments WHERE id = ?", id).Scan(&investment.UserId, &investment.Instrument, &investment.Qty, &investment.Avg)
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

		err = db.QueryRow("SELECT id, qty, avg FROM portfolio WHERE user_id = ? AND instrument = ?", investment.UserId, investment.Instrument).Scan(&portId, &tot_units, &avg_price)
		if err != nil {
			http.Error(w, "Error fetching portfolio details", http.StatusInternalServerError)
			return
		}

		avg_price = (avg_price*float64(tot_units) - investment.Avg*float64(investment.Qty)) / (float64(tot_units) - float64(investment.Qty))
		tot_units -= investment.Qty
		tot_amt := avg_price * float64(tot_units)

		if tot_units == 0 {
			_, err = db.Query("DELETE FROM portfolio WHERE id = ?", portId)
			if err != nil {
				http.Error(w, "Error delete investment", http.StatusInternalServerError)
				return
			}
		} else {
			_, err = db.Exec(
				"UPDATE portfolio SET qty = ?, avg = ?, tot_amt = ? WHERE id = ?",
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

		cacheKey := "user:" + strconv.Itoa(investment.UserId)

		ctx := context.Background()
		err = redisClient.Del(ctx, cacheKey).Err()
		if err != nil {
			log.Println("Error invalidating cache: ", err)
		} else {
			log.Println("Cache invalidated for user: ", investment.UserId)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Investment deleted successfully"})
	}
}
