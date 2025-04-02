package sequence

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/mcms"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset"
	feemanager "github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/fee-manager"
	rewardmanager "github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/reward-manager"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/verification"
	dsutil "github.com/smartcontractkit/chainlink/deployment/data-streams/utils"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/mcmsutil"
)

// DeployDataStreamsChangeset deploys the entire data streams stack to a new chain. It should be kept up to date
// with the latest contract versions and deployment logic.
var DeployDataStreamsChangeset = deployment.CreateChangeSet(deployDataStreamsLogic, deployDataStreamsPrecondition)

type DeployDataStreamsConfig struct {
	ChainsToDeploy map[uint64]DeployDataStreams
}

type BillingConfig struct {
	LinkTokenAddress   common.Address
	NativeTokenAddress common.Address
}

type BillingFeature struct {
	Enabled bool
	Config  *BillingConfig
}

type DeployDataStreams struct {
	VerifierConfig verification.SetConfig

	Billing   BillingFeature
	Ownership types.OwnershipFeature
}

func deployDataStreamsLogic(e deployment.Environment, cc DeployDataStreamsConfig) (deployment.ChangesetOutput, error) {
	ab := deployment.NewMemoryAddressBook() // changeset output expects only new addresses
	var timelockProposals []mcms.TimelockProposal
	existingAddresses := e.ExistingAddresses
	for chain, cfg := range cc.ChainsToDeploy {
		// Deploy MCMS
		if cfg.Ownership.DeployMCMS {
			mcmsDeployCfg := changeset.DeployMCMSConfig{
				ChainsToDeploy: []uint64{chain},
				Ownership:      cfg.Ownership.AsSettings(),
				Config:         cfg.Ownership.DeployMCMSConfig,
			}
			mcmsDeployOut, err := changeset.DeployAndTransferMCMSChangeset.Apply(e, mcmsDeployCfg)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to deploy MCMS on chain %d: %w", chain, err)
			}
			if err := ab.Merge(mcmsDeployOut.AddressBook); err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("address book merge failed after MCMS deployment: %w", err)
			}
			timelockProposals = append(timelockProposals, mcmsDeployOut.MCMSTimelockProposals...)
			addressesWithMCMS := deployment.NewMemoryAddressBook()
			if err := addressesWithMCMS.Merge(ab); err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to merge addressesWithMCMS with ab: %w", err)
			}
			if err := addressesWithMCMS.Merge(e.ExistingAddresses); err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to merge existing addresses with existing: %w", err)
			}
			e.ExistingAddresses = addressesWithMCMS
		}

		// Deploy Verifier Proxy
		verifierProxyCfg := verification.DeployVerifierProxyConfig{
			ChainsToDeploy: map[uint64]verification.DeployVerifierProxy{
				chain: {}, // Implement AccessController as needed
			},
			Version:   deployment.Version0_5_0,
			Ownership: cfg.Ownership.AsSettings(),
		}
		proxyOut, err := verification.DeployVerifierProxyChangeset.Apply(e, verifierProxyCfg)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to deploy verifier proxy on chain: %d err %w", chain, err)
		}

		if err := ab.Merge(proxyOut.AddressBook); err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("address book merge failed after verifier proxy deployment: %w", err)
		}

		verifierProxyAddr, err := dsutil.MaybeFindEthAddress(ab, chain, types.VerifierProxy)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to find verifier proxy address: %s", err)
		}
		timelockProposals = append(timelockProposals, proxyOut.MCMSTimelockProposals...)

		// Deploy Verifier
		verifierCfg := verification.DeployVerifierConfig{
			ChainsToDeploy: map[uint64]verification.DeployVerifier{
				chain: {VerifierProxyAddress: verifierProxyAddr},
			},
			Ownership: cfg.Ownership.AsSettings(),
		}

		verifierOut, err := verification.DeployVerifierChangeset.Apply(e, verifierCfg)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to deploy verifier on chain %d: %w", chain, err)
		}
		if err := ab.Merge(verifierOut.AddressBook); err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("address book merge failed after verifier deployment: %w", err)
		}

		verifierAddr, err := dsutil.MaybeFindEthAddress(ab, chain, types.Verifier)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to find verifier address: %s", err)
		}

		timelockProposals = append(timelockProposals, verifierOut.MCMSTimelockProposals...)

		// Initialize Verifier on VerifierProxy
		initVerifierCfg := verification.VerifierProxyInitializeVerifierConfig{
			ConfigPerChain: map[uint64][]verification.InitializeVerifierConfig{
				chain: {{
					ContractAddress: verifierProxyAddr,
					VerifierAddress: verifierAddr,
				}},
			},
		}

		_, err = verification.InitializeVerifierChangeset.Apply(e, initVerifierCfg)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to initialize verifier on chain %d: %w", chain, err)
		}

		// SetConfig
		setCfg := verification.SetConfigConfig{
			ConfigsByChain: map[uint64][]verification.SetConfig{
				chain: {verification.SetConfig{
					VerifierAddress:            verifierAddr,
					ConfigDigest:               cfg.VerifierConfig.ConfigDigest,
					Signers:                    cfg.VerifierConfig.Signers,
					F:                          cfg.VerifierConfig.F,
					RecipientAddressesAndProps: cfg.VerifierConfig.RecipientAddressesAndProps,
				}},
			},
		}

		_, err = verification.SetConfigChangeset.Apply(e, setCfg)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to set config on chain %d: %w", chain, err)
		}

		// Deploy FeeManager & Reward Manager (optional)
		if cfg.Billing.Enabled && cfg.Billing.Config != nil {
			// Deploy RewardManager
			rewardMgrCfg := rewardmanager.DeployRewardManagerConfig{
				ChainsToDeploy: map[uint64]rewardmanager.DeployRewardManager{
					chain: {LinkTokenAddress: cfg.Billing.Config.LinkTokenAddress},
				},
				Ownership: cfg.Ownership.AsSettings(),
			}

			rmOut, err := rewardmanager.DeployRewardManagerChangeset.Apply(e, rewardMgrCfg)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to deploy reward manager on chain %d: %w", chain, err)
			}
			if err := ab.Merge(rmOut.AddressBook); err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("address book merge failed after reward manager deployment: %w", err)
			}

			rewardMgrAddr, err := dsutil.MaybeFindEthAddress(ab, chain, types.RewardManager)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to find reward manager address: %s", err)
			}

			timelockProposals = append(timelockProposals, rmOut.MCMSTimelockProposals...)

			// Deploy FeeManager
			feeMgrCfg := feemanager.DeployFeeManagerConfig{
				ChainsToDeploy: map[uint64]feemanager.DeployFeeManager{
					chain: {
						LinkTokenAddress:     cfg.Billing.Config.LinkTokenAddress,
						NativeTokenAddress:   cfg.Billing.Config.NativeTokenAddress,
						VerifierProxyAddress: verifierProxyAddr,
						RewardManagerAddress: rewardMgrAddr,
					},
				},
				Ownership: cfg.Ownership.AsSettings(),
			}

			fmOut, err := feemanager.DeployFeeManagerChangeset.Apply(e, feeMgrCfg)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to deploy fee manager on chain %d: %w", chain, err)
			}
			if err := ab.Merge(fmOut.AddressBook); err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("address book merge failed after fee manager deployment: %w", err)
			}
			if _, err := dsutil.MaybeFindEthAddress(ab, chain, types.FeeManager); err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("fee manager address not found for chain %d: %w", chain, err)
			}

			timelockProposals = append(timelockProposals, fmOut.MCMSTimelockProposals...)
			// reset the address book to the original state
			e.ExistingAddresses = existingAddresses
		}
	}

	mergedTimelockProposal, err := mcmsutil.MergeSimilarTimelockProposals(timelockProposals)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to merge timelock proposals: %w", err)
	}

	return deployment.ChangesetOutput{
		AddressBook:           ab,
		MCMSTimelockProposals: []mcms.TimelockProposal{mergedTimelockProposal},
	}, nil
}

func deployDataStreamsPrecondition(_ deployment.Environment, cc DeployDataStreamsConfig) error {
	if err := cc.Validate(); err != nil {
		return fmt.Errorf("invalid DeployDataStreams config: %w", err)
	}
	return nil
}

func (cc DeployDataStreamsConfig) Validate() error {
	if len(cc.ChainsToDeploy) == 0 {
		return errors.New("ChainsToDeploy is empty")
	}
	for chain := range cc.ChainsToDeploy {
		if err := deployment.IsValidChainSelector(chain); err != nil {
			return fmt.Errorf("invalid chain selector: %d - %w", chain, err)
		}
	}
	return nil
}
