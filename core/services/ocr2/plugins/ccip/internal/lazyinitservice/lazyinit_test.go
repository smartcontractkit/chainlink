package lazyinitservice

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/test-go/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

var errInit = errors.New("boom")

var errDummyStart = errors.New("dummy start error")
var errDummyClose = errors.New("dummy close error")

type dummyService struct {
	startError     error
	closeError     error
	startCallCount int
	completeStart  chan struct{}
	closeCallCount int
}

func newDummyService() *dummyService {
	return &dummyService{
		completeStart: make(chan struct{}, 1),
	}
}

func (s *dummyService) AwaitCompleteStart() {
	<-s.completeStart
}

func (s *dummyService) Start(context.Context) error {
	s.startCallCount++
	s.completeStart <- struct{}{}
	return s.startError
}

func (s *dummyService) Close() error {
	s.closeCallCount++
	return s.closeError
}

func TestLazyInitService_AsyncInit(t *testing.T) {
	dummy := newDummyService()

	s := New(logger.NullLogger, "dummy", func(context.Context) (job.ServiceCtx, error) {
		return dummy, nil
	})
	assert.Equal(t, "dummy", s.Name())
	assert.NoError(t, s.Start(context.Background()))
	dummy.AwaitCompleteStart()
	assert.Equal(t, nil, s.Ready())
	assert.NoError(t, s.Close())
	assert.Equal(t, 1, dummy.startCallCount)
	assert.Equal(t, 1, dummy.closeCallCount)
	assert.Equal(t, ErrClosed, s.Ready())
}

func TestLazyInitService_NoStartOnUnrecoverableFailure(t *testing.T) {
	t.Parallel()
	tries := 0
	ch := make(chan struct{})
	s := New(logger.NullLogger, "dummy", func(context.Context) (job.ServiceCtx, error) {
		tries++
		close(ch)
		return nil, Unrecoverable(errInit)
	})
	assert.NoError(t, s.Start(context.Background()))
	<-ch
	assert.True(t, errors.Is(s.Ready(), errInit), "expected %v, got %v", errInit, s.Ready())
	assert.NoError(t, s.Close())
	assert.Equal(t, 1, tries)
}

func TestLazyInitService_RetryOnRecoverableFailure(t *testing.T) {
	tries := 0
	var errs []error
	dummy := newDummyService()
	s := New(logger.NullLogger, "dummy", func(context.Context) (job.ServiceCtx, error) {
		tries++
		if tries <= 3 {
			return nil, errInit
		}
		return dummy, nil
	}, WithLogErrorFunc(func(msg error) { errs = append(errs, msg) }))
	assert.NoError(t, s.Start(context.Background()))
	dummy.AwaitCompleteStart()
	assert.NoError(t, s.Close())
	assert.Equal(t, 1, dummy.startCallCount)
	assert.Equal(t, 1, dummy.closeCallCount)
	assert.Equal(t, 4, tries)
	assert.Equal(t, 3, len(errs))
}

func TestLazyInitService_ParentContextCancel(t *testing.T) {
	dummy := newDummyService()
	s := New(logger.NullLogger, "dummy", func(context.Context) (job.ServiceCtx, error) {
		return dummy, nil
	})
	ctx, cancelFunc := context.WithCancel(context.Background())
	assert.NoError(t, s.Start(ctx))
	cancelFunc()
	dummy.AwaitCompleteStart()

	assert.NoError(t, s.Close())
	assert.Equal(t, 1, dummy.startCallCount)
	assert.Equal(t, 1, dummy.closeCallCount)
}

func TestLazyInitService_FaultyInitFunction(t *testing.T) {
	var errs []error
	var wg sync.WaitGroup

	wg.Add(1)
	s := New(logger.NullLogger, "dummy", func(context.Context) (job.ServiceCtx, error) {
		return nil, nil
	}, WithLogErrorFunc(func(err error) {
		defer wg.Done()
		errs = append(errs, err)
	}))

	assert.NoError(t, s.Start(context.Background()))
	wg.Wait()
	require.Equal(t, 1, len(errs))
	assert.Equal(t, ErrNoService, errs[0])
}

func TestLazyInitService_ReportStartErrors(t *testing.T) {
	dummy := newDummyService()
	dummy.startError = errDummyStart

	var errs []error
	var wg sync.WaitGroup

	wg.Add(1)
	s := New(logger.NullLogger, "dummy", func(context.Context) (job.ServiceCtx, error) {
		return dummy, nil
	}, WithLogErrorFunc(func(err error) {
		defer wg.Done()
		errs = append(errs, err)
	}))

	assert.NoError(t, s.Start(context.Background()))
	wg.Wait()
	assert.Equal(t, 1, len(errs))
	assert.True(t, errors.Is(errs[0], errDummyStart), "expected %v, got %v", errDummyStart, errs[0])
}

func TestLazyInitService_ReportCloseErrors(t *testing.T) {
	dummy := newDummyService()
	dummy.closeError = errDummyClose

	s := New(logger.NullLogger, "dummy", func(context.Context) (job.ServiceCtx, error) {
		return dummy, nil
	})

	assert.NoError(t, s.Start(context.Background()))
	dummy.AwaitCompleteStart()
	assert.Equal(t, errDummyClose, s.Close())
	assert.Equal(t, errDummyClose, s.Ready())
}
