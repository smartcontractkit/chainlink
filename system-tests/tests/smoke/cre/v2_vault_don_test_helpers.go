package cre

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"io"
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
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/vault"
	creworkflow "github.com/smartcontractkit/chainlink/system-tests/lib/cre/workflow"
	vaultsecret_config "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/vaultsecret/config"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaultutils"
)

const (
	vaultDefaultConfigPath        = "/configs/workflow-gateway-capabilities-don.toml"
	vaultJWTAuthEnabledConfigPath = "/configs/workflow-gateway-capabilities-don-vault-jwt_auth-enabled.toml"
	vaultJWTIssuerListenAddr      = "0.0.0.0:18123"
)

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

		if !isGatewayNotAllowlistedError(body) {
			return statusCode, body
		}

		if attempt < maxRetries {
			framework.L.Warn().Msgf("Request not yet allowlisted, retrying in %s (attempt %d/%d)...", retryInterval, attempt+1, maxRetries)
			time.Sleep(retryInterval)
		}
	}

	return statusCode, body
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
	Issuer         *vault.TestJWTIssuer
	LinkingService *vault.TestLinkingService
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

func isVaultJWTAuthEnabledTopology(topologyName string) bool {
	return strings.Contains(topologyName, "vault-jwt_auth-enabled")
}

func setupVaultScenarioFixture(t *testing.T, baseConfig *ttypes.TestConfig, usePerTestKeys bool) *vaultScenarioFixture {
	t.Helper()

	issuer, err := vault.NewTestJWTIssuerOnAddr(vaultJWTIssuerListenAddr)
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, issuer.Close())
	})

	linkingService, err := vault.EnsureSharedTestLinkingServiceStarted()
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
						packageCount, err := vault.GetResultPackageCount(t.Context(), i, nodeSet.DbInput.Port)
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

func newJWTVaultRequestAuth(issuer *vault.TestJWTIssuer, orgID, workflowOwner string) vaultRequestAuth {
	return vaultRequestAuth{
		requestOwner: orgID,
		authorize: func(t *testing.T, req *jsonrpc.Request[json.RawMessage]) {
			req.Auth = mustMintVaultJWTForRequest(t, issuer, *req, orgID, workflowOwner)
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

	framework.L.Info().Msgf("Creating secrets (namespaces=%v)...", namespaces)

	uniqueRequestID := uuid.New().String()
	secretsCreateRequest := vault_helpers.CreateSecretsRequest{
		RequestId:        uniqueRequestID,
		EncryptedSecrets: buildEncryptedSecrets(secretID, auth.requestOwner, encryptedSecret, namespaces),
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
	for _, namespace := range namespaces {
		result, ok := respByNs[namespace]
		require.True(t, ok, "missing response for namespace %s", namespace)
		require.Empty(t, result.GetError())
		require.Equal(t, secretID, result.GetId().Key)
		require.Equal(t, expectedResponseOwner, result.GetId().Owner)
	}
}

func executeVaultSecretsUpdateWithAuth(t *testing.T, auth vaultRequestAuth, encryptedSecret, secretID, expectedResponseOwner, gatewayURL string, namespaces []string) {
	t.Helper()

	framework.L.Info().Msgf("Updating secrets (namespaces=%v)...", namespaces)
	require.NotEmpty(t, namespaces, "namespaces must not be empty")

	encryptedSecrets := buildEncryptedSecrets(secretID, auth.requestOwner, encryptedSecret, namespaces)
	encryptedSecrets = append(encryptedSecrets, &vault_helpers.EncryptedSecret{
		Id: &vault_helpers.SecretIdentifier{
			Key:       "invalid",
			Owner:     auth.requestOwner,
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

	framework.L.Info().Msgf("Listing secrets (namespace=%s)...", namespace)

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

func executeVaultSecretsDeleteWithAuth(t *testing.T, auth vaultRequestAuth, secretID, expectedResponseOwner, gatewayURL string, namespaces []string) {
	t.Helper()

	framework.L.Info().Msgf("Deleting secrets (namespaces=%v)...", namespaces)
	require.NotEmpty(t, namespaces, "namespaces must not be empty")

	deleteIDs := buildSecretIdentifiers(secretID, auth.requestOwner, namespaces)
	deleteIDs = append(deleteIDs, &vault_helpers.SecretIdentifier{
		Key:       "invalid",
		Owner:     auth.requestOwner,
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
	auth := newAllowlistVaultRequestAuth(requestOwner, sethClient, wfRegistryContract)
	executeVaultSecretsCreateWithAuth(t, auth, encryptedSecret, secretID, expectedResponseOwner, gatewayURL, namespaces)
}

func executeVaultJWTSecretsCreateTest(t *testing.T, issuer *vault.TestJWTIssuer, encryptedSecret, secretID, orgID, workflowOwner, gatewayURL string, namespaces []string) {
	t.Helper()

	auth := newJWTVaultRequestAuth(issuer, orgID, workflowOwner)
	executeVaultSecretsCreateWithAuth(t, auth, encryptedSecret, secretID, orgID, gatewayURL, namespaces)
}

func executeVaultJWTSecretsListTest(t *testing.T, issuer *vault.TestJWTIssuer, secretID, orgID, workflowOwner, gatewayURL, namespace string) {
	t.Helper()

	auth := newJWTVaultRequestAuth(issuer, orgID, workflowOwner)
	executeVaultSecretsListWithAuth(t, auth, []string{secretID}, orgID, gatewayURL, namespace)
}

func executeVaultJWTSecretsDeleteTest(t *testing.T, issuer *vault.TestJWTIssuer, secretID, orgID, workflowOwner, gatewayURL string, namespaces []string) {
	t.Helper()

	auth := newJWTVaultRequestAuth(issuer, orgID, workflowOwner)
	executeVaultSecretsDeleteWithAuth(t, auth, secretID, orgID, gatewayURL, namespaces)
}

func mustMintVaultJWTForRequest(t *testing.T, issuer *vault.TestJWTIssuer, req jsonrpc.Request[json.RawMessage], orgID, workflowOwner string) string {
	t.Helper()

	outboundReq := outboundRequestWithoutAuth(req)
	requestDigest, err := outboundReq.Digest()
	require.NoError(t, err, "failed to compute request digest")

	token, err := issuer.MintToken(vault.JWTTokenClaims{
		KeyID:         vault.DefaultJWTIssuerKeyID,
		Issuer:        issuer.DockerIssuerURL(),
		Audience:      vault.DefaultJWTAudience,
		OrgID:         orgID,
		WorkflowOwner: workflowOwner,
		RequestDigest: requestDigest,
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
	issuer *vault.TestJWTIssuer,
	vaultPublicKey, orgID, workflowOwner, gatewayURL string,
	expectedAuthError string,
) {
	t.Helper()

	secretID := strconv.Itoa(rand.Intn(10000))
	encryptedSecret, err := vaultutils.EncryptSecretWithOrgID("secret-jwt-disabled", mustVaultPublicKey(t, vaultPublicKey), orgID)
	require.NoError(t, err)

	uniqueRequestID := uuid.New().String()
	secretsCreateRequest := vault_helpers.CreateSecretsRequest{
		RequestId: uniqueRequestID,
		EncryptedSecrets: []*vault_helpers.EncryptedSecret{{
			Id: &vault_helpers.SecretIdentifier{
				Key:       secretID,
				Owner:     orgID,
				Namespace: "main",
			},
			EncryptedValue: encryptedSecret,
		}},
	}
	jsonRequest := newVaultJSONRequest(t, uniqueRequestID, vaulttypes.MethodSecretsCreate, &secretsCreateRequest)
	jsonRequest.Auth = mustMintVaultJWTForRequest(t, issuer, jsonRequest, orgID, workflowOwner)

	jsonResponse := sendVaultJWTRequestToGatewayExpectError(t, gatewayURL, jsonRequest, http.StatusBadRequest)
	require.Equal(t, uniqueRequestID, jsonResponse.ID)
	require.Empty(t, jsonResponse.Method)
	require.NotNil(t, jsonResponse.Error)
	require.Contains(t, jsonResponse.Error.Error(), "request not authorized")
	require.Contains(t, jsonResponse.Error.Error(), expectedAuthError)
}

func executeVaultSecretsWorkflowChecksTest(
	t *testing.T, testEnv *ttypes.TestEnvironment,
	workflowBaseName string,
	checks []vaultWorkflowCheck,
	userLogsCh chan *workflowevents.UserLogs, baseMessageCh chan *commonevents.BaseMessage,
) {
	t.Helper()

	workflowID := startVaultSecretsWorkflowPhasesTest(t, testEnv, workflowBaseName, []vaultWorkflowPhase{{
		Name:   workflowBaseName,
		Checks: checks,
	}})
	waitForVaultWorkflowPhase(t, workflowID, workflowBaseName, userLogsCh, baseMessageCh)
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
		4*time.Minute,
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

			vaultCfg := map[string]interface{}{
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
