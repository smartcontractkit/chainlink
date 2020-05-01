package cmd

import (
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

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/gracefulpanic"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	strpkg "github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/migrations"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/store/presenters"
	"github.com/smartcontractkit/chainlink/core/utils"

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
		return err
	}

	updateConfig(cli.Config, c.Bool("debug"), c.Int64("replay-from-block"))
	logger.SetLogger(cli.Config.CreateProductionLogger())
	logger.Infow("Starting Chainlink Node " + strpkg.Version + " at commit " + strpkg.Sha)

	err = InitEnclave()
	if err != nil {
		return cli.errorOut(fmt.Errorf("error initializing SGX enclave: %+v", err))
	}

	app := cli.AppFactory.NewApplication(cli.Config, func(app chainlink.Application) {
		store := app.GetStore()
		logNodeBalance(store)
		logIfNonceOutOfSync(store)
	})
	store := app.GetStore()
	if e := checkFilePermissions(cli.Config.RootDir()); e != nil {
		logger.Warn(e)
	}
	pwd, err := passwordFromFile(c.String("password"))
	if err != nil {
		return cli.errorOut(fmt.Errorf("error reading password: %+v", err))
	}
	_, err = cli.KeyStoreAuthenticator.Authenticate(store, pwd)
	if err != nil {
		return cli.errorOut(fmt.Errorf("error authenticating keystore: %+v", err))
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

	return cli.errorOut(cli.Runner.Run(app))
}

func loggedStop(app chainlink.Application) {
	logger.WarnIf(app.Stop())
}

func checkFilePermissions(rootDir string) error {
	errorMsg := "%s has overly permissive file permissions, should be atleast %s"
	keysDir := filepath.Join(rootDir, "tempkeys")
	protectedFiles := []string{"secret", "cookie"}
	err := filepath.Walk(keysDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			fileMode := info.Mode().Perm()
			if fileMode&^ownerPermsMask != 0 {
				return fmt.Errorf(errorMsg, path, ownerPermsMask)
			}
			return nil
		})
	if err != nil {
		return err
	}
	for _, fileName := range protectedFiles {
		fileInfo, err := os.Lstat(filepath.Join(rootDir, fileName))
		if err != nil {
			return err
		}
		perm := fileInfo.Mode().Perm()
		if perm&^ownerPermsMask != 0 {
			return fmt.Errorf(errorMsg, fileName, ownerPermsMask)
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

func logIfNonceOutOfSync(store *strpkg.Store) {
	account := store.TxManager.NextActiveAccount()
	if account == nil {
		return
	}
	lastNonce, err := store.GetLastNonce(account.Address)
	if err != nil {
		logger.Error("database error when checking nonce: ", err)
		return
	}

	if localNonceIsNotCurrent(lastNonce, account.Nonce()) {
		logger.Warn("The account is being used by another wallet and is not safe to use with chainlink")
	}
}

func localNonceIsNotCurrent(lastNonce, nonce uint64) bool {
	return lastNonce+1 < nonce
}

func updateConfig(config *orm.Config, debug bool, replayFromBlock int64) {
	if debug {
		config.Set("LOG_LEVEL", zapcore.DebugLevel.String())
	}
	if replayFromBlock >= 0 {
		config.Set(orm.EnvVarName("ReplayFromBlock"), replayFromBlock)
	}
}

func logNodeBalance(store *strpkg.Store) {
	accounts, err := presenters.ShowEthBalance(store)
	logger.WarnIf(err)
	for _, a := range accounts {
		logger.Infow(a["message"], "address", a["address"], "ethBalance", a["balance"])
	}

	accounts, err = presenters.ShowLinkBalance(store)
	logger.WarnIf(err)
	for _, a := range accounts {
		logger.Infow(a["message"], "address", a["address"], "linkBalance", a["balance"])
	}
}

func logConfigVariables(store *strpkg.Store) error {
	wlc, err := presenters.NewConfigWhitelist(store)
	if err != nil {
		return err
	}

	logger.Debug("Environment variables\n", wlc)
	return nil
}

// RebroadcastTransactions run locally to force manual rebroadcasting of
// transactions in a given nonce range. This MUST NOT be run concurrently with
// the node. Currently the advisory lock in FindAllTxsInNonceRange prevents
// this.
func (cli *Client) RebroadcastTransactions(c *clipkg.Context) error {
	beginningNonce := c.Uint("beginningNonce")
	endingNonce := c.Uint("endingNonce")
	gasPriceWei := c.Uint64("gasPriceWei")
	overrideGasLimit := c.Uint64("gasLimit")

	logger.SetLogger(cli.Config.CreateProductionLogger())
	app := cli.AppFactory.NewApplication(cli.Config)
	defer app.Stop()

	store := app.GetStore()

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
		return err
	}

	lastHead, err := store.LastHead()
	if err != nil {
		return err
	}
	err = store.TxManager.Connect(lastHead)
	if err != nil {
		return err
	}

	transactions, err := store.FindAllTxsInNonceRange(beginningNonce, endingNonce)
	if err != nil {
		return err
	}
	n := len(transactions)
	for i, tx := range transactions {
		var gasLimit uint64
		if overrideGasLimit == 0 {
			gasLimit = tx.GasLimit
		} else {
			gasLimit = overrideGasLimit
		}
		logger.Infow("Rebroadcasting transaction", "idx", i, "of", n, "nonce", tx.Nonce, "id", tx.ID)

		gasPrice := big.NewInt(int64(gasPriceWei))
		rawTx, err := store.TxManager.SignedRawTxWithBumpedGas(tx, gasLimit, *gasPrice)
		if err != nil {
			logger.Error(err)
			continue
		}

		hash, err := store.TxManager.SendRawTx(rawTx)
		if err != nil {
			logger.Error(err)
			continue
		}

		logger.Infow("Sent transaction", "idx", i, "of", n, "nonce", tx.Nonce, "id", tx.ID, "hash", hash)

		jobRunID, err := models.NewIDFromString(tx.SurrogateID.ValueOrZero())
		if err != nil {
			logger.Errorw("could not get UUID from surrogate ID", "SurrogateID", tx.SurrogateID.ValueOrZero())
			continue
		}
		jobRun, err := store.FindJobRun(jobRunID)
		if err != nil {
			logger.Errorw("could not find job run", "id", jobRunID)
			continue
		}
		for taskIndex := range jobRun.TaskRuns {
			taskRun := &jobRun.TaskRuns[taskIndex]
			if taskRun.Status == models.RunStatusPendingOutgoingConfirmations {
				taskRun.Status = models.RunStatusErrored
			}
		}
		jobRun.SetStatus(models.RunStatusErrored)

		err = store.ORM.SaveJobRun(&jobRun)
		if err != nil {
			logger.Errorw("error saving job run", "id", jobRunID)
			continue
		}
	}
	return nil
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

func dropAndCreateDB(parsed url.URL) error {
	// Cannot drop the database if we are connected to it, so we must connect
	// to a different one. template1 should be present on all postgres installations
	dbname := parsed.Path[1:]
	parsed.Path = "/template1"
	db, err := sql.Open(string(orm.DialectPostgres), parsed.String())
	if err != nil {
		return fmt.Errorf("unable to open postgres database for creating test db: %+v", err)
	}
	defer db.Close()

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

func insertFixtures(config *orm.Config) error {
	db, err := sql.Open(string(orm.DialectPostgres), config.DatabaseURL())
	if err != nil {
		return fmt.Errorf("unable to open postgres database for creating test db: %+v", err)
	}
	defer db.Close()

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
func (cli *Client) DeleteUser(c *clipkg.Context) error {
	logger.SetLogger(cli.Config.CreateProductionLogger())
	app := cli.AppFactory.NewApplication(cli.Config)
	defer app.Stop()
	store := app.GetStore()
	user, err := store.DeleteUser()
	if err == nil {
		logger.Info("Deleted API user ", user.Email)
	}
	return err
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

	err := utils.EnsureDirAndPerms(dstDirPath, 0700|os.ModeDir)
	if err != nil {
		return cli.errorOut(err)
	}

	err = utils.CopyFileWithPerms(srcKeyPath, dstKeyPath, 0600)
	if err != nil {
		return cli.errorOut(err)
	}

	return app.GetStore().SyncDiskKeyStoreToDB()
}
