package db

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mahanth/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)

	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.CreatedAt, user2.CreatedAt)
	require.True(t, user1.LastPasswordChangedAt.Time.IsZero())
	require.True(t, user2.LastPasswordChangedAt.Time.IsZero())
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
}

func TestUpdateUserFullName(t *testing.T) {
	user1 := createRandomUser(t)

	arg := UpdateUserParams{
		Username: user1.Username,
		FullName: pgtype.Text{String: "updated_" + user1.FullName, Valid: true},
	}
	updatedUser, err := testQueries.UpdateUser(context.Background(), arg)

	require.NoError(t, err)
	require.Equal(t, updatedUser.FullName, "updated_"+user1.FullName)
	require.Equal(t, updatedUser.Username, user1.Username)
	require.Equal(t, updatedUser.Email, updatedUser.Email)
	require.Equal(t, updatedUser.HashedPassword, user1.HashedPassword)

}

func TestUpdateUserEmail(t *testing.T) {
	user1 := createRandomUser(t)

	arg := UpdateUserParams{
		Username: user1.Username,
		Email:    pgtype.Text{String: "updated_" + user1.Email, Valid: true},
	}
	updatedUser, err := testQueries.UpdateUser(context.Background(), arg)

	require.NoError(t, err)
	require.Equal(t, updatedUser.Email, "updated_"+user1.Email)
	require.Equal(t, updatedUser.Username, user1.Username)
	require.Equal(t, updatedUser.FullName, updatedUser.FullName)
	require.Equal(t, updatedUser.HashedPassword, user1.HashedPassword)

}

func createRandomUser(t *testing.T) User {
	password := util.RandomString(6)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)
	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	require.True(t, user.LastPasswordChangedAt.Time.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}
