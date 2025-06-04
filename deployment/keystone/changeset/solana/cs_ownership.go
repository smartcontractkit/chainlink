package solana

import (
	"fmt"

	"github.com/gagliardetto/solana-go"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	ks_forwarder "github.com/smartcontractkit/chainlink-solana/contracts/generated/keystone_forwarder"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	commonstate "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/helpers"
	"github.com/smartcontractkit/mcms"
	"github.com/smartcontractkit/mcms/sdk"
	mcmsSolana "github.com/smartcontractkit/mcms/sdk/solana"
	mcmsTypes "github.com/smartcontractkit/mcms/types"
)

type TransferOwnershipForwarderRequest struct {
	ChainSel                    uint64
	CurrentOwner, ProposedOwner solana.PublicKey

	// MCMSCfg is for the accept ownership proposal
	MCMSCfg proposalutils.TimelockConfig
}

var _ cldf.ChangeSet[*TransferOwnershipForwarderRequest] = TransferOwnershipForwarder

func TransferOwnershipForwarder(env cldf.Environment, req *TransferOwnershipForwarderRequest) (cldf.ChangesetOutput, error) {
	solChain, ok := env.BlockChains.SolanaChains()[req.ChainSel]
	if !ok {
		return cldf.ChangesetOutput{}, fmt.Errorf("chain not found for selector %d", req.ChainSel)
	}

	addresses, _ := env.ExistingAddresses.AddressesForChain(req.ChainSel)

	mcmState, _ := commonstate.MaybeLoadMCMSWithTimelockChainStateSolana(solChain, addresses)
	if mcmState.TimelockProgram.IsZero() {
		return cldf.ChangesetOutput{}, fmt.Errorf("timelock is not found")
	}

	currentOwner := solChain.DeployerKey.PublicKey()
	if !req.CurrentOwner.IsZero() {
		currentOwner = req.CurrentOwner
	}

	timelockSigner := commonstate.GetTimelockSignerPDA(mcmState.TimelockProgram, mcmState.TimelockSeed)
	proposedOwner := timelockSigner

	if !req.ProposedOwner.IsZero() {
		proposedOwner = req.ProposedOwner
	}

	state, err := loadOnchainState(env, req.ChainSel)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	if state.forwarderProgramID.IsZero() {
		return cldf.ChangesetOutput{}, fmt.Errorf("forwarder not found for chain selector %d", req.ChainSel)
	}

	if state.forwarderState.IsZero() {
		return cldf.ChangesetOutput{}, fmt.Errorf("forwarder state not found for chain selector %d", req.ChainSel)
	}

	buildTransfer := func(
		proposedAuthority solana.PublicKey,
		configPDA solana.PublicKey,
		authority solana.PublicKey,
	) (solana.Instruction, error) {
		ks_forwarder.SetProgramID(state.forwarderProgramID)
		ix, err := ks_forwarder.NewTransferOwnershipInstruction(proposedAuthority, configPDA, authority).ValidateAndBuild()
		if err != nil {
			return nil, err
		}
		for _, acc := range ix.Accounts() {
			if acc.PublicKey == timelockSigner {
				acc.IsSigner = false
			}
		}
		return ix, nil
	}

	buildAccept := func(config, newOwnerAuthority solana.PublicKey) (solana.Instruction, error) {
		ks_forwarder.SetProgramID(state.forwarderProgramID)
		ix, err := ks_forwarder.NewAcceptOwnershipInstruction(
			config, newOwnerAuthority,
		).ValidateAndBuild()
		if err != nil {
			return nil, err
		}
		for _, acc := range ix.Accounts() {
			if acc.PublicKey == timelockSigner {
				acc.IsSigner = false
			}
		}
		return ix, nil
	}

	mcmsTx, err := helpers.TransferAndWrapAcceptOwnership(
		buildTransfer,
		buildAccept,
		state.forwarderProgramID,
		proposedOwner,        // timelock PDA
		state.forwarderState, // state (not PDA but forwarder owned acc)
		currentOwner,
		solChain,
		shared.Router,
		timelockSigner, // the timelock signer PDA
	)

	batch := mcmsTypes.BatchOperation{
		ChainSelector: mcmsTypes.ChainSelector(req.ChainSel),
		Transactions:  []mcmsTypes.Transaction{mcmsTx},
	}

	timelocks := map[uint64]string{}
	proposers := map[uint64]string{}
	inspectors := map[uint64]sdk.Inspector{}
	timelocks[solChain.Selector] = mcmsSolana.ContractAddress(
		mcmState.TimelockProgram,
		mcmsSolana.PDASeed(mcmState.TimelockSeed),
	)
	proposers[solChain.Selector] = mcmsSolana.ContractAddress(mcmState.McmProgram, mcmsSolana.PDASeed(mcmState.ProposerMcmSeed))
	inspectors[solChain.Selector] = mcmsSolana.NewInspector(solChain.Client)

	proposal, err := proposalutils.BuildProposalFromBatchesV2(
		env,
		timelocks,
		proposers,
		inspectors,
		[]mcmsTypes.BatchOperation{batch},
		"proposal to transfer ownership of keystone forwarder contract to timelock",
		req.MCMSCfg)

	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
	}

	return cldf.ChangesetOutput{MCMSTimelockProposals: []mcms.TimelockProposal{*proposal}}, nil
}
