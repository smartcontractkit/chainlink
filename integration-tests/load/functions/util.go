package loadfunctions

// StringToByte32 transforms a single string into a [32]byte value
func StringToByte32(s string) [32]byte {
	var result [32]byte

	for i, ch := range []byte(s) {
		if i > 31 {
			break
		}
		result[i] = ch
	}

	return result
}
