package web

import (
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

	return engine
}
