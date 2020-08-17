package vrf

// Logic for providing the precomputed values required by the solidity verifier,
// in binary-blob format.

import (
	"fmt"
	"math/big"

	"github.com/smartcontractkit/chainlink/core/services/signatures/secp256k1"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"go.dedis.ch/kyber/v3"
)

// SolidityProof contains precalculations which VRF.sol needs to verifiy proofs
type SolidityProof struct {
	P                           *Proof         // The core proof
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
	return secp256k1Curve.Point()
}

// SolidityPrecalculations returns the precomputed values needed by the solidity
// verifier, or an error on failure.
func (p *Proof) SolidityPrecalculations() (*SolidityProof, error) {
	var rv SolidityProof
	rv.P = p
	c := secp256k1.IntToScalar(p.C)
	s := secp256k1.IntToScalar(p.S)
	u := point().Add(point().Mul(c, p.PublicKey), point().Mul(s, Generator))
	var err error
	rv.UWitness = secp256k1.EthereumAddress(u)
	rv.CGammaWitness = point().Mul(c, p.Gamma)
	hash, err := HashToCurve(p.PublicKey, p.Seed, func(*big.Int) {})
	if err != nil {
		return nil, err
	}
	rv.SHashWitness = point().Mul(s, hash)
	_, _, z := ProjectiveECAdd(rv.CGammaWitness, rv.SHashWitness)
	rv.ZInv = z.ModInverse(z, FieldSize)
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
	write(uint256ToBytes32(p.P.C))
	write(uint256ToBytes32(p.P.S))
	write(uint256ToBytes32(p.P.Seed))
	write(make([]byte, 12)) // Left-pad address to 32 bytes, with zeros
	write(p.UWitness[:])
	write(secp256k1.LongMarshal(p.CGammaWitness))
	write(secp256k1.LongMarshal(p.SHashWitness))
	write(uint256ToBytes32(p.ZInv))
	if len(cursor) != ProofLength {
		panic(fmt.Errorf("wrong proof length: %d", len(proof)))
	}
	return proof
}

// MarshalForSolidityVerifier renders p as required by randomValueFromVRFProof
func (p *Proof) MarshalForSolidityVerifier() (MarshaledProof, error) {
	var rv MarshaledProof
	solidityProof, err := p.SolidityPrecalculations()
	if err != nil {
		return rv, err
	}
	return solidityProof.MarshalForSolidityVerifier(), nil
}

func UnmarshalSolidityProof(proof []byte) (rv Proof, err error) {
	failedProof := Proof{}
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
	rv.C = i().SetBytes(proof[128:160])
	rv.S = i().SetBytes(proof[160:192])
	rv.Seed = i().SetBytes(proof[192:224])
	rv.Output = utils.MustHash(string(vrfRandomOutputHashPrefix) +
		string(rawGamma)).Big()
	return rv, nil
}
