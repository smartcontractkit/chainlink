package vrf_test

import (
	"math/big"
	"strings"
	"testing"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/blockhash_store"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/solidity_vrf_consumer_interface"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/solidity_vrf_consumer_interface_v08"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/solidity_vrf_coordinator_interface"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/solidity_vrf_request_id"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/solidity_vrf_request_id_v08"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	proof2 "github.com/smartcontractkit/chainlink/core/services/vrf/proof"

	"github.com/ethereum/go-ethereum/eth/ethconfig"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	evmconfig "github.com/smartcontractkit/chainlink/core/chains/evm/config"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/vrfkey"
	"github.com/smartcontractkit/chainlink/core/services/signatures/secp256k1"
	"github.com/smartcontractkit/chainlink/core/services/vrf"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
)

// coordinatorUniverse represents the universe in which a randomness request occurs and
// is fulfilled.
type coordinatorUniverse struct {
	// Golang wrappers ofr solidity contracts
	rootContract               *solidity_vrf_coordinator_interface.VRFCoordinator
	linkContract               *link_token_interface.LinkToken
	bhsContract                *blockhash_store.BlockhashStore
	consumerContract           *solidity_vrf_consumer_interface.VRFConsumer
	requestIDBase              *solidity_vrf_request_id.VRFRequestIDBaseTestHelper
	consumerContractV08        *solidity_vrf_consumer_interface_v08.VRFConsumer
	requestIDBaseV08           *solidity_vrf_request_id_v08.VRFRequestIDBaseTestHelper
	rootContractAddress        common.Address
	consumerContractAddress    common.Address
	consumerContractAddressV08 common.Address
	linkContractAddress        common.Address
	bhsContractAddress         common.Address

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

var oneEth = big.NewInt(1000000000000000000) // 1e18 wei

func newVRFCoordinatorUniverseWithV08Consumer(t *testing.T, key ethkey.KeyV2) coordinatorUniverse {
	cu := newVRFCoordinatorUniverse(t, key)
	consumerContractAddress, _, consumerContract, err :=
		solidity_vrf_consumer_interface_v08.DeployVRFConsumer(
			cu.carol, cu.backend, cu.rootContractAddress, cu.linkContractAddress)
	require.NoError(t, err, "failed to deploy v08 VRFConsumer contract to simulated ethereum blockchain")
	_, _, requestIDBase, err :=
		solidity_vrf_request_id_v08.DeployVRFRequestIDBaseTestHelper(cu.neil, cu.backend)
	require.NoError(t, err, "failed to deploy v08 VRFRequestIDBaseTestHelper contract to simulated ethereum blockchain")
	cu.consumerContractAddressV08 = consumerContractAddress
	cu.requestIDBaseV08 = requestIDBase
	cu.consumerContractV08 = consumerContract
	_, err = cu.linkContract.Transfer(cu.sergey, consumerContractAddress, oneEth) // Actually, LINK
	require.NoError(t, err, "failed to send LINK to VRFConsumer contract on simulated ethereum blockchain")
	cu.backend.Commit()
	return cu
}

// newVRFCoordinatorUniverse sets up all identities and contracts associated with
// testing the solidity VRF contracts involved in randomness request workflow
func newVRFCoordinatorUniverse(t *testing.T, keys ...ethkey.KeyV2) coordinatorUniverse {
	var oracleTransactors []*bind.TransactOpts
	for _, key := range keys {
		oracleTransactor, err := bind.NewKeyedTransactorWithChainID(key.ToEcdsaPrivKey(), testutils.SimulatedChainID)
		require.NoError(t, err)
		oracleTransactors = append(oracleTransactors, oracleTransactor)
	}

	var (
		sergey = testutils.MustNewSimTransactor(t)
		neil   = testutils.MustNewSimTransactor(t)
		ned    = testutils.MustNewSimTransactor(t)
		carol  = testutils.MustNewSimTransactor(t)
	)
	genesisData := core.GenesisAlloc{
		sergey.From: {Balance: assets.Ether(1000).ToInt()},
		neil.From:   {Balance: assets.Ether(1000).ToInt()},
		ned.From:    {Balance: assets.Ether(1000).ToInt()},
		carol.From:  {Balance: assets.Ether(1000).ToInt()},
	}

	for _, t := range oracleTransactors {
		genesisData[t.From] = core.GenesisAccount{Balance: assets.Ether(1000).ToInt()}
	}

	gasLimit := uint32(ethconfig.Defaults.Miner.GasCeil)
	consumerABI, err := abi.JSON(strings.NewReader(
		solidity_vrf_consumer_interface.VRFConsumerABI))
	require.NoError(t, err)
	coordinatorABI, err := abi.JSON(strings.NewReader(
		solidity_vrf_coordinator_interface.VRFCoordinatorABI))
	require.NoError(t, err)
	backend := cltest.NewSimulatedBackend(t, genesisData, gasLimit)
	linkAddress, _, linkContract, err := link_token_interface.DeployLinkToken(
		sergey, backend)
	require.NoError(t, err, "failed to deploy link contract to simulated ethereum blockchain")
	bhsAddress, _, bhsContract, err := blockhash_store.DeployBlockhashStore(neil, backend)
	require.NoError(t, err, "failed to deploy BlockhashStore contract to simulated ethereum blockchain")
	coordinatorAddress, _, coordinatorContract, err :=
		solidity_vrf_coordinator_interface.DeployVRFCoordinator(
			neil, backend, linkAddress, bhsAddress)
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
	return coordinatorUniverse{
		rootContract:            coordinatorContract,
		rootContractAddress:     coordinatorAddress,
		linkContract:            linkContract,
		linkContractAddress:     linkAddress,
		bhsContract:             bhsContract,
		bhsContractAddress:      bhsAddress,
		consumerContract:        consumerContract,
		requestIDBase:           requestIDBase,
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

func TestRequestIDMatches(t *testing.T) {
	keyHash := common.HexToHash("0x01")
	key := cltest.MustGenerateRandomKey(t)
	baseContract := newVRFCoordinatorUniverse(t, key).requestIDBase
	var seed = big.NewInt(1)
	solidityRequestID, err := baseContract.MakeRequestId(nil, keyHash, seed)
	require.NoError(t, err, "failed to calculate VRF requestID on simulated ethereum blockchain")
	goRequestLog := &vrf.RandomnessRequestLog{KeyHash: keyHash, Seed: seed}
	assert.Equal(t, common.Hash(solidityRequestID), goRequestLog.ComputedRequestID(),
		"solidity VRF requestID differs from golang requestID!")
}

var (
	rawSecretKey = big.NewInt(1) // never do this in production!
	secretKey    = vrfkey.MustNewV2XXXTestingOnly(rawSecretKey)
	publicKey    = (&secp256k1.Secp256k1{}).Point().Mul(secp256k1.IntToScalar(
		rawSecretKey), nil)
	hardcodedSeed = big.NewInt(0)
	vrfFee        = big.NewInt(7)
)

// registerProvingKey registers keyHash to neil in the VRFCoordinator universe
// represented by coordinator, with the given jobID and fee.
func registerProvingKey(t *testing.T, coordinator coordinatorUniverse) (
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
	key := cltest.MustGenerateRandomKey(t)
	coord := newVRFCoordinatorUniverse(t, key)
	keyHash, jobID, fee := registerProvingKey(t, coord)
	log, err := coord.rootContract.FilterNewServiceAgreement(nil)
	require.NoError(t, err, "failed to subscribe to NewServiceAgreement logs on simulated ethereum blockchain")
	logCount := 0
	for log.Next() {
		logCount++
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

func TestFailToRegisterProvingKeyFromANonOwnerAddress(t *testing.T) {
	key := cltest.MustGenerateRandomKey(t)
	coordinator := newVRFCoordinatorUniverse(t, key)

	var jobID [32]byte
	copy(jobID[:], []byte("exactly 32 characters in length."))
	_, err := coordinator.rootContract.RegisterProvingKey(
		coordinator.ned, vrfFee, coordinator.neil.From, pair(secp256k1.Coordinates(publicKey)), jobID)

	require.Error(t, err, "expected an error")
	require.Contains(t, err.Error(), "Ownable: caller is not the owner")
}

// requestRandomness sends a randomness request via Carol's consuming contract,
// in the VRFCoordinator universe represented by coordinator, specifying the
// given keyHash and seed, and paying the given fee. It returns the log emitted
// from the VRFCoordinator in response to the request
func requestRandomness(t *testing.T, coordinator coordinatorUniverse,
	keyHash common.Hash, fee *big.Int) *vrf.RandomnessRequestLog {
	_, err := coordinator.consumerContract.TestRequestRandomness(coordinator.carol,
		keyHash, fee)
	require.NoError(t, err, "problem during initial VRF randomness request")
	coordinator.backend.Commit()
	log, err := coordinator.rootContract.FilterRandomnessRequest(nil, nil)
	require.NoError(t, err, "failed to subscribe to RandomnessRequest logs")
	logCount := 0
	for log.Next() {
		logCount++
	}
	require.Equal(t, 1, logCount, "unexpected log generated by randomness request to VRFCoordinator")
	return vrf.RawRandomnessRequestLogToRandomnessRequestLog(
		(*vrf.RawRandomnessRequestLog)(log.Event))
}

func requestRandomnessV08(t *testing.T, coordinator coordinatorUniverse,
	keyHash common.Hash, fee *big.Int) *vrf.RandomnessRequestLog {
	_, err := coordinator.consumerContractV08.DoRequestRandomness(coordinator.carol,
		keyHash, fee)
	require.NoError(t, err, "problem during initial VRF randomness request")
	coordinator.backend.Commit()
	log, err := coordinator.rootContract.FilterRandomnessRequest(nil, nil)
	require.NoError(t, err, "failed to subscribe to RandomnessRequest logs")
	logCount := 0
	for log.Next() {
		if log.Event.Sender == coordinator.consumerContractAddressV08 {
			logCount++
		}
	}
	require.Equal(t, 1, logCount, "unexpected log generated by randomness request to VRFCoordinator")
	return vrf.RawRandomnessRequestLogToRandomnessRequestLog(
		(*vrf.RawRandomnessRequestLog)(log.Event))
}

func TestRandomnessRequestLog(t *testing.T) {
	key := cltest.MustGenerateRandomKey(t)
	coord := newVRFCoordinatorUniverseWithV08Consumer(t, key)
	keyHash_, jobID_, fee := registerProvingKey(t, coord)
	keyHash := common.BytesToHash(keyHash_[:])
	jobID := common.BytesToHash(jobID_[:])
	var tt = []struct {
		rr func(t *testing.T, coordinator coordinatorUniverse,
			keyHash common.Hash, fee *big.Int) *vrf.RandomnessRequestLog
		ms              func() (*big.Int, error)
		consumerAddress common.Address
	}{
		{
			rr: requestRandomness,
			ms: func() (*big.Int, error) {
				return coord.requestIDBase.MakeVRFInputSeed(nil, keyHash, hardcodedSeed, coord.consumerContractAddress, big.NewInt(0))
			},
			consumerAddress: coord.consumerContractAddress,
		},
		{
			rr: requestRandomnessV08,
			ms: func() (*big.Int, error) {
				return coord.requestIDBaseV08.MakeVRFInputSeed(nil, keyHash, hardcodedSeed, coord.consumerContractAddressV08, big.NewInt(0))
			},
			consumerAddress: coord.consumerContractAddressV08,
		},
	}
	for _, tc := range tt {
		log := tc.rr(t, coord, keyHash, fee)
		assert.Equal(t, keyHash, log.KeyHash, "VRFCoordinator logged wrong KeyHash for randomness request")
		nonce := big.NewInt(0)
		actualSeed, err := tc.ms()
		require.NoError(t, err, "failure while using VRFCoordinator to calculate actual VRF input seed")
		assert.True(t, actualSeed.Cmp(log.Seed) == 0,
			"VRFCoordinator logged wrong actual input seed from randomness request")
		golangSeed := utils.MustHash(string(append(append(append(
			keyHash[:],
			common.BigToHash(hardcodedSeed).Bytes()...),
			tc.consumerAddress.Hash().Bytes()...),
			common.BigToHash(nonce).Bytes()...)))
		assert.Equal(t, golangSeed, common.BigToHash((log.Seed)), "VRFCoordinator logged different actual input seed than expected by golang code!")
		assert.Equal(t, jobID, log.JobID, "VRFCoordinator logged different JobID from randomness request!")
		assert.Equal(t, tc.consumerAddress, log.Sender, "VRFCoordinator logged different requester address from randomness request!")
		assert.True(t, fee.Cmp((*big.Int)(log.Fee)) == 0, "VRFCoordinator logged different fee from randomness request!")
		parsedLog, err := vrf.ParseRandomnessRequestLog(log.Raw.Raw)
		assert.NoError(t, err, "could not parse randomness request log generated by VRFCoordinator")
		assert.True(t, parsedLog.Equal(*log), "got a different randomness request log by parsing the raw data than reported by simulated backend")
	}
}

// fulfillRandomnessRequest is neil fulfilling randomness requested by log.
func fulfillRandomnessRequest(t *testing.T, coordinator coordinatorUniverse, log vrf.RandomnessRequestLog) vrfkey.Proof {
	preSeed, err := proof2.BigToSeed(log.Seed)
	require.NoError(t, err, "pre-seed %x out of range", preSeed)
	s := proof2.PreSeedData{
		PreSeed:   preSeed,
		BlockHash: log.Raw.Raw.BlockHash,
		BlockNum:  log.Raw.Raw.BlockNumber,
	}
	seed := proof2.FinalSeed(s)
	proof, err := secretKey.GenerateProofWithNonce(seed, big.NewInt(1) /* nonce */)
	require.NoError(t, err)
	proofBlob, err := GenerateProofResponseFromProof(proof, s)
	require.NoError(t, err, "could not generate VRF proof!")
	// Seems to be a bug in the simulated backend: without this extra Commit, the
	// EVM seems to think it's still on the block in which the request was made,
	// which means that the relevant blockhash is unavailable.
	coordinator.backend.Commit()
	// This is simulating a node response, so set the gas limit as chainlink does
	var neil bind.TransactOpts = *coordinator.neil
	neil.GasLimit = uint64(evmconfig.DefaultGasLimit)
	_, err = coordinator.rootContract.FulfillRandomnessRequest(&neil, proofBlob[:])
	require.NoError(t, err, "failed to fulfill randomness request!")
	coordinator.backend.Commit()
	return proof
}

func TestFulfillRandomness(t *testing.T) {
	key := cltest.MustGenerateRandomKey(t)
	coordinator := newVRFCoordinatorUniverse(t, key)
	keyHash, _, fee := registerProvingKey(t, coordinator)
	randomnessRequestLog := requestRandomness(t, coordinator, keyHash, fee)
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
	assert.Equal(t, randomnessRequestLog.RequestID, common.Hash(requestID),
		"VRFConsumer has different request ID than logged from randomness request!")
	neilBalance, err := coordinator.rootContract.WithdrawableTokens(
		nil, coordinator.neil.From)
	require.NoError(t, err, "failed to get neil's token balance, after he "+
		"successfully fulfilled a randomness request")
	assert.True(t, neilBalance.Cmp(fee) == 0, "neil's balance on VRFCoordinator "+
		"was not paid his fee, despite successful fulfillment of randomness request!")
}

func TestWithdraw(t *testing.T) {
	key := cltest.MustGenerateRandomKey(t)
	coordinator := newVRFCoordinatorUniverse(t, key)
	keyHash, _, fee := registerProvingKey(t, coordinator)
	log := requestRandomness(t, coordinator, keyHash, fee)
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
