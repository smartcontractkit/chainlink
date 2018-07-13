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

	"github.com/manyminds/api2go/jsonapi"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/presenters"
	"github.com/smartcontractkit/chainlink/utils"
	"github.com/smartcontractkit/chainlink/web"
	"github.com/tidwall/gjson"
	clipkg "github.com/urfave/cli"
	"go.uber.org/multierr"
)

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
	cfg := cli.Config
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

type RemoteClient interface {
	Get(string) (*http.Response, error)
	Post(string) (*http.Response, error)
}

type httpPrompterClient struct {
	config   store.Config
	client   *http.Client
	prompter Prompter
}

func NewHttpPrompterClient(cfg store.Config, prompter Prompter) RemoteClient {
	return &httpPrompterClient{
		config:   cfg,
		client:   &http.Client{},
		prompter: prompter,
	}
}

func (h *httpPrompterClient) Get(path string) (*http.Response, error) {
	cookie, err := h.login()
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("GET", h.config.ClientNodeURL+path, nil)
	if err != nil {
		return nil, err
	}

	request.AddCookie(cookie)
	return h.client.Do(request)
}

func (h *httpPrompterClient) Post(path string) (*http.Response, error) {
	return nil, errors.New("bogus")
}

func (h *httpPrompterClient) login() (*http.Cookie, error) {
	url := h.config.ClientNodeURL + "/sessions"
	email := h.prompter.Prompt("Enter email: ")
	pwd := h.prompter.PasswordPrompt("Enter password: ")
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

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	_, err = parseResponse(resp)
	if err != nil {
		return nil, err
	}

	cookies := resp.Cookies()
	fmt.Println(cookies)
	fmt.Println("response", resp)
	if len(cookies) < 1 {
		return nil, errors.New("Did not receive cookie with session id")
	}
	return cookies[0], nil
}
