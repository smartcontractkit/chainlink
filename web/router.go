package web

import (
	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink-go/logger"
	storelib "github.com/smartcontractkit/chainlink-go/store"
	"github.com/smartcontractkit/chainlink-go/web/controllers"
)

func Router(store storelib.Store) *gin.Engine {
	engine := gin.New()
	engine.Use(gin.LoggerWithWriter(logger.ForGin()), gin.Recovery())

	j := controllers.JobsController{store}
	engine.POST("/jobs", j.Create)
	engine.GET("/jobs/:id", j.Show)

	jr := controllers.JobRunsController{store}
	engine.GET("/jobs/:id/runs", jr.Index)

	return engine
}
