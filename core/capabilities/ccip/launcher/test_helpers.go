package launcher

import (
	"fmt"
	"math/big"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/registrysyncer"

	ragep2ptypes "github.com/smartcontractkit/libocr/ragep2p/types"
)

const (
	ccipCapVersion    = "v1.0.0"
	ccipCapNewVersion = "v1.1.0"
	ccipCapName       = "ccip"
)

var (
	defaultCapability = getCapability(ccipCapName, ccipCapVersion)
	newCapability     = getCapability(ccipCapName, ccipCapNewVersion)
	p2pID1            = getP2PID(1)
	p2pID2            = getP2PID(2)
	p2pID3            = getP2PID(3)
	p2pID4            = getP2PID(4)
	defaultCapCfgs    = map[string]registrysyncer.CapabilityConfiguration{
		defaultCapability.ID: {},
	}
	defaultRegistryDon = registrysyncer.DON{
		DON:                      getDON(1, []ragep2ptypes.PeerID{p2pID1}, 0),
		CapabilityConfigurations: defaultCapCfgs,
	}
)

func getP2PID(id uint32) ragep2ptypes.PeerID {
	return ragep2ptypes.PeerID(p2pkey.MustNewV2XXXTestingOnly(big.NewInt(int64(id))).PeerID())
}

func getCapability(ccipCapName, ccipCapVersion string) registrysyncer.Capability {
	id := fmt.Sprintf("%s@%s", ccipCapName, ccipCapVersion)
	return registrysyncer.Capability{
		CapabilityType: capabilities.CapabilityTypeTarget,
		ID:             id,
	}
}

func getDON(id uint32, members []ragep2ptypes.PeerID, cfgVersion uint32) capabilities.DON {
	return capabilities.DON{
		ID:               id,
		ConfigVersion:    cfgVersion,
		F:                uint8(1),
		IsPublic:         true,
		AcceptsWorkflows: true,
		Members:          members,
	}
}
