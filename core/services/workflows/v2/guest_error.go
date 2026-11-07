package v2

import (
	"fmt"
	"unicode/utf8"
)

// maxGuestErrorBytes bounds a workflow-returned error before it reaches logs and telemetry.
const maxGuestErrorBytes = 4 * 1024

// truncateGuestError caps s at maxGuestErrorBytes, cutting on a UTF-8 rune boundary.
func truncateGuestError(s string) string {
	if len(s) <= maxGuestErrorBytes {
		return s
	}
	cut := maxGuestErrorBytes
	for cut > 0 && !utf8.RuneStart(s[cut]) {
		cut--
	}
	return fmt.Sprintf("%s... [truncated, %d bytes total]", s[:cut], len(s))
}
