package web

import (
	"net/http"

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
	if err := ctx.ShouldBindJSON(&request); err != nil {
		publicError(ctx, http.StatusUnprocessableEntity, err)
	} else if err := c.App.GetStore().KeyStore.Unlock(request.CurrentPassword); err != nil {
		publicError(ctx, http.StatusUnauthorized, err)
	} else if account, err := c.App.GetStore().KeyStore.NewAccount(request.CurrentPassword); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	} else if err := c.App.GetStore().SyncDiskKeyStoreToDB(); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	} else if doc, err := jsonapi.Marshal(&presenters.NewAccount{Account: &account}); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	} else {
		ctx.Data(http.StatusCreated, MediaType, doc)
	}
}
