// Code generated by MockGen. DO NOT EDIT.
// Source: service.go
//
// Generated by this command:
//
//	mockgen -source=service.go -destination=service_mock_test.go -package=transaction
//

// Package transaction is a generated GoMock package.
package transaction

import (
	context "context"
	external "ewallet-transaction/external"
	models "ewallet-transaction/internal/models"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// Mockrepository is a mock of repository interface.
type Mockrepository struct {
	ctrl     *gomock.Controller
	recorder *MockrepositoryMockRecorder
	isgomock struct{}
}

// MockrepositoryMockRecorder is the mock recorder for Mockrepository.
type MockrepositoryMockRecorder struct {
	mock *Mockrepository
}

// NewMockrepository creates a new mock instance.
func NewMockrepository(ctrl *gomock.Controller) *Mockrepository {
	mock := &Mockrepository{ctrl: ctrl}
	mock.recorder = &MockrepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockrepository) EXPECT() *MockrepositoryMockRecorder {
	return m.recorder
}

// CreateTransaction mocks base method.
func (m *Mockrepository) CreateTransaction(ctx context.Context, trx *models.Transaction) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTransaction", ctx, trx)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateTransaction indicates an expected call of CreateTransaction.
func (mr *MockrepositoryMockRecorder) CreateTransaction(ctx, trx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTransaction", reflect.TypeOf((*Mockrepository)(nil).CreateTransaction), ctx, trx)
}

// GetTransaction mocks base method.
func (m *Mockrepository) GetTransaction(ctx context.Context, userID uint64) ([]models.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransaction", ctx, userID)
	ret0, _ := ret[0].([]models.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTransaction indicates an expected call of GetTransaction.
func (mr *MockrepositoryMockRecorder) GetTransaction(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransaction", reflect.TypeOf((*Mockrepository)(nil).GetTransaction), ctx, userID)
}

// GetTransactionByReference mocks base method.
func (m *Mockrepository) GetTransactionByReference(arg0 context.Context, arg1 string, arg2 bool) (models.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransactionByReference", arg0, arg1, arg2)
	ret0, _ := ret[0].(models.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTransactionByReference indicates an expected call of GetTransactionByReference.
func (mr *MockrepositoryMockRecorder) GetTransactionByReference(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransactionByReference", reflect.TypeOf((*Mockrepository)(nil).GetTransactionByReference), arg0, arg1, arg2)
}

// UpdateStatusTransaction mocks base method.
func (m *Mockrepository) UpdateStatusTransaction(ctx context.Context, reference, status, additionalInfo string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateStatusTransaction", ctx, reference, status, additionalInfo)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateStatusTransaction indicates an expected call of UpdateStatusTransaction.
func (mr *MockrepositoryMockRecorder) UpdateStatusTransaction(ctx, reference, status, additionalInfo any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateStatusTransaction", reflect.TypeOf((*Mockrepository)(nil).UpdateStatusTransaction), ctx, reference, status, additionalInfo)
}

// MockIExternal is a mock of IExternal interface.
type MockIExternal struct {
	ctrl     *gomock.Controller
	recorder *MockIExternalMockRecorder
	isgomock struct{}
}

// MockIExternalMockRecorder is the mock recorder for MockIExternal.
type MockIExternalMockRecorder struct {
	mock *MockIExternal
}

// NewMockIExternal creates a new mock instance.
func NewMockIExternal(ctrl *gomock.Controller) *MockIExternal {
	mock := &MockIExternal{ctrl: ctrl}
	mock.recorder = &MockIExternalMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIExternal) EXPECT() *MockIExternalMockRecorder {
	return m.recorder
}

// CreditBalance mocks base method.
func (m *MockIExternal) CreditBalance(ctx context.Context, token string, req external.UpdateBalance) (*external.UpdateBalanceResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreditBalance", ctx, token, req)
	ret0, _ := ret[0].(*external.UpdateBalanceResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreditBalance indicates an expected call of CreditBalance.
func (mr *MockIExternalMockRecorder) CreditBalance(ctx, token, req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreditBalance", reflect.TypeOf((*MockIExternal)(nil).CreditBalance), ctx, token, req)
}

// DebitBalance mocks base method.
func (m *MockIExternal) DebitBalance(ctx context.Context, token string, req external.UpdateBalance) (*external.UpdateBalanceResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DebitBalance", ctx, token, req)
	ret0, _ := ret[0].(*external.UpdateBalanceResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DebitBalance indicates an expected call of DebitBalance.
func (mr *MockIExternalMockRecorder) DebitBalance(ctx, token, req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DebitBalance", reflect.TypeOf((*MockIExternal)(nil).DebitBalance), ctx, token, req)
}

// SendNotification mocks base method.
func (m *MockIExternal) SendNotification(ctx context.Context, recipient, templateName string, placeHolder map[string]string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendNotification", ctx, recipient, templateName, placeHolder)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendNotification indicates an expected call of SendNotification.
func (mr *MockIExternalMockRecorder) SendNotification(ctx, recipient, templateName, placeHolder any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendNotification", reflect.TypeOf((*MockIExternal)(nil).SendNotification), ctx, recipient, templateName, placeHolder)
}

// ValidateToken mocks base method.
func (m *MockIExternal) ValidateToken(ctx context.Context, token string) (models.TokenData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateToken", ctx, token)
	ret0, _ := ret[0].(models.TokenData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateToken indicates an expected call of ValidateToken.
func (mr *MockIExternalMockRecorder) ValidateToken(ctx, token any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateToken", reflect.TypeOf((*MockIExternal)(nil).ValidateToken), ctx, token)
}
