package handlers

import (
	"example/common"
	"example/module/user_v1/transport/ginuser"
	core_service "github.com/cesc1802/core-service"
	"github.com/cesc1802/core-service/events"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// PublicHandler place to set up public handler
func PublicHandler(app core_service.Service) func(e *gin.Engine) {
	return func(e *gin.Engine) {
		user := e.Group("users")
		{
			user.GET("", ginuser.ListUser(app))
			user.POST("", ginuser.CreateUser(app))
		}
		pubsub := e.Group("pubsub")
		{
			pubsub.GET("", func(c *gin.Context) {
				ps := app.MustGet(common.KeyPubSub).(events.Stream)

				if err := ps.Publish("test", map[string]interface{}{
					"email": "test@gmail.com",
				}); err != nil {
					log.Println("error =====================================", err)
				}

				c.JSON(http.StatusOK, gin.H{
					"data": true,
				})
			})
		}
	}
}

// PrivateHandler place to set up private handler
func PrivateHandler(app core_service.Service) func(e *gin.Engine) {
	return func(e *gin.Engine) {
		e.Group("/admin")
	}
}
