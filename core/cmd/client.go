package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/web"

	"github.com/gin-gonic/gin"
	clipkg "github.com/urfave/cli"
	"go.uber.org/multierr"
	"golang.org/x/sync/errgroup"
)

var (
	// ErrorNoAPICredentialsAvailable is returned when not run from a terminal
	// and no API credentials have been provided
	ErrorNoAPICredentialsAvailable = errors.New("API credentials must be supplied")
)

// Client is the shell for the node, local commands and remote commands.
type Client struct {
	Renderer
	Config                         *orm.Config
	AppFactory                     AppFactory
	KeyStoreAuthenticator          KeyStoreAuthenticator
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
	NewApplication(*orm.Config, ...func(chainlink.Application)) chainlink.Application
}

// ChainlinkAppFactory is used to create a new Application.
type ChainlinkAppFactory struct{}

// NewApplication returns a new instance of the node with the given config.
func (n ChainlinkAppFactory) NewApplication(config *orm.Config, onConnectCallbacks ...func(chainlink.Application)) chainlink.Application {
	var ethClient eth.Client
	if config.EthereumDisabled() {
		logger.Info("ETH_DISABLED is set, using Null eth.Client")
		ethClient = &eth.NullClient{}
	} else {
		var err error
		ethClient, err = eth.NewClient(config.EthereumURL(), config.EthereumSecondaryURL())
		if err != nil {
			logger.Fatal(fmt.Sprintf("Unable to create ETH client: %+v", err))
		}
	}

	advisoryLock := postgres.NewAdvisoryLock(config.DatabaseURL())
	return chainlink.NewApplication(config, ethClient, advisoryLock, onConnectCallbacks...)
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
	gin.SetMode(app.GetStore().Config.LogLevel().ForGin())
	handler := web.Router(app.(*chainlink.ChainlinkApplication))
	config := app.GetStore().Config
	var g errgroup.Group

	if config.Port() == 0 && config.TLSPort() == 0 {
		log.Fatal("You must specify at least one port to listen on")
	}

	if config.Port() != 0 {
		g.Go(func() error { return runServer(handler, config.Port()) })
	}

	if config.TLSPort() != 0 {
		g.Go(func() error {
			return runServerTLS(
				handler,
				config.TLSPort(),
				config.CertFile(),
				config.KeyFile())
		})
	}

	return g.Wait()
}

func runServer(handler *gin.Engine, port uint16) error {
	logger.Infof("Listening and serving HTTP on port %d", port)
	server := createServer(handler, port)
	err := server.ListenAndServe()
	logger.ErrorIf(err)
	return err
}

func runServerTLS(handler *gin.Engine, port uint16, certFile, keyFile string) error {
	logger.Infof("Listening and serving HTTPS on port %d", port)
	server := createServer(handler, port)
	err := server.ListenAndServeTLS(certFile, keyFile)
	logger.ErrorIf(err)
	return err
}

func createServer(handler *gin.Engine, port uint16) *http.Server {
	url := fmt.Sprintf(":%d", port)
	s := &http.Server{
		Addr:           url,
		Handler:        handler,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
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
	config     orm.ConfigReader
	client     *http.Client
	cookieAuth CookieAuthenticator
}

// NewAuthenticatedHTTPClient uses the CookieAuthenticator to generate a sessionID
// which is then used for all subsequent HTTP API requests.
func NewAuthenticatedHTTPClient(config orm.ConfigReader, cookieAuth CookieAuthenticator) HTTPClient {
	return &authenticatedHTTPClient{
		config:     config,
		client:     &http.Client{},
		cookieAuth: cookieAuth,
	}
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
	cookie, err := h.cookieAuth.Cookie()
	if err != nil {
		return nil, err
	}

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
	request.AddCookie(cookie)
	return h.client.Do(request)
}

// CookieAuthenticator is the interface to generating a cookie to authenticate
// future HTTP requests.
type CookieAuthenticator interface {
	Cookie() (*http.Cookie, error)
	Authenticate(models.SessionRequest) (*http.Cookie, error)
}

// SessionCookieAuthenticator is a concrete implementation of CookieAuthenticator
// that retrieves a session id for the user with credentials from the session request.
type SessionCookieAuthenticator struct {
	config *orm.Config
	store  CookieStore
}

// NewSessionCookieAuthenticator creates a SessionCookieAuthenticator using the passed config
// and builder.
func NewSessionCookieAuthenticator(config *orm.Config, store CookieStore) CookieAuthenticator {
	return &SessionCookieAuthenticator{config: config, store: store}
}

// Cookie Returns the previously saved authentication cookie.
func (t *SessionCookieAuthenticator) Cookie() (*http.Cookie, error) {
	return t.store.Retrieve()
}

// Authenticate retrieves a session ID via a cookie and saves it to disk.
func (t *SessionCookieAuthenticator) Authenticate(sessionRequest models.SessionRequest) (*http.Cookie, error) {
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

	client := http.Client{Timeout: 30 * time.Second}
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

// DiskCookieStore saves a single cookie in the local cli working directory.
type DiskCookieStore struct {
	Config *orm.Config
}

// Save stores a cookie.
func (d DiskCookieStore) Save(cookie *http.Cookie) error {
	return ioutil.WriteFile(d.cookiePath(), []byte(cookie.String()), 0600)
}

// Retrieve returns any Saved cookies.
func (d DiskCookieStore) Retrieve() (*http.Cookie, error) {
	b, err := ioutil.ReadFile(d.cookiePath())
	if err != nil {
		return nil, multierr.Append(errors.New("unable to retrieve credentials, have you logged in?"), err)
	}
	header := http.Header{}
	header.Add("Cookie", string(b))
	request := http.Request{Header: header}
	cookies := request.Cookies()
	if len(cookies) == 0 {
		return nil, errors.New("Cookie not in file, have you logged in?")
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
	Build(flag string) (models.SessionRequest, error)
}

type promptingSessionRequestBuilder struct {
	prompter Prompter
}

// NewPromptingSessionRequestBuilder uses a prompter, often via terminal,
// to solicit information from a user to generate the SessionRequest.
func NewPromptingSessionRequestBuilder(prompter Prompter) SessionRequestBuilder {
	return promptingSessionRequestBuilder{prompter}
}

func (p promptingSessionRequestBuilder) Build(string) (models.SessionRequest, error) {
	email := p.prompter.Prompt("Enter email: ")
	pwd := p.prompter.PasswordPrompt("Enter password: ")
	return models.SessionRequest{Email: email, Password: pwd}, nil
}

type fileSessionRequestBuilder struct{}

// NewFileSessionRequestBuilder pulls credentials from a file to generate a SessionRequest.
func NewFileSessionRequestBuilder() SessionRequestBuilder {
	return fileSessionRequestBuilder{}
}

func (f fileSessionRequestBuilder) Build(file string) (models.SessionRequest, error) {
	return credentialsFromFile(file)
}

// APIInitializer is the interface used to create the API User credentials
// needed to access the API. Does nothing if API user already exists.
type APIInitializer interface {
	// Initialize creates a new user for API access, or does nothing if one exists.
	Initialize(store *store.Store) (models.User, error)
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
func (t *promptingAPIInitializer) Initialize(store *store.Store) (models.User, error) {
	if user, err := store.FindUser(); err == nil {
		return user, err
	}

	if !t.prompter.IsTerminal() {
		return models.User{}, ErrorNoAPICredentialsAvailable
	}

	for {
		email := t.prompter.Prompt("Enter API Email: ")
		pwd := t.prompter.PasswordPrompt("Enter API Password: ")
		user, err := models.NewUser(email, pwd)
		if err != nil {
			fmt.Println("Error creating API user: ", err)
			continue
		}
		if err = store.SaveUser(&user); err != nil {
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

func (f fileAPIInitializer) Initialize(store *store.Store) (models.User, error) {
	if user, err := store.FindUser(); err == nil {
		return user, err
	}

	request, err := credentialsFromFile(f.file)
	if err != nil {
		return models.User{}, err
	}

	user, err := models.NewUser(request.Email, request.Password)
	if err != nil {
		return user, err
	}
	return user, store.SaveUser(&user)
}

var errNoCredentialFile = errors.New("no API user credential file was passed")

func credentialsFromFile(file string) (models.SessionRequest, error) {
	if len(file) == 0 {
		return models.SessionRequest{}, errNoCredentialFile
	}

	logger.Debug("Initializing API credentials from ", file)
	dat, err := ioutil.ReadFile(file)
	if err != nil {
		return models.SessionRequest{}, err
	}
	lines := strings.Split(string(dat), "\n")
	if len(lines) < 2 {
		return models.SessionRequest{}, fmt.Errorf("malformed API credentials file does not have at least two lines at %s", file)
	}
	credentials := models.SessionRequest{
		Email:    strings.TrimSpace(lines[0]),
		Password: strings.TrimSpace(lines[1]),
	}
	return credentials, nil
}

// ChangePasswordPrompter is an interface primarily used for DI to obtain a
// password change request from the User.
type ChangePasswordPrompter interface {
	Prompt() (models.ChangePasswordRequest, error)
}

// NewChangePasswordPrompter returns the production password change request prompter
func NewChangePasswordPrompter() ChangePasswordPrompter {
	prompter := NewTerminalPrompter()
	return changePasswordPrompter{prompter: prompter}
}

type changePasswordPrompter struct {
	prompter Prompter
}

func (c changePasswordPrompter) Prompt() (models.ChangePasswordRequest, error) {
	fmt.Println("Changing your chainlink account password.")
	fmt.Println("NOTE: This will terminate any other sessions.")
	oldPassword := c.prompter.PasswordPrompt("Password:")

	fmt.Println("Now enter your **NEW** password")
	newPassword := c.prompter.PasswordPrompt("Password:")
	confirmPassword := c.prompter.PasswordPrompt("Confirmation:")

	if newPassword != confirmPassword {
		return models.ChangePasswordRequest{}, errors.New("new password and confirmation did not match")
	}

	return models.ChangePasswordRequest{
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

func confirmAction() bool {
	prompt := NewTerminalPrompter()
	var answer string
	for {
		answer = prompt.Prompt("Are you sure? This action is irreversible! (yes/no)")
		if answer == "yes" {
			return true
		} else if answer == "no" {
			return false
		} else {
			fmt.Printf("%s is not valid. Please type yes or no\n", answer)
		}
	}
}
