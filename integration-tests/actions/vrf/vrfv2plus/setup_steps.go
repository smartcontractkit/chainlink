package vrfv2plus

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/shopspring/decimal"
	"golang.org/x/sync/errgroup"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-testing-framework/utils/conversions"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"
	actions_seth "github.com/smartcontractkit/chainlink/integration-tests/actions/seth"
	vrfcommon "github.com/smartcontractkit/chainlink/integration-tests/actions/vrf/common"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
	"github.com/smartcontractkit/chainlink/integration-tests/types/config/node"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	vrfv2plus_config "github.com/smartcontractkit/chainlink/integration-tests/testconfig/vrfv2plus"
	"github.com/smartcontractkit/chainlink/integration-tests/types"
)

func CreateVRFV2PlusJob(
	chainlinkNode *client.ChainlinkClient,
	vrfJobSpecConfig vrfcommon.VRFJobSpecConfig,
) (*client.Job, error) {
	jobUUID := uuid.New()
	os := &client.VRFV2PlusTxPipelineSpec{
		Address:               vrfJobSpecConfig.CoordinatorAddress,
		EstimateGasMultiplier: vrfJobSpecConfig.EstimateGasMultiplier,
		FromAddress:           vrfJobSpecConfig.FromAddresses[0],
		SimulationBlock:       vrfJobSpecConfig.SimulationBlock,
	}
	ost, err := os.String()
	if err != nil {
		return nil, fmt.Errorf(vrfcommon.ErrGenericFormat, vrfcommon.ErrParseJob, err)
	}

	jobSpec := client.VRFV2PlusJobSpec{
		Name:                          fmt.Sprintf("vrf-v2-plus-%s", jobUUID),
		CoordinatorAddress:            vrfJobSpecConfig.CoordinatorAddress,
		FromAddresses:                 vrfJobSpecConfig.FromAddresses,
		EVMChainID:                    vrfJobSpecConfig.EVMChainID,
		MinIncomingConfirmations:      vrfJobSpecConfig.MinIncomingConfirmations,
		PublicKey:                     vrfJobSpecConfig.PublicKey,
		ExternalJobID:                 jobUUID.String(),
		ObservationSource:             ost,
		BatchFulfillmentEnabled:       vrfJobSpecConfig.BatchFulfillmentEnabled,
		BatchFulfillmentGasMultiplier: vrfJobSpecConfig.BatchFulfillmentGasMultiplier,
		PollPeriod:                    vrfJobSpecConfig.PollPeriod,
		RequestTimeout:                vrfJobSpecConfig.RequestTimeout,
	}

	if vrfJobSpecConfig.BatchFulfillmentEnabled {
		jobSpec.BatchCoordinatorAddress = vrfJobSpecConfig.BatchCoordinatorAddress
	}

	job, err := chainlinkNode.MustCreateJob(&jobSpec)
	if err != nil {
		return nil, fmt.Errorf(vrfcommon.ErrGenericFormat, ErrCreatingVRFv2PlusJob, err)
	}

	return job, nil
}

// SetupVRFV2_5Environment will create specified number of subscriptions and add the same conumer/s to each of them
func SetupVRFV2_5Environment(
	ctx context.Context,
	env *test_env.CLClusterTestEnv,
	chainID int64,
	nodesToCreate []vrfcommon.VRFNodeType,
	vrfv2PlusTestConfig types.VRFv2PlusTestConfig,
	linkToken contracts.LinkToken,
	mockNativeLINKFeed contracts.VRFMockETHLINKFeed,
	numberOfTxKeysToCreate int,
	l zerolog.Logger,
) (*vrfcommon.VRFContracts, *vrfcommon.VRFKeyData, map[vrfcommon.VRFNodeType]*vrfcommon.VRFNode, error) {
	l.Info().Msg("Starting VRFV2 Plus environment setup")
	configGeneral := vrfv2PlusTestConfig.GetVRFv2PlusConfig().General
	vrfContracts, err := SetupVRFV2PlusContracts(
		env,
		chainID,
		linkToken,
		mockNativeLINKFeed,
		configGeneral,
		l,
	)
	if err != nil {
		return nil, nil, nil, err
	}

	nodeTypeToNodeMap, err := vrfcommon.CreateNodeTypeToNodeMap(env.ClCluster, nodesToCreate)
	if err != nil {
		return nil, nil, nil, err
	}
	vrfKey, pubKeyCompressed, err := vrfcommon.CreateVRFKeyOnVRFNode(nodeTypeToNodeMap[vrfcommon.VRF], l)
	if err != nil {
		return nil, nil, nil, err
	}
	l.Info().Str("Coordinator", vrfContracts.CoordinatorV2Plus.Address()).Msg("Registering Proving Key")
	provingKey, err := VRFV2_5RegisterProvingKey(vrfKey, vrfContracts.CoordinatorV2Plus, uint64(assets.GWei(*configGeneral.CLNodeMaxGasPriceGWei).Int64()))
	if err != nil {
		return nil, nil, nil, fmt.Errorf(vrfcommon.ErrGenericFormat, vrfcommon.ErrRegisteringProvingKey, err)
	}
	keyHash, err := vrfContracts.CoordinatorV2Plus.HashOfKey(ctx, provingKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf(vrfcommon.ErrGenericFormat, vrfcommon.ErrCreatingProvingKeyHash, err)
	}

	sethClient, err := env.GetSethClient(chainID)
	if err != nil {
		return nil, nil, nil, err
	}

	vrfTXKeyAddressStrings, _, err := vrfcommon.CreateFundAndGetSendingKeys(
		l,
		sethClient,
		nodeTypeToNodeMap[vrfcommon.VRF],
		*vrfv2PlusTestConfig.GetCommonConfig().ChainlinkNodeFunding,
		numberOfTxKeysToCreate,
		big.NewInt(chainID),
	)
	if err != nil {
		return nil, nil, nil, err
	}

	nodeTypeToNodeMap[vrfcommon.VRF].TXKeyAddressStrings = vrfTXKeyAddressStrings

	g := errgroup.Group{}
	if vrfNode, exists := nodeTypeToNodeMap[vrfcommon.VRF]; exists {
		g.Go(func() error {
			err := setupVRFNode(vrfContracts, big.NewInt(chainID), configGeneral, pubKeyCompressed, l, vrfNode)
			if err != nil {
				return err
			}
			return nil
		})
	}

	if bhsNode, exists := nodeTypeToNodeMap[vrfcommon.BHS]; exists {
		g.Go(func() error {
			err := vrfcommon.SetupBHSNode(
				env,
				configGeneral.General,
				numberOfTxKeysToCreate,
				big.NewInt(chainID),
				vrfContracts.CoordinatorV2Plus.Address(),
				vrfContracts.BHS.Address(),
				*vrfv2PlusTestConfig.GetCommonConfig().ChainlinkNodeFunding,
				l,
				bhsNode,
			)
			if err != nil {
				return err
			}
			return nil
		})
	}

	if bhfNode, exists := nodeTypeToNodeMap[vrfcommon.BHF]; exists {
		g.Go(func() error {
			err := vrfcommon.SetupBHFNode(
				env,
				configGeneral.General,
				numberOfTxKeysToCreate,
				big.NewInt(chainID),
				vrfContracts.CoordinatorV2Plus.Address(),
				vrfContracts.BHS.Address(),
				vrfContracts.BatchBHS.Address(),
				*vrfv2PlusTestConfig.GetCommonConfig().ChainlinkNodeFunding,
				l,
				bhfNode,
			)
			if err != nil {
				return err
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, nil, nil, fmt.Errorf("VRF node setup ended up with an error: %w", err)
	}

	vrfKeyData := vrfcommon.VRFKeyData{
		VRFKey:            vrfKey,
		EncodedProvingKey: provingKey,
		KeyHash:           keyHash,
		PubKeyCompressed:  pubKeyCompressed,
	}

	l.Info().Msg("VRFV2 Plus environment setup is finished")
	return vrfContracts, &vrfKeyData, nodeTypeToNodeMap, nil
}

func setupVRFNode(contracts *vrfcommon.VRFContracts, chainID *big.Int, config *vrfv2plus_config.General, pubKeyCompressed string, l zerolog.Logger, vrfNode *vrfcommon.VRFNode) error {
	vrfJobSpecConfig := vrfcommon.VRFJobSpecConfig{
		ForwardingAllowed:             *config.VRFJobForwardingAllowed,
		CoordinatorAddress:            contracts.CoordinatorV2Plus.Address(),
		BatchCoordinatorAddress:       contracts.BatchCoordinatorV2Plus.Address(),
		FromAddresses:                 vrfNode.TXKeyAddressStrings,
		EVMChainID:                    chainID.String(),
		MinIncomingConfirmations:      int(*config.MinimumConfirmations),
		PublicKey:                     pubKeyCompressed,
		EstimateGasMultiplier:         *config.VRFJobEstimateGasMultiplier,
		BatchFulfillmentEnabled:       *config.VRFJobBatchFulfillmentEnabled,
		BatchFulfillmentGasMultiplier: *config.VRFJobBatchFulfillmentGasMultiplier,
		PollPeriod:                    config.VRFJobPollPeriod.Duration,
		RequestTimeout:                config.VRFJobRequestTimeout.Duration,
		SimulationBlock:               config.VRFJobSimulationBlock,
		VRFOwnerConfig:                nil,
	}

	l.Info().Msg("Creating VRFV2 Plus Job")
	job, err := CreateVRFV2PlusJob(
		vrfNode.CLNode.API,
		vrfJobSpecConfig,
	)
	if err != nil {
		return fmt.Errorf(vrfcommon.ErrGenericFormat, ErrCreateVRFV2PlusJobs, err)
	}
	vrfNode.Job = job

	// this part is here because VRFv2 can work with only a specific key
	// [[EVM.KeySpecific]]
	//	Key = '...'
	nodeConfig := node.NewConfig(vrfNode.CLNode.NodeConfig,
		node.WithKeySpecificMaxGasPrice(vrfNode.TXKeyAddressStrings, *config.CLNodeMaxGasPriceGWei),
	)
	l.Info().Msg("Restarting Node with new sending key PriceMax configuration")
	err = vrfNode.CLNode.Restart(nodeConfig)
	if err != nil {
		return fmt.Errorf(vrfcommon.ErrGenericFormat, vrfcommon.ErrRestartCLNode, err)
	}
	return nil
}

func SetupVRFV2PlusWrapperEnvironment(
	ctx context.Context,
	l zerolog.Logger,
	env *test_env.CLClusterTestEnv,
	chainID int64,
	vrfv2PlusTestConfig types.VRFv2PlusTestConfig,
	linkToken contracts.LinkToken,
	mockNativeLINKFeed contracts.MockETHLINKFeed,
	coordinator contracts.VRFCoordinatorV2_5,
	keyHash [32]byte,
	wrapperConsumerContractsAmount int,
) (*VRFV2PlusWrapperContracts, *big.Int, error) {
	// external EOA has to create a subscription for the wrapper first
	wrapperSubId, err := CreateSubAndFindSubID(ctx, env, chainID, coordinator)
	if err != nil {
		return nil, nil, err
	}

	vrfv2PlusConfig := vrfv2PlusTestConfig.GetVRFv2PlusConfig().General

	sethClient, err := env.GetSethClient(chainID)
	if err != nil {
		return nil, nil, err
	}

	wrapperContracts, err := DeployVRFV2PlusDirectFundingContracts(
		sethClient,
		linkToken.Address(),
		mockNativeLINKFeed.Address(),
		coordinator,
		wrapperConsumerContractsAmount,
		wrapperSubId,
	)
	if err != nil {
		return nil, nil, fmt.Errorf(vrfcommon.ErrGenericFormat, vrfcommon.ErrWaitTXsComplete, err)
	}

	// once the wrapper is deployed, wrapper address will become consumer of external EOA subscription
	err = coordinator.AddConsumer(wrapperSubId, wrapperContracts.VRFV2PlusWrapper.Address())
	if err != nil {
		return nil, nil, err
	}

	err = wrapperContracts.VRFV2PlusWrapper.SetConfig(
		*vrfv2PlusConfig.WrapperGasOverhead,
		*vrfv2PlusConfig.CoordinatorGasOverheadNative,
		*vrfv2PlusConfig.CoordinatorGasOverheadLink,
		*vrfv2PlusConfig.CoordinatorGasOverheadPerWord,
		*vrfv2PlusConfig.CoordinatorNativePremiumPercentage,
		*vrfv2PlusConfig.CoordinatorLinkPremiumPercentage,
		keyHash,
		*vrfv2PlusConfig.WrapperMaxNumberOfWords,
		*vrfv2PlusConfig.StalenessSeconds,
		decimal.RequireFromString(*vrfv2PlusConfig.FallbackWeiPerUnitLink).BigInt(),
		*vrfv2PlusConfig.FulfillmentFlatFeeNativePPM,
		*vrfv2PlusConfig.FulfillmentFlatFeeLinkDiscountPPM,
	)
	if err != nil {
		return nil, nil, err
	}

	//fund sub
	wrapperSubID, err := wrapperContracts.VRFV2PlusWrapper.GetSubID(ctx)
	if err != nil {
		return nil, nil, err
	}

	err = FundSubscriptions(
		big.NewFloat(*vrfv2PlusTestConfig.GetVRFv2PlusConfig().General.SubscriptionFundingAmountNative),
		big.NewFloat(*vrfv2PlusTestConfig.GetVRFv2PlusConfig().General.SubscriptionFundingAmountLink),
		linkToken,
		coordinator,
		[]*big.Int{wrapperSubID},
	)
	if err != nil {
		return nil, nil, err
	}

	//fund consumer with Link
	err = linkToken.Transfer(
		wrapperContracts.LoadTestConsumers[0].Address(),
		big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(*vrfv2PlusConfig.WrapperConsumerFundingAmountLink)),
	)
	if err != nil {
		return nil, nil, err
	}

	//fund consumer with Eth (native token)
	_, err = actions_seth.SendFunds(l, sethClient, actions_seth.FundsToSendPayload{
		ToAddress:  common.HexToAddress(wrapperContracts.LoadTestConsumers[0].Address()),
		Amount:     conversions.EtherToWei(big.NewFloat(*vrfv2PlusConfig.WrapperConsumerFundingAmountNativeToken)),
		PrivateKey: sethClient.PrivateKeys[0],
	})
	if err != nil {
		return nil, nil, err
	}

	wrapperConsumerBalanceBeforeRequestWei, err := sethClient.Client.BalanceAt(ctx, common.HexToAddress(wrapperContracts.LoadTestConsumers[0].Address()), nil)
	if err != nil {
		return nil, nil, err
	}
	l.Info().
		Str("WrapperConsumerBalanceBeforeRequestWei", wrapperConsumerBalanceBeforeRequestWei.String()).
		Str("WrapperConsumerAddress", wrapperContracts.LoadTestConsumers[0].Address()).
		Msg("WrapperConsumerBalanceBeforeRequestWei")

	return wrapperContracts, wrapperSubID, nil
}

func SetupVRFV2PlusUniverse(
	ctx context.Context,
	t *testing.T,
	envConfig vrfcommon.VRFEnvConfig,
	newEnvConfig vrfcommon.NewEnvConfig,
	l zerolog.Logger,
) (*test_env.CLClusterTestEnv, *vrfcommon.VRFContracts, *vrfcommon.VRFKeyData, map[vrfcommon.VRFNodeType]*vrfcommon.VRFNode, error) {
	var (
		env            *test_env.CLClusterTestEnv
		vrfContracts   *vrfcommon.VRFContracts
		vrfKey         *vrfcommon.VRFKeyData
		nodeTypeToNode map[vrfcommon.VRFNodeType]*vrfcommon.VRFNode
		err            error
	)
	if *envConfig.TestConfig.VRFv2Plus.General.UseExistingEnv {
		vrfContracts, vrfKey, env, err = SetupVRFV2PlusForExistingEnv(t, envConfig, l)
		if err != nil {
			return nil, nil, nil, nil, fmt.Errorf("%s, err: %w", "Error setting up VRF V2 Plus for Existing env", err)
		}
	} else {
		vrfContracts, vrfKey, env, nodeTypeToNode, err = SetupVRFV2PlusForNewEnv(ctx, t, envConfig, newEnvConfig, l)
		if err != nil {
			return nil, nil, nil, nil, fmt.Errorf("%s, err: %w", "Error setting up VRF V2 Plus for New env", err)
		}
	}
	return env, vrfContracts, vrfKey, nodeTypeToNode, nil
}

func SetupVRFV2PlusForNewEnv(
	ctx context.Context,
	t *testing.T,
	envConfig vrfcommon.VRFEnvConfig,
	newEnvConfig vrfcommon.NewEnvConfig,
	l zerolog.Logger,
) (*vrfcommon.VRFContracts, *vrfcommon.VRFKeyData, *test_env.CLClusterTestEnv, map[vrfcommon.VRFNodeType]*vrfcommon.VRFNode, error) {
	network, err := actions.EthereumNetworkConfigFromConfig(l, &envConfig.TestConfig)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("%s, err: %w", "Error building ethereum network config for V2Plus", err)
	}
	env, err := vrfcommon.BuildNewCLEnvForVRF(t, envConfig, newEnvConfig, network)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	sethClient, err := env.GetSethClient(envConfig.ChainID)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	mockETHLinkFeed, err := contracts.DeployVRFMockETHLINKFeed(sethClient, big.NewInt(*envConfig.TestConfig.VRFv2Plus.General.LinkNativeFeedResponse))
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("%s, err: %w", "error deploying mock ETH/LINK feed", err)
	}
	linkToken, err := contracts.DeployLinkTokenContract(l, sethClient)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("%s, err: %w", "error deploying LINK contract", err)
	}
	vrfContracts, vrfKey, nodeTypeToNode, err := SetupVRFV2_5Environment(
		ctx,
		env,
		envConfig.ChainID,
		newEnvConfig.NodesToCreate,
		&envConfig.TestConfig,
		linkToken,
		mockETHLinkFeed,
		newEnvConfig.NumberOfTxKeysToCreate,
		l,
	)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("%s, err: %w", "error setting up VRF v2_5 env", err)
	}
	return vrfContracts, vrfKey, env, nodeTypeToNode, nil
}

func SetupVRFV2PlusForExistingEnv(t *testing.T, envConfig vrfcommon.VRFEnvConfig, l zerolog.Logger) (*vrfcommon.VRFContracts, *vrfcommon.VRFKeyData, *test_env.CLClusterTestEnv, error) {
	commonExistingEnvConfig := envConfig.TestConfig.VRFv2Plus.ExistingEnvConfig.ExistingEnvConfig
	env, err := test_env.NewCLTestEnvBuilder().
		WithTestInstance(t).
		WithTestConfig(&envConfig.TestConfig).
		WithCustomCleanup(envConfig.CleanupFn).
		WithSeth().
		Build()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%s, err: %w", "error creating test env", err)
	}
	sethClient, err := env.GetSethClient(envConfig.ChainID)
	if err != nil {
		return nil, nil, nil, err
	}
	coordinator, err := contracts.LoadVRFCoordinatorV2_5(sethClient, *commonExistingEnvConfig.CoordinatorAddress)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%s, err: %w", "error loading VRFCoordinator2_5", err)
	}
	linkToken, err := contracts.LoadLinkTokenContract(l, sethClient, common.HexToAddress(*commonExistingEnvConfig.LinkAddress))
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%s, err: %w", "error loading LinkToken", err)
	}
	err = vrfcommon.FundNodesIfNeeded(testcontext.Get(t), commonExistingEnvConfig, sethClient, l)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("err: %w", err)
	}
	blockHashStoreAddress, err := coordinator.GetBlockHashStoreAddress(testcontext.Get(t))
	if err != nil {
		return nil, nil, nil, fmt.Errorf("err: %w", err)
	}
	blockHashStore, err := contracts.LoadBlockHashStore(sethClient, blockHashStoreAddress.String())
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%s, err: %w", "error loading BlockHashStore", err)
	}
	vrfContracts := &vrfcommon.VRFContracts{
		CoordinatorV2Plus: coordinator,
		VRFV2PlusConsumer: nil,
		LinkToken:         linkToken,
		BHS:               blockHashStore,
	}
	vrfKey := &vrfcommon.VRFKeyData{
		VRFKey:            nil,
		EncodedProvingKey: [2]*big.Int{},
		KeyHash:           common.HexToHash(*commonExistingEnvConfig.KeyHash),
	}
	return vrfContracts, vrfKey, env, nil
}

func SetupSubsAndConsumersForExistingEnv(
	ctx context.Context,
	env *test_env.CLClusterTestEnv,
	chainID int64,
	coordinator contracts.VRFCoordinatorV2_5,
	linkToken contracts.LinkToken,
	numberOfConsumerContractsToDeployAndAddToSub int,
	numberOfSubToCreate int,
	testConfig tc.TestConfig,
	l zerolog.Logger,
) ([]*big.Int, []contracts.VRFv2PlusLoadTestConsumer, error) {
	var (
		subIDs    []*big.Int
		consumers []contracts.VRFv2PlusLoadTestConsumer
		err       error
	)
	sethClient, err := env.GetSethClient(chainID)
	if err != nil {
		return nil, nil, err
	}
	if *testConfig.VRFv2Plus.General.UseExistingEnv {
		commonExistingEnvConfig := testConfig.VRFv2Plus.ExistingEnvConfig.ExistingEnvConfig
		if *commonExistingEnvConfig.CreateFundSubsAndAddConsumers {
			consumers, subIDs, err = SetupNewConsumersAndSubs(
				ctx,
				env,
				chainID,
				coordinator,
				testConfig,
				linkToken,
				numberOfConsumerContractsToDeployAndAddToSub,
				numberOfSubToCreate,
				l,
			)
			if err != nil {
				return nil, nil, fmt.Errorf("err: %w", err)
			}
		} else {
			consumer, err := contracts.LoadVRFv2PlusLoadTestConsumer(sethClient, *commonExistingEnvConfig.ConsumerAddress)
			if err != nil {
				return nil, nil, fmt.Errorf("err: %w", err)
			}
			consumers = append(consumers, consumer)
			var ok bool
			subID, ok := new(big.Int).SetString(*testConfig.VRFv2Plus.ExistingEnvConfig.SubID, 10)
			if !ok {
				return nil, nil, fmt.Errorf("unable to parse subID: %s %w", *testConfig.VRFv2Plus.ExistingEnvConfig.SubID, err)
			}
			subIDs = append(subIDs, subID)
		}
	} else {
		consumers, subIDs, err = SetupNewConsumersAndSubs(
			ctx,
			env,
			chainID,
			coordinator,
			testConfig,
			linkToken,
			numberOfConsumerContractsToDeployAndAddToSub,
			numberOfSubToCreate,
			l,
		)
		if err != nil {
			return nil, nil, fmt.Errorf("err: %w", err)
		}
	}
	return subIDs, consumers, nil
}
