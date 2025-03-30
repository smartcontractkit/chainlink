package sequence

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/deployment"
	feemanager "github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/fee-manager"
	rewardmanager "github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/reward-manager"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/verification"
	dsutil "github.com/smartcontractkit/chainlink/deployment/data-streams/utils"
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

type DeployDataStreams struct {
	VerifierConfig verification.SetConfig
	BillingConfig  *BillingConfig
}

func deployDataStreamsLogic(e deployment.Environment, cc DeployDataStreamsConfig) (deployment.ChangesetOutput, error) {
	ab := deployment.NewMemoryAddressBook() // changeset output expects only new addresses

	for chain, cfg := range cc.ChainsToDeploy {
		// Deploy Verifier Proxy
		verifierProxyCfg := verification.DeployVerifierProxyConfig{
			ChainsToDeploy: map[uint64]verification.DeployVerifierProxy{
				chain: {}, // Implement AccessController as needed
			},
			Version: deployment.Version0_5_0,
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

		// Deploy Verifier
		verifierCfg := verification.DeployVerifierConfig{
			ChainsToDeploy: map[uint64]verification.DeployVerifier{
				chain: {VerifierProxyAddress: verifierProxyAddr},
			},
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
		if cfg.BillingConfig != nil {
			// Deploy RewardManager
			rewardMgrCfg := rewardmanager.DeployRewardManagerConfig{
				ChainsToDeploy: map[uint64]rewardmanager.DeployRewardManager{
					chain: {LinkTokenAddress: cfg.BillingConfig.LinkTokenAddress},
				},
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

			// Deploy FeeManager
			feeMgrCfg := feemanager.DeployFeeManagerConfig{
				ChainsToDeploy: map[uint64]feemanager.DeployFeeManager{
					chain: {
						LinkTokenAddress:     cfg.BillingConfig.LinkTokenAddress,
						NativeTokenAddress:   cfg.BillingConfig.NativeTokenAddress,
						VerifierProxyAddress: verifierProxyAddr,
						RewardManagerAddress: rewardMgrAddr,
					},
				},
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
		}
	}

	return deployment.ChangesetOutput{
		AddressBook: ab,
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
