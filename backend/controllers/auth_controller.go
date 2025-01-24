package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/dcode-github/EquiTrack/backend/utils"
)

func Login(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var credentials struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var hashedPassword, userID string
		err := db.QueryRow("SELECT id, password FROM users WHERE username = ?", credentials.Username).Scan(&userID, &hashedPassword)
		if err != nil {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}
		if !utils.CheckPasswordHash(credentials.Password, hashedPassword) {
			http.Error(w, "Invalid password", http.StatusUnauthorized)
			return
		}
		token, err := utils.GenerateJWT(credentials.Username)
		if err != nil {
			http.Error(w, "Error generating token", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"token": token,
			"id":    userID,
		})
	}
}

func Register(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var credentials struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var existingUser string
		err := db.QueryRow("SELECT username FROM users WHERE username = ?", credentials.Username).Scan(&existingUser)
		if err == nil {
			http.Error(w, "Username already taken", http.StatusConflict)
			return
		}

		var existingEmail string
		err = db.QueryRow("SELECT email FROM users WHERE email = ?", credentials.Email).Scan(&existingEmail)
		if err == nil {
			http.Error(w, "Email address already in use", http.StatusConflict)
			return
		}

		hashedPassword, err := utils.HashPassword(credentials.Password)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}

		_, err = db.Exec(
			"INSERT INTO users (username, email, password) VALUES (?, ?, ?)",
			credentials.Username,
			credentials.Email,
			hashedPassword,
		)
		if err != nil {
			http.Error(w, "Error registering user", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
	}
}
