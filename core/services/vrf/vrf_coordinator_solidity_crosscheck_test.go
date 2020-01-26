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

	"chainlink/core/services/vrf/generated/link_token_interface"
	"chainlink/core/services/vrf/generated/solidity_vrf_coordinator_interface"
)

type coordinator struct {
	rootContract *solidity_vrf_coordinator_interface.VRFCoordinator
	linkContract *link_token_interface.LinkToken
	sergey       *bind.TransactOpts
}

func deployCoordinator(t *testing.T) coordinator {
	key, _ := crypto.GenerateKey()
	sergey := bind.NewKeyedTransactor(key)
	genesisData := core.GenesisAlloc{sergey.From: {Balance: bi(1000000000)}}
	gasLimit := eth.DefaultConfig.Miner.GasCeil
	backend := backends.NewSimulatedBackend(genesisData, gasLimit)
	linkAddress, _, linkContract, err := link_token_interface.DeployLinkToken(
		sergey, backend)
	require.NoError(t, err)
	_, _, coordinatorContract, err :=
		solidity_vrf_coordinator_interface.DeployVRFCoordinator(
			sergey, backend, linkAddress)
	require.NoError(t, err)
	backend.Commit()
	return coordinator{coordinatorContract, linkContract, sergey}
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
