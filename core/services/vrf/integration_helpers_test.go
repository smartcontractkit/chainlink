package vrf_test

import (
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/assets"
	v2 "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/v2"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_consumer_v2"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

func testSingleConsumerHappyPath(
	t *testing.T,
	ownerKey ethkey.KeyV2,
	uni coordinatorV2Universe,
	consumer *bind.TransactOpts,
	consumerContract *vrf_consumer_v2.VRFConsumerV2,
	consumerContractAddress common.Address,
	coordinator vrf_coordinator_v2.VRFCoordinatorV2Interface,
	coordinatorAddress common.Address,
	batchCoordinatorAddress common.Address,
	assertions ...func(
		t *testing.T,
		coordinator vrf_coordinator_v2.VRFCoordinatorV2Interface,
		rwfe *vrf_coordinator_v2.VRFCoordinatorV2RandomWordsFulfilled),
) {
	key1 := cltest.MustGenerateRandomKey(t)
	key2 := cltest.MustGenerateRandomKey(t)
	gasLanePriceWei := assets.GWei(10)
	config, db := heavyweight.FullTestDBV2(t, "vrfv2_singleconsumer_happypath", func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, assets.GWei(10), v2.KeySpecific{
			// Gas lane.
			Key:          ptr(key1.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		}, v2.KeySpecific{
			// Gas lane.
			Key:          ptr(key2.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		})(c, s)
		c.EVM[0].MinIncomingConfirmations = ptr[uint32](2)
	})
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, uni.backend, ownerKey, key1, key2)

	// Create a subscription and fund with 5 LINK.
	subID := subscribeAndAssertSubscriptionCreatedEvent(t, consumerContract, consumer, consumerContractAddress, big.NewInt(5e18), coordinator, uni)

	// Fund gas lanes.
	sendEth(t, ownerKey, uni.backend, key1.Address, 10)
	sendEth(t, ownerKey, uni.backend, key2.Address, 10)
	require.NoError(t, app.Start(testutils.Context(t)))

	// Create VRF job using key1 and key2 on the same gas lane.
	jbs := createVRFJobs(
		t,
		[][]ethkey.KeyV2{{key1, key2}},
		app,
		coordinator,
		coordinatorAddress,
		batchCoordinatorAddress,
		uni,
		false,
		gasLanePriceWei)
	keyHash := jbs[0].VRFSpec.PublicKey.MustHash()

	// Make the first randomness request.
	numWords := uint32(20)
	requestID1, _ := requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, keyHash, subID, numWords, 500_000, coordinator, uni)

	// Wait for fulfillment to be queued.
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		t.Log("runs", len(runs))
		return len(runs) == 1
	}, testutils.WaitTimeout(t), time.Second).Should(gomega.BeTrue())

	// Mine the fulfillment that was queued.
	mine(t, requestID1, subID, uni, db)

	// Assert correct state of RandomWordsFulfilled event.
	// In particular:
	// * success should be true
	// * payment should be exactly the amount specified as the premium in the coordinator fee config
	rwfe := assertRandomWordsFulfilled(t, requestID1, true, coordinator)
	if len(assertions) > 0 {
		assertions[0](t, coordinator, rwfe)
	}

	// Make the second randomness request and assert fulfillment is successful
	requestID2, _ := requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, keyHash, subID, numWords, 500_000, coordinator, uni)
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		t.Log("runs", len(runs))
		return len(runs) == 2
	}, testutils.WaitTimeout(t), time.Second).Should(gomega.BeTrue())
	mine(t, requestID2, subID, uni, db)

	// Assert correct state of RandomWordsFulfilled event.
	// In particular:
	// * success should be true
	// * payment should be exactly the amount specified as the premium in the coordinator fee config
	rwfe = assertRandomWordsFulfilled(t, requestID2, true, coordinator)
	if len(assertions) > 0 {
		assertions[0](t, coordinator, rwfe)
	}

	// Assert correct number of random words sent by coordinator.
	assertNumRandomWords(t, consumerContract, numWords)

	// Assert that both send addresses were used to fulfill the requests
	n, err := uni.backend.PendingNonceAt(testutils.Context(t), key1.Address)
	require.NoError(t, err)
	require.EqualValues(t, 1, n)

	n, err = uni.backend.PendingNonceAt(testutils.Context(t), key2.Address)
	require.NoError(t, err)
	require.EqualValues(t, 1, n)

	t.Log("Done!")
}

func testMultipleConsumersNeedBHS(
	t *testing.T,
	ownerKey ethkey.KeyV2,
	uni coordinatorV2Universe,
	consumers []*bind.TransactOpts,
	consumerContracts []*vrf_consumer_v2.VRFConsumerV2,
	consumerContractAddresses []common.Address,
	coordinator vrf_coordinator_v2.VRFCoordinatorV2Interface,
	coordinatorAddress common.Address,
	batchCoordinatorAddress common.Address,
	assertions ...func(
		t *testing.T,
		coordinator vrf_coordinator_v2.VRFCoordinatorV2Interface,
		rwfe *vrf_coordinator_v2.VRFCoordinatorV2RandomWordsFulfilled),
) {
	nConsumers := len(consumers)
	vrfKey := cltest.MustGenerateRandomKey(t)
	sendEth(t, ownerKey, uni.backend, vrfKey.Address, 10)

	// generate n BHS keys to make sure BHS job rotates sending keys
	var bhsKeys []ethkey.KeyV2
	var bhsKeyAddresses []string
	var keySpecificOverrides []v2.KeySpecific
	var keys []interface{}
	gasLanePriceWei := assets.GWei(10)
	for i := 0; i < nConsumers; i++ {
		bhsKey := cltest.MustGenerateRandomKey(t)
		bhsKeys = append(bhsKeys, bhsKey)
		bhsKeyAddresses = append(bhsKeyAddresses, bhsKey.Address.String())
		keys = append(keys, bhsKey)
		keySpecificOverrides = append(keySpecificOverrides, v2.KeySpecific{
			Key:          ptr(bhsKey.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		})
		sendEth(t, ownerKey, uni.backend, bhsKey.Address, 10)
	}
	keySpecificOverrides = append(keySpecificOverrides, v2.KeySpecific{
		// Gas lane.
		Key:          ptr(vrfKey.EIP55Address),
		GasEstimator: v2.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
	})

	config, db := heavyweight.FullTestDBV2(t, "vrfv2_needs_blockhash_store", func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, assets.GWei(10), keySpecificOverrides...)(c, s)
		c.EVM[0].MinIncomingConfirmations = ptr[uint32](2)
		c.Feature.LogPoller = ptr(true)
		c.EVM[0].FinalityDepth = ptr[uint32](2)
		c.EVM[0].LogPollInterval = models.MustNewDuration(time.Second)
	})
	keys = append(keys, ownerKey, vrfKey)
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, uni.backend, keys...)
	require.NoError(t, app.Start(testutils.Context(t)))

	// Create VRF job.
	vrfJobs := createVRFJobs(
		t,
		[][]ethkey.KeyV2{{vrfKey}},
		app,
		coordinator,
		coordinatorAddress,
		batchCoordinatorAddress,
		uni,
		false,
		gasLanePriceWei)
	keyHash := vrfJobs[0].VRFSpec.PublicKey.MustHash()

	_ = createAndStartBHSJob(
		t, bhsKeyAddresses, app, uni.bhsContractAddress.String(), "",
		coordinatorAddress.String())

	// Ensure log poller is ready and has all logs.
	require.NoError(t, app.Chains.EVM.Chains()[0].LogPoller().Ready())
	require.NoError(t, app.Chains.EVM.Chains()[0].LogPoller().Replay(testutils.Context(t), 1))

	for i := 0; i < nConsumers; i++ {
		consumer := consumers[i]
		consumerContract := consumerContracts[i]

		// Create a subscription and fund with 0 LINK.
		_, subID := subscribeVRF(t, consumer, consumerContract, coordinator, uni.backend, new(big.Int))
		require.Equal(t, uint64(i+1), subID)

		// Make the randomness request. It will not yet succeed since it is underfunded.
		numWords := uint32(20)

		requestID, requestBlock := requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, keyHash, subID, numWords, 500_000, coordinator, uni)

		// Wait 101 blocks.
		for i := 0; i < 100; i++ {
			uni.backend.Commit()
		}
		verifyBlockhashStored(t, uni, requestBlock)

		// Wait another 160 blocks so that the request is outside of the 256 block window
		for i := 0; i < 160; i++ {
			uni.backend.Commit()
		}

		// Fund the subscription
		_, err := consumerContract.TopUpSubscription(consumer, big.NewInt(5e18 /* 5 LINK */))
		require.NoError(t, err)

		// Wait for fulfillment to be queued.
		gomega.NewGomegaWithT(t).Eventually(func() bool {
			uni.backend.Commit()
			runs, err := app.PipelineORM().GetAllRuns()
			require.NoError(t, err)
			t.Log("runs", len(runs))
			return len(runs) == 1
		}, testutils.WaitTimeout(t), time.Second).Should(gomega.BeTrue())

		mine(t, requestID, subID, uni, db)

		rwfe := assertRandomWordsFulfilled(t, requestID, true, coordinator)
		if len(assertions) > 0 {
			assertions[0](t, coordinator, rwfe)
		}

		// Assert correct number of random words sent by coordinator.
		assertNumRandomWords(t, consumerContract, numWords)
	}

	for i := 0; i < len(bhsKeys); i++ {
		n, err := uni.backend.PendingNonceAt(testutils.Context(t), bhsKeys[i].Address)
		require.NoError(t, err)
		require.EqualValues(t, 1, n)
	}
}

func verifyBlockhashStored(
	t *testing.T,
	uni coordinatorV2Universe,
	requestBlock uint64,
) {
	// Wait for the blockhash to be stored
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		uni.backend.Commit()
		callOpts := &bind.CallOpts{
			Pending:     false,
			From:        common.Address{},
			BlockNumber: nil,
			Context:     nil,
		}
		_, err := uni.bhsContract.GetBlockhash(callOpts, big.NewInt(int64(requestBlock)))
		if err == nil {
			return true
		} else if strings.Contains(err.Error(), "execution reverted") {
			return false
		} else {
			t.Fatal(err)
			return false
		}
	}, testutils.WaitTimeout(t), time.Second).Should(gomega.BeTrue())
}

func testSingleConsumerHappyPathBatchFulfillment(
	t *testing.T,
	ownerKey ethkey.KeyV2,
	uni coordinatorV2Universe,
	consumer *bind.TransactOpts,
	consumerContract *vrf_consumer_v2.VRFConsumerV2,
	consumerContractAddress common.Address,
	coordinator vrf_coordinator_v2.VRFCoordinatorV2Interface,
	coordinatorAddress common.Address,
	batchCoordinatorAddress common.Address,
	numRequests int,
	bigGasCallback bool,
	assertions ...func(
		t *testing.T,
		coordinator vrf_coordinator_v2.VRFCoordinatorV2Interface,
		rwfe *vrf_coordinator_v2.VRFCoordinatorV2RandomWordsFulfilled),
) {
	key1 := cltest.MustGenerateRandomKey(t)
	gasLanePriceWei := assets.GWei(10)
	config, db := heavyweight.FullTestDBV2(t, "vrfv2_singleconsumer_batch_happypath", func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, assets.GWei(10), v2.KeySpecific{
			// Gas lane.
			Key:          ptr(key1.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		})(c, s)
		c.EVM[0].GasEstimator.LimitDefault = ptr[uint32](5_000_000)
		c.EVM[0].MinIncomingConfirmations = ptr[uint32](2)
	})
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, uni.backend, ownerKey, key1)

	// Create a subscription and fund with 5 LINK.
	subID := subscribeAndAssertSubscriptionCreatedEvent(t, consumerContract, consumer, consumerContractAddress, big.NewInt(5e18), coordinator, uni)

	// Fund gas lane.
	sendEth(t, ownerKey, uni.backend, key1.Address, 10)
	require.NoError(t, app.Start(testutils.Context(t)))

	// Create VRF job using key1 and key2 on the same gas lane.
	jbs := createVRFJobs(
		t,
		[][]ethkey.KeyV2{{key1}},
		app,
		coordinator,
		coordinatorAddress,
		batchCoordinatorAddress,
		uni,
		true,
		gasLanePriceWei)
	keyHash := jbs[0].VRFSpec.PublicKey.MustHash()

	// Make some randomness requests.
	numWords := uint32(2)
	var reqIDs []*big.Int
	for i := 0; i < numRequests; i++ {
		requestID, _ := requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, keyHash, subID, numWords, 500_000, coordinator, uni)
		reqIDs = append(reqIDs, requestID)
	}

	if bigGasCallback {
		// Make one randomness request with the max callback gas limit.
		// It should live in a batch on it's own.
		requestID, _ := requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, keyHash, subID, numWords, 2_500_000, coordinator, uni)
		reqIDs = append(reqIDs, requestID)
	}

	// Wait for fulfillment to be queued.
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		t.Log("runs", len(runs))
		if bigGasCallback {
			return len(runs) == (numRequests + 1)
		}
		return len(runs) == numRequests
	}, testutils.WaitTimeout(t), time.Second).Should(gomega.BeTrue())

	mineBatch(t, reqIDs, subID, uni, db)

	for i, requestID := range reqIDs {
		// Assert correct state of RandomWordsFulfilled event.
		// The last request will be the successful one because of the way the example
		// contract is written.
		var rwfe *vrf_coordinator_v2.VRFCoordinatorV2RandomWordsFulfilled
		if i == (len(reqIDs) - 1) {
			rwfe = assertRandomWordsFulfilled(t, requestID, true, coordinator)
		} else {
			rwfe = assertRandomWordsFulfilled(t, requestID, false, coordinator)
		}
		if len(assertions) > 0 {
			assertions[0](t, coordinator, rwfe)
		}
	}

	// Assert correct number of random words sent by coordinator.
	assertNumRandomWords(t, consumerContract, numWords)
}

func testSingleConsumerNeedsTopUp(
	t *testing.T,
	ownerKey ethkey.KeyV2,
	uni coordinatorV2Universe,
	consumer *bind.TransactOpts,
	consumerContract *vrf_consumer_v2.VRFConsumerV2,
	consumerContractAddress common.Address,
	coordinator vrf_coordinator_v2.VRFCoordinatorV2Interface,
	coordinatorAddress common.Address,
	batchCoordinatorAddress common.Address,
	initialFundingAmount *big.Int,
	topUpAmount *big.Int,
	assertions ...func(
		t *testing.T,
		coordinator vrf_coordinator_v2.VRFCoordinatorV2Interface,
		rwfe *vrf_coordinator_v2.VRFCoordinatorV2RandomWordsFulfilled),
) {
	key := cltest.MustGenerateRandomKey(t)
	gasLanePriceWei := assets.GWei(1000)
	config, db := heavyweight.FullTestDBV2(t, "vrfv2_singleconsumer_needstopup", func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, assets.GWei(1000), v2.KeySpecific{
			// Gas lane.
			Key:          ptr(key.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		})(c, s)
		c.EVM[0].MinIncomingConfirmations = ptr[uint32](2)
	})
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, uni.backend, ownerKey, key)

	// Create and fund a subscription
	subID := subscribeAndAssertSubscriptionCreatedEvent(t, consumerContract, consumer, consumerContractAddress, initialFundingAmount, coordinator, uni)

	// Fund expensive gas lane.
	sendEth(t, ownerKey, uni.backend, key.Address, 10)
	require.NoError(t, app.Start(testutils.Context(t)))

	// Create VRF job.
	jbs := createVRFJobs(
		t,
		[][]ethkey.KeyV2{{key}},
		app,
		coordinator,
		coordinatorAddress,
		batchCoordinatorAddress,
		uni,
		false,
		gasLanePriceWei)
	keyHash := jbs[0].VRFSpec.PublicKey.MustHash()

	numWords := uint32(20)
	requestID, _ := requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, keyHash, subID, numWords, 500_000, coordinator, uni)

	// Fulfillment will not be enqueued because subscriber doesn't have enough LINK.
	gomega.NewGomegaWithT(t).Consistently(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		t.Log("assert 1", "runs", len(runs))
		return len(runs) == 0
	}, 5*time.Second, 1*time.Second).Should(gomega.BeTrue())

	// Top up subscription with enough LINK to see the job through.
	_, err := consumerContract.TopUpSubscription(consumer, topUpAmount)
	require.NoError(t, err)

	// Wait for fulfillment to go through.
	gomega.NewWithT(t).Eventually(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		t.Log("assert 2", "runs", len(runs))
		return len(runs) == 1
	}, testutils.WaitTimeout(t), 1*time.Second).Should(gomega.BeTrue())

	// Mine the fulfillment. Need to wait for Txm to mark the tx as confirmed
	// so that we can actually see the event on the simulated chain.
	mine(t, requestID, subID, uni, db)

	// Assert the state of the RandomWordsFulfilled event.
	rwfe := assertRandomWordsFulfilled(t, requestID, true, coordinator)
	if len(assertions) > 0 {
		assertions[0](t, coordinator, rwfe)
	}

	// Assert correct number of random words sent by coordinator.
	assertNumRandomWords(t, consumerContract, numWords)
}

// testBlockHeaderFeeder starts VRF and block header feeder jobs
// subscription is unfunded initially and funded after 256 blocks
// the function makes sure the block header feeder stored blockhash for
// a block older than 256 blocks
func testBlockHeaderFeeder(
	t *testing.T,
	ownerKey ethkey.KeyV2,
	uni coordinatorV2Universe,
	consumers []*bind.TransactOpts,
	consumerContracts []*vrf_consumer_v2.VRFConsumerV2,
	consumerContractAddresses []common.Address,
	coordinator vrf_coordinator_v2.VRFCoordinatorV2Interface,
	coordinatorAddress common.Address,
	batchCoordinatorAddress common.Address,
	assertions ...func(
		t *testing.T,
		coordinator vrf_coordinator_v2.VRFCoordinatorV2Interface,
		rwfe *vrf_coordinator_v2.VRFCoordinatorV2RandomWordsFulfilled),
) {
	nConsumers := len(consumers)

	vrfKey := cltest.MustGenerateRandomKey(t)
	bhfKey := cltest.MustGenerateRandomKey(t)
	bhfKeys := []string{bhfKey.Address.String()}

	sendEth(t, ownerKey, uni.backend, bhfKey.Address, 10)
	sendEth(t, ownerKey, uni.backend, vrfKey.Address, 10)

	gasLanePriceWei := assets.GWei(10)

	config, db := heavyweight.FullTestDBV2(t, "vrfv2_test_block_header_feeder", func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, gasLanePriceWei, v2.KeySpecific{
			// Gas lane.
			Key:          ptr(vrfKey.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		})(c, s)
		c.EVM[0].MinIncomingConfirmations = ptr[uint32](2)
		c.Feature.LogPoller = ptr(true)
		c.EVM[0].FinalityDepth = ptr[uint32](2)
		c.EVM[0].LogPollInterval = models.MustNewDuration(time.Second)
	})
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, uni.backend, ownerKey, vrfKey, bhfKey)
	require.NoError(t, app.Start(testutils.Context(t)))

	// Create VRF job.
	vrfJobs := createVRFJobs(
		t,
		[][]ethkey.KeyV2{{vrfKey}},
		app,
		coordinator,
		coordinatorAddress,
		batchCoordinatorAddress,
		uni,
		false,
		gasLanePriceWei)
	keyHash := vrfJobs[0].VRFSpec.PublicKey.MustHash()

	_ = createAndStartBlockHeaderFeederJob(
		t, bhfKeys, app, uni.bhsContractAddress.String(), uni.batchBHSContractAddress.String(), "",
		coordinatorAddress.String())

	// Ensure log poller is ready and has all logs.
	require.NoError(t, app.Chains.EVM.Chains()[0].LogPoller().Ready())
	require.NoError(t, app.Chains.EVM.Chains()[0].LogPoller().Replay(testutils.Context(t), 1))

	for i := 0; i < nConsumers; i++ {
		consumer := consumers[i]
		consumerContract := consumerContracts[i]

		// Create a subscription and fund with 0 LINK.
		_, subID := subscribeVRF(t, consumer, consumerContract, coordinator, uni.backend, new(big.Int))
		require.Equal(t, uint64(i+1), subID)

		// Make the randomness request. It will not yet succeed since it is underfunded.
		numWords := uint32(20)

		requestID, requestBlock := requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, keyHash, subID, numWords, 500_000, coordinator, uni)

		// Wait 256 blocks.
		for i := 0; i < 256; i++ {
			uni.backend.Commit()
		}
		verifyBlockhashStored(t, uni, requestBlock)

		// Fund the subscription
		_, err := consumerContract.TopUpSubscription(consumer, big.NewInt(5e18 /* 5 LINK */))
		require.NoError(t, err)

		// Wait for fulfillment to be queued.
		gomega.NewGomegaWithT(t).Eventually(func() bool {
			uni.backend.Commit()
			runs, err := app.PipelineORM().GetAllRuns()
			require.NoError(t, err)
			t.Log("runs", len(runs))
			return len(runs) == 1
		}, testutils.WaitTimeout(t), time.Second).Should(gomega.BeTrue())

		mine(t, requestID, subID, uni, db)

		rwfe := assertRandomWordsFulfilled(t, requestID, true, coordinator)
		if len(assertions) > 0 {
			assertions[0](t, coordinator, rwfe)
		}

		// Assert correct number of random words sent by coordinator.
		assertNumRandomWords(t, consumerContract, numWords)
	}
}
