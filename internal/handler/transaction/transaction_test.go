package transaction

import (
	"bytes"
	"encoding/json"
	"ewallet-transaction/constants"
	"ewallet-transaction/helpers"
	"ewallet-transaction/internal/models"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func TestHandler_CreateTransaction(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockSvc := NewMockService(ctrlMock)
	mockExt := NewMockExternal(ctrlMock)
	mockMdw := NewMockMiddleware(ctrlMock)

	trx := models.Transaction{
		UserID:            1,
		Amount:            200000,
		TransactionType:   constants.TransactionTypePurchase,
		Reference:         "REFERENCE",
		Description:       "DESC",
		TransactionStatus: "PENDING",
		AdditionalInfo:    "ADDINFO",
	}

	tests := []struct {
		name               string
		expectedStatusCode int
		expectedBody       helpers.Response
		wantErr            bool
		mockFn             func()
	}{
		{
			name:               "success",
			expectedStatusCode: http.StatusOK,
			expectedBody: helpers.Response{
				Message: constants.SuccessMessage,
				Data: map[string]interface{}{
					"reference":          "REFERENCE",
					"transaction_status": "PENDING",
				},
			},
			wantErr: false,
			mockFn: func() {
				mockMdw.EXPECT().MiddlewareValidateToken(gomock.Any()).Do(func(c *gin.Context) {
					tokenData := models.TokenData{
						UserID:   1,
						Username: "USERNAME",
						Fullname: "FULLNAME",
						Token:    "TOKEN",
						Email:    "EMAIL",
					}

					c.Set("token", tokenData)
					c.Next()
				})

				mockSvc.EXPECT().CreateTransaction(gomock.Any(), &trx).Return(models.CreateTransactionResponse{
					Reference:         "REFERENCE",
					TransactionStatus: "PENDING",
				}, nil)
			},
		},
		{
			name:               "error",
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody: helpers.Response{
				Message: constants.SuccessMessage,
			},
			wantErr: true,
			mockFn: func() {
				mockMdw.EXPECT().MiddlewareValidateToken(gomock.Any()).Do(func(c *gin.Context) {
					tokenData := models.TokenData{
						UserID:   1,
						Username: "USERNAME",
						Fullname: "FULLNAME",
						Token:    "TOKEN",
						Email:    "EMAIL",
					}

					c.Set("token", tokenData)
					c.Next()
				})

				mockSvc.EXPECT().CreateTransaction(gomock.Any(), &trx).Return(models.CreateTransactionResponse{}, assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()

			api := gin.New()
			h := &Handler{
				Engine:     api,
				Service:    mockSvc,
				External:   mockExt,
				Middleware: mockMdw,
			}
			h.RegisterRoute()
			w := httptest.NewRecorder()

			endpoint := "/transaction/v1/create"
			val, err := json.Marshal(trx)
			assert.NoError(t, err)

			body := bytes.NewReader(val)
			req, err := http.NewRequest(http.MethodPost, endpoint, body)
			assert.NoError(t, err)
			req.Header.Set("Authorization", "authorization")

			h.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			if !tt.wantErr {
				res := w.Result()
				defer res.Body.Close()

				response := helpers.Response{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				assert.Equal(t, tt.expectedBody, response)
			}
		})
	}
}

func TestHandler_UpdateStatusTransaction(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockSvc := NewMockService(ctrlMock)
	mockExt := NewMockExternal(ctrlMock)
	mockMdw := NewMockMiddleware(ctrlMock)

	tokenData := models.TokenData{
		UserID:   1,
		Username: "USERNAME",
		Fullname: "FULLNAME",
		Token:    "TOKEN",
		Email:    "EMAIL",
	}

	req := models.UpdateStatusTransaction{
		Reference:         "REFERENCE",
		TransactionStatus: "SUCCESS",
		AdditionalInfo:    "ADDINFO",
	}

	tests := []struct {
		name               string
		expectedStatusCode int
		expectedBody       helpers.Response
		wantErr            bool
		mockFn             func()
	}{
		{
			name:               "success",
			expectedStatusCode: http.StatusOK,
			expectedBody: helpers.Response{
				Message: constants.SuccessMessage,
			},
			wantErr: false,
			mockFn: func() {
				mockMdw.EXPECT().MiddlewareValidateToken(gomock.Any()).Do(func(c *gin.Context) {
					tokenData := tokenData

					c.Set("token", tokenData)
					c.Next()
				})

				mockSvc.EXPECT().UpdateStatusTransaction(gomock.Any(), tokenData, &req).Return(nil)
			},
		},
		{
			name:               "error",
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody: helpers.Response{
				Message: constants.ErrServerError,
			},
			wantErr: true,
			mockFn: func() {
				mockMdw.EXPECT().MiddlewareValidateToken(gomock.Any()).Do(func(c *gin.Context) {
					tokenData := tokenData

					c.Set("token", tokenData)
					c.Next()
				})

				mockSvc.EXPECT().UpdateStatusTransaction(gomock.Any(), tokenData, &req).Return(assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()

			api := gin.New()
			h := &Handler{
				Engine:     api,
				Service:    mockSvc,
				External:   mockExt,
				Middleware: mockMdw,
			}
			h.RegisterRoute()
			w := httptest.NewRecorder()

			endpoint := fmt.Sprintf("/transaction/v1/update-status/%s", req.Reference)
			val, err := json.Marshal(req)
			assert.NoError(t, err)

			body := bytes.NewReader(val)
			req, err := http.NewRequest(http.MethodPut, endpoint, body)
			assert.NoError(t, err)
			req.Header.Set("Authorization", "authorization")

			h.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			if !tt.wantErr {
				res := w.Result()
				defer res.Body.Close()

				response := helpers.Response{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				assert.Equal(t, tt.expectedBody, response)
			}
		})
	}
}

func TestHandler_GetTransaction(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockSvc := NewMockService(ctrlMock)
	mockExt := NewMockExternal(ctrlMock)
	mockMdw := NewMockMiddleware(ctrlMock)

	now := time.Now()

	tokenData := models.TokenData{
		UserID:   1,
		Username: "USERNAME",
		Fullname: "FULLNAME",
		Token:    "TOKEN",
		Email:    "EMAIL",
	}

	tests := []struct {
		name               string
		expectedStatusCode int
		expectedBody       helpers.Response
		wantErr            bool
		mockFn             func()
	}{
		{
			name:               "success",
			expectedStatusCode: http.StatusOK,
			expectedBody: helpers.Response{
				Message: constants.SuccessMessage,
				Data: []interface{}{
					map[string]interface{}{
						"id":                 float64(1),
						"user_id":            float64(1),
						"amount":             float64(100000),
						"transaction_type":   constants.TransactionTypeTopup,
						"transaction_status": constants.TransactionStatusSuccess,
						"reference":          "REFERENCE",
						"description":        "DESC",
						"additional_info":    "ADDINFO",
						"created_at":         now.Format(time.RFC3339Nano),
						"updated_at":         now.Format(time.RFC3339Nano),
					},
					map[string]interface{}{
						"id":                 float64(2),
						"user_id":            float64(1),
						"amount":             float64(300000),
						"transaction_type":   constants.TransactionTypePurchase,
						"transaction_status": constants.TransactionStatusPending,
						"reference":          "REFERENCE",
						"description":        "DESC",
						"additional_info":    "ADDINFO",
						"created_at":         now.Format(time.RFC3339Nano),
						"updated_at":         now.Format(time.RFC3339Nano),
					},
				},
			},
			wantErr: false,
			mockFn: func() {
				mockMdw.EXPECT().MiddlewareValidateToken(gomock.Any()).Do(func(c *gin.Context) {
					tokenData := tokenData

					c.Set("token", tokenData)
					c.Next()
				})

				mockSvc.EXPECT().GetTransaction(gomock.Any(), tokenData.UserID).Return([]models.Transaction{
					{
						ID:                1,
						UserID:            1,
						Amount:            100000,
						TransactionType:   constants.TransactionTypeTopup,
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
				}, nil)
			},
		},
		{
			name:               "error",
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody: helpers.Response{
				Message: constants.ErrServerError,
			},
			wantErr: true,
			mockFn: func() {
				mockMdw.EXPECT().MiddlewareValidateToken(gomock.Any()).Do(func(c *gin.Context) {
					tokenData := tokenData

					c.Set("token", tokenData)
					c.Next()
				})

				mockSvc.EXPECT().GetTransaction(gomock.Any(), tokenData.UserID).Return(nil, assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()

			api := gin.New()
			h := &Handler{
				Engine:     api,
				Service:    mockSvc,
				External:   mockExt,
				Middleware: mockMdw,
			}
			h.RegisterRoute()
			w := httptest.NewRecorder()
			endpoint := "/transaction/v1/"

			req, err := http.NewRequest(http.MethodGet, endpoint, nil)
			assert.NoError(t, err)
			req.Header.Set("Authorization", "authorization")

			h.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			if !tt.wantErr {
				res := w.Result()
				defer res.Body.Close()

				response := helpers.Response{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				assert.Equal(t, tt.expectedBody, response)
			}
		})
	}
}

func TestHandler_GetTransactionDetail(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockSvc := NewMockService(ctrlMock)
	mockExt := NewMockExternal(ctrlMock)
	mockMdw := NewMockMiddleware(ctrlMock)

	now := time.Now()

	tokenData := models.TokenData{
		UserID:   1,
		Username: "USERNAME",
		Fullname: "FULLNAME",
		Token:    "TOKEN",
		Email:    "EMAIL",
	}

	tests := []struct {
		name               string
		expectedStatusCode int
		expectedBody       helpers.Response
		wantErr            bool
		mockFn             func()
	}{
		{
			name:               "success",
			expectedStatusCode: http.StatusOK,
			expectedBody: helpers.Response{
				Message: constants.SuccessMessage,
				Data: map[string]interface{}{
					"id":                 float64(1),
					"user_id":            float64(1),
					"amount":             float64(100000),
					"transaction_type":   constants.TransactionTypeTopup,
					"transaction_status": constants.TransactionStatusSuccess,
					"reference":          "REFERENCE",
					"description":        "DESC",
					"additional_info":    "ADDINFO",
					"created_at":         now.Format(time.RFC3339Nano),
					"updated_at":         now.Format(time.RFC3339Nano),
				},
			},
			wantErr: false,
			mockFn: func() {
				mockMdw.EXPECT().MiddlewareValidateToken(gomock.Any()).Do(func(c *gin.Context) {
					tokenData := tokenData

					c.Set("token", tokenData)
					c.Next()
				})

				mockSvc.EXPECT().GetTransactionDetail(gomock.Any(), "REFERENCE").Return(models.Transaction{
					ID:                1,
					UserID:            1,
					Amount:            100000,
					TransactionType:   constants.TransactionTypeTopup,
					TransactionStatus: constants.TransactionStatusSuccess,
					Reference:         "REFERENCE",
					Description:       "DESC",
					AdditionalInfo:    "ADDINFO",
					CreatedAt:         now,
					UpdatedAt:         now,
				}, nil)
			},
		},
		{
			name:               "error",
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody: helpers.Response{
				Message: constants.ErrServerError,
			},
			wantErr: true,
			mockFn: func() {
				mockMdw.EXPECT().MiddlewareValidateToken(gomock.Any()).Do(func(c *gin.Context) {
					tokenData := tokenData

					c.Set("token", tokenData)
					c.Next()
				})

				mockSvc.EXPECT().GetTransactionDetail(gomock.Any(), "REFERENCE").Return(models.Transaction{}, assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()

			api := gin.New()
			h := &Handler{
				Engine:     api,
				Service:    mockSvc,
				External:   mockExt,
				Middleware: mockMdw,
			}
			h.RegisterRoute()
			w := httptest.NewRecorder()

			endpoint := "/transaction/v1/REFERENCE"

			req, err := http.NewRequest(http.MethodGet, endpoint, nil)
			assert.NoError(t, err)
			req.Header.Set("Authorization", "authorization")

			h.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			if !tt.wantErr {
				res := w.Result()
				defer res.Body.Close()

				response := helpers.Response{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				assert.Equal(t, tt.expectedBody, response)
			}
		})
	}
}

func TestHandler_RefundTransaction(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockSvc := NewMockService(ctrlMock)
	mockExt := NewMockExternal(ctrlMock)
	mockMdw := NewMockMiddleware(ctrlMock)

	tokenData := models.TokenData{
		UserID:   1,
		Username: "USERNAME",
		Fullname: "FULLNAME",
		Token:    "TOKEN",
		Email:    "EMAIL",
	}

	req := models.RefundTransaction{
		Reference:      "REFERENCE",
		Description:    "DESC",
		AdditionalInfo: "ADDINFO",
	}

	tests := []struct {
		name               string
		expectedStatusCode int
		expectedBody       helpers.Response
		wantErr            bool
		mockFn             func()
	}{
		{
			name:               "success",
			expectedStatusCode: http.StatusOK,
			expectedBody: helpers.Response{
				Message: constants.SuccessMessage,
				Data: map[string]interface{}{
					"reference":          "REFERENCE",
					"transaction_status": constants.TransactionStatusSuccess,
				},
			},
			wantErr: false,
			mockFn: func() {
				mockMdw.EXPECT().MiddlewareValidateToken(gomock.Any()).Do(func(c *gin.Context) {
					tokenData := tokenData

					c.Set("token", tokenData)
					c.Next()
				})

				mockSvc.EXPECT().RefundTransaction(gomock.Any(), tokenData, &req).Return(models.CreateTransactionResponse{
					Reference:         "REFERENCE",
					TransactionStatus: constants.TransactionStatusSuccess,
				}, nil)
			},
		},
		{
			name:               "error",
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody: helpers.Response{
				Message: constants.ErrServerError,
			},
			wantErr: true,
			mockFn: func() {
				mockMdw.EXPECT().MiddlewareValidateToken(gomock.Any()).Do(func(c *gin.Context) {
					tokenData := tokenData

					c.Set("token", tokenData)
					c.Next()
				})

				mockSvc.EXPECT().RefundTransaction(gomock.Any(), tokenData, &req).Return(models.CreateTransactionResponse{}, assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()

			api := gin.New()
			h := &Handler{
				Engine:     api,
				Service:    mockSvc,
				External:   mockExt,
				Middleware: mockMdw,
			}
			h.RegisterRoute()
			w := httptest.NewRecorder()
			endpoint := "/transaction/v1/refund"
			val, err := json.Marshal(req)
			assert.NoError(t, err)

			body := bytes.NewReader(val)
			req, err := http.NewRequest(http.MethodPost, endpoint, body)
			assert.NoError(t, err)
			req.Header.Set("Authorization", "authorization")

			h.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			if !tt.wantErr {
				res := w.Result()
				defer res.Body.Close()

				response := helpers.Response{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				assert.Equal(t, tt.expectedBody, response)
			}
		})
	}
}
