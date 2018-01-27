package web

import (
	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/models"
)

type CustomTaskTypesController struct {
	App services.Application
}

func (ttc *CustomTaskTypesController) Create(c *gin.Context) {
	tt := models.NewCustomTaskType()

	if err := c.ShouldBindJSON(tt); err != nil {
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
	} else if err = ttc.App.GetStore().Save(tt); err != nil {
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
	} else {
		c.JSON(200, gin.H{"id": tt.ID})
	}
}
