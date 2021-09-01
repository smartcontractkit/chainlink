package vrf_test

import (
	"context"
	"math/big"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_malicious_consumer_v2"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	evmconfig "github.com/smartcontractkit/chainlink/core/chains/evm/config"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/mock_v3_aggregator_contract"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_consumer_v2"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/signatures/secp256k1"
	"github.com/smartcontractkit/chainlink/core/services/vrf"
	"github.com/smartcontractkit/chainlink/core/services/vrf/proof"
	"github.com/smartcontractkit/chainlink/core/testdata/testspecs"
)

type coordinatorV2Universe struct {
	// Golang wrappers ofr solidity contracts
	rootContract                     *vrf_coordinator_v2.VRFCoordinatorV2
	linkContract                     *link_token_interface.LinkToken
	consumerContract                 *vrf_consumer_v2.VRFConsumerV2
	maliciousConsumerContract        *vrf_malicious_consumer_v2.VRFMaliciousConsumerV2
	rootContractAddress              common.Address
	consumerContractAddress          common.Address
	maliciousConsumerContractAddress common.Address
	linkContractAddress              common.Address
	// Abstraction representation of the ethereum blockchain
	backend        *backends.SimulatedBackend
	coordinatorABI *abi.ABI
	consumerABI    *abi.ABI
	// Cast of participants
	sergey  *bind.TransactOpts // Owns all the LINK initially
	neil    *bind.TransactOpts // Node operator running VRF service
	ned     *bind.TransactOpts // Secondary node operator
	carol   *bind.TransactOpts // Author of consuming contract which requests randomness
	nallory *bind.TransactOpts // Author of consuming contract which requests randomness
}

var (
	weiPerUnitLink = decimal.RequireFromString("10000000000000000")
)

func newVRFCoordinatorV2Universe(t *testing.T, key ethkey.KeyV2) coordinatorV2Universe {
	oracleTransactor := cltest.MustNewSimulatedBackendKeyedTransactor(t, key.ToEcdsaPrivKey())
	var (
		sergey  = newIdentity(t)
		neil    = newIdentity(t)
		ned     = newIdentity(t)
		carol   = newIdentity(t)
		nallory = oracleTransactor
	)
	genesisData := core.GenesisAlloc{
		sergey.From:  {Balance: oneEth},
		neil.From:    {Balance: oneEth},
		ned.From:     {Balance: oneEth},
		carol.From:   {Balance: oneEth},
		nallory.From: {Balance: oneEth},
	}
	gasLimit := ethconfig.Defaults.Miner.GasCeil
	consumerABI, err := abi.JSON(strings.NewReader(
		vrf_consumer_v2.VRFConsumerV2ABI))
	require.NoError(t, err)
	coordinatorABI, err := abi.JSON(strings.NewReader(
		vrf_coordinator_v2.VRFCoordinatorV2ABI))
	require.NoError(t, err)
	backend := cltest.NewSimulatedBackend(t, genesisData, gasLimit)
	// Deploy link
	linkAddress, _, linkContract, err := link_token_interface.DeployLinkToken(
		sergey, backend)
	require.NoError(t, err, "failed to deploy link contract to simulated ethereum blockchain")
	// Deploy feed
	linkEthFeed, _, _, err :=
		mock_v3_aggregator_contract.DeployMockV3AggregatorContract(
			carol, backend, 18, weiPerUnitLink.BigInt()) // 0.01 eth per link
	require.NoError(t, err)
	// Deploy coordinator
	coordinatorAddress, _, coordinatorContract, err :=
		vrf_coordinator_v2.DeployVRFCoordinatorV2(
			neil, backend, linkAddress, common.Address{} /*blockHash store*/, linkEthFeed /* linkEth*/)
	require.NoError(t, err, "failed to deploy VRFCoordinator contract to simulated ethereum blockchain")
	backend.Commit()
	// Deploy consumer it has 1 LINK
	consumerContractAddress, _, consumerContract, err :=
		vrf_consumer_v2.DeployVRFConsumerV2(
			carol, backend, coordinatorAddress, linkAddress)
	require.NoError(t, err, "failed to deploy VRFConsumer contract to simulated ethereum blockchain")
	_, err = linkContract.Transfer(sergey, consumerContractAddress, oneEth) // Actually, LINK
	require.NoError(t, err, "failed to send LINK to VRFConsumer contract on simulated ethereum blockchain")
	backend.Commit()

	// Deploy malicious consumer with 1 link
	maliciousConsumerContractAddress, _, maliciousConsumerContract, err :=
		vrf_malicious_consumer_v2.DeployVRFMaliciousConsumerV2(
			carol, backend, coordinatorAddress, linkAddress)
	require.NoError(t, err, "failed to deploy VRFMaliciousConsumer contract to simulated ethereum blockchain")
	_, err = linkContract.Transfer(sergey, maliciousConsumerContractAddress, oneEth) // Actually, LINK
	require.NoError(t, err, "failed to send LINK to VRFMaliciousConsumer contract on simulated ethereum blockchain")
	backend.Commit()

	// Set the configuration on the coordinator.
	_, err = coordinatorContract.SetConfig(neil,
		uint16(1),    // minRequestConfirmations
		uint32(1000), // 0.0001 link flat fee
		uint32(1000000),
		uint32(60*60*24),                       // stalenessSeconds
		uint32(vrf.GasAfterPaymentCalculation), // gasAfterPaymentCalculation
		big.NewInt(10000000000000000),          // 0.01 eth per link fallbackLinkPrice
		big.NewInt(1000000000000000000),        // Minimum subscription balance 0.01 link
	)
	require.NoError(t, err, "failed to set coordinator configuration")
	backend.Commit()

	return coordinatorV2Universe{
		rootContract:                     coordinatorContract,
		rootContractAddress:              coordinatorAddress,
		linkContract:                     linkContract,
		linkContractAddress:              linkAddress,
		consumerContract:                 consumerContract,
		consumerContractAddress:          consumerContractAddress,
		maliciousConsumerContract:        maliciousConsumerContract,
		maliciousConsumerContractAddress: maliciousConsumerContractAddress,
		backend:                          backend,
		coordinatorABI:                   &coordinatorABI,
		consumerABI:                      &consumerABI,
		sergey:                           sergey,
		neil:                             neil,
		ned:                              ned,
		carol:                            carol,
		nallory:                          nallory,
	}
}

func TestIntegrationVRFV2(t *testing.T) {
	config, _, cleanupDB := heavyweight.FullTestORM(t, "vrf_v2_integration", true, true)
	defer cleanupDB()
	key := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, key)
	config.Overrides.GlobalEvmGasLimitDefault = null.IntFrom(2000000)

	gasPrice := decimal.NewFromBigInt(evmconfig.DefaultGasPrice, 0)

	app, cleanup := cltest.NewApplicationWithConfigAndKeyOnSimulatedBlockchain(t, config, uni.backend, key)
	defer cleanup()
	require.NoError(t, app.Start())

	vrfkey, err := app.GetKeyStore().VRF().Create()
	require.NoError(t, err)

	jid := uuid.NewV4()
	incomingConfs := 2
	s := testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{
		JobID:              jid.String(),
		Name:               "vrf-primary",
		CoordinatorAddress: uni.rootContractAddress.String(),
		Confirmations:      incomingConfs,
		PublicKey:          vrfkey.PublicKey.String(),
		V2:                 true,
	}).Toml()
	jb, err := vrf.ValidatedVRFSpec(s)
	require.NoError(t, err)
	jb, err = app.JobORM().CreateJob(context.Background(), &jb, jb.Pipeline)
	require.NoError(t, err)

	// Register a proving key associated with the VRF job.
	p, err := vrfkey.PublicKey.Point()
	require.NoError(t, err)
	_, err = uni.rootContract.RegisterProvingKey(
		uni.neil, uni.nallory.From, pair(secp256k1.Coordinates(p)))
	require.NoError(t, err)
	uni.backend.Commit()

	// Create and fund a subscription.
	// We should see that our subscription has 1 link.
	AssertLinkBalances(t, uni.linkContract, []common.Address{
		uni.consumerContractAddress,
		uni.rootContractAddress,
	}, []*big.Int{
		big.NewInt(1000000000000000000), // 1 link
		big.NewInt(0),                   // 0 link
	})
	subFunding := decimal.RequireFromString("1000000000000000000")
	_, err = uni.consumerContract.TestCreateSubscriptionAndFund(uni.carol,
		subFunding.BigInt())
	require.NoError(t, err)
	uni.backend.Commit()
	AssertLinkBalances(t, uni.linkContract, []common.Address{
		uni.consumerContractAddress,
		uni.rootContractAddress,
		uni.nallory.From, // Oracle's own address should have nothing
	}, []*big.Int{
		big.NewInt(0),
		big.NewInt(1000000000000000000),
		big.NewInt(0),
	})
	subId, err := uni.consumerContract.SSubId(nil)
	require.NoError(t, err)
	subStart, err := uni.rootContract.GetSubscription(nil, subId)
	require.NoError(t, err)

	// Make a request for random words.
	// By requesting 500k callback with a configured eth gas limit default of 500k,
	// we ensure that the job is indeed adjusting the gaslimit to suit the users request.
	gasRequested := 500000
	nw := 10
	requestedIncomingConfs := 3
	_, err = uni.consumerContract.TestRequestRandomness(uni.carol, vrfkey.PublicKey.MustHash(), subId, uint16(requestedIncomingConfs), uint32(gasRequested), uint32(nw))
	require.NoError(t, err)

	// Oracle tries to withdraw before its fullfilled should fail
	_, err = uni.rootContract.OracleWithdraw(uni.nallory, uni.nallory.From, big.NewInt(1000))
	require.Error(t, err)

	for i := 0; i < requestedIncomingConfs; i++ {
		uni.backend.Commit()
	}

	// We expect the request to be serviced
	// by the node.
	var runs []pipeline.Run
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		runs, err = app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		// It possible that we send the test request
		// before the job spawner has started the vrf services, which is fine
		// the lb will backfill the logs. However we need to
		// keep blocks coming in for the lb to send the backfilled logs.
		uni.backend.Commit()
		return len(runs) == 1 && runs[0].State == pipeline.RunStatusCompleted
	}, 5*time.Second, 1*time.Second).Should(gomega.BeTrue())

	// Wait for the request to be fulfilled on-chain.
	var rf []*vrf_coordinator_v2.VRFCoordinatorV2RandomWordsFulfilled
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		rfIterator, err2 := uni.rootContract.FilterRandomWordsFulfilled(nil, nil)
		require.NoError(t, err2, "failed to logs")
		for rfIterator.Next() {
			rf = append(rf, rfIterator.Event)
		}
		return len(rf) == 1
	}, 5*time.Second, 500*time.Millisecond).Should(gomega.BeTrue())
	assert.True(t, rf[0].Success, "expected callback to succeed")
	fulfillReceipt, err := uni.backend.TransactionReceipt(context.Background(), rf[0].Raw.TxHash)
	require.NoError(t, err)

	// Assert all the random words received by the consumer are different and non-zero.
	seen := make(map[string]struct{})
	var rw *big.Int
	for i := 0; i < nw; i++ {
		rw, err = uni.consumerContract.SRandomWords(nil, big.NewInt(int64(i)))
		require.NoError(t, err)
		_, ok := seen[rw.String()]
		assert.False(t, ok)
		seen[rw.String()] = struct{}{}
	}

	// We should have exactly as much gas as we requested
	// after accounting for function look up code, argument decoding etc.
	// which should be fixed in this test.
	ga, err := uni.consumerContract.SGasAvailable(nil)
	require.NoError(t, err)
	gaDecoding := big.NewInt(0).Add(ga, big.NewInt(1556))
	assert.Equal(t, 0, gaDecoding.Cmp(big.NewInt(int64(gasRequested))), "expected gas available %v to exceed gas requested %v", gaDecoding, gasRequested)
	t.Log("gas available", ga.String())

	// Assert that we were only charged for how much gas we actually used.
	// We should be charged for the verification + our callbacks execution in link.
	subEnd, err := uni.rootContract.GetSubscription(nil, subId)
	require.NoError(t, err)
	var (
		end   = decimal.RequireFromString(subEnd.Balance.String())
		start = decimal.RequireFromString(subStart.Balance.String())
		wei   = decimal.RequireFromString("1000000000000000000")
		gwei  = decimal.RequireFromString("1000000000")
	)
	t.Log("end balance", end)
	linkWeiCharged := start.Sub(end)
	// Remove flat fee of 0.0001 to get fee for just gas.
	linkCharged := linkWeiCharged.Sub(decimal.RequireFromString("100000000000000")).Div(wei)
	t.Logf("subscription charged %s with gas prices of %s gwei and %s ETH per LINK\n", linkCharged, gasPrice.Div(gwei), weiPerUnitLink.Div(wei))
	expected := decimal.RequireFromString(strconv.Itoa(int(fulfillReceipt.GasUsed))).Mul(gasPrice).Div(weiPerUnitLink)
	t.Logf("expected sub charge gas use %v %v off by %v", fulfillReceipt.GasUsed, expected, expected.Sub(linkCharged))
	// The expected sub charge should be within 200 gas of the actual gas usage.
	// wei/link * link / wei/gas = wei / (wei/gas) = gas
	gasDiff := linkCharged.Sub(expected).Mul(weiPerUnitLink).Div(gasPrice).Abs().IntPart()
	t.Log("gasDiff", gasDiff)
	assert.Less(t, gasDiff, int64(200))

	// Oracle tries to withdraw move than it was paid should fail
	_, err = uni.rootContract.OracleWithdraw(uni.nallory, uni.nallory.From, linkWeiCharged.Add(decimal.NewFromInt(1)).BigInt())
	require.Error(t, err)

	// Assert the oracle can withdraw its payment.
	_, err = uni.rootContract.OracleWithdraw(uni.nallory, uni.nallory.From, linkWeiCharged.BigInt())
	require.NoError(t, err)
	uni.backend.Commit()
	AssertLinkBalances(t, uni.linkContract, []common.Address{
		uni.consumerContractAddress,
		uni.rootContractAddress,
		uni.nallory.From, // Oracle's own address should have nothing
	}, []*big.Int{
		big.NewInt(0),
		subFunding.Sub(linkWeiCharged).BigInt(),
		linkWeiCharged.BigInt(),
	})

	// We should see the response count present
	counts := vrf.GetStartingResponseCountsV2(app.Store.DB, app.Logger)
	t.Log(counts, rf[0].RequestId.String())
	assert.Equal(t, uint64(1), counts[rf[0].RequestId.String()])
}

func TestMaliciousConsumer(t *testing.T) {
	config, _, cleanupDB := heavyweight.FullTestORM(t, "vrf_v2_integration_malicious", true, true)
	defer cleanupDB()
	key := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, key)
	config.Overrides.GlobalEvmGasLimitDefault = null.IntFrom(2000000)

	app, cleanup := cltest.NewApplicationWithConfigAndKeyOnSimulatedBlockchain(t, config, uni.backend, key)
	defer cleanup()
	require.NoError(t, app.Start())

	err := app.GetKeyStore().Unlock(cltest.Password)
	require.NoError(t, err)
	vrfkey, err := app.GetKeyStore().VRF().Create()
	require.NoError(t, err)

	jid := uuid.NewV4()
	incomingConfs := 2
	s := testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{
		JobID:              jid.String(),
		Name:               "vrf-primary",
		CoordinatorAddress: uni.rootContractAddress.String(),
		Confirmations:      incomingConfs,
		PublicKey:          vrfkey.PublicKey.String(),
		V2:                 true,
	}).Toml()
	jb, err := vrf.ValidatedVRFSpec(s)
	require.NoError(t, err)
	jb, err = app.JobORM().CreateJob(context.Background(), &jb, jb.Pipeline)
	require.NoError(t, err)
	time.Sleep(1 * time.Second)

	// Register a proving key associated with the VRF job.
	p, err := vrfkey.PublicKey.Point()
	require.NoError(t, err)
	_, err = uni.rootContract.RegisterProvingKey(
		uni.neil, uni.nallory.From, pair(secp256k1.Coordinates(p)))
	require.NoError(t, err)

	_, err = uni.maliciousConsumerContract.SetKeyHash(uni.carol,
		vrfkey.PublicKey.MustHash())
	require.NoError(t, err)
	subFunding := decimal.RequireFromString("1000000000000000000")
	_, err = uni.maliciousConsumerContract.TestCreateSubscriptionAndFund(uni.carol,
		subFunding.BigInt())
	require.NoError(t, err)
	uni.backend.Commit()

	// Send a re-entrant request
	_, err = uni.maliciousConsumerContract.TestRequestRandomness(uni.carol)
	require.NoError(t, err)

	// We expect the request to be serviced
	// by the node.
	var attempts []bulletprooftxmanager.EthTxAttempt
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		//runs, err = app.PipelineORM().GetAllRuns()
		attempts, _, err = app.GetStore().EthTxAttempts(0, 1000)
		require.NoError(t, err)
		// It possible that we send the test request
		// before the job spawner has started the vrf services, which is fine
		// the lb will backfill the logs. However we need to
		// keep blocks coming in for the lb to send the backfilled logs.
		uni.backend.Commit()
		return len(attempts) == 1 && attempts[0].EthTx.State == bulletprooftxmanager.EthTxConfirmed
	}, 5*time.Second, 1*time.Second).Should(gomega.BeTrue())

	// The fulfillment tx should succeed
	ch, err := app.GetChainSet().Default()
	require.NoError(t, err)
	r, err := ch.Client().TransactionReceipt(context.Background(), attempts[0].Hash)
	require.NoError(t, err)
	require.Equal(t, uint64(1), r.Status)

	// The user callback should have errored
	it, err := uni.rootContract.FilterRandomWordsFulfilled(nil, nil)
	require.NoError(t, err)
	var fulfillments []*vrf_coordinator_v2.VRFCoordinatorV2RandomWordsFulfilled
	for it.Next() {
		fulfillments = append(fulfillments, it.Event)
	}
	require.Equal(t, 1, len(fulfillments))
	require.Equal(t, false, fulfillments[0].Success)

	// It should not have succeeded in placing another request.
	it2, err2 := uni.rootContract.FilterRandomWordsRequested(nil, nil, nil)
	require.NoError(t, err2)
	var requests []*vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested
	for it2.Next() {
		requests = append(requests, it2.Event)
	}
	require.Equal(t, 1, len(requests))
}

func TestRequestCost(t *testing.T) {
	key := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, key)

	cfg := cltest.NewTestGeneralConfig(t)
	app, cleanup := cltest.NewApplicationWithConfigAndKeyOnSimulatedBlockchain(t, cfg, uni.backend, key)
	defer cleanup()
	require.NoError(t, app.Start())

	vrfkey, err := app.GetKeyStore().VRF().Create()
	require.NoError(t, err)
	p, err := vrfkey.PublicKey.Point()
	require.NoError(t, err)
	_, err = uni.rootContract.RegisterProvingKey(
		uni.neil, uni.neil.From, pair(secp256k1.Coordinates(p)))
	require.NoError(t, err)
	uni.backend.Commit()
	_, err = uni.consumerContract.TestCreateSubscriptionAndFund(uni.carol,
		big.NewInt(1000000000000000000)) // 0.1 LINK
	require.NoError(t, err)
	uni.backend.Commit()
	subId, err := uni.consumerContract.SSubId(nil)
	require.NoError(t, err)
	// Ensure even with large number of consumers its still cheap
	var addrs []common.Address
	for i := 0; i < 99; i++ {
		addrs = append(addrs, cltest.NewAddress())
	}
	_, err = uni.consumerContract.UpdateSubscription(uni.carol,
		addrs) // 0.1 LINK
	require.NoError(t, err)
	estimate := estimateGas(t, uni.backend, common.Address{},
		uni.consumerContractAddress, uni.consumerABI,
		"testRequestRandomness", vrfkey.PublicKey.MustHash(), subId, uint16(2), uint32(10000), uint32(1))
	t.Log(estimate)
	// V2 should be at least (87000-134000)/134000 = 35% cheaper
	// Note that a second call drops further to 68998 gas, but would also drop in V1.
	assert.Less(t, estimate, uint64(85000),
		"requestRandomness tx gas cost more than expected")
}

func TestMaxConsumersCost(t *testing.T) {
	key := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, key)

	cfg := cltest.NewTestGeneralConfig(t)
	app, cleanup := cltest.NewApplicationWithConfigAndKeyOnSimulatedBlockchain(t, cfg, uni.backend, key)
	defer cleanup()
	require.NoError(t, app.Start())
	_, err := uni.consumerContract.TestCreateSubscriptionAndFund(uni.carol,
		big.NewInt(1000000000000000000)) // 0.1 LINK
	require.NoError(t, err)
	uni.backend.Commit()
	subId, err := uni.consumerContract.SSubId(nil)
	require.NoError(t, err)
	var addrs []common.Address
	for i := 0; i < 98; i++ {
		addrs = append(addrs, cltest.NewAddress())
	}
	_, err = uni.consumerContract.UpdateSubscription(uni.carol, addrs)
	// Ensure even with max number of consumers its still reasonable gas costs.
	require.NoError(t, err)
	estimate := estimateGas(t, uni.backend, uni.consumerContractAddress,
		uni.rootContractAddress, uni.coordinatorABI,
		"removeConsumer", subId, uni.consumerContractAddress)
	t.Log(estimate)
	assert.Less(t, estimate, uint64(265000))
	estimate = estimateGas(t, uni.backend, uni.consumerContractAddress,
		uni.rootContractAddress, uni.coordinatorABI,
		"addConsumer", subId, cltest.NewAddress())
	t.Log(estimate)
	assert.Less(t, estimate, uint64(100000))
}

func TestFulfillmentCost(t *testing.T) {
	key := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, key)

	cfg := cltest.NewTestGeneralConfig(t)
	app, cleanup := cltest.NewApplicationWithConfigAndKeyOnSimulatedBlockchain(t, cfg, uni.backend, key)
	defer cleanup()
	require.NoError(t, app.Start())

	vrfkey, err := app.GetKeyStore().VRF().Create()
	require.NoError(t, err)
	p, err := vrfkey.PublicKey.Point()
	require.NoError(t, err)
	_, err = uni.rootContract.RegisterProvingKey(
		uni.neil, uni.neil.From, pair(secp256k1.Coordinates(p)))
	require.NoError(t, err)
	uni.backend.Commit()
	_, err = uni.consumerContract.TestCreateSubscriptionAndFund(uni.carol,
		big.NewInt(1000000000000000000)) // 0.1 LINK
	require.NoError(t, err)
	uni.backend.Commit()
	subId, err := uni.consumerContract.SSubId(nil)
	require.NoError(t, err)

	gasRequested := 50000
	nw := 1
	requestedIncomingConfs := 3
	_, err = uni.consumerContract.TestRequestRandomness(uni.carol, vrfkey.PublicKey.MustHash(), subId, uint16(requestedIncomingConfs), uint32(gasRequested), uint32(nw))
	require.NoError(t, err)
	for i := 0; i < requestedIncomingConfs; i++ {
		uni.backend.Commit()
	}

	requestLog := FindLatestRandomnessRequestedLog(t, uni.rootContract, vrfkey.PublicKey.MustHash())
	s, err := proof.BigToSeed(requestLog.PreSeed)
	require.NoError(t, err)
	proof, rc, err := proof.GenerateProofResponseV2(app.GetKeyStore().VRF(), vrfkey.ID(), proof.PreSeedDataV2{
		PreSeed:          s,
		BlockHash:        requestLog.Raw.BlockHash,
		BlockNum:         requestLog.Raw.BlockNumber,
		SubId:            subId,
		CallbackGasLimit: uint32(gasRequested),
		NumWords:         uint32(nw),
		Sender:           uni.consumerContractAddress,
	})
	require.NoError(t, err)
	estimate := estimateGas(t, uni.backend, common.Address{},
		uni.rootContractAddress, uni.coordinatorABI,
		"fulfillRandomWords", proof, rc)
	t.Log("estimate", estimate)
	// Establish very rough bounds on fulfillment cost
	assert.Greater(t, estimate, uint64(130000))
	assert.Less(t, estimate, uint64(500000))
}

func FindLatestRandomnessRequestedLog(t *testing.T,
	coordContract *vrf_coordinator_v2.VRFCoordinatorV2,
	keyHash [32]byte) *vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested {
	var rf []*vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		rfIterator, err2 := coordContract.FilterRandomWordsRequested(nil, [][32]byte{keyHash}, []common.Address{})
		require.NoError(t, err2, "failed to logs")
		for rfIterator.Next() {
			rf = append(rf, rfIterator.Event)
		}
		return len(rf) >= 1
	}, 5*time.Second, 500*time.Millisecond).Should(gomega.BeTrue())
	latest := len(rf) - 1
	return rf[latest]
}

func AssertLinkBalances(t *testing.T, linkContract *link_token_interface.LinkToken, addresses []common.Address, balances []*big.Int) {
	require.Equal(t, len(addresses), len(balances))
	for i, a := range addresses {
		b, err := linkContract.BalanceOf(nil, a)
		require.NoError(t, err)
		assert.Equal(t, balances[i].String(), b.String(), "invalid balance for %v", a)
	}
}
