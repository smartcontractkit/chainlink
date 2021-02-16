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
	"strings"

	"github.com/manyminds/api2go/jsonapi"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	clipkg "github.com/urfave/cli"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/models/ocrkey"
	"github.com/smartcontractkit/chainlink/core/store/models/p2pkey"
	"github.com/smartcontractkit/chainlink/core/store/presenters"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/web"
	webPresenter "github.com/smartcontractkit/chainlink/core/web/presenters"
)

var errUnauthorized = errors.New(http.StatusText(http.StatusUnauthorized))

// CreateServiceAgreement creates a ServiceAgreement based on JSON input
func (cli *Client) CreateServiceAgreement(c *clipkg.Context) (err error) {
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
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	var sa presenters.ServiceAgreement
	err = cli.renderAPIResponse(resp, &sa)
	return err
}

// CreateExternalInitiator adds an external initiator
func (cli *Client) CreateExternalInitiator(c *clipkg.Context) (err error) {
	if c.NArg() != 1 && c.NArg() != 2 {
		return cli.errorOut(errors.New("create expects 1 - 2 arguments: a name and a url (optional)"))
	}

	var request models.ExternalInitiatorRequest
	request.Name = c.Args().Get(0)

	// process optional URL
	if c.NArg() == 2 {
		var reqURL *url.URL
		reqURL, err = url.ParseRequestURI(c.Args().Get(1))
		if err != nil {
			return cli.errorOut(err)
		}
		request.URL = (*models.WebURL)(reqURL)
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
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	var ei presenters.ExternalInitiatorAuthentication
	err = cli.renderAPIResponse(resp, &ei)
	return err
}

// DeleteExternalInitiator removes an external initiator
func (cli *Client) DeleteExternalInitiator(c *clipkg.Context) (err error) {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass the name of the external initiator to delete"))
	}

	resp, err := cli.HTTP.Delete("/v2/external_initiators/" + c.Args().First())
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()
	_, err = cli.parseResponse(resp)
	return err
}

// ShowJobRun returns the status of the given Jobrun.
func (cli *Client) ShowJobRun(c *clipkg.Context) (err error) {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass the RunID to show"))
	}
	resp, err := cli.HTTP.Get("/v2/runs/" + c.Args().First())
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()
	var job presenters.JobRun
	err = cli.renderAPIResponse(resp, &job)
	return err
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
func (cli *Client) ShowJobSpec(c *clipkg.Context) (err error) {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass the job id to be shown"))
	}
	resp, err := cli.HTTP.Get("/v2/specs/" + c.Args().First())
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()
	var job presenters.JobSpec
	err = cli.renderAPIResponse(resp, &job)
	return err
}

// IndexJobSpecs returns all job specs.
func (cli *Client) IndexJobSpecs(c *clipkg.Context) error {
	return cli.getPage("/v2/specs", c.Int("page"), &[]models.JobSpec{})
}

// CreateJobSpec creates a JobSpec based on JSON input
func (cli *Client) CreateJobSpec(c *clipkg.Context) (err error) {
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
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	var js presenters.JobSpec
	err = cli.renderAPIResponse(resp, &js)
	return err
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

// ListJobsV2 lists all v2 jobs
func (cli *Client) ListJobsV2(c *clipkg.Context) (err error) {
	return cli.getPage("/v2/jobs", c.Int("page"), &[]Job{})
}

// CreateJobV2 creates a V2 job
// Valid input is a TOML string or a path to TOML file
func (cli *Client) CreateJobV2(c *clipkg.Context) (err error) {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass in TOML or filepath"))
	}

	tomlString, err := getTOMLString(c.Args().First())
	if err != nil {
		return cli.errorOut(err)
	}

	request, err := json.Marshal(models.CreateJobSpecRequest{
		TOML: tomlString,
	})
	if err != nil {
		return cli.errorOut(err)
	}

	resp, err := cli.HTTP.Post("/v2/jobs", bytes.NewReader(request))
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	if resp.StatusCode >= 400 {
		body, rerr := ioutil.ReadAll(resp.Body)
		if err != nil {
			err = multierr.Append(err, rerr)
			return cli.errorOut(err)
		}
		fmt.Printf("Error : %v\n", string(body))
		return cli.errorOut(err)
	}

	var js Job
	err = cli.renderAPIResponse(resp, &js, "Job created")
	return err
}

// DeleteJobV2 deletes a V2 job
func (cli *Client) DeleteJobV2(c *clipkg.Context) error {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass the job id to be archived"))
	}
	resp, err := cli.HTTP.Delete("/v2/jobs/" + c.Args().First())
	if err != nil {
		return cli.errorOut(err)
	}
	_, err = cli.parseResponse(resp)
	if err != nil {
		return cli.errorOut(err)
	}

	fmt.Printf("Job %v Deleted\n", c.Args().First())
	return nil
}

// TriggerPipelineRun triggers a V2 job run based on a job ID
func (cli *Client) TriggerPipelineRun(c *clipkg.Context) error {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass the job id to trigger a run"))
	}
	resp, err := cli.HTTP.Post("/v2/jobs/"+c.Args().First()+"/runs", nil)
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	var run pipeline.Run
	err = cli.renderAPIResponse(resp, &run, "Pipeline run successfully triggered")
	return err
}

// CreateJobRun creates job run based on SpecID and optional JSON
func (cli *Client) CreateJobRun(c *clipkg.Context) (err error) {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass in SpecID [JSON blob | JSON filepath]"))
	}

	buf := bytes.NewBufferString("")
	if c.NArg() > 1 {
		var jbuf *bytes.Buffer
		jbuf, err = getBufferFromJSON(c.Args().Get(1))
		if err != nil {
			return cli.errorOut(err)
		}
		buf = jbuf
	}

	resp, err := cli.HTTP.Post("/v2/specs/"+c.Args().First()+"/runs", buf)
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()
	var run presenters.JobRun
	err = cli.renderAPIResponse(resp, &run)
	return err
}

// CreateBridge adds a new bridge to the chainlink node
func (cli *Client) CreateBridge(c *clipkg.Context) (err error) {
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
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	var bridge models.BridgeTypeAuthentication
	err = cli.renderAPIResponse(resp, &bridge)
	return err
}

// IndexBridges returns all bridges.
func (cli *Client) IndexBridges(c *clipkg.Context) (err error) {
	return cli.getPage("/v2/bridge_types", c.Int("page"), &[]models.BridgeType{})
}

func (cli *Client) getPage(requestURI string, page int, model interface{}) (err error) {
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
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	err = cli.deserializeAPIResponse(resp, model, &jsonapi.Links{})
	if err != nil {
		return err
	}
	err = cli.errorOut(cli.Render(model))
	return err
}

// ShowBridge returns the info for the given Bridge name.
func (cli *Client) ShowBridge(c *clipkg.Context) (err error) {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass the name of the bridge to be shown"))
	}
	bridgeName := c.Args().First()
	resp, err := cli.HTTP.Get("/v2/bridge_types/" + bridgeName)
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()
	var bridge models.BridgeType
	return cli.renderAPIResponse(resp, &bridge)
}

// RemoveBridge removes a specific Bridge by name.
func (cli *Client) RemoveBridge(c *clipkg.Context) (err error) {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass the name of the bridge to be removed"))
	}
	bridgeName := c.Args().First()
	resp, err := cli.HTTP.Delete("/v2/bridge_types/" + bridgeName)
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()
	var bridge models.BridgeType
	err = cli.renderAPIResponse(resp, &bridge)
	return err
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

// SendEther transfers ETH from the node's account to a specified address.
func (cli *Client) SendEther(c *clipkg.Context) (err error) {
	if c.NArg() < 3 {
		return cli.errorOut(errors.New("sendether expects three arguments: amount, fromAddress and toAddress"))
	}

	amount, err := assets.NewEthValueS(c.Args().Get(0))
	if err != nil {
		return cli.errorOut(multierr.Combine(
			errors.New("while parsing ETH transfer amount"), err))
	}

	unparsedFromAddress := c.Args().Get(1)
	fromAddress, err := utils.ParseEthereumAddress(unparsedFromAddress)
	if err != nil {
		return cli.errorOut(multierr.Combine(
			fmt.Errorf("while parsing withdrawal source address %v",
				unparsedFromAddress), err))
	}

	unparsedDestinationAddress := c.Args().Get(2)
	destinationAddress, err := utils.ParseEthereumAddress(unparsedDestinationAddress)
	if err != nil {
		return cli.errorOut(multierr.Combine(
			fmt.Errorf("while parsing withdrawal destination address %v",
				unparsedDestinationAddress), err))
	}

	request := models.SendEtherRequest{
		DestinationAddress: destinationAddress,
		FromAddress:        fromAddress,
		Amount:             amount,
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
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	var tx presenters.EthTx
	err = cli.renderAPIResponse(resp, &tx)
	return err
}

// ChangePassword prompts the user for the old password and a new one, then
// posts it to Chainlink to change the password.
func (cli *Client) ChangePassword(c *clipkg.Context) (err error) {
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
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

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
	return cli.getPage("/v2/transactions", c.Int("page"), &[]presenters.EthTx{})
}

// ShowTransaction returns the info for the given transaction hash
func (cli *Client) ShowTransaction(c *clipkg.Context) (err error) {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass the hash of the transaction"))
	}
	hash := c.Args().First()
	resp, err := cli.HTTP.Get("/v2/transactions/" + hash)
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()
	var tx presenters.EthTx
	err = cli.renderAPIResponse(resp, &tx)
	return err
}

// IndexTxAttempts returns the list of transactions in descending order,
// taking an optional page parameter
func (cli *Client) IndexTxAttempts(c *clipkg.Context) error {
	return cli.getPage("/v2/tx_attempts", c.Int("page"), &[]presenters.EthTx{})
}

func (cli *Client) buildSessionRequest(flag string) (models.SessionRequest, error) {
	if len(flag) > 0 {
		return cli.FileSessionRequestBuilder.Build(flag)
	}
	return cli.PromptingSessionRequestBuilder.Build("")
}

func getTOMLString(s string) (string, error) {
	var val interface{}
	err := toml.Unmarshal([]byte(s), &val)
	if err == nil {
		return s, nil
	}

	buf, err := fromFile(s)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("invalid TOML or file not found '%s'", s)
	} else if err != nil {
		return "", fmt.Errorf("error reading from file '%s': %v", s, err)
	}
	return buf.String(), nil
}

func (cli *Client) parseResponse(resp *http.Response) ([]byte, error) {
	b, err := parseResponse(resp)
	if err == errUnauthorized {
		return nil, cli.errorOut(multierr.Append(err, fmt.Errorf("you must first login through the CLI")))
	}
	if err != nil {
		jae := models.JSONAPIErrors{}
		unmarshalErr := json.Unmarshal(b, &jae)
		return nil, cli.errorOut(multierr.Combine(err, unmarshalErr, &jae))
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

func (cli *Client) renderAPIResponse(resp *http.Response, dst interface{}, headers ...string) error {
	var links jsonapi.Links
	if err := cli.deserializeAPIResponse(resp, dst, &links); err != nil {
		return cli.errorOut(err)
	}

	return cli.errorOut(cli.Render(dst, headers...))
}

// SetMinimumGasPrice specifies the minimum gas price to use for outgoing transactions
func (cli *Client) SetMinimumGasPrice(c *clipkg.Context) (err error) {
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
	defer func() {
		if cerr := response.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	patchResponse := web.ConfigPatchResponse{}
	if err = cli.deserializeAPIResponse(response, &patchResponse, &jsonapi.Links{}); err != nil {
		return err
	}

	err = cli.errorOut(cli.Render(&patchResponse))
	return err
}

// GetConfiguration gets the nodes environment variables
func (cli *Client) GetConfiguration(c *clipkg.Context) (err error) {
	resp, err := cli.HTTP.Get("/v2/config")
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()
	cwl := presenters.ConfigPrinter{}
	err = cli.renderAPIResponse(resp, &cwl)
	return err
}

// CancelJobRun cancels a running job,
// Run ID must be passed
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

// CreateETHKey creates a new ethereum key with the same password
// as the one used to unlock the existing key.
func (cli *Client) CreateETHKey(c *clipkg.Context) (err error) {
	resp, err := cli.HTTP.Post("/v2/keys/eth", nil)
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	var keys webPresenter.ETHKeyResource
	return cli.renderAPIResponse(resp, &keys, "ETH key created.\n\n🔑 New key")
}

// ListETHKeys renders the active account address with its ETH & LINK balance
func (cli *Client) ListETHKeys(c *clipkg.Context) (err error) {
	resp, err := cli.HTTP.Get("/v2/keys/eth")
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	var keys []webPresenter.ETHKeyResource
	return cli.renderAPIResponse(resp, &keys, "🔑 ETH keys")
}

// DeleteETHKey deletes an Etherium key,
// address of key must be passed
func (cli *Client) DeleteETHKey(c *clipkg.Context) (err error) {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass the address of the key to be deleted"))
	}

	if c.Bool("hard") && !confirmAction(c) {
		return nil
	}

	var queryStr string
	var confirmationMsg string
	if c.Bool("hard") {
		queryStr = "?hard=true"
		confirmationMsg = "Deleted ETH key"
	} else {
		confirmationMsg = "Archived ETH key"
	}

	address := c.Args().Get(0)
	resp, err := cli.HTTP.Delete(fmt.Sprintf("/v2/keys/eth/%s%s", address, queryStr))
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	var key webPresenter.ETHKeyResource
	return cli.renderAPIResponse(resp, &key, fmt.Sprintf("🔑 %s", confirmationMsg))
}

// ImportETHKey imports an Etherium key,
// file path must be passed
func (cli *Client) ImportETHKey(c *clipkg.Context) (err error) {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass the filepath of the key to be imported"))
	}

	oldPasswordFile := c.String("oldpassword")
	if len(oldPasswordFile) == 0 {
		return cli.errorOut(errors.New("Must specify --oldpassword/-p flag"))
	}
	oldPassword, err := ioutil.ReadFile(oldPasswordFile)
	if err != nil {
		return cli.errorOut(errors.Wrap(err, "Could not read password file"))
	}

	filepath := c.Args().Get(0)
	keyJSON, err := ioutil.ReadFile(filepath)
	if err != nil {
		return cli.errorOut(err)
	}

	normalizedPassword := normalizePassword(string(oldPassword))
	resp, err := cli.HTTP.Post("/v2/keys/eth/import?oldpassword="+normalizedPassword, bytes.NewReader(keyJSON))
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	var key webPresenter.ETHKeyResource
	return cli.renderAPIResponse(resp, &key, "🔑 Imported ETH key")
}

// ExportETHKey exports an ETH key,
// address must be passed
func (cli *Client) ExportETHKey(c *clipkg.Context) (err error) {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass the address of the key to export"))
	}

	newPasswordFile := c.String("newpassword")
	if len(newPasswordFile) == 0 {
		return cli.errorOut(errors.New("Must specify --newpassword/-p flag"))
	}
	newPassword, err := ioutil.ReadFile(newPasswordFile)
	if err != nil {
		return cli.errorOut(errors.Wrap(err, "Could not read password file"))
	}

	filepath := c.String("output")
	if len(newPassword) == 0 {
		return cli.errorOut(errors.New("Must specify --output/-o flag"))
	}

	address := c.Args().Get(0)

	normalizedPassword := normalizePassword(string(newPassword))
	resp, err := cli.HTTP.Post("/v2/keys/eth/export/"+address+"?newpassword="+normalizedPassword, nil)
	if err != nil {
		return cli.errorOut(errors.Wrap(err, "Could not make HTTP request"))
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return cli.errorOut(errors.New("Error exporting"))
	}

	keyJSON, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return cli.errorOut(errors.Wrap(err, "Could not read response body"))
	}

	err = utils.WriteFileWithMaxPerms(filepath, keyJSON, 0600)
	if err != nil {
		return cli.errorOut(errors.Wrapf(err, "Could not write %v", filepath))
	}

	_, err = os.Stderr.WriteString("🔑 Exported ETH key " + address + " to " + filepath + "\n")
	if err != nil {
		return cli.errorOut(err)
	}

	return nil
}

// CreateP2PKey stores a P2P keypair
func (cli *Client) CreateP2PKey(c *clipkg.Context) (err error) {
	resp, err := cli.HTTP.Post("/v2/keys/p2p", nil)
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	var key p2pkey.EncryptedP2PKey
	return cli.renderAPIResponse(resp, &key, "Created P2P keypair")
}

// ListP2PKeys retrieves a list of all P2P keys
func (cli *Client) ListP2PKeys(c *clipkg.Context) (err error) {
	resp, err := cli.HTTP.Get("/v2/keys/p2p", nil)
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	var keys []p2pkey.EncryptedP2PKey
	return cli.renderAPIResponse(resp, &keys)
}

// DeleteP2PKey deletes a P2P key,
// key ID must be passed
func (cli *Client) DeleteP2PKey(c *clipkg.Context) (err error) {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass the key ID to be deleted"))
	}
	id, err := strconv.ParseUint(c.Args().Get(0), 10, 32)
	if err != nil {
		return cli.errorOut(err)
	}

	if !confirmAction(c) {
		return nil
	}

	var queryStr string
	if c.Bool("hard") {
		queryStr = "?hard=true"
	}

	resp, err := cli.HTTP.Delete(fmt.Sprintf("/v2/keys/p2p/%d%s", id, queryStr))
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	var key p2pkey.EncryptedP2PKey
	return cli.renderAPIResponse(resp, &key, "P2P key deleted")
}

// ImportP2PKey imports and stores a P2P key,
// path to key must be passed
func (cli *Client) ImportP2PKey(c *clipkg.Context) (err error) {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass the filepath of the key to be imported"))
	}

	oldPasswordFile := c.String("oldpassword")
	if len(oldPasswordFile) == 0 {
		return cli.errorOut(errors.New("Must specify --oldpassword/-p flag"))
	}
	oldPassword, err := ioutil.ReadFile(oldPasswordFile)
	if err != nil {
		return cli.errorOut(errors.Wrap(err, "Could not read password file"))
	}

	filepath := c.Args().Get(0)
	keyJSON, err := ioutil.ReadFile(filepath)
	if err != nil {
		return cli.errorOut(err)
	}

	normalizedPassword := normalizePassword(string(oldPassword))
	resp, err := cli.HTTP.Post("/v2/keys/p2p/import?oldpassword="+normalizedPassword, bytes.NewReader(keyJSON))
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	var key p2pkey.EncryptedP2PKey
	return cli.renderAPIResponse(resp, &key, "🔑 Imported P2P key")
}

// ExportP2PKey exports a P2P key,
// key ID must be passed
func (cli *Client) ExportP2PKey(c *clipkg.Context) (err error) {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass the ID of the key to export"))
	}

	newPasswordFile := c.String("newpassword")
	if len(newPasswordFile) == 0 {
		return cli.errorOut(errors.New("Must specify --newpassword/-p flag"))
	}
	newPassword, err := ioutil.ReadFile(newPasswordFile)
	if err != nil {
		return cli.errorOut(errors.Wrap(err, "Could not read password file"))
	}

	filepath := c.String("output")
	if len(filepath) == 0 {
		return cli.errorOut(errors.New("Must specify --output/-o flag"))
	}

	ID := c.Args().Get(0)

	normalizedPassword := normalizePassword(string(newPassword))
	resp, err := cli.HTTP.Post("/v2/keys/p2p/export/"+ID+"?newpassword="+normalizedPassword, nil)
	if err != nil {
		return cli.errorOut(errors.Wrap(err, "Could not make HTTP request"))
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return cli.errorOut(errors.New("Error exporting"))
	}

	keyJSON, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return cli.errorOut(errors.Wrap(err, "Could not read response body"))
	}

	err = utils.WriteFileWithMaxPerms(filepath, keyJSON, 0600)
	if err != nil {
		return cli.errorOut(errors.Wrapf(err, "Could not write %v", filepath))
	}

	_, err = os.Stderr.WriteString(fmt.Sprintf("🔑 Exported P2P key %s to %s\n", ID, filepath))
	if err != nil {
		return cli.errorOut(err)
	}

	return nil
}

// CreateOCRKeyBundle creates a key and inserts it into encrypted_ocr_key_bundles,
// protected by the password in the password file
func (cli *Client) CreateOCRKeyBundle(c *clipkg.Context) error {
	resp, err := cli.HTTP.Post("/v2/keys/ocr", nil)
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	var key ocrkey.EncryptedKeyBundle
	return cli.renderAPIResponse(resp, &key, "Created OCR key bundle")
}

// ListOCRKeyBundles lists the available OCR Key Bundles
func (cli *Client) ListOCRKeyBundles(c *clipkg.Context) error {
	resp, err := cli.HTTP.Get("/v2/keys/ocr", nil)
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	var keys []ocrkey.EncryptedKeyBundle
	return cli.renderAPIResponse(resp, &keys)
}

// DeleteOCRKeyBundle creates a key and inserts it into encrypted_ocr_keys,
// protected by the password in the password file
func (cli *Client) DeleteOCRKeyBundle(c *clipkg.Context) error {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass the key ID to be deleted"))
	}
	id, err := models.Sha256HashFromHex(c.Args().Get(0))
	if err != nil {
		return cli.errorOut(err)
	}

	if !confirmAction(c) {
		return nil
	}

	var queryStr string
	if c.Bool("hard") {
		queryStr = "?hard=true"
	}

	resp, err := cli.HTTP.Delete(fmt.Sprintf("/v2/keys/ocr/%s%s", id, queryStr))
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	var key ocrkey.EncryptedKeyBundle
	return cli.renderAPIResponse(resp, &key, "OCR key bundle deleted")
}

// ImportOCRKey imports OCR key bundle,
// file path must be passed
func (cli *Client) ImportOCRKey(c *clipkg.Context) (err error) {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass the filepath of the key to be imported"))
	}

	oldPasswordFile := c.String("oldpassword")
	if len(oldPasswordFile) == 0 {
		return cli.errorOut(errors.New("Must specify --oldpassword/-p flag"))
	}
	oldPassword, err := ioutil.ReadFile(oldPasswordFile)
	if err != nil {
		return cli.errorOut(errors.Wrap(err, "Could not read password file"))
	}

	filepath := c.Args().Get(0)
	keyJSON, err := ioutil.ReadFile(filepath)
	if err != nil {
		return cli.errorOut(err)
	}

	normalizedPassword := normalizePassword(string(oldPassword))
	resp, err := cli.HTTP.Post("/v2/keys/ocr/import?oldpassword="+normalizedPassword, bytes.NewReader(keyJSON))
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	var key ocrkey.EncryptedKeyBundle
	return cli.renderAPIResponse(resp, &key, "Imported OCR key bundle")
}

// ExportOCRKey exports OCR key bundles by ID
// ID of the key must be passed
func (cli *Client) ExportOCRKey(c *clipkg.Context) (err error) {
	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass the ID of the key to export"))
	}

	newPasswordFile := c.String("newpassword")
	if len(newPasswordFile) == 0 {
		return cli.errorOut(errors.New("Must specify --newpassword/-p flag"))
	}
	newPassword, err := ioutil.ReadFile(newPasswordFile)
	if err != nil {
		return cli.errorOut(errors.Wrap(err, "Could not read password file"))
	}

	filepath := c.String("output")
	if len(filepath) == 0 {
		return cli.errorOut(errors.New("Must specify --output/-o flag"))
	}

	ID := c.Args().Get(0)

	normalizedPassword := normalizePassword(string(newPassword))
	resp, err := cli.HTTP.Post("/v2/keys/ocr/export/"+ID+"?newpassword="+normalizedPassword, nil)
	if err != nil {
		return cli.errorOut(errors.Wrap(err, "Could not make HTTP request"))
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return cli.errorOut(errors.New("Error exporting"))
	}

	keyJSON, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return cli.errorOut(errors.Wrap(err, "Could not read response body"))
	}

	err = utils.WriteFileWithMaxPerms(filepath, keyJSON, 0600)
	if err != nil {
		return cli.errorOut(errors.Wrapf(err, "Could not write %v", filepath))
	}

	_, err = os.Stderr.WriteString(fmt.Sprintf("Exported OCR key bundle %s to %s", ID, filepath))
	if err != nil {
		return cli.errorOut(err)
	}

	return nil
}

func normalizePassword(password string) string {
	return url.PathEscape(strings.TrimSpace(password))
}

// SetLogLevel sets the log level on the node
func (cli *Client) SetLogLevel(c *clipkg.Context) (err error) {
	logLevel := c.String("level")
	request := web.LogPatchRequest{Level: logLevel}
	requestData, err := json.Marshal(request)
	if err != nil {
		return cli.errorOut(err)
	}

	buf := bytes.NewBuffer(requestData)
	resp, err := cli.HTTP.Patch("/v2/log", buf)
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	var lR webPresenter.LogResource
	err = cli.renderAPIResponse(resp, &lR)
	return err
}

// SetLogSQL enables or disables the log sql statemnts
func (cli *Client) SetLogSQL(c *clipkg.Context) (err error) {

	// Enforces selection of --enable or --disable
	if !c.Bool("enable") && !c.Bool("disable") {
		return cli.errorOut(errors.New("Must set logSql --enabled || --disable"))
	}

	// Sets logSql to true || false based on the --enabled flag
	logSql := c.Bool("enable")

	if err != nil {
		return cli.errorOut(err)
	}
	request := web.LogPatchRequest{SqlEnabled: &logSql}
	requestData, err := json.Marshal(request)
	if err != nil {
		return cli.errorOut(err)
	}

	buf := bytes.NewBuffer(requestData)
	resp, err := cli.HTTP.Patch("/v2/log", buf)
	if err != nil {
		return cli.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	var lR webPresenter.LogResource
	err = cli.renderAPIResponse(resp, &lR)
	return err
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

func parseResponse(resp *http.Response) ([]byte, error) {
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return b, multierr.Append(errors.New(resp.Status), err)
	}
	if resp.StatusCode == http.StatusUnauthorized {
		return b, errUnauthorized
	} else if resp.StatusCode >= http.StatusBadRequest {
		return b, errors.New("Error")
	}
	return b, err
}
