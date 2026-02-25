package workflow

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDeployArtifactsRemoteValidation(t *testing.T) {
	err := DeployArtifacts(context.Background(), DeployArtifactsOptions{
		Mode: ArtifactDeployModeRemote,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "remote artifact deployer is required")

	err = DeployArtifacts(context.Background(), DeployArtifactsOptions{
		Mode:           ArtifactDeployModeRemote,
		RemoteDeployer: func(context.Context, string, string, []string) error { return nil },
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "nodeset name is required")
}

func TestDeployArtifactsRemoteCallsDeployer(t *testing.T) {
	called := false
	err := DeployArtifacts(context.Background(), DeployArtifactsOptions{
		Mode:        ArtifactDeployModeRemote,
		NodeSetName: "workflow",
		Files:       []string{"a.wasm"},
		RemoteDeployer: func(_ context.Context, nodeSetName, targetDir string, files []string) error {
			called = true
			require.Equal(t, "workflow", nodeSetName)
			require.Equal(t, "/home/chainlink/workflows", targetDir)
			require.Equal(t, []string{"a.wasm"}, files)
			return nil
		},
		ContainerTargetDir: "/home/chainlink/workflows",
	})
	require.NoError(t, err)
	require.True(t, called, "expected remote deployer to be called")
}

func TestDeployArtifactsLocalValidation(t *testing.T) {
	err := DeployArtifacts(context.Background(), DeployArtifactsOptions{
		Mode:                 ArtifactDeployModeLocal,
		ContainerTargetDir:   "/tmp",
		ContainerNamePattern: "",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "container name pattern is required")

	err = DeployArtifacts(context.Background(), DeployArtifactsOptions{
		Mode:               "",
		ContainerTargetDir: "/tmp",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "container name pattern is required")
}

func TestDeployArtifactsRemotePropagatesDeployerError(t *testing.T) {
	err := DeployArtifacts(context.Background(), DeployArtifactsOptions{
		Mode:        ArtifactDeployModeRemote,
		NodeSetName: "workflow",
		RemoteDeployer: func(context.Context, string, string, []string) error {
			return errors.New("deploy failed")
		},
	})
	require.EqualError(t, err, "deploy failed")
}
