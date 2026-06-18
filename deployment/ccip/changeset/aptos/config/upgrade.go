package config

import cldfproposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"

type UpgradeAptosChainConfig struct {
	ChainSelector  uint64
	UpgradeCCIP    bool
	UpgradeOffRamp bool
	UpgradeOnRamp  bool
	UpgradeRouter  bool
	MCMS           *cldfproposalutils.TimelockConfig
}
