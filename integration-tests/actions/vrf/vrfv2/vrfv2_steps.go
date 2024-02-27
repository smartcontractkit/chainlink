package vrfv2

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"

	commonassets "github.com/smartcontractkit/chainlink-common/pkg/assets"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/conversions"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	testconfig "github.com/smartcontractkit/chainlink/integration-tests/testconfig/vrfv2"
	"github.com/smartcontractkit/chainlink/integration-tests/types/config/node"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_owner"

	"github.com/google/uuid"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	chainlinkutils "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	vrfcommon "github.com/smartcontractkit/chainlink/integration-tests/actions/vrf/common"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"

	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
	"github.com/smartcontractkit/chainlink/integration-tests/types"
)

func DeployVRFV2Contracts(
	env *test_env.CLClusterTestEnv,
	linkTokenContract contracts.LinkToken,
	linkEthFeedContract contracts.MockETHLINKFeed,
	consumerContractsAmount int,
	useVRFOwner bool,
	useTestCoordinator bool,
) (*vrfcommon.VRFContracts, error) {
	bhs, err := env.ContractDeployer.DeployBlockhashStore()
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrDeployBlockHashStore, err)
	}
	err = env.EVMClient.WaitForEvents()
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrWaitTXsComplete, err)
	}

	var coordinatorAddress string
	if useTestCoordinator {
		testCoordinator, err := env.ContractDeployer.DeployVRFCoordinatorTestV2(linkTokenContract.Address(), bhs.Address(), linkEthFeedContract.Address())
		if err != nil {
			return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrDeployCoordinator, err)
		}
		err = env.EVMClient.WaitForEvents()
		if err != nil {
			return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrWaitTXsComplete, err)
		}
		coordinatorAddress = testCoordinator.Address()
	} else {
		coordinator, err := env.ContractDeployer.DeployVRFCoordinatorV2(linkTokenContract.Address(), bhs.Address(), linkEthFeedContract.Address())
		if err != nil {
			return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrDeployCoordinator, err)
		}
		err = env.EVMClient.WaitForEvents()
		if err != nil {
			return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrWaitTXsComplete, err)
		}
		coordinatorAddress = coordinator.Address()
	}
	consumers, err := DeployVRFV2Consumers(env.ContractDeployer, coordinatorAddress, consumerContractsAmount)
	if err != nil {
		return nil, err
	}
	err = env.EVMClient.WaitForEvents()
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrWaitTXsComplete, err)
	}

	coordinator, err := env.ContractLoader.LoadVRFCoordinatorV2(coordinatorAddress)
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrLoadingCoordinator, err)
	}
	if useVRFOwner {
		vrfOwner, err := env.ContractDeployer.DeployVRFOwner(coordinatorAddress)
		if err != nil {
			return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrDeployCoordinator, err)
		}
		err = env.EVMClient.WaitForEvents()
		if err != nil {
			return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrWaitTXsComplete, err)
		}
		return &vrfcommon.VRFContracts{
			CoordinatorV2: coordinator,
			VRFOwner:      vrfOwner,
			BHS:           bhs,
			VRFV2Consumer: consumers,
		}, nil
	}
	return &vrfcommon.VRFContracts{
		CoordinatorV2: coordinator,
		VRFOwner:      nil,
		BHS:           bhs,
		VRFV2Consumer: consumers,
	}, nil
}

func DeployVRFV2Consumers(contractDeployer contracts.ContractDeployer, coordinatorAddress string, consumerContractsAmount int) ([]contracts.VRFv2LoadTestConsumer, error) {
	var consumers []contracts.VRFv2LoadTestConsumer
	for i := 1; i <= consumerContractsAmount; i++ {
		loadTestConsumer, err := contractDeployer.DeployVRFv2LoadTestConsumer(coordinatorAddress)
		if err != nil {
			return nil, fmt.Errorf("%s, err %w", ErrAdvancedConsumer, err)
		}
		consumers = append(consumers, loadTestConsumer)
	}
	return consumers, nil
}

func DeployVRFV2WrapperConsumers(contractDeployer contracts.ContractDeployer, linkTokenAddress string, vrfV2Wrapper contracts.VRFV2Wrapper, consumerContractsAmount int) ([]contracts.VRFv2WrapperLoadTestConsumer, error) {
	var consumers []contracts.VRFv2WrapperLoadTestConsumer
	for i := 1; i <= consumerContractsAmount; i++ {
		loadTestConsumer, err := contractDeployer.DeployVRFV2WrapperLoadTestConsumer(linkTokenAddress, vrfV2Wrapper.Address())
		if err != nil {
			return nil, fmt.Errorf("%s, err %w", ErrAdvancedConsumer, err)
		}
		consumers = append(consumers, loadTestConsumer)
	}
	return consumers, nil
}

func DeployVRFV2DirectFundingContracts(
	contractDeployer contracts.ContractDeployer,
	chainClient blockchain.EVMClient,
	linkTokenAddress string,
	linkEthFeedAddress string,
	coordinator contracts.VRFCoordinatorV2,
	consumerContractsAmount int,
) (*VRFV2WrapperContracts, error) {
	vrfv2Wrapper, err := contractDeployer.DeployVRFV2Wrapper(linkTokenAddress, linkEthFeedAddress, coordinator.Address())
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", ErrDeployVRFV2Wrapper, err)
	}
	err = chainClient.WaitForEvents()
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrWaitTXsComplete, err)
	}

	consumers, err := DeployVRFV2WrapperConsumers(contractDeployer, linkTokenAddress, vrfv2Wrapper, consumerContractsAmount)
	if err != nil {
		return nil, err
	}
	err = chainClient.WaitForEvents()
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrWaitTXsComplete, err)
	}
	return &VRFV2WrapperContracts{vrfv2Wrapper, consumers}, nil
}

func CreateVRFV2Job(
	chainlinkNode *client.ChainlinkClient,
	vrfJobSpecConfig vrfcommon.VRFJobSpecConfig,
) (*client.Job, error) {
	jobUUID := uuid.New()
	os := &client.VRFV2TxPipelineSpec{
		Address:               vrfJobSpecConfig.CoordinatorAddress,
		EstimateGasMultiplier: vrfJobSpecConfig.EstimateGasMultiplier,
		FromAddress:           vrfJobSpecConfig.FromAddresses[0],
		SimulationBlock:       vrfJobSpecConfig.SimulationBlock,
	}
	ost, err := os.String()
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrParseJob, err)
	}

	spec := &client.VRFV2JobSpec{
		Name:                          fmt.Sprintf("vrf-v2-%s", jobUUID),
		ForwardingAllowed:             vrfJobSpecConfig.ForwardingAllowed,
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
	if vrfJobSpecConfig.VRFOwnerConfig.UseVRFOwner {
		spec.VRFOwner = vrfJobSpecConfig.VRFOwnerConfig.OwnerAddress
		spec.UseVRFOwner = true
	}

	if err != nil {
		return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrParseJob, err)

	}
	job, err := chainlinkNode.MustCreateJob(spec)
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", ErrCreatingVRFv2Job, err)
	}
	return job, nil
}

func VRFV2RegisterProvingKey(
	vrfKey *client.VRFKey,
	oracleAddress string,
	coordinator contracts.VRFCoordinatorV2,
) (vrfcommon.VRFEncodedProvingKey, error) {
	provingKey, err := actions.EncodeOnChainVRFProvingKey(*vrfKey)
	if err != nil {
		return vrfcommon.VRFEncodedProvingKey{}, fmt.Errorf("%s, err %w", vrfcommon.ErrEncodingProvingKey, err)
	}
	err = coordinator.RegisterProvingKey(
		oracleAddress,
		provingKey,
	)
	if err != nil {
		return vrfcommon.VRFEncodedProvingKey{}, fmt.Errorf("%s, err %w", vrfcommon.ErrRegisterProvingKey, err)
	}
	return provingKey, nil
}

func FundVRFCoordinatorV2Subscription(
	linkToken contracts.LinkToken,
	coordinator contracts.VRFCoordinatorV2,
	chainClient blockchain.EVMClient,
	subscriptionID uint64,
	linkFundingAmountJuels *big.Int,
) error {
	encodedSubId, err := chainlinkutils.ABIEncode(`[{"type":"uint64"}]`, subscriptionID)
	if err != nil {
		return fmt.Errorf("%s, err %w", vrfcommon.ErrABIEncodingFunding, err)
	}
	_, err = linkToken.TransferAndCall(coordinator.Address(), linkFundingAmountJuels, encodedSubId)
	if err != nil {
		return fmt.Errorf("%s, err %w", vrfcommon.ErrSendingLinkToken, err)
	}
	return chainClient.WaitForEvents()
}

// SetupVRFV2Environment will create specified number of subscriptions and add the same conumer/s to each of them
func SetupVRFV2Environment(
	env *test_env.CLClusterTestEnv,
	nodesToCreate []vrfcommon.VRFNodeType,
	vrfv2TestConfig types.VRFv2TestConfig,
	useVRFOwner bool,
	useTestCoordinator bool,
	linkToken contracts.LinkToken,
	mockNativeLINKFeed contracts.MockETHLINKFeed,
	registerProvingKeyAgainstAddress string,
	numberOfTxKeysToCreate int,
	numberOfConsumers int,
	numberOfSubToCreate int,
	l zerolog.Logger,
) (*vrfcommon.VRFContracts, []uint64, *vrfcommon.VRFKeyData, map[vrfcommon.VRFNodeType]*vrfcommon.VRFNode, error) {
	l.Info().Msg("Starting VRFV2 environment setup")
	configGeneral := vrfv2TestConfig.GetVRFv2Config().General
	vrfContracts, subIDs, err := SetupVRFV2Contracts(
		env,
		linkToken,
		mockNativeLINKFeed,
		numberOfConsumers,
		useVRFOwner,
		useTestCoordinator,
		configGeneral,
		numberOfSubToCreate,
		l,
	)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	nodeTypeToNodeMap := vrfcommon.CreateNodeTypeToNodeMap(env.ClCluster, nodesToCreate)
	vrfKey, pubKeyCompressed, err := vrfcommon.CreateVRFKeyOnVRFNode(nodeTypeToNodeMap[vrfcommon.VRF], l)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	l.Info().Str("Coordinator", vrfContracts.CoordinatorV2.Address()).Msg("Registering Proving Key")
	provingKey, err := VRFV2RegisterProvingKey(vrfKey, registerProvingKeyAgainstAddress, vrfContracts.CoordinatorV2)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("%s, err %w", vrfcommon.ErrRegisteringProvingKey, err)
	}
	keyHash, err := vrfContracts.CoordinatorV2.HashOfKey(context.Background(), provingKey)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("%s, err %w", vrfcommon.ErrCreatingProvingKeyHash, err)
	}

	chainID := env.EVMClient.GetChainID()
	vrfTXKeyAddressStrings, vrfTXKeyAddresses, err := vrfcommon.CreateFundAndGetSendingKeys(
		env.EVMClient,
		nodeTypeToNodeMap[vrfcommon.VRF],
		*vrfv2TestConfig.GetCommonConfig().ChainlinkNodeFunding,
		numberOfTxKeysToCreate,
		chainID,
	)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	err = env.EVMClient.WaitForEvents()
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("%s, err %w", vrfcommon.ErrWaitTXsComplete, err)
	}

	nodeTypeToNodeMap[vrfcommon.VRF].TXKeyAddressStrings = vrfTXKeyAddressStrings

	var vrfOwnerConfig *vrfcommon.VRFOwnerConfig
	if useVRFOwner {
		err := setupVRFOwnerContract(env, vrfContracts, vrfTXKeyAddressStrings, vrfTXKeyAddresses, l)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		vrfOwnerConfig = &vrfcommon.VRFOwnerConfig{
			OwnerAddress: vrfContracts.VRFOwner.Address(),
			UseVRFOwner:  useVRFOwner,
		}
	} else {
		vrfOwnerConfig = &vrfcommon.VRFOwnerConfig{
			OwnerAddress: "",
			UseVRFOwner:  useVRFOwner,
		}
	}

	g := errgroup.Group{}
	if vrfNode, exists := nodeTypeToNodeMap[vrfcommon.VRF]; exists {
		g.Go(func() error {
			err := setupVRFNode(vrfContracts, chainID, configGeneral, pubKeyCompressed, vrfOwnerConfig, l, vrfNode)
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
				chainID,
				vrfContracts.CoordinatorV2.Address(),
				vrfContracts.BHS.Address(),
				*vrfv2TestConfig.GetCommonConfig().ChainlinkNodeFunding,
				l,
				bhsNode,
			)
			if err != nil {
				return err
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, nil, nil, nil, fmt.Errorf("VRF node setup ended up with an error: %w", err)
	}

	vrfKeyData := vrfcommon.VRFKeyData{
		VRFKey:            vrfKey,
		EncodedProvingKey: provingKey,
		KeyHash:           keyHash,
	}

	l.Info().Msg("VRFV2 environment setup is finished")
	return vrfContracts, subIDs, &vrfKeyData, nodeTypeToNodeMap, nil
}

func setupVRFNode(contracts *vrfcommon.VRFContracts, chainID *big.Int, vrfv2Config *testconfig.General, pubKeyCompressed string, vrfOwnerConfig *vrfcommon.VRFOwnerConfig, l zerolog.Logger, vrfNode *vrfcommon.VRFNode) error {
	vrfJobSpecConfig := vrfcommon.VRFJobSpecConfig{
		ForwardingAllowed:             *vrfv2Config.VRFJobForwardingAllowed,
		CoordinatorAddress:            contracts.CoordinatorV2.Address(),
		FromAddresses:                 vrfNode.TXKeyAddressStrings,
		EVMChainID:                    chainID.String(),
		MinIncomingConfirmations:      int(*vrfv2Config.MinimumConfirmations),
		PublicKey:                     pubKeyCompressed,
		EstimateGasMultiplier:         *vrfv2Config.VRFJobEstimateGasMultiplier,
		BatchFulfillmentEnabled:       *vrfv2Config.VRFJobBatchFulfillmentEnabled,
		BatchFulfillmentGasMultiplier: *vrfv2Config.VRFJobBatchFulfillmentGasMultiplier,
		PollPeriod:                    vrfv2Config.VRFJobPollPeriod.Duration,
		RequestTimeout:                vrfv2Config.VRFJobRequestTimeout.Duration,
		SimulationBlock:               vrfv2Config.VRFJobSimulationBlock,
		VRFOwnerConfig:                vrfOwnerConfig,
	}

	l.Info().Msg("Creating VRFV2 Job")
	vrfV2job, err := CreateVRFV2Job(
		vrfNode.CLNode.API,
		vrfJobSpecConfig,
	)
	if err != nil {
		return fmt.Errorf("%s, err %w", ErrCreateVRFV2Jobs, err)
	}
	vrfNode.Job = vrfV2job

	// this part is here because VRFv2 can work with only a specific key
	// [[EVM.KeySpecific]]
	//	Key = '...'
	nodeConfig := node.NewConfig(vrfNode.CLNode.NodeConfig,
		node.WithLogPollInterval(1*time.Second),
		node.WithVRFv2EVMEstimator(vrfNode.TXKeyAddressStrings, *vrfv2Config.CLNodeMaxGasPriceGWei),
	)
	l.Info().Msg("Restarting Node with new sending key PriceMax configuration")
	err = vrfNode.CLNode.Restart(nodeConfig)
	if err != nil {
		return fmt.Errorf("%s, err %w", vrfcommon.ErrRestartCLNode, err)
	}
	return nil
}

func SetupVRFV2Contracts(
	env *test_env.CLClusterTestEnv,
	linkToken contracts.LinkToken,
	mockNativeLINKFeed contracts.MockETHLINKFeed,
	numberOfConsumers int,
	useVRFOwner bool,
	useTestCoordinator bool,
	vrfv2Config *testconfig.General,
	numberOfSubToCreate int,
	l zerolog.Logger,
) (*vrfcommon.VRFContracts, []uint64, error) {
	l.Info().Msg("Deploying VRFV2 contracts")
	vrfContracts, err := DeployVRFV2Contracts(
		env,
		linkToken,
		mockNativeLINKFeed,
		numberOfConsumers,
		useVRFOwner,
		useTestCoordinator,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("%s, err %w", ErrDeployVRFV2Contracts, err)
	}

	vrfCoordinatorV2FeeConfig := vrf_coordinator_v2.VRFCoordinatorV2FeeConfig{
		FulfillmentFlatFeeLinkPPMTier1: *vrfv2Config.FulfillmentFlatFeeLinkPPMTier1,
		FulfillmentFlatFeeLinkPPMTier2: *vrfv2Config.FulfillmentFlatFeeLinkPPMTier2,
		FulfillmentFlatFeeLinkPPMTier3: *vrfv2Config.FulfillmentFlatFeeLinkPPMTier3,
		FulfillmentFlatFeeLinkPPMTier4: *vrfv2Config.FulfillmentFlatFeeLinkPPMTier4,
		FulfillmentFlatFeeLinkPPMTier5: *vrfv2Config.FulfillmentFlatFeeLinkPPMTier5,
		ReqsForTier2:                   big.NewInt(*vrfv2Config.ReqsForTier2),
		ReqsForTier3:                   big.NewInt(*vrfv2Config.ReqsForTier3),
		ReqsForTier4:                   big.NewInt(*vrfv2Config.ReqsForTier4),
		ReqsForTier5:                   big.NewInt(*vrfv2Config.ReqsForTier5)}

	l.Info().Str("Coordinator", vrfContracts.CoordinatorV2.Address()).Msg("Setting Coordinator Config")
	err = vrfContracts.CoordinatorV2.SetConfig(
		*vrfv2Config.MinimumConfirmations,
		*vrfv2Config.MaxGasLimitCoordinatorConfig,
		*vrfv2Config.StalenessSeconds,
		*vrfv2Config.GasAfterPaymentCalculation,
		big.NewInt(*vrfv2Config.FallbackWeiPerUnitLink),
		vrfCoordinatorV2FeeConfig,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("%s, err %w", vrfcommon.ErrSetVRFCoordinatorConfig, err)
	}
	err = env.EVMClient.WaitForEvents()
	if err != nil {
		return nil, nil, fmt.Errorf("%s, err %w", vrfcommon.ErrWaitTXsComplete, err)
	}
	l.Info().
		Str("Coordinator", vrfContracts.CoordinatorV2.Address()).
		Int("Number of Subs to create", numberOfSubToCreate).
		Msg("Creating and funding subscriptions, adding consumers")
	subIDs, err := CreateFundSubsAndAddConsumers(
		env,
		big.NewFloat(*vrfv2Config.SubscriptionFundingAmountLink),
		linkToken,
		vrfContracts.CoordinatorV2, vrfContracts.VRFV2Consumer, numberOfSubToCreate)
	if err != nil {
		return nil, nil, err
	}
	return vrfContracts, subIDs, nil
}

func setupVRFOwnerContract(env *test_env.CLClusterTestEnv, contracts *vrfcommon.VRFContracts, allNativeTokenKeyAddressStrings []string, allNativeTokenKeyAddresses []common.Address, l zerolog.Logger) error {
	l.Info().Msg("Setting up VRFOwner contract")
	l.Info().
		Str("Coordinator", contracts.CoordinatorV2.Address()).
		Str("VRFOwner", contracts.VRFOwner.Address()).
		Msg("Transferring ownership of Coordinator to VRFOwner")
	err := contracts.CoordinatorV2.TransferOwnership(common.HexToAddress(contracts.VRFOwner.Address()))
	if err != nil {
		return nil
	}
	err = env.EVMClient.WaitForEvents()
	if err != nil {
		return nil
	}
	l.Info().
		Str("VRFOwner", contracts.VRFOwner.Address()).
		Msg("Accepting VRF Ownership")
	err = contracts.VRFOwner.AcceptVRFOwnership()
	if err != nil {
		return nil
	}
	err = env.EVMClient.WaitForEvents()
	if err != nil {
		return nil
	}
	l.Info().
		Strs("Authorized Senders", allNativeTokenKeyAddressStrings).
		Str("VRFOwner", contracts.VRFOwner.Address()).
		Msg("Setting authorized senders for VRFOwner contract")
	err = contracts.VRFOwner.SetAuthorizedSenders(allNativeTokenKeyAddresses)
	if err != nil {
		return nil
	}
	err = env.EVMClient.WaitForEvents()
	if err != nil {
		return fmt.Errorf("%s, err %w", vrfcommon.ErrWaitTXsComplete, err)
	}
	return err
}

func SetupVRFV2WrapperEnvironment(
	env *test_env.CLClusterTestEnv,
	vrfv2TestConfig tc.VRFv2TestConfig,
	linkToken contracts.LinkToken,
	mockNativeLINKFeed contracts.MockETHLINKFeed,
	coordinator contracts.VRFCoordinatorV2,
	keyHash [32]byte,
	wrapperConsumerContractsAmount int,
) (*VRFV2WrapperContracts, *uint64, error) {
	// Deploy VRF v2 direct funding contracts
	wrapperContracts, err := DeployVRFV2DirectFundingContracts(
		env.ContractDeployer,
		env.EVMClient,
		linkToken.Address(),
		mockNativeLINKFeed.Address(),
		coordinator,
		wrapperConsumerContractsAmount,
	)
	if err != nil {
		return nil, nil, err
	}
	err = env.EVMClient.WaitForEvents()
	if err != nil {
		return nil, nil, fmt.Errorf("%s, err %w", vrfcommon.ErrWaitTXsComplete, err)
	}

	vrfv2Config := vrfv2TestConfig.GetVRFv2Config()

	// Configure VRF v2 wrapper contract
	err = wrapperContracts.VRFV2Wrapper.SetConfig(
		*vrfv2Config.General.WrapperGasOverhead,
		*vrfv2Config.General.CoordinatorGasOverhead,
		*vrfv2Config.General.WrapperPremiumPercentage,
		keyHash,
		*vrfv2Config.General.WrapperMaxNumberOfWords,
	)
	if err != nil {
		return nil, nil, err
	}
	err = env.EVMClient.WaitForEvents()
	if err != nil {
		return nil, nil, fmt.Errorf("%s, err %w", vrfcommon.ErrWaitTXsComplete, err)
	}

	// Fetch wrapper subscription ID
	wrapperSubID, err := wrapperContracts.VRFV2Wrapper.GetSubID(context.Background())
	if err != nil {
		return nil, nil, err
	}
	err = env.EVMClient.WaitForEvents()
	if err != nil {
		return nil, nil, fmt.Errorf("%s, err %w", vrfcommon.ErrWaitTXsComplete, err)
	}

	// Fund wrapper subscription
	err = FundSubscriptions(env, big.NewFloat(*vrfv2Config.General.SubscriptionFundingAmountLink), linkToken, coordinator, []uint64{wrapperSubID})
	if err != nil {
		return nil, nil, err
	}

	// Fund consumer with LINK
	err = linkToken.Transfer(
		wrapperContracts.LoadTestConsumers[0].Address(),
		big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(*vrfv2Config.General.WrapperConsumerFundingAmountLink)),
	)
	if err != nil {
		return nil, nil, err
	}
	err = env.EVMClient.WaitForEvents()
	if err != nil {
		return nil, nil, fmt.Errorf("%s, err %w", vrfcommon.ErrWaitTXsComplete, err)
	}

	return wrapperContracts, &wrapperSubID, nil
}

func CreateFundSubsAndAddConsumers(
	env *test_env.CLClusterTestEnv,
	subscriptionFundingAmountLink *big.Float,
	linkToken contracts.LinkToken,
	coordinator contracts.VRFCoordinatorV2,
	consumers []contracts.VRFv2LoadTestConsumer,
	numberOfSubToCreate int,
) ([]uint64, error) {
	subIDs, err := CreateSubsAndFund(env, subscriptionFundingAmountLink, linkToken, coordinator, numberOfSubToCreate)
	if err != nil {
		return nil, err
	}
	subToConsumersMap := map[uint64][]contracts.VRFv2LoadTestConsumer{}

	//each subscription will have the same consumers
	for _, subID := range subIDs {
		subToConsumersMap[subID] = consumers
	}

	err = AddConsumersToSubs(
		subToConsumersMap,
		coordinator,
	)
	if err != nil {
		return nil, err
	}

	err = env.EVMClient.WaitForEvents()
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrWaitTXsComplete, err)
	}
	return subIDs, nil
}

func CreateSubsAndFund(
	env *test_env.CLClusterTestEnv,
	subscriptionFundingAmountLink *big.Float,
	linkToken contracts.LinkToken,
	coordinator contracts.VRFCoordinatorV2,
	subAmountToCreate int,
) ([]uint64, error) {
	subs, err := CreateSubs(env, coordinator, subAmountToCreate)
	if err != nil {
		return nil, err
	}
	err = env.EVMClient.WaitForEvents()
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrWaitTXsComplete, err)
	}
	err = FundSubscriptions(env, subscriptionFundingAmountLink, linkToken, coordinator, subs)
	if err != nil {
		return nil, err
	}
	return subs, nil
}

func CreateSubs(
	env *test_env.CLClusterTestEnv,
	coordinator contracts.VRFCoordinatorV2,
	subAmountToCreate int,
) ([]uint64, error) {
	var subIDArr []uint64

	for i := 0; i < subAmountToCreate; i++ {
		subID, err := CreateSubAndFindSubID(env, coordinator)
		if err != nil {
			return nil, err
		}
		subIDArr = append(subIDArr, subID)
	}
	return subIDArr, nil
}

func AddConsumersToSubs(
	subToConsumerMap map[uint64][]contracts.VRFv2LoadTestConsumer,
	coordinator contracts.VRFCoordinatorV2,
) error {
	for subID, consumers := range subToConsumerMap {
		for _, consumer := range consumers {
			err := coordinator.AddConsumer(subID, consumer.Address())
			if err != nil {
				return fmt.Errorf("%s, err %w", vrfcommon.ErrAddConsumerToSub, err)
			}
		}
	}
	return nil
}

func CreateSubAndFindSubID(env *test_env.CLClusterTestEnv, coordinator contracts.VRFCoordinatorV2) (uint64, error) {
	tx, err := coordinator.CreateSubscription()
	if err != nil {
		return 0, fmt.Errorf("%s, err %w", vrfcommon.ErrCreateVRFSubscription, err)
	}
	err = env.EVMClient.WaitForEvents()
	if err != nil {
		return 0, fmt.Errorf("%s, err %w", vrfcommon.ErrWaitTXsComplete, err)
	}

	receipt, err := env.EVMClient.GetTxReceipt(tx.Hash())
	if err != nil {
		return 0, fmt.Errorf("%s, err %w", vrfcommon.ErrWaitTXsComplete, err)
	}

	//SubscriptionsCreated Log should be emitted with the subscription ID
	subID := receipt.Logs[0].Topics[1].Big().Uint64()

	return subID, nil
}

func FundSubscriptions(
	env *test_env.CLClusterTestEnv,
	subscriptionFundingAmountLink *big.Float,
	linkAddress contracts.LinkToken,
	coordinator contracts.VRFCoordinatorV2,
	subIDs []uint64,
) error {
	for _, subID := range subIDs {
		//Link Billing
		amountJuels := conversions.EtherToWei(subscriptionFundingAmountLink)
		err := FundVRFCoordinatorV2Subscription(linkAddress, coordinator, env.EVMClient, subID, amountJuels)
		if err != nil {
			return fmt.Errorf("%s, err %w", vrfcommon.ErrFundSubWithLinkToken, err)
		}
	}
	err := env.EVMClient.WaitForEvents()
	if err != nil {
		return fmt.Errorf("%s, err %w", vrfcommon.ErrWaitTXsComplete, err)
	}
	return nil
}

func DirectFundingRequestRandomnessAndWaitForFulfillment(
	l zerolog.Logger,
	consumer contracts.VRFv2WrapperLoadTestConsumer,
	coordinator contracts.VRFCoordinatorV2,
	subID uint64,
	vrfv2KeyData *vrfcommon.VRFKeyData,
	minimumConfirmations uint16,
	callbackGasLimit uint32,
	numberOfWords uint32,
	randomnessRequestCountPerRequest uint16,
	randomnessRequestCountPerRequestDeviation uint16,
	randomWordsFulfilledEventTimeout time.Duration,
) (*vrf_coordinator_v2.VRFCoordinatorV2RandomWordsFulfilled, error) {
	logRandRequest(l, consumer.Address(), coordinator.Address(), subID, minimumConfirmations, callbackGasLimit, numberOfWords, randomnessRequestCountPerRequest, randomnessRequestCountPerRequestDeviation)
	_, err := consumer.RequestRandomness(
		minimumConfirmations,
		callbackGasLimit,
		numberOfWords,
		randomnessRequestCountPerRequest,
	)
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrRequestRandomness, err)
	}
	wrapperAddress, err := consumer.GetWrapper(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error getting wrapper address, err: %w", err)
	}
	fulfillmentEvents, err := WaitForRequestAndFulfillmentEvents(
		wrapperAddress.String(),
		coordinator,
		vrfv2KeyData,
		subID,
		randomWordsFulfilledEventTimeout,
		l,
	)
	return fulfillmentEvents, err
}

func RequestRandomnessAndWaitForFulfillment(
	l zerolog.Logger,
	consumer contracts.VRFv2LoadTestConsumer,
	coordinator contracts.VRFCoordinatorV2,
	subID uint64,
	vrfKeyData *vrfcommon.VRFKeyData,
	minimumConfirmations uint16,
	callbackGasLimit uint32,
	numberOfWords uint32,
	randomnessRequestCountPerRequest uint16,
	randomnessRequestCountPerRequestDeviation uint16,
	randomWordsFulfilledEventTimeout time.Duration,
) (*vrf_coordinator_v2.VRFCoordinatorV2RandomWordsFulfilled, error) {
	logRandRequest(l, consumer.Address(), coordinator.Address(), subID, minimumConfirmations, callbackGasLimit, numberOfWords, randomnessRequestCountPerRequest, randomnessRequestCountPerRequestDeviation)
	_, err := consumer.RequestRandomness(
		vrfKeyData.KeyHash,
		subID,
		minimumConfirmations,
		callbackGasLimit,
		numberOfWords,
		randomnessRequestCountPerRequest,
	)
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrRequestRandomness, err)
	}

	fulfillmentEvents, err := WaitForRequestAndFulfillmentEvents(
		consumer.Address(),
		coordinator,
		vrfKeyData,
		subID,
		randomWordsFulfilledEventTimeout,
		l,
	)
	return fulfillmentEvents, err
}

func RequestRandomnessWithForceFulfillAndWaitForFulfillment(
	l zerolog.Logger,
	consumer contracts.VRFv2LoadTestConsumer,
	coordinator contracts.VRFCoordinatorV2,
	vrfOwner contracts.VRFOwner,
	vrfv2KeyData *vrfcommon.VRFKeyData,
	minimumConfirmations uint16,
	callbackGasLimit uint32,
	numberOfWords uint32,
	randomnessRequestCountPerRequest uint16,
	randomnessRequestCountPerRequestDeviation uint16,
	subTopUpAmount *big.Int,
	linkAddress common.Address,
	randomWordsFulfilledEventTimeout time.Duration,
) (*vrf_coordinator_v2.VRFCoordinatorV2ConfigSet, *vrf_coordinator_v2.VRFCoordinatorV2RandomWordsFulfilled, *vrf_owner.VRFOwnerRandomWordsForced, error) {
	logRandRequest(l, consumer.Address(), coordinator.Address(), 0, minimumConfirmations, callbackGasLimit, numberOfWords, randomnessRequestCountPerRequest, randomnessRequestCountPerRequestDeviation)
	_, err := consumer.RequestRandomWordsWithForceFulfill(
		vrfv2KeyData.KeyHash,
		minimumConfirmations,
		callbackGasLimit,
		numberOfWords,
		randomnessRequestCountPerRequest,
		subTopUpAmount,
		linkAddress,
	)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%s, err %w", vrfcommon.ErrRequestRandomness, err)
	}

	randomWordsRequestedEvent, err := coordinator.WaitForRandomWordsRequestedEvent(
		[][32]byte{vrfv2KeyData.KeyHash},
		nil,
		[]common.Address{common.HexToAddress(consumer.Address())},
		time.Minute*1,
	)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%s, err %w", vrfcommon.ErrWaitRandomWordsRequestedEvent, err)
	}
	LogRandomnessRequestedEvent(l, coordinator, randomWordsRequestedEvent)

	errorChannel := make(chan error)
	configSetEventChannel := make(chan *vrf_coordinator_v2.VRFCoordinatorV2ConfigSet)
	randWordsFulfilledEventChannel := make(chan *vrf_coordinator_v2.VRFCoordinatorV2RandomWordsFulfilled)
	randWordsForcedEventChannel := make(chan *vrf_owner.VRFOwnerRandomWordsForced)

	go func() {
		configSetEvent, err := coordinator.WaitForConfigSetEvent(
			randomWordsFulfilledEventTimeout,
		)
		if err != nil {
			l.Error().Err(err).Msg("error waiting for ConfigSetEvent")
			errorChannel <- err
		}
		configSetEventChannel <- configSetEvent
	}()

	go func() {
		randomWordsFulfilledEvent, err := coordinator.WaitForRandomWordsFulfilledEvent(
			[]*big.Int{randomWordsRequestedEvent.RequestId},
			randomWordsFulfilledEventTimeout,
		)
		if err != nil {
			l.Error().Err(err).Msg("error waiting for RandomWordsFulfilledEvent")
			errorChannel <- err
		}
		randWordsFulfilledEventChannel <- randomWordsFulfilledEvent
	}()

	go func() {
		randomWordsForcedEvent, err := vrfOwner.WaitForRandomWordsForcedEvent(
			[]*big.Int{randomWordsRequestedEvent.RequestId},
			nil,
			nil,
			randomWordsFulfilledEventTimeout,
		)
		if err != nil {
			l.Error().Err(err).Msg("error waiting for RandomWordsForcedEvent")
			errorChannel <- err
		}
		randWordsForcedEventChannel <- randomWordsForcedEvent
	}()

	var configSetEvent *vrf_coordinator_v2.VRFCoordinatorV2ConfigSet
	var randomWordsFulfilledEvent *vrf_coordinator_v2.VRFCoordinatorV2RandomWordsFulfilled
	var randomWordsForcedEvent *vrf_owner.VRFOwnerRandomWordsForced
	for i := 0; i < 3; i++ {
		select {
		case err = <-errorChannel:
			return nil, nil, nil, err
		case configSetEvent = <-configSetEventChannel:
		case randomWordsFulfilledEvent = <-randWordsFulfilledEventChannel:
			LogRandomWordsFulfilledEvent(l, coordinator, randomWordsFulfilledEvent)
		case randomWordsForcedEvent = <-randWordsForcedEventChannel:
			LogRandomWordsForcedEvent(l, vrfOwner, randomWordsForcedEvent)
		case <-time.After(randomWordsFulfilledEventTimeout):
			err = fmt.Errorf("timeout waiting for ConfigSet, RandomWordsFulfilled and RandomWordsForced events")
		}
	}
	return configSetEvent, randomWordsFulfilledEvent, randomWordsForcedEvent, err
}

func WaitForRequestAndFulfillmentEvents(
	consumerAddress string,
	coordinator contracts.VRFCoordinatorV2,
	vrfv2KeyData *vrfcommon.VRFKeyData,
	subID uint64,
	randomWordsFulfilledEventTimeout time.Duration,
	l zerolog.Logger,
) (*vrf_coordinator_v2.VRFCoordinatorV2RandomWordsFulfilled, error) {
	randomWordsRequestedEvent, err := coordinator.WaitForRandomWordsRequestedEvent(
		[][32]byte{vrfv2KeyData.KeyHash},
		[]uint64{subID},
		[]common.Address{common.HexToAddress(consumerAddress)},
		time.Minute*1,
	)
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrWaitRandomWordsRequestedEvent, err)
	}
	LogRandomnessRequestedEvent(l, coordinator, randomWordsRequestedEvent)

	randomWordsFulfilledEvent, err := coordinator.WaitForRandomWordsFulfilledEvent(
		[]*big.Int{randomWordsRequestedEvent.RequestId},
		randomWordsFulfilledEventTimeout,
	)
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrWaitRandomWordsFulfilledEvent, err)
	}
	LogRandomWordsFulfilledEvent(l, coordinator, randomWordsFulfilledEvent)
	return randomWordsFulfilledEvent, err
}

func LogSubDetails(l zerolog.Logger, subscription vrf_coordinator_v2.GetSubscription, subID uint64, coordinator contracts.VRFCoordinatorV2) {
	l.Debug().
		Str("Coordinator", coordinator.Address()).
		Str("Link Balance", (*commonassets.Link)(subscription.Balance).Link()).
		Uint64("Subscription ID", subID).
		Str("Subscription Owner", subscription.Owner.String()).
		Interface("Subscription Consumers", subscription.Consumers).
		Msg("Subscription Data")
}

func LogRandomnessRequestedEvent(
	l zerolog.Logger,
	coordinator contracts.VRFCoordinatorV2,
	randomWordsRequestedEvent *vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested,
) {
	l.Info().
		Str("Coordinator", coordinator.Address()).
		Str("Request ID", randomWordsRequestedEvent.RequestId.String()).
		Uint64("Subscription ID", randomWordsRequestedEvent.SubId).
		Str("Sender Address", randomWordsRequestedEvent.Sender.String()).
		Str("Keyhash", fmt.Sprintf("0x%x", randomWordsRequestedEvent.KeyHash)).
		Uint32("Callback Gas Limit", randomWordsRequestedEvent.CallbackGasLimit).
		Uint32("Number of Words", randomWordsRequestedEvent.NumWords).
		Uint16("Minimum Request Confirmations", randomWordsRequestedEvent.MinimumRequestConfirmations).
		Str("TX Hash", randomWordsRequestedEvent.Raw.TxHash.String()).
		Uint64("BlockNumber", randomWordsRequestedEvent.Raw.BlockNumber).
		Str("BlockHash", randomWordsRequestedEvent.Raw.BlockHash.String()).
		Msg("RandomnessRequested Event")
}

func LogRandomWordsFulfilledEvent(
	l zerolog.Logger,
	coordinator contracts.VRFCoordinatorV2,
	randomWordsFulfilledEvent *vrf_coordinator_v2.VRFCoordinatorV2RandomWordsFulfilled,
) {
	l.Info().
		Str("Coordinator", coordinator.Address()).
		Str("Total Payment", randomWordsFulfilledEvent.Payment.String()).
		Str("TX Hash", randomWordsFulfilledEvent.Raw.TxHash.String()).
		Str("Request ID", randomWordsFulfilledEvent.RequestId.String()).
		Bool("Success", randomWordsFulfilledEvent.Success).
		Uint64("BlockNumber", randomWordsFulfilledEvent.Raw.BlockNumber).
		Str("BlockHash", randomWordsFulfilledEvent.Raw.BlockHash.String()).
		Msg("RandomWordsFulfilled Event (TX metadata)")
}

func LogRandomWordsForcedEvent(
	l zerolog.Logger,
	vrfOwner contracts.VRFOwner,
	randomWordsForcedEvent *vrf_owner.VRFOwnerRandomWordsForced,
) {
	l.Debug().
		Str("VRFOwner", vrfOwner.Address()).
		Uint64("Sub ID", randomWordsForcedEvent.SubId).
		Str("TX Hash", randomWordsForcedEvent.Raw.TxHash.String()).
		Str("Request ID", randomWordsForcedEvent.RequestId.String()).
		Str("Sender", randomWordsForcedEvent.Sender.String()).
		Msg("RandomWordsForced Event (TX metadata)")
}

func logRandRequest(
	l zerolog.Logger,
	consumer string,
	coordinator string,
	subID uint64,
	minimumConfirmations uint16,
	callbackGasLimit uint32,
	numberOfWords uint32,
	randomnessRequestCountPerRequest uint16,
	randomnessRequestCountPerRequestDeviation uint16,
) {
	l.Info().
		Str("Consumer", consumer).
		Str("Coordinator", coordinator).
		Uint64("SubID", subID).
		Uint16("MinimumConfirmations", minimumConfirmations).
		Uint32("CallbackGasLimit", callbackGasLimit).
		Uint32("NumberOfWords", numberOfWords).
		Uint16("RandomnessRequestCountPerRequest", randomnessRequestCountPerRequest).
		Uint16("RandomnessRequestCountPerRequestDeviation", randomnessRequestCountPerRequestDeviation).
		Msg("Requesting randomness")
}
