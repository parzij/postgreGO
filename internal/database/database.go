package database

import (
	"fmt"

	"github.com/parzij/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.HostDB,
		cfg.UserDB,
		cfg.PasswordDB,
		cfg.NameDB,
		cfg.PortDB,
		cfg.SSLmodeDB,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
