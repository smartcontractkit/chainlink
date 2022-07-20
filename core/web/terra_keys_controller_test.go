package web_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/web"
	"github.com/smartcontractkit/chainlink/core/web/presenters"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTerraKeysController_Index_HappyPath(t *testing.T) {
	t.Parallel()

	client, keyStore := setupTerraKeysControllerTests(t)
	keys, _ := keyStore.Terra().GetAll()

	response, cleanup := client.Get("/v2/keys/terra")
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, response, http.StatusOK)

	resources := []presenters.TerraKeyResource{}
	err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &resources)
	assert.NoError(t, err)

	require.Len(t, resources, len(keys))

	assert.Equal(t, keys[0].ID(), resources[0].ID)
	assert.Equal(t, keys[0].PublicKeyStr(), resources[0].PubKey)
}

func TestTerraKeysController_Create_HappyPath(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplicationEVMDisabled(t)
	require.NoError(t, app.Start(testutils.Context(t)))
	client := app.NewHTTPClient()
	keyStore := app.GetKeyStore()

	response, cleanup := client.Post("/v2/keys/terra", nil)
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, response, http.StatusOK)

	keys, _ := keyStore.Terra().GetAll()
	require.Len(t, keys, 1)

	resource := presenters.TerraKeyResource{}
	err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &resource)
	assert.NoError(t, err)

	assert.Equal(t, keys[0].ID(), resource.ID)
	assert.Equal(t, keys[0].PublicKeyStr(), resource.PubKey)

	_, err = keyStore.Terra().Get(resource.ID)
	require.NoError(t, err)
}

func TestTerraKeysController_Delete_NonExistentTerraKeyID(t *testing.T) {
	t.Parallel()

	client, _ := setupTerraKeysControllerTests(t)

	nonExistentTerraKeyID := "foobar"
	response, cleanup := client.Delete("/v2/keys/terra/" + nonExistentTerraKeyID)
	t.Cleanup(cleanup)
	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}

func TestTerraKeysController_Delete_HappyPath(t *testing.T) {
	t.Parallel()

	client, keyStore := setupTerraKeysControllerTests(t)

	keys, _ := keyStore.Terra().GetAll()
	initialLength := len(keys)
	key, _ := keyStore.Terra().Create()

	response, cleanup := client.Delete(fmt.Sprintf("/v2/keys/terra/%s", key.ID()))
	t.Cleanup(cleanup)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Error(t, utils.JustError(keyStore.Terra().Get(key.ID())))

	keys, _ = keyStore.Terra().GetAll()
	assert.Equal(t, initialLength, len(keys))
}

func setupTerraKeysControllerTests(t *testing.T) (cltest.HTTPClientCleaner, keystore.Master) {
	t.Helper()

	app := cltest.NewApplication(t)
	require.NoError(t, app.Start(testutils.Context(t)))
	app.KeyStore.Terra().Add(cltest.DefaultTerraKey)

	client := app.NewHTTPClient()

	return client, app.GetKeyStore()
}
