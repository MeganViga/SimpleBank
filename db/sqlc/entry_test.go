package db

import (
	"context"
	"testing"

	"github.com/meganviga/simple_bank/util"
	"github.com/stretchr/testify/require"
)

func createEntryForAccountID(t *testing.T,id int64)Entry{
	args := CreateEntryParams{
		AccountID: id,
		Amount: util.RandomMoney(),
	}
	entry, err := testQueries.CreateEntry(context.Background(),args)

	require.NoError(t,err)
	require.NotEmpty(t,entry)
	require.Equal(t,args.AccountID,entry.AccountID)
	require.Equal(t,args.Amount,entry.Amount)
	return entry
}
func TestCreateEntry(t *testing.T){
	account := CreateRandomAccount(t)
	createEntryForAccountID(t,account.ID)
}

func TestGetEntry(t *testing.T){
	account := CreateRandomAccount(t)

	entry1 := createEntryForAccountID(t,account.ID)
	entry2, err := testQueries.GetEntry(context.Background(),entry1.ID)

	require.NoError(t,err)
	require.NotEmpty(t,entry2)
	require.Equal(t,account.ID,entry2.AccountID)
	require.Equal(t,entry1.Amount,entry2.Amount)

}

func TestListEntries(t *testing.T){

	for i :=0;i < 10;i++{
		account := CreateRandomAccount(t)

		createEntryForAccountID(t,account.ID)
	}
	args := ListEntriesParams{
		Limit: 5,
		Offset: 5,
	}
	entries, err := testQueries.ListEntries(context.Background(),args)
	require.NoError(t,err)
	require.Len(t,entries,5)
	for _, entry := range entries{
		require.NotEmpty(t,entry)
	}
	require.NotEmpty(t,entries)
}