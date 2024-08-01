package fp

import "math/big"

// The following code is _almost_ the original code from:
// https://github.com/GottfriedHerold/Bandersnatch/blob/f665f90b64892b9c4c89cff3219e70456bb431e5/bandersnatch/fieldElements/field_element_square_root.go
//
// We had to do some changes to make it work with gnark:
// - The type `feType_SquareRoot` was aliased to `Element` so everything looks the same. These types didn't have the exact
//   same underlying representation, so it leaded to some minor adjustements. (e.g: accessing the limbs)
// - Original APIs regarding finite-field multiplications (e.g: MulEq) were adjusted to use gnark Mul APIs.
// - The original code had to explicitly do `Normalize()` after field element operations, but this isn't needed in gnark.
// - The primitive 2^32-root-of unity value (see init()) was pulled from gnark FFT domain code.
// - The original code used anonymous functions to define global vars, but we changed to use a init() function.
//   This was required since we have other init() in the package that configure other globals (e.g: _modulus).
//   By the way init() functions execution order works, we'll have these configured before the sqrt init() is called,
//   compared with the original anonymous function global calls.

type feType_SquareRoot = Element

const (
	BaseField2Adicity              = 32
	sqrtParam_TotalBits            = BaseField2Adicity // (p-1) = n^Q. 2^S with Q odd, leads to S = 32.
	sqrtParam_BlockSize            = 8                 // 8 bit window per chunk
	sqrtParam_Blocks               = sqrtParam_TotalBits / sqrtParam_BlockSize
	sqrtParam_FirstBlockUnusedBits = sqrtParam_Blocks*sqrtParam_BlockSize - sqrtParam_TotalBits // number of unused bits in the first reconstructed block.
	sqrtParam_BitMask              = (1 << sqrtParam_BlockSize) - 1                             // bitmask to pick up the last sqrtParam_BlockSize bits.
)

// NOTE: These "variables" are actually pre-computed constants that must not change.
var (
	// sqrtPrecomp_PrimitiveDyadicRoots[i] equals DyadicRootOfUnity^(2^i) for 0 <= i <= 32
	//
	// This means that it is a 32-i'th primitive root of unitity, obtained by repeatedly squaring a 2^32th primitive root of unity [DyadicRootOfUnity_fe].
	sqrtPrecomp_PrimitiveDyadicRoots [BaseField2Adicity + 1]feType_SquareRoot

	// primitive root of unity of order 2^sqrtParam_BlockSize
	sqrtPrecomp_ReconstructionDyadicRoot feType_SquareRoot

	// sqrtPrecomp_dlogLUT is a lookup table used to implement the map sqrtPrecompt_reconstructionDyadicRoot^a -> -a
	sqrtPrecomp_dlogLUT map[uint16]uint
)

func init() {
	sqrtPrecomp_PrimitiveDyadicRoots = func() (ret [BaseField2Adicity + 1]feType_SquareRoot) {
		if _, err := ret[0].SetString("10238227357739495823651030575849232062558860180284477541189508159991286009131"); err != nil {
			panic(err)
		}
		for i := 1; i <= BaseField2Adicity; i++ { // Note <= here
			ret[i].Square(&ret[i-1])
		}
		// 31th one must be -1. We check that here.
		x := big.NewInt(0)
		ret[BaseField2Adicity-1].BigInt(x)
		if ret[BaseField2Adicity-1].String() != "-1" {
			panic("something is wrong with the dyadic roots of unity")
		}
		return
	}() // immediately invoked lambda
	sqrtPrecomp_ReconstructionDyadicRoot = sqrtPrecomp_PrimitiveDyadicRoots[BaseField2Adicity-sqrtParam_BlockSize]
	sqrtPrecomp_PrecomputedBlocks = func() (blocks [sqrtParam_Blocks][1 << sqrtParam_BlockSize]feType_SquareRoot) {
		for i := 0; i < sqrtParam_Blocks; i++ {
			blocks[i][0].SetOne()
			for j := 1; j < (1 << sqrtParam_BlockSize); j++ {
				blocks[i][j].Mul(&blocks[i][j-1], &sqrtPrecomp_PrimitiveDyadicRoots[i*sqrtParam_BlockSize])
			}
		}
		return
	}() // immediately invoked lambda

	sqrtPrecomp_dlogLUT = func() (ret map[uint16]uint) {
		const LUTSize = 1 << sqrtParam_BlockSize // 256
		ret = make(map[uint16]uint, LUTSize)

		var rootOfUnity feType_SquareRoot
		rootOfUnity.SetOne()
		for i := 0; i < LUTSize; i++ {
			const mask = LUTSize - 1
			// the LUTSize many roots of unity all (by chance) have distinct values for .words[0]&0xFFFF. Note that this uses the Montgomery representation.
			ret[uint16(rootOfUnity[0]&0xFFFF)] = uint((-i) & mask)
			rootOfUnity.Mul(&rootOfUnity, &sqrtPrecomp_ReconstructionDyadicRoot)
		}
		// This effectively checks the above claim (that .words[0]&0xFFFF is distinct).
		// Note that this might fail if we adjust the sqrtParam_BlockSize parameter and this check will alert us.
		if len(ret) != LUTSize {
			panic("failed to store all appropriate roots of unity in a map")
		}
		return
	}() // immediately invoked lambda
}

// sqrtAlg_NegDlogInSmallDyadicSubgroup takes a (not necessarily primitive) root of unity x of order 2^sqrtParam_BlockSize.
// x has the form sqrtPrecomp_ReconstructionDyadicRoot^a and returns its negative dlog -a.
//
// The returned value is only meaningful modulo 1<<sqrtParam_BlockSize and is fully reduced, i.e. in [0, 1<<sqrtParam_BlockSize )
//
// NOTE: If x is not a root of unity as asserted, the behaviour is undefined.
func sqrtAlg_NegDlogInSmallDyadicSubgroup(x *feType_SquareRoot) uint {
	return sqrtPrecomp_dlogLUT[uint16(x[0]&0xFFFF)]
}

// sqrtAlg_GetPrecomputedRootOfUnity sets target to g^(multiplier << (order * sqrtParam_BlockSize)), where g is the fixed primitive 2^32th root of unity.
//
// We assume that order 0 <= order*sqrtParam_BlockSize <= 32 and that multiplier is in [0, 1 <<sqrtParam_BlockSize)
func sqrtAlg_GetPrecomputedRootOfUnity(target *feType_SquareRoot, multiplier int, order uint) {
	*target = sqrtPrecomp_PrecomputedBlocks[order][multiplier]
}

// sqrtPrecomp_PrecomputedBlocks[i][j] == g^(j << (i* BlockSize)), where g is the fixed primitive 2^32th root of unity.
// This means that the exponent is equal to 0x00000...0000jjjjjj0000....0000, where only the i'th least significant block of size BlockSize is set
// and that value is j.
//
// Note: accessed through sqrtAlg_getPrecomputedRootOfUnity
var sqrtPrecomp_PrecomputedBlocks [sqrtParam_Blocks][1 << sqrtParam_BlockSize]feType_SquareRoot

func SqrtPrecomp(x *Element) *Element {
	res := Zero()
	if x.IsZero() {
		return &res
	}
	var xCopy feType_SquareRoot = *x
	var candidate, rootOfUnity feType_SquareRoot
	sqrtAlg_ComputeRelevantPowers(&xCopy, &candidate, &rootOfUnity)
	if !invSqrtEqDyadic(&rootOfUnity) {
		return nil
	}

	return res.Mul(&candidate, &rootOfUnity)
}

func invSqrtEqDyadic(z *Element) bool {
	// The algorithm works by essentially computing the dlog of z and then halving it.

	// negExponent is intended to hold the negative of the dlog of z.
	// We determine this 32-bit value (usually) _sqrtBlockSize many bits at a time, starting with the least-significant bits.
	//
	// If _sqrtBlockSize does not divide 32, the *first* iteration will determine fewer bits.
	var negExponent uint

	var temp, temp2 feType_SquareRoot

	// set powers[i] to z^(1<< (i*blocksize))
	var powers [sqrtParam_Blocks]feType_SquareRoot
	powers[0] = *z
	for i := 1; i < sqrtParam_Blocks; i++ {
		powers[i] = powers[i-1]
		for j := 0; j < sqrtParam_BlockSize; j++ {
			powers[i].Square(&powers[i])
		}
	}

	// looking at the dlogs, powers[i] is essentially the wanted exponent, left-shifted by i*_sqrtBlockSize and taken mod 1<<32
	// dlogHighDyadicRootNeg essentially (up to sign) reads off the _sqrtBlockSize many most significant bits. (returned as low-order bits)

	// first iteration may be slightly special if BlockSize does not divide 32
	negExponent = sqrtAlg_NegDlogInSmallDyadicSubgroup(&powers[sqrtParam_Blocks-1])
	negExponent >>= sqrtParam_FirstBlockUnusedBits

	// if the exponent we just got is odd, there is no square root, no point in determining the other bits.
	if negExponent&1 == 1 {
		return false
	}

	// Get remaining bits
	for i := 1; i < sqrtParam_Blocks; i++ {
		temp2 = powers[sqrtParam_Blocks-1-i]
		// We essentially un-set the bits we already know from powers[_sqrtNumBlocks-1-i]
		for j := 0; j < i; j++ {
			sqrtAlg_GetPrecomputedRootOfUnity(&temp, int((negExponent>>(j*sqrtParam_BlockSize))&sqrtParam_BitMask), uint(j+sqrtParam_Blocks-1-i))
			temp2.Mul(&temp2, &temp)
		}
		newBits := sqrtAlg_NegDlogInSmallDyadicSubgroup(&temp2)
		negExponent |= newBits << (sqrtParam_BlockSize*i - sqrtParam_FirstBlockUnusedBits)
	}

	// var tmp _FESquareRoot

	// negExponent is now the negative dlog of z.

	// Take the square root
	negExponent >>= 1
	// Write to z:
	z.SetOne()
	for i := 0; i < sqrtParam_Blocks; i++ {
		sqrtAlg_GetPrecomputedRootOfUnity(&temp, int((negExponent>>(i*sqrtParam_BlockSize))&sqrtParam_BitMask), uint(i))
		z.Mul(z, &temp)
	}

	return true
}

func sqrtAlg_ComputeRelevantPowers(z *Element, squareRootCandidate *feType_SquareRoot, rootOfUnity *feType_SquareRoot) {
	SquareEqNTimes := func(z *feType_SquareRoot, n int) {
		for i := 0; i < n; i++ {
			z.Square(z)
		}
	}

	// hand-crafted sliding window-type algorithm with window-size 5
	// Note that we precompute and use z^255 multiple times (even though it's not size 5)
	// and some windows actually overlap(!)

	var z2, z3, z7, z6, z9, z11, z13, z19, z21, z25, z27, z29, z31, z255 feType_SquareRoot
	var acc feType_SquareRoot
	z2.Square(z)             // 0b10
	z3.Mul(z, &z2)           // 0b11
	z6.Square(&z3)           // 0b110
	z7.Mul(z, &z6)           // 0b111
	z9.Mul(&z7, &z2)         // 0b1001
	z11.Mul(&z9, &z2)        // 0b1011
	z13.Mul(&z11, &z2)       // 0b1101
	z19.Mul(&z13, &z6)       // 0b10011
	z21.Mul(&z2, &z19)       // 0b10101
	z25.Mul(&z19, &z6)       // 0b11001
	z27.Mul(&z25, &z2)       // 0b11011
	z29.Mul(&z27, &z2)       // 0b11101
	z31.Mul(&z29, &z2)       // 0b11111
	acc.Mul(&z27, &z29)      // 56
	acc.Square(&acc)         // 112
	acc.Square(&acc)         // 224
	z255.Mul(&acc, &z31)     // 0b11111111 = 255
	acc.Square(&acc)         // 448
	acc.Square(&acc)         // 896
	acc.Mul(&acc, &z31)      // 0b1110011111 = 927
	SquareEqNTimes(&acc, 6)  // 0b1110011111000000
	acc.Mul(&acc, &z27)      // 0b1110011111011011
	SquareEqNTimes(&acc, 6)  // 0b1110011111011011000000
	acc.Mul(&acc, &z19)      // 0b1110011111011011010011
	SquareEqNTimes(&acc, 5)  // 0b111001111101101101001100000
	acc.Mul(&acc, &z21)      // 0b111001111101101101001110101
	SquareEqNTimes(&acc, 7)  // 0b1110011111011011010011101010000000
	acc.Mul(&acc, &z25)      // 0b1110011111011011010011101010011001
	SquareEqNTimes(&acc, 6)  // 0b1110011111011011010011101010011001000000
	acc.Mul(&acc, &z19)      // 0b1110011111011011010011101010011001010011
	SquareEqNTimes(&acc, 5)  // 0b111001111101101101001110101001100101001100000
	acc.Mul(&acc, &z7)       // 0b111001111101101101001110101001100101001100111
	SquareEqNTimes(&acc, 5)  // 0b11100111110110110100111010100110010100110011100000
	acc.Mul(&acc, &z11)      // 0b11100111110110110100111010100110010100110011101011
	SquareEqNTimes(&acc, 5)  // 0b1110011111011011010011101010011001010011001110101100000
	acc.Mul(&acc, &z29)      // 0b1110011111011011010011101010011001010011001110101111101
	SquareEqNTimes(&acc, 5)  // 0b111001111101101101001110101001100101001100111010111110100000
	acc.Mul(&acc, &z9)       // 0b111001111101101101001110101001100101001100111010111110101001
	SquareEqNTimes(&acc, 7)  // 0b1110011111011011010011101010011001010011001110101111101010010000000
	acc.Mul(&acc, &z3)       // 0b1110011111011011010011101010011001010011001110101111101010010000011
	SquareEqNTimes(&acc, 7)  // 0b11100111110110110100111010100110010100110011101011111010100100000110000000
	acc.Mul(&acc, &z25)      // 0b11100111110110110100111010100110010100110011101011111010100100000110011001
	SquareEqNTimes(&acc, 5)  // 0b1110011111011011010011101010011001010011001110101111101010010000011001100100000
	acc.Mul(&acc, &z25)      // 0b1110011111011011010011101010011001010011001110101111101010010000011001100111001
	SquareEqNTimes(&acc, 5)  // 0b111001111101101101001110101001100101001100111010111110101001000001100110011100100000
	acc.Mul(&acc, &z27)      // 0b111001111101101101001110101001100101001100111010111110101001000001100110011100111011
	SquareEqNTimes(&acc, 8)  // 0b11100111110110110100111010100110010100110011101011111010100100000110011001110011101100000000
	acc.Mul(&acc, z)         // 0b11100111110110110100111010100110010100110011101011111010100100000110011001110011101100000001
	SquareEqNTimes(&acc, 8)  // 0b1110011111011011010011101010011001010011001110101111101010010000011001100111001110110000000100000000
	acc.Mul(&acc, z)         // 0b1110011111011011010011101010011001010011001110101111101010010000011001100111001110110000000100000001
	SquareEqNTimes(&acc, 6)  // 0b1110011111011011010011101010011001010011001110101111101010010000011001100111001110110000000100000001000000
	acc.Mul(&acc, &z13)      // 0b1110011111011011010011101010011001010011001110101111101010010000011001100111001110110000000100000001001101
	SquareEqNTimes(&acc, 7)  // 0b11100111110110110100111010100110010100110011101011111010100100000110011001110011101100000001000000010011010000000
	acc.Mul(&acc, &z7)       // 0b11100111110110110100111010100110010100110011101011111010100100000110011001110011101100000001000000010011010000111
	SquareEqNTimes(&acc, 3)  // 0b11100111110110110100111010100110010100110011101011111010100100000110011001110011101100000001000000010011010000111000
	acc.Mul(&acc, &z3)       // 0b11100111110110110100111010100110010100110011101011111010100100000110011001110011101100000001000000010011010000111011
	SquareEqNTimes(&acc, 13) // 0b111001111101101101001110101001100101001100111010111110101001000001100110011100111011000000010000000100110100001110110000000000000
	acc.Mul(&acc, &z21)      // 0b111001111101101101001110101001100101001100111010111110101001000001100110011100111011000000010000000100110100001110110000000010101
	SquareEqNTimes(&acc, 5)  // 0b11100111110110110100111010100110010100110011101011111010100100000110011001110011101100000001000000010011010000111011000000001010100000
	acc.Mul(&acc, &z9)       // 0b11100111110110110100111010100110010100110011101011111010100100000110011001110011101100000001000000010011010000111011000000001010101001
	SquareEqNTimes(&acc, 5)  // 0b1110011111011011010011101010011001010011001110101111101010010000011001100111001110110000000100000001001101000011101100000000101010100100000
	acc.Mul(&acc, &z27)      // 0b1110011111011011010011101010011001010011001110101111101010010000011001100111001110110000000100000001001101000011101100000000101010100111011
	SquareEqNTimes(&acc, 5)  // 0b111001111101101101001110101001100101001100111010111110101001000001100110011100111011000000010000000100110100001110110000000010101010011101100000
	acc.Mul(&acc, &z27)      // 0b111001111101101101001110101001100101001100111010111110101001000001100110011100111011000000010000000100110100001110110000000010101010011101111011
	SquareEqNTimes(&acc, 5)  // 0b11100111110110110100111010100110010100110011101011111010100100000110011001110011101100000001000000010011010000111011000000001010101001110111101100000
	acc.Mul(&acc, &z9)       // 0b11100111110110110100111010100110010100110011101011111010100100000110011001110011101100000001000000010011010000111011000000001010101001110111101101001
	SquareEqNTimes(&acc, 10) // 0b111001111101101101001110101001100101001100111010111110101001000001100110011100111011000000010000000100110100001110110000000010101010011101111011010010000000000
	acc.Mul(&acc, z)         // 0b111001111101101101001110101001100101001100111010111110101001000001100110011100111011000000010000000100110100001110110000000010101010011101111011010010000000001
	SquareEqNTimes(&acc, 7)  // 0b1110011111011011010011101010011001010011001110101111101010010000011001100111001110110000000100000001001101000011101100000000101010100111011110110100100000000010000000
	acc.Mul(&acc, &z255)     // 0b1110011111011011010011101010011001010011001110101111101010010000011001100111001110110000000100000001001101000011101100000000101010100111011110110100100000000101111111
	SquareEqNTimes(&acc, 8)  // 0b111001111101101101001110101001100101001100111010111110101001000001100110011100111011000000010000000100110100001110110000000010101010011101111011010010000000010111111100000000
	acc.Mul(&acc, &z255)     // 0b111001111101101101001110101001100101001100111010111110101001000001100110011100111011000000010000000100110100001110110000000010101010011101111011010010000000010111111111111111
	SquareEqNTimes(&acc, 6)  // 0b111001111101101101001110101001100101001100111010111110101001000001100110011100111011000000010000000100110100001110110000000010101010011101111011010010000000010111111111111111000000
	acc.Mul(&acc, &z11)      // 0b111001111101101101001110101001100101001100111010111110101001000001100110011100111011000000010000000100110100001110110000000010101010011101111011010010000000010111111111111111001011
	SquareEqNTimes(&acc, 9)  // 0b111001111101101101001110101001100101001100111010111110101001000001100110011100111011000000010000000100110100001110110000000010101010011101111011010010000000010111111111111111001011000000000
	acc.Mul(&acc, &z255)     // 0b111001111101101101001110101001100101001100111010111110101001000001100110011100111011000000010000000100110100001110110000000010101010011101111011010010000000010111111111111111001011011111111
	SquareEqNTimes(&acc, 2)  // 0b11100111110110110100111010100110010100110011101011111010100100000110011001110011101100000001000000010011010000111011000000001010101001110111101101001000000001011111111111111100101101111111100
	acc.Mul(&acc, z)         // 0b11100111110110110100111010100110010100110011101011111010100100000110011001110011101100000001000000010011010000111011000000001010101001110111101101001000000001011111111111111100101101111111101
	SquareEqNTimes(&acc, 7)  // 0b111001111101101101001110101001100101001100111010111110101001000001100110011100111011000000010000000100110100001110110000000010101010011101111011010010000000010111111111111111001011011111111010000000
	acc.Mul(&acc, &z255)     // 0b111001111101101101001110101001100101001100111010111110101001000001100110011100111011000000010000000100110100001110110000000010101010011101111011010010000000010111111111111111001011011111111101111111
	SquareEqNTimes(&acc, 8)  // 0b11100111110110110100111010100110010100110011101011111010100100000110011001110011101100000001000000010011010000111011000000001010101001110111101101001000000001011111111111111100101101111111110111111100000000
	acc.Mul(&acc, &z255)     // 0b11100111110110110100111010100110010100110011101011111010100100000110011001110011101100000001000000010011010000111011000000001010101001110111101101001000000001011111111111111100101101111111110111111111111111
	SquareEqNTimes(&acc, 8)  // 0b1110011111011011010011101010011001010011001110101111101010010000011001100111001110110000000100000001001101000011101100000000101010100111011110110100100000000101111111111111110010110111111111011111111111111100000000
	acc.Mul(&acc, &z255)     // 0b1110011111011011010011101010011001010011001110101111101010010000011001100111001110110000000100000001001101000011101100000000101010100111011110110100100000000101111111111111110010110111111111011111111111111111111111
	SquareEqNTimes(&acc, 8)  // 0b111001111101101101001110101001100101001100111010111110101001000001100110011100111011000000010000000100110100001110110000000010101010011101111011010010000000010111111111111111001011011111111101111111111111111111111100000000
	acc.Mul(&acc, &z255)     // 0b111001111101101101001110101001100101001100111010111110101001000001100110011100111011000000010000000100110100001110110000000010101010011101111011010010000000010111111111111111001011011111111101111111111111111111111111111111
	// acc is now z^((BaseFieldMultiplicativeOddOrder - 1)/2)
	rootOfUnity.Square(&acc)         // BaseFieldMultiplicativeOddOrder - 1
	rootOfUnity.Mul(rootOfUnity, z)  // BaseFieldMultiplicativeOddOrder
	squareRootCandidate.Mul(&acc, z) // (BaseFieldMultiplicativeOddOrder + 1)/2
}
