package repository

import (
	"context"
	"database/sql"
	"errors"

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
	query := `SELECT id, login, password_hash, created_at 
	FROM users 
	WHERE login = $1`
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
	var isOrderExists sql.NullInt64
	err := l.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT user_id FROM orders WHERE number = $1)`, order.Number).Scan(&isOrderExists)
	if err != nil && err != sql.ErrNoRows {
		logger.Log.Error("Failed to check if order exists", zap.Error(err))
		return err
	}
	if err != sql.ErrNoRows {
		if isOrderExists.Valid && isOrderExists.Int64 == int64(order.UserID) {
			return errors.New("order already exists")
		}
		return errors.New("order does not exist")
	}

	query := `INSERT INTO orders (number, user_id, status) VALUES ($1, $2, $3)`
	_, err = l.db.ExecContext(ctx, query, order.Number, order.UserID, order.Status)
	if err != nil {
		logger.Log.Error("Failed to insert order", zap.Error(err))
		return err
	}

	logger.Log.Info("Order uploaded successfully", zap.String("orderNumber", order.Number))
	return nil
}

func (l *Loyalty) GetOrder(ctx context.Context, orderNumber string) (*models.Order, error) {
	query := `SELECT number, user_id, status, accrual, uploaded_at 
	FROM orders 
	WHERE number = $1`

	order := &models.Order{}
	err := l.db.QueryRowContext(ctx, query, orderNumber).Scan(&order.Number, &order.UserID, &order.Status, &order.Accrual, &order.UploadedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Log.Error("No rows found", zap.Error(err))
			return nil, err
		}
		logger.Log.Error("Failed to execute query", zap.Error(err))
		return nil, err
	}
	return order, nil
}

func (l *Loyalty) GetOrders(ctx context.Context, userID int) ([]models.Order, error) {
	query := `SELECT number, user_id, status, accrual, uploaded_at
	 FROM orders 
	 WHERE user_id = $1
	 ORDER BY uploaded_at DESC`

	rows, err := l.db.QueryContext(ctx, query, userID)
	if err != nil {
		logger.Log.Error("Failed to execute query", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		order := models.Order{}
		err := rows.Scan(&order.Number, &order.UserID, &order.Status, &order.Accrual, &order.UploadedAt)
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

func (l *Loyalty) UpdateOrderStatus(ctx context.Context, orderNumber string, status models.Status, accrual float64) error {
	query := `UPDATE orders 
	SET status = $1, accrual = $2 
	WHERE number = $3`

	var accrualValue interface{}
	if accrual > 0 {
		accrualValue = accrual
	} else {
		accrualValue = nil
	}
	resStatus, err := l.db.ExecContext(ctx, query, status, accrualValue, orderNumber)
	if err != nil {
		logger.Log.Error("Failed to update order status", zap.Error(err))
		return err
	}

	rowsAffected, err := resStatus.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("no rows updated")
	}
	logger.Log.Info("Order status updated successfully")
	return nil
}

func (l *Loyalty) GetBalance(ctx context.Context, userID int) (*models.Balance, error) {
	query := `SELECT current_balance, withdrawn_balance 
	FROM user 
	WHERE user_id = $1`
	balance := &models.Balance{}
	err := l.db.QueryRowContext(ctx, query, userID).Scan(&balance.Current, &balance.Withdrawn)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Log.Error("No rows found", zap.Error(err))
			return &models.Balance{}, err
		}
		logger.Log.Error("Failed to execute query", zap.Error(err))
		return nil, err
	}

	logger.Log.Info("Balance retrieved successfully", zap.Int("userID", userID))
	return balance, nil
}

func (l *Loyalty) UpdateBalance(ctx context.Context, userID int, delta float64) error {
	tx, err := l.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var currentBalance float64
	query := `SELECT current_balance FROM user WHERE user_id = $1 FOR UPDATE`
	err = tx.QueryRowContext(ctx, query, userID).Scan(&currentBalance)
	if err != nil {
		return err
	}

	if delta < 0 && currentBalance+delta < 0 {
		return errors.New("insufficient balance")
	}

	query = `UPDATE user SET current_balance = current_balance + $1 WHERE user_id = $2`
	_, err = l.db.ExecContext(ctx, query, delta, userID)
	if err != nil {
		logger.Log.Error("Failed to update balance", zap.Error(err))
		return err
	}

	logger.Log.Info("Balance updated successfully", zap.Int("userID", userID), zap.Float64("delta", delta))
	return tx.Commit()
}

func (l *Loyalty) CreateWithdrawal(ctx context.Context, userID int, orderNumber int, sum float64) error {
	tx, err := l.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var currentBalance float64
	query := `SELECT current_balance FROM user WHERE user_id = $1 FOR UPDATE`
	err = tx.QueryRowContext(ctx, query, userID).Scan(&currentBalance)
	if err != nil {
		return err
	}
	if currentBalance < sum {
		return errors.New("insufficient balance")
	}

	query = `INSERT INTO withdrawals (user_id, order_number, sum) VALUES ($1, $2, $3)`
	_, err = l.db.ExecContext(ctx, query, userID, orderNumber, sum)
	if err != nil {
		logger.Log.Error("Failed to create withdrawal", zap.Error(err))
		return err
	}
	logger.Log.Info("Withdrawal created successfully", zap.Int("userID", userID), zap.Int("orderNumber", orderNumber), zap.Float64("sum", sum))
	_, err = tx.ExecContext(ctx, `UPDATE user SET current_balance = current_balance - $1 WHERE user_id = $2`, sum, userID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (l *Loyalty) GetWithdrawals(ctx context.Context, userID int) ([]models.Withdrawal, error) {
	query := `SELECT order_number, user_id, sum, processed_at
	 FROM withdrawals 
	 WHERE user_id = $1
	 ORDER BY processed_at DESC`
	rows, err := l.db.QueryContext(ctx, query, userID)
	if err != nil {
		logger.Log.Error("Failed to execute query", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var withdrawals []models.Withdrawal
	for rows.Next() {
		withdrawal := models.Withdrawal{}
		err := rows.Scan(&withdrawal.UserID, &withdrawal.OrderNumber, &withdrawal.Sum, &withdrawal.ProcessedAt)
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

func (l *Loyalty) Ping(ctx context.Context) error {
	err := l.db.PingContext(ctx)
	if err != nil {
		logger.Log.Error("Failed to ping database", zap.Error(err))
		return err
	}
	logger.Log.Info("Database ping successful")
	return nil
}

func (l *Loyalty) Close() error {
	err := l.db.Close()
	if err != nil {
		logger.Log.Error("Failed to close database connection", zap.Error(err))
		return err
	}
	logger.Log.Info("Database connection closed successfully")
	return nil
}
