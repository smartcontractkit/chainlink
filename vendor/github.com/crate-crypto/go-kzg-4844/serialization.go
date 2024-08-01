package gokzg4844

import (
	bls12381 "github.com/consensys/gnark-crypto/ecc/bls12-381"
	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
	"github.com/crate-crypto/go-kzg-4844/internal/kzg"
	"github.com/crate-crypto/go-kzg-4844/internal/utils"
)

// CompressedG1Size is the number of bytes needed to represent a group element in G1 when compressed.
const CompressedG1Size = 48

// CompressedG2Size is the number of bytes needed to represent a group element in G2 when compressed.
const CompressedG2Size = 96

// SerializedScalarSize is the number of bytes needed to represent a field element corresponding to the order of the G1
// group.
//
// It matches [BYTES_PER_FIELD_ELEMENT] in the spec.
//
// [BYTES_PER_FIELD_ELEMENT]: https://github.com/ethereum/consensus-specs/blob/017a8495f7671f5fff2075a9bfc9238c1a0982f8/specs/deneb/polynomial-commitments.md#constants
const SerializedScalarSize = 32

// ScalarsPerBlob is the number of serialized scalars in a blob.
//
// It matches [FIELD_ELEMENTS_PER_BLOB] in the spec.
//
// Note: These scalars are not guaranteed to be valid (a value less than [BLS_MODULUS]). If any of the scalars in a blob
// are invalid (non-canonical), an error will be returned on deserialization.
//
// [BLS_MODULUS]: https://github.com/ethereum/consensus-specs/blob/017a8495f7671f5fff2075a9bfc9238c1a0982f8/specs/deneb/polynomial-commitments.md#constants
// [FIELD_ELEMENTS_PER_BLOB]: https://github.com/ethereum/consensus-specs/blob/017a8495f7671f5fff2075a9bfc9238c1a0982f8/specs/deneb/polynomial-commitments.md#blob
const ScalarsPerBlob = 4096

type (
	// G1Point matches [G1Point] in the spec.
	//
	// [G1Point]: https://github.com/ethereum/consensus-specs/blob/017a8495f7671f5fff2075a9bfc9238c1a0982f8/specs/deneb/polynomial-commitments.md#custom-types
	G1Point [CompressedG1Size]byte

	// G2Point matches [G2Point] in the spec.
	//
	// [G2Point]: https://github.com/ethereum/consensus-specs/blob/017a8495f7671f5fff2075a9bfc9238c1a0982f8/specs/deneb/polynomial-commitments.md#custom-types
	G2Point [CompressedG2Size]byte

	// Scalar matches [BLSFieldElement] in the spec.
	//
	// [BLSFieldElement]: https://github.com/ethereum/consensus-specs/blob/017a8495f7671f5fff2075a9bfc9238c1a0982f8/specs/deneb/polynomial-commitments.md#custom-types
	Scalar [SerializedScalarSize]byte

	// Blob is a flattened representation of a serialized polynomial.
	//
	// It matches [Blob] in the spec.
	//
	// [Blob]: https://github.com/ethereum/consensus-specs/blob/017a8495f7671f5fff2075a9bfc9238c1a0982f8/specs/deneb/polynomial-commitments.md#custom-types
	Blob [ScalarsPerBlob * SerializedScalarSize]byte

	// KZGProof is a serialized commitment to the quotient polynomial.
	//
	// It matches [KZGProof] in the spec.
	//
	// [KZGProof]: https://github.com/ethereum/consensus-specs/blob/017a8495f7671f5fff2075a9bfc9238c1a0982f8/specs/deneb/polynomial-commitments.md#custom-types
	KZGProof G1Point

	// KZGCommitment is a serialized commitment to a polynomial.
	//
	// It matches [KZGCommitment] in the spec.
	//
	// [KZGCommitment]: https://github.com/ethereum/consensus-specs/blob/017a8495f7671f5fff2075a9bfc9238c1a0982f8/specs/deneb/polynomial-commitments.md#custom-types
	KZGCommitment G1Point
)

// SerializeG1Point converts a [bls12381.G1Affine] to [G1Point].
func SerializeG1Point(affine bls12381.G1Affine) G1Point {
	return affine.Bytes()
}

// deserializeG1Point converts a [G1Point] to the internal [bls12381.G1Affine] type. It will return an error if the
// point is not on the group or if the point is not in the correct subgroup.
//
// It implements [validate_kzg_g1].
//
// [validate_kzg_g1]: https://github.com/ethereum/consensus-specs/blob/017a8495f7671f5fff2075a9bfc9238c1a0982f8/specs/deneb/polynomial-commitments.md#validate_kzg_g1
func deserializeG1Point(serPoint G1Point) (bls12381.G1Affine, error) {
	var point bls12381.G1Affine
	_, err := point.SetBytes(serPoint[:])
	if err != nil {
		return bls12381.G1Affine{}, err
	}
	return point, nil
}

// DeserializeKZGCommitment implements [bytes_to_kzg_commitment].
//
// [bytes_to_kzg_commitment]: https://github.com/ethereum/consensus-specs/blob/017a8495f7671f5fff2075a9bfc9238c1a0982f8/specs/deneb/polynomial-commitments.md#bytes_to_kzg_commitment
func DeserializeKZGCommitment(commitment KZGCommitment) (bls12381.G1Affine, error) {
	return deserializeG1Point(G1Point(commitment))
}

// DeserializeKZGProof implements [bytes_to_kzg_proof].
//
// [bytes_to_kzg_proof]: https://github.com/ethereum/consensus-specs/blob/017a8495f7671f5fff2075a9bfc9238c1a0982f8/specs/deneb/polynomial-commitments.md#bytes_to_kzg_proof
func DeserializeKZGProof(proof KZGProof) (bls12381.G1Affine, error) {
	return deserializeG1Point(G1Point(proof))
}

// DeserializeBlob implements [blob_to_polynomial].
//
// [blob_to_polynomial]: https://github.com/ethereum/consensus-specs/blob/017a8495f7671f5fff2075a9bfc9238c1a0982f8/specs/deneb/polynomial-commitments.md#blob_to_polynomial
func DeserializeBlob(blob Blob) (kzg.Polynomial, error) {
	poly := make(kzg.Polynomial, ScalarsPerBlob)
	for i := 0; i < ScalarsPerBlob; i++ {
		chunk := blob[i*SerializedScalarSize : (i+1)*SerializedScalarSize]
		serScalar := (*Scalar)(chunk)
		scalar, err := DeserializeScalar(*serScalar)
		if err != nil {
			return nil, err
		}
		poly[i] = scalar
	}
	return poly, nil
}

// DeserializeScalar implements [bytes_to_bls_field].
//
// Note: Returns an error if the scalar is not in the range [0, p-1] (inclusive) where `p` is the prime associated with the scalar field.
//
// [bytes_to_bls_field]: https://github.com/ethereum/consensus-specs/blob/017a8495f7671f5fff2075a9bfc9238c1a0982f8/specs/deneb/polynomial-commitments.md#bytes_to_bls_field
func DeserializeScalar(serScalar Scalar) (fr.Element, error) {
	scalar, err := utils.ReduceCanonicalBigEndian(serScalar[:])
	if err != nil {
		return fr.Element{}, ErrNonCanonicalScalar
	}
	return scalar, nil
}

// SerializeScalar converts a [fr.Element] to [Scalar].
func SerializeScalar(element fr.Element) Scalar {
	return element.Bytes()
}

// SerializePoly converts a [kzg.Polynomial] to [Blob].
//
// Note: This method is never used in the API because we always expect a byte array and will never receive deserialized
// field elements. We include it so that upstream fuzzers do not need to reimplement it.
func SerializePoly(poly kzg.Polynomial) Blob {
	var blob Blob
	for i := 0; i < ScalarsPerBlob; i++ {
		chunk := blob[i*SerializedScalarSize : (i+1)*SerializedScalarSize]
		serScalar := SerializeScalar(poly[i])
		copy(chunk, serScalar[:])
	}
	return blob
}
