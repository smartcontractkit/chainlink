package changeset

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/internal"
	evminternal "github.com/smartcontractkit/chainlink/deployment/common/changeset/internal/evm"
	solanainternal "github.com/smartcontractkit/chainlink/deployment/common/changeset/internal/solana"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
)

var (
	_ deployment.ChangeSet[map[uint64]types.MCMSWithTimelockConfig]   = DeployMCMSWithTimelock
	_ deployment.ChangeSet[map[uint64]types.MCMSWithTimelockConfigV2] = DeployMCMSWithTimelockV2
)

// DeployMCMSWithTimelock deploys and initializes the MCM and Timelock contracts
// Deprecated: use DeployMCMSWithTimelockV2 instead
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

// DeployMCMSWithTimelockV2 deploys and initializes the MCM and Timelock contracts
func DeployMCMSWithTimelockV2(
	env deployment.Environment, cfgByChain map[uint64]types.MCMSWithTimelockConfigV2,
) (deployment.ChangesetOutput, error) {
	newAddresses := deployment.NewMemoryAddressBook()

	for chainSel, cfg := range cfgByChain {
		family, err := chain_selectors.GetSelectorFamily(chainSel)
		if err != nil {
			return deployment.ChangesetOutput{AddressBook: newAddresses}, err
		}

		switch family {
		case chain_selectors.FamilyEVM:
			_, err := evminternal.DeployMCMSWithTimelockContractsEVM(env.Logger, env.Chains[chainSel], newAddresses, cfg)
			if err != nil {
				return deployment.ChangesetOutput{AddressBook: newAddresses}, err
			}

		case chain_selectors.FamilySolana:
			_, err := solanainternal.DeployMCMSWithTimelockProgramsSolana(env, env.SolChains[chainSel], newAddresses, cfg)
			if err != nil {
				return deployment.ChangesetOutput{AddressBook: newAddresses}, err
			}

		default:
			err = fmt.Errorf("unsupported chain family: %s", family)
			return deployment.ChangesetOutput{AddressBook: newAddresses}, err
		}
	}

	return deployment.ChangesetOutput{AddressBook: newAddresses}, nil
}

func ValidateOwnership(ctx context.Context, mcms bool, deployerKey, timelock common.Address, contract Ownable) error {
	owner, err := contract.Owner(&bind.CallOpts{Context: ctx})
	if err != nil {
		return fmt.Errorf("failed to get owner: %w", err)
	}
	if mcms && owner != timelock {
		return fmt.Errorf("%s not owned by timelock", contract.Address())
	} else if !mcms && owner != deployerKey {
		return fmt.Errorf("%s not owned by deployer key", contract.Address())
	}
	return nil
}

func ValidateOwnershipSolanaCommon(mcms bool, deployerKey solana.PublicKey, timelockSignerPDA solana.PublicKey, programOwner solana.PublicKey) error {
	if !mcms {
		if deployerKey.String() != programOwner.String() {
			return fmt.Errorf("deployer key %s does not match owner %s", deployerKey.String(), programOwner.String())
		}
	} else {
		if timelockSignerPDA.String() != programOwner.String() {
			return fmt.Errorf("timelock signer PDA %s does not match owner %s", timelockSignerPDA.String(), programOwner.String())
		}
	}
	return nil
}
