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
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/csakey"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func TestCSAKeyPresenter_RenderTable(t *testing.T) {
	t.Parallel()

	var (
		pubKey = "somepubkey"
		buffer = bytes.NewBufferString("")
		r      = cmd.RendererTable{Writer: buffer}
	)

	p := cmd.CSAKeyPresenter{
		JAID: cmd.JAID{ID: pubKey},
		CSAKeyResource: presenters.CSAKeyResource{
			JAID:   presenters.NewJAID(pubKey),
			PubKey: pubKey,
		},
	}

	// Render a single resource
	require.NoError(t, p.RenderTable(r))

	output := buffer.String()
	assert.Contains(t, output, pubKey)

	// Render many resources
	buffer.Reset()
	ps := cmd.CSAKeyPresenters{p}
	require.NoError(t, ps.RenderTable(r))

	output = buffer.String()
	assert.Contains(t, output, pubKey)
}

func TestShell_ListCSAKeys(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	app := startNewApplicationV2(t, nil)
	key, err := keystore.GetDefault(ctx, app.GetKeyStore().CSA())
	require.NoError(t, err)

	requireCSAKeyCount(t, app, 1)

	client, r := app.NewShellAndRenderer()

	assert.NoError(t, client.ListCSAKeys(cltest.EmptyCLIContext()))
	require.Len(t, r.Renders, 1)
	keys := *r.Renders[0].(*cmd.CSAKeyPresenters)
	assert.Equal(t, "csa_"+key.PublicKeyString(), keys[0].PubKey)
}

func TestShell_CreateCSAKey(t *testing.T) {
	t.Parallel()

	app := startNewApplicationV2(t, nil)
	client, _ := app.NewShellAndRenderer()

	requireCSAKeyCount(t, app, 1)

	require.Error(t, client.CreateCSAKey(nilContext))

	requireCSAKeyCount(t, app, 1)
}

func TestShell_ImportExportCsaKey(t *testing.T) {
	t.Parallel()

	defer deleteKeyExportFile(t)
	ctx := testutils.Context(t)

	app := startNewApplicationV2(t, nil)

	client, _ := app.NewShellAndRenderer()
	err := app.GetKeyStore().CSA().EnsureKey(ctx)
	require.NoError(t, err)

	keys := requireCSAKeyCount(t, app, 1)
	key := keys[0]
	keyName := keyNameForTest(t)

	// Export test invalid id
	set := flag.NewFlagSet("test CSA export", 0)
	flagSetApplyFromAction(client.ExportCSAKey, set, "")

	require.NoError(t, set.Parse([]string{"0"}))
	require.NoError(t, set.Set("new-password", "../internal/fixtures/incorrect_password.txt"))
	require.NoError(t, set.Set("output", keyName))

	c := cli.NewContext(nil, set, nil)
	err = client.ExportCSAKey(c)
	require.Error(t, err, "Error exporting")
	require.Error(t, utils.JustError(os.Stat(keyName)))

	// Export test
	set = flag.NewFlagSet("test CSA export", 0)
	flagSetApplyFromAction(client.ExportCSAKey, set, "")

	require.NoError(t, set.Parse([]string{key.ID()}))
	require.NoError(t, set.Set("new-password", "../internal/fixtures/incorrect_password.txt"))
	require.NoError(t, set.Set("output", keyName))

	c = cli.NewContext(nil, set, nil)

	require.NoError(t, client.ExportCSAKey(c))
	require.NoError(t, utils.JustError(os.Stat(keyName)))

	require.NoError(t, utils.JustError(app.GetKeyStore().CSA().Delete(ctx, key.ID())))
	requireCSAKeyCount(t, app, 0)

	// Import test
	set = flag.NewFlagSet("test CSA import", 0)
	flagSetApplyFromAction(client.ImportCSAKey, set, "")

	require.NoError(t, set.Parse([]string{keyName}))
	require.NoError(t, set.Set("old-password", "../internal/fixtures/incorrect_password.txt"))

	c = cli.NewContext(nil, set, nil)
	require.NoError(t, client.ImportCSAKey(c))

	requireCSAKeyCount(t, app, 1)
}

func requireCSAKeyCount(t *testing.T, app chainlink.Application, length int) []csakey.KeyV2 {
	t.Helper()

	keys, err := app.GetKeyStore().CSA().GetAll()
	require.NoError(t, err)
	require.Len(t, keys, length)
	return keys
}
