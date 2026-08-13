package registrysyncer

import (
	"context"

	"github.com/smartcontractkit/capabilities/libs/x/registrysyncer"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
)

// registrySyncerStatesTable is where core keeps its registry snapshots. The shared ORM is told its
// table rather than assuming one, because the processes that keep a registry do not share a
// database: this is core's, and crecore has its own.
const registrySyncerStatesTable = "registry_syncer_states"

// ORM persists registry snapshots.
//
// AddLocalRegistry takes a snapshot by value, as core's callers and mocks have always passed one.
// The shared ORM takes a pointer - a LocalRegistry carries the mutex guarding its local-node cache -
// so this interface exists to keep that difference off core's call sites.
type ORM interface {
	AddLocalRegistry(ctx context.Context, localRegistry LocalRegistry) error
	LatestLocalRegistry(ctx context.Context) (*LocalRegistry, error)
}

// NewORM returns the shared ORM pointed at core's own table, behind core's own by-value interface.
func NewORM(ds sqlutil.DataSource, lggr logger.Logger) ORM {
	return byValueORM{ORM: registrysyncer.NewORM(ds, lggr, registrySyncerStatesTable)}
}

// byValueORM adapts the shared ORM to the by-value AddLocalRegistry core is written against.
// LatestLocalRegistry needs no adapting, so it comes straight from the embedded ORM.
type byValueORM struct {
	registrysyncer.ORM
}

var _ ORM = byValueORM{}

func (o byValueORM) AddLocalRegistry(ctx context.Context, localRegistry LocalRegistry) error {
	return o.ORM.AddLocalRegistry(ctx, &localRegistry)
}
