package cmd

import "github.com/smartcontractkit/chainlink/v2/core/logger"

// CheckRemoteBuildCompatibility exposes checkRemoteBuildCompatibility for testing.
func (s *Shell) CheckRemoteBuildCompatibility(lggr logger.Logger, onlyWarn bool, cliVersion, cliSha string) error {
	return s.checkRemoteBuildCompatibility(lggr, onlyWarn, cliVersion, cliSha)
}

// ConfigV2Str exposes configV2Str for testing.
func (s *Shell) ConfigV2Str(userOnly bool) (string, error) {
	return s.configV2Str(userOnly)
}

type RelayerFactory = relayerFactory
