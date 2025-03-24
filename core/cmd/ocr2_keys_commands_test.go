package cmd_test

import (
	"bytes"
	"context"
	"encoding/hex"
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
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func TestOCR2KeyBundlePresenter_RenderTable(t *testing.T) {
	t.Parallel()

	var (
		bundleID = "f5bf259689b26f1374efb3c9a9868796953a0f814bb2d39b968d0e61b58620a5"
		buffer   = bytes.NewBufferString("")
		r        = cmd.RendererTable{Writer: buffer}
	)

	key := cltest.DefaultOCR2Key
	pubKeyConfig := key.ConfigEncryptionPublicKey()
	pubKey := key.OffchainPublicKey()
	p := cmd.OCR2KeyBundlePresenter{
		JAID: cmd.NewJAID(bundleID),
		OCR2KeysBundleResource: presenters.OCR2KeysBundleResource{
			JAID:              presenters.NewJAID(key.ID()),
			ChainType:         "evm",
			OnchainPublicKey:  key.OnChainPublicKey(),
			OffChainPublicKey: hex.EncodeToString(pubKey[:]),
			ConfigPublicKey:   hex.EncodeToString(pubKeyConfig[:]),
		},
	}

	// Render a single resource
	require.NoError(t, p.RenderTable(r))

	output := buffer.String()
	assert.Contains(t, output, bundleID)
	assert.Contains(t, output, key.ChainType())
	assert.Contains(t, output, key.OnChainPublicKey())
	assert.Contains(t, output, hex.EncodeToString(pubKey[:]))
	assert.Contains(t, output, hex.EncodeToString(pubKeyConfig[:]))

	// Render many resources
	buffer.Reset()
	ps := cmd.OCR2KeyBundlePresenters{p}
	require.NoError(t, ps.RenderTable(r))

	output = buffer.String()
	assert.Contains(t, output, bundleID)
	assert.Contains(t, output, key.OnChainPublicKey())
	assert.Contains(t, output, hex.EncodeToString(pubKey[:]))
	pubKeyConfig = key.ConfigEncryptionPublicKey()
	assert.Contains(t, output, hex.EncodeToString(pubKeyConfig[:]))
}

func TestShell_OCR2Keys(t *testing.T) {
	t.Parallel()

	app := startNewApplicationV2(t, nil)
	ks := app.GetKeyStore().OCR2()
	cleanup := func() {
		ctx := context.Background()
		keys, err := app.GetKeyStore().OCR2().GetAll()
		require.NoError(t, err)
		for _, key := range keys {
			require.NoError(t, ks.Delete(ctx, key.ID()))
		}
		requireOCR2KeyCount(t, app, 0)
	}

	t.Run("ListOCR2KeyBundles", func(tt *testing.T) {
		defer cleanup()
		ctx := testutils.Context(t)
		client, r := app.NewShellAndRenderer()

		key, err := app.GetKeyStore().OCR2().Create(ctx, "evm")
		require.NoError(t, err)
		requireOCR2KeyCount(t, app, 1)
		assert.NoError(t, client.ListOCR2KeyBundles(cltest.EmptyCLIContext()))
		require.Len(t, r.Renders, 1)
		output := *r.Renders[0].(*cmd.OCR2KeyBundlePresenters)
		require.Equal(t, key.ID(), output[0].ID)
	})

	t.Run("CreateOCR2KeyBundle", func(tt *testing.T) {
		defer cleanup()
		client, r := app.NewShellAndRenderer()

		set := flag.NewFlagSet("test", 0)
		flagSetApplyFromAction(client.CreateOCR2KeyBundle, set, "")

		require.NoError(tt, set.Parse([]string{"evm"}))

		c := cli.NewContext(nil, set, nil)
		require.NoError(t, client.CreateOCR2KeyBundle(c))
		keys, err := app.GetKeyStore().OCR2().GetAll()
		require.NoError(t, err)
		require.Len(t, keys, 1)
		require.Len(t, r.Renders, 1)
		output := (*r.Renders[0].(*cmd.OCR2KeyBundlePresenter))
		require.Equal(t, output.ID, keys[0].ID())
	})

	t.Run("DeleteOCR2KeyBundle", func(tt *testing.T) {
		defer cleanup()
		ctx := testutils.Context(t)
		client, r := app.NewShellAndRenderer()

		key, err := app.GetKeyStore().OCR2().Create(ctx, "evm")
		require.NoError(t, err)
		requireOCR2KeyCount(t, app, 1)
		set := flag.NewFlagSet("test", 0)
		flagSetApplyFromAction(client.DeleteOCR2KeyBundle, set, "")

		require.NoError(tt, set.Parse([]string{key.ID()}))
		require.NoError(tt, set.Set("yes", "true"))

		c := cli.NewContext(nil, set, nil)
		require.NoError(t, client.DeleteOCR2KeyBundle(c))
		requireOCR2KeyCount(t, app, 0)
		require.Len(t, r.Renders, 1)
		output := *r.Renders[0].(*cmd.OCR2KeyBundlePresenter)
		assert.Equal(t, key.ID(), output.ID)
	})

	t.Run("ImportExportOCR2Key", func(tt *testing.T) {
		defer cleanup()
		defer deleteKeyExportFile(t)
		ctx := testutils.Context(t)
		client, _ := app.NewShellAndRenderer()

		err := app.KeyStore.OCR2().Add(ctx, cltest.DefaultOCR2Key)
		require.NoError(t, err)

		keys := requireOCR2KeyCount(t, app, 1)
		key := keys[0]
		keyName := keyNameForTest(t)

		// Export test invalid id
		set := flag.NewFlagSet("test OCR2 export", 0)
		flagSetApplyFromAction(client.ExportOCR2Key, set, "")

		require.NoError(tt, set.Parse([]string{"0"}))
		require.NoError(tt, set.Set("new-password", "../internal/fixtures/new_password.txt"))
		require.NoError(tt, set.Set("output", keyName))

		c := cli.NewContext(nil, set, nil)
		err = client.ExportOCR2Key(c)
		require.Error(t, err, "Error exporting")
		require.Error(t, utils.JustError(os.Stat(keyName)))

		// Export
		set = flag.NewFlagSet("test OCR2 export", 0)
		flagSetApplyFromAction(client.ExportOCR2Key, set, "")

		require.NoError(tt, set.Parse([]string{key.ID()}))
		require.NoError(tt, set.Set("new-password", "../internal/fixtures/new_password.txt"))
		require.NoError(tt, set.Set("output", keyName))

		c = cli.NewContext(nil, set, nil)

		require.NoError(t, client.ExportOCR2Key(c))
		require.NoError(t, utils.JustError(os.Stat(keyName)))

		require.NoError(t, app.GetKeyStore().OCR2().Delete(ctx, key.ID()))
		requireOCR2KeyCount(t, app, 0)

		set = flag.NewFlagSet("test OCR2 import", 0)
		flagSetApplyFromAction(client.ImportOCR2Key, set, "")

		require.NoError(tt, set.Parse([]string{keyName}))
		require.NoError(tt, set.Set("old-password", "../internal/fixtures/new_password.txt"))

		c = cli.NewContext(nil, set, nil)
		require.NoError(t, client.ImportOCR2Key(c))

		requireOCR2KeyCount(t, app, 1)
	})
}

func requireOCR2KeyCount(t *testing.T, app chainlink.Application, length int) []ocr2key.KeyBundle {
	keys, err := app.GetKeyStore().OCR2().GetAll()
	require.NoError(t, err)
	require.Len(t, keys, length)
	return keys
}
