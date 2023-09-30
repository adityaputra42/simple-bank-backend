package db

import (
	"context"
	"database/sql"
	"simple-bank/util"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func CreateRandomTransfer(t *testing.T) Transfer {
	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)
	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomBalance(),
	}

	Transfer, err := testQuery.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, Transfer)
	assert.Equal(t, arg.FromAccountID, Transfer.FromAccountID)
	assert.Equal(t, arg.ToAccountID, Transfer.ToAccountID)
	assert.Equal(t, arg.Amount, Transfer.Amount)

	require.NotZero(t, Transfer.ID)
	require.NotZero(t, Transfer.CreatedAt)
	return Transfer
}

func TestCreateTransfer(t *testing.T) {
	CreateRandomTransfer(t)
}

func TestGetTransfer(t *testing.T) {
	transfer1 := CreateRandomTransfer(t)
	transfer2, err := testQuery.GetTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)
	assert.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	assert.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	assert.Equal(t, transfer1.Amount, transfer2.Amount)
	require.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, time.Second)
}

func TestUpdateTransfer(t *testing.T) {
	Transfer1 := CreateRandomTransfer(t)
	account := CreateRandomAccount(t)
	arg := UpdateTransferParams{
		ID:          Transfer1.ID,
		ToAccountID: account.ID,
		Amount:      util.RandomBalance(),
	}
	Transfer2, err := testQuery.UpdateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, Transfer2)
	assert.Equal(t, arg.ID, Transfer2.ID)
	assert.Equal(t, arg.ToAccountID, Transfer2.ToAccountID)
	assert.Equal(t, arg.Amount, Transfer2.Amount)
	require.WithinDuration(t, Transfer1.CreatedAt, Transfer2.CreatedAt, time.Second)

}

func TestDeleteTransfer(t *testing.T) {
	Transfer1 := CreateRandomTransfer(t)
	err := testQuery.DeleteTransfer(context.Background(), Transfer1.ID)
	require.NoError(t, err)
	Transfer2, err := testQuery.GetTransfer(context.Background(), Transfer1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, Transfer2)

}

func TestListTransfer(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateRandomTransfer(t)
	}

	arg := ListTransferParams{
		Limit:  5,
		Offset: 5,
	}
	Transfers, err := testQuery.ListTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, Transfers, 5)
	for _, Transfer := range Transfers {
		require.NotEmpty(t, Transfer)
	}
}
