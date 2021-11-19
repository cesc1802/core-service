package ginuser

import (
	goservice "github.com/cesc1802/core-service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func ListUser(sc goservice.Service) gin.HandlerFunc {
	return func(c *gin.Context) {

		db := sc.MustGet("portal").(*gorm.DB)
		var user []user

		db.Table("users").Find(&user)
		c.JSON(http.StatusOK, gin.H{
			"data": user,
		})
	}
}
