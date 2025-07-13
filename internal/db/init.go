package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

/*
InitDatabse initializes and returns a database connection.
Loads environment variables from .env file and establishes connection to PostgreSQL.

Returns:
- *sql.DB: Database connection object
- error: Any error that occurred during initialization

The function:
1. Loads environment variables from .env file
2. Gets the database URL from environment variables
3. Opens a connection to PostgreSQL using the provided URL
4. Pings the database to verify the connection
5. Returns the database connection object

Environment Variables Required:
- DB_URL: PostgreSQL connection string (e.g., "postgres://user:password@localhost:5432/dbname?sslmode=disable")

Possible errors:
- "error connecting to the database": Failed to open database connection
- "error pinging database": Failed to verify database connection

Example DB_URL format:
postgres://username:password@localhost:5432/placement_log?sslmode=disable
*/
func InitDatabse() (*sql.DB, error) {
	_ = godotenv.Load(".env")

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
