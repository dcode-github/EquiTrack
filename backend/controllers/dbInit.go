package controllers

import (
	"database/sql"
	"log"
)

var DB *sql.DB

func InitDatabase(db *sql.DB) {
	DB = db

	// Optional: Run migrations or ensure tables exist
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id INT AUTO_INCREMENT PRIMARY KEY,
            username VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL,
            password VARCHAR(255) NOT NULL
        );

        CREATE TABLE IF NOT EXISTS investments (
            id INT AUTO_INCREMENT PRIMARY KEY,
            user_id INT NOT NULL,
            stock VARCHAR(255) NOT NULL,
            units INT NOT NULL,
            price FLOAT NOT NULL,
            FOREIGN KEY (user_id) REFERENCES users(id)
        );
    `)

	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	log.Println("Database initialized successfully!")
}
