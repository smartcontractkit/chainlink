package vrf

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/stretchr/testify/require"

	chainlink_eth "chainlink/core/eth"
	"chainlink/core/services/signatures/secp256k1"
	"chainlink/core/utils"

	"chainlink/core/services/vrf/generated/link_token_interface"
	"chainlink/core/services/vrf/generated/solidity_request_id"
	"chainlink/core/services/vrf/generated/solidity_vrf_consumer_interface"
	"chainlink/core/services/vrf/generated/solidity_vrf_coordinator_interface"
)

type coordinator struct {
	rootContract            *solidity_vrf_coordinator_interface.VRFCoordinator
	linkContract            *link_token_interface.LinkToken
	consumerContract        *solidity_vrf_consumer_interface.VRFConsumer
	requestIDBase           *solidity_request_id.VRFRequestIDBaseTestHelper
	consumerContractAddress common.Address
	backend                 *backends.SimulatedBackend
	sergey                  *bind.TransactOpts
	neil                    *bind.TransactOpts
	carol                   *bind.TransactOpts
}

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}
func newIdentity() *bind.TransactOpts {
	key, err := crypto.GenerateKey()
	panicErr(err)
	return bind.NewKeyedTransactor(key)
}

func deployCoordinator() coordinator {
	var (
		sergey = newIdentity()
		neil   = newIdentity()
		carol  = newIdentity()
	)
	oneEth := bi(1000000000000000000)
	genesisData := core.GenesisAlloc{
		sergey.From: {Balance: oneEth},
		neil.From:   {Balance: oneEth},
		carol.From:  {Balance: oneEth},
	}
	gasLimit := eth.DefaultConfig.Miner.GasCeil
	backend := backends.NewSimulatedBackend(genesisData, gasLimit)
	linkAddress, _, linkContract, err := link_token_interface.DeployLinkToken(
		sergey, backend)
	panicErr(err)
	coordinatorAddress, _, coordinatorContract, err :=
		solidity_vrf_coordinator_interface.DeployVRFCoordinator(
			neil, backend, linkAddress)
	panicErr(err)
	consumerContractAddress, _, consumerContract, err :=
		solidity_vrf_consumer_interface.DeployVRFConsumer(
			carol, backend, coordinatorAddress, linkAddress)
	panicErr(err)
	_, _, requestIDBase, err :=
		solidity_request_id.DeployVRFRequestIDBaseTestHelper(neil, backend)
	panicErr(err)
	_, err = linkContract.Transfer(sergey, consumerContractAddress, oneEth) // Actually, LINK
	panicErr(err)
	backend.Commit()
	return coordinator{
		rootContract:            coordinatorContract,
		linkContract:            linkContract,
		consumerContract:        consumerContract,
		requestIDBase:           requestIDBase,
		consumerContractAddress: consumerContractAddress,
		backend:                 backend,
		sergey:                  sergey,
		neil:                    neil,
		carol:                   carol,
	}
}

func TestRequestIDMatches(t *testing.T) {
	keyHash := common.HexToHash("0x01")
	seed := big.NewInt(1)
	baseContract := deployCoordinator().requestIDBase
	solidityRequestID, err := baseContract.MakeRequestId(nil, keyHash, seed)
	require.NoError(t, err)
	goRequestID := (&RandomnessRequestLog{KeyHash: keyHash, Seed: seed}).RequestID()
	require.Equal(t, common.BytesToHash(solidityRequestID[:]), goRequestID)
}

var (
	secretKey = one // never do this in production!
	publicKey = secp256k1Curve.Point().Mul(secp256k1.IntToScalar(secretKey), nil)
	seed      = one
)

func registerProvingKey(coordinator coordinator) (
	keyHash [32]byte, jobID [32]byte, fee *big.Int) {
	fee = seven
	copy(jobID[:], []byte("exactly 32 characters in length."))
	_, err := coordinator.rootContract.RegisterProvingKey(
		coordinator.neil, fee, pair(secp256k1.Coordinates(publicKey)), jobID)
	panicErr(err)
	coordinator.backend.Commit()
	copy(keyHash[:], utils.MustHash(string(secp256k1.LongMarshal(publicKey))).Bytes())
	return keyHash, jobID, fee
}

func TestRegisterProvingKey(t *testing.T) {
	coordinator := deployCoordinator()
	keyHash, jobID, fee := registerProvingKey(coordinator)
	log, err := coordinator.rootContract.FilterNewServiceAgreement(nil)
	require.NoError(t, err)
	logCount := 0
	for log.Next() {
		logCount += 1
		require.Equal(t, log.Event.KeyHash, keyHash)
		require.True(t, equal(fee, log.Event.Fee))
	}
	require.Equal(t, 1, logCount)
	serviceAgreement, err := coordinator.rootContract.ServiceAgreements(nil, keyHash)
	require.NoError(t, err)
	require.Equal(t, coordinator.neil.From, serviceAgreement.VRFOracle)
	require.Equal(t, jobID, serviceAgreement.JobID)
	require.True(t, equal(fee, serviceAgreement.Fee))
}

func requestRandomness(t *testing.T, coordinator coordinator,
	keyHash common.Hash, jobID common.Hash, fee, seed *big.Int) RandomnessRequestLog {
	_, err := coordinator.consumerContract.RequestRandomness(coordinator.carol,
		keyHash, fee, seed)
	require.NoError(t, err)
	coordinator.backend.Commit()
	log, err := coordinator.rootContract.FilterRandomnessRequest(nil, nil)
	require.NoError(t, err)
	logCount := 0
	for log.Next() {
		logCount += 1
	}
	require.Equal(t, 1, logCount)
	return rawRandomnessRequestLogToRandomnessRequestLog(RawRandomnessRequestLog(*log.Event))
}

func TestRandomnessRequestLog(t *testing.T) {
	coordinator := deployCoordinator()
	keyHash_, jobID_, fee := registerProvingKey(coordinator)
	keyHash := common.BytesToHash(keyHash_[:])
	jobID := common.BytesToHash(jobID_[:])
	log := requestRandomness(t, coordinator, keyHash, jobID, fee, seed)
	require.Equal(t, keyHash, log.KeyHash)
	actualSeed, err := coordinator.requestIDBase.MakeVRFInputSeed(nil, keyHash,
		one, coordinator.consumerContractAddress, zero)
	require.NoError(t, err)
	require.True(t, equal(actualSeed, log.Seed))
	require.Equal(t, jobID, log.JobID)
	require.Equal(t, coordinator.consumerContractAddress, log.Sender)
	require.True(t, equal(fee, (*big.Int)(log.Fee)))
	parsedLog, err := ParseRandomnessRequestLog(chainlink_eth.Log(log.Raw.Raw))
	require.NoError(t, err)
	require.True(t, parsedLog.Equal(log))
}

func fulfillRandomnessRequest(t *testing.T, coordinator coordinator,
	log RandomnessRequestLog) *Proof {
	proof, err := generateProofWithNonce(secretKey, log.Seed, one /* nonce */)
	require.NoError(t, err)
	proofBlob, err := proof.MarshalForSolidityVerifier()
	require.NoError(t, err)
	_, err = coordinator.rootContract.FulfillRandomnessRequest(
		coordinator.neil, proofBlob[:])
	require.NoError(t, err)
	coordinator.backend.Commit()
	return proof
}

func TestFulfillRandomness(t *testing.T) {
	coordinator := deployCoordinator()
	keyHash, jobID, fee := registerProvingKey(coordinator)
	log := requestRandomness(t, coordinator, keyHash, jobID, fee, seed)
	proof := fulfillRandomnessRequest(t, coordinator, log)
	output, err := coordinator.consumerContract.RandomnessOutput(nil)
	require.NoError(t, err)
	require.True(t, equal(proof.Output, output))
	requestID, err := coordinator.consumerContract.RequestId(nil)
	require.NoError(t, err)
	require.Equal(t, log.RequestID(), common.BytesToHash(requestID[:]))
	neilBalance, err := coordinator.rootContract.WithdrawableTokens(
		nil, coordinator.neil.From)
	require.NoError(t, err)
	require.True(t, equal(neilBalance, fee))
}

func TestWithdraw(t *testing.T) {
	coordinator := deployCoordinator()
	keyHash, jobID, fee := registerProvingKey(coordinator)
	log := requestRandomness(t, coordinator, keyHash, jobID, fee, seed)
	fulfillRandomnessRequest(t, coordinator, log)
	payment := big.NewInt(5)
	peteThePunter := common.HexToAddress("0xdeadfa11deadfa11deadfa11deadfa11deadfa11")
	_, err := coordinator.rootContract.Withdraw(coordinator.neil, peteThePunter, payment)
	coordinator.backend.Commit()
	require.NoError(t, err)
	peteBalance, err := coordinator.linkContract.BalanceOf(nil, peteThePunter)
	require.NoError(t, err)
	require.True(t, equal(payment, peteBalance))
	neilBalance, err := coordinator.rootContract.WithdrawableTokens(
		nil, coordinator.neil.From)
	require.NoError(t, err)
	require.True(t, equal(i().Sub(fee, payment), neilBalance))
	_, err = coordinator.rootContract.Withdraw(coordinator.neil, peteThePunter, fee)
	require.Error(t, err, "coordinator allowed overdraft")
}
