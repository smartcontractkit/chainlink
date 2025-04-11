package v2_test

import (
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/google/uuid"
	"github.com/onsi/gomega"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrf_consumer_v2_upgradeable_example"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrf_external_sub_owner_example"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrfv2_transparent_upgradeable_proxy"
	"github.com/smartcontractkit/chainlink-evm/pkg/assets"
	v2 "github.com/smartcontractkit/chainlink-evm/pkg/config/toml"
	"github.com/smartcontractkit/chainlink-evm/pkg/types"
	evmutils "github.com/smartcontractkit/chainlink-evm/pkg/utils"
	ubig "github.com/smartcontractkit/chainlink-evm/pkg/utils/big"
	txmgrcommon "github.com/smartcontractkit/chainlink-framework/chains/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	v22 "github.com/smartcontractkit/chainlink/v2/core/services/vrf/v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/vrfcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/vrftesthelpers"
	"github.com/smartcontractkit/chainlink/v2/core/testdata/testspecs"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
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
	ctx := testutils.Context(t)
	key1 := cltest.MustGenerateRandomKey(t)
	key2 := cltest.MustGenerateRandomKey(t)
	gasLanePriceWei := assets.GWei(10)
	config, db := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, assets.GWei(10), v2.KeySpecific{
			// Gas lane.
			Key:          ptr[types.EIP55Address](key1.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		}, v2.KeySpecific{
			// Gas lane.
			Key:          ptr[types.EIP55Address](key2.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		})(c, s)
		c.EVM[0].MinIncomingConfirmations = ptr[uint32](2)
		c.Feature.LogPoller = ptr(true)
		c.EVM[0].LogPollInterval = commonconfig.MustNewDuration(1 * time.Second)
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
		gasLanePriceWei)
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
	}, testutils.WaitTimeout(t), time.Second).Should(gomega.BeTrue())

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
	}, testutils.WaitTimeout(t), time.Second).Should(gomega.BeTrue())
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
	consumerContractAddresses []common.Address,
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
	ctx := testutils.Context(t)
	nConsumers := len(consumers)
	vrfKey := cltest.MustGenerateRandomKey(t)
	sendEth(t, ownerKey, uni.backend, vrfKey.Address, 10)

	// generate n BHS keys to make sure BHS job rotates sending keys
	var bhsKeyAddresses []string
	var keySpecificOverrides []v2.KeySpecific
	var keys []interface{}
	gasLanePriceWei := assets.GWei(10)
	for i := 0; i < nConsumers; i++ {
		bhsKey := cltest.MustGenerateRandomKey(t)
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

	config, db := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, assets.GWei(10), keySpecificOverrides...)(c, s)
		c.EVM[0].MinIncomingConfirmations = ptr[uint32](2)
		c.Feature.LogPoller = ptr(true)
		c.EVM[0].LogPollInterval = commonconfig.MustNewDuration(1 * time.Second)
		c.EVM[0].FinalityDepth = ptr[uint32](2)
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
		gasLanePriceWei)
	keyHash := vrfJobs[0].VRFSpec.PublicKey.MustHash()

	var (
		v2CoordinatorAddress     string
		v2PlusCoordinatorAddress string
	)

	if vrfVersion == vrfcommon.V2 {
		v2CoordinatorAddress = coordinatorAddress.String()
	} else if vrfVersion == vrfcommon.V2Plus {
		v2PlusCoordinatorAddress = coordinatorAddress.String()
	}

	_ = vrftesthelpers.CreateAndStartBHSJob(
		t, bhsKeyAddresses, app, uni.bhsContractAddress.String(), "",
		v2CoordinatorAddress, v2PlusCoordinatorAddress, "", 0, 200, 0, 100)

	// Ensure log poller is ready and has all logs.
	require.NoError(t, app.GetRelayers().LegacyEVMChains().Slice()[0].LogPoller().Ready())
	require.NoError(t, app.GetRelayers().LegacyEVMChains().Slice()[0].LogPoller().Replay(ctx, 1))

	for i := 0; i < nConsumers; i++ {
		consumer := consumers[i]
		consumerContract := consumerContracts[i]

		// Create a subscription and fund with 0 LINK.
		_, subID := subscribeVRF(t, consumer, consumerContract, coordinator, uni.backend, new(big.Int), nativePayment)
		if vrfVersion == vrfcommon.V2 {
			require.Equal(t, uint64(i+1), subID.Uint64())
		}

		// Make the randomness request. It will not yet succeed since it is underfunded.
		numWords := uint32(20)

		requestID, requestBlock := requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, keyHash, subID, numWords, 500_000, coordinator, uni.backend, nativePayment)

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
		topUpSubscription(t, consumer, consumerContract, uni.backend, big.NewInt(5e18 /* 5 LINK */), nativePayment)

		// Wait for fulfillment to be queued.
		gomega.NewGomegaWithT(t).Eventually(func() bool {
			uni.backend.Commit()
			runs, err := app.PipelineORM().GetAllRuns(ctx)
			require.NoError(t, err)
			t.Log("runs", len(runs))
			return len(runs) == 1
		}, testutils.WaitTimeout(t), time.Second).Should(gomega.BeTrue())

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
	consumerContractAddresses []common.Address,
	coordinator v22.CoordinatorV2_X,
	coordinatorAddress common.Address,
	batchCoordinatorAddress common.Address,
	vrfVersion vrfcommon.Version,
	nativePayment bool,
	addedDelay bool,
	assertions ...func(
		t *testing.T,
		coordinator v22.CoordinatorV2_X,
		rwfe v22.RandomWordsFulfilled),
) {
	ctx := testutils.Context(t)
	nConsumers := len(consumers)
	vrfKey := cltest.MustGenerateRandomKey(t)
	sendEth(t, ownerKey, uni.backend, vrfKey.Address, 10)

	// generate n BHS keys to make sure BHS job rotates sending keys
	var bhsKeyAddresses []common.Address
	var bhsKeyAddressesStrings []string
	var keySpecificOverrides []v2.KeySpecific
	var keys []interface{}
	gasLanePriceWei := assets.GWei(10)
	for i := 0; i < nConsumers; i++ {
		bhsKey := cltest.MustGenerateRandomKey(t)
		bhsKeyAddressesStrings = append(bhsKeyAddressesStrings, bhsKey.Address.String())
		bhsKeyAddresses = append(bhsKeyAddresses, bhsKey.Address)
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

	// Whitelist vrf key for trusted BHS.
	{
		_, err := uni.trustedBhsContract.SetWhitelist(uni.neil, bhsKeyAddresses)
		require.NoError(t, err)
		uni.backend.Commit()
	}

	config, db := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, assets.GWei(10), keySpecificOverrides...)(c, s)
		c.EVM[0].MinIncomingConfirmations = ptr[uint32](2)
		c.EVM[0].GasEstimator.LimitDefault = ptr(uint64(5_000_000))
		c.Feature.LogPoller = ptr(true)
		c.EVM[0].LogPollInterval = commonconfig.MustNewDuration(1 * time.Second)
		c.EVM[0].FinalityDepth = ptr[uint32](2)
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
		gasLanePriceWei)
	keyHash := vrfJobs[0].VRFSpec.PublicKey.MustHash()

	var (
		v2CoordinatorAddress     string
		v2PlusCoordinatorAddress string
	)

	if vrfVersion == vrfcommon.V2 {
		v2CoordinatorAddress = coordinatorAddress.String()
	} else if vrfVersion == vrfcommon.V2Plus {
		v2PlusCoordinatorAddress = coordinatorAddress.String()
	}

	waitBlocks := 100
	if addedDelay {
		waitBlocks = 400
	}
	_ = vrftesthelpers.CreateAndStartBHSJob(
		t, bhsKeyAddressesStrings, app, "", "",
		v2CoordinatorAddress, v2PlusCoordinatorAddress, uni.trustedBhsContractAddress.String(), 20, 1000, 0, waitBlocks)

	// Ensure log poller is ready and has all logs.
	chain := app.GetRelayers().LegacyEVMChains().Slice()[0]
	require.NoError(t, chain.LogPoller().Ready())
	require.NoError(t, chain.LogPoller().Replay(ctx, 1))

	for i := 0; i < nConsumers; i++ {
		consumer := consumers[i]
		consumerContract := consumerContracts[i]

		// Create a subscription and fund with 0 LINK.
		_, subID := subscribeVRF(t, consumer, consumerContract, coordinator, uni.backend, new(big.Int), nativePayment)
		if vrfVersion == vrfcommon.V2 {
			require.Equal(t, uint64(i+1), subID.Uint64())
		}

		// Make the randomness request. It will not yet succeed since it is underfunded.
		numWords := uint32(20)

		requestID, requestBlock := requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, keyHash, subID, numWords, 500_000, coordinator, uni.backend, nativePayment)

		// Wait 101 blocks.
		for i := 0; i < 100; i++ {
			uni.backend.Commit()
		}

		// For an added delay, we even go beyond the EVM lookback limit. This is not a problem in a trusted BHS setup.
		if addedDelay {
			for i := 0; i < 300; i++ {
				uni.backend.Commit()
			}
		}

		verifyBlockhashStoredTrusted(t, uni, requestBlock)

		// Wait another 160 blocks so that the request is outside of the 256 block window
		for i := 0; i < 160; i++ {
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
		}, testutils.WaitTimeout(t), time.Second)

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
		}
		t.Fatal(err)
		return false
	}, testutils.WaitTimeout(t), time.Second).Should(gomega.BeTrue())
}

func verifyBlockhashStoredTrusted(
	t *testing.T,
	uni coordinatorV2PlusUniverse,
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
		_, err := uni.trustedBhsContract.GetBlockhash(callOpts, big.NewInt(int64(requestBlock)))
		if err == nil {
			return true
		} else if strings.Contains(err.Error(), "execution reverted") {
			return false
		}
		t.Fatal(err)
		return false
	}, time.Second*300, time.Second).Should(gomega.BeTrue())
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
	ctx := testutils.Context(t)
	key1 := cltest.MustGenerateRandomKey(t)
	gasLanePriceWei := assets.GWei(10)
	config, db := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, assets.GWei(10), v2.KeySpecific{
			// Gas lane.
			Key:          ptr(key1.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		})(c, s)
		c.EVM[0].GasEstimator.LimitDefault = ptr[uint64](5_000_000)
		c.EVM[0].MinIncomingConfirmations = ptr[uint32](2)
		c.EVM[0].ChainID = (*ubig.Big)(testutils.SimulatedChainID)
		c.Feature.LogPoller = ptr(true)
		c.EVM[0].LogPollInterval = commonconfig.MustNewDuration(1 * time.Second)
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
		gasLanePriceWei)
	keyHash := jbs[0].VRFSpec.PublicKey.MustHash()

	// Make some randomness requests.
	numWords := uint32(2)
	var reqIDs []*big.Int
	for i := 0; i < numRequests; i++ {
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
	}, testutils.WaitTimeout(t), time.Second).Should(gomega.BeTrue())

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
	ctx := testutils.Context(t)
	key := cltest.MustGenerateRandomKey(t)
	gasLanePriceWei := assets.GWei(1000)
	config, db := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, assets.GWei(1000), v2.KeySpecific{
			// Gas lane.
			Key:          ptr(key.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		})(c, s)
		c.EVM[0].MinIncomingConfirmations = ptr[uint32](2)
		c.Feature.LogPoller = ptr(true)
		c.EVM[0].LogPollInterval = commonconfig.MustNewDuration(1 * time.Second)
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
		gasLanePriceWei)
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
	}, testutils.WaitTimeout(t), 1*time.Second)

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
	consumerContractAddresses []common.Address,
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
	ctx := testutils.Context(t)
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
			Key:          ptr(vrfKey.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		})(c, s)
		c.EVM[0].MinIncomingConfirmations = ptr[uint32](2)
		c.Feature.LogPoller = ptr(true)
		c.EVM[0].LogPollInterval = commonconfig.MustNewDuration(1 * time.Second)
		c.EVM[0].FinalityDepth = ptr[uint32](2)
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
		gasLanePriceWei)
	keyHash := vrfJobs[0].VRFSpec.PublicKey.MustHash()
	var (
		v2coordinatorAddress     string
		v2plusCoordinatorAddress string
	)
	if vrfVersion == vrfcommon.V2 {
		v2coordinatorAddress = coordinatorAddress.String()
	} else if vrfVersion == vrfcommon.V2Plus {
		v2plusCoordinatorAddress = coordinatorAddress.String()
	}

	_ = vrftesthelpers.CreateAndStartBlockHeaderFeederJob(
		t, bhfKeys, app, uni.bhsContractAddress.String(), uni.batchBHSContractAddress.String(), "",
		v2coordinatorAddress, v2plusCoordinatorAddress)

	// Ensure log poller is ready and has all logs.
	require.NoError(t, app.GetRelayers().LegacyEVMChains().Slice()[0].LogPoller().Ready())
	require.NoError(t, app.GetRelayers().LegacyEVMChains().Slice()[0].LogPoller().Replay(ctx, 1))

	for i := 0; i < nConsumers; i++ {
		consumer := consumers[i]
		consumerContract := consumerContracts[i]

		// Create a subscription and fund with 0 LINK.
		_, subID := subscribeVRF(t, consumer, consumerContract, coordinator, uni.backend, new(big.Int), nativePayment)
		if vrfVersion == vrfcommon.V2 {
			require.Equal(t, uint64(i+1), subID.Uint64())
		}

		// Make the randomness request. It will not yet succeed since it is underfunded.
		numWords := uint32(20)

		requestID, requestBlock := requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, keyHash, subID, numWords, 500_000, coordinator, uni.backend, nativePayment)

		// Wait 256 blocks.
		for i := 0; i < 256; i++ {
			uni.backend.Commit()
		}
		verifyBlockhashStored(t, uni, requestBlock)

		// Fund the subscription
		topUpSubscription(t, consumer, consumerContract, uni.backend, big.NewInt(5e18), nativePayment)

		// Wait for fulfillment to be queued.
		gomega.NewGomegaWithT(t).Eventually(func() bool {
			uni.backend.Commit()
			runs, err := app.PipelineORM().GetAllRuns(ctx)
			require.NoError(t, err)
			t.Log("runs", len(runs))
			return len(runs) == 1
		}, testutils.WaitTimeout(t), time.Second).Should(gomega.BeTrue())

		mine(t, requestID, subID, uni.backend, db, vrfVersion, testutils.SimulatedChainID)

		rwfe := assertRandomWordsFulfilled(t, requestID, true, coordinator, nativePayment)
		if len(assertions) > 0 {
			assertions[0](t, coordinator, rwfe)
		}

		// Assert correct number of random words sent by coordinator.
		assertNumRandomWords(t, consumerContract, numWords)
	}
}

func createSubscriptionAndGetSubscriptionCreatedEvent(
	t *testing.T,
	subOwner *bind.TransactOpts,
	coordinator v22.CoordinatorV2_X,
	backend types.Backend,
) v22.SubscriptionCreated {
	_, err := coordinator.CreateSubscription(subOwner)
	require.NoError(t, err)
	backend.Commit()

	iter, err := coordinator.FilterSubscriptionCreated(nil, nil)
	require.NoError(t, err)
	require.True(t, iter.Next(), "could not find SubscriptionCreated event for subID")
	return iter.Event()
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
	event := createSubscriptionAndGetSubscriptionCreatedEvent(t, subOwner, coordinator, uni.backend)
	subID = event.SubID()

	_, err := coordinator.AddConsumer(subOwner, subID, consumerAddress)
	require.NoError(t, err, "failed to add consumer")
	uni.backend.Commit()

	if vrfVersion == vrfcommon.V2Plus {
		b, err2 := evmutils.ABIEncode(`[{"type":"uint256"}]`, subID)
		require.NoError(t, err2)
		_, err2 = uni.linkContract.TransferAndCall(
			uni.sergey, coordinatorAddress, fundingAmount, b)
		require.NoError(t, err2, "failed to fund sub")
		uni.backend.Commit()
		return
	}
	b, err := evmutils.ABIEncode(`[{"type":"uint64"}]`, subID.Uint64())
	require.NoError(t, err)
	_, err = uni.linkContract.TransferAndCall(
		uni.sergey, coordinatorAddress, fundingAmount, b)
	require.NoError(t, err, "failed to fund sub")
	uni.backend.Commit()
	return
}

func testSingleConsumerForcedFulfillment(
	t *testing.T,
	ownerKey ethkey.KeyV2,
	uni coordinatorV2Universe,
	coordinator v22.CoordinatorV2_X,
	coordinatorAddress common.Address,
	batchCoordinatorAddress common.Address,
	batchEnabled bool,
	vrfVersion vrfcommon.Version,
) {
	ctx := testutils.Context(t)
	key1 := cltest.MustGenerateRandomKey(t)
	key2 := cltest.MustGenerateRandomKey(t)
	gasLanePriceWei := assets.GWei(10)
	config, db := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
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
		c.Feature.LogPoller = ptr(true)
		c.EVM[0].LogPollInterval = commonconfig.MustNewDuration(1 * time.Second)
	})
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, uni.backend, ownerKey, key1, key2)

	eoaConsumerAddr, _, eoaConsumer, err := vrf_external_sub_owner_example.DeployVRFExternalSubOwnerExample(
		uni.neil,
		uni.backend.Client(),
		uni.oldRootContractAddress,
		uni.linkContractAddress,
	)
	require.NoError(t, err, "failed to deploy eoa consumer")
	uni.backend.Commit()

	// Create a subscription and fund with 5 LINK.
	subID := setupAndFundSubscriptionAndConsumer(
		t,
		uni.coordinatorV2UniverseCommon,
		uni.oldRootContract,
		uni.oldRootContractAddress,
		uni.neil,
		eoaConsumerAddr,
		vrfVersion,
		assets.Ether(5).ToInt(),
	)

	// Check the subscription state
	sub, err := uni.oldRootContract.GetSubscription(nil, subID)
	require.NoError(t, err, "failed to get subscription with id %d", subID)
	require.Equal(t, assets.Ether(5).ToInt(), sub.Balance())
	require.Len(t, sub.Consumers(), 1)
	require.Equal(t, eoaConsumerAddr, sub.Consumers()[0])
	require.Equal(t, uni.neil.From, sub.Owner())

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
		uni.coordinatorV2UniverseCommon,
		ptr(uni.vrfOwnerAddress),
		vrfVersion,
		batchEnabled,
		gasLanePriceWei)
	keyHash := jbs[0].VRFSpec.PublicKey.MustHash()

	// Transfer ownership of the VRF coordinator to the VRF owner,
	// which is critical for this test.
	_, err = uni.oldRootContract.TransferOwnership(uni.neil, uni.vrfOwnerAddress)
	require.NoError(t, err, "unable to TransferOwnership of VRF coordinator to VRFOwner")
	uni.backend.Commit()

	_, err = uni.vrfOwner.AcceptVRFOwnership(uni.neil)
	require.NoError(t, err, "unable to Accept VRF Ownership")
	uni.backend.Commit()

	actualCoordinatorAddr, err := uni.vrfOwner.GetVRFCoordinator(nil)
	require.NoError(t, err)
	require.Equal(t, uni.oldRootContractAddress, actualCoordinatorAddr)

	t.Log("vrf owner address:", uni.vrfOwnerAddress)

	// Add allowed callers so that the oracle can call fulfillRandomWords
	// on VRFOwner.
	_, err = uni.vrfOwner.SetAuthorizedSenders(uni.neil, []common.Address{
		key1.EIP55Address.Address(),
		key2.EIP55Address.Address(),
	})
	require.NoError(t, err, "unable to update authorized senders in VRFOwner")
	uni.backend.Commit()

	// Make the randomness request.
	// Give it a larger number of confs so that we have enough time to remove the consumer
	// and cause a 0 balance to the sub.
	numWords := 3
	confs := 10
	_, err = eoaConsumer.RequestRandomWords(uni.neil, subID.Uint64(), 500_000, uint16(confs), uint32(numWords), keyHash)
	require.NoError(t, err, "failed to request randomness from consumer")
	uni.backend.Commit()

	requestID, err := eoaConsumer.SRequestId(nil)
	require.NoError(t, err)

	// Remove consumer and cancel the sub before the request can be fulfilled
	_, err = uni.oldRootContract.RemoveConsumer(uni.neil, subID, eoaConsumerAddr)
	require.NoError(t, err, "RemoveConsumer tx failed")
	uni.backend.Commit()
	_, err = uni.oldRootContract.CancelSubscription(uni.neil, subID, uni.neil.From)
	require.NoError(t, err, "CancelSubscription tx failed")
	uni.backend.Commit()

	// Wait for force-fulfillment to be queued.
	require.Eventually(t, func() bool {
		uni.backend.Commit()
		commitment, err2 := uni.oldRootContract.GetCommitment(nil, requestID)
		require.NoError(t, err2)
		t.Log("commitment is:", hexutil.Encode(commitment[:]))
		it, err2 := uni.vrfOwner.FilterRandomWordsForced(nil, []*big.Int{requestID}, []uint64{subID.Uint64()}, []common.Address{eoaConsumerAddr})
		require.NoError(t, err2)
		i := 0
		for it.Next() {
			i++
			require.Equal(t, requestID.String(), it.Event.RequestId.String())
			require.Equal(t, subID.Uint64(), it.Event.SubId)
			require.Equal(t, eoaConsumerAddr.String(), it.Event.Sender.String())
		}
		t.Log("num RandomWordsForced logs:", i)
		return utils.IsEmpty(commitment[:])
	}, testutils.WaitTimeout(t), time.Second)

	// Mine the fulfillment that was queued.
	mine(t, requestID, subID, uni.backend, db, vrfVersion, testutils.SimulatedChainID)

	// Assert correct state of RandomWordsFulfilled event.
	// In this particular case:
	// * success should be true
	// * payment should be zero (forced fulfillment)
	rwfe := assertRandomWordsFulfilled(t, requestID, true, coordinator, false)
	require.Equal(t, "0", rwfe.Payment().String())

	// Check that the RandomWordsForced event is emitted correctly.
	it, err := uni.vrfOwner.FilterRandomWordsForced(nil, []*big.Int{requestID}, []uint64{subID.Uint64()}, []common.Address{eoaConsumerAddr})
	require.NoError(t, err)
	i := 0
	for it.Next() {
		i++
		require.Equal(t, requestID.String(), it.Event.RequestId.String())
		require.Equal(t, subID.Uint64(), it.Event.SubId)
		require.Equal(t, eoaConsumerAddr.String(), it.Event.Sender.String())
	}
	require.Positive(t, i)

	t.Log("Done!")
}

func testSingleConsumerEIP150(
	t *testing.T,
	ownerKey ethkey.KeyV2,
	uni coordinatorV2UniverseCommon,
	batchCoordinatorAddress common.Address,
	batchEnabled bool,
	vrfVersion vrfcommon.Version,
	nativePayment bool,
) {
	ctx := testutils.Context(t)
	callBackGasLimit := int64(2_500_000) // base callback gas.

	key1 := cltest.MustGenerateRandomKey(t)
	gasLanePriceWei := assets.GWei(10)
	config, _ := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, assets.GWei(10), v2.KeySpecific{
			// Gas lane.
			Key:          ptr(key1.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		})(c, s)
		c.EVM[0].GasEstimator.LimitDefault = ptr(uint64(3.5e6))
		c.EVM[0].MinIncomingConfirmations = ptr[uint32](2)
		c.Feature.LogPoller = ptr(true)
		c.EVM[0].LogPollInterval = commonconfig.MustNewDuration(1 * time.Second)
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
		gasLanePriceWei)
	keyHash := jbs[0].VRFSpec.PublicKey.MustHash()

	// Make the first randomness request.
	numWords := uint32(1)
	requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, keyHash, subID, numWords, uint32(callBackGasLimit), uni.rootContract, uni.backend, nativePayment)

	// Wait for simulation to pass.
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns(ctx)
		require.NoError(t, err)
		t.Log("runs", len(runs))
		return len(runs) == 1
	}, testutils.WaitTimeout(t), time.Second).Should(gomega.BeTrue())

	t.Log("Done!")
}

func testSingleConsumerEIP150Revert(
	t *testing.T,
	ownerKey ethkey.KeyV2,
	uni coordinatorV2UniverseCommon,
	batchCoordinatorAddress common.Address,
	batchEnabled bool,
	vrfVersion vrfcommon.Version,
	nativePayment bool,
) {
	ctx := testutils.Context(t)
	callBackGasLimit := int64(2_500_000)            // base callback gas.
	eip150Fee := int64(0)                           // no premium given for callWithExactGas
	coordinatorFulfillmentOverhead := int64(90_000) // fixed gas used in coordinator fulfillment
	gasLimit := callBackGasLimit + eip150Fee + coordinatorFulfillmentOverhead

	key1 := cltest.MustGenerateRandomKey(t)
	gasLanePriceWei := assets.GWei(10)
	config, _ := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, assets.GWei(10), v2.KeySpecific{
			// Gas lane.
			Key:          ptr(key1.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		})(c, s)
		c.EVM[0].GasEstimator.LimitDefault = ptr(uint64(gasLimit))
		c.EVM[0].MinIncomingConfirmations = ptr[uint32](2)
		c.Feature.LogPoller = ptr(true)
		c.EVM[0].LogPollInterval = commonconfig.MustNewDuration(1 * time.Second)
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
		gasLanePriceWei)
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
	ctx := testutils.Context(t)
	key1 := cltest.MustGenerateRandomKey(t)
	gasLanePriceWei := assets.GWei(100)
	config, db := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, assets.GWei(100), v2.KeySpecific{
			// Gas lane.
			Key:          ptr(key1.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		})(c, s)
		c.EVM[0].GasEstimator.LimitDefault = ptr[uint64](5_000_000)
		c.EVM[0].MinIncomingConfirmations = ptr[uint32](2)
		c.Feature.LogPoller = ptr(true)
		c.EVM[0].LogPollInterval = commonconfig.MustNewDuration(1 * time.Second)
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
		gasLanePriceWei)
	keyHash := jbs[0].VRFSpec.PublicKey.MustHash()

	// Make some randomness requests, each one block apart, which contain a single low-gas request sandwiched between two high-gas requests.
	numWords := uint32(2)
	reqIDs := []*big.Int{}
	callbackGasLimits := []uint32{2_500_000, 50_000, 1_500_000}
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
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns(ctx)
		require.NoError(t, err)
		t.Log("runs", len(runs))
		return len(runs) == 1
	}, testutils.WaitTimeout(t), time.Second).Should(gomega.BeTrue())

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
	}, 5*time.Second, 1*time.Second).Should(gomega.BeTrue())

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
	ctx := testutils.Context(t)
	cheapKey := cltest.MustGenerateRandomKey(t)
	expensiveKey := cltest.MustGenerateRandomKey(t)
	cheapGasLane := assets.GWei(10)
	expensiveGasLane := assets.GWei(1000)
	config, db := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, assets.GWei(10), v2.KeySpecific{
			// Cheap gas lane.
			Key:          ptr(cheapKey.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: cheapGasLane},
		}, v2.KeySpecific{
			// Expensive gas lane.
			Key:          ptr(expensiveKey.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: expensiveGasLane},
		})(c, s)
		c.EVM[0].MinIncomingConfirmations = ptr[uint32](2)
		c.EVM[0].GasEstimator.LimitDefault = ptr[uint64](5_000_000)
		c.Feature.LogPoller = ptr(true)
		c.EVM[0].LogPollInterval = commonconfig.MustNewDuration(1 * time.Second)
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
		cheapGasLane, expensiveGasLane)
	cheapHash := jbs[0].VRFSpec.PublicKey.MustHash()
	expensiveHash := jbs[1].VRFSpec.PublicKey.MustHash()

	numWords := uint32(20)
	cheapRequestID, _ :=
		requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, cheapHash, subID, numWords, 500_000, uni.rootContract, uni.backend, nativePayment)

	// Wait for fulfillment to be queued for cheap key hash.
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns(ctx)
		require.NoError(t, err)
		t.Log("assert 1", "runs", len(runs))
		return len(runs) == 1
	}, testutils.WaitTimeout(t), 1*time.Second).Should(gomega.BeTrue())

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
	}, testutils.WaitTimeout(t), 1*time.Second).Should(gomega.BeTrue())

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

func testSingleConsumerAlwaysRevertingCallbackStillFulfilled(
	t *testing.T,
	ownerKey ethkey.KeyV2,
	uni coordinatorV2UniverseCommon,
	batchCoordinatorAddress common.Address,
	batchEnabled bool,
	vrfVersion vrfcommon.Version,
	nativePayment bool,
) {
	ctx := testutils.Context(t)
	key := cltest.MustGenerateRandomKey(t)
	gasLanePriceWei := assets.GWei(10)
	config, db := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, assets.GWei(10), v2.KeySpecific{
			// Gas lane.
			Key:          ptr(key.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		})(c, s)
		c.EVM[0].MinIncomingConfirmations = ptr[uint32](2)
		c.Feature.LogPoller = ptr(true)
		c.EVM[0].LogPollInterval = commonconfig.MustNewDuration(1 * time.Second)
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
		gasLanePriceWei)
	keyHash := jbs[0].VRFSpec.PublicKey.MustHash()

	// Make the randomness request.
	numWords := uint32(20)
	requestID, _ := requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumer, keyHash, subID, numWords, 500_000, uni.rootContract, uni.backend, nativePayment)

	// Wait for fulfillment to be queued.
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns(ctx)
		require.NoError(t, err)
		t.Log("runs", len(runs))
		return len(runs) == 1
	}, testutils.WaitTimeout(t), 1*time.Second).Should(gomega.BeTrue())

	// Mine the fulfillment that was queued.
	mine(t, requestID, subID, uni.backend, db, vrfVersion, testutils.SimulatedChainID)

	// Assert correct state of RandomWordsFulfilled event.
	assertRandomWordsFulfilled(t, requestID, false, uni.rootContract, nativePayment)
	t.Log("Done!")
}

func testConsumerProxyHappyPath(
	t *testing.T,
	ownerKey ethkey.KeyV2,
	uni coordinatorV2UniverseCommon,
	batchCoordinatorAddress common.Address,
	vrfVersion vrfcommon.Version,
	nativePayment bool,
) {
	ctx := testutils.Context(t)
	key1 := cltest.MustGenerateRandomKey(t)
	key2 := cltest.MustGenerateRandomKey(t)
	gasLanePriceWei := assets.GWei(10)
	config, db := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, assets.GWei(10), v2.KeySpecific{
			// Gas lane.
			Key:          ptr(key1.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		}, v2.KeySpecific{
			Key:          ptr(key2.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		})(c, s)
		c.EVM[0].MinIncomingConfirmations = ptr[uint32](2)
		c.Feature.LogPoller = ptr(true)
		c.EVM[0].LogPollInterval = commonconfig.MustNewDuration(1 * time.Second)
	})
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, uni.backend, ownerKey, key1, key2)
	consumerOwner := uni.neil
	consumerContract := uni.consumerProxyContract
	consumerContractAddress := uni.consumerProxyContractAddress

	// Create a subscription and fund with 5 LINK.
	subID := subscribeAndAssertSubscriptionCreatedEvent(
		t, consumerContract, consumerOwner, consumerContractAddress,
		assets.Ether(5).ToInt(), uni.rootContract, uni.backend, nativePayment)

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
		gasLanePriceWei)
	keyHash := jbs[0].VRFSpec.PublicKey.MustHash()

	// Make the first randomness request.
	numWords := uint32(20)
	requestID1, _ := requestRandomnessAndAssertRandomWordsRequestedEvent(
		t, consumerContract, consumerOwner, keyHash, subID, numWords, 750_000, uni.rootContract, uni.backend, nativePayment)

	// Wait for fulfillment to be queued.
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns(ctx)
		require.NoError(t, err)
		t.Log("runs", len(runs))
		return len(runs) == 1
	}, testutils.WaitTimeout(t), time.Second).Should(gomega.BeTrue())

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
		t, consumerContract, consumerOwner, keyHash, subID, numWords, 750_000, uni.rootContract, uni.backend, nativePayment)
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns(ctx)
		require.NoError(t, err)
		t.Log("runs", len(runs))
		return len(runs) == 2
	}, testutils.WaitTimeout(t), time.Second).Should(gomega.BeTrue())
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

	t.Log("Done!")
}

func testConsumerProxyCoordinatorZeroAddress(
	t *testing.T,
	uni coordinatorV2UniverseCommon,
) {
	// Deploy another upgradeable consumer, proxy, and proxy admin
	// to test vrfCoordinator != 0x0 condition.
	upgradeableConsumerAddress, _, _, err := vrf_consumer_v2_upgradeable_example.DeployVRFConsumerV2UpgradeableExample(uni.neil, uni.backend.Client())
	require.NoError(t, err, "failed to deploy upgradeable consumer to simulated ethereum blockchain")
	uni.backend.Commit()

	// Deployment should revert if we give the 0x0 address for the coordinator.
	upgradeableAbi, err := vrf_consumer_v2_upgradeable_example.VRFConsumerV2UpgradeableExampleMetaData.GetAbi()
	require.NoError(t, err)
	initializeCalldata, err := upgradeableAbi.Pack("initialize",
		common.BytesToAddress(common.LeftPadBytes([]byte{}, 20)), // zero address for the coordinator
		uni.linkContractAddress)
	require.NoError(t, err)
	_, _, _, err = vrfv2_transparent_upgradeable_proxy.DeployVRFV2TransparentUpgradeableProxy(
		uni.neil, uni.backend.Client(), upgradeableConsumerAddress, uni.proxyAdminAddress, initializeCalldata)
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
	ctx := testutils.Context(t)
	config, _ := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].GasEstimator.LimitDefault = ptr[uint64](2_000_000)
		c.EVM[0].GasEstimator.PriceMax = assets.GWei(1)
		c.EVM[0].GasEstimator.PriceDefault = assets.GWei(1)
		c.EVM[0].GasEstimator.FeeCapDefault = assets.GWei(1)
		c.EVM[0].ChainID = (*ubig.Big)(testutils.SimulatedChainID)
		c.Feature.LogPoller = ptr(true)
		c.EVM[0].LogPollInterval = commonconfig.MustNewDuration(1 * time.Second)
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
	time.Sleep(1 * time.Second)

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
	}, testutils.WaitTimeout(t), 1*time.Second)

	// The fulfillment tx should succeed
	ch, err := app.GetRelayers().LegacyEVMChains().Get(evmtest.MustGetDefaultChainID(t, config.EVMConfigs()).String())
	require.NoError(t, err)
	r, err := ch.Client().TransactionReceipt(ctx, attempts[0].Hash)
	require.NoError(t, err)
	require.Equal(t, uint64(1), r.Status)

	// The user callback should have errored
	it, err := uni.rootContract.FilterRandomWordsFulfilled(nil, nil, nil)
	require.NoError(t, err)
	var fulfillments []v22.RandomWordsFulfilled
	for it.Next() {
		fulfillments = append(fulfillments, it.Event())
	}
	require.Len(t, fulfillments, 1)
	require.False(t, fulfillments[0].Success())

	// It should not have succeeded in placing another request.
	it2, err2 := uni.rootContract.FilterRandomWordsRequested(nil, nil, nil, nil)
	require.NoError(t, err2)
	var requests []v22.RandomWordsRequested
	for it2.Next() {
		requests = append(requests, it2.Event())
	}
	require.Len(t, requests, 1)
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
		subID *big.Int),
) {
	ctx := testutils.Context(t)
	sendingKey := cltest.MustGenerateRandomKey(t)
	gasLanePriceWei := assets.GWei(10)
	config, _ := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, assets.GWei(10), v2.KeySpecific{
			// Gas lane.
			Key:          ptr(sendingKey.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		})(c, s)
		c.EVM[0].MinIncomingConfirmations = ptr[uint32](2)
		c.Feature.LogPoller = ptr(true)
		c.EVM[0].LogPollInterval = commonconfig.MustNewDuration(1 * time.Second)
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
	for i := 0; i < nBlocks; i++ {
		uni.backend.Commit()
	}

	config, db := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, assets.GWei(10), v2.KeySpecific{
			// Gas lane.
			Key:          ptr(sendingKey.EIP55Address),
			GasEstimator: v2.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		})(c, s)
		c.EVM[0].MinIncomingConfirmations = ptr[uint32](2)
		c.Feature.LogPoller = ptr(true)
		c.EVM[0].LogPollInterval = commonconfig.MustNewDuration(1 * time.Second)
	})

	// Start a new app and create VRF job using the same VRF key created above
	app = cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, uni.backend, ownerKey, sendingKey)

	require.NoError(t, app.Start(ctx))

	vrfKey, err := app.GetKeyStore().VRF().Import(ctx, encodedVrfKey, testutils.Password)
	require.NoError(t, err)

	incomingConfs := 2
	var vrfOwnerString string
	if vrfOwnerAddress != nil {
		vrfOwnerString = vrfOwnerAddress.Hex()
	}

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
		VRFOwnerAddress:          vrfOwnerString,
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
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns(ctx)
		require.NoError(t, err)
		t.Log("runs", len(runs))
		return len(runs) == 1
	}, testutils.WaitTimeout(t), time.Second).Should(gomega.BeTrue())

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
}
