package web

import (
	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink-go/logger"
	"github.com/smartcontractkit/chainlink-go/web/controllers"
)

func Router() *gin.Engine {
	engine := gin.New()
	engine.Use(gin.LoggerWithWriter(logger.ForGin()), gin.Recovery())

	j := controllers.JobsController{}
	engine.POST("/jobs", j.Create)
	engine.GET("/jobs/:id", j.Show)

	jr := controllers.JobRunsController{}
	engine.GET("/jobs/:id/runs", jr.Index)

	return engine
}
