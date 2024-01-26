package testconfig

import (
	"bytes"
	_ "embed"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/AlekSi/pointer"
	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-testing-framework/docker/test_env"
	"github.com/smartcontractkit/chainlink-testing-framework/networks"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/osutil"

	"github.com/smartcontractkit/chainlink-common/pkg/config"
	ctfconfig "github.com/smartcontractkit/chainlink-testing-framework/config"
)

const (
	ErrReadConfig             = "failed to read TOML config"
	ErrUnmarshalConfig        = "failed to unmarshal TOML config"
	Load               string = "load"
	Chaos              string = "chaos"
	Smoke              string = "smoke"
	ProductCCIP               = "CCIP"
)

var (
	//go:embed tomls/ccip-default.toml
	DefaultConfig []byte
)

func GlobalTestConfig() *Config {
	var err error
	cfg, err := NewConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}
	return cfg
}

// GenericConfig is an interface for all product based config types to implement
type GenericConfig interface {
	ReadSecrets() error
	Validate() error
	ApplyOverrides(from interface{}) error
}

// Config is the top level config struct. It contains config for all product based tests.
type Config struct {
	CCIP *CCIP `toml:",omitempty"`
}

func (c *Config) Validate() error {
	return c.CCIP.Validate()
}

func (c *Config) TOMLString() string {
	buf := new(bytes.Buffer)
	err := toml.NewEncoder(buf).Encode(c)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to encode config to TOML")
	}
	return buf.String()
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	override := &Config{}
	// load config from default file
	err := config.DecodeTOML(bytes.NewReader(DefaultConfig), cfg)
	if err != nil {
		return nil, errors.Wrap(err, ErrReadConfig)
	}

	// load config from env var if specified
	rawConfig, _ := osutil.GetEnv("BASE64_TEST_CONFIG_OVERRIDE")
	if rawConfig != "" {
		log.Info().Msg("Found BASE64_TEST_CONFIG_OVERRIDE env var, overriding default config")
		d, err := base64.StdEncoding.DecodeString(rawConfig)
		if err != nil {
			return nil, errors.Wrap(err, ErrReadConfig)
		}
		err = toml.Unmarshal(d, &override)
		if err != nil {
			return nil, errors.Wrap(err, ErrUnmarshalConfig)
		}
		log.Info().Interface("override", override).Msg("Applied overrides")
	}
	if override != nil {
		// apply overrides for all products
		if override.CCIP != nil {
			log.Info().Interface("override", override).Msg("Applying overrides for CCIP")
			if cfg.CCIP == nil {
				cfg.CCIP = override.CCIP
			} else {
				err = cfg.CCIP.ApplyOverrides(override.CCIP)
				if err != nil {
					return nil, err
				}
			}
		}
	}
	// read secrets for all products
	if cfg.CCIP != nil {
		err = cfg.CCIP.ReadSecrets()
		if err != nil {
			return nil, err
		}
		// validate all products
		err = cfg.CCIP.Validate()
		if err != nil {
			return nil, err
		}
	}
	fmt.Println("running test with config", cfg.TOMLString())
	return cfg, nil
}

// Common is the generic config struct which can be used with product specific configs.
// It contains generic DON and networks config which can be applied to all product based tests.
type Common struct {
	EnvUser   string                   `toml:",omitempty"`
	TTL       *config.Duration         `toml:",omitempty"`
	Chainlink *Chainlink               `toml:",omitempty"`
	Network   *ctfconfig.NetworkConfig `toml:",omitempty"`
	Logging   *ctfconfig.LoggingConfig `toml:"Logging"`
}

func (p *Common) ReadSecrets() error {
	return p.Chainlink.ReadSecrets()
}

func (p *Common) ApplyOverrides(from *Common) error {
	if from == nil {
		return nil
	}
	if from.EnvUser != "" {
		p.EnvUser = from.EnvUser
	}
	if from.TTL != nil {
		p.TTL = from.TTL
	}
	if from.Network != nil {
		p.Network = from.Network
	}
	if from.Chainlink != nil {
		if p.Chainlink == nil {
			p.Chainlink = &Chainlink{}
		}
		p.Chainlink.ApplyOverrides(from.Chainlink)
	}
	return nil
}

func (p *Common) Validate() error {
	if err := p.Logging.Validate(); err != nil {
		return fmt.Errorf("error validating logging config %w", err)
	}
	if p.Network == nil {
		return errors.New("no networks specified")
	}
	if err := p.Network.Validate(); err != nil {
		return fmt.Errorf("error validating networks config %w", err)
	}
	return p.Chainlink.Validate()
}

func (p *Common) EVMNetworks() ([]blockchain.EVMNetwork, error) {
	p.Network.UpperCaseNetworkNames()
	err := p.Network.Default()
	if err != nil {
		return nil, fmt.Errorf("error reading default network config %w", err)
	}
	return networks.MustSetNetworks(*p.Network), nil
}

func (p *Common) GetLoggingConfig() *ctfconfig.LoggingConfig {
	return p.Logging
}

func (p *Common) GetChainlinkImageConfig() *ctfconfig.ChainlinkImageConfig {
	return p.Chainlink.Common.ChainlinkImage
}

func (p *Common) GetPyroscopeConfig() *ctfconfig.PyroscopeConfig {
	return nil
}

func (p *Common) GetPrivateEthereumNetworkConfig() *test_env.EthereumNetwork {
	return nil
}

func (p *Common) GetNetworkConfig() *ctfconfig.NetworkConfig {
	return p.Network
}

// Returns Grafana URL from Logging config
func (p *Common) GetGrafanaBaseURL() (string, error) {
	if p.Logging.Grafana == nil || p.Logging.Grafana.BaseUrl == nil {
		return "", errors.New("grafana base url not set")
	}

	return strings.TrimSuffix(*p.Logging.Grafana.BaseUrl, "/"), nil
}

// Returns Grafana Dashboard URL from Logging config
func (p *Common) GetGrafanaDashboardURL() (string, error) {
	if p.Logging.Grafana == nil || p.Logging.Grafana.DashboardUrl == nil {
		return "", errors.New("grafana dashboard url not set")
	}

	url := *p.Logging.Grafana.DashboardUrl
	if !strings.HasPrefix(url, "/") {
		url = "/" + url
	}

	return url, nil
}

type Chainlink struct {
	Common     *Node    `toml:",omitempty"`
	NodeMemory string   `toml:",omitempty"`
	NodeCPU    string   `toml:",omitempty"`
	DBMemory   string   `toml:",omitempty"`
	DBCPU      string   `toml:",omitempty"`
	DBCapacity string   `toml:",omitempty"`
	IsStateful *bool    `toml:",omitempty"`
	DBArgs     []string `toml:",omitempty"`
	NoOfNodes  *int     `toml:",omitempty"`
	Nodes      []*Node  `toml:",omitempty"` // to be mentioned only if diff nodes follow diff configs; not required if all nodes follow CommonConfig
}

func (c *Chainlink) ApplyOverrides(from *Chainlink) {
	if from == nil {
		return
	}
	if from.NoOfNodes != nil {
		c.NoOfNodes = from.NoOfNodes
	}
	if from.Common != nil {
		c.Common.ApplyOverrides(from.Common)
	}
	if from.Nodes != nil {
		for i, node := range from.Nodes {
			if len(c.Nodes) > i {
				c.Nodes[i].ApplyOverrides(node)
			} else {
				c.Nodes = append(c.Nodes, node)
			}
			c.Nodes[i].Merge(c.Common)
		}
	}
	if from.NodeMemory != "" {
		c.NodeMemory = from.NodeMemory
	}
	if from.NodeCPU != "" {
		c.NodeCPU = from.NodeCPU
	}
	if from.DBMemory != "" {
		c.DBMemory = from.DBMemory
	}
	if from.DBCPU != "" {
		c.DBCPU = from.DBCPU
	}
	if from.DBArgs != nil {
		c.DBArgs = from.DBArgs
	}
	if from.DBCapacity != "" {
		c.DBCapacity = from.DBCapacity
	}
	if from.IsStateful != nil {
		c.IsStateful = from.IsStateful
	}
}

// TODO change this later to directly reading from toml instead of reading from env vars
func (c *Chainlink) ReadSecrets() error {
	image, _ := osutil.GetEnv("CHAINLINK_IMAGE")
	tag, _ := osutil.GetEnv("CHAINLINK_VERSION")
	if image != "" && tag != "" {
		c.Common.ChainlinkImage = &ctfconfig.ChainlinkImageConfig{
			Image:   &image,
			Version: &tag,
		}
	}

	upgradeImage, _ := osutil.GetEnv("UPGRADE_IMAGE")
	upgradeTag, _ := osutil.GetEnv("UPGRADE_VERSION")
	if upgradeImage != "" && upgradeTag != "" {
		c.Common.ChainlinkUpgradeImage = &ctfconfig.ChainlinkImageConfig{
			Image:   &upgradeImage,
			Version: &upgradeTag,
		}
	}

	for i := range c.Nodes {
		image, _ := osutil.GetEnv(fmt.Sprintf("CHAINLINK_IMAGE-%d", i+1))
		tag, _ := osutil.GetEnv(fmt.Sprintf("CHAINLINK_VERSION-%d", i+1))
		if image != "" && tag != "" {
			c.Nodes[i].ChainlinkImage = &ctfconfig.ChainlinkImageConfig{
				Image:   &image,
				Version: &tag,
			}
		} else {
			c.Nodes[i].ChainlinkImage = c.Common.ChainlinkImage
		}
		upgradeImage, _ := osutil.GetEnv(fmt.Sprintf("UPGRADE_IMAGE-%d", i+1))
		upgradeTag, _ := osutil.GetEnv(fmt.Sprintf("UPGRADE_VERSION-%d", i+1))
		if upgradeImage != "" && upgradeTag != "" {
			c.Nodes[i].ChainlinkUpgradeImage = &ctfconfig.ChainlinkImageConfig{
				Image:   &upgradeImage,
				Version: &upgradeTag,
			}
		} else {
			c.Nodes[i].ChainlinkUpgradeImage = c.Common.ChainlinkUpgradeImage
		}
	}
	return nil
}

func (c *Chainlink) Validate() error {
	if c.Common == nil {
		return errors.New("common config can't be empty")
	}
	if err := c.Common.ChainlinkImage.Validate(); err != nil {
		return err
	}
	if c.Common.DBImage == "" || c.Common.DBTag == "" {
		return errors.New("must provide db image and tag")
	}
	if c.NoOfNodes == nil {
		return errors.New("chainlink config is invalid, NoOfNodes should be specified")
	}
	if c.Nodes != nil && len(c.Nodes) > 0 {
		noOfNodes := pointer.GetInt(c.NoOfNodes)
		if noOfNodes != len(c.Nodes) {
			return errors.New("chainlink config is invalid, NoOfNodes and Nodes length mismatch")
		}
	}
	return nil
}

type Node struct {
	Name                   string                          `toml:",omitempty"`
	ChainlinkImage         *ctfconfig.ChainlinkImageConfig `toml:"ChainlinkImage"`
	ChainlinkUpgradeImage  *ctfconfig.ChainlinkImageConfig `toml:"ChainlinkUpgradeImage"`
	BaseConfigTOML         string                          `toml:",omitempty"`
	CommonChainConfigTOML  string                          `toml:",omitempty"`
	ChainConfigTOMLByChain map[string]string               `toml:",omitempty"` // key is chainID
	DBImage                string                          `toml:",omitempty"`
	DBTag                  string                          `toml:",omitempty"`
}

func (n *Node) Merge(from *Node) {
	if from == nil || n == nil {
		return
	}
	if n.Name == "" {
		n.Name = from.Name
	}
	if n.ChainlinkImage == nil {
		if from.ChainlinkImage != nil {
			n.ChainlinkImage = &ctfconfig.ChainlinkImageConfig{
				Image:   from.ChainlinkImage.Image,
				Version: from.ChainlinkImage.Version,
			}
		}
	} else {
		if n.ChainlinkImage.Image == nil && from.ChainlinkImage != nil {
			n.ChainlinkImage.Image = from.ChainlinkImage.Image
		}
		if n.ChainlinkImage.Version == nil && from.ChainlinkImage != nil {
			n.ChainlinkImage.Version = from.ChainlinkImage.Version
		}
	}

	if n.ChainlinkUpgradeImage == nil {
		if from.ChainlinkUpgradeImage != nil {
			n.ChainlinkUpgradeImage = &ctfconfig.ChainlinkImageConfig{
				Image:   from.ChainlinkUpgradeImage.Image,
				Version: from.ChainlinkUpgradeImage.Version,
			}
		}
	} else {
		if n.ChainlinkUpgradeImage.Image == nil && from.ChainlinkUpgradeImage != nil {
			n.ChainlinkUpgradeImage.Image = from.ChainlinkUpgradeImage.Image
		}
		if n.ChainlinkUpgradeImage.Version == nil && from.ChainlinkUpgradeImage != nil {
			n.ChainlinkUpgradeImage.Version = from.ChainlinkUpgradeImage.Version
		}
	}

	if n.DBImage == "" {
		n.DBImage = from.DBImage
	}
	if n.DBTag == "" {
		n.DBTag = from.DBTag
	}
	if n.BaseConfigTOML == "" {
		n.BaseConfigTOML = from.BaseConfigTOML
	}
	if n.CommonChainConfigTOML == "" {
		n.CommonChainConfigTOML = from.CommonChainConfigTOML
	}
	if n.ChainConfigTOMLByChain == nil {
		n.ChainConfigTOMLByChain = from.ChainConfigTOMLByChain
	} else {
		for k, v := range from.ChainConfigTOMLByChain {
			if _, ok := n.ChainConfigTOMLByChain[k]; !ok {
				n.ChainConfigTOMLByChain[k] = v
			}
		}
	}
}

func (n *Node) ApplyOverrides(from *Node) {
	if from == nil {
		return
	}
	if n == nil {
		return
	}
	if from.Name != "" {
		n.Name = from.Name
	}
	if from.ChainlinkImage != nil {
		if n.ChainlinkImage == nil {
			n.ChainlinkImage = from.ChainlinkImage
		} else {
			if from.ChainlinkImage.Image != nil {
				n.ChainlinkImage.Image = from.ChainlinkImage.Image
			}
			if from.ChainlinkImage.Version != nil {
				n.ChainlinkImage.Version = from.ChainlinkImage.Version
			}
		}
	}
	if from.ChainlinkUpgradeImage != nil {
		if n.ChainlinkUpgradeImage == nil {
			n.ChainlinkUpgradeImage = from.ChainlinkUpgradeImage
		} else {
			if from.ChainlinkUpgradeImage.Image != nil {
				n.ChainlinkUpgradeImage.Image = from.ChainlinkUpgradeImage.Image
			}
			if from.ChainlinkUpgradeImage.Version != nil {
				n.ChainlinkUpgradeImage.Version = from.ChainlinkUpgradeImage.Version
			}
		}
	}
	if from.DBImage != "" {
		n.DBImage = from.DBImage
	}
	if from.DBTag != "" {
		n.DBTag = from.DBTag
	}
	if from.BaseConfigTOML != "" {
		n.BaseConfigTOML = from.BaseConfigTOML
	}
	if from.CommonChainConfigTOML != "" {
		n.CommonChainConfigTOML = from.CommonChainConfigTOML
	}
	if from.ChainConfigTOMLByChain != nil {
		if n.ChainConfigTOMLByChain == nil {
			n.ChainConfigTOMLByChain = from.ChainConfigTOMLByChain
		} else {
			for chainID, cfg := range from.ChainConfigTOMLByChain {
				n.ChainConfigTOMLByChain[chainID] = cfg
			}
		}
	}
}
