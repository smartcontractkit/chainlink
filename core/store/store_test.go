package store_test

import (
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestStore_Start(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()

	store := app.Store
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	txmMock := mocks.NewMockTxManager(ctrl)
	store.TxManager = txmMock
	txmMock.EXPECT().Register(gomock.Any())
	assert.NoError(t, store.Start())
}

func TestStore_Close(t *testing.T) {
	t.Parallel()

	s, cleanup := cltest.NewStore(t)
	defer cleanup()

	s.RunChannel.Send(models.NewID())
	s.RunChannel.Send(models.NewID())

	_, open := <-s.RunChannel.Receive()
	assert.True(t, open)

	_, open = <-s.RunChannel.Receive()
	assert.True(t, open)

	assert.NoError(t, s.Close())

	rr, open := <-s.RunChannel.Receive()
	assert.Equal(t, store.RunRequest{}, rr)
	assert.False(t, open)
}

func TestStore_SyncDiskKeyStoreToDB_HappyPath(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t)
	defer cleanup()
	store := app.GetStore()

	// create key on disk
	pwd := "p@ssword"
	acc, err := store.KeyStore.NewAccount(pwd)
	require.NoError(t, err)

	// assert creation on disk is successful
	files, err := utils.FilesInDir(app.Config.KeysDir())
	require.NoError(t, err)
	require.Len(t, files, 1)

	// sync
	require.NoError(t, store.SyncDiskKeyStoreToDB())

	// assert creation in db is successful
	keys, err := store.Keys()
	require.NoError(t, err)
	require.Len(t, keys, 1)
	key := keys[0]
	require.Equal(t, acc.Address.Hex(), key.Address.String())

	// assert contents are the same
	content, err := utils.FileContents(filepath.Join(app.Config.KeysDir(), files[0]))
	require.NoError(t, err)
	require.Equal(t, keys[0].JSON.String(), content)
}

func TestStore_SyncDiskKeyStoreToDB_MultipleKeys(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t)
	app.AddUnlockedKey() // second account
	defer cleanup()

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
		require.Equal(t, content, payloads[payloadAddress])
	}
}

func TestStore_SyncDiskKeyStoreToDB_DBKeyAlreadyExists(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()
	require.NoError(t, utils.JustError(app.MockStartAndConnect()))
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

func TestQueuedRunChannel_Send(t *testing.T) {
	t.Parallel()

	rq := store.NewQueuedRunChannel()

	assert.NoError(t, rq.Send(models.NewID()))
	rr1 := <-rq.Receive()
	assert.NotNil(t, rr1)
}

func TestQueuedRunChannel_Send_afterClose(t *testing.T) {
	t.Parallel()

	rq := store.NewQueuedRunChannel()
	rq.Close()

	assert.Error(t, rq.Send(models.NewID()))
}
