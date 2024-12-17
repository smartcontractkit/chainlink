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

	streamIDs []StreamID
}

func (m *mockPipeline) Run(ctx context.Context) (*pipeline.Run, pipeline.TaskRunResults, error) {
	return m.run, m.trrs, m.err
}

func (m *mockPipeline) StreamIDs() []StreamID {
	return m.streamIDs
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

		// registers new pipeline with multiple stream IDs
		assert.Empty(t, sr.pipelines)
		// err := sr.Register(job.Job{PipelineSpec: &pipeline.Spec{ID: 32, DotDagSource: "source"}}, nil)
		// TODO: what if the dag is unparseable?
		// err := sr.Register(1, pipeline.Spec{ID: 32, DotDagSource: "source"}, nil)
		err := sr.Register(job.Job{ID: 100, Type: job.Stream, PipelineSpec: &pipeline.Spec{ID: 32, DotDagSource: `
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
		assert.Len(t, sr.pipelines, 3) // three streams, one pipeline
		assert.Contains(t, sr.pipelines, StreamID(1))
		assert.Contains(t, sr.pipelines, StreamID(2))
		assert.Contains(t, sr.pipelines, StreamID(3))
		p := sr.pipelines[1]
		assert.Equal(t, p, sr.pipelines[2])
		assert.Equal(t, p, sr.pipelines[3])

		v, exists := sr.Get(1)
		require.True(t, exists)
		msp := v.(*multiStreamPipeline)
		assert.Equal(t, []StreamID{1, 2, 3}, msp.StreamIDs())
		assert.Equal(t, int32(32), msp.spec.ID)

		// errors when attempt to re-register a stream with an existing job ID
		err = sr.Register(job.Job{ID: 100, Type: job.Stream, PipelineSpec: &pipeline.Spec{ID: 33, DotDagSource: `
result1          [type=memo value="900.0022"];
		`}}, nil)
		require.EqualError(t, err, "cannot register job with ID: 100; it is already registered")

		// errors when attempt to register a new job with duplicates stream IDs within ig
		err = sr.Register(job.Job{ID: 101, StreamID: ptr(StreamID(100)), Type: job.Stream, PipelineSpec: &pipeline.Spec{ID: 33, DotDagSource: `
result1          [type=memo value="900.0022" streamID=100];
		`}}, nil)
		require.EqualError(t, err, "cannot register job with ID: 101; invalid stream IDs: duplicate stream ID: 100")

		// errors with unparseable pipeline
		err = sr.Register(job.Job{ID: 101, Type: job.Stream, PipelineSpec: &pipeline.Spec{ID: 33, DotDagSource: "source"}}, nil)
		require.Error(t, err)
		require.EqualError(t, err, "cannot register job with ID: 101; unparseable pipeline: UnmarshalTaskFromMap: unknown task type: \"\"")

		// errors when attempt to re-register a stream with an existing streamID at top-level
		err = sr.Register(job.Job{ID: 101, StreamID: ptr(StreamID(3)), Type: job.Stream, PipelineSpec: &pipeline.Spec{ID: 33, DotDagSource: `
result1          [type=memo value="900.0022"];
multiply2 	  	 [type=multiply times=1 streamID=4 index=0]; // force conversion to decimal
result2          [type=bridge name="foo-bridge" requestData="{\"data\":{\"data\":\"foo\"}}"];
result2_parse    [type=jsonparse path="result" streamID=5 index=1];
result3          [type=bridge name="bar-bridge" requestData="{\"data\":{\"data\":\"bar\"}}"];
result3_parse    [type=jsonparse path="result"];
multiply3 	  	 [type=multiply times=1 streamID=6 index=2]; // force conversion to decimal
result1 -> multiply2;
result2 -> result2_parse;
result3 -> result3_parse -> multiply3;
`}}, nil)
		require.Error(t, err)
		require.EqualError(t, err, "cannot register job with ID: 101; stream id 3 is already registered")

		// errors when attempt to re-register a stream with an existing streamID in DAG
		err = sr.Register(job.Job{ID: 101, StreamID: ptr(StreamID(4)), Type: job.Stream, PipelineSpec: &pipeline.Spec{ID: 33, DotDagSource: `
result1          [type=memo value="900.0022"];
multiply2 	  	 [type=multiply times=1 streamID=1 index=0]; // force conversion to decimal
result2          [type=bridge name="foo-bridge" requestData="{\"data\":{\"data\":\"foo\"}}"];
result2_parse    [type=jsonparse path="result" streamID=5 index=1];
result3          [type=bridge name="bar-bridge" requestData="{\"data\":{\"data\":\"bar\"}}"];
result3_parse    [type=jsonparse path="result"];
multiply3 	  	 [type=multiply times=1 streamID=6 index=2]; // force conversion to decimal
result1 -> multiply2;
result2 -> result2_parse;
result3 -> result3_parse -> multiply3;
`}}, nil)
		require.Error(t, err)
		require.EqualError(t, err, "cannot register job with ID: 101; stream id 1 is already registered")

		// registers new job with all new stream IDs
		err = sr.Register(job.Job{ID: 101, StreamID: ptr(StreamID(4)), Type: job.Stream, PipelineSpec: &pipeline.Spec{ID: 33, DotDagSource: `
result1          [type=memo value="900.0022"];
multiply2 	  	 [type=multiply times=1 streamID=5 index=0]; // force conversion to decimal
result2          [type=bridge name="foo-bridge" requestData="{\"data\":{\"data\":\"foo\"}}"];
result2_parse    [type=jsonparse path="result" streamID=6 index=1];
result3          [type=bridge name="bar-bridge" requestData="{\"data\":{\"data\":\"bar\"}}"];
result3_parse    [type=jsonparse path="result"];
multiply3 	  	 [type=multiply times=1 streamID=7 index=2]; // force conversion to decimal
result1 -> multiply2;
result2 -> result2_parse;
result3 -> result3_parse -> multiply3;
`}}, nil)
		require.NoError(t, err)

		// did not overwrite existing stream
		assert.Len(t, sr.pipelines, 7)
		assert.Equal(t, p, sr.pipelines[1])
		assert.Equal(t, p, sr.pipelines[2])
		assert.Equal(t, p, sr.pipelines[3])
		p2 := sr.pipelines[4]
		assert.NotEqual(t, p, p2)
		assert.Equal(t, p2, sr.pipelines[5])
		assert.Equal(t, p2, sr.pipelines[6])
		assert.Equal(t, p2, sr.pipelines[7])

		v, exists = sr.Get(1)
		require.True(t, exists)
		msp = v.(*multiStreamPipeline)
		assert.ElementsMatch(t, []StreamID{1, 2, 3}, msp.StreamIDs())
		assert.Equal(t, int32(32), msp.spec.ID)

		v, exists = sr.Get(4)
		require.True(t, exists)
		msp = v.(*multiStreamPipeline)
		assert.ElementsMatch(t, []StreamID{4, 5, 6, 7}, msp.StreamIDs())
		assert.Equal(t, int32(33), msp.spec.ID)
	})
	t.Run("Unregister", func(t *testing.T) {
		sr := newRegistry(lggr, runner)

		err := sr.Register(job.Job{ID: 100, StreamID: ptr(StreamID(1)), Type: job.Stream, PipelineSpec: &pipeline.Spec{ID: 33, DotDagSource: `
result1          [type=memo value="900.0022" streamID=2];
		`}}, nil)
		require.NoError(t, err)
		err = sr.Register(job.Job{ID: 101, Type: job.Stream, PipelineSpec: &pipeline.Spec{ID: 33, DotDagSource: `
result1          [type=memo value="900.0022" streamID=3];
		`}}, nil)
		require.NoError(t, err)
		err = sr.Register(job.Job{ID: 102, Type: job.Stream, PipelineSpec: &pipeline.Spec{ID: 33, DotDagSource: `
result1          [type=memo value="900.0022" streamID=4];
		`}}, nil)
		require.NoError(t, err)

		t.Run("unregisters a stream", func(t *testing.T) {
			assert.Len(t, sr.pipelines, 4)

			sr.Unregister(100)

			assert.Len(t, sr.pipelines, 2)
			_, exists := sr.pipelines[1]
			assert.False(t, exists)
			_, exists = sr.pipelines[2]
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
