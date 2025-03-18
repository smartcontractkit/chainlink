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

func TestStarkNetKeysController_Index_HappyPath(t *testing.T) {
	t.Parallel()

	client, keyStore := setupStarkNetKeysControllerTests(t)
	keys, _ := keyStore.StarkNet().GetAll()

	response, cleanup := client.Get("/v2/keys/starknet")
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, response, http.StatusOK)

	resources := []presenters.StarkNetKeyResource{}
	err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &resources)
	assert.NoError(t, err)

	require.Len(t, resources, len(keys))

	assert.Equal(t, keys[0].ID(), resources[0].ID)
	assert.Equal(t, keys[0].StarkKeyStr(), resources[0].StarkKey)
}

func TestStarkNetKeysController_Create_HappyPath(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplicationEVMDisabled(t)
	require.NoError(t, app.Start(testutils.Context(t)))
	client := app.NewHTTPClient(nil)
	keyStore := app.GetKeyStore()

	response, cleanup := client.Post("/v2/keys/starknet", nil)
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, response, http.StatusOK)

	keys, _ := keyStore.StarkNet().GetAll()
	require.Len(t, keys, 1)

	resource := presenters.StarkNetKeyResource{}
	err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &resource)
	assert.NoError(t, err)

	assert.Equal(t, keys[0].ID(), resource.ID)
	assert.Equal(t, keys[0].StarkKeyStr(), resource.StarkKey)

	_, err = keyStore.StarkNet().Get(resource.ID)
	require.NoError(t, err)
}

func TestStarkNetKeysController_Delete_NonExistentStarkNetKeyID(t *testing.T) {
	t.Parallel()

	client, _ := setupStarkNetKeysControllerTests(t)

	nonExistentStarkNetKeyID := "foobar"
	response, cleanup := client.Delete("/v2/keys/starknet/" + nonExistentStarkNetKeyID)
	t.Cleanup(cleanup)
	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}

func TestStarkNetKeysController_Delete_HappyPath(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	client, keyStore := setupStarkNetKeysControllerTests(t)

	keys, _ := keyStore.StarkNet().GetAll()
	initialLength := len(keys)
	key, _ := keyStore.StarkNet().Create(ctx)

	response, cleanup := client.Delete("/v2/keys/starknet/" + key.ID())
	t.Cleanup(cleanup)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Error(t, utils.JustError(keyStore.StarkNet().Get(key.ID())))

	keys, _ = keyStore.StarkNet().GetAll()
	assert.Len(t, keys, initialLength)
}

func setupStarkNetKeysControllerTests(t *testing.T) (cltest.HTTPClientCleaner, keystore.Master) {
	t.Helper()
	ctx := testutils.Context(t)

	app := cltest.NewApplication(t)
	require.NoError(t, app.Start(ctx))
	require.NoError(t, app.KeyStore.OCR().Add(ctx, cltest.DefaultOCRKey))
	require.NoError(t, app.KeyStore.StarkNet().Add(ctx, cltest.DefaultStarkNetKey))

	client := app.NewHTTPClient(nil)

	return client, app.GetKeyStore()
}
