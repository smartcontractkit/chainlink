package gokzg4844

import (
	"crypto/sha256"
	"encoding/binary"

	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
)

// DomSepProtocol is a Domain Separator to identify the protocol.
//
// It matches [FIAT_SHAMIR_PROTOCOL_DOMAIN] in the spec.
//
// [FIAT_SHAMIR_PROTOCOL_DOMAIN]: https://github.com/ethereum/consensus-specs/blob/017a8495f7671f5fff2075a9bfc9238c1a0982f8/specs/deneb/polynomial-commitments.md#blob
const DomSepProtocol = "FSBLOBVERIFY_V1_"

// computeChallenge is provided to match the spec at [compute_challenge].
//
// [compute_challenge]: https://github.com/ethereum/consensus-specs/blob/017a8495f7671f5fff2075a9bfc9238c1a0982f8/specs/deneb/polynomial-commitments.md#compute_challenge
func computeChallenge(blob Blob, commitment KZGCommitment) fr.Element {
	polyDegreeBytes := u64ToByteArray16(ScalarsPerBlob)
	data := append([]byte(DomSepProtocol), polyDegreeBytes...)
	data = append(data, blob[:]...)
	data = append(data, commitment[:]...)

	return hashToBLSField(data)
}

// hashToBLSField hashed the given binary data to a field element according to [hash_to_bls_field].
//
// [hash_to_bls_field]: https://github.com/ethereum/consensus-specs/blob/017a8495f7671f5fff2075a9bfc9238c1a0982f8/specs/deneb/polynomial-commitments.md#hash_to_bls_field
func hashToBLSField(data []byte) fr.Element {
	digest := sha256.Sum256(data)

	// Now interpret those bytes as a field element
	var challenge fr.Element
	challenge.SetBytes(digest[:])

	return challenge
}

// u64ToByteArray16 converts a uint64 to a byte slice of length 16 in big endian format. This implies that the first 8 bytes of the result are always 0.
func u64ToByteArray16(number uint64) []byte {
	bytes := make([]byte, 16)
	binary.BigEndian.PutUint64(bytes[8:], number)
	return bytes
}
