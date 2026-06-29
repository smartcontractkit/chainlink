package cmd_test

import (
	"bytes"
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/stellarkey"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"

	"github.com/smartcontractkit/chainlink/v2/core/cmd"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func TestStellarKeyPresenter_RenderTable(t *testing.T) {
	t.Parallel()

	var (
		id     = "GDTERXGZ7J6NQCSTG7ZNDR7RYV6QD2K4U4XG7K6Y4GD5EQLNMXS6K4C5"
		pubKey = "somepubkey"
		buffer = bytes.NewBufferString("")
		r      = cmd.RendererTable{Writer: buffer}
	)

	p := cmd.StellarKeyPresenter{
		JAID: cmd.JAID{ID: id},
		StellarKeyResource: presenters.StellarKeyResource{
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
	ps := cmd.StellarKeyPresenters{p}
	require.NoError(t, ps.RenderTable(r))

	output = buffer.String()
	assert.Contains(t, output, id)
	assert.Contains(t, output, pubKey)
}

//nolint:paralleltest // subtests share a single keystore/application instance
func TestShell_StellarKeys(t *testing.T) {
	app := startNewApplicationV2(t, nil)
	ks := app.GetKeyStore().Stellar()
	cleanup := func() {
		ctx := t.Context()
		keys, err := ks.GetAll()
		require.NoError(t, err)
		for _, key := range keys {
			require.NoError(t, utils.JustError(ks.Delete(ctx, key.ID())))
		}
		requireStellarKeyCount(t, app, 0)
	}

	//nolint:paralleltest // subtests share a single keystore/application instance
	t.Run("ListStellarKeys", func(tt *testing.T) {
		defer cleanup()
		ctx := tt.Context()
		client, r := app.NewShellAndRenderer()
		key, err := app.GetKeyStore().Stellar().Create(ctx)
		require.NoError(t, err)
		requireStellarKeyCount(t, app, 1)
		require.NoError(t, cmd.NewStellarKeysClient(client).ListKeys(cltest.EmptyCLIContext()))
		require.Len(t, r.Renders, 1)
		keys := *r.Renders[0].(*cmd.StellarKeyPresenters)
		assert.Equal(t, key.PublicKeyStr(), keys[0].PubKey)
	})

	//nolint:paralleltest // subtests share a single keystore/application instance
	t.Run("CreateStellarKey", func(tt *testing.T) {
		defer cleanup()
		client, _ := app.NewShellAndRenderer()
		require.NoError(t, cmd.NewStellarKeysClient(client).CreateKey(nilContext))
		keys, err := app.GetKeyStore().Stellar().GetAll()
		require.NoError(t, err)
		require.Len(t, keys, 1)
	})

	//nolint:paralleltest // subtests share a single keystore/application instance
	t.Run("DeleteStellarKey", func(tt *testing.T) {
		defer cleanup()
		ctx := tt.Context()
		client, _ := app.NewShellAndRenderer()
		key, err := app.GetKeyStore().Stellar().Create(ctx)
		require.NoError(t, err)
		requireStellarKeyCount(t, app, 1)
		set := flag.NewFlagSet("test", 0)
		flagSetApplyFromAction(cmd.NewStellarKeysClient(client).DeleteKey, set, "stellar")

		require.NoError(tt, set.Set("yes", "true"))

		strID := key.ID()
		err = set.Parse([]string{strID})
		require.NoError(t, err)
		c := cli.NewContext(nil, set, nil)
		err = cmd.NewStellarKeysClient(client).DeleteKey(c)
		require.NoError(t, err)
		requireStellarKeyCount(t, app, 0)
	})

	//nolint:paralleltest // subtests share a single keystore/application instance
	t.Run("ImportExportStellarKey", func(tt *testing.T) {
		defer cleanup()
		defer deleteKeyExportFile(t)
		ctx := tt.Context()
		client, _ := app.NewShellAndRenderer()

		_, err := app.GetKeyStore().Stellar().Create(ctx)
		require.NoError(t, err)

		keys := requireStellarKeyCount(t, app, 1)
		key := keys[0]
		keyName := keyNameForTest(t)

		// Export test invalid id
		set := flag.NewFlagSet("test Stellar export", 0)
		flagSetApplyFromAction(cmd.NewStellarKeysClient(client).ExportKey, set, "stellar")

		require.NoError(tt, set.Parse([]string{"0"}))
		require.NoError(tt, set.Set("new-password", "../internal/fixtures/incorrect_password.txt"))
		require.NoError(tt, set.Set("output", keyName))

		c := cli.NewContext(nil, set, nil)
		err = cmd.NewStellarKeysClient(client).ExportKey(c)
		require.Error(t, err, "Error exporting")
		require.Error(t, utils.JustError(os.Stat(keyName)))

		// Export test
		set = flag.NewFlagSet("test Stellar export", 0)
		flagSetApplyFromAction(cmd.NewStellarKeysClient(client).ExportKey, set, "stellar")

		require.NoError(tt, set.Parse([]string{key.ID()}))
		require.NoError(tt, set.Set("new-password", "../internal/fixtures/incorrect_password.txt"))
		require.NoError(tt, set.Set("output", keyName))

		c = cli.NewContext(nil, set, nil)

		require.NoError(t, cmd.NewStellarKeysClient(client).ExportKey(c))
		require.NoError(t, utils.JustError(os.Stat(keyName)))

		require.NoError(t, utils.JustError(app.GetKeyStore().Stellar().Delete(ctx, key.ID())))
		requireStellarKeyCount(t, app, 0)

		set = flag.NewFlagSet("test Stellar import", 0)
		flagSetApplyFromAction(cmd.NewStellarKeysClient(client).ImportKey, set, "stellar")

		require.NoError(tt, set.Parse([]string{keyName}))
		require.NoError(tt, set.Set("old-password", "../internal/fixtures/incorrect_password.txt"))
		c = cli.NewContext(nil, set, nil)
		require.NoError(t, cmd.NewStellarKeysClient(client).ImportKey(c))

		requireStellarKeyCount(t, app, 1)
	})
}

func requireStellarKeyCount(t *testing.T, app chainlink.Application, length int) []stellarkey.Key {
	t.Helper()
	keys, err := app.GetKeyStore().Stellar().GetAll()
	require.NoError(t, err)
	require.Len(t, keys, length)
	return keys
}
