package web_test

import (
	"net/http"
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/utils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/web"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOCRKeysController_Index_HappyPath(t *testing.T) {
	client, OCRKeyStore := setupOCRKeysControllerTests(t)

	keys, _ := OCRKeyStore.GetAll()

	response, cleanup := client.Get("/v2/keys/ocr")
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, response, http.StatusOK)

	resources := []presenters.OCRKeysBundleResource{}
	err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &resources)
	require.NoError(t, err)

	require.Len(t, resources, len(keys))
	assert.Equal(t, keys[0].ID(), resources[0].ID)
}

func TestOCRKeysController_Create_HappyPath(t *testing.T) {
	client, OCRKeyStore := setupOCRKeysControllerTests(t)

	keys, _ := OCRKeyStore.GetAll()
	initialLength := len(keys)

	response, cleanup := client.Post("/v2/keys/ocr", nil)
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, response, http.StatusOK)

	keys, _ = OCRKeyStore.GetAll()
	require.Len(t, keys, initialLength+1)

	resource := presenters.OCRKeysBundleResource{}
	err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &resource)
	assert.NoError(t, err)

	var ids []string
	for _, key := range keys {
		ids = append(ids, key.ID())
	}
	require.Contains(t, ids, resource.ID)

	_, err = OCRKeyStore.Get(resource.ID)
	require.NoError(t, err)
}

func TestOCRKeysController_Delete_NonExistentOCRKeyID(t *testing.T) {
	client, _ := setupOCRKeysControllerTests(t)

	nonExistentOCRKeyID := "eb81f4a35033ac8dd68b9d33a039a713d6fd639af6852b81f47ffeda1c95de54"
	response, cleanup := client.Delete("/v2/keys/ocr/" + nonExistentOCRKeyID)
	t.Cleanup(cleanup)
	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}

func TestOCRKeysController_Delete_HappyPath(t *testing.T) {
	ctx := testutils.Context(t)
	client, OCRKeyStore := setupOCRKeysControllerTests(t)

	keys, _ := OCRKeyStore.GetAll()
	initialLength := len(keys)
	key, _ := OCRKeyStore.Create(ctx)

	response, cleanup := client.Delete("/v2/keys/ocr/" + key.ID())
	t.Cleanup(cleanup)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Error(t, utils.JustError(OCRKeyStore.Get(key.ID())))

	keys, _ = OCRKeyStore.GetAll()
	assert.Len(t, keys, initialLength)
}

func setupOCRKeysControllerTests(t *testing.T) (cltest.HTTPClientCleaner, keystore.OCR) {
	t.Parallel()
	ctx := testutils.Context(t)

	app := cltest.NewApplicationEVMDisabled(t)
	require.NoError(t, app.Start(testutils.Context(t)))
	client := app.NewHTTPClient(nil)

	require.NoError(t, app.KeyStore.OCR().Add(ctx, cltest.DefaultOCRKey))

	return client, app.GetKeyStore().OCR()
}
