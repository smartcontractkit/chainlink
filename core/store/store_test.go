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
	"github.com/smartcontractkit/chainlink/core/store/migrationsv2"

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
	err := migrationsv2.MigrateUp(db, "")
	require.NoError(t, err)
	err = store.CheckSquashUpgrade(db)
	require.NoError(t, err)
	err = migrationsv2.MigrateDown(db)
	require.NoError(t, err)

	// Newer app version with older migrations should fail.
	err = migrations.MigrateTo(db, "1611388693") // 1 before S-1
	require.NoError(t, err)
	err = store.CheckSquashUpgrade(db)
	t.Log(err)
	require.Error(t, err)

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

	err = store.ImportKey([]byte(`{"address":"3cb8e3FD9d27e39a5e9e6852b0e96160061fd4ea","crypto":{"cipher":"aes-128-ctr","ciphertext":"7515678239ccbeeaaaf0b103f0fba46a979bf6b2a52260015f35b9eb5fed5c17","cipherparams":{"iv":"87e5a5db334305e1e4fb8b3538ceea12"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"d89ac837b5dcdce5690af764762fe349d8162bb0086cea2bc3a4289c47853f96"},"mac":"57a7f4ada10d3d89644f541c91f89b5bde73e15e827ee40565e2d1f88bb0ac96"},"id":"c8cb9bc7-0a51-43bd-8348-8a67fd1ec52c","version":3}`), cltest.Password)
	require.NoError(t, err)

	keys, err = store.AllKeys()
	require.NoError(t, err)
	require.Len(t, keys, 1)

	var addrs []common.Address
	for _, key := range keys {
		addrs = append(addrs, key.Address.Address())
	}
	require.Contains(t, addrs, common.HexToAddress("0x3cb8e3FD9d27e39a5e9e6852b0e96160061fd4ea"))
}

func TestStore_ExportKey(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	err := store.KeyStore.Unlock(cltest.Password)
	require.NoError(t, err)

	keys, err := store.AllKeys()
	require.NoError(t, err)
	require.Len(t, keys, 0)

	keyJSON := []byte(`{"address":"3cb8e3FD9d27e39a5e9e6852b0e96160061fd4ea","crypto":{"cipher":"aes-128-ctr","ciphertext":"7515678239ccbeeaaaf0b103f0fba46a979bf6b2a52260015f35b9eb5fed5c17","cipherparams":{"iv":"87e5a5db334305e1e4fb8b3538ceea12"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"d89ac837b5dcdce5690af764762fe349d8162bb0086cea2bc3a4289c47853f96"},"mac":"57a7f4ada10d3d89644f541c91f89b5bde73e15e827ee40565e2d1f88bb0ac96"},"id":"c8cb9bc7-0a51-43bd-8348-8a67fd1ec52c","version":3}`)

	err = store.ImportKey(keyJSON, cltest.Password)
	require.NoError(t, err)

	keys, err = store.AllKeys()
	require.NoError(t, err)
	require.Len(t, keys, 1)

	bytes, err := store.KeyStore.Export(common.HexToAddress("0x3cb8e3FD9d27e39a5e9e6852b0e96160061fd4ea"), cltest.Password)
	require.NoError(t, err)

	var addr struct {
		Address string `json:"address"`
	}
	err = json.Unmarshal(bytes, &addr)
	require.NoError(t, err)

	require.Equal(t, common.HexToAddress("0x3cb8e3FD9d27e39a5e9e6852b0e96160061fd4ea"), common.HexToAddress("0x"+addr.Address))
}
