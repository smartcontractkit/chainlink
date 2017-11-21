package web

import (
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()
	asgn := AssignmentsController{}
	r.POST("/assignments", asgn.Create)
	return r
}
