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
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/dkgsignkey"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

func TestDKGSignKeyPresenter_RenderTable(t *testing.T) {
	t.Parallel()

	var (
		id     = "1"
		pubKey = "somepubkey"
		buffer = bytes.NewBufferString("")
		r      = cmd.RendererTable{Writer: buffer}
	)

	p := cmd.DKGSignKeyPresenter{
		JAID: cmd.JAID{ID: id},
		DKGSignKeyResource: presenters.DKGSignKeyResource{
			JAID:      presenters.NewJAID(id),
			PublicKey: pubKey,
		},
	}

	// Render a single resource
	require.NoError(t, p.RenderTable(r))

	output := buffer.String()
	assert.Contains(t, output, id)
	assert.Contains(t, output, pubKey)

	// Render many resources
	buffer.Reset()
	ps := cmd.DKGSignKeyPresenters{p}
	require.NoError(t, ps.RenderTable(r))

	output = buffer.String()
	assert.Contains(t, output, id)
	assert.Contains(t, output, pubKey)
}

func TestClient_DKGSignKeys(t *testing.T) {
	app := startNewApplicationV2(t, nil)
	ks := app.GetKeyStore().DKGSign()
	cleanup := func() {
		keys, err := ks.GetAll()
		assert.NoError(t, err)
		for _, key := range keys {
			assert.NoError(t, utils.JustError(ks.Delete(key.ID())))
		}
		requireDKGSignKeyCount(t, app, 0)
	}

	t.Run("ListDKGSignKeys", func(tt *testing.T) {
		defer cleanup()
		client, r := app.NewClientAndRenderer()
		key, err := app.GetKeyStore().DKGSign().Create()
		assert.NoError(tt, err)
		requireDKGSignKeyCount(t, app, 1)
		assert.Nil(t, cmd.NewDKGSignKeysClient(client).ListKeys(cltest.EmptyCLIContext()))
		assert.Equal(t, 1, len(r.Renders))
		keys := *r.Renders[0].(*cmd.DKGSignKeyPresenters)
		assert.True(t, key.PublicKeyString() == keys[0].PublicKey)
	})

	t.Run("CreateDKGSignKey", func(tt *testing.T) {
		defer cleanup()
		client, _ := app.NewClientAndRenderer()
		assert.NoError(tt, cmd.NewDKGSignKeysClient(client).CreateKey(nilContext))
		keys, err := app.GetKeyStore().DKGSign().GetAll()
		assert.NoError(tt, err)
		assert.Len(t, keys, 1)
	})

	t.Run("DeleteDKGSignKey", func(tt *testing.T) {
		defer cleanup()
		client, _ := app.NewClientAndRenderer()
		key, err := app.GetKeyStore().DKGSign().Create()
		assert.NoError(tt, err)
		requireDKGSignKeyCount(tt, app, 1)
		set := flag.NewFlagSet("test", 0)
		cltest.FlagSetApplyFromAction(cmd.NewDKGSignKeysClient(client).DeleteKey, set, "")

		require.NoError(tt, set.Set("yes", "true"))
		strID := key.ID()
		set.Parse([]string{strID})
		c := cli.NewContext(nil, set, nil)
		err = cmd.NewDKGSignKeysClient(client).DeleteKey(c)
		assert.NoError(tt, err)
		requireDKGSignKeyCount(tt, app, 0)
	})

	t.Run("ImportExportDKGSignKey", func(tt *testing.T) {
		defer cleanup()
		defer deleteKeyExportFile(tt)
		client, _ := app.NewClientAndRenderer()

		_, err := app.GetKeyStore().DKGSign().Create()
		require.NoError(tt, err)

		keys := requireDKGSignKeyCount(tt, app, 1)
		key := keys[0]
		t.Log("key id:", key.ID())
		keyName := keyNameForTest(t)

		// Export test invalid id
		set := flag.NewFlagSet("test DKGSign export", 0)
		cltest.FlagSetApplyFromAction(cmd.NewDKGSignKeysClient(client).ExportKey, set, "")

		require.NoError(tt, set.Parse([]string{"0"}))
		require.NoError(tt, set.Set("newpassword", "../internal/fixtures/incorrect_password.txt"))
		require.NoError(tt, set.Set("output", keyName))

		c := cli.NewContext(nil, set, nil)
		err = cmd.NewDKGSignKeysClient(client).ExportKey(c)
		require.Error(tt, err)
		require.Error(tt, utils.JustError(os.Stat(keyName)))

		// Export test
		set = flag.NewFlagSet("test DKGSign export", 0)
		cltest.FlagSetApplyFromAction(cmd.NewDKGSignKeysClient(client).ExportKey, set, "")

		require.NoError(tt, set.Parse([]string{fmt.Sprint(key.ID())}))
		require.NoError(tt, set.Set("newpassword", "../internal/fixtures/incorrect_password.txt"))
		require.NoError(tt, set.Set("output", keyName))

		c = cli.NewContext(nil, set, nil)

		require.NoError(tt, cmd.NewDKGSignKeysClient(client).ExportKey(c))
		require.NoError(tt, utils.JustError(os.Stat(keyName)))

		require.NoError(tt, utils.JustError(app.GetKeyStore().DKGSign().Delete(key.ID())))
		requireDKGSignKeyCount(tt, app, 0)

		set = flag.NewFlagSet("test DKGSign import", 0)
		cltest.FlagSetApplyFromAction(cmd.NewDKGSignKeysClient(client).ImportKey, set, "")

		require.NoError(tt, set.Parse([]string{keyName}))
		require.NoError(tt, set.Set("oldpassword", "../internal/fixtures/incorrect_password.txt"))

		c = cli.NewContext(nil, set, nil)
		require.NoError(tt, cmd.NewDKGSignKeysClient(client).ImportKey(c))

		requireDKGSignKeyCount(tt, app, 1)
	})
}

func requireDKGSignKeyCount(t *testing.T, app chainlink.Application, length int) []dkgsignkey.Key {
	t.Helper()
	keys, err := app.GetKeyStore().DKGSign().GetAll()
	require.NoError(t, err)
	require.Len(t, keys, length)
	return keys
}
