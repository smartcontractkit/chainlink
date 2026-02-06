package fakes

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
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

// AESGCMEncryptionKeyName is the magic secret key name that triggers AES-GCM encryption.
// This must match the constant in confidential-compute/types/types.go
const AESGCMEncryptionKeyName = "san_marino_aes_gcm_encryption_key"

// FakeAESGCMEncryptionKey is a well-known 32-byte test key used by this fake for AES-256-GCM encryption.
// This allows testing encrypt_output functionality without real VaultDON secrets.
// WARNING: This key is for testing only and must not be used in production.
var FakeAESGCMEncryptionKey = []byte("test-key-for-fake-encrypt-32-by") // exactly 32 bytes for AES-256

var directConfidentialHTTPActionInfo = commonCap.MustNewCapabilityInfo(
	ConfidentialHTTPActionID,
	commonCap.CapabilityTypeCombined,
	"An action that makes a confidential HTTP request with secrets",
)

// aesGCMEncrypt encrypts plaintext using AES-GCM with the provided key.
// Returns nonce || ciphertext || tag.
func aesGCMEncrypt(plaintext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// hasEncryptionSecret checks if the VaultDonSecrets contain the magic AES-GCM encryption key.
func hasEncryptionSecret(secrets []*confidentialhttp.SecretIdentifier) bool {
	for _, s := range secrets {
		if s.GetKey() == AESGCMEncryptionKeyName {
			return true
		}
	}
	return false
}

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

func (fh *DirectConfidentialHTTPAction) SendRequest(ctx context.Context, metadata commonCap.RequestMetadata, input *confidentialhttp.ConfidentialHTTPRequest) (*commonCap.ResponseAndMetadata[*confidentialhttp.HTTPResponse], caperrors.Error) {
	fh.eng.Infow("Confidential HTTP Action SendRequest Started", "input", input, "secretsCount", len(input.GetVaultDonSecrets()))

	// Warn if secrets are provided - this fake does not handle secret resolution
	if len(input.GetVaultDonSecrets()) > 0 {
		fh.eng.Warnw("This fake does not handle secrets - VaultDonSecrets will be ignored. Template variables like {{.secretName}} will not be resolved.", "secretsCount", len(input.GetVaultDonSecrets()))
	}

	req := input.GetRequest()
	if req == nil {
		return nil, caperrors.NewPublicUserError(errors.New("request cannot be nil"), caperrors.InvalidArgument)
	}

	fh.eng.Infow("Processing confidential HTTP request", "url", req.GetUrl(), "method", req.GetMethod())

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
	if bodyStr := req.GetBodyString(); bodyStr != "" {
		body = bytes.NewReader([]byte(bodyStr))
	} else if bodyBytes := req.GetBodyBytes(); len(bodyBytes) > 0 {
		body = bytes.NewReader(bodyBytes)
	}

	// Create the HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, method, req.GetUrl(), body)
	if err != nil {
		fh.eng.Errorw("Failed to create HTTP request", "error", err)
		return nil, caperrors.NewPublicUserError(fmt.Errorf("failed to create HTTP request: %w", err), caperrors.InvalidArgument)
	}

	// Add headers from multi_headers (map[string]*HeaderValues)
	for name, headerValues := range req.GetMultiHeaders() {
		if headerValues != nil {
			for _, value := range headerValues.GetValues() {
				httpReq.Header.Add(name, value)
			}
		}
	}

	// Make the HTTP request
	resp, err := client.Do(httpReq)
	if err != nil {
		fh.eng.Errorw("Failed to execute confidential HTTP request", "error", err)
		return nil, caperrors.NewPublicUserError(fmt.Errorf("failed to execute HTTP request: %w", err), caperrors.InvalidArgument)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fh.eng.Errorw("Failed to read response body", "error", err)
		return nil, caperrors.NewPublicUserError(fmt.Errorf("failed to read response body: %w", err), caperrors.InvalidArgument)
	}

	// Encrypt response if encrypt_output is true and the magic secret is present
	if input.GetEncryptOutput() && hasEncryptionSecret(input.GetVaultDonSecrets()) {
		fh.eng.Infow("Encrypting response body with fake AES-GCM key", "originalSize", len(respBody))
		encryptedBody, encErr := aesGCMEncrypt(respBody, FakeAESGCMEncryptionKey)
		if encErr != nil {
			fh.eng.Errorw("Failed to encrypt response body", "error", encErr)
			return nil, caperrors.NewPublicUserError(fmt.Errorf("failed to encrypt response body: %w", encErr), caperrors.Internal)
		}
		respBody = encryptedBody
		fh.eng.Infow("Response body encrypted", "encryptedSize", len(respBody))
	} else if input.GetEncryptOutput() {
		fh.eng.Warnw("encrypt_output is true but no AES-GCM encryption secret found - response will not be encrypted. TDH2 encryption is not supported by this fake.", "secretsCount", len(input.GetVaultDonSecrets()))
	}

	// Convert response headers to map[string]*HeaderValues
	responseHeaders := make(map[string]*confidentialhttp.HeaderValues)
	for name, values := range resp.Header {
		responseHeaders[name] = &confidentialhttp.HeaderValues{
			Values: values,
		}
	}

	// Create response
	response := &confidentialhttp.HTTPResponse{
		StatusCode:   uint32(resp.StatusCode), //nolint:gosec // HTTP status codes are always positive (100-599)
		Body:         respBody,
		MultiHeaders: responseHeaders,
	}

	responseAndMetadata := commonCap.ResponseAndMetadata[*confidentialhttp.HTTPResponse]{
		Response:         response,
		ResponseMetadata: commonCap.ResponseMetadata{},
	}

	fh.eng.Infow("Confidential HTTP Action Finished", "status", resp.StatusCode, "url", req.GetUrl())
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
