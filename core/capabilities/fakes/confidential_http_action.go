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
	"sync"
	"time"

	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"

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
var FakeAESGCMEncryptionKey = []byte("test-key-for-fake-encrypt-32-byt") // exactly 32 bytes for AES-256

// TDH2 test keys for threshold encryption when AES-GCM key is not provided.
// These are lazily initialized on first use.
var (
	fakeTDH2InitOnce      sync.Once
	fakeTDH2PublicKey     *tdh2easy.PublicKey
	fakeTDH2PrivateShares []*tdh2easy.PrivateShare
	fakeTDH2InitErr       error
)

// FakeTDH2Threshold is the threshold for the fake TDH2 key setup (k of n).
const FakeTDH2Threshold = 2

// FakeTDH2TotalShares is the total number of shares for the fake TDH2 key setup.
const FakeTDH2TotalShares = 3

// initFakeTDH2Keys lazily initializes the fake TDH2 keys.
func initFakeTDH2Keys() (*tdh2easy.PublicKey, []*tdh2easy.PrivateShare, error) {
	fakeTDH2InitOnce.Do(func() {
		_, pubKey, privShares, err := tdh2easy.GenerateKeys(FakeTDH2Threshold, FakeTDH2TotalShares)
		if err != nil {
			fakeTDH2InitErr = fmt.Errorf("failed to generate fake TDH2 keys: %w", err)
			return
		}
		fakeTDH2PublicKey = pubKey
		fakeTDH2PrivateShares = privShares
	})
	return fakeTDH2PublicKey, fakeTDH2PrivateShares, fakeTDH2InitErr
}

// GetFakeTDH2PublicKey returns the fake TDH2 public key used by this fake for encryption.
// Tests can use this to verify that encryption occurred or to set up decryption.
func GetFakeTDH2PublicKey() (*tdh2easy.PublicKey, error) {
	pubKey, _, err := initFakeTDH2Keys()
	return pubKey, err
}

// GetFakeTDH2PrivateShares returns the fake TDH2 private key shares used by this fake.
// Tests can use these to decrypt responses encrypted with TDH2.
func GetFakeTDH2PrivateShares() ([]*tdh2easy.PrivateShare, error) {
	_, privShares, err := initFakeTDH2Keys()
	return privShares, err
}

// DecryptFakeTDH2Ciphertext decrypts a TDH2 ciphertext that was encrypted by this fake.
// This is a helper for tests that need to verify the decrypted content.
func DecryptFakeTDH2Ciphertext(ciphertextBytes []byte) ([]byte, error) {
	pubKey, privShares, err := initFakeTDH2Keys()
	if err != nil {
		return nil, err
	}

	ciphertext := &tdh2easy.Ciphertext{}
	if err := ciphertext.UnmarshalVerify(ciphertextBytes, pubKey); err != nil {
		return nil, fmt.Errorf("failed to unmarshal TDH2 ciphertext: %w", err)
	}

	// Create decryption shares using threshold number of private key shares
	shares := make([]*tdh2easy.DecryptionShare, FakeTDH2Threshold)
	for i := 0; i < FakeTDH2Threshold; i++ {
		share, err := tdh2easy.Decrypt(ciphertext, privShares[i])
		if err != nil {
			return nil, fmt.Errorf("failed to create decryption share %d: %w", i, err)
		}
		shares[i] = share
	}

	// Aggregate shares to recover plaintext (n = total participants)
	plaintext, err := tdh2easy.Aggregate(ciphertext, shares, FakeTDH2TotalShares)
	if err != nil {
		return nil, fmt.Errorf("failed to aggregate TDH2 shares: %w", err)
	}

	return plaintext, nil
}

// tdh2Encrypt encrypts plaintext using TDH2 with the fake public key.
func tdh2Encrypt(plaintext []byte) ([]byte, error) {
	pubKey, _, err := initFakeTDH2Keys()
	if err != nil {
		return nil, err
	}

	ciphertext, err := tdh2easy.Encrypt(pubKey, plaintext)
	if err != nil {
		return nil, fmt.Errorf("failed to TDH2 encrypt: %w", err)
	}

	ciphertextBytes, err := ciphertext.Marshal()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal TDH2 ciphertext: %w", err)
	}

	return ciphertextBytes, nil
}

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

	// Note: This fake does not resolve secrets from VaultDON.
	// Template variables like {{.secretName}} will NOT be substituted.
	// However, we DO check for the magic AES-GCM encryption key name to determine encryption mode.
	if len(input.GetVaultDonSecrets()) > 0 {
		fh.eng.Infow("Secrets requested (template variables will NOT be resolved, but encryption key detection works)",
			"secretsCount", len(input.GetVaultDonSecrets()),
			"hasAESKey", hasEncryptionSecret(input.GetVaultDonSecrets()))
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

	// Encrypt response if encrypt_output is true
	if input.GetEncryptOutput() {
		if hasEncryptionSecret(input.GetVaultDonSecrets()) {
			// Use AES-GCM encryption with the magic key
			fh.eng.Infow("Encrypting response body with fake AES-GCM key", "originalSize", len(respBody))
			encryptedBody, encErr := aesGCMEncrypt(respBody, FakeAESGCMEncryptionKey)
			if encErr != nil {
				fh.eng.Errorw("Failed to encrypt response body with AES-GCM", "error", encErr)
				return nil, caperrors.NewPublicUserError(fmt.Errorf("failed to encrypt response body: %w", encErr), caperrors.Internal)
			}
			respBody = encryptedBody
			fh.eng.Infow("Response body encrypted with AES-GCM", "encryptedSize", len(respBody))
		} else {
			// Use TDH2 encryption with the fake master public key (mimics real enclave behavior)
			fh.eng.Infow("Encrypting response body with fake TDH2 key", "originalSize", len(respBody))
			encryptedBody, encErr := tdh2Encrypt(respBody)
			if encErr != nil {
				fh.eng.Errorw("Failed to encrypt response body with TDH2", "error", encErr)
				return nil, caperrors.NewPublicUserError(fmt.Errorf("failed to encrypt response body: %w", encErr), caperrors.Internal)
			}
			respBody = encryptedBody
			fh.eng.Infow("Response body encrypted with TDH2", "encryptedSize", len(respBody))
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
