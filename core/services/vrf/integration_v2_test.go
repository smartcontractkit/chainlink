package vrf_test

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"testing"
	"time"

	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/vrfkey"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"

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

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/mock_v3_aggregator_contract"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_consumer_v2"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
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
		sergey.From:  {Balance: assets.Ether(1000)},
		neil.From:    {Balance: assets.Ether(1000)},
		ned.From:     {Balance: assets.Ether(1000)},
		carol.From:   {Balance: assets.Ether(1000)},
		nallory.From: {Balance: assets.Ether(1000)},
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
	// Deploy consumer it has 10 LINK
	consumerContractAddress, _, consumerContract, err :=
		vrf_consumer_v2.DeployVRFConsumerV2(
			carol, backend, coordinatorAddress, linkAddress)
	require.NoError(t, err, "failed to deploy VRFConsumer contract to simulated ethereum blockchain")
	_, err = linkContract.Transfer(sergey, consumerContractAddress, assets.Ether(10)) // Actually, LINK
	require.NoError(t, err, "failed to send LINK to VRFConsumer contract on simulated ethereum blockchain")
	backend.Commit()

	// Deploy malicious consumer with 1 link
	maliciousConsumerContractAddress, _, maliciousConsumerContract, err :=
		vrf_malicious_consumer_v2.DeployVRFMaliciousConsumerV2(
			carol, backend, coordinatorAddress, linkAddress)
	require.NoError(t, err, "failed to deploy VRFMaliciousConsumer contract to simulated ethereum blockchain")
	_, err = linkContract.Transfer(sergey, maliciousConsumerContractAddress, assets.Ether(1)) // Actually, LINK
	require.NoError(t, err, "failed to send LINK to VRFMaliciousConsumer contract on simulated ethereum blockchain")
	backend.Commit()

	// Set the configuration on the coordinator.
	_, err = coordinatorContract.SetConfig(neil,
		uint16(1),                              // minRequestConfirmations
		uint32(1000000),                        // gas limit
		uint32(1000),                           // 0.001 link flat fee
		uint32(100),                            // 0.0001 link flat fee
		uint32(10),                             // 0.00001 link flat fee
		uint16(1056),                           // 00000100 00100000 // bound1=10^3, bound2=10^6
		uint32(60*60*24),                       // stalenessSeconds
		uint32(vrf.GasAfterPaymentCalculation), // gasAfterPaymentCalculation
		big.NewInt(10000000000000000),          // 0.01 eth per link fallbackLinkPrice
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

// Send eth from prefunded account.
// Amount is number of ETH not wei.
func sendEth(t *testing.T, key ethkey.KeyV2, ec *backends.SimulatedBackend, to common.Address, eth int) {
	nonce, err := ec.PendingNonceAt(context.Background(), key.Address.Address())
	require.NoError(t, err)
	tx := gethtypes.NewTx(&gethtypes.DynamicFeeTx{
		ChainID:   big.NewInt(1337),
		Nonce:     nonce,
		GasTipCap: big.NewInt(1),
		GasFeeCap: big.NewInt(10000000000), // block base fee in sim
		Gas:       uint64(21000),
		To:        &to,
		Value:     big.NewInt(0).Mul(big.NewInt(int64(eth)), big.NewInt(1000000000000000000)),
		Data:      nil,
	})
	signedTx, err := gethtypes.SignTx(tx, gethtypes.NewLondonSigner(big.NewInt(1337)), key.ToEcdsaPrivKey())
	require.NoError(t, err)
	err = ec.SendTransaction(context.Background(), signedTx)
	require.NoError(t, err)
	ec.Commit()
}

func TestIntegrationVRFV2_OffchainSimulation(t *testing.T) {
	config, _, _ := heavyweight.FullTestDB(t, "vrf_v2_integration_sim", true, true)
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, ownerKey)
	app := cltest.NewApplicationWithConfigAndKeyOnSimulatedBlockchain(t, config, uni.backend, ownerKey)
	config.Overrides.GlobalEvmGasLimitDefault = null.NewInt(0, false)
	config.Overrides.GlobalMinIncomingConfirmations = null.IntFrom(2)

	// Lets create 2 gas lanes
	key1, err := app.KeyStore.Eth().Create(big.NewInt(1337))
	require.NoError(t, err)
	sendEth(t, ownerKey, uni.backend, key1.Address.Address(), 10)
	key2, err := app.KeyStore.Eth().Create(big.NewInt(1337))
	require.NoError(t, err)
	sendEth(t, ownerKey, uni.backend, key2.Address.Address(), 10)

	gasPrice := decimal.NewFromBigInt(big.NewInt(10000000000), 0) // Default is 10 gwei
	configureSimChain(app, map[string]types.ChainCfg{
		key1.Address.String(): {
			EvmMaxGasPriceWei: utils.NewBig(big.NewInt(10000000000)), // 10 gwei
		},
		key2.Address.String(): {
			EvmMaxGasPriceWei: utils.NewBig(big.NewInt(100000000000)), // 100 gwei
		},
	}, gasPrice.BigInt())
	require.NoError(t, app.Start())

	var jbs []job.Job
	// Create separate jobs for each gas lane and register their keys
	for i, key := range []ethkey.KeyV2{key1, key2} {
		vrfkey, err := app.GetKeyStore().VRF().Create()
		require.NoError(t, err)

		jid := uuid.NewV4()
		incomingConfs := 2
		s := testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{
			JobID:              jid.String(),
			Name:               fmt.Sprintf("vrf-primary-%d", i),
			CoordinatorAddress: uni.rootContractAddress.String(),
			Confirmations:      incomingConfs,
			PublicKey:          vrfkey.PublicKey.String(),
			FromAddress:        key.Address.String(),
			V2:                 true,
		}).Toml()
		jb, err := vrf.ValidatedVRFSpec(s)
		t.Log(jb.VRFSpec.PublicKey.MustHash(), vrfkey.PublicKey.MustHash())
		require.NoError(t, err)
		jb, err = app.JobSpawner().CreateJob(context.Background(), jb, jb.Name)
		require.NoError(t, err)
		registerProvingKeyHelper(t, uni, vrfkey)
		jbs = append(jbs, jb)
	}
	// Wait until all jobs are active and listening for logs
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		jbs := app.JobSpawner().ActiveJobs()
		return len(jbs) == 2
	}, 5*time.Second, 100*time.Millisecond).Should(gomega.BeTrue())
	// Unfortunately the lb needs heads to be able to backfill logs to new subscribers.
	// To avoid confirming
	// TODO: it could just backfill immediately upon receiving a new subscriber? (though would
	// only be useful for tests, probably a more robust way is to have the job spawner accept a signal that a
	// job is fully up and running and not add it to the active jobs list before then)
	time.Sleep(2 * time.Second)

	subFunding := decimal.RequireFromString("1000000000000000000")
	_, err = uni.consumerContract.TestCreateSubscriptionAndFund(uni.carol,
		subFunding.BigInt())
	require.NoError(t, err)
	uni.backend.Commit()
	sub, err := uni.rootContract.GetSubscription(nil, uint64(1))
	require.NoError(t, err)
	t.Log("Sub balance", sub.Balance)
	for i := 0; i < 5; i++ {
		// Request 20 words (all get saved) so we use the full 300k
		_, err := uni.consumerContract.TestRequestRandomness(uni.carol, jbs[0].VRFSpec.PublicKey.MustHash(), uint64(1), uint16(2), uint32(300000), uint32(20))
		require.NoError(t, err)
	}
	// Send a requests to the high gas price max keyhash, should remain queued until
	// a significant topup
	for i := 0; i < 1; i++ {
		_, err := uni.consumerContract.TestRequestRandomness(uni.carol, jbs[1].VRFSpec.PublicKey.MustHash(), uint64(1), uint16(2), uint32(300000), uint32(20))
		require.NoError(t, err)
	}
	// Confirm all those requests
	for i := 0; i < 3; i++ {
		uni.backend.Commit()
	}
	// Now we should see ONLY 2 requests enqueued to the bptxm
	// since we only have 2 requests worth of link at the max keyhash
	// gas price.
	var runs []pipeline.Run
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		runs, err = app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		t.Log("runs", len(runs))
		return len(runs) == 2
	}, 5*time.Second, 1*time.Second).Should(gomega.BeTrue())
	// As we send new blocks, we should observe the fulfllments goes through the balance
	// be reduced.
	gomega.NewGomegaWithT(t).Consistently(func() bool {
		runs, err = app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		uni.backend.Commit()
		return len(runs) == 2
	}, 5*time.Second, 100*time.Millisecond).Should(gomega.BeTrue())
	sub, err = uni.rootContract.GetSubscription(nil, uint64(1))
	require.NoError(t, err)
	t.Log("Sub balance should be near zero", sub.Balance)
	etxes, n, err := app.BPTXMORM().EthTransactionsWithAttempts(0, 1000)
	require.Equal(t, 2, n) // Only sent 2 transactions
	// Should have max link set
	require.NotNil(t, etxes[0].Meta)
	require.NotNil(t, etxes[1].Meta)
	md := bulletprooftxmanager.EthTxMeta{}
	require.NoError(t, json.Unmarshal(*etxes[0].Meta, &md))
	require.NotEqual(t, "", md.MaxLink)
	// Now lets top up and see the next batch go through
	_, err = uni.consumerContract.TopUpSubscription(uni.carol, assets.Ether(1))
	require.NoError(t, err)
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		runs, err = app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		t.Log("runs", len(runs))
		uni.backend.Commit()
		return len(runs) == 4
	}, 10*time.Second, 1*time.Second).Should(gomega.BeTrue())
	// One more time for the final tx
	_, err = uni.consumerContract.TopUpSubscription(uni.carol, assets.Ether(1))
	require.NoError(t, err)
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		runs, err = app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		t.Log("runs", len(runs))
		uni.backend.Commit()
		return len(runs) == 5
	}, 10*time.Second, 1*time.Second).Should(gomega.BeTrue())

	// Send a huge topup and observe the high max gwei go through.
	_, err = uni.consumerContract.TopUpSubscription(uni.carol, assets.Ether(7))
	require.NoError(t, err)
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		runs, err = app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		t.Log("runs", len(runs))
		uni.backend.Commit()
		return len(runs) == 6
	}, 10*time.Second, 1*time.Second).Should(gomega.BeTrue())
}

func configureSimChain(app *cltest.TestApplication, ks map[string]types.ChainCfg, defaultGasPrice *big.Int) {
	zero := models.MustMakeDuration(0 * time.Millisecond)
	reaperThreshold := models.MustMakeDuration(100 * time.Millisecond)
	app.ChainSet.Configure(
		big.NewInt(1337),
		true,
		types.ChainCfg{
			GasEstimatorMode:                 null.StringFrom("FixedPrice"),
			EvmGasPriceDefault:               utils.NewBig(defaultGasPrice),
			EvmHeadTrackerMaxBufferSize:      null.IntFrom(100),
			EvmHeadTrackerSamplingInterval:   &zero, // Head sampling disabled
			EthTxResendAfterThreshold:        &zero,
			EvmFinalityDepth:                 null.IntFrom(15),
			EthTxReaperThreshold:             &reaperThreshold,
			MinIncomingConfirmations:         null.IntFrom(1),
			MinRequiredOutgoingConfirmations: null.IntFrom(1),
			MinimumContractPayment:           assets.NewLinkFromJuels(100),
			EvmGasLimitDefault:               null.NewInt(2000000, true),
			KeySpecific:                      ks,
		},
	)
}

func registerProvingKeyHelper(t *testing.T, uni coordinatorV2Universe, vrfkey vrfkey.KeyV2) {
	// Register a proving key associated with the VRF job.
	p, err := vrfkey.PublicKey.Point()
	require.NoError(t, err)
	_, err = uni.rootContract.RegisterProvingKey(
		uni.neil, uni.nallory.From, pair(secp256k1.Coordinates(p)))
	require.NoError(t, err)
	uni.backend.Commit()
}

func TestIntegrationVRFV2(t *testing.T) {
	config, _, _ := heavyweight.FullTestDB(t, "vrf_v2_integration", true, true)
	key := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, key)

	app := cltest.NewApplicationWithConfigAndKeyOnSimulatedBlockchain(t, config, uni.backend, key)
	config.Overrides.GlobalEvmGasLimitDefault = null.NewInt(0, false)
	config.Overrides.GlobalMinIncomingConfirmations = null.IntFrom(2)
	keys, err := app.KeyStore.Eth().SendingKeys()

	// Reconfigure the sim chain with a default gas price of 1 gwei,
	// max gas limit of 2M and a key specific max 10 gwei price.
	// Keep the prices low so we can operate with small link balance subscriptions.
	gasPrice := decimal.NewFromBigInt(big.NewInt(1000000000), 0)
	configureSimChain(app, map[string]types.ChainCfg{
		keys[0].Address.String(): {
			EvmMaxGasPriceWei: utils.NewBig(big.NewInt(10000000000)),
		},
	}, gasPrice.BigInt())

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
		FromAddress:        keys[0].Address.String(),
		V2:                 true,
	}).Toml()
	jb, err := vrf.ValidatedVRFSpec(s)
	require.NoError(t, err)
	jb, err = app.JobSpawner().CreateJob(context.Background(), jb, jb.Name)
	require.NoError(t, err)

	registerProvingKeyHelper(t, uni, vrfkey)

	// Create and fund a subscription.
	// We should see that our subscription has 1 link.
	AssertLinkBalances(t, uni.linkContract, []common.Address{
		uni.consumerContractAddress,
		uni.rootContractAddress,
	}, []*big.Int{
		assets.Ether(10), // 10 link
		big.NewInt(0),    // 0 link
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
		assets.Ether(9),
		assets.Ether(1),
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
	}, 10*time.Second, 1*time.Second).Should(gomega.BeTrue())

	// Wait for the request to be fulfilled on-chain.
	var rf []*vrf_coordinator_v2.VRFCoordinatorV2RandomWordsFulfilled
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		rfIterator, err2 := uni.rootContract.FilterRandomWordsFulfilled(nil, nil)
		require.NoError(t, err2, "failed to logs")
		uni.backend.Commit()
		for rfIterator.Next() {
			rf = append(rf, rfIterator.Event)
		}
		return len(rf) == 1
	}, 10*time.Second, 500*time.Millisecond).Should(gomega.BeTrue())
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
	// Remove flat fee of 0.001 to get fee for just gas.
	linkCharged := linkWeiCharged.Sub(decimal.RequireFromString("1000000000000000")).Div(wei)
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
		assets.Ether(9),
		subFunding.Sub(linkWeiCharged).BigInt(),
		linkWeiCharged.BigInt(),
	})

	// We should see the response count present
	counts := vrf.GetStartingResponseCountsV2(app.GetDB(), app.Logger)
	t.Log(counts, rf[0].RequestId.String())
	assert.Equal(t, uint64(1), counts[rf[0].RequestId.String()])
}

func TestMaliciousConsumer(t *testing.T) {
	config, _, _ := heavyweight.FullTestDB(t, "vrf_v2_integration_malicious", true, true)
	key := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, key)
	config.Overrides.GlobalEvmGasLimitDefault = null.IntFrom(2000000)
	config.Overrides.GlobalEvmMaxGasPriceWei = big.NewInt(1000000000)  // 1 gwei
	config.Overrides.GlobalEvmGasPriceDefault = big.NewInt(1000000000) // 1 gwei

	app := cltest.NewApplicationWithConfigAndKeyOnSimulatedBlockchain(t, config, uni.backend, key)
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
	jb, err = app.JobSpawner().CreateJob(context.Background(), jb, jb.Name)
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
		attempts, _, err = app.BPTXMORM().EthTxAttempts(0, 1000)
		require.NoError(t, err)
		// It possible that we send the test request
		// before the job spawner has started the vrf services, which is fine
		// the lb will backfill the logs. However we need to
		// keep blocks coming in for the lb to send the backfilled logs.
		t.Log("attempts", attempts)
		uni.backend.Commit()
		return len(attempts) == 1 && attempts[0].EthTx.State == bulletprooftxmanager.EthTxConfirmed
	}, 10*time.Second, 1*time.Second).Should(gomega.BeTrue())

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
	app := cltest.NewApplicationWithConfigAndKeyOnSimulatedBlockchain(t, cfg, uni.backend, key)
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
	app := cltest.NewApplicationWithConfigAndKeyOnSimulatedBlockchain(t, cfg, uni.backend, key)
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
	app := cltest.NewApplicationWithConfigAndKeyOnSimulatedBlockchain(t, cfg, uni.backend, key)
	require.NoError(t, app.Start())
	defer app.Stop()

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
	assert.Greater(t, estimate, uint64(120000))
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
