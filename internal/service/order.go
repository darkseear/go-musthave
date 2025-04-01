package service

import (
	"context"

	"github.com/darkseear/go-musthave/internal/models"
	"github.com/darkseear/go-musthave/internal/repository"
)

type Order struct {
	store repository.LoyaltyRepository
}

func NewOrder(store repository.LoyaltyRepository) *Order {
	return &Order{store: store}
}

func (o *Order) UserUploadsOrder(ctx context.Context, order models.Order) error {
	return o.store.UploadOrder(ctx, order)
}

func (o *Order) UserGetOrder(ctx context.Context, userID int) ([]models.Order, error) {
	orders, err := o.store.GetOrders(ctx, userID)
	if err != nil {
		return nil, err
	}
	return orders, nil
}
