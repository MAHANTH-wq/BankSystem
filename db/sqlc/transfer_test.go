package db

import (
	"context"
	"testing"

	"github.com/mahanth/simplebank/util"

	"github.com/stretchr/testify/require"
)

func TestCreateTransfer(t *testing.T) {
	// create two accounts to transfer between
	account1, err := testQueries.CreateAccount(context.Background(), CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurrency(),
	})
	require.NoError(t, err)
	require.NotEmpty(t, account1)
	require.NotZero(t, account1.ID)
	account2, err := testQueries.CreateAccount(context.Background(), CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurrency(),
	})
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.NotZero(t, account2.ID)
	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomBalance(),
	}
	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)
	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)
}

func TestGetTransfer(t *testing.T) {
	// create two accounts to transfer between
	account1, err := testQueries.CreateAccount(context.Background(), CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurrency(),
	})
	require.NoError(t, err)
	require.NotEmpty(t, account1)
	require.NotZero(t, account1.ID)
	account2, err := testQueries.CreateAccount(context.Background(), CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurrency(),
	})
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.NotZero(t, account2.ID)

	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomBalance(),
	}

	transfer1, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer1)

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, transfer1.Amount, transfer2.Amount)
}

func TestListTransfers(t *testing.T) {
	// create two accounts to transfer between
	account1, err := testQueries.CreateAccount(context.Background(), CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurrency(),
	})
	require.NoError(t, err)
	require.NotEmpty(t, account1)
	require.NotZero(t, account1.ID)
	account2, err := testQueries.CreateAccount(context.Background(), CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurrency(),
	})
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.NotZero(t, account2.ID)

	for i := 0; i < 10; i++ {
		arg := CreateTransferParams{
			FromAccountID: account1.ID,
			ToAccountID:   account2.ID,
			Amount:        util.RandomBalance(),
		}
		_, err := testQueries.CreateTransfer(context.Background(), arg)
		require.NoError(t, err)
	}

	arg := ListTransfersParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Limit:         5,
		Offset:        5,
	}
	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)
	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)
		require.NotZero(t, transfer.FromAccountID)
		require.NotZero(t, transfer.ToAccountID)
		require.NotZero(t, transfer.Amount)
	}
}
