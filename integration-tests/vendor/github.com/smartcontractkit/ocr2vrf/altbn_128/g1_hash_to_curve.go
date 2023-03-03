package altbn_128

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/smartcontractkit/ocr2vrf/gethwrappers/vrf"
	"github.com/smartcontractkit/ocr2vrf/internal/util"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"

	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/group/mod"
)

var (
	curve  = (&PairingSuite{}).G1()
	p      = bn256.P
	curveB = m(3)

	j            = big.NewInt
	uint256Bound = j(0).Lsh(j(1), 256)
	uint256ModP  = j(0).Mod(uint256Bound, p)
	maxP         = j(0).Sub(uint256Bound, uint256ModP)

	sqrMinus3, sqrMinus3Valid = m(0), sqrMinus3.Sqrt(m(-3))
	thirdModP                 = m(0).Exp(m(3), j(-1))
	sqrpwr                    = j(0).Rsh(j(0).Add(p, j(1)), 2)

	vConst, _ = m(0), vConst.Div(vConst.Add(m(-1), sqrMinus3), m(2))
)

type fProof struct {
	point         kyber.Point
	t             *big.Int
	interimValues vrf.HashToCurveFProof
}

var zeroHashPoint = hexutil.MustDecode(
	"0x000000000000000059e26bcea0d48bacd4f263f1acdb5c4f5763473177fffffe",
)

func newFProof(t *big.Int) *fProof {
	if t.Cmp(j(0)) == 0 {
		return &fProof{point: curve.Point().Null(), t: t}
	}
	if t.Cmp(maxP) >= 0 {

		panic("input must be less than maxP")
	}
	it := mod.NewInt(t, p)
	tmp := m(0)
	iTSquared := m(0).Mul(it, it)
	denom := tmp.Add(m(1), tmp.Add(curveB, iTSquared))
	denomInv := m(0).Inv(denom)
	tmp = m(0)
	v := tmp.Sub(vConst, tmp.Mul(tmp.Mul(sqrMinus3, iTSquared), denomInv))
	pseudoSqrtY := func(x *mod.Int) (isActualSquareRoot bool, psqrt *mod.Int) {
		tmp := m(0)
		ySquare := tmp.Add(tmp.Mul(tmp.Mul(x, x), x), curveB).(*mod.Int)

		pseudoY := m(0).Exp(ySquare, sqrpwr)
		pseudoYSquare := m(0).Mul(pseudoY, pseudoY)
		valid := pseudoYSquare.Equal(ySquare)
		if !valid {
			negPseudoYSquare := m(0).Neg(pseudoYSquare)
			if !negPseudoYSquare.Equal(ySquare) {
				panic(fmt.Sprintln("failed to compute correct pseudo square root of", ySquare))
			}
		}
		return valid, pseudoY.(*mod.Int)
	}
	rv := &fProof{
		t: t,
		interimValues: vrf.HashToCurveFProof{
			DenomInv: &denomInv.(*mod.Int).V,

			TInvSquared: j(0), Y1: j(0), Y2: j(0), Y3: j(0),
		},
	}

	x1 := v.Clone().(*mod.Int)
	valid, y1 := pseudoSqrtY(x1)
	rv.interimValues.Y1 = &y1.V
	if valid {
		rv.point = coordinatesToG1(x1, y1, it)
		return rv
	}

	x2 := m(0).Neg(m(0).Add(x1, m(1))).(*mod.Int)
	valid, y2 := pseudoSqrtY(x2)
	rv.interimValues.Y2 = &y2.V
	if valid {
		rv.point = coordinatesToG1(x2, y2, it)
		return rv
	}

	tInvSquared := m(0).Exp(it, j(-2)).(*mod.Int)
	numSquared := m(0).Mul(denom, denom)
	x3 := m(0).Sub(m(1), m(0).Mul(numSquared, m(0).Mul(tInvSquared, thirdModP))).(*mod.Int)
	valid, y3 := pseudoSqrtY(x3)
	rv.interimValues.Y3 = &y3.V
	rv.interimValues.TInvSquared = &tInvSquared.V
	if valid {
		rv.point = coordinatesToG1(x3, y3, it)
		return rv
	}

	panic(
		"one of x1, x2, x3 should have been the x ordinate of a point on G1: " +
			rv.String(),
	)
}

func coordinatesToG1(x, y, t *mod.Int) *g1Point {
	if x.M.Cmp(p) != 0 || y.M.Cmp(p) != 0 || t.M.Cmp(p) != 0 {
		panic("inputs are not in base field")
	}
	xBin, err := x.MarshalBinary()
	if err != nil {
		panic(err)
	}
	targetParity := t.V.Bit(0)
	if y.V.Bit(0) != targetParity {
		_ = y.Neg(y)
		if y.V.Bit(0) != targetParity {
			panic("failed to set target parity for output")
		}
	}
	yBin, err := y.MarshalBinary()
	if err != nil {
		panic(err)
	}
	pt := newG1Point()
	if _, err := pt.G1.Unmarshal(append(xBin, yBin...)); err != nil {
		panic(err)
	}
	return pt
}

func CoordinatesToG1(x, y *mod.Int) (kyber.Point, error) {
	if x.M.Cmp(p) != 0 {
		return nil, fmt.Errorf("x ordinate %+v is not in base field", x)
	}
	if y.M.Cmp(p) != 0 {
		return nil, fmt.Errorf("y ordinate %+v is not in base field", y)
	}
	xBin, err := x.MarshalBinary()
	if err != nil {
		return nil, util.WrapError(err, "could not marshal x ordinate")
	}
	yBin, err := y.MarshalBinary()
	if err != nil {
		return nil, util.WrapError(err, "could not marshal y ordinate")
	}
	pt := newG1Point()
	combinedOrdinates := append(xBin, yBin...)
	if _, err := pt.G1.Unmarshal(combinedOrdinates); err != nil {
		return nil, util.WrapErrorf(
			err,
			"could not unmarshal combined ordinates 0x%x",
			combinedOrdinates,
		)
	}
	return pt, nil
}

type HashProof struct {
	msg           [32]byte
	HashPoint     kyber.Point
	SummandProofs [2]*fProof
}

func NewHashProof(msg [32]byte) *HashProof {
	rv := &HashProof{msg: msg, HashPoint: newG1Point().Null()}
	for hashCount := 0; hashCount < 2; {
		nmsg := crypto.Keccak256(msg[:])
		copy(msg[:], nmsg)
		t := j(0).SetBytes(msg[:])
		if t.Cmp(maxP) < 0 {
			t = t.Mod(t, p)
			rv.SummandProofs[hashCount] = newFProof(t)
			rv.HashPoint = newG1Point().Add(rv.HashPoint, rv.SummandProofs[hashCount].point)
			hashCount++
		}
	}
	return rv
}

func init() {
	if !sqrMinus3Valid {
		panic("-3 is not a square in ℤ/pℤ")
	}
	if j(0).Lsh(sqrpwr, 2).Cmp(j(0).Add(p, j(1))) != 0 {
		panic("p ≢ 3 mod 4")
	}
}

func SolidityVRFProof(pubKey, output kyber.Point, hp HashProof) (*vrf.VRFProof, error) {
	pubKeyB, err := pubKey.MarshalBinary()
	if err != nil {
		return nil, err
	}
	var pubKeySol vrf.ECCArithmeticG2Point
	for i := range pubKeySol.P {
		pubKeySol.P[i] = big.NewInt(0).SetBytes(pubKeyB[i*32 : (i+1)*32])
	}
	outputB := LongMarshal(output)
	var outputSol vrf.ECCArithmeticG1Point
	for i := range outputSol.P {
		outputSol.P[i] = big.NewInt(0).SetBytes(outputB[i*32 : (i+1)*32])
	}
	proof := vrf.VRFProof{
		pubKeySol,
		outputSol,
		hp.SummandProofs[0].interimValues,
		hp.SummandProofs[1].interimValues,
	}
	return &proof, nil
}

func (hp *HashProof) EqualFProofs(f1, f2 vrf.HashToCurveFProof) bool {
	s := hp.SummandProofs
	return f1.DenomInv.Cmp(s[0].interimValues.DenomInv) == 0 &&
		f1.TInvSquared.Cmp(s[0].interimValues.TInvSquared) == 0 &&
		f1.Y1.Cmp(s[0].interimValues.Y1) == 0 &&
		f1.Y2.Cmp(s[0].interimValues.Y2) == 0 &&
		f1.Y3.Cmp(s[0].interimValues.Y3) == 0 &&
		f2.DenomInv.Cmp(s[1].interimValues.DenomInv) == 0 &&
		f2.TInvSquared.Cmp(s[1].interimValues.TInvSquared) == 0 &&
		f2.Y1.Cmp(s[1].interimValues.Y1) == 0 &&
		f2.Y2.Cmp(s[1].interimValues.Y2) == 0 &&
		f2.Y3.Cmp(s[1].interimValues.Y3) == 0
}

func (f *fProof) String() string {
	i := f.interimValues
	return fmt.Sprintf(
		"&fProof{point: %s, t: 0x%x, interimValues: "+
			"{DenomInv: 0x%x, TInvSquared: 0x%x, Y1: 0x%x, Y2: 0x%x, Y3: 0x%x}"+
			"}",
		f.point, f.t, i.DenomInv, i.TInvSquared, i.Y1, i.Y2, i.Y3,
	)
}
