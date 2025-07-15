package registrysyncer

import (
	"context"
	"errors"
	"fmt"

	"github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
)

type DonID uint32

type DON struct {
	capabilities.DON
	CapabilityConfigurations map[string]CapabilityConfiguration
}

type CapabilityConfiguration struct {
	Config []byte
}

type Capability struct {
	ID             string
	CapabilityType capabilities.CapabilityType
}

type LocalRegistry struct {
	lggr              logger.Logger
	getPeerID         func() (types.PeerID, error)
	IDsToDONs         map[DonID]DON
	IDsToNodes        map[types.PeerID]kcr.INodeInfoProviderNodeInfo
	IDsToCapabilities map[string]Capability
}

func NewLocalRegistry(
	lggr logger.Logger,
	getPeerID func() (types.PeerID, error),
	IDsToDONs map[DonID]DON,
	IDsToNodes map[types.PeerID]kcr.INodeInfoProviderNodeInfo,
	IDsToCapabilities map[string]Capability,
) LocalRegistry {
	return LocalRegistry{
		lggr:              logger.Named(lggr, "LocalRegistry"),
		getPeerID:         getPeerID,
		IDsToDONs:         IDsToDONs,
		IDsToNodes:        IDsToNodes,
		IDsToCapabilities: IDsToCapabilities,
	}
}

func (l *LocalRegistry) LocalNode(ctx context.Context) (capabilities.Node, error) {
	// Load the current nodes PeerWrapper, this gets us the current node's
	// PeerID, allowing us to contextualize registry information in terms of DON ownership
	// (eg. get my current DON configuration, etc).
	pid, err := l.getPeerID()
	if err != nil {
		return capabilities.Node{}, errors.New("unable to get local node: peerWrapper hasn't started yet")
	}

	return l.NodeByPeerID(ctx, pid)
}

func (l *LocalRegistry) NodeByPeerID(ctx context.Context, peerID types.PeerID) (capabilities.Node, error) {
	err := l.ensureNotEmpty()
	if err != nil {
		return capabilities.Node{}, err
	}
	nodeInfo, ok := l.IDsToNodes[peerID]
	if !ok {
		return capabilities.Node{}, errors.New("could not find peerID " + peerID.String())
	}

	var workflowDON capabilities.DON
	var capabilityDONs []capabilities.DON
	for _, d := range l.IDsToDONs {
		for _, p := range d.Members {
			if p == peerID {
				if d.AcceptsWorkflows {
					// The CapabilitiesRegistry enforces that the DON ID is strictly
					// greater than 0, so if the ID is 0, it means we've not set `workflowDON` initialized above yet.
					if workflowDON.ID == 0 {
						workflowDON = d.DON
						l.lggr.Debug("Workflow DON identified: %+v", workflowDON)
					} else {
						l.lggr.Errorf("Configuration error: node %s belongs to more than one workflowDON", peerID)
					}
				}

				capabilityDONs = append(capabilityDONs, d.DON)
			}
		}
	}

	return capabilities.Node{
		PeerID:              &peerID,
		NodeOperatorID:      nodeInfo.NodeOperatorId,
		Signer:              nodeInfo.Signer,
		EncryptionPublicKey: nodeInfo.EncryptionPublicKey,
		WorkflowDON:         workflowDON,
		CapabilityDONs:      capabilityDONs,
	}, nil
}

func (l *LocalRegistry) ConfigForCapability(ctx context.Context, capabilityID string, donID uint32) (CapabilityConfiguration, error) {
	err := l.ensureNotEmpty()
	if err != nil {
		return CapabilityConfiguration{}, err
	}
	d, ok := l.IDsToDONs[DonID(donID)]
	if !ok {
		return CapabilityConfiguration{}, fmt.Errorf("could not find don %d", donID)
	}

	cc, ok := d.CapabilityConfigurations[capabilityID]
	if !ok {
		return CapabilityConfiguration{}, fmt.Errorf("could not find capability configuration for capability %s and donID %d", capabilityID, donID)
	}

	return cc, nil
}

func (l *LocalRegistry) ensureNotEmpty() error {
	if len(l.IDsToDONs) == 0 {
		return errors.New("empty local registry. no DONs registered in the local registry")
	}
	if len(l.IDsToNodes) == 0 {
		return errors.New("empty local registry. no nodes registered in the local registry")
	}
	if len(l.IDsToCapabilities) == 0 {
		return errors.New("empty local registry. no capabilities registered in the local registry")
	}
	return nil
}
