package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Inheritance from Queries struct
type Store struct {
	q  *Queries
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	q := New(db)
	return &Store{
		q:  q,
		db: db,
	}
}

// execTx executes a function within a transaction context.
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.Begin(ctx)
	if err != nil {
		return err
	}
	//It is the same as the New(db) function and it works because the transaction implements the DBTX interface
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return rbErr
		}
		return err
	}
	if err = tx.Commit(ctx); err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return rbErr
		}
		return err
	}
	return nil
}

// TransferTxParams contains the parameters for the TransferTx function.
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult contains the result of the TransferTx function.
// It includes the transfer details, the accounts involved, and the entries created for the transaction.
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TranferTx performs a money transfer from one account to another within a transaction context.
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromAccount, err = q.GetAccount(ctx, arg.FromAccountID)

		if err != nil {
			return err
		}

		result.ToAccount, err = q.GetAccount(ctx, arg.ToAccountID)

		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})

		//update accounts' balance and handle deadlock avoidance very important

		result.FromAccount, result.ToAccount, err = q.transferMoney(ctx,
			arg.FromAccountID, arg.ToAccountID, arg.Amount)

		return err
	})

	if err != nil {
		return TransferTxResult{}, err
	}
	return result, nil
}

func (q *Queries) transferMoney(ctx context.Context,
	fromAccountID, toAccountID int64, amount int64) (account1 Account, account2 Account, err error) {
	if fromAccountID > toAccountID {

		account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     fromAccountID,
			Amount: -amount,
		})

		if err != nil {
			return
		}

		account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     toAccountID,
			Amount: amount,
		})

		if err != nil {
			return
		}

	} else {

		account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     toAccountID,
			Amount: amount,
		})

		if err != nil {
			return
		}

		account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     fromAccountID,
			Amount: -amount,
		})

		if err != nil {
			return
		}

	}
	return
}
