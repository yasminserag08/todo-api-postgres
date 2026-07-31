package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv" // to avoid hardcoding the password
	_ "github.com/lib/pq"
)

func connect() (*sql.DB, error) {
	godotenv.Load() // reads .env into environment variables

	// in case there's no .env file
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using system environment variables")
	}

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("postgres", connStr) // validate that connStr looks correct & set up connection pool

	if err != nil {
		return nil, err
	}

	// actual attempt to connect (error is returned if password is incorrect, db doesn't exist, psql is not running, etc)
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
