package routes

import (
	"database/sql"

	"github.com/dcode-github/EquiTrack/backend/controllers"
	"github.com/dcode-github/EquiTrack/backend/middleware"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
)

func Routes(router *mux.Router, db *sql.DB, redisClient *redis.Client) {
	// controllers.InitDatabase(db)

	router.HandleFunc("/login", controllers.Login(db)).Methods("POST")
	router.HandleFunc("/register", controllers.Register(db)).Methods("POST")
	router.HandleFunc("/price", controllers.LivePrice(db)).Methods("GET")
	router.HandleFunc("/priceWebSocket", controllers.LivePriceWebSocket())

	protectedRouter := router.PathPrefix("/").Subrouter()
	protectedRouter.Use(middleware.JWTAuthMiddleware)

	protectedRouter.HandleFunc("/individualInvestments", controllers.GetIndvInvestment(db)).Methods("GET")
	protectedRouter.HandleFunc("/investments", controllers.AddInvestment(db, redisClient)).Methods("POST")
	protectedRouter.HandleFunc("/investments", controllers.GetInvestment(db, redisClient)).Methods("GET")
	// protectedRouter.HandleFunc("/investments/{id}", controllers.UpdateInvestment(db)).Methods("PUT")
	protectedRouter.HandleFunc("/investments", controllers.DeleteInvestment(db, redisClient)).Methods("DELETE")
}
