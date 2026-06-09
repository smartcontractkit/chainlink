package cre

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"io"
	"maps"
	"math/big"
	"math/rand"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-protos/cre/go/values"

	vault_helpers "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"
	workflow_registry_v2_wrapper "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v2"
	commonevents "github.com/smartcontractkit/chainlink-protos/workflows/go/common"
	workflowevents "github.com/smartcontractkit/chainlink-protos/workflows/go/events"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	ctfblockchain "github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/evm"
	stvault "github.com/smartcontractkit/chainlink/system-tests/lib/cre/vault"
	creworkflow "github.com/smartcontractkit/chainlink/system-tests/lib/cre/workflow"
	vaultsecret_config "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/vaultsecret/config"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
	vaultjwt "github.com/smartcontractkit/chainlink/v2/core/capabilities/vault"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaultutils"
)

const (
	vaultDefaultConfigPath              = "/configs/workflow-gateway-capabilities-don.toml"
	vaultJWTAuthEnabledConfigPath       = "/configs/workflow-gateway-capabilities-don-vault-jwt_auth-enabled.toml"
	vaultOptimizationsEnabledConfigPath = "/configs/workflow-gateway-capabilities-don-vault-optimizations-enabled.toml"
	vaultJWTIssuerListenAddr            = "0.0.0.0:18123"
	// vaultJWTTestTenantID is the tenant_id / urn:chainlink:tenant_id claim for Vault JWT tests and
	// matches the org_id passed to DeriveJWTAuthorizedVaultWorkflowOwner.
	vaultJWTTestTenantID uint64 = 1
)

// mustDeriveJWTVaultWorkflowOwner returns the JWT-derived vault workflow owner address for an org_id.
func mustDeriveJWTVaultWorkflowOwner(t *testing.T, orgID string) string {
	t.Helper()
	derived, err := vaultjwt.DeriveJWTAuthorizedVaultWorkflowOwner(orgID, vaultJWTTestTenantID, "")
	require.NoError(t, err)
	return derived
}

func FetchVaultPublicKey(t *testing.T, gatewayURL string) (publicKey string) {
	framework.L.Info().Msg("Fetching Vault Public Key...")

	uniqueRequestID := uuid.New().String()

	getPublicKeyRequest := jsonrpc.Request[vault_helpers.GetPublicKeyRequest]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      uniqueRequestID,
		Method:  vaulttypes.MethodPublicKeyGet,
		Params:  &vault_helpers.GetPublicKeyRequest{},
	}
	requestBody, err := json.Marshal(getPublicKeyRequest)
	require.NoError(t, err, "failed to marshal public key request")

	require.Eventually(t, func() bool {
		statusCode, _ := sendVaultRequestToGateway(t, gatewayURL, requestBody)
		return statusCode == http.StatusOK
	}, time.Second*120, time.Second*5)
	statusCode, httpResponseBody := sendVaultRequestToGateway(t, gatewayURL, requestBody)
	require.Equal(t, http.StatusOK, statusCode, "Gateway endpoint should respond with 200 OK")

	framework.L.Info().Msg("Checking jsonResponse structure...")
	var jsonResponse jsonrpc.Response[vault_helpers.GetPublicKeyResponse]
	err = json.Unmarshal(httpResponseBody, &jsonResponse)
	require.NoError(t, err, "failed to unmarshal GetPublicKeyResponse")
	framework.L.Info().Msgf("JSON Body: %v", jsonResponse)
	if jsonResponse.Error != nil {
		require.Empty(t, jsonResponse.Error.Error())
	}
	require.Equal(t, jsonrpc.JsonRpcVersion, jsonResponse.Version)
	require.Equal(t, uniqueRequestID, jsonResponse.ID)
	require.Equal(t, vaulttypes.MethodPublicKeyGet, jsonResponse.Method)

	publicKeyResponse := jsonResponse.Result
	framework.L.Info().Msgf("Public Key: %s", publicKeyResponse.PublicKey)
	return publicKeyResponse.PublicKey
}

func mustVaultPublicKey(t *testing.T, publicKey string) *tdh2easy.PublicKey {
	t.Helper()

	publicKeyBytes, err := hex.DecodeString(publicKey)
	require.NoError(t, err, "failed to decode vault public key")

	parsed := &tdh2easy.PublicKey{}
	err = parsed.Unmarshal(publicKeyBytes)
	require.NoError(t, err, "failed to unmarshal vault public key")

	return parsed
}

func sendVaultRequestToGateway(t *testing.T, gatewayURL string, requestBody []byte) (statusCode int, body []byte) {
	return sendVaultRequestToGatewayWithHeaders(t, gatewayURL, requestBody, nil)
}

func sendVaultRequestToGatewayWithHeaders(t *testing.T, gatewayURL string, requestBody []byte, headers map[string]string) (statusCode int, body []byte) {
	const maxRetries = 7
	const retryInterval = 2 * time.Second

	framework.L.Info().Msgf("Request Body: %s", string(requestBody))

	for attempt := range maxRetries + 1 {
		req, err := http.NewRequestWithContext(t.Context(), "POST", gatewayURL, bytes.NewBuffer(requestBody))
		require.NoError(t, err, "failed to create request")

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		for key, value := range headers {
			req.Header.Set(key, value)
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err, "failed to execute request")

		body, err = io.ReadAll(resp.Body)
		resp.Body.Close()
		require.NoError(t, err, "failed to read http response body")
		statusCode = resp.StatusCode

		framework.L.Info().Msgf("HTTP Response Body: %s", string(body))

		if !shouldRetryGatewayRequest(statusCode, body) {
			return statusCode, body
		}

		if attempt < maxRetries {
			framework.L.Warn().Msgf("Gateway request retryable (status=%d), retrying in %s (attempt %d/%d)...", statusCode, retryInterval, attempt+1, maxRetries)
			time.Sleep(retryInterval)
		}
	}

	return statusCode, body
}

func shouldRetryGatewayRequest(statusCode int, body []byte) bool {
	if isGatewayNotAllowlistedError(body) {
		return true
	}
	switch statusCode {
	case http.StatusServiceUnavailable, http.StatusBadGateway, http.StatusGatewayTimeout:
		// Gateway-to-DON timeout: the gateway gave up relaying the response, but the DON likely
		// already processed the request. Retrying the same request body just hits the vault
		// replay guard ("request was already authorized previously"). Don't retry these.
		if bytes.Contains(body, []byte("Request timed out")) {
			return false
		}
		return true
	default:
		return false
	}
}

// isGatewayNotAllowlistedError checks whether the response is a gateway-level
// "request not allowlisted" rejection (method is empty, error code -32600).
// Node-level rejections (method is set, code -32603) have a different format
// and must not be retried because the gateway has already consumed the request.
func isGatewayNotAllowlistedError(body []byte) bool {
	var resp struct {
		Method string `json:"method"`
		Error  *struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if json.Unmarshal(body, &resp) != nil {
		return false
	}
	return resp.Method == "" && resp.Error != nil &&
		strings.Contains(resp.Error.Message, "request not allowlisted")
}

type vaultScenarioFixture struct {
	TestEnv        *ttypes.TestEnvironment
	Issuer         *stvault.TestJWTIssuer
	LinkingService *stvault.TestLinkingService
	GatewayURL     *url.URL
	VaultPublicKey string
}

type vaultWorkflowCheck struct {
	Name            string
	SecretKey       string
	SecretNamespace string
	ExpectedValue   string
	ExpectNotFound  bool
}

type vaultWorkflowPhase struct {
	Name   string
	Checks []vaultWorkflowCheck
}

type vaultRequestAuth struct {
	requestOwner string
	authorize    func(t *testing.T, req *jsonrpc.Request[json.RawMessage])
}

func getVaultJWTAuthEnabledTestConfig(t *testing.T) *ttypes.TestConfig {
	t.Helper()

	return t_helpers.GetTestConfig(t, vaultJWTAuthEnabledConfigPath)
}

func getVaultDefaultTestConfig(t *testing.T) *ttypes.TestConfig {
	t.Helper()

	return t_helpers.GetTestConfig(t, vaultDefaultConfigPath)
}

func getVaultOptimizationsEnabledTestConfig(t *testing.T) *ttypes.TestConfig {
	t.Helper()

	return t_helpers.GetTestConfig(t, vaultOptimizationsEnabledConfigPath)
}

func isVaultJWTAuthEnabledTopology(topologyName string) bool {
	return strings.Contains(topologyName, "vault-jwt_auth-enabled")
}

func isVaultOptimizationsEnabledTopology(topologyName string) bool {
	return strings.Contains(topologyName, "vault-optimizations-enabled")
}

func setupVaultScenarioFixture(t *testing.T, baseConfig *ttypes.TestConfig, usePerTestKeys bool) *vaultScenarioFixture {
	t.Helper()

	issuer, err := stvault.NewTestJWTIssuerOnAddr(vaultJWTIssuerListenAddr)
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, issuer.Close())
	})

	linkingService, err := stvault.EnsureSharedTestLinkingServiceStarted()
	require.NoError(t, err)

	var testEnv *ttypes.TestEnvironment
	if usePerTestKeys {
		testEnv = t_helpers.SetupTestEnvironmentWithPerTestKeys(t, baseConfig)
	} else {
		testEnv = t_helpers.SetupTestEnvironmentWithConfig(t, baseConfig)
	}

	ensureVaultDKGResultPackages(t, testEnv)
	gatewayURL := mustVaultGatewayURL(t, testEnv)
	vaultPublicKey := FetchVaultPublicKey(t, gatewayURL.String())
	updateVaultCapabilityConfigInRegistry(t, testEnv, vaultPublicKey)

	return &vaultScenarioFixture{
		TestEnv:        testEnv,
		Issuer:         issuer,
		LinkingService: linkingService,
		GatewayURL:     gatewayURL,
		VaultPublicKey: vaultPublicKey,
	}
}

func setupVaultSharedScenarioFixture(t *testing.T, baseConfig *ttypes.TestConfig) *vaultScenarioFixture {
	t.Helper()

	return setupVaultScenarioFixture(t, baseConfig, false)
}

func ensureVaultDKGResultPackages(t *testing.T, testEnv *ttypes.TestEnvironment) {
	t.Helper()

	framework.L.Info().Msg("Ensuring DKG result packages are present...")
	require.Eventually(t, func() bool {
		for _, nodeSet := range testEnv.Config.NodeSets {
			if slices.Contains(nodeSet.Capabilities, cre.VaultCapability) {
				for i, node := range nodeSet.NodeSpecs {
					if !slices.Contains(node.Roles, cre.BootstrapNode) {
						packageCount, err := stvault.GetResultPackageCount(t.Context(), i, nodeSet.DbInput.Port)
						if err != nil || packageCount != 1 {
							return false
						}
					}
				}
				return true
			}
		}
		return false
	}, time.Second*300, time.Second*5)
}

func requireVaultLinkOwner(t *testing.T, sc *seth.Client, workflowRegistryAddr common.Address, version *semver.Version) {
	t.Helper()

	err := creworkflow.LinkOwner(sc, workflowRegistryAddr, version)
	if err != nil && !strings.Contains(err.Error(), "OwnershipLinkAlreadyExists") {
		require.NoError(t, err)
	}
}

func mustVaultGatewayURL(t *testing.T, testEnv *ttypes.TestEnvironment) *url.URL {
	t.Helper()

	framework.L.Info().Msg("Getting gateway configuration...")
	require.NotEmpty(t, testEnv.Dons.GatewayConnectors.Configurations, "expected at least one gateway configuration")
	gatewayURL, err := url.Parse(testEnv.Dons.GatewayConnectors.Configurations[0].Incoming.Protocol + "://" + testEnv.Dons.GatewayConnectors.Configurations[0].Incoming.Host + ":" + strconv.Itoa(testEnv.Dons.GatewayConnectors.Configurations[0].Incoming.ExternalPort) + testEnv.Dons.GatewayConnectors.Configurations[0].Incoming.Path)
	require.NoError(t, err, "failed to parse gateway URL")
	framework.L.Info().Msgf("Gateway URL: %s", gatewayURL.String())
	return gatewayURL
}

func newAllowlistVaultRequestAuth(requestOwner string, sethClient *seth.Client, wfRegistryContract *workflow_registry_v2_wrapper.WorkflowRegistry) vaultRequestAuth {
	return vaultRequestAuth{
		requestOwner: requestOwner,
		authorize: func(t *testing.T, req *jsonrpc.Request[json.RawMessage]) {
			allowlistRequest(t, requestOwner, *req, sethClient, wfRegistryContract)
		},
	}
}

// newJWTVaultRequestAuth builds auth for Vault JWT requests. derivedWorkflowOwner must equal
// mustDeriveJWTVaultWorkflowOwner(orgID) — it is used as SecretIdentifier.Owner for all operations.
func newJWTVaultRequestAuth(issuer *stvault.TestJWTIssuer, orgID, derivedWorkflowOwner string) vaultRequestAuth {
	return vaultRequestAuth{
		requestOwner: derivedWorkflowOwner,
		authorize: func(t *testing.T, req *jsonrpc.Request[json.RawMessage]) {
			req.Auth = mustMintVaultJWTForRequest(t, issuer, *req, orgID)
		},
	}
}

func (a vaultRequestAuth) apply(t *testing.T, req *jsonrpc.Request[json.RawMessage]) {
	t.Helper()
	if a.authorize != nil {
		a.authorize(t, req)
	}
}

func newVaultJSONRequest(t *testing.T, requestID, method string, params any) jsonrpc.Request[json.RawMessage] {
	t.Helper()

	requestBody, err := json.Marshal(params)
	require.NoError(t, err, "failed to marshal secrets request")
	requestBodyJSON := json.RawMessage(requestBody)

	return jsonrpc.Request[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      requestID,
		Method:  method,
		Params:  &requestBodyJSON,
	}
}

func buildEncryptedSecrets(secretID, owner, encryptedSecret string, namespaces []string) []*vault_helpers.EncryptedSecret {
	encryptedSecrets := make([]*vault_helpers.EncryptedSecret, 0, len(namespaces))
	for _, namespace := range namespaces {
		encryptedSecrets = append(encryptedSecrets, &vault_helpers.EncryptedSecret{
			Id: &vault_helpers.SecretIdentifier{
				Key:       secretID,
				Owner:     owner,
				Namespace: namespace,
			},
			EncryptedValue: encryptedSecret,
		})
	}

	return encryptedSecrets
}

func buildSecretIdentifiers(secretID, owner string, namespaces []string) []*vault_helpers.SecretIdentifier {
	identifiers := make([]*vault_helpers.SecretIdentifier, 0, len(namespaces))
	for _, namespace := range namespaces {
		identifiers = append(identifiers, &vault_helpers.SecretIdentifier{
			Key:       secretID,
			Owner:     owner,
			Namespace: namespace,
		})
	}

	return identifiers
}

func sendVaultSignedOCRRequestToGateway(t *testing.T, gatewayURL string, jsonRequest jsonrpc.Request[json.RawMessage]) jsonrpc.Response[vaulttypes.SignedOCRResponse] {
	t.Helper()

	authToken := jsonRequest.Auth
	jsonRequest = outboundRequestWithoutAuth(jsonRequest)

	requestBody, err := json.Marshal(jsonRequest)
	require.NoError(t, err, "failed to marshal vault request")

	headers := map[string]string{}
	if authToken != "" {
		headers["Authorization"] = "Bearer " + authToken
	}

	statusCode, httpResponseBody := sendVaultRequestToGatewayWithHeaders(t, gatewayURL, requestBody, headers)
	require.Equal(t, http.StatusOK, statusCode, "Gateway endpoint should respond with 200 OK")

	var jsonResponse jsonrpc.Response[vaulttypes.SignedOCRResponse]
	err = json.Unmarshal(httpResponseBody, &jsonResponse)
	require.NoError(t, err, "failed to unmarshal gateway response")
	if jsonResponse.Error != nil {
		require.Empty(t, jsonResponse.Error.Error())
	}

	require.Equal(t, jsonrpc.JsonRpcVersion, jsonResponse.Version)

	return jsonResponse
}

func executeVaultSecretsCreateWithAuth(t *testing.T, auth vaultRequestAuth, encryptedSecret, secretID, expectedResponseOwner, gatewayURL string, namespaces []string) {
	t.Helper()

	executeVaultSecretsCreateWithAuthExpectOwners(t, auth, encryptedSecret, secretID, []string{expectedResponseOwner}, gatewayURL, namespaces)
}

func executeVaultSecretsCreateWithAuthExpectOwners(t *testing.T, auth vaultRequestAuth, encryptedSecret, secretID string, expectedResponseOwners []string, gatewayURL string, namespaces []string) string {
	t.Helper()

	return executeVaultSecretsCreateWithAuthExpectOwnersAndIdentifierOwner(t, auth, auth.requestOwner, encryptedSecret, secretID, expectedResponseOwners, gatewayURL, namespaces)
}

func executeVaultSecretsCreateWithAuthExpectOwnersAndIdentifierOwner(t *testing.T, auth vaultRequestAuth, identifierOwner, encryptedSecret, secretID string, expectedResponseOwners []string, gatewayURL string, namespaces []string) string {
	t.Helper()

	framework.L.Info().Msgf("Creating secrets (namespaces=%v)...", namespaces)
	require.NotEmpty(t, expectedResponseOwners, "expected response owners must not be empty")

	uniqueRequestID := uuid.New().String()
	secretsCreateRequest := vault_helpers.CreateSecretsRequest{
		RequestId:        uniqueRequestID,
		EncryptedSecrets: buildEncryptedSecrets(secretID, identifierOwner, encryptedSecret, namespaces),
	}
	jsonRequest := newVaultJSONRequest(t, uniqueRequestID, vaulttypes.MethodSecretsCreate, &secretsCreateRequest)
	auth.apply(t, &jsonRequest)

	jsonResponse := sendVaultSignedOCRRequestToGateway(t, gatewayURL, jsonRequest)
	require.Equal(t, uniqueRequestID, jsonResponse.ID)
	require.Equal(t, vaulttypes.MethodSecretsCreate, jsonResponse.Method)

	createSecretsResponse := vault_helpers.CreateSecretsResponse{}
	err := protojson.Unmarshal(jsonResponse.Result.Payload, &createSecretsResponse)
	require.NoError(t, err, "failed to decode payload into CreateSecretsResponse proto")

	require.Len(t, createSecretsResponse.Responses, len(namespaces), "Expected one item in the response per namespace")
	respByNs := make(map[string]*vault_helpers.CreateSecretResponse, len(namespaces))
	for _, r := range createSecretsResponse.GetResponses() {
		respByNs[r.GetId().GetNamespace()] = r
	}
	actualResponseOwner := ""
	for _, namespace := range namespaces {
		result, ok := respByNs[namespace]
		require.True(t, ok, "missing response for namespace %s", namespace)
		require.Empty(t, result.GetError())
		require.Equal(t, secretID, result.GetId().Key)
		require.Contains(t, expectedResponseOwners, result.GetId().Owner)
		if actualResponseOwner == "" {
			actualResponseOwner = result.GetId().Owner
			continue
		}
		require.Equal(t, actualResponseOwner, result.GetId().Owner)
	}

	return actualResponseOwner
}

func executeVaultSecretsUpdateWithAuth(t *testing.T, auth vaultRequestAuth, encryptedSecret, secretID, expectedResponseOwner, gatewayURL string, namespaces []string) {
	t.Helper()

	executeVaultSecretsUpdateWithAuthAndIdentifierOwner(t, auth, auth.requestOwner, encryptedSecret, secretID, expectedResponseOwner, gatewayURL, namespaces)
}

// tryJWTSignedVaultSecretsUpdate sends a bearer JWT Signed-OCR update without requiring HTTP success.
// The gateway may refuse cross-principal operations before quorum; workflows assert the authoritative row is unchanged.
func tryJWTSignedVaultSecretsUpdate(t *testing.T, jwtAuth vaultRequestAuth, identifierOwner, encryptedSecret, secretID, gatewayURL string, namespaces []string) {
	t.Helper()

	require.NotEmpty(t, namespaces)

	encryptedSecrets := buildEncryptedSecrets(secretID, identifierOwner, encryptedSecret, namespaces)
	uniqueRequestID := uuid.New().String()
	secretsUpdateRequest := vault_helpers.UpdateSecretsRequest{
		RequestId:        uniqueRequestID,
		EncryptedSecrets: encryptedSecrets,
	}
	jsonRequest := newVaultJSONRequest(t, uniqueRequestID, vaulttypes.MethodSecretsUpdate, &secretsUpdateRequest)
	jwtAuth.apply(t, &jsonRequest)

	authToken := jsonRequest.Auth
	outboundReq := outboundRequestWithoutAuth(jsonRequest)

	requestBody, err := json.Marshal(outboundReq)
	require.NoError(t, err)

	headers := map[string]string{}
	if authToken != "" {
		headers["Authorization"] = "Bearer " + authToken
	}
	statusCode, body := sendVaultRequestToGatewayWithHeaders(t, gatewayURL, requestBody, headers)
	framework.L.Info().Msgf("tryJWTSignedVaultSecretsUpdate HTTP status=%d body=%s", statusCode, string(body))
}

func executeVaultSecretsUpdateWithAuthAndIdentifierOwner(t *testing.T, auth vaultRequestAuth, identifierOwner, encryptedSecret, secretID, expectedResponseOwner, gatewayURL string, namespaces []string) {
	t.Helper()

	framework.L.Info().Msgf("Updating secrets (namespaces=%v)...", namespaces)
	require.NotEmpty(t, namespaces, "namespaces must not be empty")

	encryptedSecrets := buildEncryptedSecrets(secretID, identifierOwner, encryptedSecret, namespaces)
	encryptedSecrets = append(encryptedSecrets, &vault_helpers.EncryptedSecret{
		Id: &vault_helpers.SecretIdentifier{
			Key:       "invalid",
			Owner:     identifierOwner,
			Namespace: namespaces[0],
		},
		EncryptedValue: encryptedSecret,
	})

	uniqueRequestID := uuid.New().String()
	secretsUpdateRequest := vault_helpers.UpdateSecretsRequest{
		RequestId:        uniqueRequestID,
		EncryptedSecrets: encryptedSecrets,
	}
	jsonRequest := newVaultJSONRequest(t, uniqueRequestID, vaulttypes.MethodSecretsUpdate, &secretsUpdateRequest)
	auth.apply(t, &jsonRequest)

	jsonResponse := sendVaultSignedOCRRequestToGateway(t, gatewayURL, jsonRequest)
	require.Equal(t, uniqueRequestID, jsonResponse.ID)
	require.Equal(t, vaulttypes.MethodSecretsUpdate, jsonResponse.Method)

	updateSecretsResponse := vault_helpers.UpdateSecretsResponse{}
	err := protojson.Unmarshal(jsonResponse.Result.Payload, &updateSecretsResponse)
	require.NoError(t, err, "failed to decode payload into UpdateSecretsResponse proto")

	require.Len(t, updateSecretsResponse.Responses, len(namespaces)+1, "Expected one updated item per namespace plus one invalid item")
	var foundInvalid bool
	updateRespByNs := make(map[string]*vault_helpers.UpdateSecretResponse, len(namespaces))
	for _, r := range updateSecretsResponse.GetResponses() {
		if r.GetId().GetKey() == "invalid" {
			require.Contains(t, r.Error, "key does not exist")
			foundInvalid = true
			continue
		}
		updateRespByNs[r.GetId().GetNamespace()] = r
	}
	require.True(t, foundInvalid, "expected an error response for the 'invalid' key")
	for _, namespace := range namespaces {
		result, ok := updateRespByNs[namespace]
		require.True(t, ok, "missing update response for namespace %s", namespace)
		require.Empty(t, result.GetError())
		require.Equal(t, secretID, result.GetId().Key)
		require.Equal(t, expectedResponseOwner, result.GetId().Owner)
	}
}

func executeVaultSecretsListWithAuth(t *testing.T, auth vaultRequestAuth, expectedKeys []string, expectedOwner, gatewayURL, namespace string) {
	t.Helper()

	executeVaultSecretsListWithAuthAndOwner(t, auth, auth.requestOwner, expectedKeys, expectedOwner, gatewayURL, namespace)
}

func executeVaultSecretsListWithAuthAndOwner(t *testing.T, auth vaultRequestAuth, requestOwner string, expectedKeys []string, expectedOwner, gatewayURL, namespace string) {
	t.Helper()

	framework.L.Info().Msgf("Listing secrets (namespace=%s)...", namespace)

	uniqueRequestID := uuid.New().String()
	secretsListRequest := vault_helpers.ListSecretIdentifiersRequest{
		RequestId: uniqueRequestID,
		Owner:     requestOwner,
		Namespace: namespace,
	}
	jsonRequest := newVaultJSONRequest(t, uniqueRequestID, vaulttypes.MethodSecretsList, &secretsListRequest)
	auth.apply(t, &jsonRequest)

	jsonResponse := sendVaultSignedOCRRequestToGateway(t, gatewayURL, jsonRequest)
	require.Equal(t, uniqueRequestID, jsonResponse.ID)
	require.Equal(t, vaulttypes.MethodSecretsList, jsonResponse.Method)

	listSecretsResponse := vault_helpers.ListSecretIdentifiersResponse{}
	err := protojson.Unmarshal(jsonResponse.Result.Payload, &listSecretsResponse)
	require.NoError(t, err, "failed to decode payload into ListSecretIdentifiersResponse proto")

	require.True(t, listSecretsResponse.Success, err)
	require.GreaterOrEqual(t, len(listSecretsResponse.Identifiers), len(expectedKeys), "Expected enough identifiers in the response")
	keys := make([]string, 0, len(listSecretsResponse.Identifiers))
	for _, identifier := range listSecretsResponse.Identifiers {
		keys = append(keys, identifier.Key)
		require.Equal(t, expectedOwner, identifier.Owner)
		require.Equal(t, namespace, identifier.Namespace)
	}
	for _, secretID := range expectedKeys {
		require.Contains(t, keys, secretID)
	}
}

func executeVaultJWTSecretsListAbsentFromNamespace(t *testing.T, issuer *stvault.TestJWTIssuer, absentKey, orgID, derivedWorkflowOwner, gatewayURL, namespace string) {
	t.Helper()

	auth := newJWTVaultRequestAuth(issuer, orgID, derivedWorkflowOwner)

	framework.L.Info().Msgf("Listing secrets expecting key %q absent (namespace=%s)...", absentKey, namespace)

	uniqueRequestID := uuid.New().String()
	secretsListRequest := vault_helpers.ListSecretIdentifiersRequest{
		RequestId: uniqueRequestID,
		Owner:     auth.requestOwner,
		Namespace: namespace,
	}
	jsonRequest := newVaultJSONRequest(t, uniqueRequestID, vaulttypes.MethodSecretsList, &secretsListRequest)
	auth.apply(t, &jsonRequest)

	jsonResponse := sendVaultSignedOCRRequestToGateway(t, gatewayURL, jsonRequest)
	require.Equal(t, uniqueRequestID, jsonResponse.ID)
	require.Equal(t, vaulttypes.MethodSecretsList, jsonResponse.Method)

	listSecretsResponse := vault_helpers.ListSecretIdentifiersResponse{}
	err := protojson.Unmarshal(jsonResponse.Result.Payload, &listSecretsResponse)
	require.NoError(t, err, "failed to decode payload into ListSecretIdentifiersResponse proto")

	require.True(t, listSecretsResponse.Success)
	keys := make([]string, 0, len(listSecretsResponse.Identifiers))
	for _, identifier := range listSecretsResponse.Identifiers {
		require.Equal(t, derivedWorkflowOwner, identifier.Owner)
		require.Equal(t, namespace, identifier.Namespace)
		keys = append(keys, identifier.Key)
	}
	require.NotContains(t, keys, absentKey, "JWT-authenticated SecretsList must not contain deleted key")
}

func executeVaultSecretsDeleteWithAuth(t *testing.T, auth vaultRequestAuth, secretID, expectedResponseOwner, gatewayURL string, namespaces []string) {
	t.Helper()

	executeVaultSecretsDeleteWithAuthAndIdentifierOwner(t, auth, auth.requestOwner, secretID, expectedResponseOwner, gatewayURL, namespaces)
}

func executeVaultSecretsDeleteWithAuthAndIdentifierOwner(t *testing.T, auth vaultRequestAuth, identifierOwner, secretID, expectedResponseOwner, gatewayURL string, namespaces []string) {
	t.Helper()

	framework.L.Info().Msgf("Deleting secrets (namespaces=%v)...", namespaces)
	require.NotEmpty(t, namespaces, "namespaces must not be empty")

	deleteIDs := buildSecretIdentifiers(secretID, identifierOwner, namespaces)
	deleteIDs = append(deleteIDs, &vault_helpers.SecretIdentifier{
		Key:       "invalid",
		Owner:     identifierOwner,
		Namespace: namespaces[0],
	})

	uniqueRequestID := uuid.New().String()
	secretsDeleteRequest := vault_helpers.DeleteSecretsRequest{
		RequestId: uniqueRequestID,
		Ids:       deleteIDs,
	}
	jsonRequest := newVaultJSONRequest(t, uniqueRequestID, vaulttypes.MethodSecretsDelete, &secretsDeleteRequest)
	auth.apply(t, &jsonRequest)

	jsonResponse := sendVaultSignedOCRRequestToGateway(t, gatewayURL, jsonRequest)
	require.Equal(t, uniqueRequestID, jsonResponse.ID)
	require.Equal(t, vaulttypes.MethodSecretsDelete, jsonResponse.Method)

	deleteSecretsResponse := vault_helpers.DeleteSecretsResponse{}
	err := protojson.Unmarshal(jsonResponse.Result.Payload, &deleteSecretsResponse)
	require.NoError(t, err, "failed to decode payload into DeleteSecretResponse proto")

	require.Len(t, deleteSecretsResponse.Responses, len(namespaces)+1, "Expected one deleted item per namespace plus one invalid item")
	var foundDeleteInvalid bool
	deleteRespByNs := make(map[string]*vault_helpers.DeleteSecretResponse, len(namespaces))
	for _, r := range deleteSecretsResponse.GetResponses() {
		if r.GetId().GetKey() == "invalid" {
			require.Contains(t, r.Error, "key does not exist")
			foundDeleteInvalid = true
			continue
		}
		deleteRespByNs[r.GetId().GetNamespace()] = r
	}
	require.True(t, foundDeleteInvalid, "expected an error response for the 'invalid' key")
	for _, namespace := range namespaces {
		result, ok := deleteRespByNs[namespace]
		require.True(t, ok, "missing delete response for namespace %s", namespace)
		require.True(t, result.Success, result.Error)
		require.Equal(t, expectedResponseOwner, result.Id.Owner)
		require.Equal(t, secretID, result.Id.Key)
	}
}

func executeVaultAllowListSecretsCreateTest(t *testing.T, encryptedSecret, secretID, requestOwner, expectedResponseOwner, gatewayURL string, namespaces []string, sethClient *seth.Client, wfRegistryContract *workflow_registry_v2_wrapper.WorkflowRegistry) {
	t.Helper()

	auth := newAllowlistVaultRequestAuth(requestOwner, sethClient, wfRegistryContract)
	executeVaultSecretsCreateWithAuth(t, auth, encryptedSecret, secretID, expectedResponseOwner, gatewayURL, namespaces)
}

func executeVaultJWTSecretsCreateTest(t *testing.T, issuer *stvault.TestJWTIssuer, encryptedSecret, secretID, orgID, derivedWorkflowOwner, gatewayURL string, namespaces []string) {
	t.Helper()

	auth := newJWTVaultRequestAuth(issuer, orgID, derivedWorkflowOwner)
	executeVaultSecretsCreateWithAuth(t, auth, encryptedSecret, secretID, derivedWorkflowOwner, gatewayURL, namespaces)
}

func executeVaultJWTSecretsListTest(t *testing.T, issuer *stvault.TestJWTIssuer, secretID, orgID, derivedWorkflowOwner, gatewayURL, namespace string) {
	t.Helper()

	auth := newJWTVaultRequestAuth(issuer, orgID, derivedWorkflowOwner)
	executeVaultSecretsListWithAuth(t, auth, []string{secretID}, derivedWorkflowOwner, gatewayURL, namespace)
}

func executeVaultJWTSecretsDeleteTest(t *testing.T, issuer *stvault.TestJWTIssuer, secretID, orgID, derivedWorkflowOwner, gatewayURL string, namespaces []string) {
	t.Helper()

	auth := newJWTVaultRequestAuth(issuer, orgID, derivedWorkflowOwner)
	executeVaultSecretsDeleteWithAuth(t, auth, secretID, derivedWorkflowOwner, gatewayURL, namespaces)
}

func mustMintVaultJWTForRequest(t *testing.T, issuer *stvault.TestJWTIssuer, req jsonrpc.Request[json.RawMessage], orgID string) string {
	t.Helper()
	return mustMintVaultJWTForRequestWithExtraClaims(t, issuer, req, orgID, nil)
}

func mustMintVaultJWTForRequestWithExtraClaims(t *testing.T, issuer *stvault.TestJWTIssuer, req jsonrpc.Request[json.RawMessage], orgID string, extraClaims map[string]any) string {
	t.Helper()

	outboundReq := outboundRequestWithoutAuth(req)
	requestDigest, err := outboundReq.Digest()
	require.NoError(t, err, "failed to compute request digest")

	oauthScope, err := vaultjwt.OAuthScopeForVaultRPCMethod(req.Method)
	require.NoError(t, err, "resolve OAuth scope for Vault method")

	merged := map[string]any{
		vaultjwt.ClaimChainlinkTenantID: strconv.FormatUint(vaultJWTTestTenantID, 10),
	}
	maps.Copy(merged, extraClaims)

	token, err := issuer.MintToken(stvault.JWTTokenClaims{
		KeyID:         stvault.DefaultJWTIssuerKeyID,
		Issuer:        issuer.DockerIssuerURL(),
		Audience:      stvault.DefaultJWTAudience,
		OrgID:         orgID,
		WorkflowOwner: "",
		RequestDigest: requestDigest,
		Scopes:        []string{oauthScope},
		ExtraClaims:   merged,
	})
	require.NoError(t, err, "failed to mint JWT")

	return token
}

func sendVaultJWTRequestToGatewayExpectError(t *testing.T, gatewayURL string, jsonRequest jsonrpc.Request[json.RawMessage], wantStatus int) jsonrpc.Response[json.RawMessage] {
	t.Helper()

	authToken := jsonRequest.Auth
	jsonRequest = outboundRequestWithoutAuth(jsonRequest)

	requestBody, err := json.Marshal(jsonRequest)
	require.NoError(t, err, "failed to marshal JWT-authenticated request")

	headers := map[string]string{}
	if authToken != "" {
		headers["Authorization"] = "Bearer " + authToken
	}

	statusCode, httpResponseBody := sendVaultRequestToGatewayWithHeaders(t, gatewayURL, requestBody, headers)
	require.Equal(t, wantStatus, statusCode, "Gateway endpoint should respond with the expected error status")

	var jsonResponse jsonrpc.Response[json.RawMessage]
	err = json.Unmarshal(httpResponseBody, &jsonResponse)
	require.NoError(t, err, "failed to unmarshal gateway error response")
	require.Equal(t, jsonrpc.JsonRpcVersion, jsonResponse.Version)

	return jsonResponse
}

func outboundRequestWithoutAuth(req jsonrpc.Request[json.RawMessage]) jsonrpc.Request[json.RawMessage] {
	req.Auth = ""
	return req
}

func executeVaultJWTSecretsCreateUnauthorizedTest(
	t *testing.T,
	issuer *stvault.TestJWTIssuer,
	vaultPublicKey, orgID, gatewayURL string,
	expectedAuthError string,
) {
	t.Helper()
	executeVaultJWTSecretsCreateUnauthorizedWithExtraClaimsTest(t, issuer, vaultPublicKey, orgID, gatewayURL, nil, expectedAuthError)
}

func executeVaultJWTSecretsCreateUnauthorizedWithExtraClaimsTest(
	t *testing.T,
	issuer *stvault.TestJWTIssuer,
	vaultPublicKey, orgID, gatewayURL string,
	extraClaims map[string]any,
	expectedAuthError string,
) {
	t.Helper()

	derived := mustDeriveJWTVaultWorkflowOwner(t, orgID)
	derivedAddr := common.HexToAddress(derived)

	secretID := strconv.Itoa(rand.Intn(10000))
	encryptedSecret, err := vaultutils.EncryptSecretWithWorkflowOwner("secret-jwt-disabled", mustVaultPublicKey(t, vaultPublicKey), derivedAddr)
	require.NoError(t, err)

	uniqueRequestID := uuid.New().String()
	secretsCreateRequest := vault_helpers.CreateSecretsRequest{
		RequestId: uniqueRequestID,
		EncryptedSecrets: []*vault_helpers.EncryptedSecret{{
			Id: &vault_helpers.SecretIdentifier{
				Key:       secretID,
				Owner:     derived,
				Namespace: "main",
			},
			EncryptedValue: encryptedSecret,
		}},
	}
	jsonRequest := newVaultJSONRequest(t, uniqueRequestID, vaulttypes.MethodSecretsCreate, &secretsCreateRequest)
	jsonRequest.Auth = mustMintVaultJWTForRequestWithExtraClaims(t, issuer, jsonRequest, orgID, extraClaims)

	jsonResponse := sendVaultJWTRequestToGatewayExpectError(t, gatewayURL, jsonRequest, http.StatusBadRequest)
	require.Equal(t, uniqueRequestID, jsonResponse.ID)
	require.Empty(t, jsonResponse.Method)
	require.NotNil(t, jsonResponse.Error)
	require.Contains(t, jsonResponse.Error.Error(), "request not authorized")
	require.Contains(t, jsonResponse.Error.Error(), expectedAuthError)
}

func startVaultSecretsWorkflowPhasesTest(
	t *testing.T, testEnv *ttypes.TestEnvironment,
	workflowBaseName string,
	phases []vaultWorkflowPhase,
) string {
	t.Helper()

	testLogger := framework.L
	testLogger.Info().
		Str("workflow_base_name", workflowBaseName).
		Int("phase_count", len(phases)).
		Msg("Starting vault workflow phase verification")

	workflowName := t_helpers.UniqueWorkflowName(testEnv, workflowBaseName)
	cfgPhases := make([]vaultsecret_config.Phase, 0, len(phases))
	for _, phase := range phases {
		cfgChecks := make([]vaultsecret_config.Check, 0, len(phase.Checks))
		for _, check := range phase.Checks {
			cfgChecks = append(cfgChecks, vaultsecret_config.Check{
				Name:            check.Name,
				SecretKey:       check.SecretKey,
				SecretNamespace: check.SecretNamespace,
				ExpectedValue:   check.ExpectedValue,
				ExpectNotFound:  check.ExpectNotFound,
			})
		}
		cfgPhases = append(cfgPhases, vaultsecret_config.Phase{
			Name:   phase.Name,
			Checks: cfgChecks,
		})
	}

	cfg := &vaultsecret_config.Config{Phases: cfgPhases}
	const workflowFileLocation = "./vaultsecret/main.go"
	return t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, workflowName, cfg, workflowFileLocation)
}

func waitForVaultWorkflowPhase(
	t *testing.T,
	workflowID, phaseName string,
	userLogsCh chan *workflowevents.UserLogs,
	baseMessageCh chan *commonevents.BaseMessage,
) {
	t.Helper()

	testLogger := framework.L
	t_helpers.WatchWorkflowLogs(
		t,
		testLogger,
		userLogsCh,
		baseMessageCh,
		t_helpers.WorkflowEngineInitErrorLog,
		"Vault secret workflow phase completed: "+phaseName,
		1*time.Minute,
		t_helpers.WithUserLogWorkflowID(workflowID),
	)
	testLogger.Info().Str("phase_name", phaseName).Msg("Vault secret workflow phase completed")
}

func executeVaultSecretsUpdateTest(t *testing.T, encryptedSecret, secretID, requestOwner, expectedResponseOwner, gatewayURL string, namespaces []string, sethClient *seth.Client, wfRegistryContract *workflow_registry_v2_wrapper.WorkflowRegistry) {
	auth := newAllowlistVaultRequestAuth(requestOwner, sethClient, wfRegistryContract)
	executeVaultSecretsUpdateWithAuth(t, auth, encryptedSecret, secretID, expectedResponseOwner, gatewayURL, namespaces)
}

func executeVaultSecretsListTest(t *testing.T, secretID, requestOwner, expectedOwner, gatewayURL, namespace string, sethClient *seth.Client, wfRegistryContract *workflow_registry_v2_wrapper.WorkflowRegistry) {
	auth := newAllowlistVaultRequestAuth(requestOwner, sethClient, wfRegistryContract)
	executeVaultSecretsListWithAuth(t, auth, []string{secretID}, expectedOwner, gatewayURL, namespace)
}

func executeVaultSecretsDeleteTest(t *testing.T, secretID, requestOwner, expectedResponseOwner, gatewayURL string, namespaces []string, sethClient *seth.Client, wfRegistryContract *workflow_registry_v2_wrapper.WorkflowRegistry) {
	auth := newAllowlistVaultRequestAuth(requestOwner, sethClient, wfRegistryContract)
	executeVaultSecretsDeleteWithAuth(t, auth, secretID, expectedResponseOwner, gatewayURL, namespaces)
}

// executeVaultBinaryEncodedSharesSmokeTest verifies a workflow can fetch a secret when the vault
// DON emits binary-encoded decryption shares (VaultOptimizationsEnabled). GetSecrets is not
// exposed on the vault gateway; binary encoding is confirmed via vault node logs after the
// workflow triggers a capability get.
func executeVaultBinaryEncodedSharesSmokeTest(
	t *testing.T,
	testEnv *ttypes.TestEnvironment,
	secretID, namespace, expectedPlaintext string,
	userLogsCh chan *workflowevents.UserLogs,
	baseMessageCh chan *commonevents.BaseMessage,
) {
	t.Helper()

	framework.L.Info().Msg("Verifying workflow secret fetch and binary share encoding when optimizations are enabled...")

	workflowID := startVaultSecretsWorkflowPhasesTest(t, testEnv, "binary-shares", []vaultWorkflowPhase{{
		Name: "binary-shares-fetch",
		Checks: []vaultWorkflowCheck{{
			Name:            "binary-shares-main",
			SecretKey:       secretID,
			SecretNamespace: namespace,
			ExpectedValue:   expectedPlaintext,
		}},
	}})
	waitForVaultWorkflowPhase(t, workflowID, "binary-shares-fetch", userLogsCh, baseMessageCh)
}

// updateVaultCapabilityConfigInRegistry updates the on-chain capabilities registry
// so that the vault@1.0.0 capability config includes DefaultConfig with VaultPublicKey
// and Threshold. This is required for workflows that call runtime.GetSecret().
// Uses the original deployer key (not per-test key) since the registry is owned by the deployer.
func updateVaultCapabilityConfigInRegistry(t *testing.T, testEnv *ttypes.TestEnvironment, vaultPublicKey string) {
	t.Helper()
	testLogger := framework.L
	testLogger.Info().Msg("Updating vault capability config in capabilities registry with VaultPublicKey...")

	capRegAddr := crecontracts.MustGetAddressFromDataStore(
		testEnv.CreEnvironment.CldfEnvironment.DataStore,
		testEnv.CreEnvironment.RegistryChainSelector,
		keystone_changeset.CapabilitiesRegistry.String(),
		testEnv.CreEnvironment.ContractVersions[keystone_changeset.CapabilitiesRegistry.String()],
		"",
	)

	require.IsType(t, &evm.Blockchain{}, testEnv.CreEnvironment.Blockchains[0])
	sethClient := testEnv.CreEnvironment.Blockchains[0].(*evm.Blockchain).SethClient

	deployerClient, err := seth.NewClientBuilder().
		WithRpcUrl(sethClient.URL).
		WithPrivateKeys([]string{ctfblockchain.DefaultAnvilPrivateKey}).
		WithProtections(false, false, seth.MustMakeDuration(time.Second)).
		Build()
	require.NoError(t, err, "failed to create deployer seth client")

	capReg, err := capabilities_registry_v2.NewCapabilitiesRegistry(
		common.HexToAddress(capRegAddr), deployerClient.Client,
	)
	require.NoError(t, err, "failed to create capabilities registry wrapper")

	allDONs, err := capReg.GetDONs(&bind.CallOpts{}, big.NewInt(0), big.NewInt(100))
	require.NoError(t, err, "failed to get DONs from registry")

	var don *capabilities_registry_v2.CapabilitiesRegistryDONInfo
	for i := range allDONs {
		for _, cc := range allDONs[i].CapabilityConfigurations {
			if cc.CapabilityId == "vault@1.0.0" {
				don = &allDONs[i]
				break
			}
		}
		if don != nil {
			break
		}
	}
	require.NotNil(t, don, "could not find a DON with vault@1.0.0 capability in the registry")
	testLogger.Info().Msgf("Found vault capability on DON %q (ID=%d)", don.Name, don.Id)

	newConfigs := make([]capabilities_registry_v2.CapabilitiesRegistryCapabilityConfiguration, 0, len(don.CapabilityConfigurations))
	for _, cc := range don.CapabilityConfigurations {
		if cc.CapabilityId == "vault@1.0.0" {
			existingConfig := &capabilitiespb.CapabilityConfig{}
			if len(cc.Config) > 0 {
				require.NoError(t, proto.Unmarshal(cc.Config, existingConfig), "failed to unmarshal existing vault capability config")
			}

			vaultCfg := map[string]any{
				"VaultPublicKey": vaultPublicKey,
				"Threshold":      1,
			}
			valueMap, wrapErr := values.WrapMap(vaultCfg)
			require.NoError(t, wrapErr, "failed to wrap vault config values")

			existingConfig.DefaultConfig = values.ProtoMap(valueMap)

			configBytes, marshalErr := proto.Marshal(existingConfig)
			require.NoError(t, marshalErr, "failed to marshal updated vault capability config")

			cc.Config = configBytes
			testLogger.Info().Msg("Injected VaultPublicKey and Threshold into vault@1.0.0 capability config")
		}
		newConfigs = append(newConfigs, cc)
	}

	updateParams := capabilities_registry_v2.CapabilitiesRegistryUpdateDONParams{
		Name:                     don.Name,
		Nodes:                    don.NodeP2PIds,
		CapabilityConfigurations: newConfigs,
		IsPublic:                 don.IsPublic,
		F:                        don.F,
		Config:                   don.Config,
	}

	_, err = deployerClient.Decode(capReg.UpdateDONByName(deployerClient.NewTXOpts(), don.Name, updateParams))
	require.NoError(t, err, "UpdateDONByName tx failed")

	testLogger.Info().Msg("Waiting for registry syncer to propagate the on-chain config change...")
	time.Sleep(15 * time.Second) // registry syncer polls every 12s; one tick + margin
}

func allowlistRequest(t *testing.T, owner string, request jsonrpc.Request[json.RawMessage], sethClient *seth.Client, wfRegistryContract *workflow_registry_v2_wrapper.WorkflowRegistry) {
	requestDigest, err := request.Digest()
	require.NoError(t, err, "failed to get digest for request")
	requestDigestBytes, err := hex.DecodeString(requestDigest)
	require.NoError(t, err, "failed to decode digest")
	reqDigestBytes := [32]byte(requestDigestBytes)
	_, err = wfRegistryContract.AllowlistRequest(sethClient.NewTXOpts(), reqDigestBytes, uint32(time.Now().Add(1*time.Hour).Unix())) //nolint:gosec // disable G115
	require.NoError(t, err, "failed to allowlist request")

	framework.L.Info().Msgf("Allowlisting request digest at contract %s, for owner: %s, digestHexStr: %s", wfRegistryContract.Address().Hex(), owner, requestDigest)
	allowedList, err := wfRegistryContract.GetAllowlistedRequests(&bind.CallOpts{}, big.NewInt(0), big.NewInt(100))
	require.NoError(t, err, "failed to validate allowlisted request")
	for _, req := range allowedList {
		if req.RequestDigest == reqDigestBytes {
			framework.L.Info().Msgf("Request digest found in allowlist")
		}
		framework.L.Info().Msgf("Allowlisted request digestHexStr: %s, owner: %s, expiry: %d", hex.EncodeToString(req.RequestDigest[:]), req.Owner.Hex(), req.ExpiryTimestamp)
	}
}
