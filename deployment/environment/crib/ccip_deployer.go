package crib

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/big"
	"os"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	"github.com/rs/zerolog"
	xerrgroup "golang.org/x/sync/errgroup"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-ccip/chainconfig"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/rmn_home"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/rmn_remote"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
	ccipChangesetSolana "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/solana"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_5_1"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	ccipops "github.com/smartcontractkit/chainlink/deployment/ccip/operation/evm/v1_6"
	ccipseq "github.com/smartcontractkit/chainlink/deployment/ccip/sequence/evm/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"

	solconfig "github.com/smartcontractkit/chainlink-ccip/chains/solana/contracts/tests/config"
	solTestTokenPool "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/test_token_pool"
	solcommon "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	solstate "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"
	soltokens "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"

	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
)

const (
	tokenApproveCheckedAmount = 1e4 * 1e9
)

// DeployHomeChainContracts deploys the home chain contracts so that the chainlink nodes can use the CR address in Capabilities.ExternalRegistry
// Afterward, we call DeployHomeChainChangeset changeset with nodeinfo ( the peer id and all)
func DeployHomeChainContracts(ctx context.Context, lggr logger.Logger, envConfig devenv.EnvironmentConfig, homeChainSel uint64, feedChainSel uint64) (deployment.CapabilityRegistryConfig, cldf.AddressBook, error) {
	e, _, err := devenv.NewEnvironment(func() context.Context { return ctx }, lggr, envConfig)
	if err != nil {
		return deployment.CapabilityRegistryConfig{}, nil, err
	}
	if e == nil {
		return deployment.CapabilityRegistryConfig{}, nil, errors.New("environment is nil")
	}

	evmChains := e.BlockChains.EVMChains()

	nodes, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	if err != nil {
		return deployment.CapabilityRegistryConfig{}, e.ExistingAddresses, fmt.Errorf("failed to get node info from env: %w", err)
	}

	// Fund the deployer
	solChainSelectors := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainselectors.FamilySolana))

	for _, selector := range solChainSelectors {
		lggr.Infof("Funding solana deployer account %v", e.BlockChains.SolanaChains()[selector].DeployerKey.PublicKey())
		err = memory.FundSolanaAccounts(e.GetContext(), []solana.PublicKey{e.BlockChains.SolanaChains()[selector].DeployerKey.PublicKey()}, 10000, e.BlockChains.SolanaChains()[selector].Client)
		if err != nil {
			return deployment.CapabilityRegistryConfig{}, nil, err
		}
	}

	p2pIds := nodes.NonBootstraps().PeerIDs()
	cfg := make(map[uint64]commontypes.MCMSWithTimelockConfigV2)
	for _, chain := range e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainselectors.FamilyEVM)) {
		mcmsConfig, err := mcmstypes.NewConfig(1, []common.Address{evmChains[chain].DeployerKey.From}, []mcmstypes.Config{})
		if err != nil {
			return deployment.CapabilityRegistryConfig{}, e.ExistingAddresses, fmt.Errorf("failed to create mcms config: %w", err)
		}
		cfg[chain] = commontypes.MCMSWithTimelockConfigV2{
			Canceller:        mcmsConfig,
			Bypasser:         mcmsConfig,
			Proposer:         mcmsConfig,
			TimelockMinDelay: big.NewInt(0),
		}
	}
	*e, err = commonchangeset.Apply(nil, *e, commonchangeset.Configure(
		cldf.CreateLegacyChangeSet(commonchangeset.DeployMCMSWithTimelockV2),
		cfg,
	), commonchangeset.Configure(
		cldf.CreateLegacyChangeSet(v1_6.DeployHomeChainChangeset),
		v1_6.DeployHomeChainConfig{
			HomeChainSel:             homeChainSel,
			RMNStaticConfig:          testhelpers.NewTestRMNStaticConfig(),
			RMNDynamicConfig:         testhelpers.NewTestRMNDynamicConfig(),
			NodeOperators:            testhelpers.NewTestNodeOperator(evmChains[homeChainSel].DeployerKey.From),
			NodeP2PIDsPerNodeOpAdmin: map[string][][32]byte{"NodeOperator": p2pIds},
		},
	))
	if err != nil {
		return deployment.CapabilityRegistryConfig{}, e.ExistingAddresses, fmt.Errorf("changeset sequence execution failed with error: %w", err)
	}
	state, err := stateview.LoadOnchainState(*e)
	if err != nil {
		return deployment.CapabilityRegistryConfig{}, e.ExistingAddresses, fmt.Errorf("failed to load on chain state: %w", err)
	}
	capRegAddr := state.Chains[homeChainSel].CapabilityRegistry.Address()
	if capRegAddr == common.HexToAddress("0x") {
		return deployment.CapabilityRegistryConfig{}, e.ExistingAddresses, fmt.Errorf("cap Reg address not found: %w", err)
	}
	capRegConfig := deployment.CapabilityRegistryConfig{
		EVMChainID:  homeChainSel,
		Contract:    state.Chains[homeChainSel].CapabilityRegistry.Address(),
		NetworkType: relay.NetworkEVM,
	}
	return capRegConfig, e.ExistingAddresses, nil
}

// DeployCCIPAndAddLanes is the actual ccip setup once the nodes are initialized.
func DeployCCIPAndAddLanes(ctx context.Context, lggr logger.Logger, envConfig devenv.EnvironmentConfig, homeChainSel, feedChainSel uint64, ab cldf.AddressBook, rmnEnabled bool) (DeployCCIPOutput, error) {
	e, don, err := devenv.NewEnvironment(func() context.Context { return ctx }, lggr, envConfig)
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to initiate new environment: %w", err)
	}
	e.ExistingAddresses = ab

	// ------ Part 1 -----
	// Setup because we only need to deploy the contracts and distribute job specs
	lggr.Infow("setting up chains...")
	*e, err = setupChains(lggr, e, homeChainSel, feedChainSel)
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to apply setting up chain changesets: %w", err)
	}

	state, err := stateview.LoadOnchainState(*e)
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	// Set up EVM lanes
	lggr.Infow("setting up EVM lanes...")
	*e, err = setupSolEvmLanes(lggr, e, state, homeChainSel, feedChainSel)
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to apply connecting lanes changesets: %w", err)
	}
	// ------ Part 1 -----

	// ----- Part 2 -----
	lggr.Infow("setting up ocr...")
	*e, err = mustOCR(e, homeChainSel, feedChainSel, true, rmnEnabled)
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to apply OCR changesets: %w", err)
	}

	// distribute funds to transmitters
	// we need to use the nodeinfo from the envConfig here, because multiAddr is not
	// populated in the environment variable
	lggr.Infow("distributing funds...")
	err = distributeTransmitterFunds(lggr, don.PluginNodes(), *e)
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to distribute funds to node transmitters: %w", err)
	}

	addresses, err := e.ExistingAddresses.Addresses()
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to convert address book to address book map: %w", err)
	}
	return DeployCCIPOutput{
		AddressBook: *cldf.NewMemoryAddressBookFromMap(addresses),
		NodeIDs:     e.NodeIDs,
	}, nil
}

// ConfigureCCIPOCR is a group of changesets used from CRIB to redeploy the chainlink don on an existing setup
func ConfigureCCIPOCR(ctx context.Context, lggr logger.Logger, envConfig devenv.EnvironmentConfig, homeChainSel, feedChainSel uint64, ab cldf.AddressBook, rmnEnabled bool) (DeployCCIPOutput, error) {
	e, don, err := devenv.NewEnvironment(func() context.Context { return ctx }, lggr, envConfig)
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to initiate new environment: %w", err)
	}
	e.ExistingAddresses = ab

	lggr.Infow("resetting ocr...")
	*e, err = mustOCR(e, homeChainSel, feedChainSel, false, rmnEnabled)
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to apply changesets for setting up OCR: %w", err)
	}
	err = distributeTransmitterFunds(lggr, don.PluginNodes(), *e)
	if err != nil {
		return DeployCCIPOutput{}, err
	}

	addresses, err := e.ExistingAddresses.Addresses()
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to get convert address book to address book map: %w", err)
	}
	return DeployCCIPOutput{
		AddressBook: *cldf.NewMemoryAddressBookFromMap(addresses),
		NodeIDs:     e.NodeIDs,
	}, nil
}

// FundCCIPTransmitters is used from CRIB to provide funds to the node transmitters
// This function sends funds from the deployer key to the chainlink node transmitters
func FundCCIPTransmitters(ctx context.Context, lggr logger.Logger, envConfig devenv.EnvironmentConfig, ab cldf.AddressBook) (DeployCCIPOutput, error) {
	e, don, err := devenv.NewEnvironment(func() context.Context { return ctx }, lggr, envConfig)
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to initiate new environment: %w", err)
	}
	e.ExistingAddresses = ab

	// distribute funds to transmitters
	// we need to use the nodeinfo from the envConfig here, because multiAddr is not
	// populated in the environment variable
	lggr.Infow("distributing funds...")
	err = distributeTransmitterFunds(lggr, don.PluginNodes(), *e)
	if err != nil {
		return DeployCCIPOutput{}, err
	}

	addresses, err := e.ExistingAddresses.Addresses()
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to get convert address book to address book map: %w", err)
	}
	return DeployCCIPOutput{
		AddressBook: *cldf.NewMemoryAddressBookFromMap(addresses),
		NodeIDs:     e.NodeIDs,
	}, nil
}

func setupChains(lggr logger.Logger, e *cldf.Environment, homeChainSel, feedChainSel uint64) (cldf.Environment, error) {
	evmChainSelectors := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainselectors.FamilyEVM))
	solChainSelectors := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainselectors.FamilySolana))
	chainConfigs := make(map[uint64]v1_6.ChainConfig)
	nodeInfo, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	if err != nil {
		return *e, fmt.Errorf("failed to get node info from env: %w", err)
	}
	prereqCfgs := make([]changeset.DeployPrerequisiteConfigPerChain, 0)
	contractParams := make(map[uint64]ccipseq.ChainContractParams)

	for _, chain := range evmChainSelectors {
		prereqCfgs = append(prereqCfgs, changeset.DeployPrerequisiteConfigPerChain{
			ChainSelector: chain,
		})
		chainConfigs[chain] = v1_6.ChainConfig{
			Readers: nodeInfo.NonBootstraps().PeerIDs(),
			// Number of nodes is 3f+1
			//nolint:gosec // this should always be less than max uint8
			FChain: uint8(len(nodeInfo.NonBootstraps().PeerIDs()) / 3),
			EncodableChainConfig: chainconfig.ChainConfig{
				GasPriceDeviationPPB:    cciptypes.BigInt{Int: big.NewInt(1000)},
				DAGasPriceDeviationPPB:  cciptypes.BigInt{Int: big.NewInt(1_000_000)},
				OptimisticConfirmations: 1,
			},
		}
		contractParams[chain] = ccipseq.ChainContractParams{
			FeeQuoterParams: ccipops.DefaultFeeQuoterParams(),
			OffRampParams:   ccipops.DefaultOffRampParams(),
		}
	}

	if len(solChainSelectors) > 0 {
		var solLinkChangesets []commonchangeset.ConfiguredChangeSet
		// TODO - Find a way to combine this into one loop with AllChainSelectors
		// Currently it seems to throw a nil pointer when run with both solana and evm and needs to be investigated
		for _, chain := range solChainSelectors {
			chainConfigs[chain] = v1_6.ChainConfig{
				Readers: nodeInfo.NonBootstraps().PeerIDs(),
				// #nosec G115 - Overflow is not a concern in this test scenario
				FChain: uint8(len(nodeInfo.NonBootstraps().PeerIDs()) / 3),
				EncodableChainConfig: chainconfig.ChainConfig{
					GasPriceDeviationPPB:    cciptypes.BigInt{Int: big.NewInt(testhelpers.DefaultGasPriceDeviationPPB)},
					DAGasPriceDeviationPPB:  cciptypes.BigInt{Int: big.NewInt(testhelpers.DefaultDAGasPriceDeviationPPB)},
					OptimisticConfirmations: globals.OptimisticConfirmations,
				},
			}

			privKey, err := solana.NewRandomPrivateKey()
			if err != nil {
				return *e, fmt.Errorf("failed to create the link token priv key: %w", err)
			}
			solLinkChangeset := commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(commonchangeset.DeploySolanaLinkToken),
				commonchangeset.DeploySolanaLinkTokenConfig{
					ChainSelector: chain,
					TokenPrivKey:  privKey,
					TokenDecimals: 9,
				},
			)
			solLinkChangesets = append(solLinkChangesets, solLinkChangeset)
		}

		*e, err = commonchangeset.Apply(nil, *e, solLinkChangesets[0], solLinkChangesets[1:]...)
		if err != nil {
			return *e, fmt.Errorf("failed to apply Solana Link token changesets: %w", err)
		}
	}

	*e, err = commonchangeset.Apply(nil, *e,
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(v1_6.UpdateChainConfigChangeset),
			v1_6.UpdateChainConfigConfig{
				HomeChainSelector: homeChainSel,
				RemoteChainAdds:   chainConfigs,
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(commonchangeset.DeployLinkToken),
			evmChainSelectors,
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(changeset.DeployPrerequisitesChangeset),
			changeset.DeployPrerequisiteConfig{
				Configs: prereqCfgs,
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(v1_6.DeployChainContractsChangeset),
			ccipseq.DeployChainContractsConfig{
				HomeChainSelector:      homeChainSel,
				ContractParamsPerChain: contractParams,
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(v1_6.SetRMNRemoteOnRMNProxyChangeset),
			v1_6.SetRMNRemoteOnRMNProxyConfig{
				ChainSelectors: evmChainSelectors,
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(v1_6.CCIPCapabilityJobspecChangeset),
			nil, // ChangeSet does not use a config.
		),
	)
	if err != nil {
		return *e, fmt.Errorf("failed to apply EVM chain changesets: %w", err)
	}

	if len(solChainSelectors) > 0 {
		deployedEnv := testhelpers.DeployedEnv{
			Env:          *e,
			HomeChainSel: homeChainSel,
			FeedChainSel: feedChainSel,
		}

		buildConfig := ccipChangesetSolana.BuildSolanaConfig{
			GitCommitSha:   "c6cd4a526da4",
			DestinationDir: deployedEnv.Env.BlockChains.SolanaChains()[solChainSelectors[0]].ProgramsPath,
			LocalBuild: ccipChangesetSolana.LocalBuildConfig{
				BuildLocally: true,
			},
		}

		solTestReceiver := commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.DeployReceiverForTest),
			ccipChangesetSolana.DeployForTestConfig{
				ChainSelector: solChainSelectors[0],
			},
		)

		lggr.Info("Starting changeset deployment, this will take long on first run due to anchor build for solana programs")
		solCs, err := testhelpers.DeployChainContractsToSolChainCS(deployedEnv, solChainSelectors[0], false, &buildConfig)
		if err != nil {
			return *e, err
		}

		solCs = append(solCs, solTestReceiver)
		*e = deployedEnv.Env

		*e, err = commonchangeset.Apply(nil, *e, solCs[0], solCs[1:]...)
		if err != nil {
			return *e, err
		}
		err = testhelpers.ValidateSolanaState(*e, solChainSelectors)
		if err != nil {
			return *e, err
		}

		lggr.Infow("setup SOL Link pools")
		*e, err = setupSolLinkPools(e)
		if err != nil {
			return *e, fmt.Errorf("failed to setup solana link pools: %w", err)
		}
	}
	lggr.Infow("setup EVM Link pools")
	return setupLinkPools(e)
}

func setupLinkPools(e *cldf.Environment) (cldf.Environment, error) {
	evmChains := e.BlockChains.EVMChains()
	state, err := stateview.LoadOnchainState(*e)
	if err != nil {
		return *e, fmt.Errorf("failed to load onchain state: %w", err)
	}
	chainSelectors := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainselectors.FamilyEVM))
	poolInput := make(map[uint64]v1_5_1.DeployTokenPoolInput)
	pools := make(map[uint64]map[shared.TokenSymbol]v1_5_1.TokenPoolInfo)

	for _, chain := range chainSelectors {
		poolInput[chain] = v1_5_1.DeployTokenPoolInput{
			Type:               shared.BurnMintTokenPool,
			LocalTokenDecimals: 18,
			AllowList:          []common.Address{},
			TokenAddress:       state.Chains[chain].LinkToken.Address(),
		}
		pools[chain] = map[shared.TokenSymbol]v1_5_1.TokenPoolInfo{
			shared.LinkSymbol: {
				Type:          shared.BurnMintTokenPool,
				Version:       deployment.Version1_5_1,
				ExternalAdmin: evmChains[chain].DeployerKey.From,
			},
		}
	}
	env, err := commonchangeset.Apply(nil, *e, commonchangeset.Configure(
		cldf.CreateLegacyChangeSet(v1_5_1.DeployTokenPoolContractsChangeset),
		v1_5_1.DeployTokenPoolContractsConfig{
			TokenSymbol: shared.LinkSymbol,
			NewPools:    poolInput,
		},
	), commonchangeset.Configure(
		cldf.CreateLegacyChangeSet(v1_5_1.ProposeAdminRoleChangeset),
		v1_5_1.TokenAdminRegistryChangesetConfig{
			Pools: pools,
		},
	), commonchangeset.Configure(
		cldf.CreateLegacyChangeSet(v1_5_1.AcceptAdminRoleChangeset),
		v1_5_1.TokenAdminRegistryChangesetConfig{
			Pools: pools,
		},
	), commonchangeset.Configure(
		cldf.CreateLegacyChangeSet(v1_5_1.SetPoolChangeset),
		v1_5_1.TokenAdminRegistryChangesetConfig{
			Pools: pools,
		},
	))

	if err != nil {
		return *e, fmt.Errorf("failed to apply changesets: %w", err)
	}

	state, err = stateview.LoadOnchainState(env)
	if err != nil {
		return *e, fmt.Errorf("failed to load onchain state: %w", err)
	}

	for _, chain := range chainSelectors {
		linkPool := state.Chains[chain].BurnMintTokenPools[shared.LinkSymbol][deployment.Version1_5_1]
		linkToken := state.Chains[chain].LinkToken
		tx, err := linkToken.GrantMintAndBurnRoles(evmChains[chain].DeployerKey, linkPool.Address())
		_, err = cldf.ConfirmIfNoError(evmChains[chain], tx, err)
		if err != nil {
			return *e, fmt.Errorf("failed to grant mint and burn roles for link pool: %w", err)
		}
	}
	return env, err
}

func setupSolLinkPools(e *cldf.Environment) (cldf.Environment, error) {
	sels := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainselectors.FamilySolana))
	state, err := stateview.LoadOnchainState(*e)
	if err != nil {
		return *e, fmt.Errorf("failed to load onchain state: %w", err)
	}
	for _, solChainSel := range sels {
		solTokenAddress := state.SolChains[solChainSel].LinkToken
		bnm := solTestTokenPool.BurnAndMint_PoolType

		*e, err = commonchangeset.Apply(nil, *e,
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(ccipChangesetSolana.CreateSolanaTokenATA),
				ccipChangesetSolana.CreateSolanaTokenATAConfig{
					ChainSelector: solChainSel,
					TokenPubkey:   solTokenAddress,
					// TODO - Seems to be nil, deployer not set properly
					ATAList: []string{e.BlockChains.SolanaChains()[solChainSel].DeployerKey.PublicKey().String()},
				},
			),
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(ccipChangesetSolana.MintSolanaToken),
				ccipChangesetSolana.MintSolanaTokenConfig{
					ChainSelector: solChainSel,
					TokenPubkey:   solTokenAddress.String(),
					AmountToAddress: map[string]uint64{
						e.BlockChains.SolanaChains()[solChainSel].DeployerKey.PublicKey().String(): math.MaxUint64,
					},
				},
			),
			// add solana token pool and token pool lookup table
			commonchangeset.Configure(
				// deploy token pool and set the burn/mint authority to the tokenPool
				cldf.CreateLegacyChangeSet(ccipChangesetSolana.E2ETokenPool),
				ccipChangesetSolana.E2ETokenPoolConfig{
					AddTokenPoolAndLookupTable: []ccipChangesetSolana.TokenPoolConfig{
						{
							ChainSelector: solChainSel,
							TokenPubKey:   solTokenAddress,
							PoolType:      &bnm,
							Metadata:      shared.CLLMetadata,
						},
					},
					RegisterTokenAdminRegistry: []ccipChangesetSolana.RegisterTokenAdminRegistryConfig{
						{
							ChainSelector:           solChainSel,
							TokenPubKey:             solTokenAddress,
							TokenAdminRegistryAdmin: e.BlockChains.SolanaChains()[solChainSel].DeployerKey.PublicKey().String(),
							RegisterType:            ccipChangesetSolana.ViaGetCcipAdminInstruction,
						},
					},
					AcceptAdminRoleTokenAdminRegistry: []ccipChangesetSolana.AcceptAdminRoleTokenAdminRegistryConfig{
						{
							ChainSelector: solChainSel,
							TokenPubKey:   solTokenAddress,
						},
					},
					SetPool: []ccipChangesetSolana.SetPoolConfig{
						{
							ChainSelector:   solChainSel,
							TokenPubKey:     solTokenAddress,
							PoolType:        &bnm,
							Metadata:        shared.CLLMetadata,
							WritableIndexes: []uint8{3, 4, 7},
						},
					},
				},
			),
		)
		if err != nil {
			return *e, fmt.Errorf("failed to apply solana setup link pool changesets: %w", err)
		}

		sourceAccount := *e.BlockChains.SolanaChains()[solChainSel].DeployerKey
		rpcClient := e.BlockChains.SolanaChains()[solChainSel].Client
		router := state.SolChains[solChainSel].Router
		tokenProgram := solana.TokenProgramID
		wSOL := solana.SolMint
		// token transfer enablement changesets
		ixAtaUser, accountWSOLAta, err := soltokens.CreateAssociatedTokenAccount(tokenProgram, wSOL, sourceAccount.PublicKey(), sourceAccount.PublicKey())
		if err != nil {
			return *e, fmt.Errorf("failed to create deployer's wSOL ata: %w", err)
		}

		// Approve CCIP to transfer the user's token for billing
		billingSignerPDA, _, err := solstate.FindFeeBillingSignerPDA(router)
		if err != nil {
			return *e, fmt.Errorf("failed to find billing signer PDA: %w", err)
		}

		ixApproveWSOL, err := soltokens.TokenApproveChecked(math.MaxUint64, 9, tokenProgram, accountWSOLAta, wSOL, billingSignerPDA, sourceAccount.PublicKey(), []solana.PublicKey{})
		if err != nil {
			return *e, fmt.Errorf("failed to create approve instruction: %w", err)
		}

		_, err = solcommon.SendAndConfirm(e.GetContext(), rpcClient, []solana.Instruction{ixAtaUser, ixApproveWSOL}, sourceAccount, solconfig.DefaultCommitment)
		if err != nil {
			return *e, fmt.Errorf("failed to confirm instructions for approving router to spend deployer's wSOL: %w", err)
		}

		// Approve CCIP to transfer the user's Link token for token transfers
		link := state.SolChains[solChainSel].LinkToken
		tokenProgramID, _ := state.SolChains[solChainSel].TokenToTokenProgram(link)
		deployerATA, _, err := soltokens.FindAssociatedTokenAddress(tokenProgramID, link, sourceAccount.PublicKey())
		if err != nil {
			return *e, fmt.Errorf("failed to find associated token address: %w", err)
		}
		ixApproveLink, err := soltokens.TokenApproveChecked(
			tokenApproveCheckedAmount,
			9,
			tokenProgramID,
			deployerATA,
			link,
			billingSignerPDA,
			sourceAccount.PublicKey(),
			[]solana.PublicKey{})
		if err != nil {
			return *e, fmt.Errorf("failed to create approve instruction: %w", err)
		}
		_, err = solcommon.SendAndConfirm(e.GetContext(), rpcClient, []solana.Instruction{ixApproveLink}, sourceAccount, solconfig.DefaultCommitment)
		if err != nil {
			return *e, fmt.Errorf("failed to confirm instructions for approving router to spend deployer's wSOL: %w", err)
		}
	}
	return *e, nil
}

func setupSolEvmLanes(lggr logger.Logger, e *cldf.Environment, state stateview.CCIPOnChainState, homeCS, feedCS uint64) (cldf.Environment, error) {
	var err error
	evmSelectors := e.BlockChains.EVMChains()
	solSelectors := e.BlockChains.SolanaChains()
	g := new(xerrgroup.Group)
	mu := sync.Mutex{}

	for _, solSelector := range solSelectors {
		solSelector := solSelector // capture range variable
		solChainState := state.SolChains[solSelector.ChainSelector()]
		poolUpdates := make(map[uint64]ccipChangesetSolana.EVMRemoteConfig)
		for _, evmSelector := range evmSelectors {
			lggr.Infow("running against evm chain", "evm", evmSelector)
			evmSelector := evmSelector // capture range variables
			g.Go(func() error {
				lggr.Infow("Setting up sol evm lanes for chains", "evmSelector", evmSelector, "solSelector", solSelector)
				laneChangesets := make([]commonchangeset.ConfiguredChangeSet, 0)
				evmChainState := state.Chains[evmSelector.ChainSelector()]

				deployedEnv := testhelpers.DeployedEnv{
					Env:          *e,
					HomeChainSel: homeCS,
					FeedChainSel: feedCS,
				}
				gasPrices := map[uint64]*big.Int{
					solSelector.ChainSelector(): testhelpers.DefaultGasPrice,
				}
				stateChainFrom := evmChainState
				tokenPrices := map[common.Address]*big.Int{
					stateChainFrom.LinkToken.Address(): testhelpers.DefaultLinkPrice,
					stateChainFrom.Weth9.Address():     testhelpers.DefaultWethPrice,
				}
				fqCfg := v1_6.DefaultFeeQuoterDestChainConfig(true, solSelector.ChainSelector())

				mu.Lock()
				poolUpdates[evmSelector.ChainSelector()] = ccipChangesetSolana.EVMRemoteConfig{
					TokenSymbol: shared.LinkSymbol,
					PoolType:    shared.BurnMintTokenPool,
					PoolVersion: shared.CurrentTokenPoolVersion,
					RateLimiterConfig: ccipChangesetSolana.RateLimiterConfig{
						Inbound:  solTestTokenPool.RateLimitConfig{},
						Outbound: solTestTokenPool.RateLimitConfig{},
					},
				}
				mu.Unlock()

				// EVM -> SOL
				cs := testhelpers.AddEVMSrcChangesets(evmSelector.ChainSelector(), solSelector.ChainSelector(), false, gasPrices, tokenPrices, fqCfg)
				laneChangesets = append(laneChangesets, cs...)
				cs = testhelpers.AddLaneSolanaChangesets(&deployedEnv, solSelector.Selector, evmSelector.Selector, chainselectors.FamilyEVM)
				laneChangesets = append(laneChangesets, cs...)

				// SOL -> EVM
				cs = testhelpers.AddEVMDestChangesets(&deployedEnv, evmSelector.Selector, solSelector.Selector, false)
				laneChangesets = append(laneChangesets, cs...)

				bnm := solTestTokenPool.BurnAndMint_PoolType
				laneChangesets = append(laneChangesets,
					commonchangeset.Configure(
						cldf.CreateLegacyChangeSet(ccipChangesetSolana.SetupTokenPoolForRemoteChain),
						ccipChangesetSolana.RemoteChainTokenPoolConfig{
							SolChainSelector: solSelector.Selector,
							SolTokenPubKey:   solChainState.LinkToken,
							SolPoolType:      &bnm,
							EVMRemoteConfigs: map[uint64]ccipChangesetSolana.EVMRemoteConfig{
								evmSelector.Selector: {
									TokenSymbol: shared.LinkSymbol,
									PoolType:    shared.BurnMintTokenPool,
									PoolVersion: shared.CurrentTokenPoolVersion,
									RateLimiterConfig: ccipChangesetSolana.RateLimiterConfig{
										Inbound:  solTestTokenPool.RateLimitConfig{},
										Outbound: solTestTokenPool.RateLimitConfig{},
									},
								},
							},
						},
					),
				)
				lggr.Infow("Applying evm <> svm lane changesets", "len", len(laneChangesets), "evmSel", evmSelector, "svmSel", solSelector)
				_, err = commonchangeset.Apply(nil, *e, laneChangesets[0], laneChangesets[1:]...)
				return err
			})
		}
		err = g.Wait()
		if err != nil {
			return *e, fmt.Errorf("failed to apply sol evm lane changesets: %w", err)
		}
	}
	return *e, nil
}

func mustOCR(e *cldf.Environment, homeChainSel uint64, feedChainSel uint64, newDons bool, rmnEnabled bool) (cldf.Environment, error) {
	chainSelectors := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainselectors.FamilyEVM))
	var commitOCRConfigPerSelector = make(map[uint64]v1_6.CCIPOCRParams)
	var execOCRConfigPerSelector = make(map[uint64]v1_6.CCIPOCRParams)
	// Should be configured in the future based on the load test scenario
	chainType := v1_6.Default

	overrides := func(params v1_6.CCIPOCRParams) v1_6.CCIPOCRParams { return params }
	if rmnEnabled {
		overrides = func(params v1_6.CCIPOCRParams) v1_6.CCIPOCRParams {
			params.CommitOffChainConfig.RMNEnabled = true
			return params
		}
	}

	for selector := range e.BlockChains.EVMChains() {
		commitOCRConfigPerSelector[selector] = v1_6.DeriveOCRParamsForCommit(chainType, feedChainSel, nil, overrides)
		execOCRConfigPerSelector[selector] = v1_6.DeriveOCRParamsForExec(chainType, nil, nil)
	}

	var commitChangeset commonchangeset.ConfiguredChangeSet
	if newDons {
		commitChangeset = commonchangeset.Configure(
			// Add the DONs and candidate commit OCR instances for the chain
			cldf.CreateLegacyChangeSet(v1_6.AddDonAndSetCandidateChangeset),
			v1_6.AddDonAndSetCandidateChangesetConfig{
				SetCandidateConfigBase: v1_6.SetCandidateConfigBase{
					HomeChainSelector: homeChainSel,
					FeedChainSelector: feedChainSel,
				},
				PluginInfo: v1_6.SetCandidatePluginInfo{
					OCRConfigPerRemoteChainSelector: commitOCRConfigPerSelector,
					PluginType:                      types.PluginTypeCCIPCommit,
				},
			},
		)
	} else {
		commitChangeset = commonchangeset.Configure(
			// Update commit OCR instances for existing chains
			cldf.CreateLegacyChangeSet(v1_6.SetCandidateChangeset),
			v1_6.SetCandidateChangesetConfig{
				SetCandidateConfigBase: v1_6.SetCandidateConfigBase{
					HomeChainSelector: homeChainSel,
					FeedChainSelector: feedChainSel,
				},
				PluginInfo: []v1_6.SetCandidatePluginInfo{
					{
						OCRConfigPerRemoteChainSelector: commitOCRConfigPerSelector,
						PluginType:                      types.PluginTypeCCIPCommit,
					},
				},
			},
		)
	}

	return commonchangeset.Apply(nil, *e, commitChangeset, commonchangeset.Configure(
		// Add the exec OCR instances for the new chains
		cldf.CreateLegacyChangeSet(v1_6.SetCandidateChangeset),
		v1_6.SetCandidateChangesetConfig{
			SetCandidateConfigBase: v1_6.SetCandidateConfigBase{
				HomeChainSelector: homeChainSel,
				FeedChainSelector: feedChainSel,
			},
			PluginInfo: []v1_6.SetCandidatePluginInfo{
				{
					OCRConfigPerRemoteChainSelector: execOCRConfigPerSelector,
					PluginType:                      types.PluginTypeCCIPExec,
				},
			},
		},
	), commonchangeset.Configure(
		// Promote everything
		cldf.CreateLegacyChangeSet(v1_6.PromoteCandidateChangeset),
		v1_6.PromoteCandidateChangesetConfig{
			HomeChainSelector: homeChainSel,
			PluginInfo: []v1_6.PromoteCandidatePluginInfo{
				{
					RemoteChainSelectors: chainSelectors,
					PluginType:           types.PluginTypeCCIPCommit,
				},
				{
					RemoteChainSelectors: chainSelectors,
					PluginType:           types.PluginTypeCCIPExec,
				},
			},
		},
	), commonchangeset.Configure(
		// Enable the OCR config on the remote chains
		cldf.CreateLegacyChangeSet(v1_6.SetOCR3OffRampChangeset),
		v1_6.SetOCR3OffRampConfig{
			HomeChainSel:       homeChainSel,
			RemoteChainSels:    chainSelectors,
			CCIPHomeConfigType: globals.ConfigTypeActive,
		},
	))
}

type RMNNodeConfig struct {
	v1_6.RMNNopConfig
	RageProxyKeystore string
	RMNKeystore       string
	Passphrase        string
}

func SetupRMNNodeOnAllChains(ctx context.Context, lggr logger.Logger, envConfig devenv.EnvironmentConfig, homeChainSel, feedChainSel uint64, ab cldf.AddressBook, nodes []RMNNodeConfig) (DeployCCIPOutput, error) {
	e, _, err := devenv.NewEnvironment(func() context.Context { return ctx }, lggr, envConfig)
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to create environment: %w", err)
	}

	e.ExistingAddresses = ab

	allChains := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainselectors.FamilyEVM))
	allUpdates := make(map[uint64]map[uint64]v1_6.OffRampSourceUpdate)
	for _, chainIdx := range allChains {
		updates := make(map[uint64]v1_6.OffRampSourceUpdate)

		for _, subChainID := range allChains {
			if subChainID == chainIdx {
				continue
			}
			updates[subChainID] = v1_6.OffRampSourceUpdate{
				IsRMNVerificationDisabled: false,
				IsEnabled:                 true,
			}
		}

		allUpdates[chainIdx] = updates
	}

	_, err = v1_6.UpdateOffRampSourcesChangeset(*e, v1_6.UpdateOffRampSourcesConfig{
		UpdatesByChain: allUpdates,
	})

	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to update dynamic off ramp config: %w", err)
	}

	rmnNodes := make([]rmn_home.RMNHomeNode, len(nodes))
	bitmap := new(big.Int)
	for i, node := range nodes {
		rmnNodes[i] = rmn_home.RMNHomeNode{
			PeerId:            node.PeerId,
			OffchainPublicKey: node.OffchainPublicKey,
		}
		bitmap.SetBit(bitmap, i, 1)
	}

	sourceChains := make([]rmn_home.RMNHomeSourceChain, len(allChains))
	for i, chain := range allChains {
		sourceChains[i] = rmn_home.RMNHomeSourceChain{
			ChainSelector:       chain,
			FObserve:            1,
			ObserverNodesBitmap: bitmap,
		}
	}

	env, err := commonchangeset.Apply(nil, *e,
		commonchangeset.Configure(
			// Enable the OCR config on the remote chains
			cldf.CreateLegacyChangeSet(v1_6.SetRMNHomeCandidateConfigChangeset),
			v1_6.SetRMNHomeCandidateConfig{
				HomeChainSelector: homeChainSel,
				RMNStaticConfig: rmn_home.RMNHomeStaticConfig{
					Nodes:          rmnNodes,
					OffchainConfig: []byte{},
				},
				RMNDynamicConfig: rmn_home.RMNHomeDynamicConfig{
					OffchainConfig: []byte{},
					SourceChains:   sourceChains,
				},
			},
		),
	)
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to set rmn node candidate: %w", err)
	}

	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to load chain state: %w", err)
	}

	configDigest, err := state.Chains[homeChainSel].RMNHome.GetCandidateDigest(nil)

	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to get rmn home candidate digest: %w", err)
	}

	env, err = commonchangeset.Apply(nil, *e,
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(v1_6.PromoteRMNHomeCandidateConfigChangeset),
			v1_6.PromoteRMNHomeCandidateConfig{
				HomeChainSelector: homeChainSel,
				DigestToPromote:   configDigest,
			},
		),
	)
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to promote rmn node candidate: %w", err)
	}

	signers := make([]rmn_remote.RMNRemoteSigner, len(nodes))
	for i, node := range nodes {
		signers[i] = node.ToRMNRemoteSigner()
	}

	g, ctx := xerrgroup.WithContext(context.Background())
	for _, chain := range allChains {
		g.Go(func() error {
			rmnRemoteConfig := map[uint64]ccipops.RMNRemoteConfig{
				chain: {
					Signers: signers,
					F:       1,
				},
			}

			_, err := v1_6.SetRMNRemoteConfigChangeset(*e, ccipseq.SetRMNRemoteConfig{
				RMNRemoteConfigs: rmnRemoteConfig,
			})
			return err
		})
	}
	if err := g.Wait(); err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to set rmn remote config: %w", err)
	}

	addresses, err := env.ExistingAddresses.Addresses()
	if err != nil {
		return DeployCCIPOutput{}, fmt.Errorf("failed to get existing addresses: %w", err)
	}
	return DeployCCIPOutput{
		AddressBook: *cldf.NewMemoryAddressBookFromMap(addresses),
		NodeIDs:     e.NodeIDs,
	}, nil
}

func GenerateRMNNodeIdentities(rmnNodeCount uint, rageProxyImageURI, rageProxyImageTag, afn2proxyImageURI,
	afn2proxyImageTag string, imagePlatform string) ([]RMNNodeConfig, error) {
	lggr := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout})
	rmnNodeConfigs := make([]RMNNodeConfig, rmnNodeCount)

	for i := uint(0); i < rmnNodeCount; i++ {
		peerID, rawKeystore, _, err := devenv.GeneratePeerID(zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}), rageProxyImageURI, rageProxyImageTag, imagePlatform)
		if err != nil {
			return nil, err
		}

		keys, rawRMNKeystore, afnPassphrase, err := devenv.GenerateRMNKeyStore(lggr, afn2proxyImageURI, afn2proxyImageTag, imagePlatform)
		if err != nil {
			return nil, err
		}

		newPeerID, err := p2pkey.MakePeerID(peerID.String())
		if err != nil {
			return nil, err
		}

		rmnNodeConfigs[i] = RMNNodeConfig{
			RMNNopConfig: v1_6.RMNNopConfig{
				NodeIndex:           uint64(i),
				OffchainPublicKey:   [32]byte(keys.OffchainPublicKey),
				EVMOnChainPublicKey: keys.EVMOnchainPublicKey,
				PeerId:              newPeerID,
			},
			RageProxyKeystore: rawKeystore,
			RMNKeystore:       rawRMNKeystore,
			Passphrase:        afnPassphrase,
		}
	}
	return rmnNodeConfigs, nil
}
