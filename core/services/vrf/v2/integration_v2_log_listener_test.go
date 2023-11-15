package v2_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	v22 "github.com/smartcontractkit/chainlink/v2/core/services/vrf/v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/vrfcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/vrftesthelpers"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

func testVRFSingleConsumerHappyPath(
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
	key1 := cltest.MustGenerateRandomKey(t)
	key2 := cltest.MustGenerateRandomKey(t)
	gasLanePriceWei := assets.GWei(10)
	config, db := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, assets.GWei(10), toml.KeySpecific{
			// Gas lane.
			Key:          ptr(key1.EIP55Address),
			GasEstimator: toml.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		}, toml.KeySpecific{
			// Gas lane.
			Key:          ptr(key2.EIP55Address),
			GasEstimator: toml.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		})(c, s)
		c.EVM[0].MinIncomingConfirmations = ptr[uint32](2)
		c.Feature.LogPoller = ptr(true)
		c.EVM[0].LogPollInterval = models.MustNewDuration(1 * time.Second)
	})
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, uni.backend, ownerKey, key1, key2)

	// Create a subscription and fund with 5 LINK.
	subID := subscribeAndAssertSubscriptionCreatedEvent(t, consumerContract, consumer, consumerContractAddress, big.NewInt(5e18), coordinator, uni.backend, nativePayment)

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
		runs, err := app.PipelineORM().GetAllRuns()
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
		runs, err := app.PipelineORM().GetAllRuns()
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
	n, err := uni.backend.PendingNonceAt(testutils.Context(t), key1.Address)
	require.NoError(t, err)
	require.EqualValues(t, 1, n)

	n, err = uni.backend.PendingNonceAt(testutils.Context(t), key2.Address)
	require.NoError(t, err)
	require.EqualValues(t, 1, n)

	t.Log("Done!")
}

func TestInitProcessedBlock_2NoVRFReqs(t *testing.T) {
	t.Parallel()
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, ownerKey, 1)
	testVRFSingleConsumerHappyPath(
		t,
		ownerKey,
		uni.coordinatorV2UniverseCommon,
		uni.vrfConsumers[0],
		uni.consumerContracts[0],
		uni.consumerContractAddresses[0],
		uni.rootContract,
		uni.rootContractAddress,
		uni.batchCoordinatorContractAddress,
		ptr(uni.vrfOwnerAddress),
		vrfcommon.V2,
		false,
		func(t *testing.T, coordinator v22.CoordinatorV2_X, rwfe v22.RandomWordsFulfilled, expectedSubID *big.Int) {
			require.PanicsWithValue(t, "VRF V2 RandomWordsFulfilled does not implement SubID", func() {
				rwfe.SubID()
			})
		},
	)
}
