package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/Yorshik/Go-Final/internal/models"
)

var DB *gorm.DB

func ConnectDatabase() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		log.Fatalf("DB_NAME is not set in the .env file")
	}
	dsn := fmt.Sprintf("%s.db", dbName)
	DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	if err := DB.AutoMigrate(&models.User{}, &models.Expression{}); err != nil {
		log.Fatalf("Ошибка миграции базы данных: %v", err)
	}
	fmt.Println("Успешное подключение к базе данных SQLite!")
}
