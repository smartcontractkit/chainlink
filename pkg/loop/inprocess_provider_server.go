package loop

// Temporary shim to allow the migration of generic plugins from using a provider connection to using a RelayerSet.
// Once that migration is complete, this shim can be removed.  Alternatively, once the migration of all relayers to
// run as LOOPP plugins is complete, this shim will no longer be required and can be removed

import (
	"context"

	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayerset/inprocessprovider"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

type ProviderServer interface {
	Start(context.Context) error
	Close() error
	GetConn() (grpc.ClientConnInterface, error)
}

func NewProviderServer(p types.PluginProvider, pType types.OCR2PluginType, lggr logger.Logger) (ProviderServer, error) {
	return inprocessprovider.NewProviderServer(p, pType, lggr)
}
