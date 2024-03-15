// Package actions enables common chainlink interactions
package actions

import "github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"

// UpgradeChainlinkNodeVersions upgrades all Chainlink nodes to a new version, and then runs the test environment
// to apply the upgrades
func UpgradeChainlinkNodeVersionsLocal(
	newImage, newVersion string,
	nodes ...*test_env.ClNode,
) error {
	for _, node := range nodes {
		if err := node.UpgradeVersion(newImage, newVersion); err != nil {
			return err
		}
	}
	return nil
}
