package typeandversion

import (
	"fmt"
	"strings"
)

const (
	UnknownContractType = "Unknown"
	DefaultVersion      = "1.0.0"
)

// ParseTypeAndVersion parses "<Type> <Semver>" strings returned by contracts.
// Empty strings are normalized to "Unknown 1.0.0".
func ParseTypeAndVersion(tvStr string) (string, string, error) {
	if tvStr == "" {
		tvStr = UnknownContractType + " " + DefaultVersion
	}
	typeAndVersionValues := strings.Split(tvStr, " ")
	if len(typeAndVersionValues) != 2 {
		return "", "", fmt.Errorf("invalid type and version %s", tvStr)
	}
	return typeAndVersionValues[0], typeAndVersionValues[1], nil
}
