package gokzg4844

import (
	"bytes"
	_ "embed"
	"encoding/hex"
	"sync"

	bls12381 "github.com/consensys/gnark-crypto/ecc/bls12-381"
)

// This library will not :
// - Check that the points are in the correct subgroup.
// - Check that setupG1Lagrange is the lagrange version of setupG1.
//
// Note: There is an embedded (via a //go:embed - compiler instruction) setup
// testKzgSetupStr, to which we do check those properties in a test function.

// JSONTrustedSetup is a struct used for serializing the trusted setup from/to JSON format.
//
// The intended use-case is that library users store the trusted setup in a JSON file and we provide such a file
// as part of the package.
type JSONTrustedSetup struct {
	SetupG2         []G2CompressedHexStr               `json:"g2_monomial"`
	SetupG1Lagrange [ScalarsPerBlob]G1CompressedHexStr `json:"g1_lagrange"`
}

// G1CompressedHexStr is a hex-string (with the 0x prefix) of a compressed G1 point.
type G1CompressedHexStr = string

// G2CompressedHexStr is a hex-string (with the 0x prefix) of a compressed G2 point.
type G2CompressedHexStr = string

// This is the test trusted setup, which SHOULD NOT BE USED IN PRODUCTION.
// The secret for this 1337.
//
//go:embed trusted_setup.json
var testKzgSetupStr string

// CheckTrustedSetupIsWellFormed checks whether the trusted setup is well-formed.
//
// To be specific, this checks that:
//   - All elements are in the correct subgroup.
func CheckTrustedSetupIsWellFormed(trustedSetup *JSONTrustedSetup) error {
	for i := 0; i < len(trustedSetup.SetupG1Lagrange); i++ {
		var point bls12381.G1Affine
		byts, err := hex.DecodeString(trim0xPrefix(trustedSetup.SetupG1Lagrange[i]))
		if err != nil {
			return err
		}
		_, err = point.SetBytes(byts)
		if err != nil {
			return err
		}
	}

	for i := 0; i < len(trustedSetup.SetupG2); i++ {
		var point bls12381.G2Affine
		byts, err := hex.DecodeString(trim0xPrefix(trustedSetup.SetupG2[i]))
		if err != nil {
			return err
		}
		_, err = point.SetBytes(byts)
		if err != nil {
			return err
		}
	}

	return nil
}

// parseTrustedSetup parses the trusted setup in `JSONTrustedSetup` format
// which contains hex encoded strings to corresponding group elements.
// Elements are assumed to be well-formed.
//
// This method wil panic if the points have not been serialized correctly.
func parseTrustedSetup(trustedSetup *JSONTrustedSetup) (bls12381.G1Affine, []bls12381.G1Affine, []bls12381.G2Affine) {
	// The G1 generator is the first element of the monomial G1 points.
	// We do not have that and so we use the fact that the setup started at
	// the canonical generator point.
	_, _, genG1, _ := bls12381.Generators()

	setupLagrangeG1Points := parseG1PointsNoSubgroupCheck(trustedSetup.SetupG1Lagrange[:])
	g2Points := parseG2PointsNoSubgroupCheck(trustedSetup.SetupG2)
	return genG1, setupLagrangeG1Points, g2Points
}

// parseG1PointNoSubgroupCheck parses a hex-string (with the 0x prefix) into a G1 point.
//
// This function performs no (expensive) subgroup checks, and should only be used
// for trusted inputs.
func parseG1PointNoSubgroupCheck(hexString string) (bls12381.G1Affine, error) {
	byts, err := hex.DecodeString(trim0xPrefix(hexString))
	if err != nil {
		return bls12381.G1Affine{}, err
	}

	var point bls12381.G1Affine
	noSubgroupCheck := bls12381.NoSubgroupChecks()
	d := bls12381.NewDecoder(bytes.NewReader(byts), noSubgroupCheck)

	return point, d.Decode(&point)
}

// parseG2PointNoSubgroupCheck parses a hex-string (with the 0x prefix) into a G2 point.
//
// This function performs no (expensive) subgroup checks, and should only be used
// for trusted inputs.
func parseG2PointNoSubgroupCheck(hexString string) (bls12381.G2Affine, error) {
	byts, err := hex.DecodeString(trim0xPrefix(hexString))
	if err != nil {
		return bls12381.G2Affine{}, err
	}

	var point bls12381.G2Affine
	noSubgroupCheck := bls12381.NoSubgroupChecks()
	d := bls12381.NewDecoder(bytes.NewReader(byts), noSubgroupCheck)

	return point, d.Decode(&point)
}

// parseG1PointsNoSubgroupCheck parses a slice hex-string (with the 0x prefix) into a
// slice of G1 points.
//
// This is essentially a parallelized version of calling [parseG1PointNoSubgroupCheck]
// on each element of the slice individually.
//
// This function performs no (expensive) subgroup checks, and should only be used
// for trusted inputs.
func parseG1PointsNoSubgroupCheck(hexStrings []string) []bls12381.G1Affine {
	numG1 := len(hexStrings)
	g1Points := make([]bls12381.G1Affine, numG1)

	var wg sync.WaitGroup
	wg.Add(numG1)
	for i := 0; i < numG1; i++ {
		go func(j int) {
			g1Point, err := parseG1PointNoSubgroupCheck(hexStrings[j])
			if err != nil {
				panic(err)
			}
			g1Points[j] = g1Point
			wg.Done()
		}(i)
	}
	wg.Wait()

	return g1Points
}

// parseG2PointsNoSubgroupCheck parses a slice hex-string (with the 0x prefix) into a
// slice of G2 points.
//
// This is essentially a parallelized version of calling [parseG2PointNoSubgroupCheck]
// on each element of the slice individually.
//
// This function performs no (expensive) subgroup checks, and should only be used
// for trusted inputs.
func parseG2PointsNoSubgroupCheck(hexStrings []string) []bls12381.G2Affine {
	numG2 := len(hexStrings)
	g2Points := make([]bls12381.G2Affine, numG2)

	var wg sync.WaitGroup
	wg.Add(numG2)
	for i := 0; i < numG2; i++ {
		go func(_i int) {
			g2Point, err := parseG2PointNoSubgroupCheck(hexStrings[_i])
			if err != nil {
				panic(err)
			}
			g2Points[_i] = g2Point
			wg.Done()
		}(i)
	}
	wg.Wait()

	return g2Points
}

// trim0xPrefix removes the "0x" from a hex-string.
func trim0xPrefix(hexString string) string {
	// Check that we are trimming off 0x
	if hexString[0:2] != "0x" {
		panic("hex string is not prefixed with 0x")
	}
	return hexString[2:]
}
