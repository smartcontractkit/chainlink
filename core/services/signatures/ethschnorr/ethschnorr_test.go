package ethschnorr

// This code is largely based on go.dedis.ch/kyber/sign/schnorr_test from
// EPFL's DEDIS

import (
	crand "crypto/rand"
	"fmt"
	"math/big"
	mrand "math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	"go.dedis.ch/kyber"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/signatures/secp256k1"
)

var numSignatures = 5

var randomStream = cltest.NewStream(&testing.T{}, 0)

var printTests = false

func printTest(t *testing.T, msg *big.Int, private kyber.Scalar,
	public kyber.Point, signature Signature) {
	privateBytes, err := private.MarshalBinary()
	require.Nil(t, err)
	pX, pY := secp256k1.Coordinates(public)
	fmt.Printf("  ['%064x',\n   '%064x',\n   '%064x',\n   '%064x',\n   "+
		"'%064x',\n   '%040x'],\n",
		msg, privateBytes, pX, pY, signature.Signature,
		signature.CommitmentPublicAddress)
}

func TestShortSchnorr_SignAndVerify(t *testing.T) {
	if printTests {
		fmt.Printf("tests = [\n")
	}
	for i := 0; i < numSignatures; i++ {
		rand := mrand.New(mrand.NewSource(0))
		msg, err := crand.Int(rand, maxUint256)
		require.NoError(t, err)
		kp := secp256k1.Generate(randomStream)
		sig, err := Sign(kp.Private, msg)
		require.NoError(t, err, "failed to sign message")
		require.NoError(t, Verify(kp.Public, msg, sig),
			"failed to validate own signature")
		require.Error(t, Verify(kp.Public, u256Cardinality, sig),
			"failed to abort on too large a message")
		require.Error(t, Verify(kp.Public, big.NewInt(0).Neg(big.NewInt(1)), sig),
			"failed to abort on negative message")
		if printTests {
			printTest(t, msg, kp.Private, kp.Public, sig)
		}
		wrongMsg := big.NewInt(0).Add(msg, big.NewInt(1))
		require.Error(t, Verify(kp.Public, wrongMsg, sig),
			"failed to reject signature with bad message")
		wrongPublic := secp256k1Group.Point().Add(kp.Public, kp.Public)
		require.Error(t, Verify(wrongPublic, msg, sig),
			"failed to reject signature with bad public key")
		wrongSignature := &signature{
			CommitmentPublicAddress: sig.CommitmentPublicAddress,
			Signature:               big.NewInt(0).Add(sig.Signature, one),
		}
		require.Error(t, Verify(kp.Public, msg, wrongSignature),
			"failed to reject bad signature")
		badPublicCommitmentAddress := &signature{Signature: sig.Signature}
		copy(badPublicCommitmentAddress.CommitmentPublicAddress[:],
			sig.CommitmentPublicAddress[:])
		badPublicCommitmentAddress.CommitmentPublicAddress[0] ^= 1 // Corrupt it
		require.Error(t, Verify(kp.Public, msg, badPublicCommitmentAddress),
			"failed to reject signature with bad public commitment")
	}
	if printTests {
		fmt.Println("]")
	}
}
