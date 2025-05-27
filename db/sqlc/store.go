package db

import (
	"context"
	"database/sql"
	"fmt"
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
	db *sql.DB
}

// NewStore membuat store baru
func NewStore(db *sql.DB) *SQLStore {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// ExecTx mengeksekusi fungsi dalam database transaction
func (store *SQLStore) ExecTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	query := New(tx)
	err = fn(query)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err : %v, rb err : %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}
