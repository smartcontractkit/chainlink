package testconfig

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// MissingImageInfoAsError return a helfpul error message when the no Chainlink image info is found in TOML config.
// If legacy env vars are found it prints ready to use TOML configuration
func MissingImageInfoAsError(errStr string) error {
	missingImage := `
Chainlink image is a secret and must be set as env var in ~/.testsecrets file or passed as env var (either E2E_TEST_CHAINLINK_IMAGE or E2E_TEST_CHAINLINK_UPGRADE_IMAGE). You might have used old configuration approach.
Please refer to integration-tests/testconfig/README.md for more information.
`
	missingVersion := `
Chainlink version must be set in toml config.
`

	if os.Getenv("E2E_TEST_CHAINLINK_IMAGE") == "" || os.Getenv("E2E_TEST_CHAINLINK_UPGRADE_IMAGE") == "" {
		return fmt.Errorf("%s\n%s", errStr, missingImage)
	}
	if os.Getenv("CHAINLINK_VERSION") == "" || os.Getenv("CHAINLINK_UPGRADE_VERSION") == "" {
		return fmt.Errorf("%s\n%s", errStr, missingVersion)
	}
	return errors.New(errStr)
}

// NoSelectedNetworkInfoAsError return a helfpul error message when the no selected network info is found in TOML config.
// If legacy env var is found it prints ready to use TOML configuration.
func NoSelectedNetworkInfoAsError(errStr string) error {
	intro := `
You might have used old configuration approach. If so, use TOML instead of env vars.
Please refer to integration-tests/testconfig/README.md for more information.
`

	if net := os.Getenv("SELECTED_NETWORKS"); net != "" {
		parts := strings.Split(net, ",")
		selectedNetworkStr := "["
		for i, network := range parts {
			selectedNetworkStr += fmt.Sprintf("\"%s\"", network)

			if i < len(parts)-1 {
				selectedNetworkStr += ", "
			}
		}
		selectedNetworkStr += "]"

		extraInfo := `
Or if you want to run your tests right now add following content to integration-tests/testconfig/overrides.toml:
[Network]
selected_networks=`
		return fmt.Errorf("%s\n%s%s%s", errStr, intro, extraInfo, selectedNetworkStr)
	}

	return fmt.Errorf("%s\n%s", errStr, intro)
}

// GetChainAndTestTypeSpecificConfig returns a TestConfig with the chain and test type specific configuration.
// extraFileNames are optional and can be used to specify additional config files to load
// in order to override certain config values (e.g NodeConfig.ChainConfigTOMLByChainID).
func GetChainAndTestTypeSpecificConfig(testType string, product Product, extraFileNames ...string) (TestConfig, error) {
	config, err := GetConfig([]string{testType}, product, extraFileNames...)
	if err != nil {
		return TestConfig{}, fmt.Errorf("error getting config: %w", err)
	}
	config, err = GetConfig(
		[]string{
			testType,
			config.GetNetworkConfig().SelectedNetworks[0],
			fmt.Sprintf("%s-%s", config.GetNetworkConfig().SelectedNetworks[0], testType),
		},
		product,
		extraFileNames...,
	)
	if err != nil {
		return TestConfig{}, fmt.Errorf("error getting config: %w", err)
	}
	return config, err
}
