package routes

import (
	"database/sql"

	"github.com/dcode-github/EquiTrack/backend/middleware"
	"github.com/gorilla/mux"
)

func Routes(router *mux.Router, db *sql.DB) {
	// Initialize any required services or configurations using the DB
	controllers.InitDatabase(db)

	router.HandleFunc("/login", controllers.Login(db)).Methods("POST")
	router.HandleFunc("/register", controllers.Register(db)).Methods("POST")

	// Protected routes with JWT authentication middleware
	protectedRouter := router.PathPrefix("/").Subrouter()
	protectedRouter.Use(middleware.JWTAuthMiddleware)

	protectedRouter.HandleFunc("/investments", controllers.GetInvestment(db)).Methods("GET")
	protectedRouter.HandleFunc("/investments", controllers.AddInvestment(db)).Methods("POST")
	protectedRouter.HandleFunc("/investments/{id}", controllers.UpdateInvestment(db)).Methods("PUT")
	protectedRouter.HandleFunc("/investments/{id}", controllers.DeleteInvestment(db)).Methods("DELETE")
}
