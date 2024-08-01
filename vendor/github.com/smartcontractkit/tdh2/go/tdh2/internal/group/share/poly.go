// Package share implements Shamir secret sharing and polynomial commitments.
// Shamir's scheme allows you to split a secret value into multiple parts, so called
// shares, by evaluating a secret sharing polynomial at certain indices. The
// shared secret can only be reconstructed (via Lagrange interpolation) if a
// threshold of the participants provide their shares. A polynomial commitment
// scheme allows a committer to commit to a secret sharing polynomial so that
// a verifier can check the claimed evaluations of the committed polynomial.
// Both schemes of this package are core building blocks for more advanced
// secret sharing techniques.
package share

import (
	"crypto/cipher"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/smartcontractkit/tdh2/go/tdh2/internal/group"
)

// PriShare represents a private share.
type PriShare struct {
	I int          // Index of the private share
	V group.Scalar // Value of the private share
}

func (p *PriShare) String() string {
	return fmt.Sprintf("{%d:%s}", p.I, p.V)
}

// PriPoly represents a secret sharing polynomial.
type PriPoly struct {
	g      group.Group    // Cryptographic group
	coeffs []group.Scalar // Coefficients of the polynomial
}

// NewPriPoly creates a new secret sharing polynomial using the provided
// cryptographic group, the secret sharing threshold t, and the secret to be
// shared s. If s is nil, a new s is chosen using the provided randomness
// stream rand.
func NewPriPoly(grp group.Group, t int, s group.Scalar, rand cipher.Stream) *PriPoly {
	coeffs := make([]group.Scalar, t)
	coeffs[0] = s
	if coeffs[0] == nil {
		coeffs[0] = grp.Scalar().Pick(rand)
	}
	for i := 1; i < t; i++ {
		coeffs[i] = grp.Scalar().Pick(rand)
	}
	return &PriPoly{g: grp, coeffs: coeffs}
}

// Secret returns the shared secret p(0), i.e., the constant term of the polynomial.
func (p *PriPoly) Secret() group.Scalar {
	return p.coeffs[0]
}

// Eval computes the private share v = p(i).
func (p *PriPoly) Eval(i int) *PriShare {
	xi := p.g.Scalar().SetInt64(1 + int64(i))
	v := p.g.Scalar().Zero()
	for j := len(p.coeffs) - 1; j >= 0; j-- {
		v.Mul(v, xi)
		v.Add(v, p.coeffs[j])
	}
	return &PriShare{i, v}
}

// Shares creates a list of n private shares p(1),...,p(n).
func (p *PriPoly) Shares(n int) []*PriShare {
	shares := make([]*PriShare, n)
	for i := range shares {
		shares[i] = p.Eval(i)
	}
	return shares
}

func (p *PriPoly) String() string {
	var strs = make([]string, len(p.coeffs))
	for i, c := range p.coeffs {
		strs[i] = c.String()
	}
	return "[ " + strings.Join(strs, ", ") + " ]"
}

// PubShare represents a public share.
type PubShare struct {
	I int         // Index of the public share
	V group.Point // Value of the public share
}

// xyCommits is the public version of xScalars.
func xyCommit(g group.Group, shares []*PubShare, t, n int) (map[int]group.Scalar, map[int]group.Point) {
	// we are sorting first the shares since the shares may be unrelated for
	// some applications. In this case, all participants needs to interpolate on
	// the exact same order shares.
	sorted := make([]*PubShare, 0, n)
	for _, share := range shares {
		if share != nil {
			sorted = append(sorted, share)
		}
	}
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].I < sorted[j].I })

	x := make(map[int]group.Scalar)
	y := make(map[int]group.Point)

	for _, s := range sorted {
		if s == nil || s.V == nil || s.I < 0 {
			continue
		}
		idx := s.I
		x[idx] = g.Scalar().SetInt64(int64(idx + 1))
		y[idx] = s.V
		if len(x) == t {
			break
		}
	}
	return x, y
}

// RecoverCommit reconstructs the secret commitment p(0) from a list of public
// shares using Lagrange interpolation.
func RecoverCommit(g group.Group, shares []*PubShare, t, n int) (group.Point, error) {
	x, y := xyCommit(g, shares, t, n)
	if len(x) < t {
		return nil, errors.New("share: not enough good public shares to reconstruct secret commitment")
	}

	num := g.Scalar()
	den := g.Scalar()
	tmp := g.Scalar()
	Acc := g.Point().Null()
	Tmp := g.Point()

	for i, xi := range x {
		num.One()
		den.One()
		for j, xj := range x {
			if i == j {
				continue
			}
			num.Mul(num, xj)
			den.Mul(den, tmp.Sub(xj, xi))
		}
		Tmp.Mul(num.Div(num, den), y[i])
		Acc.Add(Acc, Tmp)
	}

	return Acc, nil
}
