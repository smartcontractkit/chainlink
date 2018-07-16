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
	"path"
	"strconv"

	"github.com/manyminds/api2go/jsonapi"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/presenters"
	"github.com/smartcontractkit/chainlink/web"
	"github.com/tidwall/gjson"
	clipkg "github.com/urfave/cli"
	"go.uber.org/multierr"
)

var errUnauthorized = errors.New("401 Unauthorized")

// DisplayAccountBalance renders a table containing the active account address
// with it's ETH & LINK balance
func (cli *Client) DisplayAccountBalance(c *clipkg.Context) error {
	resp, err := cli.RemoteClient.Get("/v2/account_balance")
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
	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass the job id to be shown"))
	}
	resp, err := cli.RemoteClient.Get("/v2/specs/" + c.Args().First())
	if err != nil {
		return cli.errorOut(err)
	}
	defer resp.Body.Close()
	var job presenters.JobSpec
	return cli.renderAPIResponse(resp, &job)
}

// GetJobSpecs returns all job specs.
func (cli *Client) GetJobSpecs(c *clipkg.Context) error {
	page := 0
	if c != nil && c.IsSet("page") {
		page = c.Int("page")
	}

	var links jsonapi.Links
	var jobs []models.JobSpec
	err := cli.getPage("/v2/specs", page, &jobs, &links)
	if err != nil {
		return err
	}
	return cli.errorOut(cli.Render(&jobs))
}

// CreateJobSpec creates job spec based on JSON input
func (cli *Client) CreateJobSpec(c *clipkg.Context) error {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass in JSON or filepath"))
	}

	buf, err := getBufferFromJSON(c.Args().First())
	if err != nil {
		return cli.errorOut(err)
	}

	resp, err := cli.RemoteClient.Post("/v2/specs", buf)
	if err != nil {
		return cli.errorOut(err)
	}
	defer resp.Body.Close()

	var jobs presenters.JobSpec
	return cli.renderResponse(resp, &jobs)
}

// CreateJobRun creates job run based on SpecID and optional JSON
func (cli *Client) CreateJobRun(c *clipkg.Context) error {
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

	resp, err := cli.RemoteClient.Post("/v2/specs/"+c.Args().First()+"/runs", buf)
	if err != nil {
		return cli.errorOut(err)
	}
	defer resp.Body.Close()
	var jobs presenters.JobSpec
	return cli.renderResponse(resp, &jobs)
}

// BackupDatabase streams a backup of the node's db to the passed filepath.
func (cli *Client) BackupDatabase(c *clipkg.Context) error {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass the path to save the backup"))
	}
	resp, err := cli.RemoteClient.Get("/v2/backup")
	if err != nil {
		return cli.errorOut(err)
	}
	defer resp.Body.Close()
	return cli.errorOut(saveBodyAsFile(resp, c.Args().First()))
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

// AddBridge adds a new bridge to the chainlink node
func (cli *Client) AddBridge(c *clipkg.Context) error {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass in the bridge's parameters [JSON blob | JSON filepath]"))
	}

	buf, err := getBufferFromJSON(c.Args().First())
	if err != nil {
		return cli.errorOut(err)
	}

	resp, err := cli.RemoteClient.Post("/v2/bridge_types", buf)
	if err != nil {
		return cli.errorOut(err)
	}
	defer resp.Body.Close()

	var bridge models.BridgeType
	return cli.deserializeResponse(resp, &bridge)
}

// GetBridges returns all bridges.
func (cli *Client) GetBridges(c *clipkg.Context) error {
	page := 0
	if c != nil && c.IsSet("page") {
		page = c.Int("page")
	}

	var links jsonapi.Links
	var bridges []models.BridgeType
	err := cli.getPage("/v2/bridge_types", page, &bridges, &links)
	if err != nil {
		return err
	}
	return cli.errorOut(cli.Render(&bridges))
}

func (cli *Client) getPage(requestURI string, page int, model interface{}, links *jsonapi.Links) error {
	uri, err := url.Parse(requestURI)
	if err != nil {
		return err
	}
	q := uri.Query()
	if page > 0 {
		q.Set("page", strconv.Itoa(page))
	}
	uri.RawQuery = q.Encode()

	resp, err := cli.RemoteClient.Get(uri.String())
	if err != nil {
		return cli.errorOut(err)
	}
	defer resp.Body.Close()

	return cli.deserializeAPIResponse(resp, model, links)
}

// ShowBridge returns the info for the given Bridge name.
func (cli *Client) ShowBridge(c *clipkg.Context) error {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass the name of the bridge to be shown"))
	}
	resp, err := cli.RemoteClient.Get("/v2/bridge_types/" + c.Args().First())
	if err != nil {
		return cli.errorOut(err)
	}
	defer resp.Body.Close()
	var bridge models.BridgeType
	return cli.renderResponse(resp, &bridge)
}

// RemoveBridge removes a specific Bridge by name.
func (cli *Client) RemoveBridge(c *clipkg.Context) error {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass the name of the bridge to be removed"))
	}
	resp, err := cli.RemoteClient.Delete("/v2/bridge_types/" + c.Args().First())
	if err != nil {
		return cli.errorOut(err)
	}
	defer resp.Body.Close()
	var bridge models.BridgeType
	return cli.renderResponse(resp, &bridge)
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
	b, err := cli.parseResponse(resp)
	if err != nil {
		return err
	}
	if err = web.ParsePaginatedResponse(b, dst, links); err != nil {
		return cli.errorOut(err)
	}
	return nil
}

func (cli *Client) deserializeResponse(resp *http.Response, dst interface{}) error {
	b, err := cli.parseResponse(resp)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(b, &dst); err != nil {
		return cli.errorOut(err)
	}
	return nil
}

func (cli *Client) parseResponse(resp *http.Response) ([]byte, error) {
	b, err := parseResponse(resp)
	if err == errUnauthorized {
		return nil, cli.errorOut(multierr.Append(err, fmt.Errorf("Retry after deleting your cookie at %s/cookie", cli.Config.RootDir)))
	}
	if err != nil {
		return nil, cli.errorOut(err)
	}
	return b, err
}

func parseResponse(resp *http.Response) ([]byte, error) {
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return b, multierr.Append(errors.New(resp.Status), err)
	}
	if resp.StatusCode == 401 {
		return b, errUnauthorized
	} else if resp.StatusCode >= 400 {
		return b, errors.New(resp.Status)
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

type RemoteClient interface {
	Get(string) (*http.Response, error)
	Post(string, io.Reader) (*http.Response, error)
	Delete(string) (*http.Response, error)
}

type authenticatedHTTPClient struct {
	config     store.Config
	client     *http.Client
	cookieAuth CookieAuthenticator
}

func NewAuthenticatedHttpClient(cfg store.Config, cookieAuth CookieAuthenticator) RemoteClient {
	return &authenticatedHTTPClient{
		config:     cfg,
		client:     &http.Client{},
		cookieAuth: cookieAuth,
	}
}

func (h *authenticatedHTTPClient) Get(path string) (*http.Response, error) {
	return h.doRequest("GET", path, nil)
}

func (h *authenticatedHTTPClient) Post(path string, body io.Reader) (*http.Response, error) {
	return h.doRequest("POST", path, body)
}

func (h *authenticatedHTTPClient) Delete(path string) (*http.Response, error) {
	return h.doRequest("DELETE", path, nil)
}

func (h *authenticatedHTTPClient) doRequest(verb, path string, body io.Reader) (*http.Response, error) {
	cookie, err := h.cookieAuth.Authenticate()
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(verb, h.config.ClientNodeURL+path, body)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.AddCookie(cookie)
	return h.client.Do(request)
}

type CookieAuthenticator interface {
	Authenticate() (*http.Cookie, error)
}

type TerminalCookieAuthenticator struct {
	config   store.Config
	prompter Prompter
}

func NewTerminalCookieAuthenticator(cfg store.Config, prompter Prompter) CookieAuthenticator {
	return &TerminalCookieAuthenticator{config: cfg, prompter: prompter}
}

func (h *TerminalCookieAuthenticator) Authenticate() (*http.Cookie, error) {
	cookie, err := h.retrieveCookieFromDisk()
	if err == nil {
		return cookie, nil
	}

	url := h.config.ClientNodeURL + "/sessions"
	email := h.prompter.Prompt("Enter email: ")
	pwd := h.prompter.PasswordPrompt("Enter password: ")
	sessionRequest := models.SessionRequest{Email: email, Password: pwd}
	b := new(bytes.Buffer)
	err = json.NewEncoder(b).Encode(sessionRequest)
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
	return cookies[0], h.saveCookieToDisk(cookies[0])
}

func (h *TerminalCookieAuthenticator) saveCookieToDisk(cookie *http.Cookie) error {
	return ioutil.WriteFile(h.cookiePath(), []byte(cookie.String()), 0660)
}

func (h *TerminalCookieAuthenticator) retrieveCookieFromDisk() (*http.Cookie, error) {
	b, err := ioutil.ReadFile(h.cookiePath())
	if err != nil {
		return nil, err
	}
	header := http.Header{}
	header.Add("Cookie", string(b))
	request := http.Request{Header: header}
	cookies := request.Cookies()
	if len(cookies) == 0 {
		return nil, errors.New("Cookie not in file")
	}
	return request.Cookies()[0], nil
}

func (h *TerminalCookieAuthenticator) cookiePath() string {
	return path.Join(h.config.RootDir, "cookie")
}
