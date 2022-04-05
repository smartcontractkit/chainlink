package utils

import "strings"

// RemovePrefix filters out any prefixes attached to keys from the core node (if present)
func RemovePrefix(s string) string {
	sArr := strings.Split(s, "_")
	return sArr[len(sArr)-1] // always use the last split (removes multiple prefixes if present)
}
