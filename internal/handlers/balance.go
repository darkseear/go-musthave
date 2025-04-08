package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/darkseear/go-musthave/internal/config"
	"github.com/darkseear/go-musthave/internal/middleware"
	"github.com/darkseear/go-musthave/internal/models"
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
	userID, err := middleware.GetUserID(r.Header.Get("Authorization"), b.cfg.SecretKey)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

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
	userID, err := middleware.GetUserID(r.Header.Get("Authorization"), b.cfg.SecretKey)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.ReqWithdraw
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	err = b.balanceService.UserWithdrawn(r.Context(), userID, req.Order, req.Sum)
	if err != nil {
		if err.Error() == "negative amount" {
			http.Error(w, "Negative amount", http.StatusPaymentRequired)
			return
		}
		if err.Error() == "insufficient funds" {
			http.Error(w, "insufficient funds", http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
func (b *BalanceHandler) UserGetWithdrawals(w http.ResponseWriter, r *http.Request) {
	authCode := r.Header.Get("Authorization")
	if authCode == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID, err := middleware.GetUserID(r.Header.Get("Authorization"), b.cfg.SecretKey)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	withdrawals, err := b.balanceService.UserGetWithdrawals(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if len(withdrawals) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if err := json.NewEncoder(w).Encode(withdrawals); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

}
