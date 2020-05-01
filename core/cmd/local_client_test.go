package cmd_test

import (
	"flag"
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/logger"
	strpkg "github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/orm"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
)

func TestClient_RunNodeShowsEnv(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	_, err := store.KeyStore.NewAccount(cltest.Password)
	require.NoError(t, err)
	require.NoError(t, store.KeyStore.Unlock(cltest.Password))

	store.Config.Set("LINK_CONTRACT_ADDRESS", "0x514910771AF9Ca656af840dff83E8264EcF986CA")
	store.Config.Set("CHAINLINK_PORT", 6688)

	app := new(mocks.Application)
	app.On("GetStore").Return(store)
	app.On("Start").Return(nil)
	app.On("Stop").Return(nil)

	auth := cltest.CallbackAuthenticator{Callback: func(*strpkg.Store, string) (string, error) { return "", nil }}
	runner := cltest.BlockedRunner{Done: make(chan struct{})}
	client := cmd.Client{
		Config:                 store.Config,
		AppFactory:             cltest.InstanceAppFactory{App: app},
		KeyStoreAuthenticator:  auth,
		FallbackAPIInitializer: &cltest.MockAPIInitializer{},
		Runner:                 runner,
	}

	set := flag.NewFlagSet("test", 0)
	set.Bool("debug", true, "")
	c := cli.NewContext(nil, set, nil)

	// Start RunNode in a goroutine, it will block until we resume the runner
	go func() {
		assert.NoError(t, client.RunNode(c))
	}()

	// Unlock the runner to the client can begin shutdown
	select {
	case runner.Done <- struct{}{}:
	case <-time.After(30 * time.Second):
		t.Fatal("Timed out waiting for runner")
	}

	logger.Sync()
	logs, err := cltest.ReadLogs(store.Config)
	require.NoError(t, err)

	assert.Contains(t, logs, "LOG_LEVEL: debug\\n")
	assert.Contains(t, logs, "LOG_TO_DISK: true")
	assert.Contains(t, logs, "JSON_CONSOLE: false")
	assert.Contains(t, logs, "ROOT: /tmp/chainlink_test/")
	assert.Contains(t, logs, "CHAINLINK_PORT: 6688\\n")
	assert.Contains(t, logs, "ETH_URL: ws://")
	assert.Contains(t, logs, "ETH_CHAIN_ID: 3\\n")
	assert.Contains(t, logs, "CLIENT_NODE_URL: http://")
	assert.Contains(t, logs, "MIN_OUTGOING_CONFIRMATIONS: 6\\n")
	assert.Contains(t, logs, "MIN_INCOMING_CONFIRMATIONS: 1\\n")
	assert.Contains(t, logs, "ETH_GAS_BUMP_THRESHOLD: 3\\n")
	assert.Contains(t, logs, "ETH_GAS_BUMP_WEI: 5000000000\\n")
	assert.Contains(t, logs, "ETH_GAS_PRICE_DEFAULT: 20000000000\\n")
	assert.Contains(t, logs, "LINK_CONTRACT_ADDRESS: 0x514910771AF9Ca656af840dff83E8264EcF986CA\\n")
	assert.Contains(t, logs, "MINIMUM_CONTRACT_PAYMENT: 0.000000000000000100\\n")
	assert.Contains(t, logs, "ORACLE_CONTRACT_ADDRESS: \\n")
	assert.Contains(t, logs, "ALLOW_ORIGINS: http://localhost:3000,http://localhost:6688\\n")
	assert.Contains(t, logs, "BRIDGE_RESPONSE_URL: http://localhost:6688\\n")

	app.AssertExpectations(t)
}

func TestClient_RunNodeWithPasswords(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		pwdfile      string
		wantUnlocked bool
	}{
		{"correct", "../internal/fixtures/correct_password.txt", true},
		{"incorrect", "../internal/fixtures/incorrect_password.txt", false},
		{"wrongfile", "doesntexist.txt", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store, cleanup := cltest.NewStore(t)
			// Clear out fixture
			store.DeleteUser()
			defer cleanup()
			_, err := store.KeyStore.NewAccount(cltest.Password)
			require.NoError(t, err)
			require.NoError(t, store.KeyStore.Unlock(cltest.Password))

			app := new(mocks.Application)
			app.On("GetStore").Return(store)
			app.On("Start").Maybe().Return(nil)
			app.On("Stop").Maybe().Return(nil)

			_, err = store.KeyStore.NewAccount("password") // matches correct_password.txt
			require.NoError(t, err)

			var unlocked bool
			callback := func(store *strpkg.Store, phrase string) (string, error) {
				err := store.KeyStore.Unlock(phrase)
				unlocked = err == nil
				return phrase, err
			}

			auth := cltest.CallbackAuthenticator{Callback: callback}
			apiPrompt := &cltest.MockAPIInitializer{}
			client := cmd.Client{
				Config:                 store.Config,
				AppFactory:             cltest.InstanceAppFactory{App: app},
				KeyStoreAuthenticator:  auth,
				FallbackAPIInitializer: apiPrompt,
				Runner:                 cltest.EmptyRunner{},
			}

			set := flag.NewFlagSet("test", 0)
			set.String("password", test.pwdfile, "")
			c := cli.NewContext(nil, set, nil)

			if test.wantUnlocked {
				assert.NoError(t, client.RunNode(c))
				assert.True(t, unlocked)
				assert.Equal(t, 1, apiPrompt.Count)
			} else {
				assert.Error(t, client.RunNode(c))
				assert.False(t, unlocked)
				assert.Equal(t, 0, apiPrompt.Count)
			}
		})
	}
}

func TestClient_RunNodeWithAPICredentialsFile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		apiFile    string
		wantPrompt bool
		wantError  bool
	}{
		{"correct", "../internal/fixtures/apicredentials", false, false},
		{"no file", "", true, false},
		{"wrong file", "doesntexist.txt", false, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config := orm.NewConfig()

			store, cleanup := cltest.NewStore(t)
			// Clear out fixture
			store.DeleteUser()
			defer cleanup()
			_, err := store.KeyStore.NewAccount(cltest.Password)
			require.NoError(t, err)
			require.NoError(t, store.KeyStore.Unlock(cltest.Password))

			app := new(mocks.Application)
			app.On("GetStore").Return(store)
			app.On("Start").Maybe().Return(nil)
			app.On("Stop").Maybe().Return(nil)

			callback := func(*strpkg.Store, string) (string, error) { return "", nil }
			noauth := cltest.CallbackAuthenticator{Callback: callback}
			apiPrompt := &cltest.MockAPIInitializer{}
			client := cmd.Client{
				Config:                 config,
				AppFactory:             cltest.InstanceAppFactory{App: app},
				KeyStoreAuthenticator:  noauth,
				FallbackAPIInitializer: apiPrompt,
				Runner:                 cltest.EmptyRunner{},
			}

			set := flag.NewFlagSet("test", 0)
			set.String("api", test.apiFile, "")
			c := cli.NewContext(nil, set, nil)

			if test.wantError {
				assert.EqualError(t, client.RunNode(c), "error creating api initializer: open doesntexist.txt: no such file or directory")
			} else {
				assert.NoError(t, client.RunNode(c))
			}
			assert.Equal(t, test.wantPrompt, apiPrompt.Count > 0)

			app.AssertExpectations(t)
		})
	}
}

func TestClient_ImportKey(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t, cltest.EthMockRegisterChainID)
	defer cleanup()
	require.NoError(t, app.Start())

	client, _ := app.NewClientAndRenderer()

	set := flag.NewFlagSet("import", 0)
	set.Parse([]string{"../internal/fixtures/keys/7fc66c61f88A61DFB670627cA715Fe808057123e.json"})
	c := cli.NewContext(nil, set, nil)
	assert.NoError(t, client.ImportKey(c))

	keys, err := app.GetStore().Keys()
	require.NoError(t, err)
	addresses := []string{}
	for _, k := range keys {
		addresses = append(addresses, k.Address.String())
	}

	sort.Strings(addresses)
	expectation := []string{"0x3cb8e3FD9d27e39a5e9e6852b0e96160061fd4ea", "0x7fc66c61f88A61DFB670627cA715Fe808057123e"}
	require.Equal(t, expectation, addresses)
}

func TestClient_LogToDiskOptionDisablesAsExpected(t *testing.T) {
	tests := []struct {
		name            string
		logToDiskValue  bool
		fileShouldExist bool
	}{
		{"LogToDisk = false => no log on disk", false, false},
		{"LogToDisk = true => log on disk (positive control)", true, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, configCleanup := cltest.NewConfig(t)
			defer configCleanup()
			config.Set("CHAINLINK_DEV", true)
			config.Set("LOG_TO_DISK", tt.logToDiskValue)
			require.NoError(t, os.MkdirAll(config.KeysDir(), os.FileMode(0700)))
			defer os.RemoveAll(config.RootDir())

			previousLogger := logger.GetLogger().Desugar()
			logger.SetLogger(config.CreateProductionLogger())
			defer logger.SetLogger(previousLogger)
			filepath := filepath.Join(config.RootDir(), "log.jsonl")
			_, err := os.Stat(filepath)
			assert.Equal(t, os.IsNotExist(err), !tt.fileShouldExist)
		})
	}
}
