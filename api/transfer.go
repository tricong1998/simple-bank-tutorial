package api

import (
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/Sotatek-CongNguyen/simple-bank-practice/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createTransferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Currency      string `json:"currency" binding:"required,currency"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req createTransferRequest
	if err := ctx.ShouldBindJSON((&req)); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if validateAccount := server.validateAccountTransfer(ctx, req.FromAccountID, req.Currency); !validateAccount {
		return
	}

	if validateAccount := server.validateAccountTransfer(ctx, req.ToAccountID, req.Currency); !validateAccount {
		return
	}

	arg := db.TransferTxParams{FromAccountID: req.FromAccountID,
		ToAccountID: req.ToAccountID, Amount: req.Amount}
	account, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func (server *Server) validateAccountTransfer(ctx *gin.Context, accountID int64, currency string) bool {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}

	if account.Currency != currency {
		err = fmt.Errorf("invalid currency, account: %d, account's currency: %s, currency input: %s", accountID, account.Currency, currency)
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return false
	}

	return true
}
