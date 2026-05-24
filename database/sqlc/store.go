package database

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		Queries: New(db),
		db:      db,
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	err := store.execTx(ctx, func(queries *Queries) error {
		var err error

		result.Transfer, err = queries.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}
		result.FromEntry, err = queries.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}
		result.ToEntry, err = queries.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// update balance
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, err = queries.AddAccountBalance(ctx, AddAccountBalanceParams{
				Amount: -arg.Amount,
				ID:     arg.FromAccountID,
			})
			if err != nil {
				return err
			}
			result.ToAccount, err = queries.AddAccountBalance(ctx, AddAccountBalanceParams{
				Amount: arg.Amount,
				ID:     arg.ToAccountID,
			})
			if err != nil {
				return err
			}
		} else {
			result.ToAccount, err = queries.AddAccountBalance(ctx, AddAccountBalanceParams{
				Amount: arg.Amount,
				ID:     arg.ToAccountID,
			})
			if err != nil {
				return err
			}
			result.FromAccount, err = queries.AddAccountBalance(ctx, AddAccountBalanceParams{
				Amount: -arg.Amount,
				ID:     arg.FromAccountID,
			})
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return result, err
	}
	return result, nil
}
