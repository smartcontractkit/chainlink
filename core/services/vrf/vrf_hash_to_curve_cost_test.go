package vrf_test

import (
	"crypto/ecdsa"
	"math/big"
	"strings"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/testutils"

	"github.com/ethereum/go-ethereum/eth/ethconfig"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/solidity_vrf_verifier_wrapper"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/vrfkey"
	"github.com/smartcontractkit/chainlink/core/services/signatures/secp256k1"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type contract struct {
	contract *bind.BoundContract
	address  common.Address
	abi      *abi.ABI
	backend  *backends.SimulatedBackend
}

// deployVRFContract returns a deployed VRF contract, with some extra attributes
// which are useful for gas measurements.
func deployVRFContract(t *testing.T) (contract, common.Address) {
	x, y := secp256k1.Coordinates(vrfkey.Generator)
	key := ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{Curve: crypto.S256(), X: x, Y: y},
		D:         big.NewInt(1),
	}
	auth := cltest.MustNewSimulatedBackendKeyedTransactor(t, &key)
	genesisData := core.GenesisAlloc{auth.From: {Balance: assets.Ether(100)}}
	gasLimit := ethconfig.Defaults.Miner.GasCeil
	backend := cltest.NewSimulatedBackend(t, genesisData, gasLimit)
	parsed, err := abi.JSON(strings.NewReader(
		solidity_vrf_verifier_wrapper.VRFTestHelperABI))
	require.NoError(t, err, "could not parse VRF ABI")
	address, _, vRFContract, err := bind.DeployContract(auth, parsed,
		common.FromHex(solidity_vrf_verifier_wrapper.VRFTestHelperBin), backend)
	require.NoError(t, err, "failed to deploy VRF contract to simulated blockchain")
	backend.Commit()
	return contract{vRFContract, address, &parsed, backend}, crypto.PubkeyToAddress(
		key.PublicKey)
}

// estimateGas returns the estimated gas cost of running the given method on the
// contract at address to, on the given backend, with the given args, and given
// that the transaction is sent from the from address.
func estimateGas(t *testing.T, backend *backends.SimulatedBackend,
	from, to common.Address, abi *abi.ABI, method string, args ...interface{},
) uint64 {
	rawData, err := abi.Pack(method, args...)
	require.NoError(t, err, "failed to construct raw %s transaction with args %s",
		method, args)
	callMsg := ethereum.CallMsg{From: from, To: &to, Data: rawData}
	estimate, err := backend.EstimateGas(testutils.Context(t), callMsg)
	require.NoError(t, err, "failed to estimate gas from %s call with args %s",
		method, args)
	return estimate
}

func measureHashToCurveGasCost(t *testing.T, contract contract,
	owner common.Address, input int64) (gasCost, numOrdinates uint64) {
	estimate := estimateGas(t, contract.backend, owner, contract.address,
		contract.abi, "hashToCurve_", pair(secp256k1.Coordinates(vrfkey.Generator)),
		big.NewInt(input))

	_, err := vrfkey.HashToCurve(vrfkey.Generator, big.NewInt(input),
		func(*big.Int) { numOrdinates += 1 })
	require.NoError(t, err, "corresponding golang HashToCurve calculation failed")
	return estimate, numOrdinates
}

var baseCost uint64 = 25000
var marginalCost uint64 = 15555

func HashToCurveGasCostBound(numOrdinates uint64) uint64 {
	return baseCost + marginalCost*numOrdinates
}

func TestMeasureHashToCurveGasCost(t *testing.T) {
	contract, owner := deployVRFContract(t)
	numSamples := int64(numSamples())
	for i := int64(0); i < numSamples; i += 1 {
		gasCost, numOrdinates := measureHashToCurveGasCost(t, contract, owner, i)
		assert.Less(t, gasCost, HashToCurveGasCostBound(numOrdinates),
			"on-chain hashToCurve gas cost exceeded estimate function")
	}
	require.Less(t, HashToCurveGasCostBound(128), uint64(2.017e6),
		"estimate for on-chain hashToCurve gas cost with 128 iterations is greater "+
			"than stated in the VRF.sol documentation")
}
