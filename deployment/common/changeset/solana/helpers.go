package solana

import (
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
)

// FundFromAddressIxs transfers SOL from the given address to each provided account and waits for confirmations.
func FundFromAddressIxs(solChain cldf_solana.Chain, from solana.PublicKey, accounts []solana.PublicKey, amount uint64) ([]solana.Instruction, error) {
	var ixs []solana.Instruction
	for _, account := range accounts {
		// Create a transfer instruction using the provided builder.
		ix, err := system.NewTransferInstruction(
			amount,
			from,    // funding account (sender)
			account, // recipient account
		).ValidateAndBuild()
		if err != nil {
			return nil, fmt.Errorf("failed to create transfer instruction: %w", err)
		}
		ixs = append(ixs, ix)
	}

	return ixs, nil
}

// FundFromDeployerKey transfers SOL from the deployer to each provided account and waits for confirmations.
func FundFromDeployerKey(solChain cldf_solana.Chain, accounts []solana.PublicKey, amount uint64) error {
	ixs, err := FundFromAddressIxs(solChain, solChain.DeployerKey.PublicKey(), accounts, amount)
	if err != nil {
		return fmt.Errorf("failed to create transfer instructions: %w", err)
	}
	err = solChain.Confirm(ixs)
	if err != nil {
		return fmt.Errorf("failed to confirm transaction: %w", err)
	}
	return nil
}
