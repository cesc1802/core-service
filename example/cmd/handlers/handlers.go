package handlers

import (
	"example/module/user_v1/transport/ginuser"
	core_service "github.com/cesc1802/core-service"
	"github.com/gin-gonic/gin"
)

// PublicHandler place to set up public handler
func PublicHandler(app core_service.Service) func(e *gin.Engine) {
	return func(e *gin.Engine) {
		user := e.Group("users")
		{
			user.GET("", ginuser.ListUser(app))
			user.POST("", ginuser.CreateUser(app))
		}
	}
}

// PrivateHandler place to set up private handler
func PrivateHandler(app core_service.Service) func(e *gin.Engine) {
	return func(e *gin.Engine) {
		e.Group("/admin")
	}
}
