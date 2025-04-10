package routes

import (
	"database/sql"

	"github.com/dcode-github/EquiTrack/backend/controllers"
	"github.com/dcode-github/EquiTrack/backend/middleware"
	"github.com/gorilla/mux"
)

func Routes(router *mux.Router, db *sql.DB) {
	// controllers.InitDatabase(db)

	router.HandleFunc("/login", controllers.Login(db)).Methods("POST")
	router.HandleFunc("/register", controllers.Register(db)).Methods("POST")
	router.HandleFunc("/individualInvestments", controllers.GetIndvInvestment(db)).Methods("GET")
	router.HandleFunc("/investments", controllers.AddInvestment(db)).Methods("POST")
	// router.HandleFunc("/investments", controllers.GetInvestment(db)).Methods("GET")
	// router.HandleFunc("/investments/{id}", controllers.UpdateInvestment(db)).Methods("PUT")
	router.HandleFunc("/investments", controllers.DeleteInvestment(db)).Methods("DELETE")
	router.HandleFunc("/price", controllers.LivePrice(db)).Methods("GET")
	router.HandleFunc("/priceWebSocket", controllers.LivePriceWebSocket())

	protectedRouter := router.PathPrefix("/").Subrouter()
	protectedRouter.Use(middleware.JWTAuthMiddleware)

	protectedRouter.HandleFunc("/investments", controllers.GetInvestment(db)).Methods("GET")
	// protectedRouter.HandleFunc("/individualInvestments", controllers.GetIndvInvestment(db)).Methods("GET")
	// protectedRouter.HandleFunc("/investments", controllers.AddInvestment(db)).Methods("POST")
	// protectedRouter.HandleFunc("/investments/{id}", controllers.UpdateInvestment(db)).Methods("PUT")
	// protectedRouter.HandleFunc("/investments/{id}", controllers.DeleteInvestment(db)).Methods("DELETE")
	// protectedRouter.HandleFunc("/price", controllers.LivePrice(db)).Methods("GET")
}
