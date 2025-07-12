package db

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	// create two accounts to transfer between
	account1 := createRandomAccount(t)
	require.NotEmpty(t, account1)
	require.NotZero(t, account1.ID)

	// Create a second account to transfer to
	account2 := createRandomAccount(t)

	require.NotEmpty(t, account2)
	require.NotZero(t, account2.ID)

	fmt.Println(">> before transfer: ", account1.Balance, account2.Balance)

	store := NewStore(testDB)
	//run a concurrent transfer transaction

	totalTransfers := 5
	errs := make(chan error)
	results := make(chan TransferTxResult)
	amount := int64(10)

	// Create a wait group to wait for all goroutines to finish
	wg := sync.WaitGroup{}

	for i := 0; i < totalTransfers; i++ {

		go func() {
			wg.Add(1)
			defer wg.Done()
			arg := TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			}

			result, err := store.TransferTx(context.Background(), arg)
			errs <- err
			results <- result

		}()
	}
	// Wait for all goroutines to finish
	wg.Wait()

	existed := make(map[int]bool)

	for i := 0; i < totalTransfers; i++ {
		err := <-errs
		require.NoError(t, err)
		result := <-results
		require.NotEmpty(t, result)

		// Check the accounts involved in the transfer
		require.Equal(t, account1.ID, result.FromAccount.ID)
		require.Equal(t, account2.ID, result.ToAccount.ID)
		require.Equal(t, account1.ID, result.Transfer.FromAccountID)
		require.Equal(t, account2.ID, result.Transfer.ToAccountID)
		require.NotZero(t, result.Transfer.ID)
		require.NotZero(t, result.Transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), result.Transfer.ID)
		require.NoError(t, err)

		//Check the entries created for the transfer
		require.NotEmpty(t, result.FromEntry)
		require.NotEmpty(t, result.ToEntry)
		require.Equal(t, result.FromEntry.AccountID, account1.ID)
		require.Equal(t, result.ToEntry.AccountID, account2.ID)
		require.Equal(t, result.FromEntry.Amount, -result.Transfer.Amount)
		require.Equal(t, result.ToEntry.Amount, result.Transfer.Amount)
		require.NotZero(t, result.FromEntry.ID)
		require.NotZero(t, result.ToEntry.ID)
		require.NotZero(t, result.FromEntry.CreatedAt)
		require.NotZero(t, result.ToEntry.CreatedAt)

		//TODO: check the final balances of the accounts
		//check accounts
		fromAccount := result.FromAccount
		toAccount := result.ToAccount
		require.NotEmpty(t, fromAccount)
		require.NotEmpty(t, toAccount)
		require.Equal(t, account1.ID, fromAccount.ID)
		require.Equal(t, account2.ID, toAccount.ID)

		//check account balances
		fmt.Println(">> tx:", fromAccount.Balance, toAccount.Balance)
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0, "amount should be a multiple of the transfer amount") // 1*amount, 2*amount, etc.

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= totalTransfers, "k should be between 1 and totalTransfers")
		require.False(t, existed[k], "k should be unique")
		existed[k] = true
	}

	// Check the final balances of the accounts
	finalAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, finalAccount1)
	require.Equal(t, account1.Balance-int64(totalTransfers)*amount, finalAccount1.Balance)
	finalAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	require.NotEmpty(t, finalAccount2)
	require.Equal(t, account2.Balance+int64(totalTransfers)*amount, finalAccount2.Balance)

	fmt.Println(">> after transfer: ", finalAccount1.Balance, finalAccount2.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	// create two accounts to transfer between
	account1 := createRandomAccount(t)
	require.NotEmpty(t, account1)
	require.NotZero(t, account1.ID)

	// Create a second account to transfer to
	account2 := createRandomAccount(t)

	require.NotEmpty(t, account2)
	require.NotZero(t, account2.ID)

	fmt.Println(">> before transfer: ", account1.Balance, account2.Balance)

	store := NewStore(testDB)
	//run a concurrent transfer transaction

	totalTransfers := 10
	errs := make(chan error)

	amount := int64(10)

	for i := 0; i < totalTransfers; i++ {

		if i%2 == 0 {
			go func() {

				arg := TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        amount,
				}

				_, err := store.TransferTx(context.Background(), arg)
				errs <- err

			}()
		} else {
			go func() {

				arg := TransferTxParams{
					FromAccountID: account2.ID,
					ToAccountID:   account1.ID,
					Amount:        amount,
				}

				_, err := store.TransferTx(context.Background(), arg)
				errs <- err

			}()
		}

	}

	for i := 0; i < totalTransfers; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	// Check the final balances of the accounts
	finalAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, finalAccount1)
	require.Equal(t, account1.Balance, finalAccount1.Balance)
	finalAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	require.NotEmpty(t, finalAccount2)
	require.Equal(t, account2.Balance, finalAccount2.Balance)

	fmt.Println(">> after transfer: ", finalAccount1.Balance, finalAccount2.Balance)
}
