package byzquorum

// Size of a byz. quorum. We assume n >= 3*f + 1.
func Size(n, f int) int {
	return (n+f)/2 + 1
}
