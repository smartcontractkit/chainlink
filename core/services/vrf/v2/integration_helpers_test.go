//go:build integration

package v2_test

import (
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/onsi/gomega"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/ethkey"
	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrf_consumer_v2_plus_upgradeable_example"
	"github.com/smartcontractkit/chainlink-evm/pkg/assets"
	"github.com/smartcontractkit/chainlink-evm/pkg/chains/legacyevm"
	v2 "github.com/smartcontractkit/chainlink-evm/pkg/config/toml"
	"github.com/smartcontractkit/chainlink-evm/pkg/txmgr"
	"github.com/smartcontractkit/chainlink-evm/pkg/types"
	evmutils "github.com/smartcontractkit/chainlink-evm/pkg/utils"
	txmgrcommon "github.com/smartcontractkit/chainlink-framework/chains/txmgr"
	txmgrtypes "github.com/smartcontractkit/chainlink-framework/chains/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	v22 "github.com/smartcontractkit/chainlink/v2/core/services/vrf/v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/vrfcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/vrftesthelpers"
	"github.com/smartcontractkit/chainlink/v2/core/testdata/testspecs"
	"github.com/smartcontractkit/chainlink/v2/core/utils/testutils/heavyweight"
)

func testSingleConsumerHappyPath(
	t *testing.T,
	ownerKey ethkey.KeyV2,
	uni coordinatorV2UniverseCommon,
	consumer *bind.TransactOpts,
	consumerContract vrftesthelpers.VRFConsumerContract,
	consumerContractAddress common.Address,
	coordinator v22.CoordinatorV2_X,
	coordinatorAddress common.Address,
	batchCoordinatorAddress common.Address,
	vrfOwnerAddress *common.Address,
	vrfVersion vrfcommon.Version,
	nativePayment bool,
	assertions ...func(
		t *testing.T,
		coordinator v22.CoordinatorV2_X,
		rwfe v22.RandomWordsFulfilled,
		subID *big.Int),
) {
	ctx := t.Context()
	key1 := cltest.MustGenerateRandomKey(t)
	key2 := cltest.MustGenerateRandomKey(t)
	gasLanePriceWei := assets.GWei(10)
	config, db := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, assets.GWei(10), v2.KeySpecific{
			// Gas lane.
			Key:          new(key1.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		}, v2.KeySpecific{
			// Gas lane.
			Key:          new(key2.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		})(c, s)
		c.EVM[0].MinIncomingConfirmations = new(uint32(2))
		c.Feature.LogPoller = new(true)
		c.EVM[0].LogPollInterval = commonconfig.MustNewDuration(100 * time.Millisecond)
	})
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, uni.backend, ownerKey, key1, key2)

	// Create a subscription and fund with 5 LINK.
	subID := subscribeAndAssertSubscriptionCreatedEvent(t, consumerContract, consumer, consumerContractAddress, big.NewInt(5e18), coordinator, uni.backend, nativePayment)

	// Fund gas lanes.
	sendEth(t, ownerKey, uni.backend, key1.Address, 10)
	sendEth(t, ownerKey, uni.backend, key2.Address, 10)
	require.NoError(t, app.Start(ctx))

	// Create VRF job using key1 and key2 on the same gas lane.
	jbs := createVRFJobs(
		t,
		[][]ethkey.KeyV2{{key1, key2}},
		app,
		coordinator,
		coordinatorAddress,
		batchCoordinatorAddress,
		uni,
		vrfOwnerAddress,
		vrfVersion,
		false,
		gasLanePriceWei,
	)
	keyHash := jbs[0].VRFSpec.PublicKey.MustHash()

	// Make the first randomness request.
	numWords := uint32(20)
	requestID1, _ := requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, keyHash, subID, numWords, 500_000, coordinator, uni.backend, nativePayment)

	// Wait for fulfillment to be queued.
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns(ctx)
		require.NoError(t, err)
		t.Log("runs", len(runs))
		return len(runs) == 1
	}, testutils.WaitTimeout(t), 100*time.Millisecond).Should(gomega.BeTrue())

	// Mine the fulfillment that was queued.
	mine(t, requestID1, subID, uni.backend, db, vrfVersion, testutils.SimulatedChainID)

	// Assert correct state of RandomWordsFulfilled event.
	// In particular:
	// * success should be true
	// * payment should be exactly the amount specified as the premium in the coordinator fee config
	rwfe := assertRandomWordsFulfilled(t, requestID1, true, coordinator, nativePayment)
	if len(assertions) > 0 {
		assertions[0](t, coordinator, rwfe, subID)
	}

	// Make the second randomness request and assert fulfillment is successful
	requestID2, _ := requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, keyHash, subID, numWords, 500_000, coordinator, uni.backend, nativePayment)
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns(ctx)
		require.NoError(t, err)
		t.Log("runs", len(runs))
		return len(runs) == 2
	}, testutils.WaitTimeout(t), 100*time.Millisecond).Should(gomega.BeTrue())
	mine(t, requestID2, subID, uni.backend, db, vrfVersion, testutils.SimulatedChainID)

	// Assert correct state of RandomWordsFulfilled event.
	// In particular:
	// * success should be true
	// * payment should be exactly the amount specified as the premium in the coordinator fee config
	rwfe = assertRandomWordsFulfilled(t, requestID2, true, coordinator, nativePayment)
	if len(assertions) > 0 {
		assertions[0](t, coordinator, rwfe, subID)
	}

	// Assert correct number of random words sent by coordinator.
	assertNumRandomWords(t, consumerContract, numWords)

	// Assert that both send addresses were used to fulfill the requests
	n, err := uni.backend.Client().PendingNonceAt(ctx, key1.Address)
	require.NoError(t, err)
	require.EqualValues(t, 1, n)

	n, err = uni.backend.Client().PendingNonceAt(ctx, key2.Address)
	require.NoError(t, err)
	require.EqualValues(t, 1, n)

	t.Log("Done!")
}

func testMultipleConsumersNeedBHS(
	t *testing.T,
	ownerKey ethkey.KeyV2,
	uni coordinatorV2UniverseCommon,
	consumers []*bind.TransactOpts,
	consumerContracts []vrftesthelpers.VRFConsumerContract,
	_ []common.Address,
	coordinator v22.CoordinatorV2_X,
	coordinatorAddress common.Address,
	batchCoordinatorAddress common.Address,
	vrfOwnerAddress *common.Address,
	vrfVersion vrfcommon.Version,
	nativePayment bool,
	assertions ...func(
		t *testing.T,
		coordinator v22.CoordinatorV2_X,
		rwfe v22.RandomWordsFulfilled,
	),
) {
	ctx := t.Context()
	nConsumers := len(consumers)
	vrfKey := cltest.MustGenerateRandomKey(t)
	sendEth(t, ownerKey, uni.backend, vrfKey.Address, 10)

	// generate n BHS keys to make sure BHS job rotates sending keys
	bhsKeyAddresses := make([]string, 0, nConsumers)
	keySpecificOverrides := make([]v2.KeySpecific, 0, nConsumers+1)
	keys := make([]any, 0, nConsumers+2)
	gasLanePriceWei := assets.GWei(10)
	for range nConsumers {
		bhsKey := cltest.MustGenerateRandomKey(t)
		bhsKeyAddresses = append(bhsKeyAddresses, bhsKey.Address.String())
		keys = append(keys, bhsKey)
		keySpecificOverrides = append(keySpecificOverrides, v2.KeySpecific{
			Key:          new(bhsKey.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		})
		sendEth(t, ownerKey, uni.backend, bhsKey.Address, 10)
	}
	keySpecificOverrides = append(keySpecificOverrides, v2.KeySpecific{
		// Gas lane.
		Key:          new(vrfKey.EIP55Address),
		GasEstimator: v2.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
	})

	config, db := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, assets.GWei(10), keySpecificOverrides...)(c, s)
		c.EVM[0].MinIncomingConfirmations = new(uint32(2))
		c.Feature.LogPoller = new(true)
		c.EVM[0].LogPollInterval = commonconfig.MustNewDuration(100 * time.Millisecond)
		c.EVM[0].FinalityDepth = new(uint32(2))
	})
	keys = append(keys, ownerKey, vrfKey)
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, uni.backend, keys...)
	require.NoError(t, app.Start(ctx))

	// Create VRF job.
	vrfJobs := createVRFJobs(
		t,
		[][]ethkey.KeyV2{{vrfKey}},
		app,
		coordinator,
		coordinatorAddress,
		batchCoordinatorAddress,
		uni,
		vrfOwnerAddress,
		vrfVersion,
		false,
		gasLanePriceWei,
	)
	keyHash := vrfJobs[0].VRFSpec.PublicKey.MustHash()

	v2PlusCoordinatorAddress := coordinatorAddress.String()

	_ = vrftesthelpers.CreateAndStartBHSJob(
		t, bhsKeyAddresses, app, uni.bhsContractAddress.String(),
		"", v2PlusCoordinatorAddress, "", 0, 200, 0, 100,
	)

	chain, ok := app.GetRelayers().LegacyEVMChains().Slice()[0].(legacyevm.Chain) //nolint:staticcheck // TODO: migrate to relayer interface
	require.True(t, ok)
	// Ensure log poller is ready and has all logs.
	require.NoError(t, chain.LogPoller().Ready())
	require.NoError(t, chain.LogPoller().Replay(ctx, 1))

	for i := range nConsumers {
		consumer := consumers[i]
		consumerContract := consumerContracts[i]

		// Create a subscription and fund with 0 LINK.
		_, subID := subscribeVRF(t, consumer, consumerContract, coordinator, uni.backend, new(big.Int), nativePayment)

		// Make the randomness request. It will not yet succeed since it is underfunded.
		numWords := uint32(20)

		requestID, requestBlock := requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, keyHash, subID, numWords, 500_000, coordinator, uni.backend, nativePayment)

		// Wait 101 blocks.
		for range 100 {
			uni.backend.Commit()
		}
		verifyBlockhashStored(t, uni, requestBlock)

		// Wait another 160 blocks so that the request is outside of the 256 block window
		for range 160 {
			uni.backend.Commit()
		}

		// Fund the subscription
		topUpSubscription(t, consumer, consumerContract, uni.backend, big.NewInt(5e18 /* 5 LINK */), nativePayment)

		// Wait for fulfillment to be queued.
		require.EventuallyWithT(t, func(c *assert.CollectT) {
			uni.backend.Commit()
			runs, err := app.PipelineORM().GetAllRuns(ctx)
			require.NoError(c, err)
			require.Len(c, runs, 1)
		}, testutils.WaitTimeout(t), 100*time.Millisecond)

		mine(t, requestID, subID, uni.backend, db, vrfVersion, testutils.SimulatedChainID)

		rwfe := assertRandomWordsFulfilled(t, requestID, true, coordinator, nativePayment)
		if len(assertions) > 0 {
			assertions[0](t, coordinator, rwfe)
		}

		// Assert correct number of random words sent by coordinator.
		assertNumRandomWords(t, consumerContract, numWords)
	}
}

func testMultipleConsumersNeedTrustedBHS(
	t *testing.T,
	ownerKey ethkey.KeyV2,
	uni coordinatorV2PlusUniverse,
	consumers []*bind.TransactOpts,
	consumerContracts []vrftesthelpers.VRFConsumerContract,
	_ []common.Address,
	coordinator v22.CoordinatorV2_X,
	coordinatorAddress common.Address,
	batchCoordinatorAddress common.Address,
	vrfVersion vrfcommon.Version,
	nativePayment bool,
	addedDelay bool,
	assertions ...func(
		t *testing.T,
		coordinator v22.CoordinatorV2_X,
		rwfe v22.RandomWordsFulfilled,
	),
) {
	ctx := t.Context()
	nConsumers := len(consumers)
	vrfKey := cltest.MustGenerateRandomKey(t)
	sendEth(t, ownerKey, uni.backend, vrfKey.Address, 10)

	// generate n BHS keys to make sure BHS job rotates sending keys
	bhsKeyAddresses := make([]common.Address, 0, nConsumers)
	bhsKeyAddressesStrings := make([]string, 0, nConsumers)
	keySpecificOverrides := make([]v2.KeySpecific, 0, nConsumers+1)
	keys := make([]any, 0, nConsumers+2)
	gasLanePriceWei := assets.GWei(10)
	for range nConsumers {
		bhsKey := cltest.MustGenerateRandomKey(t)
		bhsKeyAddressesStrings = append(bhsKeyAddressesStrings, bhsKey.Address.String())
		bhsKeyAddresses = append(bhsKeyAddresses, bhsKey.Address)
		keys = append(keys, bhsKey)
		keySpecificOverrides = append(keySpecificOverrides, v2.KeySpecific{
			Key:          new(bhsKey.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		})
		sendEth(t, ownerKey, uni.backend, bhsKey.Address, 10)
	}
	keySpecificOverrides = append(keySpecificOverrides, v2.KeySpecific{
		// Gas lane.
		Key:          new(vrfKey.EIP55Address),
		GasEstimator: v2.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
	})

	// Whitelist vrf key for trusted BHS.
	{
		_, err := uni.trustedBhsContract.SetWhitelist(uni.neil, bhsKeyAddresses)
		require.NoError(t, err)
		uni.backend.Commit()
	}

	config, db := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, assets.GWei(10), keySpecificOverrides...)(c, s)
		c.EVM[0].MinIncomingConfirmations = new(uint32(2))
		c.EVM[0].GasEstimator.LimitDefault = new(uint64(5_000_000))
		c.Feature.LogPoller = new(true)
		c.EVM[0].LogPollInterval = commonconfig.MustNewDuration(100 * time.Millisecond)
		c.EVM[0].FinalityDepth = new(uint32(2))
	})
	keys = append(keys, ownerKey, vrfKey)
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, uni.backend, keys...)
	require.NoError(t, app.Start(ctx))

	// Create VRF job.
	vrfJobs := createVRFJobs(
		t,
		[][]ethkey.KeyV2{{vrfKey}},
		app,
		coordinator,
		coordinatorAddress,
		batchCoordinatorAddress,
		uni.coordinatorV2UniverseCommon,
		nil,
		vrfVersion,
		false,
		gasLanePriceWei,
	)
	keyHash := vrfJobs[0].VRFSpec.PublicKey.MustHash()

	v2PlusCoordinatorAddress := coordinatorAddress.String()

	waitBlocks := 100
	if addedDelay {
		waitBlocks = 400
	}
	_ = vrftesthelpers.CreateAndStartBHSJob(
		t, bhsKeyAddressesStrings, app, "", "", v2PlusCoordinatorAddress, uni.trustedBhsContractAddress.String(), 20, 1000, 0, waitBlocks,
	)

	// Ensure log poller is ready and has all logs.
	chain, ok := app.GetRelayers().LegacyEVMChains().Slice()[0].(legacyevm.Chain) //nolint:staticcheck // TODO: migrate to relayer interface
	require.True(t, ok)
	require.NoError(t, chain.LogPoller().Ready())
	require.NoError(t, chain.LogPoller().Replay(ctx, 1))

	for i := range nConsumers {
		consumer := consumers[i]
		consumerContract := consumerContracts[i]

		// Create a subscription and fund with 0 LINK.
		_, subID := subscribeVRF(t, consumer, consumerContract, coordinator, uni.backend, new(big.Int), nativePayment)

		// Make the randomness request. It will not yet succeed since it is underfunded.
		numWords := uint32(20)

		requestID, requestBlock := requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, keyHash, subID, numWords, 500_000, coordinator, uni.backend, nativePayment)

		// Wait 101 blocks.
		for range 100 {
			uni.backend.Commit()
		}

		// For an added delay, we even go beyond the EVM lookback limit. This is not a problem in a trusted BHS setup.
		if addedDelay {
			for range 300 {
				uni.backend.Commit()
			}
		}

		verifyBlockhashStoredTrusted(t, uni, requestBlock)

		// Wait another 160 blocks so that the request is outside of the 256 block window
		for range 160 {
			uni.backend.Commit()
		}

		// Fund the subscription
		topUpSubscription(t, consumer, consumerContract, uni.backend, big.NewInt(5e18 /* 5 LINK */), nativePayment)

		// Wait for fulfillment to be queued.
		require.Eventually(t, func() bool {
			uni.backend.Commit()
			runs, err := app.PipelineORM().GetAllRuns(ctx)
			require.NoError(t, err)
			t.Log("runs", len(runs))
			return len(runs) >= 1
		}, testutils.WaitTimeout(t), 100*time.Millisecond)

		mine(t, requestID, subID, uni.backend, db, vrfVersion, testutils.SimulatedChainID)

		rwfe := assertRandomWordsFulfilled(t, requestID, true, coordinator, nativePayment)
		if len(assertions) > 0 {
			assertions[0](t, coordinator, rwfe)
		}

		// Assert correct number of random words sent by coordinator.
		assertNumRandomWords(t, consumerContract, numWords)
	}
}

func verifyBlockhashStored(
	t *testing.T,
	uni coordinatorV2UniverseCommon,
	requestBlock uint64,
) {
	// Wait for the blockhash to be stored
	require.Eventually(t, func() bool {
		uni.backend.Commit()
		callOpts := &bind.CallOpts{
			Pending:     false,
			From:        common.Address{},
			BlockNumber: nil,
			Context:     nil,
		}
		_, err := uni.bhsContract.GetBlockhash(callOpts, new(big.Int).SetUint64(requestBlock))
		if err == nil {
			return true
		}
		if strings.Contains(err.Error(), "execution reverted") {
			return false
		}
		require.FailNowf(t, "GetBlockhash: %v", err.Error())
		return false
	}, testutils.WaitTimeoutCustom(t, 5*time.Minute), 100*time.Millisecond)
}

func verifyBlockhashStoredTrusted(
	t *testing.T,
	uni coordinatorV2PlusUniverse,
	requestBlock uint64,
) {
	// Wait for the blockhash to be stored
	require.Eventually(t, func() bool {
		uni.backend.Commit()
		callOpts := &bind.CallOpts{
			Pending:     false,
			From:        common.Address{},
			BlockNumber: nil,
			Context:     nil,
		}
		_, err := uni.trustedBhsContract.GetBlockhash(callOpts, new(big.Int).SetUint64(requestBlock))
		if err == nil {
			return true
		}
		if strings.Contains(err.Error(), "execution reverted") {
			return false
		}
		require.FailNowf(t, "GetBlockhash (trusted BHS): %v", err.Error())
		return false
	}, testutils.WaitTimeoutCustom(t, 5*time.Minute), 100*time.Millisecond)
}

func testSingleConsumerHappyPathBatchFulfillment(
	t *testing.T,
	ownerKey ethkey.KeyV2,
	uni coordinatorV2UniverseCommon,
	consumer *bind.TransactOpts,
	consumerContract vrftesthelpers.VRFConsumerContract,
	consumerContractAddress common.Address,
	coordinator v22.CoordinatorV2_X,
	coordinatorAddress common.Address,
	batchCoordinatorAddress common.Address,
	vrfOwnerAddress *common.Address,
	numRequests int,
	bigGasCallback bool,
	vrfVersion vrfcommon.Version,
	nativePayment bool,
	assertions ...func(
		t *testing.T,
		coordinator v22.CoordinatorV2_X,
		rwfe v22.RandomWordsFulfilled,
		subID *big.Int),
) {
	ctx := t.Context()
	key1 := cltest.MustGenerateRandomKey(t)
	gasLanePriceWei := assets.GWei(10)
	config, db := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, assets.GWei(10), v2.KeySpecific{
			// Gas lane.
			Key:          new(key1.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		})(c, s)
		c.EVM[0].GasEstimator.LimitDefault = new(uint64(5_000_000))
		c.EVM[0].MinIncomingConfirmations = new(uint32(2))
		c.EVM[0].ChainID = (*sqlutil.Big)(testutils.SimulatedChainID)
		c.Feature.LogPoller = new(true)
		c.EVM[0].LogPollInterval = commonconfig.MustNewDuration(100 * time.Millisecond)
	})
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, uni.backend, ownerKey, key1)

	// Create a subscription and fund with 5 LINK.
	subID := subscribeAndAssertSubscriptionCreatedEvent(t, consumerContract, consumer, consumerContractAddress, big.NewInt(5e18), coordinator, uni.backend, nativePayment)

	// Fund gas lane.
	sendEth(t, ownerKey, uni.backend, key1.Address, 10)
	require.NoError(t, app.Start(ctx))

	// Create VRF job using key1 and key2 on the same gas lane.
	jbs := createVRFJobs(
		t,
		[][]ethkey.KeyV2{{key1}},
		app,
		coordinator,
		coordinatorAddress,
		batchCoordinatorAddress,
		uni,
		vrfOwnerAddress,
		vrfVersion,
		true,
		gasLanePriceWei,
	)
	keyHash := jbs[0].VRFSpec.PublicKey.MustHash()

	// Make some randomness requests.
	numWords := uint32(2)
	var reqIDs []*big.Int
	for range numRequests {
		requestID, _ := requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, keyHash, subID, numWords, 500_000, coordinator, uni.backend, nativePayment)
		reqIDs = append(reqIDs, requestID)
	}

	if bigGasCallback {
		// Make one randomness request with the max callback gas limit.
		// It should live in a batch on it's own.
		requestID, _ := requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, keyHash, subID, numWords, 2_500_000, coordinator, uni.backend, nativePayment)
		reqIDs = append(reqIDs, requestID)
	}

	// Wait for fulfillment to be queued.
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns(ctx)
		require.NoError(t, err)
		t.Log("runs", len(runs))
		if bigGasCallback {
			return len(runs) == (numRequests + 1)
		}
		return len(runs) == numRequests
	}, testutils.WaitTimeout(t), 100*time.Millisecond).Should(gomega.BeTrue())

	mineBatch(t, reqIDs, subID, uni.backend, db, vrfVersion, testutils.SimulatedChainID)

	for i, requestID := range reqIDs {
		// Assert correct state of RandomWordsFulfilled event.
		// The last request will be the successful one because of the way the example
		// contract is written.
		var rwfe v22.RandomWordsFulfilled
		if i == (len(reqIDs) - 1) {
			rwfe = assertRandomWordsFulfilled(t, requestID, true, coordinator, nativePayment)
		} else {
			rwfe = assertRandomWordsFulfilled(t, requestID, false, coordinator, nativePayment)
		}
		if len(assertions) > 0 {
			assertions[0](t, coordinator, rwfe, subID)
		}
	}

	// Assert correct number of random words sent by coordinator.
	assertNumRandomWords(t, consumerContract, numWords)
}

func testSingleConsumerNeedsTopUp(
	t *testing.T,
	ownerKey ethkey.KeyV2,
	uni coordinatorV2UniverseCommon,
	consumer *bind.TransactOpts,
	consumerContract vrftesthelpers.VRFConsumerContract,
	consumerContractAddress common.Address,
	coordinator v22.CoordinatorV2_X,
	coordinatorAddress common.Address,
	batchCoordinatorAddress common.Address,
	vrfOwnerAddress *common.Address,
	initialFundingAmount *big.Int,
	topUpAmount *big.Int,
	vrfVersion vrfcommon.Version,
	nativePayment bool,
	assertions ...func(
		t *testing.T,
		coordinator v22.CoordinatorV2_X,
		rwfe v22.RandomWordsFulfilled),
) {
	ctx := t.Context()
	key := cltest.MustGenerateRandomKey(t)
	gasLanePriceWei := assets.GWei(1000)
	config, db := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, assets.GWei(1000), v2.KeySpecific{
			// Gas lane.
			Key:          new(key.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		})(c, s)
		c.EVM[0].MinIncomingConfirmations = new(uint32(2))
		c.Feature.LogPoller = new(true)
		c.EVM[0].LogPollInterval = commonconfig.MustNewDuration(100 * time.Millisecond)
	})
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, uni.backend, ownerKey, key)

	// Create and fund a subscription
	subID := subscribeAndAssertSubscriptionCreatedEvent(t, consumerContract, consumer, consumerContractAddress, initialFundingAmount, coordinator, uni.backend, nativePayment)

	// Fund expensive gas lane.
	sendEth(t, ownerKey, uni.backend, key.Address, 10)
	require.NoError(t, app.Start(ctx))

	// Create VRF job.
	jbs := createVRFJobs(
		t,
		[][]ethkey.KeyV2{{key}},
		app,
		coordinator,
		coordinatorAddress,
		batchCoordinatorAddress,
		uni,
		vrfOwnerAddress,
		vrfVersion,
		false,
		gasLanePriceWei,
	)
	keyHash := jbs[0].VRFSpec.PublicKey.MustHash()

	numWords := uint32(20)
	requestID, _ := requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, keyHash, subID, numWords, 500_000, coordinator, uni.backend, nativePayment)

	// Fulfillment will not be enqueued because subscriber doesn't have enough LINK.
	gomega.NewGomegaWithT(t).Consistently(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns(ctx)
		require.NoError(t, err)
		t.Log("assert 1", "runs", len(runs))
		return len(runs) == 0
	}, 5*time.Second, 1*time.Second).Should(gomega.BeTrue())

	// Top up subscription with enough LINK to see the job through.
	topUpSubscription(t, consumer, consumerContract, uni.backend, topUpAmount, nativePayment)
	uni.backend.Commit()

	// Wait for fulfillment to go through.
	require.Eventually(t, func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns(ctx)
		require.NoError(t, err)
		t.Log("assert 2", "runs", len(runs))
		return len(runs) == 1
	}, testutils.WaitTimeout(t), 100*time.Millisecond)

	// Mine the fulfillment. Need to wait for Txm to mark the tx as confirmed
	// so that we can actually see the event on the simulated chain.
	mine(t, requestID, subID, uni.backend, db, vrfVersion, testutils.SimulatedChainID)

	// Assert the state of the RandomWordsFulfilled event.
	rwfe := assertRandomWordsFulfilled(t, requestID, true, coordinator, nativePayment)
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
	uni coordinatorV2UniverseCommon,
	consumers []*bind.TransactOpts,
	consumerContracts []vrftesthelpers.VRFConsumerContract,
	_ []common.Address,
	coordinator v22.CoordinatorV2_X,
	coordinatorAddress common.Address,
	batchCoordinatorAddress common.Address,
	vrfOwnerAddress *common.Address,
	vrfVersion vrfcommon.Version,
	nativePayment bool,
	assertions ...func(
		t *testing.T,
		coordinator v22.CoordinatorV2_X,
		rwfe v22.RandomWordsFulfilled),
) {
	ctx := t.Context()
	nConsumers := len(consumers)

	vrfKey := cltest.MustGenerateRandomKey(t)
	bhfKey := cltest.MustGenerateRandomKey(t)
	bhfKeys := []string{bhfKey.Address.String()}

	sendEth(t, ownerKey, uni.backend, bhfKey.Address, 10)
	sendEth(t, ownerKey, uni.backend, vrfKey.Address, 10)

	gasLanePriceWei := assets.GWei(10)

	config, db := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, gasLanePriceWei, v2.KeySpecific{
			// Gas lane.
			Key:          new(vrfKey.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		})(c, s)
		c.EVM[0].MinIncomingConfirmations = new(uint32(2))
		c.Feature.LogPoller = new(true)
		c.EVM[0].LogPollInterval = commonconfig.MustNewDuration(100 * time.Millisecond)
		c.EVM[0].FinalityDepth = new(uint32(2))
	})
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, uni.backend, ownerKey, vrfKey, bhfKey)
	require.NoError(t, app.Start(ctx))

	// Create VRF job.
	vrfJobs := createVRFJobs(
		t,
		[][]ethkey.KeyV2{{vrfKey}},
		app,
		coordinator,
		coordinatorAddress,
		batchCoordinatorAddress,
		uni,
		vrfOwnerAddress,
		vrfVersion,
		false,
		gasLanePriceWei,
	)
	keyHash := vrfJobs[0].VRFSpec.PublicKey.MustHash()
	v2plusCoordinatorAddress := coordinatorAddress.String()

	_ = vrftesthelpers.CreateAndStartBlockHeaderFeederJob(
		t, bhfKeys, app, uni.bhsContractAddress.String(), uni.batchBHSContractAddress.String(),
		"", v2plusCoordinatorAddress,
	)

	// Ensure log poller is ready and has all logs.
	chain, ok := app.GetRelayers().LegacyEVMChains().Slice()[0].(legacyevm.Chain) //nolint:staticcheck // TODO: migrate to relayer interface
	require.True(t, ok)
	require.NoError(t, chain.LogPoller().Ready())
	require.NoError(t, chain.LogPoller().Replay(ctx, 1))

	for i := range nConsumers {
		consumer := consumers[i]
		consumerContract := consumerContracts[i]

		// Create a subscription and fund with 0 LINK.
		_, subID := subscribeVRF(t, consumer, consumerContract, coordinator, uni.backend, new(big.Int), nativePayment)

		// Make the randomness request. It will not yet succeed since it is underfunded.
		numWords := uint32(20)

		requestID, requestBlock := requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, keyHash, subID, numWords, 500_000, coordinator, uni.backend, nativePayment)

		// Wait 256 blocks.
		for range 256 {
			uni.backend.Commit()
		}
		verifyBlockhashStored(t, uni, requestBlock)

		// Fund the subscription
		topUpSubscription(t, consumer, consumerContract, uni.backend, big.NewInt(5e18), nativePayment)

		// Wait for fulfillment to be queued.
		require.EventuallyWithT(t, func(c *assert.CollectT) {
			uni.backend.Commit()
			runs, err := app.PipelineORM().GetAllRuns(ctx)
			require.NoError(c, err)
			t.Log("runs", len(runs))
			require.Len(c, runs, 1)
		}, testutils.WaitTimeout(t), 100*time.Millisecond)

		mine(t, requestID, subID, uni.backend, db, vrfVersion, testutils.SimulatedChainID)

		rwfe := assertRandomWordsFulfilled(t, requestID, true, coordinator, nativePayment)
		if len(assertions) > 0 {
			assertions[0](t, coordinator, rwfe)
		}

		// Assert correct number of random words sent by coordinator.
		assertNumRandomWords(t, consumerContract, numWords)
	}
}

func createSubscriptionAndGetSubID(
	t *testing.T,
	subOwner *bind.TransactOpts,
	coordinator v22.CoordinatorV2_X,
	backend types.Backend,
) *big.Int {
	tx, err := coordinator.CreateSubscription(subOwner)
	require.NoError(t, err)
	backend.Commit()

	receipt, err := backend.Client().TransactionReceipt(t.Context(), tx.Hash())
	require.NoError(t, err)
	require.Equal(t, uint64(1), receipt.Status)
	for _, log := range receipt.Logs {
		if log.Address != coordinator.Address() {
			continue
		}
		// SubscriptionCreated(uint64 indexed subId, address owner): Topics[1] = subId
		return new(big.Int).SetBytes(log.Topics[1].Bytes())
	}
	require.FailNow(t, "no SubscriptionCreated log from coordinator in CreateSubscription receipt")
	return nil
}

func setupAndFundSubscriptionAndConsumer(
	t *testing.T,
	uni coordinatorV2UniverseCommon,
	coordinator v22.CoordinatorV2_X,
	coordinatorAddress common.Address,
	subOwner *bind.TransactOpts,
	consumerAddress common.Address,
	vrfVersion vrfcommon.Version,
	fundingAmount *big.Int,
) (subID *big.Int) {
	subID = createSubscriptionAndGetSubID(t, subOwner, coordinator, uni.backend)

	_, err := coordinator.AddConsumer(subOwner, subID, consumerAddress)
	require.NoError(t, err, "failed to add consumer")
	uni.backend.Commit()

	b, err := evmutils.ABIEncode(`[{"type":"uint256"}]`, subID)
	require.NoError(t, err)
	_, err = uni.linkContract.TransferAndCall(
		uni.sergey, coordinatorAddress, fundingAmount, b,
	)
	require.NoError(t, err, "failed to fund sub")
	uni.backend.Commit()
	return subID
}

func testSingleConsumerEIP150(
	t *testing.T,
	ownerKey ethkey.KeyV2,
	uni coordinatorV2UniverseCommon,
	batchCoordinatorAddress common.Address,
	_ bool,
	vrfVersion vrfcommon.Version,
	nativePayment bool,
) {
	ctx := t.Context()
	callBackGasLimit := int64(2_500_000) // base callback gas.

	key1 := cltest.MustGenerateRandomKey(t)
	gasLanePriceWei := assets.GWei(10)
	config, _ := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, assets.GWei(10), v2.KeySpecific{
			// Gas lane.
			Key:          new(key1.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		})(c, s)
		c.EVM[0].GasEstimator.LimitDefault = new(uint64(3.5e6))
		c.EVM[0].MinIncomingConfirmations = new(uint32(2))
		c.Feature.LogPoller = new(true)
		c.EVM[0].LogPollInterval = commonconfig.MustNewDuration(100 * time.Millisecond)
	})
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, uni.backend, ownerKey, key1)
	consumer := uni.vrfConsumers[0]
	consumerContract := uni.consumerContracts[0]
	consumerContractAddress := uni.consumerContractAddresses[0]
	// Create a subscription and fund with 500 LINK.
	subAmount := big.NewInt(1).Mul(big.NewInt(5e18), big.NewInt(100))
	subID := subscribeAndAssertSubscriptionCreatedEvent(t, consumerContract, consumer, consumerContractAddress, subAmount, uni.rootContract, uni.backend, nativePayment)

	// Fund gas lane.
	sendEth(t, ownerKey, uni.backend, key1.Address, 10)
	require.NoError(t, app.Start(ctx))

	// Create VRF job.
	jbs := createVRFJobs(
		t,
		[][]ethkey.KeyV2{{key1}},
		app,
		uni.rootContract,
		uni.rootContractAddress,
		batchCoordinatorAddress,
		uni,
		nil,
		vrfVersion,
		false,
		gasLanePriceWei,
	)
	keyHash := jbs[0].VRFSpec.PublicKey.MustHash()

	// Make the first randomness request.
	numWords := uint32(1)
	requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, keyHash, subID, numWords, uint32(callBackGasLimit), uni.rootContract, uni.backend, nativePayment)

	// Wait for simulation to pass.
	require.EventuallyWithT(t, func(c *assert.CollectT) {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns(ctx)
		require.NoError(c, err)
		t.Log("runs", len(runs))
		require.Len(c, runs, 1)
	}, testutils.WaitTimeout(t), 100*time.Millisecond)

	t.Log("Done!")
}

func testSingleConsumerEIP150Revert(
	t *testing.T,
	ownerKey ethkey.KeyV2,
	uni coordinatorV2UniverseCommon,
	batchCoordinatorAddress common.Address,
	_ bool,
	vrfVersion vrfcommon.Version,
	nativePayment bool,
) {
	ctx := t.Context()
	callBackGasLimit := uint64(2_500_000)            // base callback gas.
	eip150Fee := uint64(0)                           // no premium given for callWithExactGas
	coordinatorFulfillmentOverhead := uint64(90_000) // fixed gas used in coordinator fulfillment
	gasLimit := callBackGasLimit + eip150Fee + coordinatorFulfillmentOverhead

	key1 := cltest.MustGenerateRandomKey(t)
	gasLanePriceWei := assets.GWei(10)
	config, _ := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, assets.GWei(10), v2.KeySpecific{
			// Gas lane.
			Key:          new(key1.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		})(c, s)
		c.EVM[0].GasEstimator.LimitDefault = new(gasLimit)
		c.EVM[0].MinIncomingConfirmations = new(uint32(2))
		c.Feature.LogPoller = new(true)
		c.EVM[0].LogPollInterval = commonconfig.MustNewDuration(100 * time.Millisecond)
	})
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, uni.backend, ownerKey, key1)
	consumer := uni.vrfConsumers[0]
	consumerContract := uni.consumerContracts[0]
	consumerContractAddress := uni.consumerContractAddresses[0]
	// Create a subscription and fund with 500 LINK.
	subAmount := big.NewInt(1).Mul(big.NewInt(5e18), big.NewInt(100))
	subID := subscribeAndAssertSubscriptionCreatedEvent(t, consumerContract, consumer, consumerContractAddress, subAmount, uni.rootContract, uni.backend, nativePayment)

	// Fund gas lane.
	sendEth(t, ownerKey, uni.backend, key1.Address, 10)
	require.NoError(t, app.Start(ctx))

	// Create VRF job.
	jbs := createVRFJobs(
		t,
		[][]ethkey.KeyV2{{key1}},
		app,
		uni.rootContract,
		uni.rootContractAddress,
		batchCoordinatorAddress,
		uni,
		nil,
		vrfVersion,
		false,
		gasLanePriceWei,
	)
	keyHash := jbs[0].VRFSpec.PublicKey.MustHash()

	// Make the first randomness request.
	numWords := uint32(1)
	requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, keyHash, subID, numWords, uint32(callBackGasLimit), uni.rootContract, uni.backend, nativePayment)

	// Simulation should not pass.
	gomega.NewGomegaWithT(t).Consistently(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns(ctx)
		require.NoError(t, err)
		t.Log("runs", len(runs))
		return len(runs) == 0
	}, 5*time.Second, time.Second).Should(gomega.BeTrue())

	t.Log("Done!")
}

func testSingleConsumerBigGasCallbackSandwich(
	t *testing.T,
	ownerKey ethkey.KeyV2,
	uni coordinatorV2UniverseCommon,
	batchCoordinatorAddress common.Address,
	vrfVersion vrfcommon.Version,
	nativePayment bool,
) {
	ctx := t.Context()
	key1 := cltest.MustGenerateRandomKey(t)
	gasLanePriceWei := assets.GWei(100)
	config, db := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, assets.GWei(100), v2.KeySpecific{
			// Gas lane.
			Key:          new(key1.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		})(c, s)
		c.EVM[0].GasEstimator.LimitDefault = new(uint64(5_000_000))
		c.EVM[0].MinIncomingConfirmations = new(uint32(2))
		c.Feature.LogPoller = new(true)
		c.EVM[0].LogPollInterval = commonconfig.MustNewDuration(100 * time.Millisecond)
	})
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, uni.backend, ownerKey, key1)
	consumer := uni.vrfConsumers[0]
	consumerContract := uni.consumerContracts[0]
	consumerContractAddress := uni.consumerContractAddresses[0]

	subID := subscribeAndAssertSubscriptionCreatedEvent(t, consumerContract, consumer, consumerContractAddress, assets.Ether(2).ToInt(), uni.rootContract, uni.backend, nativePayment)

	// Fund gas lane.
	sendEth(t, ownerKey, uni.backend, key1.Address, 10)
	require.NoError(t, app.Start(ctx))

	// Create VRF job.
	jbs := createVRFJobs(
		t,
		[][]ethkey.KeyV2{{key1}},
		app,
		uni.rootContract,
		uni.rootContractAddress,
		batchCoordinatorAddress,
		uni,
		nil,
		vrfVersion,
		false,
		gasLanePriceWei,
	)
	keyHash := jbs[0].VRFSpec.PublicKey.MustHash()

	// Make some randomness requests, each one block apart, which contain a single low-gas request sandwiched between two high-gas requests.
	numWords := uint32(2)
	callbackGasLimits := []uint32{2_500_000, 50_000, 1_500_000}
	reqIDs := make([]*big.Int, 0, len(callbackGasLimits))
	for _, limit := range callbackGasLimits {
		requestID, _ := requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, keyHash, subID, numWords, limit, uni.rootContract, uni.backend, nativePayment)
		reqIDs = append(reqIDs, requestID)
		uni.backend.Commit()
	}

	// Assert that we've completed 0 runs before adding 3 new requests.
	{
		runs, err := app.PipelineORM().GetAllRuns(ctx)
		require.NoError(t, err)
		assert.Empty(t, runs)
		assert.Len(t, reqIDs, 3)
	}

	// Wait for the 50_000 gas randomness request to be enqueued.
	require.EventuallyWithT(t, func(c *assert.CollectT) {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns(ctx)
		require.NoError(c, err)
		require.Len(c, runs, 1)
	}, testutils.WaitTimeout(t), 100*time.Millisecond)

	// After the first successful request, no more will be enqueued.
	gomega.NewGomegaWithT(t).Consistently(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns(ctx)
		require.NoError(t, err)
		t.Log("assert 1", "runs", len(runs))
		return len(runs) == 1
	}, 3*time.Second, 1*time.Second).Should(gomega.BeTrue())

	// Mine the fulfillment that was queued.
	mine(t, reqIDs[1], subID, uni.backend, db, vrfVersion, testutils.SimulatedChainID)

	// Assert the random word was fulfilled
	assertRandomWordsFulfilled(t, reqIDs[1], false, uni.rootContract, nativePayment)

	// Assert that we've still only completed 1 run before adding new requests.
	{
		runs, err := app.PipelineORM().GetAllRuns(ctx)
		require.NoError(t, err)
		assert.Len(t, runs, 1)
	}

	// Make some randomness requests, each one block apart, this time without a low-gas request present in the callbackGasLimit slice.
	callbackGasLimits = []uint32{2_500_000, 2_500_000, 2_500_000}
	for _, limit := range callbackGasLimits {
		_, _ = requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, keyHash, subID, numWords, limit, uni.rootContract, uni.backend, nativePayment)
		uni.backend.Commit()
	}

	// Fulfillment will not be enqueued because subscriber doesn't have enough LINK for any of the requests.
	gomega.NewGomegaWithT(t).Consistently(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns(ctx)
		require.NoError(t, err)
		t.Log("assert 1", "runs", len(runs))
		return len(runs) == 1
	}, testutils.WaitTimeoutCustom(t, 5*time.Second), 1*time.Second).Should(gomega.BeTrue())

	t.Log("Done!")
}

func testSingleConsumerMultipleGasLanes(
	t *testing.T,
	ownerKey ethkey.KeyV2,
	uni coordinatorV2UniverseCommon,
	batchCoordinatorAddress common.Address,
	vrfVersion vrfcommon.Version,
	nativePayment bool,
) {
	ctx := t.Context()
	cheapKey := cltest.MustGenerateRandomKey(t)
	expensiveKey := cltest.MustGenerateRandomKey(t)
	cheapGasLane := assets.GWei(10)
	expensiveGasLane := assets.GWei(1000)
	config, db := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, assets.GWei(10), v2.KeySpecific{
			// Cheap gas lane.
			Key:          new(cheapKey.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: cheapGasLane},
		}, v2.KeySpecific{
			// Expensive gas lane.
			Key:          new(expensiveKey.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: expensiveGasLane},
		})(c, s)
		c.EVM[0].MinIncomingConfirmations = new(uint32(2))
		c.EVM[0].GasEstimator.LimitDefault = new(uint64(5_000_000))
		c.Feature.LogPoller = new(true)
		c.EVM[0].LogPollInterval = commonconfig.MustNewDuration(100 * time.Millisecond)
	})

	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, uni.backend, ownerKey, cheapKey, expensiveKey)
	consumer := uni.vrfConsumers[0]
	consumerContract := uni.consumerContracts[0]
	consumerContractAddress := uni.consumerContractAddresses[0]

	// Create a subscription and fund with 5 LINK.
	subID := subscribeAndAssertSubscriptionCreatedEvent(t, consumerContract, consumer, consumerContractAddress, big.NewInt(5e18), uni.rootContract, uni.backend, nativePayment)

	// Fund gas lanes.
	sendEth(t, ownerKey, uni.backend, cheapKey.Address, 10)
	sendEth(t, ownerKey, uni.backend, expensiveKey.Address, 10)
	require.NoError(t, app.Start(ctx))

	// Create VRF jobs.
	jbs := createVRFJobs(
		t,
		[][]ethkey.KeyV2{{cheapKey}, {expensiveKey}},
		app,
		uni.rootContract,
		uni.rootContractAddress,
		batchCoordinatorAddress,
		uni,
		nil,
		vrfVersion,
		false,
		cheapGasLane, expensiveGasLane,
	)
	cheapHash := jbs[0].VRFSpec.PublicKey.MustHash()
	expensiveHash := jbs[1].VRFSpec.PublicKey.MustHash()

	numWords := uint32(20)
	cheapRequestID, _ := requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, cheapHash, subID, numWords, 500_000, uni.rootContract, uni.backend, nativePayment)

	// Wait for fulfillment to be queued for cheap key hash.
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns(ctx)
		require.NoError(t, err)
		t.Log("assert 1", "runs", len(runs))
		return len(runs) == 1
	}, testutils.WaitTimeout(t), 100*time.Millisecond).Should(gomega.BeTrue())

	// Mine the fulfillment that was queued.
	mine(t, cheapRequestID, subID, uni.backend, db, vrfVersion, testutils.SimulatedChainID)

	// Assert correct state of RandomWordsFulfilled event.
	assertRandomWordsFulfilled(t, cheapRequestID, true, uni.rootContract, nativePayment)

	// Assert correct number of random words sent by coordinator.
	assertNumRandomWords(t, consumerContract, numWords)

	expensiveRequestID, _ := requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, expensiveHash, subID, numWords, 500_000, uni.rootContract, uni.backend, nativePayment)

	// We should not have any new fulfillments until a top up.
	gomega.NewWithT(t).Consistently(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns(ctx)
		require.NoError(t, err)
		t.Log("assert 2", "runs", len(runs))
		return len(runs) == 1
	}, 5*time.Second, 1*time.Second).Should(gomega.BeTrue())

	// Top up subscription with enough LINK to see the job through. 100 LINK should do the trick.
	topUpSubscription(t, consumer, consumerContract, uni.backend, decimal.RequireFromString("100e18").BigInt(), nativePayment)

	// Wait for fulfillment to be queued for expensive key hash.
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns(ctx)
		require.NoError(t, err)
		t.Log("assert 1", "runs", len(runs))
		return len(runs) == 2
	}, testutils.WaitTimeout(t), 100*time.Millisecond).Should(gomega.BeTrue())

	// Mine the fulfillment that was queued.
	mine(t, expensiveRequestID, subID, uni.backend, db, vrfVersion, testutils.SimulatedChainID)

	// Assert correct state of RandomWordsFulfilled event.
	assertRandomWordsFulfilled(t, expensiveRequestID, true, uni.rootContract, nativePayment)

	// Assert correct number of random words sent by coordinator.
	assertNumRandomWords(t, consumerContract, numWords)
}

func topUpSubscription(t *testing.T, consumer *bind.TransactOpts, consumerContract vrftesthelpers.VRFConsumerContract, backend types.Backend, fundingAmount *big.Int, nativePayment bool) {
	if nativePayment {
		_, err := consumerContract.TopUpSubscriptionNative(consumer, fundingAmount)
		require.NoError(t, err)
	} else {
		_, err := consumerContract.TopUpSubscription(consumer, fundingAmount)
		require.NoError(t, err)
	}
	backend.Commit()
}

// createVRFJobWithOverrides creates a single direct-mode VRF job whose spec
// overrides the given optional settings.
func createVRFJobWithOverrides(
	t *testing.T,
	fromKey ethkey.KeyV2,
	app *cltest.TestApplication,
	coordinator v22.CoordinatorV2_X,
	coordinatorAddress common.Address,
	uni coordinatorV2UniverseCommon,
	vrfVersion vrfcommon.Version,
	gasLanePrice *assets.Wei,
	gasOverhead *uint32,
	premiumPercentage *uint32,
) job.Job {
	ctx := t.Context()
	externalJobID := uuid.New()
	vrfkey, err := app.GetKeyStore().VRF().Create(ctx)
	require.NoError(t, err)
	spec := testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{
		JobID:                    externalJobID.String(),
		Name:                     "vrf-overrides",
		VRFVersion:               vrfVersion,
		CoordinatorAddress:       coordinatorAddress.Hex(),
		MinIncomingConfirmations: 2,
		PublicKey:                vrfkey.PublicKey.String(),
		FromAddresses:            []string{fromKey.Address.String()},
		BackoffInitialDelay:      10 * time.Millisecond,
		BackoffMaxDelay:          time.Second,
		V2:                       true,
		GasLanePrice:             gasLanePrice,
		EVMChainID:               testutils.SimulatedChainID.String(),
	}).Toml()
	if gasOverhead != nil {
		spec += fmt.Sprintf("\ngasOverhead = %d\n", *gasOverhead)
	}
	if premiumPercentage != nil {
		spec += fmt.Sprintf("\npremiumPercentage = %d\n", *premiumPercentage)
	}

	jb, err := vrfcommon.ValidatedVRFSpec(spec)
	require.NoError(t, err)
	require.Equal(t, gasOverhead, jb.VRFSpec.GasOverhead)
	require.Equal(t, premiumPercentage, jb.VRFSpec.PremiumPercentage)
	err = app.JobSpawner().CreateJob(ctx, nil, &jb)
	require.NoError(t, err)
	registerProvingKeyHelper(t, uni, coordinator, vrfkey, new(gasLanePrice.ToInt().Uint64()))
	require.Eventually(t, func() bool {
		for _, j := range app.JobSpawner().ActiveJobs() {
			if j.Type == job.VRF && j.ExternalJobID == externalJobID {
				return true
			}
		}
		return false
	}, testutils.WaitTimeout(t), 100*time.Millisecond)
	return jb
}

// testSingleConsumerWorstCaseGasFloor checks that the fulfillment gas limit is
// floored at the worst case. With a small callbackGasLimit the estimate lands
// below the floor, so the floor binds on the enqueued transaction.
func testSingleConsumerWorstCaseGasFloor(
	t *testing.T,
	ownerKey ethkey.KeyV2,
	uni coordinatorV2UniverseCommon,
	batchCoordinatorAddress common.Address,
	_ bool,
	vrfVersion vrfcommon.Version,
	nativePayment bool,
) {
	ctx := t.Context()
	callBackGasLimit := uint32(100_000)

	key1 := cltest.MustGenerateRandomKey(t)
	gasLanePriceWei := assets.GWei(10)
	config, db := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, assets.GWei(10), v2.KeySpecific{
			// Gas lane.
			Key:          new(key1.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		})(c, s)
		c.EVM[0].GasEstimator.LimitDefault = new(uint64(3.5e6))
		c.EVM[0].MinIncomingConfirmations = new(uint32(2))
		c.Feature.LogPoller = new(true)
		c.EVM[0].LogPollInterval = commonconfig.MustNewDuration(100 * time.Millisecond)
	})
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, uni.backend, ownerKey, key1)
	consumer := uni.vrfConsumers[0]
	consumerContract := uni.consumerContracts[0]
	consumerContractAddress := uni.consumerContractAddresses[0]
	// Create a subscription and fund with 5 LINK.
	subID := subscribeAndAssertSubscriptionCreatedEvent(t, consumerContract, consumer, consumerContractAddress, big.NewInt(5e18), uni.rootContract, uni.backend, nativePayment)

	// Fund gas lane.
	sendEth(t, ownerKey, uni.backend, key1.Address, 10)
	require.NoError(t, app.Start(ctx))

	// Create VRF job.
	jbs := createVRFJobs(
		t,
		[][]ethkey.KeyV2{{key1}},
		app,
		uni.rootContract,
		uni.rootContractAddress,
		batchCoordinatorAddress,
		uni,
		nil,
		vrfVersion,
		false,
		gasLanePriceWei,
	)
	keyHash := jbs[0].VRFSpec.PublicKey.MustHash()

	expectedFloor := v22.WorstCaseGasFloor(callBackGasLimit, v22.GasProofVerification)

	// The reserved fee must cover the worst-case charge: the full callback gas
	// limit at the gas lane maximum price including the premium.
	chain, ok := app.GetRelayers().LegacyEVMChains().Slice()[0].(legacyevm.Chain) //nolint:staticcheck // TODO: migrate to relayer interface
	require.True(t, ok)
	gasPriceWei := chain.Config().EVM().GasEstimator().PriceMaxKey(key1.Address)
	expectedReservation, resErr := v22.EstimateFeeJuels(callBackGasLimit, gasPriceWei.ToInt(), vrftesthelpers.WeiPerUnitLink.BigInt(), v22.DefaultPremiumPercentage)
	require.NoError(t, resErr)

	// Make the randomness request.
	numWords := uint32(1)
	requestID, _ := requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, keyHash, subID, numWords, callBackGasLimit, uni.rootContract, uni.backend, nativePayment)

	// Wait for the fulfillment to be queued.
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns(ctx)
		require.NoError(t, err)
		t.Log("runs", len(runs))
		return len(runs) == 1
	}, testutils.WaitTimeout(t), 100*time.Millisecond).Should(gomega.BeTrue())

	// Mine the fulfillment that was queued.
	require.True(t, mine(t, requestID, subID, uni.backend, db, vrfVersion, testutils.SimulatedChainID))

	assertRandomWordsFulfilled(t, requestID, true, uni.rootContract, nativePayment)
	assertNumRandomWords(t, consumerContract, numWords)

	// The mined fulfillment's gas limit must be at least the floor.
	txstore := txmgr.NewTxStore(db, logger.TestLogger(t))
	txes, err := txstore.FindTxesByMetaFieldAndStates(ctx, "GlobalSubId", subID.String(), []txmgrtypes.TxState{txmgrcommon.TxConfirmed, txmgrcommon.TxFinalized}, testutils.SimulatedChainID)
	require.NoError(t, err)
	requestIDHash := common.BytesToHash(requestID.Bytes()).String()
	found := false
	for _, tx := range txes {
		meta, err := tx.GetMeta()
		require.NoError(t, err)
		if meta.RequestID != nil && meta.RequestID.String() == requestIDHash {
			found = true
			t.Logf("fulfillment gas limit %d, floor %d", tx.FeeLimit, expectedFloor)
			require.GreaterOrEqual(t, tx.FeeLimit, expectedFloor)
			reserved := meta.MaxLink
			if nativePayment {
				reserved = meta.MaxEth
			}
			require.NotNil(t, reserved)
			reservedFee, ok := new(big.Int).SetString(*reserved, 10)
			require.True(t, ok)
			t.Logf("reserved fee %s, worst-case cost %s", reservedFee, expectedReservation)
			require.GreaterOrEqual(t, reservedFee.Cmp(expectedReservation), 0)
		}
	}
	require.True(t, found, "no fulfillment transaction found for request %s", requestIDHash)

	t.Log("Done!")
}

// testSingleConsumerWorstCaseGasFloorOverride checks that overrides set in the
// job spec raise the fulfillment gas floor and the reserved worst-case fee.
// With a cheap callback the estimate lands below even the default floor, so a
// gas limit at the override floor can only come from the job spec override.
func testSingleConsumerWorstCaseGasFloorOverride(
	t *testing.T,
	ownerKey ethkey.KeyV2,
	uni coordinatorV2UniverseCommon,
	_ common.Address,
	vrfVersion vrfcommon.Version,
	nativePayment bool,
) {
	ctx := t.Context()
	callBackGasLimit := uint32(100_000)
	gasOverheadOverride := uint32(250_000)
	premiumOverride := uint32(50)

	key1 := cltest.MustGenerateRandomKey(t)
	gasLanePriceWei := assets.GWei(10)
	config, db := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, assets.GWei(10), v2.KeySpecific{
			Key:          new(key1.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		})(c, s)
		c.EVM[0].GasEstimator.LimitDefault = new(uint64(3.5e6))
		c.EVM[0].MinIncomingConfirmations = new(uint32(2))
		c.Feature.LogPoller = new(true)
		c.EVM[0].LogPollInterval = commonconfig.MustNewDuration(100 * time.Millisecond)
	})
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, uni.backend, ownerKey, key1)
	consumer := uni.vrfConsumers[0]
	consumerContract := uni.consumerContracts[0]
	consumerContractAddress := uni.consumerContractAddresses[0]
	subID := subscribeAndAssertSubscriptionCreatedEvent(t, consumerContract, consumer, consumerContractAddress, big.NewInt(5e18), uni.rootContract, uni.backend, nativePayment)

	sendEth(t, ownerKey, uni.backend, key1.Address, 10)
	require.NoError(t, app.Start(ctx))

	jb := createVRFJobWithOverrides(t, key1, app, uni.rootContract, uni.rootContractAddress, uni, vrfVersion, gasLanePriceWei, &gasOverheadOverride, &premiumOverride)
	keyHash := jb.VRFSpec.PublicKey.MustHash()

	expectedFloor := v22.WorstCaseGasFloor(callBackGasLimit, gasOverheadOverride)
	defaultFloor := v22.WorstCaseGasFloor(callBackGasLimit, v22.GasProofVerification)
	// The override floor must exceed the default floor for this test to prove
	// the override took effect.
	require.Greater(t, expectedFloor, defaultFloor)

	chain, ok := app.GetRelayers().LegacyEVMChains().Slice()[0].(legacyevm.Chain) //nolint:staticcheck // TODO: migrate to relayer interface
	require.True(t, ok)
	gasPriceWei := chain.Config().EVM().GasEstimator().PriceMaxKey(key1.Address)
	expectedReservation, expErr := v22.EstimateFeeJuels(callBackGasLimit, gasPriceWei.ToInt(), vrftesthelpers.WeiPerUnitLink.BigInt(), premiumOverride)
	require.NoError(t, expErr)
	// The premium override differs from the default, so a reservation at the
	// override value can only come from the job spec.
	defaultReservation, defErr := v22.EstimateFeeJuels(callBackGasLimit, gasPriceWei.ToInt(), vrftesthelpers.WeiPerUnitLink.BigInt(), v22.DefaultPremiumPercentage)
	require.NoError(t, defErr)
	require.Negative(t, expectedReservation.Cmp(defaultReservation)) // override premium is lower than the default

	numWords := uint32(1)
	requestID, _ := requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, keyHash, subID, numWords, callBackGasLimit, uni.rootContract, uni.backend, nativePayment)

	gomega.NewGomegaWithT(t).Eventually(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns(ctx)
		require.NoError(t, err)
		return len(runs) == 1
	}, testutils.WaitTimeout(t), 100*time.Millisecond).Should(gomega.BeTrue())

	require.True(t, mine(t, requestID, subID, uni.backend, db, vrfVersion, testutils.SimulatedChainID))

	assertRandomWordsFulfilled(t, requestID, true, uni.rootContract, nativePayment)
	assertNumRandomWords(t, consumerContract, numWords)

	txstore := txmgr.NewTxStore(db, logger.TestLogger(t))
	txes, err := txstore.FindTxesByMetaFieldAndStates(ctx, "GlobalSubId", subID.String(), []txmgrtypes.TxState{txmgrcommon.TxConfirmed, txmgrcommon.TxFinalized}, testutils.SimulatedChainID)
	require.NoError(t, err)
	requestIDHash := common.BytesToHash(requestID.Bytes()).String()
	found := false
	for _, tx := range txes {
		meta, err := tx.GetMeta()
		require.NoError(t, err)
		if meta.RequestID != nil && meta.RequestID.String() == requestIDHash {
			found = true
			t.Logf("fulfillment gas limit %d, override floor %d, default floor %d", tx.FeeLimit, expectedFloor, defaultFloor)
			require.GreaterOrEqual(t, tx.FeeLimit, expectedFloor)
			reserved := meta.MaxLink
			if nativePayment {
				reserved = meta.MaxEth
			}
			require.NotNil(t, reserved)
			reservedFee, ok := new(big.Int).SetString(*reserved, 10)
			require.True(t, ok)
			t.Logf("reserved fee %s, premium-override cost %s, default cost %s", reservedFee, expectedReservation, defaultReservation)
			require.GreaterOrEqual(t, reservedFee.Cmp(expectedReservation), 0)
			require.Negative(t, reservedFee.Cmp(defaultReservation))
		}
	}
	require.True(t, found, "no fulfillment transaction found for request %s", requestIDHash)

	t.Log("Done!")
}

// testSingleConsumerWorstCaseReservation checks that the node does not process
// a request whose subscription balance is below the worst-case fulfillment
// cost, and that topping the subscription up unblocks it.
func testSingleConsumerWorstCaseReservation(
	t *testing.T,
	ownerKey ethkey.KeyV2,
	uni coordinatorV2UniverseCommon,
	batchCoordinatorAddress common.Address,
	vrfVersion vrfcommon.Version,
	_ bool,
) {
	ctx := t.Context()
	callBackGasLimit := uint32(100_000)

	key1 := cltest.MustGenerateRandomKey(t)
	gasLanePriceWei := assets.GWei(10)
	config, db := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, assets.GWei(10), v2.KeySpecific{
			Key:          new(key1.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		})(c, s)
		c.EVM[0].MinIncomingConfirmations = new(uint32(2))
		c.Feature.LogPoller = new(true)
		c.EVM[0].LogPollInterval = commonconfig.MustNewDuration(100 * time.Millisecond)
	})
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, uni.backend, ownerKey, key1)
	consumer := uni.vrfConsumers[0]
	consumerContract := uni.consumerContracts[0]
	consumerContractAddress := uni.consumerContractAddresses[0]

	sendEth(t, ownerKey, uni.backend, key1.Address, 10)
	require.NoError(t, app.Start(ctx))

	chain, ok := app.GetRelayers().LegacyEVMChains().Slice()[0].(legacyevm.Chain) //nolint:staticcheck // TODO: migrate to relayer interface
	require.True(t, ok)
	gasPriceWei := chain.Config().EVM().GasEstimator().PriceMaxKey(key1.Address)
	maxCost, err := v22.EstimateFeeJuels(callBackGasLimit, gasPriceWei.ToInt(), vrftesthelpers.WeiPerUnitLink.BigInt(), v22.DefaultPremiumPercentage)
	require.NoError(t, err)

	// Fund below the worst-case cost but above the simulated charge.
	initialFund := new(big.Int).Mul(maxCost, big.NewInt(70))
	initialFund.Div(initialFund, big.NewInt(100))
	subID := subscribeAndAssertSubscriptionCreatedEvent(t, consumerContract, consumer, consumerContractAddress, initialFund, uni.rootContract, uni.backend, false)

	jbs := createVRFJobs(t, [][]ethkey.KeyV2{{key1}}, app, uni.rootContract, uni.rootContractAddress, batchCoordinatorAddress, uni, nil, vrfVersion, false, gasLanePriceWei)
	keyHash := jbs[0].VRFSpec.PublicKey.MustHash()

	numWords := uint32(1)
	requestID, _ := requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, keyHash, subID, numWords, callBackGasLimit, uni.rootContract, uni.backend, false)

	// The request must not be fulfilled while the balance is below the
	// worst-case cost.
	txstore := txmgr.NewTxStore(db, logger.TestLogger(t))
	allStates := []txmgrtypes.TxState{
		txmgrcommon.TxUnstarted, txmgrcommon.TxInProgress, txmgrcommon.TxUnconfirmed,
		txmgrcommon.TxConfirmed, txmgrcommon.TxFinalized,
	}
	require.Never(t, func() bool {
		uni.backend.Commit()
		txes, txErr := txstore.FindTxesByMetaFieldAndStates(ctx, "GlobalSubId", subID.String(), allStates, testutils.SimulatedChainID)
		require.NoError(t, txErr)
		return len(txes) > 0
	}, 15*time.Second, 250*time.Millisecond)

	// Top up to the worst-case cost; the request is fulfilled on the next pass.
	_, err = consumerContract.TopUpSubscription(consumer, maxCost)
	require.NoError(t, err)
	uni.backend.Commit()
	require.True(t, mine(t, requestID, subID, uni.backend, db, vrfVersion, testutils.SimulatedChainID))

	assertRandomWordsFulfilled(t, requestID, true, uni.rootContract, false)
	assertNumRandomWords(t, consumerContract, numWords)

	t.Log("Done!")
}

func testSingleConsumerAlwaysRevertingCallbackStillFulfilled(
	t *testing.T,
	ownerKey ethkey.KeyV2,
	uni coordinatorV2UniverseCommon,
	batchCoordinatorAddress common.Address,
	_ bool,
	vrfVersion vrfcommon.Version,
	nativePayment bool,
) {
	ctx := t.Context()
	key := cltest.MustGenerateRandomKey(t)
	gasLanePriceWei := assets.GWei(10)
	config, db := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, assets.GWei(10), v2.KeySpecific{
			// Gas lane.
			Key:          new(key.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		})(c, s)
		c.EVM[0].MinIncomingConfirmations = new(uint32(2))
		c.Feature.LogPoller = new(true)
		c.EVM[0].LogPollInterval = commonconfig.MustNewDuration(100 * time.Millisecond)
	})
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, uni.backend, ownerKey, key)
	consumer := uni.reverter
	consumerContract := uni.revertingConsumerContract
	consumerContractAddress := uni.revertingConsumerContractAddress

	// Create a subscription and fund with 5 LINK.
	subID := subscribeAndAssertSubscriptionCreatedEvent(t, consumerContract, consumer, consumerContractAddress, big.NewInt(5e18), uni.rootContract, uni.backend, nativePayment)

	// Fund gas lane.
	sendEth(t, ownerKey, uni.backend, key.Address, 10)
	require.NoError(t, app.Start(ctx))

	// Create VRF job.
	jbs := createVRFJobs(
		t,
		[][]ethkey.KeyV2{{key}},
		app,
		uni.rootContract,
		uni.rootContractAddress,
		batchCoordinatorAddress,
		uni,
		nil,
		vrfVersion,
		false,
		gasLanePriceWei,
	)
	keyHash := jbs[0].VRFSpec.PublicKey.MustHash()

	// Make the randomness request.
	numWords := uint32(20)
	requestID, _ := requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, keyHash, subID, numWords, 500_000, uni.rootContract, uni.backend, nativePayment)

	// Wait for fulfillment to be queued.
	require.EventuallyWithT(t, func(c *assert.CollectT) {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns(ctx)
		require.NoError(c, err)
		require.Len(c, runs, 1)
	}, testutils.WaitTimeout(t), 100*time.Millisecond)

	// Mine the fulfillment that was queued.
	mine(t, requestID, subID, uni.backend, db, vrfVersion, testutils.SimulatedChainID)

	// Assert correct state of RandomWordsFulfilled event.
	assertRandomWordsFulfilled(t, requestID, false, uni.rootContract, nativePayment)
}

func testConsumerProxyHappyPath(
	t *testing.T,
	ownerKey ethkey.KeyV2,
	uni coordinatorV2UniverseCommon,
	batchCoordinatorAddress common.Address,
	vrfVersion vrfcommon.Version,
	nativePayment bool,
) {
	ctx := t.Context()
	key1 := cltest.MustGenerateRandomKey(t)
	key2 := cltest.MustGenerateRandomKey(t)
	gasLanePriceWei := assets.GWei(10)
	config, db := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, assets.GWei(10), v2.KeySpecific{
			// Gas lane.
			Key:          new(key1.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		}, v2.KeySpecific{
			Key:          new(key2.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		})(c, s)
		c.EVM[0].MinIncomingConfirmations = new(uint32(2))
		c.Feature.LogPoller = new(true)
		c.EVM[0].LogPollInterval = commonconfig.MustNewDuration(100 * time.Millisecond)
	})
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, uni.backend, ownerKey, key1, key2)
	consumerOwner := uni.neil
	consumerContract := uni.consumerProxyContract
	consumerContractAddress := uni.consumerProxyContractAddress

	// Create a subscription and fund with 5 LINK.
	subID := subscribeAndAssertSubscriptionCreatedEvent(
		t, consumerContract, consumerOwner, consumerContractAddress,
		assets.Ether(5).ToInt(), uni.rootContract, uni.backend, nativePayment,
	)

	// Create gas lane.
	sendEth(t, ownerKey, uni.backend, key1.Address, 10)
	sendEth(t, ownerKey, uni.backend, key2.Address, 10)
	require.NoError(t, app.Start(ctx))

	// Create VRF job using key1 and key2 on the same gas lane.
	jbs := createVRFJobs(
		t,
		[][]ethkey.KeyV2{{key1, key2}},
		app,
		uni.rootContract,
		uni.rootContractAddress,
		batchCoordinatorAddress,
		uni,
		nil,
		vrfVersion,
		false,
		gasLanePriceWei,
	)
	keyHash := jbs[0].VRFSpec.PublicKey.MustHash()

	// Make the first randomness request.
	numWords := uint32(20)
	requestID1, _ := requestRandomnessAndAssertRandomWordsRequestedEvent(
		t, consumerContract, consumerOwner, keyHash, subID, numWords, 750_000, uni.rootContract, uni.backend, nativePayment,
	)

	// Wait for fulfillment to be queued.
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns(ctx)
		require.NoError(t, err)
		t.Log("runs", len(runs))
		return len(runs) == 1
	}, testutils.WaitTimeout(t), 100*time.Millisecond).Should(gomega.BeTrue())

	// Mine the fulfillment that was queued.
	mine(t, requestID1, subID, uni.backend, db, vrfVersion, testutils.SimulatedChainID)

	// Assert correct state of RandomWordsFulfilled event.
	assertRandomWordsFulfilled(t, requestID1, true, uni.rootContract, nativePayment)

	// Gas available will be around 724,385, which means that 750,000 - 724,385 = 25,615 gas was used.
	// This is ~20k more than what the non-proxied consumer uses.
	// So to be safe, users should probably over-estimate their fulfillment gas by ~25k.
	{
		gasAvailable, err := consumerContract.SGasAvailable(nil)
		require.NoError(t, err)
		t.Log("gas available after proxied callback:", gasAvailable)
	}

	// Make the second randomness request and assert fulfillment is successful
	requestID2, _ := requestRandomnessAndAssertRandomWordsRequestedEvent(
		t, consumerContract, consumerOwner, keyHash, subID, numWords, 750_000, uni.rootContract, uni.backend, nativePayment,
	)
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns(ctx)
		require.NoError(t, err)
		t.Log("runs", len(runs))
		return len(runs) == 2
	}, testutils.WaitTimeout(t), 100*time.Millisecond).Should(gomega.BeTrue())
	mine(t, requestID2, subID, uni.backend, db, vrfVersion, testutils.SimulatedChainID)
	assertRandomWordsFulfilled(t, requestID2, true, uni.rootContract, nativePayment)

	// Assert correct number of random words sent by coordinator.
	assertNumRandomWords(t, consumerContract, numWords)

	// Assert that both send addresses were used to fulfill the requests
	n, err := uni.backend.Client().PendingNonceAt(ctx, key1.Address)
	require.NoError(t, err)
	require.EqualValues(t, 1, n)

	n, err = uni.backend.Client().PendingNonceAt(ctx, key2.Address)
	require.NoError(t, err)
	require.EqualValues(t, 1, n)
}

func testConsumerProxyCoordinatorZeroAddress(
	t *testing.T,
	uni coordinatorV2UniverseCommon,
) {
	// Deploy another upgradeable consumer, proxy, and proxy admin
	// to test vrfCoordinator != 0x0 condition.
	upgradeableConsumerAddress, _, _, err := vrf_consumer_v2_plus_upgradeable_example.DeployVRFConsumerV2PlusUpgradeableExample(uni.neil, uni.backend.Client())
	require.NoError(t, err, "failed to deploy upgradeable consumer to simulated ethereum blockchain")
	uni.backend.Commit()

	// Deployment should revert if we give the 0x0 address for the coordinator.
	upgradeableAbi, err := vrf_consumer_v2_plus_upgradeable_example.VRFConsumerV2PlusUpgradeableExampleMetaData.GetAbi()
	require.NoError(t, err)
	initializeCalldata, err := upgradeableAbi.Pack("initialize",
		common.BytesToAddress(common.LeftPadBytes([]byte{}, 20)), // zero address for the coordinator
		uni.linkContractAddress)
	require.NoError(t, err)
	_, _, err = deployOZTransparentUpgradeableProxy(
		uni.neil, uni.backend.Client(), upgradeableConsumerAddress, uni.proxyAdminAddress, initializeCalldata,
	)
	require.Error(t, err)
}

func testMaliciousConsumer(
	t *testing.T,
	ownerKey ethkey.KeyV2,
	uni coordinatorV2UniverseCommon,
	batchCoordinatorAddress common.Address,
	batchEnabled bool,
	vrfVersion vrfcommon.Version,
) {
	ctx := t.Context()
	config, _ := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].GasEstimator.LimitDefault = new(uint64(2_000_000))
		c.EVM[0].GasEstimator.PriceMax = assets.GWei(1)
		c.EVM[0].GasEstimator.PriceDefault = assets.GWei(1)
		c.EVM[0].GasEstimator.FeeCapDefault = assets.GWei(1)
		c.EVM[0].ChainID = (*sqlutil.Big)(testutils.SimulatedChainID)
		c.Feature.LogPoller = new(true)
		c.EVM[0].LogPollInterval = commonconfig.MustNewDuration(100 * time.Millisecond)
	})
	carol := uni.vrfConsumers[0]

	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, uni.backend, ownerKey)
	require.NoError(t, app.Start(ctx))

	err := app.GetKeyStore().Unlock(ctx, cltest.Password)
	require.NoError(t, err)
	vrfkey, err := app.GetKeyStore().VRF().Create(ctx)
	require.NoError(t, err)

	jid := uuid.New()
	incomingConfs := 2
	s := testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{
		JobID:                    jid.String(),
		Name:                     "vrf-primary",
		VRFVersion:               vrfVersion,
		FromAddresses:            []string{ownerKey.Address.String()},
		CoordinatorAddress:       uni.rootContractAddress.String(),
		BatchCoordinatorAddress:  batchCoordinatorAddress.String(),
		MinIncomingConfirmations: incomingConfs,
		GasLanePrice:             assets.GWei(1),
		PublicKey:                vrfkey.PublicKey.String(),
		V2:                       true,
		EVMChainID:               testutils.SimulatedChainID.String(),
	}).Toml()
	jb, err := vrfcommon.ValidatedVRFSpec(s)
	require.NoError(t, err)
	err = app.JobSpawner().CreateJob(ctx, nil, &jb)
	require.NoError(t, err)
	require.Eventually(t, func() bool {
		jbs := app.JobSpawner().ActiveJobs()
		for _, j := range jbs {
			if j.ID == jb.ID {
				return true
			}
		}
		return false
	}, testutils.WaitTimeout(t), 100*time.Millisecond)

	// Register a proving key associated with the VRF job.
	registerProvingKeyHelper(t, uni, uni.rootContract, vrfkey, &defaultMaxGasPrice)

	subFunding := decimal.RequireFromString("1000000000000000000")
	_, err = uni.maliciousConsumerContract.CreateSubscriptionAndFund(carol,
		subFunding.BigInt())
	require.NoError(t, err)
	uni.backend.Commit()

	// Send a re-entrant request
	// subID, nConfs, callbackGas, numWords are hard-coded within the contract, so setting them to 0 here
	_, err = uni.maliciousConsumerContract.RequestRandomness(carol, vrfkey.PublicKey.MustHash(), big.NewInt(0), 0, 0, 0, false)
	require.NoError(t, err)

	// We expect the request to be serviced
	// by the node.
	var attempts []txmgr.TxAttempt
	require.Eventually(t, func() bool {
		attempts, _, err = app.TxmStorageService().TxAttempts(ctx, 0, 1000)
		require.NoError(t, err)
		// It possible that we send the test request
		// before the job spawner has started the vrf services, which is fine
		// the lb will backfill the logs. However we need to
		// keep blocks coming in for the lb to send the backfilled logs.
		t.Log("attempts", attempts)
		uni.backend.Commit()
		return len(attempts) == 1 && attempts[0].Tx.State == txmgrcommon.TxConfirmed
	}, testutils.WaitTimeout(t), 100*time.Millisecond)

	// The fulfillment tx should succeed
	cs, err := app.GetRelayers().LegacyEVMChains().Get(evmtest.MustGetDefaultChainID(t, config.EVMConfigs()).String()) //nolint:staticcheck // TODO: migrate to relayer interface
	require.NoError(t, err)
	ch, ok := cs.(legacyevm.Chain)
	require.True(t, ok)
	r, err := ch.Client().TransactionReceipt(ctx, attempts[0].Hash)
	require.NoError(t, err)
	require.Equal(t, uint64(1), r.Status)

	// The user callback should have errored; parse the fulfillment event directly from the receipt.
	var fulfillment v22.RandomWordsFulfilled
	for _, log := range r.Logs {
		event, parseErr := uni.rootContract.ParseRandomWordsFulfilled(*log)
		if parseErr == nil {
			fulfillment = event
			break
		}
	}
	require.NotNil(t, fulfillment, "no RandomWordsFulfilled event in fulfillment receipt")
	require.False(t, fulfillment.Success())
	// A reverted callback reverts all state and events within it, including any re-entrant
	// requestRandomness call, so there is no separate check needed for request count.
}

func testReplayOldRequestsOnStartUp(
	t *testing.T,
	ownerKey ethkey.KeyV2,
	uni coordinatorV2UniverseCommon,
	consumer *bind.TransactOpts,
	consumerContract vrftesthelpers.VRFConsumerContract,
	consumerContractAddress common.Address,
	coordinator v22.CoordinatorV2_X,
	coordinatorAddress common.Address,
	batchCoordinatorAddress common.Address,
	vrfOwnerAddress *common.Address,
	vrfVersion vrfcommon.Version,
	nativePayment bool,
	assertions ...func(
		t *testing.T,
		coordinator v22.CoordinatorV2_X,
		rwfe v22.RandomWordsFulfilled,
		subID *big.Int,
	),
) {
	ctx := t.Context()
	sendingKey := cltest.MustGenerateRandomKey(t)
	gasLanePriceWei := assets.GWei(10)
	config, _ := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, assets.GWei(10), v2.KeySpecific{
			// Gas lane.
			Key:          new(sendingKey.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		})(c, s)
		c.EVM[0].MinIncomingConfirmations = new(uint32(2))
		c.Feature.LogPoller = new(true)
		c.EVM[0].LogPollInterval = commonconfig.MustNewDuration(100 * time.Millisecond)
	})
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, uni.backend, ownerKey, sendingKey)

	// Create a subscription and fund with 5 LINK.
	subID := subscribeAndAssertSubscriptionCreatedEvent(t, consumerContract, consumer, consumerContractAddress, big.NewInt(5e18), coordinator, uni.backend, nativePayment)

	// Fund gas lanes.
	sendEth(t, ownerKey, uni.backend, sendingKey.Address, 10)
	require.NoError(t, app.Start(ctx))

	// Create VRF Key, register it to coordinator and export
	vrfkey, err := app.GetKeyStore().VRF().Create(ctx)
	require.NoError(t, err)
	registerProvingKeyHelper(t, uni, coordinator, vrfkey, &defaultMaxGasPrice)
	keyHash := vrfkey.PublicKey.MustHash()

	encodedVrfKey, err := app.GetKeyStore().VRF().Export(vrfkey.ID(), testutils.Password)
	require.NoError(t, err)

	// Shut down the node before making the randomness request
	require.NoError(t, app.Stop())

	// Make the first randomness request.
	numWords := uint32(20)
	requestID1, _ := requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, keyHash, subID, numWords, 500_000, coordinator, uni.backend, nativePayment)

	// number of blocks to mine before restarting the node
	nBlocks := 100
	for range nBlocks {
		uni.backend.Commit()
	}

	config, db := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, assets.GWei(10), v2.KeySpecific{
			// Gas lane.
			Key:          new(sendingKey.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		})(c, s)
		c.EVM[0].MinIncomingConfirmations = new(uint32(2))
		c.Feature.LogPoller = new(true)
		c.EVM[0].LogPollInterval = commonconfig.MustNewDuration(100 * time.Millisecond)
	})

	// Start a new app and create VRF job using the same VRF key created above
	app = cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, uni.backend, ownerKey, sendingKey)

	require.NoError(t, app.Start(ctx))

	vrfKey, err := app.GetKeyStore().VRF().Import(ctx, encodedVrfKey, testutils.Password)
	require.NoError(t, err)

	incomingConfs := 2

	spec := testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{
		Name:                     "vrf-primary",
		VRFVersion:               vrfVersion,
		CoordinatorAddress:       coordinatorAddress.Hex(),
		BatchCoordinatorAddress:  batchCoordinatorAddress.Hex(),
		MinIncomingConfirmations: incomingConfs,
		PublicKey:                vrfKey.PublicKey.String(),
		FromAddresses:            []string{sendingKey.Address.String()},
		BackoffInitialDelay:      10 * time.Millisecond,
		BackoffMaxDelay:          time.Second,
		V2:                       true,
		GasLanePrice:             gasLanePriceWei,
		EVMChainID:               testutils.SimulatedChainID.String(),
	}).Toml()

	jb, err := vrfcommon.ValidatedVRFSpec(spec)
	require.NoError(t, err)
	t.Log(jb.VRFSpec.PublicKey.MustHash(), vrfKey.PublicKey.MustHash())
	err = app.JobSpawner().CreateJob(ctx, nil, &jb)
	require.NoError(t, err)

	// Wait until all jobs are active and listening for logs
	require.Eventually(t, func() bool {
		jbs := app.JobSpawner().ActiveJobs()
		for _, jb := range jbs {
			if jb.Type == job.VRF {
				return true
			}
		}
		return false
	}, testutils.WaitTimeout(t), 100*time.Millisecond)

	// Wait for fulfillment to be queued.
	require.EventuallyWithT(t, func(c *assert.CollectT) {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns(ctx)
		require.NoError(c, err)
		require.Len(c, runs, 1)
	}, testutils.WaitTimeout(t), 100*time.Millisecond)

	// Mine the fulfillment that was queued.
	mine(t, requestID1, subID, uni.backend, db, vrfVersion, testutils.SimulatedChainID)

	// Assert correct state of RandomWordsFulfilled event.
	// In particular:
	// * success should be true
	// * payment should be exactly the amount specified as the premium in the coordinator fee config
	rwfe := assertRandomWordsFulfilled(t, requestID1, true, coordinator, nativePayment)
	for _, assertion := range assertions {
		assertion(t, coordinator, rwfe, subID)
	}
}
