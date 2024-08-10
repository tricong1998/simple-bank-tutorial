package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Printf("Before transfer, account 1: %v, account 2: %v \n", account1.Balance, account2.Balance)
	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx %d", i+1)

		go func() {
			arg := TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			}
			ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := store.TransferTx(ctx, arg)

			errs <- err
			results <- result
		}()
	}

	mapExisted := make(map[int]bool)

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		// check results
		result := <-results
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)
		require.Equal(t, transfer.FromAccountID, account1.ID)
		require.Equal(t, transfer.ToAccountID, account2.ID)
		require.Equal(t, transfer.Amount, amount)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check from entry
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)
		require.Equal(t, fromEntry.AccountID, account1.ID)
		require.Equal(t, fromEntry.Amount, -amount)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		// check to entry
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)
		require.Equal(t, toEntry.AccountID, account2.ID)
		require.Equal(t, toEntry.Amount, amount)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check account balance
		fromAccount := result.FromAccount
		toAccount := result.ToAccount
		diffFrom := account1.Balance - fromAccount.Balance
		diffTo := toAccount.Balance - account2.Balance
		fmt.Printf("After transfer, account 1: %v, account 2: %v\n", fromAccount.Balance, toAccount.Balance)

		require.Equal(t, diffFrom, diffTo)
		require.True(t, diffFrom > 0)
		require.True(t, diffFrom%amount == 0)
		k := int(diffFrom / amount)
		require.NotContains(t, mapExisted, k)
		mapExisted[k] = true
	}

	updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NotEmpty(t, updatedAccount1)
	require.NoError(t, err)
	require.Equal(t, updatedAccount1.Balance, account1.Balance-(amount*int64(n)))

	updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)
	require.NotEmpty(t, updatedAccount2)
	require.NoError(t, err)
	require.Equal(t, updatedAccount2.Balance, account2.Balance+(amount*int64(n)))

	fmt.Printf("After transfer, account 1: %v, account 2: %v", updatedAccount1.Balance, updatedAccount2.Balance)
}

func TestTransferTxDeadLock(t *testing.T) {
	store := NewStore(testDB)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Printf("Before transfer, account 1: %v, account 2: %v \n", account1.Balance, account2.Balance)
	n := 10
	amount := int64(10)

	errs := make(chan error)

	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx %d", i+1)
		fromAccountId := account1.ID
		toAccountId := account2.ID
		if i%2 == 1 {
			fromAccountId = account2.ID
			toAccountId = account1.ID
		}
		go func() {
			arg := TransferTxParams{
				FromAccountID: fromAccountId,
				ToAccountID:   toAccountId,
				Amount:        amount,
			}
			ctx := context.WithValue(context.Background(), txKey, txName)
			_, err := store.TransferTx(ctx, arg)

			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

	}

	updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NotEmpty(t, updatedAccount1)
	require.NoError(t, err)
	require.Equal(t, updatedAccount1.Balance, account1.Balance)

	updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)
	require.NotEmpty(t, updatedAccount2)
	require.NoError(t, err)
	require.Equal(t, updatedAccount2.Balance, account2.Balance)

	fmt.Printf("After transfer, account 1: %v, account 2: %v", updatedAccount1.Balance, updatedAccount2.Balance)
}
