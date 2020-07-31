package vrf_test

import (
	mrand "math/rand"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/services/signatures/secp256k1"
	"github.com/smartcontractkit/chainlink/core/services/vrf"
)

func TestMeasureRandomValueFromVRFProofGasCost(t *testing.T) {
	r := mrand.New(mrand.NewSource(10))
	sk := randomScalar(t, r)
	skNum := secp256k1.ToInt(sk)
	nonce := randomScalar(t, r)
	seed := randomUint256(t, r)
	proof, err := vrf.GenerateProofWithNonce(skNum, seed, secp256k1.ToInt(nonce))
	require.NoError(t, err, "failed to generate VRF proof")
	mproof, err := proof.MarshalForSolidityVerifier()
	require.NoError(t, err, "failed to marshal VRF proof for on-chain verification")
	contract, _ := deployVRFContract(t)

	estimate := estimateGas(t, contract.backend, common.Address{},
		contract.address, contract.abi, "randomValueFromVRFProof_", mproof[:])

	require.NoError(t, err, "failed to estimate gas cost for VRF verification")
	require.Less(t, estimate, uint64(100000))
}
