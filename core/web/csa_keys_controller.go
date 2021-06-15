package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

// CSAKeysController manages CSA keys
type CSAKeysController struct {
	App chainlink.Application
}

// Index lists P2P keys
// Example:
// "GET <application>/keys/csa"
func (ctrl *CSAKeysController) Index(c *gin.Context) {
	keys, err := ctrl.App.GetKeyStore().CSA().ListCSAKeys()
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	jsonAPIResponse(c, presenters.NewResources(keys), "csaKeys")
}

// Create and return a P2P key
// Example:
// "POST <application>/keys/csa"
func (ctrl *CSAKeysController) Create(c *gin.Context) {
	key, err := ctrl.App.GetKeyStore().CSA().CreateCSAKey()
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	jsonAPIResponse(c, presenters.NewResource(*key), "csaKeys")
}
