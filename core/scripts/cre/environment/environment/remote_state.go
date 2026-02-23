package environment

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
)

const (
	remoteStateDirname  = "core/scripts/cre/environment/state_remote"
	remoteStateFilename = "remote_components.toml"
	remoteAgentFilename = "remote_agent.toml"
)

type remoteAgentState struct {
	Mode         string `toml:"mode,omitempty"`
	LocalURL     string `toml:"local_url,omitempty"`
	EC2URL       string `toml:"ec2_url,omitempty"`
	EC2InstanceID string `toml:"ec2_instance_id,omitempty"`
	EC2AgentPort string `toml:"ec2_agent_port,omitempty"`
	AWSProfile   string `toml:"aws_profile,omitempty"`
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
	cfg := &envconfig.Config{}
	if err := cfg.Load(remoteStateFileAbsPath(relativePathToRepoRoot)); err != nil {
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
		return fmt.Errorf("cannot store nil remote stop config")
	}
	stopCfg := &envconfig.Config{
		Blockchains: []*envconfig.Blockchain{},
		NodeSets:    []*cre.NodeSet{},
	}
	for _, configuredBlockchain := range cfg.Blockchains {
		if configuredBlockchain != nil && configuredBlockchain.Target == envconfig.TargetRemote {
			stopCfg.Blockchains = append(stopCfg.Blockchains, configuredBlockchain)
		}
	}
	for _, nodeSet := range cfg.NodeSets {
		if nodeSet != nil && strings.TrimSpace(nodeSet.Target) == string(envconfig.TargetRemote) {
			stopCfg.NodeSets = append(stopCfg.NodeSets, nodeSet)
		}
	}
	if cfg.JD != nil && cfg.JD.Target == envconfig.TargetRemote {
		stopCfg.JD = cfg.JD
	}
	if err := stopCfg.Store(remoteStateFileAbsPath(relativePathToRepoRoot)); err != nil {
		return err
	}
	agentEnvelope := &remoteAgentStateEnvelope{
		Agent: remoteAgentState{
			Mode:          os.Getenv("CRE_AGENT_MODE"),
			LocalURL:      os.Getenv("CRE_LOCAL_AGENT_URL"),
			EC2URL:        os.Getenv("CRE_EC2_AGENT_URL"),
			EC2InstanceID: os.Getenv("CRE_EC2_INSTANCE_ID"),
			EC2AgentPort:  os.Getenv("CRE_EC2_AGENT_PORT"),
			AWSProfile:    firstNonEmpty(os.Getenv("CRE_AWS_PROFILE"), os.Getenv("AWS_PROFILE")),
		},
	}
	data, err := toml.Marshal(agentEnvelope)
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
	for _, path := range []string{
		remoteStateFileAbsPath(relativePathToRepoRoot),
		remoteAgentFileAbsPath(relativePathToRepoRoot),
	} {
		err := os.Remove(path)
		if err == nil || os.IsNotExist(err) {
			continue
		}
		return err
	}
	return nil
}

func remoteAgentFileAbsPath(relativePathToRepoRoot string) string {
	absPath, err := filepath.Abs(filepath.Join(relativePathToRepoRoot, remoteStateDirname, remoteAgentFilename))
	if err != nil {
		panic(fmt.Errorf("failed to get absolute path for remote agent state file: %w", err))
	}
	return absPath
}
