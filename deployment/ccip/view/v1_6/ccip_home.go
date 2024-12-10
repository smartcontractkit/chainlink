package v1_6

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/view/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/ccip_home"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
)

// https://github.com/smartcontractkit/chainlink/blob/develop/contracts/src/v0.8/ccip/libraries/Internal.sol#L190
const (
	CommitPluginType = 0
	ExecPluginType   = 1
)

type DonView struct {
	DonID         uint32                  `json:"donID"`
	CommitConfigs ccip_home.GetAllConfigs `json:"commitConfigs"`
	ExecConfigs   ccip_home.GetAllConfigs `json:"execConfigs"`
}

type CCIPHomeView struct {
	types.ContractMetaData
	ChainConfigs       []ccip_home.CCIPHomeChainConfigArgs `json:"chainConfigs"`
	CapabilityRegistry common.Address                      `json:"capabilityRegistry"`
	Dons               []DonView                           `json:"dons"`
}

func GenerateCCIPHomeView(c deployment.OnchainClient, ch *ccip_home.CCIPHome) (CCIPHomeView, error) {
	if ch == nil {
		return CCIPHomeView{}, fmt.Errorf("cannot generate view for nil CCIPHomeView")
	}
	meta, err := types.NewContractMetaData(ch, ch.Address())
	if err != nil {
		return CCIPHomeView{}, fmt.Errorf("failed to generate contract metadata for CCIPHomeView: %w", err)
	}
	numChains, err := ch.GetNumChainConfigurations(nil)
	if err != nil {
		return CCIPHomeView{}, fmt.Errorf("failed to get number of chain configurations for CCIPHomeView: %w", err)
	}
	// Pagination shouldn't be required here, but we can add it if needed.
	chains, err := ch.GetAllChainConfigs(nil, big.NewInt(0), numChains)
	if err != nil {
		return CCIPHomeView{}, fmt.Errorf("failed to get all chain configs for CCIPHomeView: %w", err)
	}
	crAddr, err := ch.GetCapabilityRegistry(nil)
	if err != nil {
		return CCIPHomeView{}, fmt.Errorf("failed to get capability registry for CCIPHomeView: %w", err)
	}
	cr, err := capabilities_registry.NewCapabilitiesRegistry(crAddr, c)
	if err != nil {
		return CCIPHomeView{}, fmt.Errorf("failed to instantiate capability registry for CCIPHomeView: %w", err)
	}
	dons, err := cr.GetDONs(nil)
	if err != nil {
		return CCIPHomeView{}, fmt.Errorf("failed to get DONs for CCIPHomeView: %w", err)
	}
	// Get every don's configuration.
	var dvs []DonView
	for _, d := range dons {
		commitConfigs, err := ch.GetAllConfigs(nil, d.Id, CommitPluginType)
		if err != nil {
			return CCIPHomeView{}, fmt.Errorf("failed to get active commit config for CCIPHomeView: %w", err)
		}
		execConfigs, err := ch.GetAllConfigs(nil, d.Id, ExecPluginType)
		if err != nil {
			return CCIPHomeView{}, fmt.Errorf("failed to get active commit config for CCIPHomeView: %w", err)
		}
		dvs = append(dvs, DonView{
			DonID:         d.Id,
			CommitConfigs: commitConfigs,
			ExecConfigs:   execConfigs,
		})
	}
	return CCIPHomeView{
		ContractMetaData:   meta,
		ChainConfigs:       chains,
		CapabilityRegistry: crAddr,
		Dons:               dvs,
	}, nil
}
