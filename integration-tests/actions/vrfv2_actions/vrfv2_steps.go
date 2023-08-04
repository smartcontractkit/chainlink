package vrfv2_actions

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	eth "github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	chainlinkutils "github.com/smartcontractkit/chainlink/v2/core/utils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

var (
	ErrCreatingProvingKey   = "error creating a keyHash from the proving key"
	ErrDeployBlockHashStore = "error deploying blockhash store"
	ErrDeployCoordinator    = "error deploying VRFv2 CoordinatorV2"
	ErrAdvancedConsumer     = "error deploying VRFv2 Advanced Consumer"
	ErrABIEncodingFunding   = "error Abi encoding subscriptionID"
	ErrSendingLinkToken     = "error sending Link token"
)

func DeployVRFV2Contracts(
	contractDeployer contracts.ContractDeployer,
	chainClient blockchain.EVMClient,
	linkTokenContract contracts.LinkToken,
	linkEthFeedContract contracts.MockETHLINKFeed,
) (*VRFV2Contracts, error) {
	bhs, err := contractDeployer.DeployBlockhashStore()
	if err != nil {
		return nil, errors.Wrap(err, ErrDeployBlockHashStore)
	}
	coordinator, err := contractDeployer.DeployVRFCoordinatorV2(linkTokenContract.Address(), bhs.Address(), linkEthFeedContract.Address())
	if err != nil {
		return nil, errors.Wrap(err, ErrDeployCoordinator)
	}
	loadTestConsumer, err := contractDeployer.DeployVRFv2LoadTestConsumer(coordinator.Address())
	if err != nil {
		return nil, errors.Wrap(err, ErrAdvancedConsumer)
	}
	err = chainClient.WaitForEvents()
	if err != nil {
		return nil, err
	}
	return &VRFV2Contracts{coordinator, bhs, loadTestConsumer}, nil
}

func CreateVRFV2Jobs(
	chainlinkNodes []*client.ChainlinkClient,
	coordinator contracts.VRFCoordinatorV2,
	c blockchain.EVMClient,
	minIncomingConfirmations uint16,
) ([]VRFV2JobInfo, error) {
	jobInfo := make([]VRFV2JobInfo, 0)
	for _, chainlinkNode := range chainlinkNodes {
		vrfKey, err := chainlinkNode.MustCreateVRFKey()
		if err != nil {
			return nil, errors.Wrap(err, "error creating VRF key")
		}
		pubKeyCompressed := vrfKey.Data.ID
		jobUUID := uuid.New()
		os := &client.VRFV2TxPipelineSpec{
			Address: coordinator.Address(),
		}
		ost, err := os.String()
		if err != nil {
			return nil, errors.Wrap(err, "Error getting job string")
		}
		nativeTokenPrimaryKeyAddress, err := chainlinkNode.PrimaryEthAddress()
		if err != nil {
			return nil, errors.Wrap(err, "Error getting node's primary ETH key")
		}
		job, err := chainlinkNode.MustCreateJob(&client.VRFV2JobSpec{
			Name:                     fmt.Sprintf("vrf-%s", jobUUID),
			CoordinatorAddress:       coordinator.Address(),
			FromAddresses:            []string{nativeTokenPrimaryKeyAddress},
			EVMChainID:               c.GetChainID().String(),
			MinIncomingConfirmations: int(minIncomingConfirmations),
			PublicKey:                pubKeyCompressed,
			ExternalJobID:            jobUUID.String(),
			ObservationSource:        ost,
			BatchFulfillmentEnabled:  false,
		})
		if err != nil {
			return nil, errors.Wrap(err, "Error creating VRFv2 job")
		}
		provingKey, err := VRFV2RegisterProvingKey(vrfKey, nativeTokenPrimaryKeyAddress, coordinator)
		if err != nil {
			return nil, errors.Wrap(err, "Error creating VRFv2 proving key")
		}
		keyHash, err := coordinator.HashOfKey(context.Background(), provingKey)
		if err != nil {
			return nil, errors.Wrap(err, "Error creating a keyHash from the proving key")
		}
		ji := VRFV2JobInfo{
			Job:               job,
			VRFKey:            vrfKey,
			EncodedProvingKey: provingKey,
			KeyHash:           keyHash,
		}
		jobInfo = append(jobInfo, ji)
	}
	return jobInfo, nil
}

func VRFV2RegisterProvingKey(
	vrfKey *client.VRFKey,
	oracleAddress string,
	coordinator contracts.VRFCoordinatorV2,
) (VRFV2EncodedProvingKey, error) {
	provingKey, err := actions.EncodeOnChainVRFProvingKey(*vrfKey)
	if err != nil {
		return VRFV2EncodedProvingKey{}, errors.Wrap(err, "Error encoding proving key")
	}
	err = coordinator.RegisterProvingKey(
		oracleAddress,
		provingKey,
	)
	if err != nil {
		return VRFV2EncodedProvingKey{}, errors.Wrap(err, "Error registering proving keys")
	}
	return provingKey, nil
}

func FundVRFCoordinatorV2Subscription(linkToken contracts.LinkToken, coordinator contracts.VRFCoordinatorV2, chainClient blockchain.EVMClient, subscriptionID uint64, linkFundingAmount *big.Int) error {
	encodedSubId, err := chainlinkutils.ABIEncode(`[{"type":"uint64"}]`, subscriptionID)
	if err != nil {
		return errors.Wrap(err, "failed to encode ABI subscriptionID")
	}
	_, err = linkToken.TransferAndCall(coordinator.Address(), big.NewInt(0).Mul(linkFundingAmount, big.NewInt(1e18)), encodedSubId)
	if err != nil {
		return errors.Wrap(err, "failed to send LINK token")
	}
	return chainClient.WaitForEvents()
}

//func SetupVRFV2Universe(
//	t *testing.T,
//	linkToken contracts.LinkToken,
//	mockETHLinkFeed contracts.MockETHLINKFeed,
//	contractDeployer contracts.ContractDeployer,
//	chainClient blockchain.EVMClient,
//	chainlinkNodes []*client.ChainlinkK8sClient,
//	testNetwork blockchain.EVMNetwork,
//	existingTestEnvironment *environment.Environment,
//	chainlinkNodeFundingAmountEth *big.Float,
//	vrfSubscriptionFundingAmountInLink *big.Int,
//	newEnvNamespacePrefix string,
//	newEnvTTL time.Duration,
//) (VRFV2Contracts, []*client.ChainlinkK8sClient, []VRFV2JobInfo, *environment.Environment) {
//
//	vrfV2Contracts, err := DeployVRFV2Contracts(
//		t,
//		contractDeployer,
//		chainClient,
//		linkToken,
//		mockETHLinkFeed,
//	)
//
//	err := actions.FundChainlinkNodes(chainlinkNodes, chainClient, chainlinkNodeFundingAmountEth)
//	require.NoError(t, err)
//	err = chainClient.WaitForEvents()
//	require.NoError(t, err)
//
//	err = vrfV2Contracts.Coordinator.SetConfig(
//		uint16(vrfv2_constants.MinimumConfirmations),
//		vrfv2_constants.MaxGasLimitVRFCoordinatorConfig,
//		vrfv2_constants.StalenessSeconds,
//		vrfv2_constants.GasAfterPaymentCalculation,
//		vrfv2_constants.LinkEthFeedResponse,
//		vrfv2_constants.VRFCoordinatorV2FeeConfig,
//	)
//	require.NoError(t, err)
//	err = chainClient.WaitForEvents()
//	require.NoError(t, err)
//
//	err = vrfV2Contracts.Coordinator.CreateSubscription()
//	require.NoError(t, err)
//	err = chainClient.WaitForEvents()
//	require.NoError(t, err)
//
//	err = vrfV2Contracts.Coordinator.AddConsumer(vrfv2_constants.SubID, vrfV2Contracts.LoadTestConsumer.Address())
//	require.NoError(t, err, "Error adding a Load Test Consumer to a subscription in VRFCoordinator contract")
//
//	FundVRFCoordinatorV2Subscription(
//		t,
//		linkToken,
//		vrfV2Contracts.Coordinator,
//		chainClient,
//		vrfv2_constants.SubID,
//		vrfSubscriptionFundingAmountInLink,
//	)
//
//	vrfV2jobs := CreateVRFV2Jobs(
//		t,
//		chainlinkNodes,
//		vrfV2Contracts.Coordinator,
//		chainClient,
//		vrfv2_constants.MinimumConfirmations,
//	)
//
//	nativeTokenPrimaryKeyAddress, err := chainlinkNodes[0].PrimaryEthAddress()
//	require.NoError(t, err, "Error getting node's primary ETH key")
//
//	evmKeySpecificConfigTemplate := `
//[[EVM.KeySpecific]]
//Key = '%s'
//
//[EVM.KeySpecific.GasEstimator]
//PriceMax = '%d gwei'
//`
//	//todo - make evmKeySpecificConfigTemplate for multiple eth keys
//	evmKeySpecificConfig := fmt.Sprintf(evmKeySpecificConfigTemplate, nativeTokenPrimaryKeyAddress, vrfv2_constants.MaxGasPriceGWei)
//	tomlConfigWithUpdates := fmt.Sprintf("%s\n%s", config.BaseVRFV2NetworkDetailTomlConfig, evmKeySpecificConfig)
//
//	//todo - this does not show up??
//	newEnvLabel := "updatedWithRollout=true"
//	testEnvironmentAfterRedeployment := SetupVRFV2Environment(
//		t,
//		testNetwork,
//		tomlConfigWithUpdates,
//		existingTestEnvironment.Cfg.Namespace,
//		newEnvNamespacePrefix,
//		newEnvLabel,
//		newEnvTTL,
//	)
//
//	err = testEnvironmentAfterRedeployment.RolloutStatefulSets()
//	require.NoError(t, err, "Error performing rollout restart for test environment")
//
//	err = testEnvironmentAfterRedeployment.Run()
//	require.NoError(t, err, "Error running test environment")
//
//	//need to get node's urls again since port changed after redeployment
//	chainlinkNodesAfterRedeployment, err := client.ConnectChainlinkNodes(testEnvironmentAfterRedeployment)
//	require.NoError(t, err)
//
//	return vrfV2Contracts, chainlinkNodesAfterRedeployment, vrfV2jobs, testEnvironmentAfterRedeployment
//}

func SetupVRFV2Environment(
	t *testing.T,
	testNetwork blockchain.EVMNetwork,
	networkDetailTomlConfig string,
	existingNamespace string,
	namespacePrefix string,
	newEnvLabel string,
	ttl time.Duration,
) (testEnvironment *environment.Environment) {
	gethChartConfig := getGethChartConfig(testNetwork)

	if existingNamespace != "" {
		testEnvironment = environment.New(&environment.Config{
			Namespace: existingNamespace,
			Test:      t,
			TTL:       ttl,
			Labels:    []string{newEnvLabel},
		})
	} else {
		testEnvironment = environment.New(&environment.Config{
			NamespacePrefix: fmt.Sprintf("%s-%s", namespacePrefix, strings.ReplaceAll(strings.ToLower(testNetwork.Name), " ", "-")),
			Test:            t,
			TTL:             ttl,
		})
	}

	cd, err := chainlink.NewDeployment(1, map[string]any{
		"toml": client.AddNetworkDetailedConfig("", networkDetailTomlConfig, testNetwork),
		//need to restart the node with updated eth key config
		"db": map[string]interface{}{
			"stateful": "true",
		},
	})
	require.NoError(t, err, "Error creating chainlink deployment")
	testEnvironment = testEnvironment.
		AddHelm(gethChartConfig).
		AddHelmCharts(cd)
	err = testEnvironment.Run()
	require.NoError(t, err, "Error running test environment")
	return testEnvironment
}

func getGethChartConfig(testNetwork blockchain.EVMNetwork) environment.ConnectedChart {
	evmConfig := eth.New(nil)
	if !testNetwork.Simulated {
		evmConfig = eth.New(&eth.Props{
			NetworkName: testNetwork.Name,
			Simulated:   testNetwork.Simulated,
			WsURLs:      testNetwork.URLs,
		})
	}
	return evmConfig
}
