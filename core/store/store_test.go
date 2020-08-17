package store_test

import (
	"math/big"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestStore_Start(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	txManager := new(mocks.TxManager)
	txManager.On("Register", mock.Anything).Return(big.NewInt(3), nil)
	store.TxManager = txManager

	assert.NoError(t, store.Start())

	txManager.AssertExpectations(t)
}

func TestStore_Close(t *testing.T) {
	t.Parallel()

	s, cleanup := cltest.NewStore(t)
	defer cleanup()

	assert.NoError(t, s.Close())
}

func TestStore_SyncDiskKeyStoreToDB_HappyPath(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t, cltest.LenientEthMock)
	defer cleanup()
	require.NoError(t, app.Start())
	store := app.GetStore()

	// create key on disk
	pwd := "p@ssword"
	acc, err := store.KeyStore.NewAccount(pwd)
	require.NoError(t, err)

	// assert creation on disk is successful
	files, err := utils.FilesInDir(app.Config.KeysDir())
	require.NoError(t, err)
	require.Len(t, files, 2)

	// sync
	require.NoError(t, store.SyncDiskKeyStoreToDB())

	// assert creation in db is successful
	keys, err := store.Keys()
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
		require.JSONEq(t, key.JSON.String(), content)
	}
}

func TestStore_SyncDiskKeyStoreToDB_MultipleKeys(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t, cltest.LenientEthMock)
	app.AddUnlockedKey() // second account
	defer cleanup()
	require.NoError(t, app.Start())

	store := app.GetStore()

	// assert creation on disk is successful
	files, err := utils.FilesInDir(app.Config.KeysDir())
	require.NoError(t, err)
	require.Len(t, files, 2)

	// sync
	require.NoError(t, store.SyncDiskKeyStoreToDB())

	// assert creation in db is successful
	keys, err := store.Keys()
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
		require.JSONEq(t, content, payloads[payloadAddress])
	}
}

func TestStore_SyncDiskKeyStoreToDB_DBKeyAlreadyExists(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t,
		cltest.EthMockRegisterGetBalance,
	)
	defer cleanup()
	app.EthMock.Context("app.Start()", func(meth *cltest.EthMock) {
		meth.Register("eth_getTransactionCount", "0x1")
		meth.Register("eth_chainId", app.Store.Config.ChainID())
	})
	require.NoError(t, app.StartAndConnect())
	store := app.GetStore()

	// assert sync worked on NewApplication
	keys, err := store.Keys()
	require.NoError(t, err)
	require.Len(t, keys, 1, "key should already exist because of Application#Start")

	// get account
	acc, err := store.KeyStore.GetFirstAccount()
	require.NoError(t, err)

	require.NoError(t, store.SyncDiskKeyStoreToDB()) // sync

	// assert no change in db
	keys, err = store.Keys()
	require.NoError(t, err)
	require.Len(t, keys, 1)
	require.Equal(t, acc.Address.Hex(), keys[0].Address.String())
}
