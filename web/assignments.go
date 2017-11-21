package web

import (
	"github.com/gin-gonic/gin"
)

type Assignments struct{}

func (a *Assignments) Create(c *gin.Context) {
	c.JSON(500, gin.H{
		"errors": []string{"Error saving to database."},
	})
}
