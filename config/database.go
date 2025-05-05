package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ No .env file found, using environment variables")
	}

	host := getEnv("DB_HOST", "127.0.0.1")
	portStr := getEnv("DB_PORT", "3306")
	user := getEnv("DB_USER", "root")
	password := getEnv("DB_PASSWORD", "")
	name := getEnv("DB_NAME", "project_db")

	// If running inside Docker, use host.docker.internal (for Docker Desktop on macOS/Windows)
	if _, err := os.Stat("/.dockerenv"); err == nil {
		// Running inside Docker
		host = "host.docker.internal"
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("❌ Invalid DB_PORT: %v", err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local", user, password, host, port, name)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		log.Fatal("❌ Failed to connect to DB: ", err)
	}

	log.Println("✅ DB Connected:", name)
	return db
}

// getEnv gets env variable or fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
