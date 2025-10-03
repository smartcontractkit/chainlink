package cre

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"google.golang.org/protobuf/encoding/protojson"

	vault_helpers "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"

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
			var vaultFound bool
			for _, cap := range nodeSet.Capabilities {
				if cap == cre.VaultCapability {
					vaultFound = true
					break
				}
			}
			if vaultFound {
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

	owner := "Owner1"
	waitUntilReady(t, owner, gatewayURL.String())

	secretID := strconv.Itoa(rand.Intn(10000)) // generate a random secret ID for testing
	secretValue := "Secret Value to be stored"
	vaultPublicKey := FetchVaultPublicKey(t, gatewayURL.String())
	encryptedSecret, err := crevault.EncryptSecret(secretValue, vaultPublicKey)
	require.NoError(t, err, "failed to encrypt secret")

	// Wait for the node to be up.
	framework.L.Info().Msg("Waiting 30 seconds for the Vault DON to be ready...")
	time.Sleep(30 * time.Second)
	executeVaultSecretsCreateTest(t, encryptedSecret, secretID, owner, gatewayURL.String())
	executeVaultSecretsGetTest(t, secretID, owner, gatewayURL.String())
	executeVaultSecretsUpdateTest(t, encryptedSecret, secretID, owner, gatewayURL.String())
	executeVaultSecretsListTest(t, secretID, owner, gatewayURL.String())
	executeVaultSecretsDeleteTest(t, secretID, owner, gatewayURL.String())
}

// waitUntilReady tries to list the keys in a loop until it succeeds, indicating that the Vault DON is ready.
func waitUntilReady(t *testing.T, owner, gatewayURL string) {
	framework.L.Info().Msg("Polling for vault DON to be ready...")

	uniqueRequestID := uuid.New().String()

	getPublicKeyRequest := jsonrpc.Request[vault_helpers.ListSecretIdentifiersRequest]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      uniqueRequestID,
		Method:  vaulttypes.MethodSecretsList,
		Params: &vault_helpers.ListSecretIdentifiersRequest{
			Owner: owner,
		},
	}
	requestBody, err := json.Marshal(getPublicKeyRequest)
	require.NoError(t, err, "failed to marshal public key request")

	statusCode, _ := sendVaultRequestToGateway(t, gatewayURL, requestBody)
	if statusCode == http.StatusGatewayTimeout {
		framework.L.Warn().Msg("Received 504 Gateway Timeout. This may be due to the Vault DON not being ready yet. Retrying 1st time in 30 seconds...")
		time.Sleep(30 * time.Second)
		statusCode, _ = sendVaultRequestToGateway(t, gatewayURL, requestBody)
		if statusCode == http.StatusGatewayTimeout {
			framework.L.Warn().Msg("Received 504 Gateway Timeout again. This may be due to the Vault DON not being ready yet. Retrying 2nd time in 30 seconds...")
			time.Sleep(30 * time.Second)
			statusCode, _ = sendVaultRequestToGateway(t, gatewayURL, requestBody)
		}
	}
	require.Equal(t, http.StatusOK, statusCode, "Gateway endpoint should respond with 200 OK")

	framework.L.Info().Msgf("Received ready response from Vault DON")
}

func executeVaultSecretsCreateTest(t *testing.T, encryptedSecret, secretID, owner, gatewayURL string) {
	framework.L.Info().Msg("Creating secret...")

	uniqueRequestID := uuid.New().String()

	secretsCreateRequest := jsonrpc.Request[vault_helpers.CreateSecretsRequest]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      uniqueRequestID,
		Method:  vaulttypes.MethodSecretsCreate,
		Params: &vault_helpers.CreateSecretsRequest{
			RequestId: uniqueRequestID,
			EncryptedSecrets: []*vault_helpers.EncryptedSecret{
				{
					Id: &vault_helpers.SecretIdentifier{
						Key:   secretID,
						Owner: owner,
						// Namespace: "main", // Uncomment if you want to use namespaces
					}, // Note: Namespace is not used in this test, but can be added if needed
					EncryptedValue: encryptedSecret,
				},
			},
		},
	}
	requestBody, err := json.Marshal(secretsCreateRequest)
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

func executeVaultSecretsUpdateTest(t *testing.T, encryptedSecret, secretID, owner, gatewayURL string) {
	framework.L.Info().Msg("Updating secret...")
	uniqueRequestID := uuid.New().String()

	secretsUpdateRequest := jsonrpc.Request[vault_helpers.UpdateSecretsRequest]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      uniqueRequestID,
		Method:  vaulttypes.MethodSecretsUpdate,
		Params: &vault_helpers.UpdateSecretsRequest{
			RequestId: uniqueRequestID,
			EncryptedSecrets: []*vault_helpers.EncryptedSecret{
				{
					Id: &vault_helpers.SecretIdentifier{
						Key:   secretID,
						Owner: owner,
					},
					EncryptedValue: encryptedSecret,
				},
				{
					Id: &vault_helpers.SecretIdentifier{
						Key:   "invalid",
						Owner: "invalid",
					},
					EncryptedValue: encryptedSecret,
				},
			},
		},
	}
	requestBody, err := json.Marshal(secretsUpdateRequest)
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

func executeVaultSecretsGetTest(t *testing.T, secretID, owner, gatewayURL string) {
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

func executeVaultSecretsListTest(t *testing.T, secretID, owner, gatewayURL string) {
	framework.L.Info().Msg("Listing secret...")
	uniqueRequestID := uuid.New().String()

	secretsListRequest := jsonrpc.Request[vault_helpers.ListSecretIdentifiersRequest]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      uniqueRequestID,
		Method:  vaulttypes.MethodSecretsList,
		Params: &vault_helpers.ListSecretIdentifiersRequest{
			RequestId: uniqueRequestID,
			Owner:     owner,
		},
	}
	requestBody, err := json.Marshal(secretsListRequest)
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

func executeVaultSecretsDeleteTest(t *testing.T, secretID, owner, gatewayURL string) {
	framework.L.Info().Msg("Deleting secret...")
	uniqueRequestID := uuid.New().String()

	secretsUpdateRequest := jsonrpc.Request[vault_helpers.DeleteSecretsRequest]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      uniqueRequestID,
		Method:  vaulttypes.MethodSecretsDelete,
		Params: &vault_helpers.DeleteSecretsRequest{
			RequestId: uniqueRequestID,
			Ids: []*vault_helpers.SecretIdentifier{
				{
					Key:   secretID,
					Owner: owner,
				},
				{
					Key:   "invalid",
					Owner: "invalid",
				},
			},
		},
	}
	requestBody, err := json.Marshal(secretsUpdateRequest)
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
