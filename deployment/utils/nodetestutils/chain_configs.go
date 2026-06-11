package nodetestutils

import (
	"fmt"
	"net/url"
	"strings"

	cldf_aptos "github.com/smartcontractkit/chainlink-deployments-framework/chain/aptos"
	cldf_sui "github.com/smartcontractkit/chainlink-deployments-framework/chain/sui"
	cldf_tron "github.com/smartcontractkit/chainlink-deployments-framework/chain/tron"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

// defaultSuiGrpcToken matches the default token used by CLDF's NewPTBClientFromNodeURL so
// the DON relayer's gRPC client targets the same endpoint as the test's client.
const defaultSuiGrpcToken = "test"

// suiGrpcTargetFromNodeURL derives the gRPC target (host:port) from a Sui fullnode HTTP URL.
// The Sui fullnode serves the gRPC v2 API on the same host:port as JSON-RPC. This mirrors
// chainlink-deployments-framework/chain/sui.grpcTargetFromNodeURL so the DON relayer connects
// to the same endpoint CLDF's client uses.
func suiGrpcTargetFromNodeURL(nodeURL string) (string, error) {
	u, err := url.Parse(nodeURL)
	if err != nil {
		return "", fmt.Errorf("parse node URL %q: %w", nodeURL, err)
	}
	host := u.Hostname()
	if host == "" {
		return "", fmt.Errorf("node URL %q has no host", nodeURL)
	}
	port := u.Port()
	if port == "" {
		if u.Scheme == "https" {
			port = "443"
		} else {
			port = "9000"
		}
	}
	if strings.Contains(host, ":") && !strings.HasPrefix(host, "[") {
		return fmt.Sprintf("[%s]:%s", host, port), nil
	}

	return fmt.Sprintf("%s:%s", host, port), nil
}

// createSuiChainConfig creates a Chainlink node defined Sui chain configuration.
func createSuiChainConfig(chainID string, chain cldf_sui.Chain) chainlink.RawConfig {
	chainConfig := chainlink.RawConfig{}

	chainConfig["Enabled"] = true
	chainConfig["ChainID"] = chainID
	chainConfig["NetworkName"] = "sui-localnet"
	chainConfig["NetworkNameFull"] = "sui-localnet"

	node := map[string]any{
		"Name": "primary",
		"URL":  chain.URL,
	}
	// The migrated Sui relayer builds its gRPC client from the node's GrpcTarget/GrpcToken
	// (it has no URL-based fallback). Without these the relayer panics on a nil deref and the
	// DON never reads the Sui source, so commit reports carry no merkle roots.
	if grpcTarget, err := suiGrpcTargetFromNodeURL(chain.URL); err == nil {
		node["GrpcTarget"] = grpcTarget
		node["GrpcToken"] = defaultSuiGrpcToken
	}
	chainConfig["Nodes"] = []any{node}

	return chainConfig
}

// createTronChainConfig creates a Chainlink node defined Tron chain configuration.
func createTronChainConfig(chainID string, chain cldf_tron.Chain) chainlink.RawConfig {
	chainConfig := chainlink.RawConfig{}

	chainConfig["Enabled"] = true
	chainConfig["ChainID"] = chainID
	chainConfig["NetworkName"] = "tron-local"
	chainConfig["NetworkNameFull"] = "tron-local"
	chainConfig["Nodes"] = []any{
		map[string]any{
			"Name": "primary",
			"URL":  chain.URL,
		},
	}

	return chainConfig
}

// createAptosChainConfig creates a Chainlink node defined Aptos chain configuration.
func createAptosChainConfig(chainID string, chain cldf_aptos.Chain) chainlink.RawConfig {
	chainConfig := chainlink.RawConfig{}

	chainConfig["Enabled"] = true
	chainConfig["ChainID"] = chainID
	chainConfig["NetworkName"] = "localnet"
	chainConfig["NetworkNameFull"] = "aptos-localnet"
	chainConfig["Nodes"] = []any{
		map[string]any{
			"Name": "primary",
			"URL":  chain.URL,
		},
	}

	return chainConfig
}
