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
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/terrakey"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
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
	app := startNewApplicationV2(t, nil)
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
		assert.Nil(t, cmd.NewTerraKeysClient(client).ListKeys(cltest.EmptyCLIContext()))
		require.Equal(t, 1, len(r.Renders))
		keys := *r.Renders[0].(*cmd.TerraKeyPresenters)
		assert.True(t, key.PublicKeyStr() == keys[0].PubKey)

	})

	t.Run("CreateTerraKey", func(tt *testing.T) {
		defer cleanup()
		client, _ := app.NewClientAndRenderer()
		require.NoError(t, cmd.NewTerraKeysClient(client).CreateKey(nilContext))
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
		cltest.CopyFlagSetFromAction(cmd.NewTerraKeysClient(client).DeleteKey, set, "terra")

		strID := key.ID()
		require.NoError(tt, set.Set("yes", "true"))
		require.NoError(tt, set.Parse([]string{strID}))

		c := cli.NewContext(nil, set, nil)
		err = cmd.NewTerraKeysClient(client).DeleteKey(c)
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
		cltest.CopyFlagSetFromAction(cmd.NewTerraKeysClient(client).ExportKey, set, "terra")

		require.NoError(tt, set.Parse([]string{"0"}))
		require.NoError(tt, set.Set("newpassword", "../internal/fixtures/incorrect_password.txt"))
		require.NoError(tt, set.Set("output", keyName))

		c := cli.NewContext(nil, set, nil)
		tclient := cmd.NewTerraKeysClient(client)
		err = tclient.ExportKey(c)
		require.Error(t, err, "Error exporting")
		require.Error(t, utils.JustError(os.Stat(keyName)))

		// Export test
		set = flag.NewFlagSet("test Terra export", 0)
		cltest.CopyFlagSetFromAction(cmd.NewTerraKeysClient(client).ExportKey, set, "terra")

		require.NoError(tt, set.Parse([]string{fmt.Sprint(key.ID())}))
		require.NoError(tt, set.Set("newpassword", "../internal/fixtures/incorrect_password.txt"))
		require.NoError(tt, set.Set("output", keyName))

		c = cli.NewContext(nil, set, nil)

		require.NoError(t, tclient.ExportKey(c))
		require.NoError(t, utils.JustError(os.Stat(keyName)))

		require.NoError(t, utils.JustError(app.GetKeyStore().Terra().Delete(key.ID())))
		requireTerraKeyCount(t, app, 0)

		set = flag.NewFlagSet("test Terra import", 0)
		cltest.CopyFlagSetFromAction(cmd.NewTerraKeysClient(client).ImportKey, set, "terra")

		require.NoError(tt, set.Parse([]string{keyName}))
		require.NoError(tt, set.Set("oldpassword", "../internal/fixtures/incorrect_password.txt"))

		c = cli.NewContext(nil, set, nil)
		require.NoError(t, tclient.ImportKey(c))

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
