package vrf

import (
	"context"
	mrand "math/rand"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/stretchr/testify/require"

	"chainlink/core/services/signatures/secp256k1"
)

func TestMeasureRandomValueFromVRFProofGasCost(t *testing.T) {
	r := mrand.New(mrand.NewSource(10))
	sk := randomScalar(t, r)
	skNum := secp256k1.ToInt(sk)
	nonce := randomScalar(t, r)
	seed := randomUint256(t, r)
	proof, err := generateProofWithNonce(skNum, seed, secp256k1.ToInt(nonce))
	require.NoError(t, err)
	mproof, err := proof.MarshalForSolidityVerifier()
	require.NoError(t, err)
	contract, owner := deployVRFContract(t)
	rawData, err := contract.abi.Pack("randomValueFromVRFProof_", mproof[:])
	require.NoError(t, err)
	callMsg := ethereum.CallMsg{From: owner, To: &contract.address, Data: rawData}
	estimate, err := contract.backend.EstimateGas(context.TODO(), callMsg)
	require.NoError(t, err)
	require.Less(t, estimate, uint64(100000))
}
