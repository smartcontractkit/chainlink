package multiproof

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"runtime"

	"github.com/crate-crypto/go-ipa/bandersnatch/fr"
	"github.com/crate-crypto/go-ipa/banderwagon"
	"github.com/crate-crypto/go-ipa/common"
	"github.com/crate-crypto/go-ipa/ipa"
)

// The following are unexported labels to be used in Fiat-Shamir during the
// multiproof protocol.
//
// The following is a short description on how they're used in the protocol:
//  1. Append the domain separator. (labelDomainSep)
//  2. For each opening, we append to the transcript:
//     a. The polynomial commitment (labelC).
//     b. The evaluation point (labelZ).
//     c. The evaluation result (labelY).
//  3. Pull a scalar-field element from the transcript to be used for
//     the random linear combination of openings. (labelR)
//  4. Append point D which is sum(r^i * (f_i(x)-y_i)/(x-z_i)). (labelD)
//  5. Pull a random scalar-field to be used as a random evaluation point. (labelT)
//  5. Append point E which is sum(r^i * f_i(x)/(t-z_i)). (labelE)
//  7. Create the IPA proof for (E-D) at point `t`. See the `ipa` package for the FS description.
//
// Note: this package must not mutate these label values, nor pass them to
// parts of the code that would mutate them.
var (
	labelC         = []byte("C")
	labelZ         = []byte("z")
	labelY         = []byte("y")
	labelD         = []byte("D")
	labelE         = []byte("E")
	labelT         = []byte("t")
	labelR         = []byte("r")
	labelDomainSep = []byte("multiproof")
)

// MultiProof is a multi-proof for several polynomials in evaluation form.
type MultiProof struct {
	IPA ipa.IPAProof
	D   banderwagon.Element
}

// CreateMultiProof creates a multi-proof for several polynomials in evaluation form.
// The list of triplets (C, Fs, Z) represents each polynomial commitment, evaluations in the domain, and evaluation
// point respectively.
func CreateMultiProof(transcript *common.Transcript, ipaConf *ipa.IPAConfig, Cs []*banderwagon.Element, fs [][]fr.Element, zs []uint8) (*MultiProof, error) {
	transcript.DomainSep(labelDomainSep)

	for _, f := range fs {
		if len(f) != common.VectorLength {
			return nil, fmt.Errorf("polynomial length = %d, while expected length = %d", len(f), common.VectorLength)
		}
	}

	if len(Cs) != len(fs) {
		return nil, fmt.Errorf("number of commitments = %d, while number of functions = %d", len(Cs), len(fs))
	}
	if len(Cs) != len(zs) {
		return nil, fmt.Errorf("number of commitments = %d, while number of points = %d", len(Cs), len(zs))
	}

	num_queries := len(Cs)
	if num_queries == 0 {
		return nil, errors.New("cannot create a multiproof with 0 queries")
	}

	if err := banderwagon.BatchNormalize(Cs); err != nil {
		return nil, fmt.Errorf("could not batch normalize commitments: %w", err)
	}

	for i := 0; i < num_queries; i++ {
		transcript.AppendPoint(Cs[i], labelC)
		var z = domainToFr(zs[i])
		transcript.AppendScalar(&z, labelZ)

		// get the `y` value

		f := fs[i]
		y := f[zs[i]]
		transcript.AppendScalar(&y, labelY)
	}

	r := transcript.ChallengeScalar(labelR)
	powersOfR := common.PowersOf(r, num_queries)

	// Compute g(x)
	// We first compute the polynomials in lagrange form grouped by evaluation point, and
	// then we compute g(X). This limit the numbers of DivideOnDomain() calls up to
	// the domain size.
	groupedFs := groupPolynomialsByEvaluationPoint(fs, powersOfR, zs)

	g_x := make([]fr.Element, common.VectorLength)
	for index, f := range groupedFs {
		// If there is no polynomial for this evaluation point, we skip it.
		if len(f) == 0 {
			continue
		}
		quotient := ipaConf.PrecomputedWeights.DivideOnDomain(uint8(index), f)
		for j := 0; j < common.VectorLength; j++ {
			g_x[j].Add(&g_x[j], &quotient[j])
		}
	}

	D := ipaConf.Commit(g_x)

	transcript.AppendPoint(&D, labelD)
	t := transcript.ChallengeScalar(labelT)

	// Calculate the denominator inverses only for referenced evaluation points.
	den_inv := make([]fr.Element, 0, common.VectorLength)
	for z, f := range groupedFs {
		if len(f) == 0 {
			continue
		}
		var z = domainToFr(uint8(z))
		var den fr.Element
		den.Sub(&t, &z)
		den_inv = append(den_inv, den)
	}
	den_inv = fr.BatchInvert(den_inv)

	// Compute h(X) = g_1(X)
	h_x := make([]fr.Element, common.VectorLength)
	denInvIdx := 0
	for _, f := range groupedFs {
		if len(f) == 0 {
			continue
		}
		for k := 0; k < common.VectorLength; k++ {
			var tmp fr.Element
			tmp.Mul(&f[k], &den_inv[denInvIdx])
			h_x[k].Add(&h_x[k], &tmp)
		}
		denInvIdx++
	}

	h_minus_g := make([]fr.Element, common.VectorLength)
	for i := 0; i < common.VectorLength; i++ {
		h_minus_g[i].Sub(&h_x[i], &g_x[i])
	}

	E := ipaConf.Commit(h_x)
	transcript.AppendPoint(&E, labelE)

	var EminusD banderwagon.Element

	EminusD.Sub(&E, &D)

	ipaProof, err := ipa.CreateIPAProof(transcript, ipaConf, EminusD, h_minus_g, t)
	if err != nil {
		return nil, fmt.Errorf("could not create IPA proof: %w", err)
	}

	return &MultiProof{
		IPA: ipaProof,
		D:   D,
	}, nil
}

// CheckMultiProof verifies a multi-proof for several polynomials in evaluation form.
// The list of triplets (C, Y, Z) represents each polynomial commitment, evaluation
// result, and evaluation point in the domain.
func CheckMultiProof(transcript *common.Transcript, ipaConf *ipa.IPAConfig, proof *MultiProof, Cs []*banderwagon.Element, ys []*fr.Element, zs []uint8) (bool, error) {
	transcript.DomainSep(labelDomainSep)

	if len(Cs) != len(ys) {
		return false, fmt.Errorf("number of commitments = %d, while number of output points = %d", len(Cs), len(ys))
	}
	if len(Cs) != len(zs) {
		return false, fmt.Errorf("number of commitments = %d, while number of input points = %d", len(Cs), len(zs))
	}

	num_queries := len(Cs)
	if num_queries == 0 {
		return false, errors.New("number of queries is zero")
	}

	for i := 0; i < num_queries; i++ {
		transcript.AppendPoint(Cs[i], labelC)
		var z = domainToFr(zs[i])
		transcript.AppendScalar(&z, labelZ)
		transcript.AppendScalar(ys[i], labelY)
	}

	r := transcript.ChallengeScalar(labelR)
	powers_of_r := common.PowersOf(r, num_queries)

	transcript.AppendPoint(&proof.D, labelD)
	t := transcript.ChallengeScalar(labelT)

	// Compute the polynomials in lagrange form grouped by evaluation point, and
	// the needed helper scalars.
	groupedEvals := make([]fr.Element, common.VectorLength)
	for i := 0; i < num_queries; i++ {
		z := zs[i]

		// r * y_i
		r := powers_of_r[i]
		var scaledEvaluation fr.Element
		scaledEvaluation.Mul(&r, ys[i])
		groupedEvals[z].Add(&groupedEvals[z], &scaledEvaluation)
	}

	// Compute helper_scalar_den. This is 1 / t - z_i
	helper_scalar_den := make([]fr.Element, common.VectorLength)
	for i := 0; i < common.VectorLength; i++ {
		// (t - z_i)
		var z = domainToFr(uint8(i))
		helper_scalar_den[i].Sub(&t, &z)
	}
	helper_scalar_den = fr.BatchInvert(helper_scalar_den)

	// Compute g_2(t) = SUM (y_i * r^i) / (t - z_i) = SUM (y_i * r) * helper_scalars_den
	g_2_t := fr.Zero()
	for i := 0; i < common.VectorLength; i++ {
		if groupedEvals[i].IsZero() {
			continue
		}
		var tmp fr.Element
		tmp.Mul(&groupedEvals[i], &helper_scalar_den[i])
		g_2_t.Add(&g_2_t, &tmp)
	}

	// Compute E = SUM C_i * (r^i / t - z_i) = SUM C_i * msm_scalars
	msm_scalars := make([]fr.Element, len(Cs))
	Csnp := make([]banderwagon.Element, len(Cs))
	for i := 0; i < len(Cs); i++ {
		Csnp[i] = *Cs[i]

		msm_scalars[i].Mul(&powers_of_r[i], &helper_scalar_den[zs[i]])
	}

	E, err := ipa.MultiScalar(Csnp, msm_scalars)
	if err != nil {
		return false, fmt.Errorf("could not compute E: %w", err)
	}
	transcript.AppendPoint(&E, labelE)

	var E_minus_D banderwagon.Element
	E_minus_D.Sub(&E, &proof.D)

	ok, err := ipa.CheckIPAProof(transcript, ipaConf, E_minus_D, proof.IPA, t, g_2_t)
	if err != nil {
		return false, fmt.Errorf("could not check IPA proof: %w", err)
	}

	return ok, nil
}

func domainToFr(in uint8) fr.Element {
	var x fr.Element
	x.SetUint64(uint64(in))
	return x
}

// Write serializes a multi-proof to a writer.
func (mp *MultiProof) Write(w io.Writer) error {
	if err := binary.Write(w, binary.BigEndian, mp.D.Bytes()); err != nil {
		return fmt.Errorf("failed to write D: %w", err)
	}
	if err := mp.IPA.Write(w); err != nil {
		return fmt.Errorf("failed to write IPA proof: %w", err)
	}
	return nil
}

// Read deserializes a multi-proof from a reader.
func (mp *MultiProof) Read(r io.Reader) error {
	D, err := common.ReadPoint(r)
	if err != nil {
		return fmt.Errorf("failed to read D: %w", err)
	}
	mp.D = *D
	if err := mp.IPA.Read(r); err != nil {
		return fmt.Errorf("failed to read IPA proof: %w", err)
	}
	// Check that the next read is EOF.
	var buf [1]byte
	if _, err := r.Read(buf[:]); err != io.EOF {
		return errors.New("expected EOF")
	}

	return nil
}

// Equal checks if two multi-proofs are equal.
func (mp MultiProof) Equal(other MultiProof) bool {
	if !mp.IPA.Equal(other.IPA) {
		return false
	}
	return mp.D.Equal(&other.D)
}

func groupPolynomialsByEvaluationPoint(fs [][]fr.Element, powersOfR []fr.Element, zs []uint8) [common.VectorLength][]fr.Element {
	workersAggregations := make(chan [common.VectorLength][]fr.Element)

	numWorkers := runtime.NumCPU()
	batchSize := (len(fs) + numWorkers - 1) / numWorkers
	for i := 0; i < numWorkers; i++ {
		go func(start, end int) {
			if end > len(fs) {
				end = len(fs)
			}
			var groupedFs [common.VectorLength][]fr.Element
			for i := start; i < end; i++ {
				z := zs[i]
				if len(groupedFs[z]) == 0 {
					groupedFs[z] = make([]fr.Element, common.VectorLength)
				}

				for j := 0; j < common.VectorLength; j++ {
					var scaledEvaluation fr.Element
					scaledEvaluation.Mul(&powersOfR[i], &fs[i][j])
					groupedFs[z][j].Add(&groupedFs[z][j], &scaledEvaluation)
				}
			}
			workersAggregations <- groupedFs
		}(i*batchSize, (i+1)*batchSize)
	}

	// Each worker has computed its own aggregation. Now we aggregate the results.
	// This is bounded to reducing a `numWorkers` sized array of `common.VectorLength` sized arrays.
	var groupedFs [common.VectorLength][]fr.Element
	for i := 0; i < numWorkers; i++ {
		workerAggregation := <-workersAggregations
		for z := range workerAggregation {
			if len(workerAggregation[z]) == 0 {
				continue
			}
			// If this is the first time we see this evaluation point, we initialize it
			// reusing the worker result.
			if groupedFs[z] == nil {
				groupedFs[z] = workerAggregation[z]
				continue
			}
			// If not, we aggregate the worker result with the previous result for this evaluation.
			for j := 0; j < common.VectorLength; j++ {
				groupedFs[z][j].Add(&groupedFs[z][j], &workerAggregation[z][j])
			}
		}
	}

	return groupedFs
}
