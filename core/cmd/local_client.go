package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/fatih/color"
	"github.com/pkg/errors"
	clipkg "github.com/urfave/cli"
	"go.uber.org/multierr"
	null "gopkg.in/guregu/null.v4"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/smartcontractkit/chainlink/core/gracefulpanic"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/health"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/static"
	"github.com/smartcontractkit/chainlink/core/store/config"
	"github.com/smartcontractkit/chainlink/core/store/dialects"
	"github.com/smartcontractkit/chainlink/core/store/migrate"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/store/presenters"
	"github.com/smartcontractkit/chainlink/core/utils"
	webPresenters "github.com/smartcontractkit/chainlink/core/web/presenters"
)

// ownerPermsMask are the file permission bits reserved for owner.
const ownerPermsMask = os.FileMode(0700)

// RunNode starts the Chainlink core.
func (cli *Client) RunNode(c *clipkg.Context) error {
	logger.SetLogger(cli.Config.CreateProductionLogger())

	err := cli.Config.Validate()
	if err != nil {
		return cli.errorOut(err)
	}

	logger.Infow(fmt.Sprintf("Starting Chainlink Node %s at commit %s", static.Version, static.Sha), "id", "boot", "Version", static.Version, "SHA", static.Sha, "InstanceUUID", static.InstanceUUID)
	if cli.Config.Dev() {
		logger.Warn("Chainlink is running in DEVELOPMENT mode. This is a security risk if enabled in production.")
	}
	if cli.Config.EthereumDisabled() {
		logger.Warn("Ethereum is disabled. Chainlink will only run services that can operate without an ethereum connection")
	}

	app, err := cli.AppFactory.NewApplication(cli.Config)
	if err != nil {
		return cli.errorOut(errors.Wrap(err, "creating application"))
	}

	store := app.GetStore()
	keyStore := app.GetKeyStore()
	err = cli.KeyStoreAuthenticator.authenticate(c, keyStore)
	if err != nil {
		return cli.errorOut(errors.Wrap(err, "error authenticating keystore"))
	}

	var vrfpwd string
	var fileErr error
	if len(c.String("vrfpassword")) != 0 {
		vrfpwd, fileErr = passwordFromFile(c.String("vrfpassword"))
		if fileErr != nil {
			return cli.errorOut(errors.Wrapf(fileErr,
				"error reading VRF password from vrfpassword file \"%s\"",
				c.String("vrfpassword")))
		}
	}

	err = keyStore.Migrate(vrfpwd)
	if err != nil {
		return cli.errorOut(errors.Wrap(err, "error migrating keystore"))
	}

	if e := checkFilePermissions(cli.Config.RootDir()); e != nil {
		logger.Warn(e)
	}

	var user models.User
	if _, err = NewFileAPIInitializer(c.String("api")).Initialize(store); err != nil && err != ErrNoCredentialFile {
		return cli.errorOut(fmt.Errorf("error creating api initializer: %+v", err))
	}
	if user, err = cli.FallbackAPIInitializer.Initialize(store); err != nil {
		if err == ErrorNoAPICredentialsAvailable {
			return cli.errorOut(err)
		}
		return cli.errorOut(fmt.Errorf("error creating fallback initializer: %+v", err))
	}

	logger.Info("API exposed for user ", user.Email)
	if e := app.Start(); e != nil {
		return cli.errorOut(fmt.Errorf("error starting app: %+v", e))
	}
	defer loggedStop(app)
	err = logConfigVariables(cli.Config)
	if err != nil {
		return cli.errorOut(err)
	}
	for _, ch := range app.GetChainSet().Chains() {
		skey, sexisted, fkey, fexisted, err := app.GetKeyStore().Eth().EnsureKeys(ch.ID())
		if err != nil {
			return cli.errorOut(err)
		}
		if !fexisted {
			logger.Infow("New funding address created", "address", fkey.Address.Hex(), "evmChainID", ch.ID())
		}
		if !sexisted {
			logger.Infow("New sending address created", "address", skey.Address.Hex(), "evmChainID", ch.ID())
		}
	}

	logger.Infof("Chainlink booted in %s", time.Since(static.InitTime))
	return cli.errorOut(cli.Runner.Run(app))
}

func loggedStop(app chainlink.Application) {
	logger.WarnIf(app.Stop())
}

func checkFilePermissions(rootDir string) error {
	// Ensure `$CLROOT/tls` directory (and children) permissions are <= `ownerPermsMask``
	tlsDir := filepath.Join(rootDir, "tls")
	_, err := os.Stat(tlsDir)
	if err != nil && !os.IsNotExist(err) {
		logger.Errorf("error checking perms of 'tls' directory: %v", err)
	} else if err == nil {
		err := utils.EnsureDirAndMaxPerms(tlsDir, ownerPermsMask)
		if err != nil {
			return err
		}

		err = filepath.Walk(tlsDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				logger.Errorf(`error checking perms of "%v": %v`, path, err)
				return err
			}
			if utils.TooPermissive(info.Mode().Perm(), ownerPermsMask) {
				newPerms := info.Mode().Perm() & ownerPermsMask
				logger.Warnf("%s has overly permissive file permissions, reducing them from %s to %s", path, info.Mode().Perm(), newPerms)
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
			logger.Warnf("%s has overly permissive file permissions, reducing them from %s to %s", path, fileInfo.Mode().Perm(), newPerms)
			err = utils.EnsureFilepathMaxPerms(path, newPerms)
			if err != nil {
				return err
			}
		}
		owned, err := utils.IsFileOwnedByChainlink(fileInfo)
		if err != nil {
			logger.Warn(err)
			continue
		}
		if !owned {
			logger.Warnf("The file %v is not owned by the user running chainlink. This will be made mandatory in the future.", path)
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

func logConfigVariables(cfg config.GeneralConfig) error {
	wlc, err := presenters.NewConfigPrinter(cfg)
	if err != nil {
		return err
	}

	logger.Debug("Environment variables\n", wlc)
	return nil
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

	logger.SetLogger(cli.Config.CreateProductionLogger())
	cli.Config.SetDialect(dialects.PostgresWithoutLock)
	app, err := cli.AppFactory.NewApplication(cli.Config)
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
	chain, err := app.GetChainSet().Get(chainID)
	if err != nil {
		return cli.errorOut(err)
	}
	store := app.GetStore()
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

	err = store.Start()
	if err != nil {
		return cli.errorOut(err)
	}

	logger.Infof("Rebroadcasting transactions from %v to %v", beginningNonce, endingNonce)

	keyStates, err := keyStore.Eth().GetStatesForChain(chain.ID())
	if err != nil {
		return cli.errorOut(err)
	}
	ec := bulletprooftxmanager.NewEthConfirmer(store.DB, ethClient, chain.Config(), keyStore.Eth(), store.AdvisoryLocker, keyStates, nil, nil, chain.Logger())
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
	case health.StatusFailing:
		status = red(p.Status)
	case health.StatusPassing:
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
	logger.SetLogger(cli.Config.CreateProductionLogger())
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
	logger.Infof("Resetting database: %#v", parsed.String())
	if err := dropAndCreateDB(parsed); err != nil {
		return cli.errorOut(err)
	}
	if err := migrateDB(cfg); err != nil {
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
	if err := insertFixtures(cfg); err != nil {
		return cli.errorOut(err)
	}
	return nil
}

// MigrateDatabase migrates the database
func (cli *Client) MigrateDatabase(c *clipkg.Context) error {
	logger.SetLogger(cli.Config.CreateProductionLogger())
	cfg := cli.Config
	parsed := cfg.DatabaseURL()
	if parsed.String() == "" {
		return cli.errorOut(errors.New("You must set DATABASE_URL env variable. HINT: If you are running this to set up your local test database, try DATABASE_URL=postgresql://postgres@localhost:5432/chainlink_test?sslmode=disable"))
	}

	logger.Infof("Migrating database: %#v", parsed.String())
	if err := migrateDB(cfg); err != nil {
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

	logger.SetLogger(cli.Config.CreateProductionLogger())
	cfg := cli.Config
	parsed := cfg.DatabaseURL()
	if parsed.String() == "" {
		return cli.errorOut(errors.New("You must set DATABASE_URL env variable. HINT: If you are running this to set up your local test database, try DATABASE_URL=postgresql://postgres@localhost:5432/chainlink_test?sslmode=disable"))
	}

	orm, err := orm.NewORM(parsed.String(), cfg.DatabaseTimeout(), gracefulpanic.NewSignal(), cfg.GetDatabaseDialectConfiguredOrDefault(), cfg.GetAdvisoryLockIDConfiguredOrDefault(), cfg.GlobalLockRetryInterval().Duration(), cfg.ORMMaxOpenConns(), cfg.ORMMaxIdleConns())
	if err != nil {
		return fmt.Errorf("failed to initialize orm: %v", err)
	}

	db := postgres.UnwrapGormDB(orm.DB).DB

	if err := migrate.Rollback(db, version); err != nil {
		return fmt.Errorf("migrateDB failed: %v", err)
	}

	return nil
}

// VersionDatabase displays the current database version.
func (cli *Client) VersionDatabase(c *clipkg.Context) error {
	logger.SetLogger(cli.Config.CreateProductionLogger())
	cfg := cli.Config
	parsed := cfg.DatabaseURL()
	if parsed.String() == "" {
		return cli.errorOut(errors.New("You must set DATABASE_URL env variable. HINT: If you are running this to set up your local test database, try DATABASE_URL=postgresql://postgres@localhost:5432/chainlink_test?sslmode=disable"))
	}

	orm, err := orm.NewORM(parsed.String(), cfg.DatabaseTimeout(), gracefulpanic.NewSignal(), cfg.GetDatabaseDialectConfiguredOrDefault(), cfg.GetAdvisoryLockIDConfiguredOrDefault(), cfg.GlobalLockRetryInterval().Duration(), cfg.ORMMaxOpenConns(), cfg.ORMMaxIdleConns())
	if err != nil {
		return fmt.Errorf("failed to initialize orm: %v", err)
	}

	db := postgres.UnwrapGormDB(orm.DB).DB

	version, err := migrate.Current(db)
	if err != nil {
		return fmt.Errorf("migrateDB failed: %v", err)
	}

	logger.Infof("Database version: %v", version)
	return nil
}

// StatusDatabase displays the database migration status
func (cli *Client) StatusDatabase(c *clipkg.Context) error {
	logger.SetLogger(cli.Config.CreateProductionLogger())
	cfg := cli.Config
	parsed := cfg.DatabaseURL()
	if parsed.String() == "" {
		return cli.errorOut(errors.New("You must set DATABASE_URL env variable. HINT: If you are running this to set up your local test database, try DATABASE_URL=postgresql://postgres@localhost:5432/chainlink_test?sslmode=disable"))
	}

	orm, err := orm.NewORM(parsed.String(), cfg.DatabaseTimeout(), gracefulpanic.NewSignal(), cfg.GetDatabaseDialectConfiguredOrDefault(), cfg.GetAdvisoryLockIDConfiguredOrDefault(), cfg.GlobalLockRetryInterval().Duration(), cfg.ORMMaxOpenConns(), cfg.ORMMaxIdleConns())
	if err != nil {
		return fmt.Errorf("failed to initialize orm: %v", err)
	}

	db := postgres.UnwrapGormDB(orm.DB).DB

	if err = migrate.Status(db); err != nil {
		return fmt.Errorf("Status failed: %v", err)
	}
	return nil
}

// CreateMigration displays the database migration status
func (cli *Client) CreateMigration(c *clipkg.Context) error {
	logger.SetLogger(cli.Config.CreateProductionLogger())
	cfg := cli.Config
	parsed := cfg.DatabaseURL()
	if parsed.String() == "" {
		return cli.errorOut(errors.New("You must set DATABASE_URL env variable. HINT: If you are running this to set up your local test database, try DATABASE_URL=postgresql://postgres@localhost:5432/chainlink_test?sslmode=disable"))
	}

	if !c.Args().Present() {
		return cli.errorOut(errors.New("You must specify a migration name"))
	}

	orm, err := orm.NewORM(parsed.String(), cfg.DatabaseTimeout(), gracefulpanic.NewSignal(), cfg.GetDatabaseDialectConfiguredOrDefault(), cfg.GetAdvisoryLockIDConfiguredOrDefault(), cfg.GlobalLockRetryInterval().Duration(), cfg.ORMMaxOpenConns(), cfg.ORMMaxIdleConns())
	if err != nil {
		return fmt.Errorf("failed to initialize orm: %v", err)
	}

	db := postgres.UnwrapGormDB(orm.DB).DB

	migrationType := c.String("type")
	if migrationType != "go" {
		migrationType = "sql"
	}

	if err = migrate.Create(db, c.Args().First(), migrationType); err != nil {
		return fmt.Errorf("Status failed: %v", err)
	}
	return nil
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

func migrateDB(config config.GeneralConfig) error {
	dbURL := config.DatabaseURL()
	orm, err := orm.NewORM(dbURL.String(), config.DatabaseTimeout(), gracefulpanic.NewSignal(), config.GetDatabaseDialectConfiguredOrDefault(), config.GetAdvisoryLockIDConfiguredOrDefault(), config.GlobalLockRetryInterval().Duration(), config.ORMMaxOpenConns(), config.ORMMaxIdleConns())
	if err != nil {
		return fmt.Errorf("failed to initialize orm: %v", err)
	}
	orm.SetLogging(config.LogSQLStatements() || config.LogSQLMigrations())

	db := postgres.UnwrapGormDB(orm.DB).DB

	if err = migrate.Migrate(db); err != nil {
		return fmt.Errorf("migrateDB failed: %v", err)
	}
	orm.SetLogging(config.LogSQLStatements())
	return orm.Close()
}

func insertFixtures(config config.GeneralConfig) (err error) {
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
	filepath := path.Join(path.Dir(filename), "../store/fixtures/fixtures.sql")
	fixturesSQL, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	_, err = db.Exec(string(fixturesSQL))
	return err
}

// DeleteUser is run locally to remove the User row from the node's database.
func (cli *Client) DeleteUser(c *clipkg.Context) (err error) {
	logger.SetLogger(cli.Config.CreateProductionLogger())
	app, err := cli.AppFactory.NewApplication(cli.Config)
	if err != nil {
		return cli.errorOut(errors.Wrap(err, "creating application"))
	}
	defer func() {
		if serr := app.Stop(); serr != nil {
			err = multierr.Append(err, serr)
		}
	}()
	store := app.GetStore()
	user, err := store.FindUser()
	if err == nil {
		logger.Info("No such API user ", user.Email)
		return err
	}
	err = store.DeleteUser()
	if err == nil {
		logger.Info("Deleted API user ", user.Email)
	}
	return err
}

// SetNextNonce manually updates the keys.next_nonce field for the given key with the given nonce value
func (cli *Client) SetNextNonce(c *clipkg.Context) error {
	addressHex := c.String("address")
	nextNonce := c.Uint64("nextNonce")
	dbURL := cli.Config.DatabaseURL()

	logger.SetLogger(cli.Config.CreateProductionLogger())
	db, err := gorm.Open(gormpostgres.New(gormpostgres.Config{
		DSN: dbURL.String(),
	}), &gorm.Config{})
	if err != nil {
		return cli.errorOut(err)
	}

	address, err := hexutil.Decode(addressHex)
	if err != nil {
		return cli.errorOut(errors.Wrap(err, "could not decode address"))
	}

	res := db.Exec(`UPDATE eth_key_states SET next_nonce = ? WHERE address = ?`, nextNonce, address)
	if res.Error != nil {
		return cli.errorOut(err)
	}
	if res.RowsAffected == 0 {
		return cli.errorOut(fmt.Errorf("no key found matching address %s", addressHex))
	}
	return nil
}
