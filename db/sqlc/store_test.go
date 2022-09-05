package db

import (
	"context"
	"testing"
	"fmt"
	"github.com/meganviga/simple_bank/util"
	"github.com/stretchr/testify/require"
)


func TestTransferTx(t *testing.T){
	store := NewStore(testDB)
	t.Log("hello")
	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)
	fmt.Println("==>before",account1.Balance,account2.Balance)
	
	n := 5
	errs := make(chan error)
	results := make(chan TransferTxResult)
	amount := util.RandomMoney()
	fmt.Println("Amount: ",amount)
	for i := 0;i < n;i++{
		txName := fmt.Sprintf("tx %d",i + 1)
		ctx := context.WithValue(context.Background(),txKey,txName)
		go func() {
			result, err := store.TransferTx(ctx,TransferTxParams{
				FROMACCOUNTID: account1.ID,
				TOACCOUNTID: account2.ID,
				AMOUNT: amount,
			})
			errs <- err
			results <- result
		}()
	}
	 existed := make(map[int]bool)
	for i := 0;i < n;i++{
		err := <- errs
		require.NoError(t,err)
		result := <-results
		transfer := result.Transfer
		require.NotEmpty(t,transfer)
		require.Equal(t,transfer.FromAccountID,account1.ID)
		require.Equal(t,transfer.ToAccountID,account2.ID)
		require.Equal(t,transfer.Amount,amount)
		require.NotZero(t,transfer.ID)
		require.NotZero(t,transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(),transfer.ID)
		require.NoError(t,err)

		fromentry := result.FromEntry
		require.NotEmpty(t,fromentry)
		require.Equal(t,account1.ID,fromentry.AccountID)
		require.Equal(t,-amount,fromentry.Amount)
		require.NotZero(t,fromentry.ID)
		require.NotZero(t,fromentry.CreatedAt)

		_, err = store.GetEntry(context.Background(),fromentry.ID)

		require.NoError(t,err)

		toentry := result.ToEntry
		require.NotEmpty(t,toentry)
		require.Equal(t,toentry.AccountID,account2.ID)
		require.Equal(t,toentry.Amount,amount)
		require.NotZero(t,toentry.ID)
		require.NotZero(t,toentry.CreatedAt)

		_, err = store.GetEntry(context.Background(),toentry.ID)

		require.NoError(t,err)

		//check account
		fromAccount := result.FromAccount
		require.NotEmpty(t,fromAccount)
		require.Equal(t,fromAccount.ID,account1.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t,toAccount)
		require.Equal(t,toAccount.ID,account2.ID)
		fmt.Println("==>tx",fromAccount.Balance,toAccount.Balance)
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance

		require.Equal(t,diff1,diff2)
		require.True(t,diff1 > 0)
		require.True(t,diff1 % amount == 0)

		k := int(diff1 / amount)
		require.True(t,k>=1,k<=n)
		require.NotContains(t,existed,k)
		existed[k]= true


	}
		updateAccount1, err := testQueries.GetAccount(context.Background(),account1.ID)
		require.NoError(t,err)
		require.NotEmpty(t,updateAccount1)
		updateAccount2, err := testQueries.GetAccount(context.Background(),account2.ID)
		require.NoError(t,err)
		require.NotEmpty(t,updateAccount2)
		fmt.Println("==>after",updateAccount1.Balance,updateAccount2.Balance)
		require.Equal(t,account1.Balance - int64(n)*amount,updateAccount1.Balance)
		require.Equal(t,account2.Balance + int64(n)*amount,updateAccount2.Balance)
}