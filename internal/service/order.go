package service

import (
	"context"
	"errors"
	"fmt"

	logger "github.com/darkseear/go-musthave/internal/logging"
	"github.com/darkseear/go-musthave/internal/models"
	"github.com/darkseear/go-musthave/internal/processor"
	"github.com/darkseear/go-musthave/internal/repository"
	"github.com/darkseear/go-musthave/internal/utils"
	"go.uber.org/zap"
)

type Order struct {
	store          repository.LoyaltyRepository
	orderProcessor *processor.Order
}

func NewOrder(store repository.LoyaltyRepository) *Order {
	return &Order{store: store, orderProcessor: processor.NewOrder(nil, store)}
}

func (o *Order) UserUploadsOrder(ctx context.Context, order models.Order) error {
	if !utils.ValidLuhn(order.Number) {
		fmt.Println(utils.ValidLuhn(order.Number))
		logger.Log.Info("Invalid format Luhn", zap.String("order_number", order.Number))
		return errors.New("invalid order")
	}
	err := o.store.UploadOrder(ctx, order)
	if err != nil {
		return err
	}
	o.orderProcessor.AddOrder(ctx, order.Number)
	return nil
}

func (o *Order) UserGetOrder(ctx context.Context, userID int) ([]models.Order, error) {
	return o.store.GetOrders(ctx, userID)
}
