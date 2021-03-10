package store_test

import (
	"encoding/json"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/static"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/migrations"

	"github.com/smartcontractkit/chainlink/core/services/eth"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestStore_SquashMigrationUpgrade(t *testing.T) {
	_, orm, cleanup := cltest.BootstrapThrowawayORM(t, "migrationssquash", false)
	defer cleanup()
	db := orm.DB

	// Latest migrations should work fine.
	static.Version = "0.9.11"
	err := migrations.MigrateUp(db, "")
	require.NoError(t, err)
	err = store.CheckSquashUpgrade(db)
	require.NoError(t, err)
	static.Version = "unset"
}

func TestStore_Start(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	assert.NoError(t, store.Start())
}

func TestStore_Close(t *testing.T) {
	t.Parallel()

	s, cleanup := cltest.NewStore(t)
	defer cleanup()

	assert.NoError(t, s.Close())
}

func TestStore_SyncDiskKeyStoreToDB_HappyPath(t *testing.T) {
	t.Parallel()

	rpcClient, gethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	app, cleanup := cltest.NewApplicationWithKey(t,
		eth.NewClientWith(rpcClient, gethClient),
	)
	defer cleanup()
	require.NoError(t, app.Start())
	store := app.GetStore()
	pwd := cltest.Password
	require.NoError(t, store.KeyStore.Unlock(pwd))

	// create key on disk
	err := store.KeyStore.Unlock(pwd)
	require.NoError(t, err)
	acc, err := store.KeyStore.NewAccount()
	require.NoError(t, err)

	// assert creation on disk is successful
	files, err := utils.FilesInDir(app.Config.KeysDir())
	require.NoError(t, err)
	require.Len(t, files, 2)

	// sync
	require.NoError(t, store.SyncDiskKeyStoreToDB())

	// assert creation in db is successful
	keys, err := store.SendKeys()
	require.NoError(t, err)
	// New key in addition to fixture key gives 2
	require.Len(t, keys, 2)
	// Newer key will always come later
	key := keys[1]
	require.Equal(t, acc.Address.Hex(), key.Address.String())

	// assert contents are the same
	require.Equal(t, len(keys), len(files))

	// Files are preceded by timestamp so sorting will put the most recent last (to match keys)
	sort.Slice(files, func(i, j int) bool {
		return strings.ToLower(files[i]) < strings.ToLower(files[j])
	})
	for _, f := range files {
		assert.Regexp(t, regexp.MustCompile(`^UTC--\d{4}-\d{2}-\d{2}T\d{2}-\d{2}-\d{2}\.\d{9}Z--[0-9a-fA-F]{40}$`), f)
	}

	for i, key := range keys {
		content, err := utils.FileContents(filepath.Join(app.Config.KeysDir(), files[i]))
		require.NoError(t, err)

		filekey, err := keystore.DecryptKey([]byte(content), cltest.Password)
		require.NoError(t, err)
		dbkey, err := keystore.DecryptKey(key.JSON.Bytes(), cltest.Password)
		require.NoError(t, err)

		require.Equal(t, dbkey, filekey)
	}
}

func TestStore_SyncDiskKeyStoreToDB_MultipleKeys(t *testing.T) {
	t.Parallel()

	rpcClient, gethClient, _, assertMocksCalled := cltest.NewEthMocks(t)
	defer assertMocksCalled()
	app, cleanup := cltest.NewApplicationWithKey(t,
		eth.NewClientWith(rpcClient, gethClient),
	)
	defer cleanup()
	cltest.MustAddRandomKeyToKeystore(t, app.Store) // second account

	store := app.GetStore()

	// assert creation on disk is successful
	files, err := utils.FilesInDir(app.Config.KeysDir())
	require.NoError(t, err)
	require.Len(t, files, 2)

	// sync
	require.NoError(t, store.SyncDiskKeyStoreToDB())

	// assert creation in db is successful
	keys, err := store.SendKeys()
	require.NoError(t, err)
	require.Len(t, keys, 2)

	accounts := store.KeyStore.Accounts()
	accountKeys := []string{}
	for _, a := range accounts {
		accountKeys = append(accountKeys, a.Address.Hex())
	}

	payloads := map[string]string{}
	addresses := []string{}
	for _, k := range keys {
		payloads[strings.ToLower(k.Address.String())[2:]] = k.JSON.String()
		addresses = append(addresses, k.Address.String())
	}
	sort.Strings(accountKeys)
	sort.Strings(addresses)
	require.Equal(t, accountKeys, addresses)

	for _, f := range files {
		path := filepath.Join(app.Config.KeysDir(), f)
		content, err := utils.FileContents(path)
		require.NoError(t, err)

		payloadAddress := gjson.Parse(content).Get("address").String()

		filekey, err := keystore.DecryptKey([]byte(content), cltest.Password)
		require.NoError(t, err)
		dbkey, err := keystore.DecryptKey([]byte(payloads[payloadAddress]), cltest.Password)
		require.NoError(t, err)

		require.Equal(t, dbkey, filekey)
	}
}

func TestStore_SyncDiskKeyStoreToDB_DBKeyAlreadyExists(t *testing.T) {
	t.Parallel()

	rpcClient, gethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	app, cleanup := cltest.NewApplicationWithKey(t,
		eth.NewClientWith(rpcClient, gethClient),
	)
	defer cleanup()
	require.NoError(t, app.StartAndConnect())
	store := app.GetStore()

	// assert sync worked on NewApplication
	keys, err := store.SendKeys()
	require.NoError(t, err)
	require.Len(t, keys, 1, "key should already exist because of Application#Start")

	// get account
	acc := store.KeyStore.Accounts()[0]
	require.NoError(t, err)

	require.NoError(t, store.SyncDiskKeyStoreToDB()) // sync

	// assert no change in db
	keys, err = store.SendKeys()
	require.NoError(t, err)
	require.Len(t, keys, 1)
	require.Equal(t, acc.Address.Hex(), keys[0].Address.String())
}

func TestStore_DeleteKey(t *testing.T) {
	rpcClient, gethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	app, cleanup := cltest.NewApplicationWithKey(t,
		eth.NewClientWith(rpcClient, gethClient),
	)
	defer cleanup()
	require.NoError(t, app.StartAndConnect())
	store := app.GetStore()

	keys, err := store.AllKeys()
	require.NoError(t, err)
	require.Len(t, keys, 1)

	err = store.DeleteKey(keys[0].Address.Address())
	require.NoError(t, err)

	keys, err = store.AllKeys()
	require.NoError(t, err)
	require.Len(t, keys, 0)
}

func TestStore_ArchiveKey(t *testing.T) {
	rpcClient, gethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	app, cleanup := cltest.NewApplicationWithKey(t,
		eth.NewClientWith(rpcClient, gethClient),
	)
	defer cleanup()
	require.NoError(t, app.StartAndConnect())
	store := app.GetStore()

	var addrs []struct {
		Address   common.Address
		DeletedAt time.Time
	}
	err := store.DB.Raw(`SELECT address, deleted_at FROM keys`).Scan(&addrs).Error
	require.NoError(t, err)

	keys, err := store.AllKeys()
	require.NoError(t, err)
	require.Len(t, keys, 1)

	err = store.ArchiveKey(keys[0].Address.Address())
	require.NoError(t, err)

	err = store.DB.Raw(`SELECT address, deleted_at FROM keys`).Scan(&addrs).Error
	require.NoError(t, err)
	require.Len(t, addrs, 1)

	keys, err = store.SendKeys()
	require.NoError(t, err)
	require.Len(t, keys, 0)

	keys, err = store.AllKeys()
	require.NoError(t, err)
	require.Len(t, keys, 0)
}

func TestStore_ImportKey(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	err := store.KeyStore.Unlock(cltest.Password)
	require.NoError(t, err)

	keys, err := store.AllKeys()
	require.NoError(t, err)
	require.Len(t, keys, 0)

	err = store.ImportKey([]byte(`{"address":"72f4f206d41339921570e47409cfef89ad528605","crypto":{"cipher":"aes-128-ctr","ciphertext":"d55d1cf27b464a7262e947fc6b09161c9c56b2efb1a2e6aef8b1ed0c22e02143","cipherparams":{"iv":"ff9effce7ce8318f54029c30e5e60c3a"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":2,"p":2,"r":8,"salt":"bdec27593d039aca0fe87047bf425bd603a6eb134b8f04ee993ef090086300f7"},"mac":"5e06e90baef19112fcc301fb708d20577af9220e8b1f72329f9f06a70aade18e"},"id":"ec04d5fc-49ce-4d98-bdce-13d1dfa89eb9","version":3}`), cltest.Password)
	require.NoError(t, err)

	keys, err = store.AllKeys()
	require.NoError(t, err)
	require.Len(t, keys, 1)

	var addrs []common.Address
	for _, key := range keys {
		addrs = append(addrs, key.Address.Address())
	}
	require.Contains(t, addrs, common.HexToAddress("0x72f4f206d41339921570e47409cfef89ad528605"))
}

func TestStore_ExportKey(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	err := store.KeyStore.Unlock(cltest.Password)
	require.NoError(t, err)

	keys, err := store.AllKeys()
	require.NoError(t, err)
	require.Len(t, keys, 0)

	keyJSON := cltest.MustReadFile(t, "../internal/fixtures/keys/"+cltest.DefaultKeyFixtureFileName)

	err = store.ImportKey(keyJSON, cltest.Password)
	require.NoError(t, err)

	keys, err = store.AllKeys()
	require.NoError(t, err)
	require.Len(t, keys, 1)

	bytes, err := store.KeyStore.Export(common.HexToAddress(cltest.DefaultKeyAddress), cltest.Password)
	require.NoError(t, err)

	var addr struct {
		Address string `json:"address"`
	}
	err = json.Unmarshal(bytes, &addr)
	require.NoError(t, err)

	require.Equal(t, common.HexToAddress(cltest.DefaultKeyAddress), common.HexToAddress("0x"+addr.Address))
}
