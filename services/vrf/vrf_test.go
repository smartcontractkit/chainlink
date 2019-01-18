package vrf

import (
	"fmt"
	"math"
	"testing"

	"github.com/ethereum/go-ethereum/crypto/bn256"
	"github.com/stretchr/testify/assert"
)

var insecureKeyPair = KeyPair{
	Public: new(bn256.G1).Add(Generator, Generator),
	secret: two.i,
}

func TestVRF_isSquare(t *testing.T) {
	assert.True(t, isSquare(four))
	assert.False(t, isSquare(P.sub(one)), "P-1 is not square in GF(P)")
}

func TestVRF_squareRoot(t *testing.T) {
	assert.Equal(t, two, four.modSqrt(), "4^{(P-1)/4} = 2 in GF(P)")
}

func TestVRF_ySquared(t *testing.T) {
	assert.Equal(t, bi(2*2*2+3), ySquared(two), "11=2^3+3 in GF(P)")
}

func TestVRF_isCurveXOrdinate(t *testing.T) {
	assert.True(t, isCurveXOrdinate(one), "2^2=1^3+3")
	assert.False(t, isCurveXOrdinate(four),
		"There's no y s.t. y^2=4^3+1 in GF(P)")
}

func TestVRF_CoordsFromPoint(t *testing.T) {
	x, y := CoordsFromPoint(Generator)
	assert.Equal(t, x, one.i, "Wrong x ordinate from Generator")
	assert.Equal(t, y, two.i, "Wrong y ordinate from Generator")
}

func bigFromHex(s string) integer {
	n, ok := i().SetString(s, 16)
	if !ok {
		panic(fmt.Errorf(`failed to convert "%s" as hex to big.Int`, s))
	}
	return integer{n}
}

func TestVRF_zqHash(t *testing.T) {
	log2Mod := 5.0
	modulus := bi(uint(math.Pow(2, log2Mod)))
	bitMask := bi(uint(math.Pow(2, log2Mod+1) - 1))
	reHashTriggeringSeed := zero
	hash, err := hashUint256s(reHashTriggeringSeed)
	if err != nil {
		panic(err)
	}
	hash = hash.and(bitMask)
	assert.True(t, hash.greaterThan(modulus),
		`need an example which hashes to something bigger than the
modulus, to test the rehash logic.`)
	zqH, err := zqHash(modulus, reHashTriggeringSeed)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, zqH, bi(13))
}

func TestVRF_hashToCurve(t *testing.T) {
	reHashTriggeringInput := []integer{one, two, bi(5)}
	x, err := zqHash(P, reHashTriggeringInput...)
	if err != nil {
		panic(err)
	}
	assert.False(t, isCurveXOrdinate(x),
		`need an example where first hash is not an x-ordinate for any
point on the curve, to exercise rehash logic.`)
	p, err := hashToCurve(reHashTriggeringInput[0],
		reHashTriggeringInput[1], reHashTriggeringInput[2])
	if err != nil {
		panic(err)
	}
	px, py := CoordsFromPoint(p)
	eX := "247154f2ce523897365341b03669e1061049e801e8750ae708e1cb02f36cb225"
	eY := "16e1157d5b94324127e094abe222a05a5c47be3124254a6aa047d5e1f2d864ea"
	assert.Equal(t, bigFromHex(eX), integer{px},
		"x ordinate of hashToCurve case tested in VRF_test.js has changed")
	assert.Equal(t, bigFromHex(eY), integer{py},
		"y ordinate of hashToCurve case tested in VRF_test.js has changed")
}

func TestVRF_scalarFromCurve(t *testing.T) {
	g := Generator
	s, err := scalarFromCurve(g, g, g, g, g)
	if err != nil {
		panic(err)
	}
	eS := "57bf013147ceec913f17ef97d3bcfad8315d99752af81f8913ad1c88493e669"
	assert.Equal(t, bigFromHex(eS), s,
		"scalarFromCurve case tested in VRF_test.js has changed")
}

func pointsEqual(p1, p2 *bn256.G1) bool {
	s1, _ := scalarFromCurve(p1)
	s2, _ := scalarFromCurve(p2)
	return s1.equal(s2)
}

func TestVRF_Marshalling(t *testing.T) {
	assert.Panics(t, func() { insecureKeyPair.Marshal([]byte{}) }, "Marshal should be disabled")
	assert.Panics(t, func() { insecureKeyPair.Unmarshal([]byte{}) }, "Unmarshal should be disabled")
	passPhrase := []byte("sekret")
	m, err := insecureKeyPair.MarshalEncryptedKeyPair(passPhrase)
	if err != nil {
		panic(err)
	}
	up, err := UnmarshalPublicKey(m)
	if err != nil {
		panic(err)
	}
	assert.True(t, pointsEqual(insecureKeyPair.Public, up),
		"Unmarshalling should return correct public key.")
	u, err := UnmarshalEncryptedKeyPair(m, passPhrase)
	if err != nil {
		panic(err)
	}
	assert.True(t, pointsEqual(insecureKeyPair.Public, u.Public) &&
		(integer{u.secret}).equal((integer{insecureKeyPair.secret})),
		"Unmarshalling should return the same key pair")
}

func TestVRF_GenerateProof(t *testing.T) {
	seed := zero
	proof, err := GenerateProof(&insecureKeyPair, seed.i)
	if err != nil {
		panic(err)
	}
	assert.True(t, pointsEqual(insecureKeyPair.Public, proof.PublicKey))
	gammaX, gammaY := CoordsFromPoint(proof.Gamma)
	gX := "26feb384a4a3f28742d0e0e0f5458474ba54ef9816d4d31f3bf538dfcf67cf3f"
	gY := "1eaed2431dd78ad75dd0c9f013cabff4f1d8c4c83cda79fff3855c988a3606d8"
	assert.Equal(t, bigFromHex(gX), integer{gammaX},
		"x ordinate of gamma tested in VRF_test.js has changed.")
	assert.Equal(t, bigFromHex(gY), integer{gammaY},
		"y ordinate of gamma tested in VRF_test.js has changed.")
	verification, err := proof.VerifyProof()
	if err != nil {
		panic(err)
	}
	assert.True(t, verification, "Generated proof should verify.")
}
