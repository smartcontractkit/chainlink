package gokzg4844

import (
	"encoding/json"

	"github.com/crate-crypto/go-kzg-4844/internal/kzg"
)

// Context holds the necessary configuration needed to create and verify proofs.
//
// Note: We could marshall this object so that clients won't need to process the SRS each time. The time to process is
// about 2-5 seconds.
type Context struct {
	domain    *kzg.Domain
	commitKey *kzg.CommitKey
	openKey   *kzg.OpeningKey
}

// BlsModulus is the bytes representation of the bls12-381 scalar field modulus.
//
// It matches [BLS_MODULUS] in the spec.
//
// [BLS_MODULUS]: https://github.com/ethereum/consensus-specs/blob/017a8495f7671f5fff2075a9bfc9238c1a0982f8/specs/deneb/polynomial-commitments.md#constants
var BlsModulus = [32]byte{
	0x73, 0xed, 0xa7, 0x53, 0x29, 0x9d, 0x7d, 0x48,
	0x33, 0x39, 0xd8, 0x08, 0x09, 0xa1, 0xd8, 0x05,
	0x53, 0xbd, 0xa4, 0x02, 0xff, 0xfe, 0x5b, 0xfe,
	0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x01,
}

// PointAtInfinity represents the serialized form of the point at infinity on the G1 group.
//
// It matches [G1_POINT_AT_INFINITY] in the spec.
//
// [G1_POINT_AT_INFINITY]: https://github.com/ethereum/consensus-specs/blob/017a8495f7671f5fff2075a9bfc9238c1a0982f8/specs/deneb/polynomial-commitments.md#constants
var PointAtInfinity = [48]byte{0xc0}

// NewContext4096Secure creates a new context object which will hold the state needed for one to use the KZG
// methods. "4096" denotes that we will only be able to commit to polynomials with at most 4096 evaluations. "Secure"
// denotes that this method is using a trusted setup file that was generated in an official
// ceremony. In particular, the trusted file being used was taken from the ethereum KZG ceremony.
func NewContext4096Secure() (*Context, error) {
	if ScalarsPerBlob != 4096 {
		// This is a library bug and so we panic.
		panic("this method is named `NewContext4096Insecure1337` we expect SCALARS_PER_BLOB to be 4096")
	}

	parsedSetup := JSONTrustedSetup{}

	err := json.Unmarshal([]byte(testKzgSetupStr), &parsedSetup)
	if err != nil {
		return nil, err
	}

	if ScalarsPerBlob != len(parsedSetup.SetupG1Lagrange) {
		// This is a library method and so we panic
		panic("this method is named `NewContext4096Insecure1337` we expect the number of G1 elements in the trusted setup to be 4096")
	}
	return NewContext4096(&parsedSetup)
}

// NewContext4096 creates a new context object which will hold the state needed for one to use the EIP-4844 methods. The
// 4096 represents the fact that without extra changes to the code, this context will only handle polynomials with 4096
// evaluations (degree 4095).
//
// Note: The G2 points do not have a fixed size. Technically, we could specify it to be 2, as this is the number of G2
// points that are required for KZG. However, the trusted setup in Ethereum has 65 since they want to use it for a
// future protocol: [Full Danksharding]. For this reason, we do not apply a fixed size, allowing the user to pass, say,
// 2 or 65.
//
// To initialize one must pass the parameters generated after the trusted setup, plus the lagrange version of the G1
// points. This function assumes that the G1 and G2 points are in order:
//
//   - G1points = {G, alpha * G, alpha^2 * G, ..., alpha^n * G}
//   - G2points = {H, alpha * H, alpha^2 * H, ..., alpha^n * H}
//   - Lagrange G1Points = {L_0(alpha^0) * G, L_1(alpha) * G, L_2(alpha^2) * G, ..., L_n(alpha^n) * G}
//
// [Full Danksharding]: https://notes.ethereum.org/@dankrad/new_sharding
func NewContext4096(trustedSetup *JSONTrustedSetup) (*Context, error) {
	// This should not happen for the ETH protocol
	// However since it's a public method, we add the check.
	if len(trustedSetup.SetupG2) < 2 {
		return nil, kzg.ErrMinSRSSize
	}

	// Parse the trusted setup from hex strings to G1 and G2 points
	genG1, setupLagrangeG1Points, setupG2Points := parseTrustedSetup(trustedSetup)

	// Get the generator points and the degree-1 element for G2 points
	// The generators are the degree-0 elements in the trusted setup
	//
	// This will never panic as we checked the minimum SRS size is >= 2
	// and `ScalarsPerBlob` is 4096
	genG2 := setupG2Points[0]
	alphaGenG2 := setupG2Points[1]

	commitKey := kzg.CommitKey{
		G1: setupLagrangeG1Points,
	}
	openingKey := kzg.OpeningKey{
		GenG1:   genG1,
		GenG2:   genG2,
		AlphaG2: alphaGenG2,
	}

	domain := kzg.NewDomain(ScalarsPerBlob)
	// Bit-Reverse the roots and the trusted setup according to the specs
	// The bit reversal is not needed for simple KZG however it was
	// implemented to make the step for full dank-sharding easier.
	commitKey.ReversePoints()
	domain.ReverseRoots()

	return &Context{
		domain:    domain,
		commitKey: &commitKey,
		openKey:   &openingKey,
	}, nil
}
