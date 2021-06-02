package web

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

// P2PKeysController manages P2P keys
type P2PKeysController struct {
	App chainlink.Application
}

// Index lists P2P keys
// Example:
// "GET <application>/keys/p2p"
func (p2pkc *P2PKeysController) Index(c *gin.Context) {
	keys, err := p2pkc.App.GetKeyStore().OCR.FindEncryptedP2PKeys()
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	jsonAPIResponse(c, presenters.NewP2PKeyResources(keys), "p2pKey")
}

// Create and return a P2P key
// Example:
// "POST <application>/keys/p2p"
func (p2pkc *P2PKeysController) Create(c *gin.Context) {
	_, key, err := p2pkc.App.GetKeyStore().OCR.GenerateEncryptedP2PKey()
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	jsonAPIResponse(c, presenters.NewP2PKeyResource(key), "p2pKey")
}

// Delete a P2P key
// Example:
// "DELETE <application>/keys/p2p/:keyID"
// "DELETE <application>/keys/p2p/:keyID?hard=true"
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
	key, err := p2pkc.App.GetKeyStore().OCR.FindEncryptedP2PKeyByID(ep2pk.ID)
	if err != nil {
		jsonAPIError(c, http.StatusNotFound, err)
		return
	}
	if hardDelete {
		err = p2pkc.App.GetKeyStore().OCR.DeleteEncryptedP2PKey(key)
	} else {
		err = p2pkc.App.GetKeyStore().OCR.ArchiveEncryptedP2PKey(key)
	}
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	jsonAPIResponse(c, presenters.NewP2PKeyResource(*key), "p2pKey")
}

// Import imports a P2P key
// Example:
// "Post <application>/keys/p2p/import"
func (p2pkc *P2PKeysController) Import(c *gin.Context) {
	defer logger.ErrorIfCalling(c.Request.Body.Close)

	bytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}
	oldPassword := c.Query("oldpassword")
	key, err := p2pkc.App.GetKeyStore().OCR.ImportP2PKey(bytes, oldPassword)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponse(c, presenters.NewP2PKeyResource(*key), "p2pKey")
}

// Export exports a P2P key
// Example:
// "Post <application>/keys/p2p/export"
func (p2pkc *P2PKeysController) Export(c *gin.Context) {
	defer logger.ErrorIfCalling(c.Request.Body.Close)

	stringID := c.Param("ID")
	id64, err := strconv.ParseInt(stringID, 10, 32)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, errors.New("invalid key ID"))
		return
	}
	id := int32(id64)
	newPassword := c.Query("newpassword")
	bytes, err := p2pkc.App.GetKeyStore().OCR.ExportP2PKey(id, newPassword)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	c.Data(http.StatusOK, MediaType, bytes)
}
