package web_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/offchainreporting"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/models/ocrkey"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOffChainReportingKeysController_Index_HappyPath(t *testing.T) {
	client, OCRKeyStore, cleanup := setup(t)
	defer cleanup()

	ocrKeys := []ocrkey.EncryptedKeyBundle{}

	keys, _ := OCRKeyStore.FindEncryptedOCRKeyBundles()

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

func TestOffChainReportingKeysController_Create_InvalidBody(t *testing.T) {
	client, _, cleanup := setup(t)
	defer cleanup()

	invalidBody, _ := json.Marshal(struct {
		BadParam string
	}{
		BadParam: "randomString",
	})
	response, cleanup := client.Post("/v2/off_chain_reporting_keys", bytes.NewBuffer(invalidBody))
	defer cleanup()
	cltest.AssertServerResponse(t, response, http.StatusUnprocessableEntity)
}

func TestOffChainReportingKeysController_Create_HappyPath(t *testing.T) {
	client, OCRKeyStore, cleanup := setup(t)
	defer cleanup()

	keys, _ := OCRKeyStore.FindEncryptedOCRKeyBundles()
	initialLength := len(keys)

	body, _ := json.Marshal(models.CreateOCRKeysRequest{
		Password: cltest.Password,
	})
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

func TestOffChainReportingKeysController_Delete_InvalidOCRKey(t *testing.T) {
	client, _, cleanup := setup(t)
	defer cleanup()

	invalidOCRKeyID := "bad_key_id"
	response, cleanup := client.Delete("/v2/off_chain_reporting_keys/" + invalidOCRKeyID)
	defer cleanup()
	assert.Equal(t, http.StatusUnprocessableEntity, response.StatusCode)
}

func TestOffChainReportingKeysController_Delete_NonExistentOCRKeyID(t *testing.T) {
	client, _, cleanup := setup(t)
	defer cleanup()

	nonExistentOCRKeyID := "eb81f4a35033ac8dd68b9d33a039a713d6fd639af6852b81f47ffeda1c95de54"
	response, cleanup := client.Delete("/v2/off_chain_reporting_keys/" + nonExistentOCRKeyID)
	defer cleanup()
	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}

func TestOffChainReportingKeysController_Delete_HappyPath(t *testing.T) {
	client, OCRKeyStore, cleanup := setup(t)
	defer cleanup()

	keys, _ := OCRKeyStore.FindEncryptedOCRKeyBundles()
	initialLength := len(keys)
	_, encryptedKeyBundle, _ := OCRKeyStore.GenerateEncryptedOCRKeyBundle(cltest.Password)

	response, cleanup := client.Delete("/v2/off_chain_reporting_keys/" + encryptedKeyBundle.ID.String())
	defer cleanup()
	assert.Equal(t, http.StatusNoContent, response.StatusCode)
	assert.Error(t, utils.JustError(OCRKeyStore.FindEncryptedOCRKeyBundleByID(encryptedKeyBundle.ID)))

	keys, _ = OCRKeyStore.FindEncryptedOCRKeyBundles()
	assert.Equal(t, initialLength, len(keys))
}

func setup(t *testing.T) (cltest.HTTPClientCleaner, *offchainreporting.KeyStore, func()) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t, cltest.LenientEthMock)
	require.NoError(t, app.Start())
	client := app.NewHTTPClient()

	OCRKeyStore := app.GetStore().OCRKeyStore
	return client, OCRKeyStore, cleanup
}
