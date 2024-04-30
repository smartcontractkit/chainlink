package cmd_test

import (
	"bytes"
	"context"
	"flag"
	"fmt"
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
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/dkgencryptkey"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func TestDKGEncryptKeyPresenter_RenderTable(t *testing.T) {
	t.Parallel()

	var (
		id     = "1"
		pubKey = "somepubkey"
		buffer = bytes.NewBufferString("")
		r      = cmd.RendererTable{Writer: buffer}
	)

	p := cmd.DKGEncryptKeyPresenter{
		JAID: cmd.JAID{ID: id},
		DKGEncryptKeyResource: presenters.DKGEncryptKeyResource{
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
	ps := cmd.DKGEncryptKeyPresenters{p}
	require.NoError(t, ps.RenderTable(r))

	output = buffer.String()
	assert.Contains(t, output, id)
	assert.Contains(t, output, pubKey)
}

func TestShell_DKGEncryptKeys(t *testing.T) {
	app := startNewApplicationV2(t, nil)
	ks := app.GetKeyStore().DKGEncrypt()
	cleanup := func() {
		ctx := context.Background()
		keys, err := ks.GetAll()
		assert.NoError(t, err)
		for _, key := range keys {
			assert.NoError(t, utils.JustError(ks.Delete(ctx, key.ID())))
		}
		requireDKGEncryptKeyCount(t, app, 0)
	}

	t.Run("ListDKGEncryptKeys", func(tt *testing.T) {
		defer cleanup()
		ctx := testutils.Context(t)
		client, r := app.NewShellAndRenderer()
		key, err := app.GetKeyStore().DKGEncrypt().Create(ctx)
		assert.NoError(tt, err)
		requireDKGEncryptKeyCount(t, app, 1)
		assert.Nil(t, cmd.NewDKGEncryptKeysClient(client).ListKeys(cltest.EmptyCLIContext()))
		assert.Equal(t, 1, len(r.Renders))
		keys := *r.Renders[0].(*cmd.DKGEncryptKeyPresenters)
		assert.True(t, key.PublicKeyString() == keys[0].PublicKey)
	})

	t.Run("CreateDKGEncryptKey", func(tt *testing.T) {
		defer cleanup()
		client, _ := app.NewShellAndRenderer()
		assert.NoError(tt, cmd.NewDKGEncryptKeysClient(client).CreateKey(nilContext))
		keys, err := app.GetKeyStore().DKGEncrypt().GetAll()
		assert.NoError(tt, err)
		assert.Len(t, keys, 1)
	})

	t.Run("DeleteDKGEncryptKey", func(tt *testing.T) {
		defer cleanup()
		ctx := testutils.Context(t)
		client, _ := app.NewShellAndRenderer()
		key, err := app.GetKeyStore().DKGEncrypt().Create(ctx)
		assert.NoError(tt, err)
		requireDKGEncryptKeyCount(tt, app, 1)
		set := flag.NewFlagSet("test", 0)
		flagSetApplyFromAction(cmd.NewDKGEncryptKeysClient(client).DeleteKey, set, "")

		require.NoError(tt, set.Set("yes", "true"))

		strID := key.ID()
		err = set.Parse([]string{strID})
		require.NoError(t, err)
		c := cli.NewContext(nil, set, nil)
		err = cmd.NewDKGEncryptKeysClient(client).DeleteKey(c)
		assert.NoError(tt, err)
		requireDKGEncryptKeyCount(tt, app, 0)
	})

	t.Run("ImportExportDKGEncryptKey", func(tt *testing.T) {
		defer cleanup()
		defer deleteKeyExportFile(tt)
		ctx := testutils.Context(t)
		client, _ := app.NewShellAndRenderer()

		_, err := app.GetKeyStore().DKGEncrypt().Create(ctx)
		require.NoError(tt, err)

		keys := requireDKGEncryptKeyCount(tt, app, 1)
		key := keys[0]
		t.Log("key id:", key.ID())
		keyName := keyNameForTest(t)

		// Export test invalid id
		set := flag.NewFlagSet("test DKGEncrypt export", 0)
		flagSetApplyFromAction(cmd.NewDKGEncryptKeysClient(client).ExportKey, set, "")

		require.NoError(tt, set.Parse([]string{"0"}))
		require.NoError(tt, set.Set("new-password", "../internal/fixtures/incorrect_password.txt"))
		require.NoError(tt, set.Set("output", keyName))

		c := cli.NewContext(nil, set, nil)
		err = cmd.NewDKGEncryptKeysClient(client).ExportKey(c)
		require.Error(tt, err, "Error exporting")
		require.Error(tt, utils.JustError(os.Stat(keyName)))

		// Export test
		set = flag.NewFlagSet("test DKGEncrypt export", 0)
		flagSetApplyFromAction(cmd.NewDKGEncryptKeysClient(client).ExportKey, set, "")

		require.NoError(tt, set.Parse([]string{fmt.Sprint(key.ID())}))
		require.NoError(tt, set.Set("new-password", "../internal/fixtures/incorrect_password.txt"))
		require.NoError(tt, set.Set("output", keyName))

		c = cli.NewContext(nil, set, nil)

		require.NoError(tt, cmd.NewDKGEncryptKeysClient(client).ExportKey(c))
		require.NoError(tt, utils.JustError(os.Stat(keyName)))

		require.NoError(tt, utils.JustError(app.GetKeyStore().DKGEncrypt().Delete(ctx, key.ID())))
		requireDKGEncryptKeyCount(tt, app, 0)

		//Import test
		set = flag.NewFlagSet("test DKGEncrypt import", 0)
		flagSetApplyFromAction(cmd.NewDKGEncryptKeysClient(client).ImportKey, set, "")

		require.NoError(tt, set.Parse([]string{keyName}))
		require.NoError(tt, set.Set("old-password", "../internal/fixtures/incorrect_password.txt"))

		c = cli.NewContext(nil, set, nil)
		require.NoError(tt, cmd.NewDKGEncryptKeysClient(client).ImportKey(c))

		requireDKGEncryptKeyCount(tt, app, 1)
	})
}

func requireDKGEncryptKeyCount(t *testing.T, app chainlink.Application, length int) []dkgencryptkey.Key {
	t.Helper()
	keys, err := app.GetKeyStore().DKGEncrypt().GetAll()
	require.NoError(t, err)
	require.Len(t, keys, length)
	return keys
}
