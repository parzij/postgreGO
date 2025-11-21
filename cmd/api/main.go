package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/parzij/internal/config"
	"github.com/parzij/internal/database"
	"github.com/parzij/internal/users"
	"gorm.io/gorm"
)


func main() {
	_ = godotenv.Load()
	
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	db, err := database.New(cfg)
	if err != nil {
		log.Fatal("cannot connect to database:", err)
	}

	if err := migrate(db); err != nil {
		log.Fatal("cannot migrate:", err)
	}

	r := chi.NewRouter()

	userHandler := users.Newhandler(db)
	userHandler.RegisterRoutes(r)

	log.Println("server listening on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}

func migrate(db *gorm.DB) error {
	return db.AutoMigrate(&users.User{})
}
