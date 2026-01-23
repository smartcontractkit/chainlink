package cre

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
)

var (
	lloFakeProviderStarted sync.Once
	lloFakeProviderOutput  *fake.Output
	lloFakeProviderErr     error
)

func init() {
	// Disable Gin request logging at package init time to reduce noise in test output
	// This must be set before any Gin router is created
	// Set both the global mode and environment variable (Gin checks env var first)
	if os.Getenv("GIN_MODE") == "" {
		os.Setenv("GIN_MODE", "release")
	}
	gin.SetMode(gin.ReleaseMode)
	// Redirect Gin's default writer to discard to prevent any logging
	// This is a more aggressive approach that ensures no Gin logs appear
	gin.DefaultWriter = io.Discard
	gin.DefaultErrorWriter = io.Discard
}

// LLOMagicNumbers defines the magic numbers for each stream ID for E2E verification
// If these values appear in workflow logs, it proves full E2E connectivity
const (
	// MAGIC_NUMBER_FORMAT5 is for ReportFormat 5 (CapabilityTrigger) - Stream 1 (TEST/USD)
	LLOMagicNumberFormat5 = 424242
	// MAGIC_NUMBER_FORMAT7 is for ReportFormat 7 (EVMABIEncodeUnpackedExpr) - Stream 4 (DATA/USD) base value
	// This value (111111) is multiplied by 5 via calculated stream expression to get 555555
	LLOMagicNumberFormat7Base = 111111
	// MAGIC_NUMBER_FORMAT7 is the expected final value after calculated stream expression (111111 * 5 = 555555)
	LLOMagicNumberFormat7 = 555555
)

// LLOPriceConfig defines the prices returned by the mock LLO price provider
type LLOPriceConfig struct {
	// TestUSD is the price for TEST/USD (Stream ID 1, Format 5)
	TestUSD float64
	// NativeUSD is the price for NATIVE/USD (Stream ID 2, Format 7 fee)
	NativeUSD float64
	// LinkUSD is the price for LINK/USD (Stream ID 3, Format 7 fee)
	LinkUSD float64
	// DataUSD is the price for DATA/USD (Stream ID 4, Format 7 data)
	DataUSD float64
}

// DefaultLLOPriceConfig returns the default configuration with magic numbers for E2E verification
func DefaultLLOPriceConfig() LLOPriceConfig {
	return LLOPriceConfig{
		TestUSD:   float64(LLOMagicNumberFormat5), // 424242
		NativeUSD: 3000.00,
		LinkUSD:   15.00,
		DataUSD:   float64(LLOMagicNumberFormat7Base), // 111111 (will be multiplied by 5 via calculated stream)
	}
}

// LLOPriceRequest represents a Chainlink External Adapter request for LLO
type LLOPriceRequest struct {
	ID   string                 `json:"id"`
	Data map[string]interface{} `json:"data"`
}

// LLOPriceResponse represents a Chainlink External Adapter response for LLO
type LLOPriceResponse struct {
	JobRunID   string      `json:"jobRunID"`
	Data       interface{} `json:"data"`
	Result     interface{} `json:"result"`
	StatusCode int         `json:"statusCode"`
}

// SetupLLOPriceProvider starts the fake data provider for LLO E2E tests.
// It serves mock prices for TEST/USD, NATIVE/USD, LINK/USD, and DATA/USD.
// The provider is started only once across all test runs (sync.Once).
// Uses NewFakeDataProvider which runs locally and is accessible via host.docker.internal.
func SetupLLOPriceProvider(testLogger zerolog.Logger, input *fake.Input, config LLOPriceConfig) (string, error) {
	// This sync.Once ensures that the fake data provider is only started once across all test runs.
	lloFakeProviderStarted.Do(func() {
		// Use local fake provider - Docker containers access via host.docker.internal
		// Gin mode is already set to ReleaseMode in init() to disable request logging
		lloFakeProviderOutput, lloFakeProviderErr = fake.NewFakeDataProvider(input)
		if lloFakeProviderErr != nil {
			testLogger.Error().Err(lloFakeProviderErr).Msg("Failed to start LLO fake data provider")
		} else {
			testLogger.Info().
				Str("baseURLDocker", lloFakeProviderOutput.BaseURLDocker).
				Str("baseURLHost", lloFakeProviderOutput.BaseURLHost).
				Msg("LLO fake data provider started successfully")
		}
	})

	if lloFakeProviderErr != nil {
		return "", errors.Wrap(lloFakeProviderErr, "failed to start fake data provider")
	}

	if lloFakeProviderOutput == nil {
		return "", errors.New("fake provider output is nil - provider may not have started correctly")
	}

	lloAPIPath := "/llo/price"
	// Use the Docker URL for container-to-container communication
	// Try to use gateway IP if host.docker.internal might not work
	baseURL := lloFakeProviderOutput.BaseURLDocker
	if gatewayIP := getDockerGatewayIP(); gatewayIP != "" {
		// Replace host.docker.internal with gateway IP for better reliability
		baseURL = strings.Replace(baseURL, "host.docker.internal", gatewayIP, 1)
		testLogger.Info().
			Str("originalURL", lloFakeProviderOutput.BaseURLDocker).
			Str("gatewayIP", gatewayIP).
			Str("updatedURL", baseURL).
			Msg("Using Docker gateway IP for price provider access")
	}
	lloFinalURL := baseURL + lloAPIPath

	// Price lookup table
	prices := map[string]float64{
		"TEST/USD":   config.TestUSD,
		"NATIVE/USD": config.NativeUSD,
		"LINK/USD":   config.LinkUSD,
		"DATA/USD":   config.DataUSD,
	}

	// Map stream IDs to trading pairs
	streamIDToPair := map[int]string{
		1: "TEST/USD",
		2: "NATIVE/USD",
		3: "LINK/USD",
		4: "DATA/USD",
	}

	err := fake.Func("POST", lloAPIPath, func(c *gin.Context) {
		var req LLOPriceRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "invalid request body"})
			return
		}

		// Extract trading pair from request
		pair := extractLLOPair(req.Data, streamIDToPair)
		if pair == "" {
			c.JSON(400, gin.H{"error": "missing or invalid trading pair"})
			return
		}

		price, ok := prices[pair]
		if !ok {
			c.JSON(404, gin.H{"error": fmt.Sprintf("unknown trading pair: %s", pair)})
			return
		}

		testLogger.Debug().
			Str("pair", pair).
			Float64("price", price).
			Str("requestID", req.ID).
			Msg("LLO price request")

		response := LLOPriceResponse{
			JobRunID: req.ID,
			Data: map[string]interface{}{
				"result":               price,
				"pair":                 pair,
				"magic_number_format5": LLOMagicNumberFormat5,
				"magic_number_format7": LLOMagicNumberFormat7,
			},
			Result:     price,
			StatusCode: 200,
		}

		c.JSON(200, response)
	})
	if err != nil {
		return "", errors.Wrap(err, "failed to set up LLO fake price provider")
	}

	// Also serve channel definitions on this provider (needed for LLO)
	channelDefsPath := "/channel-definitions.json"
	channelDefsBaseURL := lloFakeProviderOutput.BaseURLDocker
	if gatewayIP := getDockerGatewayIP(); gatewayIP != "" {
		channelDefsBaseURL = strings.Replace(channelDefsBaseURL, "host.docker.internal", gatewayIP, 1)
	}
	channelDefsURL := channelDefsBaseURL + channelDefsPath

	channelDefs := llotypes.ChannelDefinitions{
		// Channel 1: Format 5 (Capability Trigger) for TEST/USD - magic number 424242
		1: {
			ReportFormat: llotypes.ReportFormatCapabilityTrigger,
			Streams: []llotypes.Stream{
				{StreamID: 1, Aggregator: llotypes.AggregatorMedian},
			},
		},
		// Channel 2: Format 7 (EVM ABI Encode Unpacked) for DATA/USD
		// Stream 4 has base value 111111, calculated stream 5 multiplies by 5 to get 555555
		2: {
			ReportFormat: llotypes.ReportFormatEVMABIEncodeUnpackedExpr,
			Streams: []llotypes.Stream{
				{StreamID: 4, Aggregator: llotypes.AggregatorMedian},
				{StreamID: 5, Aggregator: llotypes.AggregatorCalculated},
			},
		},
	}

	channelDefsJSON, err := json.Marshal(channelDefs)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal channel definitions")
	}

	err = fake.Func("GET", channelDefsPath, func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.String(200, string(channelDefsJSON))
	})
	if err != nil {
		return "", errors.Wrap(err, "failed to set up channel definitions endpoint")
	}

	testLogger.Info().
		Str("priceURL", lloFinalURL).
		Str("channelDefsURL", channelDefsURL).
		Float64("TEST/USD", config.TestUSD).
		Float64("NATIVE/USD", config.NativeUSD).
		Float64("LINK/USD", config.LinkUSD).
		Float64("DATA/USD", config.DataUSD).
		Msg("LLO fake price provider configured")

	return lloFinalURL, nil
}

// GetLLOProviderDockerURL returns the Docker URL of the LLO fake provider
// This URL should be used by Docker containers to reach the provider
// Uses gateway IP if available for better reliability
func GetLLOProviderDockerURL() string {
	if lloFakeProviderOutput == nil {
		return ""
	}
	baseURL := lloFakeProviderOutput.BaseURLDocker
	if gatewayIP := getDockerGatewayIP(); gatewayIP != "" {
		baseURL = strings.Replace(baseURL, "host.docker.internal", gatewayIP, 1)
	}
	return baseURL
}

// GetLLOProviderPriceURL returns the full URL for the price endpoint
func GetLLOProviderPriceURL() string {
	if lloFakeProviderOutput == nil {
		return ""
	}
	return GetLLOProviderDockerURL() + "/llo/price"
}

// GetLLOProviderChannelDefsURL returns the full URL for the channel definitions endpoint
func GetLLOProviderChannelDefsURL() string {
	if lloFakeProviderOutput == nil {
		return ""
	}
	return GetLLOProviderDockerURL() + "/channel-definitions.json"
}

// getDockerGatewayIP attempts to get the Docker bridge network gateway IP
// which containers can use to reach the host when host.docker.internal doesn't work
func getDockerGatewayIP() string {
	// Try to get gateway IP from Docker bridge network
	cmd := exec.Command("docker", "network", "inspect", "bridge", "--format", "{{range .IPAM.Config}}{{.Gateway}}{{end}}")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	gatewayIP := strings.TrimSpace(string(output))
	if gatewayIP != "" && gatewayIP != "<no value>" {
		return gatewayIP
	}
	return ""
}

// extractLLOPair extracts the trading pair from an LLO price request
func extractLLOPair(data map[string]interface{}, streamIDToPair map[int]string) string {
	// Try base/quote format
	if base, ok := data["base"].(string); ok {
		if quote, ok := data["quote"].(string); ok {
			return base + "/" + quote
		}
	}

	// Try pair format
	if pair, ok := data["pair"].(string); ok {
		return pair
	}

	// Try from/to format
	if from, ok := data["from"].(string); ok {
		if to, ok := data["to"].(string); ok {
			return from + "/" + to
		}
	}

	// Try streamId format (used by LLO)
	if streamID, ok := data["streamId"]; ok {
		switch id := streamID.(type) {
		case float64:
			if pair, ok := streamIDToPair[int(id)]; ok {
				return pair
			}
		case int:
			if pair, ok := streamIDToPair[id]; ok {
				return pair
			}
		case string:
			var intID int
			if _, err := fmt.Sscanf(id, "%d", &intID); err == nil {
				if pair, ok := streamIDToPair[intID]; ok {
					return pair
				}
			}
		}
	}

	return ""
}
