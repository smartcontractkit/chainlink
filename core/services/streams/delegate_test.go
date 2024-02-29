package streams

import (
	"testing"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockRegistry struct{}

func (m *mockRegistry) Get(streamID StreamID) (strm Stream, exists bool) { return }
func (m *mockRegistry) Register(streamID StreamID, spec pipeline.Spec, rrs ResultRunSaver) error {
	return nil
}
func (m *mockRegistry) Unregister(streamID StreamID) {}

type mockDelegateConfig struct{}

func (m *mockDelegateConfig) MaxSuccessfulRuns() uint64     { return 0 }
func (m *mockDelegateConfig) ResultWriteQueueDepth() uint64 { return 0 }

func Test_Delegate(t *testing.T) {
	lggr := logger.TestLogger(t)
	registry := &mockRegistry{}
	runner := &mockRunner{}
	cfg := &mockDelegateConfig{}
	d := NewDelegate(lggr, registry, runner, cfg)

	t.Run("ServicesForSpec", func(t *testing.T) {
		jb := job.Job{PipelineSpec: &pipeline.Spec{ID: 1}}
		t.Run("errors if job is missing streamID", func(t *testing.T) {
			_, err := d.ServicesForSpec(testutils.Context(t), jb)
			assert.EqualError(t, err, "streamID is required to be present for stream specs")
		})
		jb.StreamID = ptr(uint32(42))
		t.Run("returns services", func(t *testing.T) {
			srvs, err := d.ServicesForSpec(testutils.Context(t), jb)
			require.NoError(t, err)

			assert.Len(t, srvs, 2)
			assert.IsType(t, &ocrcommon.RunResultSaver{}, srvs[0])

			strmSrv := srvs[1].(*StreamService)
			assert.Equal(t, registry, strmSrv.registry)
			assert.Equal(t, StreamID(42), strmSrv.id)
			assert.Equal(t, jb.PipelineSpec, strmSrv.spec)
			assert.NotNil(t, strmSrv.lggr)
			assert.Equal(t, srvs[0], strmSrv.rrs)
		})
	})
}

func Test_ValidatedStreamSpec(t *testing.T) {
	var tt = []struct {
		name      string
		toml      string
		assertion func(t *testing.T, os job.Job, err error)
	}{
		{
			name: "minimal stream spec",
			toml: `
type               = "stream"
streamID 		   = 12345
name 			   = "voter-turnout"
schemaVersion      = 1
observationSource  = """
ds1          [type=bridge name=voter_turnout];
ds1_parse    [type=jsonparse path="one,two"];
ds1_multiply [type=multiply times=1.23];
ds1 -> ds1_parse -> ds1_multiply -> answer1;
answer1      [type=median index=0];
"""
`,
			assertion: func(t *testing.T, jb job.Job, err error) {
				require.NoError(t, err)
				assert.Equal(t, job.Type("stream"), jb.Type)
				assert.Equal(t, uint32(1), jb.SchemaVersion)
				assert.True(t, jb.Name.Valid)
				require.NotNil(t, jb.StreamID)
				assert.Equal(t, uint32(12345), *jb.StreamID)
				assert.Equal(t, "voter-turnout", jb.Name.String)
			},
		},
		{
			name: "unparseable toml",
			toml: `not toml`,
			assertion: func(t *testing.T, jb job.Job, err error) {
				assert.EqualError(t, err, "toml unmarshal error on job: toml: expected character =")
			},
		},
		{
			name: "invalid field type",
			toml: `
type               = "stream"
name 			   = "voter-turnout"
schemaVersion      = "should be integer"
`,
			assertion: func(t *testing.T, jb job.Job, err error) {
				assert.EqualError(t, err, "toml unmarshal error on job: toml: cannot decode TOML string into struct field job.Job.SchemaVersion of type uint32")
			},
		},
		{
			name: "invalid fields",
			toml: `
type               = "stream"
name 			   = "voter-turnout"
notAValidField     = "some value"
schemaVersion      = 1
`,
			assertion: func(t *testing.T, jb job.Job, err error) {
				assert.EqualError(t, err, "toml unmarshal error on job: strict mode: fields in the document are missing in the target struct")
			},
		},
		{
			name: "wrong type",
			toml: `
type               = "not a valid type"
name 			   = "voter-turnout"
schemaVersion      = 1
observationSource  = """
ds1          [type=bridge name=voter_turnout];
ds1_parse    [type=jsonparse path="one,two"];
ds1_multiply [type=multiply times=1.23];
ds1 -> ds1_parse -> ds1_multiply -> answer1;
answer1      [type=median index=0];
"""
`,
			assertion: func(t *testing.T, jb job.Job, err error) {
				assert.EqualError(t, err, "unsupported type: \"not a valid type\"")
			},
		},
		{
			name: "no error if missing name",
			toml: `
type               = "stream"
schemaVersion      = 1
streamID 		   = 12345
observationSource  = """
ds1          [type=bridge name=voter_turnout];
ds1_parse    [type=jsonparse path="one,two"];
ds1_multiply [type=multiply times=1.23];
ds1 -> ds1_parse -> ds1_multiply -> answer1;
answer1      [type=median index=0];
"""
`,
			assertion: func(t *testing.T, jb job.Job, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "error if missing streamID",
			toml: `
type               = "stream"
schemaVersion      = 1
observationSource  = """
ds1          [type=bridge name=voter_turnout];
ds1_parse    [type=jsonparse path="one,two"];
ds1_multiply [type=multiply times=1.23];
ds1 -> ds1_parse -> ds1_multiply -> answer1;
answer1      [type=median index=0];
"""
`,
			assertion: func(t *testing.T, jb job.Job, err error) {
				assert.EqualError(t, err, "jobs of type 'stream' require streamID to be specified")
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s, err := ValidatedStreamSpec(tc.toml)
			tc.assertion(t, s, err)
		})
	}
}
func ptr[T any](t T) *T { return &t }
