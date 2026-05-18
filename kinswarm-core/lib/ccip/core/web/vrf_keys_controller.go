package web

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/smartcontractkit/chainlink/v2/core/logger/audit"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

// VRFKeysController manages VRF keys
type VRFKeysController struct {
	App chainlink.Application
}

// Index lists VRF keys
// Example:
// "GET <application>/keys/vrf"
func (vrfkc *VRFKeysController) Index(c *gin.Context) {
	keys, err := vrfkc.App.GetKeyStore().VRF().GetAll()
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	jsonAPIResponse(c, presenters.NewVRFKeyResources(keys, vrfkc.App.GetLogger()), "vrfKey")
}

// Create and return a VRF key
// Example:
// "POST <application>/keys/vrf"
func (vrfkc *VRFKeysController) Create(c *gin.Context) {
	ctx := c.Request.Context()
	pk, err := vrfkc.App.GetKeyStore().VRF().Create(ctx)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	vrfkc.App.GetAuditLogger().Audit(audit.KeyCreated, map[string]interface{}{
		"type":                "vrf",
		"id":                  pk.ID(),
		"vrfPublicKey":        pk.PublicKey,
		"vrfPublicKeyAddress": pk.PublicKey.Address(),
	})

	jsonAPIResponse(c, presenters.NewVRFKeyResource(pk, vrfkc.App.GetLogger()), "vrfKey")
}

// Delete a VRF key
// Example:
// "DELETE <application>/keys/vrf/:keyID"
// "DELETE <application>/keys/vrf/:keyID?hard=true"
func (vrfkc *VRFKeysController) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	keyID := c.Param("keyID")
	key, err := vrfkc.App.GetKeyStore().VRF().Get(keyID)
	if err != nil {
		jsonAPIError(c, http.StatusNotFound, err)
		return
	}
	_, err = vrfkc.App.GetKeyStore().VRF().Delete(ctx, keyID)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	vrfkc.App.GetAuditLogger().Audit(audit.KeyDeleted, map[string]interface{}{
		"type": "vrf",
		"id":   keyID,
	})

	jsonAPIResponse(c, presenters.NewVRFKeyResource(key, vrfkc.App.GetLogger()), "vrfKey")
}

// Import imports a VRF key
// Example:
// "Post <application>/keys/vrf/import"
func (vrfkc *VRFKeysController) Import(c *gin.Context) {
	defer vrfkc.App.GetLogger().ErrorIfFn(c.Request.Body.Close, "Error closing Import request body")
	ctx := c.Request.Context()

	bytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}
	oldPassword := c.Query("oldpassword")
	key, err := vrfkc.App.GetKeyStore().VRF().Import(ctx, bytes, oldPassword)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	vrfkc.App.GetAuditLogger().Audit(audit.KeyImported, map[string]interface{}{
		"type":                "vrf",
		"id":                  key.ID(),
		"vrfPublicKey":        key.PublicKey,
		"vrfPublicKeyAddress": key.PublicKey.Address(),
	})

	jsonAPIResponse(c, presenters.NewVRFKeyResource(key, vrfkc.App.GetLogger()), "vrfKey")
}

// Export exports a VRF key
// Example:
// "Post <application>/keys/vrf/export/:keyID"
func (vrfkc *VRFKeysController) Export(c *gin.Context) {
	defer vrfkc.App.GetLogger().ErrorIfFn(c.Request.Body.Close, "Error closing Export request body")

	keyID := c.Param("keyID")
	// New password to re-encrypt the export with
	newPassword := c.Query("newpassword")
	bytes, err := vrfkc.App.GetKeyStore().VRF().Export(keyID, newPassword)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	vrfkc.App.GetAuditLogger().Audit(audit.KeyExported, map[string]interface{}{
		"type": "vrf",
		"id":   keyID,
	})

	c.Data(http.StatusOK, MediaType, bytes)
}
