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

func (tt *TestTopology) MakeNodes(count int) []p2ptypes.PeerID {
	nodes := make([]p2ptypes.PeerID, count)
	for i := range nodes {
		nodes[i] = randomWord()
	}
	return nodes
}

func (tt *TestTopology) DonMaker(dID uint32, donNodes []p2ptypes.PeerID, acceptWorkflow bool) capabilities.DON {
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

func (tt *TestTopology) MakeLocalRegistry(dID uint32, capDonID uint32, triggerCapID [32]byte, fullTriggerCapID string) *registrysyncer.LocalRegistry {
	return &registrysyncer.LocalRegistry{
		IDsToDONs: map[registrysyncer.DonID]registrysyncer.DON{
			registrysyncer.DonID(dID): {
				DON: tt.DonMaker(dID, tt.workflowDonNodes, true),
			},
			registrysyncer.DonID(capDonID): {
				DON: tt.DonMaker(capDonID, tt.capabilityDonNodes, false),
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
	tt.workflowDonNodes = tt.MakeNodes(workflowNodesCount)
	tt.capabilityDonNodes = tt.MakeNodes(capabilityNodesCount)
	tt.capabilityDonNodes[0] = pid
	return &tt
}
