package cmd_test

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"

	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/csakey"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
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

func TestClient_ListCSAKeys(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	key, err := app.GetKeyStore().CSA().Create()
	require.NoError(t, err)

	requireCSAKeyCount(t, app, 1)

	client, r := app.NewClientAndRenderer()

	assert.Nil(t, client.ListCSAKeys(cltest.EmptyCLIContext()))
	require.Equal(t, 1, len(r.Renders))
	keys := *r.Renders[0].(*cmd.CSAKeyPresenters)
	assert.Equal(t, key.PublicKeyString(), keys[0].PubKey)
}

func TestClient_CreateCSAKey(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, _ := app.NewClientAndRenderer()

	requireCSAKeyCount(t, app, 0)

	require.NoError(t, client.CreateCSAKey(nilContext))

	requireCSAKeyCount(t, app, 1)
}

func TestClient_ImportExportCsaKey(t *testing.T) {
	t.Parallel()

	defer deleteKeyExportFile(t)

	app := startNewApplication(t)
	client, _ := app.NewClientAndRenderer()
	_, err := app.GetKeyStore().CSA().Create()
	require.NoError(t, err)

	keys := requireCSAKeyCount(t, app, 1)
	key := keys[0]
	keyName := keyNameForTest(t)

	// Export test invalid id
	set := flag.NewFlagSet("test CSA export", 0)
	set.Parse([]string{"0"})
	set.String("newpassword", "../internal/fixtures/incorrect_password.txt", "")
	set.String("output", keyName, "")
	c := cli.NewContext(nil, set, nil)
	err = client.ExportCSAKey(c)
	require.Error(t, err, "Error exporting")
	require.Error(t, utils.JustError(os.Stat(keyName)))

	// Export test
	set = flag.NewFlagSet("test CSA export", 0)
	set.Parse([]string{fmt.Sprint(key.ID())})
	set.String("newpassword", "../internal/fixtures/incorrect_password.txt", "")
	set.String("output", keyName, "")
	c = cli.NewContext(nil, set, nil)

	require.NoError(t, client.ExportCSAKey(c))
	require.NoError(t, utils.JustError(os.Stat(keyName)))

	require.NoError(t, utils.JustError(app.GetKeyStore().CSA().Delete(key.ID())))
	requireCSAKeyCount(t, app, 0)

	set = flag.NewFlagSet("test CSA import", 0)
	set.Parse([]string{keyName})
	set.String("oldpassword", "../internal/fixtures/incorrect_password.txt", "")
	c = cli.NewContext(nil, set, nil)
	require.NoError(t, client.ImportCSAKey(c))

	requireCSAKeyCount(t, app, 1)
}

func requireCSAKeyCount(t *testing.T, app chainlink.Application, length int) []csakey.KeyV2 {
	t.Helper()

	keys, err := app.GetKeyStore().CSA().GetAll()
	require.NoError(t, err)
	require.Equal(t, length, len(keys))
	return keys
}
