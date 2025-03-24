package cmd_test

import (
	"bytes"
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"

	"github.com/smartcontractkit/chainlink-common/pkg/utils"
	"github.com/smartcontractkit/chainlink/v2/core/cmd"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func TestP2PKeyPresenter_RenderTable(t *testing.T) {
	t.Parallel()

	var (
		id     = "1"
		peerID = configtest.DefaultPeerID
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

func TestShell_ListP2PKeys(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	app := startNewApplicationV2(t, nil)
	key, err := app.GetKeyStore().P2P().Create(ctx)
	require.NoError(t, err)

	requireP2PKeyCount(t, app, 1)

	client, r := app.NewShellAndRenderer()

	assert.NoError(t, client.ListP2PKeys(cltest.EmptyCLIContext()))
	require.Len(t, r.Renders, 1)
	keys := *r.Renders[0].(*cmd.P2PKeyPresenters)
	assert.Equal(t, key.PublicKeyHex(), keys[0].PubKey)
}

func TestShell_CreateP2PKey(t *testing.T) {
	t.Parallel()

	app := startNewApplicationV2(t, nil)
	client, _ := app.NewShellAndRenderer()

	require.NoError(t, client.CreateP2PKey(nilContext))

	keys, err := app.GetKeyStore().P2P().GetAll()
	require.NoError(t, err)

	require.Len(t, keys, 1)
}

func TestShell_DeleteP2PKey(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	app := startNewApplicationV2(t, nil)
	client, _ := app.NewShellAndRenderer()

	key, err := app.GetKeyStore().P2P().Create(ctx)
	require.NoError(t, err)

	requireP2PKeyCount(t, app, 1)

	set := flag.NewFlagSet("test", 0)
	flagSetApplyFromAction(client.DeleteP2PKey, set, "")

	require.NoError(t, set.Set("yes", "true"))

	strID := key.ID()
	err = set.Parse([]string{strID})
	require.NoError(t, err)
	c := cli.NewContext(nil, set, nil)
	err = client.DeleteP2PKey(c)
	require.NoError(t, err)

	requireP2PKeyCount(t, app, 0)
}

func TestShell_ImportExportP2PKeyBundle(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	defer deleteKeyExportFile(t)

	app := startNewApplicationV2(t, nil)
	client, _ := app.NewShellAndRenderer()
	_, err := app.GetKeyStore().P2P().Create(ctx)
	require.NoError(t, err)

	keys := requireP2PKeyCount(t, app, 1)
	key := keys[0]
	keyName := keyNameForTest(t)

	// Export test invalid id
	set := flag.NewFlagSet("test P2P export", 0)
	flagSetApplyFromAction(client.ExportP2PKey, set, "")

	require.NoError(t, set.Parse([]string{"0"}))
	require.NoError(t, set.Set("new-password", "../internal/fixtures/incorrect_password.txt"))
	require.NoError(t, set.Set("output", keyName))

	c := cli.NewContext(nil, set, nil)
	err = client.ExportP2PKey(c)
	require.Error(t, err, "Error exporting")
	require.Error(t, utils.JustError(os.Stat(keyName)))

	// Export test
	set = flag.NewFlagSet("test P2P export", 0)
	flagSetApplyFromAction(client.ExportP2PKey, set, "")

	require.NoError(t, set.Parse([]string{key.ID()}))
	require.NoError(t, set.Set("new-password", "../internal/fixtures/incorrect_password.txt"))
	require.NoError(t, set.Set("output", keyName))

	c = cli.NewContext(nil, set, nil)

	require.NoError(t, client.ExportP2PKey(c))
	require.NoError(t, utils.JustError(os.Stat(keyName)))

	require.NoError(t, utils.JustError(app.GetKeyStore().P2P().Delete(ctx, key.PeerID())))
	requireP2PKeyCount(t, app, 0)

	set = flag.NewFlagSet("test P2P import", 0)
	flagSetApplyFromAction(client.ImportP2PKey, set, "")

	require.NoError(t, set.Parse([]string{keyName}))
	require.NoError(t, set.Set("old-password", "../internal/fixtures/incorrect_password.txt"))

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
