package cmd

import "github.com/smartcontractkit/chainlink/v2/core/logger"

// CheckRemoteBuildCompatibility exposes checkRemoteBuildCompatibility for testing.
func (cli *Client) CheckRemoteBuildCompatibility(lggr logger.Logger, onlyWarn bool, cliVersion, cliSha string) error {
	return cli.checkRemoteBuildCompatibility(lggr, onlyWarn, cliVersion, cliSha)
}

// ConfigV2Str exposes configV2Str for testing.
func (cli *Client) ConfigV2Str(userOnly bool) (string, error) {
	return cli.configV2Str(userOnly)
}
