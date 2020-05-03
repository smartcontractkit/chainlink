package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/presenters"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/web"

	"github.com/ethereum/go-ethereum/common"
	"github.com/manyminds/api2go/jsonapi"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
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
		return cli.errorOut(errors.Wrap(err, "while extracting json to buffer"))
	}

	resp, err := cli.HTTP.Post("/v2/service_agreements", buf)
	if err != nil {
		return cli.errorOut(errors.Wrap(err, "from initializing service-agreement-creation request"))
	}
	defer resp.Body.Close()

	var sa presenters.ServiceAgreement
	return cli.renderAPIResponse(resp, &sa)
}

// CreateExternalInitiator adds an external initiator
func (cli *Client) CreateExternalInitiator(c *clipkg.Context) error {
	if c.NArg() != 1 && c.NArg() != 2 {
		return cli.errorOut(errors.New("create expects 1 - 2 arguments: a name and a url (optional)"))
	}

	var request models.ExternalInitiatorRequest
	request.Name = c.Args().Get(0)

	// process optional URL
	if c.NArg() == 2 {
		url, err := url.ParseRequestURI(c.Args().Get(1))
		if err != nil {
			return cli.errorOut(err)
		}
		request.URL = (*models.WebURL)(url)
	}

	requestData, err := json.Marshal(request)
	if err != nil {
		return cli.errorOut(err)
	}

	buf := bytes.NewBuffer(requestData)
	resp, err := cli.HTTP.Post("/v2/external_initiators", buf)
	if err != nil {
		return cli.errorOut(err)
	}
	defer resp.Body.Close()

	var ei presenters.ExternalInitiatorAuthentication
	return cli.renderAPIResponse(resp, &ei)
}

// DeleteExternalInitiator removes an external initiator
func (cli *Client) DeleteExternalInitiator(c *clipkg.Context) error {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass the name of the external initiator to delete"))
	}

	resp, err := cli.HTTP.Delete("/v2/external_initiators/" + c.Args().First())
	if err != nil {
		return cli.errorOut(err)
	}
	defer resp.Body.Close()
	_, err = cli.parseResponse(resp)
	return err
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

// IndexJobRuns returns the list of all job runs for a specific job
// if no jobid is passed, defaults to returning all jobruns
func (cli *Client) IndexJobRuns(c *clipkg.Context) error {
	jobID := c.String("jobid")
	if jobID != "" {
		return cli.getPage("/v2/runs?jobSpecId="+jobID, c.Int("page"), &[]presenters.JobRun{})
	}
	return cli.getPage("/v2/runs", c.Int("page"), &[]presenters.JobRun{})
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

// IndexJobSpecs returns all job specs.
func (cli *Client) IndexJobSpecs(c *clipkg.Context) error {
	return cli.getPage("/v2/specs", c.Int("page"), &[]models.JobSpec{})
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

// ArchiveJobSpec soft deletes a job and its associated runs.
func (cli *Client) ArchiveJobSpec(c *clipkg.Context) error {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass the job id to be archived"))
	}
	resp, err := cli.HTTP.Delete("/v2/specs/" + c.Args().First())
	if err != nil {
		return cli.errorOut(err)
	}
	_, err = cli.parseResponse(resp)
	if err != nil {
		return cli.errorOut(err)
	}
	return nil
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

// CreateBridge adds a new bridge to the chainlink node
func (cli *Client) CreateBridge(c *clipkg.Context) error {
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

	var bridge models.BridgeTypeAuthentication
	return cli.renderAPIResponse(resp, &bridge)
}

// IndexBridges returns all bridges.
func (cli *Client) IndexBridges(c *clipkg.Context) error {
	return cli.getPage("/v2/bridge_types", c.Int("page"), &[]models.BridgeType{})
}

func (cli *Client) getPage(requestURI string, page int, model interface{}) error {
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

	err = cli.deserializeAPIResponse(resp, model, &jsonapi.Links{})
	if err != nil {
		return err
	}
	return cli.errorOut(cli.Render(model))
}

// ShowBridge returns the info for the given Bridge name.
func (cli *Client) ShowBridge(c *clipkg.Context) error {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass the name of the bridge to be shown"))
	}
	bridgeName := c.Args().First()
	resp, err := cli.HTTP.Get("/v2/bridge_types/" + bridgeName)
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
	bridgeName := c.Args().First()
	resp, err := cli.HTTP.Delete("/v2/bridge_types/" + bridgeName)
	if err != nil {
		return cli.errorOut(err)
	}
	defer resp.Body.Close()
	var bridge models.BridgeType
	return cli.renderAPIResponse(resp, &bridge)
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

// IndexTransactions returns the list of transactions in descending order,
// taking an optional page parameter
func (cli *Client) IndexTransactions(c *clipkg.Context) error {
	return cli.getPage("/v2/transactions", c.Int("page"), &[]presenters.Tx{})
}

// ShowTransaction returns the info for the given transaction hash
func (cli *Client) ShowTransaction(c *clipkg.Context) error {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass the hash of the transaction"))
	}
	hash := c.Args().First()
	resp, err := cli.HTTP.Get("/v2/transactions/" + hash)
	if err != nil {
		return cli.errorOut(err)
	}
	defer resp.Body.Close()
	var tx presenters.Tx
	return cli.renderAPIResponse(resp, &tx)
}

// IndexTxAttempts returns the list of transactions in descending order,
// taking an optional page parameter
func (cli *Client) IndexTxAttempts(c *clipkg.Context) error {
	return cli.getPage("/v2/tx_attempts", c.Int("page"), &[]models.TxAttempt{})
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
		return nil, fmt.Errorf("invalid JSON or file not found '%s'", s)
	} else if err != nil {
		return nil, fmt.Errorf("error reading from file '%s': %v", s, err)
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
		return errors.Wrap(err, "parseResponse error")
	}
	if err = web.ParsePaginatedResponse(b, dst, links); err != nil {
		return cli.errorOut(err)
	}
	return nil
}

func (cli *Client) parseResponse(resp *http.Response) ([]byte, error) {
	b, err := parseResponse(resp)
	if err == errUnauthorized {
		return nil, cli.errorOut(multierr.Append(err, fmt.Errorf("try logging in")))
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
	if resp.StatusCode == http.StatusUnauthorized {
		return b, errUnauthorized
	} else if resp.StatusCode >= http.StatusBadRequest {
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
		CurrentPassword: password,
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

// SetMinimumGasPrice specifies the minimum gas price to use for outgoing transactions
func (cli *Client) SetMinimumGasPrice(c *clipkg.Context) error {
	if c.NArg() != 1 {
		return cli.errorOut(errors.New("expecting an amount"))
	}

	value := c.Args().Get(0)
	amount, ok := new(big.Float).SetString(value)
	if !ok {
		return cli.errorOut(fmt.Errorf("invalid ethereum amount %s", value))
	}

	if c.IsSet("gwei") {
		amount.Mul(amount, big.NewFloat(1000000000))
	}

	adjustedAmount, _ := amount.Int(nil)
	request := struct {
		EthGasPriceDefault string `json:"ethGasPriceDefault"`
	}{EthGasPriceDefault: adjustedAmount.String()}
	requestData, err := json.Marshal(request)
	if err != nil {
		return cli.errorOut(err)
	}

	buf := bytes.NewBuffer(requestData)
	response, err := cli.HTTP.Patch("/v2/config", buf)
	if err != nil {
		return cli.errorOut(err)
	}
	defer response.Body.Close()

	patchResponse := web.ConfigPatchResponse{}
	if err := cli.deserializeAPIResponse(response, &patchResponse, &jsonapi.Links{}); err != nil {
		return err
	}

	return cli.errorOut(cli.Render(&patchResponse))
}

// GetConfiguration gets the nodes environment variables
func (cli *Client) GetConfiguration(c *clipkg.Context) error {
	resp, err := cli.HTTP.Get("/v2/config")
	if err != nil {
		return cli.errorOut(err)
	}
	defer resp.Body.Close()
	cwl := presenters.ConfigWhitelist{}
	return cli.renderAPIResponse(resp, &cwl)
}

// CancelJob cancels a running job
func (cli *Client) CancelJobRun(c *clipkg.Context) error {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass the run id to be cancelled"))
	}

	response, err := cli.HTTP.Put(fmt.Sprintf("/v2/runs/%s/cancellation", c.Args().First()), nil)
	if err != nil {
		return cli.errorOut(errors.Wrap(err, "HTTP.Put"))
	}
	_, err = cli.parseResponse(response)
	if err != nil {
		return cli.errorOut(errors.Wrap(err, "cli.parseResponse"))
	}
	return nil
}
