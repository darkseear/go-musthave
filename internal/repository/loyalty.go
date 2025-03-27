package repository

import (
	"context"
	"database/sql"
)

type Loyalty struct {
	db  *sql.DB
	ctx context.Context
}

func NewLoyalty(db *sql.DB, ctx context.Context) *Loyalty {
	return &Loyalty{db: db, ctx: ctx}
}
