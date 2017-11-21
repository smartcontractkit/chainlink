package web

import (
	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink-go/web/controllers"
)

func Router() *gin.Engine {
	r := gin.Default()
	j := controllers.JobsController{}
	r.POST("/jobs", j.Create)
	return r
}
