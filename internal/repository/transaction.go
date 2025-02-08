package repository

import (
	"context"
	"ewallet-transaction/constants"
	"ewallet-transaction/internal/models"

	"gorm.io/gorm"
)

type TransactionRepo struct {
	DB *gorm.DB
}

func (r *TransactionRepo) CreateTransaction(ctx context.Context, trx *models.Transaction) error {
	return r.DB.Create(trx).Error
}

func (r *TransactionRepo) GetTransactionByReference(ctx context.Context, reference string, includeRefund bool) (models.Transaction, error) {
	var (
		resp models.Transaction
	)
	sql := r.DB.Where("reference = ?", reference)
	if !includeRefund {
		sql = sql.Where("transaction_type != ?", constants.TransactionTypeRefund)
	}
	err := sql.Last(&resp).Error

	return resp, err
}

func (r *TransactionRepo) UpdateStatusTransaction(ctx context.Context, reference string, status string, additionalInfo string) error {
	return r.DB.Exec("UPDATE transactions SET transaction_status = ?, additional_info = ? WHERE reference = ?", status, additionalInfo, reference).Error
}

func (r *TransactionRepo) GetTransaction(ctx context.Context, userID uint64) ([]models.Transaction, error) {
	var (
		resp []models.Transaction
	)
	err := r.DB.Order("id DESC").Where("user_id = ?", userID).Find(&resp).Error

	return resp, err
}
