package proof

// Logic for providing the precomputed values required by the solidity verifier,
// in binary-blob format.

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/vrfkey"
	"github.com/smartcontractkit/chainlink/core/services/signatures/secp256k1"
	"github.com/smartcontractkit/chainlink/core/utils"
	bm "github.com/smartcontractkit/chainlink/core/utils/big_math"
	"go.dedis.ch/kyber/v3"
)

// SolidityProof contains precalculations which VRF.sol needs to verify proofs
type SolidityProof struct {
	P                           *vrfkey.Proof  // The core proof
	UWitness                    common.Address // Address of P.C*P.PK+P.S*G
	CGammaWitness, SHashWitness kyber.Point    // P.C*P.Gamma, P.S*HashToCurve(P.Seed)
	ZInv                        *big.Int       // Inverse of Z coord from ProjectiveECAdd(CGammaWitness, SHashWitness)
}

// String returns the values in p, in hexadecimal format
func (p *SolidityProof) String() string {
	return fmt.Sprintf(
		"SolidityProof{P: %s, UWitness: %x, CGammaWitness: %s, SHashWitness: %s, ZInv: %x}",
		p.P, p.UWitness, p.CGammaWitness, p.SHashWitness, p.ZInv)
}

func point() kyber.Point {
	return vrfkey.Secp256k1Curve.Point()
}

// SolidityPrecalculations returns the precomputed values needed by the solidity
// verifier, or an error on failure.
func SolidityPrecalculations(p *vrfkey.Proof) (*SolidityProof, error) {
	var rv SolidityProof
	rv.P = p
	c := secp256k1.IntToScalar(p.C)
	s := secp256k1.IntToScalar(p.S)
	u := point().Add(point().Mul(c, p.PublicKey), point().Mul(s, vrfkey.Generator))
	var err error
	rv.UWitness = secp256k1.EthereumAddress(u)
	rv.CGammaWitness = point().Mul(c, p.Gamma)
	hash, err := vrfkey.HashToCurve(p.PublicKey, p.Seed, func(*big.Int) {})
	if err != nil {
		return nil, err
	}
	rv.SHashWitness = point().Mul(s, hash)
	_, _, z := vrfkey.ProjectiveECAdd(rv.CGammaWitness, rv.SHashWitness)
	rv.ZInv = z.ModInverse(z, vrfkey.FieldSize)
	return &rv, nil
}

// Length of marshaled proof, in bytes
const ProofLength = 64 + // PublicKey
	64 + // Gamma
	32 + // C
	32 + // S
	32 + // Seed
	32 + // uWitness (gets padded to 256 bits, even though it's only 160)
	64 + // cGammaWitness
	64 + // sHashWitness
	32 // zInv  (Leave Output out, because that can be efficiently calculated)

// MarshaledProof contains a VRF proof for randomValueFromVRFProof.
//
// NB: when passing one of these to randomValueFromVRFProof via the geth
// blockchain simulator, it must be passed as a slice ("proof[:]"). Passing it
// as-is sends hundreds of single bytes, each padded to their own 32-byte word.
type MarshaledProof [ProofLength]byte

// String returns m as 0x-hex bytes
func (m MarshaledProof) String() string {
	return fmt.Sprintf("0x%x", [ProofLength]byte(m))
}

// MarshalForSolidityVerifier renders p as required by randomValueFromVRFProof
func (p *SolidityProof) MarshalForSolidityVerifier() (proof MarshaledProof) {
	cursor := proof[:0]
	write := func(b []byte) { cursor = append(cursor, b...) }
	write(secp256k1.LongMarshal(p.P.PublicKey))
	write(secp256k1.LongMarshal(p.P.Gamma))
	write(utils.Uint256ToBytes32(p.P.C))
	write(utils.Uint256ToBytes32(p.P.S))
	write(utils.Uint256ToBytes32(p.P.Seed))
	write(make([]byte, 12)) // Left-pad address to 32 bytes, with zeros
	write(p.UWitness[:])
	write(secp256k1.LongMarshal(p.CGammaWitness))
	write(secp256k1.LongMarshal(p.SHashWitness))
	write(utils.Uint256ToBytes32(p.ZInv))
	if len(cursor) != ProofLength {
		panic(fmt.Errorf("wrong proof length: %d", len(proof)))
	}
	return proof
}

// MarshalForSolidityVerifier renders p as required by randomValueFromVRFProof
func MarshalForSolidityVerifier(p *vrfkey.Proof) (MarshaledProof, error) {
	var rv MarshaledProof
	solidityProof, err := SolidityPrecalculations(p)
	if err != nil {
		return rv, err
	}
	return solidityProof.MarshalForSolidityVerifier(), nil
}

func UnmarshalSolidityProof(proof []byte) (rv vrfkey.Proof, err error) {
	failedProof := vrfkey.Proof{}
	if len(proof) != ProofLength {
		return failedProof, fmt.Errorf(
			"VRF proof is %d bytes long, should be %d: \"%x\"", len(proof),
			ProofLength, proof)
	}
	if rv.PublicKey, err = secp256k1.LongUnmarshal(proof[:64]); err != nil {
		return failedProof, errors.Wrapf(err, "while reading proof public key")
	}
	rawGamma := proof[64:128]
	if rv.Gamma, err = secp256k1.LongUnmarshal(rawGamma); err != nil {
		return failedProof, errors.Wrapf(err, "while reading proof gamma")
	}
	rv.C = bm.I().SetBytes(proof[128:160])
	rv.S = bm.I().SetBytes(proof[160:192])
	rv.Seed = bm.I().SetBytes(proof[192:224])
	rv.Output = utils.MustHash(string(vrfkey.RandomOutputHashPrefix) +
		string(rawGamma)).Big()
	return rv, nil
}
