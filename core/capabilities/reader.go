package capabilities

import (
	"context"
	"encoding/json"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/keystone_capability_registry"
	evmrelaytypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

type remoteRegistryReader struct {
	r types.ContractReader
}

var _ reader = (*remoteRegistryReader)(nil)

type state struct {
	DONs              []kcr.CapabilityRegistryDONInfo
	IDsToNodes        map[[32]byte]kcr.CapabilityRegistryNodeInfo
	IDsToCapabilities map[[32]byte]kcr.CapabilityRegistryCapability
}

func (r *remoteRegistryReader) state(ctx context.Context) (state, error) {
	dons := []kcr.CapabilityRegistryDONInfo{}
	err := r.r.GetLatestValue(ctx, "capabilityRegistry", "getDONs", nil, &dons)
	if err != nil {
		return state{}, err
	}

	caps := kcr.GetCapabilities{}
	err = r.r.GetLatestValue(ctx, "capabilityRegistry", "getCapabilities", nil, &caps)
	if err != nil {
		return state{}, err
	}

	idsToCapabilities := map[[32]byte]kcr.CapabilityRegistryCapability{}
	for i, c := range caps.Capabilities {
		idsToCapabilities[caps.HashedCapabilityIds[i]] = c
	}

	nodes := &kcr.GetNodes{}
	err = r.r.GetLatestValue(ctx, "capabilityRegistry", "getNodes", nil, &nodes)
	if err != nil {
		return state{}, err
	}

	idsToNodes := map[[32]byte]kcr.CapabilityRegistryNodeInfo{}
	for _, node := range nodes.NodeInfo {
		idsToNodes[node.P2pId] = node
	}

	return state{DONs: dons, IDsToCapabilities: idsToCapabilities, IDsToNodes: idsToNodes}, nil
}

type contractReaderFactory interface {
	NewContractReader(context.Context, []byte) (types.ContractReader, error)
}

func newRemoteRegistryReader(ctx context.Context, relayer contractReaderFactory, remoteRegistryAddress string) (*remoteRegistryReader, error) {
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

	return &remoteRegistryReader{r: cr}, err
}
