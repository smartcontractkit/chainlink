package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
)

// OffChainReportingKeysController manages OCR key bundles
type OffChainReportingKeysController struct {
	App chainlink.Application
}

// Index lists OCR key bundles
// Example:
//  "<application>/off-chain-reporting-keys"
func (ocrkbc *OffChainReportingKeysController) Index(c *gin.Context) {
	keys, err := ocrkbc.App.GetStore().OCRKeyStore.FindEncryptedOCRKeyBundles()

	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponse(c, keys, "offChainReportingKeyBundle")
}
