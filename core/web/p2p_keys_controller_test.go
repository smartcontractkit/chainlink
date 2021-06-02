package web_test

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/web"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestP2PKeysController_Index_HappyPath(t *testing.T) {
	t.Parallel()

	client, OCRKeyStore := setupP2PKeysControllerTests(t)
	keys, _ := OCRKeyStore.FindEncryptedP2PKeys()

	response, cleanup := client.Get("/v2/keys/p2p")
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, response, http.StatusOK)

	resources := []presenters.P2PKeyResource{}
	err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &resources)
	assert.NoError(t, err)

	require.Len(t, resources, len(keys))

	assert.Equal(t, strconv.Itoa(int(keys[0].ID)), resources[0].ID)
	assert.Equal(t, keys[0].PubKey.String(), resources[0].PubKey)
	assert.Equal(t, keys[0].PeerID.String(), resources[0].PeerID)
}

func TestP2PKeysController_Create_HappyPath(t *testing.T) {
	t.Parallel()

	client, OCRKeyStore := setupP2PKeysControllerTests(t)

	keys, _ := OCRKeyStore.FindEncryptedP2PKeys()
	initialLength := len(keys)

	response, cleanup := client.Post("/v2/keys/p2p", nil)
	t.Cleanup(cleanup)
	cltest.AssertServerResponse(t, response, http.StatusOK)

	keys, _ = OCRKeyStore.FindEncryptedP2PKeys()
	require.Len(t, keys, initialLength+1)

	resource := presenters.P2PKeyResource{}
	err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &resource)
	assert.NoError(t, err)

	lastKeyIndex := len(keys) - 1
	assert.Equal(t, strconv.Itoa(int(keys[lastKeyIndex].ID)), resource.ID)
	assert.Equal(t, keys[lastKeyIndex].PubKey.String(), resource.PubKey)
	assert.Equal(t, keys[lastKeyIndex].PeerID.String(), resource.PeerID)

	var peerID p2pkey.PeerID
	peerID.UnmarshalText([]byte(resource.PeerID))
	_, exists := OCRKeyStore.DecryptedP2PKey(peer.ID(peerID))
	assert.Equal(t, exists, true)
}

func TestP2PKeysController_Delete_InvalidP2PKey(t *testing.T) {
	t.Parallel()

	client, _ := setupP2PKeysControllerTests(t)

	invalidP2PKeyID := "bad_key_id"
	response, cleanup := client.Delete("/v2/keys/p2p/" + invalidP2PKeyID)
	t.Cleanup(cleanup)
	assert.Equal(t, http.StatusUnprocessableEntity, response.StatusCode)
}

func TestP2PKeysController_Delete_NonExistentP2PKeyID(t *testing.T) {
	t.Parallel()

	client, _ := setupP2PKeysControllerTests(t)

	nonExistentP2PKeyID := "1234567890"
	response, cleanup := client.Delete("/v2/keys/p2p/" + nonExistentP2PKeyID)
	t.Cleanup(cleanup)
	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}

func TestP2PKeysController_Delete_HappyPath(t *testing.T) {
	t.Parallel()

	client, OCRKeyStore := setupP2PKeysControllerTests(t)
	require.NoError(t, OCRKeyStore.Unlock(cltest.Password))

	keys, _ := OCRKeyStore.FindEncryptedP2PKeys()
	initialLength := len(keys)
	_, encryptedKeyBundle, _ := OCRKeyStore.GenerateEncryptedP2PKey()

	response, cleanup := client.Delete(fmt.Sprintf("/v2/keys/p2p/%d", encryptedKeyBundle.ID))
	t.Cleanup(cleanup)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Error(t, utils.JustError(OCRKeyStore.FindEncryptedP2PKeyByID(encryptedKeyBundle.ID)))

	keys, _ = OCRKeyStore.FindEncryptedP2PKeys()
	assert.Equal(t, initialLength, len(keys))
}

func setupP2PKeysControllerTests(t *testing.T) (cltest.HTTPClientCleaner, *keystore.OCR) {
	t.Helper()

	app, cleanup := cltest.NewApplication(t)
	t.Cleanup(cleanup)
	require.NoError(t, app.Start())

	client := app.NewHTTPClient()

	OCRKeyStore := app.GetKeyStore().OCR
	return client, OCRKeyStore
}
