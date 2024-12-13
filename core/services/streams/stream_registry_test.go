package streams

import (
	"context"
	"testing"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var _ Pipeline = &mockPipeline{}

type mockPipeline struct {
	run  *pipeline.Run
	trrs pipeline.TaskRunResults
	err  error
}

func (m *mockPipeline) Run(ctx context.Context) (*pipeline.Run, pipeline.TaskRunResults, error) {
	return m.run, m.trrs, m.err
}

func (m *mockPipeline) StreamIDs() []StreamID {
	return nil
}

func Test_Registry(t *testing.T) {
	lggr := logger.TestLogger(t)
	runner := &mockRunner{}

	t.Run("Get", func(t *testing.T) {
		sr := newRegistry(lggr, runner)

		sr.pipelines[1] = &mockPipeline{run: &pipeline.Run{ID: 1}}
		sr.pipelines[2] = &mockPipeline{run: &pipeline.Run{ID: 2}}
		sr.pipelines[3] = &mockPipeline{run: &pipeline.Run{ID: 3}}

		v, exists := sr.Get(1)
		assert.True(t, exists)
		assert.Equal(t, sr.pipelines[1], v)

		v, exists = sr.Get(2)
		assert.True(t, exists)
		assert.Equal(t, sr.pipelines[2], v)

		v, exists = sr.Get(3)
		assert.True(t, exists)
		assert.Equal(t, sr.pipelines[3], v)

		v, exists = sr.Get(4)
		assert.Nil(t, v)
		assert.False(t, exists)
	})
	t.Run("Register", func(t *testing.T) {
		sr := newRegistry(lggr, runner)

		t.Run("registers new pipeline with multiple stream IDs", func(t *testing.T) {
			assert.Len(t, sr.pipelines, 0)
			// err := sr.Register(job.Job{PipelineSpec: &pipeline.Spec{ID: 32, DotDagSource: "source"}}, nil)
			// TODO: what if the dag is unparseable?
			// err := sr.Register(1, pipeline.Spec{ID: 32, DotDagSource: "source"}, nil)
			err := sr.Register(job.Job{PipelineSpec: &pipeline.Spec{ID: 32, DotDagSource: `
result1          [type=memo value="900.0022"];
multiply2 	  	 [type=multiply times=1 streamID=1 index=0]; // force conversion to decimal
result2          [type=bridge name="foo-bridge" requestData="{\"data\":{\"data\":\"foo\"}}"];
result2_parse    [type=jsonparse path="result" streamID=2 index=1];
result3          [type=bridge name="bar-bridge" requestData="{\"data\":{\"data\":\"bar\"}}"];
result3_parse    [type=jsonparse path="result"];
multiply3 	  	 [type=multiply times=1 streamID=3 index=2]; // force conversion to decimal
result1 -> multiply2;
result2 -> result2_parse;
result3 -> result3_parse -> multiply3;
`}}, nil)
			require.NoError(t, err)
			assert.Len(t, sr.pipelines, 1)

			v, exists := sr.Get(1)
			require.True(t, exists)
			msp := v.(*multiStreamPipeline)
			assert.Equal(t, "foo", msp.StreamIDs())
			assert.Equal(t, int32(32), msp.spec.ID)
		})

		t.Run("errors when attempt to re-register a stream with an existing ID", func(t *testing.T) {
			assert.Len(t, sr.pipelines, 1)
			err := sr.Register(1, pipeline.Spec{ID: 33, DotDagSource: "source"}, nil)
			require.Error(t, err)
			assert.Len(t, sr.pipelines, 1)
			assert.EqualError(t, err, "stream already registered for id: 1")

			v, exists := sr.Get(1)
			require.True(t, exists)
			msp := v.(*multiStreamPipeline)
			assert.Equal(t, StreamID(1), msp.id)
			assert.Equal(t, int32(32), msp.spec.ID)
		})
	})
	t.Run("Unregister", func(t *testing.T) {
		sr := newRegistry(lggr, runner)

		sr.pipelines[1] = &mockPipeline{run: &pipeline.Run{ID: 1}}
		sr.pipelines[2] = &mockPipeline{run: &pipeline.Run{ID: 2}}
		sr.pipelines[3] = &mockPipeline{run: &pipeline.Run{ID: 3}}

		t.Run("unregisters a stream", func(t *testing.T) {
			assert.Len(t, sr.pipelines, 3)

			sr.Unregister(1)

			assert.Len(t, sr.pipelines, 2)
			_, exists := sr.pipelines[1]
			assert.False(t, exists)
		})
		t.Run("no effect when unregistering a non-existent stream", func(t *testing.T) {
			assert.Len(t, sr.pipelines, 2)

			sr.Unregister(1)

			assert.Len(t, sr.pipelines, 2)
			_, exists := sr.pipelines[1]
			assert.False(t, exists)
		})
	})
}
