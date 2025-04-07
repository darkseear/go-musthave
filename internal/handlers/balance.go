package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/darkseear/go-musthave/internal/config"
	"github.com/darkseear/go-musthave/internal/middleware"
	"github.com/darkseear/go-musthave/internal/service"
)

type BalanceHandler struct {
	balanceService *service.Balance
	cfg            *config.Config
}

func NewBalanceHandler(balanceService *service.Balance, cfg *config.Config) *BalanceHandler {
	return &BalanceHandler{balanceService: balanceService, cfg: cfg}
}

func (b *BalanceHandler) UserGetBalance(w http.ResponseWriter, r *http.Request) {
	authCode := r.Header.Get("Authorization")
	if authCode == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := middleware.GetUserID(r.Header.Get("Authorization"), b.cfg.SecretKey)

	balance, err := b.balanceService.UserGetBalance(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(balance); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	// w.Write([]byte(`{"balance": ` + fmt.Sprintf("%f", balance) + `}`))
}
func (b *BalanceHandler) UserWithdrawBalance(w http.ResponseWriter, r *http.Request) {
	authCode := r.Header.Get("Authorization")
	if authCode == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := middleware.GetUserID(r.Header.Get("Authorization"), b.cfg.SecretKey)

	err := b.balanceService.UserUpdateBalance(r.Context(), userID, -10.0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "success"}`))
}
func (b *BalanceHandler) UserGetWithdrawals(w http.ResponseWriter, r *http.Request) {
	authCode := r.Header.Get("Authorization")
	if authCode == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := middleware.GetUserID(r.Header.Get("Authorization"), b.cfg.SecretKey)
	fmt.Println(userID)
}
