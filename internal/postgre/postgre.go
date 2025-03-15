package postgre

import (
	"fmt"

	"github.com/RazuOff/NotifyTwitchBot/internal/config"
	"github.com/RazuOff/NotifyTwitchBot/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB(config *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", config.DBHost, config.User, config.Password, config.DBName, config.DBPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("database connect error: %w", err)
	}

	if err := db.AutoMigrate(&models.Chat{}, &models.Follow{}, &models.StreamerAccount{}); err != nil {
		return nil, fmt.Errorf("database connect error: %w", err)
	}

	return db, nil
}
