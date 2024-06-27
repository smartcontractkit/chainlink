package web_test

import (
	"fmt"
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

func TestDKGEncryptKeysController_Index_HappyPath(t *testing.T) {
	t.Parallel()

	client, keyStore := setupDKGEncryptKeysControllerTests(t)
	keys, _ := keyStore.DKGEncrypt().GetAll()

	response, cleanup := client.Get("/v2/keys/dkgencrypt")
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, response, http.StatusOK)

	resources := []presenters.DKGEncryptKeyResource{}
	err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &resources)
	assert.NoError(t, err)

	assert.Len(t, resources, len(keys))

	assert.Equal(t, keys[0].ID(), resources[0].ID)
	assert.Equal(t, keys[0].PublicKeyString(), resources[0].PublicKey)
}

func TestDKGEncryptKeysController_Create_HappyPath(t *testing.T) {
	t.Parallel()

	client, keyStore := setupDKGEncryptKeysControllerTests(t)

	response, cleanup := client.Post("/v2/keys/dkgencrypt", nil)
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, response, http.StatusOK)

	keys, _ := keyStore.DKGEncrypt().GetAll()
	assert.Len(t, keys, 2)

	resource := presenters.DKGEncryptKeyResource{}
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

	_, err = keyStore.DKGEncrypt().Get(resource.ID)
	assert.NoError(t, err)
}

func TestDKGEncryptKeysController_Delete_NonExistentDKGEncryptKeyID(t *testing.T) {
	t.Parallel()

	client, _ := setupDKGEncryptKeysControllerTests(t)

	response, cleanup := client.Delete("/v2/keys/dkgencrypt/" + "nonexistentKey")
	t.Cleanup(cleanup)
	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}

func TestDKGEncryptKeysController_Delete_HappyPath(t *testing.T) {
	t.Parallel()

	client, keyStore := setupDKGEncryptKeysControllerTests(t)

	keys, _ := keyStore.DKGEncrypt().GetAll()
	initialLength := len(keys)

	response, cleanup := client.Delete(fmt.Sprintf("/v2/keys/dkgencrypt/%s", keys[0].ID()))
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, response, http.StatusOK)
	assert.Error(t, utils.JustError(keyStore.DKGEncrypt().Get(keys[0].ID())))

	afterKeys, err := keyStore.DKGEncrypt().GetAll()
	assert.NoError(t, err)
	assert.Equal(t, initialLength-1, len(afterKeys))
}

func setupDKGEncryptKeysControllerTests(t *testing.T) (cltest.HTTPClientCleaner, keystore.Master) {
	t.Helper()

	app := cltest.NewApplication(t)
	ctx := testutils.Context(t)
	require.NoError(t, app.Start(ctx))
	require.NoError(t, app.KeyStore.DKGEncrypt().Add(ctx, cltest.DefaultDKGEncryptKey))

	client := app.NewHTTPClient(nil)

	return client, app.GetKeyStore()
}
