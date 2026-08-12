package helpers

import (
	"context"
	"encoding/binary"
	"fmt"
	"math/big"
	"time"

	"github.com/gagliardetto/solana-go"
	solRpc "github.com/gagliardetto/solana-go/rpc"
	solstate "github.com/smartcontractkit/cld-changesets/legacy/pkg/family/solana"
	pdasol "github.com/smartcontractkit/cld-changesets/pkg/family/solana"

	proposeutils "github.com/smartcontractkit/cld-changesets/legacy/mcms/proposeutils"

	cldfsol "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldfproposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"

	"github.com/smartcontractkit/mcms"
	"github.com/smartcontractkit/mcms/sdk"
	mcmsSolana "github.com/smartcontractkit/mcms/sdk/solana"
	mcmsTypes "github.com/smartcontractkit/mcms/types"
)

func BuildMCMSTxn(ixn solana.Instruction, programID string, contractType cldf.ContractType) (*mcmsTypes.Transaction, error) {
	data, err := ixn.Data()
	if err != nil {
		return nil, fmt.Errorf("failed to extract data: %w", err)
	}
	for _, account := range ixn.Accounts() {
		if account.IsSigner {
			account.IsSigner = false
		}
	}
	tx, err := mcmsSolana.NewTransaction(
		programID,
		data,
		big.NewInt(0),        // e.g. value
		ixn.Accounts(),       // pass along needed accounts
		string(contractType), // some string identifying the target
		[]string{},           // any relevant metadata
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}
	return &tx, nil
}

// ProgramDataAddress derives the account holding the executable bytes of an upgradeable program.
func ProgramDataAddress(programID solana.PublicKey) solana.PublicKey {
	address, _, _ := solana.FindProgramAddress([][]byte{programID.Bytes()}, solana.BPFLoaderUpgradeableProgramID)

	return address
}

// setUpgradeAuthority creates a transaction to set the upgrade authority for a program
func SetUpgradeAuthority(
	e *cldf.Environment,
	programID solana.PublicKey,
	currentUpgradeAuthority solana.PublicKey,
	newUpgradeAuthority solana.PublicKey,
	isBuffer bool,
) solana.Instruction {
	e.Logger.Infow("Setting upgrade authority", "programID", programID.String(), "currentUpgradeAuthority", currentUpgradeAuthority.String(), "newUpgradeAuthority", newUpgradeAuthority.String())
	// Buffers use the program account as the program data account
	programDataSlice := solana.NewAccountMeta(programID, true, false)
	if !isBuffer {
		// Actual program accounts use the program data account
		programDataAddress, _, _ := solana.FindProgramAddress([][]byte{programID.Bytes()}, solana.BPFLoaderUpgradeableProgramID)
		programDataSlice = solana.NewAccountMeta(programDataAddress, true, false)
	}

	keys := solana.AccountMetaSlice{
		programDataSlice, // Program account (writable)
		solana.NewAccountMeta(currentUpgradeAuthority, false, true), // Current upgrade authority (signer)
		solana.NewAccountMeta(newUpgradeAuthority, false, false),    // New upgrade authority
	}

	instruction := solana.NewInstruction(
		solana.BPFLoaderUpgradeableProgramID,
		keys,
		// https://github.com/solana-playground/solana-playground/blob/2998d4cf381aa319d26477c5d4e6d15059670a75/vscode/src/commands/deploy/bpf-upgradeable/bpf-upgradeable.ts#L72
		[]byte{4, 0, 0, 0}, // 4-byte SetAuthority instruction identifier
	)

	return instruction
}

// UpgradeProgram creates a transaction that replaces the executable bytes of programID with the
// contents of bufferAddress. The buffer must already be owned by the same upgrade authority as the
// program. The buffer is deallocated by the instruction and its lamports are credited to
// spillAddress.
func UpgradeProgram(
	e *cldf.Environment,
	programID solana.PublicKey,
	bufferAddress solana.PublicKey,
	spillAddress solana.PublicKey,
	upgradeAuthority solana.PublicKey,
) solana.Instruction {
	e.Logger.Infow("Upgrading program",
		"programID", programID.String(), "buffer", bufferAddress.String(), "upgradeAuthority", upgradeAuthority.String())

	keys := solana.AccountMetaSlice{
		solana.NewAccountMeta(ProgramDataAddress(programID), true, false), // Program data account (writable)
		solana.NewAccountMeta(programID, true, false),                     // Program account (writable)
		solana.NewAccountMeta(bufferAddress, true, false),                 // Buffer account (writable)
		solana.NewAccountMeta(spillAddress, true, false),                  // Spill account (writable)
		solana.NewAccountMeta(solana.SysVarRentPubkey, false, false),
		solana.NewAccountMeta(solana.SysVarClockPubkey, false, false),
		solana.NewAccountMeta(upgradeAuthority, false, true), // Current upgrade authority (signer)
	}

	return solana.NewInstruction(
		solana.BPFLoaderUpgradeableProgramID,
		keys,
		// https://github.com/solana-playground/solana-playground/blob/2998d4cf381aa319d26477c5d4e6d15059670a75/vscode/src/commands/deploy/bpf-upgradeable/bpf-upgradeable.ts#L66
		[]byte{3, 0, 0, 0}, // 4-byte Upgrade instruction identifier
	)
}

// CloseBuffer creates a transaction that closes a buffer account, refunding its rent to recipient.
func CloseBuffer(
	e *cldf.Environment,
	bufferAddress solana.PublicKey,
	recipient solana.PublicKey,
	upgradeAuthority solana.PublicKey,
) solana.Instruction {
	e.Logger.Infow("Closing buffer", "buffer", bufferAddress.String(), "recipient", recipient.String())

	keys := solana.AccountMetaSlice{
		solana.Meta(bufferAddress).WRITE(),
		solana.Meta(recipient).WRITE(),
		solana.Meta(upgradeAuthority).SIGNER(),
	}

	return solana.NewInstruction(
		solana.BPFLoaderUpgradeableProgramID,
		keys,
		// https://github.com/solana-playground/solana-playground/blob/2998d4cf381aa319d26477c5d4e6d15059670a75/vscode/src/commands/deploy/bpf-upgradeable/bpf-upgradeable.ts#L78
		[]byte{5, 0, 0, 0}, // 4-byte Close instruction identifier
	)
}

// ExtendProgram creates a transaction that grows the program data account of programID by
// extraBytes, which is required when an upgrade brings a binary larger than the deployed one.
// Extending is permissionless, so payer does not have to be the upgrade authority. This matters for
// timelock owned programs, since a PDA cannot pay for the extra rent.
func ExtendProgram(programID solana.PublicKey, payer solana.PublicKey, extraBytes uint32) solana.Instruction {
	// https://github.com/solana-labs/solana/blob/7700cb3128c1f19820de67b81aa45d18f73d2ac0/sdk/program/src/loader_upgradeable_instruction.rs#L146
	data := binary.LittleEndian.AppendUint32([]byte{}, 6) // 4-byte Extend instruction identifier
	data = binary.LittleEndian.AppendUint32(data, extraBytes)

	keys := solana.AccountMetaSlice{
		solana.NewAccountMeta(ProgramDataAddress(programID), true, false), // Program data account (writable)
		solana.NewAccountMeta(programID, true, false),                     // Program account (writable)
		solana.NewAccountMeta(solana.SystemProgramID, false, false),       // System program
		solana.NewAccountMeta(payer, true, true),                          // Payer for rent
	}

	return solana.NewInstruction(solana.BPFLoaderUpgradeableProgramID, keys, data)
}

// GetSolAccountSize returns the length of the data held by an account.
func GetSolAccountSize(ctx context.Context, chain cldfsol.Chain, account solana.PublicKey) (int, error) {
	accountInfo, err := chain.Client.GetAccountInfoWithOpts(ctx, account, &solRpc.GetAccountInfoOpts{
		Commitment: cldfsol.SolDefaultCommitment,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to get account info for %s: %w", account, err)
	}
	if accountInfo == nil || accountInfo.Value == nil {
		return 0, fmt.Errorf("account %s not found", account)
	}

	return len(accountInfo.Value.Data.GetBinary()), nil
}

func BuildProposalsForTxns(
	e cldf.Environment,
	chainSelector uint64,
	description string,
	minDelay time.Duration,
	txns []mcmsTypes.Transaction) (*mcms.TimelockProposal, error) {
	timelocks := map[uint64]string{}
	proposers := map[uint64]string{}
	inspectors := map[uint64]sdk.Inspector{}
	batches := make([]mcmsTypes.BatchOperation, 0)
	chain := e.BlockChains.SolanaChains()[chainSelector]
	addresses := e.DataStore.Addresses().Filter(datastore.AddressRefByChainSelector(chainSelector))
	mcmState, _ := solstate.MaybeLoadMCMSWithTimelockChainStateV2(addresses)

	timelocks[chainSelector] = mcmsSolana.ContractAddress(
		mcmState.TimelockProgram,
		mcmsSolana.PDASeed(mcmState.TimelockSeed),
	)
	proposers[chainSelector] = mcmsSolana.ContractAddress(mcmState.McmProgram, mcmsSolana.PDASeed(mcmState.ProposerMcmSeed))
	inspectors[chainSelector] = mcmsSolana.NewInspector(chain.Client)
	batches = append(batches, mcmsTypes.BatchOperation{
		ChainSelector: mcmsTypes.ChainSelector(chainSelector),
		Transactions:  txns,
	})
	proposal, err := proposeutils.BuildProposalFromBatchesV2(
		e,
		timelocks,
		proposers,
		inspectors,
		batches,
		description,
		cldfproposalutils.TimelockConfig{MinDelay: minDelay})
	if err != nil {
		return nil, fmt.Errorf("failed to build proposal: %w", err)
	}
	return proposal, nil
}

func FetchTimelockSigner(refs []datastore.AddressRef) (solana.PublicKey, error) {
	mcmState, err := solstate.MaybeLoadMCMSWithTimelockChainStateV2(refs)
	if err != nil {
		return solana.PublicKey{}, fmt.Errorf("failed to load mcm state: %w", err)
	}
	timelockSignerPDA := pdasol.GetTimelockSignerPDA(mcmState.TimelockProgram, mcmState.TimelockSeed)
	return timelockSignerPDA, nil
}
