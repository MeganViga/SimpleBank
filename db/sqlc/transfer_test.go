package db

import (
	"context"
	"testing"

	"github.com/meganviga/simple_bank/util"
	"github.com/stretchr/testify/require"
)

func createTransferEntry(t *testing.T)Transfer{
	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)
	args := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID: account2.ID,
		Amount: util.RandomMoney(),
	}
	transfer, err := testQueries.CreateTransfer(context.Background(),args)
	require.NoError(t, err)
	require.NotEmpty(t,transfer)
	require.Equal(t,account1.ID,transfer.FromAccountID)
	require.Equal(t,account2.ID,transfer.ToAccountID)
	return transfer
}
func TestCreateTransfer(t *testing.T){
	
	createTransferEntry(t)
}

func TestGetTransfer(t *testing.T){
	transfer1 := createTransferEntry(t)

	transfer2, err := testQueries.GetTransfer(context.Background(),transfer1.ID)
	require.NoError(t,err)
	require.NotEmpty(t,transfer2)
	require.Equal(t,transfer1.ID,transfer2.ID)
	require.Equal(t,transfer1.FromAccountID,transfer2.FromAccountID)
	require.Equal(t,transfer1.ToAccountID,transfer2.ToAccountID)
	require.Equal(t,transfer1.Amount,transfer2.Amount)
	require.Equal(t,transfer1.CreatedAt,transfer2.CreatedAt)
}

func TestListTransfers(t *testing.T){
	for i:=0;i< 10;i++{
		 createTransferEntry(t)
	}
	args := ListTransfersParams{
		Limit: 5,
		Offset: 5,
	}
	transfers, err := testQueries.ListTransfers(context.Background(),args)
	require.NoError(t,err)
	require.Len(t,transfers,5)
	for _, transfer := range transfers{
		require.NotEmpty(t,transfer)
	}
	require.NotEmpty(t,transfers)
}