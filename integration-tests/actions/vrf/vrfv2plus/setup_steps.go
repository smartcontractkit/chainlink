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

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	"github.com/smartcontractkit/chainlink/deployment/environment/nodeclient"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	vrfcommon "github.com/smartcontractkit/chainlink/integration-tests/actions/vrf/common"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
	vrfv2plusconfig "github.com/smartcontractkit/chainlink/integration-tests/testconfig/vrfv2plus"
	"github.com/smartcontractkit/chainlink/integration-tests/types"
	"github.com/smartcontractkit/chainlink/integration-tests/types/config/node"

	"github.com/smartcontractkit/chainlink-evm/pkg/assets"
)

func CreateVRFV2PlusJob(
	chainlinkNode *nodeclient.ChainlinkClient,
	vrfJobSpecConfig vrfcommon.VRFJobSpecConfig,
) (*nodeclient.Job, error) {
	jobUUID := uuid.New()
	os := &nodeclient.VRFV2PlusTxPipelineSpec{
		Address:               vrfJobSpecConfig.CoordinatorAddress,
		EstimateGasMultiplier: vrfJobSpecConfig.EstimateGasMultiplier,
		FromAddress:           vrfJobSpecConfig.FromAddresses[0],
		SimulationBlock:       vrfJobSpecConfig.SimulationBlock,
	}
	ost, err := os.String()
	if err != nil {
		return nil, fmt.Errorf(vrfcommon.ErrGenericFormat, vrfcommon.ErrParseJob, err)
	}

	jobSpec := nodeclient.VRFV2PlusJobSpec{
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
	sethClient *seth.Client,
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
		sethClient,
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
	provingKey, err := VRFV2_5RegisterProvingKey(vrfKey, vrfContracts.CoordinatorV2Plus, assets.GWei(*configGeneral.CLNodeMaxGasPriceGWei).ToInt().Uint64())
	if err != nil {
		return nil, nil, nil, fmt.Errorf(vrfcommon.ErrGenericFormat, vrfcommon.ErrRegisteringProvingKey, err)
	}
	keyHash, err := vrfContracts.CoordinatorV2Plus.HashOfKey(ctx, provingKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf(vrfcommon.ErrGenericFormat, vrfcommon.ErrCreatingProvingKeyHash, err)
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
				sethClient,
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
				sethClient,
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

func setupVRFNode(contracts *vrfcommon.VRFContracts, chainID *big.Int, config *vrfv2plusconfig.General, pubKeyCompressed string, l zerolog.Logger, vrfNode *vrfcommon.VRFNode) error {
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
	l.Info().
		Strs("Sending Keys", vrfNode.TXKeyAddressStrings).
		Int64("Price Max Setting", *config.CLNodeMaxGasPriceGWei).
		Msg("Restarting Node with new sending key PriceMax configuration")
	err = vrfNode.CLNode.Restart(nodeConfig)
	if err != nil {
		return fmt.Errorf(vrfcommon.ErrGenericFormat, vrfcommon.ErrRestartCLNode, err)
	}
	return nil
}

func SetupVRFV2PlusWrapperForExistingEnv(
	ctx context.Context,
	sethClient *seth.Client,
	vrfContracts *vrfcommon.VRFContracts,
	keyHash [32]byte,
	vrfv2PlusTestConfig types.VRFv2PlusTestConfig,
	numberOfConsumerContracts int,
	l zerolog.Logger,
) (*VRFV2PlusWrapperContracts, *big.Int, error) {
	config := *vrfv2PlusTestConfig.GetVRFv2PlusConfig()
	var wrapper contracts.VRFV2PlusWrapper
	var err error
	if *config.ExistingEnvConfig.UseExistingWrapper {
		wrapper, err = contracts.LoadVRFV2PlusWrapper(sethClient, *config.ExistingEnvConfig.WrapperAddress)
		if err != nil {
			return nil, nil, fmt.Errorf(vrfcommon.ErrGenericFormat, "error loading VRFV2PlusWrapper", err)
		}
	} else {
		wrapperSubId, err := CreateSubAndFindSubID(ctx, sethClient, vrfContracts.CoordinatorV2Plus)
		if err != nil {
			return nil, nil, err
		}
		wrapper, err = contracts.DeployVRFV2PlusWrapper(sethClient, vrfContracts.LinkToken.Address(), vrfContracts.LinkNativeFeedAddress, vrfContracts.CoordinatorV2Plus.Address(), wrapperSubId)
		if err != nil {
			return nil, nil, fmt.Errorf(vrfcommon.ErrGenericFormat, ErrDeployWrapper, err)
		}
		err = FundSubscriptions(
			big.NewFloat(*config.General.SubscriptionFundingAmountNative),
			big.NewFloat(*config.General.SubscriptionFundingAmountLink),
			vrfContracts.LinkToken,
			vrfContracts.CoordinatorV2Plus,
			[]*big.Int{wrapperSubId},
			*config.General.SubscriptionBillingType,
		)
		if err != nil {
			return nil, nil, err
		}
		err = vrfContracts.CoordinatorV2Plus.AddConsumer(wrapperSubId, wrapper.Address())
		if err != nil {
			return nil, nil, err
		}
		err = wrapper.SetConfig(
			*config.General.WrapperGasOverhead,
			*config.General.CoordinatorGasOverheadNative,
			*config.General.CoordinatorGasOverheadLink,
			*config.General.CoordinatorGasOverheadPerWord,
			*config.General.CoordinatorNativePremiumPercentage,
			*config.General.CoordinatorLinkPremiumPercentage,
			keyHash,
			*config.General.WrapperMaxNumberOfWords,
			*config.General.StalenessSeconds,
			decimal.RequireFromString(*config.General.FallbackWeiPerUnitLink).BigInt(),
			*config.General.FulfillmentFlatFeeNativePPM,
			*config.General.FulfillmentFlatFeeLinkDiscountPPM,
		)
		if err != nil {
			return nil, nil, err
		}
	}
	wrapperSubID, err := wrapper.GetSubID(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf(vrfcommon.ErrGenericFormat, "error getting subID", err)
	}
	var wrapperConsumers []contracts.VRFv2PlusWrapperLoadTestConsumer
	if *config.ExistingEnvConfig.CreateFundAddWrapperConsumers {
		wrapperConsumers, err = DeployVRFV2PlusWrapperConsumers(sethClient, wrapper, numberOfConsumerContracts)
		if err != nil {
			return nil, nil, err
		}
	} else {
		wrapperConsumer, err := contracts.LoadVRFV2WrapperLoadTestConsumer(sethClient, *config.ExistingEnvConfig.WrapperConsumerAddress)
		if err != nil {
			return nil, nil, fmt.Errorf(vrfcommon.ErrGenericFormat, "error loading VRFV2WrapperLoadTestConsumer", err)
		}
		wrapperConsumers = append(wrapperConsumers, wrapperConsumer)
	}
	wrapperContracts := &VRFV2PlusWrapperContracts{wrapper, wrapperConsumers}
	for _, consumer := range wrapperConsumers {
		err = FundWrapperConsumer(
			sethClient,
			*config.General.SubscriptionBillingType,
			vrfContracts.LinkToken,
			consumer,
			config.General,
			l,
		)
		if err != nil {
			return nil, nil, err
		}
	}
	return wrapperContracts, wrapperSubID, nil
}

func SetupVRFV2PlusWrapperForNewEnv(
	ctx context.Context,
	sethClient *seth.Client,
	vrfv2PlusTestConfig types.VRFv2PlusTestConfig,
	vrfContracts *vrfcommon.VRFContracts,
	keyHash [32]byte,
	wrapperConsumerContractsAmount int,
	l zerolog.Logger,
) (*VRFV2PlusWrapperContracts, *big.Int, error) {
	// external EOA has to create a subscription for the wrapper first
	wrapperSubId, err := CreateSubAndFindSubID(ctx, sethClient, vrfContracts.CoordinatorV2Plus)
	if err != nil {
		return nil, nil, err
	}
	vrfv2PlusConfig := vrfv2PlusTestConfig.GetVRFv2PlusConfig().General
	wrapperContracts, err := DeployVRFV2PlusDirectFundingContracts(
		sethClient,
		vrfContracts.LinkToken.Address(),
		vrfContracts.MockETHLINKFeed.Address(),
		vrfContracts.CoordinatorV2Plus,
		wrapperConsumerContractsAmount,
		wrapperSubId,
		vrfv2PlusConfig,
	)
	if err != nil {
		return nil, nil, err
	}
	// once the wrapper is deployed, wrapper address will become consumer of external EOA subscription
	err = vrfContracts.CoordinatorV2Plus.AddConsumer(wrapperSubId, wrapperContracts.VRFV2PlusWrapper.Address())
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
		vrfContracts.LinkToken,
		vrfContracts.CoordinatorV2Plus,
		[]*big.Int{wrapperSubID},
		*vrfv2PlusConfig.SubscriptionBillingType,
	)
	if err != nil {
		return nil, nil, err
	}
	for _, consumer := range wrapperContracts.WrapperConsumers {
		err = FundWrapperConsumer(
			sethClient,
			*vrfv2PlusConfig.SubscriptionBillingType,
			vrfContracts.LinkToken,
			consumer,
			vrfv2PlusConfig,
			l,
		)
		if err != nil {
			return nil, nil, err
		}
	}
	return wrapperContracts, wrapperSubID, nil
}

func SetupVRFV2PlusUniverse(
	ctx context.Context,
	t *testing.T,
	envConfig vrfcommon.VRFEnvConfig,
	newEnvConfig vrfcommon.NewEnvConfig,
	l zerolog.Logger,
) (*test_env.CLClusterTestEnv, *vrfcommon.VRFContracts, *vrfcommon.VRFKeyData, map[vrfcommon.VRFNodeType]*vrfcommon.VRFNode, *seth.Client, error) {
	var (
		env            *test_env.CLClusterTestEnv
		vrfContracts   *vrfcommon.VRFContracts
		vrfKey         *vrfcommon.VRFKeyData
		nodeTypeToNode map[vrfcommon.VRFNodeType]*vrfcommon.VRFNode
		sethClient     *seth.Client
		err            error
	)
	if *envConfig.TestConfig.VRFv2Plus.General.UseExistingEnv {
		vrfContracts, vrfKey, env, sethClient, err = SetupVRFV2PlusForExistingEnv(t, envConfig, l)
		if err != nil {
			return nil, nil, nil, nil, nil, fmt.Errorf("%s, err: %w", "Error setting up VRF V2 Plus for Existing env", err)
		}
	} else {
		vrfContracts, vrfKey, env, nodeTypeToNode, sethClient, err = SetupVRFV2PlusForNewEnv(ctx, t, envConfig, newEnvConfig, l)
		if err != nil {
			return nil, nil, nil, nil, nil, fmt.Errorf("%s, err: %w", "Error setting up VRF V2 Plus for New env", err)
		}
	}

	return env, vrfContracts, vrfKey, nodeTypeToNode, sethClient, nil
}

func SetupVRFV2PlusForNewEnv(
	ctx context.Context,
	t *testing.T,
	envConfig vrfcommon.VRFEnvConfig,
	newEnvConfig vrfcommon.NewEnvConfig,
	l zerolog.Logger,
) (*vrfcommon.VRFContracts, *vrfcommon.VRFKeyData, *test_env.CLClusterTestEnv, map[vrfcommon.VRFNodeType]*vrfcommon.VRFNode, *seth.Client, error) {
	network, err := actions.EthereumNetworkConfigFromConfig(l, &envConfig.TestConfig)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("%s, err: %w", "Error building ethereum network config for V2Plus", err)
	}
	env, sethClient, err := vrfcommon.BuildNewCLEnvForVRF(l, t, envConfig, newEnvConfig, network)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	mockETHLinkFeed, err := contracts.DeployVRFMockETHLINKFeed(sethClient, big.NewInt(*envConfig.TestConfig.VRFv2Plus.General.LinkNativeFeedResponse))
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("%s, err: %w", "error deploying mock ETH/LINK feed", err)
	}
	linkToken, err := contracts.DeployLinkTokenContract(l, sethClient)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("%s, err: %w", "error deploying LINK contract", err)
	}
	vrfContracts, vrfKey, nodeTypeToNode, err := SetupVRFV2_5Environment(
		ctx,
		sethClient,
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
		return nil, nil, nil, nil, nil, fmt.Errorf("%s, err: %w", "error setting up VRF v2_5 env", err)
	}

	return vrfContracts, vrfKey, env, nodeTypeToNode, sethClient, nil
}

func SetupVRFV2PlusForExistingEnv(t *testing.T, envConfig vrfcommon.VRFEnvConfig, l zerolog.Logger) (*vrfcommon.VRFContracts, *vrfcommon.VRFKeyData, *test_env.CLClusterTestEnv, *seth.Client, error) {
	commonExistingEnvConfig := envConfig.TestConfig.VRFv2Plus.ExistingEnvConfig.ExistingEnvConfig
	env, sethClient, err := vrfcommon.LoadExistingCLEnvForVRF(
		t,
		envConfig,
		commonExistingEnvConfig,
		l,
	)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("%s, err: %w", "error loading existing CL env", err)
	}
	coordinator, err := contracts.LoadVRFCoordinatorV2_5(sethClient, *commonExistingEnvConfig.CoordinatorAddress)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("%s, err: %w", "error loading VRFCoordinator2_5", err)
	}
	linkAddress, err := coordinator.GetLinkAddress(testcontext.Get(t))
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("%s, err: %w", "error getting Link address from Coordinator", err)
	}
	linkToken, err := contracts.LoadLinkTokenContract(l, sethClient, common.HexToAddress(linkAddress.String()))
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("%s, err: %w", "error loading LinkToken", err)
	}
	linkNativeFeedAddress, err := coordinator.GetLinkNativeFeed(testcontext.Get(t))
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("%s, err: %w", "error getting Link address from Coordinator", err)
	}
	blockHashStoreAddress, err := coordinator.GetBlockHashStoreAddress(testcontext.Get(t))
	if err != nil {
		return nil, nil, nil, nil, err
	}
	blockHashStore, err := contracts.LoadBlockHashStore(sethClient, blockHashStoreAddress.String())
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("%s, err: %w", "error loading BlockHashStore", err)
	}
	vrfContracts := &vrfcommon.VRFContracts{
		CoordinatorV2Plus:     coordinator,
		VRFV2PlusConsumer:     nil,
		LinkToken:             linkToken,
		BHS:                   blockHashStore,
		LinkNativeFeedAddress: linkNativeFeedAddress.String(),
	}
	vrfKey := &vrfcommon.VRFKeyData{
		VRFKey:            nil,
		EncodedProvingKey: [2]*big.Int{},
		KeyHash:           common.HexToHash(*commonExistingEnvConfig.KeyHash),
	}
	return vrfContracts, vrfKey, env, sethClient, nil
}

func SetupSubsAndConsumersForExistingEnv(
	ctx context.Context,
	sethClient *seth.Client,
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
	if *testConfig.VRFv2Plus.General.UseExistingEnv {
		commonExistingEnvConfig := testConfig.VRFv2Plus.ExistingEnvConfig.ExistingEnvConfig
		if *commonExistingEnvConfig.CreateFundSubsAndAddConsumers {
			consumers, subIDs, err = SetupNewConsumersAndSubs(
				ctx,
				sethClient,
				coordinator,
				testConfig,
				linkToken,
				numberOfConsumerContractsToDeployAndAddToSub,
				numberOfSubToCreate,
				l,
			)
			if err != nil {
				return nil, nil, err
			}
		} else {
			consumer, err := contracts.LoadVRFv2PlusLoadTestConsumer(sethClient, *commonExistingEnvConfig.ConsumerAddress)
			if err != nil {
				return nil, nil, err
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
			sethClient,
			coordinator,
			testConfig,
			linkToken,
			numberOfConsumerContractsToDeployAndAddToSub,
			numberOfSubToCreate,
			l,
		)
		if err != nil {
			return nil, nil, err
		}
	}
	return subIDs, consumers, nil
}

func SelectBillingTypeWithDistribution(billingType string, distributionFn func() bool) (bool, error) {
	switch vrfv2plusconfig.BillingType(billingType) {
	case vrfv2plusconfig.BillingType_Link:
		return false, nil
	case vrfv2plusconfig.BillingType_Native:
		return true, nil
	case vrfv2plusconfig.BillingType_Link_and_Native:
		return distributionFn(), nil
	default:
		return false, fmt.Errorf("invalid billing type: %s", billingType)
	}
}

func SetupVRFV2PlusWrapperUniverse(
	ctx context.Context,
	sethClient *seth.Client,
	vrfContracts *vrfcommon.VRFContracts,
	config *tc.TestConfig,
	keyHash [32]byte,
	numberOfConsumerContracts int,
	l zerolog.Logger,
) (*VRFV2PlusWrapperContracts, *big.Int, error) {
	var (
		wrapperContracts *VRFV2PlusWrapperContracts
		wrapperSubID     *big.Int
		err              error
	)
	if *config.VRFv2Plus.General.UseExistingEnv {
		wrapperContracts, wrapperSubID, err = SetupVRFV2PlusWrapperForExistingEnv(
			ctx,
			sethClient,
			vrfContracts,
			keyHash,
			config,
			numberOfConsumerContracts,
			l,
		)
		if err != nil {
			return nil, nil, err
		}
	} else {
		wrapperContracts, wrapperSubID, err = SetupVRFV2PlusWrapperForNewEnv(
			ctx,
			sethClient,
			config,
			vrfContracts,
			keyHash,
			numberOfConsumerContracts,
			l,
		)
		if err != nil {
			return nil, nil, err
		}
	}
	return wrapperContracts, wrapperSubID, nil
}
