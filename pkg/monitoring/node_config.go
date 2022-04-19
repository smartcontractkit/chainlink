package monitoring

import (
	"io"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

type NodesParser func(buf io.ReadCloser) ([]NodeConfig, error)

type NodeConfig interface {
	GetName() string
	GetAccount() types.Account
}
