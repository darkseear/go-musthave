package service

import (
	"context"

	"github.com/darkseear/go-musthave/internal/repository"
)

type Balance struct {
	store repository.LoyaltyRepository
}

func NewBalance(store repository.LoyaltyRepository) *Balance {
	return &Balance{store: store}
}

func (b *Balance) UserGetBalance(ctx context.Context, userID int) (float64, error) {
	balance, err := b.store.GetBalance(ctx, userID)
	if err != nil {
		return 0, err
	}
	return balance.Current, nil
}

func (b *Balance) UserUpdateBalance(ctx context.Context, userID int, delta float64) error {
	return b.store.UpdateBalance(ctx, userID, delta)
}

func (b *Balance) UserCreateWithdrawal(ctx context.Context, userID int, orderNumber int, sum float64) error {
	return b.store.CreateWithdrawal(ctx, userID, orderNumber, sum)
}
