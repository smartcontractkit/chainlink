package web

import (
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
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
func (ocrkc *OCRKeysController) Index(c *gin.Context) {
	keys, err := ocrkc.App.GetStore().OCRKeyStore.FindEncryptedOCRKeyBundles()
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	jsonAPIResponse(c, keys, "offChainReportingKeyBundle")
}

// Create and return an OCR key bundle
// Example:
// "POST <application>/keys/ocr"
func (ocrkc *OCRKeysController) Create(c *gin.Context) {
	_, encryptedKeyBundle, err := ocrkc.App.GetStore().OCRKeyStore.GenerateEncryptedOCRKeyBundle()
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
func (ocrkc *OCRKeysController) Delete(c *gin.Context) {
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
	ekb, err := ocrkc.App.GetStore().OCRKeyStore.FindEncryptedOCRKeyBundleByID(id)
	if err != nil {
		jsonAPIError(c, http.StatusNotFound, err)
		return
	}
	if hardDelete {
		err = ocrkc.App.GetStore().OCRKeyStore.DeleteEncryptedOCRKeyBundle(&ekb)
	} else {
		err = ocrkc.App.GetStore().OCRKeyStore.ArchiveEncryptedOCRKeyBundle(&ekb)
	}
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	jsonAPIResponse(c, ekb, "offChainReportingKeyBundle")
}

// Import imports an OCR key bundle
// Example:
// "Post <application>/keys/ocr/import"
func (ocrkc *OCRKeysController) Import(c *gin.Context) {
	defer logger.ErrorIfCalling(c.Request.Body.Close)

	store := ocrkc.App.GetStore()
	bytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}
	oldPassword := c.Query("oldpassword")
	encryptedOCRKeyBundle, err := store.OCRKeyStore.ImportOCRKeyBundle(bytes, oldPassword)
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
	defer logger.ErrorIfCalling(c.Request.Body.Close)

	stringID := c.Param("ID")
	id, err := models.Sha256HashFromHex(stringID)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, errors.New("invalid key ID"))
	}
	newPassword := c.Query("newpassword")
	bytes, err := ocrkc.App.GetStore().OCRKeyStore.ExportOCRKeyBundle(id, newPassword)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	c.Data(http.StatusOK, MediaType, bytes)
}
