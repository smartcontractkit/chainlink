package cmd_test

import (
	"bytes"
	"flag"
	"os"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ocrkey"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
)

func TestOCRKeyBundlePresenter_RenderTable(t *testing.T) {
	t.Parallel()

	var (
		createdAt = time.Now()
		updatedAt = time.Now().Add(time.Second)
		deletedAt = time.Now().Add(2 * time.Second)
		bundleID  = "7f993fb701b3410b1f6e8d4d93a7462754d24609b9b31a4fe64a0cb475a4d934"
		buffer    = bytes.NewBufferString("")
		r         = cmd.RendererTable{Writer: buffer}
	)

	pk, err := ocrkey.NewKeyBundle()
	require.NoError(t, err)
	pkEncrypted, err := pk.Encrypt("p4SsW0rD1!@#_", utils.FastScryptParams)
	require.NoError(t, err)

	p := cmd.OCRKeyBundlePresenter{
		JAID: cmd.JAID{ID: bundleID},
		OCRKeysBundleResource: presenters.OCRKeysBundleResource{
			JAID:                  presenters.NewJAID(bundleID),
			OnChainSigningAddress: pkEncrypted.OnChainSigningAddress,
			OffChainPublicKey:     pkEncrypted.OffChainPublicKey,
			ConfigPublicKey:       pkEncrypted.ConfigPublicKey,
			CreatedAt:             createdAt,
			UpdatedAt:             updatedAt,
			DeletedAt:             &deletedAt,
		},
	}

	// Render a single resource
	require.NoError(t, p.RenderTable(r))

	output := buffer.String()
	assert.Contains(t, output, bundleID)
	assert.Contains(t, output, pkEncrypted.OnChainSigningAddress.String())
	assert.Contains(t, output, pkEncrypted.OffChainPublicKey.String())
	assert.Contains(t, output, pkEncrypted.ConfigPublicKey.String())
	assert.Contains(t, output, createdAt.String())
	assert.Contains(t, output, updatedAt.String())
	assert.Contains(t, output, deletedAt.String())

	// Render many resources
	buffer.Reset()
	ps := cmd.OCRKeyBundlePresenters{p}
	require.NoError(t, ps.RenderTable(r))

	output = buffer.String()
	assert.Contains(t, output, bundleID)
	assert.Contains(t, output, pkEncrypted.OnChainSigningAddress.String())
	assert.Contains(t, output, pkEncrypted.OffChainPublicKey.String())
	assert.Contains(t, output, pkEncrypted.ConfigPublicKey.String())
	assert.Contains(t, output, createdAt.String())
	assert.Contains(t, output, updatedAt.String())
	assert.Contains(t, output, deletedAt.String())
}

func TestClient_ListOCRKeyBundles(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, r := app.NewClientAndRenderer()

	app.GetKeyStore().OCR.Unlock(cltest.Password)

	key, err := ocrkey.NewKeyBundle()
	require.NoError(t, err)
	encKey, err := key.Encrypt(cltest.Password, utils.FastScryptParams)
	require.NoError(t, err)
	err = app.GetKeyStore().OCR.CreateEncryptedOCRKeyBundle(encKey)
	require.NoError(t, err)

	requireOCRKeyCount(t, app, 2) // Created key + fixture key

	assert.Nil(t, client.ListOCRKeyBundles(cltest.EmptyCLIContext()))
	require.Equal(t, 1, len(r.Renders))
	output := *r.Renders[0].(*cmd.OCRKeyBundlePresenters)
	assert.Equal(t, encKey.ID.String(), output[1].ID)
}

func TestClient_CreateOCRKeyBundle(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, r := app.NewClientAndRenderer()

	app.GetKeyStore().OCR.Unlock(cltest.Password)

	requireOCRKeyCount(t, app, 1) // The initial fixture key

	require.NoError(t, client.CreateOCRKeyBundle(nilContext))

	keys, err := app.GetKeyStore().OCR.FindEncryptedOCRKeyBundles()
	require.NoError(t, err)
	require.Len(t, keys, 2)

	// Check we can decrypt the created key
	for _, e := range keys {
		_, err = e.Decrypt(cltest.Password)
		require.NoError(t, err)
	}

	require.Equal(t, 1, len(r.Renders))
	output := *r.Renders[0].(*cmd.OCRKeyBundlePresenter)
	assert.Equal(t, keys[1].ID.String(), output.ID)
}

func TestClient_DeleteOCRKeyBundle(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, r := app.NewClientAndRenderer()

	app.GetKeyStore().OCR.Unlock(cltest.Password)

	key, err := ocrkey.NewKeyBundle()
	require.NoError(t, err)
	encKey, err := key.Encrypt(cltest.Password, utils.FastScryptParams)
	require.NoError(t, err)
	err = app.GetKeyStore().OCR.CreateEncryptedOCRKeyBundle(encKey)
	require.NoError(t, err)

	requireOCRKeyCount(t, app, 2) // Created key + fixture key

	set := flag.NewFlagSet("test", 0)
	set.Parse([]string{key.ID.String()})
	set.Bool("yes", true, "")
	c := cli.NewContext(nil, set, nil)

	require.NoError(t, client.DeleteOCRKeyBundle(c))
	requireOCRKeyCount(t, app, 1) // Only fixture key remains

	require.Equal(t, 1, len(r.Renders))
	output := *r.Renders[0].(*cmd.OCRKeyBundlePresenter)
	assert.Equal(t, key.ID.String(), output.ID)
}

func TestClient_ImportExportOCRKeyBundle(t *testing.T) {
	defer deleteKeyExportFile(t)

	app := startNewApplication(t)
	client, r := app.NewClientAndRenderer()

	app.GetKeyStore().OCR.Unlock(cltest.Password)

	keys := requireOCRKeyCount(t, app, 1)
	key := keys[0]
	keyName := keyNameForTest(t)

	// Export test invalid id
	set := flag.NewFlagSet("test OCR export", 0)
	set.Parse([]string{"0"})
	set.String("newpassword", "../internal/fixtures/apicredentials", "")
	set.String("output", keyName, "")
	c := cli.NewContext(nil, set, nil)
	err := client.ExportOCRKey(c)
	require.Error(t, err, "Error exporting")
	require.Error(t, utils.JustError(os.Stat(keyName)))

	// Export
	set = flag.NewFlagSet("test OCR export", 0)
	set.Parse([]string{key.ID.String()})
	set.String("newpassword", "../internal/fixtures/apicredentials", "")
	set.String("output", keyName, "")
	c = cli.NewContext(nil, set, nil)

	require.NoError(t, client.ExportOCRKey(c))
	require.NoError(t, utils.JustError(os.Stat(keyName)))

	require.NoError(t, app.GetKeyStore().OCR.DeleteEncryptedOCRKeyBundle(&key))
	requireOCRKeyCount(t, app, 0)

	set = flag.NewFlagSet("test OCR import", 0)
	set.Parse([]string{keyName})
	set.String("oldpassword", "../internal/fixtures/apicredentials", "")
	c = cli.NewContext(nil, set, nil)
	require.NoError(t, client.ImportOCRKey(c))

	requireOCRKeyCount(t, app, 1)

	require.Equal(t, 1, len(r.Renders))
	output := *r.Renders[0].(*cmd.OCRKeyBundlePresenter)
	assert.Equal(t, key.ID.String(), output.ID)
}

func requireOCRKeyCount(t *testing.T, app chainlink.Application, length int) []ocrkey.EncryptedKeyBundle {
	keys, err := app.GetKeyStore().OCR.FindEncryptedOCRKeyBundles()
	require.NoError(t, err)
	require.Len(t, keys, length)
	return keys
}
