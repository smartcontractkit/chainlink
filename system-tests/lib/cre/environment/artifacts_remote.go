package environment

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	pkgerrors "github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/agent"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/tunnel"
)

func DeployArtifactsToRemoteNodeSet(
	ctx context.Context,
	lggr zerolog.Logger,
	tunnelManager tunnel.Manager,
	nodeSetName string,
	containerTargetDir string,
	files []string,
) error {
	if nodeSetName == "" {
		return fmt.Errorf("nodeset name is required")
	}
	if containerTargetDir == "" {
		return fmt.Errorf("container target dir is required")
	}

	if tunnelManager == nil {
		return fmt.Errorf("tunnel manager is required for remote artifact deploy")
	}

	startClient, err := newStartComponentClient(lggr, tunnelManager)
	if err != nil {
		return pkgerrors.Wrap(err, "failed to initialize remote component client for artifact deploy")
	}

	payloadFiles := make([]agent.DeployArtifactsFile, 0, len(files))
	for _, path := range files {
		if path == "" {
			continue
		}
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return pkgerrors.Wrapf(readErr, "failed to read artifact file %s", path)
		}
		payloadFiles = append(payloadFiles, agent.DeployArtifactsFile{
			Name:          filepath.Base(path),
			ContentBase64: base64.StdEncoding.EncodeToString(data),
		})
	}
	if len(payloadFiles) == 0 {
		return fmt.Errorf("no artifact files to deploy")
	}

	payloadBytes, err := json.Marshal(agent.DeployArtifactsPayload{
		NodeSetName: nodeSetName,
		TargetDir:   containerTargetDir,
		Files:       payloadFiles,
	})
	if err != nil {
		return pkgerrors.Wrap(err, "failed to encode deploy artifacts payload")
	}

	response, err := startClient.StartComponent(ctx, agent.StartComponentEnvelope{
		SchemaVersion: agent.SchemaVersionV1,
		Operation:     agent.OperationDeployArtifacts,
		Payload:       payloadBytes,
	})
	if err != nil {
		return pkgerrors.Wrapf(err, "failed to deploy artifacts to remote nodeset %s", nodeSetName)
	}

	for _, logLine := range response.AgentLogs {
		pretty := prettifyAgentLogLine(logLine)
		if pretty == "" {
			continue
		}
		lggr.Info().Msgf("[agent] %s", pretty)
	}
	return nil
}
