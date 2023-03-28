package web_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
	"github.com/smartcontractkit/chainlink/v2/core/web"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestP2PKeysController_Index_HappyPath(t *testing.T) {
	t.Parallel()

	client, keyStore := setupP2PKeysControllerTests(t)
	keys, _ := keyStore.P2P().GetAll()

	response, cleanup := client.Get("/v2/keys/p2p")
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, response, http.StatusOK)

	resources := []presenters.P2PKeyResource{}
	err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &resources)
	assert.NoError(t, err)

	require.Len(t, resources, len(keys))

	assert.Equal(t, keys[0].ID(), resources[0].ID)
	assert.Equal(t, keys[0].PublicKeyHex(), resources[0].PubKey)
	assert.Equal(t, keys[0].PeerID().String(), resources[0].PeerID)
}

func TestP2PKeysController_Create_HappyPath(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplicationEVMDisabled(t)
	require.NoError(t, app.Start(testutils.Context(t)))
	client := app.NewHTTPClient(cltest.APIEmailAdmin)
	keyStore := app.GetKeyStore()

	response, cleanup := client.Post("/v2/keys/p2p", nil)
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, response, http.StatusOK)

	keys, _ := keyStore.P2P().GetAll()
	require.Len(t, keys, 1)

	resource := presenters.P2PKeyResource{}
	err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &resource)
	assert.NoError(t, err)

	assert.Equal(t, keys[0].ID(), resource.ID)
	assert.Equal(t, keys[0].PublicKeyHex(), resource.PubKey)
	assert.Equal(t, keys[0].PeerID().String(), resource.PeerID)

	var peerID p2pkey.PeerID
	require.NoError(t, peerID.UnmarshalText([]byte(resource.PeerID)))
	_, err = keyStore.P2P().Get(peerID)
	require.NoError(t, err)
}

func TestP2PKeysController_Delete_NonExistentP2PKeyID(t *testing.T) {
	t.Parallel()

	client, _ := setupP2PKeysControllerTests(t)

	nonExistentP2PKeyID := "12D3KooWPjceQrSwdWXPyLLeABRXmuqt69Rg3sBYbU1Nft9HyQ6a"
	response, cleanup := client.Delete("/v2/keys/p2p/" + nonExistentP2PKeyID)
	t.Cleanup(cleanup)
	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}

func TestP2PKeysController_Delete_InvalidPeerID(t *testing.T) {
	t.Parallel()

	client, _ := setupP2PKeysControllerTests(t)

	nonExistentP2PKeyID := "1234567890"
	response, cleanup := client.Delete("/v2/keys/p2p/" + nonExistentP2PKeyID)
	t.Cleanup(cleanup)
	assert.Equal(t, http.StatusUnprocessableEntity, response.StatusCode)
}

func TestP2PKeysController_Delete_HappyPath(t *testing.T) {
	t.Parallel()

	client, keyStore := setupP2PKeysControllerTests(t)

	keys, _ := keyStore.P2P().GetAll()
	initialLength := len(keys)
	key, _ := keyStore.P2P().Create()

	response, cleanup := client.Delete(fmt.Sprintf("/v2/keys/p2p/%s", key.ID()))
	t.Cleanup(cleanup)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Error(t, utils.JustError(keyStore.P2P().Get(key.PeerID())))

	keys, _ = keyStore.P2P().GetAll()
	assert.Equal(t, initialLength, len(keys))
}

func setupP2PKeysControllerTests(t *testing.T) (cltest.HTTPClientCleaner, keystore.Master) {
	t.Helper()

	app := cltest.NewApplication(t)
	require.NoError(t, app.Start(testutils.Context(t)))
	require.NoError(t, app.KeyStore.OCR().Add(cltest.DefaultOCRKey))
	require.NoError(t, app.KeyStore.P2P().Add(cltest.DefaultP2PKey))

	client := app.NewHTTPClient(cltest.APIEmailAdmin)

	return client, app.GetKeyStore()
}
