package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/presenters"
)

// KeysController manages account keys
type KeysController struct {
	App services.Application
}

// Create adds a new account
// Example:
//  "<application>/keys"
func (kc *KeysController) Create(c *gin.Context) {
	request := models.CreateKeyRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
	} else if err := kc.App.GetStore().KeyStore.Unlock(request.CurrentPassword); err != nil {
		jsonAPIError(c, http.StatusUnauthorized, err)
	} else if account, err := kc.App.GetStore().KeyStore.NewAccount(request.CurrentPassword); err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
	} else if err := kc.App.GetStore().SyncDiskKeyStoreToDB(); err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
	} else {
		jsonAPIResponseWithStatus(c, presenters.NewAccount{Account: &account}, "account", http.StatusCreated)
	}
}
