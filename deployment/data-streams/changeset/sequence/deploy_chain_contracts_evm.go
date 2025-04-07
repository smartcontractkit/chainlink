package sequence

import (
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
)

// deployChainComponentsEVM deploys all necessary components for a single evm chain
func deployChainComponentsEVM(env deployment.Environment, chain uint64, cfg DeployDataStreams, newAddresses deployment.AddressBook) ([]mcms.TimelockProposal, error) {
	var timelockProposals []mcms.TimelockProposal

	// Step 1: Deploy MCMS if configured
	if cfg.Ownership.ShouldDeployMCMS {
		mcmsProposals, err := deployMCMS(env, chain, cfg, newAddresses)
		if err != nil {
			return nil, fmt.Errorf("failed to deploy MCMS: %w", err)
		}
		timelockProposals = append(timelockProposals, mcmsProposals...)
	}

	// Step 2: Deploy VerifierProxy
	verifierProxyAddr, verifierProxyProposals, err := deployVerifierProxy(env, chain, cfg, newAddresses)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy verifier proxy: %w", err)
	}
	timelockProposals = append(timelockProposals, verifierProxyProposals...)

	// Step 3: Deploy Verifier
	verifierAddr, verifierProposals, err := deployVerifier(env, chain, cfg, verifierProxyAddr, newAddresses)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy verifier: %w", err)
	}
	timelockProposals = append(timelockProposals, verifierProposals...)

	// Step 4: Initialize Verifier on VerifierProxy
	if err := initializeVerifier(env, chain, verifierProxyAddr, verifierAddr); err != nil {
		return nil, fmt.Errorf("failed to initialize verifier: %w", err)
	}

	// Step 5: Set Verifier Config
	if err := setVerifierConfig(env, chain, cfg, verifierAddr); err != nil {
		return nil, fmt.Errorf("failed to set verifier config: %w", err)
	}

	// Step 6: Deploy and configure billing components if enabled
	if cfg.Billing.Enabled && cfg.Billing.Config != nil {
		billingProposals, err := deployBillingComponents(env, chain, cfg, verifierProxyAddr, newAddresses)
		if err != nil {
			return nil, fmt.Errorf("failed to deploy billing components: %w", err)
		}
		timelockProposals = append(timelockProposals, billingProposals...)
	}

	return timelockProposals, nil
}

// deployVerifierProxy deploys VerifierProxy contract
func deployVerifierProxy(env deployment.Environment, chain uint64, cfg DeployDataStreams, newAddresses deployment.AddressBook) (common.Address, []mcms.TimelockProposal, error) {
	verifierProxyCfg := verification.DeployVerifierProxyConfig{
		ChainsToDeploy: map[uint64]verification.DeployVerifierProxy{
			chain: {}, // Implement AccessController as needed
		},
		Version:   deployment.Version0_5_0,
		Ownership: cfg.Ownership.AsSettings(),
	}

	proxyOut, err := verification.DeployVerifierProxyChangeset.Apply(env, verifierProxyCfg)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("failed to deploy verifier proxy on chain: %d err %w", chain, err)
	}

	if err := newAddresses.Merge(proxyOut.AddressBook); err != nil {
		return common.Address{}, nil, fmt.Errorf("address book merge failed after verifier proxy deployment: %w", err)
	}

	verifierProxyAddr, err := dsutil.MaybeFindEthAddress(newAddresses, chain, types.VerifierProxy)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("failed to find verifier proxy address: %w", err)
	}

	if err := env.ExistingAddresses.Merge(proxyOut.AddressBook); err != nil {
		return common.Address{}, nil, fmt.Errorf("failed to merge verifier proxy address: %w", err)
	}

	return verifierProxyAddr, proxyOut.MCMSTimelockProposals, nil
}

// deployVerifier deploys Verifier contract
func deployVerifier(env deployment.Environment, chain uint64, cfg DeployDataStreams, verifierProxyAddr common.Address, newAddresses deployment.AddressBook) (common.Address, []mcms.TimelockProposal, error) {
	verifierCfg := verification.DeployVerifierConfig{
		ChainsToDeploy: map[uint64]verification.DeployVerifier{
			chain: {VerifierProxyAddress: verifierProxyAddr},
		},
		Ownership: cfg.Ownership.AsSettings(),
	}

	verifierOut, err := verification.DeployVerifierChangeset.Apply(env, verifierCfg)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("failed to deploy verifier on chain %d: %w", chain, err)
	}

	if err := newAddresses.Merge(verifierOut.AddressBook); err != nil {
		return common.Address{}, nil, fmt.Errorf("address book merge failed after verifier deployment: %w", err)
	}

	verifierAddr, err := dsutil.MaybeFindEthAddress(newAddresses, chain, types.Verifier)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("failed to find verifier address: %w", err)
	}

	if err := env.ExistingAddresses.Merge(verifierOut.AddressBook); err != nil {
		return common.Address{}, nil, fmt.Errorf("failed to merge in verifier address: %w", err)
	}

	return verifierAddr, verifierOut.MCMSTimelockProposals, nil
}

// initializeVerifier initializes the Verifier in VerifierProxy
func initializeVerifier(env deployment.Environment, chain uint64, verifierProxyAddr, verifierAddr common.Address) error {
	initVerifierCfg := verification.VerifierProxyInitializeVerifierConfig{
		ConfigPerChain: map[uint64][]verification.InitializeVerifierConfig{
			chain: {{
				ContractAddress: verifierProxyAddr,
				VerifierAddress: verifierAddr,
			}},
		},
	}

	_, err := verification.InitializeVerifierChangeset.Apply(env, initVerifierCfg)
	if err != nil {
		return fmt.Errorf("failed to initialize verifier on chain %d: %w", chain, err)
	}

	return nil
}

// setVerifierConfig sets the configuration for the Verifier
func setVerifierConfig(env deployment.Environment, chain uint64, cfg DeployDataStreams, verifierAddr common.Address) error {
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

	_, err := verification.SetConfigChangeset.Apply(env, setCfg)
	if err != nil {
		return fmt.Errorf("failed to set config on chain %d: %w", chain, err)
	}

	return nil
}

// deployBillingComponents deploys and configures RewardManager and FeeManager
func deployBillingComponents(env deployment.Environment, chain uint64, cfg DeployDataStreams, verifierProxyAddr common.Address, newAddresses deployment.AddressBook) ([]mcms.TimelockProposal, error) {
	var timelockProposals []mcms.TimelockProposal

	// Step 1: Deploy RewardManager
	rewardMgrAddr, rmProposals, err := deployRewardManager(env, chain, cfg, newAddresses)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy reward manager: %w", err)
	}
	timelockProposals = append(timelockProposals, rmProposals...)

	// Step 2: Deploy FeeManager
	feeManagerAddr, fmProposals, err := deployFeeManager(env, chain, cfg, verifierProxyAddr, rewardMgrAddr, newAddresses)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy fee manager: %w", err)
	}
	timelockProposals = append(timelockProposals, fmProposals...)

	// Step 3: Configure native surcharge on FeeManager
	if err := setNativeSurcharge(env, chain, cfg, feeManagerAddr); err != nil {
		return nil, fmt.Errorf("failed to set native surcharge: %w", err)
	}

	// Step 4: Set FeeManager on VerifierProxy
	if err := setFeeManagerOnVerifierProxy(env, chain, verifierProxyAddr, feeManagerAddr); err != nil {
		return nil, fmt.Errorf("failed to set fee manager on verifier proxy: %w", err)
	}

	// Step 5: Set FeeManager on RewardManager
	if err := setFeeManagerOnRewardManager(env, chain, rewardMgrAddr, feeManagerAddr); err != nil {
		return nil, fmt.Errorf("failed to set fee manager on reward manager: %w", err)
	}

	return timelockProposals, nil
}

// deployRewardManager deploys the RewardManager contract
func deployRewardManager(env deployment.Environment, chain uint64, cfg DeployDataStreams, newAddresses deployment.AddressBook) (common.Address, []mcms.TimelockProposal, error) {
	rewardMgrCfg := rewardmanager.DeployRewardManagerConfig{
		ChainsToDeploy: map[uint64]rewardmanager.DeployRewardManager{
			chain: {LinkTokenAddress: cfg.Billing.Config.LinkTokenAddress},
		},
		Ownership: cfg.Ownership.AsSettings(),
	}

	rmOut, err := rewardmanager.DeployRewardManagerChangeset.Apply(env, rewardMgrCfg)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("failed to deploy reward manager on chain %d: %w", chain, err)
	}

	if err := newAddresses.Merge(rmOut.AddressBook); err != nil {
		return common.Address{}, nil, fmt.Errorf("address book merge failed after reward manager deployment: %w", err)
	}

	rewardMgrAddr, err := dsutil.MaybeFindEthAddress(newAddresses, chain, types.RewardManager)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("failed to find reward manager address: %w", err)
	}

	if err := env.ExistingAddresses.Merge(rmOut.AddressBook); err != nil {
		return common.Address{}, nil, fmt.Errorf("failed to merge in reward manager address: %w", err)
	}

	return rewardMgrAddr, rmOut.MCMSTimelockProposals, nil
}

// deployFeeManager deploys the FeeManager contract
func deployFeeManager(env deployment.Environment, chain uint64, cfg DeployDataStreams, verifierProxyAddr, rewardMgrAddr common.Address, newAddresses deployment.AddressBook) (common.Address, []mcms.TimelockProposal, error) {
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

	fmOut, err := feemanager.DeployFeeManagerChangeset.Apply(env, feeMgrCfg)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("failed to deploy fee manager on chain %d: %w", chain, err)
	}

	if err := newAddresses.Merge(fmOut.AddressBook); err != nil {
		return common.Address{}, nil, fmt.Errorf("address book merge failed after fee manager deployment: %w", err)
	}

	feeManagerAddress, err := dsutil.MaybeFindEthAddress(newAddresses, chain, types.FeeManager)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("fee manager address not found for chain %d: %w", chain, err)
	}

	if err := env.ExistingAddresses.Merge(fmOut.AddressBook); err != nil {
		return common.Address{}, nil, fmt.Errorf("failed to merge in fee manager address: %w", err)
	}

	return feeManagerAddress, fmOut.MCMSTimelockProposals, nil
}

// setNativeSurcharge sets the native surcharge on the FeeManager
func setNativeSurcharge(env deployment.Environment, chain uint64, cfg DeployDataStreams, feeManagerAddr common.Address) error {
	setNativeCfg := feemanager.SetNativeSurchargeConfig{
		ConfigPerChain: map[uint64][]feemanager.SetNativeSurcharge{
			chain: {
				feemanager.SetNativeSurcharge{
					FeeManagerAddress: feeManagerAddr,
					Surcharge:         cfg.Billing.Config.Surcharge,
				},
			},
		},
	}

	_, err := feemanager.SetNativeSurchargeChangeset.Apply(env, setNativeCfg)
	if err != nil {
		return fmt.Errorf("failed to set native surcharge on chain %d: %w", chain, err)
	}

	return nil
}

// setFeeManagerOnVerifierProxy sets the FeeManager address on the VerifierProxy
func setFeeManagerOnVerifierProxy(env deployment.Environment, chain uint64, verifierProxyAddr, feeManagerAddr common.Address) error {
	setFeeManagerCfg := verification.VerifierProxySetFeeManagerConfig{
		ConfigPerChain: map[uint64][]verification.SetFeeManagerConfig{
			chain: {
				verification.SetFeeManagerConfig{
					ContractAddress:   verifierProxyAddr,
					FeeManagerAddress: feeManagerAddr,
				},
			},
		},
	}

	_, err := verification.SetFeeManagerChangeset.Apply(env, setFeeManagerCfg)
	if err != nil {
		return fmt.Errorf("failed to set fee manager on verifier proxy on chain %d: %w", chain, err)
	}

	return nil
}

// setFeeManagerOnRewardManager sets the FeeManager address on the RewardManager
func setFeeManagerOnRewardManager(env deployment.Environment, chain uint64, rewardMgrAddr, feeManagerAddr common.Address) error {
	rmSetFeeManagerCfg := rewardmanager.SetFeeManagerConfig{
		ConfigsByChain: map[uint64][]rewardmanager.SetFeeManager{
			chain: {
				rewardmanager.SetFeeManager{
					FeeManagerAddress:    feeManagerAddr,
					RewardManagerAddress: rewardMgrAddr,
				},
			},
		},
	}

	_, err := rewardmanager.SetFeeManagerChangeset.Apply(env, rmSetFeeManagerCfg)
	if err != nil {
		return fmt.Errorf("failed to set fee manager on reward manager on chain %d: %w", chain, err)
	}

	return nil
}

// deployMCMS deploys Multi-Chain Management System
func deployMCMS(env deployment.Environment, chain uint64, cfg DeployDataStreams, newAddresses deployment.AddressBook) ([]mcms.TimelockProposal, error) {
	mcmsDeployCfg := changeset.DeployMCMSConfig{
		ChainsToDeploy: []uint64{chain},
		Ownership:      cfg.Ownership.AsSettings(),
		Config:         *cfg.Ownership.DeployMCMSConfig,
	}

	mcmsDeployOut, err := changeset.DeployAndTransferMCMSChangeset.Apply(env, mcmsDeployCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy MCMS on chain %d: %w", chain, err)
	}

	if err := newAddresses.Merge(mcmsDeployOut.AddressBook); err != nil {
		return nil, fmt.Errorf("address book merge failed after MCMS deployment: %w", err)
	}

	if err := env.ExistingAddresses.Merge(mcmsDeployOut.AddressBook); err != nil {
		return nil, fmt.Errorf("failed to merge existing addresses with mcms addresses: %w", err)
	}

	return mcmsDeployOut.MCMSTimelockProposals, nil
}
