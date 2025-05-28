package db

import (
	"context"
	"fmt"
)

// ExecTx mengeksekusi fungsi dalam database transaction
func (store *SQLStore) ExecTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.Begin(ctx)
	if err != nil {
		return err
	}
	query := New(tx)
	err = fn(query)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err : %v, rb err : %v", err, rbErr)
		}
		return err
	}
	return tx.Commit(ctx)
}
