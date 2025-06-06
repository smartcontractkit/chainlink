package solana

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"
	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"

	solOffRamp "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_offramp"
	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"

	"github.com/smartcontractkit/chainlink/deployment/ccip/view/shared"
)

type OffRampView struct {
	PDA                        string                              `json:"pda,omitempty"`
	Version                    uint8                               `json:"version,omitempty"`
	DefaultCodeVersion         uint8                               `json:"defaultCodeVersion,omitempty"`
	SvmChainSelector           uint64                              `json:"svmChainSelector,omitempty"`
	EnableManualExecutionAfter int64                               `json:"enableManualExecutionAfter,omitempty"`
	Owner                      string                              `json:"owner,omitempty"`
	ProposedOwner              string                              `json:"proposedOwner,omitempty"`
	ReferenceAddresses         OffRampReferenceAddresses           `json:"referenceAddresses,omitempty"`
	SourceChains               map[uint64]OffRampSourceChainConfig `json:"sourceChains,omitempty"`
}

type OffRampReferenceAddresses struct {
	PDA                string `json:"pda,omitempty"`
	Version            uint8  `json:"version,omitempty"`
	Router             string `json:"router,omitempty"`
	FeeQuoter          string `json:"feeQuoter,omitempty"`
	OfframpLookupTable string `json:"offrampLookupTable,omitempty"`
	RmnRemote          string `json:"rmnRemote,omitempty"`
}

type OffRampSourceChainConfig struct {
	PDA                       string `json:"pda,omitempty"`
	IsEnabled                 bool   `json:"isEnabled,omitempty"`
	IsRmnVerificationDisabled bool   `json:"isRmnVerificationDisabled,omitempty"`
	LaneCodeVersion           string `json:"laneCodeVersion,omitempty"`
	OnRamp                    string `json:"onRamp,omitempty"`
}

func GenerateOffRampView(chain cldf_solana.Chain, program solana.PublicKey, remoteChains []uint64, tokens []solana.PublicKey) (OffRampView, error) {
	view := OffRampView{}
	var config solOffRamp.Config
	configPDA, _, _ := solState.FindOfframpConfigPDA(program)
	err := chain.GetAccountDataBorshInto(context.Background(), configPDA, &config)
	if err != nil {
		return view, fmt.Errorf("config not found in existing state, initialize the off ramp first %d", chain.Selector)
	}
	view.PDA = configPDA.String()
	view.DefaultCodeVersion = config.DefaultCodeVersion
	view.SvmChainSelector = config.SvmChainSelector
	view.Owner = config.Owner.String()
	view.ProposedOwner = config.ProposedOwner.String()
	view.EnableManualExecutionAfter = config.EnableManualExecutionAfter
	view.SourceChains = make(map[uint64]OffRampSourceChainConfig)

	var referenceAddressesAccount solOffRamp.ReferenceAddresses
	offRampReferenceAddressesPDA, _, _ := solState.FindOfframpReferenceAddressesPDA(program)
	if err = chain.GetAccountDataBorshInto(context.Background(), offRampReferenceAddressesPDA, &referenceAddressesAccount); err != nil {
		return view, fmt.Errorf("failed to get offramp reference addresses: %w", err)
	}
	view.ReferenceAddresses = OffRampReferenceAddresses{
		PDA:                offRampReferenceAddressesPDA.String(),
		Version:            referenceAddressesAccount.Version,
		Router:             referenceAddressesAccount.Router.String(),
		FeeQuoter:          referenceAddressesAccount.FeeQuoter.String(),
		OfframpLookupTable: referenceAddressesAccount.OfframpLookupTable.String(),
		RmnRemote:          referenceAddressesAccount.RmnRemote.String(),
	}
	for _, remote := range remoteChains {
		remoteChainPDA, _, err := solState.FindOfframpSourceChainPDA(remote, program)
		if err != nil {
			return view, fmt.Errorf("failed to find source chain state pda for remote chain %d: %w", remote, err)
		}
		var chainStateAccount solOffRamp.SourceChain
		if err = chain.GetAccountDataBorshInto(context.Background(), remoteChainPDA, &chainStateAccount); err != nil {
			return view, fmt.Errorf("remote %d is not configured on solana chain %d", remote, chain.Selector)
		}
		onRamp := chainStateAccount.Config.OnRamp
		view.SourceChains[remote] = OffRampSourceChainConfig{
			PDA:                       remoteChainPDA.String(),
			IsEnabled:                 chainStateAccount.Config.IsEnabled,
			IsRmnVerificationDisabled: chainStateAccount.Config.IsRmnVerificationDisabled,
			LaneCodeVersion:           chainStateAccount.Config.LaneCodeVersion.String(),
			OnRamp:                    shared.GetAddressFromBytes(remote, onRamp.Bytes[:onRamp.Len]),
		}
	}
	return view, nil
}
