package api

import (
	"ewallet-transaction/constants"
	"ewallet-transaction/helpers"
	"ewallet-transaction/internal/interfaces"
	"ewallet-transaction/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransactionAPI struct {
	TransactionService interfaces.ITransactionService
}

func (api *TransactionAPI) CreateTransaction(c *gin.Context) {
	var (
		log = helpers.Logger
		req models.Transaction
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error("failed to parse request, ", err)
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
		return
	}

	if err := req.Validate(); err != nil {
		log.Error("failed to validate request, ", err)
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
		return
	}

	token, ok := c.Get("token")
	if !ok {
		log.Error("failed to get token data")
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	tokenData, ok := token.(models.TokenData)
	if !ok {
		log.Error("failed to parse token data")
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	if !constants.MapTransactionType[req.TransactionType] {
		log.Error("invalid transaction type")
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
		return
	}

	req.UserID = tokenData.UserID

	resp, err := api.TransactionService.CreateTransaction(c.Request.Context(), &req)
	if err != nil {
		log.Error("failed to create transaction, ", err)
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	helpers.SendResponseHTTP(c, http.StatusOK, constants.SuccessMessage, resp)
}

func (api *TransactionAPI) UpdateStatusTransaction(c *gin.Context) {
	var (
		log = helpers.Logger
		req models.UpdateStatusTransaction
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error("failed to parse request, ", err)
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
		return
	}

	req.Reference = c.Param("reference")
	if err := req.Validate(); err != nil {
		log.Error("failed to validate request, ", err)
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
		return
	}

	token, ok := c.Get("token")
	if !ok {
		log.Error("failed to get token data")
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	tokenData, ok := token.(models.TokenData)
	if !ok {
		log.Error("failed to parse token data")
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	err := api.TransactionService.UpdateStatusTransaction(c.Request.Context(), tokenData, &req)

	if err != nil {
		log.Error("failed to update transaction, ", err)
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	helpers.SendResponseHTTP(c, http.StatusOK, constants.SuccessMessage, nil)

}
