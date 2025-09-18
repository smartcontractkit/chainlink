package cre

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"google.golang.org/protobuf/encoding/protojson"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	crevault "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/vault"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/vault"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

func ExecuteVaultTest(t *testing.T, testEnv *TestEnvironment) {
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
	require.NotEmpty(t, testEnv.FullCldEnvOutput.DonTopology.GatewayConnectorOutput.Configurations, "expected at least one gateway configuration")
	gatewayURL, err := url.Parse(testEnv.FullCldEnvOutput.DonTopology.GatewayConnectorOutput.Configurations[0].Incoming.Protocol + "://" + testEnv.FullCldEnvOutput.DonTopology.GatewayConnectorOutput.Configurations[0].Incoming.Host + ":" + strconv.Itoa(testEnv.FullCldEnvOutput.DonTopology.GatewayConnectorOutput.Configurations[0].Incoming.ExternalPort) + testEnv.FullCldEnvOutput.DonTopology.GatewayConnectorOutput.Configurations[0].Incoming.Path)
	require.NoError(t, err, "failed to parse gateway URL")
	testLogger.Info().Msgf("Gateway URL: %s", gatewayURL.String())

	owner := "Owner1"
	waitUntilReady(t, owner, gatewayURL.String())

	secretID := strconv.Itoa(rand.Intn(10000)) // generate a random secret ID for testing
	secretValue := "Secret Value to be stored"
	vaultPublicKey := fetchVaultPublicKey(t, gatewayURL.String())
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

	getPublicKeyRequest := jsonrpc.Request[vaultcommon.ListSecretIdentifiersRequest]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      uniqueRequestID,
		Method:  vaulttypes.MethodSecretsList,
		Params: &vaultcommon.ListSecretIdentifiersRequest{
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

func fetchVaultPublicKey(t *testing.T, gatewayURL string) (publicKey string) {
	framework.L.Info().Msg("Fetching Vault Public Key...")

	uniqueRequestID := uuid.New().String()

	getPublicKeyRequest := jsonrpc.Request[vaultcommon.GetPublicKeyRequest]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      uniqueRequestID,
		Method:  vaulttypes.MethodPublicKeyGet,
		Params:  &vaultcommon.GetPublicKeyRequest{},
	}
	requestBody, err := json.Marshal(getPublicKeyRequest)
	require.NoError(t, err, "failed to marshal public key request")

	statusCode, httpResponseBody := sendVaultRequestToGateway(t, gatewayURL, requestBody)
	require.Equal(t, http.StatusOK, statusCode, "Gateway endpoint should respond with 200 OK")

	framework.L.Info().Msg("Checking jsonResponse structure...")
	var jsonResponse jsonrpc.Response[vaultcommon.GetPublicKeyResponse]
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

func executeVaultSecretsCreateTest(t *testing.T, encryptedSecret, secretID, owner, gatewayURL string) {
	framework.L.Info().Msg("Creating secret...")

	uniqueRequestID := uuid.New().String()

	secretsCreateRequest := jsonrpc.Request[vaultcommon.CreateSecretsRequest]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      uniqueRequestID,
		Method:  vaulttypes.MethodSecretsCreate,
		Params: &vaultcommon.CreateSecretsRequest{
			RequestId: uniqueRequestID,
			EncryptedSecrets: []*vaultcommon.EncryptedSecret{
				{
					Id: &vaultcommon.SecretIdentifier{
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
	createSecretsResponse := vaultcommon.CreateSecretsResponse{}
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

	secretsUpdateRequest := jsonrpc.Request[vaultcommon.UpdateSecretsRequest]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      uniqueRequestID,
		Method:  vaulttypes.MethodSecretsUpdate,
		Params: &vaultcommon.UpdateSecretsRequest{
			RequestId: uniqueRequestID,
			EncryptedSecrets: []*vaultcommon.EncryptedSecret{
				{
					Id: &vaultcommon.SecretIdentifier{
						Key:   secretID,
						Owner: owner,
					},
					EncryptedValue: encryptedSecret,
				},
				{
					Id: &vaultcommon.SecretIdentifier{
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

	updateSecretsResponse := vaultcommon.UpdateSecretsResponse{}
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
	secretsGetRequest := jsonrpc.Request[vaultcommon.GetSecretsRequest]{
		Version: jsonrpc.JsonRpcVersion,
		Method:  vaulttypes.MethodSecretsGet,
		Params: &vaultcommon.GetSecretsRequest{
			Requests: []*vaultcommon.SecretRequest{
				{
					Id: &vaultcommon.SecretIdentifier{
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
	 * The json unmarshaling is not compatible with the proto oneof in vaultcommon.SecretResponse
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
		ID    *vaultcommon.SecretIdentifier `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
		Data  *SecretData                   `protobuf:"bytes,2,opt,name=data,proto3"`
		Error string                        `protobuf:"bytes,3,opt,name=error,proto3"`
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

	secretsListRequest := jsonrpc.Request[vaultcommon.ListSecretIdentifiersRequest]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      uniqueRequestID,
		Method:  vaulttypes.MethodSecretsList,
		Params: &vaultcommon.ListSecretIdentifiersRequest{
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

	listSecretsResponse := vaultcommon.ListSecretIdentifiersResponse{}
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

	secretsUpdateRequest := jsonrpc.Request[vaultcommon.DeleteSecretsRequest]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      uniqueRequestID,
		Method:  vaulttypes.MethodSecretsDelete,
		Params: &vaultcommon.DeleteSecretsRequest{
			RequestId: uniqueRequestID,
			Ids: []*vaultcommon.SecretIdentifier{
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

	deleteSecretsResponse := vaultcommon.DeleteSecretsResponse{}
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

func sendVaultRequestToGateway(t *testing.T, gatewayURL string, requestBody []byte) (statusCode int, body []byte) {
	framework.L.Info().Msgf("Request Body: %s", string(requestBody))
	req, err := http.NewRequestWithContext(context.Background(), "POST", gatewayURL, bytes.NewBuffer(requestBody))
	require.NoError(t, err, "failed to create request")

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err, "failed to execute request")
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	require.NoError(t, err, "failed to read http response body")
	framework.L.Info().Msgf("HTTP Response Body: %s", string(body))
	return resp.StatusCode, body
}
