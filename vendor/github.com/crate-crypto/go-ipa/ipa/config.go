package ipa

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
	"runtime"

	"github.com/crate-crypto/go-ipa/bandersnatch/fp"
	"github.com/crate-crypto/go-ipa/bandersnatch/fr"
	"github.com/crate-crypto/go-ipa/banderwagon"
	"github.com/crate-crypto/go-ipa/common"
)

// IPAConfig contains all the necessary information to create an IPA related proofs,
// such as the SRS, Q, and precomputed weights for the barycentric formula.
type IPAConfig struct {
	SRS []banderwagon.Element
	Q   banderwagon.Element

	PrecompMSM         banderwagon.MSMPrecomp
	PrecomputedWeights *PrecomputedWeights
	// The number of rounds the prover and verifier must complete
	// in the IPA argument, this will be log2 of the size of the input vectors
	// since the vector is halved on each round
	numRounds uint32
}

// NewIPASettings generates the SRS, Q and precomputed weights for the barycentric formula.
// The SRS is generated as common.VectorLength random points where the relative discrete log is
// not known between each generator.
func NewIPASettings() (*IPAConfig, error) {
	srs := GenerateRandomPoints(common.VectorLength)
	precompMSM, err := banderwagon.NewPrecompMSM(srs)
	if err != nil {
		return nil, fmt.Errorf("creating precomputed MSM: %s", err)
	}
	return &IPAConfig{
		SRS:                srs,
		Q:                  banderwagon.Generator,
		PrecompMSM:         precompMSM,
		PrecomputedWeights: NewPrecomputedWeights(),
		numRounds:          computeNumRounds(common.VectorLength),
	}, nil
}

// MultiScalar computes the multi scalar multiplication of points and scalars.
func MultiScalar(points []banderwagon.Element, scalars []fr.Element) (banderwagon.Element, error) {
	var result banderwagon.Element
	result.SetIdentity()

	res, err := result.MultiExp(points, scalars, banderwagon.MultiExpConfig{NbTasks: runtime.NumCPU(), ScalarsMont: true})
	if err != nil {
		return banderwagon.Element{}, fmt.Errorf("mult exponentiation was not successful: %w", err)
	}

	return *res, nil
}

// Commit calculates the Pedersen Commitment of a polynomial polynomial
// in evaluation form using the SRS.
func (ic *IPAConfig) Commit(polynomial []fr.Element) banderwagon.Element {
	return ic.PrecompMSM.MSM(polynomial)
}

// commit commits to a polynomial using the input group elements
func commit(groupElements []banderwagon.Element, polynomial []fr.Element) (banderwagon.Element, error) {
	if len(groupElements) != len(polynomial) {
		return banderwagon.Element{}, fmt.Errorf("group elements and polynomial are different sizes, %d != %d", len(groupElements), len(polynomial))
	}
	return MultiScalar(groupElements, polynomial)
}

// InnerProd computes the inner product of a and b.
func InnerProd(a []fr.Element, b []fr.Element) (fr.Element, error) {
	if len(a) != len(b) {
		return fr.Element{}, fmt.Errorf("a and b are different sizes, %d != %d", len(a), len(b))
	}

	result := fr.Zero()
	for i := 0; i < len(a); i++ {
		var tmp fr.Element

		tmp.Mul(&a[i], &b[i])
		result.Add(&result, &tmp)
	}

	return result, nil
}

// Computes c[i] =a[i] + b[i] * x
// returns c
func foldScalars(a []fr.Element, b []fr.Element, x fr.Element) ([]fr.Element, error) {
	if len(a) != len(b) {
		return nil, fmt.Errorf("slices not equal length")
	}
	result := make([]fr.Element, len(a))
	for i := 0; i < len(a); i++ {
		var bx fr.Element
		bx.Mul(&x, &b[i])
		result[i].Add(&bx, &a[i])
	}

	return result, nil
}

// Computes c[i] =a[i] + b[i] * x
// returns c
func foldPoints(a []banderwagon.Element, b []banderwagon.Element, x fr.Element) ([]banderwagon.Element, error) {
	if len(a) != len(b) {
		return nil, fmt.Errorf("slices not equal length")
	}

	result := make([]banderwagon.Element, len(a))
	for i := 0; i < len(a); i++ {
		var bx banderwagon.Element
		bx.ScalarMul(&b[i], &x)
		result[i].Add(&bx, &a[i])
	}
	return result, nil
}

// Splits a slice of scalars into two slices of equal length
// Eg [S1,S2,S3,S4] becomes [S1,S2] , [S3,S4]
func splitScalars(x []fr.Element) ([]fr.Element, []fr.Element, error) {
	if len(x)%2 != 0 {
		return nil, nil, fmt.Errorf("slice should have an even length")
	}

	mid := len(x) / 2
	return x[:mid], x[mid:], nil
}

// Splits a slice of points into two slices of equal length
// Eg [P1,P2,P3,P4,P5,P6] becomes [P1,P2,P3] , [P4,P5,P6]
func splitPoints(x []banderwagon.Element) ([]banderwagon.Element, []banderwagon.Element, error) {
	if len(x)%2 != 0 {
		return nil, nil, fmt.Errorf("slice should have an even length")
	}
	mid := len(x) / 2

	return x[:mid], x[mid:], nil
}

// This function does log2(vector_size)
//
// Since we do not allow for 0 size vectors, this is checked
// since we also do not allow for vectors which are not powers of 2, this is also checked
//
// It is okay to panic here, because the input is a constant, so it will panic before
// any proofs are made.
func computeNumRounds(vectorSize uint32) uint32 {
	// Check if this number is 0
	// zero is not a valid input to this function for our usecase
	if vectorSize == 0 {
		panic("zero is not a valid input")
	}

	// See: https://stackoverflow.com/a/600306
	isPow2 := (vectorSize & (vectorSize - 1)) == 0

	if !isPow2 {
		panic("non power of 2 numbers are not valid inputs")
	}

	res := math.Log2(float64(vectorSize))

	return uint32(res)
}

// GenerateRandomPoints generates numPoints random points on the curve using
// hardcoded seed.
func GenerateRandomPoints(numPoints uint64) []banderwagon.Element {
	seed := "eth_verkle_oct_2021"

	points := []banderwagon.Element{}

	var increment uint64 = 0

	for uint64(len(points)) != numPoints {

		digest := sha256.New()
		digest.Write([]byte(seed))

		b := make([]byte, 8)
		binary.BigEndian.PutUint64(b, increment)
		digest.Write(b)

		hash := digest.Sum(nil)

		var x fp.Element
		x.SetBytes(hash)

		increment++

		x_as_bytes := x.Bytes()
		var point_found banderwagon.Element
		err := point_found.SetBytes(x_as_bytes[:])
		if err != nil {
			// This point is not in the correct subgroup or on the curve
			continue
		}
		points = append(points, point_found)

	}

	return points
}
