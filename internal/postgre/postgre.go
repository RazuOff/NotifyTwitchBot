package postgre

import (
	"fmt"
	"os"

	"github.com/RazuOff/NotifyTwitchBot/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() error {
	host, exists := os.LookupEnv("DBHOST")
	if !exists {
		return fmt.Errorf("env par HOST not found")
	}

	user, exists := os.LookupEnv("USER")
	if !exists {
		return fmt.Errorf("env par HOST not found")
	}

	password, exists := os.LookupEnv("PASSWORD")
	if !exists {
		return fmt.Errorf("env par HOST not found")
	}

	dbname, exists := os.LookupEnv("DBNAME")
	if !exists {
		return fmt.Errorf("env par HOST not found")
	}

	port, exists := os.LookupEnv("DBPORT")
	if !exists {
		return fmt.Errorf("env par HOST not found")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("database connect error: %w", err)
	}

	if err := db.AutoMigrate(&models.Chat{}, &models.Follow{}); err != nil {
		return fmt.Errorf("database connect error: %w", err)
	}

	DB = db

	return nil
}
