package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// InitDB initializes PostgreSQL connection
func InitDB() {

	dsn := os.Getenv("DATABASE_URL")

	if dsn == "" {
		// fallback local dev
		dsn = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			getEnv("DB_HOST", "localhost"),
			getEnv("DB_PORT", "5432"),
			getEnv("DB_USER", "postgres"),
			getEnv("DB_PASSWORD", "123123"),
			getEnv("DB_NAME", "doc_tracker"),
		)
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("[DB] open error:", err)
	}

	// Connection pool config (AMAN)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(30 * time.Minute)

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatal("[DB] ping error:", err)
	}

	DB = db
	log.Println("[DB] PostgreSQL connected")
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
