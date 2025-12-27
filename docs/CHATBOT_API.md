# Perplexity AI Chatbot API

This document describes the Perplexity AI Chatbot API integration in the Chainlink node.

## Overview

The Chainlink node now includes an API endpoint that allows authenticated users to interact with Perplexity AI's chatbot service. This enables AI-powered question answering and conversational capabilities within the Chainlink ecosystem.

## API Endpoint

### POST /v2/chatbot

Send a message to the Perplexity AI chatbot and receive a response.

#### Authentication

This endpoint requires authentication with either:
- Session-based authentication
- API token authentication

The endpoint also requires a valid Perplexity AI API key to be provided in the request headers.

#### Request Headers

- `Content-Type`: `application/json` (required)
- `X-Perplexity-API-Key`: Your Perplexity AI API key (required)
- Authentication headers (session or API token)

#### Request Body

```json
{
  "message": "Your question or message here",
  "model": "llama-3.1-sonar-small-128k-online"
}
```

Fields:
- `message` (string, required): The message or question to send to the chatbot
- `model` (string, optional): The Perplexity AI model to use. Defaults to `llama-3.1-sonar-small-128k-online` if not specified

#### Response

Success (200 OK):
```json
{
  "response": "The chatbot's response to your message",
  "model": "llama-3.1-sonar-small-128k-online"
}
```

Error responses:
- `400 Bad Request`: Missing API key or invalid request
- `401 Unauthorized`: Missing or invalid authentication
- `422 Unprocessable Entity`: Invalid request body (missing required fields)
- `500 Internal Server Error`: Server error or Perplexity API error

## Example Usage

### Using cURL

```bash
curl -X POST https://your-chainlink-node.com/v2/chatbot \
  -H "Content-Type: application/json" \
  -H "X-Chainlink-EA-AccessKey: your-access-key" \
  -H "X-Chainlink-EA-Secret: your-secret" \
  -H "X-Perplexity-API-Key: your-perplexity-api-key" \
  -d '{
    "message": "What is Chainlink?",
    "model": "llama-3.1-sonar-small-128k-online"
  }'
```

### Response Example

```json
{
  "response": "Chainlink is a decentralized oracle network that enables smart contracts to securely interact with real-world data, events, and payments. It provides reliable, tamper-proof inputs and outputs for complex smart contracts on any blockchain.",
  "model": "llama-3.1-sonar-small-128k-online"
}
```

## Available Models

Perplexity AI offers various models. Some common options:
- `llama-3.1-sonar-small-128k-online` (default)
- `llama-3.1-sonar-large-128k-online`
- `llama-3.1-sonar-huge-128k-online`

Refer to [Perplexity AI documentation](https://docs.perplexity.ai/) for the complete list of available models.

## Security Considerations

1. **API Key Protection**: The Perplexity API key is passed via request headers and is not stored by the Chainlink node. Ensure your API keys are kept secure and never committed to version control.

2. **Authentication Required**: All requests to this endpoint must be authenticated. This prevents unauthorized access to the chatbot functionality.

3. **Audit Logging**: All successful chatbot interactions are logged via the audit logging system for compliance and monitoring purposes.

4. **Rate Limiting**: The endpoint is subject to the same rate limiting rules as other authenticated endpoints in the Chainlink node.

## Audit Events

Successful chatbot interactions generate an audit event:
- Event ID: `CHATBOT_INTERACTION_SUCCESS`
- Data includes: message and model used

These events can be forwarded to your configured audit logging service.

## Notes

- The endpoint forwards requests to Perplexity AI's API at `https://api.perplexity.ai/chat/completions`
- Network connectivity to Perplexity AI's API is required
- Response times depend on Perplexity AI's service performance
- The implementation follows Chainlink's standard controller patterns for consistency
