package transaction

import (
	"context"
	"ewallet-transaction/constants"
	"ewallet-transaction/external"
	"ewallet-transaction/internal/models"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_service_CreateTransaction(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockRepo := NewMockrepository(ctrlMock)
	mockExt := NewMockIExternal(ctrlMock)

	type args struct {
		ctx context.Context
		req *models.Transaction
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		mockfn  func(args args)
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				req: &models.Transaction{
					UserID:            1,
					Amount:            100000,
					TransactionType:   "DEBIT",
					TransactionStatus: "PENDING",
					Description:       "DESCRIPTION",
					AdditionalInfo:    "{\"purchase\":\"testing purchase\"}",
				},
			},
			wantErr: false,
			mockfn: func(args args) {
				mockRepo.EXPECT().CreateTransaction(gomock.Any(), args.req).Return(nil)
			},
		},
		{
			name: "error",
			args: args{
				ctx: context.Background(),
				req: &models.Transaction{
					UserID:            1,
					Amount:            100000,
					TransactionType:   "DEBIT",
					TransactionStatus: "PENDING",
					Description:       "DESCRIPTION",
					AdditionalInfo:    "{\"purchase\":\"testing purchase\"}",
				},
			},
			wantErr: true,
			mockfn: func(args args) {
				mockRepo.EXPECT().CreateTransaction(gomock.Any(), args.req).Return(assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockfn(tt.args)
			s := &service{
				repository: mockRepo,
				external:   mockExt,
			}
			got, err := s.CreateTransaction(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.CreateTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.NotEmpty(t, got)
			} else {
				assert.Empty(t, got)
			}
		})
	}
}

func Test_service_UpdateStatusTransaction(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockRepo := NewMockrepository(ctrlMock)
	mockExt := NewMockIExternal(ctrlMock)

	now := time.Now()

	type args struct {
		ctx       context.Context
		tokenData models.TokenData
		req       *models.UpdateStatusTransaction
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "success update status from pending to success for topup",
			args: args{
				ctx: context.Background(),
				tokenData: models.TokenData{
					UserID:   1,
					Username: "USERNAME",
					Fullname: "FULLNAME",
					Token:    "TOKENDATA",
					Email:    "email@gmail.com",
				},
				req: &models.UpdateStatusTransaction{
					Reference:         "REFERENCE",
					TransactionStatus: "SUCCESS",
					AdditionalInfo:    "{\"purchase\":\"testing purchase\"}",
				},
			},
			wantErr: false,
			mockFn: func(args args) {
				mockRepo.EXPECT().GetTransactionByReference(gomock.Any(), args.req.Reference, false).Return(models.Transaction{
					ID:                1,
					UserID:            1,
					Amount:            100000,
					TransactionType:   constants.TransactionTypeTopup,
					TransactionStatus: "PENDING",
					Reference:         "REFERENCE",
					Description:       "DESCRIPTION",
					AdditionalInfo:    "{\"purchase\":\"testing purchase\"}",
					CreatedAt:         now,
					UpdatedAt:         now,
				}, nil)

				mockExt.EXPECT().CreditBalance(gomock.Any(), args.tokenData.Token, external.UpdateBalance{
					Reference: args.req.Reference,
					Amount:    100000,
				}).Return(&external.UpdateBalanceResponse{
					Message: constants.SuccessMessage,
					Amount:  100000,
				}, nil)

				mockRepo.EXPECT().UpdateStatusTransaction(gomock.Any(), args.req.Reference, args.req.TransactionStatus, args.req.AdditionalInfo).Return(nil)

			},
		},
		{
			name: "success update status from pending to success for purchase",
			args: args{
				ctx: context.Background(),
				tokenData: models.TokenData{
					UserID:   1,
					Username: "USERNAME",
					Fullname: "FULLNAME",
					Token:    "TOKENDATA",
					Email:    "email@gmail.com",
				},
				req: &models.UpdateStatusTransaction{
					Reference:         "REFERENCE",
					TransactionStatus: "SUCCESS",
					AdditionalInfo:    "{\"purchase\":\"testing purchase\"}",
				},
			},
			wantErr: false,
			mockFn: func(args args) {
				mockRepo.EXPECT().GetTransactionByReference(gomock.Any(), args.req.Reference, false).Return(models.Transaction{
					ID:                1,
					UserID:            1,
					Amount:            100000,
					TransactionType:   constants.TransactionTypePurchase,
					TransactionStatus: "PENDING",
					Reference:         "REFERENCE",
					Description:       "DESCRIPTION",
					AdditionalInfo:    "{\"purchase\":\"testing purchase\"}",
					CreatedAt:         now,
					UpdatedAt:         now,
				}, nil)

				mockExt.EXPECT().DebitBalance(gomock.Any(), args.tokenData.Token, external.UpdateBalance{
					Reference: args.req.Reference,
					Amount:    100000,
				}).Return(&external.UpdateBalanceResponse{
					Message: constants.SuccessMessage,
					Amount:  100000,
				}, nil)

				mockRepo.EXPECT().UpdateStatusTransaction(gomock.Any(), args.req.Reference, args.req.TransactionStatus, args.req.AdditionalInfo).Return(nil)

				mockExt.EXPECT().SendNotification(gomock.Any(), args.tokenData.Email, gomock.Any(), gomock.Any())
			},
		},
		{
			name: "success update status from pending to failed for topup",
			args: args{
				ctx: context.Background(),
				tokenData: models.TokenData{
					UserID:   1,
					Username: "USERNAME",
					Fullname: "FULLNAME",
					Token:    "TOKENDATA",
					Email:    "email@gmail.com",
				},
				req: &models.UpdateStatusTransaction{
					Reference:         "REFERENCE",
					TransactionStatus: "FAILED",
					AdditionalInfo:    "{\"purchase\":\"testing purchase\"}",
				},
			},
			wantErr: false,
			mockFn: func(args args) {
				mockRepo.EXPECT().GetTransactionByReference(gomock.Any(), args.req.Reference, false).Return(models.Transaction{
					ID:                1,
					UserID:            1,
					Amount:            100000,
					TransactionType:   constants.TransactionTypeTopup,
					TransactionStatus: "PENDING",
					Reference:         "REFERENCE",
					Description:       "DESCRIPTION",
					AdditionalInfo:    "{\"purchase\":\"testing purchase\"}",
					CreatedAt:         now,
					UpdatedAt:         now,
				}, nil)

				mockRepo.EXPECT().UpdateStatusTransaction(gomock.Any(), args.req.Reference, args.req.TransactionStatus, args.req.AdditionalInfo).Return(nil)
			},
		},
		{
			name: "success update status from pending to failed for purchase",
			args: args{
				ctx: context.Background(),
				tokenData: models.TokenData{
					UserID:   1,
					Username: "USERNAME",
					Fullname: "FULLNAME",
					Token:    "TOKENDATA",
					Email:    "email@gmail.com",
				},
				req: &models.UpdateStatusTransaction{
					Reference:         "REFERENCE",
					TransactionStatus: "FAILED",
					AdditionalInfo:    "{\"purchase\":\"testing purchase\"}",
				},
			},
			wantErr: false,
			mockFn: func(args args) {
				mockRepo.EXPECT().GetTransactionByReference(gomock.Any(), args.req.Reference, false).Return(models.Transaction{
					ID:                1,
					UserID:            1,
					Amount:            100000,
					TransactionType:   constants.TransactionTypePurchase,
					TransactionStatus: "PENDING",
					Reference:         "REFERENCE",
					Description:       "DESCRIPTION",
					AdditionalInfo:    "{\"purchase\":\"testing purchase\"}",
					CreatedAt:         now,
					UpdatedAt:         now,
				}, nil)

				mockRepo.EXPECT().UpdateStatusTransaction(gomock.Any(), args.req.Reference, args.req.TransactionStatus, args.req.AdditionalInfo).Return(nil)
			},
		},
		{
			name: "success update status from success to reversed for topup",
			args: args{
				ctx: context.Background(),
				tokenData: models.TokenData{
					UserID:   1,
					Username: "USERNAME",
					Fullname: "FULLNAME",
					Token:    "TOKENDATA",
					Email:    "email@gmail.com",
				},
				req: &models.UpdateStatusTransaction{
					Reference:         "REFERENCE",
					TransactionStatus: "REVERSED",
					AdditionalInfo:    "{\"purchase\":\"testing purchase\"}",
				},
			},
			wantErr: false,
			mockFn: func(args args) {
				mockRepo.EXPECT().GetTransactionByReference(gomock.Any(), args.req.Reference, false).Return(models.Transaction{
					ID:                1,
					UserID:            1,
					Amount:            100000,
					TransactionType:   constants.TransactionTypeTopup,
					TransactionStatus: "SUCCESS",
					Reference:         "REFERENCE",
					Description:       "DESCRIPTION",
					AdditionalInfo:    "{\"purchase\":\"testing purchase\"}",
					CreatedAt:         now,
					UpdatedAt:         now,
				}, nil)

				mockExt.EXPECT().DebitBalance(gomock.Any(), args.tokenData.Token, external.UpdateBalance{
					Reference: fmt.Sprintf("REVERSED-%s", args.req.Reference),
					Amount:    100000,
				}).Return(&external.UpdateBalanceResponse{
					Message: constants.SuccessMessage,
					Amount:  100000,
				}, nil)

				mockRepo.EXPECT().UpdateStatusTransaction(gomock.Any(), args.req.Reference, args.req.TransactionStatus, args.req.AdditionalInfo).Return(nil)

			},
		},
		{
			name: "success update status from success to reversed for purchase",
			args: args{
				ctx: context.Background(),
				tokenData: models.TokenData{
					UserID:   1,
					Username: "USERNAME",
					Fullname: "FULLNAME",
					Token:    "TOKENDATA",
					Email:    "email@gmail.com",
				},
				req: &models.UpdateStatusTransaction{
					Reference:         "REFERENCE",
					TransactionStatus: "REVERSED",
					AdditionalInfo:    "{\"purchase\":\"testing purchase\"}",
				},
			},
			wantErr: false,
			mockFn: func(args args) {
				mockRepo.EXPECT().GetTransactionByReference(gomock.Any(), args.req.Reference, false).Return(models.Transaction{
					ID:                1,
					UserID:            1,
					Amount:            100000,
					TransactionType:   constants.TransactionTypePurchase,
					TransactionStatus: "SUCCESS",
					Reference:         "REFERENCE",
					Description:       "DESCRIPTION",
					AdditionalInfo:    "{\"purchase\":\"testing purchase\"}",
					CreatedAt:         now,
					UpdatedAt:         now,
				}, nil)

				mockExt.EXPECT().CreditBalance(gomock.Any(), args.tokenData.Token, external.UpdateBalance{
					Reference: fmt.Sprintf("REVERSED-%s", args.req.Reference),
					Amount:    100000,
				}).Return(&external.UpdateBalanceResponse{
					Message: constants.SuccessMessage,
					Amount:  100000,
				}, nil)

				mockRepo.EXPECT().UpdateStatusTransaction(gomock.Any(), args.req.Reference, args.req.TransactionStatus, args.req.AdditionalInfo).Return(nil)

			},
		},
		{
			name: "success update status from failed to success for topup",
			args: args{
				ctx: context.Background(),
				tokenData: models.TokenData{
					UserID:   1,
					Username: "USERNAME",
					Fullname: "FULLNAME",
					Token:    "TOKENDATA",
					Email:    "email@gmail.com",
				},
				req: &models.UpdateStatusTransaction{
					Reference:         "REFERENCE",
					TransactionStatus: "SUCCESS",
					AdditionalInfo:    "{\"purchase\":\"testing purchase\"}",
				},
			},
			wantErr: false,
			mockFn: func(args args) {
				mockRepo.EXPECT().GetTransactionByReference(gomock.Any(), args.req.Reference, false).Return(models.Transaction{
					ID:                1,
					UserID:            1,
					Amount:            100000,
					TransactionType:   constants.TransactionTypeTopup,
					TransactionStatus: "FAILED",
					Reference:         "REFERENCE",
					Description:       "DESCRIPTION",
					AdditionalInfo:    "{\"purchase\":\"testing purchase\"}",
					CreatedAt:         now,
					UpdatedAt:         now,
				}, nil)

				mockExt.EXPECT().CreditBalance(gomock.Any(), args.tokenData.Token, external.UpdateBalance{
					Reference: args.req.Reference,
					Amount:    100000,
				}).Return(&external.UpdateBalanceResponse{
					Message: constants.SuccessMessage,
					Amount:  100000,
				}, nil)

				mockRepo.EXPECT().UpdateStatusTransaction(gomock.Any(), args.req.Reference, args.req.TransactionStatus, args.req.AdditionalInfo).Return(nil)

			},
		},
		{
			name: "error update status from pending to reversed",
			args: args{
				ctx: context.Background(),
				tokenData: models.TokenData{
					UserID:   1,
					Username: "USERNAME",
					Fullname: "FULLNAME",
					Token:    "TOKENDATA",
					Email:    "email@gmail.com",
				},
				req: &models.UpdateStatusTransaction{
					Reference:         "REFERENCE",
					TransactionStatus: "REVERSED",
					AdditionalInfo:    "{\"purchase\":\"testing purchase\"}",
				},
			},
			wantErr: true,
			mockFn: func(args args) {
				mockRepo.EXPECT().GetTransactionByReference(gomock.Any(), args.req.Reference, false).Return(models.Transaction{
					ID:                1,
					UserID:            1,
					Amount:            100000,
					TransactionType:   constants.TransactionTypeTopup,
					TransactionStatus: "PENDING",
					Reference:         "REFERENCE",
					Description:       "DESCRIPTION",
					AdditionalInfo:    "{\"purchase\":\"testing purchase\"}",
					CreatedAt:         now,
					UpdatedAt:         now,
				}, nil)
			},
		},
		{
			name: "error update status from pending to pending",
			args: args{
				ctx: context.Background(),
				tokenData: models.TokenData{
					UserID:   1,
					Username: "USERNAME",
					Fullname: "FULLNAME",
					Token:    "TOKENDATA",
					Email:    "email@gmail.com",
				},
				req: &models.UpdateStatusTransaction{
					Reference:         "REFERENCE",
					TransactionStatus: "PENDING",
					AdditionalInfo:    "{\"purchase\":\"testing purchase\"}",
				},
			},
			wantErr: true,
			mockFn: func(args args) {
				mockRepo.EXPECT().GetTransactionByReference(gomock.Any(), args.req.Reference, false).Return(models.Transaction{
					ID:                1,
					UserID:            1,
					Amount:            100000,
					TransactionType:   constants.TransactionTypeTopup,
					TransactionStatus: "PENDING",
					Reference:         "REFERENCE",
					Description:       "DESCRIPTION",
					AdditionalInfo:    "{\"purchase\":\"testing purchase\"}",
					CreatedAt:         now,
					UpdatedAt:         now,
				}, nil)
			},
		},
		{
			name: "error update status from success to failed",
			args: args{
				ctx: context.Background(),
				tokenData: models.TokenData{
					UserID:   1,
					Username: "USERNAME",
					Fullname: "FULLNAME",
					Token:    "TOKENDATA",
					Email:    "email@gmail.com",
				},
				req: &models.UpdateStatusTransaction{
					Reference:         "REFERENCE",
					TransactionStatus: "FAILED",
					AdditionalInfo:    "{\"purchase\":\"testing purchase\"}",
				},
			},
			wantErr: true,
			mockFn: func(args args) {
				mockRepo.EXPECT().GetTransactionByReference(gomock.Any(), args.req.Reference, false).Return(models.Transaction{
					ID:                1,
					UserID:            1,
					Amount:            100000,
					TransactionType:   constants.TransactionTypeTopup,
					TransactionStatus: "SUCCESS",
					Reference:         "REFERENCE",
					Description:       "DESCRIPTION",
					AdditionalInfo:    "{\"purchase\":\"testing purchase\"}",
					CreatedAt:         now,
					UpdatedAt:         now,
				}, nil)
			},
		},
		{
			name: "error update status from success to success",
			args: args{
				ctx: context.Background(),
				tokenData: models.TokenData{
					UserID:   1,
					Username: "USERNAME",
					Fullname: "FULLNAME",
					Token:    "TOKENDATA",
					Email:    "email@gmail.com",
				},
				req: &models.UpdateStatusTransaction{
					Reference:         "REFERENCE",
					TransactionStatus: "SUCCESS",
					AdditionalInfo:    "{\"purchase\":\"testing purchase\"}",
				},
			},
			wantErr: true,
			mockFn: func(args args) {
				mockRepo.EXPECT().GetTransactionByReference(gomock.Any(), args.req.Reference, false).Return(models.Transaction{
					ID:                1,
					UserID:            1,
					Amount:            100000,
					TransactionType:   constants.TransactionTypeTopup,
					TransactionStatus: "SUCCESS",
					Reference:         "REFERENCE",
					Description:       "DESCRIPTION",
					AdditionalInfo:    "{\"purchase\":\"testing purchase\"}",
					CreatedAt:         now,
					UpdatedAt:         now,
				}, nil)
			},
		},
		{
			name: "error update status from failed to reversed",
			args: args{
				ctx: context.Background(),
				tokenData: models.TokenData{
					UserID:   1,
					Username: "USERNAME",
					Fullname: "FULLNAME",
					Token:    "TOKENDATA",
					Email:    "email@gmail.com",
				},
				req: &models.UpdateStatusTransaction{
					Reference:         "REFERENCE",
					TransactionStatus: "REVERSED",
					AdditionalInfo:    "{\"purchase\":\"testing purchase\"}",
				},
			},
			wantErr: true,
			mockFn: func(args args) {
				mockRepo.EXPECT().GetTransactionByReference(gomock.Any(), args.req.Reference, false).Return(models.Transaction{
					ID:                1,
					UserID:            1,
					Amount:            100000,
					TransactionType:   constants.TransactionTypeTopup,
					TransactionStatus: "FAILED",
					Reference:         "REFERENCE",
					Description:       "DESCRIPTION",
					AdditionalInfo:    "{\"purchase\":\"testing purchase\"}",
					CreatedAt:         now,
					UpdatedAt:         now,
				}, nil)
			},
		},
		{
			name: "error update status from failed to failed",
			args: args{
				ctx: context.Background(),
				tokenData: models.TokenData{
					UserID:   1,
					Username: "USERNAME",
					Fullname: "FULLNAME",
					Token:    "TOKENDATA",
					Email:    "email@gmail.com",
				},
				req: &models.UpdateStatusTransaction{
					Reference:         "REFERENCE",
					TransactionStatus: "FAILED",
					AdditionalInfo:    "{\"purchase\":\"testing purchase\"}",
				},
			},
			wantErr: true,
			mockFn: func(args args) {
				mockRepo.EXPECT().GetTransactionByReference(gomock.Any(), args.req.Reference, false).Return(models.Transaction{
					ID:                1,
					UserID:            1,
					Amount:            100000,
					TransactionType:   constants.TransactionTypeTopup,
					TransactionStatus: "FAILED",
					Reference:         "REFERENCE",
					Description:       "DESCRIPTION",
					AdditionalInfo:    "{\"purchase\":\"testing purchase\"}",
					CreatedAt:         now,
					UpdatedAt:         now,
				}, nil)
			},
		},
		{
			name: "error update status from success to reversal because time limit",
			args: args{
				ctx: context.Background(),
				tokenData: models.TokenData{
					UserID:   1,
					Username: "USERNAME",
					Fullname: "FULLNAME",
					Token:    "TOKENDATA",
					Email:    "email@gmail.com",
				},
				req: &models.UpdateStatusTransaction{
					Reference:         "REFERENCE",
					TransactionStatus: "REVERSED",
					AdditionalInfo:    "{\"purchase\":\"testing purchase\"}",
				},
			},
			wantErr: true,
			mockFn: func(args args) {
				mockRepo.EXPECT().GetTransactionByReference(gomock.Any(), args.req.Reference, false).Return(models.Transaction{
					ID:                1,
					UserID:            1,
					Amount:            100000,
					TransactionType:   constants.TransactionTypeTopup,
					TransactionStatus: "SUCCESS",
					Reference:         "REFERENCE",
					Description:       "DESCRIPTION",
					AdditionalInfo:    "{\"purchase\":\"testing purchase\"}",
					CreatedAt:         now.Add(time.Hour * -36),
					UpdatedAt:         now,
				}, nil)
			},
		},
		{
			name: "error credit balance for topup",
			args: args{
				ctx: context.Background(),
				tokenData: models.TokenData{
					UserID:   1,
					Username: "USERNAME",
					Fullname: "FULLNAME",
					Token:    "TOKENDATA",
					Email:    "email@gmail.com",
				},
				req: &models.UpdateStatusTransaction{
					Reference:         "REFERENCE",
					TransactionStatus: "SUCCESS",
					AdditionalInfo:    "{\"purchase\":\"testing purchase\"}",
				},
			},
			wantErr: true,
			mockFn: func(args args) {
				mockRepo.EXPECT().GetTransactionByReference(gomock.Any(), args.req.Reference, false).Return(models.Transaction{
					ID:                1,
					UserID:            1,
					Amount:            100000,
					TransactionType:   constants.TransactionTypeTopup,
					TransactionStatus: "PENDING",
					Reference:         "REFERENCE",
					Description:       "DESCRIPTION",
					AdditionalInfo:    "{\"purchase\":\"testing purchase\"}",
					CreatedAt:         now,
					UpdatedAt:         now,
				}, nil)

				mockExt.EXPECT().CreditBalance(gomock.Any(), args.tokenData.Token, external.UpdateBalance{
					Reference: args.req.Reference,
					Amount:    100000,
				}).Return(&external.UpdateBalanceResponse{
					Message: constants.SuccessMessage,
					Amount:  100000,
				}, assert.AnError)
			},
		},
		{
			name: "error debit balance for topup",
			args: args{
				ctx: context.Background(),
				tokenData: models.TokenData{
					UserID:   1,
					Username: "USERNAME",
					Fullname: "FULLNAME",
					Token:    "TOKENDATA",
					Email:    "email@gmail.com",
				},
				req: &models.UpdateStatusTransaction{
					Reference:         "REFERENCE",
					TransactionStatus: "REVERSED",
					AdditionalInfo:    "{\"purchase\":\"testing purchase\"}",
				},
			},
			wantErr: true,
			mockFn: func(args args) {
				mockRepo.EXPECT().GetTransactionByReference(gomock.Any(), args.req.Reference, false).Return(models.Transaction{
					ID:                1,
					UserID:            1,
					Amount:            100000,
					TransactionType:   constants.TransactionTypeTopup,
					TransactionStatus: "SUCCESS",
					Reference:         "REFERENCE",
					Description:       "DESCRIPTION",
					AdditionalInfo:    "{\"purchase\":\"testing purchase\"}",
					CreatedAt:         now,
					UpdatedAt:         now,
				}, nil)

				mockExt.EXPECT().DebitBalance(gomock.Any(), args.tokenData.Token, external.UpdateBalance{
					Reference: fmt.Sprintf("REVERSED-%s", args.req.Reference),
					Amount:    100000,
				}).Return(&external.UpdateBalanceResponse{
					Message: constants.SuccessMessage,
					Amount:  100000,
				}, assert.AnError)
			},
		},
		{
			name: "error debit balance for purchase",
			args: args{
				ctx: context.Background(),
				tokenData: models.TokenData{
					UserID:   1,
					Username: "USERNAME",
					Fullname: "FULLNAME",
					Token:    "TOKENDATA",
					Email:    "email@gmail.com",
				},
				req: &models.UpdateStatusTransaction{
					Reference:         "REFERENCE",
					TransactionStatus: "SUCCESS",
					AdditionalInfo:    "{\"purchase\":\"testing purchase\"}",
				},
			},
			wantErr: true,
			mockFn: func(args args) {
				mockRepo.EXPECT().GetTransactionByReference(gomock.Any(), args.req.Reference, false).Return(models.Transaction{
					ID:                1,
					UserID:            1,
					Amount:            100000,
					TransactionType:   constants.TransactionTypePurchase,
					TransactionStatus: "PENDING",
					Reference:         "REFERENCE",
					Description:       "DESCRIPTION",
					AdditionalInfo:    "{\"purchase\":\"testing purchase\"}",
					CreatedAt:         now,
					UpdatedAt:         now,
				}, nil)

				mockExt.EXPECT().DebitBalance(gomock.Any(), args.tokenData.Token, external.UpdateBalance{
					Reference: args.req.Reference,
					Amount:    100000,
				}).Return(&external.UpdateBalanceResponse{
					Message: constants.SuccessMessage,
					Amount:  100000,
				}, nil)

				mockRepo.EXPECT().UpdateStatusTransaction(gomock.Any(), args.req.Reference, args.req.TransactionStatus, args.req.AdditionalInfo).Return(assert.AnError)
			},
		},
		{
			name: "error credit balance for purchase",
			args: args{
				ctx: context.Background(),
				tokenData: models.TokenData{
					UserID:   1,
					Username: "USERNAME",
					Fullname: "FULLNAME",
					Token:    "TOKENDATA",
					Email:    "email@gmail.com",
				},
				req: &models.UpdateStatusTransaction{
					Reference:         "REFERENCE",
					TransactionStatus: "REVERSED",
					AdditionalInfo:    "{\"purchase\":\"testing purchase\"}",
				},
			},
			wantErr: true,
			mockFn: func(args args) {
				mockRepo.EXPECT().GetTransactionByReference(gomock.Any(), args.req.Reference, false).Return(models.Transaction{
					ID:                1,
					UserID:            1,
					Amount:            100000,
					TransactionType:   constants.TransactionTypePurchase,
					TransactionStatus: "SUCCESS",
					Reference:         "REFERENCE",
					Description:       "DESCRIPTION",
					AdditionalInfo:    "{\"purchase\":\"testing purchase\"}",
					CreatedAt:         now,
					UpdatedAt:         now,
				}, nil)

				mockExt.EXPECT().CreditBalance(gomock.Any(), args.tokenData.Token, external.UpdateBalance{
					Reference: fmt.Sprintf("REVERSED-%s", args.req.Reference),
					Amount:    100000,
				}).Return(&external.UpdateBalanceResponse{
					Message: constants.SuccessMessage,
					Amount:  100000,
				}, assert.AnError)
			},
		},
		{
			name: "error updated status",
			args: args{
				ctx: context.Background(),
				tokenData: models.TokenData{
					UserID:   1,
					Username: "USERNAME",
					Fullname: "FULLNAME",
					Token:    "TOKENDATA",
					Email:    "email@gmail.com",
				},
				req: &models.UpdateStatusTransaction{
					Reference:         "REFERENCE",
					TransactionStatus: "SUCCESS",
					AdditionalInfo:    "{\"purchase\":\"testing purchase\"}",
				},
			},
			wantErr: true,
			mockFn: func(args args) {
				mockRepo.EXPECT().GetTransactionByReference(gomock.Any(), args.req.Reference, false).Return(models.Transaction{
					ID:                1,
					UserID:            1,
					Amount:            100000,
					TransactionType:   constants.TransactionTypeTopup,
					TransactionStatus: "PENDING",
					Reference:         "REFERENCE",
					Description:       "DESCRIPTION",
					AdditionalInfo:    "{\"purchase\":\"testing purchase\"}",
					CreatedAt:         now,
					UpdatedAt:         now,
				}, nil)

				mockExt.EXPECT().CreditBalance(gomock.Any(), args.tokenData.Token, external.UpdateBalance{
					Reference: args.req.Reference,
					Amount:    100000,
				}).Return(&external.UpdateBalanceResponse{
					Message: constants.SuccessMessage,
					Amount:  100000,
				}, nil)

				mockRepo.EXPECT().UpdateStatusTransaction(gomock.Any(), args.req.Reference, args.req.TransactionStatus, args.req.AdditionalInfo).Return(assert.AnError)

			},
		},
		{
			name: "error send notif",
			args: args{
				ctx: context.Background(),
				tokenData: models.TokenData{
					UserID:   1,
					Username: "USERNAME",
					Fullname: "FULLNAME",
					Token:    "TOKENDATA",
					Email:    "email@gmail.com",
				},
				req: &models.UpdateStatusTransaction{
					Reference:         "REFERENCE",
					TransactionStatus: "SUCCESS",
					AdditionalInfo:    "{\"purchase\":\"testing purchase\"}",
				},
			},
			wantErr: false,
			mockFn: func(args args) {
				mockRepo.EXPECT().GetTransactionByReference(gomock.Any(), args.req.Reference, false).Return(models.Transaction{
					ID:                1,
					UserID:            1,
					Amount:            100000,
					TransactionType:   constants.TransactionTypePurchase,
					TransactionStatus: "PENDING",
					Reference:         "REFERENCE",
					Description:       "DESCRIPTION",
					AdditionalInfo:    "{\"purchase\":\"testing purchase\"}",
					CreatedAt:         now,
					UpdatedAt:         now,
				}, nil)

				mockExt.EXPECT().DebitBalance(gomock.Any(), args.tokenData.Token, external.UpdateBalance{
					Reference: args.req.Reference,
					Amount:    100000,
				}).Return(&external.UpdateBalanceResponse{
					Message: constants.SuccessMessage,
					Amount:  100000,
				}, nil)

				mockRepo.EXPECT().UpdateStatusTransaction(gomock.Any(), args.req.Reference, args.req.TransactionStatus, args.req.AdditionalInfo).Return(nil)

				mockExt.EXPECT().SendNotification(gomock.Any(), args.tokenData.Email, gomock.Any(), gomock.Any()).Return(assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			s := &service{
				repository: mockRepo,
				external:   mockExt,
			}
			if err := s.UpdateStatusTransaction(tt.args.ctx, tt.args.tokenData, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("service.UpdateStatusTransaction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_service_GetTransactionDetail(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockRepo := NewMockrepository(ctrlMock)
	mockExt := NewMockIExternal(ctrlMock)

	now := time.Now()
	transaction := models.Transaction{
		ID:                1,
		UserID:            1,
		Amount:            200000,
		TransactionType:   constants.TransactionTypeTopup,
		TransactionStatus: "PENDING",
		Reference:         "REFERENCE",
		Description:       "DESCRIPTION",
		AdditionalInfo:    "{\"purchase\":\"testing purchase\"}",
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	type args struct {
		ctx       context.Context
		reference string
	}
	tests := []struct {
		name    string
		args    args
		want    models.Transaction
		wantErr bool
		mockfn  func(args args)
	}{
		{
			name: "success",
			args: args{
				ctx:       context.Background(),
				reference: "REFERENCE",
			},
			want:    transaction,
			wantErr: false,
			mockfn: func(args args) {
				mockRepo.EXPECT().GetTransactionByReference(gomock.Any(), args.reference, true).Return(transaction, nil)
			},
		},
		{
			name: "error",
			args: args{
				ctx:       context.Background(),
				reference: "REFERENCE",
			},
			want:    transaction,
			wantErr: true,
			mockfn: func(args args) {
				mockRepo.EXPECT().GetTransactionByReference(gomock.Any(), args.reference, true).Return(transaction, assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockfn(tt.args)
			s := &service{
				repository: mockRepo,
				external:   mockExt,
			}
			got, err := s.GetTransactionDetail(tt.args.ctx, tt.args.reference)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.GetTransactionDetail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.GetTransactionDetail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_GetTransaction(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockRepo := NewMockrepository(ctrlMock)
	mockExt := NewMockIExternal(ctrlMock)
	now := time.Now()

	transactions := []models.Transaction{
		{
			ID:                1,
			UserID:            1,
			Amount:            200000,
			TransactionType:   constants.TransactionTypePurchase,
			TransactionStatus: constants.TransactionStatusSuccess,
			Reference:         "REFERENCE",
			Description:       "DESC",
			AdditionalInfo:    "ADDINFO",
			CreatedAt:         now,
			UpdatedAt:         now,
		},
		{
			ID:                2,
			UserID:            1,
			Amount:            300000,
			TransactionType:   constants.TransactionTypePurchase,
			TransactionStatus: constants.TransactionStatusPending,
			Reference:         "REFERENCE",
			Description:       "DESC",
			AdditionalInfo:    "ADDINFO",
			CreatedAt:         now,
			UpdatedAt:         now,
		},
	}

	type args struct {
		ctx    context.Context
		userID uint64
	}
	tests := []struct {
		name    string
		args    args
		want    []models.Transaction
		wantErr bool
		mockfn  func(args args)
	}{
		{
			name: "success",
			args: args{
				ctx:    context.Background(),
				userID: 1,
			},
			want:    transactions,
			wantErr: false,
			mockfn: func(args args) {
				mockRepo.EXPECT().GetTransaction(gomock.Any(), args.userID).Return(transactions, nil)
			},
		},
		{
			name: "error",
			args: args{
				ctx:    context.Background(),
				userID: 1,
			},
			want:    nil,
			wantErr: true,
			mockfn: func(args args) {
				mockRepo.EXPECT().GetTransaction(gomock.Any(), args.userID).Return(nil, assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockfn(tt.args)
			s := &service{
				repository: mockRepo,
				external:   mockExt,
			}
			got, err := s.GetTransaction(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.GetTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.GetTransaction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_RefundTransaction(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockRepo := NewMockrepository(ctrlMock)
	mockExt := NewMockIExternal(ctrlMock)
	now := time.Now()

	type args struct {
		ctx       context.Context
		tokenData models.TokenData
		req       *models.RefundTransaction
	}
	tests := []struct {
		name    string
		args    args
		want    models.CreateTransactionResponse
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				tokenData: models.TokenData{
					UserID:   1,
					Username: "USERNAME",
					Fullname: "FULLNAME",
					Token:    "TOKEN",
					Email:    "EMAIL@gmail.com",
				},
				req: &models.RefundTransaction{
					Reference:      "REFERENCE",
					Description:    "DESCRIPTION",
					AdditionalInfo: "ADDINFO",
				},
			},
			want: models.CreateTransactionResponse{
				Reference:         "REFUND-REFERENCE",
				TransactionStatus: constants.TransactionStatusReversed,
			},
			wantErr: false,
			mockFn: func(args args) {
				trx := models.Transaction{
					ID:                1,
					UserID:            1,
					Amount:            200000,
					TransactionType:   constants.TransactionTypePurchase,
					TransactionStatus: constants.TransactionStatusSuccess,
					Reference:         "REFERENCE",
					Description:       "DESCRIPTION",
					AdditionalInfo:    "ADDINFO",
					CreatedAt:         now,
					UpdatedAt:         now,
				}
				mockRepo.EXPECT().GetTransactionByReference(gomock.Any(), args.req.Reference, false).Return(trx, nil)

				mockExt.EXPECT().CreditBalance(gomock.Any(), args.tokenData.Token, external.UpdateBalance{
					Reference: "REFUND-REFERENCE",
					Amount:    200000,
				}).Return(&external.UpdateBalanceResponse{
					Message: constants.SuccessMessage,
					Amount:  200000,
				}, nil)

				mockRepo.EXPECT().CreateTransaction(gomock.Any(), &models.Transaction{
					UserID:            1,
					Amount:            trx.Amount,
					TransactionType:   constants.TransactionTypeRefund,
					TransactionStatus: constants.TransactionStatusSuccess,
					Reference:         "REFUND-REFERENCE",
					Description:       args.req.Description,
					AdditionalInfo:    args.req.AdditionalInfo,
				}).Do(func(ctx context.Context, trx *models.Transaction) {
					trx.TransactionStatus = constants.TransactionStatusReversed
				}).Return(nil)

			},
		},
		{
			name: "error when transaction not success",
			args: args{
				ctx: context.Background(),
				tokenData: models.TokenData{
					UserID:   1,
					Username: "USERNAME",
					Fullname: "FULLNAME",
					Token:    "TOKEN",
					Email:    "EMAIL@gmail.com",
				},
				req: &models.RefundTransaction{
					Reference:      "REFERENCE",
					Description:    "DESCRIPTION",
					AdditionalInfo: "ADDINFO",
				},
			},
			want:    models.CreateTransactionResponse{},
			wantErr: true,
			mockFn: func(args args) {
				trx := models.Transaction{
					ID:                1,
					UserID:            1,
					Amount:            200000,
					TransactionType:   constants.TransactionTypePurchase,
					TransactionStatus: constants.TransactionStatusFailed,
					Reference:         "REFERENCE",
					Description:       "DESCRIPTION",
					AdditionalInfo:    "ADDINFO",
					CreatedAt:         now,
					UpdatedAt:         now,
				}
				mockRepo.EXPECT().GetTransactionByReference(gomock.Any(), args.req.Reference, false).Return(trx, nil)

			},
		},
		{
			name: "error when transaction not purchase",
			args: args{
				ctx: context.Background(),
				tokenData: models.TokenData{
					UserID:   1,
					Username: "USERNAME",
					Fullname: "FULLNAME",
					Token:    "TOKEN",
					Email:    "EMAIL@gmail.com",
				},
				req: &models.RefundTransaction{
					Reference:      "REFERENCE",
					Description:    "DESCRIPTION",
					AdditionalInfo: "ADDINFO",
				},
			},
			want:    models.CreateTransactionResponse{},
			wantErr: true,
			mockFn: func(args args) {
				trx := models.Transaction{
					ID:                1,
					UserID:            1,
					Amount:            200000,
					TransactionType:   constants.TransactionTypeRefund,
					TransactionStatus: constants.TransactionStatusSuccess,
					Reference:         "REFERENCE",
					Description:       "DESCRIPTION",
					AdditionalInfo:    "ADDINFO",
					CreatedAt:         now,
					UpdatedAt:         now,
				}
				mockRepo.EXPECT().GetTransactionByReference(gomock.Any(), args.req.Reference, false).Return(trx, nil)

			},
		},
		{
			name: "error credit balance",
			args: args{
				ctx: context.Background(),
				tokenData: models.TokenData{
					UserID:   1,
					Username: "USERNAME",
					Fullname: "FULLNAME",
					Token:    "TOKEN",
					Email:    "EMAIL@gmail.com",
				},
				req: &models.RefundTransaction{
					Reference:      "REFERENCE",
					Description:    "DESCRIPTION",
					AdditionalInfo: "ADDINFO",
				},
			},
			want:    models.CreateTransactionResponse{},
			wantErr: true,
			mockFn: func(args args) {
				trx := models.Transaction{
					ID:                1,
					UserID:            1,
					Amount:            200000,
					TransactionType:   constants.TransactionTypePurchase,
					TransactionStatus: constants.TransactionStatusSuccess,
					Reference:         "REFERENCE",
					Description:       "DESCRIPTION",
					AdditionalInfo:    "ADDINFO",
					CreatedAt:         now,
					UpdatedAt:         now,
				}
				mockRepo.EXPECT().GetTransactionByReference(gomock.Any(), args.req.Reference, false).Return(trx, nil)

				mockExt.EXPECT().CreditBalance(gomock.Any(), args.tokenData.Token, external.UpdateBalance{
					Reference: "REFUND-REFERENCE",
					Amount:    200000,
				}).Return(&external.UpdateBalanceResponse{
					Message: constants.SuccessMessage,
					Amount:  200000,
				}, assert.AnError)

			},
		},
		{
			name: "error create transaction",
			args: args{
				ctx: context.Background(),
				tokenData: models.TokenData{
					UserID:   1,
					Username: "USERNAME",
					Fullname: "FULLNAME",
					Token:    "TOKEN",
					Email:    "EMAIL@gmail.com",
				},
				req: &models.RefundTransaction{
					Reference:      "REFERENCE",
					Description:    "DESCRIPTION",
					AdditionalInfo: "ADDINFO",
				},
			},
			want:    models.CreateTransactionResponse{},
			wantErr: true,
			mockFn: func(args args) {
				trx := models.Transaction{
					ID:                1,
					UserID:            1,
					Amount:            200000,
					TransactionType:   constants.TransactionTypePurchase,
					TransactionStatus: constants.TransactionStatusSuccess,
					Reference:         "REFERENCE",
					Description:       "DESCRIPTION",
					AdditionalInfo:    "ADDINFO",
					CreatedAt:         now,
					UpdatedAt:         now,
				}
				mockRepo.EXPECT().GetTransactionByReference(gomock.Any(), args.req.Reference, false).Return(trx, nil)

				mockExt.EXPECT().CreditBalance(gomock.Any(), args.tokenData.Token, external.UpdateBalance{
					Reference: "REFUND-REFERENCE",
					Amount:    200000,
				}).Return(&external.UpdateBalanceResponse{
					Message: constants.SuccessMessage,
					Amount:  200000,
				}, nil)

				mockRepo.EXPECT().CreateTransaction(gomock.Any(), &models.Transaction{
					UserID:            1,
					Amount:            trx.Amount,
					TransactionType:   constants.TransactionTypeRefund,
					TransactionStatus: constants.TransactionStatusSuccess,
					Reference:         "REFUND-REFERENCE",
					Description:       args.req.Description,
					AdditionalInfo:    args.req.AdditionalInfo,
				}).Return(assert.AnError)

			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			s := &service{
				repository: mockRepo,
				external:   mockExt,
			}
			got, err := s.RefundTransaction(tt.args.ctx, tt.args.tokenData, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.RefundTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.RefundTransaction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_sendNotification(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockRepo := NewMockrepository(ctrlMock)
	mockExt := NewMockIExternal(ctrlMock)
	now := time.Now()

	type args struct {
		ctx       context.Context
		tokenData models.TokenData
		trx       models.Transaction
	}
	tests := []struct {
		name   string
		args   args
		mockFn func(args args)
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				tokenData: models.TokenData{
					UserID:   1,
					Username: "USERNAME",
					Fullname: "FULLNAME",
					Token:    "TOKEN",
					Email:    "EMAIL.@gmail.com",
				},
				trx: models.Transaction{
					ID:                1,
					UserID:            1,
					Amount:            200000,
					TransactionType:   constants.TransactionTypePurchase,
					TransactionStatus: constants.TransactionStatusSuccess,
					Reference:         "REFERENCE",
					Description:       "DESCRIPTION",
					AdditionalInfo:    "ADDINFO",
					CreatedAt:         now,
					UpdatedAt:         now,
				},
			},
			mockFn: func(args args) {
				mockExt.EXPECT().SendNotification(gomock.Any(), args.tokenData.Email, "purchase_success", map[string]string{
					"full_name":   args.tokenData.Fullname,
					"description": args.trx.Description,
					"reference":   args.trx.Reference,
					"date":        args.trx.CreatedAt.Format("2006-01-02 15:04:05"),
				}).Return(nil)
			},
		},
		{
			name: "error",
			args: args{
				ctx: context.Background(),
				tokenData: models.TokenData{
					UserID:   1,
					Username: "USERNAME",
					Fullname: "FULLNAME",
					Token:    "TOKEN",
					Email:    "EMAIL.@gmail.com",
				},
				trx: models.Transaction{
					ID:                1,
					UserID:            1,
					Amount:            200000,
					TransactionType:   constants.TransactionTypePurchase,
					TransactionStatus: constants.TransactionStatusSuccess,
					Reference:         "REFERENCE",
					Description:       "DESCRIPTION",
					AdditionalInfo:    "ADDINFO",
					CreatedAt:         now,
					UpdatedAt:         now,
				},
			},
			mockFn: func(args args) {
				mockExt.EXPECT().SendNotification(gomock.Any(), args.tokenData.Email, "purchase_success", map[string]string{
					"full_name":   args.tokenData.Fullname,
					"description": args.trx.Description,
					"reference":   args.trx.Reference,
					"date":        args.trx.CreatedAt.Format("2006-01-02 15:04:05"),
				}).Return(assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			s := &service{
				repository: mockRepo,
				external:   mockExt,
			}
			s.sendNotification(tt.args.ctx, tt.args.tokenData, tt.args.trx)
		})
	}
}
