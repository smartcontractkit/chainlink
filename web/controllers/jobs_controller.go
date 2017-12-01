package controllers

import (
	"github.com/asdine/storm"
	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink-go/models"
	"github.com/smartcontractkit/chainlink-go/store"
)

type JobsController struct {
	Store store.Store
}

func (self *JobsController) Create(c *gin.Context) {
	j := models.NewJob()

	if err := c.ShouldBindJSON(&j); err != nil {
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
	} else if err = self.Store.Save(&j); err != nil {
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
	} else {
		c.JSON(200, gin.H{"id": j.ID})
	}
}

func (self *JobsController) Show(c *gin.Context) {
	id := c.Param("id")
	var j models.Job
	err := self.Store.One("ID", id, &j)

	if err == storm.ErrNotFound {
		c.JSON(404, gin.H{
			"errors": []string{"Job not found."},
		})
	} else if err != nil {
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
	} else {
		c.JSON(200, j)
	}
}
