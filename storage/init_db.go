package storage

import (
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

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

	gormDB, err := gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{},
	)
	if err != nil {
		log.Fatal("[DB] open error:", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Fatal("[DB] get sql.DB error:", err)
	}

	// Connection pool config (AMAN)
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		log.Fatal("[DB] ping error:", err)
	}

	DB = gormDB
	log.Println("[DB] PostgreSQL connected")
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
