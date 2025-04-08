package service

import (
	"context"
	"errors"

	"github.com/darkseear/go-musthave/internal/models"
	"github.com/darkseear/go-musthave/internal/repository"
	"github.com/darkseear/go-musthave/internal/utils"
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

func (b *Balance) UserGetWithdrawals(ctx context.Context, userID int) ([]models.Withdrawal, error) {
	return b.store.GetWithdrawals(ctx, userID)
}

func (b *Balance) UserWithdrawn(ctx context.Context, userID int, orderNumber string, amount float64) error {
	if !utils.ValidLuhn(orderNumber) {
		return errors.New("invalid order number")
	}
	if amount <= 0 {
		return errors.New("negative amount")
	}

	err := b.store.CreateWithdrawal(ctx, userID, orderNumber, amount)
	if err != nil {
		if errors.Is(err, errors.New("insufficient funds")) {
			return errors.New("insufficient funds")
		}
		return err
	}
	return nil
}
