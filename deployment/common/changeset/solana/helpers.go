package solana

import (
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"

	"github.com/smartcontractkit/chainlink/deployment"
)

// FundFromDeployerKey transfers SOL from the deployer to each provided account and waits for confirmations.
func FundFromDeployerKey(solChain deployment.SolChain, accounts []solana.PublicKey, amount uint64) error {
	var ixs []solana.Instruction
	for _, account := range accounts {
		// Create a transfer instruction using the provided builder.
		ix, err := system.NewTransferInstruction(
			amount,
			solChain.DeployerKey.PublicKey(), // funding account (sender)
			account,                          // recipient account
		).ValidateAndBuild()
		if err != nil {
			return fmt.Errorf("failed to create transfer instruction: %w", err)
		}
		ixs = append(ixs, ix)
	}

	err := solChain.Confirm(ixs)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}
	return nil
}
