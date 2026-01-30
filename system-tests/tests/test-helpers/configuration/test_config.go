package configuration

import (
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
)

// TestConfig holds common test specific configurations related to the test execution
// These configurations are not meant to impact the actual test logic
type TestConfig struct {
	RelativePathToRepoRoot   string
	EnvironmentConfigPath    string
	EnvironmentDirPath       string
	EnvironmentStateFile     string
	EnvironmentArtifactPaths string
	ChipIngressGRPCPort      string
}

// TestEnvironment holds references to the main test components
type TestEnvironment struct {
	Config         *envconfig.Config
	TestConfig     *TestConfig
	Logger         zerolog.Logger
	CreEnvironment *cre.Environment
	Dons           *cre.Dons
}
