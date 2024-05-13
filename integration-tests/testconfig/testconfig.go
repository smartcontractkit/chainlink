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
	"github.com/google/uuid"
	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/smartcontractkit/seth"

	ctf_config "github.com/smartcontractkit/chainlink-testing-framework/config"
	k8s_config "github.com/smartcontractkit/chainlink-testing-framework/k8s/config"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/conversions"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/osutil"
	a_config "github.com/smartcontractkit/chainlink/integration-tests/testconfig/automation"
	f_config "github.com/smartcontractkit/chainlink/integration-tests/testconfig/functions"
	keeper_config "github.com/smartcontractkit/chainlink/integration-tests/testconfig/keeper"
	lp_config "github.com/smartcontractkit/chainlink/integration-tests/testconfig/log_poller"
	ocr_config "github.com/smartcontractkit/chainlink/integration-tests/testconfig/ocr"
	ocr2_config "github.com/smartcontractkit/chainlink/integration-tests/testconfig/ocr2"
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
	GetOCRConfig() *ocr_config.Config
}

type Ocr2TestConfig interface {
	GetOCR2Config() *ocr2_config.Config
}

type TestConfig struct {
	ctf_config.TestConfig

	Common     *Common                  `toml:"Common"`
	Automation *a_config.Config         `toml:"Automation"`
	Functions  *f_config.Config         `toml:"Functions"`
	Keeper     *keeper_config.Config    `toml:"Keeper"`
	LogPoller  *lp_config.Config        `toml:"LogPoller"`
	OCR        *ocr_config.Config       `toml:"OCR"`
	OCR2       *ocr2_config.Config      `toml:"OCR2"`
	VRF        *vrf_config.Config       `toml:"VRF"`
	VRFv2      *vrfv2_config.Config     `toml:"VRFv2"`
	VRFv2Plus  *vrfv2plus_config.Config `toml:"VRFv2Plus"`

	ConfigurationName string `toml:"-"`
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

func (c TestConfig) GetConfigurationName() string {
	return c.ConfigurationName
}

func (c TestConfig) GetSethConfig() *seth.Config {
	return c.Seth
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
	DirectRequest Product = "direct_request"
	Flux          Product = "flux"
	ForwarderOcr  Product = "forwarder_ocr"
	ForwarderOcr2 Product = "forwarder_ocr2"
	Functions     Product = "functions"
	Keeper        Product = "keeper"
	LogPoller     Product = "log_poller"
	Node          Product = "node"
	OCR           Product = "ocr"
	OCR2          Product = "ocr2"
	OCR2VRF       Product = "ocr2vrf"
	RunLog        Product = "runlog"
	VRF           Product = "vrf"
	VRFv2         Product = "vrfv2"
	VRFv2Plus     Product = "vrfv2plus"
)

var TestTypesWithLoki = []string{"Load", "Soak", "Stress", "Spike", "Volume"}

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
	NoKey                    = "NO_KEY"
)

func GetConfig(configurationName string, product Product) (TestConfig, error) {
	logger := logging.GetTestLogger(nil)

	configurationName = strings.ReplaceAll(configurationName, "/", "_")
	configurationName = strings.ReplaceAll(configurationName, " ", "_")
	configurationName = cases.Title(language.English, cases.NoLower).String(configurationName)
	fileNames := []string{
		"default.toml",
		fmt.Sprintf("%s.toml", product),
		"overrides.toml",
	}

	testConfig := TestConfig{}
	testConfig.ConfigurationName = configurationName
	logger.Debug().Msgf("Will apply configuration named '%s' if it is found in any of the configs", configurationName)

	var handleSpecialOverrides = func(logger zerolog.Logger, filename, configurationName string, target *TestConfig, content []byte, product Product) error {
		switch product {
		case Automation:
			return handleAutomationConfigOverride(logger, filename, configurationName, target, content)
		default:
			return handleDefaultConfigOverride(logger, filename, configurationName, target, content)
		}
	}

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

			err = handleSpecialOverrides(logger, fileName, configurationName, &testConfig, file, product)
			if err != nil {
				return TestConfig{}, errors.Wrapf(err, "error unmarshalling embedded config")
			}
		}
	}

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

		err = handleSpecialOverrides(logger, fileName, configurationName, &testConfig, content, product)
		if err != nil {
			return TestConfig{}, errors.Wrapf(err, "error reading file %s", filePath)
		}
	}

	logger.Info().Msg("Reading configs from Base64 override env var")
	configEncoded, isSet := os.LookupEnv(Base64OverrideEnvVarName)
	if isSet && configEncoded != "" {
		logger.Debug().Msgf("Found base64 config override environment variable '%s' found", Base64OverrideEnvVarName)
		decoded, err := base64.StdEncoding.DecodeString(configEncoded)
		if err != nil {
			return TestConfig{}, err
		}

		err = handleSpecialOverrides(logger, Base64OverrideEnvVarName, configurationName, &testConfig, decoded, product)
		if err != nil {
			return TestConfig{}, errors.Wrapf(err, "error unmarshaling base64 config")
		}
	} else {
		logger.Debug().Msg("Base64 config override from environment variable not found")
	}

	// it neede some custom logic, so we do it separately
	err := testConfig.readNetworkConfiguration()
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

	isAnySimulated := false
	for _, network := range testConfig.Network.SelectedNetworks {
		if strings.Contains(strings.ToUpper(network), "SIMULATED") {
			isAnySimulated = true
			break
		}
	}

	if testConfig.Seth != nil && !isAnySimulated && (testConfig.Seth.EphemeralAddrs != nil && *testConfig.Seth.EphemeralAddrs != 0) {
		testConfig.Seth.EphemeralAddrs = new(int64)
		logger.Warn().
			Msg("Ephemeral addresses were enabled, but test was setup to run on a live network. Ephemeral addresses will be disabled.")
	}

	if testConfig.Seth != nil && (testConfig.Seth.EphemeralAddrs != nil && *testConfig.Seth.EphemeralAddrs != 0) {
		rootBuffer := testConfig.Seth.RootKeyFundsBuffer
		if rootBuffer == nil {
			rootBuffer = big.NewInt(0)
		}
		clNodeFunding := testConfig.Common.ChainlinkNodeFunding
		if clNodeFunding == nil {
			zero := 0.0
			clNodeFunding = &zero
		}
		minRequiredFunds := big.NewFloat(0).Mul(big.NewFloat(*clNodeFunding), big.NewFloat(6.0))

		//add buffer to the minimum required funds, this isn't even a rough estimate, because we don't know how many contracts will be deployed from root key, but it's here to let you know that you should have some buffer
		minRequiredFundsBuffered := big.NewFloat(0).Mul(minRequiredFunds, big.NewFloat(1.2))
		minRequiredFundsBufferedInt, _ := minRequiredFundsBuffered.Int(nil)

		rootBuffer64, _ := rootBuffer.Float64()

		if big.NewFloat(rootBuffer64).Cmp(minRequiredFundsBuffered) <= 0 {
			msg := `
The funds allocated to the root key buffer are below the minimum requirement, which could lead to insufficient funds for performing contract deployments. Please review and adjust your TOML configuration file to ensure that the root key buffer has adequate funds. Increase the fund settings as necessary to meet this requirement.

Example:
[Seth]
root_key_funds_buffer = 1_000
`

			logger.Warn().
				Str("Root key buffer (wei/ether)", fmt.Sprintf("%s/%s", rootBuffer.String(), conversions.WeiToEther(rootBuffer).Text('f', -1))).
				Str("Minimum required funds (wei/ether)", fmt.Sprintf("%s/%s", minRequiredFundsBuffered.String(), conversions.WeiToEther(minRequiredFundsBufferedInt).Text('f', -1))).
				Msg(msg)
		}
	}

	logger.Debug().Msg("Correct test config constructed successfully")
	return testConfig, nil
}

func (c *TestConfig) readNetworkConfiguration() error {
	// currently we need to read that kind of secrets only for network configuration
	if c.Network == nil {
		c.Network = &ctf_config.NetworkConfig{}
	}

	c.Network.UpperCaseNetworkNames()
	c.Network.OverrideURLsAndKeysFromEVMNetwork()
	err := c.Network.Default()
	if err != nil {
		return errors.Wrapf(err, "error reading default network config")
	}

	// this is the only value we need to generate dynamically before starting a new simulated chain
	if c.PrivateEthereumNetwork != nil && c.PrivateEthereumNetwork.EthereumChainConfig != nil {
		c.PrivateEthereumNetwork.EthereumChainConfig.GenerateGenesisTimestamp()
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

func handleAutomationConfigOverride(logger zerolog.Logger, filename, configurationName string, target *TestConfig, content []byte) error {
	logger.Debug().Msgf("Handling automation config override for %s", filename)
	oldConfig := MustCopy(target)
	newConfig := TestConfig{}

	err := ctf_config.BytesToAnyTomlStruct(logger, filename, configurationName, &target, content)
	if err != nil {
		return errors.Wrapf(err, "error reading file %s", filename)
	}

	err = ctf_config.BytesToAnyTomlStruct(logger, filename, configurationName, &newConfig, content)
	if err != nil {
		return errors.Wrapf(err, "error reading file %s", filename)
	}

	// override instead of merging
	if (newConfig.Automation != nil && len(newConfig.Automation.Load) > 0) && (oldConfig != nil && oldConfig.Automation != nil && len(oldConfig.Automation.Load) > 0) {
		target.Automation.Load = newConfig.Automation.Load
	}

	return nil
}

func handleDefaultConfigOverride(logger zerolog.Logger, filename, configurationName string, target *TestConfig, content []byte) error {
	logger.Debug().Msgf("Handling default config override for %s", filename)
	oldConfig := MustCopy(target)
	newConfig := TestConfig{}

	err := ctf_config.BytesToAnyTomlStruct(logger, filename, configurationName, &target, content)
	if err != nil {
		return errors.Wrapf(err, "error reading file %s", filename)
	}

	err = ctf_config.BytesToAnyTomlStruct(logger, filename, configurationName, &newConfig, content)
	if err != nil {
		return errors.Wrapf(err, "error reading file %s", filename)
	}

	// temporary fix for Duration not being correctly copied
	if oldConfig != nil && oldConfig.Seth != nil && oldConfig.Seth.Networks != nil {
		for i, old_network := range oldConfig.Seth.Networks {
			for _, target_network := range target.Seth.Networks {
				if old_network.ChainID == target_network.ChainID {
					oldConfig.Seth.Networks[i].TxnTimeout = target_network.TxnTimeout
				}
			}
		}
	}

	// override instead of merging
	if (newConfig.Seth != nil && len(newConfig.Seth.Networks) > 0) && (oldConfig != nil && oldConfig.Seth != nil && len(oldConfig.Seth.Networks) > 0) {
		networksToUse := map[string]*seth.Network{}
		for i, old_network := range oldConfig.Seth.Networks {
			for _, new_network := range newConfig.Seth.Networks {
				if old_network.ChainID == new_network.ChainID {
					oldConfig.Seth.Networks[i] = new_network
					break
				}
				if _, ok := networksToUse[new_network.ChainID]; !ok {
					networksToUse[new_network.ChainID] = new_network
				}
			}
			networksToUse[old_network.ChainID] = oldConfig.Seth.Networks[i]
		}
		target.Seth.Networks = []*seth.Network{}
		for _, network := range networksToUse {
			target.Seth.Networks = append(target.Seth.Networks, network)
		}
	}

	return nil
}
