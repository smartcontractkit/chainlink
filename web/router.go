package web

import (
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()
	asgns := Assignments{}
	r.GET("/assignments", asgns.Index)
	return r
}
