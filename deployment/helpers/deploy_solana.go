package helpers

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/Masterminds/semver/v3"
	solBinary "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	solRpc "github.com/gagliardetto/solana-go/rpc"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/mcms"
	"github.com/smartcontractkit/mcms/sdk"
	mcmsSolana "github.com/smartcontractkit/mcms/sdk/solana"
	mcmsTypes "github.com/smartcontractkit/mcms/types"
)

// DeployAndMaybeSaveToAddressBook deploys a program to the Solana chain and saves it to the address book
// if it is not an upgrade. It returns the program ID of the deployed program.
func DeployAndMaybeSaveToAddressBook(
	e cldf.Environment,
	chain cldf.SolChain,
	ab cldf.AddressBook,
	contractType cldf.ContractType,
	version semver.Version,
	isUpgrade bool,
	metadata string) (solana.PublicKey, error) {
	programName := getTypeToProgramDeployName()[contractType]
	overallocate := true
	// by default we want to overallocate buffers, but if metadata is set (i.e. we're managing partner programs)
	// we want to set the overallocate flag to false
	if metadata != "" && metadata != shared.CLLMetadata {
		overallocate = false
	}
	programID, err := chain.DeployProgram(e.Logger, cldf.SolProgramInfo{
		Name:  programName,
		Bytes: deployment.SolanaProgramBytes[programName],
	}, isUpgrade, overallocate)
	if err != nil {
		return solana.PublicKey{}, fmt.Errorf("failed to deploy program: %w", err)
	}
	address := solana.MustPublicKeyFromBase58(programID)

	e.Logger.Infow("Deployed program", "Program", contractType, "addr", programID, "chain", chain.String(), "isUpgrade", isUpgrade)

	if !isUpgrade {
		tv := cldf.NewTypeAndVersion(contractType, version)
		if metadata != "" {
			tv.AddLabel(metadata)
		}
		err = ab.Save(chain.Selector, programID, tv)
		if err != nil {
			return solana.PublicKey{}, fmt.Errorf("failed to save address: %w", err)
		}
	}
	return address, nil
}

// UPGRADE FUNCTIONS
func GenerateUpgradeTxns(
	e cldf.Environment,
	chain cldf.SolChain,
	ab cldf.AddressBook,
	spillAddress solana.PublicKey,
	upgradeAuthority solana.PublicKey,
	newVersion *semver.Version,
	programID solana.PublicKey,
	contractType cldf.ContractType,
) ([]mcmsTypes.Transaction, error) {
	e.Logger.Infow("Generating instruction for upgrading contract", "contractType", contractType)
	txns := make([]mcmsTypes.Transaction, 0)
	bufferProgram, err := DeployAndMaybeSaveToAddressBook(e, chain, ab, contractType, *newVersion, true, "")
	if err != nil {
		return txns, fmt.Errorf("failed to deploy program: %w", err)
	}

	ixn := SetUpgradeAuthority(&e, bufferProgram, chain.DeployerKey.PublicKey(), upgradeAuthority, true)
	if err := chain.Confirm([]solana.Instruction{ixn}); err != nil {
		return txns, fmt.Errorf("failed to confirm setUpgradeAuthority: %w", err)
	}
	upgradeIxn, err := generateUpgradeIxn(
		programID,
		bufferProgram,
		spillAddress,
		upgradeAuthority,
	)
	if err != nil {
		return txns, fmt.Errorf("failed to generate upgrade instruction: %w", err)
	}
	closeIxn, err := generateCloseBufferIxn(
		bufferProgram,
		spillAddress,
		upgradeAuthority,
	)
	if err != nil {
		return txns, fmt.Errorf("failed to generate close buffer instruction: %w", err)
	}
	addresses, err := e.ExistingAddresses.AddressesForChain(chain.Selector)
	if err != nil {
		return txns, fmt.Errorf("failed to get existing addresses: %w", err)
	}
	mcmState, err := state.MaybeLoadMCMSWithTimelockChainStateSolana(chain, addresses)
	if err != nil {
		return txns, fmt.Errorf("failed to load MCMS with timelock chain state: %w", err)
	}
	timelockSignerPDA := state.GetTimelockSignerPDA(mcmState.TimelockProgram, mcmState.TimelockSeed)
	// if we're not upgrading via timelock, execute the raw ixns
	if upgradeAuthority != timelockSignerPDA {
		ixns := []solana.Instruction{upgradeIxn}
		extendIxn, err := generateExtendIxn(
			&e,
			chain,
			programID,
			bufferProgram,
			upgradeAuthority,
		)
		if err != nil {
			return txns, fmt.Errorf("failed to generate extend buffer instruction: %w", err)
		}
		if extendIxn != nil {
			ixns = append(ixns, extendIxn)
		}
		ixns = append(ixns, closeIxn)
		if err := chain.Confirm(ixns); err != nil {
			return txns, fmt.Errorf("failed to confirm instructions: %w", err)
		}
		return []mcmsTypes.Transaction{}, nil
	}
	upgradeTx, err := BuildMCMSTxn(upgradeIxn, solana.BPFLoaderUpgradeableProgramID.String(), contractType)
	if err != nil {
		return txns, fmt.Errorf("failed to create upgrade transaction: %w", err)
	}
	closeTx, err := BuildMCMSTxn(closeIxn, solana.BPFLoaderUpgradeableProgramID.String(), contractType)
	if err != nil {
		return txns, fmt.Errorf("failed to create close transaction: %w", err)
	}
	// We do not support extend as part of upgrades due to MCMS limitations
	// https://docs.google.com/document/d/1Fk76lOeyS2z2X6MokaNX_QTMFAn5wvSZvNXJluuNV1E/edit?tab=t.0#heading=h.uij286zaarkz
	txns = append(txns, *upgradeTx, *closeTx)
	return txns, nil
}

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

func generateUpgradeIxn(
	programID solana.PublicKey,
	bufferAddress solana.PublicKey,
	spillAddress solana.PublicKey,
	upgradeAuthority solana.PublicKey,
) (solana.Instruction, error) {
	// Derive the program data address
	programDataAccount, _, _ := solana.FindProgramAddress([][]byte{programID.Bytes()}, solana.BPFLoaderUpgradeableProgramID)

	// Accounts involved in the transaction
	keys := solana.AccountMetaSlice{
		solana.NewAccountMeta(programDataAccount, true, false), // Program account (writable)
		solana.NewAccountMeta(programID, true, false),
		solana.NewAccountMeta(bufferAddress, true, false),             // Buffer account (writable)
		solana.NewAccountMeta(spillAddress, true, false),              // Spill account (writable)
		solana.NewAccountMeta(solana.SysVarRentPubkey, false, false),  // System program
		solana.NewAccountMeta(solana.SysVarClockPubkey, false, false), // System program
		solana.NewAccountMeta(upgradeAuthority, false, true),          // Current upgrade authority (signer)
	}

	instruction := solana.NewInstruction(
		solana.BPFLoaderUpgradeableProgramID,
		keys,
		// https://github.com/solana-playground/solana-playground/blob/2998d4cf381aa319d26477c5d4e6d15059670a75/vscode/src/commands/deploy/bpf-upgradeable/bpf-upgradeable.ts#L66
		[]byte{3, 0, 0, 0}, // 4-byte Upgrade instruction identifier
	)

	return instruction, nil
}

func generateExtendIxn(
	e *cldf.Environment,
	chain cldf.SolChain,
	programID solana.PublicKey,
	bufferAddress solana.PublicKey,
	payer solana.PublicKey,
) (*solana.GenericInstruction, error) {
	// Derive the program data address
	programDataAccount, _, _ := solana.FindProgramAddress([][]byte{programID.Bytes()}, solana.BPFLoaderUpgradeableProgramID)

	programDataSize, err := GetSolProgramSize(e, chain, programDataAccount)
	if err != nil {
		return nil, fmt.Errorf("failed to get program size: %w", err)
	}
	e.Logger.Debugw("Program data size", "programDataSize", programDataSize)

	bufferSize, err := GetSolProgramSize(e, chain, bufferAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get buffer size: %w", err)
	}
	e.Logger.Debugw("Buffer account size", "bufferSize", bufferSize)
	if bufferSize <= programDataSize {
		e.Logger.Debugf("Buffer account size %d is less than program account size %d", bufferSize, programDataSize)
		return nil, nil
	}
	extraBytes := bufferSize - programDataSize
	if extraBytes > math.MaxUint32 {
		return nil, fmt.Errorf("extra bytes %d exceeds maximum value %d", extraBytes, math.MaxUint32)
	}
	// https://github.com/solana-labs/solana/blob/7700cb3128c1f19820de67b81aa45d18f73d2ac0/sdk/program/src/loader_upgradeable_instruction.rs#L146
	data := binary.LittleEndian.AppendUint32([]byte{}, 6) // 4-byte Extend instruction identifier
	//nolint:gosec // G115 we check for overflow above
	data = binary.LittleEndian.AppendUint32(data, uint32(extraBytes+1024)) // add some padding

	keys := solana.AccountMetaSlice{
		solana.NewAccountMeta(programDataAccount, true, false),      // Program data account (writable)
		solana.NewAccountMeta(programID, true, false),               // Program account (writable)
		solana.NewAccountMeta(solana.SystemProgramID, false, false), // System program
		solana.NewAccountMeta(payer, true, true),                    // Payer for rent
	}

	ixn := solana.NewInstruction(
		solana.BPFLoaderUpgradeableProgramID,
		keys,
		data,
	)

	return ixn, nil
}

func generateCloseBufferIxn(
	bufferAddress solana.PublicKey,
	recipient solana.PublicKey,
	upgradeAuthority solana.PublicKey,
) (solana.Instruction, error) {
	keys := solana.AccountMetaSlice{
		solana.NewAccountMeta(bufferAddress, true, false),
		solana.NewAccountMeta(recipient, true, false),
		solana.NewAccountMeta(upgradeAuthority, false, true),
	}

	instruction := solana.NewInstruction(
		solana.BPFLoaderUpgradeableProgramID,
		keys,
		// https://github.com/solana-playground/solana-playground/blob/2998d4cf381aa319d26477c5d4e6d15059670a75/vscode/src/commands/deploy/bpf-upgradeable/bpf-upgradeable.ts#L78
		[]byte{5, 0, 0, 0}, // 4-byte Close instruction identifier
	)

	return instruction, nil
}

// HELPER FUNCTIONS
func GetSolProgramSize(e *cldf.Environment, chain cldf.SolChain, programID solana.PublicKey) (int, error) {
	accountInfo, err := chain.Client.GetAccountInfoWithOpts(e.GetContext(), programID, &rpc.GetAccountInfoOpts{
		Commitment: cldf.SolDefaultCommitment,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to get account info: %w", err)
	}
	if accountInfo == nil {
		return 0, fmt.Errorf("program account not found: %w", err)
	}
	programBytes := len(accountInfo.Value.Data.GetBinary())
	return programBytes, nil
}

func GetSolProgramData(e cldf.Environment, chain cldf.SolChain, programID solana.PublicKey) (struct {
	DataType uint32
	Address  solana.PublicKey
}, error) {
	var programData struct {
		DataType uint32
		Address  solana.PublicKey
	}
	data, err := chain.Client.GetAccountInfoWithOpts(e.GetContext(), programID, &solRpc.GetAccountInfoOpts{
		Commitment: solRpc.CommitmentConfirmed,
	})
	if err != nil {
		return programData, fmt.Errorf("failed to deploy program: %w", err)
	}

	err = solBinary.UnmarshalBorsh(&programData, data.Bytes())
	if err != nil {
		return programData, fmt.Errorf("failed to unmarshal program data: %w", err)
	}
	return programData, nil
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
	chain := e.SolChains[chainSelector]
	addresses, _ := e.ExistingAddresses.AddressesForChain(chainSelector)
	mcmState, _ := state.MaybeLoadMCMSWithTimelockChainStateSolana(chain, addresses)

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
	proposal, err := proposalutils.BuildProposalFromBatchesV2(
		e,
		timelocks,
		proposers,
		inspectors,
		batches,
		description,
		proposalutils.TimelockConfig{MinDelay: minDelay})
	if err != nil {
		return nil, fmt.Errorf("failed to build proposal: %w", err)
	}
	return proposal, nil
}

func FetchTimelockSigner(e cldf.Environment, chainSelector uint64) (solana.PublicKey, error) {
	addresses, err := e.ExistingAddresses.AddressesForChain(chainSelector)
	if err != nil {
		return solana.PublicKey{}, fmt.Errorf("failed to load addresses for chain %d: %w", chainSelector, err)
	}
	mcmState, err := state.MaybeLoadMCMSWithTimelockChainStateSolana(e.SolChains[chainSelector], addresses)
	if err != nil {
		return solana.PublicKey{}, fmt.Errorf("failed to load mcm state: %w", err)
	}
	timelockSignerPDA := state.GetTimelockSignerPDA(mcmState.TimelockProgram, mcmState.TimelockSeed)
	return timelockSignerPDA, nil
}

type CloseBuffersConfig struct {
	ChainSelector uint64
	Buffers       []string
}

func CloseBuffersChangeset(e cldf.Environment, cfg CloseBuffersConfig) (cldf.ChangesetOutput, error) {
	for _, buffer := range cfg.Buffers {
		if err := e.SolChains[cfg.ChainSelector].CloseBuffers(e.Logger, buffer); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to close buffer: %w", err)
		}
	}
	return cldf.ChangesetOutput{}, nil
}
