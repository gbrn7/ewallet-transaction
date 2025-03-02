package transaction

import (
	"context"
	"ewallet-transaction/external"
	"ewallet-transaction/internal/models"
)

//go:generate mockgen -source=external.go -destination=external_mock_test.go -package=transaction
type External interface {
	ValidateToken(ctx context.Context, token string) (models.TokenData, error)
	CreditBalance(ctx context.Context, token string, req external.UpdateBalance) (*external.UpdateBalanceResponse, error)
	DebitBalance(ctx context.Context, token string, req external.UpdateBalance) (*external.UpdateBalanceResponse, error)
	SendNotification(ctx context.Context, recipient string, templateName string, placeHolder map[string]string) error
}
