package common

func isHexCharacter(c byte) bool {
	return ('0' <= c && c <= '9') || ('a' <= c && c <= 'f') || ('A' <= c && c <= 'F')
}

func IsHex(s string) bool {
	if has0xPrefix(s) {
		s = s[2:]
	}

	for i := 0; i < len(s); i++ {
		if !isHexCharacter(s[i]) {
			return false
		}
	}
	return true
}
