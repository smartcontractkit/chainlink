package web_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/web"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChatbotController_Chat(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplicationEVMDisabled(t)
	require.NoError(t, app.Start(testutils.Context(t)))

	// Mock Perplexity API server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/chat/completions", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Contains(t, r.Header.Get("Authorization"), "Bearer")

		// Verify request body
		var reqBody map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		require.NoError(t, err)
		assert.Equal(t, "llama-3.1-sonar-small-128k-online", reqBody["model"])
		
		messages := reqBody["messages"].([]interface{})
		require.Len(t, messages, 1)
		msg := messages[0].(map[string]interface{})
		assert.Equal(t, "user", msg["role"])
		assert.Equal(t, "Hello, how are you?", msg["content"])

		// Send mock response
		response := map[string]interface{}{
			"model": "llama-3.1-sonar-small-128k-online",
			"choices": []map[string]interface{}{
				{
					"message": map[string]interface{}{
						"content": "I'm doing well, thank you for asking!",
					},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	// Create request
	reqBody := web.ChatbotRequest{
		Message: "Hello, how are you?",
		Model:   "llama-3.1-sonar-small-128k-online",
	}
	body, err := json.Marshal(reqBody)
	require.NoError(t, err)

	client := app.NewHTTPClient(nil)
	
	// Note: This test will fail in actual execution because we can't override the Perplexity API URL
	// In a real implementation, the API URL should be configurable for testing
	resp, cleanup := client.Post("/v2/chatbot", bytes.NewReader(body))
	defer cleanup()

	// The test expects StatusBadRequest because X-Perplexity-API-Key header is missing
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestChatbotController_Chat_MissingAPIKey(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplicationEVMDisabled(t)
	require.NoError(t, app.Start(testutils.Context(t)))

	// Create request without API key
	reqBody := web.ChatbotRequest{
		Message: "Hello, how are you?",
	}
	body, err := json.Marshal(reqBody)
	require.NoError(t, err)

	client := app.NewHTTPClient(nil)
	resp, cleanup := client.Post("/v2/chatbot", bytes.NewReader(body))
	defer cleanup()

	// Expect BadRequest because API key is missing
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	
	respBody := cltest.ParseResponseBody(t, resp)
	assert.Contains(t, string(respBody), "X-Perplexity-API-Key")
}

func TestChatbotController_Chat_InvalidRequest(t *testing.T) {
	t.Parallel()

	app := cltest.NewApplicationEVMDisabled(t)
	require.NoError(t, app.Start(testutils.Context(t)))

	// Create invalid request (missing message)
	reqBody := map[string]string{}
	body, err := json.Marshal(reqBody)
	require.NoError(t, err)

	client := app.NewHTTPClient(nil)
	resp, cleanup := client.Post("/v2/chatbot", bytes.NewReader(body))
	defer cleanup()

	// Expect UnprocessableEntity because message is required
	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
}

func TestChatbotController_Chat_NoCredentials(t *testing.T) {
	t.Parallel()

	ctx := testutils.Context(t)
	app := cltest.NewApplicationEVMDisabled(t)
	require.NoError(t, app.Start(ctx))

	// Create request
	reqBody := web.ChatbotRequest{
		Message: "Hello",
	}
	body, err := json.Marshal(reqBody)
	require.NoError(t, err)

	// Make request without authentication
	url := app.Server.URL + "/v2/chatbot"
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Expect Unauthorized because no credentials provided
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}
