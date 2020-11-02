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

	"go.uber.org/multierr"

	"github.com/pkg/errors"
	"github.com/tidwall/gjson"

	"github.com/smartcontractkit/chainlink/core/gracefulpanic"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/static"
	strpkg "github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/migrations"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/store/presenters"
	"github.com/smartcontractkit/chainlink/core/utils"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/jinzhu/gorm"
	clipkg "github.com/urfave/cli"
	"go.uber.org/zap/zapcore"
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
	logger.Infow("Starting Chainlink Node " + static.Version + " at commit " + static.Sha)

	app := cli.AppFactory.NewApplication(cli.Config, func(app chainlink.Application) {
		store := app.GetStore()
		checkAccountsForExternalUse(store)
	})
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

	if authErr := cli.KeyStoreAuthenticator.AuthenticateOCRKey(store, keyStorePwd); authErr != nil {
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
	if _, err = NewFileAPIInitializer(c.String("api")).Initialize(store); err != nil && err != errNoCredentialFile {
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
		fundingKey, currentBalance, err := setupFundingKey(context.TODO(), app.GetStore(), keyStorePwd)
		if err != nil {
			return cli.errorOut(errors.Wrap(err, "failed to generate a funding address"))
		}
		if currentBalance.Cmp(big.NewInt(0)) == 0 {
			logger.Infow("The backup funding address does not have sufficient funds", "address", fundingKey.Address, "balance", currentBalance)
		} else {
			logger.Infow("Funding address ready", "address", fundingKey.Address, "current-balance", currentBalance)
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

func checkAccountsForExternalUse(store *strpkg.Store) {
	keys, err := store.AllKeys()
	if err != nil {
		logger.Error("database error while retrieving send keys:", err)
		return
	}
	for _, key := range keys {
		logIfNonceOutOfSync(store, key)
	}
}

func logIfNonceOutOfSync(store *strpkg.Store, key models.Key) {
	onChainNonce, err := store.EthClient.PendingNonceAt(context.TODO(), key.Address.Address())
	if err != nil {
		logger.Error(fmt.Sprintf("error determining nonce for address %s: %v", key.Address.Hex(), err))
		return
	}
	var nonce int64
	if key.NextNonce != nil {
		nonce = *key.NextNonce
	}
	if nonce < int64(onChainNonce) {
		logger.Warn(fmt.Sprintf("The account %s is being used by another wallet and is not safe to use with chainlink", key.Address.Hex()))
	}
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

func setupFundingKey(ctx context.Context, str *strpkg.Store, pwd string) (*models.Key, *big.Int, error) {
	key := models.Key{}
	err := str.DB.Where("is_funding = TRUE").First(&key).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, nil, err
	}
	if err == nil && key.ID != 0 {
		// TODO How to make sure the EthClient is connected?
		balance, ethErr := str.EthClient.BalanceAt(ctx, gethCommon.HexToAddress(string(key.Address)), nil)
		return &key, balance, ethErr
	}
	// Key record not found so create one.
	ethAccount, err := str.KeyStore.NewAccount()
	if err != nil {
		return nil, nil, err
	}
	exportedJSON, err := str.KeyStore.Export(ethAccount.Address, pwd)
	if err != nil {
		return nil, nil, err
	}
	var firstNonce int64 = 0
	key = models.Key{
		Address:   models.EIP55Address(ethAccount.Address.Hex()),
		IsFunding: true,
		JSON: models.JSON{
			Result: gjson.ParseBytes(exportedJSON),
		},
		NextNonce: &firstNonce,
	}
	// The key does not exist at this point, so we're only creating it here.
	if err = str.CreateKeyIfNotExists(key); err != nil {
		return nil, nil, err
	}
	logger.Infow("New funding address created", "address", key.Address, "balance", 0)
	return &key, big.NewInt(0), nil
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
	cli.Config.Dialect = orm.DialectPostgresWithoutLock
	app := cli.AppFactory.NewApplication(cli.Config)
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

	ec := bulletprooftxmanager.NewEthConfirmer(store, cli.Config)
	err = ec.ForceRebroadcast(beginningNonce, endingNonce, gasPriceWei, address, overrideGasLimit)
	return cli.errorOut(err)
}

// HardReset will remove all non-started transactions if any are found.
func (cli *Client) HardReset(c *clipkg.Context) error {
	logger.SetLogger(cli.Config.CreateProductionLogger())

	if !confirmAction(c) {
		return nil
	}

	app, cleanupFn := cli.makeApp()
	defer cleanupFn()
	storeInstance := app.GetStore()
	ormInstance := storeInstance.ORM

	// Ensure that the CL node is down by trying to acquire the global advisory lock.
	// This method will panic if it can't get the lock.
	logger.Info("Make sure the Chainlink node is not running")
	ormInstance.MustEnsureAdvisoryLock()

	if err := ormInstance.RemoveUnstartedTransactions(); err != nil {
		logger.Errorw("failed to remove unstarted transactions", "error", err)
		return err
	}

	logger.Info("successfully reset the node state in the database")
	return nil
}

func (cli *Client) makeApp() (chainlink.Application, func()) {
	app := cli.AppFactory.NewApplication(cli.Config)
	return app, func() {
		if err := app.Stop(); err != nil {
			logger.Errorw("Failed to stop the application on hard reset", "error", err)
		}
	}
}

// ResetDatabase drops, creates and migrates the database specified by DATABASE_URL
// This is useful to setup the database for testing
func (cli *Client) ResetDatabase(c *clipkg.Context) error {
	logger.SetLogger(cli.Config.CreateProductionLogger())
	config := orm.NewConfig()
	if config.DatabaseURL() == "" {
		return cli.errorOut(errors.New("You must set DATABASE_URL env variable. HINT: If you are running this to set up your local test database, try DATABASE_URL=postgresql://postgres@localhost:5432/chainlink_test?sslmode=disable"))
	}
	parsed, err := url.Parse(config.DatabaseURL())
	if err != nil {
		return cli.errorOut(err)
	}

	dbname := parsed.Path[1:]
	if !strings.HasSuffix(dbname, "_test") {
		return cli.errorOut(fmt.Errorf("cannot reset database named `%s`. This command can only be run against databases with a name that ends in `_test`, to prevent accidental data loss", dbname))
	}
	logger.Infof("Resetting database: %#v", config.DatabaseURL())
	if err := dropAndCreateDB(*parsed); err != nil {
		return cli.errorOut(err)
	}
	if err := migrateTestDB(config); err != nil {
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

func dropAndCreateDB(parsed url.URL) (err error) {
	// Cannot drop the database if we are connected to it, so we must connect
	// to a different one. template1 should be present on all postgres installations
	dbname := parsed.Path[1:]
	parsed.Path = "/template1"
	db, err := sql.Open(string(orm.DialectPostgres), parsed.String())
	if err != nil {
		return fmt.Errorf("unable to open postgres database for creating test db: %+v", err)
	}
	defer func() {
		if cerr := db.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	_, err = db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbname))
	if err != nil {
		return fmt.Errorf("unable to drop postgres database: %v", err)
	}
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbname))
	if err != nil {
		return fmt.Errorf("unable to create postgres database: %v", err)
	}
	return nil
}

func migrateTestDB(config *orm.Config) error {
	orm, err := orm.NewORM(config.DatabaseURL(), config.DatabaseTimeout(), gracefulpanic.NewSignal(), config.GetDatabaseDialectConfiguredOrDefault(), config.GetAdvisoryLockIDConfiguredOrDefault())
	if err != nil {
		return fmt.Errorf("failed to initialize orm: %v", err)
	}
	orm.SetLogging(config.LogSQLStatements() || config.LogSQLMigrations())
	err = orm.RawDB(func(db *gorm.DB) error {
		return migrations.GORMMigrate(db)
	})
	if err != nil {
		return fmt.Errorf("migrateTestDB failed: %v", err)
	}
	orm.SetLogging(config.LogSQLStatements())
	return orm.Close()
}

func insertFixtures(config *orm.Config) (err error) {
	db, err := sql.Open(string(orm.DialectPostgres), config.DatabaseURL())
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
	filepath := path.Join(path.Dir(filename), "../store/testdata/fixtures.sql")
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
	app := cli.AppFactory.NewApplication(cli.Config)
	defer func() {
		if serr := app.Stop(); serr != nil {
			err = multierr.Append(err, serr)
		}
	}()
	store := app.GetStore()
	user, err := store.DeleteUser()
	if err == nil {
		logger.Info("Deleted API user ", user.Email)
	}
	return err
}

// SetNextNonce manually updates the keys.next_nonce field for the given key with the given nonce value
func (cli *Client) SetNextNonce(c *clipkg.Context) error {
	addressHex := c.String("address")
	nextNonce := c.Uint64("nextNonce")

	logger.SetLogger(cli.Config.CreateProductionLogger())
	db, err := gorm.Open(string(orm.DialectPostgres), cli.Config.DatabaseURL())
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
func (cli *Client) ImportKey(c *clipkg.Context) error {
	logger.SetLogger(cli.Config.CreateProductionLogger())
	app := cli.AppFactory.NewApplication(cli.Config)

	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass in filepath to key"))
	}

	var (
		srcKeyPath = c.Args().First()                      // ex: ./keys/mykey
		srcKeyFile = filepath.Base(srcKeyPath)             // ex: mykey
		dstDirPath = cli.Config.KeysDir()                  // ex: /clroot/keys
		dstKeyPath = filepath.Join(dstDirPath, srcKeyFile) // ex: /clroot/keys/mykey
	)

	err := utils.EnsureDirAndMaxPerms(dstDirPath, 0700|os.ModeDir)
	if err != nil {
		return cli.errorOut(err)
	}

	err = utils.CopyFileWithMaxPerms(srcKeyPath, dstKeyPath, 0600)
	if err != nil {
		return cli.errorOut(err)
	}

	return app.GetStore().SyncDiskKeyStoreToDB()
}
