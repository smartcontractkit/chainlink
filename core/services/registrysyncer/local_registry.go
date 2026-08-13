package registrysyncer

import (
	"maps"

	"github.com/smartcontractkit/capabilities/libs/x/registrysyncer"
	"github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

// The registry snapshot, and the lookups over it, now live in the capabilities repo's
// libs/x/registrysyncer, which crecore shares. They are aliased rather than wrapped so that a
// snapshot handed to a launcher here and one served by crecore are the same type answering the same
// way, and so nothing in core that names these types has to change.
//
// What stays here is only what core alone needs: the copy a launcher is handed (DeepCopyLocalRegistry
// below), the ORM's own table (orm.go), and the on-chain read that fills a snapshot in (syncer.go).
type (
	DonID                   = registrysyncer.DonID
	DON                     = registrysyncer.DON
	CapabilityConfiguration = registrysyncer.CapabilityConfiguration
	Capability              = registrysyncer.Capability
	NodeInfo                = registrysyncer.NodeInfo
	LocalRegistry           = registrysyncer.LocalRegistry
)

// NewLocalRegistry returns a snapshot by value, as core's callers have always taken it.
//
// The shared constructor returns a pointer - a LocalRegistry carries the mutex guarding its
// local-node cache, so handing one back by value copies a lock - and this cannot call it and
// dereference the result for the same reason. It builds the value directly instead, which is what
// the shared constructor does too.
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
		// A config is a value wrapping a byte slice that is never written through, so copying the
		// values into a map of this DON's own is copy enough.
		capCfgs := make(map[string]CapabilityConfiguration, len(don.CapabilityConfigurations))
		maps.Copy(capCfgs, don.CapabilityConfigurations)
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
			WorkflowDONID:       node.WorkflowDONID,
			Signer:              node.Signer,
			P2pID:               node.P2pID,
			EncryptionPublicKey: node.EncryptionPublicKey,
			CapabilityIDs:       make([]string, len(node.CapabilityIDs)),
			CsaKey:              node.CsaKey,
		}
		copy(nodeInfo.CapabilityIDs, node.CapabilityIDs)
		lrCopy.IDsToNodes[id] = nodeInfo
	}

	return lrCopy
}
