package sequence

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/mcms"

	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/mcmsutil"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset"
	feemanager "github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/fee-manager"
	rewardmanager "github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/reward-manager"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/verification"
	dsutil "github.com/smartcontractkit/chainlink/deployment/data-streams/utils"
)

// DeployDataStreamsDestinationChainChangeset deploys the entire data streams destination chain contracts. It should be kept up to date
// with the latest contract versions and deployment logic.
var DeployDataStreamsDestinationChainChangeset = deployment.CreateChangeSet(deployDataStreamsLogic, deployDataStreamsPrecondition)

type DeployDataStreamsConfig struct {
	ChainsToDeploy map[uint64]DeployDataStreams
}

type BillingConfig struct {
	LinkTokenAddress   common.Address
	NativeTokenAddress common.Address
	Surcharge          uint64
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
	newAddresses := deployment.NewMemoryAddressBook() // changeset output expects only new addresses
	// Clone env. to avoid mutation the changeset Applier expects no changes.
	existingAddresses, err := e.ExistingAddresses.Addresses()
	abClone := deployment.NewMemoryAddressBookFromMap(existingAddresses)
	cloneEnv := deployment.Environment{
		Name:              e.Name,
		Logger:            e.Logger,
		ExistingAddresses: abClone,
		Chains:            e.Chains,
		SolChains:         e.SolChains,
		NodeIDs:           e.NodeIDs,
		Offchain:          e.Offchain,
		OCRSecrets:        e.OCRSecrets,
		GetContext:        e.GetContext,
	}

	var timelockProposals []mcms.TimelockProposal

	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get existing addresses: %w", err)
	}
	for chain, cfg := range cc.ChainsToDeploy {
		// Deploy MCMS
		if cfg.Ownership.DeployMCMS {
			mcmsDeployCfg := changeset.DeployMCMSConfig{
				ChainsToDeploy: []uint64{chain},
				Ownership:      cfg.Ownership.AsSettings(),
				Config:         cfg.Ownership.DeployMCMSConfig,
			}
			mcmsDeployOut, err := changeset.DeployAndTransferMCMSChangeset.Apply(cloneEnv, mcmsDeployCfg)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to deploy MCMS on chain %d: %w", chain, err)
			}
			if err := newAddresses.Merge(mcmsDeployOut.AddressBook); err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("address book merge failed after MCMS deployment: %w", err)
			}
			timelockProposals = append(timelockProposals, mcmsDeployOut.MCMSTimelockProposals...)
			if err := cloneEnv.ExistingAddresses.Merge(mcmsDeployOut.AddressBook); err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to merge existing addresses with mcms addresses: %w", err)
			}
		}

		// Deploy Verifier Proxy
		verifierProxyCfg := verification.DeployVerifierProxyConfig{
			ChainsToDeploy: map[uint64]verification.DeployVerifierProxy{
				chain: {}, // Implement AccessController as needed
			},
			Version:   deployment.Version0_5_0,
			Ownership: cfg.Ownership.AsSettings(),
		}
		proxyOut, err := verification.DeployVerifierProxyChangeset.Apply(cloneEnv, verifierProxyCfg)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to deploy verifier proxy on chain: %d err %w", chain, err)
		}
		timelockProposals = append(timelockProposals, proxyOut.MCMSTimelockProposals...)
		if err := newAddresses.Merge(proxyOut.AddressBook); err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("address book merge failed after verifier proxy deployment: %w", err)
		}

		verifierProxyAddr, err := dsutil.MaybeFindEthAddress(newAddresses, chain, types.VerifierProxy)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to find verifier proxy address: %w", err)
		}

		if err := cloneEnv.ExistingAddresses.Merge(proxyOut.AddressBook); err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to merge verifier proxy address: %w", err)
		}

		// Deploy Verifier
		verifierCfg := verification.DeployVerifierConfig{
			ChainsToDeploy: map[uint64]verification.DeployVerifier{
				chain: {VerifierProxyAddress: verifierProxyAddr},
			},
			Ownership: cfg.Ownership.AsSettings(),
		}

		verifierOut, err := verification.DeployVerifierChangeset.Apply(cloneEnv, verifierCfg)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to deploy verifier on chain %d: %w", chain, err)
		}
		timelockProposals = append(timelockProposals, verifierOut.MCMSTimelockProposals...)
		if err := newAddresses.Merge(verifierOut.AddressBook); err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("address book merge failed after verifier deployment: %w", err)
		}

		verifierAddr, err := dsutil.MaybeFindEthAddress(newAddresses, chain, types.Verifier)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to find verifier address: %w", err)
		}

		if err := cloneEnv.ExistingAddresses.Merge(verifierOut.AddressBook); err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to merge in verifier address: %w", err)
		}

		// Initialize Verifier on VerifierProxy
		initVerifierCfg := verification.VerifierProxyInitializeVerifierConfig{
			ConfigPerChain: map[uint64][]verification.InitializeVerifierConfig{
				chain: {{
					ContractAddress: verifierProxyAddr,
					VerifierAddress: verifierAddr,
				}},
			},
		}

		_, err = verification.InitializeVerifierChangeset.Apply(cloneEnv, initVerifierCfg)
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

		_, err = verification.SetConfigChangeset.Apply(cloneEnv, setCfg)
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

			rmOut, err := rewardmanager.DeployRewardManagerChangeset.Apply(cloneEnv, rewardMgrCfg)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to deploy reward manager on chain %d: %w", chain, err)
			}
			timelockProposals = append(timelockProposals, rmOut.MCMSTimelockProposals...)
			if err := newAddresses.Merge(rmOut.AddressBook); err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("address book merge failed after reward manager deployment: %w", err)
			}
			rewardMgrAddr, err := dsutil.MaybeFindEthAddress(newAddresses, chain, types.RewardManager)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to find reward manager address: %w", err)
			}
			if err := cloneEnv.ExistingAddresses.Merge(rmOut.AddressBook); err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to merge in fm address: %w", err)
			}
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

			fmOut, err := feemanager.DeployFeeManagerChangeset.Apply(cloneEnv, feeMgrCfg)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to deploy fee manager on chain %d: %w", chain, err)
			}
			if err := newAddresses.Merge(fmOut.AddressBook); err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("address book merge failed after fee manager deployment: %w", err)
			}

			feeManagerAddress, err := dsutil.MaybeFindEthAddress(newAddresses, chain, types.FeeManager)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("fee manager address not found for chain %d: %w", chain, err)
			}

			if err := cloneEnv.ExistingAddresses.Merge(fmOut.AddressBook); err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to merge in fm address: %w", err)
			}

			timelockProposals = append(timelockProposals, fmOut.MCMSTimelockProposals...)

			// set the native surcharge on the fee manager
			setNativeCfg := feemanager.SetNativeSurchargeConfig{
				ConfigPerChain: map[uint64][]feemanager.SetNativeSurcharge{
					chain: {
						feemanager.SetNativeSurcharge{
							FeeManagerAddress: feeManagerAddress,
							Surcharge:         cfg.Billing.Config.Surcharge,
						},
					},
				},
			}

			_, err = feemanager.SetNativeSurchargeChangeset.Apply(cloneEnv, setNativeCfg)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to set native surcharge on chain %d: %w", chain, err)
			}

			// Update VerifierProxy to set the FeeManager address
			setFeeManagerCfg := verification.VerifierProxySetFeeManagerConfig{
				ConfigPerChain: map[uint64][]verification.SetFeeManagerConfig{
					chain: {
						verification.SetFeeManagerConfig{
							ContractAddress:   verifierProxyAddr,
							FeeManagerAddress: feeManagerAddress,
						},
					},
				},
			}

			_, err = verification.SetFeeManagerChangeset.Apply(cloneEnv, setFeeManagerCfg)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to set fee manager on verifier proxy on chain %d: %w", chain, err)
			}

			// Update RewardManager to set the FeeManager address
			rmSetFeeManagerCfg := rewardmanager.SetFeeManagerConfig{
				ConfigsByChain: map[uint64][]rewardmanager.SetFeeManager{
					chain: {
						rewardmanager.SetFeeManager{
							FeeManagerAddress:    feeManagerAddress,
							RewardManagerAddress: rewardMgrAddr,
						},
					},
				},
			}
			_, err = rewardmanager.SetFeeManagerChangeset.Apply(cloneEnv, rmSetFeeManagerCfg)
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to set fee manager on reward manager on chain %d: %w", chain, err)
			}
		}
	}

	mergedTimelockProposal, err := mcmsutil.MergeSimilarTimelockProposals(timelockProposals)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to merge timelock proposals: %w", err)
	}

	return deployment.ChangesetOutput{
		AddressBook:           newAddresses,
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
