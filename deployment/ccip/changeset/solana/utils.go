package solana

import (
	"fmt"
	"math/big"

	"github.com/gagliardetto/solana-go"
	mcmsSolana "github.com/smartcontractkit/mcms/sdk/solana"
	mcmsTypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink/deployment"
	cs "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
)

func ValidateMCMSConfig(e deployment.Environment, chainSelector uint64, mcms *cs.MCMSConfig) error {
	if mcms != nil {
		// If there is no timelock and mcms proposer on the chain, the transfer will fail.
		timelockID, err := deployment.SearchAddressBook(e.ExistingAddresses, chainSelector, types.RBACTimelock)
		if err != nil {
			return fmt.Errorf("timelock not present on the chain %w", err)
		}
		proposerID, err := deployment.SearchAddressBook(e.ExistingAddresses, chainSelector, types.ProposerManyChainMultisig)
		if err != nil {
			return fmt.Errorf("mcms proposer not present on the chain %w", err)
		}
		// Make sure addresses are correctly parsed. Format is: "programID.PDASeed"
		_, _, err = mcmsSolana.ParseContractAddress(timelockID)
		if err != nil {
			return fmt.Errorf("failed to parse timelock address: %w", err)
		}
		_, _, err = mcmsSolana.ParseContractAddress(proposerID)
		if err != nil {
			return fmt.Errorf("failed to parse proposer address: %w", err)
		}
	}
	return nil
}

func BuildMCMSTxn(ixn solana.Instruction, programID string, contractType deployment.ContractType) (*mcmsTypes.Transaction, error) {
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
