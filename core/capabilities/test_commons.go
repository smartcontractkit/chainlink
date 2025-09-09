package capabilities

import (
	"crypto/rand"
	"math/big"

	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"

	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/registrysyncer"
)

// RandomUTF8BytesWord generates a [32]byte array containing random UTF-8 encoded characters.
func RandomUTF8BytesWord() [32]byte {
	var result [32]byte
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	for i := 0; i < 32; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			panic(err)
		}
		result[i] = letters[num.Int64()]
	}
	return result
}

type TestTopology struct {
	workflowDonNodes   []p2ptypes.PeerID
	capabilityDonNodes []p2ptypes.PeerID
}

func MakeNodes(count int) []p2ptypes.PeerID {
	nodes := make([]p2ptypes.PeerID, count)
	for i := range nodes {
		nodes[i] = RandomUTF8BytesWord()
	}
	return nodes
}

func DonMaker(dID uint32, donNodes []p2ptypes.PeerID, acceptWorkflow bool) capabilities.DON {
	return capabilities.DON{
		ID:               dID,
		ConfigVersion:    uint32(0),
		F:                uint8(1),
		IsPublic:         true,
		AcceptsWorkflows: acceptWorkflow,
		Members:          donNodes,
	}
}

func (tt *TestTopology) IDsToNodesMaker(triggerCapID [32]byte) map[p2ptypes.PeerID]registrysyncer.NodeInfo {
	IDsToNodes := map[p2ptypes.PeerID]registrysyncer.NodeInfo{}
	for i := range tt.capabilityDonNodes {
		IDsToNodes[tt.capabilityDonNodes[i]] = registrysyncer.NodeInfo{
			NodeOperatorID:      1,
			Signer:              RandomUTF8BytesWord(),
			P2pID:               tt.capabilityDonNodes[i],
			EncryptionPublicKey: RandomUTF8BytesWord(),
			HashedCapabilityIDs: [][32]byte{triggerCapID},
			CapabilitiesDONIds:  nil,
		}
	}
	for i := range tt.workflowDonNodes {
		IDsToNodes[tt.workflowDonNodes[i]] = registrysyncer.NodeInfo{
			NodeOperatorID:      1,
			Signer:              RandomUTF8BytesWord(),
			P2pID:               tt.workflowDonNodes[i],
			EncryptionPublicKey: RandomUTF8BytesWord(),
		}
	}
	return IDsToNodes
}

// MakeLocalRegistry Function creates LocalRegistry structure populated with 3 DONs:
//   - workflow DON (4 members)
//   - capabilities only DON (4 members)
//   - workflow & capabilities DON (1 member: selected capability DON accepting workflows)
func (tt *TestTopology) MakeLocalRegistry(
	workflowDONID uint32,
	capabilitiesDONID uint32,
	workflowNCapabilitiesDONID uint32,
	triggerCapID [32]byte,
	fullTriggerCapID string,
) *registrysyncer.LocalRegistry {
	return &registrysyncer.LocalRegistry{
		IDsToDONs: map[registrysyncer.DonID]registrysyncer.DON{
			registrysyncer.DonID(workflowDONID): {
				DON: DonMaker(workflowDONID, tt.workflowDonNodes, true),
			},
			registrysyncer.DonID(capabilitiesDONID): {
				DON: DonMaker(capabilitiesDONID, tt.capabilityDonNodes, false),
				CapabilityConfigurations: map[string]registrysyncer.CapabilityConfiguration{
					fullTriggerCapID: {},
				},
			},
			registrysyncer.DonID(workflowNCapabilitiesDONID): {
				DON: DonMaker(workflowNCapabilitiesDONID, tt.capabilityDonNodes[2:3], true),
				CapabilityConfigurations: map[string]registrysyncer.CapabilityConfiguration{
					fullTriggerCapID: {},
				},
			},
		},
		IDsToCapabilities: map[string]registrysyncer.Capability{
			fullTriggerCapID: {
				ID:             fullTriggerCapID,
				CapabilityType: capabilities.CapabilityTypeTrigger,
			},
		},
		IDsToNodes: tt.IDsToNodesMaker(triggerCapID),
	}
}

func NewTestTopology(pid ragetypes.PeerID, workflowNodesCount int, capabilityNodesCount int) *TestTopology {
	tt := TestTopology{}
	tt.workflowDonNodes = MakeNodes(workflowNodesCount)
	tt.capabilityDonNodes = MakeNodes(capabilityNodesCount)
	tt.capabilityDonNodes[0] = pid
	return &tt
}
