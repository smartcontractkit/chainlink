package crib

import (
	"fmt"
	"math"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	xerrgroup "golang.org/x/sync/errgroup"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-ccip/chainconfig"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/token_pool"
	evm_fee_quoter "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/fee_quoter"
	solconfig "github.com/smartcontractkit/chainlink-ccip/chains/solana/contracts/tests/config"
	solTestTokenPool "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_0/test_token_pool"
	solcommon "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	solstate "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"
	soltokens "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
	ccipChangesetSolana "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/solana_v0_1_0"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_5_1"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	ccipops "github.com/smartcontractkit/chainlink/deployment/ccip/operation/evm/v1_6"
	ccipseq "github.com/smartcontractkit/chainlink/deployment/ccip/sequence/evm/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"

	"github.com/smartcontractkit/chainlink/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
)

const (
	tokenApproveCheckedAmount = 1e4 * 1e9
)

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
			GitCommitSha:   "6aaf88e0848a",
			DestinationDir: deployedEnv.Env.BlockChains.SolanaChains()[solChainSelectors[0]].ProgramsPath,
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
					AddTokenPoolAndLookupTable: []ccipChangesetSolana.AddTokenPoolAndLookupTableConfig{
						{
							ChainSelector: solChainSel,
							TokenPoolConfigs: []ccipChangesetSolana.TokenPoolConfig{
								{
									TokenPubKey: solTokenAddress,
									PoolType:    shared.BurnMintTokenPool,
									Metadata:    shared.CLLMetadata,
								},
							},
						},
					},
					RegisterTokenAdminRegistry: []ccipChangesetSolana.RegisterTokenAdminRegistryConfig{
						{
							ChainSelector: solChainSel,
							RegisterTokenConfigs: []ccipChangesetSolana.RegisterTokenConfig{
								{
									TokenPubKey:             solTokenAddress,
									TokenAdminRegistryAdmin: e.BlockChains.SolanaChains()[solChainSel].DeployerKey.PublicKey(),
									RegisterType:            ccipChangesetSolana.ViaGetCcipAdminInstruction,
								},
							},
						},
					},
					AcceptAdminRoleTokenAdminRegistry: []ccipChangesetSolana.AcceptAdminRoleTokenAdminRegistryConfig{
						{
							ChainSelector: solChainSel,
							AcceptAdminRoleTokenConfigs: []ccipChangesetSolana.AcceptAdminRoleTokenConfig{
								{
									TokenPubKey: solTokenAddress,
								},
							},
						},
					},
					SetPool: []ccipChangesetSolana.SetPoolConfig{
						{
							ChainSelector:   solChainSel,
							WritableIndexes: []uint8{3, 4, 7},
							SetPoolTokenConfigs: []ccipChangesetSolana.SetPoolTokenConfig{
								{
									TokenPubKey: solTokenAddress,
									PoolType:    shared.BurnMintTokenPool,
									Metadata:    shared.CLLMetadata,
								},
							},
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

func hasLaneFromTo(lanes []LaneConfig, from, to uint64) bool {
	for _, lane := range lanes {
		if lane.SourceChain == from && lane.DestinationChain == to {
			return true
		}
	}
	return false
}

func setupSolEvmLanes(lggr logger.Logger, e *cldf.Environment, state stateview.CCIPOnChainState, homeCS, feedCS uint64, laneConfig *LaneConfiguration) (cldf.Environment, error) {
	var err error

	lanes, err := laneConfig.GetLanes()
	if err != nil {
		return *e, fmt.Errorf("failed to get lanes from lane configuration: %w", err)
	}

	evmSelectors := e.BlockChains.EVMChains()
	solSelectors := e.BlockChains.SolanaChains()
	g := new(xerrgroup.Group)
	mu := sync.Mutex{}

	// Filter lanes to only include Sol <-> EVM combinations
	evmChainSet := make(map[uint64]bool)
	solChainSet := make(map[uint64]bool)

	for _, evmSelector := range evmSelectors {
		evmChainSet[evmSelector.ChainSelector()] = true
	}
	for _, solSelector := range solSelectors {
		solChainSet[solSelector.ChainSelector()] = true
	}

	lanesBySolChain := make(map[uint64][]LaneConfig)
	for _, lane := range lanes {
		if solChainSet[lane.SourceChain] && evmChainSet[lane.DestinationChain] {
			lanesBySolChain[lane.SourceChain] = append(lanesBySolChain[lane.SourceChain], lane)
		}

		if evmChainSet[lane.SourceChain] && solChainSet[lane.DestinationChain] {
			lanesBySolChain[lane.DestinationChain] = append(lanesBySolChain[lane.DestinationChain], lane)
		}
	}

	for _, solSelector := range solSelectors {
		solSelector := solSelector // capture range variable
		solChainSel := solSelector.ChainSelector()
		relevantLanes := lanesBySolChain[solChainSel]

		if len(relevantLanes) == 0 {
			continue // Skip if no lanes involve this Solana chain
		}

		solChainState := state.SolChains[solChainSel]
		poolUpdates := make(map[uint64]ccipChangesetSolana.EVMRemoteConfig)

		for _, evmSelector := range evmSelectors {
			evmChainSel := evmSelector.ChainSelector()

			// Check if there's a lane between this Sol and EVM chain
			hasLane := false
			for _, lane := range relevantLanes {
				if (lane.SourceChain == solChainSel && lane.DestinationChain == evmChainSel) ||
					(lane.SourceChain == evmChainSel && lane.DestinationChain == solChainSel) {
					hasLane = true
					break
				}
			}

			if !hasLane {
				continue // Skip if no lane exists between these chains
			}

			lggr.Infow("running against evm chain", "evm", evmChainSel)
			evmSelector := evmSelector
			g.Go(func() error {
				lggr.Infow("Setting up sol evm lanes for chains", "evmSelector", evmChainSel, "solSelector", solChainSel)
				laneChangesets := make([]commonchangeset.ConfiguredChangeSet, 0)
				evmChainState := state.Chains[evmChainSel]

				deployedEnv := testhelpers.DeployedEnv{
					Env:          *e,
					HomeChainSel: homeCS,
					FeedChainSel: feedCS,
				}
				gasPrices := map[uint64]*big.Int{
					solChainSel: testhelpers.DefaultGasPrice,
				}
				stateChainFrom := evmChainState
				tokenPrices := map[common.Address]*big.Int{
					stateChainFrom.LinkToken.Address(): testhelpers.DefaultLinkPrice,
					stateChainFrom.Weth9.Address():     testhelpers.DefaultWethPrice,
				}
				fqCfg := v1_6.DefaultFeeQuoterDestChainConfig(true, solChainSel)

				mu.Lock()
				poolUpdates[evmChainSel] = ccipChangesetSolana.EVMRemoteConfig{
					TokenSymbol: shared.LinkSymbol,
					PoolType:    shared.BurnMintTokenPool,
					PoolVersion: shared.CurrentTokenPoolVersion,
					RateLimiterConfig: ccipChangesetSolana.RateLimiterConfig{
						Inbound:  solTestTokenPool.RateLimitConfig{},
						Outbound: solTestTokenPool.RateLimitConfig{},
					},
				}
				mu.Unlock()

				// TODO: Maybe use maps to make it more efficient (for the n chains/lanes we use now it doesn't really
				//  matter
				// EVM -> SOL (only if lane exists)
				if hasLaneFromTo(relevantLanes, evmChainSel, solChainSel) {
					cs := testhelpers.AddEVMSrcChangesets(evmChainSel, solChainSel, false, gasPrices, tokenPrices, fqCfg)
					laneChangesets = append(laneChangesets, cs...)
					cs = testhelpers.AddLaneSolanaChangesetsV0_1_1(&deployedEnv, solSelector.Selector, evmSelector.Selector, chainselectors.FamilyEVM)
					laneChangesets = append(laneChangesets, cs...)
				}

				// SOL -> EVM (only if lane exists)
				if hasLaneFromTo(relevantLanes, solChainSel, evmChainSel) {
					cs := testhelpers.AddEVMDestChangesets(&deployedEnv, evmSelector.Selector, solSelector.Selector, false)
					laneChangesets = append(laneChangesets, cs...)
				}

				laneChangesets = append(laneChangesets,
					commonchangeset.Configure(
						cldf.CreateLegacyChangeSet(ccipChangesetSolana.SetupTokenPoolForRemoteChain),
						ccipChangesetSolana.SetupTokenPoolForRemoteChainConfig{
							SolChainSelector: solSelector.Selector,
							RemoteTokenPoolConfigs: []ccipChangesetSolana.RemoteChainTokenPoolConfig{
								{
									SolTokenPubKey: solChainState.LinkToken,
									SolPoolType:    shared.BurnMintTokenPool,
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
							},
						},
					),
				)
				lggr.Infow("Applying evm <> svm lane changesets", "len", len(laneChangesets), "evmSel", evmChainSel, "svmSel", solChainSel)
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

func setupEVM2EVMLanes(e *cldf.Environment, state stateview.CCIPOnChainState, laneConfig *LaneConfiguration) (cldf.Environment, error) {
	lanes, err := laneConfig.GetLanes()
	if err != nil {
		return *e, fmt.Errorf("failed to get lanes from config: %w", err)
	}

	poolUpdates := make(map[uint64]v1_5_1.TokenPoolConfig)
	rateLimitPerChain := make(v1_5_1.RateLimiterPerChain)
	evmChains := e.BlockChains.EVMChains()

	// Filter to only include EVM chains
	evmLanes := make([]LaneConfig, 0)

	for _, lane := range lanes {
		_, srcExists := evmChains[lane.SourceChain]
		_, dstExists := evmChains[lane.DestinationChain]
		if srcExists && dstExists {
			evmLanes = append(evmLanes, lane)
		}
	}

	eg := new(xerrgroup.Group)
	mu := sync.Mutex{}

	globalUpdateOffRampSources := make(map[uint64]map[uint64]v1_6.OffRampSourceUpdate)
	globalUpdateRouterChanges := make(map[uint64]v1_6.RouterUpdates)

	// Initialize maps for all chains that will be used
	for chainID := range evmChains {
		globalUpdateOffRampSources[chainID] = make(map[uint64]v1_6.OffRampSourceUpdate)
		globalUpdateRouterChanges[chainID] = v1_6.RouterUpdates{
			OffRampUpdates: make(map[uint64]bool),
			OnRampUpdates:  make(map[uint64]bool),
		}
	}

	// Group lanes by source chain for parallel processing
	lanesBySource := make(map[uint64][]LaneConfig)
	for _, lane := range evmLanes {
		lanesBySource[lane.SourceChain] = append(lanesBySource[lane.SourceChain], lane)
	}

	for src := range evmChains {
		src := src
		lanesFromSrc := lanesBySource[src]
		if len(lanesFromSrc) == 0 {
			continue // Skip chains that don't have any outgoing lanes
		}

		eg.Go(func() error {
			onRampUpdatesByChain := make(map[uint64]map[uint64]v1_6.OnRampDestinationUpdate)
			pricesByChain := make(map[uint64]v1_6.FeeQuoterPriceUpdatePerSource)
			feeQuoterDestsUpdatesByChain := make(map[uint64]map[uint64]evm_fee_quoter.FeeQuoterDestChainConfig)

			onRampUpdatesByChain[src] = make(map[uint64]v1_6.OnRampDestinationUpdate)
			pricesByChain[src] = v1_6.FeeQuoterPriceUpdatePerSource{
				TokenPrices: map[common.Address]*big.Int{
					state.Chains[src].LinkToken.Address(): testhelpers.DefaultLinkPrice,
					state.Chains[src].Weth9.Address():     testhelpers.DefaultWethPrice,
				},
				GasPrices: map[uint64]*big.Int{},
			}
			feeQuoterDestsUpdatesByChain[src] = make(map[uint64]evm_fee_quoter.FeeQuoterDestChainConfig)

			mu.Lock()
			poolUpdates[src] = v1_5_1.TokenPoolConfig{
				Type:         shared.BurnMintTokenPool,
				Version:      deployment.Version1_5_1,
				ChainUpdates: rateLimitPerChain,
			}
			mu.Unlock()

			// Only configure lanes that actually exist in our configuration
			for _, lane := range lanesFromSrc {
				dst := lane.DestinationChain

				onRampUpdatesByChain[src][dst] = v1_6.OnRampDestinationUpdate{
					IsEnabled:        true,
					AllowListEnabled: false,
				}
				pricesByChain[src].GasPrices[dst] = testhelpers.DefaultGasPrice
				feeQuoterDestsUpdatesByChain[src][dst] = v1_6.DefaultFeeQuoterDestChainConfig(true)

				mu.Lock()
				// Use the pre-initialized global maps
				globalUpdateOffRampSources[dst][src] = v1_6.OffRampSourceUpdate{
					IsEnabled:                 true,
					IsRMNVerificationDisabled: true,
				}

				globalUpdateRouterChanges[dst].OffRampUpdates[src] = true
				globalUpdateRouterChanges[src].OnRampUpdates[dst] = true

				rateLimitPerChain[dst] = v1_5_1.RateLimiterConfig{
					Inbound: token_pool.RateLimiterConfig{
						IsEnabled: false,
						Capacity:  big.NewInt(0),
						Rate:      big.NewInt(0),
					},
					Outbound: token_pool.RateLimiterConfig{
						IsEnabled: false,
						Capacity:  big.NewInt(0),
						Rate:      big.NewInt(0),
					},
				}
				mu.Unlock()
			}

			appliedChangesets := []commonchangeset.ConfiguredChangeSet{
				commonchangeset.Configure(
					cldf.CreateLegacyChangeSet(v1_6.UpdateOnRampsDestsChangeset),
					v1_6.UpdateOnRampDestsConfig{
						UpdatesByChain: onRampUpdatesByChain,
					},
				),
				commonchangeset.Configure(
					cldf.CreateLegacyChangeSet(v1_6.UpdateFeeQuoterPricesChangeset),
					v1_6.UpdateFeeQuoterPricesConfig{
						PricesByChain: pricesByChain,
					},
				),
				commonchangeset.Configure(
					cldf.CreateLegacyChangeSet(v1_6.UpdateFeeQuoterDestsChangeset),
					v1_6.UpdateFeeQuoterDestsConfig{
						UpdatesByChain: feeQuoterDestsUpdatesByChain,
					},
				),
			}
			_, err := commonchangeset.Apply(nil, *e, appliedChangesets[0], appliedChangesets[1:]...)
			return err
		})
	}

	err = eg.Wait()
	if err != nil {
		return *e, err
	}

	// Apply the global updates after all goroutines complete
	finalChangesets := []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(v1_6.UpdateOffRampSourcesChangeset),
			v1_6.UpdateOffRampSourcesConfig{
				UpdatesByChain: globalUpdateOffRampSources,
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(v1_6.UpdateRouterRampsChangeset),
			v1_6.UpdateRouterRampsConfig{
				UpdatesByChain: globalUpdateRouterChanges,
			},
		),
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(v1_5_1.ConfigureTokenPoolContractsChangeset),
			v1_5_1.ConfigureTokenPoolContractsConfig{
				TokenSymbol: shared.LinkSymbol,
				PoolUpdates: poolUpdates,
			},
		),
	}

	_, err = commonchangeset.Apply(nil, *e, finalChangesets[0], finalChangesets[1:]...)
	return *e, err
}

func mustOCR(e *cldf.Environment, homeChainSel uint64, feedChainSel uint64, newDons bool, rmnEnabled bool) (cldf.Environment, error) {
	evmSelectors := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainselectors.FamilyEVM))
	solSelectors := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainselectors.FamilySolana))
	// need to have extra definition here for golint
	var allSelectors = make([]uint64, 0)
	allSelectors = append(allSelectors, evmSelectors...)
	allSelectors = append(allSelectors, solSelectors...)
	var commitOCRConfigPerSelector = make(map[uint64]v1_6.CCIPOCRParams)
	var execOCRConfigPerSelector = make(map[uint64]v1_6.CCIPOCRParams)
	// Should be configured in the future based on the load test scenario
	chainType := v1_6.Default
	_, err := testhelpers.DeployFeeds(e.Logger, e.ExistingAddresses, e.BlockChains.EVMChains()[feedChainSel], testhelpers.DefaultLinkPrice, testhelpers.DefaultWethPrice)
	if err != nil {
		return *e, fmt.Errorf("failed to deploy feeds: %w", err)
	}
	state, err := stateview.LoadOnchainState(*e)
	if err != nil {
		return *e, fmt.Errorf("failed to load onchain state: %w", err)
	}

	overrides := func(params v1_6.CCIPOCRParams) v1_6.CCIPOCRParams { return params }

	tokenConfig := shared.NewTestTokenConfig(state.Chains[feedChainSel].USDFeeds)
	var tokenDataProviders []pluginconfig.TokenDataObserverConfig

	for _, selector := range evmSelectors {
		commitOCRConfigPerSelector[selector] = v1_6.DeriveOCRParamsForCommit(chainType, feedChainSel, nil, overrides)
		execOCRConfigPerSelector[selector] = v1_6.DeriveOCRParamsForExec(chainType, tokenDataProviders, nil)
	}

	for _, selector := range solSelectors {
		// TODO: this is a workaround for tokenConfig.GetTokenInfo
		tokenInfo := map[cciptypes.UnknownEncodedAddress]pluginconfig.TokenInfo{}
		tokenInfo[cciptypes.UnknownEncodedAddress(state.SolChains[selector].LinkToken.String())] = tokenConfig.TokenSymbolToInfo[shared.LinkSymbol]
		// TODO: point this to proper SOL feed, apparently 0 signified SOL
		tokenInfo[cciptypes.UnknownEncodedAddress(solana.SolMint.String())] = tokenConfig.TokenSymbolToInfo[shared.WethSymbol]
		commitOCRConfigPerSelector[selector] = v1_6.DeriveOCRParamsForCommit(chainType, feedChainSel, tokenInfo,
			func(params v1_6.CCIPOCRParams) v1_6.CCIPOCRParams {
				params.OCRParameters.MaxDurationQuery = 100 * time.Millisecond
				params.OCRParameters.DeltaRound = 5 * time.Second
				params.CommitOffChainConfig.RMNEnabled = false
				params.CommitOffChainConfig.RMNSignaturesTimeout = 100 * time.Millisecond
				params.CommitOffChainConfig.MultipleReportsEnabled = true
				params.CommitOffChainConfig.MaxMerkleRootsPerReport = 1
				params.CommitOffChainConfig.MaxPricesPerReport = 3
				params.CommitOffChainConfig.MaxMerkleTreeSize = 1
				params.CommitOffChainConfig.MerkleRootAsyncObserverDisabled = true
				return params
			})
		execOCRConfigPerSelector[selector] = v1_6.DeriveOCRParamsForExec(chainType, tokenDataProviders,
			func(params v1_6.CCIPOCRParams) v1_6.CCIPOCRParams {
				params.ExecuteOffChainConfig.MaxSingleChainReports = 1
				params.ExecuteOffChainConfig.BatchGasLimit = 1000000
				params.ExecuteOffChainConfig.MaxReportMessages = 1
				params.ExecuteOffChainConfig.MultipleReportsEnabled = true

				return params
			})
		commitOCRConfigPerSelector[selector].CommitOffChainConfig.RMNEnabled = false
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
					RemoteChainSelectors: allSelectors,
					PluginType:           types.PluginTypeCCIPCommit,
				},
				{
					RemoteChainSelectors: allSelectors,
					PluginType:           types.PluginTypeCCIPExec,
				},
			},
		},
	), commonchangeset.Configure(
		// Enable the OCR config on the remote chains
		cldf.CreateLegacyChangeSet(v1_6.SetOCR3OffRampChangeset),
		v1_6.SetOCR3OffRampConfig{
			HomeChainSel:       homeChainSel,
			RemoteChainSels:    evmSelectors,
			CCIPHomeConfigType: globals.ConfigTypeActive,
		},
	), commonchangeset.Configure(
		// Enable the OCR config on the remote chains.
		cldf.CreateLegacyChangeSet(ccipChangesetSolana.SetOCR3ConfigSolana),
		v1_6.SetOCR3OffRampConfig{
			HomeChainSel:       homeChainSel,
			RemoteChainSels:    solSelectors,
			CCIPHomeConfigType: globals.ConfigTypeActive,
		},
	),
	)
}

type RMNNodeConfig struct {
	v1_6.RMNNopConfig
	RageProxyKeystore string
	RMNKeystore       string
	Passphrase        string
}
