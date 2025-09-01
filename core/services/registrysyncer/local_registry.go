package registrysyncer

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/smartcontractkit/libocr/ragep2p/types"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
)

type DonID uint32

type DON struct {
	capabilities.DON
	CapabilityConfigurations map[string]CapabilityConfiguration
}

type CapabilityConfiguration struct {
	Config []byte
}

func (c CapabilityConfiguration) Unmarshal() (capabilities.CapabilityConfiguration, error) {
	cconf := &capabilitiespb.CapabilityConfig{}
	err := proto.Unmarshal(c.Config, cconf)
	if err != nil {
		return capabilities.CapabilityConfiguration{}, fmt.Errorf("failed to unmarshal capability configuration: %w", err)
	}

	var remoteTriggerConfig *capabilities.RemoteTriggerConfig
	var remoteTargetConfig *capabilities.RemoteTargetConfig

	switch cconf.GetRemoteConfig().(type) {
	case *capabilitiespb.CapabilityConfig_RemoteTriggerConfig:
		prtc := cconf.GetRemoteTriggerConfig()
		remoteTriggerConfig = &capabilities.RemoteTriggerConfig{
			RegistrationRefresh:     prtc.RegistrationRefresh.AsDuration(),
			RegistrationExpiry:      prtc.RegistrationExpiry.AsDuration(),
			MinResponsesToAggregate: prtc.MinResponsesToAggregate,
			MessageExpiry:           prtc.MessageExpiry.AsDuration(),
		}
	case *capabilitiespb.CapabilityConfig_RemoteTargetConfig:
		prtc := cconf.GetRemoteTargetConfig()
		remoteTargetConfig = &capabilities.RemoteTargetConfig{
			RequestHashExcludedAttributes: prtc.RequestHashExcludedAttributes,
		}
	}

	dc, err := values.FromMapValueProto(cconf.DefaultConfig)
	if err != nil {
		return capabilities.CapabilityConfiguration{}, fmt.Errorf("failed to unmarshal capability configuration: %w", err)
	}

	rc, err := values.FromMapValueProto(cconf.RestrictedConfig)
	if err != nil {
		return capabilities.CapabilityConfiguration{}, fmt.Errorf("failed to unmarshal capability configuration: %w", err)
	}

	return capabilities.CapabilityConfiguration{
		DefaultConfig:       dc,
		RestrictedKeys:      cconf.RestrictedKeys,
		RestrictedConfig:    rc,
		RemoteTriggerConfig: remoteTriggerConfig,
		RemoteTargetConfig:  remoteTargetConfig,
	}, nil
}

type Capability struct {
	ID             string
	CapabilityType capabilities.CapabilityType
}

type NodeInfo struct {
	NodeOperatorID      uint32
	ConfigCount         uint32
	WorkflowDONId       uint32
	Signer              [32]byte
	P2pID               [32]byte
	EncryptionPublicKey [32]byte
	CapabilitiesDONIds  []*big.Int
	HashedCapabilityIDs [][32]byte
	CapabilityIDs       []string

	// V2 specific fields
	CsaKey [32]byte
}

type LocalRegistry struct {
	core.UnimplementedCapabilitiesRegistryMetadata

	Logger            logger.Logger
	GetPeerID         func() (types.PeerID, error)
	IDsToDONs         map[DonID]DON
	IDsToNodes        map[types.PeerID]NodeInfo
	IDsToCapabilities map[string]Capability
}

func NewLocalRegistry(
	lggr logger.Logger,
	getPeerID func() (types.PeerID, error),
	idsToDONs map[DonID]DON,
	idsToNodes map[types.PeerID]NodeInfo,
	idsToCapabilities map[string]Capability,
) LocalRegistry {
	return LocalRegistry{
		Logger:            logger.Named(lggr, "LocalRegistry"),
		GetPeerID:         getPeerID,
		IDsToDONs:         idsToDONs,
		IDsToNodes:        idsToNodes,
		IDsToCapabilities: idsToCapabilities,
	}
}

func (l *LocalRegistry) LocalNode(ctx context.Context) (capabilities.Node, error) {
	// Load the current nodes PeerWrapper, this gets us the current node's
	// PeerID, allowing us to contextualize registry information in terms of DON ownership
	// (eg. get my current DON configuration, etc).
	pid, err := l.GetPeerID()
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
						l.Logger.Debug("Workflow DON identified: %+v", workflowDON)
					} else {
						l.Logger.Errorf("Configuration error: node %s belongs to more than one workflowDON", peerID)
					}
				}

				capabilityDONs = append(capabilityDONs, d.DON)
			}
		}
	}

	return capabilities.Node{
		PeerID:              &peerID,
		NodeOperatorID:      nodeInfo.NodeOperatorID,
		Signer:              nodeInfo.Signer,
		EncryptionPublicKey: nodeInfo.EncryptionPublicKey,
		WorkflowDON:         workflowDON,
		CapabilityDONs:      capabilityDONs,
	}, nil
}

func (l *LocalRegistry) DONsForCapability(ctx context.Context, capabilityID string) ([]capabilities.DONWithNodes, error) {
	err := l.ensureNotEmpty()
	if err != nil {
		return []capabilities.DONWithNodes{}, err
	}

	foundDONs := []capabilities.DONWithNodes{}
	for _, don := range l.IDsToDONs {
		for cid := range don.CapabilityConfigurations {
			if cid == capabilityID {
				nodes, err := l.nodesForDON(ctx, don.DON)
				if err != nil {
					return nil, fmt.Errorf("could not fetch nodes for DON %d: %w", don.ID, err)
				}
				donWithNodes := capabilities.DONWithNodes{DON: don.DON, Nodes: nodes}
				foundDONs = append(foundDONs, donWithNodes)
			}
		}
	}

	if len(foundDONs) == 0 {
		return nil, fmt.Errorf("could not find DON for capability %s", capabilityID)
	}

	for _, d := range foundDONs {
		nodes := []capabilities.Node{}
		for _, n := range d.DON.Members {
			node, err := l.NodeByPeerID(ctx, n)
			if err != nil {
				return nil, fmt.Errorf("could not find node for peerID %s: %w", n.String(), err)
			}

			nodes = append(nodes, node)
		}
		(&d).Nodes = nodes
	}

	return foundDONs, nil
}

func (l *LocalRegistry) nodesForDON(ctx context.Context, don capabilities.DON) ([]capabilities.Node, error) {
	nodes := []capabilities.Node{}
	for _, n := range don.Members {
		node, err := l.NodeByPeerID(ctx, n)
		if err != nil {
			return nil, fmt.Errorf("could not find node for peerID %s: %w", n.String(), err)
		}

		nodes = append(nodes, node)
	}
	return nodes, nil
}

func (l *LocalRegistry) ConfigForCapability(ctx context.Context, capabilityID string, donID uint32) (capabilities.CapabilityConfiguration, error) {
	err := l.ensureNotEmpty()
	if err != nil {
		return capabilities.CapabilityConfiguration{}, err
	}
	d, ok := l.IDsToDONs[DonID(donID)]
	if !ok {
		return capabilities.CapabilityConfiguration{}, fmt.Errorf("could not find don %d", donID)
	}

	cc, ok := d.CapabilityConfigurations[capabilityID]
	if !ok {
		return capabilities.CapabilityConfiguration{}, fmt.Errorf("could not find capability configuration for capability %s and donID %d", capabilityID, donID)
	}

	return cc.Unmarshal()
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

func DeepCopyLocalRegistry(lr *LocalRegistry) LocalRegistry {
	var lrCopy LocalRegistry
	lrCopy.Logger = lr.Logger
	lrCopy.GetPeerID = lr.GetPeerID
	lrCopy.IDsToDONs = make(map[DonID]DON, len(lr.IDsToDONs))
	for id, don := range lr.IDsToDONs {
		d := capabilities.DON{
			Name:             don.Name,
			ID:               don.ID,
			Families:         don.Families,
			ConfigVersion:    don.ConfigVersion,
			Members:          make([]types.PeerID, len(don.Members)),
			F:                don.F,
			IsPublic:         don.IsPublic,
			AcceptsWorkflows: don.AcceptsWorkflows,
			Config:           don.Config,
		}
		copy(d.Members, don.Members)
		capCfgs := make(map[string]CapabilityConfiguration, len(don.CapabilityConfigurations))
		for capID, capCfg := range don.CapabilityConfigurations {
			capCfgs[capID] = CapabilityConfiguration{
				Config: capCfg.Config,
			}
		}
		lrCopy.IDsToDONs[id] = DON{
			DON:                      d,
			CapabilityConfigurations: capCfgs,
		}
	}

	lrCopy.IDsToCapabilities = make(map[string]Capability, len(lr.IDsToCapabilities))
	for id, capability := range lr.IDsToCapabilities {
		cp := capability
		lrCopy.IDsToCapabilities[id] = cp
	}

	lrCopy.IDsToNodes = make(map[types.PeerID]NodeInfo, len(lr.IDsToNodes))
	for id, node := range lr.IDsToNodes {
		nodeInfo := NodeInfo{
			NodeOperatorID:      node.NodeOperatorID,
			ConfigCount:         node.ConfigCount,
			WorkflowDONId:       node.WorkflowDONId,
			Signer:              node.Signer,
			P2pID:               node.P2pID,
			EncryptionPublicKey: node.EncryptionPublicKey,
			HashedCapabilityIDs: make([][32]byte, len(node.HashedCapabilityIDs)),
			CapabilityIDs:       make([]string, len(node.CapabilityIDs)),
			CapabilitiesDONIds:  make([]*big.Int, len(node.CapabilitiesDONIds)),
			CsaKey:              node.CsaKey,
		}
		copy(nodeInfo.HashedCapabilityIDs, node.HashedCapabilityIDs)
		copy(nodeInfo.CapabilityIDs, node.CapabilityIDs)
		copy(nodeInfo.CapabilitiesDONIds, node.CapabilitiesDONIds)
		lrCopy.IDsToNodes[id] = nodeInfo
	}

	return lrCopy
}
