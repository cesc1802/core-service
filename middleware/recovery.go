package middleware

import (
	"github.com/cesc1802/core-service/constant"
	"github.com/cesc1802/core-service/i18n"
	"github.com/cesc1802/core-service/model/app_error"
	"github.com/cesc1802/core-service/util"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func Recovery(i18n *i18n.I18n) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.Header("Content-Type", "application/json")

				if ve, ok := err.(validator.ValidationErrors); ok {
					appVE := util.HandleValidationErrors(c.GetHeader(constant.HeaderAcceptLanguage), i18n, ve)
					c.AbortWithStatusJSON(appVE.StatusCode, appVE)
					return
				}

				appErr := app_error.ErrInternal(err.(error))
				c.AbortWithStatusJSON(appErr.StatusCode, appErr)
				panic(err)
				return
			}
		}()

		c.Next()
	}
}
