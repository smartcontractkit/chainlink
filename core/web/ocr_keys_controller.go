package web

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/smartcontractkit/chainlink/v2/core/logger/audit"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

// OCRKeysController manages OCR key bundles
type OCRKeysController struct {
	App chainlink.Application
}

// Index lists OCR key bundles
// Example:
// "GET <application>/keys/ocr"
func (ocrkc *OCRKeysController) Index(c *gin.Context) {
	ekbs, err := ocrkc.App.GetKeyStore().OCR().GetAll()
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	jsonAPIResponse(c, presenters.NewOCRKeysBundleResources(ekbs), "offChainReportingKeyBundle")
}

// Create and return an OCR key bundle
// Example:
// "POST <application>/keys/ocr"
func (ocrkc *OCRKeysController) Create(c *gin.Context) {
	ctx := c.Request.Context()
	key, err := ocrkc.App.GetKeyStore().OCR().Create(ctx)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	ocrkc.App.GetAuditLogger().Audit(audit.OCRKeyBundleCreated, map[string]interface{}{
		"ocrKeyBundleID":                      key.ID(),
		"ocrKeyBundlePublicKeyAddressOnChain": key.PublicKeyAddressOnChain(),
	})
	jsonAPIResponse(c, presenters.NewOCRKeysBundleResource(key), "offChainReportingKeyBundle")
}

// Delete an OCR key bundle
// Example:
// "DELETE <application>/keys/ocr/:keyID"
// "DELETE <application>/keys/ocr/:keyID?hard=true"
func (ocrkc *OCRKeysController) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("keyID")
	key, err := ocrkc.App.GetKeyStore().OCR().Get(id)
	if err != nil {
		jsonAPIError(c, http.StatusNotFound, err)
		return
	}
	_, err = ocrkc.App.GetKeyStore().OCR().Delete(ctx, id)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	ocrkc.App.GetAuditLogger().Audit(audit.OCRKeyBundleDeleted, map[string]interface{}{"id": id})
	jsonAPIResponse(c, presenters.NewOCRKeysBundleResource(key), "offChainReportingKeyBundle")
}

// Import imports an OCR key bundle
// Example:
// "Post <application>/keys/ocr/import"
func (ocrkc *OCRKeysController) Import(c *gin.Context) {
	defer ocrkc.App.GetLogger().ErrorIfFn(c.Request.Body.Close, "Error closing Import request body")
	ctx := c.Request.Context()

	bytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}
	oldPassword := c.Query("oldpassword")
	encryptedOCRKeyBundle, err := ocrkc.App.GetKeyStore().OCR().Import(ctx, bytes, oldPassword)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	ocrkc.App.GetAuditLogger().Audit(audit.OCRKeyBundleImported, map[string]interface{}{
		"OCRID":                      encryptedOCRKeyBundle.GetID(),
		"OCRPublicKeyAddressOnChain": encryptedOCRKeyBundle.PublicKeyAddressOnChain(),
		"OCRPublicKeyOffChain":       encryptedOCRKeyBundle.PublicKeyOffChain(),
	})

	jsonAPIResponse(c, encryptedOCRKeyBundle, "offChainReportingKeyBundle")
}

// Export exports an OCR key bundle
// Example:
// "Post <application>/keys/ocr/export"
func (ocrkc *OCRKeysController) Export(c *gin.Context) {
	defer ocrkc.App.GetLogger().ErrorIfFn(c.Request.Body.Close, "Error closing Export response body")

	stringID := c.Param("ID")
	newPassword := c.Query("newpassword")
	bytes, err := ocrkc.App.GetKeyStore().OCR().Export(stringID, newPassword)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	ocrkc.App.GetAuditLogger().Audit(audit.OCRKeyBundleExported, map[string]interface{}{"keyID": stringID})
	c.Data(http.StatusOK, MediaType, bytes)
}
