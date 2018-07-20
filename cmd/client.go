package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/logger"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/web"
	clipkg "github.com/urfave/cli"
	"go.uber.org/multierr"
	"golang.org/x/sync/errgroup"
)

// Client is the shell for the node, local commands and remote commands.
type Client struct {
	Renderer
	Config                         store.Config
	AppFactory                     AppFactory
	KeyStoreAuthenticator          KeyStoreAuthenticator
	FallbackAPIInitializer         APIInitializer
	Runner                         Runner
	HTTP                           HTTPClient
	CookieAuthenticator            CookieAuthenticator
	FileSessionRequestBuilder      SessionRequestBuilder
	PromptingSessionRequestBuilder SessionRequestBuilder
}

func (cli *Client) errorOut(err error) error {
	if err != nil {
		return clipkg.NewExitError(err.Error(), 1)
	}
	return nil
}

// AppFactory implements the NewApplication method.
type AppFactory interface {
	NewApplication(store.Config) services.Application
}

// ChainlinkAppFactory is used to create a new Application.
type ChainlinkAppFactory struct{}

// NewApplication returns a new instance of the node with the given config.
func (n ChainlinkAppFactory) NewApplication(config store.Config) services.Application {
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

// HTTPClient encapsulates all methods used to interact with a chainlink node API.
type HTTPClient interface {
	Get(string, ...map[string]string) (*http.Response, error)
	Post(string, io.Reader) (*http.Response, error)
	Patch(string, io.Reader) (*http.Response, error)
	Delete(string) (*http.Response, error)
}

type authenticatedHTTPClient struct {
	config     store.Config
	client     *http.Client
	cookieAuth CookieAuthenticator
}

// NewAuthenticatedHTTPClient uses the CookieAuthenticator to generate a sessionID
// which is then used for all subsequent HTTP API requests.
func NewAuthenticatedHTTPClient(cfg store.Config, cookieAuth CookieAuthenticator) HTTPClient {
	return &authenticatedHTTPClient{
		config:     cfg,
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

// Patch performs an HTTP Patch using the authenticated HTTP client's cookie.
func (h *authenticatedHTTPClient) Patch(path string, body io.Reader) (*http.Response, error) {
	return h.doRequest("PATCH", path, body)
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

	request, err := http.NewRequest(verb, h.config.ClientNodeURL+path, body)
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
	config store.Config
}

// NewSessionCookieAuthenticator creates a SessionCookieAuthenticator using the passed config
// and builder.
func NewSessionCookieAuthenticator(cfg store.Config) CookieAuthenticator {
	return &SessionCookieAuthenticator{config: cfg}
}

// Cookie Returns the previously saved authentication cookie.
func (t *SessionCookieAuthenticator) Cookie() (*http.Cookie, error) {
	return t.retrieveCookieFromDisk()
}

// Authenticate retrieves a session ID via a cookie and saves it to disk.
func (t *SessionCookieAuthenticator) Authenticate(sessionRequest models.SessionRequest) (*http.Cookie, error) {
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(sessionRequest)
	if err != nil {
		return nil, err
	}
	url := t.config.ClientNodeURL + "/sessions"
	req, err := http.NewRequest("POST", url, b)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	_, err = parseResponse(resp)
	if err != nil {
		return nil, err
	}

	cookies := resp.Cookies()
	if len(cookies) == 0 {
		return nil, errors.New("Did not receive cookie with session id")
	}
	return cookies[0], t.saveCookieToDisk(cookies[0])
}

func (t *SessionCookieAuthenticator) saveCookieToDisk(cookie *http.Cookie) error {
	return ioutil.WriteFile(t.cookiePath(), []byte(cookie.String()), 0660)
}

func (t *SessionCookieAuthenticator) retrieveCookieFromDisk() (*http.Cookie, error) {
	b, err := ioutil.ReadFile(t.cookiePath())
	if err != nil {
		return nil, multierr.Append(errors.New("Unable to retrieve credentials, have you logged in?"), err)
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

func (t *SessionCookieAuthenticator) cookiePath() string {
	return path.Join(t.config.RootDir, "cookie")
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

	for {
		email := t.prompter.Prompt("Enter API Email: ")
		pwd := t.prompter.PasswordPrompt("Enter API Password: ")
		user, err := models.NewUser(email, pwd)
		if err != nil {
			fmt.Println("Error creating API user: ", err)
			continue
		}
		if err = store.Save(&user); err != nil {
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
	return user, store.Save(&user)
}

var errNoCredentialFile = errors.New("No API user credential file was passed")

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
		return models.SessionRequest{}, fmt.Errorf("Malformed API credentials file does not have at least two lines at %s", file)
	}
	credentials := models.SessionRequest{
		Email:    strings.TrimSpace(lines[0]),
		Password: strings.TrimSpace(lines[1]),
	}
	return credentials, nil
}
