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
	"strings"

	"github.com/fatih/color"
	"github.com/smartcontractkit/chainlink/core/store/dialects"
	"github.com/smartcontractkit/chainlink/core/store/migrations"

	gormpostgres "gorm.io/driver/postgres"

	"go.uber.org/multierr"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/gracefulpanic"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/health"
	"github.com/smartcontractkit/chainlink/core/static"
	strpkg "github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/store/presenters"
	"github.com/smartcontractkit/chainlink/core/utils"
	webPresenters "github.com/smartcontractkit/chainlink/core/web/presenters"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	clipkg "github.com/urfave/cli"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
)

// ownerPermsMask are the file permission bits reserved for owner.
const ownerPermsMask = os.FileMode(0700)

// RunNode starts the Chainlink core.
func (cli *Client) RunNode(c *clipkg.Context) error {
	err := cli.Config.Validate()
	if err != nil {
		return cli.errorOut(err)
	}

	updateConfig(cli.Config, c.Bool("debug"), c.Int64("replay-from-block"))
	logger.SetLogger(cli.Config.CreateProductionLogger())
	logger.Infow(fmt.Sprintf("Starting Chainlink Node %s at commit %s", static.Version, static.Sha), "id", "boot", "Version", static.Version, "SHA", static.Sha, "InstanceUUID", static.InstanceUUID)
	if cli.Config.Dev() {
		logger.Warn("Chainlink is running in DEVELOPMENT mode. This is a security risk if enabled in production.")
	}

	app, err := cli.AppFactory.NewApplication(cli.Config)
	if err != nil {
		return cli.errorOut(errors.Wrap(err, "creating application"))
	}
	store := app.GetStore()
	if e := checkFilePermissions(cli.Config.RootDir()); e != nil {
		logger.Warn(e)
	}
	pwd, err := passwordFromFile(c.String("password"))
	if err != nil {
		return cli.errorOut(fmt.Errorf("error reading password: %+v", err))
	}
	keyStorePwd, err := cli.KeyStoreAuthenticator.Authenticate(store, pwd)
	if err != nil {
		return cli.errorOut(fmt.Errorf("error authenticating keystore: %+v", err))
	}

	if authErr := cli.KeyStoreAuthenticator.AuthenticateOCRKey(app, keyStorePwd); authErr != nil {
		return cli.errorOut(errors.Wrapf(authErr, "while authenticating with OCR password"))
	}

	if len(c.String("vrfpassword")) != 0 {
		vrfpwd, fileErr := passwordFromFile(c.String("vrfpassword"))
		if fileErr != nil {
			return cli.errorOut(errors.Wrapf(fileErr,
				"error reading VRF password from vrfpassword file \"%s\"",
				c.String("vrfpassword")))
		}
		if authErr := cli.KeyStoreAuthenticator.AuthenticateVRFKey(store, vrfpwd); authErr != nil {
			return cli.errorOut(errors.Wrapf(authErr, "while authenticating with VRF password"))
		}
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
	err = logConfigVariables(store)
	if err != nil {
		return err
	}

	if !store.Config.EthereumDisabled() {
		key, currentBalance, err := setupFundingKey(context.TODO(), app.GetStore(), keyStorePwd)
		if err != nil {
			return cli.errorOut(errors.Wrap(err, "failed to generate a funding address"))
		}
		if currentBalance.Cmp(big.NewInt(0)) == 0 {
			logger.Infow("The backup funding address does not have sufficient funds", "address", key.Address.Hex(), "balance", currentBalance)
		} else {
			logger.Infow("Funding address ready", "address", key.Address.Hex(), "current-balance", currentBalance)
		}
	}

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

func updateConfig(config *orm.Config, debug bool, replayFromBlock int64) {
	if debug {
		config.Set("LOG_LEVEL", zapcore.DebugLevel.String())
	}
	if replayFromBlock >= 0 {
		config.Set(orm.EnvVarName("ReplayFromBlock"), replayFromBlock)
	}
}

func logConfigVariables(store *strpkg.Store) error {
	wlc, err := presenters.NewConfigPrinter(store)
	if err != nil {
		return err
	}

	logger.Debug("Environment variables\n", wlc)
	return nil
}

func setupFundingKey(ctx context.Context, str *strpkg.Store, pwd string) (key models.Key, balance *big.Int, err error) {
	key, existed, err := str.KeyStore.EnsureFundingKey()
	if err != nil {
		return key, nil, err
	}
	if existed {
		// TODO How to make sure the EthClient is connected?
		balance, ethErr := str.EthClient.BalanceAt(ctx, key.Address.Address(), nil)
		return key, balance, ethErr
	}
	logger.Infow("New funding address created", "address", key.Address.Hex(), "balance", 0)
	return key, big.NewInt(0), nil
}

// RebroadcastTransactions run locally to force manual rebroadcasting of
// transactions in a given nonce range.
func (cli *Client) RebroadcastTransactions(c *clipkg.Context) (err error) {
	beginningNonce := c.Uint("beginningNonce")
	endingNonce := c.Uint("endingNonce")
	gasPriceWei := c.Uint64("gasPriceWei")
	overrideGasLimit := c.Uint64("gasLimit")
	addressHex := c.String("address")

	addressBytes, err := hexutil.Decode(addressHex)
	if err != nil {
		return cli.errorOut(errors.Wrap(err, "could not decode address"))
	}
	address := gethCommon.BytesToAddress(addressBytes)

	logger.SetLogger(cli.Config.CreateProductionLogger())
	cli.Config.Dialect = dialects.PostgresWithoutLock
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

	err = store.EthClient.Dial(context.TODO())
	if err != nil {
		return err
	}

	pwd, err := passwordFromFile(c.String("password"))
	if err != nil {
		return cli.errorOut(fmt.Errorf("error reading password: %+v", err))
	}
	_, err = cli.KeyStoreAuthenticator.Authenticate(store, pwd)
	if err != nil {
		return cli.errorOut(fmt.Errorf("error authenticating keystore: %+v", err))
	}

	err = store.Start()
	if err != nil {
		return cli.errorOut(err)
	}

	logger.Infof("Rebroadcasting transactions from %v to %v", beginningNonce, endingNonce)

	allKeys, err := store.KeyStore.AllKeys()
	if err != nil {
		return cli.errorOut(err)
	}
	ec := bulletprooftxmanager.NewEthConfirmer(store.DB, store.EthClient, cli.Config, store.KeyStore, store.AdvisoryLocker, allKeys)
	err = ec.ForceRebroadcast(beginningNonce, endingNonce, gasPriceWei, address, overrideGasLimit)
	return cli.errorOut(err)
}

// HardReset will remove all non-started transactions if any are found.
func (cli *Client) HardReset(c *clipkg.Context) error {
	logger.SetLogger(cli.Config.CreateProductionLogger())

	fmt.Print("/// WARNING WARNING WARNING ///\n\n\n")
	fmt.Print("Do not run this while a Chainlink node is currently using the DB as it could cause undefined behavior.\n\n")
	if !confirmAction(c) {
		return nil
	}

	app, cleanupFn, err := cli.makeApp()
	if err != nil {
		logger.Errorw("error while creating application", "error", err)
		return err
	}
	defer cleanupFn()
	storeInstance := app.GetStore()
	ormInstance := storeInstance.ORM

	if err := ormInstance.RemoveUnstartedTransactions(); err != nil {
		logger.Errorw("failed to remove unstarted transactions", "error", err)
		return err
	}

	logger.Info("successfully reset the node state in the database")
	return nil
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

func (cli *Client) makeApp() (chainlink.Application, func(), error) {
	app, err := cli.AppFactory.NewApplication(cli.Config)
	if err != nil {
		return nil, nil, err
	}
	return app, func() {
		if err := app.Stop(); err != nil {
			logger.Errorw("Failed to stop the application on hard reset", "error", err)
		}
	}, nil
}

// ResetDatabase drops, creates and migrates the database specified by DATABASE_URL
// This is useful to setup the database for testing
func (cli *Client) ResetDatabase(c *clipkg.Context) error {
	logger.SetLogger(cli.Config.CreateProductionLogger())
	config := orm.NewConfig()
	parsed := config.DatabaseURL()
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
	if err := migrateDB(config); err != nil {
		return cli.errorOut(err)
	}
	return nil
}

// PrepareTestDatabase calls ResetDatabase then loads fixtures required for tests
func (cli *Client) PrepareTestDatabase(c *clipkg.Context) error {
	if err := cli.ResetDatabase(c); err != nil {
		return cli.errorOut(err)
	}
	config := orm.NewConfig()
	if err := insertFixtures(config); err != nil {
		return cli.errorOut(err)
	}
	return nil
}

// MigrateDatabase migrates the database
func (cli *Client) MigrateDatabase(c *clipkg.Context) error {
	logger.SetLogger(cli.Config.CreateProductionLogger())
	config := orm.NewConfig()
	parsed := config.DatabaseURL()
	if parsed.String() == "" {
		return cli.errorOut(errors.New("You must set DATABASE_URL env variable. HINT: If you are running this to set up your local test database, try DATABASE_URL=postgresql://postgres@localhost:5432/chainlink_test?sslmode=disable"))
	}

	logger.Infof("Migrating database: %#v", parsed.String())
	if err := migrateDB(config); err != nil {
		return cli.errorOut(err)
	}
	return nil
}

// VersionDatabase displays the current database version.
func (cli *Client) VersionDatabase(c *clipkg.Context) error {
	logger.SetLogger(cli.Config.CreateProductionLogger())
	config := orm.NewConfig()
	parsed := config.DatabaseURL()
	if parsed.String() == "" {
		return cli.errorOut(errors.New("You must set DATABASE_URL env variable. HINT: If you are running this to set up your local test database, try DATABASE_URL=postgresql://postgres@localhost:5432/chainlink_test?sslmode=disable"))
	}

	orm, err := orm.NewORM(parsed.String(), config.DatabaseTimeout(), gracefulpanic.NewSignal(), config.GetDatabaseDialectConfiguredOrDefault(), config.GetAdvisoryLockIDConfiguredOrDefault(), config.GlobalLockRetryInterval().Duration(), config.ORMMaxOpenConns(), config.ORMMaxIdleConns())
	if err != nil {
		return fmt.Errorf("failed to initialize orm: %v", err)
	}

	version, err := migrations.Current(orm.DB)
	if err != nil {
		return fmt.Errorf("migrateDB failed: %v", err)
	}

	logger.Infof("Database version: %v", version.ID)
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

func migrateDB(config *orm.Config) error {
	dbURL := config.DatabaseURL()
	orm, err := orm.NewORM(dbURL.String(), config.DatabaseTimeout(), gracefulpanic.NewSignal(), config.GetDatabaseDialectConfiguredOrDefault(), config.GetAdvisoryLockIDConfiguredOrDefault(), config.GlobalLockRetryInterval().Duration(), config.ORMMaxOpenConns(), config.ORMMaxIdleConns())
	if err != nil {
		return fmt.Errorf("failed to initialize orm: %v", err)
	}
	orm.SetLogging(config.LogSQLStatements() || config.LogSQLMigrations())

	from, err := migrations.Current(orm.DB)
	if err != nil {
		from = &migrations.Migration{
			ID: "(none)",
		}
	}

	to := migrations.Migrations[len(migrations.Migrations)-1]

	logger.Infof("Migrating from %v to %v", from.ID, to.ID)

	err = migrations.Migrate(orm.DB)
	if err != nil {
		return fmt.Errorf("migrateDB failed: %v", err)
	}
	orm.SetLogging(config.LogSQLStatements())
	return orm.Close()
}

func insertFixtures(config *orm.Config) (err error) {
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

	res := db.Exec(`UPDATE keys SET next_nonce = ? WHERE address = ?`, nextNonce, address)
	if res.Error != nil {
		return cli.errorOut(err)
	}
	if res.RowsAffected == 0 {
		return cli.errorOut(fmt.Errorf("no key found matching address %s", addressHex))
	}
	return nil
}

// ImportKey imports a key to be used with the chainlink node
// NOTE: This should not be run concurrently with a running chainlink node.
// If you do run it concurrently, it will not take effect until the next reboot.
func (cli *Client) ImportKey(c *clipkg.Context) error {
	logger.SetLogger(cli.Config.CreateProductionLogger())
	app, err := cli.AppFactory.NewApplication(cli.Config)
	if err != nil {
		return cli.errorOut(errors.Wrap(err, "creating application"))
	}

	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass in filepath to key"))
	}

	srcKeyPath := c.Args().First() // e.g. ./keys/mykey

	_, err = app.GetStore().KeyStore.ImportKeyFileToDB(srcKeyPath)
	return cli.errorOut(err)
}
