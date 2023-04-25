package cmd

import (
	"bytes"
	"context"
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Depado/ginprom"
	"github.com/Masterminds/semver/v3"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	clipkg "github.com/urfave/cli"
	"go.uber.org/multierr"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/v2/core/chains/cosmos"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	"github.com/smartcontractkit/chainlink/v2/core/chains/solana"
	"github.com/smartcontractkit/chainlink/v2/core/chains/starknet"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/logger/audit"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/periodicbackup"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/services/versioning"
	"github.com/smartcontractkit/chainlink/v2/core/services/webhook"
	"github.com/smartcontractkit/chainlink/v2/core/sessions"
	"github.com/smartcontractkit/chainlink/v2/core/static"
	"github.com/smartcontractkit/chainlink/v2/core/store/migrate"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
	clhttp "github.com/smartcontractkit/chainlink/v2/core/utils/http"
	"github.com/smartcontractkit/chainlink/v2/core/web"
)

var prometheus *ginprom.Prometheus
var initOnce sync.Once

func initPrometheus(cfg config.GeneralConfig) {
	prometheus = ginprom.New(ginprom.Namespace("service"), ginprom.Token(cfg.PrometheusAuthToken()))
}

var (
	// ErrorNoAPICredentialsAvailable is returned when not run from a terminal
	// and no API credentials have been provided
	ErrorNoAPICredentialsAvailable = errors.New("API credentials must be supplied")
)

// Client is the shell for the node, local commands and remote commands.
type Client struct {
	Renderer
	Config                         chainlink.GeneralConfig // initialized in Before
	Logger                         logger.Logger           // initialized in Before
	CloseLogger                    func() error            // called in After
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

	configInitialized bool
}

func (cli *Client) errorOut(err error) error {
	if err != nil {
		return clipkg.NewExitError(err.Error(), 1)
	}
	return nil
}

func (cli *Client) setConfig(opts *chainlink.GeneralConfigOpts, configFiles []string, secretsFile string) error {
	if err := loadOpts(opts, configFiles...); err != nil {
		return err
	}

	secretsTOML := ""
	if secretsFile != "" {
		b, err := os.ReadFile(secretsFile)
		if err != nil {
			return errors.Wrapf(err, "failed to read secrets file: %s", secretsFile)
		}
		secretsTOML = string(b)
	}
	err := opts.ParseSecrets(secretsTOML)
	if err != nil {
		return err
	}
	if cfg, lggr, closeLggr, err := opts.NewAndLogger(); err != nil {
		return err
	} else {
		cli.Config = cfg
		cli.Logger = lggr
		cli.CloseLogger = closeLggr
	}
	return nil
}

// AppFactory implements the NewApplication method.
type AppFactory interface {
	NewApplication(ctx context.Context, cfg chainlink.GeneralConfig, appLggr logger.Logger, db *sqlx.DB) (chainlink.Application, error)
}

// ChainlinkAppFactory is used to create a new Application.
type ChainlinkAppFactory struct{}

// NewApplication returns a new instance of the node with the given config.
func (n ChainlinkAppFactory) NewApplication(ctx context.Context, cfg chainlink.GeneralConfig, appLggr logger.Logger, db *sqlx.DB) (app chainlink.Application, err error) {
	// Avoid double initializations.
	initOnce.Do(func() {
		initPrometheus(cfg)
	})

	keyStore := keystore.New(db, utils.GetScryptParams(cfg), appLggr, cfg)

	// Set up the versioning Configs
	verORM := versioning.NewORM(db, appLggr, cfg.DatabaseDefaultQueryTimeout())

	if static.Version != static.Unset {
		var appv, dbv *semver.Version
		appv, dbv, err = versioning.CheckVersion(db, appLggr, static.Version)
		if err != nil {
			// Exit immediately and don't touch the database if the app version is too old
			return nil, errors.Wrap(err, "CheckVersion")
		}

		// Take backup if app version is newer than DB version
		// Need to do this BEFORE migration
		if cfg.DatabaseBackupMode() != config.DatabaseBackupModeNone && cfg.DatabaseBackupOnVersionUpgrade() {
			if err = takeBackupIfVersionUpgrade(cfg, appLggr, appv, dbv); err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					appLggr.Debugf("Failed to find any node version in the DB: %w", err)
				} else if strings.Contains(err.Error(), "relation \"node_versions\" does not exist") {
					appLggr.Debugf("Failed to find any node version in the DB, the node_versions table does not exist yet: %w", err)
				} else {
					return nil, errors.Wrap(err, "initializeORM#FindLatestNodeVersion")
				}
			}
		}
	}

	// Migrate the database
	if cfg.MigrateDatabase() {
		if err = migrate.Migrate(db.DB, appLggr); err != nil {
			return nil, errors.Wrap(err, "initializeORM#Migrate")
		}
	}

	// Update to latest version
	if static.Version != static.Unset {
		version := versioning.NewNodeVersion(static.Version)
		if err = verORM.UpsertNodeVersion(version); err != nil {
			return nil, errors.Wrap(err, "UpsertNodeVersion")
		}
	}

	mailMon := utils.NewMailboxMonitor(cfg.AppID().String())

	// Upsert EVM chains/nodes from ENV, necessary for backwards compatibility
	if cfg.EVMEnabled() {
		var ids []utils.Big
		for _, c := range cfg.EVMConfigs() {
			c := c
			ids = append(ids, *c.ChainID)
		}
		if len(ids) > 0 {
			if err = evm.EnsureChains(db, appLggr, cfg, ids); err != nil {
				return nil, errors.Wrap(err, "failed to setup EVM chains")
			}
		}
	}

	eventBroadcaster := pg.NewEventBroadcaster(cfg.DatabaseURL(), cfg.DatabaseListenerMinReconnectInterval(), cfg.DatabaseListenerMaxReconnectDuration(), appLggr, cfg.AppID())
	ccOpts := evm.ChainSetOpts{
		Config:           cfg,
		Logger:           appLggr,
		DB:               db,
		KeyStore:         keyStore.Eth(),
		EventBroadcaster: eventBroadcaster,
		MailMon:          mailMon,
	}
	var chains chainlink.Chains
	chains.EVM, err = evm.NewTOMLChainSet(ctx, ccOpts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load EVM chainset")
	}

	if cfg.CosmosEnabled() {
		cosmosLggr := appLggr.Named("Cosmos")
		opts := cosmos.ChainSetOpts{
			Config:           cfg,
			Logger:           cosmosLggr,
			DB:               db,
			KeyStore:         keyStore.Cosmos(),
			EventBroadcaster: eventBroadcaster,
		}
		cfgs := cfg.CosmosConfigs()
		opts.Configs = cosmos.NewConfigs(cfgs)
		chains.Cosmos, err = cosmos.NewChainSet(opts, cfgs)
		if err != nil {
			return nil, errors.Wrap(err, "failed to load Cosmos chainset")
		}
	}

	if cfg.SolanaEnabled() {
		solLggr := appLggr.Named("Solana")
		opts := solana.ChainSetOpts{
			Logger:   solLggr,
			KeyStore: &keystore.SolanaSigner{keyStore.Solana()},
		}
		cfgs := cfg.SolanaConfigs()
		var ids []string
		for _, c := range cfgs {
			c := c
			ids = append(ids, *c.ChainID)
		}
		if len(ids) > 0 {
			if err = solana.EnsureChains(db, solLggr, cfg, ids); err != nil {
				return nil, errors.Wrap(err, "failed to setup Solana chains")
			}
		}
		opts.Configs = solana.NewConfigs(cfgs)
		chains.Solana, err = solana.NewChainSet(opts, cfgs)
		if err != nil {
			return nil, errors.Wrap(err, "failed to load Solana chainset")
		}
	}

	if cfg.StarkNetEnabled() {
		starkLggr := appLggr.Named("StarkNet")
		opts := starknet.ChainSetOpts{
			Config:   cfg,
			Logger:   starkLggr,
			KeyStore: keyStore.StarkNet(),
		}
		cfgs := cfg.StarknetConfigs()
		var ids []string
		for _, c := range cfgs {
			c := c
			ids = append(ids, *c.ChainID)
		}
		if len(ids) > 0 {
			if err = starknet.EnsureChains(db, starkLggr, cfg, ids); err != nil {
				return nil, errors.Wrap(err, "failed to setup StarkNet chains")
			}
		}
		opts.Configs = starknet.NewConfigs(cfgs)
		chains.StarkNet, err = starknet.NewChainSet(opts, cfgs)
		if err != nil {
			return nil, errors.Wrap(err, "failed to load StarkNet chainset")
		}
	}

	// Configure and optionally start the audit log forwarder service
	auditLogger, err := audit.NewAuditLogger(appLggr, cfg)
	if err != nil {
		return nil, err
	}

	restrictedClient := clhttp.NewRestrictedHTTPClient(cfg, appLggr)
	unrestrictedClient := clhttp.NewUnrestrictedHTTPClient()
	externalInitiatorManager := webhook.NewExternalInitiatorManager(db, unrestrictedClient, appLggr, cfg)
	return chainlink.NewApplication(chainlink.ApplicationOpts{
		Config:                   cfg,
		SqlxDB:                   db,
		KeyStore:                 keyStore,
		Chains:                   chains,
		EventBroadcaster:         eventBroadcaster,
		MailMon:                  mailMon,
		Logger:                   appLggr,
		AuditLogger:              auditLogger,
		ExternalInitiatorManager: externalInitiatorManager,
		Version:                  static.Version,
		RestrictedHTTPClient:     restrictedClient,
		UnrestrictedHTTPClient:   unrestrictedClient,
		SecretGenerator:          chainlink.FilePersistedSecretGenerator{},
	})
}

func takeBackupIfVersionUpgrade(cfg periodicbackup.Config, lggr logger.Logger, appv, dbv *semver.Version) (err error) {
	if appv == nil {
		lggr.Debug("Application version is missing, skipping automatic DB backup.")
		return nil
	}
	if dbv == nil {
		lggr.Debug("Database version is missing, skipping automatic DB backup.")
		return nil
	}
	if !appv.GreaterThan(dbv) {
		lggr.Debugf("Application version %s is older or equal to database version %s, skipping automatic DB backup.", appv.String(), dbv.String())
		return nil
	}
	lggr.Infof("Upgrade detected: application version %s is newer than database version %s, taking automatic DB backup. To skip automatic database backup before version upgrades, set Database.Backup.OnVersionUpgrade=false. To disable backups entirely set Database.Backup.Mode=none.", appv.String(), dbv.String())

	databaseBackup, err := periodicbackup.NewDatabaseBackup(cfg, lggr)
	if err != nil {
		return errors.Wrap(err, "takeBackupIfVersionUpgrade failed")
	}
	return databaseBackup.RunBackup(appv.String())
}

// Runner implements the Run method.
type Runner interface {
	Run(context.Context, chainlink.Application) error
}

// ChainlinkRunner is used to run the node application.
type ChainlinkRunner struct{}

// Run sets the log level based on config and starts the web router to listen
// for input and return data.
func (n ChainlinkRunner) Run(ctx context.Context, app chainlink.Application) error {
	config := app.GetConfig()

	mode := gin.ReleaseMode
	if config.Dev() && config.LogLevel() < zapcore.InfoLevel {
		mode = gin.DebugMode
	}
	gin.SetMode(mode)
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		app.GetLogger().Debugf("%-6s %-25s --> %s (%d handlers)", httpMethod, absolutePath, handlerName, nuHandlers)
	}

	if err := sentryInit(config); err != nil {
		return errors.Wrap(err, "failed to initialize sentry")
	}

	if config.Port() == 0 && config.TLSPort() == 0 {
		return errors.New("You must specify at least one port to listen on")
	}

	handler, err := web.NewRouter(app, prometheus)
	if err != nil {
		return errors.Wrap(err, "failed to create web router")
	}
	server := server{handler: handler, lggr: app.GetLogger()}

	g, gCtx := errgroup.WithContext(ctx)
	if config.Port() != 0 {
		go tryRunServerUntilCancelled(gCtx, app.GetLogger(), config, func() error {
			return server.run(config.Port(), config.HTTPServerWriteTimeout())
		})
	}

	if config.TLSPort() != 0 {
		go tryRunServerUntilCancelled(gCtx, app.GetLogger(), config, func() error {
			return server.runTLS(
				config.TLSPort(),
				config.CertFile(),
				config.KeyFile(),
				config.HTTPServerWriteTimeout())
		})
	}

	g.Go(func() error {
		<-gCtx.Done()
		var err error
		if server.httpServer != nil {
			err = errors.WithStack(server.httpServer.Shutdown(context.Background()))
		}
		if server.tlsServer != nil {
			err = multierr.Combine(err, errors.WithStack(server.tlsServer.Shutdown(context.Background())))
		}
		return err
	})

	return errors.WithStack(g.Wait())
}

func sentryInit(cfg config.BasicConfig) error {
	sentrydsn := cfg.SentryDSN()
	if sentrydsn == "" {
		// Do not initialize sentry at all if the DSN is missing
		return nil
	}

	var sentryenv string
	if env := cfg.SentryEnvironment(); env != "" {
		sentryenv = env
	} else if cfg.Dev() {
		sentryenv = "dev"
	} else {
		sentryenv = "prod"
	}

	var sentryrelease string
	if release := cfg.SentryRelease(); release != "" {
		sentryrelease = release
	} else {
		sentryrelease = static.Version
	}

	return sentry.Init(sentry.ClientOptions{
		// AttachStacktrace is needed to send stacktrace alongside panics
		AttachStacktrace: true,
		Dsn:              sentrydsn,
		Environment:      sentryenv,
		Release:          sentryrelease,
		Debug:            cfg.SentryDebug(),
	})
}

func tryRunServerUntilCancelled(ctx context.Context, lggr logger.Logger, cfg config.BasicConfig, runServer func() error) {
	for {
		// try calling runServer() and log error if any
		if err := runServer(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				lggr.Criticalf("Error starting server: %v", err)
			}
		}
		// if ctx is cancelled, we must leave the loop
		select {
		case <-ctx.Done():
			return
		case <-time.After(cfg.DefaultHTTPTimeout().Duration()):
			// pause between attempts, default 15s
		}
	}
}

type server struct {
	httpServer *http.Server
	tlsServer  *http.Server
	handler    *gin.Engine
	lggr       logger.Logger
}

func (s *server) run(port uint16, writeTimeout time.Duration) error {
	s.lggr.Infof("Listening and serving HTTP on port %d", port)
	s.httpServer = createServer(s.handler, port, writeTimeout)
	err := s.httpServer.ListenAndServe()
	return errors.Wrap(err, "failed to run plaintext HTTP server")
}

func (s *server) runTLS(port uint16, certFile, keyFile string, writeTimeout time.Duration) error {
	s.lggr.Infof("Listening and serving HTTPS on port %d", port)
	s.tlsServer = createServer(s.handler, port, writeTimeout)
	err := s.tlsServer.ListenAndServeTLS(certFile, keyFile)
	return errors.Wrap(err, "failed to run TLS server (NOTE: you can disable TLS server completely and silence these errors by setting WebServer.TLS.HTTPSPort=0 in your config)")
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

type authenticatedHTTPClient struct {
	client         *http.Client
	cookieAuth     CookieAuthenticator
	sessionRequest sessions.SessionRequest
	remoteNodeURL  url.URL
}

// NewAuthenticatedHTTPClient uses the CookieAuthenticator to generate a sessionID
// which is then used for all subsequent HTTP API requests.
func NewAuthenticatedHTTPClient(lggr logger.Logger, clientOpts ClientOpts, cookieAuth CookieAuthenticator, sessionRequest sessions.SessionRequest) HTTPClient {
	return &authenticatedHTTPClient{
		client:         newHttpClient(lggr, clientOpts.InsecureSkipVerify),
		cookieAuth:     cookieAuth,
		sessionRequest: sessionRequest,
		remoteNodeURL:  clientOpts.RemoteNodeURL,
	}
}

func newHttpClient(lggr logger.Logger, insecureSkipVerify bool) *http.Client {
	tr := &http.Transport{
		// User enables this at their own risk!
		// #nosec G402
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecureSkipVerify},
	}
	if insecureSkipVerify {
		lggr.Warn("InsecureSkipVerify is on, skipping SSL certificate verification.")
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

	request, err := http.NewRequest(verb, h.remoteNodeURL.String()+path, body)
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
	Logout() error
}

type ClientOpts struct {
	RemoteNodeURL      url.URL
	InsecureSkipVerify bool
}

// SessionCookieAuthenticator is a concrete implementation of CookieAuthenticator
// that retrieves a session id for the user with credentials from the session request.
type SessionCookieAuthenticator struct {
	config ClientOpts
	store  CookieStore
	lggr   logger.SugaredLogger
}

// NewSessionCookieAuthenticator creates a SessionCookieAuthenticator using the passed config
// and builder.
func NewSessionCookieAuthenticator(config ClientOpts, store CookieStore, lggr logger.Logger) CookieAuthenticator {
	return &SessionCookieAuthenticator{config: config, store: store, lggr: logger.Sugared(lggr)}
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
	url := t.config.RemoteNodeURL.String() + "/sessions"
	req, err := http.NewRequest("POST", url, b)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := newHttpClient(t.lggr, t.config.InsecureSkipVerify)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer t.lggr.ErrorIfFn(resp.Body.Close, "Error closing Authenticate response body")

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

// Deletes any stored session
func (t *SessionCookieAuthenticator) Logout() error {
	return t.store.Reset()
}

// CookieStore is a place to store and retrieve cookies.
type CookieStore interface {
	Save(cookie *http.Cookie) error
	Retrieve() (*http.Cookie, error)
	Reset() error
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

// Removes any stored cookie.
func (m *MemoryCookieStore) Reset() error {
	m.Cookie = nil
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
	return os.WriteFile(d.cookiePath(), []byte(cookie.String()), 0600)
}

// Removes any stored cookie.
func (d DiskCookieStore) Reset() error {
	// Write empty bytes
	return os.WriteFile(d.cookiePath(), []byte(""), 0600)
}

// Retrieve returns any Saved cookies.
func (d DiskCookieStore) Retrieve() (*http.Cookie, error) {
	b, err := os.ReadFile(d.cookiePath())
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

type UserCache struct {
	dir        string
	lggr       func() logger.Logger // func b/c we don't have the final logger at construction time
	ensureOnce sync.Once
}

func NewUserCache(subdir string, lggr func() logger.Logger) (*UserCache, error) {
	cd, err := os.UserCacheDir()
	if err != nil {
		return nil, err
	}
	return &UserCache{dir: filepath.Join(cd, "chainlink", subdir), lggr: lggr}, nil
}

func (cs *UserCache) ensure() {
	if err := os.MkdirAll(cs.dir, 0700); err != nil {
		cs.lggr().Errorw("Failed to make user cache dir", "dir", cs.dir, "err", err)
	}
}

func (cs *UserCache) RootDir() string {
	cs.ensureOnce.Do(cs.ensure)
	return cs.dir
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

type fileSessionRequestBuilder struct {
	lggr logger.Logger
}

// NewFileSessionRequestBuilder pulls credentials from a file to generate a SessionRequest.
func NewFileSessionRequestBuilder(lggr logger.Logger) SessionRequestBuilder {
	return &fileSessionRequestBuilder{lggr: lggr}
}

func (f *fileSessionRequestBuilder) Build(file string) (sessions.SessionRequest, error) {
	return credentialsFromFile(file, f.lggr.With("file", file))
}

// APIInitializer is the interface used to create the API User credentials
// needed to access the API. Does nothing if API user already exists.
type APIInitializer interface {
	// Initialize creates a new user for API access, or does nothing if one exists.
	Initialize(orm sessions.ORM, lggr logger.Logger) (sessions.User, error)
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
func (t *promptingAPIInitializer) Initialize(orm sessions.ORM, lggr logger.Logger) (sessions.User, error) {
	// Load list of users to determine which to assume, or if a user needs to be created
	dbUsers, err := orm.ListUsers()
	if err != nil {
		return sessions.User{}, err
	}

	// If there are no users in the database, prompt for initial admin user creation
	if len(dbUsers) == 0 {
		if !t.prompter.IsTerminal() {
			return sessions.User{}, ErrorNoAPICredentialsAvailable
		}

		for {
			email := t.prompter.Prompt("Enter API Email: ")
			pwd := t.prompter.PasswordPrompt("Enter API Password: ")
			// On a fresh DB, create an admin user
			user, err2 := sessions.NewUser(email, pwd, sessions.UserRoleAdmin)
			if err2 != nil {
				lggr.Errorw("Error creating API user", "err", err2)
				continue
			}
			if err = orm.CreateUser(&user); err != nil {
				lggr.Errorf("Error creating API user: ", err, "err")
			}
			return user, err
		}
	}

	// Attempt to contextually return the correct admin user, CLI access here implies admin
	if adminUser, found := attemptAssumeAdminUser(dbUsers, lggr); found {
		return adminUser, nil
	}

	// Otherwise, multiple admin users exist, prompt for which to use
	email := t.prompter.Prompt("Enter email of API user account to assume: ")
	user, err := orm.FindUser(email)

	if err != nil {
		return sessions.User{}, err
	}
	return user, nil
}

type fileAPIInitializer struct {
	file string
}

// NewFileAPIInitializer creates a concrete instance of APIInitializer
// that pulls API user credentials from the passed file path.
func NewFileAPIInitializer(file string) APIInitializer {
	return fileAPIInitializer{file: file}
}

func (f fileAPIInitializer) Initialize(orm sessions.ORM, lggr logger.Logger) (sessions.User, error) {
	request, err := credentialsFromFile(f.file, lggr)
	if err != nil {
		return sessions.User{}, err
	}

	// Load list of users to determine which to assume, or if a user needs to be created
	dbUsers, err := orm.ListUsers()
	if err != nil {
		return sessions.User{}, err
	}

	// If there are no users in the database, create initial admin user from session request from file creds
	if len(dbUsers) == 0 {
		user, err2 := sessions.NewUser(request.Email, request.Password, sessions.UserRoleAdmin)
		if err2 != nil {
			return user, errors.Wrap(err2, "failed to instantiate new user")
		}
		return user, orm.CreateUser(&user)
	}

	// Attempt to contextually return the correct admin user, CLI access here implies admin
	if adminUser, found := attemptAssumeAdminUser(dbUsers, lggr); found {
		return adminUser, nil
	}

	// Otherwise, multiple admin users exist, attempt to load email specified in session request
	user, err := orm.FindUser(request.Email)
	if err != nil {
		return sessions.User{}, err
	}
	return user, nil
}

func attemptAssumeAdminUser(users []sessions.User, lggr logger.Logger) (sessions.User, bool) {
	if len(users) == 0 {
		return sessions.User{}, false
	}

	// If there is only a single DB user, select it within the context of CLI
	if len(users) == 1 {
		lggr.Infow("Defaulted to assume single DB API User", "email", users[0].Email)
		return users[0], true
	}

	// If there is only one admin user, use it within the context of CLI
	var singleAdmin sessions.User
	populatedUser := false
	for _, user := range users {
		if user.Role == sessions.UserRoleAdmin {
			// If multiple admin users found, don't use assume any and clear to continue to prompt
			if populatedUser {
				// Clear flag to skip return
				populatedUser = false
				break
			}
			singleAdmin = user
			populatedUser = true
		}
	}
	if populatedUser {
		lggr.Infow("Defaulted to assume single DB admin API User", "email", singleAdmin)
		return singleAdmin, true
	}

	return sessions.User{}, false
}

var ErrNoCredentialFile = errors.New("no API user credential file was passed")

func credentialsFromFile(file string, lggr logger.Logger) (sessions.SessionRequest, error) {
	if len(file) == 0 {
		return sessions.SessionRequest{}, ErrNoCredentialFile
	}

	lggr.Debug("Initializing API credentials")
	dat, err := os.ReadFile(file)
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
