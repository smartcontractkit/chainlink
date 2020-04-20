package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	strpkg "github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/store/presenters"
	"github.com/smartcontractkit/chainlink/core/utils"

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
	if err := checkFilePermissions(cli.Config.RootDir()); err != nil {
		logger.Warn(err)
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
		vrfpwd, err := passwordFromFile(c.String("vrfpassword"))
		if err != nil {
			return cli.errorOut(errors.Wrapf(err,
				"error reading VRF password from vrfpassword file \"%s\"",
				c.String("vrfpassword")))
		}
		if err := cli.KeyStoreAuthenticator.AuthenticateVRFKey(store, vrfpwd); err != nil {
			return cli.errorOut(errors.Wrapf(err, "while authenticating with VRF password"))
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
	if err := app.Start(); err != nil {
		return cli.errorOut(fmt.Errorf("error starting app: %+v", err))
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
	if lastNonce+1 < nonce {
		return true
	}

	return false
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
			if taskRun.Status == models.RunStatusPendingConfirmations {
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

	src := c.Args().First()
	kdir := cli.Config.KeysDir()

	if !utils.FileExists(kdir) {
		err := os.MkdirAll(kdir, os.FileMode(0700))
		if err != nil {
			return cli.errorOut(err)
		}
	}

	if i := strings.LastIndex(src, "/"); i < 0 {
		kdir += "/" + src
	} else {
		kdir += src[strings.LastIndex(src, "/"):]
	}

	if err := copyFile(src, kdir); err != nil {
		return cli.errorOut(err)
	}

	return app.GetStore().SyncDiskKeyStoreToDB()
}

func copyFile(src, dst string) error {
	from, err := os.Open(src)
	if err != nil {
		return err
	}
	defer from.Close()

	to, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer to.Close()

	_, err = io.Copy(to, from)

	return err
}
