package vrf_test

import (
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/signatures/secp256k1"
	"github.com/smartcontractkit/chainlink/core/services/vrf"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/solidity_vrf_consumer_interface"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/solidity_vrf_coordinator_interface"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/solidity_vrf_request_id"
)

// coordinator represents the universe in which a randomness request occurs and
// is fulfilled.
type coordinator struct {
	// Golang wrappers ofr solidity contracts
	rootContract            *solidity_vrf_coordinator_interface.VRFCoordinator
	linkContract            *link_token_interface.LinkToken
	consumerContract        *solidity_vrf_consumer_interface.VRFConsumer
	requestIDBase           *solidity_vrf_request_id.VRFRequestIDBaseTestHelper
	rootContractAddress     common.Address
	consumerContractAddress common.Address
	// Abstraction representation of the ethereum blockchain
	backend        *backends.SimulatedBackend
	coordinatorABI *abi.ABI
	consumerABI    *abi.ABI
	// Cast of participants
	sergey *bind.TransactOpts // Owns all the LINK initially
	neil   *bind.TransactOpts // Node operator running VRF service
	carol  *bind.TransactOpts // Author of consuming contract which requests randomness
}

var seed = big.NewInt(1)

// newIdentity returns a go-ethereum abstraction of an ethereum account for
// interacting with contract golang wrappers
func newIdentity(t *testing.T) *bind.TransactOpts {
	key, err := crypto.GenerateKey()
	require.NoError(t, err, "failed to generate ethereum identity")
	return bind.NewKeyedTransactor(key)
}

// deployCoordinator sets up all identities and contracts associated with
// testing the solidity VRF contracts involved in randomness request workflow
func deployCoordinator(t *testing.T) coordinator {
	var (
		sergey = newIdentity(t)
		neil   = newIdentity(t)
		carol  = newIdentity(t)
	)
	oneEth := big.NewInt(1000000000000000000)
	genesisData := core.GenesisAlloc{
		sergey.From: {Balance: oneEth},
		neil.From:   {Balance: oneEth},
		carol.From:  {Balance: oneEth},
	}
	gasLimit := eth.DefaultConfig.Miner.GasCeil
	consumerABI, err := abi.JSON(strings.NewReader(
		solidity_vrf_consumer_interface.VRFConsumerABI))
	require.NoError(t, err)
	coordinatorABI, err := abi.JSON(strings.NewReader(
		solidity_vrf_coordinator_interface.VRFCoordinatorABI))
	require.NoError(t, err)
	backend := backends.NewSimulatedBackend(genesisData, gasLimit)
	linkAddress, _, linkContract, err := link_token_interface.DeployLinkToken(
		sergey, backend)
	require.NoError(t, err, "failed to deploy link contract to simulated ethereum blockchain")
	coordinatorAddress, _, coordinatorContract, err :=
		solidity_vrf_coordinator_interface.DeployVRFCoordinator(
			neil, backend, linkAddress, common.Address{} /* BlockHashStore address */)
	require.NoError(t, err, "failed to deploy VRFCoordinator contract to simulated ethereum blockchain")
	consumerContractAddress, _, consumerContract, err :=
		solidity_vrf_consumer_interface.DeployVRFConsumer(
			carol, backend, coordinatorAddress, linkAddress)
	require.NoError(t, err, "failed to deploy VRFConsumer contract to simulated ethereum blockchain")
	_, _, requestIDBase, err :=
		solidity_vrf_request_id.DeployVRFRequestIDBaseTestHelper(neil, backend)
	require.NoError(t, err, "failed to deploy VRFRequestIDBaseTestHelper contract to simulated ethereum blockchain")
	_, err = linkContract.Transfer(sergey, consumerContractAddress, oneEth) // Actually, LINK
	require.NoError(t, err, "failed to send LINK to VRFConsumer contract on simulated ethereum blockchain")
	backend.Commit()
	return coordinator{
		rootContract:            coordinatorContract,
		rootContractAddress:     coordinatorAddress,
		linkContract:            linkContract,
		consumerContract:        consumerContract,
		requestIDBase:           requestIDBase,
		consumerContractAddress: consumerContractAddress,
		backend:                 backend,
		coordinatorABI:          &coordinatorABI,
		consumerABI:             &consumerABI,
		sergey:                  sergey,
		neil:                    neil,
		carol:                   carol,
	}
}

func TestRequestIDMatches(t *testing.T) {
	keyHash := common.HexToHash("0x01")
	baseContract := deployCoordinator(t).requestIDBase
	solidityRequestID, err := baseContract.MakeRequestId(nil, keyHash, seed)
	require.NoError(t, err, "failed to calculate VRF requestID on simulated ethereum blockchain")
	goRequestLog := &models.RandomnessRequestLog{KeyHash: keyHash, Seed: seed}
	assert.Equal(t, common.Hash(solidityRequestID), goRequestLog.RequestID(),
		"solidity VRF requestID differs from golang requestID!")
}

var (
	rawSecretKey = big.NewInt(1) // never do this in production!
	publicKey    = (&secp256k1.Secp256k1{}).Point().Mul(secp256k1.IntToScalar(
		rawSecretKey), nil)
	rawSeed = big.NewInt(2)
	vrfFee  = big.NewInt(7)
)

// registerProvingKey registers keyHash to neil in the VRFCoordinator universe
// represented by coordinator, with the given jobID and fee.
func registerProvingKey(t *testing.T, coordinator coordinator) (
	keyHash [32]byte, jobID [32]byte, fee *big.Int) {
	copy(jobID[:], []byte("exactly 32 characters in length."))
	_, err := coordinator.rootContract.RegisterProvingKey(
		coordinator.neil, vrfFee, coordinator.neil.From, pair(secp256k1.Coordinates(publicKey)), jobID)
	require.NoError(t, err, "failed to register VRF proving key on VRFCoordinator contract")
	coordinator.backend.Commit()
	keyHash = utils.MustHash(string(secp256k1.LongMarshal(publicKey)))
	return keyHash, jobID, vrfFee
}

func TestRegisterProvingKey(t *testing.T) {
	coord := deployCoordinator(t)
	keyHash, jobID, fee := registerProvingKey(t, coord)
	log, err := coord.rootContract.FilterNewServiceAgreement(nil)
	require.NoError(t, err, "failed to subscribe to NewServiceAgreement logs on simulated ethereum blockchain")
	logCount := 0
	for log.Next() {
		logCount += 1
		assert.Equal(t, log.Event.KeyHash, keyHash, "VRFCoordinator logged a different keyHash than was registered")
		assert.True(t, fee.Cmp(log.Event.Fee) == 0, "VRFCoordinator logged a different fee than was registered")
	}
	require.Equal(t, 1, logCount, "unexpected NewServiceAgreement log generated by key VRF key registration")
	serviceAgreement, err := coord.rootContract.ServiceAgreements(nil, keyHash)
	require.NoError(t, err, "failed to retrieve previously registered VRF service agreement from VRFCoordinator")
	assert.Equal(t, coord.neil.From, serviceAgreement.VRFOracle,
		"VRFCoordinator registered wrong provider, on service agreement!")
	assert.Equal(t, jobID, serviceAgreement.JobID,
		"VRFCoordinator registered wrong jobID, on service agreement!")
	assert.True(t, fee.Cmp(serviceAgreement.Fee) == 0,
		"VRFCoordinator registered wrong fee, on service agreement!")
}

// requestRandomness sends a randomness request via Carol's consuming contract,
// in the VRFCoordinator universe represented by coordinator, specifying the
// given keyHash and seed, and paying the given fee. It returns the log emitted
// from the VRFCoordinator in response to the request
func requestRandomness(t *testing.T, coordinator coordinator,
	keyHash common.Hash, fee, seed *big.Int) *models.RandomnessRequestLog {
	_, err := coordinator.consumerContract.RequestRandomness(coordinator.carol,
		keyHash, fee, seed)
	require.NoError(t, err, "problem during initial VRF randomness request")
	coordinator.backend.Commit()
	log, err := coordinator.rootContract.FilterRandomnessRequest(nil, nil)
	require.NoError(t, err, "failed to subscribe to RandomnessRequest logs")
	logCount := 0
	for log.Next() {
		logCount += 1
	}
	require.Equal(t, 1, logCount, "unexpected log generated by randomness request to VRFCoordinator")
	return models.RawRandomnessRequestLogToRandomnessRequestLog(
		(*models.RawRandomnessRequestLog)(log.Event))
}

func TestRandomnessRequestLog(t *testing.T) {
	coord := deployCoordinator(t)
	keyHash_, jobID_, fee := registerProvingKey(t, coord)
	keyHash := common.BytesToHash(keyHash_[:])
	jobID := common.BytesToHash(jobID_[:])
	log := requestRandomness(t, coord, keyHash, fee, rawSeed)
	assert.Equal(t, keyHash, log.KeyHash, "VRFCoordinator logged wrong KeyHash for randomness request")
	nonce := big.NewInt(0)
	actualSeed, err := coord.requestIDBase.MakeVRFInputSeed(nil, keyHash,
		rawSeed, coord.consumerContractAddress, nonce)
	require.NoError(t, err, "failure while using VRFCoordinator to calculate actual VRF input seed")
	assert.True(t, actualSeed.Cmp(log.Seed) == 0,
		"VRFCoordinator logged wrong actual input seed from randomness request")
	golangSeed := utils.MustHash(string(append(append(append(
		keyHash[:],
		common.BigToHash(rawSeed).Bytes()...),
		coord.consumerContractAddress.Hash().Bytes()...),
		common.BigToHash(nonce).Bytes()...)))
	assert.Equal(t, golangSeed, common.BigToHash((log.Seed)), "VRFCoordinator logged different actual input seed than expected by golang code!")
	assert.Equal(t, jobID, log.JobID, "VRFCoordinator logged different JobID from randomness request!")
	assert.Equal(t, coord.consumerContractAddress, log.Sender, "VRFCoordinator logged different requester address from randomness request!")
	assert.True(t, fee.Cmp((*big.Int)(log.Fee)) == 0, "VRFCoordinator logged different fee from randomness request!")
	parsedLog, err := models.ParseRandomnessRequestLog(cltest.ChainlinkEthLogFromGethLog(log.Raw.Raw))
	assert.NoError(t, err, "could not parse randomness request log generated by VRFCoordinator")
	assert.True(t, parsedLog.Equal(*log), "got a different randomness request log by parsing the raw data than reported by simulated backend")
}

// fulfillRandomnessRequest is neil fulfilling randomness requested by log.
func fulfillRandomnessRequest(t *testing.T, coordinator coordinator,
	log models.RandomnessRequestLog) vrf.Proof {
	preseed, err := vrf.BigToSeed(log.Seed)
	require.NoError(t, err, "preseed %x out of range", preseed)
	s := vrf.PreSeedData{
		PreSeed:   preseed,
		BlockHash: log.Raw.Raw.BlockHash,
		BlockNum:  log.Raw.Raw.BlockNumber,
	}
	proofBlob, err := vrf.GenerateProofResponseWithNonce(rawSecretKey, s, big.NewInt(1) /* nonce */)
	require.NoError(t, err, "could not generate VRF proof!")
	// Seems to be a bug in the simulated backend: without this extra Commit, the
	// EVM seems to think it's still on the block in which the request was made,
	// which means that the relevant blockhash is unavailable.
	coordinator.backend.Commit()
	// This is simulating a node response, so set the gas limit as chainlink does
	var neil bind.TransactOpts = *coordinator.neil
	neil.GasLimit = orm.NewConfig().EthGasLimitDefault()
	_, err = coordinator.rootContract.FulfillRandomnessRequest(&neil, proofBlob[:])
	require.NoError(t, err, "failed to fulfill randomness request!")
	coordinator.backend.Commit()
	goProofResponse, err := vrf.UnmarshalProofResponse(proofBlob)
	require.NoError(t, err,
		"could not rehydrate proof from blob sent to fulfillRandomnessRequest")
	proof, err := goProofResponse.ActualProof(s)
	require.NoError(t, err, "could not construct actual proof from proof response")
	return proof
}

func TestFulfillRandomness(t *testing.T) {
	coordinator := deployCoordinator(t)
	keyHash, _, fee := registerProvingKey(t, coordinator)
	randomnessRequestLog := requestRandomness(t, coordinator, keyHash, fee, seed)
	proof := fulfillRandomnessRequest(t, coordinator, *randomnessRequestLog)
	output, err := coordinator.consumerContract.RandomnessOutput(nil)
	require.NoError(t, err, "failed to get VRF output from consuming contract, "+
		"after randomness request was fulfilled")
	assert.True(t, proof.Output.Cmp(output) == 0, "VRF output from randomness "+
		"request fulfillment was different than provided! Expected %d, got %d. "+
		"This can happen if you update the VRFCoordinator wrapper without a "+
		"corresponding update to the VRFConsumer", proof.Output, output)
	requestID, err := coordinator.consumerContract.RequestId(nil)
	require.NoError(t, err, "failed to get requestId from VRFConsumer")
	assert.Equal(t, randomnessRequestLog.RequestID(), common.Hash(requestID),
		"VRFConsumer has different request ID than logged from randomness request!")
	neilBalance, err := coordinator.rootContract.WithdrawableTokens(
		nil, coordinator.neil.From)
	require.NoError(t, err, "failed to get neil's token balance, after he "+
		"successfully fulfilled a randomness request")
	assert.True(t, neilBalance.Cmp(fee) == 0, "neil's balance on VRFCoordinator "+
		"was not paid his fee, despite successful fulfillment of randomness request!")
}

func TestWithdraw(t *testing.T) {
	coordinator := deployCoordinator(t)
	keyHash, _, fee := registerProvingKey(t, coordinator)
	log := requestRandomness(t, coordinator, keyHash, fee, rawSeed)
	fulfillRandomnessRequest(t, coordinator, *log)
	payment := big.NewInt(4)
	peteThePunter := common.HexToAddress("0xdeadfa11deadfa11deadfa11deadfa11deadfa11")
	_, err := coordinator.rootContract.Withdraw(coordinator.neil, peteThePunter, payment)
	require.NoError(t, err, "failed to withdraw LINK from neil's balance")
	coordinator.backend.Commit()
	peteBalance, err := coordinator.linkContract.BalanceOf(nil, peteThePunter)
	require.NoError(t, err, "failed to get balance of payee on LINK contract, after payment")
	assert.True(t, payment.Cmp(peteBalance) == 0,
		"LINK balance is wrong, following payment")
	neilBalance, err := coordinator.rootContract.WithdrawableTokens(
		nil, coordinator.neil.From)
	require.NoError(t, err, "failed to get neil's balance on VRFCoordinator")
	assert.True(t, big.NewInt(0).Sub(fee, payment).Cmp(neilBalance) == 0,
		"neil's VRFCoordinator balance is wrong, after he's made a withdrawal!")
	_, err = coordinator.rootContract.Withdraw(coordinator.neil, peteThePunter, fee)
	assert.Error(t, err, "VRFcoordinator allowed overdraft")
}
