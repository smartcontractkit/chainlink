package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/v2/core/logger/audit"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

const (
	// PerplexityAPIURL is the base URL for Perplexity AI API
	PerplexityAPIURL = "https://api.perplexity.ai/chat/completions"
)

// ChatbotController handles Perplexity AI chatbot requests
type ChatbotController struct {
	App        chainlink.Application
	httpClient *http.Client
}

// NewChatbotController creates a new ChatbotController instance
func NewChatbotController(app chainlink.Application) *ChatbotController {
	return &ChatbotController{
		App:        app,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// ChatbotRequest represents the request payload for chatbot
type ChatbotRequest struct {
	Message string `json:"message" binding:"required"`
	Model   string `json:"model,omitempty"`
}

// ChatbotResponse represents the response from the chatbot
type ChatbotResponse struct {
	Response string `json:"response"`
	Model    string `json:"model"`
}

// PerplexityAPIRequest represents the request to Perplexity API
type PerplexityAPIRequest struct {
	Model    string                   `json:"model"`
	Messages []PerplexityAPIMessage   `json:"messages"`
}

// PerplexityAPIMessage represents a message in the Perplexity API request
type PerplexityAPIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// PerplexityAPIResponse represents the response from Perplexity API
type PerplexityAPIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Model string `json:"model"`
}

// Chat handles chatbot requests and forwards them to Perplexity AI API
func (cc *ChatbotController) Chat(c *gin.Context) {
	ctx := c.Request.Context()
	request := &ChatbotRequest{}

	if err := c.ShouldBindJSON(request); err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	// Default model if not specified
	model := request.Model
	if model == "" {
		model = "llama-3.1-sonar-small-128k-online"
	}

	// Prepare Perplexity API request
	perplexityReq := PerplexityAPIRequest{
		Model: model,
		Messages: []PerplexityAPIMessage{
			{
				Role:    "user",
				Content: request.Message,
			},
		},
	}

	reqBody, err := json.Marshal(perplexityReq)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	// Make request to Perplexity API
	httpReq, err := http.NewRequestWithContext(ctx, "POST", PerplexityAPIURL, bytes.NewBuffer(reqBody))
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	// Get API key from request headers
	apiKey := c.GetHeader("X-Perplexity-API-Key")
	if apiKey == "" {
		jsonAPIError(c, http.StatusBadRequest, errors.New("X-Perplexity-API-Key header is required"))
		return
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := cc.httpClient.Do(httpReq)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		cc.App.GetLogger().Errorw("Perplexity API error", "status", resp.StatusCode, "body", string(body))
		jsonAPIError(c, resp.StatusCode, fmt.Errorf("Perplexity API error: %s", string(body)))
		return
	}

	// Parse Perplexity API response
	var perplexityResp PerplexityAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&perplexityResp); err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	// Extract response
	var responseContent string
	if len(perplexityResp.Choices) > 0 {
		responseContent = perplexityResp.Choices[0].Message.Content
	}

	response := ChatbotResponse{
		Response: responseContent,
		Model:    perplexityResp.Model,
	}

	cc.App.GetAuditLogger().Audit(audit.ChatbotInteractionSuccess, map[string]interface{}{
		"message": request.Message,
		"model":   model,
	})

	c.JSON(http.StatusOK, response)
}
