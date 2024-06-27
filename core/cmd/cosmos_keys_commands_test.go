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
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/cosmoskey"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func TestCosmosKeyPresenter_RenderTable(t *testing.T) {
	t.Parallel()

	var (
		id     = "1"
		pubKey = "somepubkey"
		buffer = bytes.NewBufferString("")
		r      = cmd.RendererTable{Writer: buffer}
	)

	p := cmd.CosmosKeyPresenter{
		JAID: cmd.JAID{ID: id},
		CosmosKeyResource: presenters.CosmosKeyResource{
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
	ps := cmd.CosmosKeyPresenters{p}
	require.NoError(t, ps.RenderTable(r))

	output = buffer.String()
	assert.Contains(t, output, id)
	assert.Contains(t, output, pubKey)
}

func TestShell_CosmosKeys(t *testing.T) {
	app := startNewApplicationV2(t, nil)
	ks := app.GetKeyStore().Cosmos()
	cleanup := func() {
		ctx := context.Background()
		keys, err := ks.GetAll()
		require.NoError(t, err)
		for _, key := range keys {
			require.NoError(t, utils.JustError(ks.Delete(ctx, key.ID())))
		}
		requireCosmosKeyCount(t, app, 0)
	}

	t.Run("ListCosmosKeys", func(tt *testing.T) {
		defer cleanup()
		ctx := testutils.Context(t)
		client, r := app.NewShellAndRenderer()
		key, err := app.GetKeyStore().Cosmos().Create(ctx)
		require.NoError(t, err)
		requireCosmosKeyCount(t, app, 1)
		assert.Nil(t, cmd.NewCosmosKeysClient(client).ListKeys(cltest.EmptyCLIContext()))
		require.Equal(t, 1, len(r.Renders))
		keys := *r.Renders[0].(*cmd.CosmosKeyPresenters)
		assert.True(t, key.PublicKeyStr() == keys[0].PubKey)
	})

	t.Run("CreateCosmosKey", func(tt *testing.T) {
		defer cleanup()
		client, _ := app.NewShellAndRenderer()
		require.NoError(t, cmd.NewCosmosKeysClient(client).CreateKey(nilContext))
		keys, err := app.GetKeyStore().Cosmos().GetAll()
		require.NoError(t, err)
		require.Len(t, keys, 1)
	})

	t.Run("DeleteCosmosKey", func(tt *testing.T) {
		defer cleanup()
		ctx := testutils.Context(t)
		client, _ := app.NewShellAndRenderer()
		key, err := app.GetKeyStore().Cosmos().Create(ctx)
		require.NoError(t, err)
		requireCosmosKeyCount(t, app, 1)
		set := flag.NewFlagSet("test", 0)
		flagSetApplyFromAction(cmd.NewCosmosKeysClient(client).DeleteKey, set, "cosmos")

		strID := key.ID()
		require.NoError(tt, set.Set("yes", "true"))
		require.NoError(tt, set.Parse([]string{strID}))

		c := cli.NewContext(nil, set, nil)
		err = cmd.NewCosmosKeysClient(client).DeleteKey(c)
		require.NoError(t, err)
		requireCosmosKeyCount(t, app, 0)
	})

	t.Run("ImportExportCosmosKey", func(tt *testing.T) {
		defer cleanup()
		defer deleteKeyExportFile(t)
		ctx := testutils.Context(t)
		client, _ := app.NewShellAndRenderer()

		_, err := app.GetKeyStore().Cosmos().Create(ctx)
		require.NoError(t, err)

		keys := requireCosmosKeyCount(t, app, 1)
		key := keys[0]
		keyName := keyNameForTest(t)

		// Export test invalid id
		set := flag.NewFlagSet("test Cosmos export", 0)
		flagSetApplyFromAction(cmd.NewCosmosKeysClient(client).ExportKey, set, "cosmos")

		require.NoError(tt, set.Parse([]string{"0"}))
		require.NoError(tt, set.Set("new-password", "../internal/fixtures/incorrect_password.txt"))
		require.NoError(tt, set.Set("output", keyName))

		c := cli.NewContext(nil, set, nil)
		tclient := cmd.NewCosmosKeysClient(client)
		err = tclient.ExportKey(c)
		require.Error(t, err, "Error exporting")
		require.Error(t, utils.JustError(os.Stat(keyName)))

		// Export test
		set = flag.NewFlagSet("test Cosmos export", 0)
		flagSetApplyFromAction(cmd.NewCosmosKeysClient(client).ExportKey, set, "cosmos")

		require.NoError(tt, set.Parse([]string{fmt.Sprint(key.ID())}))
		require.NoError(tt, set.Set("new-password", "../internal/fixtures/incorrect_password.txt"))
		require.NoError(tt, set.Set("output", keyName))

		c = cli.NewContext(nil, set, nil)

		require.NoError(t, tclient.ExportKey(c))
		require.NoError(t, utils.JustError(os.Stat(keyName)))

		require.NoError(t, utils.JustError(app.GetKeyStore().Cosmos().Delete(ctx, key.ID())))
		requireCosmosKeyCount(t, app, 0)

		set = flag.NewFlagSet("test Cosmos import", 0)
		flagSetApplyFromAction(cmd.NewCosmosKeysClient(client).ImportKey, set, "cosmos")

		require.NoError(tt, set.Parse([]string{keyName}))
		require.NoError(tt, set.Set("old-password", "../internal/fixtures/incorrect_password.txt"))

		c = cli.NewContext(nil, set, nil)
		require.NoError(t, tclient.ImportKey(c))

		requireCosmosKeyCount(t, app, 1)
	})
}

func requireCosmosKeyCount(t *testing.T, app chainlink.Application, length int) []cosmoskey.Key {
	t.Helper()
	keys, err := app.GetKeyStore().Cosmos().GetAll()
	require.NoError(t, err)
	require.Len(t, keys, length)
	return keys
}
