package mcmsnew

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/deployment"
	commonChangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/internal"
	mcmEvm "github.com/smartcontractkit/chainlink/deployment/common/changeset/mcmsnew/evm"
	mcmSolana "github.com/smartcontractkit/chainlink/deployment/common/changeset/mcmsnew/solana"
	mcmsNewTypes "github.com/smartcontractkit/chainlink/deployment/common/changeset/mcmsnew/types"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
)

var _ deployment.ChangeSet[map[uint64]types.MCMSWithTimelockConfig] = DeployMCMSWithTimelock

func DeployMCMSWithTimelock(e deployment.Environment, cfgByChain map[uint64]types.MCMSWithTimelockConfig) (deployment.ChangesetOutput, error) {
	newAddresses := deployment.NewMemoryAddressBook()
	err := internal.DeployMCMSWithTimelockContractsBatch(
		e.Logger, e.Chains, newAddresses, cfgByChain,
	)
	if err != nil {
		return deployment.ChangesetOutput{AddressBook: newAddresses}, err
	}
	return deployment.ChangesetOutput{AddressBook: newAddresses}, nil
}

// DeployMCMSWithTimelockContractsBatch is a helper function to deploy MCMS contracts with timelock
// on multiple chains provided on the chains parameter.
func DeployMCMSWithTimelockContractsBatch(
	lggr logger.Logger,
	chains deployment.MultiFamilyChains,
	ab deployment.AddressBook,
	cfgByChain map[uint64]mcmsNewTypes.MCMSWithTimelockConfig,
) error {
	for chainSel, cfg := range cfgByChain {
		family, err := chain_selectors.GetSelectorFamily(chainSel)
		if err != nil {
			return err
		}
		switch family {
		case chain_selectors.FamilyEVM:
			_, err := mcmEvm.DeployMCMSWithTimelockContractsEVM(lggr, chains.EVMChains[chainSel], ab, cfg)
			if err != nil {
				return err
			}
		case chain_selectors.FamilySolana:
			_, err := mcmSolana.DeployMCMSWithTimelockProgramsSolana(lggr, chains.SolChains[chainSel], ab, cfg)
			if err != nil {
				return err
			}
		}

	}

	return nil
}

func ValidateOwnership(ctx context.Context, mcms bool, deployerKey, timelock common.Address, contract commonChangeset.Ownable) error {
	owner, err := contract.Owner(&bind.CallOpts{Context: ctx})
	if err != nil {
		return fmt.Errorf("failed to get owner: %w", err)
	}
	if mcms && owner != timelock {
		return fmt.Errorf("%s not owned by deployer key", contract.Address())
	} else if !mcms && owner != deployerKey {
		return fmt.Errorf("%s not owned by deployer key", contract.Address())
	}
	return nil
}

// TODO: SOLANA_CCIP
func ValidateOwnershipSolana(ctx context.Context, mcms bool, deployerKey, timelock, ccipRouter solana.PublicKey) error {
	return nil
}
