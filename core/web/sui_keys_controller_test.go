package web_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/utils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/web"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func TestSuiKeysController_Index_HappyPath(t *testing.T) {
	t.Parallel()

	client, keyStore := setupSuiKeysControllerTests(t)
	keys, _ := keyStore.Sui().GetAll()

	response, cleanup := client.Get("/v2/keys/sui")
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, response, http.StatusOK)

	resources := []presenters.SuiKeyResource{}
	err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &resources)
	require.NoError(t, err)

	require.Len(t, resources, len(keys))

	assert.Equal(t, keys[0].ID(), resources[0].ID)
	assert.Equal(t, keys[0].PublicKeyStr(), resources[0].PubKey)
}

func TestSuiKeysController_Create_HappyPath(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplicationEVMDisabled(t)
	require.NoError(t, app.Start(testutils.Context(t)))
	client := app.NewHTTPClient(nil)
	keyStore := app.GetKeyStore()

	response, cleanup := client.Post("/v2/keys/sui", nil)
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, response, http.StatusOK)

	keys, _ := keyStore.Sui().GetAll()
	require.Len(t, keys, 1)

	resource := presenters.SuiKeyResource{}
	err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &resource)
	require.NoError(t, err)

	assert.Equal(t, keys[0].ID(), resource.ID)
	assert.Equal(t, keys[0].PublicKeyStr(), resource.PubKey)

	_, err = keyStore.Sui().Get(resource.ID)
	require.NoError(t, err)
}

func TestSuiKeysController_Delete_NonExistentSuiKeyID(t *testing.T) {
	t.Parallel()

	client, _ := setupSuiKeysControllerTests(t)

	nonExistentSuiKeyID := "foobar"
	response, cleanup := client.Delete("/v2/keys/sui/" + nonExistentSuiKeyID)
	t.Cleanup(cleanup)
	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}

func TestSuiKeysController_Delete_HappyPath(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	client, keyStore := setupSuiKeysControllerTests(t)

	keys, _ := keyStore.Sui().GetAll()
	initialLength := len(keys)
	key, _ := keyStore.Sui().Create(ctx)

	response, cleanup := client.Delete("/v2/keys/sui/" + key.ID())
	t.Cleanup(cleanup)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	require.Error(t, utils.JustError(keyStore.Sui().Get(key.ID())))

	keys, _ = keyStore.Sui().GetAll()
	assert.Len(t, keys, initialLength)
}

func setupSuiKeysControllerTests(t *testing.T) (cltest.HTTPClientCleaner, keystore.Master) {
	t.Helper()
	ctx := testutils.Context(t)

	app := cltest.NewApplication(t)
	require.NoError(t, app.Start(ctx))
	require.NoError(t, app.KeyStore.OCR().Add(ctx, cltest.DefaultOCRKey))
	suiKeyStore := app.GetKeyStore().Sui()
	require.NotNil(t, suiKeyStore)
	require.NoError(t, suiKeyStore.Add(ctx, cltest.DefaultSuiKey))

	client := app.NewHTTPClient(nil)

	return client, app.GetKeyStore()
}
