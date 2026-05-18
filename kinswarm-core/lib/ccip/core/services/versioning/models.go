package versioning

import (
	"time"
)

type NodeVersion struct {
	Version   string
	CreatedAt time.Time
}

func NewNodeVersion(version string) NodeVersion {
	return NodeVersion{
		Version:   version,
		CreatedAt: time.Now(),
	}
}
