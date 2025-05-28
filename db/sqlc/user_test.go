package db

import (
	"context"
	"simple-bank/util"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func CreateRandomUser(t *testing.T) User {

	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testStore.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	assert.Equal(t, arg.Username, user.Username)
	assert.Equal(t, arg.HashedPassword, user.HashedPassword)
	assert.Equal(t, arg.FullName, user.FullName)
	assert.Equal(t, arg.Email, user.Email)
	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)
	return user
}

func TestCreateUser(t *testing.T) {
	CreateRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := CreateRandomUser(t)
	user2, err := testStore.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	assert.Equal(t, user1.Username, user2.Username)
	assert.Equal(t, user1.HashedPassword, user2.HashedPassword)
	assert.Equal(t, user1.FullName, user2.FullName)
	assert.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

func TestUpdateUserOnlyFullName(t *testing.T) {
	oldUser := CreateRandomUser(t)
	newFullName := util.RandomOwner()
	newUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{FullName: pgtype.Text{
		String: newFullName, Valid: true,
	}, Username: oldUser.Username})

	require.NoError(t, err)
	require.NotEmpty(t, newUser)
	assert.Equal(t, oldUser.Username, newUser.Username)
	assert.Equal(t, oldUser.Email, newUser.Email)
	assert.Equal(t, oldUser.HashedPassword, newUser.HashedPassword)
	assert.NotEqual(t, oldUser.FullName, newUser.FullName)
}

func TestUpdateUserOnlyHashPassword(t *testing.T) {
	oldUser := CreateRandomUser(t)
	newPassword := util.RandomString(10)
	newHashPassword, err := util.HashPassword(newPassword)
	require.NoError(t, err)
	newUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{HashedPassword: pgtype.Text{
		String: newHashPassword, Valid: true,
	}, Username: oldUser.Username})

	require.NoError(t, err)
	require.NotEmpty(t, newUser)
	assert.Equal(t, oldUser.Username, newUser.Username)
	assert.Equal(t, oldUser.Email, newUser.Email)
	assert.NotEqual(t, oldUser.HashedPassword, newUser.HashedPassword)
	assert.Equal(t, oldUser.FullName, newUser.FullName)
}

func TestUpdateUserOnlyEmail(t *testing.T) {
	oldUser := CreateRandomUser(t)
	newEmail := util.RandomEmail()
	newUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{Email: pgtype.Text{
		String: newEmail, Valid: true,
	}, Username: oldUser.Username})

	require.NoError(t, err)
	require.NotEmpty(t, newUser)
	assert.Equal(t, oldUser.Username, newUser.Username)
	assert.NotEqual(t, oldUser.Email, newUser.Email)
	assert.Equal(t, oldUser.HashedPassword, newUser.HashedPassword)
	assert.Equal(t, oldUser.FullName, newUser.FullName)
}
