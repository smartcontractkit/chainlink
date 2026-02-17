package tunnel

import (
	"fmt"
	"strings"
)

type ComponentKind string

const (
	KindBlockchain ComponentKind = "blockchain"
	KindNodeSet    ComponentKind = "nodeset"
	KindJD         ComponentKind = "jd"
)

func CanonicalComponentID(kind ComponentKind, index int, name string) string {
	if name == "" {
		return fmt.Sprintf("%s:%d", kind, index)
	}

	normalized := strings.ToLower(strings.TrimSpace(name))
	if normalized == "" {
		return fmt.Sprintf("%s:%d", kind, index)
	}

	return fmt.Sprintf("%s:%d:%s", kind, index, normalized)
}
