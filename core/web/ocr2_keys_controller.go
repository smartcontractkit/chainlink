package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

// OCRKeysController manages OCR key bundles
type OCR2KeysController struct {
	App chainlink.Application
}

// Index lists OCR key bundles
// Example:
// "GET <application>/keys/ocr"
func (ocr2kc *OCR2KeysController) Index(c *gin.Context) {
	ekbs, err := ocr2kc.App.GetKeyStore().OCR2().GetAll()
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	jsonAPIResponse(c, presenters.NewOCR2KeysBundleResources(ekbs), "offChainReportingKeyBundle")
}

// Create and return an OCR key bundle
// Example:
// "POST <application>/keys/ocr"
func (ocr2kc *OCR2KeysController) Create(c *gin.Context) {
	key, err := ocr2kc.App.GetKeyStore().OCR2().Create()
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	jsonAPIResponse(c, presenters.NewOCR2KeysBundleResource(key), "offChainReportingKeyBundle")
}

// Delete an OCR key bundle
// Example:
// "DELETE <application>/keys/ocr/:keyID"
func (ocr2kc *OCR2KeysController) Delete(c *gin.Context) {
	id := c.Param("keyID")
	key, err := ocr2kc.App.GetKeyStore().OCR2().Get(id)
	if err != nil {
		jsonAPIError(c, http.StatusNotFound, err)
		return
	}
	err = ocr2kc.App.GetKeyStore().OCR2().Delete(id)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	jsonAPIResponse(c, presenters.NewOCR2KeysBundleResource(key), "offChainReportingKeyBundle")
}
