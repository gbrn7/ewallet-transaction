package middleware

import (
	"ewallet-transaction/helpers"
	"ewallet-transaction/internal/handler/transaction"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ExternalDependency struct {
	External transaction.External
}

func (d *ExternalDependency) MiddlewareValidateToken(c *gin.Context) {
	var (
		log = helpers.Logger
	)
	auth := c.Request.Header.Get("Authorization")
	if auth == "" {
		log.Println("authorization empty")
		helpers.SendResponseHTTP(c, http.StatusUnauthorized, "unauthorized", nil)
		c.Abort()
		return
	}

	tokenData, err := d.External.ValidateToken(c.Request.Context(), auth)
	if err != nil {
		log.Error(err)
		helpers.SendResponseHTTP(c, http.StatusUnauthorized, "unauthorized", nil)
		c.Abort()

	}

	tokenData.Token = auth

	c.Set("token", tokenData)

	c.Next()
}
