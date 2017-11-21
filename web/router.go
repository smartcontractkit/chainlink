package web

import (
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()
	j := JobsController{}
	r.POST("/jobs", j.Create)
	return r
}
