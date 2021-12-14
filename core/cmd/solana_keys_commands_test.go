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
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/solkey"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
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

func TestClient_SolanaKeys(t *testing.T) {
	app := startNewApplication(t)
	ks := app.GetKeyStore().Solana()
	cleanup := func() {
		keys, err := ks.GetAll()
		require.NoError(t, err)
		for _, key := range keys {
			require.NoError(t, utils.JustError(ks.Delete(key.ID())))
		}
		requireSolanaKeyCount(t, app, 0)
	}

	t.Run("ListSolanaKeys", func(tt *testing.T) {
		defer cleanup()
		client, r := app.NewClientAndRenderer()
		key, err := app.GetKeyStore().Solana().Create()
		require.NoError(t, err)
		requireSolanaKeyCount(t, app, 1)
		assert.Nil(t, client.ListSolanaKeys(cltest.EmptyCLIContext()))
		require.Equal(t, 1, len(r.Renders))
		keys := *r.Renders[0].(*cmd.SolanaKeyPresenters)
		assert.True(t, key.PublicKeyStr() == keys[0].PubKey)

	})

	t.Run("CreateSolanaKey", func(tt *testing.T) {
		defer cleanup()
		client, _ := app.NewClientAndRenderer()
		require.NoError(t, client.CreateSolanaKey(nilContext))
		keys, err := app.GetKeyStore().Solana().GetAll()
		require.NoError(t, err)
		require.Len(t, keys, 1)
	})

	t.Run("DeleteSolanaKey", func(tt *testing.T) {
		defer cleanup()
		client, _ := app.NewClientAndRenderer()
		key, err := app.GetKeyStore().Solana().Create()
		require.NoError(t, err)
		requireSolanaKeyCount(t, app, 1)
		set := flag.NewFlagSet("test", 0)
		set.Bool("yes", true, "")
		strID := key.ID()
		set.Parse([]string{strID})
		c := cli.NewContext(nil, set, nil)
		err = client.DeleteSolanaKey(c)
		require.NoError(t, err)
		requireSolanaKeyCount(t, app, 0)
	})

	t.Run("ImportExportSolanaKey", func(tt *testing.T) {
		defer cleanup()
		defer deleteKeyExportFile(t)
		client, _ := app.NewClientAndRenderer()

		_, err := app.GetKeyStore().Solana().Create()
		require.NoError(t, err)

		keys := requireSolanaKeyCount(t, app, 1)
		key := keys[0]
		keyName := keyNameForTest(t)

		// Export test invalid id
		set := flag.NewFlagSet("test Solana export", 0)
		set.Parse([]string{"0"})
		set.String("newpassword", "../internal/fixtures/incorrect_password.txt", "")
		set.String("output", keyName, "")
		c := cli.NewContext(nil, set, nil)
		err = client.ExportSolanaKey(c)
		require.Error(t, err, "Error exporting")
		require.Error(t, utils.JustError(os.Stat(keyName)))

		// Export test
		set = flag.NewFlagSet("test Solana export", 0)
		set.Parse([]string{fmt.Sprint(key.ID())})
		set.String("newpassword", "../internal/fixtures/incorrect_password.txt", "")
		set.String("output", keyName, "")
		c = cli.NewContext(nil, set, nil)

		require.NoError(t, client.ExportSolanaKey(c))
		require.NoError(t, utils.JustError(os.Stat(keyName)))

		require.NoError(t, utils.JustError(app.GetKeyStore().Solana().Delete(key.ID())))
		requireSolanaKeyCount(t, app, 0)

		set = flag.NewFlagSet("test Solana import", 0)
		set.Parse([]string{keyName})
		set.String("oldpassword", "../internal/fixtures/incorrect_password.txt", "")
		c = cli.NewContext(nil, set, nil)
		require.NoError(t, client.ImportSolanaKey(c))

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
