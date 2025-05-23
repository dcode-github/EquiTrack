package controllers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func GetInvestment(db *sql.DB, redisClient *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr := r.URL.Query().Get("userId")
		user_id, err := strconv.Atoi(userIDStr)
		if err != nil || user_id <= 0 {
			http.Error(w, "Invalid user_id", http.StatusBadRequest)
			return
		}

		var investments []Investment
		cacheKey := "user:" + userIDStr

		ctx := context.Background()
		cachedInvestments, err := redisClient.Get(ctx, cacheKey).Result()
		if err == nil && cachedInvestments != "" {
			log.Println("cache hit")
			fmt.Println(cachedInvestments)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(cachedInvestments))
			return
		}
		log.Println("cache miss")

		rows, err := db.Query("SELECT id, user_id, instrument, qty, avg FROM portfolio WHERE user_id = ?", user_id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error querying portfolio: %v", err), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var investment Investment
			var id int
			if err := rows.Scan(&id, &investment.UserId, &investment.Instrument, &investment.Qty, &investment.Avg); err != nil {
				http.Error(w, fmt.Sprintf("Error scanning investment data: %v", err), http.StatusInternalServerError)
				return
			}
			investment.Avg = roundToTwoDecimalPlaces(investment.Avg)

			investment.TotInvestment = float64(investment.Qty) * investment.Avg

			investments = append(investments, investment)
		}

		if err := rows.Err(); err != nil {
			http.Error(w, fmt.Sprintf("Error iterating over rows: %v", err), http.StatusInternalServerError)
			return
		}

		if len(investments) == 0 {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"message": "No investments found"})
		}

		response := struct {
			Investments []Investment `json:"investments"`
		}{
			Investments: investments,
		}

		investmentJSON, err := json.Marshal(response)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error marshalling response: %v", err), http.StatusInternalServerError)
			return
		}

		err = redisClient.Set(ctx, cacheKey, string(investmentJSON), 5*time.Minute).Err()
		if err != nil {
			log.Println("Error caching investments in Redis: ", err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(investmentJSON)
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
