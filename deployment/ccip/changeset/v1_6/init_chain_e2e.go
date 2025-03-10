package v1_6

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
)

var (
	_ deployment.ChangeSet[InitChaine2eConfig] = InitChaine2eChangeset
)

type InitChaine2eConfig struct {
	HomeChainSelector uint64
	McmsConfig        map[uint64]types.MCMSWithTimelockConfigV2
	// RMNRemoteConfigs       map[uint64]RMNRemoteConfig
	PreReqConfig           changeset.DeployPrerequisiteConfig
	ContractParamsPerChain DeployChainContractsConfig
	UpdateChainConfig      UpdateChainConfigConfig
}

func InitChaine2eChangeset(env deployment.Environment, cfg InitChaine2eConfig) (deployment.ChangesetOutput, error) {
	// TODO add custom validation

	// Searches through the addressbook and check if that contract exist onchain (type and version)
	// We need a way to pass in previous deployed addressBook here
	// State only returns contracts from addressbook and check it's online state
	// state, err := changeset.LoadOnchainState(env)

	addressBook := deployment.NewMemoryAddressBook()
	// batches := make([]mcmstypes.BatchOperation, 0)

	// Correct ordering for new chain integration
	// err := commonchangeset.DeployInternalMCMSWithTimelockV2ForEVM(env, env.Logger, addressBook, cfg.McmsConfig)
	// if err != nil {
	// 	return deployment.ChangesetOutput{AddressBook: addressBook}, err
	// }

	// err = MergeAddress(env, env.ExistingAddresses, addressBook)
	// if err != nil {
	// 	return deployment.ChangesetOutput{}, err
	// }

	// addressBook1 := deployment.NewMemoryAddressBook()
	// err = changeset.DeployPrerequisiteChainContracts(env, addressBook1, cfg.PreReqConfig)
	// if err != nil {
	// 	env.Logger.Errorw("Failed to deploy prerequisite contracts", "err", err, "addressBook", addressBook1)
	// 	return deployment.ChangesetOutput{
	// 		AddressBook: addressBook1,
	// 	}, fmt.Errorf("failed to deploy prerequisite contracts: %w", err)
	// }

	// err = MergeAddress(env, env.ExistingAddresses, addressBook1)
	// if err != nil {
	// 	return deployment.ChangesetOutput{}, err
	// }

	// // Note: simpler to declare new memory addressbook for each changeset rather than filtering address from addressbook to merge
	// addressBook2 := deployment.NewMemoryAddressBook()
	// err = deployChainContractsForChains(env, addressBook2, cfg.HomeChainSelector, cfg.ContractParamsPerChain)
	// if err != nil {
	// 	env.Logger.Errorw("Failed to deploy CCIP contracts", "err", err, "newAddresses", addressBook2)
	// 	return deployment.ChangesetOutput{AddressBook: addressBook2}, deployment.MaybeDataErr(err)
	// }

	// err = MergeAddress(env, env.ExistingAddresses, addressBook2)
	// if err != nil {
	// 	return deployment.ChangesetOutput{}, err
	// }
	fmt.Println("ENV ADDRESS 0: ", env.ExistingAddresses)

	output, err := changeset.DeployPrerequisitesChangeset(env, cfg.PreReqConfig)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("Error running DeployPrerequisiteChainContracts: ", err)
	}

	fmt.Println("ENV ADDRESS 1: ", env.ExistingAddresses)
	err = env.ExistingAddresses.Merge(output.AddressBook)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	err = addressBook.Merge(output.AddressBook)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	fmt.Println("ENV ADDRESS 2: ", env.ExistingAddresses)
	output, err = DeployChainContractsChangeset(env, cfg.ContractParamsPerChain)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("Error running DeployChainContractsChangeset: ", err)
	}
	err = env.ExistingAddresses.Merge(output.AddressBook)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	err = addressBook.Merge(output.AddressBook)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	fmt.Println("ENV ADDRESS 3: ", env.ExistingAddresses)
	// Generate an MCMs proposal
	// output, err := UpdateChainConfigChangeset(env, cfg.UpdateChainConfig)
	// if err != nil {
	// 	return deployment.ChangesetOutput{}, fmt.Errorf("Error running UpdateChainConfigChangeset: ", err)
	// }

	// fmt.Println(output)
	// TODO handle MCMS proposals

	return deployment.ChangesetOutput{
		Proposals:   []timelock.MCMSWithTimelockProposal{},
		AddressBook: addressBook}, nil
}

func MergeAddress(env deployment.Environment, existingAddressBook, newAddresses deployment.AddressBook) error {
	if err := existingAddressBook.Merge(newAddresses); err != nil {
		return err
	}

	addrs, err := existingAddressBook.Addresses()
	if err != nil {
		return err
	}

	// TODO: donot hardcode this value
	// TODO: can we trigger this from chainlink-deployments repo, one repo changing another repos file seems weird
	// TODO: Also not sure how this is going to run on pipeline
	WriteFile("/Users/stackman/Desktop/chainlink-deployments/domains/ccip/staging-bix/addresses.json", addrs)

	return nil
}

func WriteFile(path string, data any) error {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, b, 0644)
}

// Order of operations
// remote, ok := rmnRemotePerChain[chain]
// if !ok {
// 	return deployment.ChangesetOutput{}, fmt.Errorf("RMNRemote contract not found for chain %d", chain)
// }

// deployer := getDeployer(e, chain, config.MCMSConfig)
// tx, err := remote.SetConfig(deployer, newConfig)

// batchOperation, err := setRMNRemoteOnRMNProxyOp(txOpts, chain, state.Chains[sel], cfg.MCMSConfig != nil)
// if err != nil {
// 	return deployment.ChangesetOutput{}, fmt.Errorf("failed to set RMNRemote on RMNProxy for chain %s: %w", chain.String(), err)
// }

// tx, err := state.Chains[cfg.HomeChainSelector].CCIPHome.ApplyChainConfigUpdates(txOpts, cfg.RemoteChainRemoves, adds)
// if cfg.MCMS == nil {
// 	_, err = deployment.ConfirmIfNoErrorWithABI(e.Chains[cfg.HomeChainSelector], tx, ccip_home.CCIPHomeABI, err)
// 	if err != nil {
// 		return deployment.ChangesetOutput{}, err
// 	}
// 	e.Logger.Infof("Updated chain config on chain %d removes %v, adds %v", cfg.HomeChainSelector, cfg.RemoteChainRemoves, cfg.RemoteChainAdds)
// 	return deployment.ChangesetOutput{}, nil
// }

// AddDonAndSetCandidateChangeset(
// 	e deployment.Environment,
// 	cfg AddDonAndSetCandidateChangesetConfig,
// ) (deployment.ChangesetOutput, error) {}

// SetCandidateChangeset(
// 	e deployment.Environment,
// 	cfg SetCandidateChangesetConfig,
// ) (deployment.ChangesetOutput, error) {}

// PromoteCandidateChangeset(
// 		e deployment.Environment,
// 		cfg PromoteCandidateChangesetConfig,
// ) (deployment.ChangesetOutput, error) {}

// SetOCR3OffRampChangeset(e deployment.Environment, cfg SetOCR3OffRampConfig) (deployment.ChangesetOutput, error) { }
