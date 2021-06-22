package web

import (
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/smartcontractkit/chainlink/core/services/signatures/secp256k1"

	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

// VRFKeysController manages VRF keys
type VRFKeysController struct {
	App chainlink.Application
}

// Index lists VRF keys
// Example:
// "GET <application>/keys/vrf"
func (vrfkc *VRFKeysController) Index(c *gin.Context) {
	keys, err := vrfkc.App.GetKeyStore().VRF().Get()
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	jsonAPIResponse(c, presenters.NewVRFKeyResources(keys), "vrfKey")
}

// Create and return a VRF key
// Example:
// "POST <application>/keys/vrf"
func (vrfkc *VRFKeysController) Create(c *gin.Context) {
	pk, err := vrfkc.App.GetKeyStore().VRF().CreateKey()
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	encKey, err := vrfkc.App.GetKeyStore().VRF().GetSpecificKey(pk)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	jsonAPIResponse(c, presenters.NewVRFKeyResource(*encKey), "vrfKey")
}

// Delete a VRF key
// Example:
// "DELETE <application>/keys/vrf/:keyID"
// "DELETE <application>/keys/vrf/:keyID?hard=true"
func (vrfkc *VRFKeysController) Delete(c *gin.Context) {
	var hardDelete bool
	var err error
	if c.Query("hard") != "" {
		hardDelete, err = strconv.ParseBool(c.Query("hard"))
		if err != nil {
			jsonAPIError(c, http.StatusUnprocessableEntity, err)
			return
		}
	}
	pk, err := secp256k1.NewPublicKeyFromHex(c.Param("keyID"))
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}
	key, err := vrfkc.App.GetKeyStore().VRF().GetSpecificKey(pk)
	if err != nil {
		jsonAPIError(c, http.StatusNotFound, err)
		return
	}
	if hardDelete {
		err = vrfkc.App.GetKeyStore().VRF().Delete(pk)
	} else {
		err = vrfkc.App.GetKeyStore().VRF().Archive(pk)
	}
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	jsonAPIResponse(c, presenters.NewVRFKeyResource(*key), "vrfKey")
}

// Import imports a VRF key
// Example:
// "Post <application>/keys/vrf/import"
func (vrfkc *VRFKeysController) Import(c *gin.Context) {
	defer logger.ErrorIfCalling(c.Request.Body.Close)

	bytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}
	oldPassword := c.Query("oldpassword")
	key, err := vrfkc.App.GetKeyStore().VRF().Import(bytes, oldPassword)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponse(c, presenters.NewVRFKeyResource(key), "vrfKey")
}

// Export exports a VRF key
// Example:
// "Post <application>/keys/vrf/export/:keyID"
func (vrfkc *VRFKeysController) Export(c *gin.Context) {
	defer logger.ErrorIfCalling(c.Request.Body.Close)

	pk, err := secp256k1.NewPublicKeyFromHex(c.Param("keyID"))
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}
	// New password to re-encrypt the export with
	newPassword := c.Query("newpassword")
	bytes, err := vrfkc.App.GetKeyStore().VRF().Export(pk, newPassword)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	c.Data(http.StatusOK, MediaType, bytes)
}
