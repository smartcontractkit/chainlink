package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

// OffChainReportingKeysController manages OCR key bundles
type OffChainReportingKeysController struct {
	App chainlink.Application
}

// Index lists OCR key bundles
// Example:
// "GET <application>/off-chain-reporting-keys"
func (ocrkbc *OffChainReportingKeysController) Index(c *gin.Context) {
	keys, err := ocrkbc.App.GetStore().OCRKeyStore.FindEncryptedOCRKeyBundles()

	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponse(c, keys, "offChainReportingKeyBundle")
}

// Create and return an OCR key bundle
// Example:
// "POST <application>/off-chain-reporting-keys"
func (ocrkbc *OffChainReportingKeysController) Create(c *gin.Context) {
	request := models.CreateOCRKeysRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	_, encryptedKeyBundle, err := ocrkbc.App.GetStore().OCRKeyStore.GenerateEncryptedOCRKeyBundle(request.Password)

	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponse(c, encryptedKeyBundle, "offChainReportingKeyBundle")
}
