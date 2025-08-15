package view

import (
	"encoding/json"
	"fmt"
	"sort"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/view/v1_0"
	"github.com/smartcontractkit/chainlink/deployment/vault/changeset"
	"github.com/smartcontractkit/chainlink/deployment/vault/changeset/types"
)

var _ cldf.ViewStateV2 = Vault

type VaultView struct {
	TimelockBalances     map[uint64]*types.TimelockNativeBalanceInfo `json:"timelock_balances"`
	WhitelistedAddresses map[uint64][]changeset.WhitelistEntry       `json:"whitelisted_addresses"`
	MCMSWithTimelock     map[uint64]v1_0.MCMSWithTimelockView        `json:"mcms_with_timelock,omitempty"`
}

func (v *VaultView) MarshalJSON() ([]byte, error) {
	type Alias VaultView
	return json.MarshalIndent((*Alias)(v), "", "  ")
}

func Vault(e cldf.Environment, _ json.Marshaler) (json.Marshaler, error) {
	lggr := e.Logger
	lggr.Info("Generating vault state view")

	chainSelectors := make([]uint64, 0)
	for chainSel := range e.BlockChains.EVMChains() {
		chainSelectors = append(chainSelectors, chainSel)
	}
	sort.Slice(chainSelectors, func(i, j int) bool {
		return chainSelectors[i] < chainSelectors[j]
	})

	if len(chainSelectors) == 0 {
		lggr.Warn("No EVM chains found in environment")
		return &VaultView{}, nil
	}

	view, err := GenerateVaultView(e, chainSelectors)
	if err != nil {
		return nil, fmt.Errorf("failed to generate vault view: %w", err)
	}

	return view, nil
}

func GenerateVaultView(e cldf.Environment, chainSelectors []uint64) (*VaultView, error) {
	view := &VaultView{
		WhitelistedAddresses: make(map[uint64][]changeset.WhitelistEntry),
		MCMSWithTimelock:     make(map[uint64]v1_0.MCMSWithTimelockView),
	}

	balances, err := changeset.GetTimelockBalances(e, chainSelectors)
	if err != nil {
		return nil, fmt.Errorf("failed to get timelock balances: %w", err)
	}
	view.TimelockBalances = balances

	addresses, err := changeset.GetWhitelistedAddresses(e, chainSelectors)
	if err != nil {
		return nil, fmt.Errorf("failed to get whitelisted addresses: %w", err)
	}
	view.WhitelistedAddresses = addresses

	mcmsStates, err := state.MaybeLoadMCMSWithTimelockStateDataStore(e, chainSelectors)
	if err != nil {
		e.Logger.Warnf("Failed to load MCMS state (this may be expected if MCMS is not deployed): %v", err)
	} else {
		for chainSelector, mcmsState := range mcmsStates {
			if mcmsState != nil {
				mcmsView, err := mcmsState.GenerateMCMSWithTimelockView()
				if err != nil {
					e.Logger.Warnf("Failed to generate MCMS view for chain %d: %v", chainSelector, err)
				} else {
					view.MCMSWithTimelock[chainSelector] = mcmsView
				}
			}
		}
	}

	return view, nil
}
