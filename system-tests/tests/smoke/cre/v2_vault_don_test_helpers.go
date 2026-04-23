package cre

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	vault_helpers "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
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

func sendVaultRequestToGateway(t *testing.T, gatewayURL string, requestBody []byte) (statusCode int, body []byte) {
	const maxRetries = 7
	const retryInterval = 2 * time.Second

	framework.L.Info().Msgf("Request Body: %s", string(requestBody))

	for attempt := range maxRetries + 1 {
		req, err := http.NewRequestWithContext(t.Context(), "POST", gatewayURL, bytes.NewBuffer(requestBody))
		require.NoError(t, err, "failed to create request")

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

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
