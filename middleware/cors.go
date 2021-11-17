package middleware

import (
	"github.com/cesc1802/core-service/config"
	"github.com/cesc1802/core-service/util"
	"github.com/gin-gonic/gin"
)

func Cors(cfg *config.CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", util.SliceStringToString(cfg.AllowOrigins, ","))
		c.Writer.Header().Set("Access-Control-Allow-Credentials", util.BoolToBoolString(cfg.AllowCredentials))
		c.Writer.Header().Set("Access-Control-Allow-Headers", util.SliceStringToString(cfg.AllowHeaders, ","))
		c.Writer.Header().Set("Access-Control-Allow-Methods", util.SliceStringToString(cfg.AllowMethods, ","))

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
