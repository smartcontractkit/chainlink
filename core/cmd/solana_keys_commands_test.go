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
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/solkey"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func TestSolanaKeyPresenter_RenderTable(t *testing.T) {
	t.Parallel()

	var (
		id     = "1"
		pubKey = "somepubkey"
		buffer = bytes.NewBufferString("")
		r      = cmd.RendererTable{Writer: buffer}
	)

	p := cmd.SolanaKeyPresenter{
		JAID: cmd.JAID{ID: id},
		SolanaKeyResource: presenters.SolanaKeyResource{
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
	ps := cmd.SolanaKeyPresenters{p}
	require.NoError(t, ps.RenderTable(r))

	output = buffer.String()
	assert.Contains(t, output, id)
	assert.Contains(t, output, pubKey)
}

func TestShell_SolanaKeys(t *testing.T) {
	app := startNewApplicationV2(t, nil)
	ks := app.GetKeyStore().Solana()
	cleanup := func() {
		ctx := context.Background()
		keys, err := ks.GetAll()
		require.NoError(t, err)
		for _, key := range keys {
			require.NoError(t, utils.JustError(ks.Delete(ctx, key.ID())))
		}
		requireSolanaKeyCount(t, app, 0)
	}

	t.Run("ListSolanaKeys", func(tt *testing.T) {
		defer cleanup()
		ctx := testutils.Context(t)
		client, r := app.NewShellAndRenderer()
		key, err := app.GetKeyStore().Solana().Create(ctx)
		require.NoError(t, err)
		requireSolanaKeyCount(t, app, 1)
		assert.Nil(t, cmd.NewSolanaKeysClient(client).ListKeys(cltest.EmptyCLIContext()))
		require.Equal(t, 1, len(r.Renders))
		keys := *r.Renders[0].(*cmd.SolanaKeyPresenters)
		assert.True(t, key.PublicKeyStr() == keys[0].PubKey)
	})

	t.Run("CreateSolanaKey", func(tt *testing.T) {
		defer cleanup()
		client, _ := app.NewShellAndRenderer()
		require.NoError(t, cmd.NewSolanaKeysClient(client).CreateKey(nilContext))
		keys, err := app.GetKeyStore().Solana().GetAll()
		require.NoError(t, err)
		require.Len(t, keys, 1)
	})

	t.Run("DeleteSolanaKey", func(tt *testing.T) {
		defer cleanup()
		ctx := testutils.Context(t)
		client, _ := app.NewShellAndRenderer()
		key, err := app.GetKeyStore().Solana().Create(ctx)
		require.NoError(t, err)
		requireSolanaKeyCount(t, app, 1)
		set := flag.NewFlagSet("test", 0)
		flagSetApplyFromAction(cmd.NewSolanaKeysClient(client).DeleteKey, set, "solana")

		require.NoError(tt, set.Set("yes", "true"))

		strID := key.ID()
		err = set.Parse([]string{strID})
		require.NoError(t, err)
		c := cli.NewContext(nil, set, nil)
		err = cmd.NewSolanaKeysClient(client).DeleteKey(c)
		require.NoError(t, err)
		requireSolanaKeyCount(t, app, 0)
	})

	t.Run("ImportExportSolanaKey", func(tt *testing.T) {
		defer cleanup()
		defer deleteKeyExportFile(t)
		ctx := testutils.Context(t)
		client, _ := app.NewShellAndRenderer()

		_, err := app.GetKeyStore().Solana().Create(ctx)
		require.NoError(t, err)

		keys := requireSolanaKeyCount(t, app, 1)
		key := keys[0]
		keyName := keyNameForTest(t)

		// Export test invalid id
		set := flag.NewFlagSet("test Solana export", 0)
		flagSetApplyFromAction(cmd.NewSolanaKeysClient(client).ExportKey, set, "solana")

		require.NoError(tt, set.Parse([]string{"0"}))
		require.NoError(tt, set.Set("new-password", "../internal/fixtures/incorrect_password.txt"))
		require.NoError(tt, set.Set("output", keyName))

		c := cli.NewContext(nil, set, nil)
		err = cmd.NewSolanaKeysClient(client).ExportKey(c)
		require.Error(t, err, "Error exporting")
		require.Error(t, utils.JustError(os.Stat(keyName)))

		// Export test
		set = flag.NewFlagSet("test Solana export", 0)
		flagSetApplyFromAction(cmd.NewSolanaKeysClient(client).ExportKey, set, "solana")

		require.NoError(tt, set.Parse([]string{fmt.Sprint(key.ID())}))
		require.NoError(tt, set.Set("new-password", "../internal/fixtures/incorrect_password.txt"))
		require.NoError(tt, set.Set("output", keyName))

		c = cli.NewContext(nil, set, nil)

		require.NoError(t, cmd.NewSolanaKeysClient(client).ExportKey(c))
		require.NoError(t, utils.JustError(os.Stat(keyName)))

		require.NoError(t, utils.JustError(app.GetKeyStore().Solana().Delete(ctx, key.ID())))
		requireSolanaKeyCount(t, app, 0)

		set = flag.NewFlagSet("test Solana import", 0)
		flagSetApplyFromAction(cmd.NewSolanaKeysClient(client).ImportKey, set, "solana")

		require.NoError(tt, set.Parse([]string{keyName}))
		require.NoError(tt, set.Set("old-password", "../internal/fixtures/incorrect_password.txt"))
		c = cli.NewContext(nil, set, nil)
		require.NoError(t, cmd.NewSolanaKeysClient(client).ImportKey(c))

		requireSolanaKeyCount(t, app, 1)
	})
}

func requireSolanaKeyCount(t *testing.T, app chainlink.Application, length int) []solkey.Key {
	t.Helper()
	keys, err := app.GetKeyStore().Solana().GetAll()
	require.NoError(t, err)
	require.Len(t, keys, length)
	return keys
}
