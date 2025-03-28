package handlers

import (
	"github.com/darkseear/go-musthave/internal/config"
	"github.com/darkseear/go-musthave/internal/middleware"
	"github.com/darkseear/go-musthave/internal/repository"
	"github.com/darkseear/go-musthave/internal/service"
	"github.com/go-chi/chi/v5"
)

type Router struct {
	Router *chi.Mux
	cfg    *config.Config
	store  *repository.Loyalty
}

func Routers(cfg *config.Config, store *repository.Loyalty, auth *service.Auth) *Router {
	r := Router{
		Router: chi.NewRouter(),
		cfg:    cfg,
		store:  store,
	}

	userService := service.NewUser(store)
	userHandler := NewUsersHandler(userService, auth)

	r.Router.Group(func(r chi.Router) {
		middleware.AuthMiddleware(r, auth)                          //middleware для авторизации
		r.Post("/api/users/register", userHandler.UserRegistration) //регистрация пользователя
		r.Post("/api/users/login", userHandler.UserRegistration)    //аутентификация пользователя
	})

	return &r
}
