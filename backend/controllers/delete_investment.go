package controllers

import (
	"database/sql"
	"net/http"
)

func DeleteInvestment(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
