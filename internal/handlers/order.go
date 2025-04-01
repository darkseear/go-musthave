package handlers

import "github.com/darkseear/go-musthave/internal/service"

type OrderHandler struct {
	orderServices *service.Order
}

func NewOrderHandler(orderServices *service.Order) *OrderHandler {
	return &OrderHandler{orderServices: orderServices}
}
