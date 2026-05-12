package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	_ "github.com/parzij/docs"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/parzij/internal/config"
	"github.com/parzij/internal/database"
	"github.com/parzij/internal/users"
	"gorm.io/gorm"
)

// @title           postgreGO API
// @version         1.0
// @description     Учебный backend-сервис на Go + GORM + PostgreSQL для работы с пользователями
// @host            localhost:9090
// @BasePath        /
// @schemes         http
func main() {
	log.Println("loading .env...")
	_ = godotenv.Load()

	log.Println("loading config...")
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	log.Println("initializing database...")
	db, err := database.New(cfg)
	if err != nil {
		log.Fatal("cannot connect to database:", err)
	}

	log.Println("running migrations...")
	if err := migrate(db); err != nil {
		log.Fatal("cannot migrate:", err)
	}

	r := chi.NewRouter()

	// Swagger UI
	r.Route("/swagger", func(r chi.Router) {
		r.Get("/*", httpSwagger.Handler(
			httpSwagger.URL("/swagger/doc.json"),
			httpSwagger.DeepLinking(true),
			httpSwagger.DocExpansion("list"),
		))
	})

	userHandler := users.Newhandler(db)
	userHandler.RegisterRoutes(r)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = ":9090"
	}

	log.Printf("server starting on %s", port)
	log.Println("Swagger UI: http://localhost" + port + "/swagger/index.html")
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatal(err)
	}
}

func migrate(db *gorm.DB) error {
	return db.AutoMigrate(&users.User{})
}