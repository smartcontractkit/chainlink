package web_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
	"github.com/smartcontractkit/chainlink/v2/core/web"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func TestDKGSignKeysController_Index_HappyPath(t *testing.T) {
	t.Parallel()

	client, keyStore := setupDKGSignKeysControllerTests(t)
	keys, _ := keyStore.DKGSign().GetAll()

	response, cleanup := client.Get("/v2/keys/dkgsign")
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, response, http.StatusOK)

	resources := []presenters.DKGSignKeyResource{}
	err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &resources)
	assert.NoError(t, err)

	assert.Len(t, resources, len(keys))

	assert.Equal(t, keys[0].ID(), resources[0].ID)
	assert.Equal(t, keys[0].PublicKeyString(), resources[0].PublicKey)
}

func TestDKGSignKeysController_Create_HappyPath(t *testing.T) {
	t.Parallel()

	client, keyStore := setupDKGSignKeysControllerTests(t)

	response, cleanup := client.Post("/v2/keys/dkgsign", nil)
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, response, http.StatusOK)

	keys, _ := keyStore.DKGSign().GetAll()
	assert.Len(t, keys, 2)

	resource := presenters.DKGSignKeyResource{}
	err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &resource)
	assert.NoError(t, err)

	// the order in which keys are returned by GetAll is non-deterministic
	// due to map iteration non-determinism.
	found := false
	for _, key := range keys {
		if key.ID() == resource.ID {
			assert.Equal(t, key.PublicKeyString(), resource.PublicKey)
			found = true
		}
	}
	assert.True(t, found)

	_, err = keyStore.DKGSign().Get(resource.ID)
	assert.NoError(t, err)
}

func TestDKGSignKeysController_Delete_NonExistentDKGSignKeyID(t *testing.T) {
	t.Parallel()

	client, _ := setupDKGSignKeysControllerTests(t)

	response, cleanup := client.Delete("/v2/keys/dkgsign/" + "nonexistentKey")
	t.Cleanup(cleanup)
	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}

func TestDKGSignKeysController_Delete_HappyPath(t *testing.T) {
	t.Parallel()

	client, keyStore := setupDKGSignKeysControllerTests(t)

	keys, _ := keyStore.DKGSign().GetAll()
	initialLength := len(keys)

	response, cleanup := client.Delete(fmt.Sprintf("/v2/keys/dkgsign/%s", keys[0].ID()))
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, response, http.StatusOK)
	assert.Error(t, utils.JustError(keyStore.DKGSign().Get(keys[0].ID())))

	afterKeys, err := keyStore.DKGSign().GetAll()
	assert.NoError(t, err)
	assert.Equal(t, initialLength-1, len(afterKeys))

}

func setupDKGSignKeysControllerTests(t *testing.T) (cltest.HTTPClientCleaner, keystore.Master) {
	t.Helper()

	app := cltest.NewApplication(t)
	require.NoError(t, app.Start(testutils.Context(t)))
	require.NoError(t, app.KeyStore.DKGSign().Add(cltest.DefaultDKGSignKey))

	client := app.NewHTTPClient(cltest.APIEmailAdmin)

	return client, app.GetKeyStore()
}
