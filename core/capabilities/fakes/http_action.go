package fakes

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"time"

	commonCap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	customhttp "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/actions/http"
	httpserver "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/actions/http/server"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

var _ httpserver.ClientCapability = (*DirectHTTPAction)(nil)
var _ services.Service = (*DirectHTTPAction)(nil)
var _ commonCap.ExecutableCapability = (*DirectHTTPAction)(nil)

const HTTPActionID = "http-action@1.0.0"
const HTTPActionServiceName = "HttpActionService"

var directHTTPActionInfo = commonCap.MustNewCapabilityInfo(
	HTTPActionID,
	commonCap.CapabilityTypeAction,
	"An action that makes a direct HTTP request",
)

type DirectHTTPAction struct {
	commonCap.CapabilityInfo
	services.Service
	eng *services.Engine

	lggr logger.Logger
}

func NewDirectHTTPAction(lggr logger.Logger) *DirectHTTPAction {
	fc := &DirectHTTPAction{
		lggr: lggr,
	}

	fc.Service, fc.eng = services.Config{
		Name:  "directHttpAction",
		Start: fc.Start,
		Close: fc.Close,
	}.NewServiceEngine(lggr)
	return fc
}

func (fh *DirectHTTPAction) SendRequest(ctx context.Context, metadata commonCap.RequestMetadata, input *customhttp.Request) (*customhttp.Response, error) {
	fh.eng.Infow("Http Action SendRequest Started", "input", input)

	// Create HTTP client with timeout
	timeout := time.Duration(30) * time.Second // default timeout
	if input.GetTimeoutMs() > 0 {
		timeout = time.Duration(input.GetTimeoutMs()) * time.Millisecond
	}

	client := &http.Client{
		Timeout: timeout,
	}

	// Determine HTTP method (default to GET if not specified)
	method := input.GetMethod()
	if method == "" {
		method = "GET"
	}
	method = strings.ToUpper(method)

	// Create request body
	var body io.Reader
	if len(input.GetBody()) > 0 {
		body = bytes.NewReader(input.GetBody())
	}

	// Create the HTTP request
	req, err := http.NewRequestWithContext(ctx, method, input.GetUrl(), body)
	if err != nil {
		fh.eng.Errorw("Failed to create HTTP request", "error", err)
		return &customhttp.Response{
			ErrorMessage: err.Error(),
			StatusCode:   0,
		}, err
	}

	// Add headers
	for k, v := range input.GetHeaders() {
		req.Header.Set(k, v)
	}

	// Make the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		fh.eng.Errorw("Failed to execute HTTP request", "error", err)
		return &customhttp.Response{
			ErrorMessage: err.Error(),
			StatusCode:   0,
		}, err
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fh.eng.Errorw("Failed to read response body", "error", err)
		return &customhttp.Response{
			ErrorMessage: err.Error(),
			StatusCode:   uint32(resp.StatusCode), //nolint:gosec // status code is always in valid range
		}, err
	}

	// Convert headers
	headers := make(map[string]string)
	for k, v := range resp.Header {
		// Join multiple header values with comma
		headers[k] = strings.Join(v, ", ")
	}

	// Create response
	response := &customhttp.Response{
		StatusCode: uint32(resp.StatusCode), //nolint:gosec // status code is always in valid range
		Headers:    headers,
		Body:       respBody,
	}

	// Add error message if status code indicates an error
	if resp.StatusCode >= 400 {
		response.ErrorMessage = resp.Status
	}

	fh.eng.Debugw("HTTP request completed", "status", resp.StatusCode, "url", input.GetUrl())
	return response, nil
}

func (fh *DirectHTTPAction) Start(ctx context.Context) error {
	fh.eng.Infow("Http Action Start Started")
	return nil
}

func (fh *DirectHTTPAction) Close() error {
	fh.eng.Infow("Http Action Close Started")
	return nil
}

func (fh *DirectHTTPAction) Name() string {
	return HTTPActionServiceName
}

func (fh *DirectHTTPAction) Description() string {
	return directHTTPActionInfo.Description
}

func (fh *DirectHTTPAction) Ready() error {
	return nil
}

func (fh *DirectHTTPAction) Initialise(ctx context.Context, config string, _ core.TelemetryService,
	_ core.KeyValueStore,
	_ core.ErrorLog,
	_ core.PipelineRunnerService,
	_ core.RelayerSet,
	_ core.OracleFactory,
	_ core.GatewayConnector) error {
	// TODO: do validation of config here

	err := fh.Start(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (fh *DirectHTTPAction) Execute(ctx context.Context, request commonCap.CapabilityRequest) (commonCap.CapabilityResponse, error) {
	fh.eng.Infow("Direct Http Action Execute Started", "request", request)
	return commonCap.CapabilityResponse{}, nil
}

func (fh *DirectHTTPAction) RegisterToWorkflow(ctx context.Context, request commonCap.RegisterToWorkflowRequest) error {
	fh.eng.Infow("Registered to Direct Http Action", "workflowID", request.Metadata.WorkflowID)
	return nil
}

func (fh *DirectHTTPAction) UnregisterFromWorkflow(ctx context.Context, request commonCap.UnregisterFromWorkflowRequest) error {
	fh.eng.Infow("Unregistered from Direct Http Action", "workflowID", request.Metadata.WorkflowID)
	return nil
}
