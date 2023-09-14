package streams

import (
	"context"
	"testing"

	"github.com/smartcontractkit/chainlink-data-streams/streams"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"

	"github.com/stretchr/testify/assert"
)

type mockRunner struct{}

func (m *mockRunner) ExecuteRun(ctx context.Context, spec pipeline.Spec, vars pipeline.Vars, l logger.Logger) (run *pipeline.Run, trrs pipeline.TaskRunResults, err error) {
	return
}

func Test_ORM(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	orm := NewORM(db)
	ctx := testutils.Context(t)
	lggr := logger.TestLogger(t)
	runner := &mockRunner{}

	t.Run("LoadStreams", func(t *testing.T) {
		t.Run("nothing in database", func(t *testing.T) {
			m := make(map[streams.StreamID]Stream)
			err := orm.LoadStreams(ctx, lggr, runner, m)
			assert.NoError(t, err)

			assert.Len(t, m, 0)
		})
		t.Run("loads streams from database", func(t *testing.T) {
			pgtest.MustExec(t, db, `
WITH pipeline_specs AS (
	INSERT INTO pipeline_specs (dot_dag_source, created_at) VALUES
	('foo', NOW()),
	('bar', NOW()),
	('baz', NOW())
	RETURNING id, dot_dag_source
)
INSERT INTO streams(id, pipeline_spec_id, created_at)
SELECT CONCAT('stream-', pipeline_specs.dot_dag_source), pipeline_specs.id, NOW()
FROM pipeline_specs
`)

			m := make(map[streams.StreamID]Stream)
			err := orm.LoadStreams(ctx, lggr, runner, m)
			assert.NoError(t, err)

			assert.Len(t, m, 3)
			assert.Contains(t, m, streams.StreamID("stream-foo"))
			assert.Contains(t, m, streams.StreamID("stream-bar"))
			assert.Contains(t, m, streams.StreamID("stream-baz"))

			// test one of the streams to ensure it got loaded correctly
			s := m["stream-foo"].(*stream)
			assert.Equal(t, streams.StreamID("stream-foo"), s.id)
			assert.NotNil(t, s.lggr)
			assert.Equal(t, "foo", s.spec.DotDagSource)
			assert.NotZero(t, s.spec.ID)
			assert.NotNil(t, s.runner)
			assert.Equal(t, runner, s.runner)
		})
	})
}
