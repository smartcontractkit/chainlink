#!/bin/bash
# Example script to test the Perplexity AI Chatbot API

# Configuration
NODE_URL="${CHAINLINK_NODE_URL:-http://localhost:6688}"
EMAIL="${CHAINLINK_EMAIL:-user@example.com}"
PASSWORD="${CHAINLINK_PASSWORD:-password}"
PERPLEXITY_API_KEY="${PERPLEXITY_API_KEY}"

if [ -z "$PERPLEXITY_API_KEY" ]; then
    echo "Error: PERPLEXITY_API_KEY environment variable is required"
    echo "Usage: PERPLEXITY_API_KEY=your-key ./chatbot_api_test.sh"
    exit 1
fi

# Login and get session cookie
echo "Logging in to Chainlink node..."
COOKIE_FILE=$(mktemp)
curl -s -c "$COOKIE_FILE" -X POST "$NODE_URL/sessions" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}" > /dev/null

# Test the chatbot API
echo "Testing chatbot API..."
RESPONSE=$(curl -s -b "$COOKIE_FILE" -X POST "$NODE_URL/v2/chatbot" \
  -H "Content-Type: application/json" \
  -H "X-Perplexity-API-Key: $PERPLEXITY_API_KEY" \
  -d '{
    "message": "What is Chainlink?",
    "model": "llama-3.1-sonar-small-128k-online"
  }')

echo "Response:"
echo "$RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$RESPONSE"

# Cleanup
rm -f "$COOKIE_FILE"
