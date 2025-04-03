package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	logger "github.com/darkseear/go-musthave/internal/logging"
	"github.com/darkseear/go-musthave/internal/models"
	"github.com/darkseear/go-musthave/internal/service"
)

type OrderHandler struct {
	orderServices *service.Order
}

func NewOrderHandler(orderServices *service.Order) *OrderHandler {
	return &OrderHandler{orderServices: orderServices}
}

func (h *OrderHandler) UploadOrder(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	orderNumber := strings.TrimSpace(string(body))
	if orderNumber == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	err = h.orderServices.UserUploadsOrder(r.Context(), models.Order{Number: orderNumber, UserID: userID, Status: models.Registered})
	if err != nil {
		logger.Log.Error("error upload")
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (h *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)

	orders, err := h.orderServices.UserGetOrder(r.Context(), userID)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(orders); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
}
