package middleware

import (
	"github.com/cesc1802/core-service/config"
	"github.com/cesc1802/core-service/util"
	"github.com/gin-gonic/gin"
)

func Cors(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", util.SliceStringToString(cfg.CORS.AllowOrigins, ","))
		c.Writer.Header().Set("Access-Control-Allow-Credentials", util.BoolToBoolString(cfg.CORS.AllowCredentials))
		//"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With"
		c.Writer.Header().Set("Access-Control-Allow-Headers", util.SliceStringToString(cfg.CORS.AllowHeaders, ","))
		//"POST, OPTIONS, GET, PUT, DELETE"
		c.Writer.Header().Set("Access-Control-Allow-Methods", util.SliceStringToString(cfg.CORS.AllowMethods, ","))

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
