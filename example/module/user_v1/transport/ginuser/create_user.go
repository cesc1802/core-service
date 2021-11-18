package ginuser

import (
	core_service "github.com/cesc1802/core-service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type user struct {
	Username string `json:"username" gorm:"column:username" binding:"required"`
	Password string `json:"password" gorm:"column:password" binding:"required"`
}

func CreateUser(sc core_service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		db := sc.MustGet("portal").(*gorm.DB)
		var u user

		if err := c.ShouldBind(&u); err != nil {
			panic(err)
		}

		if err := db.Table("users").Create(&u).Error; err != nil {
			panic(err)
		}
		c.JSON(http.StatusOK, gin.H{
			"data": true,
		})
	}
}
