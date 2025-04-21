package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/yokeTH/chat-app-backend/internal/domain"
	"github.com/yokeTH/chat-app-backend/pkg/db"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Unable to load .env file: %s", err)
	}

	portStr := os.Getenv("POSTGRES_PORT")
	port, err := strconv.ParseInt(portStr, 10, 64)
	if err != nil {
		log.Fatalf("Failed to parse port number: %v", err)
	}

	config := db.DBConfig{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     int(port),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBName:   os.Getenv("POSTGRES_NAME"),
		SSLMode:  os.Getenv("POSTGRES_SSL_MODE"),
	}

	db, err := db.New(config)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	if err := db.AutoMigrate(
		&domain.Book{},
		&domain.File{},
		&domain.User{},
		&domain.Conversation{},
		&domain.Message{},
		&domain.Reaction{},
	); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	fmt.Println("Migration completed")
}
