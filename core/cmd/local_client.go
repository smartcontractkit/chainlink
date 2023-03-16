package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/fatih/color"
	"github.com/kylelemons/godebug/diff"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	clipkg "github.com/urfave/cli"
	"go.uber.org/multierr"
	"golang.org/x/sync/errgroup"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/sessions"
	"github.com/smartcontractkit/chainlink/core/shutdown"
	"github.com/smartcontractkit/chainlink/core/static"
	"github.com/smartcontractkit/chainlink/core/store/dialects"
	"github.com/smartcontractkit/chainlink/core/store/migrate"
	"github.com/smartcontractkit/chainlink/core/utils"
	webPresenters "github.com/smartcontractkit/chainlink/core/web/presenters"
)

func initLocalSubCmds(client *Client, devMode bool) []cli.Command {
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
			Action: client.RunNode,
		},
		{
			Name:   "rebroadcast-transactions",
			Usage:  "Manually rebroadcast txs matching nonce range with the specified gas price. This is useful in emergencies e.g. high gas prices and/or network congestion to forcibly clear out the pending TX queue",
			Action: client.RebroadcastTransactions,
			Flags: []cli.Flag{
				cli.Uint64Flag{
					Name:  "beginningNonce, b",
					Usage: "beginning of nonce range to rebroadcast",
				},
				cli.Uint64Flag{
					Name:  "endingNonce, e",
					Usage: "end of nonce range to rebroadcast (inclusive)",
				},
				cli.Uint64Flag{
					Name:  "gasPriceWei, g",
					Usage: "gas price (in Wei) to rebroadcast transactions at",
				},
				cli.StringFlag{
					Name:  "password, p",
					Usage: "text file holding the password for the node's account",
				},
				cli.StringFlag{
					Name:  "address, a",
					Usage: "The address (in hex format) for the key which we want to rebroadcast transactions",
				},
				cli.StringFlag{
					Name:  "evmChainID",
					Usage: "Chain ID for which to rebroadcast transactions. If left blank, ETH_CHAIN_ID will be used.",
				},
				cli.Uint64Flag{
					Name:  "gasLimit",
					Usage: "OPTIONAL: gas limit to use for each transaction ",
				},
			},
		},
		{
			Name:   "status",
			Usage:  "Displays the health of various services running inside the node.",
			Action: client.Status,
			Flags:  []cli.Flag{},
		},
		{
			Name:   "profile",
			Usage:  "Collects profile metrics from the node.",
			Action: client.Profile,
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
		},
		{
			Name:        "db",
			Usage:       "Commands for managing the database.",
			Description: "Potentially destructive commands for managing the database.",
			Subcommands: []cli.Command{
				{
					Name:   "reset",
					Usage:  "Drop, create and migrate database. Useful for setting up the database in order to run tests or resetting the dev database. WARNING: This will ERASE ALL DATA for the specified DATABASE_URL.",
					Hidden: !devMode,
					Action: client.ResetDatabase,
					Flags: []cli.Flag{
						cli.BoolFlag{
							Name:  "dangerWillRobinson",
							Usage: "set to true to enable dropping non-test databases",
						},
					},
				},
				{
					Name:   "preparetest",
					Usage:  "Reset database and load fixtures.",
					Hidden: !devMode,
					Action: client.PrepareTestDatabase,
					Flags: []cli.Flag{
						cli.BoolFlag{
							Name:  "user-only",
							Usage: "only include test user fixture",
						},
					},
				},
				{
					Name:   "version",
					Usage:  "Display the current database version.",
					Action: client.VersionDatabase,
					Flags:  []cli.Flag{},
				},
				{
					Name:   "status",
					Usage:  "Display the current database migration status.",
					Action: client.StatusDatabase,
					Flags:  []cli.Flag{},
				},
				{
					Name:   "migrate",
					Usage:  "Migrate the database to the latest version.",
					Action: client.MigrateDatabase,
					Flags:  []cli.Flag{},
				},
				{
					Name:   "rollback",
					Usage:  "Roll back the database to a previous <version>. Rolls back a single migration if no version specified.",
					Action: client.RollbackDatabase,
					Flags:  []cli.Flag{},
				},
				{
					Name:   "create-migration",
					Usage:  "Create a new migration.",
					Hidden: !devMode,
					Action: client.CreateMigration,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "type",
							Usage: "set to `go` to generate a .go migration (instead of .sql)",
						},
					},
				},
			},
		},
	}
}

// ownerPermsMask are the file permission bits reserved for owner.
const ownerPermsMask = os.FileMode(0o700)

// PristineDBName is a clean copy of test DB with migrations.
// Used by heavyweight.FullTestDB* functions.
const (
	PristineDBName   = "chainlink_test_pristine"
	TestDBNamePrefix = "chainlink_test_"
)

// RunNode starts the Chainlink core.
func (cli *Client) RunNode(c *clipkg.Context) error {
	if err := cli.runNode(c); err != nil {
		return cli.errorOut(err)
	}
	return nil
}

func (cli *Client) runNode(c *clipkg.Context) error {
	lggr := logger.Sugared(cli.Logger.Named("RunNode"))

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

	cli.Config.SetPasswords(pwd, vrfpwd)

	err := cli.Config.Validate()
	if err != nil {
		return errors.Wrap(err, "config validation failed")
	}

	lggr.Infow(fmt.Sprintf("Starting Chainlink Node %s at commit %s", static.Version, static.Sha), "Version", static.Version, "SHA", static.Sha)

	if cli.Config.Dev() {
		lggr.Warn("Chainlink is running in DEVELOPMENT mode. This is a security risk if enabled in production.")
	}

	ldb := pg.NewLockedDB(cli.Config, lggr)

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
		case <-time.After(cli.Config.ShutdownGracePeriod()):
		}

		lggr.Criticalf("Shutdown grace period of %v exceeded, closing DB and exiting...", cli.Config.ShutdownGracePeriod())
		// LockedDB.Close() will release DB locks and close DB connection
		// Executing this explicitly because defers are not executed in case of os.Exit()
		if err = ldb.Close(); err != nil {
			lggr.Criticalf("Failed to close LockedDB: %v", err)
		}
		if err = cli.CloseLogger(); err != nil {
			log.Printf("Failed to close Logger: %v", err)
		}

		os.Exit(-1)
	})

	// Try opening DB connection and acquiring DB locks at once
	if err = ldb.Open(rootCtx); err != nil {
		// If not successful, we know neither locks nor connection remains opened
		return cli.errorOut(errors.Wrap(err, "opening db"))
	}
	defer lggr.ErrorIfFn(ldb.Close, "Error closing db")

	// From now on, DB locks and DB connection will be released on every return.
	// Keep watching on logger.Fatal* calls and os.Exit(), because defer will not be executed.

	app, err := cli.AppFactory.NewApplication(rootCtx, cli.Config, cli.Logger, ldb.DB())
	if err != nil {
		return cli.errorOut(errors.Wrap(err, "fatal error instantiating application"))
	}

	sessionORM := app.SessionORM()
	keyStore := app.GetKeyStore()
	err = cli.KeyStoreAuthenticator.authenticate(keyStore, cli.Config)
	if err != nil {
		return errors.Wrap(err, "error authenticating keystore")
	}

	evmChainSet := app.GetChains().EVM
	// By passing in a function we can be lazy trying to look up a default
	// chain - if there are no existing keys, there is no need to check for
	// a chain ID
	DefaultEVMChainIDFunc := func() (*big.Int, error) {
		def, err2 := evmChainSet.Default()
		if err2 != nil {
			return nil, errors.Wrap(err2, "cannot get default EVM chain ID; no default EVM chain available")
		}
		return def.ID(), nil
	}
	err = keyStore.Migrate(cli.Config.VRFPassword(), DefaultEVMChainIDFunc)

	if cli.Config.EVMEnabled() {
		if err != nil {
			return errors.Wrap(err, "error migrating keystore")
		}

		for _, ch := range evmChainSet.Chains() {
			err2 := app.GetKeyStore().Eth().EnsureKeys(ch.ID())
			if err2 != nil {
				return errors.Wrap(err2, "failed to ensure keystore keys")
			}
		}
	}

	if cli.Config.FeatureOffchainReporting() {
		err2 := app.GetKeyStore().OCR().EnsureKey()
		if err2 != nil {
			return errors.Wrap(err2, "failed to ensure ocr key")
		}
	}
	if cli.Config.FeatureOffchainReporting2() {
		err2 := app.GetKeyStore().OCR2().EnsureKeys()
		if err2 != nil {
			return errors.Wrap(err2, "failed to ensure ocr key")
		}
	}
	if cli.Config.P2PEnabled() {
		err2 := app.GetKeyStore().P2P().EnsureKey()
		if err2 != nil {
			return errors.Wrap(err2, "failed to ensure p2p key")
		}
	}
	if cli.Config.SolanaEnabled() {
		err2 := app.GetKeyStore().Solana().EnsureKey()
		if err2 != nil {
			return errors.Wrap(err2, "failed to ensure solana key")
		}
	}
	if cli.Config.StarkNetEnabled() {
		err2 := app.GetKeyStore().StarkNet().EnsureKey()
		if err2 != nil {
			return errors.Wrap(err2, "failed to ensure starknet key")
		}
	}

	err2 := app.GetKeyStore().CSA().EnsureKey()
	if err2 != nil {
		return errors.Wrap(err2, "failed to ensure CSA key")
	}

	if e := checkFilePermissions(lggr, cli.Config.RootDir()); e != nil {
		lggr.Warn(e)
	}

	var user sessions.User
	if user, err = NewFileAPIInitializer(c.String("api")).Initialize(sessionORM, lggr); err != nil {
		if !errors.Is(err, ErrNoCredentialFile) {
			return errors.Wrap(err, "error creating api initializer")
		}
		if user, err = cli.FallbackAPIInitializer.Initialize(sessionORM, lggr); err != nil {
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

	cli.Config.LogConfiguration(lggr.Debug)

	lggr.Infow(fmt.Sprintf("Chainlink booted in %.2fs", time.Since(static.InitTime).Seconds()), "appID", app.ID())

	grp.Go(func() error {
		errInternal := cli.Runner.Run(grpCtx, app)
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
	// Ensure `$CLROOT/tls` directory (and children) permissions are <= `ownerPermsMask``
	tlsDir := filepath.Join(rootDir, "tls")
	_, err := os.Stat(tlsDir)
	if err != nil && !os.IsNotExist(err) {
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

	// Ensure `$CLROOT/{secret,cookie}` files' permissions are <= `ownerPermsMask``
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
func (cli *Client) RebroadcastTransactions(c *clipkg.Context) (err error) {
	beginningNonce := c.Uint("beginningNonce")
	endingNonce := c.Uint("endingNonce")
	gasPriceWei := c.Uint64("gasPriceWei")
	overrideGasLimit := c.Uint("gasLimit")
	addressHex := c.String("address")
	chainIDStr := c.String("evmChainID")

	addressBytes, err := hexutil.Decode(addressHex)
	if err != nil {
		return cli.errorOut(errors.Wrap(err, "could not decode address"))
	}
	address := gethCommon.BytesToAddress(addressBytes)

	var chainID *big.Int
	if chainIDStr != "" {
		var ok bool
		chainID, ok = big.NewInt(0).SetString(chainIDStr, 10)
		if !ok {
			return cli.errorOut(errors.New("invalid evmChainID"))
		}
	}

	lggr := logger.Sugared(cli.Logger.Named("RebroadcastTransactions"))
	db, err := pg.OpenUnlockedDB(cli.Config)
	if err != nil {
		return cli.errorOut(errors.Wrap(err, "opening DB"))
	}
	defer lggr.ErrorIfFn(db.Close, "Error closing db")

	app, err := cli.AppFactory.NewApplication(context.TODO(), cli.Config, lggr, db)
	if err != nil {
		return cli.errorOut(errors.Wrap(err, "fatal error instantiating application"))
	}
	defer func() {
		if serr := app.Stop(); serr != nil {
			err = multierr.Append(err, serr)
		}
	}()
	pwd, err := utils.PasswordFromFile(c.String("password"))
	if err != nil {
		return cli.errorOut(fmt.Errorf("error reading password: %+v", err))
	}
	chain, err := app.GetChains().EVM.Get(chainID)
	if err != nil {
		return cli.errorOut(err)
	}
	keyStore := app.GetKeyStore()

	ethClient := chain.Client()

	err = ethClient.Dial(context.TODO())
	if err != nil {
		return err
	}

	err = keyStore.Unlock(pwd)
	if err != nil {
		return cli.errorOut(errors.Wrap(err, "error authenticating keystore"))
	}

	cli.Logger.Infof("Rebroadcasting transactions from %v to %v", beginningNonce, endingNonce)

	keyStates, err := keyStore.Eth().GetStatesForChain(chain.ID())
	if err != nil {
		return cli.errorOut(err)
	}

	orm := txmgr.NewORM(app.GetSqlxDB(), lggr, cli.Config)
	ec := txmgr.NewEthConfirmer(orm, ethClient, chain.Config(), keyStore.Eth(), keyStates, nil, nil, chain.Logger())
	err = ec.ForceRebroadcast(beginningNonce, endingNonce, gasPriceWei, address, uint32(overrideGasLimit))
	return cli.errorOut(err)
}

type HealthCheckPresenter struct {
	webPresenters.Check
}

func (p *HealthCheckPresenter) ToRow() []string {
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	var status string

	switch p.Status {
	case services.StatusFailing:
		status = red(p.Status)
	case services.StatusPassing:
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

// Status will display the health of various services
func (cli *Client) Status(c *clipkg.Context) error {
	resp, err := cli.HTTP.Get("/health?full=1", nil)
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	return cli.renderAPIResponse(resp, &HealthCheckPresenters{})
}

// ResetDatabase drops, creates and migrates the database specified by DATABASE_URL
// This is useful to setup the database for testing
func (cli *Client) ResetDatabase(c *clipkg.Context) error {
	cfg := cli.Config
	parsed := cfg.DatabaseURL()
	if parsed.String() == "" {
		return cli.errorOut(errors.New("You must set DATABASE_URL env variable. HINT: If you are running this to set up your local test database, try DATABASE_URL=postgresql://postgres@localhost:5432/chainlink_test?sslmode=disable"))
	}

	dangerMode := c.Bool("dangerWillRobinson")

	dbname := parsed.Path[1:]
	if !dangerMode && !strings.HasSuffix(dbname, "_test") {
		return cli.errorOut(fmt.Errorf("cannot reset database named `%s`. This command can only be run against databases with a name that ends in `_test`, to prevent accidental data loss. If you REALLY want to reset this database, pass in the -dangerWillRobinson option", dbname))
	}
	lggr := cli.Logger
	lggr.Infof("Resetting database: %#v", parsed.String())
	lggr.Debugf("Dropping and recreating database: %#v", parsed.String())
	if err := dropAndCreateDB(parsed); err != nil {
		return cli.errorOut(err)
	}
	lggr.Debugf("Migrating database: %#v", parsed.String())
	if err := migrateDB(cfg, lggr); err != nil {
		return cli.errorOut(err)
	}
	schema, err := dumpSchema(parsed)
	if err != nil {
		return cli.errorOut(err)
	}
	lggr.Debugf("Testing rollback and re-migrate for database: %#v", parsed.String())
	var baseVersionID int64 = 54
	if err := downAndUpDB(cfg, lggr, baseVersionID); err != nil {
		return cli.errorOut(err)
	}
	if err := checkSchema(parsed, schema); err != nil {
		return cli.errorOut(err)
	}
	return nil
}

// PrepareTestDatabase calls ResetDatabase then loads fixtures required for tests
func (cli *Client) PrepareTestDatabase(c *clipkg.Context) error {
	if err := cli.ResetDatabase(c); err != nil {
		return cli.errorOut(err)
	}
	cfg := cli.Config

	// Creating pristine DB copy to speed up FullTestDB
	dbUrl := cfg.DatabaseURL()
	db, err := sqlx.Open(string(dialects.Postgres), dbUrl.String())
	if err != nil {
		return cli.errorOut(err)
	}
	defer db.Close()
	templateDB := strings.Trim(dbUrl.Path, "/")
	if err = dropAndCreatePristineDB(db, templateDB); err != nil {
		return cli.errorOut(err)
	}

	userOnly := c.Bool("user-only")
	fixturePath := "../store/fixtures/fixtures.sql"
	if userOnly {
		fixturePath = "../store/fixtures/users_only_fixture.sql"
	}
	if err := insertFixtures(dbUrl, fixturePath); err != nil {
		return cli.errorOut(err)
	}

	return cli.errorOut(dropDanglingTestDBs(cli.Logger, db))
}

func dropDanglingTestDBs(lggr logger.Logger, db *sqlx.DB) (err error) {
	// Drop all old dangling databases
	var dbs []string
	if err = db.Select(&dbs, `SELECT datname FROM pg_database WHERE datistemplate = false;`); err != nil {
		return err
	}

	// dropping database is very slow in postgres so we parallelise it here
	nWorkers := 25
	ch := make(chan string)
	var wg sync.WaitGroup
	wg.Add(nWorkers)
	errCh := make(chan error, len(dbs))
	for i := 0; i < nWorkers; i++ {
		go func() {
			defer wg.Done()
			for dbname := range ch {
				lggr.Infof("Dropping old, dangling test database: %q", dbname)
				gerr := utils.JustError(db.Exec(fmt.Sprintf(`DROP DATABASE IF EXISTS %s`, dbname)))
				errCh <- gerr
			}
		}()
	}
	for _, dbname := range dbs {
		if strings.HasPrefix(dbname, TestDBNamePrefix) && !strings.HasSuffix(dbname, "_pristine") {
			ch <- dbname
		}
	}
	close(ch)
	wg.Wait()
	close(errCh)
	for gerr := range errCh {
		err = multierr.Append(err, gerr)
	}
	return
}

// PrepareTestDatabase calls ResetDatabase then loads fixtures required for local
// testing against testnets. Does not include fake chain fixtures.
func (cli *Client) PrepareTestDatabaseUserOnly(c *clipkg.Context) error {
	if err := cli.ResetDatabase(c); err != nil {
		return cli.errorOut(err)
	}
	cfg := cli.Config
	if err := insertFixtures(cfg.DatabaseURL(), "../store/fixtures/users_only_fixtures.sql"); err != nil {
		return cli.errorOut(err)
	}
	return nil
}

// MigrateDatabase migrates the database
func (cli *Client) MigrateDatabase(c *clipkg.Context) error {
	cfg := cli.Config
	parsed := cfg.DatabaseURL()
	if parsed.String() == "" {
		return cli.errorOut(errors.New("You must set DATABASE_URL env variable. HINT: If you are running this to set up your local test database, try DATABASE_URL=postgresql://postgres@localhost:5432/chainlink_test?sslmode=disable"))
	}

	cli.Logger.Infof("Migrating database: %#v", parsed.String())
	if err := migrateDB(cfg, cli.Logger); err != nil {
		return cli.errorOut(err)
	}
	return nil
}

// VersionDatabase displays the current database version.
func (cli *Client) RollbackDatabase(c *clipkg.Context) error {
	var version null.Int
	if c.Args().Present() {
		arg := c.Args().First()
		numVersion, err := strconv.ParseInt(arg, 10, 64)
		if err != nil {
			return cli.errorOut(errors.Errorf("Unable to parse %v as integer", arg))
		}
		version = null.IntFrom(numVersion)
	}

	db, err := newConnection(cli.Config)
	if err != nil {
		return fmt.Errorf("failed to initialize orm: %v", err)
	}

	if err := migrate.Rollback(db.DB, cli.Logger, version); err != nil {
		return fmt.Errorf("migrateDB failed: %v", err)
	}

	return nil
}

// VersionDatabase displays the current database version.
func (cli *Client) VersionDatabase(c *clipkg.Context) error {
	db, err := newConnection(cli.Config)
	if err != nil {
		return fmt.Errorf("failed to initialize orm: %v", err)
	}

	version, err := migrate.Current(db.DB, cli.Logger)
	if err != nil {
		return fmt.Errorf("migrateDB failed: %v", err)
	}

	cli.Logger.Infof("Database version: %v", version)
	return nil
}

// StatusDatabase displays the database migration status
func (cli *Client) StatusDatabase(c *clipkg.Context) error {
	db, err := newConnection(cli.Config)
	if err != nil {
		return fmt.Errorf("failed to initialize orm: %v", err)
	}

	if err = migrate.Status(db.DB, cli.Logger); err != nil {
		return fmt.Errorf("Status failed: %v", err)
	}
	return nil
}

// CreateMigration displays the database migration status
func (cli *Client) CreateMigration(c *clipkg.Context) error {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("You must specify a migration name"))
	}
	db, err := newConnection(cli.Config)
	if err != nil {
		return fmt.Errorf("failed to initialize orm: %v", err)
	}

	migrationType := c.String("type")
	if migrationType != "go" {
		migrationType = "sql"
	}

	if err = migrate.Create(db.DB, c.Args().First(), migrationType); err != nil {
		return fmt.Errorf("Status failed: %v", err)
	}
	return nil
}

type dbConfig interface {
	pg.ConnectionConfig
	DatabaseURL() url.URL
	GetDatabaseDialectConfiguredOrDefault() dialects.DialectName
}

func newConnection(cfg dbConfig) (*sqlx.DB, error) {
	parsed := cfg.DatabaseURL()
	if parsed.String() == "" {
		return nil, errors.New("You must set DATABASE_URL env variable. HINT: If you are running this to set up your local test database, try DATABASE_URL=postgresql://postgres@localhost:5432/chainlink_test?sslmode=disable")
	}
	return pg.NewConnection(parsed.String(), cfg.GetDatabaseDialectConfiguredOrDefault(), cfg)
}

func dropAndCreateDB(parsed url.URL) (err error) {
	// Cannot drop the database if we are connected to it, so we must connect
	// to a different one. template1 should be present on all postgres installations
	dbname := parsed.Path[1:]
	parsed.Path = "/template1"
	db, err := sql.Open(string(dialects.Postgres), parsed.String())
	if err != nil {
		return fmt.Errorf("unable to open postgres database for creating test db: %+v", err)
	}
	defer func() {
		if cerr := db.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	_, err = db.Exec(fmt.Sprintf(`DROP DATABASE IF EXISTS "%s"`, dbname))
	if err != nil {
		return fmt.Errorf("unable to drop postgres database: %v", err)
	}
	_, err = db.Exec(fmt.Sprintf(`CREATE DATABASE "%s"`, dbname))
	if err != nil {
		return fmt.Errorf("unable to create postgres database: %v", err)
	}
	return nil
}

func dropAndCreatePristineDB(db *sqlx.DB, template string) (err error) {
	_, err = db.Exec(fmt.Sprintf(`DROP DATABASE IF EXISTS "%s"`, PristineDBName))
	if err != nil {
		return fmt.Errorf("unable to drop postgres database: %v", err)
	}
	_, err = db.Exec(fmt.Sprintf(`CREATE DATABASE "%s" WITH TEMPLATE "%s"`, PristineDBName, template))
	if err != nil {
		return fmt.Errorf("unable to create postgres database: %v", err)
	}
	return nil
}

func migrateDB(config dbConfig, lggr logger.Logger) error {
	db, err := newConnection(config)
	if err != nil {
		return fmt.Errorf("failed to initialize orm: %v", err)
	}
	if err = migrate.Migrate(db.DB, lggr); err != nil {
		return fmt.Errorf("migrateDB failed: %v", err)
	}
	return db.Close()
}

func downAndUpDB(cfg dbConfig, lggr logger.Logger, baseVersionID int64) error {
	db, err := newConnection(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize orm: %v", err)
	}
	if err = migrate.Rollback(db.DB, lggr, null.IntFrom(baseVersionID)); err != nil {
		return fmt.Errorf("test rollback failed: %v", err)
	}
	if err = migrate.Migrate(db.DB, lggr); err != nil {
		return fmt.Errorf("second migrateDB failed: %v", err)
	}
	return db.Close()
}

func dumpSchema(dbURL url.URL) (string, error) {
	args := []string{
		dbURL.String(),
		"--schema-only",
	}
	cmd := exec.Command(
		"pg_dump", args...,
	)

	schema, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to dump schema: %v", err)
	}
	return string(schema), nil
}

func checkSchema(dbURL url.URL, prevSchema string) error {
	newSchema, err := dumpSchema(dbURL)
	if err != nil {
		return err
	}
	df := diff.Diff(prevSchema, newSchema)
	if len(df) > 0 {
		fmt.Println(df)
		return errors.New("schema pre- and post- rollback does not match (ctrl+f for '+' or '-' to find the changed lines)")
	}
	return nil
}

func insertFixtures(dbURL url.URL, pathToFixtures string) (err error) {
	db, err := sql.Open(string(dialects.Postgres), dbURL.String())
	if err != nil {
		return fmt.Errorf("unable to open postgres database for creating test db: %+v", err)
	}
	defer func() {
		if cerr := db.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		return errors.New("could not get runtime.Caller(1)")
	}
	filepath := path.Join(path.Dir(filename), pathToFixtures)
	fixturesSQL, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}
	_, err = db.Exec(string(fixturesSQL))
	return err
}
