package api

import (
	//"log"
	"database/sql"
	"net/http"

	//"strconv"

	//"fmt"

	"github.com/gin-gonic/gin"
	db "github.com/meganviga/simple_bank/db/sqlc"
)

type createAccountRequest struct{
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD INR"`
}
func (s *Server)createAccount(ctx *gin.Context){
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil{
		ctx.JSON(http.StatusBadRequest,errResponse(err))
		return
	}
	arg := db.CreateAccountParams{
		Owner: req.Owner,
		Currency: req.Currency,
	}
	account, err := s.store.CreateAccount(ctx,arg)

	if err != nil{
		ctx.JSON(http.StatusNonAuthoritativeInfo,errResponse(err))
		return
	}
	ctx.JSON(http.StatusOK,account)


}

type getAccountParams struct{
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (s *Server)getAccount(ctx *gin.Context){
	var req getAccountParams

	if err := ctx.ShouldBindUri(&req);err != nil{
		ctx.JSON(http.StatusNotFound,errResponse(err))
		return
	}
	
	account, err := s.store.GetAccount(ctx,req.ID)
	if err != nil{
		if err == sql.ErrNoRows{
			ctx.JSON(http.StatusNotFound,errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError,errResponse(err))
		return
	}
	ctx.JSON(http.StatusOK,account)
}
type listAccountsParams struct{
	PageID int64 `form:"page_id" binding:"required,min=1"`
	PageSize   int64 `form:"page_size" binding:"required,min=5,max=10"`
}
func (s * Server)listAccounts(ctx *gin.Context){
	var req listAccountsParams
	if err := ctx.ShouldBindQuery(&req); err != nil{
		ctx.JSON(http.StatusNotFound,errResponse(err))
		return
	}
	arg := db.ListAccountsParams{
		Limit: int32(req.PageSize),
		Offset: (int32(req.PageID) - 1) * int32(req.PageSize),
	}
	accounts, err := s.store.ListAccounts(ctx,arg)
	if err != nil{
		if err == sql.ErrNoRows{
			ctx.JSON(http.StatusNotExtended,errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return

	}

	ctx.JSON(http.StatusOK,accounts)


}