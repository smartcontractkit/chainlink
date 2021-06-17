package vrf_test

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/mock_v3_aggregator_contract"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_consumer_v2"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/signatures/secp256k1"
	"github.com/smartcontractkit/chainlink/core/services/vrf"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/testdata/testspecs"
	"github.com/stretchr/testify/require"
	"math/big"
	"strings"
	"testing"
	"time"
)

type coordinatorV2Universe struct {
	// Golang wrappers ofr solidity contracts
	rootContract               *vrf_coordinator_v2.VRFCoordinatorV2
	linkContract               *link_token_interface.LinkToken
	consumerContract           *vrf_consumer_v2.VRFConsumerV2
	rootContractAddress        common.Address
	consumerContractAddress    common.Address
	linkContractAddress        common.Address
	// Abstraction representation of the ethereum blockchain
	backend        *backends.SimulatedBackend
	coordinatorABI *abi.ABI
	consumerABI    *abi.ABI
	// Cast of participants
	sergey *bind.TransactOpts // Owns all the LINK initially
	neil   *bind.TransactOpts // Node operator running VRF service
	ned    *bind.TransactOpts // Secondary node operator
	carol  *bind.TransactOpts // Author of consuming contract which requests randomness
}

func newVRFCoordinatorV2Universe(t *testing.T, key models.Key) coordinatorV2Universe {
	k, err := keystore.DecryptKey(key.JSON.RawMessage[:], cltest.Password)
	require.NoError(t, err)
	oracleTransactor := cltest.MustNewSimulatedBackendKeyedTransactor(t, k.PrivateKey)
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
	backend := backends.NewSimulatedBackend(genesisData, gasLimit)
	// Deploy link
	linkAddress, _, linkContract, err := link_token_interface.DeployLinkToken(
		sergey, backend)
	require.NoError(t, err, "failed to deploy link contract to simulated ethereum blockchain")
	// Deploy feeds
	fastGasFeed, _, _, err :=
		mock_v3_aggregator_contract.DeployMockV3AggregatorContract(
			carol, backend, 0, big.NewInt(1000000000)) // 1 gwei per unit gas
	linkEthFeed, _, _, err :=
		mock_v3_aggregator_contract.DeployMockV3AggregatorContract(
			carol, backend, 18, big.NewInt(10000000000000000)) // 0.01 eth per link
	// Deploy coordinator
	coordinatorAddress, _, coordinatorContract, err :=
		vrf_coordinator_v2.DeployVRFCoordinatorV2(
			neil, backend, linkAddress, common.Address{} /*blockHash store*/, fastGasFeed /* gasPrices */, linkEthFeed /* linkEth*/)
	require.NoError(t, err, "failed to deploy VRFCoordinator contract to simulated ethereum blockchain")
	// Deploy consumer
	consumerContractAddress, _, consumerContract, err :=
		vrf_consumer_v2.DeployVRFConsumerV2(
			carol, backend, coordinatorAddress, linkAddress)
	require.NoError(t, err, "failed to deploy VRFConsumer contract to simulated ethereum blockchain")
	_, err = linkContract.Transfer(sergey, consumerContractAddress, oneEth) // Actually, LINK
	require.NoError(t, err, "failed to send LINK to VRFConsumer contract on simulated ethereum blockchain")
	backend.Commit()
	return coordinatorV2Universe{
		rootContract:            coordinatorContract,
		rootContractAddress:     coordinatorAddress,
		linkContract:            linkContract,
		linkContractAddress:     linkAddress,
		consumerContract:        consumerContract,
		consumerContractAddress: consumerContractAddress,
		backend:                 backend,
		coordinatorABI:          &coordinatorABI,
		consumerABI:             &consumerABI,
		sergey:                  sergey,
		neil:                    neil,
		ned:                     ned,
		carol:                   carol,
	}
}

func TestIntegrationVRFV2(t *testing.T) {
	config, _, cleanupDB := heavyweight.FullTestORM(t, "vrf_v2", true)
	defer cleanupDB()
	key := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, key)
	t.Log(uni)

	app, cleanup := cltest.NewApplicationWithConfigAndKeyOnSimulatedBlockchain(t, config, uni.backend, key)
	defer cleanup()
	require.NoError(t, app.StartAndConnect())

	vrfkey, err := app.Store.VRFKeyStore.CreateKey(cltest.Password)
	require.NoError(t, err)
	unlocked, err := app.Store.VRFKeyStore.Unlock(cltest.Password)
	require.NoError(t, err)
	jid := uuid.NewV4()
	incomingConfs := 2
	s := testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{
		JobID:              jid.String(),
		Name:               "vrf-primary",
		CoordinatorAddress: uni.rootContractAddress.String(),
		Confirmations:      incomingConfs,
		PublicKey:          unlocked[0].String()}).Toml()
	jb, err := vrf.ValidatedVRFSpec(s)
	require.NoError(t, err)
	require.NoError(t, app.JobORM().CreateJob(context.Background(), &jb, jb.Pipeline))
	t.Log(vrfkey)

	p, err := vrfkey.Point()
	require.NoError(t, err)
	_, err = uni.rootContract.RegisterProvingKey(
		uni.neil, uni.neil.From, pair(secp256k1.Coordinates(p)))
	require.NoError(t, err)
	uni.backend.Commit()
	_, err = uni.consumerContract.TestCreateSubscriptionAndFund(uni.carol,
		big.NewInt(100))
	require.NoError(t, err)
	uni.backend.Commit()
	t.Log("Funded account")
	subId, err := uni.consumerContract.SubId(nil)
	require.NoError(t, err)
	t.Log("subscription ID", subId)
	_, err = uni.consumerContract.TestRequestRandomness(uni.carol, vrfkey.MustHash(), subId, 2, 10000,big.NewInt(1))
	require.NoError(t, err)
	// Mine the required number of blocks
	// So our request gets confirmed.
	for i := 0; i < incomingConfs; i++ {
		uni.backend.Commit()
	}
	reqID, err := uni.consumerContract.RequestId(nil)
	require.NoError(t, err)
	t.Log(reqID)
	callback, err := uni.rootContract.SCallbacks(nil, reqID)
	require.NoError(t, err)
	t.Log(callback)
	var runs []pipeline.Run
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		runs, err = app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		// It possible that we send the test request
		// before the job spawner has started the vrf services, which is fine
		// the lb will backfill the logs. However we need to
		// keep blocks coming in for the lb to send the backfilled logs.
		uni.backend.Commit()
		return len(runs) == 1
	}, 5*time.Second, 1*time.Second).Should(gomega.BeTrue())
	t.Log(runs[0])

	// Assert the request was fulfilled on-chain.
	var rf []*vrf_coordinator_v2.VRFCoordinatorV2RandomWordsFulfilled
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		rfIterator, err := uni.rootContract.FilterRandomWordsFulfilled(nil)
		require.NoError(t, err, "failed to logs")
		for rfIterator.Next() {
			rf = append(rf, rfIterator.Event)
		}
		return len(rf) == 1
	}, 5*time.Second, 500*time.Millisecond).Should(gomega.BeTrue())
	t.Log("randomness fulfilled req ID", rf[0].RequestId.String())
}
