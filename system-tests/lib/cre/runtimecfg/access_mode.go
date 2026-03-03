package runtimecfg

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

const (
	EnvRemoteHostIP             = "CRE_REMOTE_HOST_IP"
	EnvLocalHostIP              = "CRE_LOCAL_HOST_IP"
	EnvRemoteAgentEC2InstanceID = "CRE_REMOTE_AGENT_EC2_INSTANCE_ID"

	defaultEC2Region = "us-west-2"
)

// IsDirectMode is retained for compatibility; CRE now only supports direct mode.
func IsDirectMode() bool {
	return true
}

func DirectHostIP() (string, error) {
	hostIP := strings.TrimSpace(os.Getenv(EnvRemoteHostIP))
	if hostIP != "" {
		return hostIP, nil
	}

	instanceID := strings.TrimSpace(os.Getenv(EnvRemoteAgentEC2InstanceID))
	if instanceID == "" {
		return "", fmt.Errorf("%s must be set (or set %s explicitly)", EnvRemoteAgentEC2InstanceID, EnvRemoteHostIP)
	}
	return discoverEC2HostIP(instanceID)
}

func LocalHostIP() string {
	raw := strings.TrimSpace(os.Getenv(EnvLocalHostIP))
	if ip := net.ParseIP(raw); ip != nil {
		return ip.String()
	}
	// Best-effort ensure the default CTF network exists before gateway discovery.
	// This avoids startup-order coupling where announce resolution runs before first container start.
	_ = framework.DefaultNetwork(nil)
	if gatewayIP := discoverDockerNetworkGatewayIP(framework.DefaultNetworkName); gatewayIP != "" {
		return gatewayIP
	}
	resolver := net.Resolver{}
	addrs, err := resolver.LookupIPAddr(context.Background(), "host.docker.internal")
	if err != nil {
		return ""
	}
	for _, addr := range addrs {
		if ipv4 := addr.IP.To4(); ipv4 != nil {
			return ipv4.String()
		}
	}
	return ""
}

func discoverDockerNetworkGatewayIP(networkName string) string {
	name := strings.TrimSpace(networkName)
	if name == "" {
		return ""
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	out, err := exec.CommandContext(ctx, "docker", "network", "inspect", name).Output()
	if err != nil {
		return ""
	}
	var inspected []struct {
		IPAM struct {
			Config []struct {
				Gateway string `json:"Gateway"`
			} `json:"Config"`
		} `json:"IPAM"`
	}
	if jsonErr := json.Unmarshal(out, &inspected); jsonErr != nil {
		return ""
	}
	for _, netCfg := range inspected {
		for _, ipamCfg := range netCfg.IPAM.Config {
			if ip := net.ParseIP(strings.TrimSpace(ipamCfg.Gateway)); ip != nil && ip.To4() != nil {
				return ip.String()
			}
		}
	}
	return ""
}

func ResolveAWSCLIProfileSelection() (string, string) {
	if hasStaticAWSKeys() {
		return "", "env-creds"
	}
	if hasWebIdentityCreds() {
		return "", "web-identity"
	}
	if profile := strings.TrimSpace(os.Getenv("AWS_PROFILE")); profile != "" {
		return profile, "profile:AWS_PROFILE"
	}
	if profile := strings.TrimSpace(os.Getenv("AWS_DEFAULT_PROFILE")); profile != "" {
		return profile, "profile:AWS_DEFAULT_PROFILE"
	}
	return "", "default-profile"
}

func hasStaticAWSKeys() bool {
	return strings.TrimSpace(os.Getenv("AWS_ACCESS_KEY_ID")) != "" && strings.TrimSpace(os.Getenv("AWS_SECRET_ACCESS_KEY")) != ""
}

func hasWebIdentityCreds() bool {
	return strings.TrimSpace(os.Getenv("AWS_WEB_IDENTITY_TOKEN_FILE")) != "" && strings.TrimSpace(os.Getenv("AWS_ROLE_ARN")) != ""
}

func awsRegion() string {
	if region := strings.TrimSpace(os.Getenv("AWS_REGION")); region != "" {
		return region
	}
	if region := strings.TrimSpace(os.Getenv("AWS_DEFAULT_REGION")); region != "" {
		return region
	}
	return defaultEC2Region
}

func discoverEC2HostIP(instanceID string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	profile, authMode := ResolveAWSCLIProfileSelection()
	args := []string{
		"ec2", "describe-instances",
		"--instance-ids", instanceID,
		"--region", awsRegion(),
		"--query", "Reservations[0].Instances[0].[PrivateIpAddress,PublicIpAddress]",
		"--output", "text",
	}
	if profile != "" {
		args = append(args, "--profile", profile)
	}

	out, err := exec.CommandContext(ctx, "aws", args...).CombinedOutput()
	if err != nil {
		msg := strings.TrimSpace(string(out))
		if msg == "" {
			msg = err.Error()
		}
		return "", fmt.Errorf("failed to resolve EC2 host IP via aws cli (auth mode=%s, instance=%s): %s", authMode, instanceID, msg)
	}

	parts := strings.Fields(strings.TrimSpace(string(out)))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" || strings.EqualFold(part, "none") {
			continue
		}
		return part, nil
	}
	return "", fmt.Errorf("no private/public IP found for instance %s", instanceID)
}
