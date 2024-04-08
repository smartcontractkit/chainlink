package medianpoc

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/jsonserializable"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink/v2/core/logger"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

type mockPipelineRunner struct {
	results types.TaskResults
	err     error
	spec    string
	vars    types.Vars
	options types.Options
}

func (m *mockPipelineRunner) ExecuteRun(ctx context.Context, spec string, vars types.Vars, options types.Options) (types.TaskResults, error) {
	m.spec = spec
	m.vars = vars
	m.options = options
	return m.results, m.err
}

func TestDataSource(t *testing.T) {
	lggr := logger.TestLogger(t)
	expect := jsonserializable.JSONSerializable{Val: int64(3), Valid: true}
	pr := &mockPipelineRunner{
		results: types.TaskResults{
			{
				TaskValue: types.TaskValue{
					Value:      expect,
					Error:      nil,
					IsTerminal: true,
				},
				Index: 2,
			},
			{
				TaskValue: types.TaskValue{
					Value:      jsonserializable.JSONSerializable{Val: int64(4), Valid: true},
					Error:      nil,
					IsTerminal: false,
				},
				Index: 1,
			},
		},
	}
	spec := "SPEC"
	ds := &DataSource{
		pipelineRunner: pr,
		spec:           spec,
		lggr:           lggr,
	}
	res, err := ds.Observe(tests.Context(t), ocrtypes.ReportTimestamp{})
	require.NoError(t, err)
	expectBN := big.NewInt(expect.Val.(int64))
	assert.Equal(t, expectBN, res)
	assert.Equal(t, spec, pr.spec)
	assert.Equal(t, expectBN, ds.current.LatestAnswer)
}

func TestDataSource_ResultErrors(t *testing.T) {
	lggr := logger.TestLogger(t)
	pr := &mockPipelineRunner{
		results: types.TaskResults{
			{
				TaskValue: types.TaskValue{
					Error:      errors.New("something went wrong"),
					IsTerminal: true,
				},
				Index: 0,
			},
		},
	}
	spec := "SPEC"
	ds := &DataSource{
		pipelineRunner: pr,
		spec:           spec,
		lggr:           lggr,
	}
	_, err := ds.Observe(tests.Context(t), ocrtypes.ReportTimestamp{})
	assert.ErrorContains(t, err, "something went wrong")
}

func TestDataSource_ResultNotAnInt(t *testing.T) {
	lggr := logger.TestLogger(t)

	expect := jsonserializable.JSONSerializable{Val: "string-result", Valid: true}
	pr := &mockPipelineRunner{
		results: types.TaskResults{
			{
				TaskValue: types.TaskValue{
					Value:      expect,
					IsTerminal: true,
				},
				Index: 0,
			},
		},
	}
	spec := "SPEC"
	ds := &DataSource{
		pipelineRunner: pr,
		spec:           spec,
		lggr:           lggr,
	}
	_, err := ds.Observe(tests.Context(t), ocrtypes.ReportTimestamp{})
	assert.ErrorContains(t, err, "cannot convert observation to decimal")
}
