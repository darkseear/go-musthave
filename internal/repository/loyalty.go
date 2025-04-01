package repository

import (
	"context"
	"database/sql"

	logger "github.com/darkseear/go-musthave/internal/logging"
	"github.com/darkseear/go-musthave/internal/models"
	"go.uber.org/zap"
)

type Loyalty struct {
	db  *sql.DB
	ctx context.Context
}

func NewLoyalty(db *sql.DB, ctx context.Context) *Loyalty {
	return &Loyalty{db: db, ctx: ctx}
}

func (l *Loyalty) GreaterUser(ctx context.Context, user models.UserInput) (*models.User, error) {
	query := `INSERT INTO users (login, password_hash) VALUES ($1, $2) RETURNING id, 
	login, password_hash, created_at`
	userUser := &models.User{}
	userDB, err := l.UserDB(ctx, userUser, query, user.Login, user.Password)
	if err != nil {
		logger.Log.Error("Failed to insert user", zap.Error(err))
		return nil, err
	}
	return userDB, nil
}

func (l *Loyalty) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	query := `SELECT id, login, password_hash, created_at FROM users WHERE login = $1`
	user := &models.User{}
	UserDB, err := l.UserDB(ctx, user, query, login, "")
	if err != nil {
		logger.Log.Error("Failed to get user by login", zap.Error(err))
		return nil, err
	}
	return UserDB, nil
}

func (l *Loyalty) UserDB(ctx context.Context, user *models.User, query, login, password string) (*models.User, error) {
	var err error
	if password == "" {
		err = l.db.QueryRowContext(ctx, query, login).Scan(
			&user.ID,
			&user.Login,
			&user.PasswordHash,
			&user.CreatedAt,
		)
	} else {
		err = l.db.QueryRowContext(ctx, query, login, password).Scan(
			&user.ID,
			&user.Login,
			&user.PasswordHash,
			&user.CreatedAt,
		)
	}
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Log.Error("No rows found", zap.Error(err))
			return nil, err
		}
		logger.Log.Error("Failed to execute query", zap.Error(err))
		return nil, err
	}

	logger.Log.Info("User retrieved successfully", zap.String("login", user.Login))
	return user, nil
}

func (l *Loyalty) UploadOrder(ctx context.Context, order models.Order) error {
	query := `INSERT INTO orders (number, user_id, status, accrual) VALUES ($1, $2, $3, $4)`
	_, err := l.db.ExecContext(ctx, query, order.Number, order.UserID, order.Status, order.Accrual)
	if err != nil {
		logger.Log.Error("Failed to insert order", zap.Error(err))
		return err
	}

	logger.Log.Info("Order uploaded successfully", zap.Int("orderNumber", order.Number))
	return nil
}

func (l *Loyalty) GetOrders(ctx context.Context, userID int) ([]models.Order, error) {
	query := `SELECT number, user_id, status, accrual FROM orders WHERE user_id = $1`
	rows, err := l.db.QueryContext(ctx, query, userID)
	if err != nil {
		logger.Log.Error("Failed to execute query", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		order := models.Order{}
		err := rows.Scan(&order.Number, &order.UserID, &order.Status, &order.Accrual)
		if err != nil {
			logger.Log.Error("Failed to scan row", zap.Error(err))
			return nil, err
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		logger.Log.Error("Error occurred during row iteration", zap.Error(err))
		return nil, err
	}

	logger.Log.Info("Orders retrieved successfully", zap.Int("userID", userID))
	return orders, nil
}

func (l *Loyalty) GetBalance(ctx context.Context, userID int) (*models.Balance, error) {
	query := `SELECT current_balance, withdrawn_balance FROM user WHERE user_id = $1`
	balance := &models.Balance{}
	err := l.db.QueryRowContext(ctx, query, userID).Scan(&balance.Current, &balance.Withdrawn)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Log.Error("No rows found", zap.Error(err))
			return nil, err
		}
		logger.Log.Error("Failed to execute query", zap.Error(err))
		return nil, err
	}

	logger.Log.Info("Balance retrieved successfully", zap.Int("userID", userID))
	return balance, nil
}

func (l *Loyalty) UpdateBalance(ctx context.Context, userID int, delta float64) error {
	query := `UPDATE user SET current_balance = current_balance + $1 WHERE user_id = $2`
	_, err := l.db.ExecContext(ctx, query, delta, userID)
	if err != nil {
		logger.Log.Error("Failed to update balance", zap.Error(err))
		return err
	}

	logger.Log.Info("Balance updated successfully", zap.Int("userID", userID), zap.Float64("delta", delta))
	return nil
}

func (l *Loyalty) CreateWithdrawal(ctx context.Context, userID int, orderNumber int, sum float64) error {
	query := `INSERT INTO withdrawals (user_id, order_number, sum) VALUES ($1, $2, $3)`
	_, err := l.db.ExecContext(ctx, query, userID, orderNumber, sum)
	if err != nil {
		logger.Log.Error("Failed to create withdrawal", zap.Error(err))
		return err
	}

	logger.Log.Info("Withdrawal created successfully", zap.Int("userID", userID), zap.Int("orderNumber", orderNumber), zap.Float64("sum", sum))
	return nil
}

func (l *Loyalty) GetWithdrawals(ctx context.Context, userID int) ([]models.Withdrawal, error) {
	query := `SELECT order_number, sum FROM withdrawals WHERE user_id = $1`
	rows, err := l.db.QueryContext(ctx, query, userID)
	if err != nil {
		logger.Log.Error("Failed to execute query", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var withdrawals []models.Withdrawal
	for rows.Next() {
		withdrawal := models.Withdrawal{}
		err := rows.Scan(&withdrawal.OrderNumber, &withdrawal.ProcessedAt, &withdrawal.Sum, &withdrawal.UserID)
		if err != nil {
			logger.Log.Error("Failed to scan row", zap.Error(err))
			return nil, err
		}
		withdrawals = append(withdrawals, withdrawal)
	}

	if err := rows.Err(); err != nil {
		logger.Log.Error("Error occurred during row iteration", zap.Error(err))
		return nil, err
	}

	logger.Log.Info("Withdrawals retrieved successfully", zap.Int("userID", userID))
	return withdrawals, nil
}
