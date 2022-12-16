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

func TestSolanaKeysController_Index_HappyPath(t *testing.T) {
	t.Parallel()

	client, keyStore := setupSolanaKeysControllerTests(t)
	keys, _ := keyStore.Solana().GetAll()

	response, cleanup := client.Get("/v2/keys/solana")
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, response, http.StatusOK)

	resources := []presenters.SolanaKeyResource{}
	err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &resources)
	assert.NoError(t, err)

	require.Len(t, resources, len(keys))

	assert.Equal(t, keys[0].ID(), resources[0].ID)
	assert.Equal(t, keys[0].PublicKeyStr(), resources[0].PubKey)
}

func TestSolanaKeysController_Create_HappyPath(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplicationEVMDisabled(t)
	require.NoError(t, app.Start(testutils.Context(t)))
	client := app.NewHTTPClient(cltest.APIEmailAdmin)
	keyStore := app.GetKeyStore()

	response, cleanup := client.Post("/v2/keys/solana", nil)
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, response, http.StatusOK)

	keys, _ := keyStore.Solana().GetAll()
	require.Len(t, keys, 1)

	resource := presenters.SolanaKeyResource{}
	err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &resource)
	assert.NoError(t, err)

	assert.Equal(t, keys[0].ID(), resource.ID)
	assert.Equal(t, keys[0].PublicKeyStr(), resource.PubKey)

	_, err = keyStore.Solana().Get(resource.ID)
	require.NoError(t, err)
}

func TestSolanaKeysController_Delete_NonExistentSolanaKeyID(t *testing.T) {
	t.Parallel()

	client, _ := setupSolanaKeysControllerTests(t)

	nonExistentSolanaKeyID := "foobar"
	response, cleanup := client.Delete("/v2/keys/solana/" + nonExistentSolanaKeyID)
	t.Cleanup(cleanup)
	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}

func TestSolanaKeysController_Delete_HappyPath(t *testing.T) {
	t.Parallel()

	client, keyStore := setupSolanaKeysControllerTests(t)

	keys, _ := keyStore.Solana().GetAll()
	initialLength := len(keys)
	key, _ := keyStore.Solana().Create()

	response, cleanup := client.Delete(fmt.Sprintf("/v2/keys/solana/%s", key.ID()))
	t.Cleanup(cleanup)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Error(t, utils.JustError(keyStore.Solana().Get(key.ID())))

	keys, _ = keyStore.Solana().GetAll()
	assert.Equal(t, initialLength, len(keys))
}

func setupSolanaKeysControllerTests(t *testing.T) (cltest.HTTPClientCleaner, keystore.Master) {
	t.Helper()

	app := cltest.NewApplication(t)
	require.NoError(t, app.Start(testutils.Context(t)))
	require.NoError(t, app.KeyStore.OCR().Add(cltest.DefaultOCRKey))
	require.NoError(t, app.KeyStore.Solana().Add(cltest.DefaultSolanaKey))

	client := app.NewHTTPClient(cltest.APIEmailAdmin)

	return client, app.GetKeyStore()
}
