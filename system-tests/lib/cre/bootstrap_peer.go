package cre

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/connectivity"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/runtimecfg"
)

func ResolveBootstrapAddress(callerTarget, bootstrapTarget, internalHost string, port int) (string, error) {
	if strings.TrimSpace(internalHost) == "" {
		return "", fmt.Errorf("bootstrap internal host is empty")
	}
	if port <= 0 || port > 65535 {
		return "", fmt.Errorf("invalid bootstrap port: %d", port)
	}

	callerPlacement, err := connectivity.PlacementFromTarget(callerTarget)
	if err != nil {
		return "", err
	}
	targetPlacement, err := connectivity.PlacementFromTarget(bootstrapTarget)
	if err != nil {
		return "", err
	}

	internal := net.JoinHostPort(strings.TrimSpace(internalHost), strconv.Itoa(port))
	external, err := resolveBootstrapExternalAddress(targetPlacement, port)
	if err != nil {
		return "", err
	}

	resolved, err := connectivity.Resolve(callerPlacement, targetPlacement, connectivity.EndpointPair{
		Name:     "ocr-bootstrap",
		Internal: internal,
		External: external,
	})
	if err != nil {
		return "", err
	}
	if !resolved.RequiresBridge {
		return resolved.URL, nil
	}
	return rewriteEndpointForRemoteCaller(resolved.URL)
}

func ResolveBootstrapPeerURL(callerTarget, bootstrapTarget, peerID, internalHost string, port int) (string, error) {
	address, err := ResolveBootstrapAddress(callerTarget, bootstrapTarget, internalHost, port)
	if err != nil {
		return "", err
	}
	trimmedPeerID := strings.TrimSpace(strings.TrimPrefix(peerID, "p2p_"))
	if trimmedPeerID == "" {
		return "", fmt.Errorf("bootstrap peerID is empty")
	}
	return trimmedPeerID + "@" + address, nil
}

func ResolveP2PAnnounceAddresses(nodePlacement string, hasRemoteNodeSets bool, advertisedPort int) ([]string, error) {
	if advertisedPort <= 0 || advertisedPort > 65535 {
		return nil, fmt.Errorf("invalid p2p announce port: %d", advertisedPort)
	}

	placement, err := connectivity.PlacementFromTarget(nodePlacement)
	if err != nil {
		return nil, err
	}

	addresses := []string{}
	add := func(addr string) {
		trimmed := strings.TrimSpace(addr)
		if trimmed != "" {
			addresses = append(addresses, trimmed)
		}
	}

	switch placement {
	case connectivity.PlacementLocal:
		if !hasRemoteNodeSets {
			// Keep AnnounceAddresses unset for local-only setups.
			// This lets libocr use default/local-network behavior.
			return addresses, nil
		}
		// In mixed mode, local nodes are reached through EC2 relay listeners.
		external, externalErr := resolveBootstrapExternalAddress(connectivity.PlacementRemote, advertisedPort)
		if externalErr != nil {
			return nil, externalErr
		}
		add(external)
	case connectivity.PlacementRemote:
		// Remote nodes advertise direct EC2-reachable host ports.
		external, externalErr := resolveBootstrapExternalAddress(connectivity.PlacementRemote, advertisedPort)
		if externalErr != nil {
			return nil, externalErr
		}
		add(external)
	default:
		return nil, fmt.Errorf("unsupported node placement: %s", nodePlacement)
	}

	return addresses, nil
}

func resolveBootstrapExternalAddress(targetPlacement connectivity.Placement, port int) (string, error) {
	if targetPlacement == connectivity.PlacementLocal {
		return net.JoinHostPort("127.0.0.1", strconv.Itoa(port)), nil
	}
	if !runtimecfg.IsDirectMode() {
		return "", fmt.Errorf("mixed DON bootstrap resolution requires direct access mode for remote bootstrap targets")
	}
	hostIP, err := runtimecfg.DirectHostIP()
	if err != nil {
		return "", err
	}
	return net.JoinHostPort(hostIP, strconv.Itoa(port)), nil
}

func rewriteEndpointForRemoteCaller(raw string) (string, error) {
	dockerHost := strings.TrimPrefix(framework.HostDockerInternal(), "http://")
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", fmt.Errorf("endpoint is empty")
	}
	if strings.Contains(trimmed, "://") {
		parsed, err := url.Parse(trimmed)
		if err != nil {
			return "", fmt.Errorf("parse url %q: %w", raw, err)
		}
		if parsed.Port() != "" {
			parsed.Host = net.JoinHostPort(dockerHost, parsed.Port())
			return parsed.String(), nil
		}
		parsed.Host = dockerHost
		return parsed.String(), nil
	}
	_, port, err := net.SplitHostPort(trimmed)
	if err != nil {
		return "", fmt.Errorf("parse host:port %q: %w", raw, err)
	}
	return net.JoinHostPort(dockerHost, port), nil
}
