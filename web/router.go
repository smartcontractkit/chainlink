package web

import (
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()
	asgn := Assignments{}
	r.POST("/assignments", asgn.Create)
	return r
}
