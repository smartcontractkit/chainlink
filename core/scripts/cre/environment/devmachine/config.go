package devmachine

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// Defaults used when the user does not supply a value.
const (
	DefaultNamespace    = "cre-staggered-rollout-workflow"
	DefaultRegion       = "us-west-2"
	DefaultInstanceType = "m8a.2xlarge"
	DefaultVolumeSize   = 200
	DefaultAZ           = "a"
	DefaultAMI          = "amazon-linux-2023"
	DefaultAgentPort    = 18080
	DefaultJDImage      = "job-distributor-amd:0.22.1"
	DefaultJDLoadTag    = "job-distributor:0.22.1"
	DefaultJDTag        = "v0.22.1"
	DefaultCLImage      = "chainlink-amd:latest"
	DefaultSSHKey       = "~/.ssh/id_ed25519"
	DefaultEnvWorkdir   = "core/scripts/cre/environment"
	DefaultEnvConfig    = "configs/workflow-gateway-don-mixed-placement.toml"
	DefaultRemoteUser   = "ec2-user"
)

// Config is both the persisted tool configuration and the run state that is
// checkpointed after each phase.
type Config struct {
	AWS     AWSConfig     `yaml:"aws"`
	Cluster ClusterConfig `yaml:"cluster"`
	Machine MachineConfig `yaml:"machine"`
	SSH     SSHConfig     `yaml:"ssh"`
	Repos   ReposConfig   `yaml:"repos"`
	Images  ImagesConfig  `yaml:"images"`
	Agent   AgentConfig   `yaml:"agent"`
	Env     EnvConfig     `yaml:"env"`
}

type AWSConfig struct {
	Profile string `yaml:"profile"`
	Region  string `yaml:"region"`
}

type ClusterConfig struct {
	Namespace  string `yaml:"namespace"`
	Privileged bool   `yaml:"privileged"`
}

type MachineConfig struct {
	Name          string `yaml:"name"`
	InstanceID    string `yaml:"instanceId"`
	PrivateIP     string `yaml:"privateIp"`
	CreatedByTool bool   `yaml:"createdByTool"`
	InstanceType  string `yaml:"instanceType"`
	VolumeSize    int    `yaml:"volumeSize"`
	AvailabilityZ string `yaml:"availabilityZone"`
	AMI           string `yaml:"ami"`
	Team          string `yaml:"team"`
	CostCenter    string `yaml:"costCenter"`
}

type SSHConfig struct {
	Alias string `yaml:"alias"`
	Key   string `yaml:"key"`
}

type ReposConfig struct {
	Chainlink      string `yaml:"chainlink"`
	JobDistributor string `yaml:"jobDistributor"`
}

type ImagesConfig struct {
	Chainlink      string `yaml:"chainlink"`
	JobDistributor string `yaml:"jobDistributor"`
	JDLoadTag      string `yaml:"jdLoadTag"`
	JDTag          string `yaml:"jdTag"`
}

type AgentConfig struct {
	RemotePath string `yaml:"remotePath"`
	Addr       string `yaml:"addr"`
	Port       int    `yaml:"port"`
}

type EnvConfig struct {
	Workdir string `yaml:"workdir"`
	Config  string `yaml:"config"`
}

// DefaultConfig returns a config populated with all built-in defaults.
func DefaultConfig() *Config {
	return &Config{
		AWS: AWSConfig{
			Region: DefaultRegion,
		},
		Cluster: ClusterConfig{
			Namespace: DefaultNamespace,
		},
		Machine: MachineConfig{
			InstanceType:  DefaultInstanceType,
			VolumeSize:    DefaultVolumeSize,
			AvailabilityZ: DefaultAZ,
			AMI:           DefaultAMI,
		},
		SSH: SSHConfig{
			Key: DefaultSSHKey,
		},
		Images: ImagesConfig{
			Chainlink:      DefaultCLImage,
			JobDistributor: DefaultJDImage,
			JDLoadTag:      DefaultJDLoadTag,
			JDTag:          DefaultJDTag,
		},
		Agent: AgentConfig{
			RemotePath: "/home/" + DefaultRemoteUser + "/remote-agent",
			Addr:       fmt.Sprintf("0.0.0.0:%d", DefaultAgentPort),
			Port:       DefaultAgentPort,
		},
		Env: EnvConfig{
			Workdir: DefaultEnvWorkdir,
			Config:  DefaultEnvConfig,
		},
	}
}

// ConfigDir returns ~/.config/cre.
func ConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrap(err, "failed to resolve home directory")
	}
	return filepath.Join(home, ".config", "cre"), nil
}

// DefaultConfigPath returns ~/.config/cre/devmachine.yaml.
func DefaultConfigPath() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "devmachine.yaml"), nil
}

// CacheDir returns ~/.cache/cre/devmachine.
func CacheDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrap(err, "failed to resolve home directory")
	}
	return filepath.Join(home, ".cache", "cre", "devmachine"), nil
}

// LoadConfig reads the config at path. A missing file is not an error: the
// built-in defaults are returned instead.
func LoadConfig(path string) (*Config, error) {
	cfg := DefaultConfig()
	if path == "" {
		var err error
		path, err = DefaultConfigPath()
		if err != nil {
			return nil, err
		}
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, errors.Wrapf(err, "failed to read config %s", path)
	}
	if err := yaml.Unmarshal(raw, cfg); err != nil {
		return nil, errors.Wrapf(err, "failed to parse config %s", path)
	}
	cfg.applyDefaults()
	return cfg, nil
}

// Save writes the config to path, creating parent directories as needed.
func (c *Config) Save(path string) error {
	if path == "" {
		var err error
		path, err = DefaultConfigPath()
		if err != nil {
			return err
		}
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return errors.Wrapf(err, "failed to create config directory for %s", path)
	}
	raw, err := yaml.Marshal(c)
	if err != nil {
		return errors.Wrap(err, "failed to marshal config")
	}
	if err := os.WriteFile(path, raw, 0o600); err != nil {
		return errors.Wrapf(err, "failed to write config %s", path)
	}
	return nil
}

// applyDefaults fills any empty fields with built-in defaults.
func (c *Config) applyDefaults() {
	def := DefaultConfig()
	if c.AWS.Region == "" {
		c.AWS.Region = def.AWS.Region
	}
	if c.Cluster.Namespace == "" {
		c.Cluster.Namespace = def.Cluster.Namespace
	}
	if c.Machine.InstanceType == "" {
		c.Machine.InstanceType = def.Machine.InstanceType
	}
	if c.Machine.VolumeSize == 0 {
		c.Machine.VolumeSize = def.Machine.VolumeSize
	}
	if c.Machine.AvailabilityZ == "" {
		c.Machine.AvailabilityZ = def.Machine.AvailabilityZ
	}
	if c.Machine.AMI == "" {
		c.Machine.AMI = def.Machine.AMI
	}
	if c.SSH.Key == "" {
		c.SSH.Key = def.SSH.Key
	}
	if c.Images.Chainlink == "" {
		c.Images.Chainlink = def.Images.Chainlink
	}
	if c.Images.JobDistributor == "" {
		c.Images.JobDistributor = def.Images.JobDistributor
	}
	if c.Images.JDLoadTag == "" {
		c.Images.JDLoadTag = def.Images.JDLoadTag
	}
	if c.Images.JDTag == "" {
		c.Images.JDTag = def.Images.JDTag
	}
	if c.Agent.RemotePath == "" {
		c.Agent.RemotePath = def.Agent.RemotePath
	}
	if c.Agent.Port == 0 {
		c.Agent.Port = def.Agent.Port
	}
	if c.Agent.Addr == "" {
		c.Agent.Addr = def.Agent.Addr
	}
	if c.Env.Workdir == "" {
		c.Env.Workdir = def.Env.Workdir
	}
	if c.Env.Config == "" {
		c.Env.Config = def.Env.Config
	}
}

// SanitizeName converts a user identifier into a DNS-1123 compatible suffix
// suitable for use in a DevMachine name.
func SanitizeName(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	lastDash := false
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
			lastDash = false
		default:
			if !lastDash {
				b.WriteByte('-')
				lastDash = true
			}
		}
	}
	return strings.Trim(b.String(), "-")
}

// MachineName builds the default DevMachine name for a user (e.g. the local
// username). The result is truncated to fit the 63-char DNS label limit.
func MachineName(user string) string {
	name := "hybrid-local-cre-" + SanitizeName(user)
	if len(name) > 63 {
		name = strings.TrimRight(name[:63], "-")
	}
	return name
}

// ExpandPath expands a leading ~ to the user's home directory.
func ExpandPath(p string) string {
	if p == "" || p[0] != '~' {
		return p
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return p
	}
	if len(p) == 1 {
		return home
	}
	if p[1] == '/' {
		return filepath.Join(home, p[2:])
	}
	return p
}
