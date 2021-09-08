package vrf_test

import (
	"math/big"
	mrand "math/rand"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/solidity_vrf_v08_verifier_wrapper"
	proof2 "github.com/smartcontractkit/chainlink/core/services/vrf/proof"

	"github.com/ethereum/go-ethereum/eth/ethconfig"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/vrfkey"
	"github.com/smartcontractkit/chainlink/core/services/signatures/secp256k1"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Note these tests are identical to the ones in vrf_solidity_crosscheck_test.go,
// (with the exception of TestVRFV08_InvalidPointCoordinates which is a new check in v0.8)
// except we are testing against the v0.8 implementation of VRF.sol.
func deployVRFV08TestHelper(t *testing.T) *solidity_vrf_v08_verifier_wrapper.VRFV08TestHelper {
	key, err := crypto.GenerateKey()
	require.NoError(t, err, "failed to create root ethereum identity")
	auth := cltest.MustNewSimulatedBackendKeyedTransactor(t, key)
	genesisData := core.GenesisAlloc{auth.From: {Balance: assets.Ether(100)}}
	gasLimit := ethconfig.Defaults.Miner.GasCeil
	backend := cltest.NewSimulatedBackend(t, genesisData, gasLimit)
	_, _, verifier, err := solidity_vrf_v08_verifier_wrapper.DeployVRFV08TestHelper(auth, backend)
	require.NoError(t, err, "failed to deploy VRF contract to simulated blockchain")
	backend.Commit()
	return verifier
}

func TestVRFV08_InvalidPointCoordinates(t *testing.T) {
	verifier := deployVRFV08TestHelper(t)
	// A value outside [0, ..., FIELD_SIZE-1] should fail
	_, err := verifier.IsOnCurve(nil,
		[2]*big.Int{big.NewInt(10), secp256k1.FieldSize})
	require.Error(t, err)
	assert.Equal(t, err.Error(), "execution reverted: invalid y-ordinate")
	_, err = verifier.IsOnCurve(nil,
		[2]*big.Int{secp256k1.FieldSize, big.NewInt(10)})
	require.Error(t, err)
	assert.Equal(t, err.Error(), "execution reverted: invalid x-ordinate")
	// Values inside should succeed
	_, err = verifier.IsOnCurve(nil,
		[2]*big.Int{big.NewInt(10), big.NewInt(0).Sub(secp256k1.FieldSize, big.NewInt(1))})
	require.NoError(t, err)
	_, err = verifier.IsOnCurve(nil,
		[2]*big.Int{big.NewInt(0).Sub(secp256k1.FieldSize, big.NewInt(1)), big.NewInt(10)})
	require.NoError(t, err)
}

func TestVRFV08_CompareProjectiveECAddToVerifier(t *testing.T) {
	t.Parallel()
	r := mrand.New(mrand.NewSource(11))
	for j := 0; j < numSamples(); j++ {
		p := randomPoint(t, r)
		q := randomPoint(t, r)
		px, py := secp256k1.Coordinates(p)
		qx, qy := secp256k1.Coordinates(q)
		actualX, actualY, actualZ := vrfkey.ProjectiveECAdd(p, q)
		verifier := deployVRFV08TestHelper(t)
		expectedX, expectedY, expectedZ, err := verifier.ProjectiveECAdd(
			nil, px, py, qx, qy)
		require.NoError(t, err, "failed to compute secp256k1 sum in projective coords")
		assert.Equal(t, [3]*big.Int{expectedX, expectedY, expectedZ},
			[3]*big.Int{actualX, actualY, actualZ},
			"got different answers on-chain vs off-chain, for ProjectiveECAdd")
	}
}

func TestVRFV08_CompareBigModExpToVerifier(t *testing.T) {
	t.Parallel()
	r := mrand.New(mrand.NewSource(0))
	for j := 0; j < numSamples(); j++ {
		base := randomUint256(t, r)
		exponent := randomUint256(t, r)
		actual, err := deployVRFV08TestHelper(t).BigModExp(nil, base, exponent)
		require.NoError(t, err, "while computing bigmodexp on-chain")
		expected := big.NewInt(0).Exp(base, exponent, vrfkey.FieldSize)
		assert.Equal(t, expected, actual,
			"%x ** %x %% %x = %x ≠ %x from solidity calculation",
			base, exponent, vrfkey.FieldSize, expected, actual)
	}
}

func TestVRFV08_CompareSquareRoot(t *testing.T) {
	t.Parallel()
	r := mrand.New(mrand.NewSource(1))
	for j := 0; j < numSamples(); j++ {
		maybeSquare := randomUint256(t, r) // Might not be square; should get same result anyway
		squareRoot, err := deployVRFV08TestHelper(t).SquareRoot(nil, maybeSquare)
		require.NoError(t, err, "failed to compute square root on-chain")
		golangSquareRoot := vrfkey.SquareRoot(maybeSquare)
		assert.Equal(t, golangSquareRoot, squareRoot,
			"expected square root in GF(fieldSize) of %x to be %x, got %x on-chain",
			maybeSquare, golangSquareRoot, squareRoot)
		assert.True(t,
			(!vrfkey.IsSquare(maybeSquare)) || big.NewInt(1).Exp(squareRoot,
				big.NewInt(2), vrfkey.FieldSize).Cmp(maybeSquare) == 0,
			"maybeSquare is a square, but failed to calculate its square root!")
		assert.NotEqual(t, vrfkey.IsSquare(maybeSquare), vrfkey.IsSquare(
			big.NewInt(1).Sub(vrfkey.FieldSize, maybeSquare)),
			"negative of a non square should be square, and vice-versa, since -1 is not a square")
	}
}

func TestVRFV08_CompareYSquared(t *testing.T) {
	t.Parallel()
	r := mrand.New(mrand.NewSource(2))
	for i := 0; i < numSamples(); i++ {
		x := randomUint256(t, r)
		actual, err := deployVRFV08TestHelper(t).YSquared(nil, x)
		require.NoError(t, err, "failed to compute y² given x, on-chain")
		assert.Equal(t, vrfkey.YSquared(x), actual,
			"different answers for y², on-chain vs off-chain")
	}
}

func TestVRFV08_CompareFieldHash(t *testing.T) {
	t.Parallel()
	r := mrand.New(mrand.NewSource(3))
	msg := make([]byte, 32)
	for j := 0; j < numSamples(); j++ {
		_, err := r.Read(msg)
		require.NoError(t, err, "failed to randomize intended hash message")
		actual, err := deployVRFV08TestHelper(t).FieldHash(nil, msg)
		require.NoError(t, err, "failed to compute fieldHash on-chain")
		expected := vrfkey.FieldHash(msg)
		require.Equal(t, expected, actual,
			"fieldHash value on-chain differs from off-chain")
	}
}

// randomKey deterministically generates a secp256k1 key.
//
// Never use this if cryptographic security is required
//func randomKey(t *testing.T, r *mrand.Rand) *ecdsa.PrivateKey {
//	secretKey := vrfkey.FieldSize
//	for secretKey.Cmp(vrfkey.FieldSize) >= 0 { // Keep picking until secretKey < fieldSize
//		secretKey = randomUint256(t, r)
//	}
//	cKey := crypto.ToECDSAUnsafe(secretKey.Bytes())
//	return cKey
//}
//
// pair returns the inputs as a length-2 big.Int array. Useful for translating
// coordinates to the uint256[2]'s VRF.sol uses to represent secp256k1 points.
//func pair(x, y *big.Int) [2]*big.Int   { return [2]*big.Int{x, y} }
//func asPair(p kyber.Point) [2]*big.Int { return pair(secp256k1.Coordinates(p)) }

func TestVRFV08_CompareHashToCurve(t *testing.T) {
	t.Parallel()
	r := mrand.New(mrand.NewSource(4))
	for i := 0; i < numSamples(); i++ {
		input := randomUint256(t, r)
		cKey := randomKey(t, r)
		pubKeyCoords := pair(cKey.X, cKey.Y)
		actual, err := deployVRFV08TestHelper(t).HashToCurve(nil, pubKeyCoords, input)
		require.NoError(t, err, "failed to compute hashToCurve on-chain")
		pubKeyPoint := secp256k1.SetCoordinates(cKey.X, cKey.Y)
		expected, err := vrfkey.HashToCurve(pubKeyPoint, input, func(*big.Int) {})
		require.NoError(t, err, "failed to compute HashToCurve in golang")
		require.Equal(t, asPair(expected), actual,
			"on-chain and off-chain calculations of HashToCurve gave different secp256k1 points")
	}
}

// randomPoint deterministically simulates a uniform sample of secp256k1 points,
// given r's seed
//
// Never use this if cryptographic security is required
//func randomPoint(t *testing.T, r *mrand.Rand) kyber.Point {
//	p, err := vrfkey.HashToCurve(vrfkey.Generator, randomUint256(t, r), func(*big.Int) {})
//	require.NoError(t, err,
//		"failed to hash random value to secp256k1 while generating random point")
//	if r.Int63n(2) == 1 { // Uniform sample of ±p
//		p.Neg(p)
//	}
//	return p
//}
//
//// randomPointWithPair returns a random secp256k1, both as a kyber.Point and as
//// a pair of *big.Int's. Useful for translating between the types needed by the
//// golang contract wrappers.
//func randomPointWithPair(t *testing.T, r *mrand.Rand) (kyber.Point, [2]*big.Int) {
//	p := randomPoint(t, r)
//	return p, asPair(p)
//}

// randomScalar deterministically simulates a uniform sample of secp256k1
// scalars, given r's seed
//
// Never use this if cryptographic security is required
//func randomScalar(t *testing.T, r *mrand.Rand) kyber.Scalar {
//	s := randomUint256(t, r)
//	for s.Cmp(secp256k1.GroupOrder) >= 0 {
//		s = randomUint256(t, r)
//	}
//	return secp256k1.IntToScalar(s)
//}

func TestVRFV08_CheckSolidityPointAddition(t *testing.T) {
	t.Parallel()
	r := mrand.New(mrand.NewSource(5))
	for j := 0; j < numSamples(); j++ {
		p1 := randomPoint(t, r)
		p2 := randomPoint(t, r)
		p1x, p1y := secp256k1.Coordinates(p1)
		p2x, p2y := secp256k1.Coordinates(p2)
		psx, psy, psz, err := deployVRFV08TestHelper(t).ProjectiveECAdd(
			nil, p1x, p1y, p2x, p2y)
		require.NoError(t, err, "failed to compute ProjectiveECAdd, on-chain")
		apx, apy, apz := vrfkey.ProjectiveECAdd(p1, p2)
		require.Equal(t, []*big.Int{apx, apy, apz}, []*big.Int{psx, psy, psz},
			"got different values on-chain and off-chain for ProjectiveECAdd")
		zInv := big.NewInt(1).ModInverse(psz, vrfkey.FieldSize)
		require.Equal(t, big.NewInt(1).Mod(big.NewInt(1).Mul(psz, zInv),
			vrfkey.FieldSize), big.NewInt(1),
			"failed to calculate correct inverse of z ordinate")
		actualSum, err := deployVRFV08TestHelper(t).AffineECAdd(
			nil, pair(p1x, p1y), pair(p2x, p2y), zInv)
		require.NoError(t, err,
			"failed to deploy VRF contract to simulated blockchain")
		assert.Equal(t, asPair((&secp256k1.Secp256k1{}).Point().Add(p1, p2)),
			actualSum,
			"got different answers, on-chain vs off-chain, for secp256k1 sum in affine coordinates")
	}
}

func TestVRFV08_CheckSolidityECMulVerify(t *testing.T) {
	t.Parallel()
	r := mrand.New(mrand.NewSource(6))
	for j := 0; j < numSamples(); j++ {
		p := randomPoint(t, r)
		pxy := pair(secp256k1.Coordinates(p))
		s := randomScalar(t, r)
		product := asPair((&secp256k1.Secp256k1{}).Point().Mul(s, p))
		actual, err := deployVRFV08TestHelper(t).EcmulVerify(nil, pxy, secp256k1.ToInt(s),
			product)
		require.NoError(t, err, "failed to check on-chain that s*p=product")
		assert.True(t, actual,
			"EcmulVerify rejected a valid secp256k1 scalar product relation")
		shouldReject, err := deployVRFV08TestHelper(t).EcmulVerify(nil, pxy,
			big.NewInt(0).Add(secp256k1.ToInt(s), big.NewInt(1)), product)
		require.NoError(t, err, "failed to check on-chain that (s+1)*p≠product")
		assert.False(t, shouldReject,
			"failed to reject a false secp256k1 scalar product relation")
	}
}

func TestVRFV08_CheckSolidityVerifyLinearCombinationWithGenerator(t *testing.T) {
	t.Parallel()
	r := mrand.New(mrand.NewSource(7))
	for j := 0; j < numSamples(); j++ {
		c := randomScalar(t, r)
		s := randomScalar(t, r)
		p := randomPoint(t, r)
		expectedPoint := (&secp256k1.Secp256k1{}).Point().Add(
			(&secp256k1.Secp256k1{}).Point().Mul(c, p),
			(&secp256k1.Secp256k1{}).Point().Mul(s, vrfkey.Generator)) // cp+sg
		expectedAddress := secp256k1.EthereumAddress(expectedPoint)
		pPair := asPair(p)
		actual, err := deployVRFV08TestHelper(t).VerifyLinearCombinationWithGenerator(nil,
			secp256k1.ToInt(c), pPair, secp256k1.ToInt(s), expectedAddress)
		require.NoError(t, err,
			"failed to check on-chain that secp256k1 linear relationship holds")
		assert.True(t, actual,
			"VerifyLinearCombinationWithGenerator rejected a valid secp256k1 linear relationship")
		shouldReject, err := deployVRFV08TestHelper(t).VerifyLinearCombinationWithGenerator(nil,
			big.NewInt(0).Add(secp256k1.ToInt(c), big.NewInt(1)), pPair,
			secp256k1.ToInt(s), expectedAddress)
		require.NoError(t, err,
			"failed to check on-chain that address((c+1)*p+s*g)≠expectedAddress")
		assert.False(t, shouldReject,
			"VerifyLinearCombinationWithGenerator accepted an invalid secp256k1 linear relationship!")
	}
}

func TestVRFV08_CheckSolidityLinearComination(t *testing.T) {
	t.Parallel()
	r := mrand.New(mrand.NewSource(8))
	for j := 0; j < numSamples(); j++ {
		c := randomScalar(t, r)
		cNum := secp256k1.ToInt(c)
		p1, p1Pair := randomPointWithPair(t, r)
		s := randomScalar(t, r)
		sNum := secp256k1.ToInt(s)
		p2, p2Pair := randomPointWithPair(t, r)
		cp1 := (&secp256k1.Secp256k1{}).Point().Mul(c, p1)
		cp1Pair := asPair(cp1)
		sp2 := (&secp256k1.Secp256k1{}).Point().Mul(s, p2)
		sp2Pair := asPair(sp2)
		expected := asPair((&secp256k1.Secp256k1{}).Point().Add(cp1, sp2))
		_, _, z := vrfkey.ProjectiveECAdd(cp1, sp2)
		zInv := big.NewInt(0).ModInverse(z, vrfkey.FieldSize)
		actual, err := deployVRFV08TestHelper(t).LinearCombination(nil, cNum, p1Pair,
			cp1Pair, sNum, p2Pair, sp2Pair, zInv)
		require.NoError(t, err, "failed to compute c*p1+s*p2, on-chain")
		assert.Equal(t, expected, actual,
			"on-chain computation of c*p1+s*p2 gave wrong answer")
		_, err = deployVRFV08TestHelper(t).LinearCombination(nil, big.NewInt(0).Add(
			cNum, big.NewInt(1)), p1Pair, cp1Pair, sNum, p2Pair, sp2Pair, zInv)
		assert.Error(t, err,
			"on-chain LinearCombination accepted a bad product relation! ((c+1)*p1)")
		assert.Contains(t, err.Error(), "First mul check failed",
			"revert message wrong.")
		_, err = deployVRFV08TestHelper(t).LinearCombination(nil, cNum, p1Pair,
			cp1Pair, big.NewInt(0).Add(sNum, big.NewInt(1)), p2Pair, sp2Pair, zInv)
		assert.Error(t, err,
			"on-chain LinearCombination accepted a bad product relation! ((s+1)*p2)")
		assert.Contains(t, err.Error(), "Second mul check failed",
			"revert message wrong.")
	}
}

func TestVRFV08_CompareSolidityScalarFromCurvePoints(t *testing.T) {
	t.Parallel()
	r := mrand.New(mrand.NewSource(9))
	for j := 0; j < numSamples(); j++ {
		hash, hashPair := randomPointWithPair(t, r)
		pk, pkPair := randomPointWithPair(t, r)
		gamma, gammaPair := randomPointWithPair(t, r)
		var uWitness [20]byte
		require.NoError(t, utils.JustError(r.Read(uWitness[:])),
			"failed to randomize uWitness")
		v, vPair := randomPointWithPair(t, r)
		expected := vrfkey.ScalarFromCurvePoints(hash, pk, gamma, uWitness, v)
		actual, err := deployVRFV08TestHelper(t).ScalarFromCurvePoints(nil, hashPair, pkPair,
			gammaPair, uWitness, vPair)
		require.NoError(t, err, "on-chain ScalarFromCurvePoints calculation failed")
		assert.Equal(t, expected, actual,
			"on-chain ScalarFromCurvePoints output does not match off-chain output!")
	}
}

func TestVRFV08_MarshalProof(t *testing.T) {
	t.Parallel()
	r := mrand.New(mrand.NewSource(10))
	for j := 0; j < numSamples(); j++ {
		sk := randomScalar(t, r)
		skNum := secp256k1.ToInt(sk)
		pk := vrfkey.MustNewV2XXXTestingOnly(skNum)
		nonce := randomScalar(t, r)
		randomSeed := randomUint256(t, r)
		proof, err := pk.GenerateProofWithNonce(randomSeed, secp256k1.ToInt(nonce))
		require.NoError(t, err, "failed to generate VRF proof!")
		require.NoError(t, err, "failed to marshal VRF proof for on-chain verification")
		seed, err := proof2.BigToSeed(randomSeed)
		require.NoError(t, err)
		// Don't care about the request commitment for this test.
		solProof, _, err := proof2.GenerateProofResponseFromProofV2(proof, proof2.PreSeedDataV2{
			PreSeed: seed,
		})
		require.NoError(t, err)
		response, err := deployVRFV08TestHelper(t).RandomValueFromVRFProof(nil, solidity_vrf_v08_verifier_wrapper.VRFProof{
			Pk:            solProof.Pk,
			Gamma:         solProof.Gamma,
			C:             solProof.C,
			S:             solProof.S,
			Seed:          solProof.Seed,
			UWitness:      solProof.UWitness,
			CGammaWitness: solProof.CGammaWitness,
			SHashWitness:  solProof.SHashWitness,
			ZInv:          solProof.ZInv,
		}, randomSeed)
		require.NoError(t, err, "failed on-chain to verify VRF proof / get its output")
		require.True(t, response.Cmp(proof.Output) == 0,
			"on-chain VRF output differs from off-chain!")
	}
}
