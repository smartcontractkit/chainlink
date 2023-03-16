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

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/web/presenters"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
	"gopkg.in/guregu/null.v4"
)

func ptr[T any](t T) *T { return &t }

func TestEthKeysPresenter_RenderTable(t *testing.T) {
	t.Parallel()

	var (
		address        = "0x5431F5F973781809D18643b87B44921b11355d81"
		ethBalance     = assets.NewEth(1)
		linkBalance    = assets.NewLinkFromJuels(2)
		isDisabled     = true
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
			Disabled:       isDisabled,
			CreatedAt:      createdAt,
			UpdatedAt:      updatedAt,
			MaxGasPriceWei: maxGasPriceWei,
		},
	}

	// Render a single resource
	require.NoError(t, p.RenderTable(r))

	output := buffer.String()
	assert.Contains(t, output, address)
	assert.Contains(t, output, ethBalance.String())
	assert.Contains(t, output, linkBalance.String())
	assert.Contains(t, output, strconv.FormatBool(isDisabled))
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
	assert.Contains(t, output, strconv.FormatBool(isDisabled))
	assert.Contains(t, output, createdAt.String())
	assert.Contains(t, output, updatedAt.String())
	assert.Contains(t, output, maxGasPriceWei.String())
}

func TestClient_ListETHKeys(t *testing.T) {
	t.Parallel()

	ethClient := newEthMock(t)
	ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Return(big.NewInt(42), nil)
	ethClient.On("GetLINKBalance", mock.Anything, mock.Anything, mock.Anything).Return(assets.NewLinkFromJuels(13), nil)
	app := startNewApplicationV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].Enabled = ptr(true)
		c.EVM[0].NonceAutoSync = ptr(false)
		c.EVM[0].BalanceMonitor.Enabled = ptr(false)
	},
		withKey(),
		withMocks(ethClient),
	)
	client, r := app.NewClientAndRenderer()

	assert.Nil(t, client.ListETHKeys(cltest.EmptyCLIContext()))
	require.Equal(t, 1, len(r.Renders))
	balances := *r.Renders[0].(*cmd.EthKeyPresenters)
	assert.Equal(t, app.Keys[0].Address.Hex(), balances[0].Address)
	assert.Equal(t, "0.000000000000000042", balances[0].EthBalance.String())
	assert.Equal(t, "13", balances[0].LinkBalance.String())
}

func TestClient_ListETHKeys_Error(t *testing.T) {
	t.Parallel()

	ethClient := newEthMock(t)
	ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("fake error"))
	ethClient.On("GetLINKBalance", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("fake error"))
	app := startNewApplicationV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].Enabled = ptr(true)
		c.EVM[0].NonceAutoSync = ptr(false)
		c.EVM[0].BalanceMonitor.Enabled = ptr(false)
	},
		withKey(),
		withMocks(ethClient),
	)
	client, r := app.NewClientAndRenderer()

	assert.Nil(t, client.ListETHKeys(cltest.EmptyCLIContext()))
	require.Equal(t, 1, len(r.Renders))
	balances := *r.Renders[0].(*cmd.EthKeyPresenters)
	assert.Equal(t, app.Keys[0].Address.Hex(), balances[0].Address)
	assert.Nil(t, balances[0].EthBalance)
	assert.Nil(t, balances[0].LinkBalance)
}

func TestClient_ListETHKeys_Disabled(t *testing.T) {
	t.Parallel()

	ethClient := newEthMock(t)
	app := startNewApplicationV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].Enabled = ptr(false)
	},
		withKey(),
		withMocks(ethClient),
	)
	client, r := app.NewClientAndRenderer()
	keys, err := app.KeyStore.Eth().GetAll()
	require.NoError(t, err)
	require.Equal(t, 1, len(keys))
	k := keys[0]

	assert.Nil(t, client.ListETHKeys(cltest.EmptyCLIContext()))
	require.Equal(t, 1, len(r.Renders))
	balances := *r.Renders[0].(*cmd.EthKeyPresenters)
	assert.Equal(t, app.Keys[0].Address.Hex(), balances[0].Address)
	assert.Nil(t, balances[0].EthBalance)
	assert.Nil(t, balances[0].LinkBalance)
	assert.Nil(t, balances[0].MaxGasPriceWei)
	assert.Equal(t, []string{
		k.Address.String(), "0", "0", "<nil>", "0", "false",
		balances[0].UpdatedAt.String(), balances[0].CreatedAt.String(), "<nil>",
	}, balances[0].ToRow())
}

func TestClient_CreateETHKey(t *testing.T) {
	t.Parallel()

	ethClient := newEthMock(t)
	ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Return(big.NewInt(42), nil)
	ethClient.On("GetLINKBalance", mock.Anything, mock.Anything, mock.Anything).Return(assets.NewLinkFromJuels(42), nil)
	app := startNewApplicationV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].Enabled = ptr(true)
		c.EVM[0].NonceAutoSync = ptr(false)
		c.EVM[0].BalanceMonitor.Enabled = ptr(false)
	},
		withKey(),
		withMocks(ethClient),
	)
	db := app.GetSqlxDB()
	client, _ := app.NewClientAndRenderer()

	cltest.AssertCount(t, db, "evm_key_states", 1) // The initial funding key
	keys, err := app.KeyStore.Eth().GetAll()
	require.NoError(t, err)
	require.Equal(t, 1, len(keys))

	// create a key on the default chain
	set := flag.NewFlagSet("test", 0)
	cltest.FlagSetApplyFromAction(client.CreateETHKey, set, "")
	c := cli.NewContext(nil, set, nil)
	assert.NoError(t, client.CreateETHKey(c))

	// create the key on a specific chainID
	id := big.NewInt(0)

	set = flag.NewFlagSet("test", 0)
	cltest.FlagSetApplyFromAction(client.CreateETHKey, set, "")

	require.NoError(t, set.Set("evmChainID", ""))

	c = cli.NewContext(nil, set, nil)
	require.NoError(t, set.Parse([]string{"-evmChainID", id.String()}))
	assert.NoError(t, client.CreateETHKey(c))

	cltest.AssertCount(t, db, "evm_key_states", 3)
	keys, err = app.KeyStore.Eth().GetAll()
	require.NoError(t, err)
	require.Equal(t, 3, len(keys))
}

// https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
func TestClient_UpdateETHKey(t *testing.T) {
	t.Parallel()

	ethClient := newEthMock(t)
	ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Return(big.NewInt(42), nil)
	ethClient.On("GetLINKBalance", mock.Anything, mock.Anything, mock.Anything).Return(assets.NewLinkFromJuels(42), nil)
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
	cltest.FlagSetApplyFromAction(client.UpdateETHKey, set, "")

	require.NoError(t, set.Set("maxGasPriceGWei", "0"))
	require.NoError(t, set.Set("maxGasPriceGWei", "12345"))
	require.NoError(t, set.Parse([]string{key.Address.Hex()}))

	c := cli.NewContext(nil, set, nil)
	require.NoError(t, client.UpdateETHKey(c))

	// Checking updated config
	chain, err := app.Chains.EVM.Get(&cltest.FixtureChainID)
	require.NoError(t, err)
	price := chain.Config().KeySpecificMaxGasPriceWei(key.Address)
	require.Equal(t, assets.GWei(12345), price)
}

func TestClient_DeleteETHKey(t *testing.T) {
	t.Parallel()

	app := startNewApplicationV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].Enabled = ptr(true)
		c.EVM[0].NonceAutoSync = ptr(false)
		c.EVM[0].BalanceMonitor.Enabled = ptr(false)
	},
		withKey(),
	)
	ethKeyStore := app.GetKeyStore().Eth()
	client, _ := app.NewClientAndRenderer()

	// Create the key
	key, err := ethKeyStore.Create(&cltest.FixtureChainID)
	require.NoError(t, err)

	// Delete the key
	set := flag.NewFlagSet("test", 0)
	cltest.FlagSetApplyFromAction(client.DeleteETHKey, set, "")

	require.NoError(t, set.Set("yes", "true"))
	require.NoError(t, set.Parse([]string{key.Address.Hex()}))

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
	ethClient.On("GetLINKBalance", mock.Anything, mock.Anything, mock.Anything).Return(assets.NewLinkFromJuels(42), nil)
	app := startNewApplicationV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].Enabled = ptr(true)
		c.EVM[0].NonceAutoSync = ptr(false)
		c.EVM[0].BalanceMonitor.Enabled = ptr(false)
	},
		withMocks(ethClient),
	)
	client, r := app.NewClientAndRenderer()
	ethKeyStore := app.GetKeyStore().Eth()

	set := flag.NewFlagSet("test", 0)
	cltest.FlagSetApplyFromAction(client.RemoteLogin, set, "")

	require.NoError(t, set.Set("file", "internal/fixtures/apicredentials"))
	require.NoError(t, set.Set("bypass-version-check", "true"))

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
	cltest.FlagSetApplyFromAction(client.ExportETHKey, set, "")

	require.NoError(t, set.Set("newpassword", "../internal/fixtures/incorrect_password.txt"))
	require.NoError(t, set.Set("output", keyfilepath))
	require.NoError(t, set.Parse([]string{address}))

	c = cli.NewContext(nil, set, nil)
	err = client.ExportETHKey(c)
	require.NoError(t, err)

	// Delete the key
	set = flag.NewFlagSet("test", 0)
	cltest.FlagSetApplyFromAction(client.DeleteETHKey, set, "")

	require.NoError(t, set.Set("yes", "true"))
	require.NoError(t, set.Parse([]string{address}))

	c = cli.NewContext(nil, set, nil)
	err = client.DeleteETHKey(c)
	require.NoError(t, err)
	_, err = ethKeyStore.Get(address)
	require.Error(t, err)

	cltest.AssertCount(t, app.GetSqlxDB(), "evm_key_states", 0)

	// Import the key
	set = flag.NewFlagSet("test", 0)
	set.String("oldpassword", "../internal/fixtures/incorrect_password.txt", "")
	set.Parse([]string{keyfilepath})
	c = cli.NewContext(nil, set, nil)
	err = client.ImportETHKey(c)
	require.NoError(t, err)

	r.Renders = nil

	set = flag.NewFlagSet("test", 0)
	cltest.FlagSetApplyFromAction(client.ListETHKeys, set, "")
	c = cli.NewContext(nil, set, nil)
	err = client.ListETHKeys(c)
	require.NoError(t, err)
	require.Len(t, *r.Renders[0].(*cmd.EthKeyPresenters), 1)
	_, err = ethKeyStore.Get(address)
	require.NoError(t, err)

	// Export test invalid id
	keyName := keyNameForTest(t)
	set = flag.NewFlagSet("test Eth export invalid id", 0)
	cltest.FlagSetApplyFromAction(client.ExportETHKey, set, "")

	require.NoError(t, set.Parse([]string{"999"}))
	require.NoError(t, set.Set("newpassword", "../internal/fixtures/apicredentials"))
	require.NoError(t, set.Set("output", "keyName"))

	c = cli.NewContext(nil, set, nil)
	err = client.ExportETHKey(c)
	require.Error(t, err, "Error exporting")
	require.Error(t, utils.JustError(os.Stat(keyName)))
}
func TestClient_ImportExportETHKey_WithChains(t *testing.T) {
	t.Parallel()

	t.Cleanup(func() { deleteKeyExportFile(t) })

	ethClient := newEthMock(t)
	app := startNewApplicationV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].Enabled = ptr(true)
		c.EVM[0].NonceAutoSync = ptr(false)
		c.EVM[0].BalanceMonitor.Enabled = ptr(false)
	},
		withMocks(ethClient),
	)
	client, r := app.NewClientAndRenderer()
	ethKeyStore := app.GetKeyStore().Eth()

	ethClient.On("Dial", mock.Anything).Maybe()
	ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Return(big.NewInt(42), nil)
	ethClient.On("GetLINKBalance", mock.Anything, mock.Anything, mock.Anything).Return(assets.NewLinkFromJuels(42), nil)

	set := flag.NewFlagSet("test", 0)
	cltest.FlagSetApplyFromAction(client.RemoteLogin, set, "")

	require.NoError(t, set.Set("file", "internal/fixtures/apicredentials"))
	require.NoError(t, set.Set("bypass-version-check", "true"))

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
	cltest.FlagSetApplyFromAction(client.ExportETHKey, set, "")

	require.NoError(t, set.Set("newpassword", "../internal/fixtures/incorrect_password.txt"))
	require.NoError(t, set.Set("output", keyfilepath))
	require.NoError(t, set.Parse([]string{address}))

	c = cli.NewContext(nil, set, nil)
	err = client.ExportETHKey(c)
	require.NoError(t, err)

	// Delete the key
	set = flag.NewFlagSet("test", 0)
	cltest.FlagSetApplyFromAction(client.DeleteETHKey, set, "")

	require.NoError(t, set.Set("yes", "true"))
	require.NoError(t, set.Parse([]string{address}))

	c = cli.NewContext(nil, set, nil)
	err = client.DeleteETHKey(c)
	require.NoError(t, err)
	_, err = ethKeyStore.Get(address)
	require.Error(t, err)

	// Import the key
	set = flag.NewFlagSet("test", 0)
	cltest.FlagSetApplyFromAction(client.ImportETHKey, set, "")

	require.NoError(t, set.Set("oldpassword", "../internal/fixtures/incorrect_password.txt"))
	require.NoError(t, set.Parse([]string{keyfilepath}))

	c = cli.NewContext(nil, set, nil)
	err = client.ImportETHKey(c)
	require.NoError(t, err)

	r.Renders = nil

	set = flag.NewFlagSet("test", 0)
	cltest.FlagSetApplyFromAction(client.ListETHKeys, set, "")
	c = cli.NewContext(nil, set, nil)
	err = client.ListETHKeys(c)
	require.NoError(t, err)
	require.Len(t, *r.Renders[0].(*cmd.EthKeyPresenters), 1)
	_, err = ethKeyStore.Get(address)
	require.NoError(t, err)

	// Export test invalid id
	keyName := keyNameForTest(t)
	set = flag.NewFlagSet("test Eth export invalid id", 0)
	cltest.FlagSetApplyFromAction(client.ExportETHKey, set, "")

	require.NoError(t, set.Parse([]string{"999"}))
	require.NoError(t, set.Set("newpassword", "../internal/fixtures/apicredentials"))
	require.NoError(t, set.Set("output", keyName))

	c = cli.NewContext(nil, set, nil)
	err = client.ExportETHKey(c)
	require.Error(t, err, "Error exporting")
	require.Error(t, utils.JustError(os.Stat(keyName)))
}

func TestClient_UpdateChainEVMKey(t *testing.T) {
	t.Parallel()

	ethClient := newEthMock(t)
	ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Return(big.NewInt(42), nil)
	ethClient.On("GetLINKBalance", mock.Anything, mock.Anything, mock.Anything).Return(assets.NewLinkFromJuels(42), nil)
	app := startNewApplicationV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].Enabled = ptr(true)
		c.EVM[0].NonceAutoSync = ptr(false)
		c.EVM[0].BalanceMonitor.Enabled = ptr(false)
	},
		withKey(),
		withMocks(ethClient),
	)
	db := app.GetSqlxDB()
	client, _ := app.NewClientAndRenderer()

	newFlagSet := func() *flag.FlagSet {
		fs := testutils.NewTestFlagSet()
		cltest.FlagSetApplyFromAction(client.UpdateChainEVMKey, fs, "")
		return fs
	}

	t.Run("resets a key nonce", func(t *testing.T) {
		fs := newFlagSet()
		require.NoError(t, fs.Set("evmChainID", "0"))
		require.NoError(t, fs.Set("setNextNonce", "42"))
		require.NoError(t, fs.Set("address", app.Keys[0].Address.Hex()))
		c := cli.NewContext(nil, fs, nil)
		assert.NoError(t, client.UpdateChainEVMKey(c))

		var nonce int64
		require.NoError(t, db.Get(&nonce, `SELECT next_nonce FROM evm_key_states`))
		require.Equal(t, int64(42), nonce)
	})

	t.Run("disables and enables a key", func(t *testing.T) {
		fs := newFlagSet()
		require.NoError(t, fs.Set("evmChainID", "0"))
		require.NoError(t, fs.Set("disable", "true"))
		require.NoError(t, fs.Set("enable", "true"))
		require.NoError(t, fs.Set("address", app.Keys[0].Address.Hex()))
		c := cli.NewContext(nil, fs, nil)

		err := client.UpdateChainEVMKey(c)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot set both --enable and --disable simultaneously")

		fs = newFlagSet()
		require.NoError(t, fs.Set("evmChainID", "0"))
		require.NoError(t, fs.Set("disable", "true"))
		require.NoError(t, fs.Set("address", app.Keys[0].Address.Hex()))
		c = cli.NewContext(nil, fs, nil)

		assert.NoError(t, client.UpdateChainEVMKey(c))

		testutils.AssertCount(t, db, "evm_key_states", 1)
		var disabled bool
		require.NoError(t, db.Get(&disabled, `SELECT disabled FROM evm_key_states WHERE address = $1`, app.Keys[0].Address))
		require.True(t, disabled)

		fs = newFlagSet()
		require.NoError(t, fs.Set("evmChainID", "0"))
		require.NoError(t, fs.Set("enable", "true"))
		require.NoError(t, fs.Set("address", app.Keys[0].Address.Hex()))
		c = cli.NewContext(nil, fs, nil)

		assert.NoError(t, client.UpdateChainEVMKey(c))

		testutils.AssertCount(t, db, "evm_key_states", 1)
		require.NoError(t, db.Get(&disabled, `SELECT disabled FROM evm_key_states WHERE address = $1`, app.Keys[0].Address))
		assert.False(t, disabled)
	})
}
