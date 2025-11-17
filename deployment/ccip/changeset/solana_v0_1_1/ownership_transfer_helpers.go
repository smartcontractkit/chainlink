package solana

import (
	"fmt"

	"github.com/gagliardetto/solana-go"
	mcmsTypes "github.com/smartcontractkit/mcms/types"

	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"

	burnmint "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/burnmint_token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/ccip_offramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/ccip_router"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/cctp_token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/fee_quoter"
	lockrelease "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/lockrelease_token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/rmn_remote"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
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

// TransferAndWrapAcceptOwnership abstracts logic of:
//   - building a “transfer ownership” instruction
//   - confirming on-chain
//   - building an “accept ownership” instruction
//   - wrapping it in an MCMS transaction
//   - returning the mcms transaction for the accept ownership
func TransferAndWrapAcceptOwnership(
	buildTransfer TransferOwnershipFn,
	buildAccept AcceptOwnershipFn,
	programID solana.PublicKey, // e.g. token_pool program or router program
	proposedOwner solana.PublicKey, // e.g. usually, the timelock signer PDA
	configPDA solana.PublicKey, // e.g. for routerConfigPDA or a token-pool config
	currentOwner solana.PublicKey, // the “from” authority
	solChain cldf_solana.Chain, // used for solChain.Confirm
	label cldf.ContractType, // e.g. "Router" or "TokenPool"
	timelockSigner solana.PublicKey, // the timelock signer PDA
) (mcmsTypes.Transaction, error) {
	// 1. Build the instruction that transfers ownership to the timelock
	ixTransfer, err := buildTransfer(proposedOwner, configPDA, currentOwner)
	if err != nil {
		return mcmsTypes.Transaction{}, fmt.Errorf("%s: failed to create transfer ownership instruction: %w", label, err)
	}

	// if the old owner is the timelock signer, we can skip the on-chain confirmation
	// We can't perform the accept ownership step here because the timelock signer is not a signer of the transaction
	// 2. Wrap in MCMS transaction or confirm on-chain
	if currentOwner.Equals(timelockSigner) {
		mcmsTx, err := BuildMCMSTxn(ixTransfer, programID.String(), label)
		if err != nil {
			return mcmsTypes.Transaction{}, fmt.Errorf("%s: failed to create MCMS transaction: %w", label, err)
		}
		return *mcmsTx, nil
	}

	if err := solChain.Confirm([]solana.Instruction{ixTransfer}); err != nil {
		return mcmsTypes.Transaction{}, fmt.Errorf("%s: failed to confirm transfer on-chain: %w", label, err)
	}

	// 3. Build the “accept ownership” instruction
	ixAccept, err := buildAccept(configPDA, proposedOwner)
	if err != nil {
		return mcmsTypes.Transaction{}, fmt.Errorf("%s: failed to create accept ownership instruction: %w", label, err)
	}

	// 4. Wrap in MCMS transaction or confirm on-chain
	if proposedOwner.Equals(timelockSigner) {
		mcmsTx, err := BuildMCMSTxn(ixAccept, programID.String(), label)
		if err != nil {
			return mcmsTypes.Transaction{}, fmt.Errorf("%s: failed to create MCMS transaction: %w", label, err)
		}

		return *mcmsTx, nil
	}

	if err := solChain.Confirm([]solana.Instruction{ixAccept}); err != nil {
		return mcmsTypes.Transaction{}, fmt.Errorf("%s: failed to confirm transfer on-chain: %w", label, err)
	}
	return mcmsTypes.Transaction{}, nil
}

// transferOwnershipRouter transfers ownership of the router to the timelock.
func transferOwnershipRouter(
	ccipState stateview.CCIPOnChainState,
	chainSelector uint64,
	solChain cldf_solana.Chain,
	currentOwner solana.PublicKey,
	proposedOwner solana.PublicKey,
	timelockSigner solana.PublicKey,
) ([]mcmsTypes.Transaction, error) {
	var result []mcmsTypes.Transaction

	state := ccipState.SolChains[chainSelector]

	// The relevant on-chain addresses
	routerProgramID := state.Router
	routerConfigPDA := state.RouterConfigPDA

	// Build specialized closures
	buildTransfer := func(newOwner, config, authority solana.PublicKey) (solana.Instruction, error) {
		ix, err := ccip_router.NewTransferOwnershipInstruction(
			newOwner, config, authority,
		).ValidateAndBuild()
		if err != nil {
			return nil, err
		}
		ixData, err := ix.Data()
		if err != nil {
			return nil, fmt.Errorf("failed to extract data payload from router transfer ownership instruction: %w", err)
		}
		transferOwnershipIx := solana.NewInstruction(routerProgramID, ix.Accounts(), ixData)
		for _, acc := range transferOwnershipIx.Accounts() {
			if acc.PublicKey == timelockSigner {
				acc.IsSigner = false
			}
		}
		return transferOwnershipIx, nil
	}
	buildAccept := func(config, newOwnerAuthority solana.PublicKey) (solana.Instruction, error) {
		// If the router has its own accept function, use that
		ix, err := ccip_router.NewAcceptOwnershipInstruction(
			config, newOwnerAuthority,
		).ValidateAndBuild()
		if err != nil {
			return nil, err
		}
		ixData, err := ix.Data()
		if err != nil {
			return nil, fmt.Errorf("failed to extract data payload from router transfer ownership instruction: %w", err)
		}
		acceptOwnershipIx := solana.NewInstruction(routerProgramID, ix.Accounts(), ixData)
		for _, acc := range acceptOwnershipIx.Accounts() {
			if acc.PublicKey == timelockSigner {
				acc.IsSigner = false
			}
		}
		return acceptOwnershipIx, nil
	}

	tx, err := TransferAndWrapAcceptOwnership(
		buildTransfer,
		buildAccept,
		routerProgramID,
		proposedOwner,   // timelock PDA
		routerConfigPDA, // config PDA
		currentOwner,
		solChain,
		shared.Router,
		timelockSigner, // the timelock signer PDA
	)

	if err != nil {
		return nil, fmt.Errorf("failed to transfer router ownership: %w", err)
	}

	result = append(result, tx)
	return result, nil
}

// transferOwnershipFeeQuoter transfers ownership of the fee quoter to the timelock.
func transferOwnershipFeeQuoter(
	ccipState stateview.CCIPOnChainState,
	chainSelector uint64,
	solChain cldf_solana.Chain,
	currentOwner solana.PublicKey,
	proposedOwner solana.PublicKey,
	timelockSigner solana.PublicKey,
) ([]mcmsTypes.Transaction, error) {
	var result []mcmsTypes.Transaction

	state := ccipState.SolChains[chainSelector]

	// The relevant on-chain addresses
	feeQuoterProgramID := state.FeeQuoter
	feeQuoterConfigPDA := state.FeeQuoterConfigPDA

	// Build specialized closures
	buildTransfer := func(proposedOwner, config, authority solana.PublicKey) (solana.Instruction, error) {
		ix, err := fee_quoter.NewTransferOwnershipInstruction(
			proposedOwner, config, authority,
		).ValidateAndBuild()
		if err != nil {
			return nil, err
		}
		ixData, err := ix.Data()
		if err != nil {
			return nil, fmt.Errorf("failed to extract data payload from fee quoter transfer ownership instruction: %w", err)
		}
		transferOwnershipIx := solana.NewInstruction(feeQuoterProgramID, ix.Accounts(), ixData)
		for _, acc := range transferOwnershipIx.Accounts() {
			if acc.PublicKey == timelockSigner {
				acc.IsSigner = false
			}
		}
		return transferOwnershipIx, nil
	}
	buildAccept := func(config, newOwnerAuthority solana.PublicKey) (solana.Instruction, error) {
		// If the router has its own accept function, use that
		ix, err := fee_quoter.NewAcceptOwnershipInstruction(
			config, newOwnerAuthority,
		).ValidateAndBuild()
		if err != nil {
			return nil, err
		}
		ixData, err := ix.Data()
		if err != nil {
			return nil, fmt.Errorf("failed to extract data payload from fee quoter accept ownership instruction: %w", err)
		}
		acceptOwnershipIx := solana.NewInstruction(feeQuoterProgramID, ix.Accounts(), ixData)
		for _, acc := range acceptOwnershipIx.Accounts() {
			if acc.PublicKey == timelockSigner {
				acc.IsSigner = false
			}
		}
		return acceptOwnershipIx, nil
	}

	tx, err := TransferAndWrapAcceptOwnership(
		buildTransfer,
		buildAccept,
		feeQuoterProgramID,
		proposedOwner,      // timelock PDA
		feeQuoterConfigPDA, // config PDA
		currentOwner,
		solChain,
		shared.FeeQuoter,
		timelockSigner, // the timelock signer PDA
	)

	if err != nil {
		return nil, fmt.Errorf("failed to transfer fee quoter ownership: %w", err)
	}

	result = append(result, tx)
	return result, nil
}

// transferOwnershipOffRamp transfers ownership of the offRamp to the timelock.
func transferOwnershipOffRamp(
	ccipState stateview.CCIPOnChainState,
	chainSelector uint64,
	solChain cldf_solana.Chain,
	currentOwner solana.PublicKey,
	proposedOwner solana.PublicKey,
	timelockSigner solana.PublicKey,
) ([]mcmsTypes.Transaction, error) {
	var result []mcmsTypes.Transaction

	state := ccipState.SolChains[chainSelector]

	// The relevant on-chain addresses
	offRampProgramID := state.OffRamp
	offRampConfigPDA := state.OffRampConfigPDA

	// Build specialized closures
	buildTransfer := func(proposedOwner, config, authority solana.PublicKey) (solana.Instruction, error) {
		ix, err := ccip_offramp.NewTransferOwnershipInstruction(
			proposedOwner, config, authority,
		).ValidateAndBuild()
		if err != nil {
			return nil, err
		}
		ixData, err := ix.Data()
		if err != nil {
			return nil, fmt.Errorf("failed to extract data payload from offramp transfer ownership instruction: %w", err)
		}
		transferOwnershipIx := solana.NewInstruction(offRampProgramID, ix.Accounts(), ixData)
		for _, acc := range transferOwnershipIx.Accounts() {
			if acc.PublicKey == timelockSigner {
				acc.IsSigner = false
			}
		}
		return transferOwnershipIx, nil
	}
	buildAccept := func(config, newOwnerAuthority solana.PublicKey) (solana.Instruction, error) {
		// If the router has its own accept function, use that
		ix, err := ccip_offramp.NewAcceptOwnershipInstruction(
			config, newOwnerAuthority,
		).ValidateAndBuild()
		if err != nil {
			return nil, err
		}
		ixData, err := ix.Data()
		if err != nil {
			return nil, fmt.Errorf("failed to extract data payload from offramp transfer ownership instruction: %w", err)
		}
		acceptOwnershipIx := solana.NewInstruction(offRampProgramID, ix.Accounts(), ixData)
		for _, acc := range acceptOwnershipIx.Accounts() {
			if acc.PublicKey == timelockSigner {
				acc.IsSigner = false
			}
		}
		return acceptOwnershipIx, nil
	}

	tx, err := TransferAndWrapAcceptOwnership(
		buildTransfer,
		buildAccept,
		offRampProgramID,
		proposedOwner,    // timelock PDA
		offRampConfigPDA, // config PDA
		currentOwner,
		solChain,
		shared.OffRamp,
		timelockSigner, // the timelock signer PDA
	)

	if err != nil {
		return nil, fmt.Errorf("failed to transfer offRamp ownership: %w", err)
	}

	result = append(result, tx)
	return result, nil
}

// transferOwnershipLockMintTokenPools transfers ownership of the lock mint token pools.
func transferOwnershipBurnMintTokenPools(
	ccipState stateview.CCIPOnChainState,
	tokenPoolConfigPDA solana.PublicKey,
	tokenMint solana.PublicKey,
	chainSelector uint64,
	solChain cldf_solana.Chain,
	tokenPoolMetadata string,
	currentOwner solana.PublicKey,
	proposedOwner solana.PublicKey,
	timelockSigner solana.PublicKey,
) ([]mcmsTypes.Transaction, error) {
	var result []mcmsTypes.Transaction

	state := ccipState.SolChains[chainSelector]

	// Build specialized closures
	buildTransfer := func(proposedOwner, config, authority solana.PublicKey) (solana.Instruction, error) {
		burnmint.SetProgramID(state.BurnMintTokenPools[tokenPoolMetadata])
		ix, err := burnmint.NewTransferOwnershipInstruction(
			proposedOwner, config, tokenMint, authority,
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
	buildAccept := func(config, newOwnerAuthority solana.PublicKey) (solana.Instruction, error) {
		burnmint.SetProgramID(state.BurnMintTokenPools[tokenPoolMetadata])
		// If the router has its own accept function, use that
		ix, err := burnmint.NewAcceptOwnershipInstruction(
			config, tokenMint, newOwnerAuthority,
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

	tx, err := TransferAndWrapAcceptOwnership(
		buildTransfer,
		buildAccept,
		state.BurnMintTokenPools[tokenPoolMetadata],
		proposedOwner,      // timelock PDA
		tokenPoolConfigPDA, // config PDA
		currentOwner,
		solChain,
		shared.BurnMintTokenPool,
		timelockSigner, // the timelock signer PDA
	)

	if err != nil {
		return nil, fmt.Errorf("failed to transfer burn-mint token pool ownership: %w", err)
	}

	result = append(result, tx)
	return result, nil
}

// transferOwnershipLockReleaseTokenPools transfers ownership of the lock mint token pools.
func transferOwnershipLockReleaseTokenPools(
	ccipState stateview.CCIPOnChainState,
	tokenPoolConfigPDA solana.PublicKey,
	tokenMint solana.PublicKey,
	chainSelector uint64,
	solChain cldf_solana.Chain,
	tokenPoolMetadata string,
	currentOwner solana.PublicKey,
	proposedOwner solana.PublicKey,
	timelockSigner solana.PublicKey,
) ([]mcmsTypes.Transaction, error) {
	var result []mcmsTypes.Transaction

	state := ccipState.SolChains[chainSelector]

	// Build specialized closures
	buildTransfer := func(proposedOwner, config, authority solana.PublicKey) (solana.Instruction, error) {
		lockrelease.SetProgramID(state.LockReleaseTokenPools[tokenPoolMetadata])
		ix, err := lockrelease.NewTransferOwnershipInstruction(
			proposedOwner, config, tokenMint, authority,
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
	buildAccept := func(config, newOwnerAuthority solana.PublicKey) (solana.Instruction, error) {
		lockrelease.SetProgramID(state.LockReleaseTokenPools[tokenPoolMetadata])
		// If the router has its own accept function, use that
		ix, err := lockrelease.NewAcceptOwnershipInstruction(
			config, tokenMint, newOwnerAuthority,
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

	tx, err := TransferAndWrapAcceptOwnership(
		buildTransfer,
		buildAccept,
		state.LockReleaseTokenPools[tokenPoolMetadata],
		proposedOwner,      // timelock PDA
		tokenPoolConfigPDA, // config PDA
		currentOwner,
		solChain,
		shared.LockReleaseTokenPool,
		timelockSigner, // the timelock signer PDA
	)

	if err != nil {
		return nil, fmt.Errorf("failed to transfer lock-release token pool ownership: %w", err)
	}

	result = append(result, tx)
	return result, nil
}

// transferOwnershipCCTPTokenPools transfers ownership of the CCTP token pool.
func transferOwnershipCCTPTokenPools(
	ccipState stateview.CCIPOnChainState,
	tokenPoolConfigPDA solana.PublicKey,
	tokenMint solana.PublicKey,
	chainSelector uint64,
	solChain cldf_solana.Chain,
	currentOwner solana.PublicKey,
	proposedOwner solana.PublicKey,
	timelockSigner solana.PublicKey,
) ([]mcmsTypes.Transaction, error) {
	var result []mcmsTypes.Transaction

	state := ccipState.SolChains[chainSelector]

	// Build specialized closures
	buildTransfer := func(proposedOwner, config, authority solana.PublicKey) (solana.Instruction, error) {
		cctp_token_pool.SetProgramID(state.CCTPTokenPool)
		ix, err := cctp_token_pool.NewTransferOwnershipInstruction(
			proposedOwner, config, tokenMint, authority,
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
	buildAccept := func(config, newOwnerAuthority solana.PublicKey) (solana.Instruction, error) {
		cctp_token_pool.SetProgramID(state.CCTPTokenPool)
		// If the router has its own accept function, use that
		ix, err := cctp_token_pool.NewAcceptOwnershipInstruction(
			config, tokenMint, newOwnerAuthority,
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

	tx, err := TransferAndWrapAcceptOwnership(
		buildTransfer,
		buildAccept,
		state.CCTPTokenPool,
		proposedOwner,      // timelock PDA
		tokenPoolConfigPDA, // config PDA
		currentOwner,
		solChain,
		shared.CCTPTokenPool,
		timelockSigner, // the timelock signer PDA
	)

	if err != nil {
		return nil, fmt.Errorf("failed to transfer CCTP token pool ownership: %w", err)
	}

	result = append(result, tx)
	return result, nil
}

// transferOwnershipRMNRemote transfers ownership of the RMNRemote to the timelock.
func transferOwnershipRMNRemote(
	ccipState stateview.CCIPOnChainState,
	chainSelector uint64,
	solChain cldf_solana.Chain,
	currentOwner solana.PublicKey,
	proposedOwner solana.PublicKey,
	timelockSigner solana.PublicKey,
) ([]mcmsTypes.Transaction, error) {
	var result []mcmsTypes.Transaction

	state := ccipState.SolChains[chainSelector]

	// The relevant on-chain addresses
	rmnRemoteProgramID := state.RMNRemote
	rmnRemoteConfigPDA := state.RMNRemoteConfigPDA
	rmnRemoteCursesPDA := state.RMNRemoteCursesPDA

	// Build specialized closures
	buildTransfer := func(newOwner, config, cursesConfig, authority solana.PublicKey) (solana.Instruction, error) {
		ix, err := rmn_remote.NewTransferOwnershipInstruction(
			newOwner, config, cursesConfig, authority,
		).ValidateAndBuild()
		if err != nil {
			return nil, err
		}
		ixData, err := ix.Data()
		if err != nil {
			return nil, fmt.Errorf("failed to extract data payload from rmn remote transfer ownership instruction: %w", err)
		}
		transferOwnershipIx := solana.NewInstruction(rmnRemoteProgramID, ix.Accounts(), ixData)
		for _, acc := range transferOwnershipIx.Accounts() {
			if acc.PublicKey == timelockSigner {
				acc.IsSigner = false
			}
		}
		return transferOwnershipIx, nil
	}
	buildAccept := func(config, newOwnerAuthority solana.PublicKey) (solana.Instruction, error) {
		// If the router has its own accept function, use that
		ix, err := rmn_remote.NewAcceptOwnershipInstruction(
			config, newOwnerAuthority,
		).ValidateAndBuild()
		if err != nil {
			return nil, err
		}
		ixData, err := ix.Data()
		if err != nil {
			return nil, fmt.Errorf("failed to extract data payload from rmn remote accept ownership instruction: %w", err)
		}
		acceptOwnershipIx := solana.NewInstruction(rmnRemoteProgramID, ix.Accounts(), ixData)
		for _, acc := range acceptOwnershipIx.Accounts() {
			if acc.PublicKey == timelockSigner {
				acc.IsSigner = false
			}
		}
		return acceptOwnershipIx, nil
	}

	programID := rmnRemoteProgramID
	configPDA := rmnRemoteConfigPDA
	label := shared.RMNRemote

	// We can't reuse the generic transferAndWrapAcceptOwnership function here
	// because the RMNRemote has an additional cursesConfig account that needs to be transferred.

	// 1. Build the instruction that transfers ownership to the timelock
	ixTransfer, err := buildTransfer(proposedOwner, configPDA, rmnRemoteCursesPDA, currentOwner)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create transfer ownership instruction: %w", label, err)
	}

	// if the old owner is the timelock signer, we need to build the accept instruction and submit it
	// We can't perform the accept ownership step here because the timelock signer is not a signer of the transaction
	if currentOwner.Equals(timelockSigner) {
		mcmsTx, err := BuildMCMSTxn(ixTransfer, programID.String(), label)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to create MCMS transaction: %w", label, err)
		}
		// we cannot accept ownership afterwards because the proposal needs to execute before accept is valid
		result = append(result, *mcmsTx)
		return result, nil
	}

	// 2. Confirm on-chain
	if err := solChain.Confirm([]solana.Instruction{ixTransfer}); err != nil {
		return nil, fmt.Errorf("%s: failed to confirm transfer on-chain: %w", label, err)
	}

	// 3. Build the “accept ownership” instruction
	ixAccept, err := buildAccept(configPDA, proposedOwner)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create accept ownership instruction: %w", label, err)
	}

	if proposedOwner.Equals(timelockSigner) {
		// 4. Wrap in MCMS transaction
		mcmsTx, err := BuildMCMSTxn(ixAccept, programID.String(), label)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to create MCMS transaction: %w", label, err)
		}
		result = append(result, *mcmsTx)
		return result, nil
	}

	// 4. Confirm on-chain
	if err := solChain.Confirm([]solana.Instruction{ixAccept}); err != nil {
		return nil, fmt.Errorf("%s: failed to confirm transfer on-chain: %w", label, err)
	}

	return result, nil
}
