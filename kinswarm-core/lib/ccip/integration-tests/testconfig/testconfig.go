package testconfig

import (
	"embed"
	"encoding/base64"
	"fmt"
	"math/big"
	"os"
	"slices"
	"strings"

	"github.com/barkimedes/go-deepcopy"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	ctf_config "github.com/smartcontractkit/chainlink-testing-framework/lib/config"
	k8s_config "github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/config"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/networks"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/conversions"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/osutil"

	a_config "github.com/smartcontractkit/chainlink/integration-tests/testconfig/automation"
	ccip_config "github.com/smartcontractkit/chainlink/integration-tests/testconfig/ccip"
	f_config "github.com/smartcontractkit/chainlink/integration-tests/testconfig/functions"
	keeper_config "github.com/smartcontractkit/chainlink/integration-tests/testconfig/keeper"
	lp_config "github.com/smartcontractkit/chainlink/integration-tests/testconfig/log_poller"
	ocr_config "github.com/smartcontractkit/chainlink/integration-tests/testconfig/ocr"
	vrf_config "github.com/smartcontractkit/chainlink/integration-tests/testconfig/vrf"
	vrfv2_config "github.com/smartcontractkit/chainlink/integration-tests/testconfig/vrfv2"
	vrfv2plus_config "github.com/smartcontractkit/chainlink/integration-tests/testconfig/vrfv2plus"
)

type UpgradeableChainlinkTestConfig interface {
	GetChainlinkUpgradeImageConfig() *ctf_config.ChainlinkImageConfig
}

type CommonTestConfig interface {
	GetCommonConfig() *Common
}

type VRFv2TestConfig interface {
	GetVRFv2Config() *vrfv2_config.Config
}

type VRFv2PlusTestConfig interface {
	GetVRFv2PlusConfig() *vrfv2plus_config.Config
}

type FunctionsTestConfig interface {
	GetFunctionsConfig() *f_config.Config
}

type KeeperTestConfig interface {
	GetKeeperConfig() *keeper_config.Config
}

type AutomationTestConfig interface {
	GetAutomationConfig() *a_config.Config
}

type OcrTestConfig interface {
	GetActiveOCRConfig() *ocr_config.Config
}

type Ocr2TestConfig interface {
	GetOCR2Config() *ocr_config.Config
}

type CCIPTestConfig interface {
	GetCCIPConfig() *ccip_config.Config
}

type LinkTokenContractConfig interface {
	LinkTokenContractAddress() (common.Address, error)
	UseExistingLinkTokenContract() bool
}

const (
	E2E_TEST_DATA_STREAMS_URL_ENV      = "E2E_TEST_DATA_STREAMS_URL"
	E2E_TEST_DATA_STREAMS_USERNAME_ENV = "E2E_TEST_DATA_STREAMS_USERNAME"
	E2E_TEST_DATA_STREAMS_PASSWORD_ENV = "E2E_TEST_DATA_STREAMS_PASSWORD"
)

type TestConfig struct {
	ctf_config.TestConfig

	Common     *Common                  `toml:"Common"`
	Automation *a_config.Config         `toml:"Automation"`
	Functions  *f_config.Config         `toml:"Functions"`
	Keeper     *keeper_config.Config    `toml:"Keeper"`
	LogPoller  *lp_config.Config        `toml:"LogPoller"`
	OCR        *ocr_config.Config       `toml:"OCR"`
	OCR2       *ocr_config.Config       `toml:"OCR2"`
	VRF        *vrf_config.Config       `toml:"VRF"`
	VRFv2      *vrfv2_config.Config     `toml:"VRFv2"`
	VRFv2Plus  *vrfv2plus_config.Config `toml:"VRFv2Plus"`
	CCIP       *ccip_config.Config      `toml:"CCIP"`

	ConfigurationNames []string `toml:"-"`
}

var embeddedConfigs embed.FS
var areConfigsEmbedded bool

func init() {
	embeddedConfigs = embeddedConfigsFs
}

// Returns Grafana URL from Logging config
func (c *TestConfig) GetGrafanaBaseURL() (string, error) {
	if c.Logging.Grafana == nil || c.Logging.Grafana.BaseUrl == nil {
		return "", errors.New("grafana base url not set")
	}

	return strings.TrimSuffix(*c.Logging.Grafana.BaseUrl, "/"), nil
}

// Returns Grafana Dashboard URL from Logging config
func (c *TestConfig) GetGrafanaDashboardURL() (string, error) {
	if c.Logging.Grafana == nil || c.Logging.Grafana.DashboardUrl == nil {
		return "", errors.New("grafana dashboard url not set")
	}

	url := *c.Logging.Grafana.DashboardUrl
	if !strings.HasPrefix(url, "/") {
		url = "/" + url
	}

	return url, nil
}

// Saves Test Config to a local file
func (c *TestConfig) Save() (string, error) {
	filePath := fmt.Sprintf("test_config-%s.toml", uuid.New())

	content, err := toml.Marshal(*c)
	if err != nil {
		return "", errors.Wrapf(err, "error marshaling test config")
	}

	err = os.WriteFile(filePath, content, 0600)
	if err != nil {
		return "", errors.Wrapf(err, "error writing test config")
	}

	return filePath, nil
}

// MustCopy Returns a deep copy of the Test Config or panics on error
func (c TestConfig) MustCopy() any {
	return deepcopy.MustAnything(c).(TestConfig)
}

// MustCopy Returns a deep copy of struct passed to it and returns a typed copy (or panics on error)
func MustCopy[T any](c T) T {
	return deepcopy.MustAnything(c).(T)
}

func (c *TestConfig) GetLoggingConfig() *ctf_config.LoggingConfig {
	return c.Logging
}

func (c TestConfig) GetNetworkConfig() *ctf_config.NetworkConfig {
	return c.Network
}

func (c TestConfig) GetChainlinkImageConfig() *ctf_config.ChainlinkImageConfig {
	return c.ChainlinkImage
}

func (c TestConfig) GetPrivateEthereumNetworkConfig() *ctf_config.EthereumNetworkConfig {
	return c.PrivateEthereumNetwork
}

func (c TestConfig) GetPyroscopeConfig() *ctf_config.PyroscopeConfig {
	return c.Pyroscope
}

func (c TestConfig) GetCommonConfig() *Common {
	return c.Common
}

func (c TestConfig) GetVRFv2Config() *vrfv2_config.Config {
	return c.VRFv2
}

func (c TestConfig) GetFunctionsConfig() *f_config.Config {
	return c.Functions
}

func (c TestConfig) GetVRFv2PlusConfig() *vrfv2plus_config.Config {
	return c.VRFv2Plus
}

func (c TestConfig) GetChainlinkUpgradeImageConfig() *ctf_config.ChainlinkImageConfig {
	return c.ChainlinkUpgradeImage
}

func (c TestConfig) GetKeeperConfig() *keeper_config.Config {
	return c.Keeper
}

func (c TestConfig) GetAutomationConfig() *a_config.Config {
	return c.Automation
}

func (c TestConfig) GetOCRConfig() *ocr_config.Config {
	return c.OCR
}

func (c TestConfig) GetCCIPConfig() *ccip_config.Config {
	return c.CCIP
}

func (c TestConfig) GetConfigurationNames() []string {
	return c.ConfigurationNames
}

func (c TestConfig) GetSethConfig() *seth.Config {
	return c.Seth
}

func (c TestConfig) GetActiveOCRConfig() *ocr_config.Config {
	if c.OCR != nil {
		return c.OCR
	}

	return c.OCR2
}

func (c *TestConfig) AsBase64() (string, error) {
	content, err := toml.Marshal(*c)
	if err != nil {
		return "", errors.Wrapf(err, "error marshaling test config")
	}

	return base64.StdEncoding.EncodeToString(content), nil
}

type Common struct {
	ChainlinkNodeFunding *float64 `toml:"chainlink_node_funding"`
}

func (c *Common) Validate() error {
	if c.ChainlinkNodeFunding != nil && *c.ChainlinkNodeFunding < 0 {
		return fmt.Errorf("chainlink node funding must be positive")
	}

	return nil
}

type Product string

const (
	Automation    Product = "automation"
	Cron          Product = "cron"
	Flux          Product = "flux"
	ForwarderOcr  Product = "forwarder_ocr"
	ForwarderOcr2 Product = "forwarder_ocr2"
	Functions     Product = "functions"
	Keeper        Product = "keeper"
	LogPoller     Product = "log_poller"
	Node          Product = "node"
	OCR           Product = "ocr"
	OCR2          Product = "ocr2"
	RunLog        Product = "runlog"
	VRF           Product = "vrf"
	VRFv2         Product = "vrfv2"
	VRFv2Plus     Product = "vrfv2plus"

	CCIP Product = "ccip"
)

const TestTypeEnvVarName = "TEST_TYPE"

func GetConfigurationNameFromEnv() (string, error) {
	testType := os.Getenv(TestTypeEnvVarName)
	if testType == "" {
		return "", fmt.Errorf("%s env var not set", TestTypeEnvVarName)
	}

	return cases.Title(language.English, cases.NoLower).String(testType), nil
}

const (
	Base64OverrideEnvVarName = k8s_config.EnvBase64ConfigOverride
)

func GetConfig(configurationNames []string, product Product) (TestConfig, error) {
	logger := logging.GetTestLogger(nil)

	for idx, configurationName := range configurationNames {
		configurationNames[idx] = strings.ReplaceAll(configurationName, "/", "_")
		configurationNames[idx] = strings.ReplaceAll(configurationName, " ", "_")
		configurationNames[idx] = cases.Title(language.English, cases.NoLower).String(configurationName)
	}

	// add unnamed (default) configuration as the first one to be read
	configurationNamesCopy := make([]string, len(configurationNames))
	copy(configurationNamesCopy, configurationNames)

	configurationNames = append([]string{""}, configurationNamesCopy...)

	fileNames := []string{
		"default.toml",
		fmt.Sprintf("%s.toml", product),
		"overrides.toml",
	}

	testConfig := TestConfig{}
	testConfig.ConfigurationNames = configurationNames

	logger.Debug().Msgf("Will apply configurations named '%s' if they are found in any of the configs", strings.Join(configurationNames, ","))

	// read embedded configs is build tag "embed" is set
	// this makes our life much easier when using a binary
	if areConfigsEmbedded {
		logger.Info().Msg("Reading embedded configs")
		embeddedFiles := []string{"default.toml", fmt.Sprintf("%s/%s.toml", product, product)}
		for _, fileName := range embeddedFiles {
			file, err := embeddedConfigs.ReadFile(fileName)
			if err != nil && errors.Is(err, os.ErrNotExist) {
				logger.Debug().Msgf("Embedded config file %s not found. Continuing", fileName)
				continue
			} else if err != nil {
				return TestConfig{}, errors.Wrapf(err, "error reading embedded config")
			}

			for _, configurationName := range configurationNames {
				err = ctf_config.BytesToAnyTomlStruct(logger, fileName, configurationName, &testConfig, file)
				if err != nil {
					return TestConfig{}, errors.Wrapf(err, "error unmarshalling embedded config %s", embeddedFiles)
				}
			}
		}
	} else {
		logger.Info().Msg("Reading configs from file system")
		for _, fileName := range fileNames {
			logger.Debug().Msgf("Looking for config file %s", fileName)
			filePath, err := osutil.FindFile(fileName, osutil.DEFAULT_STOP_FILE_NAME, 3)

			if err != nil && errors.Is(err, os.ErrNotExist) {
				logger.Debug().Msgf("Config file %s not found", fileName)
				continue
			} else if err != nil {
				return TestConfig{}, errors.Wrapf(err, "error looking for file %s", filePath)
			}
			logger.Debug().Str("location", filePath).Msgf("Found config file %s", fileName)

			content, err := readFile(filePath)
			if err != nil {
				return TestConfig{}, errors.Wrapf(err, "error reading file %s", filePath)
			}

			_ = checkSecretsInToml(content)

			for _, configurationName := range configurationNames {
				err = ctf_config.BytesToAnyTomlStruct(logger, fileName, configurationName, &testConfig, content)
				if err != nil {
					return TestConfig{}, errors.Wrapf(err, "error reading file %s", filePath)
				}
			}
		}
	}

	logger.Info().Msg("Setting env vars from testsecrets dot-env files")
	err := ctf_config.LoadSecretEnvsFromFiles()
	if err != nil {
		return TestConfig{}, errors.Wrapf(err, "error reading test config values from ~/.testsecrets file")
	}

	logger.Info().Msg("Reading config values from existing env vars")
	err = testConfig.ReadFromEnvVar()
	if err != nil {
		return TestConfig{}, errors.Wrapf(err, "error reading test config values from env vars")
	}

	logger.Info().Msgf("Overriding config from %s env var", Base64OverrideEnvVarName)
	configEncoded, isSet := os.LookupEnv(Base64OverrideEnvVarName)
	if isSet && configEncoded != "" {
		logger.Debug().Msgf("Found base64 config override environment variable '%s' found", Base64OverrideEnvVarName)
		decoded, err := base64.StdEncoding.DecodeString(configEncoded)
		if err != nil {
			return TestConfig{}, err
		}

		_ = checkSecretsInToml(decoded)

		for _, configurationName := range configurationNames {
			err = ctf_config.BytesToAnyTomlStruct(logger, Base64OverrideEnvVarName, configurationName, &testConfig, decoded)
			if err != nil {
				return TestConfig{}, errors.Wrapf(err, "error unmarshaling base64 config")
			}
		}
	} else {
		logger.Debug().Msg("Base64 config override from environment variable not found")
	}

	err = testConfig.readNetworkConfiguration()
	if err != nil {
		return TestConfig{}, errors.Wrapf(err, "error reading network config")
	}

	logger.Debug().Msg("Validating test config")
	err = testConfig.Validate()
	if err != nil {
		logger.Error().
			Msg("Error validating test config. You might want refer to integration-tests/testconfig/README.md for more information.")
		return TestConfig{}, errors.Wrapf(err, "error validating test config")
	}

	if testConfig.Common == nil {
		testConfig.Common = &Common{}
	}

	testConfig.logRiskySettings(logger)

	logger.Debug().Msg("Correct test config constructed successfully")
	return testConfig, nil
}

// Read config values from environment variables
func (c *TestConfig) ReadFromEnvVar() error {
	logger := logging.GetTestLogger(nil)

	// Read values for base config from env vars
	err := c.TestConfig.ReadFromEnvVar()
	if err != nil {
		return errors.Wrapf(err, "error reading test config values from env vars")
	}

	dsURL := ctf_config.MustReadEnvVar_String(E2E_TEST_DATA_STREAMS_URL_ENV)
	if dsURL != "" {
		if c.Automation == nil {
			c.Automation = &a_config.Config{}
		}
		if c.Automation.DataStreams == nil {
			c.Automation.DataStreams = &a_config.DataStreams{}
		}
		logger.Debug().Msgf("Using %s env var to override Automation.DataStreams.URL", E2E_TEST_DATA_STREAMS_URL_ENV)
		c.Automation.DataStreams.URL = &dsURL
	}

	dsUsername := ctf_config.MustReadEnvVar_String(E2E_TEST_DATA_STREAMS_USERNAME_ENV)
	if dsUsername != "" {
		if c.Automation == nil {
			c.Automation = &a_config.Config{}
		}
		if c.Automation.DataStreams == nil {
			c.Automation.DataStreams = &a_config.DataStreams{}
		}
		logger.Debug().Msgf("Using %s env var to override Automation.DataStreams.Username", E2E_TEST_DATA_STREAMS_USERNAME_ENV)
		c.Automation.DataStreams.Username = &dsUsername
	}

	dsPassword := ctf_config.MustReadEnvVar_String(E2E_TEST_DATA_STREAMS_PASSWORD_ENV)
	if dsPassword != "" {
		if c.Automation == nil {
			c.Automation = &a_config.Config{}
		}
		if c.Automation.DataStreams == nil {
			c.Automation.DataStreams = &a_config.DataStreams{}
		}
		logger.Debug().Msgf("Using %s env var to override Automation.DataStreams.Password", E2E_TEST_DATA_STREAMS_PASSWORD_ENV)
		c.Automation.DataStreams.Password = &dsPassword
	}

	return nil
}

func (c *TestConfig) logRiskySettings(logger zerolog.Logger) {
	isAnySimulated := false
	for _, network := range c.Network.SelectedNetworks {
		if strings.Contains(strings.ToUpper(network), "SIMULATED") {
			isAnySimulated = true
			break
		}
	}

	if c.Seth != nil && !isAnySimulated && (c.Seth.EphemeralAddrs != nil && *c.Seth.EphemeralAddrs != 0) {
		c.Seth.EphemeralAddrs = new(int64)
		logger.Warn().
			Msg("Ephemeral addresses were enabled, but test was setup to run on a live network. Ephemeral addresses will be disabled.")
	}

	if c.Seth != nil && (c.Seth.EphemeralAddrs != nil && *c.Seth.EphemeralAddrs != 0) {
		rootBuffer := c.Seth.RootKeyFundsBuffer
		zero := int64(0)
		if rootBuffer == nil {
			rootBuffer = &zero
		}
		clNodeFunding := c.Common.ChainlinkNodeFunding
		if clNodeFunding == nil {
			zero := 0.0
			clNodeFunding = &zero
		}
		minRequiredFunds := big.NewFloat(0).Mul(big.NewFloat(*clNodeFunding), big.NewFloat(6.0))

		//add buffer to the minimum required funds, this isn't even a rough estimate, because we don't know how many contracts will be deployed from root key, but it's here to let you know that you should have some buffer
		minRequiredFundsBuffered := big.NewFloat(0).Mul(minRequiredFunds, big.NewFloat(1.2))
		minRequiredFundsBufferedInt, _ := minRequiredFundsBuffered.Int(nil)

		if *rootBuffer < minRequiredFundsBufferedInt.Int64() {
			msg := `
The funds allocated to the root key buffer are below the minimum requirement, which could lead to insufficient funds for performing contract deployments. Please review and adjust your TOML configuration file to ensure that the root key buffer has adequate funds. Increase the fund settings as necessary to meet this requirement.

Example:
[Seth]
root_key_funds_buffer = 1_000
`

			logger.Warn().
				Str("Root key buffer (wei/ether)", fmt.Sprintf("%s/%s", fmt.Sprint(rootBuffer), conversions.WeiToEther(big.NewInt(*rootBuffer)).Text('f', -1))).
				Str("Minimum required funds (wei/ether)", fmt.Sprintf("%s/%s", minRequiredFundsBuffered.String(), conversions.WeiToEther(minRequiredFundsBufferedInt).Text('f', -1))).
				Msg(msg)
		}
	}

	var customChainSettings []string
	for _, network := range networks.MustGetSelectedNetworkConfig(c.Network) {
		if c.NodeConfig != nil && len(c.NodeConfig.ChainConfigTOMLByChainID) > 0 {
			if _, ok := c.NodeConfig.ChainConfigTOMLByChainID[fmt.Sprint(network.ChainID)]; ok {
				logger.Warn().Msgf("You have provided custom Chainlink Node configuration for network '%s' (chain id: %d). Chainlink Node's default settings won't be used", network.Name, network.ChainID)
				customChainSettings = append(customChainSettings, fmt.Sprint(network.ChainID))
			}
		}
	}

	if len(customChainSettings) == 0 && c.NodeConfig != nil && c.NodeConfig.CommonChainConfigTOML != "" {
		logger.Warn().Msg("***** You have provided your own default Chainlink Node configuration for all networks. Chainlink Node's default settings for selected networks won't be used *****")
	}

}

// checkSecretsInToml checks if the TOML file contains secrets and shows error logs if it does
// This is a temporary and will be removed after migration to test secrets from env vars
func checkSecretsInToml(content []byte) error {
	logger := logging.GetTestLogger(nil)
	data := make(map[string]interface{})

	// Decode the TOML data
	err := toml.Unmarshal(content, &data)
	if err != nil {
		return errors.Wrapf(err, "error decoding TOML file")
	}

	logError := func(key, envVar string) {
		logger.Error().Msgf("Error in TOML test config!! TOML cannot have '%s' key. Remove it and set %s env in ~/.testsecrets instead", key, envVar)
	}

	if data["ChainlinkImage"] != nil {
		chainlinkImage := data["ChainlinkImage"].(map[string]interface{})
		if chainlinkImage["image"] != nil {
			logError("ChainlinkImage.image", "E2E_TEST_CHAINLINK_IMAGE")
		}
	}

	if data["ChainlinkUpgradeImage"] != nil {
		chainlinkUpgradeImage := data["ChainlinkUpgradeImage"].(map[string]interface{})
		if chainlinkUpgradeImage["image"] != nil {
			logError("ChainlinkUpgradeImage.image", "E2E_TEST_CHAINLINK_UPGRADE_IMAGE")
		}
	}

	if data["Network"] != nil {
		network := data["Network"].(map[string]interface{})
		if network["RpcHttpUrls"] != nil {
			logError("Network.RpcHttpUrls", "`E2E_TEST_(.+)_RPC_HTTP_URL$` like E2E_TEST_ARBITRUM_SEPOLIA_RPC_HTTP_URL")
		}
		if network["RpcWsUrls"] != nil {
			logError("Network.RpcWsUrls", "`E2E_TEST_(.+)_RPC_WS_URL$` like E2E_TEST_ARBITRUM_SEPOLIA_RPC_WS_URL")
		}
		if network["WalletKeys"] != nil {
			logError("Network.wallet_keys", "`E2E_TEST_(.+)_WALLET_KEY$` E2E_TEST_ARBITRUM_SEPOLIA_WALLET_KEY")
		}
	}

	return nil
}

func (c *TestConfig) readNetworkConfiguration() error {
	// currently we need to read that kind of secrets only for network configuration
	if c.Network == nil {
		c.Network = &ctf_config.NetworkConfig{}
	}

	c.Network.UpperCaseNetworkNames()
	c.Network.OverrideURLsAndKeysFromEVMNetwork()

	// this is the only value we need to generate dynamically before starting a new simulated chain
	if c.PrivateEthereumNetwork != nil && c.PrivateEthereumNetwork.EthereumChainConfig != nil {
		c.PrivateEthereumNetwork.EthereumChainConfig.GenerateGenesisTimestamp()
	}

	for _, network := range networks.MustGetSelectedNetworkConfig(c.Network) {
		for _, key := range network.PrivateKeys {
			address, err := conversions.PrivateKeyHexToAddress(key)
			if err != nil {
				return errors.Wrapf(err, "error converting private key to address")
			}
			c.PrivateEthereumNetwork.EthereumChainConfig.AddressesToFund = append(
				c.PrivateEthereumNetwork.EthereumChainConfig.AddressesToFund, address.Hex(),
			)
		}
	}
	return nil
}

func (c *TestConfig) Validate() error {
	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Errorf("panic during test config validation: '%v'. Most probably due to presence of partial product config", r))
		}
	}()

	if c.ChainlinkImage == nil {
		return MissingImageInfoAsError("chainlink image config must be set")
	}
	if c.ChainlinkImage != nil {
		if err := c.ChainlinkImage.Validate(); err != nil {
			return MissingImageInfoAsError(fmt.Sprintf("chainlink image config validation failed: %s", err.Error()))
		}
	}
	if c.ChainlinkUpgradeImage != nil {
		if err := c.ChainlinkUpgradeImage.Validate(); err != nil {
			return MissingImageInfoAsError(fmt.Sprintf("chainlink upgrade image config validation failed: %s", err.Error()))
		}
	}
	if err := c.Network.Validate(); err != nil {
		return NoSelectedNetworkInfoAsError(fmt.Sprintf("network config validation failed: %s", err.Error()))
	}

	if c.Logging == nil {
		return fmt.Errorf("logging config must be set")
	}

	if err := c.Logging.Validate(); err != nil {
		return errors.Wrapf(err, "logging config validation failed")
	}

	if c.Logging.Loki != nil {
		if err := c.Logging.Loki.Validate(); err != nil {
			return errors.Wrapf(err, "loki config validation failed")
		}
	}

	if c.Logging.LogStream != nil && slices.Contains(c.Logging.LogStream.LogTargets, "loki") {
		if c.Logging.Loki == nil {
			return fmt.Errorf("in order to use Loki as logging target you must set Loki config in logging config")
		}

		if err := c.Logging.Loki.Validate(); err != nil {
			return errors.Wrapf(err, "loki config validation failed")
		}
	}

	if c.Pyroscope != nil {
		if err := c.Pyroscope.Validate(); err != nil {
			return errors.Wrapf(err, "pyroscope config validation failed")
		}
	}

	if c.PrivateEthereumNetwork != nil {
		if err := c.PrivateEthereumNetwork.Validate(); err != nil {
			return errors.Wrapf(err, "private ethereum network config validation failed")
		}
	}

	if c.Common != nil {
		if err := c.Common.Validate(); err != nil {
			return errors.Wrapf(err, "Common config validation failed")
		}
	}

	if c.Automation != nil {
		if err := c.Automation.Validate(); err != nil {
			return errors.Wrapf(err, "Automation config validation failed")
		}
	}

	if c.Functions != nil {
		if err := c.Functions.Validate(); err != nil {
			return errors.Wrapf(err, "Functions config validation failed")
		}
	}

	if c.Keeper != nil {
		if err := c.Keeper.Validate(); err != nil {
			return errors.Wrapf(err, "Keeper config validation failed")
		}
	}

	if c.LogPoller != nil {
		if err := c.LogPoller.Validate(); err != nil {
			return errors.Wrapf(err, "LogPoller config validation failed")
		}
	}

	if c.OCR != nil {
		if err := c.OCR.Validate(); err != nil {
			return errors.Wrapf(err, "OCR config validation failed")
		}
	}

	if c.OCR2 != nil {
		if err := c.OCR2.Validate(); err != nil {
			return errors.Wrapf(err, "OCR2 config validation failed")
		}
	}

	if c.VRF != nil {
		if err := c.VRF.Validate(); err != nil {
			return errors.Wrapf(err, "VRF config validation failed")
		}
	}

	if c.VRFv2 != nil {
		if err := c.VRFv2.Validate(); err != nil {
			return errors.Wrapf(err, "VRFv2 config validation failed")
		}
	}

	if c.VRFv2Plus != nil {
		if err := c.VRFv2Plus.Validate(); err != nil {
			return errors.Wrapf(err, "VRFv2Plus config validation failed")
		}
	}

	if c.WaspConfig != nil {
		if err := c.WaspConfig.Validate(); err != nil {
			return errors.Wrapf(err, "WaspAutoBuildConfig validation failed")
		}
	}
	return nil
}

func readFile(filePath string) ([]byte, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "error reading file %s", filePath)
	}

	return content, nil
}
