package repository

import (
	"context"

	"github.com/darkseear/go-musthave/internal/models"
)

type LoyaltyRepository interface {
	GreaterUser(ctx context.Context, user models.UserInput) (*models.User, error)
	// GetUserByLogin(ctx context.Context, login string) (*models.User, error)
}
