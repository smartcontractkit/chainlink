package vrf

import (
	"crypto/ecdsa"
	"math/big"
	mrand "math/rand"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"go.dedis.ch/kyber/v3"

	"chainlink/core/services/vrf/generated/solidity_verifier_wrapper"

	"chainlink/core/services/signatures/secp256k1"
)

// Cross-checks of golang implementation details vs corresponding solidity
// details.
//
// It's worth automatically checking these implementation details because they
// can help to quickly locate any disparity between the solidity and golang
// implementations.

var verifier *solidity_verifier_wrapper.VRFTestHelper

// init initializes the wrapper of the EVM verifier contract.
//
// NOTE: If persistent state is ever added to the verifier contract, a separate
// verifier must be initialized for each test.
//
// NB: For changes to the VRF solidity code to be reflected here, "go generate"
// must be run in core/services/vrf.
func init() {
	key, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}
	auth := bind.NewKeyedTransactor(key)
	genesisData := core.GenesisAlloc{auth.From: {Balance: bi(1000000000)}}
	gasLimit := eth.DefaultConfig.Miner.GasCeil
	backend := backends.NewSimulatedBackend(genesisData, gasLimit)
	// err must already be declared, or next line will shadow global verifier
	_, _, verifier, err = solidity_verifier_wrapper.DeployVRFTestHelper(auth, backend)
	if err != nil {
		panic(errors.Wrapf(err, "while initializing EVM contract wrapper"))
	}
	backend.Commit()
}

// randomUint256 deterministically simulates a uniform sample of uint256's,
// given r's seed
//
// Never use this if cryptographic security is required
func randomUint256(t *testing.T, r *mrand.Rand) *big.Int {
	b := make([]byte, 32)
	_, err := r.Read(b)
	require.NoError(t, err)
	return i().SetBytes(b)
}

func numSamples() int {
	if testing.Short() {
		return 10
	}
	return 1000
}

func TestVRF_CompareProjectiveECAddToVerifier(t *testing.T) {
	t.Parallel()
	r := mrand.New(mrand.NewSource(11))
	numSamps := numSamples()
	for j := 0; j < numSamps; j++ {
		p := randomPoint(t, r)
		q := randomPoint(t, r)
		px, py := secp256k1.Coordinates(p)
		qx, qy := secp256k1.Coordinates(q)
		actualX, actualY, actualZ := ProjectiveECAdd(p, q)
		expectedX, expectedY, expectedZ, err := verifier.ProjectiveECAdd(
			nil, px, py, qx, qy)
		require.NoError(t, err)
		require.Equal(t, expectedX, actualX)
		require.Equal(t, expectedY, actualY)
		require.Equal(t, expectedZ, actualZ)
	}
}

func TestVRF_CompareBigModExpToVerifier(t *testing.T) {
	t.Parallel()
	r := mrand.New(mrand.NewSource(0))
	numSamps := numSamples()
	for j := 0; j < numSamps; j++ {
		base := randomUint256(t, r)
		exponent := randomUint256(t, r)
		actual, err := verifier.BigModExp(nil, base, exponent)
		require.NoError(t, err, "while computing bigmodexp")
		expected := i().Exp(base, exponent, fieldSize)
		require.Equal(t, expected, actual, "%d ** %d %% %d = %d ≠ %d",
			base, exponent, fieldSize, expected, actual)
	}
}

func TestVRF_CompareSquareRoot(t *testing.T) {
	t.Parallel()
	r := mrand.New(mrand.NewSource(1))
	numSamps := numSamples()
	for i := 0; i < numSamps; i++ {
		square := randomUint256(t, r)
		squareRoot, err := verifier.SquareRoot(nil, square)
		require.NoError(t, err)
		require.Equal(t, SquareRoot(square), squareRoot)
	}
}

func TestVRF_CompareYSquared(t *testing.T) {
	t.Parallel()
	r := mrand.New(mrand.NewSource(2))
	numSamps := numSamples()
	for i := 0; i < numSamps; i++ {
		x := randomUint256(t, r)
		actual, err := verifier.YSquared(nil, x)
		require.NoError(t, err)
		require.Equal(t, YSquared(x), actual)
	}
}

func TestVRF_CompareFieldHash(t *testing.T) {
	t.Parallel()
	r := mrand.New(mrand.NewSource(3))
	msg := make([]byte, 32)
	numSamps := numSamples()
	for j := 0; j < numSamps; j++ {
		r.Read(msg)
		msgAsNum := i().SetBytes(msg)
		actual, err := verifier.FieldHash(nil, msgAsNum)
		require.NoError(t, err)
		expected := fieldHash(msg)
		require.Equal(t, expected, actual)
	}
}

func randomKey(t *testing.T, r *mrand.Rand) *ecdsa.PrivateKey {
	secretKey := fieldSize
	for secretKey.Cmp(fieldSize) != -1 { // Keep picking until secretKey < P
		secretKey = randomUint256(t, r)
	}
	cKey := crypto.ToECDSAUnsafe(secretKey.Bytes())
	return cKey
}

func pair(x, y *big.Int) [2]*big.Int { return [2]*big.Int{x, y} }

func TestVRF_CompareHashToCurve(t *testing.T) {
	t.Parallel()
	r := mrand.New(mrand.NewSource(4))
	numSamps := numSamples()
	for i := 0; i < numSamps; i++ {
		input := randomUint256(t, r)
		cKey := randomKey(t, r)
		pubKeyCoords := pair(cKey.X, cKey.Y)
		actual, err := verifier.HashToCurve(nil, pubKeyCoords, input)
		require.NoError(t, err)
		pubKeyPoint := secp256k1.SetCoordinates(cKey.X, cKey.Y)
		expected, err := HashToCurve(pubKeyPoint, input, func(*big.Int) {})
		require.NoError(t, err)
		require.Equal(t, expected, secp256k1.SetCoordinates(actual[0], actual[1]))
	}
}

// randomPoint deterministically simulates a uniform sample of secp256k1 points,
// given r's seed
//
// Never use this if cryptographic security is required
func randomPoint(t *testing.T, r *mrand.Rand) kyber.Point {
	p, err := HashToCurve(Generator, randomUint256(t, r), func(*big.Int) {})
	require.NoError(t, err)
	if r.Int63n(2) == 1 { // Uniform sample of ±p
		p.Neg(p)
	}
	return p
}

// randomScalar deterministically simulates a uniform of secp256k1
// scalars, given r's seed
//
// Never use this if cryptographic security is required
func randomScalar(t *testing.T, r *mrand.Rand) kyber.Scalar {
	s := randomUint256(t, r)
	for s.Cmp(groupOrder) >= 0 {
		s = randomUint256(t, r)
	}
	return secp256k1.IntToScalar(s)
}

func TestVRF_CheckSolidityPointAddition(t *testing.T) {
	t.Parallel()
	r := mrand.New(mrand.NewSource(5))
	numSamps := numSamples()
	for j := 0; j < numSamps; j++ {
		p1 := randomPoint(t, r)
		p2 := randomPoint(t, r)
		p1x, p1y := secp256k1.Coordinates(p1)
		p2x, p2y := secp256k1.Coordinates(p2)
		_, _, psz, err := verifier.ProjectiveECAdd(nil, p1x, p1y, p2x, p2y)
		require.NoError(t, err)
		zInv := i().ModInverse(psz, fieldSize)
		require.Equal(t, i().Mod(i().Mul(psz, zInv), fieldSize), one) // (sz * zInv) % fieldSize = 1
		actualSum, err := verifier.AffineECAdd(
			nil, pair(p1x, p1y), pair(p2x, p2y), zInv)
		require.NoError(t, err)
		actual := secp256k1.SetCoordinates(actualSum[0], actualSum[1])
		expected := point().Add(p1, p2)
		require.Equal(t, expected, actual)
	}
}

func TestVRF_CheckSolidityECMulVerify(t *testing.T) {
	t.Parallel()
	r := mrand.New(mrand.NewSource(6))
	numSamps := numSamples()
	for j := 0; j < numSamps; j++ {
		p := randomPoint(t, r)
		x, y := secp256k1.Coordinates(p)
		s := randomScalar(t, r)
		product := point().Mul(s, p)
		px, py := secp256k1.Coordinates(product)
		actual, err := verifier.EcmulVerify(nil, pair(x, y), secp256k1.ToInt(s),
			pair(px, py))
		require.NoError(t, err)
		require.True(t, actual)
	}
}

func TestVRF_CheckSolidityVerifyLinearCombinationWithGenerator(t *testing.T) {
	t.Parallel()
	r := mrand.New(mrand.NewSource(7))
	numSamps := numSamples()
	for j := 0; j < numSamps; j++ {
		c := randomScalar(t, r)
		s := randomScalar(t, r)
		p := randomPoint(t, r)
		sg := point().Mul(s, Generator)
		expectedPoint := point().Add(point().Mul(c, p), sg)
		expectedAddress := address(t, expectedPoint)
		px, py := secp256k1.Coordinates(p)
		actual, err := verifier.VerifyLinearCombinationWithGenerator(nil,
			secp256k1.ToInt(c), pair(px, py), secp256k1.ToInt(s), expectedAddress)
		require.NoError(t, err)
		require.True(t, actual)
	}
}

func asPair(p kyber.Point) [2]*big.Int {
	x, y := secp256k1.Coordinates(p)
	return pair(x, y)
}

func TestVRF_CheckSolidityLinearComination(t *testing.T) {
	t.Parallel()
	r := mrand.New(mrand.NewSource(8))
	numSamps := numSamples()
	for j := 0; j < numSamps; j++ {
		c := randomScalar(t, r)
		cNum := secp256k1.ToInt(c)
		p1 := randomPoint(t, r)
		p1Pair := asPair(p1)
		s := randomScalar(t, r)
		sNum := secp256k1.ToInt(s)
		p2 := randomPoint(t, r)
		p2Pair := asPair(p2)
		cp1 := point().Mul(c, p1)
		cp1x, cp1y := secp256k1.Coordinates(cp1)
		cp1Pair := pair(cp1x, cp1y)
		sp2 := point().Mul(s, p2)
		sp2x, sp2y := secp256k1.Coordinates(sp2)
		sp2Pair := pair(sp2x, sp2y)
		expected := point().Add(cp1, sp2)
		_, _, z, err := verifier.ProjectiveECAdd(nil, cp1x, cp1y, sp2x, sp2y)
		require.NoError(t, err)
		zInv := i().ModInverse(z, fieldSize)
		actual, err := verifier.LinearCombination(nil, cNum, p1Pair, cp1Pair, sNum,
			p2Pair, sp2Pair, zInv)
		require.NoError(t, err)
		actualPoint := secp256k1.SetCoordinates(actual[0], actual[1])
		require.Equal(t, expected, actualPoint)
	}
}

func TestVRF_CompareSolidityScalarFromCurve(t *testing.T) {
	t.Parallel()
	r := mrand.New(mrand.NewSource(9))
	numSamps := numSamples()
	for j := 0; j < numSamps; j++ {
		hash := randomPoint(t, r)
		hashPair := asPair(hash)
		pk := randomPoint(t, r)
		pkPair := asPair(pk)
		gamma := randomPoint(t, r)
		gammaPair := asPair(gamma)
		var uWitness [20]byte
		_, err := r.Read(uWitness[:])
		require.NoError(t, err)
		v := randomPoint(t, r)
		vPair := asPair(v)
		expected := ScalarFromCurve(hash, pk, gamma, uWitness, v)
		actual, err := verifier.ScalarFromCurvePoints(nil, hashPair, pkPair,
			gammaPair, uWitness, vPair)
		require.NoError(t, err)
		require.Equal(t, expected, actual)
	}
}

func TestVRF_MarshalProof(t *testing.T) {
	t.Parallel()
	r := mrand.New(mrand.NewSource(10))
	numSamps := numSamples()
	for j := 0; j < numSamps; j++ {
		sk := randomScalar(t, r)
		skNum := secp256k1.ToInt(sk)
		nonce := randomScalar(t, r)
		seed := randomUint256(t, r)
		proof, err := generateProofWithNonce(skNum, seed, secp256k1.ToInt(nonce))
		require.NoError(t, err)
		mproof, err := proof.MarshalForSolidityVerifier()
		require.NoError(t, err)
		response, err := verifier.RandomValueFromVRFProof(nil, mproof[:])
		require.NoError(t, err)
		require.True(t, response.Cmp(proof.Output) == 0)
	}
}
