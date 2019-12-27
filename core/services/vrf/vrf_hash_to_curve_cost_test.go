package vrf

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"strings"
	"testing"

	"chainlink/core/services/signatures/secp256k1"
	"chainlink/core/services/vrf/generated/solidity_verifier_wrapper"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/stretchr/testify/require"
)

type contract struct {
	contract *bind.BoundContract
	address  common.Address
	abi      *abi.ABI
	backend  *backends.SimulatedBackend
}

// deployVRFContract operates like
// solidity_verifier_wrapper.DeployVRFTestHelper, except that it exposes the
// actual contract, which is useful for gas measurements.
func deployVRFContract(t *testing.T) (contract, common.Address) {
	x, y := secp256k1.Coordinates(Generator)
	key := ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{Curve: crypto.S256(), X: x, Y: y},
		D:         big.NewInt(1),
	}
	auth := bind.NewKeyedTransactor(&key)
	genesisData := core.GenesisAlloc{auth.From: {Balance: bi(1000000000)}}
	gasLimit := eth.DefaultConfig.Miner.GasCeil
	backend := backends.NewSimulatedBackend(genesisData, gasLimit)
	parsed, err := abi.JSON(strings.NewReader(
		solidity_verifier_wrapper.VRFTestHelperABI))
	require.NoError(t, err)
	address, _, vRFContract, err := bind.DeployContract(auth, parsed,
		common.FromHex(solidity_verifier_wrapper.VRFTestHelperBin), backend)
	require.NoError(t, err)
	backend.Commit()
	return contract{vRFContract, address, &parsed, backend}, crypto.PubkeyToAddress(
		key.PublicKey)
}

func measureHashToCurveGasCost(t *testing.T, contract contract,
	owner common.Address, input int64) (gasCost, numOrdinates uint64) {
	rawData, err := contract.abi.Pack("hashToCurve_", pair(secp256k1.Coordinates(Generator)),
		big.NewInt(input))
	require.NoError(t, err)
	callMsg := ethereum.CallMsg{From: owner, To: &contract.address, Data: rawData}
	estimate, err := contract.backend.EstimateGas(context.TODO(), callMsg)
	require.NoError(t, err)
	_, err = HashToCurve(Generator, big.NewInt(input),
		func(*big.Int) { numOrdinates += 1 })
	require.NoError(t, err)
	return estimate, numOrdinates
}

var baseCost uint64 = 24000
var marginalCost uint64 = 15555

func HashToCurveGasCostBound(numOrdinates uint64) uint64 {
	return baseCost + marginalCost*numOrdinates
}

func TestMeasureHashToCurveGasCost(t *testing.T) {
	contract, owner := deployVRFContract(t)
	numSamples := int64(10) // Holds for first 1,000 samples, but set to 10 for speed.
	for i := int64(0); i < numSamples; i += 1 {
		gasCost, numOrdinates := measureHashToCurveGasCost(t, contract, owner, i)
		require.Less(t, gasCost, HashToCurveGasCostBound(numOrdinates))
	}
	require.Less(t, HashToCurveGasCostBound(128), uint64(2.016e6))
}
