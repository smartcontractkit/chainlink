package v1_6

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
	mcmslib "github.com/smartcontractkit/mcms"
	mcmssdk "github.com/smartcontractkit/mcms/sdk"
	mcmstypes "github.com/smartcontractkit/mcms/types"
)

var (
	_ deployment.ChangeSet[InitChaine2eConfig] = InitChaine2eChangeset
)

type InitChaine2eConfig struct {
	HomeChainSelector uint64
	McmsConfig        map[uint64]types.MCMSWithTimelockConfigV2
	// RMNRemoteConfigs       map[uint64]RMNRemoteConfig
	PreReqConfig             changeset.DeployPrerequisiteConfig
	ContractParamsPerChain   DeployChainContractsConfig
	UpdateChainConfig        UpdateChainConfigConfig
	AddDonSetCandidateConfig AddDonAndSetCandidateChangesetConfig
	SetCandidateConfig       SetCandidateChangesetConfig
	PromoteCandidateConfig   PromoteCandidateChangesetConfig
	Ocr3Config               SetOCR3OffRampConfig
}

func InitChaine2eChangeset(env deployment.Environment, cfg InitChaine2eConfig) (deployment.ChangesetOutput, error) {
	// TODO add custom validation

	// Searches through the addressbook and check if that contract exist onchain (type and version)
	// We need a way to pass in previous deployed addressBook here
	// State only returns contracts from addressbook and check it's online state

	state, err := changeset.LoadOnchainState(env)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	addressBook := deployment.NewMemoryAddressBook()
	// batches := make([]mcmstypes.BatchOperation, 0)

	// NEW APPROACH
	// fmt.Println("ENV ADDRESS 0: ", env.ExistingAddresses)

	// output, err := changeset.DeployPrerequisitesChangeset(env, cfg.PreReqConfig)
	// if err != nil {
	// 	return deployment.ChangesetOutput{}, fmt.Errorf("Error running DeployPrerequisiteChainContracts: ", err)
	// }

	// fmt.Println("ENV ADDRESS 1: ", env.ExistingAddresses)
	// err = env.ExistingAddresses.Merge(output.AddressBook)
	// if err != nil {
	// 	return deployment.ChangesetOutput{}, err
	// }

	// err = addressBook.Merge(output.AddressBook)
	// if err != nil {
	// 	return deployment.ChangesetOutput{}, err
	// }

	// fmt.Println("ENV ADDRESS 2: ", env.ExistingAddresses)
	// output, err = DeployChainContractsChangeset(env, cfg.ContractParamsPerChain)
	// if err != nil {
	// 	return deployment.ChangesetOutput{}, fmt.Errorf("Error running DeployChainContractsChangeset: ", err)
	// }
	// err = env.ExistingAddresses.Merge(output.AddressBook)
	// if err != nil {
	// 	return deployment.ChangesetOutput{}, err
	// }

	// err = addressBook.Merge(output.AddressBook)
	// if err != nil {
	// 	return deployment.ChangesetOutput{}, err
	// }

	// fmt.Println("ENV ADDRESS 3: ", env.ExistingAddresses)
	// Generate an MCMs proposal
	// output, err := UpdateChainConfigChangeset(env, cfg.UpdateChainConfig)
	// if err != nil {
	// 	return deployment.ChangesetOutput{}, fmt.Errorf("error running UpdateChainConfigChangeset: %w", err)
	// }

	// fmt.Println("MCMS output UpdateChainConfigChangeset: ", output)
	// TODO handle MCMS proposals

	// Aggregate all proposals.
	var batches []mcmstypes.BatchOperation

	output, err := AddDonAndSetCandidateChangeset(env, cfg.AddDonSetCandidateConfig)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("error running AddDonAndSetCandidateChangeset: %w", err)
	}

	batches, err = BatchProposal(output, batches)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	output, err = SetCandidateChangeset(env, cfg.SetCandidateConfig)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("error running SetCandidateChangeset: %w", err)
	}

	batches, err = BatchProposal(output, batches)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	output, err = PromoteCandidateChangeset(env, cfg.PromoteCandidateConfig)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("error running PromoteCandidateChangeset: %w", err)
	}

	batches, err = BatchProposal(output, batches)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	output, err = SetOCR3OffRampChangeset(env, cfg.Ocr3Config)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("error running SetOCR3OffRampChangeset: %w", err)
	}

	batches, err = BatchProposal(output, batches)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	// Store the timelocks, proposers, and inspectors for each chain.
	timelocks := make(map[uint64]string)
	proposers := make(map[uint64]string)
	inspectors := make(map[uint64]mcmssdk.Inspector)
	for _, op := range batches {
		chainSel := uint64(op.ChainSelector)
		timelocks[chainSel] = state.Chains[chainSel].Timelock.Address().Hex()
		proposers[chainSel] = state.Chains[chainSel].ProposerMcm.Address().Hex()
		inspectors[chainSel], err = proposalutils.McmsInspectorForChain(env, chainSel)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to get MCMS inspector for chain with selector %d: %w", chainSel, err)
		}
	}

	proposal, err := proposalutils.BuildProposalFromBatchesV2(
		env,
		timelocks,
		proposers,
		inspectors,
		batches,
		fmt.Sprintf("Running OCR3Config + setCandidate, PromoteCandidate on HomeChain"),
		time.Duration(0),
	)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
	}

	fmt.Println("Full batch proposal: ", proposal)

	return deployment.ChangesetOutput{
		MCMSTimelockProposals: []mcmslib.TimelockProposal{*proposal},
		AddressBook:           addressBook,
	}, nil
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

func BatchProposal(output deployment.ChangesetOutput, batches []mcmstypes.BatchOperation) ([]mcmstypes.BatchOperation, error) {
	for _, proposal := range output.MCMSTimelockProposals {
		for _, p := range proposal.Operations {
			for _, batchTx := range p.Transactions {
				batchOperation, err := proposalutils.BatchOperationForChain(
					uint64(p.ChainSelector),
					batchTx.To,
					batchTx.Data,
					big.NewInt(0),
					batchTx.ContractType,
					batchTx.Tags,
				)
				if err != nil {
					return nil, fmt.Errorf("failed to create batch operation on chain with selector %d: %w", p.ChainSelector, err)
				}
				batches = append(batches, batchOperation)
			}
		}
	}

	return batches, nil
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
