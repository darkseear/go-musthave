package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	// "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/darkseear/go-musthave/internal/config"
	"github.com/darkseear/go-musthave/internal/database"
	logger "github.com/darkseear/go-musthave/internal/logging"
	"github.com/darkseear/go-musthave/internal/middleware"
	"github.com/darkseear/go-musthave/internal/service"
)

func main() {
	if err := run(); err != nil {
		logger.Log.Error("Start server anormal")
		log.Fatal(err)
	}
}

func run() error {
	config := config.New()
	LogLevel := config.LogLevel
	if err := logger.Initialize(LogLevel); err != nil {
		return err
	}

	//инициализировать дб
	db, err := database.InitDB(config.Database)
	if err != nil {
		logger.Log.Error("Failed to initialize database")
		log.Fatal(err)
	}
	defer db.Close()

	//миграции
	err = database.RunMigrations(db)
	if err != nil {
		logger.Log.Error("Failed to run migrations")
		log.Fatal(err)
	}

	auth := service.NewAuth(config.SecretKey)
	//создать роутер
	r := chi.NewRouter()
	// r.Use(middleware.Recoverer)

	// r.Group(func(r chi.Router){
	// 	r.Post("/api/users/register", )//решистрация пользователя
	// })

	logger.Log.Info("Running server", zap.String("address", config.Address))
	return http.ListenAndServe(config.Address, middleware.AuthMiddleware(r, auth))
}
