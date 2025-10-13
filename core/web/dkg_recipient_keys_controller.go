package web

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

const (
	dkgRecipientKeysControllerName = "dkgRecipientKey"
)

// WorkflowKeysController exposes the workflow key
type DKGRecipientKeysController struct {
	App chainlink.Application
}

// Index lists workflow keys
// Example:
// "GET <application>/keys/dkgrecipient"
func (drkc *DKGRecipientKeysController) Index(c *gin.Context) {
	keys, err := drkc.App.GetKeyStore().DKGRecipient().GetAll()
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	jsonAPIResponse(c, presenters.NewDKGRecipientKeyResources(keys), dkgRecipientKeysControllerName)
}
