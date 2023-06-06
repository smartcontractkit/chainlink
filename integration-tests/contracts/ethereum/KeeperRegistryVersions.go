package ethereum

import (
	"github.com/ethereum/go-ethereum/common"
)

// AbigenLog is an interface for abigen generated log topics
type AbigenLog interface {
	Topic() common.Hash
}

type KeeperRegistryVersion int32

const (
	RegistryVersion_1_0 KeeperRegistryVersion = iota
	RegistryVersion_1_1
	RegistryVersion_1_2
	RegistryVersion_1_3
	RegistryVersion_2_0
)
