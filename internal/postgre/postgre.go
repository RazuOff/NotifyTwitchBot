package postgre

import (
	"fmt"
	"os"

	"github.com/RazuOff/NotifyTwitchBot/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init() (*gorm.DB, error) {
	host, exists := os.LookupEnv("DBHOST")
	if !exists {
		return nil, fmt.Errorf("env par HOST not found")
	}

	user, exists := os.LookupEnv("USER")
	if !exists {
		return nil, fmt.Errorf("env par HOST not found")
	}

	password, exists := os.LookupEnv("PASSWORD")
	if !exists {
		return nil, fmt.Errorf("env par HOST not found")
	}

	dbname, exists := os.LookupEnv("DBNAME")
	if !exists {
		return nil, fmt.Errorf("env par HOST not found")
	}

	port, exists := os.LookupEnv("DBPORT")
	if !exists {
		return nil, fmt.Errorf("env par HOST not found")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("database connect error: %w", err)
	}

	if err := db.AutoMigrate(&models.Chat{}, &models.Follow{}); err != nil {
		return nil, fmt.Errorf("database connect error: %w", err)
	}

	return db, nil
}
