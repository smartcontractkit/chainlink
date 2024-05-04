package testconfig

import (
	"fmt"
	"os"
	"strings"
)

// MissingImageInfoAsError return a helfpul error message when the no Chainlink image info is found in TOML config.
// If legacy env vars are found it prints ready to use TOML configuration
func MissingImageInfoAsError(errStr string) error {
	intro := `
You might have used old configuration approach. If so, use TOML instead of env vars.
Please refer to integration-tests/testconfig/README.md for more information.
`

	var imgStr, versionStr string

	if img := os.Getenv("CHAINLINK_IMAGE"); img != "" {
		imgStr = fmt.Sprintf("image = \"%s\"\n", img)
	}

	if version := os.Getenv("CHAINLINK_VERSION"); version != "" {
		versionStr = fmt.Sprintf("version = \"%s\"\n", version)
	}

	finalErrStr := fmt.Sprintf("%s\n%s", errStr, intro)

	if imgStr != "" && versionStr != "" {
		extraInfo := `
Or if you want to run your tests right now add following content to integration-tests/testconfig/overrides.toml:
[ChainlinkImage]
`
		finalErrStr = fmt.Sprintf("%s\n%s%s%s%s", errStr, intro, extraInfo, imgStr, versionStr)
	}

	return fmt.Errorf(finalErrStr)
}

// NoSelectedNetworkInfoAsError return a helfpul error message when the no selected network info is found in TOML config.
// If legacy env var is found it prints ready to use TOML configuration.
func NoSelectedNetworkInfoAsError(errStr string) error {
	intro := `
You might have used old configuration approach. If so, use TOML instead of env vars.
Please refer to integration-tests/testconfig/README.md for more information.
`

	finalErrStr := fmt.Sprintf("%s\n%s", errStr, intro)

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
		finalErrStr = fmt.Sprintf("%s\n%s%s%s", errStr, intro, extraInfo, selectedNetworkStr)
	}

	return fmt.Errorf(finalErrStr)
}
