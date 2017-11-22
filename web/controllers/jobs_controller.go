package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink-go/models"
	"github.com/smartcontractkit/chainlink-go/orm"
)

type JobsController struct{}

func (jc *JobsController) Create(c *gin.Context) {
	db := orm.GetDB()
	j := models.NewJob()

	if err := c.ShouldBindJSON(&j); err != nil {
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
	} else if _, err = j.Valid(); err != nil {
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
	} else if err = db.Save(&j); err != nil {
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
	} else {
		c.JSON(200, gin.H{"id": j.ID})
	}
}
