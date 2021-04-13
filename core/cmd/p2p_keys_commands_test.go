package cmd_test

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models/p2pkey"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
)

func TestClient_ListP2PKeys(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	app.Store.OCRKeyStore.Unlock(cltest.Password)

	key, err := p2pkey.CreateKey()
	require.NoError(t, err)
	encKey, err := key.ToEncryptedP2PKey(cltest.Password, utils.FastScryptParams)
	require.NoError(t, err)
	err = app.Store.OCRKeyStore.UpsertEncryptedP2PKey(&encKey)
	require.NoError(t, err)

	requireP2PKeyCount(t, app.Store, 2) // Created  + fixture key

	client, r := app.NewClientAndRenderer()

	assert.Nil(t, client.ListP2PKeys(cltest.EmptyCLIContext()))
	require.Equal(t, 1, len(r.Renders))
	keys := *r.Renders[0].(*[]p2pkey.EncryptedP2PKey)
	assert.Equal(t, encKey.PubKey, keys[1].PubKey)
}

func TestClient_CreateP2PKey(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, _ := app.NewClientAndRenderer()

	app.Store.OCRKeyStore.Unlock(cltest.Password)

	require.NoError(t, client.CreateP2PKey(nilContext))

	keys, err := app.GetStore().OCRKeyStore.FindEncryptedP2PKeys()
	require.NoError(t, err)

	// Created + fixture key
	require.Len(t, keys, 2)

	for _, e := range keys {
		_, err = e.Decrypt(cltest.Password)
		require.NoError(t, err)
	}
}

func TestClient_DeleteP2PKey(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, _ := app.NewClientAndRenderer()

	app.Store.OCRKeyStore.Unlock(cltest.Password)

	key, err := p2pkey.CreateKey()
	require.NoError(t, err)
	encKey, err := key.ToEncryptedP2PKey(cltest.Password, utils.FastScryptParams)
	require.NoError(t, err)
	err = app.Store.OCRKeyStore.UpsertEncryptedP2PKey(&encKey)
	require.NoError(t, err)

	requireP2PKeyCount(t, app.Store, 2) // Created  + fixture key

	set := flag.NewFlagSet("test", 0)
	set.Bool("yes", true, "")
	strID := strconv.FormatInt(int64(encKey.ID), 10)
	set.Parse([]string{strID})
	c := cli.NewContext(nil, set, nil)
	err = client.DeleteP2PKey(c)
	require.NoError(t, err)

	requireP2PKeyCount(t, app.Store, 1) // fixture key only
}

func TestClient_ImportExportP2PKeyBundle(t *testing.T) {
	t.Parallel()

	defer deleteKeyExportFile(t)

	app := startNewApplication(t)
	client, _ := app.NewClientAndRenderer()
	store := app.GetStore()

	store.OCRKeyStore.Unlock(cltest.Password)

	keys := requireP2PKeyCount(t, store, 1)
	key := keys[0]
	keyName := keyNameForTest(t)

	// Export test invalid id
	set := flag.NewFlagSet("test P2P export", 0)
	set.Parse([]string{"0"})
	set.String("newpassword", "../internal/fixtures/apicredentials", "")
	set.String("output", keyName, "")
	c := cli.NewContext(nil, set, nil)
	err := client.ExportP2PKey(c)
	require.Error(t, err, "Error exporting")
	require.Error(t, utils.JustError(os.Stat(keyName)))

	// Export test
	set = flag.NewFlagSet("test P2P export", 0)
	set.Parse([]string{fmt.Sprint(key.ID)})
	set.String("newpassword", "../internal/fixtures/apicredentials", "")
	set.String("output", keyName, "")
	c = cli.NewContext(nil, set, nil)

	require.NoError(t, client.ExportP2PKey(c))
	require.NoError(t, utils.JustError(os.Stat(keyName)))

	require.NoError(t, store.OCRKeyStore.DeleteEncryptedP2PKey(&key))
	requireP2PKeyCount(t, store, 0)

	set = flag.NewFlagSet("test P2P import", 0)
	set.Parse([]string{keyName})
	set.String("oldpassword", "../internal/fixtures/apicredentials", "")
	c = cli.NewContext(nil, set, nil)
	require.NoError(t, client.ImportP2PKey(c))

	requireP2PKeyCount(t, store, 1)
}

func requireP2PKeyCount(t *testing.T, store *store.Store, length int) []p2pkey.EncryptedP2PKey {
	t.Helper()

	keys, err := store.OCRKeyStore.FindEncryptedP2PKeys()
	require.NoError(t, err)
	require.Len(t, keys, length)
	return keys
}
