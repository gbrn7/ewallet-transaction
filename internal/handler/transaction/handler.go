package transaction

import (
	"context"
	"ewallet-transaction/internal/models"

	"github.com/gin-gonic/gin"
)

//go:generate mockgen -source=handler.go -destination=handler_mock_test.go -package=transaction
type Service interface {
	CreateTransaction(ctx context.Context, req *models.Transaction) (models.CreateTransactionResponse, error)
	UpdateStatusTransaction(ctx context.Context, tokenData models.TokenData, req *models.UpdateStatusTransaction) error
	GetTransactionDetail(ctx context.Context, reference string) (models.Transaction, error)
	GetTransaction(ctx context.Context, userID uint64) ([]models.Transaction, error)
	RefundTransaction(ctx context.Context, tokenData models.TokenData, req *models.RefundTransaction) (models.CreateTransactionResponse, error)
}

type Handler struct {
	*gin.Engine
	Service    Service
	External   External
	Middleware Middleware
}

func NewHandler(api *gin.Engine, service Service, ext External, mdw Middleware) *Handler {
	return &Handler{
		api,
		service,
		ext,
		mdw,
	}
}

func (h *Handler) RegisterRoute() {
	transactionV1 := h.Group("/transaction/v1")
	transactionV1.POST("/create", h.Middleware.MiddlewareValidateToken, h.CreateTransaction)
	transactionV1.PUT("/update-status/:reference", h.Middleware.MiddlewareValidateToken, h.UpdateStatusTransaction)
	transactionV1.GET("/", h.Middleware.MiddlewareValidateToken, h.GetTransaction)
	transactionV1.GET("/:reference", h.Middleware.MiddlewareValidateToken, h.GetTransactionDetail)
	transactionV1.POST("/refund", h.Middleware.MiddlewareValidateToken, h.RefundTransaction)
}
