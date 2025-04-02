package handlers

import (
	"net/http"

	"github.com/darkseear/go-musthave/internal/service"
)

type OrderHandler struct {
	orderServices *service.Order
}

func NewOrderHandler(orderServices *service.Order) *OrderHandler {
	return &OrderHandler{orderServices: orderServices}
}

func (h *OrderHandler) UploadOrder(w http.ResponseWriter, r *http.Request) {
	// Handler logic for uploading order
}

func (h *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	// Handler logic for getting orders
}
