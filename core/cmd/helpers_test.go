package cmd

import "github.com/smartcontractkit/chainlink/core/logger"

// CheckRemoteBuildCompatibility exposes checkRemoteBuildCompatibility for testing.
func (cli *Client) CheckRemoteBuildCompatibility(lggr logger.Logger, onlyWarn bool, cliVersion, cliSha string) error {
	return cli.checkRemoteBuildCompatibility(lggr, onlyWarn, cliVersion, cliSha)
}
