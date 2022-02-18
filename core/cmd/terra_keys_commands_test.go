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
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/terrakey"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
)

func TestTerraKeyPresenter_RenderTable(t *testing.T) {
	t.Parallel()

	var (
		id     = "1"
		pubKey = "somepubkey"
		buffer = bytes.NewBufferString("")
		r      = cmd.RendererTable{Writer: buffer}
	)

	p := cmd.TerraKeyPresenter{
		JAID: cmd.JAID{ID: id},
		TerraKeyResource: presenters.TerraKeyResource{
			JAID:   presenters.NewJAID(id),
			PubKey: pubKey,
		},
	}

	// Render a single resource
	require.NoError(t, p.RenderTable(r))

	output := buffer.String()
	assert.Contains(t, output, id)
	assert.Contains(t, output, pubKey)

	// Render many resources
	buffer.Reset()
	ps := cmd.TerraKeyPresenters{p}
	require.NoError(t, ps.RenderTable(r))

	output = buffer.String()
	assert.Contains(t, output, id)
	assert.Contains(t, output, pubKey)
}

func TestClient_TerraKeys(t *testing.T) {
	app := startNewApplication(t)
	ks := app.GetKeyStore().Terra()
	cleanup := func() {
		keys, err := ks.GetAll()
		require.NoError(t, err)
		for _, key := range keys {
			require.NoError(t, utils.JustError(ks.Delete(key.ID())))
		}
		requireTerraKeyCount(t, app, 0)
	}

	t.Run("ListTerraKeys", func(tt *testing.T) {
		defer cleanup()
		client, r := app.NewClientAndRenderer()
		key, err := app.GetKeyStore().Terra().Create()
		require.NoError(t, err)
		requireTerraKeyCount(t, app, 1)
		assert.Nil(t, client.ListTerraKeys(cltest.EmptyCLIContext()))
		require.Equal(t, 1, len(r.Renders))
		keys := *r.Renders[0].(*cmd.TerraKeyPresenters)
		assert.True(t, key.PublicKeyStr() == keys[0].PubKey)

	})

	t.Run("CreateTerraKey", func(tt *testing.T) {
		defer cleanup()
		client, _ := app.NewClientAndRenderer()
		require.NoError(t, client.CreateTerraKey(nilContext))
		keys, err := app.GetKeyStore().Terra().GetAll()
		require.NoError(t, err)
		require.Len(t, keys, 1)
	})

	t.Run("DeleteTerraKey", func(tt *testing.T) {
		defer cleanup()
		client, _ := app.NewClientAndRenderer()
		key, err := app.GetKeyStore().Terra().Create()
		require.NoError(t, err)
		requireTerraKeyCount(t, app, 1)
		set := flag.NewFlagSet("test", 0)
		set.Bool("yes", true, "")
		strID := key.ID()
		set.Parse([]string{strID})
		c := cli.NewContext(nil, set, nil)
		err = client.DeleteTerraKey(c)
		require.NoError(t, err)
		requireTerraKeyCount(t, app, 0)
	})

	t.Run("ImportExportTerraKey", func(tt *testing.T) {
		defer cleanup()
		defer deleteKeyExportFile(t)
		client, _ := app.NewClientAndRenderer()

		_, err := app.GetKeyStore().Terra().Create()
		require.NoError(t, err)

		keys := requireTerraKeyCount(t, app, 1)
		key := keys[0]
		keyName := keyNameForTest(t)

		// Export test invalid id
		set := flag.NewFlagSet("test Terra export", 0)
		set.Parse([]string{"0"})
		set.String("newpassword", "../internal/fixtures/incorrect_password.txt", "")
		set.String("output", keyName, "")
		c := cli.NewContext(nil, set, nil)
		err = client.ExportTerraKey(c)
		require.Error(t, err, "Error exporting")
		require.Error(t, utils.JustError(os.Stat(keyName)))

		// Export test
		set = flag.NewFlagSet("test Terra export", 0)
		set.Parse([]string{fmt.Sprint(key.ID())})
		set.String("newpassword", "../internal/fixtures/incorrect_password.txt", "")
		set.String("output", keyName, "")
		c = cli.NewContext(nil, set, nil)

		require.NoError(t, client.ExportTerraKey(c))
		require.NoError(t, utils.JustError(os.Stat(keyName)))

		require.NoError(t, utils.JustError(app.GetKeyStore().Terra().Delete(key.ID())))
		requireTerraKeyCount(t, app, 0)

		set = flag.NewFlagSet("test Terra import", 0)
		set.Parse([]string{keyName})
		set.String("oldpassword", "../internal/fixtures/incorrect_password.txt", "")
		c = cli.NewContext(nil, set, nil)
		require.NoError(t, client.ImportTerraKey(c))

		requireTerraKeyCount(t, app, 1)
	})
}

func requireTerraKeyCount(t *testing.T, app chainlink.Application, length int) []terrakey.Key {
	t.Helper()
	keys, err := app.GetKeyStore().Terra().GetAll()
	require.NoError(t, err)
	require.Len(t, keys, length)
	return keys
}
