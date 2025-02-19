package solana

import (
	"fmt"

	mcmsSolana "github.com/smartcontractkit/mcms/sdk/solana"

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
