package cmd_test

import (
	"flag"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
)

func TestClient_ListETHKeys(t *testing.T) {
	t.Parallel()

	rpcClient, gethClient := newEthMocks(t)
	app := startNewApplication(t,
		withKey(),
		withMocks(eth.NewClientWith(rpcClient, gethClient)),
	)
	client, r := app.NewClientAndRenderer()

	gethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Maybe().Return(big.NewInt(42), nil)
	rpcClient.On("Call", mock.Anything, "eth_call", mock.Anything, "latest").Return(nil)

	assert.Nil(t, client.ListETHKeys(cltest.EmptyCLIContext()))
	require.Equal(t, 1, len(r.Renders))
	balances := *r.Renders[0].(*[]presenters.ETHKeyResource)
	assert.Equal(t, app.Key.Address.Hex(), balances[0].Address)
}

func TestClient_CreateETHKey(t *testing.T) {
	t.Parallel()

	rpcClient, gethClient := newEthMocks(t)
	app := startNewApplication(t,
		withKey(),
		withMocks(eth.NewClientWith(rpcClient, gethClient)),
	)
	store := app.GetStore()
	client, _ := app.NewClientAndRenderer()

	gethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Return(big.NewInt(42), nil)
	rpcClient.On("Call", mock.Anything, "eth_call", mock.Anything, "latest").Return(nil)

	requireEthKeysCount(t, store, 1) // The initial funding key

	assert.NoError(t, client.CreateETHKey(nilContext))

	requireEthKeysCount(t, store, 2)
}

func TestClient_DeleteEthKey(t *testing.T) {
	t.Parallel()

	rpcClient, gethClient := newEthMocks(t)
	app := startNewApplication(t,
		withKey(),
		withMocks(eth.NewClientWith(rpcClient, gethClient)),
	)
	store := app.GetStore()
	client, _ := app.NewClientAndRenderer()

	gethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Maybe().Return(big.NewInt(42), nil)
	rpcClient.On("Call", mock.Anything, "eth_call", mock.Anything, "latest").Return(nil)

	// Create the key
	account, err := store.KeyStore.NewAccount()
	require.NoError(t, err)
	require.NoError(t, store.SyncDiskKeyStoreToDB())

	// Delete the key
	set := flag.NewFlagSet("test", 0)
	set.Bool("yes", true, "")
	set.Parse([]string{account.Address.Hex()})
	c := cli.NewContext(nil, set, nil)
	err = client.DeleteETHKey(c)
	require.NoError(t, err)

	_, err = store.KeyByAddress(account.Address)
	assert.Error(t, err)
}

func TestClient_ImportExportETHKey(t *testing.T) {
	t.Parallel()

	t.Cleanup(func() { deleteKeyExportFile(t) })

	rpcClient, gethClient := newEthMocks(t)
	app := startNewApplication(t,
		withMocks(eth.NewClientWith(rpcClient, gethClient)),
	)
	client, r := app.NewClientAndRenderer()

	gethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Return(big.NewInt(42), nil)
	rpcClient.On("Call", mock.Anything, "eth_call", mock.Anything, "latest").Return(nil)

	set := flag.NewFlagSet("test", 0)
	set.String("file", "internal/fixtures/apicredentials", "")
	c := cli.NewContext(nil, set, nil)
	err := client.RemoteLogin(c)
	assert.NoError(t, err)

	err = app.Store.KeyStore.Unlock(cltest.Password)
	assert.NoError(t, err)

	err = client.ListETHKeys(c)
	assert.NoError(t, err)
	require.Len(t, *r.Renders[0].(*[]presenters.ETHKeyResource), 0)

	r.Renders = nil

	set = flag.NewFlagSet("test", 0)
	set.String("oldpassword", "../internal/fixtures/correct_password.txt", "")
	set.Parse([]string{"../internal/fixtures/keys/testkey-0x69Ca211a68100E18B40683E96b55cD217AC95006.json"})
	c = cli.NewContext(nil, set, nil)
	err = client.ImportETHKey(c)
	assert.NoError(t, err)

	r.Renders = nil

	set = flag.NewFlagSet("test", 0)
	c = cli.NewContext(nil, set, nil)
	err = client.ListETHKeys(c)
	assert.NoError(t, err)
	require.Len(t, *r.Renders[0].(*[]presenters.ETHKeyResource), 1)

	ethkeys := *r.Renders[0].(*[]presenters.ETHKeyResource)
	addr := common.HexToAddress("0x69Ca211a68100E18B40683E96b55cD217AC95006")
	assert.Equal(t, addr.Hex(), ethkeys[0].Address)

	testdir := filepath.Join(os.TempDir(), t.Name())
	err = os.MkdirAll(testdir, 0700|os.ModeDir)
	assert.NoError(t, err)
	defer os.RemoveAll(testdir)

	keyfilepath := filepath.Join(testdir, "key")
	set = flag.NewFlagSet("test", 0)
	set.String("oldpassword", "../internal/fixtures/correct_password.txt", "")
	set.String("newpassword", "../internal/fixtures/incorrect_password.txt", "")
	set.String("output", keyfilepath, "")
	set.Parse([]string{addr.Hex()})
	c = cli.NewContext(nil, set, nil)
	err = client.ExportETHKey(c)
	assert.NoError(t, err)

	// Now, make sure that the keyfile can be imported with the `newpassword` and yields the correct address
	keyJSON, err := ioutil.ReadFile(keyfilepath)
	assert.NoError(t, err)
	oldpassword, err := ioutil.ReadFile("../internal/fixtures/correct_password.txt")
	assert.NoError(t, err)
	newpassword, err := ioutil.ReadFile("../internal/fixtures/incorrect_password.txt")
	assert.NoError(t, err)

	keystoreDir := filepath.Join(os.TempDir(), t.Name(), "keystore")
	err = os.MkdirAll(keystoreDir, 0700|os.ModeDir)
	assert.NoError(t, err)

	scryptParams := utils.GetScryptParams(app.Store.Config)
	keystore := store.NewKeyStore(keystoreDir, scryptParams)
	err = keystore.Unlock(string(oldpassword))
	assert.NoError(t, err)
	acct, err := keystore.Import(keyJSON, strings.TrimSpace(string(newpassword)))
	assert.NoError(t, err)
	assert.Equal(t, addr.Hex(), acct.Address.Hex())

	// Export test invalid id
	keyName := keyNameForTest(t)
	set = flag.NewFlagSet("test Eth export invalid id", 0)
	set.Parse([]string{"999"})
	set.String("newpassword", "../internal/fixtures/apicredentials", "")
	set.String("output", keyName, "")
	c = cli.NewContext(nil, set, nil)
	err = client.ExportETHKey(c)
	require.Error(t, err, "Error exporting")
	require.Error(t, utils.JustError(os.Stat(keyName)))
}

func requireEthKeysCount(t *testing.T, store *store.Store, length int) []models.Key {
	keys, err := store.AllKeys()
	require.NoError(t, err)
	require.Len(t, keys, length)
	return keys
}
