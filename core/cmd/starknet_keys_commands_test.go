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

	starkkey "github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/keys"

	"github.com/smartcontractkit/chainlink/v2/core/cmd"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func TestStarkNetKeyPresenter_RenderTable(t *testing.T) {
	t.Parallel()

	var (
		id          = "1"
		accountAddr = "someaccountaddress"
		starkKey    = "somestarkkey"
		buffer      = bytes.NewBufferString("")
		r           = cmd.RendererTable{Writer: buffer}
	)

	p := cmd.StarkNetKeyPresenter{
		JAID: cmd.JAID{ID: id},
		StarkNetKeyResource: presenters.StarkNetKeyResource{
			JAID:        presenters.NewJAID(id),
			AccountAddr: accountAddr,
			StarkKey:    starkKey,
		},
	}

	// Render a single resource
	require.NoError(t, p.RenderTable(r))

	output := buffer.String()
	assert.Contains(t, output, id)
	assert.Contains(t, output, accountAddr)

	// Render many resources
	buffer.Reset()
	ps := cmd.StarkNetKeyPresenters{p}
	require.NoError(t, ps.RenderTable(r))

	output = buffer.String()
	assert.Contains(t, output, id)
	assert.Contains(t, output, accountAddr)
	assert.Contains(t, output, starkKey)
}

func TestClient_StarkNetKeys(t *testing.T) {
	app := startNewApplicationV2(t, nil)
	ks := app.GetKeyStore().StarkNet()
	cleanup := func() {
		keys, err := ks.GetAll()
		require.NoError(t, err)
		for _, key := range keys {
			require.NoError(t, utils.JustError(ks.Delete(key.ID())))
		}
		requireStarkNetKeyCount(t, app, 0)
	}

	t.Run("ListStarkNetKeys", func(tt *testing.T) {
		defer cleanup()
		client, r := app.NewClientAndRenderer()
		key, err := app.GetKeyStore().StarkNet().Create()
		require.NoError(t, err)
		requireStarkNetKeyCount(t, app, 1)
		assert.Nil(t, cmd.NewStarkNetKeysClient(client).ListKeys(cltest.EmptyCLIContext()))
		require.Equal(t, 1, len(r.Renders))
		keys := *r.Renders[0].(*cmd.StarkNetKeyPresenters)
		assert.True(t, key.AccountAddressStr() == keys[0].AccountAddr)
		assert.True(t, key.StarkKeyStr() == keys[0].StarkKey)

	})

	t.Run("CreateStarkNetKey", func(tt *testing.T) {
		defer cleanup()
		client, _ := app.NewClientAndRenderer()
		require.NoError(t, cmd.NewStarkNetKeysClient(client).CreateKey(nilContext))
		keys, err := app.GetKeyStore().StarkNet().GetAll()
		require.NoError(t, err)
		require.Len(t, keys, 1)
	})

	t.Run("DeleteStarkNetKey", func(tt *testing.T) {
		defer cleanup()
		client, _ := app.NewClientAndRenderer()
		key, err := app.GetKeyStore().StarkNet().Create()
		require.NoError(t, err)
		requireStarkNetKeyCount(t, app, 1)
		set := flag.NewFlagSet("test", 0)
		cltest.FlagSetApplyFromAction(cmd.NewStarkNetKeysClient(client).DeleteKey, set, "starknet")

		require.NoError(tt, set.Set("yes", "true"))

		strID := key.ID()
		set.Parse([]string{strID})
		c := cli.NewContext(nil, set, nil)
		err = cmd.NewStarkNetKeysClient(client).DeleteKey(c)
		require.NoError(t, err)
		requireStarkNetKeyCount(t, app, 0)
	})

	t.Run("ImportExportStarkNetKey", func(tt *testing.T) {
		defer cleanup()
		defer deleteKeyExportFile(t)
		client, _ := app.NewClientAndRenderer()

		_, err := app.GetKeyStore().StarkNet().Create()
		require.NoError(t, err)

		keys := requireStarkNetKeyCount(t, app, 1)
		key := keys[0]
		keyName := keyNameForTest(t)

		// Export test invalid id
		set := flag.NewFlagSet("test StarkNet export", 0)
		cltest.FlagSetApplyFromAction(cmd.NewStarkNetKeysClient(client).ExportKey, set, "starknet")

		require.NoError(tt, set.Parse([]string{"0"}))
		require.NoError(tt, set.Set("new-password", "../internal/fixtures/incorrect_password.txt"))
		require.NoError(tt, set.Set("output", keyName))

		c := cli.NewContext(nil, set, nil)
		err = cmd.NewStarkNetKeysClient(client).ExportKey(c)
		require.Error(t, err, "Error exporting")
		require.Error(t, utils.JustError(os.Stat(keyName)))

		// Export test
		set = flag.NewFlagSet("test StarkNet export", 0)
		cltest.FlagSetApplyFromAction(cmd.NewStarkNetKeysClient(client).ExportKey, set, "starknet")

		require.NoError(tt, set.Parse([]string{fmt.Sprint(key.ID())}))
		require.NoError(tt, set.Set("new-password", "../internal/fixtures/incorrect_password.txt"))
		require.NoError(tt, set.Set("output", keyName))

		c = cli.NewContext(nil, set, nil)

		require.NoError(t, cmd.NewStarkNetKeysClient(client).ExportKey(c))
		require.NoError(t, utils.JustError(os.Stat(keyName)))

		require.NoError(t, utils.JustError(app.GetKeyStore().StarkNet().Delete(key.ID())))
		requireStarkNetKeyCount(t, app, 0)

		set = flag.NewFlagSet("test StarkNet import", 0)
		cltest.FlagSetApplyFromAction(cmd.NewStarkNetKeysClient(client).ImportKey, set, "starknet")

		require.NoError(tt, set.Parse([]string{keyName}))
		require.NoError(tt, set.Set("old-password", "../internal/fixtures/incorrect_password.txt"))

		c = cli.NewContext(nil, set, nil)
		require.NoError(t, cmd.NewStarkNetKeysClient(client).ImportKey(c))

		requireStarkNetKeyCount(t, app, 1)
	})
}

func requireStarkNetKeyCount(t *testing.T, app chainlink.Application, length int) []starkkey.Key {
	t.Helper()
	keys, err := app.GetKeyStore().StarkNet().GetAll()
	require.NoError(t, err)
	require.Len(t, keys, length)
	return keys
}
