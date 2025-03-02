package transaction

import (
	"ewallet-transaction/constants"
	"ewallet-transaction/helpers"
	"ewallet-transaction/internal/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateTransaction(c *gin.Context) {
	var (
		req models.Transaction
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("failed to parse request, ", err)
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
		return
	}

	if err := req.Validate(); err != nil {
		fmt.Println("failed to validate request, ", err)
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
		return
	}

	token, ok := c.Get("token")
	if !ok {
		fmt.Println("failed to get token data")
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	tokenData, ok := token.(models.TokenData)
	if !ok {
		fmt.Println("failed to parse token data")
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	if !constants.MapTransactionType[req.TransactionType] {
		fmt.Println("invalid transaction type")
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
		return
	}

	req.UserID = tokenData.UserID

	resp, err := h.Service.CreateTransaction(c.Request.Context(), &req)
	if err != nil {
		fmt.Println("failed to create transaction, ", err)
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	helpers.SendResponseHTTP(c, http.StatusOK, constants.SuccessMessage, resp)
}

func (h *Handler) UpdateStatusTransaction(c *gin.Context) {
	var (
		req models.UpdateStatusTransaction
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("failed to parse request, ", err)
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
		return
	}

	req.Reference = c.Param("reference")
	if err := req.Validate(); err != nil {
		fmt.Println("failed to validate request, ", err)
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
		return
	}

	token, ok := c.Get("token")
	if !ok {
		fmt.Println("failed to get token data")
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	tokenData, ok := token.(models.TokenData)
	if !ok {
		fmt.Println("failed to parse token data")
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	err := h.Service.UpdateStatusTransaction(c.Request.Context(), tokenData, &req)

	if err != nil {
		fmt.Println("failed to update transaction, ", err)
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	helpers.SendResponseHTTP(c, http.StatusOK, constants.SuccessMessage, nil)
}

func (h *Handler) GetTransaction(c *gin.Context) {

	token, ok := c.Get("token")
	if !ok {
		fmt.Println("failed to get token data")
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	tokenData, ok := token.(models.TokenData)
	if !ok {
		fmt.Println("failed to parse token data")
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	resp, err := h.Service.GetTransaction(c.Request.Context(), uint64(tokenData.UserID))

	if err != nil {
		fmt.Println("failed to update transaction, ", err)
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	helpers.SendResponseHTTP(c, http.StatusOK, constants.SuccessMessage, resp)
}

func (h *Handler) GetTransactionDetail(c *gin.Context) {

	reference := c.Param("reference")
	if reference == "" {
		fmt.Println("failed to get reference")
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
	}

	token, ok := c.Get("token")
	if !ok {
		fmt.Println("failed to get token data")
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	_, ok = token.(models.TokenData)
	if !ok {
		fmt.Println("failed to parse token data")
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	resp, err := h.Service.GetTransactionDetail(c.Request.Context(), reference)

	if err != nil {
		fmt.Println("failed to update transaction, ", err)
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	helpers.SendResponseHTTP(c, http.StatusOK, constants.SuccessMessage, resp)
}

func (h *Handler) RefundTransaction(c *gin.Context) {
	var (
		req models.RefundTransaction
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("failed to parse request, ", err)
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
		return
	}

	if err := req.Validate(); err != nil {
		fmt.Println("failed to validate request, ", err)
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
		return
	}

	token, ok := c.Get("token")
	if !ok {
		fmt.Println("failed to get token data")
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	tokenData, ok := token.(models.TokenData)
	if !ok {
		fmt.Println("failed to parse token data")
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	resp, err := h.Service.RefundTransaction(c.Request.Context(), tokenData, &req)
	if err != nil {
		fmt.Println("failed to refund transaction, ", err)
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	helpers.SendResponseHTTP(c, http.StatusOK, constants.SuccessMessage, resp)
}
