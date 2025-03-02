package middleware

import (
	models "ewallet-transaction/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func TestExternalDependency_MiddlewareValidateToken(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockExt := NewMockExternal(ctrlMock)

	auth := "Authorization"

	tests := []struct {
		name               string
		wantErr            bool
		mockFn             func()
		expectedStatusCode int
	}{
		{
			name:    "success",
			wantErr: false,
			mockFn: func() {
				mockExt.EXPECT().ValidateToken(gomock.Any(), auth).Return(models.TokenData{
					UserID:   1,
					Username: "username",
					Fullname: "fullname",
					Email:    "email@gmail.com",
				}, nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:    "error",
			wantErr: true,
			mockFn: func() {
				mockExt.EXPECT().ValidateToken(gomock.Any(), auth).Return(models.TokenData{}, assert.AnError)
			},
			expectedStatusCode: http.StatusUnauthorized,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			api := gin.New()

			d := &ExternalDependency{
				External: mockExt,
			}

			w := httptest.NewRecorder()
			endPoint := "/validate-token"
			api.GET(endPoint, d.MiddlewareValidateToken)

			req, err := http.NewRequest(http.MethodGet, endPoint, nil)
			assert.NoError(t, err)
			req.Header.Set("Authorization", auth)

			api.ServeHTTP(w, req)
			assert.Equal(t, tt.expectedStatusCode, w.Code)
		})
	}
}
