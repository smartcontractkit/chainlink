package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/manyminds/api2go/jsonapi"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/smartcontractkit/chainlink/logger"
	"github.com/smartcontractkit/chainlink/services"
	strpkg "github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/presenters"
	"github.com/smartcontractkit/chainlink/utils"
	"github.com/smartcontractkit/chainlink/web"
	"github.com/tidwall/gjson"
	clipkg "github.com/urfave/cli"
	"go.uber.org/multierr"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"
)

// Client is the shell for the node. It has fields for the Renderer,
// Config, AppFactory (the services application), Authenticator, and Runner.
type Client struct {
	Renderer
	Config          strpkg.Config
	AppFactory      AppFactory
	Auth            Authenticator
	UserInitializer UserInitializer
	Runner          Runner
}

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
		return cli.errorOut(fmt.Errorf("error starting app: %+v", err))
	}
	err = cli.Auth.Authenticate(store, pwd)
	if err != nil {
		return cli.errorOut(fmt.Errorf("error starting app: %+v", err))
	}
	var user models.User
	if user, err = cli.UserInitializer.Initialize(store); err != nil {
		return cli.errorOut(fmt.Errorf("error starting app: %+v", err))
	}
	logger.Info("API exposed for user ", user.Email)
	if err := app.Start(); err != nil {
		return cli.errorOut(fmt.Errorf("error starting app: %+v", err))
	}
	defer app.Stop()
	logNodeBalance(store)
	logConfigVariables(config)
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
	account := store.TxManager.GetActiveAccount()
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
	balance, err := presenters.ShowEthBalance(store)
	logger.WarnIf(err)
	logger.Infow(balance)
	balance, err = presenters.ShowLinkBalance(store)
	logger.WarnIf(err)
	logger.Infow(balance)
}

func logConfigVariables(config strpkg.Config) {
	wlc := presenters.NewConfigWhitelist(config)
	logger.Debug("Environment variables\n", wlc)
}

// DisplayAccountBalance renders a table containing the active account address
// with it's ETH & LINK balance
func (cli *Client) DisplayAccountBalance(c *clipkg.Context) error {
	cfg := cli.Config
	resp, err := utils.BasicAuthGet(
		cfg.BasicAuthUsername,
		cfg.BasicAuthPassword,
		cfg.ClientNodeURL+"/v2/account_balance",
	)
	if err != nil {
		return cli.errorOut(err)
	}
	defer resp.Body.Close()

	var links jsonapi.Links
	a := presenters.AccountBalance{}
	if err = cli.deserializeAPIResponse(resp, &a, &links); err != nil {
		return err
	}
	return cli.errorOut(cli.Render(&a))
}

// ShowJobSpec returns the status of the given JobID.
func (cli *Client) ShowJobSpec(c *clipkg.Context) error {
	cfg := cli.Config
	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass the job id to be shown"))
	}
	resp, err := utils.BasicAuthGet(
		cfg.BasicAuthUsername,
		cfg.BasicAuthPassword,
		cfg.ClientNodeURL+"/v2/specs/"+c.Args().First(),
	)
	if err != nil {
		return cli.errorOut(err)
	}
	defer resp.Body.Close()
	var job presenters.JobSpec
	return cli.renderAPIResponse(resp, &job)
}

// GetJobSpecs returns all job specs.
func (cli *Client) GetJobSpecs(c *clipkg.Context) error {
	cfg := cli.Config
	requestURI := cfg.ClientNodeURL + "/v2/specs"

	page := 0
	if c != nil && c.IsSet("page") {
		page = c.Int("page")
	}

	var links jsonapi.Links
	var jobs []models.JobSpec
	err := cli.getPage(requestURI, page, &jobs, &links)
	if err != nil {
		return err
	}
	return cli.errorOut(cli.Render(&jobs))
}

// CreateJobSpec creates job spec based on JSON input
func (cli *Client) CreateJobSpec(c *clipkg.Context) error {
	cfg := cli.Config
	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass in JSON or filepath"))
	}

	buf, err := getBufferFromJSON(c.Args().First())
	if err != nil {
		return cli.errorOut(err)
	}

	resp, err := utils.BasicAuthPost(
		cfg.BasicAuthUsername,
		cfg.BasicAuthPassword,
		cfg.ClientNodeURL+"/v2/specs",
		"application/json",
		buf,
	)
	if err != nil {
		return cli.errorOut(err)
	}
	defer resp.Body.Close()

	var jobs presenters.JobSpec
	return cli.renderResponse(resp, &jobs)
}

// CreateJobRun creates job run based on SpecID and optional JSON
func (cli *Client) CreateJobRun(c *clipkg.Context) error {
	cfg := cli.Config
	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass in SpecID [JSON blob | JSON filepath]"))
	}

	buf := bytes.NewBufferString("")
	if c.NArg() > 1 {
		jbuf, err := getBufferFromJSON(c.Args().Get(1))
		if err != nil {
			return cli.errorOut(err)
		}
		buf = jbuf
	}

	resp, err := utils.BasicAuthPost(
		cfg.BasicAuthUsername,
		cfg.BasicAuthPassword,
		cfg.ClientNodeURL+"/v2/specs/"+c.Args().First()+"/runs",
		"application/json",
		buf,
	)
	if err != nil {
		return cli.errorOut(err)
	}
	defer resp.Body.Close()
	var jobs presenters.JobSpec
	return cli.renderResponse(resp, &jobs)
}

// BackupDatabase streams a backup of the node's db to the passed filepath.
func (cli *Client) BackupDatabase(c *clipkg.Context) error {
	cfg := cli.Config
	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass the path to save the backup"))
	}
	resp, err := utils.BasicAuthGet(
		cfg.BasicAuthUsername,
		cfg.BasicAuthPassword,
		cfg.ClientNodeURL+"/v2/backup",
	)
	if err != nil {
		return cli.errorOut(err)
	}
	defer resp.Body.Close()
	return cli.errorOut(saveBodyAsFile(resp, c.Args().First()))
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

// AddBridge adds a new bridge to the chainlink node
func (cli *Client) AddBridge(c *clipkg.Context) error {
	cfg := cli.Config
	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass in the bridge's parameters [JSON blob | JSON filepath]"))
	}

	buf, err := getBufferFromJSON(c.Args().First())
	if err != nil {
		return cli.errorOut(err)
	}

	resp, err := utils.BasicAuthPost(
		cfg.BasicAuthUsername,
		cfg.BasicAuthPassword,
		cfg.ClientNodeURL+"/v2/bridge_types",
		"application/json",
		buf,
	)
	if err != nil {
		return cli.errorOut(err)
	}
	defer resp.Body.Close()

	var bridge models.BridgeType
	return cli.deserializeResponse(resp, &bridge)
}

// GetBridges returns all bridges.
func (cli *Client) GetBridges(c *clipkg.Context) error {
	cfg := cli.Config
	requestURI := cfg.ClientNodeURL + "/v2/bridge_types"

	page := 0
	if c != nil && c.IsSet("page") {
		page = c.Int("page")
	}

	var links jsonapi.Links
	var bridges []models.BridgeType
	err := cli.getPage(requestURI, page, &bridges, &links)
	if err != nil {
		return err
	}
	return cli.errorOut(cli.Render(&bridges))
}

func (cli *Client) getPage(requestURI string, page int, model interface{}, links *jsonapi.Links) error {
	cfg := cli.Config

	uri, err := url.Parse(requestURI)
	if err != nil {
		return err
	}
	q := uri.Query()
	if page > 0 {
		q.Set("page", strconv.Itoa(page))
	}
	uri.RawQuery = q.Encode()

	resp, err := utils.BasicAuthGet(cfg.BasicAuthUsername, cfg.BasicAuthPassword, uri.String())
	if err != nil {
		return cli.errorOut(err)
	}
	defer resp.Body.Close()

	return cli.deserializeAPIResponse(resp, model, links)
}

// ShowBridge returns the info for the given Bridge name.
func (cli *Client) ShowBridge(c *clipkg.Context) error {
	cfg := cli.Config
	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass the name of the bridge to be shown"))
	}
	resp, err := utils.BasicAuthGet(
		cfg.BasicAuthUsername,
		cfg.BasicAuthPassword,
		cfg.ClientNodeURL+"/v2/bridge_types/"+c.Args().First(),
	)
	if err != nil {
		return cli.errorOut(err)
	}
	defer resp.Body.Close()
	var bridge models.BridgeType
	return cli.renderResponse(resp, &bridge)
}

// RemoveBridge removes a specific Bridge by name.
func (cli *Client) RemoveBridge(c *clipkg.Context) error {
	cfg := cli.Config
	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass the name of the bridge to be removed"))
	}
	resp, err := utils.BasicAuthDelete(
		cfg.BasicAuthUsername,
		cfg.BasicAuthPassword,
		cfg.ClientNodeURL+"/v2/bridge_types/"+c.Args().First(),
		"application/json",
		nil,
	)
	if err != nil {
		return cli.errorOut(err)
	}
	defer resp.Body.Close()
	var bridge models.BridgeType
	return cli.renderResponse(resp, &bridge)
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

func saveBodyAsFile(resp *http.Response, dst string) error {
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}

func getBufferFromJSON(s string) (buf *bytes.Buffer, err error) {
	if gjson.Valid(s) {
		buf, err = bytes.NewBufferString(s), nil
	} else if buf, err = fromFile(s); err != nil {
		buf, err = nil, multierr.Append(errors.New("Must pass in JSON or filepath"), err)
	}
	return
}

func fromFile(arg string) (*bytes.Buffer, error) {
	dir, err := homedir.Expand(arg)
	if err != nil {
		return nil, err
	}
	file, err := ioutil.ReadFile(dir)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(file), nil
}

// deserializeAPIResponse is distinct from deserializeResponse in that it supports JSONAPI responses with Links
func (cli *Client) deserializeAPIResponse(resp *http.Response, dst interface{}, links *jsonapi.Links) error {
	b, err := parseResponse(resp)
	if err != nil {
		return cli.errorOut(err)
	}
	if err = web.ParsePaginatedResponse(b, dst, links); err != nil {
		return cli.errorOut(err)
	}
	return nil
}

func (cli *Client) deserializeResponse(resp *http.Response, dst interface{}) error {
	b, err := parseResponse(resp)
	if err != nil {
		return cli.errorOut(err)
	}
	if err = json.Unmarshal(b, &dst); err != nil {
		return cli.errorOut(err)
	}
	return nil
}

func parseResponse(resp *http.Response) ([]byte, error) {
	b, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		if err != nil {
			return b, errors.New(resp.Status)
		}
		return b, errors.New(string(b))
	}
	return b, err
}

func (cli *Client) renderResponse(resp *http.Response, dst interface{}) error {
	err := cli.deserializeResponse(resp, dst)
	if err != nil {
		return cli.errorOut(err)
	}
	return cli.errorOut(cli.Render(dst))
}

func (cli *Client) renderAPIResponse(resp *http.Response, dst interface{}) error {
	var links jsonapi.Links
	if err := cli.deserializeAPIResponse(resp, dst, &links); err != nil {
		return cli.errorOut(err)
	}
	return cli.errorOut(cli.Render(dst))
}

func (cli *Client) errorOut(err error) error {
	if err != nil {
		return clipkg.NewExitError(err.Error(), 1)
	}
	return nil
}

// AppFactory implements the NewApplication method.
type AppFactory interface {
	NewApplication(strpkg.Config) services.Application
}

// ChainlinkAppFactory is used to create a new Application.
type ChainlinkAppFactory struct{}

// NewApplication returns a new instance of the node with the given config.
func (n ChainlinkAppFactory) NewApplication(config strpkg.Config) services.Application {
	return services.NewApplication(config)
}

// Runner implements the Run method.
type Runner interface {
	Run(services.Application) error
}

// ChainlinkRunner is used to run the node application.
type ChainlinkRunner struct{}

// Run sets the log level based on config and starts the web router to listen
// for input and return data.
func (n ChainlinkRunner) Run(app services.Application) error {
	gin.SetMode(app.GetStore().Config.LogLevel.ForGin())
	server := web.Router(app.(*services.ChainlinkApplication))
	config := app.GetStore().Config
	var g errgroup.Group

	if config.Dev {
		g.Go(func() error { return server.Run(":" + config.Port) })
	} else {
		certFile := config.CertFile()
		keyFile := config.KeyFile()
		g.Go(func() error { return server.RunTLS(":"+config.Port, certFile, keyFile) })
	}
	return g.Wait()
}
