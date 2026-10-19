package operation

import (
	"fmt"
	"slices"

	"github.com/gagliardetto/solana-go"
	"github.com/smartcontractkit/mcms"
	"github.com/smartcontractkit/mcms/sdk"
	mcmsSolana "github.com/smartcontractkit/mcms/sdk/solana"
	mcmsTypes "github.com/smartcontractkit/mcms/types"

	solanaUtils "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldfproposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	df_store "github.com/smartcontractkit/chainlink-solana/contracts/generated/data_feeds_store"
	proposeutils "github.com/smartcontractkit/cld-changesets/legacy/mcms/proposeutils"
	solstate "github.com/smartcontractkit/cld-changesets/legacy/pkg/family/solana"

	commonOps "github.com/smartcontractkit/chainlink/deployment/common/changeset/solana/operations"
	"github.com/smartcontractkit/chainlink/deployment/helpers"
)

var (
	DeployStoreOp = operations.NewOperation(
		"deploy-data-feeds-store-op",
		Version1_0_0,
		"Deploys the DataFeeds Store program for Solana Chain",
		commonOps.Deploy,
	)
	InitStoreOp = operations.NewOperation(
		"init-data-feeds-store-op",
		Version1_0_0,
		"Initialize DataFeeds Store for Solana Chain",
		initStore,
	)
	SendInstructionsOp = operations.NewOperation(
		"send-data-feeds-store-instructions-op",
		Version1_0_0,
		"Sends instructions with the deployer, or proposes them to the timelock",
		sendInstructions,
	)
)

// maxInstructionsPerBatch keeps each timelock batch operation within a single Solana transaction.
const maxInstructionsPerBatch = 1

type (
	InitStoreInput struct {
		ProgramID                 solana.PublicKey
		ForwarderProgramID        solana.PublicKey
		ChainSel                  uint64
		FeedAdmins                []solana.PublicKey // strictly increasing
		DefaultMaxStoredReports   uint64
		MaxObservationForwardSecs uint32
	}

	InitStoreOutput struct {
		StatePubKey solana.PublicKey
	}

	// Instruction is a serializable Solana instruction.
	Instruction struct {
		ProgramID solana.PublicKey      `json:"programId"`
		Accounts  []*solana.AccountMeta `json:"accounts"`
		Data      []byte                `json:"data"`
	}

	SendInstructionsInput struct {
		ChainSel     uint64
		MCMS         *cldfproposalutils.TimelockConfig // if set, the instructions are proposed to the timelock
		ContractType cldf.ContractType
		Description  string
		Instructions []Instruction
	}

	SendInstructionsOutput struct {
		Proposals []mcms.TimelockProposal // returned when a timelock config is passed
	}
)

func initStore(b operations.Bundle, deps Deps, in InitStoreInput) (InitStoreOutput, error) {
	var out InitStoreOutput

	stateKey, err := solana.NewRandomPrivateKey()
	if err != nil {
		return out, fmt.Errorf("failed to create random keys: %w", err)
	}

	ix, err := df_store.NewInitializeInstruction(
		in.FeedAdmins,
		in.DefaultMaxStoredReports,
		in.MaxObservationForwardSecs,
		deps.Chain.DeployerKey.PublicKey(),
		stateKey.PublicKey(),
		in.ForwarderProgramID,
		solana.SystemProgramID,
	)
	if err != nil {
		return out, fmt.Errorf("failed to build initialize instruction: %w", err)
	}
	data, err := ix.Data()
	if err != nil {
		return out, fmt.Errorf("failed to encode initialize instruction: %w", err)
	}

	if err = deps.Chain.Confirm([]solana.Instruction{solana.NewInstruction(in.ProgramID, ix.Accounts(), data)}, solanaUtils.AddSigners(stateKey)); err != nil {
		return out, fmt.Errorf("failed to confirm initialize: %w", err)
	}

	out.StatePubKey = stateKey.PublicKey()
	return out, nil
}

func sendInstructions(b operations.Bundle, deps Deps, in SendInstructionsInput) (SendInstructionsOutput, error) {
	var out SendInstructionsOutput

	ixs := make([]solana.Instruction, len(in.Instructions))
	for i, ix := range in.Instructions {
		ixs[i] = solana.NewInstruction(ix.ProgramID, ix.Accounts, ix.Data)
	}

	if in.MCMS == nil {
		for _, ix := range ixs {
			if err := deps.Chain.Confirm([]solana.Instruction{ix}); err != nil {
				return out, fmt.Errorf("failed to confirm instruction: %w", err)
			}
		}
		return out, nil
	}

	txns := make([]mcmsTypes.Transaction, 0, len(ixs))
	for _, ix := range ixs {
		tx, err := helpers.BuildMCMSTxn(ix, ix.ProgramID().String(), in.ContractType)
		if err != nil {
			return out, fmt.Errorf("failed to create transaction: %w", err)
		}
		txns = append(txns, *tx)
	}

	mcmState, err := solstate.MaybeLoadMCMSWithTimelockChainStateV2(
		deps.Env.DataStore.Addresses().Filter(datastore.AddressRefByChainSelector(in.ChainSel)))
	if err != nil {
		return out, fmt.Errorf("failed to load mcms state: %w", err)
	}

	var batches []mcmsTypes.BatchOperation
	for batch := range slices.Chunk(txns, maxInstructionsPerBatch) {
		batches = append(batches, mcmsTypes.BatchOperation{
			ChainSelector: mcmsTypes.ChainSelector(in.ChainSel),
			Transactions:  batch,
		})
	}

	proposal, err := proposeutils.BuildProposalFromBatchesV2(
		deps.Env,
		map[uint64]string{in.ChainSel: mcmsSolana.ContractAddress(mcmState.TimelockProgram, mcmsSolana.PDASeed(mcmState.TimelockSeed))},
		map[uint64]string{in.ChainSel: mcmsSolana.ContractAddress(mcmState.McmProgram, mcmsSolana.PDASeed(mcmState.ProposerMcmSeed))},
		map[uint64]sdk.Inspector{in.ChainSel: mcmsSolana.NewInspector(deps.Chain.Client)},
		batches,
		in.Description,
		*in.MCMS,
	)
	if err != nil {
		return out, fmt.Errorf("failed to build proposal: %w", err)
	}
	out.Proposals = []mcms.TimelockProposal{*proposal}
	return out, nil
}
