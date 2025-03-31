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
