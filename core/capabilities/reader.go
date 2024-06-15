package capabilities

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/keystone_capability_registry"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
	evmrelaytypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

type remoteRegistryReader struct {
	r           types.ContractReader
	peerWrapper p2ptypes.PeerWrapper
	lggr        logger.Logger
}

var _ reader = (*remoteRegistryReader)(nil)

type hashedCapabilityID [32]byte
type donID uint32

type state struct {
	IDsToDONs         map[donID]kcr.CapabilityRegistryDONInfo
	IDsToNodes        map[p2ptypes.PeerID]kcr.CapabilityRegistryNodeInfo
	IDsToCapabilities map[hashedCapabilityID]kcr.CapabilityRegistryCapability
}

func (r *remoteRegistryReader) LocalNode(ctx context.Context) (capabilities.Node, error) {
	if r.peerWrapper.GetPeer() == nil {
		return capabilities.Node{}, errors.New("unable to get peer: peerWrapper hasn't started yet")
	}

	pid := r.peerWrapper.GetPeer().ID()

	readerState, err := r.state(ctx)
	if err != nil {
		return capabilities.Node{}, fmt.Errorf("failed to get state from registry to determine don ownership: %w", err)
	}

	var workflowDON capabilities.DON
	capabilityDONs := []capabilities.DON{}
	for _, d := range readerState.IDsToDONs {
		for _, p := range d.NodeP2PIds {
			if p == pid {
				if d.AcceptsWorkflows {
					if workflowDON.ID == "" {
						workflowDON = *toDONInfo(d)
					} else {
						r.lggr.Errorf("Configuration error: node %s belongs to more than one workflowDON", pid)
					}
				}

				capabilityDONs = append(capabilityDONs, *toDONInfo(d))
			}
		}
	}

	return capabilities.Node{
		PeerID:         &pid,
		WorkflowDON:    workflowDON,
		CapabilityDONs: capabilityDONs,
	}, nil
}

func (r *remoteRegistryReader) state(ctx context.Context) (state, error) {
	dons := []kcr.CapabilityRegistryDONInfo{}
	err := r.r.GetLatestValue(ctx, "capabilityRegistry", "getDONs", nil, &dons)
	if err != nil {
		return state{}, err
	}

	idsToDONs := map[donID]kcr.CapabilityRegistryDONInfo{}
	for _, d := range dons {
		idsToDONs[donID(d.Id)] = d
	}

	caps := kcr.GetCapabilities{}
	err = r.r.GetLatestValue(ctx, "capabilityRegistry", "getCapabilities", nil, &caps)
	if err != nil {
		return state{}, err
	}

	idsToCapabilities := map[hashedCapabilityID]kcr.CapabilityRegistryCapability{}
	for i, c := range caps.Capabilities {
		idsToCapabilities[caps.HashedCapabilityIds[i]] = c
	}

	nodes := &kcr.GetNodes{}
	err = r.r.GetLatestValue(ctx, "capabilityRegistry", "getNodes", nil, &nodes)
	if err != nil {
		return state{}, err
	}

	idsToNodes := map[p2ptypes.PeerID]kcr.CapabilityRegistryNodeInfo{}
	for _, node := range nodes.NodeInfo {
		idsToNodes[node.P2pId] = node
	}

	return state{IDsToDONs: idsToDONs, IDsToCapabilities: idsToCapabilities, IDsToNodes: idsToNodes}, nil
}

type contractReaderFactory interface {
	NewContractReader(context.Context, []byte) (types.ContractReader, error)
}

func newRemoteRegistryReader(ctx context.Context, lggr logger.Logger, peerWrapper p2ptypes.PeerWrapper, relayer contractReaderFactory, remoteRegistryAddress string) (*remoteRegistryReader, error) {
	contractReaderConfig := evmrelaytypes.ChainReaderConfig{
		Contracts: map[string]evmrelaytypes.ChainContractReader{
			"capabilityRegistry": {
				ContractABI: kcr.CapabilityRegistryABI,
				Configs: map[string]*evmrelaytypes.ChainReaderDefinition{
					"getDONs": {
						ChainSpecificName: "getDONs",
					},
					"getCapabilities": {
						ChainSpecificName: "getCapabilities",
					},
					"getNodes": {
						ChainSpecificName: "getNodes",
					},
				},
			},
		},
	}

	contractReaderConfigEncoded, err := json.Marshal(contractReaderConfig)
	if err != nil {
		return nil, err
	}

	cr, err := relayer.NewContractReader(ctx, contractReaderConfigEncoded)
	if err != nil {
		return nil, err
	}

	err = cr.Bind(ctx, []types.BoundContract{
		{
			Address: remoteRegistryAddress,
			Name:    "capabilityRegistry",
		},
	})
	if err != nil {
		return nil, err
	}

	return &remoteRegistryReader{
		r:           cr,
		peerWrapper: peerWrapper,
		lggr:        lggr,
	}, err
}
