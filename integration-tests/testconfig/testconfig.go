package testconfig

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"

	ctf_config "github.com/smartcontractkit/chainlink-testing-framework/config"
	ctf_test_env "github.com/smartcontractkit/chainlink-testing-framework/docker/test_env"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/osutil"
	a_config "github.com/smartcontractkit/chainlink/integration-tests/testconfig/automation"
	f_config "github.com/smartcontractkit/chainlink/integration-tests/testconfig/functions"
	keeper_config "github.com/smartcontractkit/chainlink/integration-tests/testconfig/keeper"
	lp_config "github.com/smartcontractkit/chainlink/integration-tests/testconfig/logpoller"
	ocr_config "github.com/smartcontractkit/chainlink/integration-tests/testconfig/ocr"
	vrf_config "github.com/smartcontractkit/chainlink/integration-tests/testconfig/vrf"
	vrfv2_config "github.com/smartcontractkit/chainlink/integration-tests/testconfig/vrfv2"
	vrfv2plus_config "github.com/smartcontractkit/chainlink/integration-tests/testconfig/vrfv2plus"
)

type TestConfig struct {
	ChainlinkImage         *ctf_config.ChainlinkImageConfig `toml:"ChainlinkImage"`
	ChainlinkUpgradeImage  *ctf_config.ChainlinkImageConfig `toml:"ChainlinkUpgradeImage"`
	Logging                *ctf_config.LoggingConfig        `toml:"Logging"`
	Network                *ctf_config.NetworkConfig        `toml:"Network"`
	Pyroscope              *ctf_config.PyroscopeConfig      `toml:"Pyroscope"`
	PrivateEthereumNetwork *ctf_test_env.EthereumNetwork    `toml:"PrivateEthereumNetwork"`

	Common     *Common                  `toml:"Common"`
	Automation *a_config.Config         `toml:"Automation"`
	Functions  *f_config.Config         `toml:"Functions"`
	Keeper     *keeper_config.Config    `toml:"Keeper"`
	LogPoller  *lp_config.Config        `toml:"LogPoller"`
	OCR        *ocr_config.Config       `toml:"OCR"`
	VRF        *vrf_config.Config       `toml:"VRF"`
	VRFv2      *vrfv2_config.Config     `toml:"VRFv2"`
	VRFv2Plus  *vrfv2plus_config.Config `toml:"VRFv2Plus"`

	TestType TestType
}

func (c *TestConfig) GetGrafanaURL() (string, error) {
	if c.Logging.Grafana == nil || c.Logging.Grafana.GrafanaUrl == nil {
		return "", errors.New("grafana url not set")
	}

	return *c.Logging.Grafana.GrafanaUrl, nil
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

func (c *Common) ApplyOverrides(from *Common) error {
	if from == nil {
		return nil
	}

	if from.ChainlinkNodeFunding != nil {
		c.ChainlinkNodeFunding = from.ChainlinkNodeFunding
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
	OCR           Product = "ocr"
	OCR2          Product = "ocr2"
	OCR2VRF       Product = "ocr2vrf"
	RunLog        Product = "run_log"
	VRF           Product = "vrf"
	VRFv2         Product = "vrfv2"
	VRFv2Plus     Product = "vrfv2plus"
)

type TestType string

const (
	Smoke       TestType = "smoke"
	Load        TestType = "load"
	Soak        TestType = "soak"
	Performance TestType = "performance"
	Benchmark   TestType = "benchmark"
	Reorg       TestType = "reorg"
	Chaos       TestType = "chaos"
	Stress      TestType = "stress"
	Spike       TestType = "spike"
)

const TestTypeEnvVarName = "TEST_TYPE"

func GetTestTypeFromEnv() (TestType, error) {
	testType := os.Getenv(TestTypeEnvVarName)
	if testType == "" {
		return "", fmt.Errorf("%s env var not set", TestTypeEnvVarName)
	}

	switch TestType(testType) {
	case Smoke, Load, Soak, Performance, Benchmark, Reorg, Chaos, Stress, Spike:
		return TestType(testType), nil
	default:
		return "", fmt.Errorf("invalid value for %s: %s", TestTypeEnvVarName, testType)
	}
}

const (
	Base64OverrideEnvVarName = "BASE64_CONFIG_OVERRIDE"
	NoTest                   = "NO_TEST"
)

// TODO should I  also add testname, so that we allow individual tests to override it?
func GetConfig(testName string, testType TestType, product Product) (TestConfig, error) {
	testName = strings.ReplaceAll(testName, "/", "_")
	testName = strings.ReplaceAll(testName, " ", "_")
	fileNames := []string{
		"default.toml",
		fmt.Sprintf("%s.toml", product),
		fmt.Sprintf("%s_%s.toml", testType, product),
		"overrides.toml",
	}

	if testName != NoTest {
		fileNames = append(fileNames, fmt.Sprintf("%s_%s_%s.toml", testName, testType, product))
	}

	testConfig := TestConfig{}
	maybeTestConfigs := []TestConfig{}

	for _, fileName := range fileNames {
		content, err := osutil.FindFile(fileName, osutil.DEFAULT_STOP_FILE_NAME)
		if err != nil && errors.Is(err, os.ErrNotExist) {
			continue
		}

		var readConfig TestConfig
		err = toml.Unmarshal(content, &readConfig)
		if err != nil {
			return TestConfig{}, errors.Wrapf(err, "error unmarshaling config")
		}

		maybeTestConfigs = append(maybeTestConfigs, readConfig)
	}

	configEncoded, isSet := os.LookupEnv(Base64OverrideEnvVarName)
	if isSet && configEncoded != "" {
		decoded, err := base64.StdEncoding.DecodeString(configEncoded)
		if err != nil {
			return TestConfig{}, err
		}

		var base64override TestConfig
		err = toml.Unmarshal(decoded, &base64override)
		if err != nil {
			return TestConfig{}, errors.Wrapf(err, "error unmarshaling base64 config")
		}

		maybeTestConfigs = append(maybeTestConfigs, base64override)
	}

	// currently we need to read that kind of secrets only for network configuration
	testConfig.Network = &ctf_config.NetworkConfig{}
	err := testConfig.Network.ApplySecrets()
	if err != nil {
		return TestConfig{}, errors.Wrapf(err, "error applying secrets to network config")
	}

	for i := range maybeTestConfigs {
		err := testConfig.ApplyOverrides(&maybeTestConfigs[i])
		if err != nil {
			return TestConfig{}, errors.Wrapf(err, "error applying overrides to test config")
		}
	}

	err = testConfig.Validate()
	if err != nil {
		return TestConfig{}, errors.Wrapf(err, "error validating test config")
	}

	return testConfig, nil
}

func (c *TestConfig) ApplyOverrides(from *TestConfig) error {
	if from == nil {
		return nil
	}

	if from.ChainlinkImage != nil {
		if c.ChainlinkImage == nil {
			c.ChainlinkImage = from.ChainlinkImage
		} else {
			err := c.ChainlinkImage.ApplyOverrides(from.ChainlinkImage)
			if err != nil {
				return errors.Wrapf(err, "error applying overrides to chainlink image config")
			}
		}
	}

	if from.ChainlinkUpgradeImage != nil {
		if c.ChainlinkUpgradeImage == nil {
			c.ChainlinkUpgradeImage = from.ChainlinkUpgradeImage
		} else {
			err := c.ChainlinkUpgradeImage.ApplyOverrides(from.ChainlinkUpgradeImage)
			if err != nil {
				return errors.Wrapf(err, "error applying overrides to chainlink image config")
			}
		}
	}

	if from.Logging != nil {
		if c.Logging == nil {
			c.Logging = from.Logging
		} else {
			err := c.Logging.ApplyOverrides(from.Logging)
			if err != nil {
				return errors.Wrapf(err, "error applying overrides to logging config")
			}
		}
	}

	if from.Network != nil {
		if c.Network == nil {
			c.Network = from.Network
		} else {
			err := c.Network.ApplyOverrides(from.Network)
			if err != nil {
				return errors.Wrapf(err, "error applying overrides to network config")
			}
		}
	}

	if from.Pyroscope != nil {
		if c.Pyroscope == nil {
			c.Pyroscope = from.Pyroscope
		} else {
			err := c.Pyroscope.ApplyOverrides(from.Pyroscope)
			if err != nil {
				return errors.Wrapf(err, "error applying overrides to pyroscope config")
			}
		}
	}

	if from.Common != nil {
		if c.Common == nil {
			c.Common = from.Common
		} else {
			err := c.Common.ApplyOverrides(from.Common)
			if err != nil {
				return errors.Wrapf(err, "error applying overrides to common config")
			}
		}
	}

	if from.OCR != nil {
		if c.OCR == nil {
			c.OCR = from.OCR
		} else {
			err := c.OCR.ApplyOverrides(from.OCR)
			if err != nil {
				return errors.Wrapf(err, "error applying overrides to OCR config")
			}
		}
	}

	if from.VRF != nil {
		if c.VRF == nil {
			c.VRF = from.VRF
		} else {
			err := c.VRF.ApplyOverrides(from.VRF)
			if err != nil {
				return errors.Wrapf(err, "error applying overrides to VRF config")
			}
		}
	}

	if from.VRFv2 != nil {
		if c.VRFv2 == nil {
			c.VRFv2 = from.VRFv2
		} else {
			err := c.VRFv2.ApplyOverrides(from.VRFv2)
			if err != nil {
				return errors.Wrapf(err, "error applying overrides to VRFv2 config")
			}
		}
	}

	if from.VRFv2Plus != nil {
		if c.VRFv2Plus == nil {
			c.VRFv2Plus = from.VRFv2Plus
		} else {
			err := c.VRFv2Plus.ApplyOverrides(from.VRFv2Plus)
			if err != nil {
				return errors.Wrapf(err, "error applying overrides to VRFv2Plus config")
			}
		}
	}

	if from.Keeper != nil {
		if c.Keeper == nil {
			c.Keeper = from.Keeper
		} else {
			err := c.Keeper.ApplyOverrides(from.Keeper)
			if err != nil {
				return errors.Wrapf(err, "error applying overrides to Keeper config")
			}
		}
	}

	if from.Automation != nil {
		if c.Automation == nil {
			c.Automation = from.Automation
		} else {
			err := c.Automation.ApplyOverrides(from.Automation)
			if err != nil {
				return errors.Wrapf(err, "error applying overrides to Automation config")
			}
		}
	}

	if from.LogPoller != nil {
		if c.LogPoller == nil {
			c.LogPoller = from.LogPoller
		} else {
			err := c.LogPoller.ApplyOverrides(from.LogPoller)
			if err != nil {
				return errors.Wrapf(err, "error applying overrides to LogPoller config")
			}
		}
	}

	if from.Functions != nil {
		if c.Functions == nil {
			c.Functions = from.Functions
		} else {
			err := c.Functions.ApplyOverrides(from.Functions)
			if err != nil {
				return errors.Wrapf(err, "error applying overrides to Functions config")
			}
		}
	}

	return nil
}

func (c *TestConfig) Validate() error {
	if c.ChainlinkImage == nil {
		return fmt.Errorf("chainlink image config must be set")
	}
	if err := c.ChainlinkImage.Validate(); err != nil {
		return errors.Wrapf(err, "chainlink image config validation failed")
	}
	if c.ChainlinkUpgradeImage != nil {
		if err := c.ChainlinkUpgradeImage.Validate(); err != nil {
			return errors.Wrapf(err, "chainlink upgrade image config validation failed")
		}
	}
	if err := c.Network.Validate(); err != nil {
		return errors.Wrapf(err, "network config validation failed")
	}

	if c.Logging == nil {
		return fmt.Errorf("logging config must be set")
	}

	if err := c.Logging.Validate(); err != nil {
		return errors.Wrapf(err, "logging config validation failed")
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

	return nil
}
