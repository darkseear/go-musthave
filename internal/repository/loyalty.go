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
	// Implement the logic to interact with the database and return a user
	query := `INSERT INTO users (login, password_hash) VALUES ($1, $2) RETURNING id, 
	login, password_hash, created_at`

	userUser := &models.User{}
	err := l.db.QueryRowContext(ctx, query, user.Login, user.Password).Scan(
		&userUser.ID,
		&userUser.Login,
		&userUser.PasswordHash,
		&userUser.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Log.Error("No rows found", zap.Error(err))
			return nil, err
		}
		logger.Log.Error("Failed to execute query", zap.Error(err))
		return nil, err
	}

	logger.Log.Info("User created successfully", zap.String("login", userUser.Login))
	return userUser, nil
}

func (l *Loyalty) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	query := `SELECT id, login, password_hash, created_at FROM users WHERE login = $1`
	user := &models.User{}
	err := l.db.QueryRowContext(ctx, query, login).Scan(
		&user.ID,
		&user.Login,
		&user.PasswordHash,
		&user.CreatedAt,
	)
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
