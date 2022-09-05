package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/meganviga/simple_bank/util"
	"github.com/stretchr/testify/require"
	//"fmt"
)
func CreateRandomAccount(t *testing.T)Account{
	arg := CreateAccountParams{
		Owner: util.RandomOwner(),
		Balance: util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
	/* arg2 := CreateAccountParams{
		Owner: "viganesh2",
		Balance: 20000,
		Currency: "INR",
	} */
	account, err := testQueries.CreateAccount(context.Background(),arg)
	//fmt.Println(err)
	require.NoError(t,err)
	require.NotEmpty(t,account)
	require.Equal(t,arg.Owner,account.Owner)
	require.Equal(t,arg.Balance,account.Balance)
	require.Equal(t,arg.Currency,account.Currency)
	require.NotZero(t,account.ID)
	require.NotZero(t,account.CreatedAt)
	return account
}
func TestCreateAccount(t *testing.T){
	CreateRandomAccount(t)
	//_, _ = testQueries.CreateAccount(context.Background(),arg2)
}

func TestGetAccount(t *testing.T){
	account1 := CreateRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(),account1.ID)
	require.NoError(t,err)
	require.NotEmpty(t,account2)
	require.Equal(t,account1.ID,account2.ID)
	require.Equal(t,account1.Balance,account2.Balance)
	require.Equal(t,account1.Owner,account2.Owner)
	require.Equal(t,account1.Currency,account2.Currency)
	require.Equal(t,account1.CreatedAt,account2.CreatedAt)

}

func TestListAccounts(t *testing.T){

	for i:=0 ; i < 10;i++{
		CreateRandomAccount(t)
	}

	listparams := ListAccountsParams{
		Limit: 5,
		Offset:5,
	}
	accounts , err := testQueries.ListAccounts(context.Background(),listparams)
	require.NoError(t,err)
	require.Len(t,accounts,5)
	for _, account := range accounts{
		require.NotEmpty(t,account)
	}
	
	//require.Equal(t,3,len(accounts))

}

func TestUpdateAccount(t *testing.T){
	account1 := CreateRandomAccount(t)
	updateAccountParams := UpdateAccountParams{
		ID:account1.ID,
		Balance: util.RandomMoney(),
	}
	account2,err := testQueries.UpdateAccount(context.Background(),updateAccountParams)
	require.NoError(t,err)
	require.NotEqual(t,account1.Balance,account2.Balance)
	require.Equal(t,account1.ID,account2.ID)
	require.Equal(t,account1.Owner,account2.Owner)
	require.Equal(t,account1.Currency,account2.Currency)
	require.Equal(t,account1.CreatedAt,account2.CreatedAt)

}

func TestDeleteAccount(t *testing.T){
	account1 := CreateRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(),account1.ID)
	require.NoError(t,err)
	account2, err := testQueries.GetAccount(context.Background(),account1.ID)
	require.Error(t,err)
	require.EqualError(t,err,sql.ErrNoRows.Error())
	require.Empty(t,account2)
}