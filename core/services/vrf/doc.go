// Package vrf provides a cryptographically secure pseudo-random number generator.

// Numbers are deterministically generated from seeds and a secret key, and are
// statistically indistinguishable from uniform sampling from {0,...,2**256-1},
// to computationally-bounded observers who know the seeds, don't know the key,
// and only see the generated numbers. But each number also comes with a proof
// that it was generated according to the procedure mandated by a public key
// associated with that secret key.
//
// See VRF.sol for design notes.
//
// Usage
// -----
//
// You should probably not be using this directly.
// chainlink/store/core/models/vrfkey.PrivateKey provides a simple, more
// misuse-resistant interface to the same functionality, via the CreateKey and
// MarshaledProof methods.
//
// Nonetheless, a secret key sk should be securely sampled uniformly from
// {0,...,Order-1}. Its public key can be calculated from it by
//
//   secp256k1.Secp256k1{}.Point().Mul(secretKey, Generator)
//
// To generate random output from a big.Int seed, pass sk and the seed to
// GenerateProof, and use the Output field of the returned Proof object.
//
// To verify a Proof object p, run p.Verify(); or to verify it on-chain pass
// p.MarshalForSolidityVerifier() to randomValueFromVRFProof on the VRF solidity
// contract.

package vrf
