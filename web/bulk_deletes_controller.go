package web

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/manyminds/api2go/jsonapi"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/orm"
)

// BulkDeletesController manages background tasks that delete resources given a query
type BulkDeletesController struct {
	App services.Application
}

// Create queues a task to delete runs based on a query
// Example:
//  "<application>/bulk_delete_runs"
func (c *BulkDeletesController) Create(ctx *gin.Context) {
	request := models.BulkDeleteRunRequest{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.AbortWithError(422, err)
	} else if task, err := models.NewBulkDeleteRunTask(request); err != nil {
		ctx.AbortWithError(422, err)
	} else if err := c.App.GetStore().Save(task); err != nil {
		ctx.AbortWithError(500, err)
	} else if doc, err := jsonapi.Marshal(task); err != nil {
		ctx.AbortWithError(500, err)
	} else {
		c.App.WakeBulkRunDeleter()
		ctx.Data(201, MediaType, doc)
	}
}

// Show returns the details of a BulkDeleteTask.
// Example:
//  "<application>/bulk_delete_runs/:RunID"
func (c *BulkDeletesController) Show(ctx *gin.Context) {
	id := ctx.Param("taskID")
	task := models.BulkDeleteRunTask{}

	if err := c.App.GetStore().One("ID", id, &task); err == orm.ErrorNotFound {
		ctx.AbortWithError(404, errors.New("Bulk delete task not found"))
	} else if err != nil {
		ctx.AbortWithError(500, err)
	} else if doc, err := jsonapi.Marshal(&task); err != nil {
		ctx.AbortWithError(500, err)
	} else {
		ctx.Data(200, MediaType, doc)
	}
}
