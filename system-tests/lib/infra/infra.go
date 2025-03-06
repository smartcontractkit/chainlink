package infra

import (
	"fmt"

	cretypes "github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
	"github.com/smartcontractkit/chainlink/system-tests/lib/types"
)

func Host(nodeIndex int, nodeType cretypes.CapabilityFlag, donName string, infraDetails types.InfraDetails) string {
	if infraDetails.InfraType == types.InfraType_CRIB {
		if nodeType == cretypes.BootstrapNode {
			return fmt.Sprintf("base-bt-%d", nodeIndex)
		}
		return fmt.Sprintf("base-%d", nodeIndex)
	}

	return fmt.Sprintf("%s-node%d", donName, nodeIndex)
}
