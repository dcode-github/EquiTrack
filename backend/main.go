package main

import (
	"log"
	"net/http"
	"os"

	"github.com/dcode-github/EquiTrack/backend/config"
	"github.com/dcode-github/EquiTrack/backend/routes"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	db := config.ConnectDB()
	defer db.Close()

	redisClient := config.ConnectRedis()
	defer redisClient.Close()

	router := mux.NewRouter()

	routes.Routes(router, db, redisClient)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
