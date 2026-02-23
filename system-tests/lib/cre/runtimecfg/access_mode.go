package runtimecfg

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	EnvRemoteAccessMode = "CRE_REMOTE_ACCESS_MODE"
	EnvEC2HostIP        = "CRE_EC2_HOST_IP"
	EnvEC2InstanceID    = "CRE_EC2_INSTANCE_ID"
	EnvAWSProfile       = "CRE_AWS_PROFILE"

	RemoteAccessModeSSM    = "ssm"
	RemoteAccessModeDirect = "direct"
	defaultEC2Region       = "us-west-2"
)

func RemoteAccessMode() string {
	mode := strings.ToLower(strings.TrimSpace(os.Getenv(EnvRemoteAccessMode)))
	if mode == "" {
		return RemoteAccessModeDirect
	}
	if mode == RemoteAccessModeDirect || mode == RemoteAccessModeSSM {
		return mode
	}
	return RemoteAccessModeDirect
}

func IsDirectMode() bool {
	return RemoteAccessMode() == RemoteAccessModeDirect
}

func DirectHostIP() (string, error) {
	hostIP := strings.TrimSpace(os.Getenv(EnvEC2HostIP))
	if hostIP != "" {
		return hostIP, nil
	}

	instanceID := strings.TrimSpace(os.Getenv(EnvEC2InstanceID))
	if instanceID == "" {
		return "", fmt.Errorf("%s must be set when %s=%s (or set %s explicitly)", EnvEC2InstanceID, EnvRemoteAccessMode, RemoteAccessModeDirect, EnvEC2HostIP)
	}
	return discoverEC2HostIP(instanceID)
}

func ResolveAWSCLIProfileSelection() (string, string) {
	if hasStaticAWSKeys() {
		return "", "env-creds"
	}
	if hasWebIdentityCreds() {
		return "", "web-identity"
	}
	if profile := strings.TrimSpace(os.Getenv(EnvAWSProfile)); profile != "" {
		return profile, "profile:CRE_AWS_PROFILE"
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
