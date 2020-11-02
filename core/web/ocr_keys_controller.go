package web

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

// OCRKeysController manages OCR key bundles
type OCRKeysController struct {
	App chainlink.Application
}

// Index lists OCR key bundles
// Example:
// "GET <application>/keys/ocr"
func (ocrkbc *OCRKeysController) Index(c *gin.Context) {
	keys, err := ocrkbc.App.GetStore().OCRKeyStore.FindEncryptedOCRKeyBundles()
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	jsonAPIResponse(c, keys, "offChainReportingKeyBundle")
}

// Create and return an OCR key bundle
// Example:
// "POST <application>/keys/ocr"
func (ocrkbc *OCRKeysController) Create(c *gin.Context) {
	_, encryptedKeyBundle, err := ocrkbc.App.GetStore().OCRKeyStore.GenerateEncryptedOCRKeyBundle()
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	jsonAPIResponse(c, encryptedKeyBundle, "offChainReportingKeyBundle")
}

// Delete an OCR key bundle
// Example:
// "DELETE <application>/keys/ocr/:keyID"
// "DELETE <application>/keys/ocr/:keyID?hard=true"
func (ocrkbc *OCRKeysController) Delete(c *gin.Context) {
	var hardDelete bool
	var err error
	if c.Query("hard") != "" {
		hardDelete, err = strconv.ParseBool(c.Query("hard"))
		if err != nil {
			jsonAPIError(c, http.StatusUnprocessableEntity, err)
			return
		}
	}

	id, err := models.Sha256HashFromHex(c.Param("keyID"))
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}
	ekb, err := ocrkbc.App.GetStore().OCRKeyStore.FindEncryptedOCRKeyBundleByID(id)
	if err != nil {
		jsonAPIError(c, http.StatusNotFound, err)
		return
	}
	if hardDelete {
		err = ocrkbc.App.GetStore().OCRKeyStore.DeleteEncryptedOCRKeyBundle(&ekb)
	} else {
		err = ocrkbc.App.GetStore().OCRKeyStore.ArchiveEncryptedOCRKeyBundle(&ekb)
	}
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	jsonAPIResponse(c, ekb, "offChainReportingKeyBundle")
}
