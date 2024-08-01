package curve

/*
	Although the library adheres to the 'elliptic/curve' interface.
	All testing has been done against library function explicity.
	It is recommended to use in the same way(i.e. `curve.Sign` and not `ecdsa.Sign`).
*/
import (
	"bytes"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"math/big"

	junoCrypto "github.com/NethermindEth/juno/core/crypto"
	"github.com/NethermindEth/juno/core/felt"
)

var Curve StarkCurve

/*
Returned stark curve includes several values above and beyond
what the 'elliptic' interface calls for to facilitate common starkware functions
*/
type StarkCurve struct {
	*elliptic.CurveParams
	EcGenX           *big.Int
	EcGenY           *big.Int
	MinusShiftPointX *big.Int
	MinusShiftPointY *big.Int
	Max              *big.Int
	Alpha            *big.Int
	ConstantPoints   [][]*big.Int
}

//go:embed pedersen_params.json
var PedersenParamsRaw []byte
var PedersenParams StarkCurvePayload

// struct definition for parsing 'pedersen_params.json'
type StarkCurvePayload struct {
	License        []string     `json:"_license"`
	Comment        string       `json:"_comment"`
	FieldPrime     *big.Int     `json:"FIELD_PRIME"`
	FieldGen       int          `json:"FIELD_GEN"`
	EcOrder        *big.Int     `json:"EC_ORDER"`
	Alpha          int64        `json:"ALPHA"`
	Beta           *big.Int     `json:"BETA"`
	ConstantPoints [][]*big.Int `json:"CONSTANT_POINTS"`
}

// init initializes the PedersenParams and Curve variables.
//
// It unmarshals the PedersenParamsRaw JSON data into the PedersenParams struct.
// If there is an error during unmarshalling, it will log a fatal error.
//
// It checks the length of the ConstantPoints field in PedersenParams. If the length is 0,
// it will panic with the message "decoding pedersen params json".
//
// It sets the CurveParams field of the Curve variable to a new elliptic.CurveParams with the name "stark-curve-with-constants".
// It sets the P, N, B, Gx, Gy, EcGenX, EcGenY, MinusShiftPointX, MinusShiftPointY, Max, Alpha, and BitSize fields of the Curve variable
// with the corresponding values from the PedersenParams struct.
//
// After that, it overrides the CurveParams field of the Curve variable with a new elliptic.CurveParams with the name "stark-curve".
// It sets the P, N, B, Gx, Gy, EcGenX, EcGenY, MinusShiftPointX, MinusShiftPointY, Max, Alpha, and BitSize fields of the Curve variable
// with the corresponding values from the PedersenParams struct.
//
// Note: Not all operations require a stark curve initialization including the provided constant points.
// This function can be used to initialize the curve without the constant points.
// 
// Parameters:
//  none
// Returns:
//  none
func init() {
	if err := json.Unmarshal(PedersenParamsRaw, &PedersenParams); err != nil {
		log.Fatalf("unmarshalling pedersen params: %v", err)
	}

	if len(PedersenParams.ConstantPoints) == 0 {
		panic("decoding pedersen params json")
	}

	Curve.CurveParams = &elliptic.CurveParams{Name: "stark-curve-with-constants"}

	Curve.P = PedersenParams.FieldPrime
	Curve.N = PedersenParams.EcOrder
	Curve.B = PedersenParams.Beta
	Curve.Gx = PedersenParams.ConstantPoints[0][0]
	Curve.Gy = PedersenParams.ConstantPoints[0][1]
	Curve.EcGenX = PedersenParams.ConstantPoints[1][0]
	Curve.EcGenY = PedersenParams.ConstantPoints[1][1]
	Curve.MinusShiftPointX, _ = new(big.Int).SetString("2089986280348253421170679821480865132823066470938446095505822317253594081284", 10) // MINUS_SHIFT_POINT = (SHIFT_POINT[0], FIELD_PRIME - SHIFT_POINT[1])
	Curve.MinusShiftPointY, _ = new(big.Int).SetString("1904571459125470836673916673895659690812401348070794621786009710606664325495", 10)
	Curve.Max, _ = new(big.Int).SetString("3618502788666131106986593281521497120414687020801267626233049500247285301248", 10) // 2 ** 251
	Curve.Alpha = big.NewInt(PedersenParams.Alpha)
	Curve.BitSize = 252
	Curve.ConstantPoints = PedersenParams.ConstantPoints

	/*
		Not all operations require a stark curve initialization
		including the provided constant points. Here you can
		initialize the curve without the constant points
	*/
	Curve.CurveParams = &elliptic.CurveParams{Name: "stark-curve"}
	Curve.P, _ = new(big.Int).SetString("3618502788666131213697322783095070105623107215331596699973092056135872020481", 10)  // Field Prime ./pedersen_json
	Curve.N, _ = new(big.Int).SetString("3618502788666131213697322783095070105526743751716087489154079457884512865583", 10)  // Order of base point ./pedersen_json
	Curve.B, _ = new(big.Int).SetString("3141592653589793238462643383279502884197169399375105820974944592307816406665", 10)  // Constant of curve equation ./pedersen_json
	Curve.Gx, _ = new(big.Int).SetString("2089986280348253421170679821480865132823066470938446095505822317253594081284", 10) // (x, _) of basepoint ./pedersen_json
	Curve.Gy, _ = new(big.Int).SetString("1713931329540660377023406109199410414810705867260802078187082345529207694986", 10) // (_, y) of basepoint ./pedersen_json
	Curve.EcGenX, _ = new(big.Int).SetString("874739451078007766457464989774322083649278607533249481151382481072868806602", 10)
	Curve.EcGenY, _ = new(big.Int).SetString("152666792071518830868575557812948353041420400780739481342941381225525861407", 10)
	Curve.MinusShiftPointX, _ = new(big.Int).SetString("2089986280348253421170679821480865132823066470938446095505822317253594081284", 10) // MINUS_SHIFT_POINT = (SHIFT_POINT[0], FIELD_PRIME - SHIFT_POINT[1])
	Curve.MinusShiftPointY, _ = new(big.Int).SetString("1904571459125470836673916673895659690812401348070794621786009710606664325495", 10) // MINUS_SHIFT_POINT = (SHIFT_POINT[0], FIELD_PRIME - SHIFT_POINT[1])
	Curve.Max, _ = new(big.Int).SetString("3618502788666131106986593281521497120414687020801267626233049500247285301248", 10)              // 2 ** 251
	Curve.Alpha = big.NewInt(1)
	Curve.BitSize = 252
}

// Add computes the sum of two points on the StarkCurve.
// Assumes affine form (x, y) is spread (x1 *big.Int, y1 *big.Int)
// (ref: https://github.com/starkware-libs/cairo-lang/blob/master/src/starkware/crypto/signature/math_utils.py#L59)
//
// Parameters:
// - x1, y1: The coordinates of the first point as pointers to big.Int on the curve
// - x2, y2: The coordinates of the second point as pointers to big.Int on the curve
// Returns:
// - x, y: two pointers to big.Int, representing the x and y coordinates of the sum of the two input points
func (sc StarkCurve) Add(x1, y1, x2, y2 *big.Int) (x, y *big.Int) {
	yDelta := new(big.Int).Sub(y1, y2)
	xDelta := new(big.Int).Sub(x1, x2)

	m := DivMod(yDelta, xDelta, sc.P)

	xm := new(big.Int).Mul(m, m)

	x = new(big.Int).Sub(xm, x1)
	x = x.Sub(x, x2)
	x = x.Mod(x, sc.P)

	y = new(big.Int).Sub(x1, x)
	y = y.Mul(m, y)
	y = y.Sub(y, y1)
	y = y.Mod(y, sc.P)

	return x, y
}

// Double calculates the double of a point on a StarkCurve (equation y^2 = x^3 + alpha*x + beta mod p).
// Assumes affine form (x, y) is spread (x1 *big.Int, y1 *big.Int)
// (ref: https://github.com/starkware-libs/cairo-lang/blob/master/src/starkware/crypto/signature/math_utils.py#L79)
//
// The function takes two pointers to big.Int values, x1 and y1, which represent the
// coordinates of the point to be doubled on the StarkCurve. It returns two pointers
// to big.Int values, x and y, which represent the coordinates of the resulting point
// after the doubling operation.
//
// Parameters:
// - x1, y1: The coordinates of the point to be doubled on the StarkCurve.
// Returns:
// - x, y: two pointers to big.Int, representing the x and y coordinates of the resulting point
func (sc StarkCurve) Double(x1, y1 *big.Int) (x, y *big.Int) {
	xin := new(big.Int).Mul(big.NewInt(3), x1)
	xin = xin.Mul(xin, x1)
	xin = xin.Add(xin, sc.Alpha)

	yin := new(big.Int).Mul(y1, big.NewInt(2))

	m := DivMod(xin, yin, sc.P)

	xout := new(big.Int).Mul(m, m)
	xmed := new(big.Int).Mul(big.NewInt(2), x1)
	xout = xout.Sub(xout, xmed)
	xout = xout.Mod(xout, sc.P)

	yout := new(big.Int).Sub(x1, xout)
	yout = yout.Mul(m, yout)
	yout = yout.Sub(yout, y1)
	yout = yout.Mod(yout, sc.P)

	return xout, yout
}

// ScalarMult performs scalar multiplication on a point (x1, y1) with a scalar value k.
//
// Parameters:
// - x1: The x-coordinate of the point to be multiplied.
// - y1: The y-coordinate of the point to be multiplied.
// - k: The scalar value to multiply the point with.
// Returns:
// - x: The x-coordinate of the resulting point.
// - y: The y-coordinate of the resulting point.
func (sc StarkCurve) ScalarMult(x1, y1 *big.Int, k []byte) (x, y *big.Int) {
	m := new(big.Int).SetBytes(k)
	x, y = sc.EcMult(m, x1, y1)
	return x, y
}

// ScalarBaseMult returns the result of multiplying the base point of the StarkCurve 
// by the given scalar value.
//
// Parameters:
// - k: The scalar value to multiply the base point by
// Returns:
// - x: The x-coordinate of the resulting point
// - y: The y-coordinate of the resulting point
func (sc StarkCurve) ScalarBaseMult(k []byte) (x, y *big.Int) {
	return sc.ScalarMult(sc.Gx, sc.Gy, k)
}

// IsOnCurve checks if the given point (x, y) lies on the curve defined by the StarkCurve instance.
//
// Parameters:
// - x: the x-coordinate of the point
// - y: the y-coordinate of the point
// Return type: bool
func (sc StarkCurve) IsOnCurve(x, y *big.Int) bool {
	left := new(big.Int).Mul(y, y)
	left = left.Mod(left, sc.P)

	right := new(big.Int).Mul(x, x)
	right = right.Mul(right, x)
	right = right.Mod(right, sc.P)

	ri := new(big.Int).Mul(big.NewInt(1), x)

	right = right.Add(right, ri)
	right = right.Add(right, sc.B)
	right = right.Mod(right, sc.P)

	if left.Cmp(right) == 0 {
		return true
	} else {
		return false
	}
}

// InvModCurveSize calculates the inverse modulus of a given big integer 'x' with respect to the StarkCurve 'sc'.
// (ref: https://github.com/starkware-libs/cairo-lang/blob/master/src/starkware/crypto/signature/math_utils.py)
//
// Parameters:
// - x: The big integer to calculate the inverse modulus for
// Returns:
// - The inverse modulus of 'x' with respect to 'sc.N'
func (sc StarkCurve) InvModCurveSize(x *big.Int) *big.Int {
	return DivMod(big.NewInt(1), x, sc.N)
}

// GetYCoordinate calculates the y-coordinate of a point on the StarkCurve.
// (ref: https://github.com/starkware-libs/cairo-lang/blob/master/src/starkware/crypto/signature/signature.py#L84)
// point (x,y) is on the curve.
// Note: the real y coordinate is either y or -y.
//
// Parameters:
// - starkX: The x-coordinate of the point
// Returns:
// - *big.Int: The calculated y-coordinate of the point
// a possible y coordinate such that together the point (x,y) is on the curve
// Note: the real y coordinate is either y or -y
func (sc StarkCurve) GetYCoordinate(starkX *big.Int) *big.Int {
	y := new(big.Int).Mul(starkX, starkX)
	y = y.Mul(y, starkX)
	yin := new(big.Int).Mul(sc.Alpha, starkX)

	y = y.Add(y, yin)
	y = y.Add(y, sc.B)
	y = y.Mod(y, sc.P)

	y = y.ModSqrt(y, sc.P)
	return y
}

// MimicEcMultAir performs a computation on the StarkCurve struct (m * point + shift_point)
// using the same steps like the AIR.
// (ref: https://github.com/starkware-libs/cairo-lang/blob/master/src/starkware/crypto/signature/signature.py#L176)
// AIR : Algebraic Intermediate Representation of computation
//
// Parameters:
// - mout: a pointer to a big.Int variable
// - x1, y1: a pointer to a big.Int point on the curve
// - x2, y2: a pointer to a big.Int point on the curve
// Returns:
// - x, y: a pointer to a big.Int point on the curve
// - err: an error if any
func (sc StarkCurve) MimicEcMultAir(mout, x1, y1, x2, y2 *big.Int) (x *big.Int, y *big.Int, err error) {
	m := new(big.Int).Set(mout)
	if m.Cmp(big.NewInt(0)) != 1 || m.Cmp(sc.Max) != -1 {
		return x, y, fmt.Errorf("too many bits %v", m.BitLen())
	}

	psx := x2
	psy := y2
	for i := 0; i < 251; i++ {
		if psx == x1 {
			return x, y, fmt.Errorf("xs are the same")
		}
		if m.Bit(0) == 1 {
			psx, psy = sc.Add(psx, psy, x1, y1)
		}
		x1, y1 = sc.Double(x1, y1)
		m = m.Rsh(m, 1)
	}
	if m.Cmp(big.NewInt(0)) != 0 {
		return psx, psy, fmt.Errorf("m doesn't equal zero")
	}
	return psx, psy, nil
}

// EcMult multiplies a point (equation y^2 = x^3 + alpha*x + beta mod p) on the StarkCurve by a scalar value.
// Assumes affine form (x, y) is spread (x1 *big.Int, y1 *big.Int) and that 0 < m < order(point).
// (ref: https://github.com/starkware-libs/cairo-lang/blob/master/src/starkware/crypto/signature/math_utils.py#L91)
//
// Parameters:
// - m: The scalar value to multiply the point by.
// - x1, y1: The coordinates of the point on the curve.
// Returns:
// - x, y: The coordinates of the resulting point after multiplication.
func (sc StarkCurve) EcMult(m, x1, y1 *big.Int) (x, y *big.Int) {
	var _ecMult func(m, x1, y1 *big.Int) (x, y *big.Int)

	_add := func(x1, y1, x2, y2 *big.Int) (x, y *big.Int) {
		yDelta := new(big.Int).Sub(y1, y2)
		xDelta := new(big.Int).Sub(x1, x2)

		m := DivMod(yDelta, xDelta, sc.P)

		xm := new(big.Int).Mul(m, m)

		x = new(big.Int).Sub(xm, x1)
		x = x.Sub(x, x2)
		x = x.Mod(x, sc.P)

		y = new(big.Int).Sub(x1, x)
		y = y.Mul(m, y)
		y = y.Sub(y, y1)
		y = y.Mod(y, sc.P)

		return x, y
	}

	// alpha is our Y
	_ecMult = func(m, x1, y1 *big.Int) (x, y *big.Int) {
		if m.BitLen() == 1 {
			return x1, y1
		}
		mk := new(big.Int).Mod(m, big.NewInt(2))
		if mk.Cmp(big.NewInt(0)) == 0 {
			h := new(big.Int).Div(m, big.NewInt(2))
			c, d := sc.Double(x1, y1)
			return _ecMult(h, c, d)
		}
		n := new(big.Int).Sub(m, big.NewInt(1))
		e, f := _ecMult(n, x1, y1)
		return _add(e, f, x1, y1)
	}

	x, y = _ecMult(m, x1, y1)
	return x, y
}

// Verify verifies the validity of the signature for a given message hash using the StarkCurve.
// (ref: https://github.com/starkware-libs/cairo-lang/blob/master/src/starkware/crypto/signature/signature.py#L217)
//
// Parameters:
// - msgHash: The message hash to be verified
// - r: The r component of the signature
// - s: The s component of the signature
// - pubX: The x-coordinate of the public key used for verification
// - pubY: The y-coordinate of the public key used for verification
// Returns:
// - bool: true if the signature is valid, false otherwise
func (sc StarkCurve) Verify(msgHash, r, s, pubX, pubY *big.Int) bool {
	w := sc.InvModCurveSize(s)

	if s.Cmp(big.NewInt(0)) != 1 || s.Cmp(sc.N) != -1 {
		return false
	}
	if r.Cmp(big.NewInt(0)) != 1 || r.Cmp(sc.Max) != -1 {
		return false
	}
	if w.Cmp(big.NewInt(0)) != 1 || w.Cmp(sc.Max) != -1 {
		return false
	}
	if msgHash.Cmp(big.NewInt(0)) != 1 || msgHash.Cmp(sc.Max) != -1 {
		return false
	}
	if !sc.IsOnCurve(pubX, pubY) {
		return false
	}

	zGx, zGy, err := sc.MimicEcMultAir(msgHash, sc.EcGenX, sc.EcGenY, sc.MinusShiftPointX, sc.MinusShiftPointY)
	if err != nil {
		return false
	}

	rQx, rQy, err := sc.MimicEcMultAir(r, pubX, pubY, sc.Gx, sc.Gy)
	if err != nil {
		return false
	}
	inX, inY := sc.Add(zGx, zGy, rQx, rQy)
	wBx, wBy, err := sc.MimicEcMultAir(w, inX, inY, sc.Gx, sc.Gy)
	if err != nil {
		return false
	}

	outX, _ := sc.Add(wBx, wBy, sc.MinusShiftPointX, sc.MinusShiftPointY)
	if r.Cmp(outX) == 0 {
		return true
	} else {
		altY := new(big.Int).Neg(pubY)

		zGx, zGy, err = sc.MimicEcMultAir(msgHash, sc.EcGenX, sc.EcGenY, sc.MinusShiftPointX, sc.MinusShiftPointY)
		if err != nil {
			return false
		}

		rQx, rQy, err = sc.MimicEcMultAir(r, pubX, new(big.Int).Set(altY), sc.Gx, sc.Gy)
		if err != nil {
			return false
		}
		inX, inY = sc.Add(zGx, zGy, rQx, rQy)
		wBx, wBy, err = sc.MimicEcMultAir(w, inX, inY, sc.Gx, sc.Gy)
		if err != nil {
			return false
		}

		outX, _ = sc.Add(wBx, wBy, sc.MinusShiftPointX, sc.MinusShiftPointY)
		if r.Cmp(outX) == 0 {
			return true
		}
	}
	return false
}

// Sign calculates the signature of a message using the StarkCurve algorithm.
// Secret is generated using a golang implementation of RFC 6979.
// Implementation does not yet include "extra entropy" or "retry gen".
// (ref: https://datatracker.ietf.org/doc/html/rfc6979)
//
// Parameters:
// - msgHash: The hash of the message to be signed
// - privKey: The private key used for signing
// - seed: (Optional) Additional seed values used for generating the secret
// Returns:
// - x, y: The coordinates of the signature point on the curve
// - err: An error if any occurred during the signing process
func (sc StarkCurve) Sign(msgHash, privKey *big.Int, seed ...*big.Int) (x, y *big.Int, err error) {
	if msgHash == nil {
		return x, y, fmt.Errorf("nil msgHash")
	}
	if privKey == nil {
		return x, y, fmt.Errorf("nil privKey")
	}
	if msgHash.Cmp(big.NewInt(0)) != 1 || msgHash.Cmp(sc.Max) != -1 {
		return x, y, fmt.Errorf("invalid bit length")
	}

	inSeed := big.NewInt(0)
	if len(seed) == 1 && inSeed != nil {
		inSeed = seed[0]
	}
	for {
		k := sc.GenerateSecret(big.NewInt(0).Set(msgHash), big.NewInt(0).Set(privKey), big.NewInt(0).Set(inSeed))
		// In case r is rejected k shall be generated with new seed
		inSeed = inSeed.Add(inSeed, big.NewInt(1))

		r, _ := sc.EcMult(k, sc.EcGenX, sc.EcGenY)

		// DIFF: in classic ECDSA, we take int(x) % n.
		if r.Cmp(big.NewInt(0)) != 1 || r.Cmp(sc.Max) != -1 {
			// Bad value. This fails with negligible probability.
			continue
		}

		agg := new(big.Int).Mul(r, privKey)
		agg = agg.Add(agg, msgHash)

		if new(big.Int).Mod(agg, sc.N).Cmp(big.NewInt(0)) == 0 {
			// Bad value. This fails with negligible probability.
			continue
		}

		w := DivMod(k, agg, sc.N)
		if w.Cmp(big.NewInt(0)) != 1 || w.Cmp(sc.Max) != -1 {
			// Bad value. This fails with negligible probability.
			continue
		}

		s := sc.InvModCurveSize(w)
		return r, s, nil
	}

	return x, y, nil
}

// SignFelt signs a message hash with a private key using the StarkCurve.
// just wraps Sign (previous function).
//
// Parameters:
// - msgHash: the message hash to be signed
// - privKey: the private key used for signing
// Returns:
// - xFelt: The x-coordinate of the signed message
// - yFelt: The y-coordinate of the signed message
// - error: An error if the signing process fails
func (sc StarkCurve) SignFelt(msgHash, privKey *felt.Felt) (*felt.Felt, *felt.Felt, error) {
	msgHashInt := msgHash.BigInt(new(big.Int))
	privKeyInt := privKey.BigInt(new(big.Int))
	x, y, err := sc.Sign(msgHashInt, privKeyInt)
	if err != nil {
		return nil, nil, err
	}
	xFelt := felt.NewFelt(new(felt.Felt).Impl().SetBigInt(x))
	yFelt := felt.NewFelt(new(felt.Felt).Impl().SetBigInt(y))
	return xFelt, yFelt, nil
}

// HashElements calculates the hash of a list of elements using the StarkCurve struct and a golang Pedersen Hash.
// (ref: https://github.com/seanjameshan/starknet.js/blob/main/src/utils/ellipticCurve.ts)
//
// Parameters:
// - elems: slice of big.Int pointers to be hashed
// Returns:
// - hash: The hash of the list of elements
// - err: An error if any
func (sc StarkCurve) HashElements(elems []*big.Int) (hash *big.Int, err error) {
	if len(elems) == 0 {
		elems = append(elems, big.NewInt(0))
	}

	hash = big.NewInt(0)
	for _, h := range elems {
		hash, err = sc.PedersenHash([]*big.Int{hash, h})
		if err != nil {
			return hash, err
		}
	}
	return hash, err
}

// ComputeHashOnElements computes the hash on the given elements using a golang Pedersen Hash implementation.
// (ref: https://github.com/starkware-libs/cairo-lang/blob/13cef109cd811474de114925ee61fd5ac84a25eb/src/starkware/cairo/common/hash_state.py#L6)
//
// The function appends the length of `elems` to the slice and then calls the `HashElements` method of the
// `Curve` struct, passing in `elems` as an argument. The resulting hash and
// any error that occurred during computation are returned.
//
// Parameters:
// - elems: slice of big.Int pointers to be hashed
// Returns:
// - hash: The hash of the list of elements
// - err: An error if any
func (sc StarkCurve) ComputeHashOnElements(elems []*big.Int) (hash *big.Int, err error) {
	elems = append(elems, big.NewInt(int64(len(elems))))
	return Curve.HashElements((elems))
}

// PedersenHash calculates the Pedersen hash of the given elements.
// NOTE: This function assumes the curve has been initialized with constant points
// (ref: https://github.com/seanjameshan/starknet.js/blob/main/src/utils/ellipticCurve.ts)
//
// The function requires that the precomputed constant points have been initiated.
// If the length of `sc.ConstantPoints` is zero, an error is returned.
// The function iterates over the elements in `elems` and performs the Pedersen hash calculation.
// For each element, it checks if the value is within the valid range.
// If the value is invalid, an error is returned.
// For each bit in the element, the function performs an addition operation on `ptx` and `pty`
// using the corresponding constant point from the precomputed constant points.
// If the constant point is a duplicate of `ptx`, an error is returned.
// The function returns the resulting hash and a nil error if the calculation is successful.
// Otherwise, it returns `ptx` and an error describing the issue encountered.
//
// Parameters:
// - elems: An array of big integers representing the elements to hash.
// Returns:
// - hash: The resulting Pedersen hash as a big integer.
// - err: An error, if any, encountered during the calculation.
func (sc StarkCurve) PedersenHash(elems []*big.Int) (hash *big.Int, err error) {
	if len(sc.ConstantPoints) == 0 {
		return hash, fmt.Errorf("must initiate precomputed constant points")
	}

	ptx := new(big.Int).Set(sc.Gx)
	pty := new(big.Int).Set(sc.Gy)
	for i, elem := range elems {
		x := new(big.Int).Set(elem)

		if x.Cmp(big.NewInt(0)) != -1 && x.Cmp(sc.P) != -1 {
			return ptx, fmt.Errorf("invalid x: %v", x)
		}

		for j := 0; j < 252; j++ {
			idx := 2 + (i * 252) + j
			xin := new(big.Int).Set(sc.ConstantPoints[idx][0])
			yin := new(big.Int).Set(sc.ConstantPoints[idx][1])
			if xin.Cmp(ptx) == 0 {
				return hash, fmt.Errorf("constant point duplication: %v %v", ptx, xin)
			}
			if x.Bit(0) == 1 {
				ptx, pty = sc.Add(ptx, pty, xin, yin)
			}
			x = x.Rsh(x, 1)
		}
	}

	return ptx, nil
}

// PoseidonArray is a function that takes a variadic number of felt.Felt pointers as parameters and
// NOTE: This function just wraps the Juno implementation
// (ref: https://github.com/NethermindEth/juno/blob/main/core/crypto/poseidon_hash.go#L74)
// calls the PoseidonArray function from the junoCrypto package with the provided parameters.
//
// Parameters:
// - felts: A variadic number of pointers to felt.Felt
// Returns:
// - *felt.Felt: pointer to a felt.Felt
func (sc StarkCurve) PoseidonArray(felts ...*felt.Felt) *felt.Felt {
	return junoCrypto.PoseidonArray(felts...)
}

// StarknetKeccak computes the Starknet Keccak hash of the given byte slice.
// NOTE: This function just wraps the Juno implementation
// (ref: https://github.com/NethermindEth/juno/blob/main/core/crypto/keccak.go#L11)
//
// Parameters:
// - b: The byte slice to hash
// Returns:
// - *felt.Felt: pointer to a felt.Felt
// - error: An error if any
func (sc StarkCurve) StarknetKeccak(b []byte) (*felt.Felt, error) {
	return junoCrypto.StarknetKeccak(b)
}

// GenerateSecret generates a secret using the StarkCurve struct.
// implementation based on https://github.com/codahale/rfc6979/blob/master/rfc6979.go
// for the specification, see https://tools.ietf.org/html/rfc6979#section-3.2
//
// Parameters:
// - msgHash: a pointer to a big.Int representing the message hash
// - privKey: a pointer to a big.Int representing the private key
// - seed: a pointer to a big.Int representing the seed
// Returns:
// - secret: a pointer to a big.Int representing the generated secret
func (sc StarkCurve) GenerateSecret(msgHash, privKey, seed *big.Int) (secret *big.Int) {
	alg := sha256.New
	holen := alg().Size()
	rolen := (sc.BitSize + 7) >> 3

	if msgHash.BitLen()%8 <= 4 && msgHash.BitLen() >= 248 {
		msgHash = msgHash.Mul(msgHash, big.NewInt(16))
	}

	by := append(int2octets(privKey, rolen), bits2octets(msgHash, sc.N, sc.BitSize, rolen)...)

	if seed.Cmp(big.NewInt(0)) == 1 {
		by = append(by, seed.Bytes()...)
	}

	v := bytes.Repeat([]byte{0x01}, holen)

	k := bytes.Repeat([]byte{0x00}, holen)

	k = mac(alg, k, append(append(v, 0x00), by...), k)

	v = mac(alg, k, v, v)

	k = mac(alg, k, append(append(v, 0x01), by...), k)

	v = mac(alg, k, v, v)

	for {
		var t []byte

		for len(t) < rolen {
			v = mac(alg, k, v, v)
			t = append(t, v...)
		}

		secret = bits2int(new(big.Int).SetBytes(t), sc.BitSize)
		// TODO: implement seed here, final gating function
		if secret.Cmp(big.NewInt(0)) == 1 && secret.Cmp(sc.N) == -1 {
			return secret
		}
		k = mac(alg, k, append(v, 0x00), k)
		v = mac(alg, k, v, v)
	}
}

// GetRandomPrivateKey generates a random private key for the StarkCurve struct.
// NOTE: to be used for testing purposes
//
// Parameters:
// - none
// Returns:
// - priv: a pointer to a big.Int representing the generated private key
// - err: an error if any
func (sc StarkCurve) GetRandomPrivateKey() (priv *big.Int, err error) {
	max := new(big.Int).Sub(sc.Max, big.NewInt(1))

	priv, err = rand.Int(rand.Reader, max)
	if err != nil {
		return priv, err
	}

	x, y, err := sc.PrivateToPoint(priv)
	if err != nil {
		return priv, err
	}

	if !sc.IsOnCurve(x, y) {
		return priv, fmt.Errorf("key gen is not on stark cruve")
	}

	return priv, nil
}

// PrivateToPoint generates a point on the StarkCurve from a private key.
//
// It takes a private key as a parameter and returns the x and y coordinates of 
// the generated point on the curve. If the private key is not within the range 
// of the curve, it returns an error.
//
// Parameters:
// - privKey: The private key used to generate the point
// Return values:
// - x: The x coordinate of the generated point
// - y: The y coordinate of the generated point
// - err: An error if the private key is not within the curve range
func (sc StarkCurve) PrivateToPoint(privKey *big.Int) (x, y *big.Int, err error) {
	if privKey.Cmp(big.NewInt(0)) != 1 || privKey.Cmp(sc.N) != -1 {
		return x, y, fmt.Errorf("private key not in curve range")
	}
	x, y = sc.EcMult(privKey, sc.EcGenX, sc.EcGenY)
	return x, y, nil
}
