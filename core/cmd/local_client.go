package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/fatih/color"
	"github.com/kylelemons/godebug/diff"
	"github.com/pkg/errors"
	clipkg "github.com/urfave/cli"
	"go.uber.org/multierr"
	"golang.org/x/sync/errgroup"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/chains/evm/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/config"
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
	"github.com/smartcontractkit/sqlx"
)

// ownerPermsMask are the file permission bits reserved for owner.
const ownerPermsMask = os.FileMode(0700)

func openDB(cfg config.GeneralConfig, lggr logger.Logger) (db *sqlx.DB, err error) {
	uri := cfg.DatabaseURL()
	appid := cfg.AppID()
	static.SetConsumerName(&uri, "App", &appid)
	dialect := cfg.GetDatabaseDialectConfiguredOrDefault()
	db, err = pg.NewConnection(uri.String(), string(dialect), pg.Config{
		Logger:       lggr,
		MaxOpenConns: cfg.ORMMaxOpenConns(),
		MaxIdleConns: cfg.ORMMaxIdleConns(),
	})
	err = errors.Wrap(err, "failed to open db")
	return
}

func applicationLockDB(cfg config.GeneralConfig, db *sqlx.DB, lggr logger.Logger) (release func(), err error) {
	lggr.Debugf("Using database locking mode: %s", cfg.DatabaseLockingMode())

	lockingMode := cfg.DatabaseLockingMode()
	release = func() {}

	// Lease will be explicitly released on application stop
	// Take the lease before any other DB operations
	var leaseLock pg.LeaseLock
	switch lockingMode {
	case "lease", "dual":
		leaseLock = pg.NewLeaseLock(db, cfg.AppID(), lggr, cfg.LeaseLockRefreshInterval(), cfg.LeaseLockDuration())
		if err = leaseLock.TakeAndHold(); err != nil {
			return nil, errors.Wrap(err, "failed to take initial lease on database")
		}
		release = leaseLock.Release
	}

	// Try to acquire an advisory lock to prevent multiple nodes starting at the same time
	var advisoryLock pg.AdvisoryLock
	switch lockingMode {
	case "advisorylock", "dual":
		advisoryLock = pg.NewAdvisoryLock(db, cfg.AdvisoryLockID(), lggr, cfg.AdvisoryLockCheckInterval())
		if err = advisoryLock.TakeAndHold(); err != nil {
			return nil, errors.Wrap(err, "error acquiring lock")
		}
		release = advisoryLock.Release
	}

	return release, nil
}

// RunNode starts the Chainlink core.
func (cli *Client) RunNode(c *clipkg.Context) error {
	if err := cli.runNode(c); err != nil {
		err = errors.Wrap(err, "Cannot boot Chainlink")
		cli.Logger.Error(err)
		if serr := cli.Logger.Sync(); serr != nil {
			err = multierr.Combine(serr, err)
		}
		return cli.errorOut(err)
	}
	return nil
}

func (cli *Client) runNode(c *clipkg.Context) error {
	lggr := cli.Logger.Named("RunNode")

	err := cli.Config.Validate()
	if err != nil {
		return errors.Wrap(err, "config validation failed")
	}

	lggr.Infow(fmt.Sprintf("Starting Chainlink Node %s at commit %s", static.Version, static.Sha), "Version", static.Version, "SHA", static.Sha)

	if cli.Config.Dev() {
		lggr.Warn("Chainlink is running in DEVELOPMENT mode. This is a security risk if enabled in production.")
	}
	if cli.Config.EthereumDisabled() {
		lggr.Warn("Ethereum is disabled. Chainlink will only run services that can operate without an ethereum connection")
	}

	// DB conn is opened before any services have started,
	// and closed after all services have stopped.
	db, err := openDB(cli.Config, lggr)
	if err != nil {
		return cli.errorOut(errors.Wrap(err, "opening db"))
	}
	defer lggr.ErrorIfClosing(db, "db")

	// rootCtx will be cancelled when SIGINT|SIGTERM is received.
	rootCtx, cancelRootCtx := context.WithCancel(context.Background())
	go shutdown.CancelOnShutdown(cancelRootCtx)

	releaseDbLock, err := applicationLockDB(cli.Config, db, lggr)
	if err != nil {
		return cli.errorOut(errors.Wrap(err, "obtaining application db lock"))
	}

	app, err := cli.AppFactory.NewApplication(cli.Config, db)
	if err != nil {
		return errors.Wrap(err, "error initializing application")
	}

	sessionORM := app.SessionORM()
	keyStore := app.GetKeyStore()
	err = cli.KeyStoreAuthenticator.authenticate(c, keyStore)
	if err != nil {
		return errors.Wrap(err, "error authenticating keystore")
	}

	var vrfpwd string
	var fileErr error
	if len(c.String("vrfpassword")) != 0 {
		vrfpwd, fileErr = passwordFromFile(c.String("vrfpassword"))
		if fileErr != nil {
			return errors.Wrapf(fileErr,
				"error reading VRF password from vrfpassword file \"%s\"",
				c.String("vrfpassword"))
		}
	}

	chainSet := app.GetChains().EVM
	dflt, err := chainSet.Default()
	if err != nil {
		return errors.Wrap(err, "failed to get default chainset")
	}
	err = keyStore.Migrate(vrfpwd, dflt.ID())
	if err != nil {
		return errors.Wrap(err, "error migrating keystore")
	}

	for _, ch := range chainSet.Chains() {
		skey, sexisted, fkey, fexisted, err2 := app.GetKeyStore().Eth().EnsureKeys(ch.ID())
		if err2 != nil {
			return errors.Wrap(err2, "failed to ensure keystore keys")
		}
		if !fexisted {
			lggr.Infow("New funding address created", "address", fkey.Address.Hex(), "evmChainID", ch.ID())
		}
		if !sexisted {
			lggr.Infow("New sending address created", "address", skey.Address.Hex(), "evmChainID", ch.ID())
		}
	}

	ocrKey, didExist, err := app.GetKeyStore().OCR().EnsureKey()
	if err != nil {
		return errors.Wrap(err, "failed to ensure ocr key")
	}
	if !didExist {
		lggr.Infof("Created OCR key with ID %s", ocrKey.ID())
	}
	ocr2Keys, keysDidExist, err := app.GetKeyStore().OCR2().EnsureKeys()
	if err != nil {
		return errors.Wrap(err, "failed to ensure ocr key")
	}
	for chainType, didExist := range keysDidExist {
		if !didExist {
			lggr.Infof("Created OCR2 key with ID %s", ocr2Keys[chainType].ID())
		}
	}
	p2pKey, didExist, err := app.GetKeyStore().P2P().EnsureKey()
	if err != nil {
		return errors.Wrap(err, "failed to ensure p2p key")
	}
	if !didExist {
		lggr.Infof("Created P2P key with ID %s", p2pKey.ID())
	}
	solanaKey, didExist, err := app.GetKeyStore().Solana().EnsureKey()
	if err != nil {
		return errors.Wrap(err, "failed to ensure solana key")
	}
	if !didExist {
		lggr.Infof("Created Solana key with ID %s", solanaKey.ID())
	}
	terraKey, didExist, err := app.GetKeyStore().Terra().EnsureKey()
	if err != nil {
		return errors.Wrap(err, "failed to ensure terra key")
	}
	if !didExist {
		lggr.Infof("Created Terra key with ID %s", terraKey.ID())
	}

	if e := checkFilePermissions(lggr, cli.Config.RootDir()); e != nil {
		lggr.Warn(e)
	}

	var user sessions.User
	if _, err = NewFileAPIInitializer(c.String("api"), lggr).Initialize(sessionORM); err != nil && err != ErrNoCredentialFile {
		return errors.Wrap(err, "error creating api initializer")
	}
	if user, err = cli.FallbackAPIInitializer.Initialize(sessionORM); err != nil {
		if err == ErrorNoAPICredentialsAvailable {
			return errors.WithStack(err)
		}
		return errors.Wrap(err, "error creating fallback initializer")
	}

	lggr.Info("API exposed for user ", user.Email)

	grp, grpCtx := errgroup.WithContext(rootCtx)

	if err = app.Start(); err != nil {
		return errors.Wrap(err, "error starting app")
	}

	grp.Go(func() error {
		<-grpCtx.Done()
		defer releaseDbLock()
		if err = app.Stop(); err != nil {
			return errors.Wrap(err, "error stopping app")
		}
		return nil
	})

	lggr.Debug("Environment variables\n", config.NewConfigPrinter(cli.Config))

	lggr.Infow(fmt.Sprintf("Chainlink booted in %.2fs", time.Since(static.InitTime).Seconds()), "appID", app.ID())

	// After SIGTERM grace period, we must fail fast, but try closing DB
	go func() {
		<-rootCtx.Done()
		<-time.After(cli.Config.ShutdownGracePeriod())

		// this code should not execute if the application finished within the grace period
		log.Printf("Shutdown grace period of %v exceeded, closing DB and exiting...\n", cli.Config.ShutdownGracePeriod())
		if err = db.Close(); err != nil {
			log.Printf("Failed to close DB: %v\n", err)
		}
		if err = lggr.Sync(); err != nil {
			log.Printf("Failed to sync Logger: %v", err)
		}

		os.Exit(-1)
	}()

	grp.Go(func() error {
		err = cli.Runner.Run(grpCtx, app)
		// In tests we have custom runners that stop the app gracefully,
		// therefore we need to cancel rootCtx when the Runner has quit.
		// In a real application, this does noop.
		select {
		case <-rootCtx.Done():
		default:
			cancelRootCtx()
		}
		return err
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

func passwordFromFile(pwdFile string) (string, error) {
	if len(pwdFile) == 0 {
		return "", nil
	}
	dat, err := ioutil.ReadFile(pwdFile)
	return strings.TrimSpace(string(dat)), err
}

// RebroadcastTransactions run locally to force manual rebroadcasting of
// transactions in a given nonce range.
func (cli *Client) RebroadcastTransactions(c *clipkg.Context) (err error) {
	beginningNonce := c.Uint("beginningNonce")
	endingNonce := c.Uint("endingNonce")
	gasPriceWei := c.Uint64("gasPriceWei")
	overrideGasLimit := c.Uint64("gasLimit")
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
			return cli.errorOut(errors.Wrap(err, "invalid evmChainID"))
		}
	}

	lggr := cli.Logger.Named("RebroadcastTransactions")
	db, err := openDB(cli.Config, lggr)
	if err != nil {
		return cli.errorOut(errors.Wrap(err, "opening db"))
	}
	defer lggr.ErrorIfClosing(db, "db")

	app, err := cli.AppFactory.NewApplication(cli.Config, db)
	if err != nil {
		return cli.errorOut(errors.Wrap(err, "creating application"))
	}
	defer func() {
		if serr := app.Stop(); serr != nil {
			err = multierr.Append(err, serr)
		}
	}()
	pwd, err := passwordFromFile(c.String("password"))
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
	ec := bulletprooftxmanager.NewEthConfirmer(app.GetSqlxDB(), ethClient, chain.Config(), keyStore.Eth(), keyStates, nil, nil, chain.Logger())
	err = ec.ForceRebroadcast(beginningNonce, endingNonce, gasPriceWei, address, overrideGasLimit)
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
	schema, err := dumpSchema(cfg)
	if err != nil {
		return cli.errorOut(err)
	}
	lggr.Debugf("Testing rollback and re-migrate for database: %#v", parsed.String())
	var baseVersionID int64 = 54
	if err := downAndUpDB(cfg, lggr, baseVersionID); err != nil {
		return cli.errorOut(err)
	}
	if err := checkSchema(cfg, schema); err != nil {
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
	userOnly := c.Bool("user-only")
	var fixturePath = "../store/fixtures/fixtures.sql"
	if userOnly {
		fixturePath = "../store/fixtures/user_only_fixture.sql"
	}
	if err := insertFixtures(cfg, fixturePath); err != nil {
		return cli.errorOut(err)
	}
	return nil
}

// PrepareTestDatabase calls ResetDatabase then loads fixtures required for local
// testing against testnets. Does not include fake chain fixtures.
func (cli *Client) PrepareTestDatabaseUserOnly(c *clipkg.Context) error {
	if err := cli.ResetDatabase(c); err != nil {
		return cli.errorOut(err)
	}
	cfg := cli.Config
	if err := insertFixtures(cfg, "../store/fixtures/user_only_fixtures.sql"); err != nil {
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

	db, err := newConnection(cli.Config, cli.Logger)
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
	db, err := newConnection(cli.Config, cli.Logger)
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
	db, err := newConnection(cli.Config, cli.Logger)
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
	db, err := newConnection(cli.Config, cli.Logger)
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

func newConnection(cfg config.GeneralConfig, lggr logger.Logger) (*sqlx.DB, error) {
	parsed := cfg.DatabaseURL()
	if parsed.String() == "" {
		return nil, errors.New("You must set DATABASE_URL env variable. HINT: If you are running this to set up your local test database, try DATABASE_URL=postgresql://postgres@localhost:5432/chainlink_test?sslmode=disable")
	}
	config := pg.Config{
		Logger:       lggr,
		MaxOpenConns: cfg.ORMMaxOpenConns(),
		MaxIdleConns: cfg.ORMMaxIdleConns(),
	}
	db, err := pg.NewConnection(parsed.String(), string(cfg.GetDatabaseDialectConfiguredOrDefault()), config)
	return db, err
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

func migrateDB(config config.GeneralConfig, lggr logger.Logger) error {
	db, err := newConnection(config, lggr)
	if err != nil {
		return fmt.Errorf("failed to initialize orm: %v", err)
	}
	if err = migrate.Migrate(db.DB, lggr); err != nil {
		return fmt.Errorf("migrateDB failed: %v", err)
	}
	return db.Close()
}

func downAndUpDB(cfg config.GeneralConfig, lggr logger.Logger, baseVersionID int64) error {
	db, err := newConnection(cfg, lggr)
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

func dumpSchema(cfg config.GeneralConfig) (string, error) {
	dbURL := cfg.DatabaseURL()
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

func checkSchema(cfg config.GeneralConfig, prevSchema string) error {
	newSchema, err := dumpSchema(cfg)
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

func insertFixtures(config config.GeneralConfig, pathToFixtures string) (err error) {
	dbURL := config.DatabaseURL()
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
	fixturesSQL, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	_, err = db.Exec(string(fixturesSQL))
	return err
}

// DeleteUser is run locally to remove the User row from the node's database.
func (cli *Client) DeleteUser(c *clipkg.Context) (err error) {
	lggr := cli.Logger.Named("DeleteUser")

	db, err := openDB(cli.Config, lggr)
	if err != nil {
		return cli.errorOut(errors.Wrap(err, "opening db"))
	}
	defer lggr.ErrorIfClosing(db, "db")

	app, err := cli.AppFactory.NewApplication(cli.Config, db)
	if err != nil {
		return cli.errorOut(errors.Wrap(err, "creating application"))
	}
	defer func() {
		if serr := app.Stop(); serr != nil {
			err = multierr.Append(err, serr)
		}
	}()
	orm := app.SessionORM()
	user, err := orm.FindUser()
	if err == nil {
		app.GetLogger().Info("No such API user ", user.Email)
		return err
	}
	err = orm.DeleteUser()
	if err == nil {
		app.GetLogger().Info("Deleted API user ", user.Email)
	}
	return err
}

// SetNextNonce manually updates the keys.next_nonce field for the given key with the given nonce value
func (cli *Client) SetNextNonce(c *clipkg.Context) error {
	addressHex := c.String("address")
	nextNonce := c.Uint64("nextNonce")

	db, err := newConnection(cli.Config, cli.Logger)
	if err != nil {
		return cli.errorOut(err)
	}

	address, err := hexutil.Decode(addressHex)
	if err != nil {
		return cli.errorOut(errors.Wrap(err, "could not decode address"))
	}

	res, err := db.Exec(`UPDATE eth_key_states SET next_nonce = $1 WHERE address = $2`, nextNonce, address)
	if err != nil {
		return cli.errorOut(err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return cli.errorOut(err)
	}
	if rowsAffected == 0 {
		return cli.errorOut(fmt.Errorf("no key found matching address %s", addressHex))
	}
	return nil
}
