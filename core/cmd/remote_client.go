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
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/tidwall/gjson"
	clipkg "github.com/urfave/cli"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/presenters"
	"github.com/smartcontractkit/chainlink/core/web"
	webpresenters "github.com/smartcontractkit/chainlink/core/web/presenters"
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

	var ei webpresenters.ExternalInitiatorAuthentication
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

// ReplayFromBlock replays chain data from the given block number until the most recent
func (cli *Client) ReplayFromBlock(c *clipkg.Context) (err error) {

	blockNumber := c.Int64("block-number")
	if blockNumber <= 0 {
		return cli.errorOut(errors.New("Must pass a positive value in '--block-number' parameter"))
	}

	buf := bytes.NewBufferString("{}")

	resp, err := cli.HTTP.Post(fmt.Sprintf("/v2/replay_from_block/%v", blockNumber), buf)
	if err != nil {
		return cli.errorOut(err)
	}

	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		bytes, err2 := cli.parseResponse(resp)
		if err2 != nil {
			return errors.Wrap(err2, "parseResponse error")
		}
		return cli.errorOut(errors.New(string(bytes)))
	}

	err = cli.printResponseBody(resp)
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

	if resp.StatusCode != http.StatusOK {
		return cli.errorOut(errors.Errorf("Unexpected status code: %v", resp.Status))
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

// MigrateJobSpec migrates a JobSpec based on JSON input
func (cli *Client) MigrateJobSpec(c *clipkg.Context) (err error) {
	result, _, err := cli.MigrateJobSpecForResult(c)
	if err != nil {
		return err
	}

	fmt.Println("")
	fmt.Println(result)
	return nil
}

func (cli *Client) MigrateJobSpecForResult(c *clipkg.Context) (s string, j *job.Job, err error) {
	if !c.Args().Present() {
		return s, nil, cli.errorOut(errors.New("Must pass in JSON or filepath"))
	}

	buf, err := getBufferFromJSON(c.Args().First())
	if err != nil {
		return s, nil, cli.errorOut(err)
	}
	var jsr models.JobSpecRequest
	err = json.Unmarshal(buf.Bytes(), &jsr)
	if err != nil {
		return s, nil, cli.errorOut(err)
	}

	js := models.NewJobFromRequest(jsr)

	jobSpec, err := MigrateJobSpec(js)
	if err != nil {
		return s, nil, cli.errorOut(err)
	}

	var flat interface{}
	if jobSpec.CronSpec != nil {
		flat = job.JobCronFlat{
			ExternalJobID:     jobSpec.ExternalJobID,
			CronSchedule:      jobSpec.CronSpec.CronSchedule,
			Type:              jobSpec.Type,
			SchemaVersion:     1,
			Name:              jobSpec.Name,
			ObservationSource: jobSpec.PipelineSpec.DotDagSource,
		}
	}
	if jobSpec.DirectRequestSpec != nil {
		flat = job.JobDirectRequestFlat{
			ExternalJobID:            jobSpec.ExternalJobID,
			ContractAddress:          jobSpec.DirectRequestSpec.ContractAddress,
			MinIncomingConfirmations: jobSpec.DirectRequestSpec.MinIncomingConfirmations,
			Requesters:               jobSpec.DirectRequestSpec.Requesters,
			Type:                     jobSpec.Type,
			SchemaVersion:            1,
			Name:                     jobSpec.Name,
			ObservationSource:        jobSpec.PipelineSpec.DotDagSource,
		}
	}
	result, err := toml.Marshal(flat)
	if err != nil {
		return s, nil, cli.errorOut(err)
	}
	return string(result), &jobSpec, nil
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

// RemoteLogin creates a cookie session to run remote commands.
func (cli *Client) RemoteLogin(c *clipkg.Context) error {
	sessionRequest, err := cli.buildSessionRequest(c.String("file"))
	if err != nil {
		return cli.errorOut(err)
	}
	_, err = cli.CookieAuthenticator.Authenticate(sessionRequest)
	return cli.errorOut(err)
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

	var svcLogConfig webpresenters.ServiceLogConfigResource
	err = cli.renderAPIResponse(resp, &svcLogConfig)
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

	var svcLogConfig webpresenters.ServiceLogConfigResource
	err = cli.renderAPIResponse(resp, &svcLogConfig)
	return err
}

// SetLogPkg sets the package log filter on the node
func (cli *Client) SetLogPkg(c *clipkg.Context) (err error) {
	pkg := strings.Split(c.String("pkg"), ",")
	level := strings.Split(c.String("level"), ",")

	serviceLogLevel := make([][2]string, len(pkg))
	for i, p := range pkg {
		serviceLogLevel[i][0] = p
		serviceLogLevel[i][1] = level[i]
	}

	request := web.LogPatchRequest{ServiceLogLevel: serviceLogLevel}
	requestData, err := json.Marshal(request)
	if err != nil {
		return cli.errorOut(err)
	}

	buf := bytes.NewBuffer(requestData)
	resp, err := cli.HTTP.Patch("/v2/log", buf)
	if err != nil {
		return cli.errorOut(errors.Wrap(err, "set pkg specific logging levels"))
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	var svcLogConfig webpresenters.ServiceLogConfigResource
	err = cli.renderAPIResponse(resp, &svcLogConfig)

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
