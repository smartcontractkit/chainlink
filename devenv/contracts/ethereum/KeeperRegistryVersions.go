package ethereum

import (
	"github.com/ethereum/go-ethereum/common"
)

// AbigenLog is an interface for abigen generated log topics
type AbigenLog interface {
	Topic() common.Hash
}

type KeeperRegistryVersion int32

//nolint: revive //we want to use underscores
const (
	RegistryVersion_1_0 KeeperRegistryVersion = iota
	RegistryVersion_1_1
	RegistryVersion_1_2
	RegistryVersion_1_3
	RegistryVersion_2_0
	RegistryVersion_2_1
	RegistryVersion_2_2
	RegistryVersion_2_3
)

func (k KeeperRegistryVersion) String() string {
	switch k {
	case RegistryVersion_1_0:
		return "1.0"
	case RegistryVersion_1_1:
		return "1.1"
	case RegistryVersion_1_2:
		return "1.2"
	case RegistryVersion_1_3:
		return "1.3"
	case RegistryVersion_2_0:
		return "2.0"
	case RegistryVersion_2_1:
		return "2.1"
	case RegistryVersion_2_2:
		return "2.2"
	case RegistryVersion_2_3:
		return "2.3"
	default:
		return "unknown"
	}
}
