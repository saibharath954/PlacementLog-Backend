package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func InitDatabse() (*sql.DB, error) {
	err := godotenv.Load(".env")

	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	dbUrl := os.Getenv("DB_URL")

	conn, err := sql.Open("postgres", dbUrl)

	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}

	if err = conn.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging database: %v", err)
	}

	log.Println("Connected to the database")

	return conn, nil
}
