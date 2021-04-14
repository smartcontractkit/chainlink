package cmd_test

import (
	"flag"
	"os"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models/ocrkey"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
)

func TestClient_ListOCRKeyBundles(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, r := app.NewClientAndRenderer()

	app.Store.OCRKeyStore.Unlock(cltest.Password)

	key, err := ocrkey.NewKeyBundle()
	require.NoError(t, err)
	encKey, err := key.Encrypt(cltest.Password, utils.FastScryptParams)
	require.NoError(t, err)
	err = app.Store.OCRKeyStore.CreateEncryptedOCRKeyBundle(encKey)
	require.NoError(t, err)

	requireOCRKeyCount(t, app.Store, 2) // Created key + fixture key

	assert.Nil(t, client.ListOCRKeyBundles(cltest.EmptyCLIContext()))
	require.Equal(t, 1, len(r.Renders))
	keys := *r.Renders[0].(*[]ocrkey.EncryptedKeyBundle)
	assert.Equal(t, encKey.ID, keys[1].ID)
}

func TestClient_CreateOCRKeyBundle(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, _ := app.NewClientAndRenderer()
	store := app.GetStore()

	app.Store.OCRKeyStore.Unlock(cltest.Password)

	requireOCRKeyCount(t, store, 1) // The initial fixture key

	require.NoError(t, client.CreateOCRKeyBundle(nilContext))

	keys, err := app.GetStore().OCRKeyStore.FindEncryptedOCRKeyBundles()
	require.NoError(t, err)
	require.Len(t, keys, 2)

	for _, e := range keys {
		_, err = e.Decrypt(cltest.Password)
		require.NoError(t, err)
	}
}

func TestClient_DeleteOCRKeyBundle(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, _ := app.NewClientAndRenderer()

	app.Store.OCRKeyStore.Unlock(cltest.Password)

	key, err := ocrkey.NewKeyBundle()
	require.NoError(t, err)
	encKey, err := key.Encrypt(cltest.Password, utils.FastScryptParams)
	require.NoError(t, err)
	err = app.Store.OCRKeyStore.CreateEncryptedOCRKeyBundle(encKey)
	require.NoError(t, err)

	requireOCRKeyCount(t, app.Store, 2) // Created key + fixture key

	set := flag.NewFlagSet("test", 0)
	set.Parse([]string{key.ID.String()})
	set.Bool("yes", true, "")
	c := cli.NewContext(nil, set, nil)

	require.NoError(t, client.DeleteOCRKeyBundle(c))
	requireOCRKeyCount(t, app.Store, 1) // Only fixture key remains
}

func TestClient_ImportExportOCRKeyBundle(t *testing.T) {
	defer deleteKeyExportFile(t)

	app := startNewApplication(t)
	client, _ := app.NewClientAndRenderer()

	store := app.GetStore()
	store.OCRKeyStore.Unlock(cltest.Password)

	keys := requireOCRKeyCount(t, store, 1)
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

	require.NoError(t, store.OCRKeyStore.DeleteEncryptedOCRKeyBundle(&key))
	requireOCRKeyCount(t, store, 0)

	set = flag.NewFlagSet("test OCR import", 0)
	set.Parse([]string{keyName})
	set.String("oldpassword", "../internal/fixtures/apicredentials", "")
	c = cli.NewContext(nil, set, nil)
	require.NoError(t, client.ImportOCRKey(c))

	requireOCRKeyCount(t, store, 1)
}

func requireOCRKeyCount(t *testing.T, store *store.Store, length int) []ocrkey.EncryptedKeyBundle {
	keys, err := store.OCRKeyStore.FindEncryptedOCRKeyBundles()
	require.NoError(t, err)
	require.Len(t, keys, length)
	return keys
}
