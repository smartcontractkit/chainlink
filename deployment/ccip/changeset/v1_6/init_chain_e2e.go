package v1_6

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
)

var (
	_ deployment.ChangeSet[InitChaine2eConfig] = InitChaine2eChangeset
)

type InitChaine2eConfig struct {
	HomeChainSelector uint64
	// RMNRemoteConfigs       map[uint64]RMNRemoteConfig
	PreReqConfig           changeset.DeployPrerequisiteConfig
	ContractParamsPerChain map[uint64]ChainContractParams
}

func InitChaine2eChangeset(env deployment.Environment, cfg InitChaine2eConfig) (deployment.ChangesetOutput, error) {
	// TODO add custom validation

	// Searches through the addressbook and check if that contract exist onchain (type and version)
	// We need a way to pass in previous deployed addressBook here
	// State only returns contracts from addressbook

	// state, err := changeset.LoadOnchainState(env)
	// if err != nil {
	// 	return deployment.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	// }

	// fmt.Println("state: ", state)

	addressBook := deployment.NewMemoryAddressBook()
	// batches := make([]mcmstypes.BatchOperation, 0)

	// correct ordering for new chain integration
	// err := internal.DeployMCMSWithTimelockContractsBatch(
	// 	e.Logger, e.Chains, newAddresses, cfgByChain,
	// )
	// if err != nil {
	// 	return deployment.ChangesetOutput{AddressBook: newAddresses}, err
	// }

	err := changeset.DeployPrerequisiteChainContracts(env, addressBook, cfg.PreReqConfig)
	if err != nil {
		env.Logger.Errorw("Failed to deploy prerequisite contracts", "err", err, "addressBook", addressBook)
		return deployment.ChangesetOutput{
			AddressBook: addressBook,
		}, fmt.Errorf("failed to deploy prerequisite contracts: %w", err)
	}

	fmt.Println("ADDRESSBOOK0: ", addressBook)
	err = MergeAddress(env, env.ExistingAddresses, addressBook)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	// Note: simpler to declare new memory addressbook for each changeset rather than filtering address from addressbook to merge
	addressBook1 := deployment.NewMemoryAddressBook()
	err = deployChainContractsForChains(env, addressBook1, cfg.HomeChainSelector, cfg.ContractParamsPerChain)
	if err != nil {
		env.Logger.Errorw("Failed to deploy CCIP contracts", "err", err, "newAddresses", addressBook1)
		return deployment.ChangesetOutput{AddressBook: addressBook1}, deployment.MaybeDataErr(err)
	}

	fmt.Println("ADDRESSBOOK1: ", addressBook1)
	err = MergeAddress(env, env.ExistingAddresses, addressBook1)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

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
	return deployment.ChangesetOutput{}, nil
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
