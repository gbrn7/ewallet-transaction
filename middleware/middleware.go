package middleware

import (
	"ewallet-transaction/helpers"
	"ewallet-transaction/internal/handler/transaction"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ExternalDependency struct {
	External transaction.External
}

func (d *ExternalDependency) MiddlewareValidateToken(c *gin.Context) {

	auth := c.Request.Header.Get("Authorization")
	if auth == "" {
		fmt.Println("authorization empty")
		helpers.SendResponseHTTP(c, http.StatusUnauthorized, "unauthorized", nil)
		c.Abort()
		return
	}

	tokenData, err := d.External.ValidateToken(c.Request.Context(), auth)
	if err != nil {
		fmt.Println(err)
		helpers.SendResponseHTTP(c, http.StatusUnauthorized, "unauthorized", nil)
		c.Abort()

	}

	tokenData.Token = auth

	c.Set("token", tokenData)

	c.Next()
}
