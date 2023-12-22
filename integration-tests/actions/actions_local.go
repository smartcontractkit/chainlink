// Package actions enables common chainlink interactions
package actions

import (
	"fmt"

	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
)

// UpgradeChainlinkNodeVersions upgrades all Chainlink nodes to a new version, and then runs the test environment
// to apply the upgrades
func UpgradeChainlinkNodeVersionsLocal(
	config *tc.TestConfig,
	nodes ...*test_env.ClNode,
) error {
	if config.ChainlinkUpgradeImage == nil {
		return fmt.Errorf("unable to upgrade node version, [ChainlinkUpgradeImage] was not set, must both a new image or a new version")
	}
	for _, node := range nodes {
		if err := node.UpgradeVersion(config); err != nil {
			return err
		}
	}
	return nil
}
