package handlers

import (
	"github.com/darkseear/go-musthave/internal/config"
	"github.com/darkseear/go-musthave/internal/repository"
	"github.com/darkseear/go-musthave/internal/service"
	"github.com/go-chi/chi/v5"
)

type Router struct {
	router *chi.Mux
	cfg    *config.Config
	store  *repository.Loyalty
}

func Routers(cfg *config.Config, store *repository.Loyalty, auth *service.Auth) *Router {
	r := Router{
		router: chi.NewRouter(),
		cfg:    cfg,
		store:  store,
	}

	userService := service.NewUser(store)

	userHandler := NewUsersHandler(userService)

	r.router.Group(func(r chi.Router) {
		r.Post("/api/users/register", userHandler.UserRegistration) //решистрация пользователя
	})

	return &r
}
