package adapters

import (
	"github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/client"
	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/config"
)

type Chain interface {
	types.ChainService

	ID() string
	Config() config.Config
	TxManager() TxManager
	// Reader returns a new Reader. If nodeName is provided, the underlying client must use that node.
	Reader(nodeName string) (client.Reader, error)
}
