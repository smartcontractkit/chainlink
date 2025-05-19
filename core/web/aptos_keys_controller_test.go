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

func TestAptosKeysController_Index_HappyPath(t *testing.T) {
	t.Parallel()

	client, keyStore := setupAptosKeysControllerTests(t)
	keys, _ := keyStore.Aptos().GetAll()

	response, cleanup := client.Get("/v2/keys/aptos")
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, response, http.StatusOK)

	resources := []presenters.AptosKeyResource{}
	err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &resources)
	assert.NoError(t, err)

	require.Len(t, resources, len(keys))

	assert.Equal(t, keys[0].ID(), resources[0].ID)
	assert.Equal(t, keys[0].PublicKeyStr(), resources[0].PubKey)
}

func TestAptosKeysController_Create_HappyPath(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplicationEVMDisabled(t)
	require.NoError(t, app.Start(testutils.Context(t)))
	client := app.NewHTTPClient(nil)
	keyStore := app.GetKeyStore()

	response, cleanup := client.Post("/v2/keys/aptos", nil)
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, response, http.StatusOK)

	keys, _ := keyStore.Aptos().GetAll()
	require.Len(t, keys, 1)

	resource := presenters.AptosKeyResource{}
	err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &resource)
	assert.NoError(t, err)

	assert.Equal(t, keys[0].ID(), resource.ID)
	assert.Equal(t, keys[0].PublicKeyStr(), resource.PubKey)

	_, err = keyStore.Aptos().Get(resource.ID)
	require.NoError(t, err)
}

func TestAptosKeysController_Delete_NonExistentAptosKeyID(t *testing.T) {
	t.Parallel()

	client, _ := setupAptosKeysControllerTests(t)

	nonExistentAptosKeyID := "foobar"
	response, cleanup := client.Delete("/v2/keys/aptos/" + nonExistentAptosKeyID)
	t.Cleanup(cleanup)
	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}

func TestAptosKeysController_Delete_HappyPath(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	client, keyStore := setupAptosKeysControllerTests(t)

	keys, _ := keyStore.Aptos().GetAll()
	initialLength := len(keys)
	key, _ := keyStore.Aptos().Create(ctx)

	response, cleanup := client.Delete("/v2/keys/aptos/" + key.ID())
	t.Cleanup(cleanup)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Error(t, utils.JustError(keyStore.Aptos().Get(key.ID())))

	keys, _ = keyStore.Aptos().GetAll()
	assert.Len(t, keys, initialLength)
}

func setupAptosKeysControllerTests(t *testing.T) (cltest.HTTPClientCleaner, keystore.Master) {
	t.Helper()
	ctx := testutils.Context(t)

	app := cltest.NewApplication(t)
	require.NoError(t, app.Start(ctx))
	require.NoError(t, app.KeyStore.OCR().Add(ctx, cltest.DefaultOCRKey))
	aptosKeyStore := app.GetKeyStore().Aptos()
	require.NotNil(t, aptosKeyStore)
	require.NoError(t, aptosKeyStore.Add(ctx, cltest.DefaultAptosKey))

	client := app.NewHTTPClient(nil)

	return client, app.GetKeyStore()
}
