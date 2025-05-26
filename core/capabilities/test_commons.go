package capabilities

import (
	"crypto/rand"

	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/registrysyncer"
)

func randomWord() [32]byte {
	word := make([]byte, 32)
	_, err := rand.Read(word)
	if err != nil {
		panic(err)
	}
	return [32]byte(word)
}

type TestTopology struct {
	workflowDonNodes   []p2ptypes.PeerID
	capabilityDonNodes []p2ptypes.PeerID
}

func MakeNodes(count int) []p2ptypes.PeerID {
	nodes := make([]p2ptypes.PeerID, count)
	for i := range nodes {
		nodes[i] = randomWord()
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

func (tt *TestTopology) IDsToNodesMaker(triggerCapID [32]byte) map[p2ptypes.PeerID]kcr.INodeInfoProviderNodeInfo {
	IDsToNodes := map[p2ptypes.PeerID]kcr.INodeInfoProviderNodeInfo{}
	for i := range tt.capabilityDonNodes {
		IDsToNodes[tt.capabilityDonNodes[i]] = kcr.INodeInfoProviderNodeInfo{
			NodeOperatorId:      1,
			Signer:              randomWord(),
			P2pId:               tt.capabilityDonNodes[i],
			EncryptionPublicKey: randomWord(),
			HashedCapabilityIds: [][32]byte{triggerCapID},
			CapabilitiesDONIds:  nil,
		}
	}
	for i := range tt.workflowDonNodes {
		IDsToNodes[tt.workflowDonNodes[i]] = kcr.INodeInfoProviderNodeInfo{
			NodeOperatorId:      1,
			Signer:              randomWord(),
			P2pId:               tt.workflowDonNodes[i],
			EncryptionPublicKey: randomWord(),
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
