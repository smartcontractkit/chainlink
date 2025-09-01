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

	crevault "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/vault"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

func ExecuteVaultTest(t *testing.T, testEnv *TestEnvironment) {
	// Skip till we figure out and fix the issues with environment startup on this test
	const skipReason = "Skip till the errors with topology TopologyWorkflowGatewayCapabilities are fixed: https://smartcontract-it.atlassian.net/browse/PRIV-160"
	t.Skipf("Skipping test for the following reason: %s", skipReason)
	/*
		BUILD ENVIRONMENT FROM SAVED STATE
	*/
	var testLogger = framework.L
	testLogger.Info().Msg("Getting gateway configuration...")
	require.NotEmpty(t, testEnv.FullCldEnvOutput.DonTopology.GatewayConnectorOutput.Configurations, "expected at least one gateway configuration")
	gatewayURL, err := url.Parse(testEnv.FullCldEnvOutput.DonTopology.GatewayConnectorOutput.Configurations[0].Incoming.Protocol + "://" + testEnv.FullCldEnvOutput.DonTopology.GatewayConnectorOutput.Configurations[0].Incoming.Host + ":" + strconv.Itoa(testEnv.FullCldEnvOutput.DonTopology.GatewayConnectorOutput.Configurations[0].Incoming.ExternalPort) + testEnv.FullCldEnvOutput.DonTopology.GatewayConnectorOutput.Configurations[0].Incoming.Path)
	require.NoError(t, err, "failed to parse gateway URL")

	testLogger.Info().Msgf("Gateway URL: %s", gatewayURL.String())

	framework.L.Info().Msgf("Sleeping 1 minute to allow the Vault DON to start...")
	// TODO: Remove this sleep https://smartcontract-it.atlassian.net/browse/PRIV-154
	time.Sleep(1 * time.Minute)
	testLogger.Info().Msgf("Sleep over. Executing test now...")

	secretID := strconv.Itoa(rand.Intn(10000)) // generate a random secret ID for testing
	owner := "Owner1"
	secretValue := "Secret Value to be stored"

	executeVaultSecretsCreateTest(t, secretValue, secretID, owner, gatewayURL.String())

	divider := "------------------------------------------------------"
	testLogger.Info().Msgf("%s \n%s \n%s \n%s \n%s", divider, divider, divider, divider, divider)

	executeVaultSecretsGetTest(t, secretValue, secretID, owner, gatewayURL.String())
	executeVaultSecretsUpdateTest(t, secretValue, secretID, owner, gatewayURL.String())
}

func executeVaultSecretsCreateTest(t *testing.T, secretValue, secretID, owner, gatewayURL string) {
	framework.L.Info().Msg("Creating secret...")
	encryptedSecret, err := crevault.EncryptSecret(secretValue)
	require.NoError(t, err, "failed to encrypt secret")

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

	httpResponseBody := sendVaultRequestToGateway(t, gatewayURL, requestBody)
	framework.L.Info().Msg("Checking jsonResponse structure...")
	var jsonResponse jsonrpc.Response[vaulttypes.SignedOCRResponse]
	err = json.Unmarshal(httpResponseBody, &jsonResponse)
	require.NoError(t, err, "failed to unmarshal getResponse")
	framework.L.Info().Msgf("JSON Body: %v", jsonResponse)
	if jsonResponse.Error != nil {
		require.Empty(t, jsonResponse.Error.Error())
	}
	require.Equal(t, jsonrpc.JsonRpcVersion, jsonResponse.Version)
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

	framework.L.Info().Msg("Secret created successfully")
}

func executeVaultSecretsUpdateTest(t *testing.T, secretValue, secretID, owner, gatewayURL string) {
	framework.L.Info().Msg("Updating secret...")
	uniqueRequestID := uuid.New().String()

	encryptedSecret, secretErr := crevault.EncryptSecret(secretValue)
	require.NoError(t, secretErr, "failed to encrypt secret")

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

	httpResponseBody := sendVaultRequestToGateway(t, gatewayURL, requestBody)
	framework.L.Info().Msg("Checking jsonResponse structure...")
	var jsonResponse jsonrpc.Response[vaulttypes.SignedOCRResponse]
	err = json.Unmarshal(httpResponseBody, &jsonResponse)
	require.NoError(t, err, "failed to unmarshal getResponse")
	framework.L.Info().Msgf("JSON Body: %v", jsonResponse)
	if jsonResponse.Error != nil {
		require.Empty(t, jsonResponse.Error.Error())
	}

	require.Equal(t, jsonrpc.JsonRpcVersion, jsonResponse.Version)
	require.Equal(t, vaulttypes.MethodSecretsUpdate, jsonResponse.Method)

	signedOCRResponse := jsonResponse.Result
	framework.L.Info().Msgf("Signed OCR Response: %s", signedOCRResponse.String())

	// TODO: Verify the authenticity of this signed report, by ensuring that the signatures indeed match the payload

	updateSecretsResponse := vaultcommon.UpdateSecretsResponse{}
	err = protojson.Unmarshal(signedOCRResponse.Payload, &updateSecretsResponse)
	require.NoError(t, err, "failed to decode payload into UpdateSecretsResponse proto")
	framework.L.Info().Msgf("UpdateSecretsResponse decoded as: %s", updateSecretsResponse.String())

	require.Len(t, updateSecretsResponse.Responses, 2, "Expected one item in the response")
	result0 := updateSecretsResponse.GetResponses()[0]
	require.Empty(t, result0.GetError())
	require.Equal(t, secretID, result0.GetId().Key)
	require.Equal(t, owner, result0.GetId().Owner)

	result1 := updateSecretsResponse.GetResponses()[1]
	require.Contains(t, result1.Error, "key does not exist")

	framework.L.Info().Msg("Secret updated successfully")
}

func executeVaultSecretsGetTest(t *testing.T, secretValue, secretID, owner, gatewayURL string) {
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
	httpResponseBody := sendVaultRequestToGateway(t, gatewayURL, requestBody)
	framework.L.Info().Msg("Checking jsonResponse structure...")
	var jsonResponse jsonrpc.Response[json.RawMessage]
	err = json.Unmarshal(httpResponseBody, &jsonResponse)
	require.NoError(t, err, "failed to unmarshal http response body")
	framework.L.Info().Msgf("JSON Body: %v", jsonResponse)
	if jsonResponse.Error != nil {
		require.Empty(t, jsonResponse.Error.Error())
	}
	require.Equal(t, jsonrpc.JsonRpcVersion, jsonResponse.Version)
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

	framework.L.Info().Msg("Secret get successful")
}

func sendVaultRequestToGateway(t *testing.T, gatewayURL string, requestBody []byte) []byte {
	framework.L.Info().Msgf("Request Body: %s", string(requestBody))
	req, err := http.NewRequestWithContext(context.Background(), "POST", gatewayURL, bytes.NewBuffer(requestBody))
	require.NoError(t, err, "failed to create request")

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err, "failed to execute request")
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "failed to read jsonResponse body")
	framework.L.Info().Msgf("Response Body: %s", string(body))

	require.Equal(t, http.StatusOK, resp.StatusCode, "Gateway endpoint should respond with 200 OK")
	return body
}
