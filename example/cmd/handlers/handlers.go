package handlers

import (
	"example/common"
	"example/module/user_v1/transport/ginuser"
	core_service "github.com/cesc1802/core-service"
	"github.com/cesc1802/core-service/events"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
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

		mongo := e.Group("mongo")
		{
			mongo.GET("", func(c *gin.Context) {
				mgoDB := app.MustGet(common.KeyMgoDB).(*mongo2.Client)
				type user struct {
					ID    primitive.ObjectID `bson:"_id"`
					Name  string             `bson:"name"`
					Email string             `bson:"email"`
					Age   int                `bson:"age"`
				}

				var u user

				var users []user
				userCollection := mgoDB.Database("example").Collection("users")

				ctx := c.Request.Context()
				insertRes, err := userCollection.Find(ctx, bson.M{})

				for insertRes.Next(ctx) {
					insertRes.Decode(&u)
					users = append(users, u)
				}

				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"data": err.Error()})
					return
				}

				c.JSON(http.StatusOK, gin.H{"data": users})
			})
			mongo.POST("", func(c *gin.Context) {
				mgoDB := app.MustGet(common.KeyMgoDB).(*mongo2.Client)
				type user struct {
					Name  string
					Email string
					Age   int
				}

				var u = user{
					Name:  "thuocnv",
					Email: "thuocnv@gmail.com",
					Age:   28,
				}
				userCollection := mgoDB.Database("example").Collection("users")
				insertRes, err := userCollection.InsertOne(c.Request.Context(), u)

				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"data": err.Error()})
					return
				}

				c.JSON(http.StatusOK, gin.H{"data": insertRes})
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
