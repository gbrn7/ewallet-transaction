package interfaces

import (
	"context"
	"ewallet-transaction/internal/models"

	"github.com/gin-gonic/gin"
)

type ITransactionRepo interface {
	CreateTransaction(ctx context.Context, trx *models.Transaction) error
	GetTransactionByReference(context.Context, string, bool) (models.Transaction, error)
	UpdateStatusTransaction(ctx context.Context, reference string, status string, additionalInfo string) error
	GetTransaction(ctx context.Context, userID uint64) ([]models.Transaction, error)
}

type ITransactionService interface {
	CreateTransaction(ctx context.Context, req *models.Transaction) (models.CreateTransactionResponse, error)
	UpdateStatusTransaction(ctx context.Context, tokenData models.TokenData, req *models.UpdateStatusTransaction) error
	GetTransactionDetail(ctx context.Context, reference string) (models.Transaction, error)
	GetTransaction(ctx context.Context, userID uint64) ([]models.Transaction, error)
}

type ITransactionAPI interface {
	CreateTransaction(c *gin.Context)
	UpdateStatusTransaction(c *gin.Context)
	GetTransaction(c *gin.Context)
	GetTransactionDetail(c *gin.Context)
}
