package generated

import (
	"github.com/ethereum/go-ethereum/common"
)

// AbigenLog is an interface for abigen generated log topics
type AbigenLog interface {
	Topic() common.Hash
}
