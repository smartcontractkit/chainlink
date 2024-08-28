package utils

import (
	"fmt"
	"testing"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	ctf_config "github.com/smartcontractkit/chainlink-testing-framework/config"
	pkg_seth "github.com/smartcontractkit/chainlink-testing-framework/seth"
	seth_utils "github.com/smartcontractkit/chainlink-testing-framework/utils/seth"
)

// DynamicArtifactDirConfigFn returns a function that sets Seth's artifacts directory to a unique directory for the test
func DynamicArtifactDirConfigFn(t *testing.T) func(*pkg_seth.Config) error {
	return func(cfg *pkg_seth.Config) error {
		cfg.ArtifactsDir = fmt.Sprintf("seth_artifacts/%s", t.Name())
		return nil
	}
}

// TestAwareSethClient returns a Seth client with the artifacts directory set to a unique directory for the test
func TestAwareSethClient(t *testing.T, sethConfig ctf_config.SethConfig, evmNetwork *blockchain.EVMNetwork) (*pkg_seth.Client, error) {
	return seth_utils.GetChainClientWithConfigFunction(sethConfig, *evmNetwork, DynamicArtifactDirConfigFn(t))
}
