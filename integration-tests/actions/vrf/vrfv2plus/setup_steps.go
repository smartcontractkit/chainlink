package vrfv2plus

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-testing-framework/utils/conversions"
	vrfcommon "github.com/smartcontractkit/chainlink/integration-tests/actions/vrf/common"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
	"github.com/smartcontractkit/chainlink/integration-tests/types/config/node"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
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
		return nil, fmt.Errorf("%s, err %w", vrfcommon.ErrParseJob, err)
	}

	job, err := chainlinkNode.MustCreateJob(&client.VRFV2PlusJobSpec{
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
	})
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", ErrCreatingVRFv2PlusJob, err)
	}

	return job, nil
}

// SetupVRFV2_5Environment will create specified number of subscriptions and add the same conumer/s to each of them
func SetupVRFV2_5Environment(
	env *test_env.CLClusterTestEnv,
	nodesToCreate []vrfcommon.VRFNodeType,
	vrfv2PlusTestConfig types.VRFv2PlusTestConfig,
	linkToken contracts.LinkToken,
	mockNativeLINKFeed contracts.MockETHLINKFeed,
	numberOfTxKeysToCreate int,
	numberOfConsumers int,
	numberOfSubToCreate int,
	l zerolog.Logger,
) (*vrfcommon.VRFContracts, []*big.Int, *vrfcommon.VRFKeyData, map[vrfcommon.VRFNodeType]*vrfcommon.VRFNode, error) {
	l.Info().Msg("Starting VRFV2 Plus environment setup")
	l.Info().Msg("Starting VRFV2 Plus environment setup")
	configGeneral := vrfv2PlusTestConfig.GetVRFv2PlusConfig().General
	vrfContracts, subIDs, err := SetupVRFV2PlusContracts(
		env,
		linkToken,
		mockNativeLINKFeed,
		configGeneral,
		numberOfSubToCreate,
		numberOfConsumers,
		l,
	)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	nodeTypeToNodeMap := vrfcommon.CreateNodeTypeToNodeMap(env.ClCluster, nodesToCreate)
	vrfKey, pubKeyCompressed, err := vrfcommon.CreateVRFKeyOnVRFNode(nodeTypeToNodeMap[vrfcommon.VRF], l)

	l.Info().Str("Coordinator", vrfContracts.CoordinatorV2Plus.Address()).Msg("Registering Proving Key")
	provingKey, err := VRFV2_5RegisterProvingKey(vrfKey, vrfContracts.CoordinatorV2Plus, uint64(assets.GWei(*configGeneral.CLNodeMaxGasPriceGWei).Int64()))
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("%s, err %w", vrfcommon.ErrRegisteringProvingKey, err)
	}
	keyHash, err := vrfContracts.CoordinatorV2Plus.HashOfKey(context.Background(), provingKey)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("%s, err %w", vrfcommon.ErrCreatingProvingKeyHash, err)
	}

	chainID := env.EVMClient.GetChainID()
	vrfTXKeyAddressStrings, _, err := vrfcommon.CreateFundAndGetSendingKeys(
		env.EVMClient,
		nodeTypeToNodeMap[vrfcommon.VRF],
		*vrfv2PlusTestConfig.GetCommonConfig().ChainlinkNodeFunding,
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

	g := errgroup.Group{}
	if vrfNode, exists := nodeTypeToNodeMap[vrfcommon.VRF]; exists {
		g.Go(func() error {
			err := setupVRFNode(vrfContracts, chainID, configGeneral, pubKeyCompressed, l, vrfNode)
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

	if err := g.Wait(); err != nil {
		return nil, nil, nil, nil, fmt.Errorf("VRF node setup ended up with an error: %w", err)
	}

	vrfKeyData := vrfcommon.VRFKeyData{
		VRFKey:            vrfKey,
		EncodedProvingKey: provingKey,
		KeyHash:           keyHash,
	}

	l.Info().Msg("VRFV2 Plus environment setup is finished")
	return vrfContracts, subIDs, &vrfKeyData, nodeTypeToNodeMap, nil
}

func setupVRFNode(contracts *vrfcommon.VRFContracts, chainID *big.Int, config *vrfv2plus_config.General, pubKeyCompressed string, l zerolog.Logger, vrfNode *vrfcommon.VRFNode) error {
	vrfJobSpecConfig := vrfcommon.VRFJobSpecConfig{
		ForwardingAllowed:             *config.VRFJobForwardingAllowed,
		CoordinatorAddress:            contracts.CoordinatorV2Plus.Address(),
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
		return fmt.Errorf("%s, err %w", ErrCreateVRFV2PlusJobs, err)
	}
	vrfNode.Job = job

	// this part is here because VRFv2 can work with only a specific key
	// [[EVM.KeySpecific]]
	//	Key = '...'
	nodeConfig := node.NewConfig(vrfNode.CLNode.NodeConfig,
		node.WithLogPollInterval(1*time.Second),
		node.WithVRFv2EVMEstimator(vrfNode.TXKeyAddressStrings, *config.CLNodeMaxGasPriceGWei),
	)
	l.Info().Msg("Restarting Node with new sending key PriceMax configuration")
	err = vrfNode.CLNode.Restart(nodeConfig)
	if err != nil {
		return fmt.Errorf("%s, err %w", vrfcommon.ErrRestartCLNode, err)
	}
	return nil
}

func SetupVRFV2PlusWrapperEnvironment(
	env *test_env.CLClusterTestEnv,
	vrfv2PlusTestConfig types.VRFv2PlusTestConfig,
	linkToken contracts.LinkToken,
	mockNativeLINKFeed contracts.MockETHLINKFeed,
	coordinator contracts.VRFCoordinatorV2_5,
	keyHash [32]byte,
	wrapperConsumerContractsAmount int,
) (*VRFV2PlusWrapperContracts, *big.Int, error) {
	vrfv2PlusConfig := vrfv2PlusTestConfig.GetVRFv2PlusConfig().General
	wrapperContracts, err := DeployVRFV2PlusDirectFundingContracts(
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
	err = wrapperContracts.VRFV2PlusWrapper.SetConfig(
		*vrfv2PlusConfig.WrapperGasOverhead,
		*vrfv2PlusConfig.CoordinatorGasOverhead,
		*vrfv2PlusConfig.WrapperPremiumPercentage,
		keyHash,
		*vrfv2PlusConfig.WrapperMaxNumberOfWords,
		*vrfv2PlusConfig.StalenessSeconds,
		big.NewInt(*vrfv2PlusConfig.FallbackWeiPerUnitLink),
		*vrfv2PlusConfig.FulfillmentFlatFeeLinkPPM,
		*vrfv2PlusConfig.FulfillmentFlatFeeNativePPM,
	)
	if err != nil {
		return nil, nil, err
	}

	err = env.EVMClient.WaitForEvents()
	if err != nil {
		return nil, nil, fmt.Errorf("%s, err %w", vrfcommon.ErrWaitTXsComplete, err)
	}

	//fund sub
	wrapperSubID, err := wrapperContracts.VRFV2PlusWrapper.GetSubID(context.Background())
	if err != nil {
		return nil, nil, err
	}

	err = env.EVMClient.WaitForEvents()
	if err != nil {
		return nil, nil, fmt.Errorf("%s, err %w", vrfcommon.ErrWaitTXsComplete, err)
	}

	err = FundSubscriptions(
		env,
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
	err = env.EVMClient.WaitForEvents()
	if err != nil {
		return nil, nil, fmt.Errorf("%s, err %w", vrfcommon.ErrWaitTXsComplete, err)
	}

	//fund consumer with Eth
	err = wrapperContracts.LoadTestConsumers[0].Fund(big.NewFloat(*vrfv2PlusConfig.WrapperConsumerFundingAmountNativeToken))
	if err != nil {
		return nil, nil, err
	}
	err = env.EVMClient.WaitForEvents()
	if err != nil {
		return nil, nil, fmt.Errorf("%s, err %w", vrfcommon.ErrWaitTXsComplete, err)
	}
	return wrapperContracts, wrapperSubID, nil
}

func ReturnFundsForFulfilledRequests(client blockchain.EVMClient, coordinator contracts.VRFCoordinatorV2_5, l zerolog.Logger) error {
	linkTotalBalance, err := coordinator.GetLinkTotalBalance(context.Background())
	if err != nil {
		return fmt.Errorf("Error getting LINK total balance, err: %w", err)
	}
	defaultWallet := client.GetDefaultWallet().Address()
	l.Info().
		Str("LINK amount", linkTotalBalance.String()).
		Str("Returning to", defaultWallet).
		Msg("Returning LINK for fulfilled requests")
	err = coordinator.Withdraw(
		common.HexToAddress(defaultWallet),
	)
	if err != nil {
		return fmt.Errorf("Error withdrawing LINK from coordinator to default wallet, err: %w", err)
	}
	nativeTotalBalance, err := coordinator.GetNativeTokenTotalBalance(context.Background())
	if err != nil {
		return fmt.Errorf("Error getting NATIVE total balance, err: %w", err)
	}
	l.Info().
		Str("Native Token amount", nativeTotalBalance.String()).
		Str("Returning to", defaultWallet).
		Msg("Returning Native Token for fulfilled requests")
	err = coordinator.WithdrawNative(
		common.HexToAddress(defaultWallet),
	)
	if err != nil {
		return fmt.Errorf("Error withdrawing NATIVE from coordinator to default wallet, err: %w", err)
	}
	return nil
}

func SetupVRFV2PlusUniverse(t *testing.T, testConfig tc.TestConfig, cleanupFn func(), newEnvConfig NewEnvConfig, l zerolog.Logger) (*test_env.CLClusterTestEnv, *vrfcommon.VRFContracts, []*big.Int, *vrfcommon.VRFKeyData, map[vrfcommon.VRFNodeType]*vrfcommon.VRFNode, error) {
	var (
		env            *test_env.CLClusterTestEnv
		vrfContracts   *vrfcommon.VRFContracts
		vrfKey         *vrfcommon.VRFKeyData
		subIDs         []*big.Int
		nodeTypeToNode map[vrfcommon.VRFNodeType]*vrfcommon.VRFNode
		err            error
	)
	if *testConfig.VRFv2Plus.General.UseExistingEnv {
		vrfContracts, subIDs, vrfKey, env, err = SetupVRFV2PlusForExistingEnv(t, testConfig, cleanupFn, l)
		if err != nil {
			return nil, nil, nil, nil, nil, fmt.Errorf("%s, err: %w", "Error setting up VRF V2 Plus for Existing env", err)
		}
	} else {
		vrfContracts, subIDs, vrfKey, env, nodeTypeToNode, err = SetupVRFV2PlusForNewEnv(t, testConfig, cleanupFn, newEnvConfig, l)
		if err != nil {
			return nil, nil, nil, nil, nil, fmt.Errorf("%s, err: %w", "Error setting up VRF V2 Plus for New env", err)
		}
	}
	return env, vrfContracts, subIDs, vrfKey, nodeTypeToNode, nil
}

func SetupVRFV2PlusForNewEnv(
	t *testing.T,
	testConfig tc.TestConfig,
	cleanupFn func(),
	newEnvConfig NewEnvConfig,
	l zerolog.Logger) (*vrfcommon.VRFContracts, []*big.Int, *vrfcommon.VRFKeyData, *test_env.CLClusterTestEnv, map[vrfcommon.VRFNodeType]*vrfcommon.VRFNode, error) {
	network, err := actions.EthereumNetworkConfigFromConfig(l, &testConfig)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("%s, err: %w", "Error building ethereum network config", err)
	}
	env, err := test_env.NewCLTestEnvBuilder().
		WithTestInstance(t).
		WithTestConfig(&testConfig).
		WithPrivateEthereumNetwork(network).
		WithCLNodes(1).
		WithFunding(big.NewFloat(*testConfig.Common.ChainlinkNodeFunding)).
		WithCustomCleanup(cleanupFn).
		Build()
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("%s, err: %w", "error creating test env", err)
	}

	env.ParallelTransactions(true)

	mockETHLinkFeed, err := actions.DeployMockETHLinkFeed(env.ContractDeployer, big.NewInt(*testConfig.VRFv2Plus.General.LinkNativeFeedResponse))
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("%s, err: %w", "error deploying mock ETH/LINK feed", err)
	}

	linkToken, err := actions.DeployLINKToken(env.ContractDeployer)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("%s, err: %w", "error deploying LINK contract", err)
	}

	vrfContracts, subIDs, vrfKey, nodeTypeToNode, err := SetupVRFV2_5Environment(
		env,
		newEnvConfig.NodesToCreate,
		&testConfig,
		linkToken,
		mockETHLinkFeed,
		newEnvConfig.NumberOfTxKeysToCreate,
		newEnvConfig.NumberOfConsumers,
		newEnvConfig.NumberOfSubToCreate,
		l,
	)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("%s, err: %w", "error setting up VRF v2_5 env", err)
	}
	return vrfContracts, subIDs, vrfKey, env, nodeTypeToNode, nil
}

func SetupVRFV2PlusForExistingEnv(t *testing.T, testConfig tc.TestConfig, cleanupFn func(), l zerolog.Logger) (*vrfcommon.VRFContracts, []*big.Int, *vrfcommon.VRFKeyData, *test_env.CLClusterTestEnv, error) {
	var subIDs []*big.Int
	env, err := test_env.NewCLTestEnvBuilder().
		WithTestInstance(t).
		WithTestConfig(&testConfig).
		WithCustomCleanup(cleanupFn).
		Build()
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("%s, err: %w", "error creating test env", err)
	}

	coordinator, err := env.ContractLoader.LoadVRFCoordinatorV2_5(*testConfig.VRFv2Plus.ExistingEnvConfig.CoordinatorAddress)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("%s, err: %w", "error loading VRFCoordinator2_5", err)
	}

	var consumers []contracts.VRFv2PlusLoadTestConsumer
	if *testConfig.VRFv2Plus.ExistingEnvConfig.CreateFundSubsAndAddConsumers {
		linkToken, err := env.ContractLoader.LoadLINKToken(*testConfig.VRFv2Plus.ExistingEnvConfig.LinkAddress)
		if err != nil {
			return nil, nil, nil, nil, fmt.Errorf("%s, err: %w", "error loading LinkToken", err)
		}
		consumers, err = DeployVRFV2PlusConsumers(env.ContractDeployer, coordinator, 1)
		if err != nil {
			return nil, nil, nil, nil, fmt.Errorf("err: %w", err)
		}
		err = env.EVMClient.WaitForEvents()
		if err != nil {
			return nil, nil, nil, nil, fmt.Errorf("%s, err: %w", vrfcommon.ErrWaitTXsComplete, err)
		}
		l.Info().
			Str("Coordinator", *testConfig.VRFv2Plus.ExistingEnvConfig.CoordinatorAddress).
			Int("Number of Subs to create", *testConfig.VRFv2Plus.General.NumberOfSubToCreate).
			Msg("Creating and funding subscriptions, deploying and adding consumers to subs")
		subIDs, err = CreateFundSubsAndAddConsumers(
			env,
			big.NewFloat(*testConfig.GetVRFv2PlusConfig().General.SubscriptionFundingAmountNative),
			big.NewFloat(*testConfig.GetVRFv2PlusConfig().General.SubscriptionFundingAmountLink),
			linkToken,
			coordinator,
			consumers,
			*testConfig.VRFv2Plus.General.NumberOfSubToCreate,
		)
		if err != nil {
			return nil, nil, nil, nil, fmt.Errorf("err: %w", err)
		}
	} else {
		consumer, err := env.ContractLoader.LoadVRFv2PlusLoadTestConsumer(*testConfig.VRFv2Plus.ExistingEnvConfig.ConsumerAddress)
		if err != nil {
			return nil, nil, nil, nil, fmt.Errorf("err: %w", err)
		}
		consumers = append(consumers, consumer)
		var ok bool
		subID, ok := new(big.Int).SetString(*testConfig.VRFv2Plus.ExistingEnvConfig.SubID, 10)
		if !ok {
			return nil, nil, nil, nil, fmt.Errorf("unable to parse subID: %s %w", *testConfig.VRFv2Plus.ExistingEnvConfig.SubID, err)
		}
		subIDs = append(subIDs, subID)
	}

	err = FundNodesIfNeeded(&testConfig, env.EVMClient, l)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("err: %w", err)
	}

	vrfContracts := &vrfcommon.VRFContracts{
		CoordinatorV2Plus: coordinator,
		VRFV2PlusConsumer: consumers,
		BHS:               nil,
	}

	vrfKey := &vrfcommon.VRFKeyData{
		VRFKey:            nil,
		EncodedProvingKey: [2]*big.Int{},
		KeyHash:           common.HexToHash(*testConfig.VRFv2Plus.ExistingEnvConfig.KeyHash),
	}
	return vrfContracts, subIDs, vrfKey, env, nil
}

func CancelSubsAndReturnFunds(vrfContracts *vrfcommon.VRFContracts, eoaWalletAddress string, subIDs []*big.Int, l zerolog.Logger) {
	for _, subID := range subIDs {
		l.Info().
			Str("Returning funds from SubID", subID.String()).
			Str("Returning funds to", eoaWalletAddress).
			Msg("Canceling subscription and returning funds to subscription owner")
		pendingRequestsExist, err := vrfContracts.CoordinatorV2Plus.PendingRequestsExist(context.Background(), subID)
		if err != nil {
			l.Error().Err(err).Msg("Error checking if pending requests exist")
		}
		if !pendingRequestsExist {
			_, err := vrfContracts.CoordinatorV2Plus.CancelSubscription(subID, common.HexToAddress(eoaWalletAddress))
			if err != nil {
				l.Error().Err(err).Msg("Error canceling subscription")
			}
		} else {
			l.Error().Str("Sub ID", subID.String()).Msg("Pending requests exist for subscription, cannot cancel subscription and return funds")
		}
	}
}

func FundNodesIfNeeded(vrfv2plusTestConfig tc.VRFv2PlusTestConfig, client blockchain.EVMClient, l zerolog.Logger) error {
	cfg := vrfv2plusTestConfig.GetVRFv2PlusConfig()
	if *cfg.ExistingEnvConfig.NodeSendingKeyFundingMin > 0 {
		for _, sendingKey := range cfg.ExistingEnvConfig.NodeSendingKeys {
			address := common.HexToAddress(sendingKey)
			sendingKeyBalance, err := client.BalanceAt(context.Background(), address)
			if err != nil {
				return err
			}
			fundingAtLeast := conversions.EtherToWei(big.NewFloat(*cfg.ExistingEnvConfig.NodeSendingKeyFundingMin))
			fundingToSendWei := new(big.Int).Sub(fundingAtLeast, sendingKeyBalance)
			fundingToSendEth := conversions.WeiToEther(fundingToSendWei)
			if fundingToSendWei.Cmp(big.NewInt(0)) == 1 {
				l.Info().
					Str("Sending Key", sendingKey).
					Str("Sending Key Current Balance", sendingKeyBalance.String()).
					Str("Should have at least", fundingAtLeast.String()).
					Str("Funding Amount in ETH", fundingToSendEth.String()).
					Msg("Funding Node's Sending Key")
				err := actions.FundAddress(client, sendingKey, fundingToSendEth)
				if err != nil {
					return err
				}
			} else {
				l.Info().
					Str("Sending Key", sendingKey).
					Str("Sending Key Current Balance", sendingKeyBalance.String()).
					Str("Should have at least", fundingAtLeast.String()).
					Msg("Skipping Node's Sending Key funding as it has enough funds")
			}
		}
	}
	return nil
}
