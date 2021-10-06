package cmd

import (
	"bytes"
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/gracefulpanic"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/periodicbackup"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/services/versioning"
	"github.com/smartcontractkit/chainlink/core/services/webhook"
	"github.com/smartcontractkit/chainlink/core/sessions"
	"github.com/smartcontractkit/chainlink/core/static"
	"github.com/smartcontractkit/chainlink/core/store/config"
	"github.com/smartcontractkit/chainlink/core/store/migrate"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/web"

	"github.com/Depado/ginprom"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	clipkg "github.com/urfave/cli"
	"go.uber.org/multierr"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"
)

var prometheus *ginprom.Prometheus

func init() {
	// ensure metrics are regsitered once per instance to avoid registering
	// metrics multiple times (panic)
	prometheus = ginprom.New(ginprom.Namespace("service"))
}

var (
	// ErrorNoAPICredentialsAvailable is returned when not run from a terminal
	// and no API credentials have been provided
	ErrorNoAPICredentialsAvailable = errors.New("API credentials must be supplied")
)

// Client is the shell for the node, local commands and remote commands.
type Client struct {
	Renderer
	Config                         config.GeneralConfig
	AppFactory                     AppFactory
	KeyStoreAuthenticator          TerminalKeyStoreAuthenticator
	FallbackAPIInitializer         APIInitializer
	Runner                         Runner
	HTTP                           HTTPClient
	CookieAuthenticator            CookieAuthenticator
	FileSessionRequestBuilder      SessionRequestBuilder
	PromptingSessionRequestBuilder SessionRequestBuilder
	ChangePasswordPrompter         ChangePasswordPrompter
	PasswordPrompter               PasswordPrompter
}

func (cli *Client) errorOut(err error) error {
	if err != nil {
		return clipkg.NewExitError(err.Error(), 1)
	}
	return nil
}

// AppFactory implements the NewApplication method.
type AppFactory interface {
	NewApplication(config.GeneralConfig) (chainlink.Application, error)
}

// ChainlinkAppFactory is used to create a new Application.
type ChainlinkAppFactory struct{}

// NewApplication returns a new instance of the node with the given config.
func (n ChainlinkAppFactory) NewApplication(cfg config.GeneralConfig) (chainlink.Application, error) {
	globalLogger := cfg.CreateProductionLogger()

	shutdownSignal := gracefulpanic.NewSignal()
	uri := cfg.DatabaseURL()
	dialect := cfg.GetDatabaseDialectConfiguredOrDefault()
	db, gormDB, err := postgres.NewConnection(uri.String(), string(dialect), postgres.Config{
		LogSQLStatements: cfg.LogSQLStatements(),
		MaxOpenConns:     cfg.ORMMaxOpenConns(),
		MaxIdleConns:     cfg.ORMMaxIdleConns(),
	})
	if err != nil {
		return nil, err
	}
	keyStore := keystore.New(gormDB, utils.GetScryptParams(cfg), globalLogger)
	cfg.SetDB(gormDB)

	// Set up the versioning ORM
	verORM := versioning.NewORM(db)

	// Set up periodic backup
	if cfg.DatabaseBackupMode() != config.DatabaseBackupModeNone {
		var version *versioning.NodeVersion
		var versionString string

		version, err = verORM.FindLatestNodeVersion()
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				globalLogger.Debugf("Failed to find any node version in the DB: %w", err)
			} else if strings.Contains(err.Error(), "relation \"node_versions\" does not exist") {
				globalLogger.Debugf("Failed to find any node version in the DB, the node_versions table does not exist yet: %w", err)
			} else {
				return nil, errors.Wrap(err, "initializeORM#FindLatestNodeVersion")
			}
		}

		if version != nil {
			versionString = version.Version
		}

		databaseBackup := periodicbackup.NewDatabaseBackup(cfg, globalLogger)
		databaseBackup.RunBackupGracefully(versionString)
	}

	// Migrate the database
	if cfg.MigrateDatabase() {
		if err = migrate.Migrate(db.DB); err != nil {
			return nil, errors.Wrap(err, "initializeORM#Migrate")
		}
	}

	// Determine node version
	nodeVersion := static.Version
	if nodeVersion == "unset" {
		nodeVersion = fmt.Sprintf("random_%d", rand.Uint32())
	}
	version := versioning.NewNodeVersion(nodeVersion)
	if err = verORM.UpsertNodeVersion(version); err != nil {
		return nil, errors.Wrap(err, "initializeORM#UpsertNodeVersion")
	}

	if cfg.ClobberNodesFromEnv() {
		if err = evm.ClobberNodesFromEnv(gormDB, cfg); err != nil {
			return nil, err
		}
	}

	eventBroadcaster := postgres.NewEventBroadcaster(cfg.DatabaseURL(), cfg.DatabaseListenerMinReconnectInterval(), cfg.DatabaseListenerMaxReconnectDuration())
	ccOpts := evm.ChainSetOpts{
		Config:           cfg,
		Logger:           globalLogger,
		GormDB:           gormDB,
		SQLxDB:           db,
		ORM:              evm.NewORM(db),
		KeyStore:         keyStore.Eth(),
		EventBroadcaster: eventBroadcaster,
	}
	chainSet, err := evm.LoadChainSet(ccOpts)
	if err != nil {
		logger.Fatal(err)
	}
	externalInitiatorManager := webhook.NewExternalInitiatorManager(gormDB, utils.UnrestrictedClient)
	return chainlink.NewApplication(chainlink.ApplicationOpts{
		Config:                   cfg,
		ShutdownSignal:           shutdownSignal,
		GormDB:                   gormDB,
		SqlxDB:                   db,
		KeyStore:                 keyStore,
		ChainSet:                 chainSet,
		EventBroadcaster:         eventBroadcaster,
		Logger:                   globalLogger,
		ExternalInitiatorManager: externalInitiatorManager,
	})
}

// Runner implements the Run method.
type Runner interface {
	Run(chainlink.Application) error
}

// ChainlinkRunner is used to run the node application.
type ChainlinkRunner struct{}

// Run sets the log level based on config and starts the web router to listen
// for input and return data.
func (n ChainlinkRunner) Run(app chainlink.Application) error {
	config := app.GetConfig()
	mode := gin.ReleaseMode
	if config.Dev() && config.LogLevel().Level < zapcore.InfoLevel {
		mode = gin.DebugMode
	}
	gin.SetMode(mode)
	handler := web.Router(app.(*chainlink.ChainlinkApplication), prometheus)
	var g errgroup.Group

	if config.Port() == 0 && config.TLSPort() == 0 {
		log.Fatal("You must specify at least one port to listen on")
	}

	if config.Port() != 0 {
		g.Go(func() error { return runServer(handler, config.Port(), config.HTTPServerWriteTimeout()) })
	}

	if config.TLSPort() != 0 {
		g.Go(func() error {
			return runServerTLS(
				handler,
				config.TLSPort(),
				config.CertFile(),
				config.KeyFile(),
				config.HTTPServerWriteTimeout())
		})
	}

	return g.Wait()
}

func runServer(handler *gin.Engine, port uint16, writeTimeout time.Duration) error {
	logger.Infof("Listening and serving HTTP on port %d", port)
	server := createServer(handler, port, writeTimeout)
	err := server.ListenAndServe()
	logger.ErrorIf(err, "Error starting server")
	return err
}

func runServerTLS(handler *gin.Engine, port uint16, certFile, keyFile string, writeTimeout time.Duration) error {
	logger.Infof("Listening and serving HTTPS on port %d", port)
	server := createServer(handler, port, writeTimeout)
	err := server.ListenAndServeTLS(certFile, keyFile)
	logger.ErrorIf(err, "Error starting TLS server")
	return err
}

func createServer(handler *gin.Engine, port uint16, writeTimeout time.Duration) *http.Server {
	url := fmt.Sprintf(":%d", port)
	s := &http.Server{
		Addr:           url,
		Handler:        handler,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   writeTimeout,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	return s
}

// HTTPClient encapsulates all methods used to interact with a chainlink node API.
type HTTPClient interface {
	Get(string, ...map[string]string) (*http.Response, error)
	Post(string, io.Reader) (*http.Response, error)
	Put(string, io.Reader) (*http.Response, error)
	Patch(string, io.Reader, ...map[string]string) (*http.Response, error)
	Delete(string) (*http.Response, error)
}

type HTTPClientConfig interface {
	SessionCookieAuthenticatorConfig
}

type authenticatedHTTPClient struct {
	config         HTTPClientConfig
	client         *http.Client
	cookieAuth     CookieAuthenticator
	sessionRequest sessions.SessionRequest
}

// NewAuthenticatedHTTPClient uses the CookieAuthenticator to generate a sessionID
// which is then used for all subsequent HTTP API requests.
func NewAuthenticatedHTTPClient(config HTTPClientConfig, cookieAuth CookieAuthenticator, sessionRequest sessions.SessionRequest) HTTPClient {
	return &authenticatedHTTPClient{
		config:         config,
		client:         newHttpClient(config),
		cookieAuth:     cookieAuth,
		sessionRequest: sessionRequest,
	}
}

func newHttpClient(config SessionCookieAuthenticatorConfig) *http.Client {
	tr := &http.Transport{
		// User enables this at their own risk!
		// #nosec G402
		TLSClientConfig: &tls.Config{InsecureSkipVerify: config.InsecureSkipVerify()},
	}
	if config.InsecureSkipVerify() {
		fmt.Println("WARNING: INSECURE_SKIP_VERIFY is set to true, skipping SSL certificate verification.")
	}
	return &http.Client{Transport: tr}
}

// Get performs an HTTP Get using the authenticated HTTP client's cookie.
func (h *authenticatedHTTPClient) Get(path string, headers ...map[string]string) (*http.Response, error) {
	return h.doRequest("GET", path, nil, headers...)
}

// Post performs an HTTP Post using the authenticated HTTP client's cookie.
func (h *authenticatedHTTPClient) Post(path string, body io.Reader) (*http.Response, error) {
	return h.doRequest("POST", path, body)
}

// Put performs an HTTP Put using the authenticated HTTP client's cookie.
func (h *authenticatedHTTPClient) Put(path string, body io.Reader) (*http.Response, error) {
	return h.doRequest("PUT", path, body)
}

// Patch performs an HTTP Patch using the authenticated HTTP client's cookie.
func (h *authenticatedHTTPClient) Patch(path string, body io.Reader, headers ...map[string]string) (*http.Response, error) {
	return h.doRequest("PATCH", path, body, headers...)
}

// Delete performs an HTTP Delete using the authenticated HTTP client's cookie.
func (h *authenticatedHTTPClient) Delete(path string) (*http.Response, error) {
	return h.doRequest("DELETE", path, nil)
}

func (h *authenticatedHTTPClient) doRequest(verb, path string, body io.Reader, headerArgs ...map[string]string) (*http.Response, error) {
	var headers map[string]string
	if len(headerArgs) > 0 {
		headers = headerArgs[0]
	} else {
		headers = map[string]string{}
	}

	request, err := http.NewRequest(verb, h.config.ClientNodeURL()+path, body)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		request.Header.Add(key, value)
	}
	cookie, err := h.cookieAuth.Cookie()
	if err != nil {
		return nil, err
	} else if cookie != nil {
		request.AddCookie(cookie)
	}

	response, err := h.client.Do(request)
	if err != nil {
		return response, err
	}
	if response.StatusCode == http.StatusUnauthorized && (h.sessionRequest.Email != "" || h.sessionRequest.Password != "") {
		var cookieerr error
		cookie, cookieerr = h.cookieAuth.Authenticate(h.sessionRequest)
		if cookieerr != nil {
			return response, err
		}
		request.Header.Set("Cookie", "")
		request.AddCookie(cookie)
		response, err = h.client.Do(request)
		if err != nil {
			return response, err
		}
	}
	return response, nil
}

// CookieAuthenticator is the interface to generating a cookie to authenticate
// future HTTP requests.
type CookieAuthenticator interface {
	Cookie() (*http.Cookie, error)
	Authenticate(sessions.SessionRequest) (*http.Cookie, error)
}

type SessionCookieAuthenticatorConfig interface {
	ClientNodeURL() string
	InsecureSkipVerify() bool
}

// SessionCookieAuthenticator is a concrete implementation of CookieAuthenticator
// that retrieves a session id for the user with credentials from the session request.
type SessionCookieAuthenticator struct {
	config SessionCookieAuthenticatorConfig
	store  CookieStore
}

// NewSessionCookieAuthenticator creates a SessionCookieAuthenticator using the passed config
// and builder.
func NewSessionCookieAuthenticator(config SessionCookieAuthenticatorConfig, store CookieStore) CookieAuthenticator {
	return &SessionCookieAuthenticator{config: config, store: store}
}

// Cookie Returns the previously saved authentication cookie.
func (t *SessionCookieAuthenticator) Cookie() (*http.Cookie, error) {
	return t.store.Retrieve()
}

// Authenticate retrieves a session ID via a cookie and saves it to disk.
func (t *SessionCookieAuthenticator) Authenticate(sessionRequest sessions.SessionRequest) (*http.Cookie, error) {
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(sessionRequest)
	if err != nil {
		return nil, err
	}
	url := t.config.ClientNodeURL() + "/sessions"
	req, err := http.NewRequest("POST", url, b)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := newHttpClient(t.config)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer logger.ErrorIfCalling(resp.Body.Close)

	_, err = parseResponse(resp)
	if err != nil {
		return nil, err
	}

	cookies := resp.Cookies()
	if len(cookies) == 0 {
		return nil, errors.New("did not receive cookie with session id")
	}
	sc := web.FindSessionCookie(cookies)
	return sc, t.store.Save(sc)
}

// CookieStore is a place to store and retrieve cookies.
type CookieStore interface {
	Save(cookie *http.Cookie) error
	Retrieve() (*http.Cookie, error)
}

// MemoryCookieStore keeps a single cookie in memory
type MemoryCookieStore struct {
	Cookie *http.Cookie
}

// Save stores a cookie.
func (m *MemoryCookieStore) Save(cookie *http.Cookie) error {
	m.Cookie = cookie
	return nil
}

// Retrieve returns any Saved cookies.
func (m *MemoryCookieStore) Retrieve() (*http.Cookie, error) {
	return m.Cookie, nil
}

type DiskCookieConfig interface {
	RootDir() string
}

// DiskCookieStore saves a single cookie in the local cli working directory.
type DiskCookieStore struct {
	Config DiskCookieConfig
}

// Save stores a cookie.
func (d DiskCookieStore) Save(cookie *http.Cookie) error {
	return ioutil.WriteFile(d.cookiePath(), []byte(cookie.String()), 0600)
}

// Retrieve returns any Saved cookies.
func (d DiskCookieStore) Retrieve() (*http.Cookie, error) {
	b, err := ioutil.ReadFile(d.cookiePath())
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, multierr.Append(errors.New("unable to retrieve credentials, you must first login through the CLI"), err)
	}
	header := http.Header{}
	header.Add("Cookie", string(b))
	request := http.Request{Header: header}
	cookies := request.Cookies()
	if len(cookies) == 0 {
		return nil, errors.New("Cookie not in file, you must first login through the CLI")
	}
	return request.Cookies()[0], nil
}

func (d DiskCookieStore) cookiePath() string {
	return path.Join(d.Config.RootDir(), "cookie")
}

// SessionRequestBuilder is an interface that returns a SessionRequest,
// abstracting how session requests are generated, whether they be from
// the prompt or from a file.
type SessionRequestBuilder interface {
	Build(flag string) (sessions.SessionRequest, error)
}

type promptingSessionRequestBuilder struct {
	prompter Prompter
}

// NewPromptingSessionRequestBuilder uses a prompter, often via terminal,
// to solicit information from a user to generate the SessionRequest.
func NewPromptingSessionRequestBuilder(prompter Prompter) SessionRequestBuilder {
	return promptingSessionRequestBuilder{prompter}
}

func (p promptingSessionRequestBuilder) Build(string) (sessions.SessionRequest, error) {
	email := p.prompter.Prompt("Enter email: ")
	pwd := p.prompter.PasswordPrompt("Enter password: ")
	return sessions.SessionRequest{Email: email, Password: pwd}, nil
}

type fileSessionRequestBuilder struct{}

// NewFileSessionRequestBuilder pulls credentials from a file to generate a SessionRequest.
func NewFileSessionRequestBuilder() SessionRequestBuilder {
	return fileSessionRequestBuilder{}
}

func (f fileSessionRequestBuilder) Build(file string) (sessions.SessionRequest, error) {
	return credentialsFromFile(file)
}

// APIInitializer is the interface used to create the API User credentials
// needed to access the API. Does nothing if API user already exists.
type APIInitializer interface {
	// Initialize creates a new user for API access, or does nothing if one exists.
	Initialize(orm sessions.ORM) (sessions.User, error)
}

type promptingAPIInitializer struct {
	prompter Prompter
}

// NewPromptingAPIInitializer creates a concrete instance of APIInitializer
// that uses the terminal to solicit credentials from the user.
func NewPromptingAPIInitializer(prompter Prompter) APIInitializer {
	return &promptingAPIInitializer{prompter: prompter}
}

// Initialize uses the terminal to get credentials that it then saves in the store.
func (t *promptingAPIInitializer) Initialize(orm sessions.ORM) (sessions.User, error) {
	if user, err := orm.FindUser(); err == nil {
		return user, err
	}

	if !t.prompter.IsTerminal() {
		return sessions.User{}, ErrorNoAPICredentialsAvailable
	}

	for {
		email := t.prompter.Prompt("Enter API Email: ")
		pwd := t.prompter.PasswordPrompt("Enter API Password: ")
		user, err := sessions.NewUser(email, pwd)
		if err != nil {
			fmt.Println("Error creating API user: ", err)
			continue
		}
		if err = orm.CreateUser(&user); err != nil {
			fmt.Println("Error creating API user: ", err)
		}
		return user, err
	}
}

type fileAPIInitializer struct {
	file string
}

// NewFileAPIInitializer creates a concrete instance of APIInitializer
// that pulls API user credentials from the passed file path.
func NewFileAPIInitializer(file string) APIInitializer {
	return fileAPIInitializer{file: file}
}

func (f fileAPIInitializer) Initialize(orm sessions.ORM) (sessions.User, error) {
	if user, err := orm.FindUser(); err == nil {
		return user, err
	}

	request, err := credentialsFromFile(f.file)
	if err != nil {
		return sessions.User{}, err
	}

	user, err := sessions.NewUser(request.Email, request.Password)
	if err != nil {
		return user, err
	}
	return user, orm.CreateUser(&user)
}

var ErrNoCredentialFile = errors.New("no API user credential file was passed")

func credentialsFromFile(file string) (sessions.SessionRequest, error) {
	if len(file) == 0 {
		return sessions.SessionRequest{}, ErrNoCredentialFile
	}

	logger.Debug("Initializing API credentials from ", file)
	dat, err := ioutil.ReadFile(file)
	if err != nil {
		return sessions.SessionRequest{}, err
	}
	lines := strings.Split(string(dat), "\n")
	if len(lines) < 2 {
		return sessions.SessionRequest{}, fmt.Errorf("malformed API credentials file does not have at least two lines at %s", file)
	}
	credentials := sessions.SessionRequest{
		Email:    strings.TrimSpace(lines[0]),
		Password: strings.TrimSpace(lines[1]),
	}
	return credentials, nil
}

// ChangePasswordPrompter is an interface primarily used for DI to obtain a
// password change request from the User.
type ChangePasswordPrompter interface {
	Prompt() (web.UpdatePasswordRequest, error)
}

// NewChangePasswordPrompter returns the production password change request prompter
func NewChangePasswordPrompter() ChangePasswordPrompter {
	prompter := NewTerminalPrompter()
	return changePasswordPrompter{prompter: prompter}
}

type changePasswordPrompter struct {
	prompter Prompter
}

func (c changePasswordPrompter) Prompt() (web.UpdatePasswordRequest, error) {
	fmt.Println("Changing your chainlink account password.")
	fmt.Println("NOTE: This will terminate any other sessions.")
	oldPassword := c.prompter.PasswordPrompt("Password:")

	fmt.Println("Now enter your **NEW** password")
	newPassword := c.prompter.PasswordPrompt("Password:")
	confirmPassword := c.prompter.PasswordPrompt("Confirmation:")

	if newPassword != confirmPassword {
		return web.UpdatePasswordRequest{}, errors.New("new password and confirmation did not match")
	}

	return web.UpdatePasswordRequest{
		OldPassword: oldPassword,
		NewPassword: newPassword,
	}, nil
}

// PasswordPrompter is an interface primarily used for DI to obtain a password
// from the User.
type PasswordPrompter interface {
	Prompt() string
}

// NewPasswordPrompter returns the production password change request prompter
func NewPasswordPrompter() PasswordPrompter {
	prompter := NewTerminalPrompter()
	return passwordPrompter{prompter: prompter}
}

type passwordPrompter struct {
	prompter Prompter
}

func (c passwordPrompter) Prompt() string {
	return c.prompter.PasswordPrompt("Password:")
}

func confirmAction(c *clipkg.Context) bool {
	if len(c.String("yes")) > 0 {
		yes, err := strconv.ParseBool(c.String("yes"))
		if err == nil && yes {
			return true
		}
	}

	prompt := NewTerminalPrompter()
	var answer string
	for {
		answer = prompt.Prompt("Are you sure? This action is irreversible! (yes/no) ")
		if answer == "yes" {
			return true
		} else if answer == "no" {
			return false
		} else {
			fmt.Printf("%s is not valid. Please type yes or no\n", answer)
		}
	}
}
