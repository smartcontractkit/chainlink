package ipa

import (
	"encoding/binary"
	"fmt"
	"io"
	"math/big"

	"github.com/crate-crypto/go-ipa/bandersnatch/fr"
	"github.com/crate-crypto/go-ipa/banderwagon"
	"github.com/crate-crypto/go-ipa/common"
)

var maxEvalPointInsideDomain fr.Element

func init() {
	maxEvalPointInsideDomain.SetUint64(common.VectorLength - 1)
}

// The following are unexported labels to be used in Fiat-Shamir during the
// inner-product argument protocol.
//
// The following is a short description on how they're used in the protocol:
// 1. Append the domain separator. (labelDomainSep)
// 2. Append the commitment to the polynomial. (labelC)
// 3. Append the input point. (labelInputPoint)
// 4. Append the output point. (labelOutputPoint)
// 5. Pull the re-scaling factor `w` to scale Q. (labelW).
// 6. For each round of the IPA protocol:
//    a. Append the resulting point C_L. (labelL)
//    b. Append the resulting point C_R. (labelR)
//    c. Pull the random scalar-field element `x`. (labelX)
//
// Note: this package must not mutate these label values, nor pass them to
// parts of the code that would mutate them.

var (
	labelDomainSep   = []byte("ipa")
	labelC           = []byte("C")
	labelInputPoint  = []byte("input point")
	labelOutputPoint = []byte("output point")
	labelW           = []byte("w")
	labelL           = []byte("L")
	labelR           = []byte("R")
	labelX           = []byte("x")
)

// IPAProof is an inner product argument proof.
type IPAProof struct {
	L        []banderwagon.Element
	R        []banderwagon.Element
	A_scalar fr.Element
}

// CreateIPAProof creates an IPA proof for a committed polynomial in evaluation form.
// `a` are the evaluation of the polynomial in the domain, and `evalPoint` represents the
// evaluation point. The evaluation of the polynomial at such point is computed automatically.
func CreateIPAProof(transcript *common.Transcript, ic *IPAConfig, commitment banderwagon.Element, a []fr.Element, evalPoint fr.Element) (IPAProof, error) {
	transcript.DomainSep(labelDomainSep)

	b := computeBVector(ic, evalPoint)

	inner_prod, err := InnerProd(a, b)
	if err != nil {
		return IPAProof{}, fmt.Errorf("could not compute inner product: %w", err)
	}

	transcript.AppendPoint(&commitment, labelC)
	transcript.AppendScalar(&evalPoint, labelInputPoint)
	transcript.AppendScalar(&inner_prod, labelOutputPoint)
	w := transcript.ChallengeScalar(labelW)

	var q banderwagon.Element
	q.ScalarMul(&ic.Q, &w)

	num_rounds := ic.numRounds

	current_basis := ic.SRS

	L := make([]banderwagon.Element, num_rounds)
	R := make([]banderwagon.Element, num_rounds)

	for i := 0; i < int(num_rounds); i++ {

		a_L, a_R, err := splitScalars(a)
		if err != nil {
			return IPAProof{}, fmt.Errorf("could not split a scalars: %w", err)
		}

		b_L, b_R, err := splitScalars(b)
		if err != nil {
			return IPAProof{}, fmt.Errorf("could not split b scalars: %w", err)
		}

		G_L, G_R, err := splitPoints(current_basis)
		if err != nil {
			return IPAProof{}, fmt.Errorf("could not split G points: %w", err)
		}

		z_L, err := InnerProd(a_R, b_L)
		if err != nil {
			return IPAProof{}, fmt.Errorf("could not compute a_r*b_L inner product: %w", err)
		}
		z_R, err := InnerProd(a_L, b_R)
		if err != nil {
			return IPAProof{}, fmt.Errorf("could not compute a_L*b_R inner product: %w", err)
		}

		C_L_1, err := commit(G_L, a_R)
		if err != nil {
			return IPAProof{}, fmt.Errorf("could not do G_L*a_R MSM: %w", err)
		}
		C_L, err := commit([]banderwagon.Element{C_L_1, q}, []fr.Element{fr.One(), z_L})
		if err != nil {
			return IPAProof{}, fmt.Errorf("could not do C_L_1+z_L*q MSM: %w", err)
		}

		C_R_1, err := commit(G_R, a_L)
		if err != nil {
			return IPAProof{}, fmt.Errorf("could not do G_R*a_L MSM: %w", err)
		}
		C_R, err := commit([]banderwagon.Element{C_R_1, q}, []fr.Element{fr.One(), z_R})
		if err != nil {
			return IPAProof{}, fmt.Errorf("could not do C_R_1+z_R*q MSM: %w", err)
		}

		L[i] = C_L
		R[i] = C_R

		transcript.AppendPoint(&C_L, labelL)
		transcript.AppendPoint(&C_R, labelR)
		x := transcript.ChallengeScalar(labelX)

		var xInv fr.Element
		xInv.Inverse(&x)

		// TODO: We could use a for loop here like in the Rust code
		a, err = foldScalars(a_L, a_R, x)
		if err != nil {
			return IPAProof{}, fmt.Errorf("could not fold a scalars a_L and a_R with x: %w", err)
		}
		b, err = foldScalars(b_L, b_R, xInv)
		if err != nil {
			return IPAProof{}, fmt.Errorf("could not fold b scalars b_L and b_R with xInv: %w", err)
		}

		current_basis, err = foldPoints(G_L, G_R, xInv)
		if err != nil {
			return IPAProof{}, fmt.Errorf("could not fold points G_L and G_R with xInv: %w", err)
		}

	}

	if len(a) != 1 {
		return IPAProof{}, fmt.Errorf("length of `a` should be 1 at the end of the reduction")
	}

	return IPAProof{
		L:        L,
		R:        R,
		A_scalar: a[0],
	}, nil
}

// Write serializes the IPA proof to the given writer.
func (ip *IPAProof) Write(w io.Writer) error {
	for _, el := range ip.L {
		if err := binary.Write(w, binary.BigEndian, el.Bytes()); err != nil {
			return fmt.Errorf("failed to write L: %w", err)
		}
	}
	for _, ar := range ip.R {
		if err := binary.Write(w, binary.BigEndian, ar.Bytes()); err != nil {
			return fmt.Errorf("failed to write R: %w", err)
		}
	}
	if err := binary.Write(w, binary.BigEndian, ip.A_scalar.BytesLE()); err != nil {
		return fmt.Errorf("failed to write A_scalar: %w", err)
	}
	return nil
}

// Read deserializes the IPA proof from the given reader.
func (ip *IPAProof) Read(r io.Reader) error {
	var L []banderwagon.Element
	for i := 0; i < 8; i++ {
		L_i, err := common.ReadPoint(r)
		if err != nil {
			return fmt.Errorf("failed to read L[%d]: %w", i, err)
		}
		L = append(L, *L_i)
	}
	ip.L = L
	var R []banderwagon.Element
	for i := 0; i < 8; i++ {
		R_i, err := common.ReadPoint(r)
		if err != nil {
			return fmt.Errorf("failed to read R[%d]: %w", i, err)
		}
		R = append(R, *R_i)
	}
	ip.R = R

	A_Scalar, err := common.ReadScalar(r)
	if err != nil {
		return fmt.Errorf("failed to read A_scalar: %w", err)
	}
	ip.A_scalar = *A_Scalar

	return nil
}

// Equal checks if two IPA proofs are equal.
func (ip IPAProof) Equal(other IPAProof) bool {
	num_rounds := 8
	if len(ip.L) != len(other.L) {
		return false
	}
	if len(ip.R) != len(other.R) {
		return false
	}
	if len(ip.L) != len(ip.R) {
		return false
	}
	if len(ip.L) != num_rounds {
		return false
	}

	for i := 0; i < num_rounds; i++ {
		expect_L_i := ip.L[i]
		expect_R_i := ip.R[i]

		got_L_i := other.L[i]
		got_R_i := other.R[i]

		if !expect_L_i.Equal(&got_L_i) {
			return false
		}
		if !expect_R_i.Equal(&got_R_i) {
			return false
		}
	}
	return ip.A_scalar.Equal(&other.A_scalar)
}

func computeBVector(ic *IPAConfig, evalPoint fr.Element) []fr.Element {
	if evalPoint.Cmp(&maxEvalPointInsideDomain) > 0 {
		return ic.PrecomputedWeights.ComputeBarycentricCoefficients(evalPoint)
	}
	// We build b = [0, 0, 0, ... , 1, .., 0] where the 1 element is at the index of the evaluation point.
	// This is correct since innerProductArgument(a, b) will return the evaluation of the polynomial at the
	// evaluation point in the domain.
	b := make([]fr.Element, common.VectorLength)
	var evalPointBI big.Int
	evalPoint.ToBigIntRegular(&evalPointBI)
	// Uint64() is safe because we checked that evalPoint is inside the domain (i.e <256).
	b[evalPointBI.Uint64()] = fr.One()

	return b
}
