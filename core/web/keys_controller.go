package web

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/manyminds/api2go/jsonapi"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/logger/audit"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/solkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/terrakey"
)

type Keystore[K keystore.Key] interface {
	Get(id string) (K, error)
	GetAll() ([]K, error)
	Create() (K, error)
	Delete(id string) (K, error)
	Import(keyJSON []byte, password string) (K, error)
	Export(id string, password string) ([]byte, error)
}

type KeysController interface {
	// Index lists keys
	Index(*gin.Context)
	// Create and return a key
	Create(*gin.Context)
	// Delete a key
	Delete(*gin.Context)
	// Import imports a key
	Import(*gin.Context)
	// Export exports a key
	Export(*gin.Context)
}

type keysController[K keystore.Key, R jsonapi.EntityNamer] struct {
	ks           Keystore[K]
	lggr         logger.Logger
	auditLogger  audit.AuditLogger
	resourceName string
	newResource  func(K) *R
	newResources func([]K) []R
}

func NewKeysController[K keystore.Key, R jsonapi.EntityNamer](ks Keystore[K], lggr logger.Logger, auditLogger audit.AuditLogger, resourceName string,
	newResource func(K) *R, newResources func([]K) []R) KeysController {
	return &keysController[K, R]{
		ks:           ks,
		lggr:         lggr,
		auditLogger:  auditLogger,
		resourceName: resourceName,
		newResource:  newResource,
		newResources: newResources,
	}
}

func (kc *keysController[K, R]) Index(c *gin.Context) {
	keys, err := kc.ks.GetAll()
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	jsonAPIResponse(c, kc.newResources(keys), kc.resourceName)
}

func (kc *keysController[K, R]) Create(c *gin.Context) {
	key, err := kc.ks.Create()
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	// Emit audit log, determine if Terra or Solana key
	switch unwrappedKey := any(key).(type) {
	case terrakey.Key:
		kc.auditLogger.Audit(audit.TerraKeyCreated, map[string]interface{}{
			"publicKey": unwrappedKey.PublicKey(),
			"id":        unwrappedKey.ID(),
		})
	case solkey.Key:
		kc.auditLogger.Audit(audit.SolanaKeyCreated, map[string]interface{}{
			"publicKey": unwrappedKey.PublicKey(),
			"id":        unwrappedKey.ID(),
		})
	}

	jsonAPIResponse(c, kc.newResource(key), kc.resourceName)
}

func (kc *keysController[K, R]) Delete(c *gin.Context) {
	keyID := c.Param("keyID")
	key, err := kc.ks.Get(keyID)
	if err != nil {
		jsonAPIError(c, http.StatusNotFound, err)
		return
	}
	_, err = kc.ks.Delete(key.ID())
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	// Emit audit log, determine if Terra or Solana key
	switch any(key).(type) {
	case terrakey.Key:
		kc.auditLogger.Audit(audit.TerraKeyDeleted, map[string]interface{}{"id": keyID})
	case solkey.Key:
		kc.auditLogger.Audit(audit.SolanaKeyDeleted, map[string]interface{}{"id": keyID})
	}

	jsonAPIResponse(c, kc.newResource(key), kc.resourceName)
}

func (kc *keysController[K, R]) Import(c *gin.Context) {
	defer kc.lggr.ErrorIfClosing(c.Request.Body, "Import ")

	bytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}
	oldPassword := c.Query("oldpassword")
	key, err := kc.ks.Import(bytes, oldPassword)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	// Emit audit log, determine if Terra or Solana key
	switch unwrappedKey := any(key).(type) {
	case terrakey.Key:
		kc.auditLogger.Audit(audit.TerraKeyImported, map[string]interface{}{
			"publicKey": unwrappedKey.PublicKey(),
			"id":        unwrappedKey.ID(),
		})
	case solkey.Key:
		kc.auditLogger.Audit(audit.SolanaKeyImported, map[string]interface{}{
			"publicKey": unwrappedKey.PublicKey(),
			"id":        unwrappedKey.ID(),
		})
	}

	jsonAPIResponse(c, kc.newResource(key), kc.resourceName)
}

func (kc *keysController[K, R]) Export(c *gin.Context) {
	defer kc.lggr.ErrorIfClosing(c.Request.Body, "Export request body")

	keyID := c.Param("ID")
	newPassword := c.Query("newpassword")
	bytes, err := kc.ks.Export(keyID, newPassword)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	if strings.HasPrefix(c.Request.URL.Path, "/v2/keys/terra") {
		kc.auditLogger.Audit(audit.TerraKeyExported, map[string]interface{}{"id": keyID})
	} else if strings.HasPrefix(c.Request.URL.Path, "/v2/keys/solana") {
		kc.auditLogger.Audit(audit.SolanaKeyExported, map[string]interface{}{"id": keyID})
	}

	c.Data(http.StatusOK, MediaType, bytes)
}
