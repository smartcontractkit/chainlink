package monitoring

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
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
