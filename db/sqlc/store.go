package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Transfer(ctx context.Context, arg TransferParams) (TransferResult, error)
	Querier
}

// provides all functions to execute db queries and transactions
type SQLStore struct {
	*Queries
	db *sql.DB
}

// creates a new Store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// excecutes a function with a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx er: %v, rollback err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// provides all of parameter of transaction
type TransferParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// provides result of a transaction
type TransferResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// This Tx perform a money transfer transaction from account A to account B like section 3 of the README file

// Step 1: Create a record of the transaction with amount = 10
// Step 2: Create an entry account for A with amount = -10
// Step 3: Create an entry account for B with amount = 10
// Step 4: Subtract 10 from the balance of A
// Step 5: Add 10 to the balance of B
// Step 1
func (store *SQLStore) Transfer(ctx context.Context, arg TransferParams) (TransferResult, error) {
	var result TransferResult

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

		//Step 2
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		//Step 3
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}
		//Step 4: Note that step 4 and step 5 will require locking protocol

		if arg.FromAccountID < arg.ToAccountID { // This for avoiding deadlock
			account1, err := q.GetAccountForUpdate(ctx, arg.FromAccountID)

			if err != nil {
				return err
			}
			result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
				ID:      arg.FromAccountID,
				Balance: account1.Balance - arg.Amount,
			})

			account2, err := q.GetAccountForUpdate(ctx, arg.FromAccountID)

			if err != nil {
				return err
			}
			result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
				ID:      arg.FromAccountID,
				Balance: account2.Balance + arg.Amount,
			})
		} else {
			account2, err := q.GetAccountForUpdate(ctx, arg.FromAccountID)

			if err != nil {
				return err
			}
			result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
				ID:      arg.FromAccountID,
				Balance: account2.Balance + arg.Amount,
			})

			account1, err := q.GetAccountForUpdate(ctx, arg.FromAccountID)

			if err != nil {
				return err
			}
			result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
				ID:      arg.FromAccountID,
				Balance: account1.Balance - arg.Amount,
			})
		}
		return nil
	})
	return result, err
}
