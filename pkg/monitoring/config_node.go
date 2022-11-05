package monitoring

import (
	"io"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

// NodesParser extracts multiple nodes' configurations from the configuration server, eg. weiwatchers.com
type NodesParser func(buf io.ReadCloser) ([]NodeConfig, error)

// NodeConfig is the subset of on-chain node operator's configuration required by the OM framework.
type NodeConfig interface {
	GetName() string
	GetAccount() types.Account
}
