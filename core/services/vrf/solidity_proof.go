package vrf

// Logic for providing the precomputed values required by the solidity verifier,
// in binary-blob format.

import (
	"fmt"
	"math/big"

	"chainlink/core/services/signatures/secp256k1"

	"github.com/ethereum/go-ethereum/common"
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
	return rcurve.Point()
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
	rv.UWitness, err = secp256k1.EthereumAddress(u)
	if err != nil {
		return nil, err
	}
	rv.CGammaWitness = point().Mul(c, p.Gamma)
	hash, err := HashToCurve(p.PublicKey, p.Seed, func(*big.Int) {})
	if err != nil {
		return nil, err
	}
	rv.SHashWitness = point().Mul(s, hash)
	_, _, z := ProjectiveECAdd(rv.CGammaWitness, rv.SHashWitness)
	rv.ZInv = z.ModInverse(z, fieldSize)
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
func (p *SolidityProof) MarshalForSolidityVerifier() (MarshaledProof, error) {
	var rv MarshaledProof
	cursor := rv[:0]
	write := func(b []byte) { cursor = append(cursor, b...) }
	write(secp256k1.LongMarshal(p.P.PublicKey))
	write(secp256k1.LongMarshal(p.P.Gamma))
	write(asUint256(p.P.C))
	write(asUint256(p.P.S))
	write(asUint256(p.P.Seed))
	write(make([]byte, 12)) // Left-pad address to 32 bytes, with zeros
	write(p.UWitness[:])
	write(secp256k1.LongMarshal(p.CGammaWitness))
	write(secp256k1.LongMarshal(p.SHashWitness))
	write(asUint256(p.ZInv))
	if len(cursor) != ProofLength {
		return MarshaledProof{}, fmt.Errorf("wrong proof length: %d", len(rv))
	}
	return rv, nil
}

// MarshalForSolidityVerifier renders p as required by randomValueFromVRFProof
func (p *Proof) MarshalForSolidityVerifier() (MarshaledProof, error) {
	var rv MarshaledProof
	solidityProof, err := p.SolidityPrecalculations()
	if err != nil {
		return rv, err
	}
	rv, err = solidityProof.MarshalForSolidityVerifier()
	if err != nil {
		return rv, err
	}
	return rv, nil
}
