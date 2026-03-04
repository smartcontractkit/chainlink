package environment

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/runtimecfg"
)

const (
	remoteStateDirname  = "core/scripts/cre/environment/state_remote"
	remoteStateFilename = "remote_components.toml"
	remoteAgentFilename = "remote_agent.toml"
	envRemoteAgentURL   = "CRE_REMOTE_AGENT_URL"
	envRemoteAgentPort  = "CRE_REMOTE_AGENT_PORT"
)

type remoteAgentState struct {
	RemoteAgentURL           string `toml:"remote_agent_url,omitempty"`
	RemoteAgentEC2InstanceID string `toml:"remote_agent_ec2_instance_id,omitempty"`
	RemoteAgentPort          string `toml:"remote_agent_port,omitempty"`
	AWSProfile               string `toml:"aws_profile,omitempty"`
}

type remoteAgentStateEnvelope struct {
	Agent remoteAgentState `toml:"agent"`
}

func remoteStateFileAbsPath(relativePathToRepoRoot string) string {
	absPath, err := filepath.Abs(filepath.Join(relativePathToRepoRoot, remoteStateDirname, remoteStateFilename))
	if err != nil {
		panic(fmt.Errorf("failed to get absolute path for remote CRE state file: %w", err))
	}
	return absPath
}

func remoteStateFileExists(relativePathToRepoRoot string) bool {
	_, statErr := os.Stat(remoteStateFileAbsPath(relativePathToRepoRoot))
	return statErr == nil
}

func loadRemoteStopConfig(relativePathToRepoRoot string) (*envconfig.Config, error) {
	data, err := os.ReadFile(remoteStateFileAbsPath(relativePathToRepoRoot))
	if err != nil {
		return nil, err
	}
	cfg := &envconfig.Config{}
	if err := toml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func loadRemoteAgentState(relativePathToRepoRoot string) (*remoteAgentState, error) {
	data, err := os.ReadFile(remoteAgentFileAbsPath(relativePathToRepoRoot))
	if err != nil {
		return nil, err
	}
	envelope := &remoteAgentStateEnvelope{}
	if err := toml.Unmarshal(data, envelope); err != nil {
		return nil, err
	}
	return &envelope.Agent, nil
}

func storeRemoteStopState(relativePathToRepoRoot string, cfg *envconfig.Config) error {
	if cfg == nil {
		return errors.New("cannot store nil remote stop config")
	}
	stopCfg := filteredRemoteStopConfig(cfg)
	if err := stopCfg.Store(remoteStateFileAbsPath(relativePathToRepoRoot)); err != nil {
		return err
	}
	agentEnvelope := &remoteAgentStateEnvelope{Agent: captureRemoteAgentState()}
	return storeRemoteAgentState(relativePathToRepoRoot, agentEnvelope)
}

func storeRemoteAgentStateSnapshot(relativePathToRepoRoot string) error {
	return storeRemoteAgentState(relativePathToRepoRoot, &remoteAgentStateEnvelope{Agent: captureRemoteAgentState()})
}

func filteredRemoteStopConfig(cfg *envconfig.Config) *envconfig.Config {
	stopCfg := &envconfig.Config{
		Blockchains: []*envconfig.Blockchain{},
		NodeSets:    []*cre.NodeSet{},
	}
	for _, configuredBlockchain := range cfg.Blockchains {
		if configuredBlockchain != nil && configuredBlockchain.Placement == envconfig.PlacementRemote {
			stopCfg.Blockchains = append(stopCfg.Blockchains, configuredBlockchain)
		}
	}
	for _, nodeSet := range cfg.NodeSets {
		if nodeSet != nil && strings.TrimSpace(nodeSet.Placement) == string(envconfig.PlacementRemote) {
			stopCfg.NodeSets = append(stopCfg.NodeSets, nodeSet)
		}
	}
	if cfg.JD != nil && cfg.JD.Placement == envconfig.PlacementRemote {
		stopCfg.JD = cfg.JD
	}
	return stopCfg
}

func captureRemoteAgentState() remoteAgentState {
	return remoteAgentState{
		RemoteAgentURL:           os.Getenv(envRemoteAgentURL),
		RemoteAgentEC2InstanceID: os.Getenv(runtimecfg.EnvRemoteAgentEC2InstanceID),
		RemoteAgentPort:          os.Getenv(envRemoteAgentPort),
		AWSProfile:               strings.TrimSpace(os.Getenv("AWS_PROFILE")),
	}
}

func storeRemoteAgentState(relativePathToRepoRoot string, envelope *remoteAgentStateEnvelope) error {
	data, err := toml.Marshal(envelope)
	if err != nil {
		return err
	}
	path := remoteAgentFileAbsPath(relativePathToRepoRoot)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			return trimmed
		}
	}
	return ""
}

func removeRemoteStopConfig(relativePathToRepoRoot string) error {
	stateDir, err := filepath.Abs(filepath.Join(relativePathToRepoRoot, remoteStateDirname))
	if err != nil {
		return err
	}
	return os.RemoveAll(stateDir)
}

func remoteAgentFileAbsPath(relativePathToRepoRoot string) string {
	absPath, err := filepath.Abs(filepath.Join(relativePathToRepoRoot, remoteStateDirname, remoteAgentFilename))
	if err != nil {
		panic(fmt.Errorf("failed to get absolute path for remote agent state file: %w", err))
	}
	return absPath
}
