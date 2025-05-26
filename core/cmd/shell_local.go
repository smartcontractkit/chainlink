package cmd

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"golang.org/x/sync/errgroup"
	"gopkg.in/guregu/null.v4"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-evm/pkg/assets"
	"github.com/smartcontractkit/chainlink-evm/pkg/gas"
	"github.com/smartcontractkit/chainlink-evm/pkg/keys"
	"github.com/smartcontractkit/chainlink-evm/pkg/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/build"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/sessions"
	"github.com/smartcontractkit/chainlink/v2/core/shutdown"
	"github.com/smartcontractkit/chainlink/v2/core/static"
	"github.com/smartcontractkit/chainlink/v2/core/store"
	"github.com/smartcontractkit/chainlink/v2/core/store/migrate"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
	"github.com/smartcontractkit/chainlink/v2/core/web"
	webPresenters "github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

var ErrProfileTooLong = errors.New("requested profile duration too large")

func initLocalSubCmds(s *Shell, safe bool) []cli.Command {
	return []cli.Command{
		{
			Name:    "start",
			Aliases: []string{"node", "n"},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "api, a",
					Usage: "text file holding the API email and password, each on a line",
				},
				cli.BoolFlag{
					Name:  "debug, d",
					Usage: "set logger level to debug",
				},
				cli.StringFlag{
					Name:  "password, p",
					Usage: "text file holding the password for the node's account",
				},
				cli.StringFlag{
					Name:  "vrfpassword, vp",
					Usage: "text file holding the password for the vrf keys; enables Chainlink VRF oracle",
				},
			},
			Usage:  "Run the Chainlink node",
			Action: s.RunNode,
		},
		{
			Name:   "rebroadcast-transactions",
			Usage:  "Manually rebroadcast txs matching nonce range with the specified gas price. This is useful in emergencies e.g. high gas prices and/or network congestion to forcibly clear out the pending TX queue",
			Action: s.RebroadcastTransactions,
			Flags: []cli.Flag{
				cli.Uint64Flag{
					Name:  "beginningNonce, beginning-nonce, b",
					Usage: "beginning of nonce range to rebroadcast",
				},
				cli.Uint64Flag{
					Name:  "endingNonce, ending-nonce, e",
					Usage: "end of nonce range to rebroadcast (inclusive)",
				},
				cli.Uint64Flag{
					Name:  "gasPriceWei, gas-price-wei, g",
					Usage: "gas price (in Wei) to rebroadcast transactions at",
				},
				cli.StringFlag{
					Name:  "password, p",
					Usage: "text file holding the password for the node's account",
				},
				cli.StringFlag{
					Name:     "address, a",
					Usage:    "The address (in hex format) for the key which we want to rebroadcast transactions",
					Required: true,
				},
				cli.StringFlag{
					Name:  "evmChainID, evm-chain-id",
					Usage: "Chain ID for which to rebroadcast transactions. If left blank, EVM.ChainID will be used.",
				},
				cli.Uint64Flag{
					Name:  "gasLimit, gas-limit",
					Usage: "OPTIONAL: gas limit to use for each transaction ",
				},
			},
		},
		{
			Name:   "status",
			Usage:  "Displays the health of various services running inside the node.",
			Action: s.Status,
			Flags:  []cli.Flag{},
			Hidden: true,
			Before: func(_ *cli.Context) error {
				s.Logger.Warnf("Command deprecated. Use `admin status` instead.")
				return nil
			},
		},
		{
			Name:   "profile",
			Usage:  "Collects profile metrics from the node.",
			Action: s.Profile,
			Flags: []cli.Flag{
				cli.Uint64Flag{
					Name:  "seconds, s",
					Usage: "duration of profile capture",
					Value: 8,
				},
				cli.StringFlag{
					Name:  "output_dir, o",
					Usage: "output directory of the captured profile",
					Value: "/tmp/",
				},
			},
			Hidden: true,
			Before: func(_ *cli.Context) error {
				s.Logger.Warnf("Command deprecated. Use `admin profile` instead.")
				return nil
			},
		},
		{
			Name:   "validate",
			Usage:  "Validate the TOML configuration and secrets that are passed as flags to the `node` command. Prints the full effective configuration, with defaults included",
			Action: s.ConfigFileValidate,
		},
		{
			Name:        "db",
			Usage:       "Commands for managing the database.",
			Description: "Potentially destructive commands for managing the database.",
			Subcommands: []cli.Command{
				{
					Name:   "reset",
					Usage:  "Drop, create and migrate database. Useful for setting up the database in order to run tests or resetting the dev database. WARNING: This will ERASE ALL DATA for the specified database, referred to by CL_DATABASE_URL env variable or by the Database.URL field in a secrets TOML config.",
					Hidden: safe,
					Action: s.ResetDatabase,
					Before: s.validateDB,
					Flags: []cli.Flag{
						cli.BoolFlag{
							Name:  "dangerWillRobinson",
							Usage: "set to true to enable dropping non-test databases",
						},
						cli.BoolFlag{
							Name:  "force",
							Usage: "set to true to force the reset by dropping any existing connections to the database",
						},
					},
				},
				{
					Name:   "preparetest",
					Usage:  "Reset database and load fixtures.",
					Hidden: safe,
					Action: s.PrepareTestDatabase,
					Before: s.validateDB,
					Flags: []cli.Flag{
						cli.BoolFlag{
							Name:  "user-only",
							Usage: "only include test user fixture",
						},
						cli.BoolFlag{
							Name:  "force",
							Usage: "set to true to force the reset by dropping any existing connections to the database",
						},
					},
				},
				{
					Name:   "version",
					Usage:  "Display the current database version.",
					Action: s.VersionDatabase,
					Before: s.validateDB,
					Flags:  []cli.Flag{},
				},
				{
					Name:   "status",
					Usage:  "Display the current database migration status.",
					Action: s.StatusDatabase,
					Before: s.validateDB,
					Flags:  []cli.Flag{},
				},
				{
					Name:   "migrate",
					Usage:  "Migrate the database to the latest version.",
					Action: s.MigrateDatabase,
					Before: s.validateDB,
					Flags:  []cli.Flag{},
				},
				{
					Name:   "rollback",
					Usage:  "Roll back the database to a previous <version>. Rolls back a single migration if no version specified.",
					Action: s.RollbackDatabase,
					Before: s.validateDB,
					Flags:  []cli.Flag{},
				},
				{
					Name:   "create-migration",
					Usage:  "Create a new migration.",
					Hidden: safe,
					Action: s.CreateMigration,
					Before: s.validateDB,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "type",
							Usage: "set to `go` to generate a .go migration (instead of .sql)",
						},
					},
				},
				{
					Name:    "delete-chain",
					Aliases: []string{},
					Usage:   "Commands for cleaning up chain specific db tables. WARNING: This will ERASE ALL chain specific data referred to by --type and --id options for the specified database, referred to by CL_DATABASE_URL env variable or by the Database.URL field in a secrets TOML config.",
					Action:  s.CleanupChainTables,
					Before:  s.validateDB,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:     "id",
							Usage:    "chain id based on which chain specific table cleanup will be done",
							Required: true,
						},
						cli.StringFlag{
							Name:     "type",
							Usage:    "chain type based on which table cleanup will be done, eg. EVM",
							Required: true,
						},
						cli.BoolFlag{
							Name:  "danger",
							Usage: "set to true to enable dropping non-test databases",
						},
					},
				},
			},
		},
		{
			Name:   "remove-blocks",
			Usage:  "Deletes block range and all associated data",
			Action: s.RemoveBlocks,
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:     "start",
					Usage:    "Beginning of block range to be deleted",
					Required: true,
				},
				cli.Int64Flag{
					Name:     "evm-chain-id",
					Usage:    "Chain ID of the EVM-based blockchain",
					Required: true,
				},
			},
		},
	}
}

// ownerPermsMask are the file permission bits reserved for owner.
const ownerPermsMask = os.FileMode(0o700)

// RunNode starts the Chainlink core.
func (s *Shell) RunNode(c *cli.Context) error {
	if err := s.runNode(c); err != nil {
		return s.errorOut(err)
	}
	return nil
}

func (s *Shell) runNode(c *cli.Context) error {
	ctx := s.ctx()
	lggr := logger.Sugared(s.Logger.Named("RunNode"))

	var pwd, vrfpwd *string
	if passwordFile := c.String("password"); passwordFile != "" {
		p, err := utils.PasswordFromFile(passwordFile)
		if err != nil {
			return errors.Wrap(err, "error reading password from file")
		}
		pwd = &p
	}
	if vrfPasswordFile := c.String("vrfpassword"); len(vrfPasswordFile) != 0 {
		p, err := utils.PasswordFromFile(vrfPasswordFile)
		if err != nil {
			return errors.Wrapf(err, "error reading VRF password from vrfpassword file \"%s\"", vrfPasswordFile)
		}
		vrfpwd = &p
	}

	s.Config.SetPasswords(pwd, vrfpwd)

	s.Config.LogConfiguration(lggr.Debugf, lggr.Warnf)

	if err := s.Config.Validate(); err != nil {
		return errors.Wrap(err, "config validation failed")
	}

	lggr.Infow(fmt.Sprintf("Starting Chainlink Node %s at commit %s", static.Version, static.Sha), "Version", static.Version, "SHA", static.Sha)

	if build.IsDev() {
		lggr.Warn("Chainlink is running in DEVELOPMENT mode. This is a security risk if enabled in production.")
	}

	if err := utils.EnsureDirAndMaxPerms(s.Config.RootDir(), os.FileMode(0700)); err != nil {
		return fmt.Errorf("failed to create root directory %q: %w", s.Config.RootDir(), err)
	}

	cfg := s.Config
	ldb := pg.NewLockedDB(cfg.AppID(), cfg.Database(), cfg.Database().Lock(), lggr)

	// rootCtx will be cancelled when SIGINT|SIGTERM is received
	rootCtx, cancelRootCtx := context.WithCancel(context.Background())

	// cleanExit is used to skip "fail fast" routine
	cleanExit := make(chan struct{})
	var shutdownStartTime time.Time
	defer func() {
		close(cleanExit)
		if !shutdownStartTime.IsZero() {
			log.Printf("Graceful shutdown time: %s", time.Since(shutdownStartTime))
		}
	}()

	go shutdown.HandleShutdown(func(sig string) {
		lggr.Infof("Shutting down due to %s signal received...", sig)

		shutdownStartTime = time.Now()
		cancelRootCtx()

		select {
		case <-cleanExit:
			return
		case <-time.After(s.Config.ShutdownGracePeriod()):
		}

		lggr.Criticalf("Shutdown grace period of %v exceeded, closing DB and exiting...", s.Config.ShutdownGracePeriod())
		// LockedDB.Close() will release DB locks and close DB connection
		// Executing this explicitly because defers are not executed in case of os.Exit()
		if err := ldb.Close(); err != nil {
			lggr.Criticalf("Failed to close LockedDB: %v", err)
		}
		lggr.Debug("Closed DB")
		if err := s.CloseLogger(); err != nil {
			log.Printf("Failed to close Logger: %v", err)
		}

		os.Exit(-1)
	})

	// Try opening DB connection and acquiring DB locks at once
	if err := ldb.Open(rootCtx); err != nil {
		// If not successful, we know neither locks nor connection remains opened
		return s.errorOut(errors.Wrap(err, "opening db"))
	}
	defer lggr.ErrorIfFn(ldb.Close, "Error closing db")

	// From now on, DB locks and DB connection will be released on every return.
	// Keep watching on logger.Fatal* calls and os.Exit(), because defer will not be executed.

	app, err := s.AppFactory.NewApplication(rootCtx, s.Config, s.Logger, s.Registerer, ldb.DB(), s.KeyStoreAuthenticator)
	if err != nil {
		return s.errorOut(errors.Wrap(err, "fatal error instantiating application"))
	}

	// Local shell initialization always uses local auth users table for admin auth
	authProviderORM := app.BasicAdminUsersORM()

	legacyEVMChains := app.GetRelayers().LegacyEVMChains()

	if s.Config.EVMEnabled() {
		// ensure any imported keys are imported
		for _, k := range s.Config.ImportedEthKeys().List() {
			lggr.Debug("Importing eth key")
			id, err2 := chain_selectors.GetChainIDFromSelector(k.ChainDetails().ChainSelector)
			if err != nil {
				return s.errorOut(errors.Wrapf(err2, "error getting chain id from selector when trying to import eth key %v", k.JSON()))
			}
			cid, _ := big.NewInt(0).SetString(id, 10)
			if cid == nil {
				return s.errorOut(fmt.Errorf("error converting chain id '%s' to big int", id))
			}
			_, err2 = app.GetKeyStore().Eth().Import(rootCtx, []byte(k.JSON()), k.Password(), cid)
			if err2 != nil {
				if errors.Is(err2, keystore.ErrKeyExists) {
					lggr.Debugf("Eth key %s already exists for chain %v", k.JSON(), k.ChainDetails())
					continue
				}
				return s.errorOut(errors.Wrap(err2, "error importing eth key"))
			}
			lggr.Debugf("Imported eth key %s for chain %v", k.JSON(), k.ChainDetails())
		}

		chainList, err2 := legacyEVMChains.List()
		if err2 != nil {
			return fmt.Errorf("error listing legacy evm chains: %w", err2)
		}
		for _, ch := range chainList {
			if ch.Config().EVM().AutoCreateKey() {
				lggr.Debugf("AutoCreateKey=true, will ensure EVM key for chain %s", ch.ID())
				err2 := app.GetKeyStore().Eth().EnsureKeys(rootCtx, ch.ID())
				if err2 != nil {
					return errors.Wrap(err2, "failed to ensure keystore keys")
				}
			} else {
				lggr.Debugf("AutoCreateKey=false, will not ensure EVM key for chain %s", ch.ID())
			}
		}
	}

	if s.Config.OCR().Enabled() {
		err2 := app.GetKeyStore().OCR().EnsureKey(rootCtx)
		if err2 != nil {
			return errors.Wrap(err2, "failed to ensure ocr key")
		}
	}
	if s.Config.OCR2().Enabled() {
		var enabledChains []chaintype.ChainType
		if s.Config.EVMEnabled() {
			enabledChains = append(enabledChains, chaintype.EVM)
		}
		if s.Config.CosmosEnabled() {
			enabledChains = append(enabledChains, chaintype.Cosmos)
		}
		if s.Config.SolanaEnabled() {
			enabledChains = append(enabledChains, chaintype.Solana)
		}
		if s.Config.StarkNetEnabled() {
			enabledChains = append(enabledChains, chaintype.StarkNet)
		}
		if s.Config.AptosEnabled() {
			enabledChains = append(enabledChains, chaintype.Aptos)
		}
		if s.Config.TronEnabled() {
			enabledChains = append(enabledChains, chaintype.Tron)
		}
		if s.Config.TONEnabled() {
			enabledChains = append(enabledChains, chaintype.TON)
		}
		err2 := app.GetKeyStore().OCR2().EnsureKeys(rootCtx, enabledChains...)
		if err2 != nil {
			return errors.Wrap(err2, "failed to ensure ocr key")
		}
	}

	if s.Config.P2P().Enabled() {
		if s.Config.ImportedP2PKey().JSON() != "" {
			lggr.Debugf("Importing p2p key %s", s.Config.ImportedP2PKey().JSON())
			_, err2 := app.GetKeyStore().P2P().Import(rootCtx, []byte(s.Config.ImportedP2PKey().JSON()), s.Config.ImportedP2PKey().Password())
			if errors.Is(err2, keystore.ErrKeyExists) {
				lggr.Debugf("P2P key already exists %s", s.Config.ImportedP2PKey().JSON())
			} else if err2 != nil {
				return s.errorOut(errors.Wrap(err2, "error importing p2p key"))
			}
		}
		err2 := app.GetKeyStore().P2P().EnsureKey(rootCtx)
		if err2 != nil {
			return errors.Wrap(err2, "failed to ensure p2p key")
		}
	}
	if s.Config.CosmosEnabled() {
		err2 := app.GetKeyStore().Cosmos().EnsureKey(rootCtx)
		if err2 != nil {
			return errors.Wrap(err2, "failed to ensure cosmos key")
		}
	}
	if s.Config.SolanaEnabled() {
		err2 := app.GetKeyStore().Solana().EnsureKey(rootCtx)
		if err2 != nil {
			return errors.Wrap(err2, "failed to ensure solana key")
		}
	}
	if s.Config.StarkNetEnabled() {
		err2 := app.GetKeyStore().StarkNet().EnsureKey(rootCtx)
		if err2 != nil {
			return errors.Wrap(err2, "failed to ensure starknet key")
		}
	}
	if s.Config.AptosEnabled() {
		err2 := app.GetKeyStore().Aptos().EnsureKey(rootCtx)
		if err2 != nil {
			return errors.Wrap(err2, "failed to ensure aptos key")
		}
	}
	if s.Config.TronEnabled() {
		err2 := app.GetKeyStore().Tron().EnsureKey(rootCtx)
		if err2 != nil {
			return errors.Wrap(err2, "failed to ensure tron key")
		}
	}
	if s.Config.TONEnabled() {
		err2 := app.GetKeyStore().TON().EnsureKey(rootCtx)
		if err2 != nil {
			return errors.Wrap(err2, "failed to ensure ton key")
		}
	}

	err2 := app.GetKeyStore().Workflow().EnsureKey(rootCtx)
	if err2 != nil {
		return errors.Wrap(err2, "failed to ensure workflow key")
	}

	err2 = app.GetKeyStore().CSA().EnsureKey(rootCtx)
	if err2 != nil {
		return errors.Wrap(err2, "failed to ensure CSA key")
	}

	if e := checkFilePermissions(lggr, s.Config.RootDir()); e != nil {
		lggr.Warn(e)
	}

	var user sessions.User
	if user, err = NewFileAPIInitializer(c.String("api")).Initialize(ctx, authProviderORM, lggr); err != nil {
		if !errors.Is(err, ErrNoCredentialFile) {
			return errors.Wrap(err, "error creating api initializer")
		}
		if user, err = s.FallbackAPIInitializer.Initialize(ctx, authProviderORM, lggr); err != nil {
			if errors.Is(err, ErrorNoAPICredentialsAvailable) {
				return errors.WithStack(err)
			}
			return errors.Wrap(err, "error creating fallback initializer")
		}
	}

	lggr.Info("API exposed for user ", user.Email)

	if err = app.Start(rootCtx); err != nil {
		// We do not try stopping any sub-services that might be started,
		// because the app will exit immediately upon return.
		// But LockedDB will be released by defer in above.
		return errors.Wrap(err, "error starting app")
	}

	grp, grpCtx := errgroup.WithContext(rootCtx)

	grp.Go(func() error {
		<-grpCtx.Done()
		if errInternal := app.Stop(); errInternal != nil {
			return errors.Wrap(errInternal, "error stopping app")
		}
		return nil
	})

	lggr.Infow(fmt.Sprintf("Chainlink booted in %.2fs", time.Since(static.InitTime).Seconds()), "appID", app.ID())

	grp.Go(func() error {
		errInternal := s.Runner.Run(grpCtx, app)
		if errors.Is(errInternal, http.ErrServerClosed) {
			errInternal = nil
		}
		// In tests we have custom runners that stop the app gracefully,
		// therefore we need to cancel rootCtx when the Runner has quit.
		cancelRootCtx()
		return errInternal
	})

	return grp.Wait()
}

func checkFilePermissions(lggr logger.Logger, rootDir string) error {
	// Ensure tls sub directory (and children) permissions are <= `ownerPermsMask``
	tlsDir := filepath.Join(rootDir, "tls")
	if _, err := os.Stat(tlsDir); err != nil && !os.IsNotExist(err) {
		lggr.Errorf("error checking perms of 'tls' directory: %v", err)
	} else if err == nil {
		err := utils.EnsureDirAndMaxPerms(tlsDir, ownerPermsMask)
		if err != nil {
			return err
		}

		err = filepath.Walk(tlsDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				lggr.Errorf(`error checking perms of "%v": %v`, path, err)
				return err
			}
			if utils.TooPermissive(info.Mode().Perm(), ownerPermsMask) {
				newPerms := info.Mode().Perm() & ownerPermsMask
				lggr.Warnf("%s has overly permissive file permissions, reducing them from %s to %s", path, info.Mode().Perm(), newPerms)
				return utils.EnsureFilepathMaxPerms(path, newPerms)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	// Ensure {secret,cookie} files' permissions are <= `ownerPermsMask``
	protectedFiles := []string{"secret", "cookie", ".password", ".env", ".api"}
	for _, fileName := range protectedFiles {
		path := filepath.Join(rootDir, fileName)
		fileInfo, err := os.Stat(path)
		if os.IsNotExist(err) {
			continue
		} else if err != nil {
			return err
		}
		if utils.TooPermissive(fileInfo.Mode().Perm(), ownerPermsMask) {
			newPerms := fileInfo.Mode().Perm() & ownerPermsMask
			lggr.Warnf("%s has overly permissive file permissions, reducing them from %s to %s", path, fileInfo.Mode().Perm(), newPerms)
			err = utils.EnsureFilepathMaxPerms(path, newPerms)
			if err != nil {
				return err
			}
		}
		owned, err := utils.IsFileOwnedByChainlink(fileInfo)
		if err != nil {
			lggr.Warn(err)
			continue
		}
		if !owned {
			lggr.Warnf("The file %v is not owned by the user running chainlink. This will be made mandatory in the future.", path)
		}
	}
	return nil
}

// RebroadcastTransactions run locally to force manual rebroadcasting of
// transactions in a given nonce range.
func (s *Shell) RebroadcastTransactions(c *cli.Context) (err error) {
	ctx := s.ctx()
	beginningNonce := c.Int64("beginningNonce")
	endingNonce := c.Int64("endingNonce")
	gasPriceWei := c.Uint64("gasPriceWei")
	overrideGasLimit := c.Uint("gasLimit")
	addressHex := c.String("address")
	chainIDStr := c.String("evmChainID")

	addressBytes, err := hexutil.Decode(addressHex)
	if err != nil {
		return s.errorOut(errors.Wrap(err, "could not decode address"))
	}
	address := gethCommon.BytesToAddress(addressBytes)

	var chainID *big.Int
	if chainIDStr != "" {
		var ok bool
		chainID, ok = big.NewInt(0).SetString(chainIDStr, 10)
		if !ok {
			return s.errorOut(errors.New("invalid evmChainID"))
		}
	}

	err = s.Config.Validate()
	if err != nil {
		return err
	}

	lggr := logger.Sugared(s.Logger.Named("RebroadcastTransactions"))
	db, err := pg.OpenUnlockedDB(ctx, s.Config.AppID(), s.Config.Database())
	if err != nil {
		return s.errorOut(errors.Wrap(err, "opening DB"))
	}
	defer lggr.ErrorIfFn(db.Close, "Error closing db")

	app, err := s.AppFactory.NewApplication(ctx, s.Config, lggr, s.Registerer, db, s.KeyStoreAuthenticator)
	if err != nil {
		return s.errorOut(errors.Wrap(err, "fatal error instantiating application"))
	}

	// TODO: BCF-2511 once the dust settles on BCF-2440/1 evaluate how the
	// [loop.Relayer] interface needs to be extended to support programming similar to
	// this pattern but in a chain-agnostic way
	chain, err := app.GetRelayers().LegacyEVMChains().Get(chainID.String())
	if err != nil {
		return s.errorOut(err)
	}
	keyStore := app.GetKeyStore()

	ethClient := chain.Client()

	err = ethClient.Dial(ctx)
	if err != nil {
		return err
	}

	if c.IsSet("password") {
		pwd, err2 := utils.PasswordFromFile(c.String("password"))
		if err2 != nil {
			return s.errorOut(fmt.Errorf("error reading password: %w", err2))
		}
		s.Config.SetPasswords(&pwd, nil)
	}

	err = s.Config.Validate()
	if err != nil {
		return s.errorOut(fmt.Errorf("error validating configuration: %w", err))
	}

	err = keyStore.Unlock(ctx, s.Config.Password().Keystore())
	if err != nil {
		return s.errorOut(errors.Wrap(err, "error authenticating keystore"))
	}

	if err = keyStore.Eth().CheckEnabled(ctx, address, chain.ID()); err != nil {
		return s.errorOut(err)
	}

	ks := keys.NewChainStore(keystore.NewEthSigner(keyStore.Eth(), chain.ID()), chain.ID())

	s.Logger.Infof("Rebroadcasting transactions from %v to %v", beginningNonce, endingNonce)

	orm := txmgr.NewTxStore(app.GetDB(), lggr)
	txBuilder := txmgr.NewEvmTxAttemptBuilder(*ethClient.ConfiguredChainID(), chain.Config().EVM().GasEstimator(), ks, nil)
	feeCfg := txmgr.NewEvmTxmFeeConfig(chain.Config().EVM().GasEstimator())
	stuckTxDetector := txmgr.NewStuckTxDetector(lggr, ethClient.ConfiguredChainID(), "", assets.NewWei(assets.NewEth(100).ToInt()), chain.Config().EVM().Transactions().AutoPurge(), nil, orm, ethClient)
	metrics, err := txmgr.NewEVMTxmMetrics(ethClient.ConfiguredChainID().String())
	if err != nil {
		return s.errorOut(err)
	}
	ec := txmgr.NewEvmConfirmer(orm, txmgr.NewEvmTxmClient(ethClient, chain.Config().EVM().NodePool().Errors()),
		feeCfg, chain.Config().EVM().Transactions(), app.GetConfig().Database(), ks, txBuilder, chain.Logger(), stuckTxDetector, metrics)
	totalNonces := endingNonce - beginningNonce + 1
	nonces := make([]evmtypes.Nonce, totalNonces)
	for i := int64(0); i < totalNonces; i++ {
		nonces[i] = evmtypes.Nonce(beginningNonce + i)
	}
	if gasPriceWei <= math.MaxInt64 {
		return s.errorOut(ec.ForceRebroadcast(ctx, nonces, gas.EvmFee{GasPrice: assets.NewWeiI(int64(gasPriceWei))}, address, uint64(overrideGasLimit)))
	}
	return s.errorOut(fmt.Errorf("integer overflow conversion error. GasPrice: %v", gasPriceWei))
}

type HealthCheckPresenter struct {
	webPresenters.Check
}

func (p *HealthCheckPresenter) ToRow() []string {
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	var status string

	switch p.Status {
	case web.HealthStatusFailing:
		status = red(p.Status)
	case web.HealthStatusPassing:
		status = green(p.Status)
	}

	return []string{
		p.Name,
		status,
		p.Output,
	}
}

type HealthCheckPresenters []HealthCheckPresenter

// RenderTable implements TableRenderer
func (ps HealthCheckPresenters) RenderTable(rt RendererTable) error {
	headers := []string{"Name", "Status", "Output"}
	rows := [][]string{}

	for _, p := range ps {
		rows = append(rows, p.ToRow())
	}

	renderList(headers, rows, rt.Writer)

	return nil
}

var errDBURLMissing = errors.New("You must set CL_DATABASE_URL env variable or provide a secrets TOML with Database.URL set. HINT: If you are running this to set up your local test database, try CL_DATABASE_URL=postgresql://postgres@localhost:5432/chainlink_test?sslmode=disable")

// ConfigValidate validate the client configuration and pretty-prints results
func (s *Shell) ConfigFileValidate(_ *cli.Context) error {
	fn := func(f string, params ...any) { fmt.Printf(f, params...) }
	s.Config.LogConfiguration(fn, fn)
	if err := s.configExitErr(s.Config.Validate); err != nil {
		return err
	}
	fmt.Println("Valid configuration.")
	return nil
}

// ValidateDB is a BeforeFunc to run prior to database sub commands
// the ctx must be that of the last subcommand to be validated
func (s *Shell) validateDB(c *cli.Context) error {
	return s.configExitErr(s.Config.ValidateDB)
}

// ctx returns a context.Context that will be cancelled when SIGINT|SIGTERM is received
func (s *Shell) ctx() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	go shutdown.HandleShutdown(func(_ string) { cancel() })
	return ctx
}

// ResetDatabase drops, creates and migrates the database specified by CL_DATABASE_URL or Database.URL
// in secrets TOML. This is useful to set up the database for testing
func (s *Shell) ResetDatabase(c *cli.Context) error {
	ctx := s.ctx()
	cfg := s.Config.Database()
	u := cfg.URL()
	if u.String() == "" {
		return s.errorOut(errDBURLMissing)
	}
	dangerMode := c.Bool("dangerWillRobinson")
	dbname := u.Path[1:]
	if !dangerMode && !strings.HasSuffix(dbname, "_test") {
		return s.errorOut(fmt.Errorf("cannot reset database named `%s`. This command can only be run against databases with a name that ends in `_test`, to prevent accidental data loss. If you REALLY want to reset this database, pass in the -dangerWillRobinson option", dbname))
	}

	force := c.Bool("force")

	if err := store.ResetDatabase(ctx, s.Logger, cfg, force); err != nil {
		return s.errorOut(err)
	}
	return nil
}

// PrepareTestDatabase calls ResetDatabase then loads fixtures required for tests
func (s *Shell) PrepareTestDatabase(c *cli.Context) error {
	if err := s.ResetDatabase(c); err != nil {
		return s.errorOut(err)
	}
	cfg := s.Config

	// Creating pristine DB copy to speed up FullTestDB
	dbUrl := cfg.Database().URL()
	userOnly := c.Bool("user-only")
	if err := store.PrepareTestDB(s.Logger, dbUrl, userOnly); err != nil {
		return s.errorOut(err)
	}
	return nil
}

// MigrateDatabase migrates the database
func (s *Shell) MigrateDatabase(_ *cli.Context) error {
	ctx := s.ctx()
	cfg := s.Config.Database()
	parsed := cfg.URL()
	if parsed.String() == "" {
		return s.errorOut(errDBURLMissing)
	}

	err := migrate.SetMigrationENVVars(s.Config.EVMConfigs())
	if err != nil {
		return err
	}

	s.Logger.Infof("Migrating database: %#v", parsed.String())
	if err := migrateDB(ctx, cfg); err != nil {
		return s.errorOut(err)
	}
	return nil
}

// RollbackDatabase rolls back the database via down migrations.
func (s *Shell) RollbackDatabase(c *cli.Context) error {
	ctx := s.ctx()
	var version null.Int
	if c.Args().Present() {
		arg := c.Args().First()
		numVersion, err := strconv.ParseInt(arg, 10, 64)
		if err != nil {
			return s.errorOut(errors.Errorf("Unable to parse %v as integer", arg))
		}
		version = null.IntFrom(numVersion)
	}

	db, err := store.NewConnection(ctx, s.Config.Database())
	if err != nil {
		return fmt.Errorf("failed to initialize orm: %w", err)
	}

	if err := migrate.Rollback(ctx, db.DB, version); err != nil {
		return fmt.Errorf("migrateDB failed: %w", err)
	}

	return nil
}

// VersionDatabase displays the current database version.
func (s *Shell) VersionDatabase(_ *cli.Context) error {
	ctx := s.ctx()
	db, err := store.NewConnection(ctx, s.Config.Database())
	if err != nil {
		return fmt.Errorf("failed to initialize orm: %w", err)
	}

	version, err := migrate.Current(ctx, db.DB)
	if err != nil {
		return fmt.Errorf("migrateDB failed: %w", err)
	}

	s.Logger.Infof("Database version: %v", version)
	return nil
}

// StatusDatabase displays the database migration status
func (s *Shell) StatusDatabase(_ *cli.Context) error {
	ctx := s.ctx()
	db, err := store.NewConnection(ctx, s.Config.Database())
	if err != nil {
		return fmt.Errorf("failed to initialize orm: %w", err)
	}

	if err = migrate.Status(ctx, db.DB); err != nil {
		return fmt.Errorf("Status failed: %w", err)
	}
	return nil
}

// CreateMigration displays the database migration status
func (s *Shell) CreateMigration(c *cli.Context) error {
	ctx := s.ctx()
	if !c.Args().Present() {
		return s.errorOut(errors.New("You must specify a migration name"))
	}
	db, err := store.NewConnection(ctx, s.Config.Database())
	if err != nil {
		return fmt.Errorf("failed to initialize orm: %w", err)
	}

	migrationType := c.String("type")
	if migrationType != "go" {
		migrationType = "sql"
	}

	if err = migrate.Create(db.DB, c.Args().First(), migrationType); err != nil {
		return fmt.Errorf("Status failed: %w", err)
	}
	return nil
}

// CleanupChainTables deletes database table rows based on chain type and chain id input.
func (s *Shell) CleanupChainTables(c *cli.Context) error {
	ctx := s.ctx()
	cfg := s.Config.Database()
	parsed := cfg.URL()
	if parsed.String() == "" {
		return s.errorOut(errDBURLMissing)
	}

	dbname := parsed.Path[1:]
	if !c.Bool("danger") && !strings.HasSuffix(dbname, "_test") {
		return s.errorOut(fmt.Errorf("cannot reset database named `%s`. This command can only be run against databases with a name that ends in `_test`, to prevent accidental data loss. If you really want to delete chain specific data from this database, pass in the --danger option", dbname))
	}

	db, err := store.NewConnection(ctx, cfg)
	if err != nil {
		return s.errorOut(errors.Wrap(err, "error connecting to the database"))
	}
	defer db.Close()

	// some tables with evm_chain_id (mostly job specs) are in public schema
	tablesToDeleteFromQuery := `SELECT table_name, table_schema FROM information_schema.columns WHERE "column_name"=$1;`
	// Delete rows from each table based on the chain_id.
	if !strings.EqualFold("EVM", c.String("type")) {
		return s.errorOut(errors.New("unknown chain type"))
	}
	rows, err := db.Query(tablesToDeleteFromQuery, "evm_chain_id")
	if err != nil {
		return err
	}
	defer rows.Close()

	var tablesToDeleteFrom []string
	for rows.Next() {
		var name string
		var schema string
		if err = rows.Scan(&name, &schema); err != nil {
			return err
		}
		tablesToDeleteFrom = append(tablesToDeleteFrom, schema+"."+name)
	}
	if rows.Err() != nil {
		return rows.Err()
	}

	for _, tableName := range tablesToDeleteFrom {
		query := fmt.Sprintf(`DELETE FROM %s WHERE "evm_chain_id"=$1;`, tableName)
		_, err = db.Exec(query, c.String("id"))
		if err != nil {
			fmt.Printf("Error deleting rows containing evm_chain_id from %s: %v\n", tableName, err)
		} else {
			fmt.Printf("Rows with evm_chain_id %s deleted from %s.\n", c.String("id"), tableName)
		}
	}
	return nil
}

func migrateDB(ctx context.Context, config store.Config) error {
	db, err := store.NewConnection(ctx, config)
	if err != nil {
		return fmt.Errorf("failed to initialize orm: %w", err)
	}

	if err = migrate.Migrate(ctx, db.DB); err != nil {
		return fmt.Errorf("migrateDB failed: %w", err)
	}
	return db.Close()
}

// RemoveBlocks - removes blocks after the specified blocks number
func (s *Shell) RemoveBlocks(c *cli.Context) error {
	start := c.Int64("start")
	if start <= 0 {
		return s.errorOut(errors.New("Must pass a positive value in '--start' parameter"))
	}

	chainID := big.NewInt(0)
	if c.IsSet("evm-chain-id") {
		err := chainID.UnmarshalText([]byte(c.String("evm-chain-id")))
		if err != nil {
			return s.errorOut(err)
		}
	}

	cfg := s.Config
	err := cfg.Validate()
	if err != nil {
		return s.errorOut(fmt.Errorf("error validating configuration: %w", err))
	}

	lggr := logger.Sugared(s.Logger.Named("RemoveBlocks"))
	ldb := pg.NewLockedDB(cfg.AppID(), cfg.Database(), cfg.Database().Lock(), lggr)
	ctx, cancel := context.WithCancel(context.Background())
	go shutdown.HandleShutdown(func(sig string) {
		cancel()
		lggr.Info("received signal to stop - closing the database and releasing lock")

		if cErr := ldb.Close(); cErr != nil {
			lggr.Criticalf("Failed to close LockedDB: %v", cErr)
		}

		if cErr := s.CloseLogger(); cErr != nil {
			log.Printf("Failed to close Logger: %v", cErr)
		}
	})

	if err = ldb.Open(ctx); err != nil {
		// If not successful, we know neither locks nor connection remains opened
		return s.errorOut(errors.Wrap(err, "opening db"))
	}
	defer lggr.ErrorIfFn(ldb.Close, "Error closing db")

	// From now on, DB locks and DB connection will be released on every return.
	// Keep watching on logger.Fatal* calls and os.Exit(), because defer will not be executed.

	app, err := s.AppFactory.NewApplication(ctx, s.Config, s.Logger, s.Registerer, ldb.DB(), s.KeyStoreAuthenticator)
	if err != nil {
		return s.errorOut(errors.Wrap(err, "fatal error instantiating application"))
	}

	err = app.DeleteLogPollerDataAfter(ctx, chainID, start)
	if err != nil {
		return s.errorOut(err)
	}

	lggr.Infof("RemoveBlocks: successfully removed blocks")

	return nil
}
