package solana

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/mcm"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/timelock"
	"github.com/smartcontractkit/chainlink/deployment/ccip/view/shared"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

type MCMSWithTimelockView struct {
	Timelock TimelockView `json:"timelock,omitempty"`
	MCMS     MCMSView     `json:"mcm,omitempty"`
}

type TimelockView struct {
	ProgramID                     string   `json:"programId,omitempty"`
	Owner                         string   `json:"owner,omitempty"`
	ProposedOwner                 string   `json:"proposedOwner,omitempty"`
	ProposerRoleAccessController  string   `json:"proposerRoleAccessController,omitempty"`
	ExecutorRoleAccessController  string   `json:"executorRoleAccessController,omitempty"`
	CancellerRoleAccessController string   `json:"cancellerRoleAccessController,omitempty"`
	BypasserRoleAccessController  string   `json:"bypasserRoleAccessController,omitempty"`
	MinDelay                      uint64   `json:"minDelay,omitempty"`
	BlockedSelectors              []string `json:"blockedSelectors,omitempty"`
}

type MCMSView struct {
	Bypasser  MCMSConfig `json:"bypasser,omitempty"`
	Proposer  MCMSConfig `json:"proposer,omitempty"`
	Canceller MCMSConfig `json:"canceller,omitempty"`
}

type MCMSConfig struct {
	ProgramID     string   `json:"programId,omitempty"`
	ChainID       uint64   `json:"chainId,omitempty"`
	MultisigID    string   `json:"multisigId,omitempty"`
	Owner         string   `json:"owner,omitempty"`
	ProposedOwner string   `json:"proposedOwner,omitempty"`
	GroupQuorums  string   `json:"groupQuorums,omitempty"`
	GroupParents  string   `json:"groupParents,omitempty"`
	Signers       []string `json:"signers,omitempty"`
}

func GenerateMCMSWithTimelockView(chain cldf.SolChain, addresses map[string]cldf.TypeAndVersion) (MCMSWithTimelockView, error) {
	view := MCMSWithTimelockView{}
	mcmState, err := state.MaybeLoadMCMSWithTimelockChainStateSolana(chain, addresses)
	if err != nil {
		return view, fmt.Errorf("failed to load mcms with timelock solana chain state: %w", err)
	}
	timelockConfigPDA := state.GetTimelockConfigPDA(mcmState.TimelockProgram, mcmState.TimelockSeed)
	var timelockData timelock.Config
	err = chain.GetAccountDataBorshInto(context.Background(), timelockConfigPDA, &timelockData)
	if err != nil {
		return view, fmt.Errorf("fee quoter config not found in existing state, initialize the fee quoter first %d", chain.Selector)
	}
	view.Timelock = TimelockView{
		Owner:                         timelockData.Owner.String(),
		ProposedOwner:                 timelockData.ProposedOwner.String(),
		ProposerRoleAccessController:  timelockData.ProposerRoleAccessController.String(),
		ExecutorRoleAccessController:  timelockData.ExecutorRoleAccessController.String(),
		CancellerRoleAccessController: timelockData.CancellerRoleAccessController.String(),
		BypasserRoleAccessController:  timelockData.BypasserRoleAccessController.String(),
		MinDelay:                      timelockData.MinDelay,
	}
	view.MCMS = MCMSView{}
	var mcmData mcm.MultisigConfig
	for _, mcmConfig := range []struct {
		name     string
		pda      solana.PublicKey
		mcmsView MCMSConfig
	}{
		{"Bypasser", state.GetMCMConfigPDA(mcmState.McmProgram, mcmState.BypasserMcmSeed), view.MCMS.Bypasser},
		{"Proposer", state.GetMCMConfigPDA(mcmState.McmProgram, mcmState.ProposerMcmSeed), view.MCMS.Proposer},
		{"Canceller", state.GetMCMConfigPDA(mcmState.McmProgram, mcmState.CancellerMcmSeed), view.MCMS.Canceller},
	} {
		err = chain.GetAccountDataBorshInto(context.Background(), mcmConfig.pda, &mcmData)
		if err != nil {
			return view, fmt.Errorf("failed to get account data for %s: %w", mcmConfig.name, err)
		}
		currConfig := MCMSConfig{
			ProgramID:     mcmState.McmProgram.String(),
			ChainID:       mcmData.ChainId,
			MultisigID:    string(mcmData.MultisigId[:]),
			Owner:         mcmData.Owner.String(),
			ProposedOwner: mcmData.ProposedOwner.String(),
			GroupQuorums:  string(mcmData.GroupQuorums[:]),
			GroupParents:  string(mcmData.GroupParents[:]),
		}
		for _, signer := range mcmData.Signers {
			mcmConfig.mcmsView.Signers = append(mcmConfig.mcmsView.Signers, shared.GetAddressFromBytes(chain_selectors.ETHEREUM_MAINNET.Selector, signer.EvmAddress[:]))
		}
		mcmConfig.mcmsView = currConfig
	}

	return view, nil
}
