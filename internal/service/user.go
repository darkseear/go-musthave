package service

import (
	"context"

	logger "github.com/darkseear/go-musthave/internal/logging"
	"github.com/darkseear/go-musthave/internal/models"
	"github.com/darkseear/go-musthave/internal/repository"
	"github.com/darkseear/go-musthave/internal/utils"
)

type User struct {
	store repository.LoyaltyRepository
}

func NewUser(store repository.LoyaltyRepository) *User {
	return &User{store: store}
}

func (u *User) UserRegistration(ctx context.Context, login, password string) (*models.User, error) {
	passwordHash := utils.HashPassword(password)
	logger.Log.Info("get passwordHash")
	user, err := u.store.GreaterUser(ctx, models.UserInput{Login: login, Password: passwordHash})
	if err != nil {
		logger.Log.Error("Failed to get user by login")
		return nil, err
	}
	return user, nil
}
