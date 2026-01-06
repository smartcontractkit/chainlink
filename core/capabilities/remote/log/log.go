package log

import (
	"encoding/hex"
	"unicode"
)

// TODO this is duplicated from utils to minimise the PR changeset - will be deduped later

const (
	maxLoggedStringLen = 256
)

func SanitizeLogString(s string) string {
	tooLongSuffix := ""
	if len(s) > maxLoggedStringLen {
		s = s[:maxLoggedStringLen]
		tooLongSuffix = " [TRUNCATED]"
	}
	for i := 0; i < len(s); i++ {
		if !unicode.IsPrint(rune(s[i])) {
			return "[UNPRINTABLE] " + hex.EncodeToString([]byte(s)) + tooLongSuffix
		}
	}
	return s + tooLongSuffix
}
