package transaction

import (
	"context"
	"ewallet-transaction/external"
	"ewallet-transaction/internal/models"
)

//go:generate mockgen -source=service.go -destination=service_mock_test.go -package=transaction
type repository interface {
	CreateTransaction(ctx context.Context, trx *models.Transaction) error
	GetTransactionByReference(context.Context, string, bool) (models.Transaction, error)
	UpdateStatusTransaction(ctx context.Context, reference string, status string, additionalInfo string) error
	GetTransaction(ctx context.Context, userID uint64) ([]models.Transaction, error)
}

type IExternal interface {
	ValidateToken(ctx context.Context, token string) (models.TokenData, error)
	CreditBalance(ctx context.Context, token string, req external.UpdateBalance) (*external.UpdateBalanceResponse, error)
	DebitBalance(ctx context.Context, token string, req external.UpdateBalance) (*external.UpdateBalanceResponse, error)
	SendNotification(ctx context.Context, recipient string, templateName string, placeHolder map[string]string) error
}

type service struct {
	repository repository
	external   IExternal
}

func NewService(repository repository, external IExternal) *service {
	return &service{
		repository: repository,
		external:   external,
	}
}
