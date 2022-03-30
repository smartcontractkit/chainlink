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
	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/web/presenters"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
	"gopkg.in/guregu/null.v4"
)

func TestEthKeysPresenter_RenderTable(t *testing.T) {
	t.Parallel()

	var (
		address        = "0x5431F5F973781809D18643b87B44921b11355d81"
		ethBalance     = assets.NewEth(1)
		linkBalance    = assets.NewLinkFromJuels(2)
		isFunding      = true
		createdAt      = time.Now()
		updatedAt      = time.Now().Add(time.Second)
		maxGasPriceWei = utils.NewBigI(12345)
		bundleID       = cltest.DefaultOCRKeyBundleID
		buffer         = bytes.NewBufferString("")
		r              = cmd.RendererTable{Writer: buffer}
	)

	p := cmd.EthKeyPresenter{
		ETHKeyResource: presenters.ETHKeyResource{
			JAID:           presenters.NewJAID(bundleID),
			Address:        address,
			EthBalance:     ethBalance,
			LinkBalance:    linkBalance,
			IsFunding:      isFunding,
			CreatedAt:      createdAt,
			UpdatedAt:      updatedAt,
			MaxGasPriceWei: *maxGasPriceWei,
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
	assert.Contains(t, output, maxGasPriceWei.String())

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
	assert.Contains(t, output, maxGasPriceWei.String())
}

func TestClient_ListETHKeys(t *testing.T) {
	t.Parallel()

	ethClient := newEthMock(t)
	ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Return(big.NewInt(42), nil)
	ethClient.On("GetLINKBalance", mock.Anything, mock.Anything).Return(assets.NewLinkFromJuels(42), nil)
	app := startNewApplication(t,
		withKey(),
		withMocks(ethClient),
		withConfigSet(func(c *configtest.TestGeneralConfig) {
			c.Overrides.EVMEnabled = null.BoolFrom(true)
			c.Overrides.GlobalEvmNonceAutoSync = null.BoolFrom(false)
			c.Overrides.GlobalBalanceMonitorEnabled = null.BoolFrom(false)
		}),
	)
	client, r := app.NewClientAndRenderer()

	assert.Nil(t, client.ListETHKeys(cltest.EmptyCLIContext()))
	require.Equal(t, 1, len(r.Renders))
	balances := *r.Renders[0].(*cmd.EthKeyPresenters)
	assert.Equal(t, app.Key.Address.Hex(), balances[0].Address)

	ethClient.AssertExpectations(t)
}

func TestClient_CreateETHKey(t *testing.T) {
	t.Parallel()

	ethClient := newEthMock(t)
	ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Return(big.NewInt(42), nil)
	ethClient.On("GetLINKBalance", mock.Anything, mock.Anything).Return(assets.NewLinkFromJuels(42), nil)
	app := startNewApplication(t,
		withKey(),
		withMocks(ethClient),
		withConfigSet(func(c *configtest.TestGeneralConfig) {
			c.Overrides.EVMEnabled = null.BoolFrom(true)
			c.Overrides.GlobalEvmNonceAutoSync = null.BoolFrom(false)
			c.Overrides.GlobalBalanceMonitorEnabled = null.BoolFrom(false)
		}),
	)
	db := app.GetSqlxDB()
	client, _ := app.NewClientAndRenderer()

	cltest.AssertCount(t, db, "eth_key_states", 1) // The initial funding key
	keys, err := app.KeyStore.Eth().GetAll()
	require.NoError(t, err)
	require.Equal(t, 1, len(keys))

	// create a key on the default chain
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)
	assert.NoError(t, client.CreateETHKey(c))

	// create the key on a specific chainID
	id := big.NewInt(0)
	// TODO: re-enable this once ChainSet is smart enough to reload chains at runtime
	// https://app.shortcut.com/chainlinklabs/story/17044/chainset-should-update-chains-when-nodes-are-changed
	// _, err = app.GetChains().EVM.Add(id, evmtypes.ChainCfg{})
	// require.NoError(t, err)

	set = flag.NewFlagSet("test", 0)
	set.String("evmChainID", "", "")
	set.Uint64("maxGasPriceGWei", 0, "")
	c = cli.NewContext(nil, set, nil)
	set.Set("maxGasPriceGWei", "12345")
	set.Parse([]string{"-evmChainID", id.String()})
	assert.NoError(t, client.CreateETHKey(c))

	cltest.AssertCount(t, db, "eth_key_states", 3)
	keys, err = app.KeyStore.Eth().GetAll()
	require.NoError(t, err)
	require.Equal(t, 3, len(keys))

	// TODO: re-enable this once ChainSet is smart enough to reload chains at runtime
	// https://app.shortcut.com/chainlinklabs/story/17044/chainset-should-update-chains-when-nodes-are-changed
	// states, err := app.KeyStore.Eth().GetStatesForChain(id)
	// require.Len(t, states, 1)
}

func TestClient_UpdateETHKey(t *testing.T) {
	t.Parallel()

	ethClient := newEthMock(t)
	ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Return(big.NewInt(42), nil)
	ethClient.On("GetLINKBalance", mock.Anything, mock.Anything).Return(assets.NewLinkFromJuels(42), nil)
	app := startNewApplication(t,
		withKey(),
		withMocks(ethClient),
		withConfigSet(func(c *configtest.TestGeneralConfig) {
			c.Overrides.EVMEnabled = null.BoolFrom(true)
			c.Overrides.GlobalEvmNonceAutoSync = null.BoolFrom(false)
			c.Overrides.GlobalBalanceMonitorEnabled = null.BoolFrom(false)
		}),
	)
	ethKeyStore := app.GetKeyStore().Eth()
	client, _ := app.NewClientAndRenderer()

	// Create the key
	key, err := ethKeyStore.Create(&cltest.FixtureChainID)
	require.NoError(t, err)

	// Update the key
	set := flag.NewFlagSet("test", 0)
	set.Uint64("maxGasPriceGWei", 0, "")
	set.Set("maxGasPriceGWei", "12345")
	set.Parse([]string{key.Address.Hex()})
	c := cli.NewContext(nil, set, nil)
	require.NoError(t, client.UpdateETHKey(c))

	// Checking updated config
	chain, err := app.Chains.EVM.Get(&cltest.FixtureChainID)
	require.NoError(t, err)
	price := chain.Config().KeySpecificMaxGasPriceWei(key.Address.Address())
	require.Equal(t, assets.GWei(12345), price)
}

func TestClient_DeleteETHKey(t *testing.T) {
	t.Parallel()

	ethClient := newEthMock(t)
	ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Return(big.NewInt(42), nil)
	ethClient.On("GetLINKBalance", mock.Anything, mock.Anything).Return(assets.NewLinkFromJuels(42), nil)
	app := startNewApplication(t,
		withKey(),
		withMocks(ethClient),
		withConfigSet(func(c *configtest.TestGeneralConfig) {
			c.Overrides.EVMEnabled = null.BoolFrom(true)
			c.Overrides.GlobalEvmNonceAutoSync = null.BoolFrom(false)
			c.Overrides.GlobalBalanceMonitorEnabled = null.BoolFrom(false)
		}),
	)
	ethKeyStore := app.GetKeyStore().Eth()
	client, _ := app.NewClientAndRenderer()

	// Create the key
	key, err := ethKeyStore.Create(&cltest.FixtureChainID)
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

func TestClient_ImportExportETHKey_NoChains(t *testing.T) {
	t.Parallel()

	t.Cleanup(func() { deleteKeyExportFile(t) })

	ethClient := newEthMock(t)
	ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Return(big.NewInt(42), nil)
	ethClient.On("GetLINKBalance", mock.Anything, mock.Anything).Return(assets.NewLinkFromJuels(42), nil)
	app := startNewApplication(t,
		withMocks(ethClient),
		withConfigSet(func(c *configtest.TestGeneralConfig) {
			c.Overrides.EVMEnabled = null.BoolFrom(true)
			c.Overrides.GlobalEvmNonceAutoSync = null.BoolFrom(false)
			c.Overrides.GlobalBalanceMonitorEnabled = null.BoolFrom(false)
		}),
	)
	client, r := app.NewClientAndRenderer()
	ethKeyStore := app.GetKeyStore().Eth()

	set := flag.NewFlagSet("test", 0)
	set.String("file", "internal/fixtures/apicredentials", "")
	set.Bool("bypass-version-check", true, "")
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

	cltest.AssertCount(t, app.GetSqlxDB(), "eth_key_states", 0)

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
	set.Bool("bypass-version-check", true, "")
	set.String("output", keyName, "")
	c = cli.NewContext(nil, set, nil)
	err = client.ExportETHKey(c)
	require.Error(t, err, "Error exporting")
	require.Error(t, utils.JustError(os.Stat(keyName)))
}
func TestClient_ImportExportETHKey_WithChains(t *testing.T) {
	t.Parallel()

	t.Cleanup(func() { deleteKeyExportFile(t) })

	ethClient := newEthMock(t)
	app := startNewApplication(t,
		withMocks(ethClient),
		withConfigSet(func(c *configtest.TestGeneralConfig) {
			c.Overrides.EVMEnabled = null.BoolFrom(true)
			c.Overrides.GlobalEvmNonceAutoSync = null.BoolFrom(false)
			c.Overrides.GlobalBalanceMonitorEnabled = null.BoolFrom(false)
		}),
	)
	client, r := app.NewClientAndRenderer()
	ethKeyStore := app.GetKeyStore().Eth()

	ethClient.On("Dial", mock.Anything).Maybe()
	ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Return(big.NewInt(42), nil)
	ethClient.On("GetLINKBalance", mock.Anything, mock.Anything).Return(assets.NewLinkFromJuels(42), nil)

	set := flag.NewFlagSet("test", 0)
	set.String("file", "internal/fixtures/apicredentials", "")
	set.Bool("bypass-version-check", true, "")
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
	set.Bool("bypass-version-check", true, "")
	set.String("output", keyName, "")
	c = cli.NewContext(nil, set, nil)
	err = client.ExportETHKey(c)
	require.Error(t, err, "Error exporting")
	require.Error(t, utils.JustError(os.Stat(keyName)))
}
