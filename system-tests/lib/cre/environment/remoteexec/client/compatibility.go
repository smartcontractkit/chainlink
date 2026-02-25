package client

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/remoteexec/agent"
)

const clientProtocolVersion = "1.0.0"

func CheckCompatibility(ctx context.Context, runtime *Runtime, requiredCapabilities []string) error {
	status, err := GetAgentStatus(ctx, runtime)
	if err != nil {
		return err
	}
	return checkCompatibilityStatus(status, requiredCapabilities)
}

func checkCompatibilityStatus(status *agent.AgentStatusResponse, requiredCapabilities []string) error {
	if status == nil {
		return fmt.Errorf("agent status is nil")
	}

	if strings.TrimSpace(status.ProtocolVersion) != "" {
		clientMajor, err := semverMajor(clientProtocolVersion)
		if err != nil {
			return err
		}
		agentMajor, err := semverMajor(status.ProtocolVersion)
		if err != nil {
			return fmt.Errorf("invalid agent protocolVersion %q: %w", status.ProtocolVersion, err)
		}
		if clientMajor != agentMajor {
			return fmt.Errorf("incompatible protocol major versions: client=%s agent=%s", clientProtocolVersion, status.ProtocolVersion)
		}
	}

	if len(requiredCapabilities) == 0 || len(status.Capabilities) == 0 {
		return nil
	}
	for _, required := range requiredCapabilities {
		normalized := strings.TrimSpace(required)
		if normalized == "" {
			continue
		}
		if !slices.Contains(status.Capabilities, normalized) {
			return fmt.Errorf("agent does not support required capability %q", normalized)
		}
	}
	return nil
}

func semverMajor(version string) (int, error) {
	parts := strings.Split(strings.TrimSpace(version), ".")
	if len(parts) == 0 || strings.TrimSpace(parts[0]) == "" {
		return 0, fmt.Errorf("invalid semver: %q", version)
	}
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("invalid semver major in %q: %w", version, err)
	}
	return major, nil
}
