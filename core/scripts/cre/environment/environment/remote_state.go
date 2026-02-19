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
)

type remoteStopState struct {
	Version     int                     `toml:"version"`
	Blockchains []*envconfig.Blockchain `toml:"blockchains"`
	NodeSets    []*cre.NodeSet          `toml:"nodesets"`
	JD          *envconfig.JobDistributor `toml:"jd"`
	Agent       remoteAgentState      `toml:"agent"`
}

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

func loadRemoteStopState(relativePathToRepoRoot string) (*remoteStopState, error) {
	data, err := os.ReadFile(remoteStateFileAbsPath(relativePathToRepoRoot))
	if err != nil {
		return nil, err
	}
	state := &remoteStopState{}
	if err := toml.Unmarshal(data, state); err != nil {
		return nil, err
	}
	if state.Version == 0 {
		state.Version = 1
	}
	if state.Blockchains == nil {
		state.Blockchains = []*envconfig.Blockchain{}
	}
	if state.NodeSets == nil {
		state.NodeSets = []*cre.NodeSet{}
	}
	return state, nil
}

func loadRemoteAgentState(relativePathToRepoRoot string) (*remoteAgentState, error) {
	data, err := os.ReadFile(remoteStateFileAbsPath(relativePathToRepoRoot))
	if err != nil {
		return nil, err
	}
	envelope := &remoteAgentStateEnvelope{}
	if err := toml.Unmarshal(data, envelope); err != nil {
		return nil, err
	}
	return &envelope.Agent, nil
}

func (s *remoteStopState) Config() *envconfig.Config {
	if s == nil {
		return nil
	}
	return &envconfig.Config{
		Blockchains: s.Blockchains,
		NodeSets:    s.NodeSets,
		JD:          s.JD,
	}
}

func storeRemoteStopState(relativePathToRepoRoot string, cfg *envconfig.Config) error {
	if cfg == nil {
		return fmt.Errorf("cannot store nil remote stop config")
	}
	state := &remoteStopState{
		Version:     1,
		Blockchains: []*envconfig.Blockchain{},
		NodeSets:    []*cre.NodeSet{},
		Agent: remoteAgentState{
			Mode:          os.Getenv("CRE_AGENT_MODE"),
			LocalURL:      os.Getenv("CRE_LOCAL_AGENT_URL"),
			EC2URL:        os.Getenv("CRE_EC2_AGENT_URL"),
			EC2InstanceID: os.Getenv("CRE_EC2_INSTANCE_ID"),
			EC2AgentPort:  os.Getenv("CRE_EC2_AGENT_PORT"),
			AWSProfile:    firstNonEmpty(os.Getenv("CRE_AWS_PROFILE"), os.Getenv("AWS_PROFILE")),
		},
	}
	for _, configuredBlockchain := range cfg.Blockchains {
		if configuredBlockchain != nil && configuredBlockchain.Target == envconfig.TargetRemote {
			state.Blockchains = append(state.Blockchains, configuredBlockchain)
		}
	}
	for _, nodeSet := range cfg.NodeSets {
		if nodeSet != nil && strings.TrimSpace(nodeSet.Target) == string(envconfig.TargetRemote) {
			state.NodeSets = append(state.NodeSets, nodeSet)
		}
	}
	if cfg.JD != nil && cfg.JD.Target == envconfig.TargetRemote {
		state.JD = cfg.JD
	}
	data, err := toml.Marshal(state)
	if err != nil {
		return err
	}
	path := remoteStateFileAbsPath(relativePathToRepoRoot)
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
	path := remoteStateFileAbsPath(relativePathToRepoRoot)
	err := os.Remove(path)
	if err == nil || os.IsNotExist(err) {
		return nil
	}
	return err
}
