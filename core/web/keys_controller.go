package web

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/manyminds/api2go/jsonapi"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/logger/audit"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
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
	lggr         logger.SugaredLogger
	auditLogger  audit.AuditLogger
	typ          string
	resourceName string
	newResource  func(K) *R
	newResources func([]K) []R
}

func NewKeysController[K keystore.Key, R jsonapi.EntityNamer](ks Keystore[K], lggr logger.Logger, auditLogger audit.AuditLogger, resourceName string,
	newResource func(K) *R, newResources func([]K) []R) KeysController {
	var k K
	typ, err := keystore.GetFieldNameForKey(k)
	if err != nil {
		panic(fmt.Errorf("unable to create keys controller: %v", err))
	}
	return &keysController[K, R]{
		ks:           ks,
		lggr:         logger.Sugared(lggr),
		auditLogger:  auditLogger,
		typ:          typ,
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

	kc.auditLogger.Audit(audit.KeyCreated, map[string]interface{}{
		"type": kc.typ,
		"id":   key.ID(),
	})

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

	kc.auditLogger.Audit(audit.KeyDeleted, map[string]interface{}{
		"type": kc.typ,
		"id":   key.ID(),
	})

	jsonAPIResponse(c, kc.newResource(key), kc.resourceName)
}

func (kc *keysController[K, R]) Import(c *gin.Context) {
	defer kc.lggr.ErrorIfFn(c.Request.Body.Close, "Error closing Import request body")

	bytes, err := io.ReadAll(c.Request.Body)
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

	kc.auditLogger.Audit(audit.KeyImported, map[string]interface{}{
		"type": kc.typ,
		"id":   key.ID(),
	})

	jsonAPIResponse(c, kc.newResource(key), kc.resourceName)
}

func (kc *keysController[K, R]) Export(c *gin.Context) {
	defer kc.lggr.ErrorIfFn(c.Request.Body.Close, "Error closing Export request body")

	keyID := c.Param("ID")
	newPassword := c.Query("newpassword")
	bytes, err := kc.ks.Export(keyID, newPassword)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	kc.auditLogger.Audit(audit.KeyExported, map[string]interface{}{
		"type": kc.typ,
		"id":   keyID,
	})

	c.Data(http.StatusOK, MediaType, bytes)
}
