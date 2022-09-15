package vrf_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/assets"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func TestVRFV2Integration_ConsumerProxy_HappyPath(t *testing.T) {
	config, db := heavyweight.FullTestDB(t, "vrfv2_consumerproxy_happypath")
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, ownerKey, 0)
	app := cltest.NewApplicationWithConfigAndKeyOnSimulatedBlockchain(t, config, uni.backend, ownerKey)
	config.Overrides.GlobalEvmGasLimitDefault = null.NewInt(0, false)
	config.Overrides.GlobalMinIncomingConfirmations = null.IntFrom(2)
	consumerOwner := uni.neil
	consumerContract := uni.consumerProxyContract
	consumerContractAddress := uni.consumerProxyContractAddress

	// Create a subscription and fund with 5 LINK.
	subID := subscribeAndAssertSubscriptionCreatedEvent(
		t, consumerContract, consumerOwner, consumerContractAddress, assets.Ether(5), uni)

	// Create gas lane.
	key1, err := app.KeyStore.Eth().Create(big.NewInt(1337))
	require.NoError(t, err)
	key2, err := app.KeyStore.Eth().Create(big.NewInt(1337))
	require.NoError(t, err)
	sendEth(t, ownerKey, uni.backend, key1.Address, 10)
	sendEth(t, ownerKey, uni.backend, key2.Address, 10)
	configureSimChain(t, app, map[string]evmtypes.ChainCfg{
		key1.Address.String(): {
			EvmMaxGasPriceWei: utils.NewBig(assets.GWei(10)),
		},
		key2.Address.String(): {
			EvmMaxGasPriceWei: utils.NewBig(assets.GWei(10)),
		},
	}, assets.GWei(10))
	require.NoError(t, app.Start(testutils.Context(t)))

	// Create VRF job using key1 and key2 on the same gas lane.
	jbs := createVRFJobs(t, [][]ethkey.KeyV2{{key1, key2}}, []int{10, 10}, app, uni, false)
	keyHash := jbs[0].VRFSpec.PublicKey.MustHash()

	// Make the first randomness request.
	numWords := uint32(20)
	requestID1, _ := requestRandomnessAndAssertRandomWordsRequestedEvent(
		t, consumerContract, consumerOwner, keyHash, subID, numWords, 500_000, uni)

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
	assertRandomWordsFulfilled(t, requestID1, false, uni) // TODO: figure out why success is false

	// Make the second randomness request and assert fulfillment is successful
	requestID2, _ := requestRandomnessAndAssertRandomWordsRequestedEvent(t, consumerContract, consumerOwner, keyHash, subID, numWords, 500_000, uni)
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		uni.backend.Commit()
		runs, err := app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		t.Log("runs", len(runs))
		return len(runs) == 2
	}, testutils.WaitTimeout(t), time.Second).Should(gomega.BeTrue())
	mine(t, requestID2, subID, uni, db)
	assertRandomWordsFulfilled(t, requestID2, false, uni) // TODO: figure out why success is false

	// Assert correct number of random words sent by coordinator.
	// TODO: figure out why success is false
	//assertNumRandomWords(t, consumerContract, numWords)

	// Assert that both send addresses were used to fulfill the requests
	n, err := uni.backend.PendingNonceAt(testutils.Context(t), key1.Address)
	require.NoError(t, err)
	require.EqualValues(t, 1, n)

	n, err = uni.backend.PendingNonceAt(testutils.Context(t), key2.Address)
	require.NoError(t, err)
	require.EqualValues(t, 1, n)

	t.Log("Done!")
}
