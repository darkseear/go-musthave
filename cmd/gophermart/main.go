package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/darkseear/go-musthave/internal/config"
	logger "github.com/darkseear/go-musthave/internal/logging"
)

func main() {
	if err := run(); err != nil {
		logger.Log.Error("Start server anormal")
		panic(err)
	}
}

func run() error {
	config := config.New()
	LogLevel := config.LogLevel
	if err := logger.Initialize(LogLevel); err != nil {
		return err
	}
	return http.ListenAndServe(config.Address, chi.NewRouter())
}
