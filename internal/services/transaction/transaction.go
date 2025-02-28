package transaction

import (
	"context"
	"encoding/json"
	"ewallet-transaction/constants"
	"ewallet-transaction/external"
	"ewallet-transaction/helpers"
	"ewallet-transaction/internal/models"
	"fmt"
	"time"

	"github.com/pkg/errors"
)

func (s *service) CreateTransaction(ctx context.Context, req *models.Transaction) (models.CreateTransactionResponse, error) {
	var (
		resp models.CreateTransactionResponse
	)

	req.TransactionStatus = constants.TransactionStatusPending
	req.Reference = helpers.GenerateReference()

	jsonAdditionalStatus := map[string]interface{}{}
	if req.AdditionalInfo != "" {
		err := json.Unmarshal([]byte(req.AdditionalInfo), &jsonAdditionalStatus)
		if err != nil {
			return resp, errors.Wrap(err, "failed to unmarshal current additional info")
		}
	}

	err := s.repository.CreateTransaction(ctx, req)
	if err != nil {
		return resp, errors.Wrap(err, "failed to create transaction")
	}

	resp.Reference = req.Reference
	resp.TransactionStatus = req.TransactionStatus
	return resp, nil
}

func (s *service) UpdateStatusTransaction(ctx context.Context, tokenData models.TokenData, req *models.UpdateStatusTransaction) error {
	// get transaction by reference
	trx, err := s.repository.GetTransactionByReference(ctx, req.Reference, false)
	if err != nil {
		return errors.Wrap(err, "failed to get transaction")
	}

	// validate transaction status flow
	statusValid := false
	mapStatusFlow := constants.MapTransactionStatusFlow[trx.TransactionStatus]
	for i := range mapStatusFlow {
		if mapStatusFlow[i] == req.TransactionStatus {
			statusValid = true
		}
	}
	if !statusValid {
		return fmt.Errorf("transaction status flow invalid. request status = %s", req.TransactionStatus)
	}

	//request update balance to ewallet-wallet
	reqUpdateBalance := external.UpdateBalance{
		Amount:    trx.Amount,
		Reference: req.Reference,
	}

	if req.TransactionStatus == constants.TransactionStatusReversed {
		reqUpdateBalance.Reference = "REVERSED-" + req.Reference

		now := time.Now()
		expiredReversalTime := trx.CreatedAt.Add(constants.MaximumReversalDuration)
		if now.After(expiredReversalTime) {
			return errors.New("reversal duration is already expired")
		}
	}

	var errUpdateBalance error

	switch trx.TransactionType {
	case constants.TransactionTypeTopup:
		if req.TransactionStatus == constants.TransactionStatusSuccess {
			_, errUpdateBalance = s.external.CreditBalance(ctx, tokenData.Token, reqUpdateBalance)
		} else if req.TransactionStatus == constants.TransactionStatusReversed {
			_, errUpdateBalance = s.external.DebitBalance(ctx, tokenData.Token, reqUpdateBalance)
		}
	case constants.TransactionTypePurchase:
		if req.TransactionStatus == constants.TransactionStatusSuccess {
			_, errUpdateBalance = s.external.DebitBalance(ctx, tokenData.Token, reqUpdateBalance)
		} else if req.TransactionStatus == constants.TransactionStatusReversed {
			_, errUpdateBalance = s.external.CreditBalance(ctx, tokenData.Token, reqUpdateBalance)
		}
	}

	if errUpdateBalance != nil {
		return errors.Wrap(errUpdateBalance, "failed to update balance")
	}

	// Update additional info
	var (
		newAdditionalInfo     = map[string]interface{}{}
		currentAdditionalInfo = map[string]interface{}{}
	)

	if trx.AdditionalInfo != "" {
		err = json.Unmarshal([]byte(trx.AdditionalInfo), &currentAdditionalInfo)
		if err != nil {
			return errors.Wrap(err, "failed to unmarshal current additional info")
		}
	}

	if req.AdditionalInfo != "" {
		err = json.Unmarshal([]byte(req.AdditionalInfo), &newAdditionalInfo)
		if err != nil {
			return errors.Wrap(err, "failed to unmarshal new additional info")
		}
	}

	for key, val := range newAdditionalInfo {
		currentAdditionalInfo[key] = val
	}

	byteAdditionalInfo, err := json.Marshal(currentAdditionalInfo)
	if err != nil {
		return errors.Wrap(err, "failed to marshal merged additional info")
	}

	// Update status in DB
	err = s.repository.UpdateStatusTransaction(ctx, req.Reference, req.TransactionStatus, string(byteAdditionalInfo))
	if err != nil {
		return errors.Wrap(err, "failed to update status transaction")
	}

	trx.TransactionStatus = req.TransactionStatus
	s.sendNotification(ctx, tokenData, trx)

	return nil
}

func (s *service) GetTransactionDetail(ctx context.Context, reference string) (models.Transaction, error) {
	return s.repository.GetTransactionByReference(ctx, reference, true)
}

func (s *service) GetTransaction(ctx context.Context, userID uint64) ([]models.Transaction, error) {
	return s.repository.GetTransaction(ctx, userID)
}

func (s *service) RefundTransaction(ctx context.Context, tokenData models.TokenData, req *models.RefundTransaction) (models.CreateTransactionResponse, error) {

	var (
		resp models.CreateTransactionResponse
	)

	trx, err := s.repository.GetTransactionByReference(ctx, req.Reference, false)
	if err != nil {
		return resp, errors.Wrap(err, "failed to get transaction")
	}

	if trx.TransactionStatus != constants.TransactionStatusSuccess {
		return resp, errors.New("current transaction status is not success")
	}

	if trx.TransactionType != constants.TransactionTypePurchase {
		return resp, errors.New("current transaction status is not purchase type")
	}

	refundReference := "REFUND-" + req.Reference
	reqCreditBalance := external.UpdateBalance{
		Reference: refundReference,
		Amount:    trx.Amount,
	}

	_, err = s.external.CreditBalance(ctx, tokenData.Token, reqCreditBalance)
	if err != nil {
		return resp, errors.Wrap(err, "failed to credit balance")
	}

	transaction := models.Transaction{
		UserID:            tokenData.UserID,
		Amount:            trx.Amount,
		TransactionType:   constants.TransactionTypeRefund,
		TransactionStatus: constants.TransactionStatusSuccess,
		Reference:         refundReference,
		Description:       req.Description,
		AdditionalInfo:    req.AdditionalInfo,
	}

	err = s.repository.CreateTransaction(ctx, &transaction)
	if err != nil {
		return resp, errors.Wrap(err, "failed to insert new transaction refund")
	}

	resp.Reference = refundReference
	resp.TransactionStatus = transaction.TransactionStatus

	return resp, nil
}

func (s *service) sendNotification(ctx context.Context, tokenData models.TokenData, trx models.Transaction) {
	if trx.TransactionType == constants.TransactionTypePurchase && trx.TransactionStatus == constants.TransactionStatusSuccess {

		s.external.SendNotification(ctx, tokenData.Email, "purchase_success", map[string]string{
			"full_name":   tokenData.Fullname,
			"description": trx.Description,
			"reference":   trx.Reference,
			"date":        trx.CreatedAt.Format("2006-01-02 15:04:05"),
		})

	}
}
