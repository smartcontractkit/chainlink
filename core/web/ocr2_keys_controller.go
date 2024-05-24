package web

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/logger/audit"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

// OCRKeysController manages OCR key bundles
type OCR2KeysController struct {
	App chainlink.Application
}

// Index lists OCR2 key bundles
// Example:
// "GET <application>/keys/ocr"
func (ocr2kc *OCR2KeysController) Index(c *gin.Context) {
	ekbs, err := ocr2kc.App.GetKeyStore().OCR2().GetAll()
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	jsonAPIResponse(c, presenters.NewOCR2KeysBundleResources(ekbs), "offChainReporting2KeyBundle")
}

// Create and return an OCR2 key bundle
// Example:
// "POST <application>/keys/ocr"
func (ocr2kc *OCR2KeysController) Create(c *gin.Context) {
	ctx := c.Request.Context()
	chainType := chaintype.ChainType(c.Param("chainType"))
	key, err := ocr2kc.App.GetKeyStore().OCR2().Create(ctx, chainType)
	if errors.Is(errors.Cause(err), chaintype.ErrInvalidChainType) {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	ocr2kc.App.GetAuditLogger().Audit(audit.OCR2KeyBundleCreated, map[string]interface{}{
		"ocr2KeyID":                        key.ID(),
		"ocr2KeyChainType":                 key.ChainType(),
		"ocr2KeyConfigEncryptionPublicKey": key.ConfigEncryptionPublicKey(),
		"ocr2KeyOffchainPublicKey":         key.OffchainPublicKey(),
		"ocr2KeyMaxSignatureLength":        key.MaxSignatureLength(),
		"ocr2KeyPublicKey":                 key.PublicKey(),
	})
	jsonAPIResponse(c, presenters.NewOCR2KeysBundleResource(key), "offChainReporting2KeyBundle")
}

// Delete an OCR2 key bundle
// Example:
// "DELETE <application>/keys/ocr/:keyID"
func (ocr2kc *OCR2KeysController) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("keyID")
	key, err := ocr2kc.App.GetKeyStore().OCR2().Get(id)
	if err != nil {
		jsonAPIError(c, http.StatusNotFound, err)
		return
	}
	err = ocr2kc.App.GetKeyStore().OCR2().Delete(ctx, id)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	ocr2kc.App.GetAuditLogger().Audit(audit.OCR2KeyBundleDeleted, map[string]interface{}{"id": id})
	jsonAPIResponse(c, presenters.NewOCR2KeysBundleResource(key), "offChainReporting2KeyBundle")
}

// Import imports an OCR2 key bundle
// Example:
// "Post <application>/keys/ocr/import"
func (ocr2kc *OCR2KeysController) Import(c *gin.Context) {
	defer ocr2kc.App.GetLogger().ErrorIfFn(c.Request.Body.Close, "Error closing Import request body")
	ctx := c.Request.Context()

	bytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}
	oldPassword := c.Query("oldpassword")
	keyBundle, err := ocr2kc.App.GetKeyStore().OCR2().Import(ctx, bytes, oldPassword)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	ocr2kc.App.GetAuditLogger().Audit(audit.OCR2KeyBundleImported, map[string]interface{}{
		"ocr2KeyID":                        keyBundle.ID(),
		"ocr2KeyChainType":                 keyBundle.ChainType(),
		"ocr2KeyConfigEncryptionPublicKey": keyBundle.ConfigEncryptionPublicKey(),
		"ocr2KeyOffchainPublicKey":         keyBundle.OffchainPublicKey(),
		"ocr2KeyMaxSignatureLength":        keyBundle.MaxSignatureLength(),
		"ocr2KeyPublicKey":                 keyBundle.PublicKey(),
	})

	jsonAPIResponse(c, presenters.NewOCR2KeysBundleResource(keyBundle), "offChainReporting2KeyBundle")
}

// Export exports an OCR2 key bundle
// Example:
// "Post <application>/keys/ocr/export"
func (ocr2kc *OCR2KeysController) Export(c *gin.Context) {
	defer ocr2kc.App.GetLogger().ErrorIfFn(c.Request.Body.Close, "Error closing Export response body")

	stringID := c.Param("ID")
	newPassword := c.Query("newpassword")
	bytes, err := ocr2kc.App.GetKeyStore().OCR2().Export(stringID, newPassword)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	ocr2kc.App.GetAuditLogger().Audit(audit.OCR2KeyBundleExported, map[string]interface{}{"keyID": stringID})
	c.Data(http.StatusOK, MediaType, bytes)
}
