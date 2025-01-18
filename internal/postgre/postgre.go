package postgre

import (
	"fmt"

	"github.com/RazuOff/NotifyTwitchBot/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB(host string, user string, password string, dbname string, port string) (*gorm.DB, error) {
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
