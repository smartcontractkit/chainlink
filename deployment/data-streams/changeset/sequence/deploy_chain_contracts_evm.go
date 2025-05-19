package sequence

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/mcms"

	ds "github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	mcms2 "github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/mcms"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/metadata"
	dsutil "github.com/smartcontractkit/chainlink/deployment/data-streams/utils"

	"github.com/smartcontractkit/chainlink/deployment"
	feemanager "github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/fee-manager"
	rewardmanager "github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/reward-manager"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/verification"
)

// deployChainComponentsEVM deploys all necessary components for a single evm chain
func deployChainComponentsEVM(env *cldf.Environment, chain uint64, cfg DeployDataStreams,
	newAddresses metadata.DataStreamsMutableDataStore) ([]mcms.TimelockProposal, error) {
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
func deployVerifierProxy(env *cldf.Environment, chain uint64, cfg DeployDataStreams, newAddresses metadata.DataStreamsMutableDataStore) (common.Address, []mcms.TimelockProposal, error) {
	verifierProxyCfg := verification.DeployVerifierProxyConfig{
		ChainsToDeploy: map[uint64]verification.DeployVerifierProxy{
			chain: {},
		},
		Version:   deployment.Version0_5_0,
		Ownership: cfg.Ownership.AsSettings(),
	}
	proxyOut, err := verification.DeployVerifierProxyChangeset.Apply(*env, verifierProxyCfg)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("failed to deploy verifier proxy on chain: %d err %w", chain, err)
	}

	if err := mergeNewAddresses(env, newAddresses, proxyOut.DataStore.Seal()); err != nil {
		return common.Address{}, nil, fmt.Errorf("failed to merge new addresses: %w", err)
	}

	// Filter without version here should be safe as we only expect 1 address
	address, err := dsutil.GetContractAddress(newAddresses.Addresses(), types.VerifierProxy)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("failed to get verifier proxy address: %w", err)
	}
	verifierProxyAddr := common.HexToAddress(address)

	return verifierProxyAddr, proxyOut.MCMSTimelockProposals, nil
}

// deployVerifier deploys Verifier contract
func deployVerifier(env *cldf.Environment, chain uint64, cfg DeployDataStreams, verifierProxyAddr common.Address, newAddresses metadata.DataStreamsMutableDataStore) (common.Address, []mcms.TimelockProposal, error) {
	verifierCfg := verification.DeployVerifierConfig{
		ChainsToDeploy: map[uint64]verification.DeployVerifier{
			chain: {VerifierProxyAddress: verifierProxyAddr},
		},
		Ownership: cfg.Ownership.AsSettings(),
	}

	verifierOut, err := verification.DeployVerifierChangeset.Apply(*env, verifierCfg)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("failed to deploy verifier on chain %d: %w", chain, err)
	}

	if err := mergeNewAddresses(env, newAddresses, verifierOut.DataStore.Seal()); err != nil {
		return common.Address{}, nil, fmt.Errorf("failed to merge new addresses: %w", err)
	}

	// Filter without version here should be safe as we only expect 1 address
	address, err := dsutil.GetContractAddress(newAddresses.Addresses(), types.Verifier)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("failed to get verifier address: %w", err)
	}
	verifierAddr := common.HexToAddress(address)

	return verifierAddr, verifierOut.MCMSTimelockProposals, nil
}

// initializeVerifier initializes the Verifier in VerifierProxy
func initializeVerifier(env *cldf.Environment, chain uint64, verifierProxyAddr, verifierAddr common.Address) error {
	initVerifierCfg := verification.VerifierProxyInitializeVerifierConfig{
		ConfigPerChain: map[uint64][]verification.InitializeVerifierConfig{
			chain: {{
				VerifierProxyAddress: verifierProxyAddr,
				VerifierAddress:      verifierAddr,
			}},
		},
	}

	_, err := verification.InitializeVerifierChangeset.Apply(*env, initVerifierCfg)
	if err != nil {
		return fmt.Errorf("failed to initialize verifier on chain %d: %w", chain, err)
	}

	return nil
}

// setVerifierConfig sets the configuration for the Verifier
func setVerifierConfig(env *cldf.Environment, chain uint64, cfg DeployDataStreams, verifierAddr common.Address) error {
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

	_, err := verification.SetConfigChangeset.Apply(*env, setCfg)
	if err != nil {
		return fmt.Errorf("failed to set config on chain %d: %w", chain, err)
	}

	return nil
}

// deployBillingComponents deploys and configures RewardManager and FeeManager
func deployBillingComponents(env *cldf.Environment, chain uint64, cfg DeployDataStreams, verifierProxyAddr common.Address, newAddresses metadata.DataStreamsMutableDataStore) ([]mcms.TimelockProposal, error) {
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
func deployRewardManager(env *cldf.Environment, chain uint64, cfg DeployDataStreams, newAddresses metadata.DataStreamsMutableDataStore) (common.Address, []mcms.TimelockProposal, error) {
	rewardMgrCfg := rewardmanager.DeployRewardManagerConfig{
		ChainsToDeploy: map[uint64]rewardmanager.DeployRewardManager{
			chain: {LinkTokenAddress: cfg.Billing.Config.LinkTokenAddress},
		},
		Ownership: cfg.Ownership.AsSettings(),
	}

	rmOut, err := rewardmanager.DeployRewardManagerChangeset.Apply(*env, rewardMgrCfg)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("failed to deploy reward manager on chain %d: %w", chain, err)
	}

	if err := mergeNewAddresses(env, newAddresses, rmOut.DataStore.Seal()); err != nil {
		return common.Address{}, nil, fmt.Errorf("failed to merge new addresses: %w", err)
	}

	// Filter without version here should be safe as we only expect 1 address
	address, err := dsutil.GetContractAddress(newAddresses.Addresses(), types.RewardManager)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("failed to get verifier proxy address: %w", err)
	}
	rmAddr := common.HexToAddress(address)

	return rmAddr, rmOut.MCMSTimelockProposals, nil
}

// deployFeeManager deploys the FeeManager contract
func deployFeeManager(env *cldf.Environment, chain uint64, cfg DeployDataStreams, verifierProxyAddr, rewardMgrAddr common.Address, newAddresses metadata.DataStreamsMutableDataStore) (common.Address, []mcms.TimelockProposal, error) {
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

	fmOut, err := feemanager.DeployFeeManagerChangeset.Apply(*env, feeMgrCfg)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("failed to deploy fee manager on chain %d: %w", chain, err)
	}

	if err := mergeNewAddresses(env, newAddresses, fmOut.DataStore.Seal()); err != nil {
		return common.Address{}, nil, fmt.Errorf("failed to merge new addresses: %w", err)
	}

	// Filter without version here should be safe as we only expect 1 address
	address, err := dsutil.GetContractAddress(newAddresses.Addresses(), types.FeeManager)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("failed to get verifier proxy address: %w", err)
	}
	fmAddr := common.HexToAddress(address)

	return fmAddr, fmOut.MCMSTimelockProposals, nil
}

// setNativeSurcharge sets the native surcharge on the FeeManager
func setNativeSurcharge(env *cldf.Environment, chain uint64, cfg DeployDataStreams, feeManagerAddr common.Address) error {
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

	_, err := feemanager.SetNativeSurchargeChangeset.Apply(*env, setNativeCfg)
	if err != nil {
		return fmt.Errorf("failed to set native surcharge on chain %d: %w", chain, err)
	}

	return nil
}

// setFeeManagerOnVerifierProxy sets the FeeManager address on the VerifierProxy
func setFeeManagerOnVerifierProxy(env *cldf.Environment, chain uint64, verifierProxyAddr, feeManagerAddr common.Address) error {
	setFeeManagerCfg := verification.VerifierProxySetFeeManagerConfig{
		ConfigPerChain: map[uint64][]verification.SetFeeManagerConfig{
			chain: {
				verification.SetFeeManagerConfig{
					VerifierProxyAddress: verifierProxyAddr,
					FeeManagerAddress:    feeManagerAddr,
				},
			},
		},
	}

	_, err := verification.SetFeeManagerChangeset.Apply(*env, setFeeManagerCfg)
	if err != nil {
		return fmt.Errorf("failed to set fee manager on verifier proxy on chain %d: %w", chain, err)
	}

	return nil
}

// setFeeManagerOnRewardManager sets the FeeManager address on the RewardManager
func setFeeManagerOnRewardManager(env *cldf.Environment, chain uint64, rewardMgrAddr, feeManagerAddr common.Address) error {
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

	_, err := rewardmanager.SetFeeManagerChangeset.Apply(*env, rmSetFeeManagerCfg)
	if err != nil {
		return fmt.Errorf("failed to set fee manager on reward manager on chain %d: %w", chain, err)
	}

	return nil
}

type DeployOutput struct {
	Addresses   []string
	Environment cldf.Environment
	Proposals   []mcms.TimelockProposal
}

// deployMCMS deploys the MCMS contracts
func deployMCMS(env *cldf.Environment, chain uint64, cfg DeployDataStreams, cumulativeAddresses metadata.DataStreamsMutableDataStore) ([]mcms.TimelockProposal, error) {
	mcmsDeployCfg := mcms2.DeployMCMSConfig{
		ChainsToDeploy: []uint64{chain},
		Ownership:      cfg.Ownership.AsSettings(),
		Config:         *cfg.Ownership.DeployMCMSConfig,
	}

	mcmsDeployOut, err := mcms2.DeployAndTransferMCMSChangeset.Apply(*env, mcmsDeployCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy MCMS on chain %d: %w", chain, err)
	}

	if err := mergeNewAddresses(env, cumulativeAddresses, mcmsDeployOut.DataStore.Seal()); err != nil {
		return nil, fmt.Errorf("datastore merges failed after MCMS deployment: %w", err)
	}

	return mcmsDeployOut.MCMSTimelockProposals, nil
}

// mergeNewAddresses merges new addresses into the existing environment and cumulative address book
// This is used when chaining together changesets and accumulating addresses while also updating
// the environment with new addresses so that downstream operations can use them
func mergeNewAddresses(env *cldf.Environment,
	cumulativeAddrs metadata.DataStreamsMutableDataStore,
	newAddrs ds.DataStore[ds.DefaultMetadata, ds.DefaultMetadata]) error {
	// Step 1 is update the existing environment to reflect newly deployed addresses
	updatedEnvironmentAddrs, err := ds.ToDefault(newAddrs)
	if err != nil {
		return fmt.Errorf("failed to convert data store to default format: %w", err)
	}
	if err = updatedEnvironmentAddrs.Merge(env.DataStore); err != nil {
		return fmt.Errorf("failed to merge new addresses into existing environment: %w", err)
	}
	env.DataStore = updatedEnvironmentAddrs.Seal()

	// Step 2 update newAddresses which is returned at the end of the entire "sequence"
	envDatastore, err := ds.FromDefault[metadata.SerializedContractMetadata, ds.DefaultMetadata](newAddrs)
	if err != nil {
		return fmt.Errorf("failed to convert data store from default format: %w", err)
	}
	if err = cumulativeAddrs.Merge(envDatastore); err != nil {
		return fmt.Errorf("failed to merge new addresses into address book: %w", err)
	}
	return nil
}
