package adapters

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/tunnel"
)

type BlockchainAdapter struct{}

func NewBlockchainAdapter() *BlockchainAdapter {
	return &BlockchainAdapter{}
}

func (a *BlockchainAdapter) DescribeEndpoints(componentID string, output *blockchain.Output) ([]tunnel.EndpointRef, error) {
	if output == nil {
		return nil, fmt.Errorf("blockchain output is nil")
	}

	refs := make([]tunnel.EndpointRef, 0, len(output.Nodes)*2)
	for idx := range output.Nodes {
		node := output.Nodes[idx]

		httpRef, err := endpointFromURL(componentID, fmt.Sprintf("node-%d-http", idx), node.ExternalHTTPUrl)
		if err != nil {
			return nil, err
		}
		if httpRef != nil {
			refs = append(refs, *httpRef)
		}

		wsRef, err := endpointFromURL(componentID, fmt.Sprintf("node-%d-ws", idx), node.ExternalWSUrl)
		if err != nil {
			return nil, err
		}
		if wsRef != nil {
			refs = append(refs, *wsRef)
		}
	}

	return refs, nil
}

func (a *BlockchainAdapter) RewriteWithBindings(output *blockchain.Output, bindings []tunnel.TunnelBinding) error {
	if output == nil {
		return fmt.Errorf("blockchain output is nil")
	}

	byName := make(map[string]tunnel.TunnelBinding, len(bindings))
	for _, b := range bindings {
		byName[b.EndpointName] = b
	}

	for idx := range output.Nodes {
		httpKey := fmt.Sprintf("node-%d-http", idx)
		if output.Nodes[idx].ExternalHTTPUrl != "" {
			b, ok := byName[httpKey]
			if !ok {
				return fmt.Errorf("missing tunnel binding for %s", httpKey)
			}
			output.Nodes[idx].ExternalHTTPUrl = b.LocalURL
		}

		wsKey := fmt.Sprintf("node-%d-ws", idx)
		if output.Nodes[idx].ExternalWSUrl != "" {
			b, ok := byName[wsKey]
			if !ok {
				return fmt.Errorf("missing tunnel binding for %s", wsKey)
			}
			output.Nodes[idx].ExternalWSUrl = b.LocalURL
		}
	}

	return nil
}

func endpointFromURL(componentID, endpointName, rawURL string) (*tunnel.EndpointRef, error) {
	if rawURL == "" {
		return nil, nil
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse endpoint url %q: %w", rawURL, err)
	}

	host := parsed.Hostname()
	if host == "" {
		return nil, fmt.Errorf("endpoint url %q has empty hostname", rawURL)
	}

	port, err := resolvePort(parsed)
	if err != nil {
		return nil, err
	}

	return &tunnel.EndpointRef{
		ComponentID:  componentID,
		EndpointName: endpointName,
		Scheme:       parsed.Scheme,
		Host:         host,
		Port:         port,
		OriginalURL:  rawURL,
	}, nil
}

func resolvePort(parsed *url.URL) (int, error) {
	if parsed.Port() != "" {
		port, err := strconv.Atoi(parsed.Port())
		if err != nil || port <= 0 || port > 65535 {
			return 0, fmt.Errorf("url %q has invalid port %q", parsed.String(), parsed.Port())
		}
		return port, nil
	}

	switch parsed.Scheme {
	case "http", "ws":
		return 80, nil
	case "https", "wss":
		return 443, nil
	default:
		return 0, fmt.Errorf("url %q has unsupported scheme %q without explicit port", parsed.String(), parsed.Scheme)
	}
}
