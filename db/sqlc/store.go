package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// store menyediakan semua fungsi untuk mengeksekusi query database dan transaction
type Store interface {
	Querier
	TransferTx(ctx context.Context, arg CreateTransferParams) (TransferTxResult, error)
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)
	VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error)
}

// store menyediakan semua fungsi untuk mengeksekusi SQL query database dan transaction
type SQLStore struct {
	*Queries
	db *pgxpool.Pool
}

// NewStore membuat store baru
func NewStore(db *pgxpool.Pool) *SQLStore {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}
