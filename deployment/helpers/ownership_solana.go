package helpers

import (
	"fmt"

	"github.com/gagliardetto/solana-go"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	mcmsTypes "github.com/smartcontractkit/mcms/types"
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
	solChain cldf.SolChain, // used for solChain.Confirm
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
