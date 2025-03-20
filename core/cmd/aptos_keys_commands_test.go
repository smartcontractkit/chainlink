package cmd_test

import (
	"bytes"
	"context"
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
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/aptoskey"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func TestAptosKeyPresenter_RenderTable(t *testing.T) {
	t.Parallel()

	var (
		id     = "1"
		pubKey = "somepubkey"
		buffer = bytes.NewBufferString("")
		r      = cmd.RendererTable{Writer: buffer}
	)

	p := cmd.AptosKeyPresenter{
		JAID: cmd.JAID{ID: id},
		AptosKeyResource: presenters.AptosKeyResource{
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
	ps := cmd.AptosKeyPresenters{p}
	require.NoError(t, ps.RenderTable(r))

	output = buffer.String()
	assert.Contains(t, output, id)
	assert.Contains(t, output, pubKey)
}

func TestShell_AptosKeys(t *testing.T) {
	app := startNewApplicationV2(t, nil)
	ks := app.GetKeyStore().Aptos()
	cleanup := func() {
		ctx := context.Background()
		keys, err := ks.GetAll()
		require.NoError(t, err)
		for _, key := range keys {
			require.NoError(t, utils.JustError(ks.Delete(ctx, key.ID())))
		}
		requireAptosKeyCount(t, app, 0)
	}

	t.Run("ListAptosKeys", func(tt *testing.T) {
		defer cleanup()
		ctx := testutils.Context(t)
		client, r := app.NewShellAndRenderer()
		key, err := app.GetKeyStore().Aptos().Create(ctx)
		require.NoError(t, err)
		requireAptosKeyCount(t, app, 1)
		assert.NoError(t, cmd.NewAptosKeysClient(client).ListKeys(cltest.EmptyCLIContext()))
		require.Len(t, r.Renders, 1)
		keys := *r.Renders[0].(*cmd.AptosKeyPresenters)
		assert.Equal(t, key.PublicKeyStr(), keys[0].PubKey)
	})

	t.Run("CreateAptosKey", func(tt *testing.T) {
		defer cleanup()
		client, _ := app.NewShellAndRenderer()
		require.NoError(t, cmd.NewAptosKeysClient(client).CreateKey(nilContext))
		keys, err := app.GetKeyStore().Aptos().GetAll()
		require.NoError(t, err)
		require.Len(t, keys, 1)
	})

	t.Run("DeleteAptosKey", func(tt *testing.T) {
		defer cleanup()
		ctx := testutils.Context(t)
		client, _ := app.NewShellAndRenderer()
		key, err := app.GetKeyStore().Aptos().Create(ctx)
		require.NoError(t, err)
		requireAptosKeyCount(t, app, 1)
		set := flag.NewFlagSet("test", 0)
		flagSetApplyFromAction(cmd.NewAptosKeysClient(client).DeleteKey, set, "aptos")

		require.NoError(tt, set.Set("yes", "true"))

		strID := key.ID()
		err = set.Parse([]string{strID})
		require.NoError(t, err)
		c := cli.NewContext(nil, set, nil)
		err = cmd.NewAptosKeysClient(client).DeleteKey(c)
		require.NoError(t, err)
		requireAptosKeyCount(t, app, 0)
	})

	t.Run("ImportExportAptosKey", func(tt *testing.T) {
		defer cleanup()
		defer deleteKeyExportFile(t)
		ctx := testutils.Context(t)
		client, _ := app.NewShellAndRenderer()

		_, err := app.GetKeyStore().Aptos().Create(ctx)
		require.NoError(t, err)

		keys := requireAptosKeyCount(t, app, 1)
		key := keys[0]
		keyName := keyNameForTest(t)

		// Export test invalid id
		set := flag.NewFlagSet("test Aptos export", 0)
		flagSetApplyFromAction(cmd.NewAptosKeysClient(client).ExportKey, set, "aptos")

		require.NoError(tt, set.Parse([]string{"0"}))
		require.NoError(tt, set.Set("new-password", "../internal/fixtures/incorrect_password.txt"))
		require.NoError(tt, set.Set("output", keyName))

		c := cli.NewContext(nil, set, nil)
		err = cmd.NewAptosKeysClient(client).ExportKey(c)
		require.Error(t, err, "Error exporting")
		require.Error(t, utils.JustError(os.Stat(keyName)))

		// Export test
		set = flag.NewFlagSet("test Aptos export", 0)
		flagSetApplyFromAction(cmd.NewAptosKeysClient(client).ExportKey, set, "aptos")

		require.NoError(tt, set.Parse([]string{key.ID()}))
		require.NoError(tt, set.Set("new-password", "../internal/fixtures/incorrect_password.txt"))
		require.NoError(tt, set.Set("output", keyName))

		c = cli.NewContext(nil, set, nil)

		require.NoError(t, cmd.NewAptosKeysClient(client).ExportKey(c))
		require.NoError(t, utils.JustError(os.Stat(keyName)))

		require.NoError(t, utils.JustError(app.GetKeyStore().Aptos().Delete(ctx, key.ID())))
		requireAptosKeyCount(t, app, 0)

		set = flag.NewFlagSet("test Aptos import", 0)
		flagSetApplyFromAction(cmd.NewAptosKeysClient(client).ImportKey, set, "aptos")

		require.NoError(tt, set.Parse([]string{keyName}))
		require.NoError(tt, set.Set("old-password", "../internal/fixtures/incorrect_password.txt"))
		c = cli.NewContext(nil, set, nil)
		require.NoError(t, cmd.NewAptosKeysClient(client).ImportKey(c))

		requireAptosKeyCount(t, app, 1)
	})
}

func requireAptosKeyCount(t *testing.T, app chainlink.Application, length int) []aptoskey.Key {
	t.Helper()
	keys, err := app.GetKeyStore().Aptos().GetAll()
	require.NoError(t, err)
	require.Len(t, keys, length)
	return keys
}
