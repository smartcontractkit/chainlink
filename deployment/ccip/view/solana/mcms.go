package solana

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gagliardetto/solana-go"
	chain_selectors "github.com/smartcontractkit/chain-selectors"

	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"

	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/mcm"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/timelock"

	solanashared "github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/view/shared"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

type MCMSWithTimelockView struct {
	Timelock TimelockView `json:"timelock"`
	MCMS     MCMSView     `json:"mcm"`
}

type TimelockView struct {
	PDA                           string   `json:"pda,omitempty"`
	UpgradeAuthority              string   `json:"upgradeAuthority,omitempty"`
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
	Bypasser  MCMSConfig `json:"bypasser"`
	Proposer  MCMSConfig `json:"proposer"`
	Canceller MCMSConfig `json:"canceller"`
}

type MCMSConfig struct {
	PDA              string   `json:"pda,omitempty"`
	UpgradeAuthority string   `json:"upgradeAuthority,omitempty"`
	ProgramID        string   `json:"programId,omitempty"`
	ChainID          uint64   `json:"chainId,omitempty"`
	MultisigID       string   `json:"multisigId,omitempty"`
	Owner            string   `json:"owner,omitempty"`
	ProposedOwner    string   `json:"proposedOwner,omitempty"`
	GroupQuorums     string   `json:"groupQuorums,omitempty"`
	GroupParents     string   `json:"groupParents,omitempty"`
	Signers          []string `json:"signers,omitempty"`
}

func GenerateMCMSWithTimelockView(chain cldf_solana.Chain, addresses map[string]cldf.TypeAndVersion) (MCMSWithTimelockView, error) {
	view := MCMSWithTimelockView{}
	mcmState, err := state.MaybeLoadMCMSWithTimelockChainStateSolana(chain, addresses)
	if err != nil {
		return view, fmt.Errorf("failed to load mcms with timelock solana chain state: %w", err)
	}
	timelockConfigPDA := state.GetTimelockConfigPDA(mcmState.TimelockProgram, mcmState.TimelockSeed)
	progDataAddr, err := solanashared.GetProgramDataAddress(chain.Client, mcmState.TimelockProgram)
	if err != nil {
		return view, fmt.Errorf("failed to get program data address for program %s: %w", mcmState.TimelockProgram.String(), err)
	}
	authority, _, err := solanashared.GetUpgradeAuthority(chain.Client, progDataAddr)
	if err != nil {
		return view, fmt.Errorf("failed to get upgrade authority for program data %s: %w", progDataAddr.String(), err)
	}
	var timelockData timelock.Config
	err = chain.GetAccountDataBorshInto(context.Background(), timelockConfigPDA, &timelockData)
	if err != nil {
		return view, fmt.Errorf("timelock config not found in existing state, initialize the timelock first %d", chain.Selector)
	}
	view.Timelock = TimelockView{
		PDA:                           timelockConfigPDA.String(),
		UpgradeAuthority:              authority.String(),
		Owner:                         timelockData.Owner.String(),
		ProposedOwner:                 timelockData.ProposedOwner.String(),
		ProposerRoleAccessController:  timelockData.ProposerRoleAccessController.String(),
		ExecutorRoleAccessController:  timelockData.ExecutorRoleAccessController.String(),
		CancellerRoleAccessController: timelockData.CancellerRoleAccessController.String(),
		BypasserRoleAccessController:  timelockData.BypasserRoleAccessController.String(),
		MinDelay:                      timelockData.MinDelay,
	}
	progDataAddr, err = solanashared.GetProgramDataAddress(chain.Client, mcmState.McmProgram)
	if err != nil {
		return view, fmt.Errorf("failed to get program data address for program %s: %w", mcmState.McmProgram.String(), err)
	}
	authority, _, err = solanashared.GetUpgradeAuthority(chain.Client, progDataAddr)
	if err != nil {
		return view, fmt.Errorf("failed to get upgrade authority for program data %s: %w", progDataAddr.String(), err)
	}
	view.MCMS = MCMSView{
		Bypasser:  MCMSConfig{},
		Proposer:  MCMSConfig{},
		Canceller: MCMSConfig{},
	}
	var mcmData mcm.MultisigConfig
	for _, mcmConfig := range []struct {
		name string
		pda  solana.PublicKey
	}{
		{"Bypasser", state.GetMCMConfigPDA(mcmState.McmProgram, mcmState.BypasserMcmSeed)},
		{"Proposer", state.GetMCMConfigPDA(mcmState.McmProgram, mcmState.ProposerMcmSeed)},
		{"Canceller", state.GetMCMConfigPDA(mcmState.McmProgram, mcmState.CancellerMcmSeed)},
	} {
		err = chain.GetAccountDataBorshInto(context.Background(), mcmConfig.pda, &mcmData)
		if err != nil {
			return view, fmt.Errorf("failed to get account data for %s: %w", mcmConfig.name, err)
		}
		currConfig := MCMSConfig{
			PDA:              mcmConfig.pda.String(),
			ProgramID:        mcmState.McmProgram.String(),
			UpgradeAuthority: authority.String(),
			ChainID:          mcmData.ChainId,
			MultisigID:       string(mcmData.MultisigId[:]),
			Owner:            mcmData.Owner.String(),
			ProposedOwner:    mcmData.ProposedOwner.String(),
			GroupQuorums:     toJSONString(mcmData.GroupQuorums),
			GroupParents:     toJSONString(mcmData.GroupParents),
			Signers:          []string{},
		}
		for _, signer := range mcmData.Signers {
			currConfig.Signers = append(currConfig.Signers, shared.GetAddressFromBytes(chain_selectors.ETHEREUM_MAINNET.Selector, signer.EvmAddress[:]))
		}
		switch mcmConfig.name {
		case "Bypasser":
			view.MCMS.Bypasser = currConfig
		case "Proposer":
			view.MCMS.Proposer = currConfig
		case "Canceller":
			view.MCMS.Canceller = currConfig
		default:
			return view, fmt.Errorf("unknown mcm config name: %s", mcmConfig.name)
		}
	}

	return view, nil
}

func toJSONString(arr [32]uint8) string {
	b, _ := json.Marshal(arr)
	return string(b)
}
