// Package curve25519 contains several implementations of Twisted Edwards Curves,
// from general and unoptimized to highly specialized and optimized.
//
// Twisted Edwards curves are elliptic curves satisfying the equation:
//
//	ax^2 + y^2 = c^2(1 + dx^2y^2)
//
// for some scalars c, d over some field K. We assume K is a (finite) prime field for a
// large prime p. We also assume c == 1 because all curves in the generalized form
// are isomorphic to curves having c == 1.
//
// For details see Bernstein et al, "Twisted Edwards Curves", http://eprint.iacr.org/2008/013.pdf
package curve25519

import (
	"math/big"

	"go.dedis.ch/kyber/v3/group/mod"
)

// Param defines a Twisted Edwards curve (TEC).
type Param struct {
	Name string // Name of curve

	P big.Int // Prime defining the underlying field
	Q big.Int // Order of the prime-order base point
	R int     // Cofactor: Q*R is the total size of the curve

	A, D big.Int // Edwards curve equation parameters

	FBX, FBY big.Int // Standard base point for full group
	PBX, PBY big.Int // Standard base point for prime-order subgroup

	Elligator1s big.Int // Optional s parameter for Elligator 1
	Elligator2u big.Int // Optional u parameter for Elligator 2
}

// Return the name of this curve.
func (p *Param) String() string {
	return p.Name
}

// Param1174 defines Curve1174, as specified in:
// Bernstein et al, "Elligator: Elliptic-curve points indistinguishable
// from uniform random strings"
// http://elligator.cr.yp.to/elligator-20130828.pdf
//
func Param1174() *Param {
	var p Param
	var mi mod.Int

	p.Name = "Curve1174"
	p.P.SetBit(zero, 251, 1).Sub(&p.P, big.NewInt(9))
	p.Q.SetString("45330879683285730139092453152713398835", 10)
	p.Q.Sub(&p.P, &p.Q).Div(&p.Q, big.NewInt(4))
	p.R = 4
	p.A.SetInt64(1)
	p.D.SetInt64(-1174)

	// Full-group generator is (4/V,3/5)
	mi.InitString("4", "19225777642111670230408712442205514783403012708409058383774613284963344096", 10, &p.P)
	p.FBX.Set(&mi.V)
	mi.InitString("3", "5", 10, &p.P)
	p.FBY.Set(&mi.V)

	// Elligator1 parameter s for Curve1174 (Elligator paper section 4.1)
	p.Elligator1s.SetString("1806494121122717992522804053500797229648438766985538871240722010849934886421", 10)

	return &p
}

// Param25519 defines the Edwards version of Curve25519, as specified in:
// Bernstein et al, "High-speed high-security signatures",
// http://ed25519.cr.yp.to/ed25519-20110926.pdf
//
func Param25519() *Param {
	var p Param
	var qs big.Int
	p.Name = "Curve25519"
	p.P.SetBit(zero, 255, 1).Sub(&p.P, big.NewInt(19))
	qs.SetString("27742317777372353535851937790883648493", 10)
	p.Q.SetBit(zero, 252, 1).Add(&p.Q, &qs)
	p.R = 8
	p.A.SetInt64(-1).Add(&p.P, &p.A)
	p.D.SetString("37095705934669439343138083508754565189542113879843219016388785533085940283555", 10)

	p.PBX.SetString("15112221349535400772501151409588531511454012693041857206046113283949847762202", 10)
	p.PBY.SetString("46316835694926478169428394003475163141307993866256225615783033603165251855960", 10)

	// Non-square u for Elligator2
	p.Elligator2u.SetInt64(2)

	return &p
}

// ParamE382 defines the E-382 curve specified in:
// Aranha et al, "A note on high-security general-purpose elliptic curves",
// http://eprint.iacr.org/2013/647.pdf
//
// and more recently in:
//
// "Additional Elliptic Curves for IETF protocols"
// http://tools.ietf.org/html/draft-ladd-safecurves-02
// (this I-D is now expired)
func ParamE382() *Param {
	var p Param
	var qs big.Int
	p.Name = "E-382"
	p.P.SetBit(zero, 382, 1).Sub(&p.P, big.NewInt(105)) // p = 2^382-105
	qs.SetString("1030303207694556153926491950732314247062623204330168346855", 10)
	p.Q.SetBit(zero, 380, 1).Sub(&p.Q, &qs)
	p.R = 8
	p.A.SetInt64(1)
	p.D.SetInt64(-67254)
	p.PBX.SetString("3914921414754292646847594472454013487047137431784830634731377862923477302047857640522480241298429278603678181725699", 10)
	p.PBY.SetString("17", 10)
	return &p
}

// Param41417 defines the Curve41417 curve, as specified in:
// Bernstein et al, "Curve41417: Karatsuba revisited",
// http://eprint.iacr.org/2014/526.pdf
func Param41417() *Param {
	var p Param
	var qs big.Int
	p.Name = "Curve41417"
	p.P.SetBit(zero, 414, 1).Sub(&p.P, big.NewInt(17))
	qs.SetString("33364140863755142520810177694098385178984727200411208589594759", 10)
	p.Q.SetBit(zero, 411, 1).Sub(&p.Q, &qs)
	p.R = 8
	p.A.SetInt64(1)
	p.D.SetInt64(3617)
	p.PBX.SetString("17319886477121189177719202498822615443556957307604340815256226171904769976866975908866528699294134494857887698432266169206165", 10)
	p.PBY.SetString("34", 10)
	return &p
}

// ParamE521 defines the E-521 curve specified in:
// Aranha et al, "A note on high-security general-purpose elliptic curves",
// http://eprint.iacr.org/2013/647.pdf
//
// and more recently included in:
// "Additional Elliptic Curves for IETF protocols"
// http://tools.ietf.org/html/draft-ladd-safecurves-02
//
func ParamE521() *Param {
	var p Param
	var qs big.Int
	p.Name = "E-521"
	p.P.SetBit(zero, 521, 1).Sub(&p.P, one)
	qs.SetString("337554763258501705789107630418782636071904961214051226618635150085779108655765", 10)
	p.Q.SetBit(zero, 519, 1).Sub(&p.Q, &qs)
	p.R = 8
	p.A.SetInt64(1)
	p.D.SetInt64(-376014)
	p.PBX.SetString("1571054894184995387535939749894317568645297350402905821437625181152304994381188529632591196067604100772673927915114267193389905003276673749012051148356041324", 10)
	p.PBY.SetString("12", 10)
	return &p
}
