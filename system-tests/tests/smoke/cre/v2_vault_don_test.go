package cre

import (
	"encoding/hex"
	"encoding/json"
	"math/big"
	"math/rand"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"google.golang.org/protobuf/encoding/protojson"

	vault_helpers "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"

	workflow_registry_v2_wrapper "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v2"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	crevault "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/vault"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/vault"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

func ExecuteVaultTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	/*
		BUILD ENVIRONMENT FROM SAVED STATE
	*/
	var testLogger = framework.L

	testLogger.Info().Msgf("Ensuring DKG result packages are present...")
	require.Eventually(t, func() bool {
		for _, nodeSet := range testEnv.Config.NodeSets {
			if slices.Contains(nodeSet.Capabilities, cre.VaultCapability) {
				for i := range nodeSet.Nodes {
					if i != nodeSet.BootstrapNodeIndex {
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

	// Wait a bit to ensure the Vault plugin is ready.
	time.Sleep(30 * time.Second)

	testLogger.Info().Msg("Getting gateway configuration...")
	require.NotEmpty(t, testEnv.CreEnvironment.DonTopology.GatewayConnectorOutput.Configurations, "expected at least one gateway configuration")
	gatewayURL, err := url.Parse(testEnv.CreEnvironment.DonTopology.GatewayConnectorOutput.Configurations[0].Incoming.Protocol + "://" + testEnv.CreEnvironment.DonTopology.GatewayConnectorOutput.Configurations[0].Incoming.Host + ":" + strconv.Itoa(testEnv.CreEnvironment.DonTopology.GatewayConnectorOutput.Configurations[0].Incoming.ExternalPort) + testEnv.CreEnvironment.DonTopology.GatewayConnectorOutput.Configurations[0].Incoming.Path)
	require.NoError(t, err, "failed to parse gateway URL")
	testLogger.Info().Msgf("Gateway URL: %s", gatewayURL.String())

	// Ignoring the deprecation warning as the suggest solution is not working in CI
	//lint:ignore SA1019 ignoring deprecation warning for this usage
	workflowRegistryAddress, _, workflowRegistryErr := crecontracts.FindAddressesForChain(
		testEnv.CreEnvironment.CldfEnvironment.ExistingAddresses, //nolint:staticcheck // SA1019 ignoring deprecation warning for this usage
		testEnv.WrappedBlockchainOutputs[0].ChainSelector, keystone_changeset.WorkflowRegistry.String())
	require.NoError(t, workflowRegistryErr, "failed to find workflow registry address for chain %d", testEnv.WrappedBlockchainOutputs[0].ChainID)

	sethClient := testEnv.WrappedBlockchainOutputs[0].SethClient
	ownerAddr := sethClient.MustGetRootKeyAddress().Hex()
	t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, "consensustest", &t_helpers.None{}, "../../../../core/scripts/cre/environment/examples/workflows/v2/node-mode/main.go")
	wfRegistryContract, err := workflow_registry_v2_wrapper.NewWorkflowRegistry(workflowRegistryAddress, sethClient.Client)
	require.NoError(t, err, "failed to get workflow registry contract wrapper")

	secretID := strconv.Itoa(rand.Intn(10000)) // generate a random secret ID for testing
	secretValue := "Secret Value to be stored"
	vaultPublicKey := FetchVaultPublicKey(t, gatewayURL.String())
	encryptedSecret, err := crevault.EncryptSecret(secretValue, vaultPublicKey)
	require.NoError(t, err, "failed to encrypt secret")

	// Wait for the node to be up.
	framework.L.Info().Msg("Waiting 30 seconds for the Vault DON to be ready...")
	time.Sleep(30 * time.Second)
	executeVaultSecretsCreateTest(t, encryptedSecret, secretID, ownerAddr, gatewayURL.String(), sethClient, wfRegistryContract)
	executeVaultSecretsGetTest(t, secretID, ownerAddr, gatewayURL.String(), sethClient, wfRegistryContract)
	executeVaultSecretsUpdateTest(t, encryptedSecret, secretID, ownerAddr, gatewayURL.String(), sethClient, wfRegistryContract)
	executeVaultSecretsListTest(t, secretID, ownerAddr, gatewayURL.String(), sethClient, wfRegistryContract)
	executeVaultSecretsDeleteTest(t, secretID, ownerAddr, gatewayURL.String(), sethClient, wfRegistryContract)
}

func executeVaultSecretsCreateTest(t *testing.T, encryptedSecret, secretID, owner, gatewayURL string, sethClient *seth.Client, wfRegistryContract *workflow_registry_v2_wrapper.WorkflowRegistry) {
	framework.L.Info().Msg("Creating secret...")

	uniqueRequestID := uuid.New().String()

	secretsCreateRequest := vault_helpers.CreateSecretsRequest{
		RequestId: uniqueRequestID,
		EncryptedSecrets: []*vault_helpers.EncryptedSecret{
			{
				Id: &vault_helpers.SecretIdentifier{
					Key:       secretID,
					Owner:     owner,
					Namespace: "main",
				},
				EncryptedValue: encryptedSecret,
			},
		},
	}
	secretsCreateRequestBody, err := json.Marshal(secretsCreateRequest) //nolint:govet // The lock field is not set on this proto
	require.NoError(t, err, "failed to marshal secrets request")
	secretsCreateRequestBodyJSON := json.RawMessage(secretsCreateRequestBody)
	jsonRequest := jsonrpc.Request[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      uniqueRequestID,
		Method:  vaulttypes.MethodSecretsCreate,
		Params:  &secretsCreateRequestBodyJSON,
	}
	allowlistRequest(t, owner, jsonRequest, sethClient, wfRegistryContract)

	requestBody, err := json.Marshal(jsonRequest)
	require.NoError(t, err, "failed to marshal secrets request")

	statusCode, httpResponseBody := sendVaultRequestToGateway(t, gatewayURL, requestBody)
	require.Equal(t, http.StatusOK, statusCode, "Gateway endpoint should respond with 200 OK")

	framework.L.Info().Msg("Checking jsonResponse structure...")
	var jsonResponse jsonrpc.Response[vaulttypes.SignedOCRResponse]
	err = json.Unmarshal(httpResponseBody, &jsonResponse)
	require.NoError(t, err, "failed to unmarshal getResponse")
	framework.L.Info().Msgf("JSON Body: %v", jsonResponse)
	if jsonResponse.Error != nil {
		require.Empty(t, jsonResponse.Error.Error())
	}
	require.Equal(t, jsonrpc.JsonRpcVersion, jsonResponse.Version)
	require.Equal(t, uniqueRequestID, jsonResponse.ID)
	require.Equal(t, vaulttypes.MethodSecretsCreate, jsonResponse.Method)

	signedOCRResponse := jsonResponse.Result
	framework.L.Info().Msgf("Signed OCR Response: %s", signedOCRResponse.String())

	// TODO: Verify the authenticity of this signed report, by ensuring that the signatures indeed match the payload
	createSecretsResponse := vault_helpers.CreateSecretsResponse{}
	err = protojson.Unmarshal(signedOCRResponse.Payload, &createSecretsResponse)
	require.NoError(t, err, "failed to decode payload into CreateSecretsResponse proto")
	framework.L.Info().Msgf("CreateSecretsResponse decoded as: %s", createSecretsResponse.String())

	require.Len(t, createSecretsResponse.Responses, 1, "Expected one item in the response")
	result0 := createSecretsResponse.GetResponses()[0]
	require.Empty(t, result0.GetError())
	require.Equal(t, secretID, result0.GetId().Key)
	require.Equal(t, owner, result0.GetId().Owner)
	require.Equal(t, vaulttypes.DefaultNamespace, result0.GetId().Namespace)

	framework.L.Info().Msg("Secret created successfully")
}

func executeVaultSecretsUpdateTest(t *testing.T, encryptedSecret, secretID, owner, gatewayURL string, sethClient *seth.Client, wfRegistryContract *workflow_registry_v2_wrapper.WorkflowRegistry) {
	framework.L.Info().Msg("Updating secret...")
	uniqueRequestID := uuid.New().String()

	secretsUpdateRequest := vault_helpers.UpdateSecretsRequest{
		RequestId: uniqueRequestID,
		EncryptedSecrets: []*vault_helpers.EncryptedSecret{
			{
				Id: &vault_helpers.SecretIdentifier{
					Key:       secretID,
					Owner:     owner,
					Namespace: "main",
				},
				EncryptedValue: encryptedSecret,
			},
			{
				Id: &vault_helpers.SecretIdentifier{
					Key:       "invalid",
					Owner:     owner,
					Namespace: "main",
				},
				EncryptedValue: encryptedSecret,
			},
		},
	}
	secretsUpdateRequestBody, err := json.Marshal(secretsUpdateRequest) //nolint:govet // The lock field is not set on this proto
	require.NoError(t, err, "failed to marshal secrets request")
	secretsUpdateRequestBodyJSON := json.RawMessage(secretsUpdateRequestBody)
	jsonRequest := jsonrpc.Request[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      uniqueRequestID,
		Method:  vaulttypes.MethodSecretsUpdate,
		Params:  &secretsUpdateRequestBodyJSON,
	}
	allowlistRequest(t, owner, jsonRequest, sethClient, wfRegistryContract)

	requestBody, err := json.Marshal(jsonRequest)
	require.NoError(t, err, "failed to marshal secrets request")

	statusCode, httpResponseBody := sendVaultRequestToGateway(t, gatewayURL, requestBody)
	require.Equal(t, http.StatusOK, statusCode, "Gateway endpoint should respond with 200 OK")

	framework.L.Info().Msg("Checking jsonResponse structure...")
	var jsonResponse jsonrpc.Response[vaulttypes.SignedOCRResponse]
	err = json.Unmarshal(httpResponseBody, &jsonResponse)
	require.NoError(t, err, "failed to unmarshal getResponse")
	framework.L.Info().Msgf("JSON Body: %v", jsonResponse)
	if jsonResponse.Error != nil {
		require.Empty(t, jsonResponse.Error.Error())
	}

	require.Equal(t, jsonrpc.JsonRpcVersion, jsonResponse.Version)
	require.Equal(t, uniqueRequestID, jsonResponse.ID)
	require.Equal(t, vaulttypes.MethodSecretsUpdate, jsonResponse.Method)

	signedOCRResponse := jsonResponse.Result
	framework.L.Info().Msgf("Signed OCR Response: %s", signedOCRResponse.String())

	// TODO: Verify the authenticity of this signed report, by ensuring that the signatures indeed match the payload

	updateSecretsResponse := vault_helpers.UpdateSecretsResponse{}
	err = protojson.Unmarshal(signedOCRResponse.Payload, &updateSecretsResponse)
	require.NoError(t, err, "failed to decode payload into UpdateSecretsResponse proto")
	framework.L.Info().Msgf("UpdateSecretsResponse decoded as: %s", updateSecretsResponse.String())

	require.Len(t, updateSecretsResponse.Responses, 2, "Expected 2 items in the response")
	result0 := updateSecretsResponse.GetResponses()[0]
	require.Empty(t, result0.GetError())
	require.Equal(t, secretID, result0.GetId().Key)
	require.Equal(t, owner, result0.GetId().Owner)
	require.Equal(t, vaulttypes.DefaultNamespace, result0.GetId().Namespace)

	result1 := updateSecretsResponse.GetResponses()[1]
	require.Contains(t, result1.Error, "key does not exist")

	framework.L.Info().Msg("Secret updated successfully")
}

func executeVaultSecretsGetTest(t *testing.T, secretID, owner, gatewayURL string, sethClient *seth.Client, wfRegistryContract *workflow_registry_v2_wrapper.WorkflowRegistry) {
	uniqueRequestID := uuid.New().String()
	framework.L.Info().Msg("Getting secret...")
	secretsGetRequest := jsonrpc.Request[vault_helpers.GetSecretsRequest]{
		Version: jsonrpc.JsonRpcVersion,
		Method:  vaulttypes.MethodSecretsGet,
		Params: &vault_helpers.GetSecretsRequest{
			Requests: []*vault_helpers.SecretRequest{
				{
					Id: &vault_helpers.SecretIdentifier{
						Key:   secretID,
						Owner: owner,
					},
				},
			},
		},
		ID: uniqueRequestID,
	}
	requestBody, err := json.Marshal(secretsGetRequest)
	require.NoError(t, err, "failed to marshal secrets request")
	statusCode, httpResponseBody := sendVaultRequestToGateway(t, gatewayURL, requestBody)
	require.Equal(t, http.StatusOK, statusCode, "Gateway endpoint should respond with 200 OK")
	framework.L.Info().Msg("Checking jsonResponse structure...")
	var jsonResponse jsonrpc.Response[json.RawMessage]
	err = json.Unmarshal(httpResponseBody, &jsonResponse)
	require.NoError(t, err, "failed to unmarshal http response body")
	framework.L.Info().Msgf("JSON Body: %v", jsonResponse)
	if jsonResponse.Error != nil {
		require.Empty(t, jsonResponse.Error.Error())
	}
	require.Equal(t, jsonrpc.JsonRpcVersion, jsonResponse.Version)
	require.Equal(t, uniqueRequestID, jsonResponse.ID)
	require.Equal(t, vaulttypes.MethodSecretsGet, jsonResponse.Method)

	/*
	 * The json unmarshaling is not compatible with the proto oneof in vault_helpers.SecretResponse
	 * The Data and Error fields are oneof fields in the proto definition, but when unmarshaling to JSON,
	 * the JSON unmarshaler does not handle oneof fields correctly, leading to issues.
	 * To work around this, we define custom response types that match the expected structure.
	 * This allows us to unmarshal the JSON response correctly and access the fields as expected.
	 */
	type EncryptedShares struct {
		Shares        []string `protobuf:"bytes,1,rep,name=shares,proto3" json:"shares,omitempty"`
		EncryptionKey string   `protobuf:"bytes,2,opt,name=encryption_key,json=encryptionKey,proto3" json:"encryption_key,omitempty"`
	}
	type SecretData struct {
		EncryptedValue               string             `protobuf:"bytes,2,opt,name=encrypted_value,json=encryptedValue,proto3" json:"encrypted_value,omitempty"`
		EncryptedDecryptionKeyShares []*EncryptedShares `protobuf:"bytes,3,rep,name=encrypted_decryption_key_shares,json=encryptedDecryptionKeyShares,proto3" json:"encrypted_decryption_key_shares,omitempty"`
	}
	type SecretResponse struct {
		ID    *vault_helpers.SecretIdentifier `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
		Data  *SecretData                     `protobuf:"bytes,2,opt,name=data,proto3"`
		Error string                          `protobuf:"bytes,3,opt,name=error,proto3"`
	}
	type GetSecretsResponse struct {
		Responses []*SecretResponse `protobuf:"bytes,1,rep,name=responses,proto3" json:"responses,omitempty"`
	}
	/*
	 *
	 *
	 *
	 *
	 */

	var getSecretsResponse GetSecretsResponse
	err = json.Unmarshal(*jsonResponse.Result, &getSecretsResponse)
	require.NoError(t, err, "failed to unmarshal getResponse")

	require.Len(t, getSecretsResponse.Responses, 1, "Expected one secret in the response")
	result0 := getSecretsResponse.Responses[0]
	require.Empty(t, result0.Error)
	require.Equal(t, secretID, result0.ID.Key)
	require.Equal(t, owner, result0.ID.Owner)
	require.Equal(t, vaulttypes.DefaultNamespace, result0.ID.Namespace)

	framework.L.Info().Msg("Secret get successful")
}

func executeVaultSecretsListTest(t *testing.T, secretID, owner, gatewayURL string, sethClient *seth.Client, wfRegistryContract *workflow_registry_v2_wrapper.WorkflowRegistry) {
	framework.L.Info().Msg("Listing secret...")
	uniqueRequestID := uuid.New().String()
	secretsListRequest := vault_helpers.ListSecretIdentifiersRequest{
		RequestId: uniqueRequestID,
		Owner:     owner,
		Namespace: "main",
	}
	secretsListRequestBody, err := json.Marshal(secretsListRequest) //nolint:govet // The lock field is not set on this proto
	require.NoError(t, err, "failed to marshal secrets request")
	secretsUpdateRequestBodyJSON := json.RawMessage(secretsListRequestBody)
	jsonRequest := jsonrpc.Request[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      uniqueRequestID,
		Method:  vaulttypes.MethodSecretsList,
		Params:  &secretsUpdateRequestBodyJSON,
	}
	allowlistRequest(t, owner, jsonRequest, sethClient, wfRegistryContract)

	// Ensure that multiple requests can be allowlisted
	uniqueRequestIDTwo := uuid.New().String()
	secretsListRequestTwo := vault_helpers.ListSecretIdentifiersRequest{
		RequestId: uniqueRequestIDTwo,
		Owner:     owner,
		Namespace: "main",
	}
	secretsListRequestBodyTwo, err := json.Marshal(secretsListRequestTwo) //nolint:govet // The lock field is not set on this proto
	require.NoError(t, err, "failed to marshal secrets request")
	secretsUpdateRequestBodyJSONTwo := json.RawMessage(secretsListRequestBodyTwo)
	jsonRequestTwo := jsonrpc.Request[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      uniqueRequestIDTwo,
		Method:  vaulttypes.MethodSecretsList,
		Params:  &secretsUpdateRequestBodyJSONTwo,
	}
	allowlistRequest(t, owner, jsonRequestTwo, sethClient, wfRegistryContract)

	// Request 1
	requestBody, err := json.Marshal(jsonRequest)
	require.NoError(t, err, "failed to marshal secrets request")

	statusCode, httpResponseBody := sendVaultRequestToGateway(t, gatewayURL, requestBody)
	require.Equal(t, http.StatusOK, statusCode, "Gateway endpoint should respond with 200 OK")
	var jsonResponse jsonrpc.Response[vaulttypes.SignedOCRResponse]
	err = json.Unmarshal(httpResponseBody, &jsonResponse)
	require.NoError(t, err, "failed to unmarshal getResponse")
	framework.L.Info().Msgf("JSON Body: %v", jsonResponse)
	if jsonResponse.Error != nil {
		require.Empty(t, jsonResponse.Error.Error())
	}

	require.Equal(t, jsonrpc.JsonRpcVersion, jsonResponse.Version)
	require.Equal(t, uniqueRequestID, jsonResponse.ID)
	require.Equal(t, vaulttypes.MethodSecretsList, jsonResponse.Method)

	signedOCRResponse := jsonResponse.Result
	framework.L.Info().Msgf("Signed OCR Response: %s", signedOCRResponse.String())

	// Request 2
	requestBodyTwo, err := json.Marshal(jsonRequestTwo)
	require.NoError(t, err, "failed to marshal secrets request")
	statusCodeTwo, httpResponseBodyTwo := sendVaultRequestToGateway(t, gatewayURL, requestBodyTwo)
	require.Equal(t, http.StatusOK, statusCodeTwo, "Gateway endpoint should respond with 200 OK")
	var jsonResponseTwo jsonrpc.Response[vaulttypes.SignedOCRResponse]
	err = json.Unmarshal(httpResponseBodyTwo, &jsonResponseTwo)
	require.NoError(t, err, "failed to unmarshal getResponse")
	framework.L.Info().Msgf("JSON Body: %v", jsonResponseTwo)
	if jsonResponseTwo.Error != nil {
		require.Empty(t, jsonResponseTwo.Error.Error())
	}
	require.Equal(t, jsonrpc.JsonRpcVersion, jsonResponseTwo.Version)
	require.Equal(t, uniqueRequestIDTwo, jsonResponseTwo.ID)
	require.Equal(t, vaulttypes.MethodSecretsList, jsonResponseTwo.Method)
	signedOCRResponseTwo := jsonResponseTwo.Result
	framework.L.Info().Msgf("Signed OCR Response: %s", signedOCRResponseTwo.String())

	// TODO: Verify the authenticity of this signed report, by ensuring that the signatures indeed match the payload

	listSecretsResponse := vault_helpers.ListSecretIdentifiersResponse{}
	err = protojson.Unmarshal(signedOCRResponse.Payload, &listSecretsResponse)
	require.NoError(t, err, "failed to decode payload into ListSecretIdentifiersResponse proto")
	framework.L.Info().Msgf("ListSecretIdentifiersResponse decoded as: %s", listSecretsResponse.String())

	require.True(t, listSecretsResponse.Success, err)
	require.GreaterOrEqual(t, len(listSecretsResponse.Identifiers), 1, "Expected at least one item in the response")
	var keys = make([]string, 0, len(listSecretsResponse.Identifiers))
	for _, identifier := range listSecretsResponse.Identifiers {
		keys = append(keys, identifier.Key)
		require.Equal(t, owner, identifier.Owner)
		require.Equal(t, vaulttypes.DefaultNamespace, identifier.Namespace)
	}
	require.Contains(t, keys, secretID)
	framework.L.Info().Msg("Secrets listed successfully")
}

func executeVaultSecretsDeleteTest(t *testing.T, secretID, owner, gatewayURL string, sethClient *seth.Client, wfRegistryContract *workflow_registry_v2_wrapper.WorkflowRegistry) {
	framework.L.Info().Msg("Deleting secret...")
	uniqueRequestID := uuid.New().String()

	secretsDeleteRequest := vault_helpers.DeleteSecretsRequest{
		RequestId: uniqueRequestID,
		Ids: []*vault_helpers.SecretIdentifier{
			{
				Key:       secretID,
				Owner:     owner,
				Namespace: "main",
			},
			{
				Key:       "invalid",
				Owner:     owner,
				Namespace: "main",
			},
		},
	}
	secretsDeleteRequestBody, err := json.Marshal(secretsDeleteRequest) //nolint:govet // The lock field is not set on this proto
	require.NoError(t, err, "failed to marshal secrets request")
	secretsDeleteRequestBodyJSON := json.RawMessage(secretsDeleteRequestBody)
	jsonRequest := jsonrpc.Request[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      uniqueRequestID,
		Method:  vaulttypes.MethodSecretsDelete,
		Params:  &secretsDeleteRequestBodyJSON,
	}
	allowlistRequest(t, owner, jsonRequest, sethClient, wfRegistryContract)

	requestBody, err := json.Marshal(jsonRequest)
	require.NoError(t, err, "failed to marshal secrets request")

	statusCode, httpResponseBody := sendVaultRequestToGateway(t, gatewayURL, requestBody)
	require.Equal(t, http.StatusOK, statusCode, "Gateway endpoint should respond with 200 OK")
	framework.L.Info().Msg("Checking jsonResponse structure...")
	var jsonResponse jsonrpc.Response[vaulttypes.SignedOCRResponse]
	err = json.Unmarshal(httpResponseBody, &jsonResponse)
	require.NoError(t, err, "failed to unmarshal getResponse")
	framework.L.Info().Msgf("JSON Body: %v", jsonResponse)
	if jsonResponse.Error != nil {
		require.Empty(t, jsonResponse.Error.Error())
	}

	require.Equal(t, jsonrpc.JsonRpcVersion, jsonResponse.Version)
	require.Equal(t, uniqueRequestID, jsonResponse.ID)
	require.Equal(t, vaulttypes.MethodSecretsDelete, jsonResponse.Method)

	signedOCRResponse := jsonResponse.Result
	framework.L.Info().Msgf("Signed OCR Response: %s", signedOCRResponse.String())

	// TODO: Verify the authenticity of this signed report, by ensuring that the signatures indeed match the payload

	deleteSecretsResponse := vault_helpers.DeleteSecretsResponse{}
	err = protojson.Unmarshal(signedOCRResponse.Payload, &deleteSecretsResponse)
	require.NoError(t, err, "failed to decode payload into DeleteSecretResponse proto")
	framework.L.Info().Msgf("DeleteSecretResponse decoded as: %s", deleteSecretsResponse.String())

	require.Len(t, deleteSecretsResponse.Responses, 2, "Expected 2 items in the response")
	result0 := deleteSecretsResponse.GetResponses()[0]
	require.True(t, result0.Success, result0.Error)
	require.Equal(t, result0.Id.Owner, owner)
	require.Equal(t, result0.Id.Key, secretID)

	result1 := deleteSecretsResponse.GetResponses()[1]
	require.Contains(t, result1.Error, "key does not exist")

	framework.L.Info().Msg("Secrets deleted successfully")
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
	time.Sleep(5 * time.Second) // wait a bit to ensure the allowlist is propagated onchain, gateway and vault don nodes
	allowedList, err := wfRegistryContract.GetAllowlistedRequests(&bind.CallOpts{}, big.NewInt(0), big.NewInt(100))
	require.NoError(t, err, "failed to validate allowlisted request")
	for _, req := range allowedList {
		if req.RequestDigest == reqDigestBytes {
			framework.L.Info().Msgf("Request digest found in allowlist")
		}
		framework.L.Info().Msgf("Allowlisted request digestHexStr: %s, owner: %s, expiry: %d", hex.EncodeToString(req.RequestDigest[:]), req.Owner.Hex(), req.ExpiryTimestamp)
	}
}
