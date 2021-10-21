package cmd_test

import (
	"bytes"
	"flag"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
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
		isFunding   = true
		createdAt   = time.Now()
		updatedAt   = time.Now().Add(time.Second)
		bundleID    = cltest.DefaultOCRKeyBundleID
		buffer      = bytes.NewBufferString("")
		r           = cmd.RendererTable{Writer: buffer}
	)

	p := cmd.EthKeyPresenter{
		ETHKeyResource: presenters.ETHKeyResource{
			JAID:        presenters.NewJAID(bundleID),
			Address:     address,
			EthBalance:  ethBalance,
			LinkBalance: linkBalance,
			IsFunding:   isFunding,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		},
	}

	// Render a single resource
	require.NoError(t, p.RenderTable(r))

	output := buffer.String()
	assert.Contains(t, output, address)
	assert.Contains(t, output, ethBalance.String())
	assert.Contains(t, output, linkBalance.String())
	assert.Contains(t, output, strconv.FormatBool(isFunding))
	assert.Contains(t, output, createdAt.String())
	assert.Contains(t, output, updatedAt.String())

	// Render many resources
	buffer.Reset()
	ps := cmd.EthKeyPresenters{p}
	require.NoError(t, ps.RenderTable(r))

	output = buffer.String()
	assert.Contains(t, output, address)
	assert.Contains(t, output, ethBalance.String())
	assert.Contains(t, output, linkBalance.String())
	assert.Contains(t, output, strconv.FormatBool(isFunding))
	assert.Contains(t, output, createdAt.String())
	assert.Contains(t, output, updatedAt.String())
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
	db := store.DB
	client, _ := app.NewClientAndRenderer()

	ethClient.On("Dial", mock.Anything).Maybe()
	ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Return(big.NewInt(42), nil)
	ethClient.On("GetLINKBalance", mock.Anything, mock.Anything).Return(assets.NewLink(42), nil)

	cltest.AssertCount(t, db, ethkey.State{}, 1) // The initial funding key
	keys, err := app.KeyStore.Eth().GetAll()
	require.NoError(t, err)
	require.Equal(t, 1, len(keys))

	assert.NoError(t, client.CreateETHKey(nilContext))

	cltest.AssertCount(t, db, ethkey.State{}, 2)
	keys, err = app.KeyStore.Eth().GetAll()
	require.NoError(t, err)
	require.Equal(t, 2, len(keys))
}

func TestClient_DeleteEthKey(t *testing.T) {
	t.Parallel()

	ethClient := newEthMock(t)
	app := startNewApplication(t,
		withKey(),
		withMocks(ethClient),
	)
	ethKeyStore := app.GetKeyStore().Eth()
	client, _ := app.NewClientAndRenderer()

	ethClient.On("Dial", mock.Anything)
	ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Maybe().Return(big.NewInt(42), nil)
	ethClient.On("GetLINKBalance", mock.Anything, mock.Anything).Maybe().Return(assets.NewLink(42), nil)

	// Create the key
	key, err := ethKeyStore.Create()
	require.NoError(t, err)

	// Delete the key
	set := flag.NewFlagSet("test", 0)
	set.Bool("hard", true, "")
	set.Bool("yes", true, "")
	set.Parse([]string{key.Address.Hex()})
	c := cli.NewContext(nil, set, nil)
	err = client.DeleteETHKey(c)
	require.NoError(t, err)

	_, err = ethKeyStore.Get(key.Address.Hex())
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
	ethKeyStore := app.GetKeyStore().Eth()

	ethClient.On("Dial", mock.Anything).Maybe()
	ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Return(big.NewInt(42), nil)
	ethClient.On("GetLINKBalance", mock.Anything, mock.Anything).Return(assets.NewLink(42), nil)

	set := flag.NewFlagSet("test", 0)
	set.String("file", "internal/fixtures/apicredentials", "")
	c := cli.NewContext(nil, set, nil)
	err := client.RemoteLogin(c)
	require.NoError(t, err)

	err = client.ListETHKeys(c)
	require.NoError(t, err)
	keys := *r.Renders[0].(*cmd.EthKeyPresenters)
	require.Len(t, keys, 1)
	address := keys[0].Address

	r.Renders = nil

	// Export the key
	testdir := filepath.Join(os.TempDir(), t.Name())
	err = os.MkdirAll(testdir, 0700|os.ModeDir)
	require.NoError(t, err)
	defer os.RemoveAll(testdir)
	keyfilepath := filepath.Join(testdir, "key")
	set = flag.NewFlagSet("test", 0)
	set.String("oldpassword", "../internal/fixtures/correct_password.txt", "")
	set.String("newpassword", "../internal/fixtures/incorrect_password.txt", "")
	set.String("output", keyfilepath, "")
	set.Parse([]string{address})
	c = cli.NewContext(nil, set, nil)
	err = client.ExportETHKey(c)
	require.NoError(t, err)

	// Delete the key
	set = flag.NewFlagSet("test", 0)
	set.Bool("hard", true, "")
	set.Bool("yes", true, "")
	set.Parse([]string{address})
	c = cli.NewContext(nil, set, nil)
	err = client.DeleteETHKey(c)
	require.NoError(t, err)
	_, err = ethKeyStore.Get(address)
	require.Error(t, err)

	// Import the key
	set = flag.NewFlagSet("test", 0)
	set.String("oldpassword", "../internal/fixtures/incorrect_password.txt", "")
	set.Parse([]string{keyfilepath})
	c = cli.NewContext(nil, set, nil)
	err = client.ImportETHKey(c)
	require.NoError(t, err)

	r.Renders = nil

	set = flag.NewFlagSet("test", 0)
	c = cli.NewContext(nil, set, nil)
	err = client.ListETHKeys(c)
	require.NoError(t, err)
	require.Len(t, *r.Renders[0].(*cmd.EthKeyPresenters), 1)
	_, err = ethKeyStore.Get(address)
	require.NoError(t, err)

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
