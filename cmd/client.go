package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/web"
	clipkg "github.com/urfave/cli"
	"go.uber.org/multierr"
	"golang.org/x/sync/errgroup"
)

// Client is the shell for the node. It has fields for the Renderer,
// Config, AppFactory (the services application), KeyStoreAuthenticator, and Runner.
type Client struct {
	Renderer
	Config              store.Config
	AppFactory          AppFactory
	Auth                KeyStoreAuthenticator
	APIInitializer      APIInitializer
	Runner              Runner
	RemoteClient        RemoteClient
	CookieAuthenticator CookieAuthenticator
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

// RemoteClient encapsulates all methods used to interact with a chainlink node API.
type RemoteClient interface {
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
func NewAuthenticatedHTTPClient(cfg store.Config, cookieAuth CookieAuthenticator) RemoteClient {
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
	Authenticate() (*http.Cookie, error)
}

// TerminalCookieAuthenticator is a concrete implementation of CookieAuthenticator
// that prompts the user for credentials via terminal.
type TerminalCookieAuthenticator struct {
	config   store.Config
	prompter Prompter
}

// NewTerminalCookieAuthenticator creates a TerminalCookieAuthenticator using the passed config
// and prompter.
func NewTerminalCookieAuthenticator(cfg store.Config, prompter Prompter) CookieAuthenticator {
	return &TerminalCookieAuthenticator{config: cfg, prompter: prompter}
}

// Cookie Returns the previously saved authentication cookie.
func (t *TerminalCookieAuthenticator) Cookie() (*http.Cookie, error) {
	return t.retrieveCookieFromDisk()
}

// Authenticate prompts the user for credentials to generate a cookie and saves
// it do disk.
func (t *TerminalCookieAuthenticator) Authenticate() (*http.Cookie, error) {
	url := t.config.ClientNodeURL + "/sessions"
	email := t.prompter.Prompt("Enter email: ")
	pwd := t.prompter.PasswordPrompt("Enter password: ")
	sessionRequest := models.SessionRequest{Email: email, Password: pwd}
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(sessionRequest)
	if err != nil {
		return nil, err
	}
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

func (t *TerminalCookieAuthenticator) saveCookieToDisk(cookie *http.Cookie) error {
	return ioutil.WriteFile(t.cookiePath(), []byte(cookie.String()), 0660)
}

func (t *TerminalCookieAuthenticator) retrieveCookieFromDisk() (*http.Cookie, error) {
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

func (h *TerminalCookieAuthenticator) cookiePath() string {
	return path.Join(h.config.RootDir, "cookie")
}
