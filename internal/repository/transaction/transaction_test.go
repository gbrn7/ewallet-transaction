package transaction

import (
	"context"
	"ewallet-transaction/constants"
	"ewallet-transaction/internal/models"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Test_repository_CreateTransaction(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})

	assert.NoError(t, err)

	type args struct {
		ctx context.Context
		trx *models.Transaction
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				trx: &models.Transaction{
					UserID:            1,
					Amount:            100000,
					TransactionType:   "DEBIT",
					TransactionStatus: "PENDING",
					Reference:         "REFERENCE",
					Description:       "DESCRIPTION",
					AdditionalInfo:    "",
				},
			},
			wantErr: false,
			mockFn: func(args args) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `transactions` (`user_id`,`amount`,`transaction_type`,`transaction_status`,`reference`,`description`,`additional_info`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?,?,?,?)")).WithArgs(
					args.trx.UserID,
					args.trx.Amount,
					args.trx.TransactionType,
					args.trx.TransactionStatus,
					args.trx.Reference,
					args.trx.Description,
					args.trx.AdditionalInfo,
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
				).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		{
			name: "error",
			args: args{
				ctx: context.Background(),
				trx: &models.Transaction{
					UserID:            1,
					Amount:            100000,
					TransactionType:   "DEBIT",
					TransactionStatus: "PENDING",
					Reference:         "REFERENCE",
					Description:       "DESCRIPTION",
					AdditionalInfo:    "",
				},
			},
			wantErr: true,
			mockFn: func(args args) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `transactions` (`user_id`,`amount`,`transaction_type`,`transaction_status`,`reference`,`description`,`additional_info`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?,?,?,?)")).WithArgs(
					args.trx.UserID,
					args.trx.Amount,
					args.trx.TransactionType,
					args.trx.TransactionStatus,
					args.trx.Reference,
					args.trx.Description,
					args.trx.AdditionalInfo,
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
				).WillReturnError(assert.AnError)
				mock.ExpectRollback()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			r := &repository{
				DB: gormDB,
			}
			if err := r.CreateTransaction(tt.args.ctx, tt.args.trx); (err != nil) != tt.wantErr {
				t.Errorf("repository.CreateTransaction() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func Test_repository_GetTransactionByReference(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})

	assert.NoError(t, err)

	now := time.Now()

	type args struct {
		ctx           context.Context
		reference     string
		includeRefund bool
	}
	tests := []struct {
		name    string
		args    args
		want    models.Transaction
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "success without include refund",
			args: args{
				ctx:           context.Background(),
				reference:     "REFERENCE",
				includeRefund: false,
			},
			want: models.Transaction{
				ID:                1,
				UserID:            1,
				Amount:            100000,
				TransactionType:   "DEBIT",
				TransactionStatus: "PENDING",
				Reference:         "REFERENCE",
				Description:       "DESCRIPTION",
				AdditionalInfo:    "ADDINFO",
				CreatedAt:         now,
				UpdatedAt:         now,
			},
			wantErr: false,
			mockFn: func(args args) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `transactions` WHERE reference = ? AND transaction_type != ? ORDER BY `transactions`.`id` DESC LIMIT ?")).WithArgs(
					args.reference,
					constants.TransactionTypeRefund,
					1,
				).WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "amount", "transaction_type", "transaction_status", "reference", "description", "additional_info", "created_at", "updated_at"}).AddRow(1, 1, 100000, "DEBIT", "PENDING", "REFERENCE", "DESCRIPTION", "ADDINFO", now, now))
			},
		},
		{
			name: "success with include refund",
			args: args{
				ctx:           context.Background(),
				reference:     "REFERENCE",
				includeRefund: true,
			},
			want: models.Transaction{
				ID:                1,
				UserID:            1,
				Amount:            100000,
				TransactionType:   "DEBIT",
				TransactionStatus: "PENDING",
				Reference:         "REFERENCE",
				Description:       "DESCRIPTION",
				AdditionalInfo:    "ADDINFO",
				CreatedAt:         now,
				UpdatedAt:         now,
			},
			wantErr: false,
			mockFn: func(args args) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `transactions` WHERE reference = ? ORDER BY `transactions`.`id` DESC LIMIT ?")).WithArgs(
					args.reference,
					1,
				).WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "amount", "transaction_type", "transaction_status", "reference", "description", "additional_info", "created_at", "updated_at"}).AddRow(1, 1, 100000, "DEBIT", "PENDING", "REFERENCE", "DESCRIPTION", "ADDINFO", now, now))
			},
		},
		{
			name: "error without include refund",
			args: args{
				ctx:           context.Background(),
				reference:     "REFERENCE",
				includeRefund: false,
			},
			want:    models.Transaction{},
			wantErr: true,
			mockFn: func(args args) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `transactions` WHERE reference = ? AND transaction_type != ? ORDER BY `transactions`.`id` DESC LIMIT ?")).WithArgs(
					args.reference,
					constants.TransactionTypeRefund,
					1,
				).WillReturnError(assert.AnError)
			},
		},
		{
			name: "error with include refund",
			args: args{
				ctx:           context.Background(),
				reference:     "reference",
				includeRefund: true,
			},
			want:    models.Transaction{},
			wantErr: true,
			mockFn: func(args args) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `transactions` WHERE reference = ? ORDER BY `transactions`.`id` DESC LIMIT ?")).WithArgs(
					args.reference,
					1,
				).WillReturnError(assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			r := &repository{
				DB: gormDB,
			}
			got, err := r.GetTransactionByReference(tt.args.ctx, tt.args.reference, tt.args.includeRefund)
			if (err != nil) != tt.wantErr {
				t.Errorf("repository.GetTransactionByReference() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("repository.GetTransactionByReference() = %v, want %v", got, tt.want)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func Test_repository_UpdateStatusTransaction(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})

	assert.NoError(t, err)

	type args struct {
		ctx            context.Context
		reference      string
		status         string
		additionalInfo string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "success",
			args: args{
				ctx:            context.Background(),
				reference:      "REFERENCE",
				status:         "SUCCESS",
				additionalInfo: "ADDINFO",
			},
			wantErr: false,
			mockFn: func(args args) {
				mock.ExpectExec(regexp.QuoteMeta("UPDATE transactions SET transaction_status = ?, additional_info = ? WHERE reference = ?")).WithArgs(
					args.status,
					args.additionalInfo,
					args.reference,
				).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "error",
			args: args{
				ctx:            context.Background(),
				reference:      "REFERENCE",
				status:         "SUCCESS",
				additionalInfo: "ADDINFO",
			},
			wantErr: true,
			mockFn: func(args args) {
				mock.ExpectExec(regexp.QuoteMeta("UPDATE transactions SET transaction_status = ?, additional_info = ? WHERE reference = ?")).WithArgs(
					args.status,
					args.additionalInfo,
					args.reference,
				).WillReturnError(assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			r := &repository{
				DB: gormDB,
			}
			if err := r.UpdateStatusTransaction(tt.args.ctx, tt.args.reference, tt.args.status, tt.args.additionalInfo); (err != nil) != tt.wantErr {
				t.Errorf("repository.UpdateStatusTransaction() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func Test_repository_GetTransaction(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})

	assert.NoError(t, err)

	now := time.Now()
	type args struct {
		ctx    context.Context
		userID uint64
	}
	tests := []struct {
		name    string
		args    args
		want    []models.Transaction
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "success",
			args: args{
				ctx:    context.Background(),
				userID: 1,
			},
			want: []models.Transaction{
				{
					ID:                1,
					UserID:            1,
					Amount:            100000,
					TransactionType:   "DEBIT",
					TransactionStatus: "SUCCESS",
					Reference:         "REFERENCE",
					Description:       "DESCRIPTION",
					AdditionalInfo:    "ADDINFO",
					CreatedAt:         now,
					UpdatedAt:         now,
				},
				{
					ID:                2,
					UserID:            1,
					Amount:            100000,
					TransactionType:   "DEBIT",
					TransactionStatus: "PENDING",
					Reference:         "REFERENCE",
					Description:       "DESCRIPTION",
					AdditionalInfo:    "ADDINFO",
					CreatedAt:         now,
					UpdatedAt:         now,
				},
			},
			wantErr: false,
			mockFn: func(args args) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `transactions` WHERE user_id = ? ORDER BY id DESC")).WithArgs(
					args.userID,
				).WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "amount", "transaction_type", "transaction_status", "reference", "description", "additional_info", "created_at", "updated_at"}).AddRow(1, 1, 100000, "DEBIT", "SUCCESS", "REFERENCE", "DESCRIPTION", "ADDINFO", now, now).AddRow(2, 1, 100000, "DEBIT", "PENDING", "REFERENCE", "DESCRIPTION", "ADDINFO", now, now))
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
			mockFn: func(args args) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `transactions` WHERE user_id = ? ORDER BY id DESC")).WithArgs(
					args.userID,
				).WillReturnError(assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			r := &repository{
				DB: gormDB,
			}
			got, err := r.GetTransaction(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("repository.GetTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("repository.GetTransaction() = %v, want %v", got, tt.want)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
