package web_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	peer "github.com/libp2p/go-libp2p-peer"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/offchainreporting"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/models/p2pkey"
	"github.com/smartcontractkit/chainlink/core/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestP2PKeysController_Index_HappyPath(t *testing.T) {
	client, OCRKeyStore, cleanup := setupP2PKeysControllerTests(t)
	defer cleanup()

	p2pKeys := []p2pkey.EncryptedP2PKey{}

	keys, _ := OCRKeyStore.FindEncryptedP2PKeys()

	response, cleanup := client.Get("/v2/p2p_keys")
	defer cleanup()
	cltest.AssertServerResponse(t, response, http.StatusOK)

	err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &p2pKeys)
	assert.NoError(t, err)

	require.Len(t, p2pKeys, len(keys))

	assert.Equal(t, keys[0].ID, p2pKeys[0].ID)
	assert.Equal(t, keys[0].PubKey, p2pKeys[0].PubKey)
	assert.Equal(t, keys[0].PeerID, p2pKeys[0].PeerID)
}


func TestP2PKeysController_Create_InvalidBody(t *testing.T) {
	client, _, cleanup := setupP2PKeysControllerTests(t)
	defer cleanup()

	invalidBody, _ := json.Marshal(struct {
		BadParam string
	}{
		BadParam: "randomString",
	})
	response, cleanup := client.Post("/v2/p2p_keys", bytes.NewBuffer(invalidBody))
	defer cleanup()
	cltest.AssertServerResponse(t, response, http.StatusUnprocessableEntity)
}

func TestP2PKeysController_Create_HappyPath(t *testing.T) {
	client, OCRKeyStore, cleanup := setupP2PKeysControllerTests(t)
	defer cleanup()

	keys, _ := OCRKeyStore.FindEncryptedP2PKeys()
	initialLength := len(keys)

	body, _ := json.Marshal(models.CreateP2PKeysRequest{
		Password: cltest.Password,
	})
	response, cleanup := client.Post("/v2/p2p_keys", bytes.NewBuffer(body))
	defer cleanup()
	cltest.AssertServerResponse(t, response, http.StatusOK)

	keys, _ = OCRKeyStore.FindEncryptedP2PKeys()
	require.Len(t, keys, initialLength+1)

	encryptedP2PKey := p2pkey.EncryptedP2PKey{}
	err := web.ParseJSONAPIResponse(cltest.ParseResponseBody(t, response), &encryptedP2PKey)
	assert.NoError(t, err)

	lastKeyIndex := len(keys) - 1
	assert.Equal(t, keys[lastKeyIndex].ID, encryptedP2PKey.ID)
	assert.Equal(t, keys[lastKeyIndex].PubKey, encryptedP2PKey.PubKey)
	assert.Equal(t, keys[lastKeyIndex].PeerID, encryptedP2PKey.PeerID)

	_, exists := OCRKeyStore.DecryptedP2PKey(peer.ID(encryptedP2PKey.PeerID))
	assert.Equal(t, exists, true)
}

func setupP2PKeysControllerTests(t *testing.T) (cltest.HTTPClientCleaner, *offchainreporting.KeyStore, func()) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t, cltest.LenientEthMock)
	require.NoError(t, app.Start())
	client := app.NewHTTPClient()

	OCRKeyStore := app.GetStore().OCRKeyStore
	return client, OCRKeyStore, cleanup
}
