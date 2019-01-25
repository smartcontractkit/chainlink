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

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/store/assets"

	"github.com/manyminds/api2go/jsonapi"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/presenters"
	"github.com/smartcontractkit/chainlink/utils"
	"github.com/smartcontractkit/chainlink/web"
	"github.com/tidwall/gjson"
	clipkg "github.com/urfave/cli"
	"go.uber.org/multierr"
)

var errUnauthorized = errors.New("401 Unauthorized")

// DisplayAccountBalance renders a table containing the active account address
// with it's ETH & LINK balance
func (cli *Client) DisplayAccountBalance(c *clipkg.Context) error {
	resp, err := cli.HTTP.Get("/v2/user/balances")
	if err != nil {
		return cli.errorOut(err)
	}
	defer resp.Body.Close()

	var links jsonapi.Links
	balances := []presenters.AccountBalance{}
	if err = cli.deserializeAPIResponse(resp, &balances, &links); err != nil {
		return err
	}
	return cli.errorOut(cli.Render(&balances))
}

// CreateServiceAgreement creates a ServiceAgreement based on JSON input
func (cli *Client) CreateServiceAgreement(c *clipkg.Context) error {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass in JSON or filepath"))
	}

	buf, err := getBufferFromJSON(c.Args().First())
	if err != nil {
		return cli.errorOut(err)
	}

	resp, err := cli.HTTP.Post("/v2/service_agreements", buf)
	if err != nil {
		return cli.errorOut(err)
	}
	defer resp.Body.Close()

	var sa presenters.ServiceAgreement
	return cli.renderResponse(resp, &sa)
}

// ShowJobRun returns the status of the given Jobrun.
func (cli *Client) ShowJobRun(c *clipkg.Context) error {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass the RunID to show"))
	}
	resp, err := cli.HTTP.Get("/v2/runs/" + c.Args().First())
	if err != nil {
		return cli.errorOut(err)
	}
	defer resp.Body.Close()
	var job presenters.JobRun
	return cli.renderAPIResponse(resp, &job)
}

// ShowJobSpec returns the status of the given JobID.
func (cli *Client) ShowJobSpec(c *clipkg.Context) error {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass the job id to be shown"))
	}
	resp, err := cli.HTTP.Get("/v2/specs/" + c.Args().First())
	if err != nil {
		return cli.errorOut(err)
	}
	defer resp.Body.Close()
	var job presenters.JobSpec
	return cli.renderAPIResponse(resp, &job)
}

// GetJobSpecs returns all job specs.
func (cli *Client) GetJobSpecs(c *clipkg.Context) error {
	var links jsonapi.Links
	var jobs []models.JobSpec
	err := cli.getPage("/v2/specs", c.Int("page"), &jobs, &links)
	if err != nil {
		return err
	}
	return cli.errorOut(cli.Render(&jobs))
}

// CreateJobSpec creates a JobSpec based on JSON input
func (cli *Client) CreateJobSpec(c *clipkg.Context) error {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass in JSON or filepath"))
	}

	buf, err := getBufferFromJSON(c.Args().First())
	if err != nil {
		return cli.errorOut(err)
	}

	resp, err := cli.HTTP.Post("/v2/specs", buf)
	if err != nil {
		return cli.errorOut(err)
	}
	defer resp.Body.Close()

	var js presenters.JobSpec
	return cli.renderAPIResponse(resp, &js)
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

	resp, err := cli.HTTP.Post("/v2/specs/"+c.Args().First()+"/runs", buf)
	if err != nil {
		return cli.errorOut(err)
	}
	defer resp.Body.Close()
	var run presenters.JobRun
	return cli.renderAPIResponse(resp, &run)
}

// BackupDatabase streams a backup of the node's db to the passed filepath.
func (cli *Client) BackupDatabase(c *clipkg.Context) error {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass the path to save the backup"))
	}
	resp, err := cli.HTTP.Get("/v2/backup")
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

	resp, err := cli.HTTP.Post("/v2/bridge_types", buf)
	if err != nil {
		return cli.errorOut(err)
	}
	defer resp.Body.Close()

	var bridge models.BridgeType
	return cli.renderAPIResponse(resp, &bridge)
}

// GetBridges returns all bridges.
func (cli *Client) GetBridges(c *clipkg.Context) error {
	var links jsonapi.Links
	var bridges []models.BridgeType
	err := cli.getPage("/v2/bridge_types", c.Int("page"), &bridges, &links)
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

	resp, err := cli.HTTP.Get(uri.String())
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
	resp, err := cli.HTTP.Get("/v2/bridge_types/" + c.Args().First())
	if err != nil {
		return cli.errorOut(err)
	}
	defer resp.Body.Close()
	var bridge models.BridgeType
	return cli.renderAPIResponse(resp, &bridge)
}

// RemoveBridge removes a specific Bridge by name.
func (cli *Client) RemoveBridge(c *clipkg.Context) error {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass the name of the bridge to be removed"))
	}
	resp, err := cli.HTTP.Delete("/v2/bridge_types/" + c.Args().First())
	if err != nil {
		return cli.errorOut(err)
	}
	defer resp.Body.Close()
	var bridge models.BridgeType
	return cli.renderResponse(resp, &bridge)
}

// RemoteLogin creates a cookie session to run remote commands.
func (cli *Client) RemoteLogin(c *clipkg.Context) error {
	sessionRequest, err := cli.buildSessionRequest(c.String("file"))
	if err != nil {
		return cli.errorOut(err)
	}
	_, err = cli.CookieAuthenticator.Authenticate(sessionRequest)
	return cli.errorOut(err)
}

// Withdraw will withdraw LINK to an address authorized by the node
func (cli *Client) Withdraw(c *clipkg.Context) error {
	if c.NArg() != 2 {
		return cli.errorOut(errors.New("withdraw expects two arguments: an address and an amount"))
	}

	linkAmount, err := strconv.ParseInt(c.Args().Get(1), 10, 64)

	if err != nil {
		return cli.errorOut(multierr.Combine(
			errors.New("while parsing LINK withdrawal amount"), err))
	}

	contractAddress := common.Address{}
	unParsedOracleContractAddress := c.String("from-oracle-contract-address")
	if unParsedOracleContractAddress != "" {
		contractAddress, err = utils.ParseEthereumAddress(
			unParsedOracleContractAddress)
		if err != nil {
			return cli.errorOut(multierr.Combine(
				errors.New("while parsing source contract withdrawal address"),
				err,
			))
		}
	}
	unparsedDestinationAddress := c.Args().First()
	destinationAddress, err := utils.ParseEthereumAddress(unparsedDestinationAddress)
	if err != nil {
		return cli.errorOut(multierr.Combine(
			fmt.Errorf("while parsing withdrawal destination address %v",
				unparsedDestinationAddress), err))
	}

	wR := models.WithdrawalRequest{
		DestinationAddress: destinationAddress,
		ContractAddress:    contractAddress,
		Amount:             assets.NewLink(linkAmount),
	}

	requestData, err := json.Marshal(wR)
	if err != nil {
		return cli.errorOut(err)
	}

	buf := bytes.NewBuffer(requestData)

	resp, err := cli.HTTP.Post("/v2/withdrawals", buf)
	if err != nil {
		return cli.errorOut(err)
	}
	defer resp.Body.Close()

	return cli.printResponseBody(resp)
}

// SendEther transfers ETH from the node's account to a specified address.
func (cli *Client) SendEther(c *clipkg.Context) error {
	if c.NArg() != 2 {
		return cli.errorOut(errors.New("sendether expects two arguments: an amount and an address"))
	}

	amount, err := strconv.ParseInt(c.Args().Get(0), 10, 64)
	if err != nil {
		return cli.errorOut(multierr.Combine(
			errors.New("while parsing ETH transfer amount"), err))
	}

	unparsedDestinationAddress := c.Args().Get(1)
	destinationAddress, err := utils.ParseEthereumAddress(unparsedDestinationAddress)
	if err != nil {
		return cli.errorOut(multierr.Combine(
			fmt.Errorf("while parsing withdrawal destination address %v",
				unparsedDestinationAddress), err))
	}

	unparsedFromAddress := c.String("from")
	fromAddress := common.Address{}
	if unparsedFromAddress != "" {
		fromAddress, err = utils.ParseEthereumAddress(unparsedFromAddress)
		if err != nil {
			return cli.errorOut(multierr.Combine(
				fmt.Errorf("while parsing withdrawal from address %v",
					unparsedFromAddress), err))
		}
	}

	request := models.SendEtherRequest{
		DestinationAddress: destinationAddress,
		FromAddress:        fromAddress,
		Amount:             assets.NewEth(amount),
	}

	requestData, err := json.Marshal(request)
	if err != nil {
		return cli.errorOut(err)
	}

	buf := bytes.NewBuffer(requestData)

	resp, err := cli.HTTP.Post("/v2/transfers", buf)
	if err != nil {
		return cli.errorOut(err)
	}
	defer resp.Body.Close()

	return cli.printResponseBody(resp)
}

// ChangePassword prompts the user for the old password and a new one, then
// posts it to Chainlink to change the password.
func (cli *Client) ChangePassword(c *clipkg.Context) error {
	req, err := cli.ChangePasswordPrompter.Prompt()
	if err != nil {
		return cli.errorOut(err)
	}

	requestData, err := json.Marshal(req)
	if err != nil {
		return cli.errorOut(err)
	}

	buf := bytes.NewBuffer(requestData)
	resp, err := cli.HTTP.Patch("/v2/user/password", buf)
	if err != nil {
		return cli.errorOut(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Password updated.")
	} else if resp.StatusCode == http.StatusConflict {
		fmt.Println("Old password did not match.")
	} else {
		return cli.printResponseBody(resp)
	}
	return nil
}

// GetTxAttempts returns the list of transaction attempts in descending order,
// taking an optional page parameter
func (cli *Client) GetTxAttempts(c *clipkg.Context) error {
	var links jsonapi.Links
	attempts := []models.TxAttempt{}
	err := cli.getPage("/v2/txattempts", c.Int("page"), &attempts, &links)
	if err != nil {
		return err
	}
	return cli.errorOut(cli.Render(&attempts))
}

func (cli *Client) buildSessionRequest(flag string) (models.SessionRequest, error) {
	if len(flag) > 0 {
		return cli.FileSessionRequestBuilder.Build(flag)
	}
	return cli.PromptingSessionRequestBuilder.Build("")
}

func getBufferFromJSON(s string) (*bytes.Buffer, error) {
	if gjson.Valid(s) {
		return bytes.NewBufferString(s), nil
	}

	buf, err := fromFile(s)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("Invalid JSON or file not found '%s'", s)
	} else if err != nil {
		return nil, fmt.Errorf("Error reading from file '%s': %v", s, err)
	}
	return buf, nil
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
		return nil, cli.errorOut(multierr.Append(err, fmt.Errorf("Try logging in")))
	}
	if err != nil {
		jae := models.JSONAPIErrors{}
		unmarshalErr := json.Unmarshal(b, &jae)
		return nil, cli.errorOut(multierr.Combine(err, unmarshalErr, &jae))
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

func (cli *Client) printResponseBody(resp *http.Response) error {
	b, err := parseResponse(resp)
	if err != nil {
		return cli.errorOut(err)
	}

	fmt.Println(string(b))
	return nil
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

// CreateExtraKey creates a new ethereum key with the same password
// as the one used to unlock the existing key.
func (cli *Client) CreateExtraKey(c *clipkg.Context) error {
	password := cli.PasswordPrompter.Prompt()
	request := models.CreateKeyRequest{
		CurrentPassword:    password,
		NewAccountPassword: password,
	}

	requestData, err := json.Marshal(request)
	if err != nil {
		return cli.errorOut(err)
	}

	buf := bytes.NewBuffer(requestData)
	resp, err := cli.HTTP.Post("/v2/keys", buf)
	if err != nil {
		return cli.errorOut(err)
	}
	defer resp.Body.Close()

	return cli.printResponseBody(resp)
}
