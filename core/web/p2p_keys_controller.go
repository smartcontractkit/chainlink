package web

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/store/models/p2pkey"
)

// P2PKeysController manages P2P keys
type P2PKeysController struct {
	App chainlink.Application
}

// Index lists P2P keys
// Example:
// "GET <application>/p2p_keys"
func (p2pkc *P2PKeysController) Index(c *gin.Context) {
	keys, err := p2pkc.App.GetStore().OCRKeyStore.FindEncryptedP2PKeys()
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	jsonAPIResponse(c, keys, "p2pKey")
}

// Create and return a P2P key
// Example:
// "POST <application>/p2p_keys"
func (p2pkc *P2PKeysController) Create(c *gin.Context) {
	_, encryptedP2PKey, err := p2pkc.App.GetStore().OCRKeyStore.GenerateEncryptedP2PKey()
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	jsonAPIResponse(c, encryptedP2PKey, "p2pKey")
}

// Delete a P2P key
// Example:
// "DELETE <application>/p2p_keys/:keyID"
// "DELETE <application>/p2p_keys/:keyID?hard=true"
func (p2pkc *P2PKeysController) Delete(c *gin.Context) {
	var hardDelete bool
	var err error
	if c.Query("hard") != "" {
		hardDelete, err = strconv.ParseBool(c.Query("hard"))
		if err != nil {
			jsonAPIError(c, http.StatusUnprocessableEntity, err)
			return
		}
	}

	ep2pk := p2pkey.EncryptedP2PKey{}
	err = ep2pk.SetID(c.Param("keyID"))
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}
	encryptedP2PKeyPointer, err := p2pkc.App.GetStore().OCRKeyStore.FindEncryptedP2PKeyByID(ep2pk.ID)
	if err != nil {
		jsonAPIError(c, http.StatusNotFound, err)
		return
	}
	if hardDelete {
		err = p2pkc.App.GetStore().OCRKeyStore.DeleteEncryptedP2PKey(encryptedP2PKeyPointer)
	} else {
		err = p2pkc.App.GetStore().OCRKeyStore.ArchiveEncryptedP2PKey(encryptedP2PKeyPointer)
	}
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	jsonAPIResponse(c, encryptedP2PKeyPointer, "offChainReportingKeyBundle")
}
