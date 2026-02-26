package runtimecfg

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDirectHostIPUsesExplicitEnv(t *testing.T) {
	t.Setenv(EnvRemoteHostIP, "203.0.113.10")
	t.Setenv(EnvRemoteAgentEC2InstanceID, "")

	hostIP, err := DirectHostIP()
	require.NoError(t, err)
	require.Equal(t, "203.0.113.10", hostIP)
}

func TestDirectHostIPRequiresInstanceWhenHostMissing(t *testing.T) {
	t.Setenv(EnvRemoteHostIP, "")
	t.Setenv(EnvRemoteAgentEC2InstanceID, "")

	_, err := DirectHostIP()
	require.Error(t, err)
	require.Contains(t, err.Error(), EnvRemoteAgentEC2InstanceID)
}

func TestLocalHostIPUsesExplicitEnv(t *testing.T) {
	t.Setenv(EnvLocalHostIP, "192.168.1.11")
	require.Equal(t, "192.168.1.11", LocalHostIP())
}

func TestResolveAWSCLIProfileSelectionOrder(t *testing.T) {
	t.Setenv("AWS_ACCESS_KEY_ID", "key")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "secret")
	profile, mode := ResolveAWSCLIProfileSelection()
	require.Equal(t, "", profile)
	require.Equal(t, "env-creds", mode)

	t.Setenv("AWS_ACCESS_KEY_ID", "")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "")
	t.Setenv("AWS_WEB_IDENTITY_TOKEN_FILE", "/tmp/token")
	t.Setenv("AWS_ROLE_ARN", "arn:aws:iam::123456789012:role/Role")
	profile, mode = ResolveAWSCLIProfileSelection()
	require.Equal(t, "", profile)
	require.Equal(t, "web-identity", mode)

	t.Setenv("AWS_WEB_IDENTITY_TOKEN_FILE", "")
	t.Setenv("AWS_ROLE_ARN", "")
	t.Setenv("AWS_PROFILE", "profile-b")
	t.Setenv("AWS_DEFAULT_PROFILE", "profile-c")
	profile, mode = ResolveAWSCLIProfileSelection()
	require.Equal(t, "profile-b", profile)
	require.Equal(t, "profile:AWS_PROFILE", mode)

	t.Setenv("AWS_PROFILE", "")
	profile, mode = ResolveAWSCLIProfileSelection()
	require.Equal(t, "profile-c", profile)
	require.Equal(t, "profile:AWS_DEFAULT_PROFILE", mode)
}

func TestAWSRegionResolution(t *testing.T) {
	t.Setenv("AWS_REGION", "eu-central-1")
	t.Setenv("AWS_DEFAULT_REGION", "us-east-1")
	require.Equal(t, "eu-central-1", awsRegion())

	t.Setenv("AWS_REGION", "")
	require.Equal(t, "us-east-1", awsRegion())

	t.Setenv("AWS_DEFAULT_REGION", "")
	require.Equal(t, defaultEC2Region, awsRegion())
}
