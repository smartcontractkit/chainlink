package connectivity

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
)

type Placement string

const (
	PlacementLocal  Placement = "local"
	PlacementRemote Placement = "remote"
)

type EndpointPair struct {
	Name     string
	Internal string
	External string
}

type Resolution struct {
	URL            string
	SelectedKind   string
	RequiresBridge bool
	BridgePort     int
}

type BridgeEnsurer func(ctx context.Context, endpoint EndpointPair, port int) error

func Resolve(caller, target Placement, endpoint EndpointPair) (*Resolution, error) {
	if caller == "" || target == "" {
		return nil, fmt.Errorf("caller and target placement must be set")
	}

	selectedKind := "internal"
	selectedURL := strings.TrimSpace(endpoint.Internal)
	if caller != target {
		selectedKind = "external"
		selectedURL = strings.TrimSpace(endpoint.External)
	}
	if selectedURL == "" {
		return nil, fmt.Errorf("missing %s url for endpoint %q", selectedKind, endpoint.Name)
	}

	res := &Resolution{URL: selectedURL, SelectedKind: selectedKind}
	if caller == PlacementRemote && target == PlacementLocal {
		port, err := endpointPort(selectedURL)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve bridge port for endpoint %q: %w", endpoint.Name, err)
		}
		res.RequiresBridge = true
		res.BridgePort = port
	}
	return res, nil
}

func ResolveAndEnsureReachable(
	ctx context.Context,
	caller, target Placement,
	endpoint EndpointPair,
	ensureBridge BridgeEnsurer,
) (*Resolution, error) {
	res, err := Resolve(caller, target, endpoint)
	if err != nil {
		return nil, err
	}
	if !res.RequiresBridge {
		return res, nil
	}
	if ensureBridge == nil {
		return nil, fmt.Errorf("bridge required for endpoint %q (remote caller -> local target) but no bridge ensurer was provided", endpoint.Name)
	}
	if err := ensureBridge(ctx, endpoint, res.BridgePort); err != nil {
		return nil, fmt.Errorf("ensure bridge for endpoint %q on port %d: %w", endpoint.Name, res.BridgePort, err)
	}
	return res, nil
}

func PlacementFromTarget(target string) (Placement, error) {
	switch strings.ToLower(strings.TrimSpace(target)) {
	case "", "local":
		return PlacementLocal, nil
	case "remote":
		return PlacementRemote, nil
	default:
		return "", fmt.Errorf("unsupported component target %q", target)
	}
}

func endpointPort(raw string) (int, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return 0, fmt.Errorf("endpoint is empty")
	}
	if strings.Contains(trimmed, "://") {
		parsed, err := url.Parse(trimmed)
		if err != nil {
			return 0, fmt.Errorf("parse url: %w", err)
		}
		if parsed.Port() == "" {
			return 0, fmt.Errorf("url has no explicit port")
		}
		port, err := strconv.Atoi(parsed.Port())
		if err != nil || port <= 0 || port > 65535 {
			return 0, fmt.Errorf("invalid port %q", parsed.Port())
		}
		return port, nil
	}
	_, portRaw, err := net.SplitHostPort(trimmed)
	if err != nil {
		return 0, fmt.Errorf("parse host:port: %w", err)
	}
	port, err := strconv.Atoi(portRaw)
	if err != nil || port <= 0 || port > 65535 {
		return 0, fmt.Errorf("invalid port %q", portRaw)
	}
	return port, nil
}
