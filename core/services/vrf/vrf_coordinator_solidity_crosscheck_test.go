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

	"chainlink/core/assets"
	chainlink_eth "chainlink/core/eth"
	"chainlink/core/services/signatures/secp256k1"
	"chainlink/core/utils"

	"chainlink/core/services/vrf/generated/link_token_interface"
	"chainlink/core/services/vrf/generated/solidity_vrf_consumer_interface"
	"chainlink/core/services/vrf/generated/solidity_vrf_coordinator_interface"
)

type coordinator struct {
	rootContract            *solidity_vrf_coordinator_interface.VRFCoordinator
	linkContract            *link_token_interface.LinkToken
	consumerContract        *solidity_vrf_consumer_interface.VRFConsumer
	consumerContractAddress common.Address
	backend                 *backends.SimulatedBackend
	sergey                  *bind.TransactOpts
	neil                    *bind.TransactOpts
	carol                   *bind.TransactOpts
}

func newIdentity(t *testing.T) *bind.TransactOpts {
	key, err := crypto.GenerateKey()
	require.NoError(t, err)
	return bind.NewKeyedTransactor(key)
}

func deployCoordinator(t *testing.T) coordinator {
	var (
		sergey = newIdentity(t)
		neil   = newIdentity(t)
		carol  = newIdentity(t)
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
	require.NoError(t, err)
	coordinatorAddress, _, coordinatorContract, err :=
		solidity_vrf_coordinator_interface.DeployVRFCoordinator(
			neil, backend, linkAddress)
	require.NoError(t, err)
	consumerContractAddress, _, consumerContract, err :=
		solidity_vrf_consumer_interface.DeployVRFConsumer(
			carol, backend, coordinatorAddress, linkAddress)
	require.NoError(t, err)
	_, err = linkContract.Transfer(sergey, consumerContractAddress, oneEth) // Actually, LINK
	require.NoError(t, err)
	backend.Commit()
	return coordinator{
		rootContract:            coordinatorContract,
		linkContract:            linkContract,
		consumerContract:        consumerContract,
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
	coordinator := deployCoordinator(t).rootContract
	solidityRequestID, err := coordinator.MakeRequestId(nil, keyHash, seed)
	require.NoError(t, err)
	goRequestID := (&RandomnessRequestLog{KeyHash: keyHash, Seed: seed}).RequestID()
	require.Equal(t, common.BytesToHash(solidityRequestID[:]), goRequestID)
}

func registerProvingKey(t *testing.T, coordinator coordinator) (
	keyHash [32]byte, jobID [32]byte, fee *big.Int) {
	sk := secp256k1.IntToScalar(one)
	pk := secp256k1Curve.Point().Mul(sk, nil)
	fee = seven
	copy(jobID[:], []byte("exactly 32 characters in length."))
	_, err := coordinator.rootContract.RegisterProvingKey(
		coordinator.neil, fee, pair(secp256k1.Coordinates(pk)), jobID)
	require.NoError(t, err)
	coordinator.backend.Commit()
	copy(keyHash[:], utils.MustHash(string(secp256k1.LongMarshal(pk))).Bytes())
	return keyHash, jobID, fee
}

func TestRegisterProvingKey(t *testing.T) {
	coordinator := deployCoordinator(t)
	keyHash, jobID, fee := registerProvingKey(t, coordinator)
	log, err := coordinator.rootContract.FilterNewServiceAgreement(nil)
	require.NoError(t, err)
	logCount := 0
	for log.Next() {
		logCount += 1
		require.Equal(t, log.Event.KeyHash, keyHash)
		require.True(t, fee.Cmp(log.Event.Fee) == 0)
	}
	require.Equal(t, 1, logCount)
	serviceAgreement, err := coordinator.rootContract.ServiceAgreements(nil, keyHash)
	require.NoError(t, err)
	require.Equal(t, coordinator.neil.From, serviceAgreement.VRFOracle)
	require.Equal(t, jobID, serviceAgreement.JobID)
	require.True(t, fee.Cmp(serviceAgreement.Fee) == 0)
}

func TestRandomnessRequestLog(t *testing.T) {
	coordinator := deployCoordinator(t)
	keyHash, jobID, fee := registerProvingKey(t, coordinator)
	coordinator.backend.Commit()
	seed := one
	_, err := coordinator.consumerContract.RequestRandomness(coordinator.carol,
		keyHash, fee, seed)
	coordinator.backend.Commit()
	require.NoError(t, err)
	log, err := coordinator.rootContract.FilterRandomnessRequest(nil, nil)
	require.NoError(t, err)
	logCount := 0
	for log.Next() {
		logCount += 1
		require.Equal(t, keyHash, log.Event.KeyHash)
		require.True(t, seed.Cmp(log.Event.Seed) == 0)
		require.Equal(t, jobID, log.Event.JobID)
		require.Equal(t, coordinator.consumerContractAddress, log.Event.Sender)
		require.Equal(t, fee, log.Event.Fee)
	}
	require.Equal(t, 1, logCount)
	parsedLog, err := ParseRandomnessRequestLog(chainlink_eth.Log(log.Event.Raw))
	require.NoError(t, err)
	require.True(t, parsedLog.Equal(RandomnessRequestLog{keyHash, seed, jobID,
		coordinator.consumerContractAddress, (*assets.Link)(fee),
		RawRandomnessRequestLog(*log.Event)}))
}
