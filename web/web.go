package web

import (
	"strconv"
	"system-service-template/database"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var engine *gin.Engine

func New() *gin.Engine {
	engine = gin.New()

	engine.GET("/", func(ctx *gin.Context) {
		log.Println("http request")
		ctx.String(200, "ok")
	})

	engine.GET("/user/add", func(ctx *gin.Context) {
		age, _ := strconv.Atoi(ctx.Query("age"))
		name := ctx.Query("name")
		phone := ctx.Query("phone")
		user := &database.User{
			Age:   age,
			Name:  name,
			Phone: phone,
		}
		err := database.ModelUserAdd(user)
		if err != nil {
			ctx.JSON(500, err.Error())
		} else {
			ctx.JSON(200, user)
		}
	})

	engine.GET("/user/list", func(ctx *gin.Context) {
		users, err := database.ModelUsersQuery()
		if err != nil {
			ctx.JSON(500, err.Error())
		} else {
			ctx.JSON(200, users)
		}
	})

	return engine
}
