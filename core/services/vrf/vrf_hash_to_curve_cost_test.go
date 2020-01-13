package vrf

import (
	"context"
	"crypto/ecdsa"
	"math"
	"math/big"
	"strings"
	"testing"

	"chainlink/core/services/signatures/secp256k1"

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
// actual contract.
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
		VRFTestHelperABI))
	require.NoError(t, err)
	address, _, vRFContract, err := bind.DeployContract(auth, parsed,
		common.FromHex(VRFTestHelperBin), backend)
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
	return estimate, numOrdinates
}

func TestMeasureHashToCurveGasCost(t *testing.T) {
	contract, owner := deployVRFContract(t)
	const maxOrds = 20
	var totalGasCosts [maxOrds]float64
	var totalGasCostSquareds [maxOrds]float64
	var ordCounts [maxOrds]float64
	numSamples := int64(10)
	for i := int64(0); i < numSamples; i += 1 {
		gasCost, numOrdinates := measureHashToCurveGasCost(t, contract, owner, i)
		totalGasCosts[numOrdinates] += float64(gasCost)
		totalGasCostSquareds[numOrdinates] += float64(gasCost) * float64(gasCost)
		ordCounts[numOrdinates] += 1
	}
	var meanGasCosts [maxOrds]float64
	var std [maxOrds]float64
	for numOrdinates := uint64(1); totalGasCosts[numOrdinates] > 0; numOrdinates += 1 {
		meanGasCosts[numOrdinates] = totalGasCosts[numOrdinates] / ordCounts[numOrdinates]
		meanSquared := meanGasCosts[numOrdinates] * meanGasCosts[numOrdinates]
		std[numOrdinates] = math.Sqrt(
			totalGasCostSquareds[numOrdinates]/ordCounts[numOrdinates] - meanSquared)
		require.Less(t, std[numOrdinates], 100.0)
	}
	var differences [maxOrds]float64
	totalDifferences := float64(0)
	totalDifferencesSquared := float64(0)
	diffCount := float64(0)
	for numOrdinates := uint64(2); totalGasCosts[numOrdinates] > 0; numOrdinates += 1 {
		differences[numOrdinates-2] = meanGasCosts[numOrdinates] - meanGasCosts[numOrdinates-1]
		totalDifferences += differences[numOrdinates-2]
		totalDifferencesSquared += differences[numOrdinates-2] * differences[numOrdinates-2]
		diffCount += 1
	}
	meanDiff := totalDifferences / diffCount
	stdDiff := math.Sqrt(totalDifferencesSquared/diffCount - meanDiff*meanDiff)
	require.Less(t, meanDiff, 15600.0)
	require.Less(t, stdDiff, 30.0)
}
