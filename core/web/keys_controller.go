package web

import (
	"net/http"

	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/presenters"

	"github.com/gin-gonic/gin"
)

// KeysController manages account keys
type KeysController struct {
	App chainlink.Application
}

// Create adds a new account
// Example:
//  "<application>/keys"
func (kc *KeysController) Create(c *gin.Context) {
	request := models.CreateKeyRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}
	if err := kc.App.GetStore().KeyStore.Unlock(request.CurrentPassword); err != nil {
		jsonAPIError(c, http.StatusUnauthorized, err)
		return
	}

	account, err := kc.App.GetStore().KeyStore.NewAccount(request.CurrentPassword)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	if err := kc.App.GetStore().SyncDiskKeyStoreToDB(); err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponseWithStatus(c, presenters.NewAccount{Account: &account}, "account", http.StatusCreated)
}
