package web

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/smartcontractkit/chainlink/v2/core/logger/audit"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

// P2PKeysController manages P2P keys
type P2PKeysController struct {
	App chainlink.Application
}

// Index lists P2P keys
// Example:
// "GET <application>/keys/p2p"
func (p2pkc *P2PKeysController) Index(c *gin.Context) {
	keys, err := p2pkc.App.GetKeyStore().P2P().GetAll()
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	jsonAPIResponse(c, presenters.NewP2PKeyResources(keys), "p2pKey")
}

const keyType = "Ed25519"

// Create and return a P2P key
// Example:
// "POST <application>/keys/p2p"
func (p2pkc *P2PKeysController) Create(c *gin.Context) {
	ctx := c.Request.Context()
	key, err := p2pkc.App.GetKeyStore().P2P().Create(ctx)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	p2pkc.App.GetAuditLogger().Audit(audit.KeyCreated, map[string]interface{}{
		"type":         "p2p",
		"id":           key.ID(),
		"p2pPublicKey": key.PublicKeyHex(),
		"p2pPeerID":    key.PeerID(),
		"p2pType":      keyType,
	})
	jsonAPIResponse(c, presenters.NewP2PKeyResource(key), "p2pKey")
}

// Delete a P2P key
// Example:
// "DELETE <application>/keys/p2p/:keyID"
// "DELETE <application>/keys/p2p/:keyID?hard=true"
func (p2pkc *P2PKeysController) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	keyID, err := p2pkey.MakePeerID(c.Param("keyID"))
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}
	key, err := p2pkc.App.GetKeyStore().P2P().Get(keyID)
	if err != nil {
		jsonAPIError(c, http.StatusNotFound, err)
		return
	}
	_, err = p2pkc.App.GetKeyStore().P2P().Delete(ctx, key.PeerID())
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	p2pkc.App.GetAuditLogger().Audit(audit.KeyDeleted, map[string]interface{}{
		"type": "p2p",
		"id":   keyID,
	})

	jsonAPIResponse(c, presenters.NewP2PKeyResource(key), "p2pKey")
}

// Import imports a P2P key
// Example:
// "Post <application>/keys/p2p/import"
func (p2pkc *P2PKeysController) Import(c *gin.Context) {
	defer p2pkc.App.GetLogger().ErrorIfFn(c.Request.Body.Close, "Error closing Import request body")
	ctx := c.Request.Context()

	bytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}
	oldPassword := c.Query("oldpassword")
	key, err := p2pkc.App.GetKeyStore().P2P().Import(ctx, bytes, oldPassword)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	p2pkc.App.GetAuditLogger().Audit(audit.KeyImported, map[string]interface{}{
		"type":         "p2p",
		"id":           key.ID(),
		"p2pPublicKey": key.PublicKeyHex(),
		"p2pPeerID":    key.PeerID(),
		"p2pType":      keyType,
	})

	jsonAPIResponse(c, presenters.NewP2PKeyResource(key), "p2pKey")
}

// Export exports a P2P key
// Example:
// "Post <application>/keys/p2p/export"
func (p2pkc *P2PKeysController) Export(c *gin.Context) {
	defer p2pkc.App.GetLogger().ErrorIfFn(c.Request.Body.Close, "Error closing Export request body")

	keyID, err := p2pkey.MakePeerID(c.Param("ID"))
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	newPassword := c.Query("newpassword")
	bytes, err := p2pkc.App.GetKeyStore().P2P().Export(keyID, newPassword)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	p2pkc.App.GetAuditLogger().Audit(audit.KeyExported, map[string]interface{}{
		"type": "p2p",
		"id":   keyID,
	})

	c.Data(http.StatusOK, MediaType, bytes)
}
