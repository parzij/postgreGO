package database

import (
	"fmt"
	"log"

	"github.com/parzij/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s connect_timeout=5",
		cfg.HostDB,
		cfg.UserDB,
		cfg.PasswordDB,
		cfg.NameDB,
		cfg.PortDB,
		cfg.SSLmodeDB,
	)

	log.Println("connecting to database...")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)

	log.Println("database connected successfully")
	return db, nil
}