package web

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
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
	jsonAPIResponse(c, presenters.NewCSAKeyResources(keys), "csaKeys")
}

// Create and return a P2P key
// Example:
// "POST <application>/keys/csa"
func (ctrl *CSAKeysController) Create(c *gin.Context) {
	key, err := ctrl.App.GetKeyStore().CSA().CreateCSAKey()
	if err != nil {
		if errors.Is(err, keystore.ErrCSAKeyExists) {
			jsonAPIError(c, http.StatusBadRequest, err)
			return
		}

		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	jsonAPIResponse(c, presenters.NewCSAKeyResource(*key), "csaKeys")
}
