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

func TestCosmosKeysController_Index_HappyPath(t *testing.T) {
	t.Parallel()

	client, keyStore := setupCosmosKeysControllerTests(t)
	keys, _ := keyStore.Cosmos().GetAll()

	response, cleanup := client.Get("/v2/keys/cosmos")
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, response, http.StatusOK)

	resources := []presenters.CosmosKeyResource{}
	err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &resources)
	assert.NoError(t, err)

	require.Len(t, resources, len(keys))

	assert.Equal(t, keys[0].ID(), resources[0].ID)
	assert.Equal(t, keys[0].PublicKeyStr(), resources[0].PubKey)
}

func TestCosmosKeysController_Create_HappyPath(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplicationEVMDisabled(t)
	ctx := testutils.Context(t)
	require.NoError(t, app.Start(ctx))
	client := app.NewHTTPClient(nil)
	keyStore := app.GetKeyStore()

	response, cleanup := client.Post("/v2/keys/cosmos", nil)
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, response, http.StatusOK)

	keys, _ := keyStore.Cosmos().GetAll()
	require.Len(t, keys, 1)

	resource := presenters.CosmosKeyResource{}
	err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &resource)
	assert.NoError(t, err)

	assert.Equal(t, keys[0].ID(), resource.ID)
	assert.Equal(t, keys[0].PublicKeyStr(), resource.PubKey)

	_, err = keyStore.Cosmos().Get(resource.ID)
	require.NoError(t, err)
}

func TestCosmosKeysController_Delete_NonExistentCosmosKeyID(t *testing.T) {
	t.Parallel()

	client, _ := setupCosmosKeysControllerTests(t)

	nonExistentCosmosKeyID := "foobar"
	response, cleanup := client.Delete("/v2/keys/cosmos/" + nonExistentCosmosKeyID)
	t.Cleanup(cleanup)
	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}

func TestCosmosKeysController_Delete_HappyPath(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	client, keyStore := setupCosmosKeysControllerTests(t)

	keys, _ := keyStore.Cosmos().GetAll()
	initialLength := len(keys)
	key, _ := keyStore.Cosmos().Create(ctx)

	response, cleanup := client.Delete("/v2/keys/cosmos/" + key.ID())
	t.Cleanup(cleanup)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Error(t, utils.JustError(keyStore.Cosmos().Get(key.ID())))

	keys, _ = keyStore.Cosmos().GetAll()
	assert.Len(t, keys, initialLength)
}

func setupCosmosKeysControllerTests(t *testing.T) (cltest.HTTPClientCleaner, keystore.Master) {
	t.Helper()
	ctx := testutils.Context(t)

	app := cltest.NewApplication(t)
	require.NoError(t, app.Start(ctx))
	require.NoError(t, app.KeyStore.Cosmos().Add(ctx, cltest.DefaultCosmosKey))

	client := app.NewHTTPClient(nil)

	return client, app.GetKeyStore()
}
