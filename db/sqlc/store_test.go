package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)
	fmt.Println("Before>>", account1.Balance, account2.Balance)
	// run concurent transfer transaction
	n := 2
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx %d", i+1)
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := store.TransferTx(ctx, CreateTransferParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()

	}

	// check result
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
		result := <-results
		require.NotEmpty(t, result)

		// check Transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check From Entry
		fromEntries := result.FromEntry
		require.NotEmpty(t, fromEntries)
		require.Equal(t, account1.ID, fromEntries.AccountID)
		require.Equal(t, -amount, fromEntries.Amount)
		require.NotZero(t, fromEntries.ID)
		require.NotZero(t, fromEntries.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntries.ID)
		require.NoError(t, err)

		// check To Entry
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check account
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		// check account balance
		fmt.Println(">> tx : ", fromAccount.Balance, toAccount.Balance)
		dif1 := account1.Balance - fromAccount.Balance
		dif2 := toAccount.Balance - account2.Balance

		require.Equal(t, dif1, dif2)
		require.True(t, dif1 > 0)
		require.True(t, dif1%amount == 0)

		k := int(dif1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true

	}
	// check final update balance

	updatedAccount1, err := testQuery.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQuery.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println("After>> ", updatedAccount1.Balance, updatedAccount2.Balance)

	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance-int64(n)*amount, updatedAccount2.Balance)
}
