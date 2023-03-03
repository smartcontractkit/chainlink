package edwards25519

// geScalarMultVartime computes h = a*B, where
//   a = a[0]+256*a[1]+...+256^31 a[31]
//   B is the Ed25519 base point (x,4/5) with x positive.
//
// Preconditions:
//   a[31] <= 127
func geScalarMultVartime(h *extendedGroupElement, a *[32]byte,
	A *extendedGroupElement) {

	var aSlide [256]int8
	var Ai [8]cachedGroupElement // A,3A,5A,7A,9A,11A,13A,15A
	var t completedGroupElement
	var u, A2 extendedGroupElement
	var r projectiveGroupElement
	var i int

	// Slide through the scalar exponent clumping sequences of bits,
	// resulting in only zero or odd multipliers between -15 and 15.
	slide(&aSlide, a)

	// Form an array of odd multiples of A from 1A through 15A,
	// in addition-ready cached group element form.
	// We only need odd multiples of A because slide()
	// produces only odd-multiple clumps of bits.
	A.ToCached(&Ai[0])
	A.Double(&t)
	t.ToExtended(&A2)
	for i := 0; i < 7; i++ {
		t.Add(&A2, &Ai[i])
		t.ToExtended(&u)
		u.ToCached(&Ai[i+1])
	}

	// Process the multiplications from most-significant bit downward
	for i = 255; ; i-- {
		if i < 0 { // no bits set
			h.Zero()
			return
		}
		if aSlide[i] != 0 {
			break
		}
	}

	// first (most-significant) nonzero clump of bits
	u.Zero()
	if aSlide[i] > 0 {
		t.Add(&u, &Ai[aSlide[i]/2])
	} else if aSlide[i] < 0 {
		t.Sub(&u, &Ai[(-aSlide[i])/2])
	}
	i--

	// remaining bits
	for ; i >= 0; i-- {
		t.ToProjective(&r)
		r.Double(&t)

		if aSlide[i] > 0 {
			t.ToExtended(&u)
			t.Add(&u, &Ai[aSlide[i]/2])
		} else if aSlide[i] < 0 {
			t.ToExtended(&u)
			t.Sub(&u, &Ai[(-aSlide[i])/2])
		}
	}

	t.ToExtended(h)
}
