package database

import (
	"SimpleBank/util"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestQueries_CreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestQueries_GetAccount(t *testing.T) {
	// create account
	recordCreate := createRandomAccount(t)
	// get account
	recordGet, err := testQueries.GetAccount(context.Background(), recordCreate.ID)
	require.NoError(t, err)
	require.NotEmpty(t, recordGet)

	require.Equal(t, recordCreate.ID, recordGet.ID)
	require.Equal(t, recordCreate.Owner, recordGet.Owner)
	require.Equal(t, recordCreate.Balance, recordGet.Balance)
	require.Equal(t, recordCreate.Currency, recordGet.Currency)
	require.WithinDurationf(t, recordCreate.CreatedAt, recordGet.CreatedAt, time.Second, "get duration create account")
}

func TestQueries_UpdateAccount(t *testing.T) {
	// create account
	recordCreate := createRandomAccount(t)
	// update account
	arg := UpdateAccountParams{
		ID:      recordCreate.ID,
		Balance: util.RandomMoney(),
	}
	recordGet, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, recordGet)

	require.Equal(t, recordCreate.ID, recordGet.ID)
	require.Equal(t, recordCreate.Owner, recordGet.Owner)
	require.Equal(t, arg.Balance, recordGet.Balance)
	require.Equal(t, recordCreate.Currency, recordGet.Currency)
	require.WithinDurationf(t, recordCreate.CreatedAt, recordGet.CreatedAt, time.Second, "get duration create account")
}

func TestQueries_ListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}
	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
