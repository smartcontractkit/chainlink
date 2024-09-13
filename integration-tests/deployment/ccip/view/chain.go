package view

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip/view/types"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip/view/types/v1_2"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip/view/types/v1_5"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip/view/types/v1_6"
)

type Chain struct {
	TokenAdminRegistry map[string]types.Snapshotter `json:"tokenAdminRegistry,omitempty"`
	FeeQuoter          map[string]types.Snapshotter `json:"feeQuoter,omitempty"`
	NonceManager       map[string]types.Snapshotter `json:"nonceManager,omitempty"`
	Router             map[string]types.Snapshotter `json:"router,omitempty"`
	RMN                map[string]types.Snapshotter `json:"rmn,omitempty"`
	OnRamp             map[string]types.Snapshotter `json:"onRamp,omitempty"`
}

type ChainContractsMetaData struct {
	TokenAdminRegistry map[string]types.ContractMetaData `json:"tokenAdminRegistry,omitempty"`
	FeeQuoter          map[string]types.ContractMetaData `json:"feeQuoter,omitempty"`
	NonceManager       map[string]types.ContractMetaData `json:"nonceManager,omitempty"`
	Router             map[string]types.ContractMetaData `json:"router,omitempty"`
	RMN                map[string]types.ContractMetaData `json:"rmn,omitempty"`
	OnRamp             map[string]types.ContractMetaData `json:"onRamp,omitempty"`
}

func (ccm *ChainContractsMetaData) SetTokenAdminRegistry(addr common.Address, client bind.ContractBackend) error {
	if ccm.TokenAdminRegistry == nil {
		ccm.TokenAdminRegistry = make(map[string]types.ContractMetaData)
	}
	ta, err := types.NewContractMetaData(addr, client)
	if err != nil {
		return fmt.Errorf("tokenAdminRegistry %s: %w", addr, err)
	}
	ccm.TokenAdminRegistry[addr.String()] = ta
	return nil
}

func (ccm *ChainContractsMetaData) SetFeeQuoter(addr common.Address, client bind.ContractBackend) error {
	if ccm.FeeQuoter == nil {
		ccm.FeeQuoter = make(map[string]types.ContractMetaData)
	}
	fq, err := types.NewContractMetaData(addr, client)
	if err != nil {
		return fmt.Errorf("feeQuoter %s: %w", addr, err)
	}
	ccm.FeeQuoter[addr.String()] = fq
	return nil
}

func (ccm *ChainContractsMetaData) SetRouter(addr common.Address, client bind.ContractBackend) error {
	if ccm.Router == nil {
		ccm.Router = make(map[string]types.ContractMetaData)
	}
	r, err := types.NewContractMetaData(addr, client)
	if err != nil {
		return fmt.Errorf("router %s: %w", addr, err)
	}
	ccm.Router[addr.String()] = r
	return nil

}
func (ccm *ChainContractsMetaData) GetRouterMetas() []types.ContractMetaData {
	routerMetas := make([]types.ContractMetaData, 0)
	for _, meta := range ccm.Router {
		routerMetas = append(routerMetas, meta)
	}
	return routerMetas
}

func (ccm *ChainContractsMetaData) GetTokenAdminRegistryMetas() []types.ContractMetaData {
	tokenAdminRegistryMetas := make([]types.ContractMetaData, 0)
	for _, meta := range ccm.TokenAdminRegistry {
		tokenAdminRegistryMetas = append(tokenAdminRegistryMetas, meta)
	}
	return tokenAdminRegistryMetas
}

func NewChainContractsMetaData() *ChainContractsMetaData {
	return &ChainContractsMetaData{
		TokenAdminRegistry: make(map[string]types.ContractMetaData),
		FeeQuoter:          make(map[string]types.ContractMetaData),
		NonceManager:       make(map[string]types.ContractMetaData),
		Router:             make(map[string]types.ContractMetaData),
		RMN:                make(map[string]types.ContractMetaData),
		OnRamp:             make(map[string]types.ContractMetaData),
	}
}

func NewChain(chainContractsMeta *ChainContractsMetaData, client bind.ContractBackend) (Chain, error) {
	var chain Chain
	for addr, tv := range chainContractsMeta.Router {
		switch tv.TypeAndVersion {
		case types.RouterTypeAndVersionV1_2:
			if chain.Router == nil {
				chain.Router = make(map[string]types.Snapshotter)
			}
			r := &v1_2.Router{}
			err := r.Snapshot(tv, nil, client)
			if err != nil {
				return Chain{}, err
			}
			chain.Router[addr] = r
		}
	}
	for addr, tv := range chainContractsMeta.TokenAdminRegistry {
		switch tv.TypeAndVersion {
		case types.TokenAdminRegistryTypeAndVersionV1_5:
			if chain.TokenAdminRegistry == nil {
				chain.TokenAdminRegistry = make(map[string]types.Snapshotter)
			}
			ta := &v1_5.TokenAdminRegistry{}
			err := ta.Snapshot(tv, nil, client)
			if err != nil {
				return Chain{}, err
			}
			chain.TokenAdminRegistry[addr] = ta
		}
	}
	for addr, tv := range chainContractsMeta.FeeQuoter {
		chainContractsMeta.FeeQuoter[addr] = types.ContractMetaData{
			TypeAndVersion: tv.TypeAndVersion,
			Address:        tv.Address,
		}
		switch tv.TypeAndVersion {
		case types.FEEQuoterTypeAndVersionV1_6:
			if chain.FeeQuoter == nil {
				chain.FeeQuoter = make(map[string]types.Snapshotter)
			}
			fq := &v1_6.FeeQuoter{}
			dependencies := append(chainContractsMeta.GetRouterMetas(), chainContractsMeta.GetTokenAdminRegistryMetas()...)
			err := fq.Snapshot(tv, dependencies, client)
			if err != nil {
				return Chain{}, err
			}
			chain.FeeQuoter[addr] = fq
		}
	}

	return chain, nil
}
