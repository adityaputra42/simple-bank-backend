package db

import (
	"context"
	"simple-bank/util"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func CreateRandomEntry(t *testing.T) Entry {
	account1 := CreateRandomAccount(t)
	arg := CreateEntryParams{
		AccountID: account1.ID,
		Amount:    util.RandomBalance(),
	}

	Entry, err := testStore.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, Entry)
	assert.Equal(t, arg.AccountID, Entry.AccountID)
	assert.Equal(t, arg.Amount, Entry.Amount)

	require.NotZero(t, Entry.ID)
	require.NotZero(t, Entry.CreatedAt)
	return Entry
}

func TestCreateEntry(t *testing.T) {
	CreateRandomEntry(t)
}

func TestGetEntry(t *testing.T) {
	Entry1 := CreateRandomEntry(t)
	Entry2, err := testStore.GetEntry(context.Background(), Entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, Entry2)
	assert.Equal(t, Entry1.AccountID, Entry2.AccountID)
	assert.Equal(t, Entry1.Amount, Entry2.Amount)
	require.WithinDuration(t, Entry1.CreatedAt, Entry2.CreatedAt, time.Second)
}

func TestUpdateEntry(t *testing.T) {
	Entry1 := CreateRandomEntry(t)
	arg := UpdateEntryParams{
		ID:     Entry1.ID,
		Amount: util.RandomBalance(),
	}
	Entry2, err := testStore.UpdateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, Entry2)
	assert.Equal(t, Entry1.AccountID, Entry2.AccountID)
	assert.Equal(t, arg.Amount, Entry2.Amount)
	require.WithinDuration(t, Entry1.CreatedAt, Entry2.CreatedAt, time.Second)

}

func TestDeleteEntry(t *testing.T) {
	Entry1 := CreateRandomEntry(t)
	err := testStore.DeleteEntry(context.Background(), Entry1.ID)
	require.NoError(t, err)
	Entry2, err := testStore.GetEntry(context.Background(), Entry1.ID)
	require.Error(t, err)
	require.EqualError(t, err, ErrRecordNotFound.Error())
	require.Empty(t, Entry2)

}

func TestListEntry(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateRandomEntry(t)
	}

	arg := ListEntryParams{
		Limit:  5,
		Offset: 5,
	}
	Entrys, err := testStore.ListEntry(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, Entrys, 5)
	for _, Entry := range Entrys {
		require.NotEmpty(t, Entry)
	}
}
