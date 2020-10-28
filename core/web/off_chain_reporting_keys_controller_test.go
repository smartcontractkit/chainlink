package web_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/models/ocrkey"
	"github.com/smartcontractkit/chainlink/core/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOffChainReportingKeysController_Index(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t, cltest.LenientEthMock)
	defer cleanup()
	require.NoError(t, app.Start())
	client := app.NewHTTPClient()

	ocrKeys := []ocrkey.EncryptedKeyBundle{}

	keys, _ := app.GetStore().OCRKeyStore.FindEncryptedOCRKeyBundles()

	response, cleanup := client.Get("/v2/off_chain_reporting_keys")
	defer cleanup()
	cltest.AssertServerResponse(t, response, http.StatusOK)

	err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &ocrKeys)
	assert.NoError(t, err)

	require.Len(t, ocrKeys, len(keys))

	assert.Equal(t, keys[0].ID, ocrKeys[0].ID)
	assert.Equal(t, keys[0].OnChainSigningAddress, ocrKeys[0].OnChainSigningAddress)
	assert.Equal(t, keys[0].OffChainPublicKey, ocrKeys[0].OffChainPublicKey)
	assert.Equal(t, keys[0].ConfigPublicKey, ocrKeys[0].ConfigPublicKey)
}

func TestOffChainReportingKeysController_Create(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t, cltest.LenientEthMock)
	defer cleanup()
	require.NoError(t, app.Start())
	client := app.NewHTTPClient()

	OCRKeyStore := app.GetStore().OCRKeyStore

	keys, _ := OCRKeyStore.FindEncryptedOCRKeyBundles()
	initialLength := len(keys)

	request := models.CreateOCRKeysRequest{
		Password: cltest.Password,
	}
	body, _ := json.Marshal(request)

	response, cleanup := client.Post("/v2/off_chain_reporting_keys", bytes.NewBuffer(body))
	defer cleanup()
	cltest.AssertServerResponse(t, response, http.StatusOK)

	keys, _ = OCRKeyStore.FindEncryptedOCRKeyBundles()
	require.Len(t, keys, initialLength+1)

	ocrKey := ocrkey.EncryptedKeyBundle{}
	err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &ocrKey)
	assert.NoError(t, err)

	lastKeyIndex := len(keys) - 1
	assert.Equal(t, keys[lastKeyIndex].ID, ocrKey.ID)
	assert.Equal(t, keys[lastKeyIndex].OnChainSigningAddress, ocrKey.OnChainSigningAddress)
	assert.Equal(t, keys[lastKeyIndex].OffChainPublicKey, ocrKey.OffChainPublicKey)
	assert.Equal(t, keys[lastKeyIndex].ConfigPublicKey, ocrKey.ConfigPublicKey)

	_, exists := OCRKeyStore.DecryptedOCRKey(ocrKey.ID)
	assert.Equal(t, exists, true)
}
