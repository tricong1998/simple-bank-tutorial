package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	db "github.com/Sotatek-CongNguyen/simple-bank-practice/db/sqlc"
	"github.com/Sotatek-CongNguyen/simple-bank-practice/token"
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

	fromAccount, validateAccount := server.validateAccountTransfer(ctx, req.FromAccountID, req.Currency)
	if !validateAccount {
		return
	}

	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if fromAccount.Owner != payload.Username {
		err := errors.New("account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	_, validateAccount = server.validateAccountTransfer(ctx, req.ToAccountID, req.Currency)
	if !validateAccount {
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

func (server *Server) validateAccountTransfer(ctx *gin.Context, accountID int64, currency string) (db.Account, bool) {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return account, false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}

	if account.Currency != currency {
		err = fmt.Errorf("invalid currency, account: %d, account's currency: %s, currency input: %s", accountID, account.Currency, currency)
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return account, false
	}

	return account, true
}
