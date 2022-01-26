package web

import (
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
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
	key, err := ocrkc.App.GetKeyStore().OCR().Create()
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	jsonAPIResponse(c, presenters.NewOCRKeysBundleResource(key), "offChainReportingKeyBundle")
}

// Delete an OCR key bundle
// Example:
// "DELETE <application>/keys/ocr/:keyID"
// "DELETE <application>/keys/ocr/:keyID?hard=true"
func (ocrkc *OCRKeysController) Delete(c *gin.Context) {
	id := c.Param("keyID")
	key, err := ocrkc.App.GetKeyStore().OCR().Get(id)
	if err != nil {
		jsonAPIError(c, http.StatusNotFound, err)
		return
	}
	_, err = ocrkc.App.GetKeyStore().OCR().Delete(id)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	jsonAPIResponse(c, presenters.NewOCRKeysBundleResource(key), "offChainReportingKeyBundle")
}

// Import imports an OCR key bundle
// Example:
// "Post <application>/keys/ocr/import"
func (ocrkc *OCRKeysController) Import(c *gin.Context) {
	defer ocrkc.App.GetLogger().ErrorIfClosing(c.Request.Body, "Import request body")

	bytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}
	oldPassword := c.Query("oldpassword")
	encryptedOCRKeyBundle, err := ocrkc.App.GetKeyStore().OCR().Import(bytes, oldPassword)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponse(c, encryptedOCRKeyBundle, "offChainReportingKeyBundle")
}

// Export exports an OCR key bundle
// Example:
// "Post <application>/keys/ocr/export"
func (ocrkc *OCRKeysController) Export(c *gin.Context) {
	defer ocrkc.App.GetLogger().ErrorIfClosing(c.Request.Body, "Export response body")

	stringID := c.Param("ID")
	newPassword := c.Query("newpassword")
	bytes, err := ocrkc.App.GetKeyStore().OCR().Export(stringID, newPassword)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	c.Data(http.StatusOK, MediaType, bytes)
}
