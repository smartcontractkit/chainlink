package cmd

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/smartcontractkit/chainlink/logger"
	strpkg "github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/presenters"
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

	app := cli.AppFactory.NewApplication(config)
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
	defer app.Stop()
	logNodeBalance(store)
	logConfigVariables(store)
	logIfNonceOutOfSync(store)

	return cli.errorOut(cli.Runner.Run(app))
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
	lastNonce, err := store.GetLastNonce(account.Address)
	if err != nil {
		logger.Error("database error when checking nonce: ", err)
		return
	}

	if localNonceIsNotCurrent(lastNonce, account.GetNonce()) {
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
		config.LogLevel = strpkg.LogLevel{Level: zapcore.DebugLevel}
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
	wlc := presenters.NewConfigWhitelist(store)
	logger.Debug("Environment variables\n", wlc)
}

// DeleteUser is run locally to remove the User row from the node's database.
func (cli *Client) DeleteUser(c *clipkg.Context) error {
	logger.SetLogger(cli.Config.CreateProductionLogger())
	app := cli.AppFactory.NewApplication(cli.Config)
	store := app.GetStore()
	user, err := store.DeleteUser()
	if err == nil {
		logger.Info("Deleted API user ", user.Email)
	}
	return err
}

// ImportKey imports a key to be used with the chainlink node
func (cli *Client) ImportKey(c *clipkg.Context) error {
	cfg := cli.Config
	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass in filepath to key"))
	}

	src := c.Args().First()
	kdir := cfg.KeysDir()

	if e, err := isDirEmpty(kdir); !e && err != nil {
		return cli.errorOut(err)
	}

	if i := strings.LastIndex(src, "/"); i < 0 {
		kdir += "/" + src
	} else {
		kdir += src[strings.LastIndex(src, "/"):]
	}
	return cli.errorOut(copyFile(src, kdir))
}

func isDirEmpty(dir string) (bool, error) {
	f, err := os.Open(dir)
	if err != nil {
		return false, err
	}
	defer f.Close()

	if _, err = f.Readdirnames(1); err == io.EOF {
		return true, nil
	}

	return false, fmt.Errorf("Account already present in keystore: %s", dir)
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

// CreateExtraKey creates a new ethereum key with the same password
// as the one used to unlock the existing key.
func (cli *Client) CreateExtraKey(c *clipkg.Context) error {
	logger.SetLogger(cli.Config.CreateProductionLogger())
	app := cli.AppFactory.NewApplication(cli.Config)

	pwd, err := passwordFromFile(c.String("password"))
	if err != nil {
		return cli.errorOut(fmt.Errorf("error reading password: %+v", err))
	}

	store := app.GetStore()
	pwd, err = cli.KeyStoreAuthenticator.Authenticate(store, pwd)
	if err != nil {
		return cli.errorOut(fmt.Errorf("error authenticating keystore: %+v", err))
	}

	account, err := store.KeyStore.NewAccount(pwd)
	if err != nil {
		return err
	}

	logger.Infow(fmt.Sprintf("Created account %v", account.Address.Hex()), "address", account.Address.Hex())
	return nil
}
