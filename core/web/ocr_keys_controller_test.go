package web_test

import (
	"net/http"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/web"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOCRKeysController_Index_HappyPath(t *testing.T) {
	client, OCRKeyStore := setupOCRKeysControllerTests(t)

	keys, _ := OCRKeyStore.FindEncryptedOCRKeyBundles()

	response, cleanup := client.Get("/v2/keys/ocr")
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, response, http.StatusOK)

	resources := []presenters.OCRKeysBundleResource{}
	err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &resources)
	require.NoError(t, err)

	require.Len(t, resources, len(keys))
	assert.Equal(t, keys[0].ID.String(), resources[0].ID)
	assert.Equal(t, keys[0].OnChainSigningAddress, resources[0].OnChainSigningAddress)
	assert.Equal(t, keys[0].OffChainPublicKey, resources[0].OffChainPublicKey)
	assert.Equal(t, keys[0].ConfigPublicKey, resources[0].ConfigPublicKey)
}

func TestOCRKeysController_Create_HappyPath(t *testing.T) {
	client, OCRKeyStore := setupOCRKeysControllerTests(t)

	keys, _ := OCRKeyStore.FindEncryptedOCRKeyBundles()
	initialLength := len(keys)

	response, cleanup := client.Post("/v2/keys/ocr", nil)
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, response, http.StatusOK)

	keys, _ = OCRKeyStore.FindEncryptedOCRKeyBundles()
	require.Len(t, keys, initialLength+1)

	resource := presenters.OCRKeysBundleResource{}
	err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &resource)
	assert.NoError(t, err)

	lastKeyIndex := len(keys) - 1
	assert.Equal(t, keys[lastKeyIndex].ID.String(), resource.ID)
	assert.Equal(t, keys[lastKeyIndex].OnChainSigningAddress, resource.OnChainSigningAddress)
	assert.Equal(t, keys[lastKeyIndex].OffChainPublicKey, resource.OffChainPublicKey)
	assert.Equal(t, keys[lastKeyIndex].ConfigPublicKey, resource.ConfigPublicKey)

	idSHA256, err := models.Sha256HashFromHex(resource.ID)
	require.NoError(t, err)
	_, exists := OCRKeyStore.DecryptedOCRKey(idSHA256)
	assert.Equal(t, exists, true)
}

func TestOCRKeysController_Delete_InvalidOCRKey(t *testing.T) {
	client, _ := setupOCRKeysControllerTests(t)

	invalidOCRKeyID := "bad_key_id"
	response, cleanup := client.Delete("/v2/keys/ocr/" + invalidOCRKeyID)
	t.Cleanup(cleanup)
	assert.Equal(t, http.StatusUnprocessableEntity, response.StatusCode)
}

func TestOCRKeysController_Delete_NonExistentOCRKeyID(t *testing.T) {
	client, _ := setupOCRKeysControllerTests(t)

	nonExistentOCRKeyID := "eb81f4a35033ac8dd68b9d33a039a713d6fd639af6852b81f47ffeda1c95de54"
	response, cleanup := client.Delete("/v2/keys/ocr/" + nonExistentOCRKeyID)
	t.Cleanup(cleanup)
	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}

func TestOCRKeysController_Delete_HappyPath(t *testing.T) {
	client, OCRKeyStore := setupOCRKeysControllerTests(t)
	require.NoError(t, OCRKeyStore.Unlock(cltest.Password))

	keys, _ := OCRKeyStore.FindEncryptedOCRKeyBundles()
	initialLength := len(keys)
	_, encryptedKeyBundle, _ := OCRKeyStore.GenerateEncryptedOCRKeyBundle()

	response, cleanup := client.Delete("/v2/keys/ocr/" + encryptedKeyBundle.ID.String())
	t.Cleanup(cleanup)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Error(t, utils.JustError(OCRKeyStore.FindEncryptedOCRKeyBundleByID(encryptedKeyBundle.ID)))

	keys, _ = OCRKeyStore.FindEncryptedOCRKeyBundles()
	assert.Equal(t, initialLength, len(keys))
}

func setupOCRKeysControllerTests(t *testing.T) (cltest.HTTPClientCleaner, *keystore.OCR) {
	t.Parallel()

	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	t.Cleanup(assertMocksCalled)
	app, cleanup := cltest.NewApplication(t,
		ethClient,
	)
	t.Cleanup(cleanup)
	require.NoError(t, app.StartAndConnect())
	client := app.NewHTTPClient()

	return client, app.GetKeyStore().OCR
}
