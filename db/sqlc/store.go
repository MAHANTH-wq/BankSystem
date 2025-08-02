package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)
	VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error)
}

// Inheritance from Queries struct
type SQLStore struct {
	q  *Queries
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) Store {
	q := New(db)
	return &SQLStore{
		q:  q,
		db: db,
	}
}

// execTx executes a function within a transaction context.
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
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

func (store *SQLStore) GetAccount(ctx context.Context, id int64) (Account, error) {
	account, err := store.q.GetAccount(ctx, id)
	if err != nil {
		return Account{}, err
	}
	return account, nil
}

func (store *SQLStore) CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error) {
	account, err := store.q.CreateAccount(ctx, arg)
	if err != nil {
		return Account{}, err
	}
	return account, nil
}

func (store *SQLStore) ListAccounts(ctx context.Context, arg ListAccountsParams) ([]Account, error) {
	accounts, err := store.q.ListAccounts(ctx, arg)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}
func (store *SQLStore) AddAccountBalance(ctx context.Context, arg AddAccountBalanceParams) (Account, error) {
	account, err := store.q.AddAccountBalance(ctx, arg)
	if err != nil {
		return Account{}, err
	}
	return account, nil
}

func (store *SQLStore) DeleteAccount(ctx context.Context, id int64) error {
	err := store.q.DeleteAccount(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
func (store *SQLStore) UpdateAccount(ctx context.Context, arg UpdateAccountParams) (Account, error) {
	account, err := store.q.UpdateAccount(ctx, arg)
	if err != nil {
		return Account{}, err
	}
	return account, nil
}

func (store *SQLStore) GetEntry(ctx context.Context, id int64) (Entry, error) {
	entry, err := store.q.GetEntry(ctx, id)
	if err != nil {
		return Entry{}, err
	}
	return entry, nil
}
func (store *SQLStore) CreateEntry(ctx context.Context, arg CreateEntryParams) (Entry, error) {
	entry, err := store.q.CreateEntry(ctx, arg)
	if err != nil {
		return Entry{}, err
	}
	return entry, nil
}
func (store *SQLStore) ListEntries(ctx context.Context, arg ListEntriesParams) ([]Entry, error) {
	entries, err := store.q.ListEntries(ctx, arg)
	if err != nil {
		return nil, err
	}
	return entries, nil
}

func (store *SQLStore) GetTransfer(ctx context.Context, id int64) (Transfer, error) {
	transfer, err := store.q.GetTransfer(ctx, id)
	if err != nil {
		return Transfer{}, err
	}
	return transfer, nil
}

func (store *SQLStore) ListTransfers(ctx context.Context, arg ListTransfersParams) ([]Transfer, error) {
	transfers, err := store.q.ListTransfers(ctx, arg)
	if err != nil {
		return nil, err
	}
	return transfers, nil
}

func (store *SQLStore) CreateTransfer(ctx context.Context, arg CreateTransferParams) (Transfer, error) {
	transfer, err := store.q.CreateTransfer(ctx, arg)
	if err != nil {
		return Transfer{}, err
	}
	return transfer, nil
}

func (store *SQLStore) GetAccountForUpdate(ctx context.Context, id int64) (Account, error) {
	account, err := store.q.GetAccountForUpdate(ctx, id)
	if err != nil {
		return Account{}, err
	}
	return account, nil
}

func (store *SQLStore) GetUser(ctx context.Context, username string) (User, error) {
	user, err := store.q.GetUser(ctx, username)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (store *SQLStore) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	user, err := store.q.CreateUser(ctx, arg)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (store *SQLStore) CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error) {
	session, err := store.q.CreateSession(ctx, arg)
	if err != nil {
		return Session{}, err
	}
	return session, err
}

func (store *SQLStore) GetSession(ctx context.Context, id pgtype.UUID) (Session, error) {
	session, err := store.q.GetSession(ctx, id)
	if err != nil {
		return Session{}, err
	}
	return session, err
}

func (store *SQLStore) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {

	user, err := store.q.UpdateUser(ctx, arg)
	if err != nil {
		return User{}, err
	}

	return user, err
}

func (store SQLStore) CreateVerifyEmail(ctx context.Context, arg CreateVerifyEmailParams) (VerifyEmail, error) {
	verifyEmail, err := store.q.CreateVerifyEmail(ctx, arg)
	if err != nil {
		return VerifyEmail{}, err
	}
	return verifyEmail, nil
}

func (store SQLStore) UpdateVerifyEmail(ctx context.Context, arg UpdateVerifyEmailParams) (VerifyEmail, error) {
	verifyEmail, err := store.q.UpdateVerifyEmail(ctx, arg)
	if err != nil {
		return VerifyEmail{}, err
	}
	return verifyEmail, nil
}
