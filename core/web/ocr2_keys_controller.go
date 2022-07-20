package web

import (
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
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
	chainType := chaintype.ChainType(c.Param("chainType"))
	key, err := ocr2kc.App.GetKeyStore().OCR2().Create(chainType)
	if errors.Is(errors.Cause(err), chaintype.ErrInvalidChainType) {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	jsonAPIResponse(c, presenters.NewOCR2KeysBundleResource(key), "offChainReporting2KeyBundle")
}

// Delete an OCR2 key bundle
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
	jsonAPIResponse(c, presenters.NewOCR2KeysBundleResource(key), "offChainReporting2KeyBundle")
}

// Import imports an OCR2 key bundle
// Example:
// "Post <application>/keys/ocr/import"
func (ocr2kc *OCR2KeysController) Import(c *gin.Context) {
	defer ocr2kc.App.GetLogger().ErrorIfClosing(c.Request.Body, "Import request body")

	bytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}
	oldPassword := c.Query("oldpassword")
	keyBundle, err := ocr2kc.App.GetKeyStore().OCR2().Import(bytes, oldPassword)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponse(c, presenters.NewOCR2KeysBundleResource(keyBundle), "offChainReporting2KeyBundle")
}

// Export exports an OCR2 key bundle
// Example:
// "Post <application>/keys/ocr/export"
func (ocr2kc *OCR2KeysController) Export(c *gin.Context) {
	defer ocr2kc.App.GetLogger().ErrorIfClosing(c.Request.Body, "Export response body")

	stringID := c.Param("ID")
	newPassword := c.Query("newpassword")
	bytes, err := ocr2kc.App.GetKeyStore().OCR2().Export(stringID, newPassword)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	c.Data(http.StatusOK, MediaType, bytes)
}
