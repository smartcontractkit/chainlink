package config

import "github.com/smartcontractkit/chainlink/deployment/common/proposalutils"

type UpgradeAptosChainConfig struct {
	ChainSelector  uint64
	UpgradeCCIP    bool
	UpgradeOffRamp bool
	UpgradeOnRamp  bool
	UpgradeRouter  bool
	MCMS           *proposalutils.TimelockConfig
}
