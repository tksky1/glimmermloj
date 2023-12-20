package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"glimmermloj/controller"
)

func Init() *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())
	r.MaxMultipartMemory = 101 << 20
	r.POST("/submit", controller.Submit)
	r.GET("/rank", controller.GetRank)
	return r
}
