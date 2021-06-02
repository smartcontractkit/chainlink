package cmd_test

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
)

func TestEthKeysPresenter_RenderTable(t *testing.T) {
	t.Parallel()

	var (
		address     = "0x5431F5F973781809D18643b87B44921b11355d81"
		ethBalance  = assets.NewEth(1)
		linkBalance = assets.NewLink(2)
		nextNonce   = int64(0)
		isFunding   = true
		createdAt   = time.Now()
		updatedAt   = time.Now().Add(time.Second)
		deletedAt   = time.Now().Add(2 * time.Second)
		bundleID    = "7f993fb701b3410b1f6e8d4d93a7462754d24609b9b31a4fe64a0cb475a4d934"
		buffer      = bytes.NewBufferString("")
		r           = cmd.RendererTable{Writer: buffer}
	)

	p := cmd.EthKeyPresenter{
		ETHKeyResource: presenters.ETHKeyResource{
			JAID:        presenters.NewJAID(bundleID),
			Address:     address,
			EthBalance:  ethBalance,
			LinkBalance: linkBalance,
			NextNonce:   nextNonce,
			IsFunding:   isFunding,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
			DeletedAt:   &deletedAt,
		},
	}

	// Render a single resource
	require.NoError(t, p.RenderTable(r))

	output := buffer.String()
	assert.Contains(t, output, address)
	assert.Contains(t, output, ethBalance.String())
	assert.Contains(t, output, linkBalance.String())
	assert.Contains(t, output, fmt.Sprintf("%d", nextNonce))
	assert.Contains(t, output, strconv.FormatBool(isFunding))
	assert.Contains(t, output, createdAt.String())
	assert.Contains(t, output, updatedAt.String())
	assert.Contains(t, output, deletedAt.String())

	// Render many resources
	buffer.Reset()
	ps := cmd.EthKeyPresenters{p}
	require.NoError(t, ps.RenderTable(r))

	output = buffer.String()
	assert.Contains(t, output, address)
	assert.Contains(t, output, ethBalance.String())
	assert.Contains(t, output, linkBalance.String())
	assert.Contains(t, output, fmt.Sprintf("%d", nextNonce))
	assert.Contains(t, output, strconv.FormatBool(isFunding))
	assert.Contains(t, output, createdAt.String())
	assert.Contains(t, output, updatedAt.String())
	assert.Contains(t, output, deletedAt.String())
}

func TestClient_ListETHKeys(t *testing.T) {
	t.Parallel()

	ethClient := newEthMock(t)
	app := startNewApplication(t,
		withKey(),
		withMocks(ethClient),
	)
	client, r := app.NewClientAndRenderer()

	ethClient.On("Dial", mock.Anything)
	ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Maybe().Return(big.NewInt(42), nil)
	ethClient.On("GetLINKBalance", mock.Anything, mock.Anything).Return(assets.NewLink(42), nil)

	assert.Nil(t, client.ListETHKeys(cltest.EmptyCLIContext()))
	require.Equal(t, 1, len(r.Renders))
	balances := *r.Renders[0].(*cmd.EthKeyPresenters)
	assert.Equal(t, app.Key.Address.Hex(), balances[0].Address)
}

func TestClient_CreateETHKey(t *testing.T) {
	t.Parallel()

	ethClient := newEthMock(t)
	app := startNewApplication(t,
		withKey(),
		withMocks(ethClient),
	)
	store := app.GetStore()
	client, _ := app.NewClientAndRenderer()

	ethClient.On("Dial", mock.Anything).Maybe()
	ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Return(big.NewInt(42), nil)
	ethClient.On("GetLINKBalance", mock.Anything, mock.Anything).Return(assets.NewLink(42), nil)

	requireEthKeysCount(t, store, 1) // The initial funding key

	assert.NoError(t, client.CreateETHKey(nilContext))

	requireEthKeysCount(t, store, 2)
}

func TestClient_DeleteEthKey(t *testing.T) {
	t.Parallel()

	ethClient := newEthMock(t)
	app := startNewApplication(t,
		withKey(),
		withMocks(ethClient),
	)
	ethKeyStore := app.GetKeyStore().Eth
	client, _ := app.NewClientAndRenderer()

	ethClient.On("Dial", mock.Anything)
	ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Maybe().Return(big.NewInt(42), nil)
	ethClient.On("GetLINKBalance", mock.Anything, mock.Anything).Return(assets.NewLink(42), nil)

	// Create the key
	key, err := ethKeyStore.CreateNewKey()
	require.NoError(t, err)

	// Delete the key
	set := flag.NewFlagSet("test", 0)
	set.Bool("yes", true, "")
	set.Parse([]string{key.Address.Hex()})
	c := cli.NewContext(nil, set, nil)
	err = client.DeleteETHKey(c)
	require.NoError(t, err)

	_, err = ethKeyStore.KeyByAddress(key.Address.Address())
	assert.Error(t, err)
}

func TestClient_ImportExportETHKey(t *testing.T) {
	t.Parallel()

	t.Cleanup(func() { deleteKeyExportFile(t) })

	ethClient := newEthMock(t)
	app := startNewApplication(t,
		withMocks(ethClient),
	)
	client, r := app.NewClientAndRenderer()

	ethClient.On("Dial", mock.Anything).Maybe()
	ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Return(big.NewInt(42), nil)
	ethClient.On("GetLINKBalance", mock.Anything, mock.Anything).Return(assets.NewLink(42), nil)

	set := flag.NewFlagSet("test", 0)
	set.String("file", "internal/fixtures/apicredentials", "")
	c := cli.NewContext(nil, set, nil)
	err := client.RemoteLogin(c)
	assert.NoError(t, err)

	err = app.GetKeyStore().Eth.Unlock(cltest.Password)
	assert.NoError(t, err)

	err = client.ListETHKeys(c)
	assert.NoError(t, err)
	require.Len(t, *r.Renders[0].(*cmd.EthKeyPresenters), 1)

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
	require.Len(t, *r.Renders[0].(*cmd.EthKeyPresenters), 2)

	ethkeys := *r.Renders[0].(*cmd.EthKeyPresenters)
	addr := common.HexToAddress("0x69Ca211a68100E18B40683E96b55cD217AC95006")
	assert.Equal(t, addr.Hex(), ethkeys[1].Address)

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
	keystore := keystore.New(app.Store.DB, scryptParams).Eth
	err = keystore.Unlock(string(oldpassword))
	assert.NoError(t, err)
	key, err := keystore.ImportKey(keyJSON, strings.TrimSpace(string(newpassword)))
	assert.NoError(t, err)
	assert.Equal(t, addr.Hex(), key.Address.Hex())

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

func requireEthKeysCount(t *testing.T, store *store.Store, length int) []ethkey.Key {
	var keys []ethkey.Key
	err := store.DB.Find(&keys).Error
	require.NoError(t, err)
	require.Len(t, keys, length)
	return keys
}
