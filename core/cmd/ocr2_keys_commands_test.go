package cmd_test

import (
	"bytes"
	"encoding/hex"
	"flag"
	"os"
	"testing"

	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
	"gopkg.in/guregu/null.v4"
)

func TestOCR2KeyBundlePresenter_RenderTable(t *testing.T) {
	t.Parallel()

	var (
		bundleID = "f5bf259689b26f1374efb3c9a9868796953a0f814bb2d39b968d0e61b58620a5"
		buffer   = bytes.NewBufferString("")
		r        = cmd.RendererTable{Writer: buffer}
	)

	key := cltest.DefaultOCR2Key
	pubKeyConfig := key.PublicKeyConfig()

	p := cmd.OCR2KeyBundlePresenter{
		JAID: cmd.JAID{ID: bundleID},
		OCR2KeysBundleResource: presenters.OCR2KeysBundleResource{
			JAID:                  presenters.NewJAID(key.ID()),
			ChainType:             "evm",
			OnChainSigningAddress: key.PublicKeyAddressOnChain(),
			OffChainPublicKey:     hex.EncodeToString(key.PublicKeyOffChain()),
			ConfigPublicKey:       hex.EncodeToString(pubKeyConfig[:]),
		},
	}

	// Render a single resource
	require.NoError(t, p.RenderTable(r))

	output := buffer.String()
	assert.Contains(t, output, bundleID)
	assert.Contains(t, output, key.ChainType)
	assert.Contains(t, output, key.PublicKeyAddressOnChain())
	assert.Contains(t, output, hex.EncodeToString(key.PublicKeyOffChain()[:]))
	assert.Contains(t, output, hex.EncodeToString(pubKeyConfig[:]))

	// Render many resources
	buffer.Reset()
	ps := cmd.OCR2KeyBundlePresenters{p}
	require.NoError(t, ps.RenderTable(r))

	output = buffer.String()
	assert.Contains(t, output, bundleID)
	assert.Contains(t, output, key.PublicKeyAddressOnChain())
	assert.Contains(t, output, hex.EncodeToString(key.PublicKeyOffChain()[:]))
	pubKeyConfig = key.PublicKeyConfig()
	assert.Contains(t, output, hex.EncodeToString(pubKeyConfig[:]))
}

func TestClient_ListOCR2KeyBundles(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, r := app.NewClientAndRenderer()

	key, err := app.GetKeyStore().OCR2().Create("evm")
	require.NoError(t, err)

	requireOCR2KeyCount(t, app, 1)

	assert.Nil(t, client.ListOCR2KeyBundles(cltest.EmptyCLIContext()))
	require.Equal(t, 1, len(r.Renders))
	output := *r.Renders[0].(*cmd.OCR2KeyBundlePresenters)
	require.Equal(t, key.ID(), output[0].ID)
}

func TestClient_CreateOCR2KeyBundle(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, r := app.NewClientAndRenderer()

	requireOCR2KeyCount(t, app, 0)

	set := flag.NewFlagSet("test", 0)
	set.Parse([]string{"evm"})
	set.Bool("yes", true, "")
	c := cli.NewContext(nil, set, nil)

	require.NoError(t, client.CreateOCR2KeyBundle(c))

	keys, err := app.GetKeyStore().OCR2().GetAll()
	require.NoError(t, err)
	require.Len(t, keys, 1)

	require.Equal(t, 1, len(r.Renders))
	output := *r.Renders[0].(*cmd.OCR2KeyBundlePresenter)
	require.Equal(t, output.ID, keys[0].ID())
}

func TestClient_DeleteOCR2KeyBundle(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, r := app.NewClientAndRenderer()

	key, err := app.GetKeyStore().OCR2().Create("evm")
	require.NoError(t, err)

	requireOCR2KeyCount(t, app, 1)

	set := flag.NewFlagSet("test", 0)
	set.Parse([]string{key.ID()})
	set.Bool("yes", true, "")
	c := cli.NewContext(nil, set, nil)

	require.NoError(t, client.DeleteOCR2KeyBundle(c))
	requireOCR2KeyCount(t, app, 0) // Only fixture key remains

	require.Equal(t, 1, len(r.Renders))
	output := *r.Renders[0].(*cmd.OCR2KeyBundlePresenter)
	assert.Equal(t, key.ID(), output.ID)
}

func TestClient_ImportExportOCR2Key(t *testing.T) {
	defer deleteKeyExportFile(t)

	app := startNewApplication(t, withConfigSet(func(c *configtest.TestGeneralConfig) {
		c.Overrides.EVMDisabled = null.BoolFrom(true)
	}))
	client, _ := app.NewClientAndRenderer()

	app.KeyStore.OCR2().Add(cltest.DefaultOCR2Key)

	keys := requireOCR2KeyCount(t, app, 1)
	key := keys[0]
	keyName := keyNameForTest(t)

	// Export test invalid id
	set := flag.NewFlagSet("test OCR2 export", 0)
	set.Parse([]string{"0"})
	set.String("newpassword", "../internal/fixtures/new_password.txt", "")
	set.String("output", keyName, "")
	c := cli.NewContext(nil, set, nil)
	err := client.ExportOCR2Key(c)
	require.Error(t, err, "Error exporting")
	require.Error(t, utils.JustError(os.Stat(keyName)))

	// Export
	set = flag.NewFlagSet("test OCR2 export", 0)
	set.Parse([]string{key.ID()})
	set.String("newpassword", "../internal/fixtures/new_password.txt", "")
	set.String("output", keyName, "")
	c = cli.NewContext(nil, set, nil)

	require.NoError(t, client.ExportOCR2Key(c))
	require.NoError(t, utils.JustError(os.Stat(keyName)))

	require.NoError(t, app.GetKeyStore().OCR2().Delete(key.ID()))
	requireOCR2KeyCount(t, app, 0)

	set = flag.NewFlagSet("test OCR2 import", 0)
	set.Parse([]string{keyName})
	set.String("oldpassword", "../internal/fixtures/new_password.txt", "")
	c = cli.NewContext(nil, set, nil)
	require.NoError(t, client.ImportOCR2Key(c))

	requireOCR2KeyCount(t, app, 1)
}

func requireOCR2KeyCount(t *testing.T, app chainlink.Application, length int) []ocr2key.KeyBundle {
	keys, err := app.GetKeyStore().OCR2().GetAll()
	require.NoError(t, err)
	require.Len(t, keys, length)
	return keys
}
