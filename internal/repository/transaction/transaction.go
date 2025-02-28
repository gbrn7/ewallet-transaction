package transaction

import (
	"context"
	"ewallet-transaction/constants"
	"ewallet-transaction/internal/models"
)

func (r *repository) CreateTransaction(ctx context.Context, trx *models.Transaction) error {
	return r.DB.Create(trx).Error
}

func (r *repository) GetTransactionByReference(ctx context.Context, reference string, includeRefund bool) (models.Transaction, error) {
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

func (r *repository) UpdateStatusTransaction(ctx context.Context, reference string, status string, additionalInfo string) error {
	return r.DB.Exec("UPDATE transactions SET transaction_status = ?, additional_info = ? WHERE reference = ?", status, additionalInfo, reference).Error
}

func (r *repository) GetTransaction(ctx context.Context, userID uint64) ([]models.Transaction, error) {
	var (
		resp []models.Transaction
	)
	err := r.DB.Order("id DESC").Where("user_id = ?", userID).Find(&resp).Error

	return resp, err
}
