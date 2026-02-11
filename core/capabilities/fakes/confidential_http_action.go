package fakes

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"

	"gopkg.in/yaml.v3"

	commonCap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	caperrors "github.com/smartcontractkit/chainlink-common/pkg/capabilities/errors"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/actions/confidentialhttp"
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

var directConfidentialHTTPActionInfo = commonCap.MustNewCapabilityInfo(
	ConfidentialHTTPActionID,
	commonCap.CapabilityTypeCombined,
	"An action that makes a confidential HTTP request with secrets",
)

type SecretsConfig struct {
	SecretsNames map[string][]string `yaml:"secretsNames"`
}

type DirectConfidentialHTTPAction struct {
	commonCap.CapabilityInfo
	services.Service
	eng *services.Engine

	secretsConfig SecretsConfig
	lggr          logger.Logger
}

func NewDirectConfidentialHTTPAction(lggr logger.Logger, secretsPath string) *DirectConfidentialHTTPAction {
	fc := &DirectConfidentialHTTPAction{
		lggr: lggr,
	}

	// Load secrets
	secretsFile := secretsPath
	if envFile := os.Getenv("SECRETS_FILE"); envFile != "" {
		secretsFile = envFile
	}

	if data, err := os.ReadFile(secretsFile); err == nil {
		if marshalErr := yaml.Unmarshal(data, &fc.secretsConfig); marshalErr != nil {
			lggr.Warnf("Failed to parse secrets file %s: %v", secretsFile, marshalErr)
		} else {
			lggr.Infof("Loaded secrets from %s", secretsFile)
			// Resolve environment variables
			for key, val := range fc.secretsConfig.SecretsNames {
				secrets := val
				for i, v := range secrets {
					if envVal := os.Getenv(v); envVal != "" {
						secrets[i] = envVal
					} else {
						lggr.Warnf("Secret environment variable %s not set", v)
					}
				}
				fc.secretsConfig.SecretsNames[key] = secrets
			}
		}
	} else {
		lggr.Infof("Could not read secrets file %s: %v. Continuing without local secrets.", secretsFile, err)
	}

	fc.Service, fc.eng = services.Config{
		Name: "directConfidentialHttpAction",
	}.NewServiceEngine(lggr)
	return fc
}

func (fh *DirectConfidentialHTTPAction) SendRequest(ctx context.Context, metadata commonCap.RequestMetadata, input *confidentialhttp.ConfidentialHTTPRequest) (*commonCap.ResponseAndMetadata[*confidentialhttp.HTTPResponse], caperrors.Error) {
	fh.eng.Infow("Confidential HTTP Action SendRequest Started", "input", input, "secretsCount", len(input.GetVaultDonSecrets()))

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

	// Prepare template data from loaded secrets
	templateData := make(map[string]interface{})
	for k, v := range fh.secretsConfig.SecretsNames {
		if len(v) == 1 {
			templateData[k] = v[0]
		} else {
			templateData[k] = v
		}
	}

	usesBody := method == "POST" || method == "PUT" || method == "PATCH" || method == "DELETE"
	bodyString := req.GetBodyString()
	bodyBytes := req.GetBodyBytes()

	var httpReq *http.Request
	var err error

	switch {
	case usesBody && bodyString != "":
		processedBody := &bytes.Buffer{}
		bodyTmpl, err2 := template.New("body").Parse(bodyString)
		if err2 != nil {
			fh.eng.Errorf("error parsing body template: %v", err2)
			return nil, caperrors.NewPublicUserError(errors.New("error parsing body template"), caperrors.InvalidArgument)
		}
		if err2 = bodyTmpl.Execute(processedBody, templateData); err2 != nil {
			fh.eng.Errorf("error executing body template: %v", err2)
			return nil, caperrors.NewPublicUserError(errors.New("error executing body template"), caperrors.InvalidArgument)
		}
		httpReq, err = http.NewRequestWithContext(ctx, method, req.GetUrl(), processedBody)
	case usesBody && len(bodyBytes) > 0:
		httpReq, err = http.NewRequestWithContext(ctx, method, req.GetUrl(), bytes.NewReader(bodyBytes))
	default:
		httpReq, err = http.NewRequestWithContext(ctx, method, req.GetUrl(), nil)
	}

	if err != nil {
		fh.eng.Errorw("Failed to create HTTP request", "error", err)
		return nil, caperrors.NewPublicUserError(fmt.Errorf("failed to create HTTP request: %w", err), caperrors.InvalidArgument)
	}

	// Add headers with template processing
	for name, headerValues := range req.GetMultiHeaders() {
		if headerValues != nil {
			for _, value := range headerValues.GetValues() {
				headerTmpl, tmplErr := template.New("header").Parse(value)
				if tmplErr != nil {
					fh.eng.Errorf("error parsing header template for %s: %v", name, tmplErr)
					return nil, caperrors.NewPublicUserError(errors.New("error parsing header template"), caperrors.InvalidArgument)
				}

				var processedHeader bytes.Buffer
				if tmplErr = headerTmpl.Execute(&processedHeader, templateData); tmplErr != nil {
					fh.eng.Errorf("error executing header template for %s: %v", name, tmplErr)
					return nil, caperrors.NewPublicUserError(errors.New("error executing header template"), caperrors.InvalidArgument)
				}

				httpReq.Header.Add(name, processedHeader.String())
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

	// Encrypt response if encrypt_output is true
	if input.GetRequest().GetEncryptOutput() {
		// Priority 1: Use real AES key from secrets.yaml if available
		if secretValues, exists := fh.secretsConfig.SecretsNames[AESGCMEncryptionKeyName]; exists && len(secretValues) > 0 {
			secretKeyStr := secretValues[0]
			var secretKey []byte

			// Try hex decoding first
			if decoded, decErr := hex.DecodeString(secretKeyStr); decErr == nil && len(decoded) == 32 {
				secretKey = decoded
			} else {
				// Use as raw bytes
				secretKey = []byte(secretKeyStr)
			}

			fh.eng.Infow("Encrypting response body with AES-GCM key from secrets.yaml", "originalSize", len(respBody))
			encryptedBody, encErr := AESGCMEncrypt(respBody, secretKey)
			if encErr != nil {
				fh.eng.Errorw("Failed to encrypt response body with AES-GCM", "error", encErr)
				return nil, caperrors.NewPublicUserError(fmt.Errorf("failed to encrypt response body: %w", encErr), caperrors.Internal)
			}
			respBody = encryptedBody
			fh.eng.Infow("Response body encrypted with AES-GCM (secrets.yaml key)", "encryptedSize", len(respBody))
		} else if hasEncryptionSecret(input.GetVaultDonSecrets()) {
			secretHex := fh.secretsConfig.SecretsNames[AESGCMEncryptionKeyName][0]
			secretKey, decErr := hex.DecodeString(secretHex)
			if decErr != nil {
				return nil, caperrors.NewPublicUserError(fmt.Errorf("could not decode secret key %w", decErr), caperrors.Internal)
			}
			encryptedBody, encErr := AESGCMEncrypt(respBody, secretKey)
			if encErr != nil {
				fh.eng.Errorw("Failed to encrypt response body with AES-GCM", "error", encErr)
				return nil, caperrors.NewPublicUserError(fmt.Errorf("failed to encrypt response body: %w", encErr), caperrors.Internal)
			}
			respBody = encryptedBody
			fh.eng.Infow("Response body encrypted with AES-GCM (fake key)", "encryptedSize", len(respBody))
		} else {
			fh.eng.Warn("encrypt_output is true but no encryption key found in secrets.yaml or vault DON secrets. Returning plaintext response.")
		}
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

// AESGCMEncrypt encrypts plaintext using AES-GCM with the provided key.
// Returns nonce || ciphertext || tag.
func AESGCMEncrypt(plaintext []byte, key []byte) ([]byte, error) {
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
