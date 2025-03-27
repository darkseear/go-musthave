package handlers

import (
	"net/http"

	logger "github.com/darkseear/go-musthave/internal/logging"
	"github.com/darkseear/go-musthave/internal/service"
	"go.uber.org/zap"
)

type UsersHandler struct {
	userService *service.User
}

func NewUsersHandler(userService *service.User) *UsersHandler {
	return &UsersHandler{userService: userService}
}

func (uh *UsersHandler) UserRegistration(w http.ResponseWriter, r *http.Request) {
	user, err := uh.userService.UserRegistration(r.Context(), r.FormValue("login"), r.FormValue("password"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logger.Log.Info("User registered", zap.Int("userID", user.ID))
}
