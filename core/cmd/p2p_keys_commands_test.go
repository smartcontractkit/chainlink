package cmd_test

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"testing"

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
		id     = "1"
		peerID = "12D3KooWPjceQrSwdWXPyLLeABRXmuqt69Rg3sBYbU1Nft9HyQ6X"
		pubKey = "somepubkey"
		buffer = bytes.NewBufferString("")
		r      = cmd.RendererTable{Writer: buffer}
	)

	p := cmd.P2PKeyPresenter{
		JAID: cmd.JAID{ID: id},
		P2PKeyResource: presenters.P2PKeyResource{
			JAID:   presenters.NewJAID(id),
			PeerID: peerID,
			PubKey: pubKey,
		},
	}

	// Render a single resource
	require.NoError(t, p.RenderTable(r))

	output := buffer.String()
	assert.Contains(t, output, id)
	assert.Contains(t, output, peerID)
	assert.Contains(t, output, pubKey)

	// Render many resources
	buffer.Reset()
	ps := cmd.P2PKeyPresenters{p}
	require.NoError(t, ps.RenderTable(r))

	output = buffer.String()
	assert.Contains(t, output, id)
	assert.Contains(t, output, peerID)
	assert.Contains(t, output, pubKey)
}

func TestClient_ListP2PKeys(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	key, err := app.GetKeyStore().P2P().Create()
	require.NoError(t, err)

	requireP2PKeyCount(t, app, 1)

	client, r := app.NewClientAndRenderer()

	assert.Nil(t, client.ListP2PKeys(cltest.EmptyCLIContext()))
	require.Equal(t, 1, len(r.Renders))
	keys := *r.Renders[0].(*cmd.P2PKeyPresenters)
	assert.True(t, key.PublicKeyHex() == keys[0].PubKey)
}

func TestClient_CreateP2PKey(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, _ := app.NewClientAndRenderer()

	require.NoError(t, client.CreateP2PKey(nilContext))

	keys, err := app.GetKeyStore().P2P().GetAll()
	require.NoError(t, err)

	require.Len(t, keys, 1)
}

func TestClient_DeleteP2PKey(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, _ := app.NewClientAndRenderer()

	key, err := app.GetKeyStore().P2P().Create()
	require.NoError(t, err)

	requireP2PKeyCount(t, app, 1)

	set := flag.NewFlagSet("test", 0)
	set.Bool("yes", true, "")
	strID := key.ID()
	set.Parse([]string{strID})
	c := cli.NewContext(nil, set, nil)
	err = client.DeleteP2PKey(c)
	require.NoError(t, err)

	requireP2PKeyCount(t, app, 0)
}

func TestClient_ImportExportP2PKeyBundle(t *testing.T) {
	t.Parallel()

	defer deleteKeyExportFile(t)

	app := startNewApplication(t)
	client, _ := app.NewClientAndRenderer()
	_, err := app.GetKeyStore().P2P().Create()
	require.NoError(t, err)

	keys := requireP2PKeyCount(t, app, 1)
	key := keys[0]
	keyName := keyNameForTest(t)

	// Export test invalid id
	set := flag.NewFlagSet("test P2P export", 0)
	set.Parse([]string{"0"})
	set.String("newpassword", "../internal/fixtures/incorrect_password.txt", "")
	set.String("output", keyName, "")
	c := cli.NewContext(nil, set, nil)
	err = client.ExportP2PKey(c)
	require.Error(t, err, "Error exporting")
	require.Error(t, utils.JustError(os.Stat(keyName)))

	// Export test
	set = flag.NewFlagSet("test P2P export", 0)
	set.Parse([]string{fmt.Sprint(key.ID())})
	set.String("newpassword", "../internal/fixtures/incorrect_password.txt", "")
	set.String("output", keyName, "")
	c = cli.NewContext(nil, set, nil)

	require.NoError(t, client.ExportP2PKey(c))
	require.NoError(t, utils.JustError(os.Stat(keyName)))

	require.NoError(t, utils.JustError(app.GetKeyStore().P2P().Delete(key.PeerID())))
	requireP2PKeyCount(t, app, 0)

	set = flag.NewFlagSet("test P2P import", 0)
	set.Parse([]string{keyName})
	set.String("oldpassword", "../internal/fixtures/incorrect_password.txt", "")
	c = cli.NewContext(nil, set, nil)
	require.NoError(t, client.ImportP2PKey(c))

	requireP2PKeyCount(t, app, 1)
}

func requireP2PKeyCount(t *testing.T, app chainlink.Application, length int) []p2pkey.KeyV2 {
	t.Helper()

	keys, err := app.GetKeyStore().P2P().GetAll()
	require.NoError(t, err)
	require.Len(t, keys, length)
	return keys
}
