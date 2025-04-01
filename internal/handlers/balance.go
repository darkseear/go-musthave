package handlers

import (
	"fmt"
	"net/http"

	"github.com/darkseear/go-musthave/internal/service"
)

type BalanceHandler struct {
	balanceService *service.Balance
}

func NewBalanceHandler(balanceService *service.Balance) *BalanceHandler {
	return &BalanceHandler{balanceService: balanceService}
}

func (b *BalanceHandler) UserGetBalance(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)                          // Replace with actual user ID from context
	balance, err := b.balanceService.UserGetBalance(r.Context(), userID) // Replace 1 with actual user ID from context
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"balance": ` + fmt.Sprintf("%f", balance) + `}`))
}
func (b *BalanceHandler) UserWithdrawBalance(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)                           // Replace with actual user ID from context
	err := b.balanceService.UserUpdateBalance(r.Context(), userID, -10.0) // Replace 1 with actual user ID from context
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "success"}`))
}
func (b *BalanceHandler) UserGetWithdrawals(w http.ResponseWriter, r *http.Request) {
	// Handler logic for getting user withdrawals
}
