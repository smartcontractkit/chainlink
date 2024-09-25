package testconfig

import (
	"bytes"
	_ "embed"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/AlekSi/pointer"
	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/docker/test_env"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/networks"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/osutil"

	"github.com/smartcontractkit/chainlink-common/pkg/config"
	ctfconfig "github.com/smartcontractkit/chainlink-testing-framework/lib/config"

	"github.com/smartcontractkit/chainlink/integration-tests/client"
)

const (
	OVERIDECONFIG             = "BASE64_CONFIG_OVERRIDE"
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

func DecodeConfig(rawConfig string, c any) error {
	d, err := base64.StdEncoding.DecodeString(rawConfig)
	if err != nil {
		return errors.Wrap(err, ErrReadConfig)
	}
	err = toml.Unmarshal(d, c)
	if err != nil {
		return errors.Wrap(err, ErrUnmarshalConfig)
	}
	return nil
}

// EncodeConfigAndSetEnv encodes the given struct to base64
// and sets env var ( if not empty) with the encoded base64 string
func EncodeConfigAndSetEnv(c any, envVar string) (string, error) {
	srcBytes, err := toml.Marshal(c)
	if err != nil {
		return "", err
	}
	encodedStr := base64.StdEncoding.EncodeToString(srcBytes)
	if envVar == "" {
		return encodedStr, nil
	}
	return encodedStr, os.Setenv(envVar, encodedStr)
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	var override *Config
	// var secrets *Config
	// load config from default file
	err := config.DecodeTOML(bytes.NewReader(DefaultConfig), cfg)
	if err != nil {
		return nil, errors.Wrap(err, ErrReadConfig)
	}

	// load config overrides from env var if specified
	// there can be multiple overrides separated by comma
	rawConfigs, _ := osutil.GetEnv(OVERIDECONFIG)
	if rawConfigs != "" {
		for _, rawConfig := range strings.Split(rawConfigs, ",") {
			err = DecodeConfig(rawConfig, &override)
			if err != nil {
				return nil, fmt.Errorf("failed to decode override config: %w", err)
			}
			if override != nil {
				// apply overrides for all products
				if override.CCIP != nil {
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
		}
	}

	// read secrets for all products
	if cfg.CCIP != nil {
		err := ctfconfig.LoadSecretEnvsFromFiles()
		if err != nil {
			return nil, errors.Wrap(err, "error loading testsecrets files")
		}
		err = cfg.CCIP.LoadFromEnv()
		if err != nil {
			return nil, errors.Wrap(err, "error loading env vars into CCIP config")
		}
		// validate all products
		err = cfg.CCIP.Validate()
		if err != nil {
			return nil, err
		}
	}

	return cfg, nil
}

// Common is the generic config struct which can be used with product specific configs.
// It contains generic DON and networks config which can be applied to all product based tests.
type Common struct {
	EnvUser                 string                                      `toml:",omitempty"`
	EnvToConnect            *string                                     `toml:",omitempty"`
	TTL                     *config.Duration                            `toml:",omitempty"`
	ExistingCLCluster       *CLCluster                                  `toml:",omitempty"` // ExistingCLCluster is the existing chainlink cluster to use, if specified it will be used instead of creating a new one
	Mockserver              *string                                     `toml:",omitempty"`
	NewCLCluster            *ChainlinkDeployment                        `toml:",omitempty"` // NewCLCluster is the new chainlink cluster to create, if specified along with ExistingCLCluster this will be ignored
	Network                 *ctfconfig.NetworkConfig                    `toml:",omitempty"`
	PrivateEthereumNetworks map[string]*ctfconfig.EthereumNetworkConfig `toml:",omitempty"`
	Logging                 *ctfconfig.LoggingConfig                    `toml:",omitempty"`
}

// ReadFromEnvVar loads selected env vars into the config
func (p *Common) ReadFromEnvVar() error {
	logger := logging.GetTestLogger(nil)

	testLogCollect := ctfconfig.MustReadEnvVar_Boolean(ctfconfig.E2E_TEST_LOG_COLLECT_ENV)
	if testLogCollect != nil {
		if p.Logging == nil {
			p.Logging = &ctfconfig.LoggingConfig{}
		}
		logger.Debug().Msgf("Using %s env var to override Logging.TestLogCollect", ctfconfig.E2E_TEST_LOG_COLLECT_ENV)
		p.Logging.TestLogCollect = testLogCollect
	}

	loggingRunID := ctfconfig.MustReadEnvVar_String(ctfconfig.E2E_TEST_LOGGING_RUN_ID_ENV)
	if loggingRunID != "" {
		if p.Logging == nil {
			p.Logging = &ctfconfig.LoggingConfig{}
		}
		logger.Debug().Msgf("Using %s env var to override Logging.RunID", ctfconfig.E2E_TEST_LOGGING_RUN_ID_ENV)
		p.Logging.RunId = &loggingRunID
	}

	logstreamLogTargets := ctfconfig.MustReadEnvVar_Strings(ctfconfig.E2E_TEST_LOG_STREAM_LOG_TARGETS_ENV, ",")
	if len(logstreamLogTargets) > 0 {
		if p.Logging == nil {
			p.Logging = &ctfconfig.LoggingConfig{}
		}
		if p.Logging.LogStream == nil {
			p.Logging.LogStream = &ctfconfig.LogStreamConfig{}
		}
		logger.Debug().Msgf("Using %s env var to override Logging.LogStream.LogTargets", ctfconfig.E2E_TEST_LOG_STREAM_LOG_TARGETS_ENV)
		p.Logging.LogStream.LogTargets = logstreamLogTargets
	}

	lokiTenantID := ctfconfig.MustReadEnvVar_String(ctfconfig.E2E_TEST_LOKI_TENANT_ID_ENV)
	if lokiTenantID != "" {
		if p.Logging == nil {
			p.Logging = &ctfconfig.LoggingConfig{}
		}
		if p.Logging.Loki == nil {
			p.Logging.Loki = &ctfconfig.LokiConfig{}
		}
		logger.Debug().Msgf("Using %s env var to override Logging.Loki.TenantId", ctfconfig.E2E_TEST_LOKI_TENANT_ID_ENV)
		p.Logging.Loki.TenantId = &lokiTenantID
	}

	lokiEndpoint := ctfconfig.MustReadEnvVar_String(ctfconfig.E2E_TEST_LOKI_ENDPOINT_ENV)
	if lokiEndpoint != "" {
		if p.Logging == nil {
			p.Logging = &ctfconfig.LoggingConfig{}
		}
		if p.Logging.Loki == nil {
			p.Logging.Loki = &ctfconfig.LokiConfig{}
		}
		logger.Debug().Msgf("Using %s env var to override Logging.Loki.Endpoint", ctfconfig.E2E_TEST_LOKI_ENDPOINT_ENV)
		p.Logging.Loki.Endpoint = &lokiEndpoint
	}

	lokiBasicAuth := ctfconfig.MustReadEnvVar_String(ctfconfig.E2E_TEST_LOKI_BASIC_AUTH_ENV)
	if lokiBasicAuth != "" {
		if p.Logging == nil {
			p.Logging = &ctfconfig.LoggingConfig{}
		}
		if p.Logging.Loki == nil {
			p.Logging.Loki = &ctfconfig.LokiConfig{}
		}
		logger.Debug().Msgf("Using %s env var to override Logging.Loki.BasicAuth", ctfconfig.E2E_TEST_LOKI_BASIC_AUTH_ENV)
		p.Logging.Loki.BasicAuth = &lokiBasicAuth
	}

	lokiBearerToken := ctfconfig.MustReadEnvVar_String(ctfconfig.E2E_TEST_LOKI_BEARER_TOKEN_ENV)
	if lokiBearerToken != "" {
		if p.Logging == nil {
			p.Logging = &ctfconfig.LoggingConfig{}
		}
		if p.Logging.Loki == nil {
			p.Logging.Loki = &ctfconfig.LokiConfig{}
		}
		logger.Debug().Msgf("Using %s env var to override Logging.Loki.BearerToken", ctfconfig.E2E_TEST_LOKI_BEARER_TOKEN_ENV)
		p.Logging.Loki.BearerToken = &lokiBearerToken
	}

	grafanaBaseUrl := ctfconfig.MustReadEnvVar_String(ctfconfig.E2E_TEST_GRAFANA_BASE_URL_ENV)
	if grafanaBaseUrl != "" {
		if p.Logging == nil {
			p.Logging = &ctfconfig.LoggingConfig{}
		}
		if p.Logging.Grafana == nil {
			p.Logging.Grafana = &ctfconfig.GrafanaConfig{}
		}
		logger.Debug().Msgf("Using %s env var to override Logging.Grafana.BaseUrl", ctfconfig.E2E_TEST_GRAFANA_BASE_URL_ENV)
		p.Logging.Grafana.BaseUrl = &grafanaBaseUrl
	}

	grafanaDashboardUrl := ctfconfig.MustReadEnvVar_String(ctfconfig.E2E_TEST_GRAFANA_DASHBOARD_URL_ENV)
	if grafanaDashboardUrl != "" {
		if p.Logging == nil {
			p.Logging = &ctfconfig.LoggingConfig{}
		}
		if p.Logging.Grafana == nil {
			p.Logging.Grafana = &ctfconfig.GrafanaConfig{}
		}
		logger.Debug().Msgf("Using %s env var to override Logging.Grafana.DashboardUrl", ctfconfig.E2E_TEST_GRAFANA_DASHBOARD_URL_ENV)
		p.Logging.Grafana.DashboardUrl = &grafanaDashboardUrl
	}

	grafanaBearerToken := ctfconfig.MustReadEnvVar_String(ctfconfig.E2E_TEST_GRAFANA_BEARER_TOKEN_ENV)
	if grafanaBearerToken != "" {
		if p.Logging == nil {
			p.Logging = &ctfconfig.LoggingConfig{}
		}
		if p.Logging.Grafana == nil {
			p.Logging.Grafana = &ctfconfig.GrafanaConfig{}
		}
		logger.Debug().Msgf("Using %s env var to override Logging.Grafana.BearerToken", ctfconfig.E2E_TEST_GRAFANA_BEARER_TOKEN_ENV)
		p.Logging.Grafana.BearerToken = &grafanaBearerToken
	}

	selectedNetworks := ctfconfig.MustReadEnvVar_Strings(ctfconfig.E2E_TEST_SELECTED_NETWORK_ENV, ",")
	if len(selectedNetworks) > 0 {
		if p.Network == nil {
			p.Network = &ctfconfig.NetworkConfig{}
		}
		logger.Debug().Msgf("Using %s env var to override Network.SelectedNetworks", ctfconfig.E2E_TEST_SELECTED_NETWORK_ENV)
		p.Network.SelectedNetworks = selectedNetworks
	}

	walletKeys := ctfconfig.ReadEnvVarGroupedMap(ctfconfig.E2E_TEST_WALLET_KEY_ENV, ctfconfig.E2E_TEST_WALLET_KEYS_ENV)
	if len(walletKeys) > 0 {
		if p.Network == nil {
			p.Network = &ctfconfig.NetworkConfig{}
		}
		logger.Debug().Msgf("Using %s and/or %s env vars to override Network.WalletKeys", ctfconfig.E2E_TEST_WALLET_KEY_ENV, ctfconfig.E2E_TEST_WALLET_KEYS_ENV)
		p.Network.WalletKeys = walletKeys
	}

	rpcHttpUrls := ctfconfig.ReadEnvVarGroupedMap(ctfconfig.E2E_TEST_RPC_HTTP_URL_ENV, ctfconfig.E2E_TEST_RPC_HTTP_URLS_ENV)
	if len(rpcHttpUrls) > 0 {
		if p.Network == nil {
			p.Network = &ctfconfig.NetworkConfig{}
		}
		logger.Debug().Msgf("Using %s and/or %s env vars to override Network.RpcHttpUrls", ctfconfig.E2E_TEST_RPC_HTTP_URL_ENV, ctfconfig.E2E_TEST_RPC_HTTP_URLS_ENV)
		p.Network.RpcHttpUrls = rpcHttpUrls
	}

	rpcWsUrls := ctfconfig.ReadEnvVarGroupedMap(ctfconfig.E2E_TEST_RPC_WS_URL_ENV, ctfconfig.E2E_TEST_RPC_WS_URLS_ENV)
	if len(rpcWsUrls) > 0 {
		if p.Network == nil {
			p.Network = &ctfconfig.NetworkConfig{}
		}
		logger.Debug().Msgf("Using %s and/or %s env vars to override Network.RpcWsUrls", ctfconfig.E2E_TEST_RPC_WS_URL_ENV, ctfconfig.E2E_TEST_RPC_WS_URLS_ENV)
		p.Network.RpcWsUrls = rpcWsUrls
	}

	chainlinkImage := ctfconfig.MustReadEnvVar_String(ctfconfig.E2E_TEST_CHAINLINK_IMAGE_ENV)
	if chainlinkImage != "" {
		if p.NewCLCluster == nil {
			p.NewCLCluster = &ChainlinkDeployment{}
		}
		if p.NewCLCluster.Common == nil {
			p.NewCLCluster.Common = &Node{}
		}
		if p.NewCLCluster.Common.ChainlinkImage == nil {
			p.NewCLCluster.Common.ChainlinkImage = &ctfconfig.ChainlinkImageConfig{}
		}

		logger.Debug().Msgf("Using %s env var to override NewCLCluster.Common.ChainlinkImage.Image", ctfconfig.E2E_TEST_CHAINLINK_IMAGE_ENV)
		p.NewCLCluster.Common.ChainlinkImage.Image = &chainlinkImage
	}

	chainlinkVersion := ctfconfig.MustReadEnvVar_String(ctfconfig.E2E_TEST_CHAINLINK_VERSION_ENV)
	if chainlinkVersion != "" {
		if p.NewCLCluster == nil {
			p.NewCLCluster = &ChainlinkDeployment{}
		}
		if p.NewCLCluster.Common == nil {
			p.NewCLCluster.Common = &Node{}
		}
		if p.NewCLCluster.Common.ChainlinkImage == nil {
			p.NewCLCluster.Common.ChainlinkImage = &ctfconfig.ChainlinkImageConfig{}
		}

		logger.Debug().Msgf("Using %s env var to override NewCLCluster.Common.ChainlinkImage.Version", ctfconfig.E2E_TEST_CHAINLINK_VERSION_ENV)
		p.NewCLCluster.Common.ChainlinkImage.Version = &chainlinkVersion
	}

	chainlinkPostgresVersion := ctfconfig.MustReadEnvVar_String(ctfconfig.E2E_TEST_CHAINLINK_POSTGRES_VERSION_ENV)
	if chainlinkPostgresVersion != "" {
		if p.NewCLCluster == nil {
			p.NewCLCluster = &ChainlinkDeployment{}
		}
		if p.NewCLCluster.Common == nil {
			p.NewCLCluster.Common = &Node{}
		}
		if p.NewCLCluster.Common.ChainlinkImage == nil {
			p.NewCLCluster.Common.ChainlinkImage = &ctfconfig.ChainlinkImageConfig{}
		}

		logger.Debug().Msgf("Using %s env var to override NewCLCluster.Common.ChainlinkImage.PostgresVersion", ctfconfig.E2E_TEST_CHAINLINK_POSTGRES_VERSION_ENV)
		p.NewCLCluster.Common.ChainlinkImage.PostgresVersion = &chainlinkPostgresVersion
	}

	chainlinkUpgradeImage := ctfconfig.MustReadEnvVar_String(ctfconfig.E2E_TEST_CHAINLINK_UPGRADE_IMAGE_ENV)
	if chainlinkUpgradeImage != "" {
		if p.NewCLCluster == nil {
			p.NewCLCluster = &ChainlinkDeployment{}
		}
		if p.NewCLCluster.Common == nil {
			p.NewCLCluster.Common = &Node{}
		}
		if p.NewCLCluster.Common.ChainlinkUpgradeImage == nil {
			p.NewCLCluster.Common.ChainlinkUpgradeImage = &ctfconfig.ChainlinkImageConfig{}
		}

		logger.Debug().Msgf("Using %s env var to override NewCLCluster.Common.ChainlinkUpgradeImage.Image", ctfconfig.E2E_TEST_CHAINLINK_UPGRADE_IMAGE_ENV)
		p.NewCLCluster.Common.ChainlinkUpgradeImage.Image = &chainlinkUpgradeImage
	}

	chainlinkUpgradeVersion := ctfconfig.MustReadEnvVar_String(ctfconfig.E2E_TEST_CHAINLINK_UPGRADE_VERSION_ENV)
	if chainlinkUpgradeVersion != "" {
		if p.NewCLCluster == nil {
			p.NewCLCluster = &ChainlinkDeployment{}
		}
		if p.NewCLCluster.Common == nil {
			p.NewCLCluster.Common = &Node{}
		}
		if p.NewCLCluster.Common.ChainlinkImage == nil {
			p.NewCLCluster.Common.ChainlinkImage = &ctfconfig.ChainlinkImageConfig{}
		}

		logger.Debug().Msgf("Using %s env var to override NewCLCluster.Common.ChainlinkUpgradeImage.Version", ctfconfig.E2E_TEST_CHAINLINK_UPGRADE_VERSION_ENV)
		p.NewCLCluster.Common.ChainlinkUpgradeImage.Version = &chainlinkUpgradeVersion
	}

	return nil
}

func (p *Common) GetNodeConfig() *ctfconfig.NodeConfig {
	return &ctfconfig.NodeConfig{
		BaseConfigTOML:           p.NewCLCluster.Common.BaseConfigTOML,
		CommonChainConfigTOML:    p.NewCLCluster.Common.CommonChainConfigTOML,
		ChainConfigTOMLByChainID: p.NewCLCluster.Common.ChainConfigTOMLByChain,
	}
}

func (p *Common) GetSethConfig() *seth.Config {
	return nil
}

func (p *Common) Validate() error {
	if err := p.Logging.Validate(); err != nil {
		return fmt.Errorf("error validating logging config %w", err)
	}
	if p.Network == nil {
		return errors.New("no networks specified")
	}
	// read the default network config, if specified
	p.Network.UpperCaseNetworkNames()
	p.Network.OverrideURLsAndKeysFromEVMNetwork()
	if err := p.Network.Validate(); err != nil {
		return fmt.Errorf("error validating networks config %w", err)
	}
	if p.NewCLCluster == nil && p.ExistingCLCluster == nil {
		return errors.New("no chainlink or existing cluster specified")
	}

	for k, v := range p.PrivateEthereumNetworks {
		// this is the only value we need to generate dynamically before starting a new simulated chain
		if v.EthereumChainConfig != nil {
			p.PrivateEthereumNetworks[k].EthereumChainConfig.GenerateGenesisTimestamp()
		}

		builder := test_env.NewEthereumNetworkBuilder()
		ethNetwork, err := builder.WithExistingConfig(*v).Build()
		if err != nil {
			return fmt.Errorf("error building private ethereum network ethNetworks %w", err)
		}

		p.PrivateEthereumNetworks[k] = &ethNetwork.EthereumNetworkConfig
	}

	if p.ExistingCLCluster != nil {
		if err := p.ExistingCLCluster.Validate(); err != nil {
			return fmt.Errorf("error validating existing chainlink cluster config %w", err)
		}
		if p.Mockserver == nil {
			return errors.New("no mockserver specified for existing chainlink cluster")
		}
		log.Warn().Msg("Using existing chainlink cluster, overriding new chainlink cluster config if specified")
		p.NewCLCluster = nil
	} else {
		if p.NewCLCluster != nil {
			if err := p.NewCLCluster.Validate(); err != nil {
				return fmt.Errorf("error validating chainlink config %w", err)
			}
		}
	}
	return nil
}

func (p *Common) EVMNetworks() ([]blockchain.EVMNetwork, []string, error) {
	evmNetworks := networks.MustGetSelectedNetworkConfig(p.Network)
	if len(p.Network.SelectedNetworks) != len(evmNetworks) {
		return nil, p.Network.SelectedNetworks, fmt.Errorf("selected networks %v do not match evm networks %v", p.Network.SelectedNetworks, evmNetworks)
	}
	return evmNetworks, p.Network.SelectedNetworks, nil
}

func (p *Common) GetLoggingConfig() *ctfconfig.LoggingConfig {
	return p.Logging
}

func (p *Common) GetChainlinkImageConfig() *ctfconfig.ChainlinkImageConfig {
	return p.NewCLCluster.Common.ChainlinkImage
}

func (p *Common) GetPyroscopeConfig() *ctfconfig.PyroscopeConfig {
	return nil
}

func (p *Common) GetPrivateEthereumNetworkConfig() *ctfconfig.EthereumNetworkConfig {
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

type CLCluster struct {
	Name        *string                   `toml:",omitempty"`
	NoOfNodes   *int                      `toml:",omitempty"`
	NodeConfigs []*client.ChainlinkConfig `toml:",omitempty"`
}

func (c *CLCluster) Validate() error {
	if c.NoOfNodes == nil || len(c.NodeConfigs) == 0 {
		return fmt.Errorf("no chainlink nodes specified")
	}
	if *c.NoOfNodes != len(c.NodeConfigs) {
		return fmt.Errorf("number of nodes %d does not match number of node configs %d", *c.NoOfNodes, len(c.NodeConfigs))
	}
	for i, nodeConfig := range c.NodeConfigs {
		if nodeConfig.URL == "" {
			return fmt.Errorf("node %d url not specified", i+1)
		}
		if nodeConfig.Password == "" {
			return fmt.Errorf("node %d password not specified", i+1)
		}
		if nodeConfig.Email == "" {
			return fmt.Errorf("node %d email not specified", i+1)
		}
		if nodeConfig.InternalIP == "" {
			return fmt.Errorf("node %d internal ip not specified", i+1)
		}
	}

	return nil
}

type ChainlinkDeployment struct {
	Common         *Node    `toml:",omitempty"`
	NodeMemory     string   `toml:",omitempty"`
	NodeCPU        string   `toml:",omitempty"`
	DBMemory       string   `toml:",omitempty"`
	DBCPU          string   `toml:",omitempty"`
	DBCapacity     string   `toml:",omitempty"`
	DBStorageClass *string  `toml:",omitempty"`
	PromPgExporter *bool    `toml:",omitempty"`
	IsStateful     *bool    `toml:",omitempty"`
	DBArgs         []string `toml:",omitempty"`
	NoOfNodes      *int     `toml:",omitempty"`
	Nodes          []*Node  `toml:",omitempty"` // to be mentioned only if diff nodes follow diff configs; not required if all nodes follow CommonConfig
}

func (c *ChainlinkDeployment) Validate() error {
	if c.Common == nil {
		return errors.New("common config can't be empty")
	}
	if c.Common.ChainlinkImage == nil {
		return errors.New("chainlink image can't be empty")
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
		for i := range c.Nodes {
			// merge common config with node specific config
			c.Nodes[i].Merge(c.Common)
			node := c.Nodes[i]
			if node.ChainlinkImage == nil {
				return fmt.Errorf("node %s: chainlink image can't be empty", node.Name)
			}
			if err := node.ChainlinkImage.Validate(); err != nil {
				return fmt.Errorf("node %s: %w", node.Name, err)
			}
			if node.DBImage == "" || node.DBTag == "" {
				return fmt.Errorf("node %s: must provide db image and tag", node.Name)
			}
		}
	}
	return nil
}

type Node struct {
	Name                   string                          `toml:",omitempty"`
	NeedsUpgrade           *bool                           `toml:",omitempty"`
	ChainlinkImage         *ctfconfig.ChainlinkImageConfig `toml:"ChainlinkImage"`
	ChainlinkUpgradeImage  *ctfconfig.ChainlinkImageConfig `toml:"ChainlinkUpgradeImage"`
	BaseConfigTOML         string                          `toml:",omitempty"`
	CommonChainConfigTOML  string                          `toml:",omitempty"`
	ChainConfigTOMLByChain map[string]string               `toml:",omitempty"` // key is chainID
	DBImage                string                          `toml:",omitempty"`
	DBTag                  string                          `toml:",omitempty"`
}

// Merge merges non-empty values
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
	// merge upgrade image only if the nodes is marked as NeedsUpgrade to true
	if pointer.GetBool(n.NeedsUpgrade) {
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
