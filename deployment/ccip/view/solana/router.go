package solana

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"
	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"

	solCommon "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_common"
	solRouter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"
)

type RouterView struct {
	PDA                    string                              `json:"pda,omitempty"`
	Version                uint8                               `json:"version,omitempty"`
	DefaultCodeVersion     string                              `json:"defaultCodeVersion,omitempty"`
	SvmChainSelector       uint64                              `json:"svmChainSelector,omitempty"`
	Owner                  string                              `json:"owner,omitempty"`
	ProposedOwner          string                              `json:"proposedOwner,omitempty"`
	FeeQuoter              string                              `json:"feeQuoter,omitempty"`
	RmnRemote              string                              `json:"rmnRemote,omitempty"`
	LinkTokenMint          string                              `json:"linkTokenMint,omitempty"`
	FeeAggregator          string                              `json:"feeAggregator,omitempty"`
	DestinationChainConfig map[uint64]RouterDestChainConfig    `json:"destinationChainConfig,omitempty"`
	TokenAdminRegistry     map[string]RouterTokenAdminRegistry `json:"tokenAdminRegistry,omitempty"`
}

type RouterDestChainConfig struct {
	PDA              string   `json:"pda,omitempty"`
	LaneCodeVersion  string   `json:"laneCodeVersion,omitempty"`
	AllowedSenders   []string `json:"allowedSenders,omitempty"`
	AllowListEnabled bool     `json:"allowListEnabled,omitempty"`
}

type RouterTokenAdminRegistry struct {
	PDA                  string   `json:"pda,omitempty"`
	Version              uint8    `json:"version,omitempty"`
	Administrator        string   `json:"administrator,omitempty"`
	PendingAdministrator string   `json:"pendingAdministrator,omitempty"`
	LookupTable          string   `json:"lookupTable,omitempty"`
	WritableIndexes      []string `json:"writableIndexes,omitempty"`
	Mint                 string   `json:"mint,omitempty"`
}

func GenerateRouterView(chain cldf_solana.Chain, program solana.PublicKey, remoteChains []uint64, tokens []solana.PublicKey) (RouterView, error) {
	view := RouterView{}
	var config solRouter.Config
	configPDA, _, _ := solState.FindConfigPDA(program)
	err := chain.GetAccountDataBorshInto(context.Background(), configPDA, &config)
	if err != nil {
		return view, fmt.Errorf("config not found in existing state, initialize the router first %d", chain.Selector)
	}
	view.PDA = configPDA.String()
	view.DefaultCodeVersion = config.DefaultCodeVersion.String()
	view.SvmChainSelector = config.SvmChainSelector
	view.Owner = config.Owner.String()
	view.ProposedOwner = config.ProposedOwner.String()
	view.FeeQuoter = config.FeeQuoter.String()
	view.RmnRemote = config.RmnRemote.String()
	view.LinkTokenMint = config.LinkTokenMint.String()
	view.FeeAggregator = config.FeeAggregator.String()
	view.DestinationChainConfig = make(map[uint64]RouterDestChainConfig)
	view.TokenAdminRegistry = make(map[string]RouterTokenAdminRegistry)
	for _, remote := range remoteChains {
		remoteChainPDA, err := solState.FindDestChainStatePDA(remote, program)
		if err != nil {
			return view, fmt.Errorf("failed to find dest chain state pda for remote chain %d: %w", remote, err)
		}
		var destChainStateAccount solRouter.DestChain
		if err = chain.GetAccountDataBorshInto(context.Background(), remoteChainPDA, &destChainStateAccount); err != nil {
			return view, fmt.Errorf("remote %d is not configured on solana chain %d", remote, chain.Selector)
		}
		view.DestinationChainConfig[remote] = RouterDestChainConfig{
			PDA:              remoteChainPDA.String(),
			LaneCodeVersion:  destChainStateAccount.Config.LaneCodeVersion.String(),
			AllowedSenders:   make([]string, len(destChainStateAccount.Config.AllowedSenders)),
			AllowListEnabled: destChainStateAccount.Config.AllowListEnabled,
		}
		for i, sender := range destChainStateAccount.Config.AllowedSenders {
			view.DestinationChainConfig[remote].AllowedSenders[i] = sender.String()
		}
	}
	// TODO: save the configured chains/tokens to the AB so we can reconstruct state without the loop
	for _, token := range tokens {
		tokenAdminRegistryPDA, _, err := solState.FindTokenAdminRegistryPDA(token, program)
		if err != nil {
			return view, fmt.Errorf("failed to find token admin registry pda (mint: %s, router: %s): %w", token.String(), program.String(), err)
		}
		var tokenAdminRegistryAccount solCommon.TokenAdminRegistry
		if err := chain.GetAccountDataBorshInto(context.Background(), tokenAdminRegistryPDA, &tokenAdminRegistryAccount); err == nil {
			view.TokenAdminRegistry[token.String()] = RouterTokenAdminRegistry{
				PDA:                  tokenAdminRegistryPDA.String(),
				Version:              tokenAdminRegistryAccount.Version,
				Administrator:        tokenAdminRegistryAccount.Administrator.String(),
				PendingAdministrator: tokenAdminRegistryAccount.PendingAdministrator.String(),
				LookupTable:          tokenAdminRegistryAccount.LookupTable.String(),
				WritableIndexes:      make([]string, len(tokenAdminRegistryAccount.WritableIndexes)),
				Mint:                 token.String(),
			}
			for i, index := range tokenAdminRegistryAccount.WritableIndexes {
				view.TokenAdminRegistry[token.String()].WritableIndexes[i] = index.String()
			}
		}
	}

	return view, nil
}
