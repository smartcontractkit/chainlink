package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/mitchellh/go-homedir"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	clipkg "github.com/urfave/cli"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/bridges"
	"github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/sessions"
	"github.com/smartcontractkit/chainlink/core/static"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/web"
	webpresenters "github.com/smartcontractkit/chainlink/core/web/presenters"
)

var errUnauthorized = errors.New(http.StatusText(http.StatusUnauthorized))

// CreateExternalInitiator adds an external initiator
func (cli *Client) CreateExternalInitiator(c *clipkg.Context) (err error) {
	if c.NArg() != 1 && c.NArg() != 2 {
		return cli.errorOut(errors.New("create expects 1 - 2 arguments: a name and a url (optional)"))
	}

	var request bridges.ExternalInitiatorRequest
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

// ReplayFromBlock replays chain data from the given block number until the most recent
func (cli *Client) ReplayFromBlock(c *clipkg.Context) (err error) {
	blockNumber := c.Int64("block-number")
	if blockNumber <= 0 {
		return cli.errorOut(errors.New("Must pass a positive value in '--block-number' parameter"))
	}

	forceBroadcast := c.Bool("force")

	buf := bytes.NewBufferString("{}")
	resp, err := cli.HTTP.Post(
		fmt.Sprintf(
			"/v2/replay_from_block/%v?force=%s",
			blockNumber,
			strconv.FormatBool(forceBroadcast),
		), buf)
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

// RemoteLogin creates a cookie session to run remote commands.
func (cli *Client) RemoteLogin(c *clipkg.Context) error {
	lggr := cli.Logger.Named("RemoteLogin")
	sessionRequest, err := cli.buildSessionRequest(c.String("file"))
	if err != nil {
		return cli.errorOut(err)
	}
	_, err = cli.CookieAuthenticator.Authenticate(sessionRequest)
	if err != nil {
		return cli.errorOut(err)
	}
	err = cli.checkRemoteBuildCompatibility(lggr, c.Bool("bypass-version-check"), static.Version, static.Sha)
	if err != nil {
		return cli.errorOut(err)
	}
	fmt.Println("Successfully Logged In.")
	return nil
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

	switch resp.StatusCode {
	case http.StatusOK:
		fmt.Println("Password updated.")
	case http.StatusConflict:
		fmt.Println("Old password did not match.")
	default:
		return cli.printResponseBody(resp)
	}
	return nil
}

// Profile will collect pprof metrics and store them in a folder.
func (cli *Client) Profile(c *clipkg.Context) error {
	seconds := c.Uint("seconds")
	baseDir := c.String("output_dir")
	if seconds >= uint(cli.Config.HTTPServerWriteTimeout().Seconds()) {
		return cli.errorOut(errors.New("profile duration should be less than server write timeout."))
	}

	genDir := filepath.Join(baseDir, fmt.Sprintf("debuginfo-%s", time.Now().Format(time.RFC3339)))

	err := os.Mkdir(genDir, 0755)
	if err != nil {
		return cli.errorOut(err)
	}
	cli.Logger.Infof("writing pprof to %s", genDir)
	var wgPprof sync.WaitGroup
	vitals := []string{
		"allocs",       // A sampling of all past memory allocations
		"block",        // Stack traces that led to blocking on synchronization primitives
		"cmdline",      // The command line invocation of the current program
		"goroutine",    // Stack traces of all current goroutines
		"heap",         // A sampling of memory allocations of live objects.
		"mutex",        // Stack traces of holders of contended mutexes
		"profile",      // CPU profile.
		"threadcreate", // Stack traces that led to the creation of new OS threads
		"trace",        // A trace of execution of the current program.
	}
	wgPprof.Add(len(vitals))
	errs := make(chan error)
	for _, vt := range vitals {
		go func(vt string) {
			defer wgPprof.Done()
			uri := fmt.Sprintf("/v2/debug/pprof/%s?seconds=%d", vt, seconds)
			cli.Logger.Infof("Collecting %s ", uri)
			resp, err := cli.HTTP.Get(uri)
			if err != nil {
				cli.Logger.Errorf("error collecting vt %s: %s", vt, err.Error())
				errs <- err
				return
			}
			defer resp.Body.Close()

			// write to file
			f, err := os.Create(filepath.Join(genDir, vt))
			if err != nil {
				cli.Logger.Errorf("error creating file for %s: %s", vt, err.Error())
				errs <- err
				return
			}
			defer f.Close()

			_, err = io.Copy(f, resp.Body)
			if err != nil {
				cli.Logger.Errorf("error writing to file for %s: %s", vt, err.Error())
				errs <- err
				return
			}
			cli.Logger.Infof("Collected %s", vt)
		}(vt)
	}
	wgPprof.Wait()
	if len(errs) > 0 {
		return cli.errorOut(errors.New("One or more profile collections failed %s"))
	}
	return nil
}

func (cli *Client) buildSessionRequest(flag string) (sessions.SessionRequest, error) {
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
	if errors.Is(err, errUnauthorized) {
		return nil, cli.errorOut(multierr.Append(err, fmt.Errorf("your credentials may be missing, invalid or you may need to login first using the CLI via 'chainlink admin login'")))
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

// SetEvmGasPriceDefault specifies the minimum gas price to use for outgoing transactions
func (cli *Client) SetEvmGasPriceDefault(c *clipkg.Context) (err error) {
	var adjustedAmount *big.Int
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
	var chainID *big.Int
	if c.IsSet("evmChainID") {
		var ok bool
		chainID, ok = new(big.Int).SetString(c.String("evmChainID"), 10)
		if !ok {
			return cli.errorOut(fmt.Errorf("invalid evmChainID %s", value))
		}
	}
	adjustedAmount, _ = amount.Int(nil)

	request := struct {
		EvmGasPriceDefault string     `json:"ethGasPriceDefault"`
		EvmChainID         *utils.Big `json:"evmChainID"`
	}{
		EvmGasPriceDefault: adjustedAmount.String(),
		EvmChainID:         utils.NewBig(chainID),
	}

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
	cwl := config.ConfigPrinter{}
	err = cli.renderAPIResponse(resp, &cwl)
	return err
}

func normalizePassword(password string) string {
	return url.QueryEscape(strings.TrimSpace(password))
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

// SetLogSQL enables or disables the log sql statements
func (cli *Client) SetLogSQL(c *clipkg.Context) (err error) {

	// Enforces selection of --enable or --disable
	if !c.Bool("enable") && !c.Bool("disable") {
		return cli.errorOut(errors.New("Must set logSql --enabled || --disable"))
	}

	// Sets logSql to true || false based on the --enabled flag
	logSql := c.Bool("enable")

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

func (cli *Client) checkRemoteBuildCompatibility(lggr logger.Logger, onlyWarn bool, cliVersion, cliSha string) error {
	resp, err := cli.HTTP.Get("/v2/build_info")
	if err != nil {
		lggr.Warnw("Got error querying for version. Remote node version is unknown and CLI may behave in unexpected ways.", "err", err)
		return nil
	}
	b, err := parseResponse(resp)
	if err != nil {
		lggr.Warnw("Got error parsing http response for remote version. Remote node version is unknown and CLI may behave in unexpected ways.", "resp", resp, "err", err)
		return nil
	}

	var remoteBuildInfo map[string]string
	if err := json.Unmarshal(b, &remoteBuildInfo); err != nil {
		lggr.Warnw("Got error json parsing bytes from remote version response. Remote node version is unknown and CLI may behave in unexpected ways.", "bytes", b, "err", err)
		return nil
	}
	remoteVersion, remoteSha := remoteBuildInfo["version"], remoteBuildInfo["commitSHA"]

	remoteSemverUnset := remoteVersion == "unset" || remoteVersion == "" || remoteSha == "unset" || remoteSha == ""
	cliRemoteSemverMismatch := remoteVersion != cliVersion || remoteSha != cliSha

	if remoteSemverUnset || cliRemoteSemverMismatch {
		// Show a warning but allow mismatch
		if onlyWarn {
			lggr.Warnf("CLI build (%s@%s) mismatches remote node build (%s@%s), it might behave in unexpected ways", remoteVersion, remoteSha, cliVersion, cliSha)
			return nil
		}
		// Don't allow usage of CLI by unsetting the session cookie to prevent further requests
		cli.CookieAuthenticator.Logout()
		return ErrIncompatible{CLIVersion: cliVersion, CLISha: cliSha, RemoteVersion: remoteVersion, RemoteSha: remoteSha}
	}
	return nil
}

// ErrIncompatible is returned when the cli and remote versions are not compatible.
type ErrIncompatible struct {
	CLIVersion, CLISha       string
	RemoteVersion, RemoteSha string
}

func (e ErrIncompatible) Error() string {
	return fmt.Sprintf("error: CLI build (%s@%s) mismatches remote node build (%s@%s). You can set flag --bypass-version-check to bypass this", e.CLIVersion, e.CLISha, e.RemoteVersion, e.RemoteSha)
}
