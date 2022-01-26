package monitoring

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestPoller(t *testing.T) {
	for _, testCase := range []struct {
		name           string
		duration       time.Duration
		waitOnRead     time.Duration
		pollInterval   time.Duration
		readTimeout    time.Duration
		processingTime time.Duration
		bufferCapacity uint32
		countLower     int
		countUpper     int
	}{
		{
			"non-overlapping polls, no buffering",
			1 * time.Second,
			100 * time.Millisecond,
			100 * time.Millisecond,
			100 * time.Millisecond,
			0,
			0,
			4,
			5,
		},
		{
			"slow fetching, quick polling, no buffering",
			1 * time.Second,
			300 * time.Millisecond,
			10 * time.Millisecond,
			10 * time.Millisecond,
			0,
			0,
			28,
			35,
		},
		{
			"fast fetch, fast polling, insufficient buffering, tons of backpressure",
			1 * time.Second,
			10 * time.Millisecond, // Producer will make 1000/(10+10)=50 messages in a second.
			10 * time.Millisecond,
			10 * time.Millisecond,
			200 * time.Millisecond, // time it gets the "consumer" to process a message. It will only be able to process 1000/200=5 updates per second.
			5,
			4,
			5,
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			defer goleak.VerifyNone(t)
			ctx, cancel := context.WithTimeout(context.Background(), testCase.duration)
			defer cancel()
			source := &fakeSourceWithWait{testCase.waitOnRead}
			poller := NewSourcePoller(
				source,
				newNullLogger(),
				testCase.pollInterval,
				testCase.readTimeout,
				testCase.bufferCapacity)
			go poller.Run(ctx)
			readCount := 0

		COUNTER:
			for {
				select {
				case <-poller.Updates():
					select {
					case <-time.After(testCase.processingTime):
						readCount += 1
					case <-ctx.Done():
						break COUNTER
					}
				case <-ctx.Done():
					break COUNTER
				}
			}
			require.GreaterOrEqual(t, readCount, testCase.countLower)
			require.LessOrEqual(t, readCount, testCase.countUpper)
		})
	}
	t.Run("resumes after a source error", func(t *testing.T) {
		defer goleak.VerifyNone(t)
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()
		source := &fakeSourceWithError{make(chan interface{}), make(chan error)}
		poller := NewSourcePoller(
			source,
			newNullLogger(),
			10*time.Millisecond, // poll interval
			10*time.Millisecond, // read timeout
			0,                   // buffer capacity
		)
		go poller.Run(ctx)

		source.updates <- "update1"
		require.Equal(t, "update1", <-poller.Updates())
		source.errors <- fmt.Errorf("error1")
		source.updates <- "update2"
		require.Equal(t, "update2", <-poller.Updates())
		source.errors <- fmt.Errorf("error2")
		source.updates <- "update3"
		require.Equal(t, "update3", <-poller.Updates())
		source.errors <- fmt.Errorf("error3")
		source.updates <- "update4"
		require.Equal(t, "update4", <-poller.Updates())
	})
}

type fakeSourceWithWait struct {
	waitOnRead time.Duration
}

func (f *fakeSourceWithWait) Fetch(ctx context.Context) (interface{}, error) {
	select {
	case <-time.After(f.waitOnRead):
		return 1, nil
	case <-ctx.Done():
		return 0, nil
	}
}

type fakeSourceWithError struct {
	updates chan interface{}
	errors  chan error
}

func (f *fakeSourceWithError) Fetch(ctx context.Context) (interface{}, error) {
	select {
	case update := <-f.updates:
		return update, nil
	case err := <-f.errors:
		return nil, err
	case <-ctx.Done():
		return nil, nil
	}
}
