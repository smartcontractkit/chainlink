package web

import (
	"github.com/gin-gonic/gin"
	"github.com/manyminds/api2go/jsonapi"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/presenters"
)

// KeysController manages account keys
type KeysController struct {
	App services.Application
}

// Create adds a new account
// Example:
//  "<application>/keys"
func (c *KeysController) Create(ctx *gin.Context) {
	request := models.CreateKeyRequest{}
	// TODO: Change CreateKeyRequest to only have one password
	// or validate that they are the same.

	if err := ctx.ShouldBindJSON(&request); err != nil {
		publicError(ctx, 422, err)
	} else if err := c.App.GetStore().KeyStore.Unlock(request.CurrentPassword); err != nil {
		publicError(ctx, 401, err)
	} else if account, err := c.App.GetStore().KeyStore.NewAccount(request.NewAccountPassword); err != nil {
		ctx.AbortWithError(500, err)
	} else if err := c.App.GetStore().SyncDiskKeyStoreToDB(); err != nil {
		ctx.AbortWithError(500, err)
	} else if doc, err := jsonapi.Marshal(&presenters.NewAccount{&account}); err != nil {
		ctx.AbortWithError(500, err)
	} else {
		ctx.Data(201, MediaType, doc)
	}
}
