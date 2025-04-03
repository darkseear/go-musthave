package service

import (
	"context"
	"errors"

	logger "github.com/darkseear/go-musthave/internal/logging"
	"github.com/darkseear/go-musthave/internal/models"
	"github.com/darkseear/go-musthave/internal/repository"
	"github.com/darkseear/go-musthave/internal/utils"
	"go.uber.org/zap"
)

type Order struct {
	store repository.LoyaltyRepository
}

func NewOrder(store repository.LoyaltyRepository) *Order {
	return &Order{store: store}
}

func (o *Order) UserUploadsOrder(ctx context.Context, order models.Order) error {
	if !utils.ValidLuhn(order.Number) {
		logger.Log.Info("Invalid format Luhn", zap.String("order_number", order.Number))
		return errors.New("invalid order")
	}
	err := o.store.UploadOrder(ctx, order)
	if err != nil {
		return err
	}
	return nil
}

func (o *Order) UserGetOrder(ctx context.Context, userID int) ([]models.Order, error) {
	orders, err := o.store.GetOrders(ctx, userID)
	if err != nil {
		return nil, err
	}
	return orders, nil
}
