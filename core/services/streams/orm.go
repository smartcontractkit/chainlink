package streams

import (
	"context"

	"github.com/smartcontractkit/chainlink-data-streams/streams"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

type ORM interface {
	StreamCacheORM
}

var _ ORM = &orm{}

type orm struct {
	q pg.Queryer
}

func NewORM(q pg.Queryer) ORM {
	return &orm{q}
}

func (o *orm) LoadStreams(ctx context.Context, lggr logger.Logger, runner Runner, m map[streams.StreamID]Stream) error {
	rows, err := o.q.QueryContext(ctx, "SELECT s.id, ps.id, ps.dot_dag_source, ps.max_task_duration FROM streams s JOIN pipeline_specs ps ON ps.id = s.pipeline_spec_id")
	if err != nil {
		// TODO: retries?
		return err
	}

	for rows.Next() {
		var strm stream
		if err := rows.Scan(&strm.id, &strm.spec.ID, &strm.spec.DotDagSource, &strm.spec.MaxTaskDuration); err != nil {
			return err
		}
		strm.lggr = lggr.Named("Stream").With("streamID", strm.id)
		strm.runner = runner

		m[strm.id] = &strm
	}
	return rows.Err()
}
