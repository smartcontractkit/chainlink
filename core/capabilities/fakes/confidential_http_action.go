package fakes

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	commonCap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	caperrors "github.com/smartcontractkit/chainlink-common/pkg/capabilities/errors"
	confidentialhttp "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/actions/confidentialhttp"
	httpserver "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/actions/confidentialhttp/server"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
)

var _ httpserver.ClientCapability = (*DirectConfidentialHTTPAction)(nil)
var _ services.Service = (*DirectConfidentialHTTPAction)(nil)
var _ commonCap.ExecutableCapability = (*DirectConfidentialHTTPAction)(nil)

const ConfidentialHTTPActionID = "confidential-http@1.0.0-alpha"
const ConfidentialHTTPActionServiceName = "ConfidentialHttpActionService"

var directConfidentialHTTPActionInfo = commonCap.MustNewCapabilityInfo(
	ConfidentialHTTPActionID,
	commonCap.CapabilityTypeCombined,
	"An action that makes a confidential HTTP request with secrets",
)

type DirectConfidentialHTTPAction struct {
	commonCap.CapabilityInfo
	services.Service
	eng *services.Engine

	lggr logger.Logger
}

func NewDirectConfidentialHTTPAction(lggr logger.Logger) *DirectConfidentialHTTPAction {
	fc := &DirectConfidentialHTTPAction{
		lggr: lggr,
	}

	fc.Service, fc.eng = services.Config{
		Name: "directConfidentialHttpAction",
	}.NewServiceEngine(lggr)
	return fc
}

func (fh *DirectConfidentialHTTPAction) SendRequests(ctx context.Context, metadata commonCap.RequestMetadata, input *confidentialhttp.EnclaveActionInput) (*commonCap.ResponseAndMetadata[*confidentialhttp.HTTPEnclaveResponseData], caperrors.Error) {
	fh.eng.Infow("Confidential HTTP Action SendRequests Started", "input", input, "secretsCount", len(input.GetVaultDonSecrets()))

	// Warn if secrets are provided - this fake does not handle secret resolution
	if len(input.GetVaultDonSecrets()) > 0 {
		fh.eng.Warnw("This fake does not handle secrets - VaultDonSecrets will be ignored. Template variables like {{.secretName}} will not be resolved.", "secretsCount", len(input.GetVaultDonSecrets()))
	}

	if input.GetInput() == nil {
		return nil, caperrors.NewPublicUserError(errors.New("input cannot be nil"), caperrors.InvalidArgument)
	}

	requests := input.GetInput().GetRequests()
	if len(requests) == 0 {
		return nil, caperrors.NewPublicUserError(errors.New("no requests provided"), caperrors.InvalidArgument)
	}

	// Process each request
	responses := make([]*confidentialhttp.ResponseTemplate, 0, len(requests))

	for i, req := range requests {
		fh.eng.Infow("Processing confidential HTTP request", "index", i, "url", req.GetUrl(), "method", req.GetMethod())

		// Create HTTP client with timeout (default 30 seconds)
		timeout := time.Duration(30) * time.Second
		client := &http.Client{
			Timeout: timeout,
		}

		// Validate HTTP method
		method := strings.TrimSpace(req.GetMethod())
		if method == "" {
			return nil, caperrors.NewPublicUserError(errors.New("http method cannot be empty"), caperrors.InvalidArgument)
		}
		method = strings.ToUpper(method)

		// Create request body
		var body io.Reader
		if len(req.GetBody()) > 0 {
			body = bytes.NewReader([]byte(req.GetBody()))
		}

		// Create the HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, method, req.GetUrl(), body)
		if err != nil {
			fh.eng.Errorw("Failed to create HTTP request", "error", err, "index", i)
			responses = append(responses, &confidentialhttp.ResponseTemplate{
				StatusCode: 0,
				Body:       []byte(err.Error()),
			})
			continue
		}

		// Add headers
		for _, header := range req.GetHeaders() {
			parts := strings.SplitN(header, ":", 2)
			if len(parts) == 2 {
				httpReq.Header.Set(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
			}
		}

		// Make the HTTP request
		resp, err := client.Do(httpReq)
		if err != nil {
			fh.eng.Errorw("Failed to execute confidential HTTP request", "error", err, "index", i)
			responses = append(responses, &confidentialhttp.ResponseTemplate{
				StatusCode: 0,
				Body:       []byte(err.Error()),
			})
			continue
		}

		// Read response body
		respBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			fh.eng.Errorw("Failed to read response body", "error", err, "index", i)
			responses = append(responses, &confidentialhttp.ResponseTemplate{
				StatusCode: int64(resp.StatusCode),
				Body:       []byte(err.Error()),
			})
			continue
		}

		// Create response template
		response := &confidentialhttp.ResponseTemplate{
			StatusCode: int64(resp.StatusCode),
			Body:       respBody,
		}
		responses = append(responses, response)
		fh.eng.Infow("Confidential HTTP request finished", "index", i, "status", resp.StatusCode, "url", req.GetUrl())
	}

	// Create response data
	responseData := &confidentialhttp.HTTPEnclaveResponseData{
		Responses: responses,
	}

	responseAndMetadata := commonCap.ResponseAndMetadata[*confidentialhttp.HTTPEnclaveResponseData]{
		Response:         responseData,
		ResponseMetadata: commonCap.ResponseMetadata{},
	}

	fh.eng.Infow("Confidential HTTP Action Finished", "requestCount", len(requests), "responseCount", len(responses))
	return &responseAndMetadata, nil
}

func (fh *DirectConfidentialHTTPAction) Description() string {
	return directConfidentialHTTPActionInfo.Description
}

func (fh *DirectConfidentialHTTPAction) Initialise(ctx context.Context, dependencies core.StandardCapabilitiesDependencies) error {
	// No config validation needed for this fake implementation
	return fh.Start(ctx)
}

func (fh *DirectConfidentialHTTPAction) Execute(ctx context.Context, request commonCap.CapabilityRequest) (commonCap.CapabilityResponse, error) {
	fh.eng.Infow("Direct Confidential Http Action Execute Started", "request", request)
	return commonCap.CapabilityResponse{}, nil
}

func (fh *DirectConfidentialHTTPAction) RegisterToWorkflow(ctx context.Context, request commonCap.RegisterToWorkflowRequest) error {
	fh.eng.Infow("Registered to Direct Confidential Http Action", "workflowID", request.Metadata.WorkflowID)
	return nil
}

func (fh *DirectConfidentialHTTPAction) UnregisterFromWorkflow(ctx context.Context, request commonCap.UnregisterFromWorkflowRequest) error {
	fh.eng.Infow("Unregistered from Direct Confidential Http Action", "workflowID", request.Metadata.WorkflowID)
	return nil
}
