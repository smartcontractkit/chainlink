package cmd_test

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
)

func TestP2PKeyPresenter_RenderTable(t *testing.T) {
	t.Parallel()

	var (
		id        = "1"
		peerID    = "12D3KooWApUJaQB2saFjyEUfq6BmysnsSnhLnY5CF9tURYVKgoXK"
		pubKey    = "somepubkey"
		createdAt = time.Now()
		updatedAt = time.Now().Add(time.Second)
		deletedAt = time.Now().Add(2 * time.Second)
		buffer    = bytes.NewBufferString("")
		r         = cmd.RendererTable{Writer: buffer}
	)

	p := cmd.P2PKeyPresenter{
		JAID: cmd.JAID{ID: id},
		P2PKeyResource: presenters.P2PKeyResource{
			JAID:      presenters.NewJAID(id),
			PeerID:    peerID,
			PubKey:    pubKey,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
			DeletedAt: &deletedAt,
		},
	}

	// Render a single resource
	require.NoError(t, p.RenderTable(r))

	output := buffer.String()
	assert.Contains(t, output, id)
	assert.Contains(t, output, peerID)
	assert.Contains(t, output, pubKey)
	assert.Contains(t, output, createdAt.String())
	assert.Contains(t, output, updatedAt.String())
	assert.Contains(t, output, deletedAt.String())

	// Render many resources
	buffer.Reset()
	ps := cmd.P2PKeyPresenters{p}
	require.NoError(t, ps.RenderTable(r))

	output = buffer.String()
	assert.Contains(t, output, id)
	assert.Contains(t, output, peerID)
	assert.Contains(t, output, pubKey)
	assert.Contains(t, output, createdAt.String())
	assert.Contains(t, output, updatedAt.String())
	assert.Contains(t, output, deletedAt.String())
}

func TestClient_ListP2PKeys(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	app.GetKeyStore().OCR.Unlock(cltest.Password)

	key, err := p2pkey.CreateKey()
	require.NoError(t, err)
	encKey, err := key.ToEncryptedP2PKey(cltest.Password, utils.FastScryptParams)
	require.NoError(t, err)
	err = app.GetKeyStore().OCR.UpsertEncryptedP2PKey(&encKey)
	require.NoError(t, err)

	requireP2PKeyCount(t, app, 2) // Created  + fixture key

	client, r := app.NewClientAndRenderer()

	assert.Nil(t, client.ListP2PKeys(cltest.EmptyCLIContext()))
	require.Equal(t, 1, len(r.Renders))
	keys := *r.Renders[0].(*cmd.P2PKeyPresenters)
	assert.Equal(t, encKey.PubKey.String(), keys[1].PubKey)
}

func TestClient_CreateP2PKey(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, _ := app.NewClientAndRenderer()

	app.GetKeyStore().OCR.Unlock(cltest.Password)

	require.NoError(t, client.CreateP2PKey(nilContext))

	keys, err := app.GetKeyStore().OCR.FindEncryptedP2PKeys()
	require.NoError(t, err)

	// Created + fixture key
	require.Len(t, keys, 2)

	for _, e := range keys {
		_, err = e.Decrypt(cltest.Password)
		require.NoError(t, err)
	}
}

func TestClient_DeleteP2PKey(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, _ := app.NewClientAndRenderer()

	app.GetKeyStore().OCR.Unlock(cltest.Password)

	key, err := p2pkey.CreateKey()
	require.NoError(t, err)
	encKey, err := key.ToEncryptedP2PKey(cltest.Password, utils.FastScryptParams)
	require.NoError(t, err)
	err = app.GetKeyStore().OCR.UpsertEncryptedP2PKey(&encKey)
	require.NoError(t, err)

	requireP2PKeyCount(t, app, 2) // Created  + fixture key

	set := flag.NewFlagSet("test", 0)
	set.Bool("yes", true, "")
	strID := strconv.FormatInt(int64(encKey.ID), 10)
	set.Parse([]string{strID})
	c := cli.NewContext(nil, set, nil)
	err = client.DeleteP2PKey(c)
	require.NoError(t, err)

	requireP2PKeyCount(t, app, 1) // fixture key only
}

func TestClient_ImportExportP2PKeyBundle(t *testing.T) {
	t.Parallel()

	defer deleteKeyExportFile(t)

	app := startNewApplication(t)
	client, _ := app.NewClientAndRenderer()

	app.GetKeyStore().OCR.Unlock(cltest.Password)

	keys := requireP2PKeyCount(t, app, 1)
	key := keys[0]
	keyName := keyNameForTest(t)

	// Export test invalid id
	set := flag.NewFlagSet("test P2P export", 0)
	set.Parse([]string{"0"})
	set.String("newpassword", "../internal/fixtures/apicredentials", "")
	set.String("output", keyName, "")
	c := cli.NewContext(nil, set, nil)
	err := client.ExportP2PKey(c)
	require.Error(t, err, "Error exporting")
	require.Error(t, utils.JustError(os.Stat(keyName)))

	// Export test
	set = flag.NewFlagSet("test P2P export", 0)
	set.Parse([]string{fmt.Sprint(key.ID)})
	set.String("newpassword", "../internal/fixtures/apicredentials", "")
	set.String("output", keyName, "")
	c = cli.NewContext(nil, set, nil)

	require.NoError(t, client.ExportP2PKey(c))
	require.NoError(t, utils.JustError(os.Stat(keyName)))

	require.NoError(t, app.GetKeyStore().OCR.DeleteEncryptedP2PKey(&key))
	requireP2PKeyCount(t, app, 0)

	set = flag.NewFlagSet("test P2P import", 0)
	set.Parse([]string{keyName})
	set.String("oldpassword", "../internal/fixtures/apicredentials", "")
	c = cli.NewContext(nil, set, nil)
	require.NoError(t, client.ImportP2PKey(c))

	requireP2PKeyCount(t, app, 1)
}

func requireP2PKeyCount(t *testing.T, app chainlink.Application, length int) []p2pkey.EncryptedP2PKey {
	t.Helper()

	keys, err := app.GetKeyStore().OCR.FindEncryptedP2PKeys()
	require.NoError(t, err)
	require.Len(t, keys, length)
	return keys
}
