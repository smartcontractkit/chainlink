package testconfig

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
	pkg_errors "github.com/pkg/errors"

	ctf_config "github.com/smartcontractkit/chainlink-testing-framework/config"
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
	*ctf_config.ChainlinkImageConfig
	*ctf_config.LoggingConfig
	*ctf_config.NetworkConfig
	*ctf_config.PyroscopeConfig

	Common     *Common                  `toml:"Common"`
	OCR        *ocr_config.Config       `toml:"OCR"`
	VRF        *vrf_config.Config       `toml:"VRF"`
	VRFv2      *vrfv2_config.Config     `toml:"VRFv2"`
	VRFv2Plus  *vrfv2plus_config.Config `toml:"VRFv2Plus"`
	Keeper     *keeper_config.Config    `toml:"Keeper"`
	Automation *a_config.Config         `toml:"Automation"`
	LogPoller  *lp_config.Config        `toml:"LogPoller"`
	Functions  *f_config.Config         `toml:"Functions"`
	TestType   TestType
}

func (c *TestConfig) GetGrafanaURL() (string, error) {
	if c.LoggingConfig.Logging.Grafana == nil || c.LoggingConfig.Logging.Grafana.GrafanaUrl == nil {
		return "", errors.New("grafana url not set")
	}

	return *c.LoggingConfig.Logging.Grafana.GrafanaUrl, nil
}

type Common struct {
	ChainlinkNodeFunding *float64 `toml:"chainlink_node_funding"`
}

func (a *Common) ApplyOverrides(from interface{}) error {
	switch asCfg := (from).(type) {
	case *Common:
		if asCfg == nil {
			return nil
		}

		if asCfg.ChainlinkNodeFunding != nil {
			a.ChainlinkNodeFunding = asCfg.ChainlinkNodeFunding
		}

		return nil
	default:
		return pkg_errors.Errorf("cannot apply overrides from unsupported type %T to common config", from)
	}
}

type Product string

const (
	OCR           Product = "ocr"
	OCR2          Product = "ocr2"
	VRF           Product = "vrf"
	VRFv2         Product = "vrfv2"
	VRFv2Plus     Product = "vrfv2plus"
	OCR2VRF       Product = "ocr2vrf"
	LogPoller     Product = "log_poller"
	Automation    Product = "automation"
	Functions     Product = "functions"
	DirectRequest Product = "direct_request"
	Cron          Product = "cron"
	Flux          Product = "flux"
	Keeper        Product = "keeper"
	ForwarderOcr  Product = "forwarder_ocr"
	ForwarderOcr2 Product = "forwarder_ocr2"
	RunLog        Product = "run_log"
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
)

// TODO add also testname, so that we allow individual tests to override it?
// TODO at the very end add base64 override for the final config
func GetConfig(testType TestType, product Product) (TestConfig, error) {
	fileNames := []string{
		"default.toml",
		fmt.Sprintf("%s.toml", product),
		fmt.Sprintf("%s_%s.toml", testType, product),
		"overrides.toml",
	}

	// imageCfg := ctf_config.ChainlinkImageConfig{}
	// loggingCfg := ctf_config.LoggingConfig{}
	// networkCfg := ctf_config.NetworkConfig{}
	// ocrCfg := ocr_config.Config{}
	// commonCfg := Common{}
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
			return TestConfig{}, pkg_errors.Wrapf(err, "error unmarshaling config")
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
			return TestConfig{}, pkg_errors.Wrapf(err, "error unmarshaling base64 config")
		}

		maybeTestConfigs = append(maybeTestConfigs, base64override)
	}

	for _, maybeConfig := range maybeTestConfigs {
		err := testConfig.ApplyOverrides(maybeConfig)
		if err != nil {
			return TestConfig{}, pkg_errors.Wrapf(err, "error applying overrides to test config")
		}

		// err := commonCfg.ApplyOverrides(config.Common)
		// if err != nil {
		// 	return TestConfig{}, pkg_errors.Wrapf(err, fmt.Sprintf("error applying overrides from %s config file to common config", fileName))
		// }

		// err = imageCfg.ApplyOverrides(config.ChainlinkImageConfig)
		// if err != nil {
		// 	return TestConfig{}, pkg_errors.Wrapf(err, fmt.Sprintf("error applying overrides from %s config file to image config", fileName))
		// }

		// err = loggingCfg.ApplyOverrides(config.LoggingConfig)
		// if err != nil {
		// 	return TestConfig{}, pkg_errors.Wrapf(err, fmt.Sprintf("error applying overrides from %s config file to logging config", fileName))
		// }

		// err = networkCfg.ApplyOverrides(config.NetworkConfig)
		// if err != nil {
		// 	return TestConfig{}, pkg_errors.Wrapf(err, fmt.Sprintf("error applying overrides from %s config file to network config", fileName))
		// }

		// err = ocrCfg.ApplyOverrides(config.OCR)
		// if err != nil {
		// 	return TestConfig{}, pkg_errors.Wrapf(err, fmt.Sprintf("error applying overrides from %s config file to ocr config", fileName))
		// }
	}

	err := testConfig.Validate()
	if err != nil {
		return TestConfig{}, pkg_errors.Wrapf(err, "error validating test config")
	}

	return testConfig, nil

	// err := imageCfg.Validate()
	// if err != nil {
	// 	return TestConfig{}, pkg_errors.Wrapf(err, "error validating image config")
	// }

	// err = loggingCfg.Validate()
	// if err != nil {
	// 	return TestConfig{}, pkg_errors.Wrapf(err, "error validating logging config")
	// }

	// err = networkCfg.Validate()
	// if err != nil {
	// 	return TestConfig{}, pkg_errors.Wrapf(err, "error validating network config")
	// }

	// err = ocrCfg.Validate()
	// if err != nil {
	// 	return TestConfig{}, pkg_errors.Wrapf(err, "error validating ocr config")
	// }

	// err = vrfCfg.Validate()
	// if err != nil {
	// 	return TestConfig{}, pkg_errors.Wrapf(err, "error validating vrf config")
	// }

	// return TestConfig{
	// 	Common:               &commonCfg,
	// 	ChainlinkImageConfig: &imageCfg,
	// 	LoggingConfig:        &loggingCfg,
	// 	NetworkConfig:        &networkCfg,
	// 	OCR:                  &ocrCfg,
	// 	// VRF:                  &vrfCfg,
	// 	TestType: testType,
	// }, nil
}

func (c *TestConfig) ApplyOverrides(from interface{}) error {
	switch asCfg := (from).(type) {
	case *TestConfig:
		if asCfg == nil {
			return nil
		}

		if asCfg.ChainlinkImageConfig != nil {
			if c.Common == nil {
				c.Common = asCfg.Common
			} else {
				err := c.Common.ApplyOverrides(asCfg.Common)
				if err != nil {
					return pkg_errors.Wrapf(err, "error applying overrides to common config")
				}
			}
		}

		if asCfg.LoggingConfig != nil {
			if c.LoggingConfig == nil {
				c.LoggingConfig = asCfg.LoggingConfig
			} else {
				err := c.LoggingConfig.ApplyOverrides(asCfg.LoggingConfig)
				if err != nil {
					return pkg_errors.Wrapf(err, "error applying overrides to logging config")
				}
			}
		}

		if asCfg.NetworkConfig != nil {
			if c.NetworkConfig == nil {
				c.NetworkConfig = asCfg.NetworkConfig
			} else {
				err := c.NetworkConfig.ApplyOverrides(asCfg.NetworkConfig)
				if err != nil {
					return pkg_errors.Wrapf(err, "error applying overrides to network config")
				}
			}
		}

		if asCfg.PyroscopeConfig != nil {
			if c.PyroscopeConfig == nil {
				c.PyroscopeConfig = asCfg.PyroscopeConfig
			} else {
				err := c.PyroscopeConfig.ApplyOverrides(asCfg.PyroscopeConfig)
				if err != nil {
					return pkg_errors.Wrapf(err, "error applying overrides to pyroscope config")
				}
			}
		}

		if asCfg.Common != nil {
			if c.Common == nil {
				c.Common = asCfg.Common
			} else {
				err := c.Common.ApplyOverrides(asCfg.Common)
				if err != nil {
					return pkg_errors.Wrapf(err, "error applying overrides to common config")
				}
			}
		}

		if asCfg.OCR != nil {
			if c.OCR == nil {
				c.OCR = asCfg.OCR
			} else {
				err := c.OCR.ApplyOverrides(asCfg.OCR)
				if err != nil {
					return pkg_errors.Wrapf(err, "error applying overrides to OCR config")
				}
			}
		}

		if asCfg.VRF != nil {
			if c.VRF == nil {
				c.VRF = asCfg.VRF
			} else {
				err := c.VRF.ApplyOverrides(asCfg.VRF)
				if err != nil {
					return pkg_errors.Wrapf(err, "error applying overrides to VRF config")
				}
			}
		}

		if asCfg.VRFv2 != nil {
			if c.VRFv2 == nil {
				c.VRFv2 = asCfg.VRFv2
			} else {
				err := c.VRFv2.ApplyOverrides(asCfg.VRFv2)
				if err != nil {
					return pkg_errors.Wrapf(err, "error applying overrides to VRFv2 config")
				}
			}
		}

		if asCfg.VRFv2Plus != nil {
			if c.VRFv2Plus == nil {
				c.VRFv2Plus = asCfg.VRFv2Plus
			} else {
				err := c.VRFv2Plus.ApplyOverrides(asCfg.VRFv2Plus)
				if err != nil {
					return pkg_errors.Wrapf(err, "error applying overrides to VRFv2Plus config")
				}
			}
		}

		if asCfg.Keeper != nil {
			if c.Keeper == nil {
				c.Keeper = asCfg.Keeper
			} else {
				err := c.Keeper.ApplyOverrides(asCfg.Keeper)
				if err != nil {
					return pkg_errors.Wrapf(err, "error applying overrides to Keeper config")
				}
			}
		}

		if asCfg.Automation != nil {
			if c.Automation == nil {
				c.Automation = asCfg.Automation
			} else {
				err := c.Automation.ApplyOverrides(asCfg.Automation)
				if err != nil {
					return pkg_errors.Wrapf(err, "error applying overrides to Automation config")
				}
			}
		}

		if asCfg.LogPoller != nil {
			if c.LogPoller == nil {
				c.LogPoller = asCfg.LogPoller
			} else {
				err := c.LogPoller.ApplyOverrides(asCfg.LogPoller)
				if err != nil {
					return pkg_errors.Wrapf(err, "error applying overrides to LogPoller config")
				}
			}
		}

		if asCfg.Functions != nil {
			if c.Functions == nil {
				c.Functions = asCfg.Functions
			} else {
				err := c.Functions.ApplyOverrides(asCfg.Functions)
				if err != nil {
					return pkg_errors.Wrapf(err, "error applying overrides to Functions config")
				}
			}
		}

		return nil
	case TestConfig:
		if asCfg.ChainlinkImageConfig != nil {
			if c.ChainlinkImageConfig == nil {
				c.ChainlinkImageConfig = asCfg.ChainlinkImageConfig
			} else {
				err := c.ChainlinkImageConfig.ApplyOverrides(asCfg.ChainlinkImageConfig)
				if err != nil {
					return pkg_errors.Wrapf(err, "error applying overrides to ChainlinkImageConfig config")
				}
			}
		}

		if asCfg.LoggingConfig != nil {
			if c.LoggingConfig == nil {
				c.LoggingConfig = asCfg.LoggingConfig
			} else {
				err := c.LoggingConfig.ApplyOverrides(asCfg.LoggingConfig)
				if err != nil {
					return pkg_errors.Wrapf(err, "error applying overrides to logging config")
				}
			}
		}

		if asCfg.NetworkConfig != nil {
			if c.NetworkConfig == nil {
				c.NetworkConfig = asCfg.NetworkConfig
			} else {
				err := c.NetworkConfig.ApplyOverrides(asCfg.NetworkConfig)
				if err != nil {
					return pkg_errors.Wrapf(err, "error applying overrides to network config")
				}
			}
		}

		if asCfg.PyroscopeConfig != nil {
			if c.PyroscopeConfig == nil {
				c.PyroscopeConfig = asCfg.PyroscopeConfig
			} else {
				err := c.PyroscopeConfig.ApplyOverrides(asCfg.PyroscopeConfig)
				if err != nil {
					return pkg_errors.Wrapf(err, "error applying overrides to pyroscope config")
				}
			}
		}

		if asCfg.Common != nil {
			if c.Common == nil {
				c.Common = asCfg.Common
			} else {
				err := c.Common.ApplyOverrides(asCfg.Common)
				if err != nil {
					return pkg_errors.Wrapf(err, "error applying overrides to common config")
				}
			}
		}

		if asCfg.OCR != nil {
			if c.OCR == nil {
				c.OCR = asCfg.OCR
			} else {
				err := c.OCR.ApplyOverrides(asCfg.OCR)
				if err != nil {
					return pkg_errors.Wrapf(err, "error applying overrides to OCR config")
				}
			}
		}

		if asCfg.VRF != nil {
			if c.VRF == nil {
				c.VRF = asCfg.VRF
			} else {
				err := c.VRF.ApplyOverrides(asCfg.VRF)
				if err != nil {
					return pkg_errors.Wrapf(err, "error applying overrides to VRF config")
				}
			}
		}

		if asCfg.VRFv2 != nil {
			if c.VRFv2 == nil {
				c.VRFv2 = asCfg.VRFv2
			} else {
				err := c.VRFv2.ApplyOverrides(asCfg.VRFv2)
				if err != nil {
					return pkg_errors.Wrapf(err, "error applying overrides to VRFv2 config")
				}
			}
		}

		if asCfg.VRFv2Plus != nil {
			if c.VRFv2Plus == nil {
				c.VRFv2Plus = asCfg.VRFv2Plus
			} else {
				err := c.VRFv2Plus.ApplyOverrides(asCfg.VRFv2Plus)
				if err != nil {
					return pkg_errors.Wrapf(err, "error applying overrides to VRFv2Plus config")
				}
			}
		}

		if asCfg.Keeper != nil {
			if c.Keeper == nil {
				c.Keeper = asCfg.Keeper
			} else {
				err := c.Keeper.ApplyOverrides(asCfg.Keeper)
				if err != nil {
					return pkg_errors.Wrapf(err, "error applying overrides to Keeper config")
				}
			}
		}

		if asCfg.Automation != nil {
			if c.Automation == nil {
				c.Automation = asCfg.Automation
			} else {
				err := c.Automation.ApplyOverrides(asCfg.Automation)
				if err != nil {
					return pkg_errors.Wrapf(err, "error applying overrides to Automation config")
				}
			}
		}

		if asCfg.LogPoller != nil {
			if c.LogPoller == nil {
				c.LogPoller = asCfg.LogPoller
			} else {
				err := c.LogPoller.ApplyOverrides(asCfg.LogPoller)
				if err != nil {
					return pkg_errors.Wrapf(err, "error applying overrides to LogPoller config")
				}
			}
		}

		if asCfg.Functions != nil {
			if c.Functions == nil {
				c.Functions = asCfg.Functions
			} else {
				err := c.Functions.ApplyOverrides(asCfg.Functions)
				if err != nil {
					return pkg_errors.Wrapf(err, "error applying overrides to Functions config")
				}
			}
		}
		return nil
	default:
		return pkg_errors.Errorf("cannot apply overrides from unsupported type %T to test config", from)
	}
}

func (c *TestConfig) Validate() error {
	//TODO implement me
	return nil
}
