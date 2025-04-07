package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/darkseear/go-musthave/internal/config"
	logger "github.com/darkseear/go-musthave/internal/logging"
	"github.com/darkseear/go-musthave/internal/middleware"
	"github.com/darkseear/go-musthave/internal/models"
	"github.com/darkseear/go-musthave/internal/service"
)

type OrderHandler struct {
	orderServices *service.Order
	cfg           *config.Config
}

func NewOrderHandler(orderServices *service.Order, cfg *config.Config) *OrderHandler {
	return &OrderHandler{orderServices: orderServices, cfg: cfg}
}

func (h *OrderHandler) UploadOrder(w http.ResponseWriter, r *http.Request) {
	authCode := r.Header.Get("Authorization")
	if authCode == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := middleware.GetUserID(r.Header.Get("Authorization"), h.cfg.SecretKey)

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
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (h *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	authCode := r.Header.Get("Authorization")
	if authCode == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := middleware.GetUserID(r.Header.Get("Authorization"), h.cfg.SecretKey)

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
