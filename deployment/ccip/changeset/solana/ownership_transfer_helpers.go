package solana

import (
	"fmt"
	"math/big"

	"github.com/gagliardetto/solana-go"
	mcmsSolana "github.com/smartcontractkit/mcms/sdk/solana"
	mcmsTypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_offramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/fee_quoter"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
)

type TransferOwnershipFn func(
	proposedAuthority solana.PublicKey,
	configPDA solana.PublicKey,
	authority solana.PublicKey,
) (solana.Instruction, error)

type AcceptOwnershipFn func(
	configPDA solana.PublicKey,
	authority solana.PublicKey,
) (solana.Instruction, error)

// transferAndWrapAcceptOwnership abstracts logic of:
//   - building a “transfer ownership” instruction
//   - confirming on-chain
//   - building an “accept ownership” instruction
//   - wrapping it in an MCMS transaction
//   - returning the mcms transaction for the accept ownership
func transferAndWrapAcceptOwnership(
	buildTransfer TransferOwnershipFn,
	buildAccept AcceptOwnershipFn,
	programID solana.PublicKey, // e.g. token_pool program or router program
	proposedOwner solana.PublicKey, // e.g. usually, the timelock signer PDA
	configPDA solana.PublicKey, // e.g. for routerConfigPDA or a token-pool config
	deployer solana.PublicKey, // the “from” authority
	solChain deployment.SolChain, // used for solChain.Confirm
	label deployment.ContractType, // e.g. "Router" or "TokenPool"
) (mcmsTypes.Transaction, error) {
	// 1. Build the instruction that transfers ownership to the timelock
	ixTransfer, err := buildTransfer(proposedOwner, configPDA, deployer)
	if err != nil {
		return mcmsTypes.Transaction{}, fmt.Errorf("%s: failed to create transfer ownership instruction: %w", label, err)
	}

	// 2. Confirm on-chain
	if err := solChain.Confirm([]solana.Instruction{ixTransfer}); err != nil {
		return mcmsTypes.Transaction{}, fmt.Errorf("%s: failed to confirm transfer on-chain: %w", label, err)
	}

	// 3. Build the “accept ownership” instruction
	ixAccept, err := buildAccept(configPDA, proposedOwner)
	if err != nil {
		return mcmsTypes.Transaction{}, fmt.Errorf("%s: failed to create accept ownership instruction: %w", label, err)
	}
	acceptData, err := ixAccept.Data()
	if err != nil {
		return mcmsTypes.Transaction{}, fmt.Errorf("%s: failed to extract accept data: %w", label, err)
	}

	// 4. Wrap in MCMS transaction
	mcmsTx, err := mcmsSolana.NewTransaction(
		programID.String(),
		acceptData,
		big.NewInt(0),       // e.g. value
		ixAccept.Accounts(), // pass along needed accounts
		string(label),       // some string identifying the target
		[]string{},          // any relevant metadata
	)
	if err != nil {
		return mcmsTypes.Transaction{}, fmt.Errorf("%s: failed to create MCMS transaction: %w", label, err)
	}

	return mcmsTx, nil
}

// transferOwnershipRouter transfers ownership of the router to the timelock.
func transferOwnershipRouter(
	ccipState changeset.CCIPOnChainState,
	chainSelector uint64,
	solChain deployment.SolChain,
	timelockProgramID solana.PublicKey,
	timelockInstanceSeed state.PDASeed,
) ([]mcmsTypes.Transaction, error) {
	var result []mcmsTypes.Transaction

	timelockSignerPDA := state.GetTimelockSignerPDA(timelockProgramID, timelockInstanceSeed)
	state := ccipState.SolChains[chainSelector]

	// The relevant on-chain addresses
	routerProgramID := state.Router
	routerConfigPDA := state.RouterConfigPDA

	// Build specialized closures
	buildTransfer := func(newOwner, config, authority solana.PublicKey) (solana.Instruction, error) {
		ccip_router.SetProgramID(routerProgramID)
		return ccip_router.NewTransferOwnershipInstruction(
			newOwner, config, authority,
		).ValidateAndBuild()
	}
	buildAccept := func(config, newOwnerAuthority solana.PublicKey) (solana.Instruction, error) {
		ccip_router.SetProgramID(routerProgramID)
		// If the router has its own accept function, use that
		ix, err := ccip_router.NewAcceptOwnershipInstruction(
			config, newOwnerAuthority,
		).ValidateAndBuild()
		if err != nil {
			return nil, err
		}
		for _, acc := range ix.Accounts() {
			if acc.PublicKey == newOwnerAuthority {
				acc.IsSigner = false
			}
		}
		return ix, nil
	}

	tx, err := transferAndWrapAcceptOwnership(
		buildTransfer,
		buildAccept,
		routerProgramID,
		timelockSignerPDA, // timelock PDA
		routerConfigPDA,   // config PDA
		solChain.DeployerKey.PublicKey(),
		solChain,
		changeset.Router,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to transfer router ownership: %w", err)
	}

	result = append(result, tx)
	return result, nil
}

// transferOwnershipFeeQuoter transfers ownership of the fee quoter to the timelock.
func transferOwnershipFeeQuoter(
	ccipState changeset.CCIPOnChainState,
	chainSelector uint64,
	solChain deployment.SolChain,
	timelockProgramID solana.PublicKey,
	timelockInstanceSeed state.PDASeed,
) ([]mcmsTypes.Transaction, error) {
	var result []mcmsTypes.Transaction

	timelockSignerPDA := state.GetTimelockSignerPDA(timelockProgramID, timelockInstanceSeed)
	state := ccipState.SolChains[chainSelector]

	// The relevant on-chain addresses
	feeQuoterProgramID := state.FeeQuoter
	feeQuoterConfigPDA := state.FeeQuoterConfigPDA

	// Build specialized closures
	buildTransfer := func(proposedOwner, config, authority solana.PublicKey) (solana.Instruction, error) {
		fee_quoter.SetProgramID(feeQuoterProgramID)
		return fee_quoter.NewTransferOwnershipInstruction(
			proposedOwner, config, authority,
		).ValidateAndBuild()
	}
	buildAccept := func(config, newOwnerAuthority solana.PublicKey) (solana.Instruction, error) {
		fee_quoter.SetProgramID(feeQuoterProgramID)
		// If the router has its own accept function, use that
		ix, err := fee_quoter.NewAcceptOwnershipInstruction(
			config, newOwnerAuthority,
		).ValidateAndBuild()
		if err != nil {
			return nil, err
		}
		for _, acc := range ix.Accounts() {
			if acc.PublicKey == newOwnerAuthority {
				acc.IsSigner = false
			}
		}
		return ix, nil
	}

	tx, err := transferAndWrapAcceptOwnership(
		buildTransfer,
		buildAccept,
		feeQuoterProgramID,
		timelockSignerPDA,  // timelock PDA
		feeQuoterConfigPDA, // config PDA
		solChain.DeployerKey.PublicKey(),
		solChain,
		changeset.FeeQuoter,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to transfer fee quoter ownership: %w", err)
	}

	result = append(result, tx)
	return result, nil
}

// transferOwnershipOffRamp transfers ownership of the offRamp to the timelock.
func transferOwnershipOffRamp(
	ccipState changeset.CCIPOnChainState,
	chainSelector uint64,
	solChain deployment.SolChain,
	timelockProgramID solana.PublicKey,
	timelockInstanceSeed state.PDASeed,
) ([]mcmsTypes.Transaction, error) {
	var result []mcmsTypes.Transaction

	timelockSignerPDA := state.GetTimelockSignerPDA(timelockProgramID, timelockInstanceSeed)
	state := ccipState.SolChains[chainSelector]

	// The relevant on-chain addresses
	offRampProgramID := state.OffRamp
	offRampConfigPDA := state.OffRampConfigPDA

	// Build specialized closures
	buildTransfer := func(proposedOwner, config, authority solana.PublicKey) (solana.Instruction, error) {
		ccip_offramp.SetProgramID(offRampProgramID)
		return ccip_offramp.NewTransferOwnershipInstruction(
			proposedOwner, config, authority,
		).ValidateAndBuild()
	}
	buildAccept := func(config, newOwnerAuthority solana.PublicKey) (solana.Instruction, error) {
		ccip_offramp.SetProgramID(offRampProgramID)
		// If the router has its own accept function, use that
		ix, err := ccip_offramp.NewAcceptOwnershipInstruction(
			config, newOwnerAuthority,
		).ValidateAndBuild()
		if err != nil {
			return nil, err
		}
		for _, acc := range ix.Accounts() {
			if acc.PublicKey == newOwnerAuthority {
				acc.IsSigner = false
			}
		}
		return ix, nil
	}

	tx, err := transferAndWrapAcceptOwnership(
		buildTransfer,
		buildAccept,
		offRampProgramID,
		timelockSignerPDA, // timelock PDA
		offRampConfigPDA,  // config PDA
		solChain.DeployerKey.PublicKey(),
		solChain,
		changeset.OffRamp,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to transfer offRamp ownership: %w", err)
	}

	result = append(result, tx)
	return result, nil
}
