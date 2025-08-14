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
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
)

var _ httpserver.ClientCapability = (*DirectHTTPAction)(nil)
var _ services.Service = (*DirectHTTPAction)(nil)
var _ commonCap.ExecutableCapability = (*DirectHTTPAction)(nil)

const HTTPActionID = "http-actions@0.1.0"
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
		Name: "directHttpAction",
	}.NewServiceEngine(lggr)
	return fc
}

func (fh *DirectHTTPAction) SendRequest(ctx context.Context, metadata commonCap.RequestMetadata, input *customhttp.Request) (*commonCap.ResponseAndMetadata[*customhttp.Response], error) {
	fh.eng.Infow("HTTP Action SendRequest Started", "input", input)

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
		httpResponse := &customhttp.Response{
			StatusCode: 0,
		}
		responseAndMetadata := commonCap.ResponseAndMetadata[*customhttp.Response]{
			Response:         httpResponse,
			ResponseMetadata: commonCap.ResponseMetadata{},
		}
		return &responseAndMetadata, err
	}

	// Add headers
	for k, v := range input.GetHeaders() {
		req.Header.Set(k, v)
	}

	// Make the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		fh.eng.Errorw("Failed to execute HTTP request", "error", err)
		httpResponse := &customhttp.Response{
			StatusCode: 0,
		}
		responseAndMetadata := commonCap.ResponseAndMetadata[*customhttp.Response]{
			Response:         httpResponse,
			ResponseMetadata: commonCap.ResponseMetadata{},
		}
		return &responseAndMetadata, err
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fh.eng.Errorw("Failed to read response body", "error", err)
		httpResponse := &customhttp.Response{
			StatusCode: uint32(resp.StatusCode), //nolint:gosec // status code is always in valid range
		}
		responseAndMetadata := commonCap.ResponseAndMetadata[*customhttp.Response]{
			Response:         httpResponse,
			ResponseMetadata: commonCap.ResponseMetadata{},
		}
		return &responseAndMetadata, err
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
	responseAndMetadata := commonCap.ResponseAndMetadata[*customhttp.Response]{
		Response:         response,
		ResponseMetadata: commonCap.ResponseMetadata{},
	}
	fh.eng.Infow("HTTP Action Finished", "Status", resp.StatusCode, "URL", input.GetUrl())
	return &responseAndMetadata, nil
}

func (fh *DirectHTTPAction) Description() string {
	return directHTTPActionInfo.Description
}

func (fh *DirectHTTPAction) Initialise(ctx context.Context, config string, _ core.TelemetryService,
	_ core.KeyValueStore,
	_ core.ErrorLog,
	_ core.PipelineRunnerService,
	_ core.RelayerSet,
	_ core.OracleFactory,
	_ core.GatewayConnector,
	_ core.Keystore) error {
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
