// Package actions enables common chainlink interactions
package actions

import (
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
)

// UpgradeChainlinkNodeVersions upgrades all Chainlink nodes to a new version, and then runs the test environment
// to apply the upgrades
func UpgradeChainlinkNodeVersionsLocal(
	newImage, newVersion string,
	nodes ...*test_env.ClNode,
) error {
	if newImage == "" && newVersion == "" {
		return errors.New("unable to upgrade node version, found empty image and version, must provide either a new image or a new version")
	}
	for _, node := range nodes {
		if err := node.UpgradeVersion(node.NodeConfig, newImage, newVersion); err != nil {
			return err
		}
	}
	return nil
}
