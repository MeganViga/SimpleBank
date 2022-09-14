package db

import (
	"context"
	"database/sql"
	"fmt"
)

/* type Store interface{
	Querier
	TransferTx(ctx context.Context,arg TransferTxParams)(TransferTxResult, error)

} */
type Store struct{
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB)*Store{
	return &Store{
		Queries: New(db),
		db:db,
	}
}

func (store *Store)execTx(ctx context.Context,fn func(*Queries)error)error{
	tx, err := store.db.BeginTx(context.Background(),nil)
	if err != nil{
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil{
		if rbErr := tx.Rollback();rbErr != nil{
			return fmt.Errorf("tx err: %v rbErr: %v",err,rbErr)
		}
	}

	return tx.Commit()
}

type TransferTxParams struct{
	FROMACCOUNTID int64 `json:"from_account_id"`
	TOACCOUNTID   int64 `json:"to_account_id"`
	AMOUNT        int64 `json:"amount"`
}

type TransferTxResult struct{
	Transfer Transfer `json:"transfer"`
	FromAccount  Account `json:"from_account"`
	ToAccount    Account `json:"to_account"`
	FromEntry    Entry   `json:"from_entry"`
	ToEntry      Entry   `json:"to_entry"`
}
var txKey = struct{}{}
func (store *Store)TransferTx(ctx context.Context,arg TransferTxParams)(TransferTxResult, error){
	var result TransferTxResult
	ctxValue := ctx.Value(txKey)
	
	err := store.execTx(context.Background(),func( q *Queries)error{
		var err error
		fmt.Println(ctxValue,"Create transfer ")
		result.Transfer, err = q.CreateTransfer(context.Background(),CreateTransferParams{
			FromAccountID: arg.FROMACCOUNTID,
			ToAccountID: arg.TOACCOUNTID,
			Amount : arg.AMOUNT,
		})
		if err != nil{
			return err
		}
		fmt.Println(ctxValue,"Create entry 1 ")
		result.FromEntry, err = q.CreateEntry(context.Background(),CreateEntryParams{
			AccountID: arg.FROMACCOUNTID,
			Amount: -arg.AMOUNT,
		})
		if err != nil{
			return err
		}
		fmt.Println(ctxValue,"Create entry 2 ")
		result.ToEntry , err = q.CreateEntry(context.Background(),CreateEntryParams{
			AccountID: arg.TOACCOUNTID,
			Amount: +arg.AMOUNT,
		})
		if err != nil{
			return err
		}

		//Todo update amount on accounts
		fmt.Println(ctxValue,"Get Account 1 ")
		_, err = q.GetAccountForUpdate(context.Background(),arg.FROMACCOUNTID)
		if err != nil{
			return err
		}
		fa := AddAccountBalanceParams{
			ID: arg.FROMACCOUNTID,
			Amount: -arg.AMOUNT,
		}
		fmt.Println(ctxValue,"Update Account 1 ")
		result.FromAccount, err = q.AddAccountBalance(context.Background(),fa)
		if err != nil{
			return err
		}
		fmt.Println(ctxValue,"Get Account 2 ")
		_, err = q.GetAccountForUpdate(context.Background(),arg.TOACCOUNTID)
		if err != nil{
			return err
		}
		ta := AddAccountBalanceParams{
			ID: arg.TOACCOUNTID,
			Amount:arg.AMOUNT,
		}
		fmt.Println(ctxValue,"Update Account 2 ")
		result.ToAccount, err = q.AddAccountBalance(context.Background(),ta)
		if err != nil{
			return err
		}

		return nil
	})
	return result, err
}