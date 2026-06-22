package view

import (
	"encoding/json"
	"fmt"
	"slices"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/vault/changeset"
	"github.com/smartcontractkit/chainlink/deployment/vault/changeset/types"
)

var _ cldf.ViewStateV2 = Vault

type VaultView struct {
	TimelockBalances     map[uint64]*types.TimelockNativeBalanceInfo `json:"timelock_balances"`
	WhitelistedAddresses map[uint64][]changeset.WhitelistEntry       `json:"whitelisted_addresses"`
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
	slices.Sort(chainSelectors)

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

	return view, nil
}
