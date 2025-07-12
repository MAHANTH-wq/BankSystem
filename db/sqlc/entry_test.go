package db

import (
	"context"
	"testing"

	"github.com/mahanth/simplebank/util"

	"github.com/stretchr/testify/require"
)

func TestCreateEntry(t *testing.T) {
	// create an account to associate with the entry
	account := createRandomAccount(t)

	require.NotEmpty(t, account)
	require.NotZero(t, account.ID)

	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    -util.RandomBalance(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)
}

func TestGetEntry(t *testing.T) {
	// create an account to associate with the entry
	account := createRandomAccount(t)
	require.NotEmpty(t, account)
	require.NotZero(t, account.ID)

	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomBalance(),
	}

	entry1, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry1)

	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.Equal(t, entry1.CreatedAt, entry2.CreatedAt)
}

func TestListEntries(t *testing.T) {
	// create an account to associate with the entries
	account := createRandomAccount(t)
	require.NotEmpty(t, account)
	require.NotZero(t, account.ID)

	for i := 0; i < 10; i++ {
		arg := CreateEntryParams{
			AccountID: account.ID,
			Amount:    util.RandomBalance(),
		}

		entry, err := testQueries.CreateEntry(context.Background(), arg)
		require.NoError(t, err)
		require.NotEmpty(t, entry)
	}

	arg := ListEntriesParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    2,
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.Equal(t, account.ID, entry.AccountID)
	}
}
