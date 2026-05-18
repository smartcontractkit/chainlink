package utils_test

import (
	"context"
	"encoding/hex"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUtils_NewBytes32ID(t *testing.T) {
	t.Parallel()

	id := utils.NewBytes32ID()
	assert.NotContains(t, id, "-")
}

func TestUtils_NewSecret(t *testing.T) {
	t.Parallel()

	tests := []struct {
		numOfBytes int
		wantStrLen int
	}{
		{12, 16}, {24, 32}, {48, 64}, {96, 128},
	}
	for _, test := range tests {
		test := test

		t.Run(fmt.Sprintf("%d_%d", test.numOfBytes, test.wantStrLen), func(t *testing.T) {
			t.Parallel()

			secret := utils.NewSecret(test.numOfBytes)
			assert.Equal(t, test.wantStrLen, len(secret))
		})
	}
}

func TestUtils_StringToHex(t *testing.T) {
	t.Parallel()

	tests := []struct {
		utf8 string
		hex  string
	}{
		{"abc", "0x616263"},
		{"Hi Mom!", "0x4869204d6f6d21"},
		{"", "0x"},
	}

	for _, test := range tests {
		test := test

		t.Run(test.utf8, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, test.hex, utils.StringToHex(test.utf8))
		})
	}
}

func TestUtils_DurationFromNow(t *testing.T) {
	t.Parallel()

	future := time.Now().Add(time.Second)
	duration := utils.DurationFromNow(future)
	assert.True(t, 0 < duration)
}

func TestBoundedQueue(t *testing.T) {
	t.Parallel()

	q := utils.NewBoundedQueue[int](3)
	require.True(t, q.Empty())
	require.False(t, q.Full())

	q.Add(1)
	require.False(t, q.Empty())
	require.False(t, q.Full())

	x := q.Take()
	require.Equal(t, 1, x)

	require.Zero(t, q.Take())
	require.True(t, q.Empty())
	require.False(t, q.Full())

	q.Add(1)
	q.Add(2)
	q.Add(3)
	q.Add(4)
	require.True(t, q.Full())

	x = q.Take()
	require.Equal(t, 2, x)
	require.False(t, q.Empty())
	require.False(t, q.Full())

	x = q.Take()
	require.Equal(t, 3, x)
	require.False(t, q.Empty())
	require.False(t, q.Full())

	x = q.Take()
	require.Equal(t, 4, x)
	require.True(t, q.Empty())
	require.False(t, q.Full())
}

func TestBoundedPriorityQueue(t *testing.T) {
	t.Parallel()

	q := utils.NewBoundedPriorityQueue[int](map[uint]int{
		1: 3,
		2: 1,
	})
	require.True(t, q.Empty())

	q.Add(1, 1)
	require.False(t, q.Empty())

	x := q.Take()
	require.Equal(t, 1, x)
	require.True(t, q.Empty())

	require.Zero(t, q.Take())
	require.True(t, q.Empty())

	q.Add(2, 1)
	q.Add(1, 2)
	q.Add(1, 3)
	q.Add(1, 4)

	x = q.Take()
	require.Equal(t, 2, x)
	require.False(t, q.Empty())

	x = q.Take()
	require.Equal(t, 3, x)
	require.False(t, q.Empty())

	x = q.Take()
	require.Equal(t, 4, x)
	require.False(t, q.Empty())

	x = q.Take()
	require.Equal(t, 1, x)
	require.True(t, q.Empty())

	require.Zero(t, q.Take())

	q.Add(2, 1)
	q.Add(2, 2)

	x = q.Take()
	require.Equal(t, 2, x)
	require.True(t, q.Empty())

	require.Zero(t, q.Take())
}

func Test_WithJitter(t *testing.T) {
	t.Parallel()

	d := 10 * time.Second

	for i := 0; i < 32; i++ {
		r := utils.WithJitter(d)
		require.GreaterOrEqual(t, int(r), int(9*time.Second))
		require.LessOrEqual(t, int(r), int(11*time.Second))
	}
}

func TestIsEmpty(t *testing.T) {
	t.Parallel()

	b := make([]byte, 32)
	assert.True(t, utils.IsEmpty(b))

	b[10] = 1
	assert.False(t, utils.IsEmpty(b))
}

func TestHashPassword(t *testing.T) {
	t.Parallel()

	h, err := utils.HashPassword("Qwerty123!")
	assert.NoError(t, err)
	assert.NotEmpty(t, h)

	ok := utils.CheckPasswordHash("Qwerty123!", h)
	assert.True(t, ok)

	ok = utils.CheckPasswordHash("God", h)
	assert.False(t, ok)
}

func TestBoxOutput(t *testing.T) {
	t.Parallel()

	output := utils.BoxOutput("some error %d %s", 123, "foo")
	const expected = "\n" +
		"↘↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↙\n" +
		"→                      ←\n" +
		"→  README README       ←\n" +
		"→                      ←\n" +
		"→  some error 123 foo  ←\n" +
		"→                      ←\n" +
		"→  README README       ←\n" +
		"→                      ←\n" +
		"↗↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↖\n" +
		"\n"
	assert.Equal(t, expected, output)
}

func TestISO8601UTC(t *testing.T) {
	t.Parallel()

	ts := time.Unix(1651818206, 0)
	str := utils.ISO8601UTC(ts)
	assert.Equal(t, "2022-05-06T06:23:26Z", str)
}

func TestFormatJSON(t *testing.T) {
	t.Parallel()

	json := `{"foo":123}`
	formatted, err := utils.FormatJSON(json)
	assert.NoError(t, err)
	assert.Equal(t, "\"{\\\"foo\\\":123}\"", string(formatted))
}

func TestMustUnmarshalToMap(t *testing.T) {
	t.Parallel()

	json := `{"foo":123.45}`
	expected := make(map[string]interface{})
	expected["foo"] = 123.45
	m := utils.MustUnmarshalToMap(json)
	assert.Equal(t, expected, m)

	assert.Panics(t, func() {
		utils.MustUnmarshalToMap("")
	})

	assert.Panics(t, func() {
		utils.MustUnmarshalToMap("123")
	})
}

func TestSha256(t *testing.T) {
	t.Parallel()

	hexHash, err := utils.Sha256("test")
	assert.NoError(t, err)

	hash, err := hex.DecodeString(hexHash)
	assert.NoError(t, err)
	assert.Len(t, hash, 32)
}

func TestWithCloseChan(t *testing.T) {
	t.Parallel()

	assertCtxCancelled := func(ctx context.Context, t *testing.T) {
		select {
		case <-ctx.Done():
		case <-time.After(testutils.WaitTimeout(t)):
			assert.FailNow(t, "context was not cancelled")
		}
	}

	t.Run("closing channel", func(t *testing.T) {
		t.Parallel()

		ch := make(chan struct{})
		ctx, cancel := utils.WithCloseChan(testutils.Context(t), ch)
		defer cancel()

		close(ch)

		assertCtxCancelled(ctx, t)
	})

	t.Run("cancelling ctx", func(t *testing.T) {
		t.Parallel()

		ch := make(chan struct{})
		defer close(ch)
		ctx, cancel := utils.WithCloseChan(testutils.Context(t), ch)
		cancel()

		assertCtxCancelled(ctx, t)
	})

	t.Run("cancelling parent ctx", func(t *testing.T) {
		t.Parallel()

		ch := make(chan struct{})
		defer close(ch)
		pctx, pcancel := context.WithCancel(testutils.Context(t))
		ctx, cancel := utils.WithCloseChan(pctx, ch)
		defer cancel()

		pcancel()

		assertCtxCancelled(ctx, t)
	})
}

func TestContextFromChan(t *testing.T) {
	t.Parallel()

	ch := make(chan struct{})
	ctx, cancel := utils.ContextFromChan(ch)
	defer cancel()

	close(ch)

	select {
	case <-ctx.Done():
	case <-time.After(testutils.WaitTimeout(t)):
		assert.FailNow(t, "context was not cancelled")
	}
}

func TestContextFromChanWithTimeout(t *testing.T) {
	t.Parallel()

	assertCtxCancelled := func(ctx context.Context, t *testing.T) {
		select {
		case <-ctx.Done():
		case <-time.After(testutils.WaitTimeout(t)):
			assert.FailNow(t, "context was not cancelled")
		}
	}

	t.Run("small deadline", func(t *testing.T) {
		t.Parallel()

		ch := make(chan struct{})
		ctx, cancel := utils.ContextFromChanWithTimeout(ch, testutils.TestInterval)
		defer cancel()

		assertCtxCancelled(ctx, t)
	})

	t.Run("stopped", func(t *testing.T) {
		t.Parallel()

		ch := make(chan struct{})
		ctx, cancel := utils.ContextFromChanWithTimeout(ch, testutils.WaitTimeout(t))
		defer cancel()

		ch <- struct{}{}

		assertCtxCancelled(ctx, t)
	})
}

func TestLeftPadBitString(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		str      string
		len      int
		expected string
	}{
		{"abc", 10, "0000000abc"},
		{"abc", 0, "abc"},
		{"abc", 2, "abc"},
		{"abc", 3, "abc"},
		{"abc", -10, "abc"},
	} {
		s := utils.LeftPadBitString(test.str, test.len)
		assert.Equal(t, test.expected, s)
	}
}

func TestKeyedMutex(t *testing.T) {
	t.Parallel()

	var km utils.KeyedMutex
	unlock1 := km.LockInt64(1)
	unlock2 := km.LockInt64(2)

	awaiter := cltest.NewAwaiter()
	go func() {
		km.LockInt64(1)()
		km.LockInt64(2)()
		awaiter.ItHappened()
	}()

	unlock2()
	unlock1()
	awaiter.AwaitOrFail(t)
}

func TestValidateCronSchedule(t *testing.T) {
	t.Parallel()

	err := utils.ValidateCronSchedule("")
	assert.Error(t, err)

	err = utils.ValidateCronSchedule("CRON_TZ=UTC 5 * * * *")
	assert.NoError(t, err)

	err = utils.ValidateCronSchedule("@every 1h30m")
	assert.NoError(t, err)

	err = utils.ValidateCronSchedule("@every xyz")
	assert.Error(t, err)
}

func TestPausableTicker(t *testing.T) {
	t.Parallel()

	var counter atomic.Int32

	pt := utils.NewPausableTicker(testutils.TestInterval)
	assert.Nil(t, pt.Ticks())
	defer pt.Destroy()

	followNTicks := func(n int32, awaiter cltest.Awaiter) {
		for range pt.Ticks() {
			if counter.Add(1) == n {
				awaiter.ItHappened()
			}
		}
	}

	pt.Resume()

	wait10 := cltest.NewAwaiter()
	go followNTicks(10, wait10)

	wait10.AwaitOrFail(t)

	pt.Pause()
	time.Sleep(10 * testutils.TestInterval)
	assert.Less(t, counter.Load(), int32(20))
	pt.Resume()

	wait20 := cltest.NewAwaiter()
	go followNTicks(20, wait20)

	wait20.AwaitOrFail(t)
}

func TestCronTicker(t *testing.T) {
	t.Parallel()

	var counter atomic.Int32

	ct, err := utils.NewCronTicker("@every 100ms")
	assert.NoError(t, err)

	awaiter := cltest.NewAwaiter()

	go func() {
		for range ct.Ticks() {
			if counter.Add(1) == 2 {
				awaiter.ItHappened()
			}
		}
	}()

	assert.True(t, ct.Start())
	assert.True(t, ct.Stop())
	assert.Zero(t, counter.Load())

	assert.True(t, ct.Start())

	awaiter.AwaitOrFail(t)

	assert.True(t, ct.Stop())
	c := counter.Load()
	time.Sleep(1 * time.Second)
	assert.Equal(t, c, counter.Load())
}

func TestErrorBuffer(t *testing.T) {
	t.Parallel()

	err1 := errors.New("err1")
	err2 := errors.New("err2")
	err3 := errors.New("err3")

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()
		buff := utils.ErrorBuffer{}
		buff.Append(err1)
		buff.Append(err2)
		combined := buff.Flush()
		errs := utils.UnwrapError(combined)
		assert.Equal(t, 2, len(errs))
		assert.Equal(t, err1.Error(), errs[0].Error())
		assert.Equal(t, err2.Error(), errs[1].Error())
	})

	t.Run("ovewrite oldest error when cap exceeded", func(t *testing.T) {
		t.Parallel()
		buff := utils.ErrorBuffer{}
		buff.SetCap(2)
		buff.Append(err1)
		buff.Append(err2)
		buff.Append(err3)
		combined := buff.Flush()
		errs := utils.UnwrapError(combined)
		assert.Equal(t, 2, len(errs))
		assert.Equal(t, err2.Error(), errs[0].Error())
		assert.Equal(t, err3.Error(), errs[1].Error())
	})

	t.Run("does not overwrite the buffer if cap == 0", func(t *testing.T) {
		t.Parallel()
		buff := utils.ErrorBuffer{}
		for i := 1; i <= 20; i++ {
			buff.Append(errors.Errorf("err#%d", i))
		}

		combined := buff.Flush()
		errs := utils.UnwrapError(combined)
		assert.Equal(t, 20, len(errs))
		assert.Equal(t, "err#20", errs[19].Error())
	})

	t.Run("UnwrapError returns the a single element err array if passed err is not a joinedError", func(t *testing.T) {
		t.Parallel()
		errs := utils.UnwrapError(err1)
		assert.Equal(t, 1, len(errs))
		assert.Equal(t, err1.Error(), errs[0].Error())
	})

	t.Run("flushing an empty err buffer is a nil error", func(t *testing.T) {
		t.Parallel()
		buff := utils.ErrorBuffer{}

		combined := buff.Flush()
		require.Nil(t, combined)
	})
}
