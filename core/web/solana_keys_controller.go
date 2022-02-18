package web

import (
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

// SolanaKeysController manages Solana keys
type SolanaKeysController struct {
	App chainlink.Application
}

// Index lists Solana keys
// Example:
// "GET <application>/keys/solana"
func (solkc *SolanaKeysController) Index(c *gin.Context) {
	keys, err := solkc.App.GetKeyStore().Solana().GetAll()
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	jsonAPIResponse(c, presenters.NewSolanaKeyResources(keys), "solanaKey")
}

// Create and return a Solana key
// Example:
// "POST <application>/keys/solana"
func (solkc *SolanaKeysController) Create(c *gin.Context) {
	key, err := solkc.App.GetKeyStore().Solana().Create()
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	jsonAPIResponse(c, presenters.NewSolanaKeyResource(key), "solanaKey")
}

// Delete a Solana key
// Example:
// "DELETE <application>/keys/solana/:keyID"
// "DELETE <application>/keys/solana/:keyID?hard=true"
func (solkc *SolanaKeysController) Delete(c *gin.Context) {
	keyID := c.Param("keyID")
	key, err := solkc.App.GetKeyStore().Solana().Get(keyID)
	if err != nil {
		jsonAPIError(c, http.StatusNotFound, err)
		return
	}
	_, err = solkc.App.GetKeyStore().Solana().Delete(key.ID())
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	jsonAPIResponse(c, presenters.NewSolanaKeyResource(key), "solanaKey")
}

// Import imports a Solana key
// Example:
// "Post <application>/keys/solana/import"
func (solkc *SolanaKeysController) Import(c *gin.Context) {
	defer solkc.App.GetLogger().ErrorIfClosing(c.Request.Body, "Import ")

	bytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}
	oldPassword := c.Query("oldpassword")
	key, err := solkc.App.GetKeyStore().Solana().Import(bytes, oldPassword)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponse(c, presenters.NewSolanaKeyResource(key), "solanaKey")
}

// Export exports a Solana key
// Example:
// "Post <application>/keys/solana/export"
func (solkc *SolanaKeysController) Export(c *gin.Context) {
	defer solkc.App.GetLogger().ErrorIfClosing(c.Request.Body, "Export request body")

	keyID := c.Param("ID")
	newPassword := c.Query("newpassword")
	bytes, err := solkc.App.GetKeyStore().Solana().Export(keyID, newPassword)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	c.Data(http.StatusOK, MediaType, bytes)
}
