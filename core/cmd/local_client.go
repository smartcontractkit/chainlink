package cmd

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services"
	strpkg "github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/presenters"
	"github.com/smartcontractkit/chainlink/core/utils"
	clipkg "github.com/urfave/cli"
	"go.uber.org/zap/zapcore"
)

// RunNode starts the Chainlink core.
func (cli *Client) RunNode(c *clipkg.Context) error {
	config := updateConfig(cli.Config, c.Bool("debug"))
	logger.SetLogger(config.CreateProductionLogger())
	logger.Infow("Starting Chainlink Node " + strpkg.Version + " at commit " + strpkg.Sha)

	err := InitEnclave()
	if err != nil {
		return cli.errorOut(fmt.Errorf("error initializing SGX enclave: %+v", err))
	}

	app := cli.AppFactory.NewApplication(config, func(app services.Application) {
		store := app.GetStore()
		logNodeBalance(store)
		logIfNonceOutOfSync(store)
	})
	store := app.GetStore()
	pwd, err := passwordFromFile(c.String("password"))
	if err != nil {
		return cli.errorOut(fmt.Errorf("error reading password: %+v", err))
	}
	_, err = cli.KeyStoreAuthenticator.Authenticate(store, pwd)
	if err != nil {
		return cli.errorOut(fmt.Errorf("error authenticating keystore: %+v", err))
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
	logConfigVariables(store)

	return cli.errorOut(cli.Runner.Run(app))
}

func loggedStop(app services.Application) {
	logger.WarnIf(app.Stop())
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

func updateConfig(config strpkg.Config, debug bool) strpkg.Config {
	if debug {
		config.Set("LOG_LEVEL", zapcore.DebugLevel.String())
	}
	return config
}

func logNodeBalance(store *strpkg.Store) {
	accounts, err := presenters.ShowEthBalance(store)
	logger.WarnIf(err)
	for _, a := range accounts {
		logAccountBalance(a)
	}

	accounts, err = presenters.ShowLinkBalance(store)
	logger.WarnIf(err)
	for _, a := range accounts {
		logAccountBalance(a)
	}
}

func logAccountBalance(kv map[string]interface{}) {
	logger.Infow(fmt.Sprint(kv["message"]), "address", kv["address"], "balance", kv["balance"])
}

func logConfigVariables(store *strpkg.Store) {
	wlc, err := presenters.NewConfigWhitelist(store)
	if err != nil {
		logger.Error("Failed to build environment variables\n", err)
	} else {
		logger.Debug("Environment variables\n", wlc)
	}
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
	defer app.Stop()

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
