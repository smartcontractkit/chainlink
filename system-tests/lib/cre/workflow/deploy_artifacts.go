package workflow

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/internal/dockerops"
)

type ArtifactDeployMode string

const (
	ArtifactDeployModeLocal  ArtifactDeployMode = "local"
	ArtifactDeployModeRemote ArtifactDeployMode = "remote"
)

type DeployArtifactsOptions struct {
	Mode                 ArtifactDeployMode
	NodeSetName          string
	ContainerNamePattern string
	ContainerTargetDir   string
	Files                []string
	RemoteDeployer       func(ctx context.Context, nodeSetName, containerTargetDir string, files []string) error
}

func DeployArtifacts(ctx context.Context, opts DeployArtifactsOptions) error {
	switch opts.Mode {
	case ArtifactDeployModeRemote:
		if opts.RemoteDeployer == nil {
			return fmt.Errorf("remote artifact deployer is required for mode=%s", opts.Mode)
		}
		if opts.NodeSetName == "" {
			return fmt.Errorf("nodeset name is required for mode=%s", opts.Mode)
		}
		return opts.RemoteDeployer(ctx, opts.NodeSetName, opts.ContainerTargetDir, opts.Files)
	case ArtifactDeployModeLocal:
		fallthrough
	default:
		if opts.ContainerNamePattern == "" {
			return fmt.Errorf("container name pattern is required for mode=%s", opts.Mode)
		}
		return dockerops.CopyFilesToContainers(ctx, opts.ContainerNamePattern, opts.ContainerTargetDir, opts.Files)
	}
}
