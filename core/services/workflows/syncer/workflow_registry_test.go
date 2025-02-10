package syncer

import (
	"context"
	"errors"
	"iter"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type mockContractReader struct {
	event     types.Sequence
	eventType string

	glvError func() error
}

func (mcr *mockContractReader) Start(ctx context.Context) error {
	return nil
}

func (mcr *mockContractReader) Close() error {
	return nil
}

func (mcr *mockContractReader) Bind(ctx context.Context, bc []types.BoundContract) error {
	return nil
}

func (mcr *mockContractReader) QueryKeys(ctx context.Context, keyQueries []types.ContractKeyFilter, limitAndSort query.LimitAndSort) (iter.Seq2[string, types.Sequence], error) {
	return func(yield func(string, types.Sequence) bool) {
		yield(mcr.eventType, mcr.event)
	}, nil
}

func (mcr *mockContractReader) GetLatestValueWithHeadData(ctx context.Context, readName string, confidenceLevel primitives.ConfidenceLevel, params any, returnVal any) (head *types.Head, err error) {
	err = mcr.glvError()
	return &types.Head{Height: "0"}, err

}

type crMock struct {
	cnt        int
	mockReader ContractReader
}

func (c *crMock) newContractReaderFn(ctx context.Context, b []byte) (ContractReader, error) {
	c.cnt++
	if c.cnt > 1 {
		return c.mockReader, nil
	}

	return nil, errors.New("not initialized")
}

type handler struct {
	onHandler chan struct{}
	event     Event
}

func (h *handler) Handle(ctx context.Context, event Event) error {
	h.onHandler <- struct{}{}
	h.event = event
	return nil
}

type mockDonNotifier struct {
}

func (d *mockDonNotifier) WaitForDon(ctx context.Context) (capabilities.DON, error) {
	return capabilities.DON{
		ID: 1,
	}, nil
}

func TestWorkflowRegistry_Start_ContinuesTryingToFetchContractReader(t *testing.T) {
	event, err := values.NewMap(map[string]any{
		"WorkflowName": "aworkflow",
		"WorkflowID":   [32]byte{0: 1},
		"BinaryURL":    "http://a-url.com",
	})
	require.NoError(t, err)

	eventSeq := types.Sequence{
		Cursor: "1_abc",
		Head: types.Head{
			Height: "1",
		},
		Data: event,
	}

	expectedEvent := &WorkflowRegistryEvent{
		Cursor: "1_abc",
		Data: WorkflowRegistryWorkflowRegisteredV1{
			WorkflowID:   [32]byte{0: 1},
			WorkflowName: "aworkflow",
			BinaryURL:    "http://a-url.com",
		},
		EventType: "WorkflowRegisteredV1",
		Head: Head{
			Height: "1",
		},
	}

	m := crMock{
		mockReader: &mockContractReader{
			event:     eventSeq,
			eventType: string(WorkflowRegisteredEvent),
			glvError:  func() error { return nil },
		},
	}
	onHandler := make(chan struct{}, 10)
	h := &handler{onHandler: onHandler}
	wr := NewWorkflowRegistry(
		logger.TestLogger(t),
		m.newContractReaderFn,
		"0x0",
		WorkflowEventPollerConfig{QueryCount: 10},
		h,
		&mockDonNotifier{},
	)

	wr.ticker = time.NewTicker(10 * time.Millisecond).C

	err = wr.Start(tests.Context(t))
	require.NoError(t, err)
	defer wr.Close()

	<-h.onHandler

	assert.Equal(t, expectedEvent, h.event)
}

func TestWorkflowRegistry_Start_ContinuesTryingToLoadWorkflows(t *testing.T) {
	event, err := values.NewMap(map[string]any{
		"WorkflowName": "aworkflow",
		"WorkflowID":   [32]byte{0: 1},
		"BinaryURL":    "http://a-url.com",
	})
	require.NoError(t, err)

	eventSeq := types.Sequence{
		Cursor: "1_abc",
		Head: types.Head{
			Height: "1",
		},
		Data: event,
	}

	expectedEvent := &WorkflowRegistryEvent{
		Cursor: "1_abc",
		Data: WorkflowRegistryWorkflowRegisteredV1{
			WorkflowID:   [32]byte{0: 1},
			WorkflowName: "aworkflow",
			BinaryURL:    "http://a-url.com",
		},
		EventType: "WorkflowRegisteredV1",
		Head: Head{
			Height: "1",
		},
	}

	loadWorkflowsErrorCount := 0
	m := crMock{
		mockReader: &mockContractReader{
			event:     eventSeq,
			eventType: string(WorkflowRegisteredEvent),
			glvError: func() error {
				loadWorkflowsErrorCount++
				if loadWorkflowsErrorCount < 1 {
					return errors.New("could not call get latest value")
				}
				return nil
			},
		},
	}
	onHandler := make(chan struct{}, 10)
	h := &handler{onHandler: onHandler}
	wr := NewWorkflowRegistry(
		logger.TestLogger(t),
		m.newContractReaderFn,
		"0x0",
		WorkflowEventPollerConfig{QueryCount: 10},
		h,
		&mockDonNotifier{},
	)

	wr.ticker = time.NewTicker(10 * time.Millisecond).C

	err = wr.Start(tests.Context(t))
	require.NoError(t, err)
	defer wr.Close()

	<-h.onHandler

	assert.Equal(t, expectedEvent, h.event)
}
