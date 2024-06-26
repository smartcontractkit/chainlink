package cmd_test

import (
	"crypto/rand"
	"flag"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"

	commoncfg "github.com/smartcontractkit/chainlink-common/pkg/config"
	solcfg "github.com/smartcontractkit/chainlink-solana/pkg/solana/config"
	stkcfg "github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/config"
	"github.com/smartcontractkit/chainlink/v2/core/cmd"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/logger/audit"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/sessions"
	"github.com/smartcontractkit/chainlink/v2/core/sessions/localauth"
	"github.com/smartcontractkit/chainlink/v2/plugins"
)

func TestTerminalCookieAuthenticator_AuthenticateWithoutSession(t *testing.T) {
	t.Parallel()

	ctx := testutils.Context(t)
	app := cltest.NewApplicationEVMDisabled(t)
	u := cltest.NewUserWithSession(t, app.AuthenticationProvider())

	tests := []struct {
		name, email, pwd string
	}{
		{"bad email", "notreal", cltest.Password},
		{"bad pwd", u.Email, "mostcommonwrongpwdever"},
		{"bad both", "notreal", "mostcommonwrongpwdever"},
		{"correct", u.Email, cltest.Password},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sr := sessions.SessionRequest{Email: test.email, Password: test.pwd}
			store := &cmd.MemoryCookieStore{}
			tca := cmd.NewSessionCookieAuthenticator(cmd.ClientOpts{}, store, logger.TestLogger(t))
			cookie, err := tca.Authenticate(ctx, sr)

			assert.Error(t, err)
			assert.Nil(t, cookie)
			cookie, err = store.Retrieve()
			assert.NoError(t, err)
			assert.Nil(t, cookie)
		})
	}
}

func TestTerminalCookieAuthenticator_AuthenticateWithSession(t *testing.T) {
	t.Parallel()

	ctx := testutils.Context(t)
	app := cltest.NewApplicationEVMDisabled(t)
	require.NoError(t, app.Start(ctx))

	u := cltest.NewUserWithSession(t, app.AuthenticationProvider())

	tests := []struct {
		name, email, pwd string
		wantError        bool
	}{
		{"bad email", "notreal", cltest.Password, true},
		{"bad pwd", u.Email, "mostcommonwrongpwdever", true},
		{"bad both", "notreal", "mostcommonwrongpwdever", true},
		{"success", u.Email, cltest.Password, false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sr := sessions.SessionRequest{Email: test.email, Password: test.pwd}
			store := &cmd.MemoryCookieStore{}
			tca := cmd.NewSessionCookieAuthenticator(app.NewClientOpts(), store, logger.TestLogger(t))
			cookie, err := tca.Authenticate(ctx, sr)

			if test.wantError {
				assert.Error(t, err)
				assert.Nil(t, cookie)

				cookie, err = store.Retrieve()
				assert.NoError(t, err)
				assert.Nil(t, cookie)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, cookie)

				retrievedCookie, err := store.Retrieve()
				assert.NoError(t, err)
				assert.Equal(t, cookie, retrievedCookie)
			}
		})
	}
}

type diskCookieStoreConfig struct{ rootdir string }

func (d diskCookieStoreConfig) RootDir() string {
	return d.rootdir
}

func TestDiskCookieStore_Retrieve(t *testing.T) {
	t.Parallel()

	cfg := diskCookieStoreConfig{}

	t.Run("missing cookie file", func(t *testing.T) {
		store := cmd.DiskCookieStore{Config: cfg}
		cookie, err := store.Retrieve()
		assert.NoError(t, err)
		assert.Nil(t, cookie)
	})

	t.Run("invalid cookie file", func(t *testing.T) {
		cfg.rootdir = "../internal/fixtures/badcookie"
		store := cmd.DiskCookieStore{Config: cfg}
		cookie, err := store.Retrieve()
		assert.Error(t, err)
		assert.Nil(t, cookie)
	})

	t.Run("valid cookie file", func(t *testing.T) {
		cfg.rootdir = "../internal/fixtures"
		store := cmd.DiskCookieStore{Config: cfg}
		cookie, err := store.Retrieve()
		assert.NoError(t, err)
		assert.NotNil(t, cookie)
	})
}

func TestTerminalAPIInitializer_InitializeWithoutAPIUser(t *testing.T) {
	email := "good@email.com"

	tests := []struct {
		name           string
		enteredStrings []string
		isTerminal     bool
		isError        bool
	}{
		{"correct", []string{email, cltest.Password}, true, false},
		{"bad pwd then correct", []string{email, "p4SsW0r", email, cltest.Password}, true, false},
		{"bad email then correct", []string{"", cltest.Password, email, cltest.Password}, true, false},
		{"not a terminal", []string{}, false, true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := testutils.Context(t)
			db := pgtest.NewSqlxDB(t)
			lggr := logger.TestLogger(t)
			orm := localauth.NewORM(db, time.Minute, lggr, audit.NoopLogger)

			mock := &cltest.MockCountingPrompter{T: t, EnteredStrings: test.enteredStrings, NotTerminal: !test.isTerminal}
			tai := cmd.NewPromptingAPIInitializer(mock)

			// Clear out fixture users/users created from the other test cases
			// This asserts that on initial run with an empty users table that the credentials file will instantiate and
			// create/run with a new admin user
			pgtest.MustExec(t, db, "DELETE FROM users;")

			user, err := tai.Initialize(ctx, orm, lggr)
			if test.isError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(test.enteredStrings), mock.Count)

				persistedUser, err := orm.FindUser(ctx, email)
				assert.NoError(t, err)

				assert.Equal(t, user.Email, persistedUser.Email)
				assert.Equal(t, user.HashedPassword, persistedUser.HashedPassword)
			}
		})
	}
}

func TestTerminalAPIInitializer_InitializeWithExistingAPIUser(t *testing.T) {
	ctx := testutils.Context(t)
	db := pgtest.NewSqlxDB(t)
	lggr := logger.TestLogger(t)
	orm := localauth.NewORM(db, time.Minute, lggr, audit.NoopLogger)

	// Clear out fixture users/users created from the other test cases
	// This asserts that on initial run with an empty users table that the credentials file will instantiate and
	// create/run with a new admin user
	_, err := db.Exec("DELETE FROM users;")
	require.NoError(t, err)

	initialUser := cltest.MustRandomUser(t)
	require.NoError(t, orm.CreateUser(ctx, &initialUser))

	mock := &cltest.MockCountingPrompter{T: t}
	tai := cmd.NewPromptingAPIInitializer(mock)

	// If there is an existing user, and we are in the Terminal prompt, no input prompts required
	user, err := tai.Initialize(ctx, orm, lggr)
	assert.NoError(t, err)
	assert.Equal(t, 0, mock.Count)

	assert.Equal(t, initialUser.Email, user.Email)
	assert.Equal(t, initialUser.HashedPassword, user.HashedPassword)
}

func TestFileAPIInitializer_InitializeWithoutAPIUser(t *testing.T) {
	tests := []struct {
		name      string
		file      string
		wantError bool
	}{
		{"correct", "../internal/fixtures/apicredentials", false},
		{"no file", "", true},
		{"incorrect file", "/tmp/doesnotexist", true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := testutils.Context(t)
			db := pgtest.NewSqlxDB(t)
			lggr := logger.TestLogger(t)
			orm := localauth.NewORM(db, time.Minute, lggr, audit.NoopLogger)

			// Clear out fixture users/users created from the other test cases
			// This asserts that on initial run with an empty users table that the credentials file will instantiate and
			// create/run with a new admin user
			pgtest.MustExec(t, db, "DELETE FROM users;")

			tfi := cmd.NewFileAPIInitializer(test.file)
			user, err := tfi.Initialize(ctx, orm, lggr)
			if test.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, cltest.APIEmailAdmin, user.Email)
				persistedUser, err := orm.FindUser(ctx, user.Email)
				assert.NoError(t, err)
				assert.Equal(t, persistedUser.Email, user.Email)
			}
		})
	}
}

func TestFileAPIInitializer_InitializeWithExistingAPIUser(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	orm := localauth.NewORM(db, time.Minute, logger.TestLogger(t), audit.NoopLogger)

	tests := []struct {
		name      string
		file      string
		wantError bool
	}{
		{"correct", "../internal/fixtures/apicredentials", false},
		{"no file", "", true},
		{"incorrect file", "/tmp/doesnotexist", true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := testutils.Context(t)
			lggr := logger.TestLogger(t)
			tfi := cmd.NewFileAPIInitializer(test.file)
			user, err := tfi.Initialize(ctx, orm, lggr)
			if test.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, cltest.APIEmailAdmin, user.Email)
			}
		})
	}
}

func TestPromptingSessionRequestBuilder(t *testing.T) {
	t.Parallel()

	tests := []struct {
		email, pwd string
	}{
		{"correct@input.com", "mypwd"},
	}

	for _, test := range tests {
		t.Run(test.email, func(t *testing.T) {
			enteredStrings := []string{test.email, test.pwd}
			prompter := &cltest.MockCountingPrompter{T: t, EnteredStrings: enteredStrings}
			builder := cmd.NewPromptingSessionRequestBuilder(prompter)

			sr, err := builder.Build("")
			require.NoError(t, err)
			assert.Equal(t, test.email, sr.Email)
			assert.Equal(t, test.pwd, sr.Password)
		})
	}
}

func TestFileSessionRequestBuilder(t *testing.T) {
	t.Parallel()

	builder := cmd.NewFileSessionRequestBuilder(logger.TestLogger(t))
	tests := []struct {
		name, file, wantEmail string
		wantError             bool
	}{
		{"empty", "", "", true},
		{"correct file", "../internal/fixtures/apicredentials", cltest.APIEmailAdmin, false},
		{"incorrect file", "/tmp/dontexist", "", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sr, err := builder.Build(test.file)
			assert.Equal(t, test.wantEmail, sr.Email)
			if test.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNewUserCache(t *testing.T) {
	r, err := rand.Int(rand.Reader, big.NewInt(256*1024*1024))
	require.NoError(t, err)
	// NewUserCache owns it's Dir.
	// invent a unique subdir that we can cleanup
	// because test.TempDir and ioutil.TempDir don't work well here
	subDir := filepath.Base(fmt.Sprintf("%s-%d", t.Name(), r.Int64()))
	lggr := logger.TestLogger(t)
	c, err := cmd.NewUserCache(subDir, func() logger.Logger { return lggr })
	require.NoError(t, err)
	defer func() {
		require.NoError(t, os.Remove(c.RootDir()))
	}()

	assert.DirExists(t, c.RootDir())
}

func TestSetupSolanaRelayer(t *testing.T) {
	lggr := logger.TestLogger(t)
	reg := plugins.NewLoopRegistry(lggr, nil)
	ks := mocks.NewSolana(t)

	// config 3 chains but only enable 2 => should only be 2 relayer
	nEnabledChains := 2
	tConfig := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.Solana = solcfg.TOMLConfigs{
			&solcfg.TOMLConfig{
				ChainID: ptr[string]("solana-id-1"),
				Enabled: ptr(true),
				Chain:   solcfg.Chain{},
				Nodes:   []*solcfg.Node{},
			},
			&solcfg.TOMLConfig{
				ChainID: ptr[string]("solana-id-2"),
				Enabled: ptr(true),
				Chain:   solcfg.Chain{},
				Nodes:   []*solcfg.Node{},
			},
			&solcfg.TOMLConfig{
				ChainID: ptr[string]("disabled-solana-id-1"),
				Enabled: ptr(false),
				Chain:   solcfg.Chain{},
				Nodes:   []*solcfg.Node{},
			},
		}
	})

	t2Config := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.Solana = solcfg.TOMLConfigs{
			&solcfg.TOMLConfig{
				ChainID: ptr[string]("solana-id-1"),
				Enabled: ptr(true),
				Chain:   solcfg.Chain{},
				Nodes:   []*solcfg.Node{},
			},
		}
	})

	rf := chainlink.RelayerFactory{
		Logger:       lggr,
		LoopRegistry: reg,
	}

	// not parallel; shared state
	t.Run("no plugin", func(t *testing.T) {
		relayers, err := rf.NewSolana(ks, tConfig.SolanaConfigs())
		require.NoError(t, err)
		require.NotNil(t, relayers)
		require.Len(t, relayers, nEnabledChains)
		// no using plugin, so registry should be empty
		require.Len(t, reg.List(), 0)
	})

	t.Run("plugin", func(t *testing.T) {
		t.Setenv("CL_SOLANA_CMD", "phony_solana_cmd")

		relayers, err := rf.NewSolana(ks, tConfig.SolanaConfigs())
		require.NoError(t, err)
		require.NotNil(t, relayers)
		require.Len(t, relayers, nEnabledChains)
		// make sure registry has the plugin
		require.Len(t, reg.List(), nEnabledChains)
	})

	// test that duplicate enabled chains is an error when
	duplicateConfig := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.Solana = solcfg.TOMLConfigs{
			&solcfg.TOMLConfig{
				ChainID: ptr[string]("dupe"),
				Enabled: ptr(true),
				Chain:   solcfg.Chain{},
				Nodes:   []*solcfg.Node{},
			},
			&solcfg.TOMLConfig{
				ChainID: ptr[string]("dupe"),
				Enabled: ptr(true),
				Chain:   solcfg.Chain{},
				Nodes:   []*solcfg.Node{},
			},
		}
	})

	// not parallel; shared state
	t.Run("no plugin, duplicate chains", func(t *testing.T) {
		_, err := rf.NewSolana(ks, duplicateConfig.SolanaConfigs())
		require.Error(t, err)
	})

	t.Run("plugin, duplicate chains", func(t *testing.T) {
		t.Setenv("CL_SOLANA_CMD", "phony_solana_cmd")
		_, err := rf.NewSolana(ks, duplicateConfig.SolanaConfigs())
		require.Error(t, err)
	})

	t.Run("plugin env parsing fails", func(t *testing.T) {
		t.Setenv("CL_SOLANA_CMD", "phony_solana_cmd")
		t.Setenv("CL_SOLANA_ENV", "fake_path")

		_, err := rf.NewSolana(ks, t2Config.SolanaConfigs())
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to parse Solana env file")
	})

	t.Run("plugin already registered", func(t *testing.T) {
		t.Setenv("CL_SOLANA_CMD", "phony_solana_cmd")

		_, err := rf.NewSolana(ks, tConfig.SolanaConfigs())
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to create Solana LOOP command")
	})
}

func TestSetupStarkNetRelayer(t *testing.T) {
	lggr := logger.TestLogger(t)
	reg := plugins.NewLoopRegistry(lggr, nil)
	ks := mocks.NewStarkNet(t)
	// config 3 chains but only enable 2 => should only be 2 relayer
	nEnabledChains := 2
	tConfig := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.Starknet = stkcfg.TOMLConfigs{
			&stkcfg.TOMLConfig{
				ChainID:   ptr[string]("starknet-id-1"),
				Enabled:   ptr(true),
				Chain:     stkcfg.Chain{},
				Nodes:     []*stkcfg.Node{},
				FeederURL: commoncfg.MustParseURL("https://feeder.url"),
			},
			&stkcfg.TOMLConfig{
				ChainID:   ptr[string]("starknet-id-2"),
				Enabled:   ptr(true),
				Chain:     stkcfg.Chain{},
				Nodes:     []*stkcfg.Node{},
				FeederURL: commoncfg.MustParseURL("https://feeder.url"),
			},
			&stkcfg.TOMLConfig{
				ChainID:   ptr[string]("disabled-starknet-id-1"),
				Enabled:   ptr(false),
				Chain:     stkcfg.Chain{},
				Nodes:     []*stkcfg.Node{},
				FeederURL: commoncfg.MustParseURL("https://feeder.url"),
			},
		}
	})

	t2Config := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.Starknet = stkcfg.TOMLConfigs{
			&stkcfg.TOMLConfig{
				ChainID:   ptr[string]("starknet-id-3"),
				Enabled:   ptr(true),
				Chain:     stkcfg.Chain{},
				Nodes:     []*stkcfg.Node{},
				FeederURL: commoncfg.MustParseURL("https://feeder.url"),
			},
		}
	})
	rf := chainlink.RelayerFactory{
		Logger:       lggr,
		LoopRegistry: reg,
	}

	// not parallel; shared state
	t.Run("no plugin", func(t *testing.T) {
		relayers, err := rf.NewStarkNet(ks, tConfig.StarknetConfigs())
		require.NoError(t, err)
		require.NotNil(t, relayers)
		require.Len(t, relayers, nEnabledChains)
		// no using plugin, so registry should be empty
		require.Len(t, reg.List(), 0)
	})

	t.Run("plugin", func(t *testing.T) {
		t.Setenv("CL_STARKNET_CMD", "phony_starknet_cmd")

		relayers, err := rf.NewStarkNet(ks, tConfig.StarknetConfigs())
		require.NoError(t, err)
		require.NotNil(t, relayers)
		require.Len(t, relayers, nEnabledChains)
		// make sure registry has the plugin
		require.Len(t, reg.List(), nEnabledChains)
	})

	// test that duplicate enabled chains is an error when
	duplicateConfig := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.Starknet = stkcfg.TOMLConfigs{
			&stkcfg.TOMLConfig{
				ChainID:   ptr[string]("dupe"),
				Enabled:   ptr(true),
				Chain:     stkcfg.Chain{},
				Nodes:     []*stkcfg.Node{},
				FeederURL: commoncfg.MustParseURL("https://feeder.url"),
			},
			&stkcfg.TOMLConfig{
				ChainID:   ptr[string]("dupe"),
				Enabled:   ptr(true),
				Chain:     stkcfg.Chain{},
				Nodes:     []*stkcfg.Node{},
				FeederURL: commoncfg.MustParseURL("https://feeder.url"),
			},
		}
	})

	// not parallel; shared state
	t.Run("no plugin, duplicate chains", func(t *testing.T) {
		_, err := rf.NewStarkNet(ks, duplicateConfig.StarknetConfigs())
		require.Error(t, err)
	})

	t.Run("plugin, duplicate chains", func(t *testing.T) {
		t.Setenv("CL_STARKNET_CMD", "phony_starknet_cmd")
		_, err := rf.NewStarkNet(ks, duplicateConfig.StarknetConfigs())
		require.Error(t, err)
	})

	t.Run("plugin env parsing fails", func(t *testing.T) {
		t.Setenv("CL_STARKNET_CMD", "phony_starknet_cmd")
		t.Setenv("CL_STARKNET_ENV", "fake_path")

		_, err := rf.NewStarkNet(ks, t2Config.StarknetConfigs())
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to parse Starknet env file")
	})

	t.Run("plugin already registered", func(t *testing.T) {
		t.Setenv("CL_STARKNET_CMD", "phony_starknet_cmd")

		_, err := rf.NewStarkNet(ks, tConfig.StarknetConfigs())
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to create StarkNet LOOP command")
	})
}

// flagSetApplyFromAction applies the flags from action to the flagSet.
// `parentCommand` will filter the app commands and only applies the flags if the command/subcommand has a parent with that name, if left empty no filtering is done
func flagSetApplyFromAction(action interface{}, flagSet *flag.FlagSet, parentCommand string) {
	cliApp := cmd.Shell{}
	app := cmd.NewApp(&cliApp)

	foundName := parentCommand == ""
	actionFuncName := getFuncName(action)

	for _, command := range app.Commands {
		flags := recursiveFindFlagsWithName(actionFuncName, command, parentCommand, foundName)

		for _, flag := range flags {
			flag.Apply(flagSet)
		}
	}
}

func recursiveFindFlagsWithName(actionFuncName string, command cli.Command, parent string, foundName bool) []cli.Flag {
	if command.Action != nil {
		if actionFuncName == getFuncName(command.Action) && foundName {
			return command.Flags
		}
	}

	for _, subcommand := range command.Subcommands {
		if !foundName {
			foundName = strings.EqualFold(subcommand.Name, parent)
		}

		found := recursiveFindFlagsWithName(actionFuncName, subcommand, parent, foundName)
		if found != nil {
			return found
		}
	}
	return nil
}

func getFuncName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
