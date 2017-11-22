package web

import (
	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink-go/web/controllers"
)

func Router() *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())

	j := controllers.JobsController{}
	engine.POST("/jobs", j.Create)
	engine.GET("/jobs/:id", j.Show)

	return engine
}
