package utils_test

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/multierr"
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

func TestUtils_IsEmptyAddress(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		addr common.Address
		want bool
	}{
		{"zero address", common.Address{}, true},
		{"non-zero address", testutils.NewAddress(), false},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			actual := utils.IsEmptyAddress(test.addr)
			assert.Equal(t, test.want, actual)
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

func TestUtils_BackoffSleeper(t *testing.T) {
	t.Parallel()

	bs := utils.NewBackoffSleeper()
	assert.Equal(t, time.Duration(0), bs.Duration(), "should initially return immediately")
	bs.Sleep()

	d := 1 * time.Nanosecond
	bs.Min = d
	bs.Factor = 2
	assert.Equal(t, d, bs.Duration())
	bs.Sleep()

	d2 := 2 * time.Nanosecond
	assert.Equal(t, d2, bs.Duration())

	bs.Reset()
	assert.Equal(t, time.Duration(0), bs.Duration(), "should initially return immediately")
}

func TestUtils_DurationFromNow(t *testing.T) {
	t.Parallel()

	future := time.Now().Add(time.Second)
	duration := utils.DurationFromNow(future)
	assert.True(t, 0 < duration)
}

func TestKeccak256(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"basic", "0xf00b", "0x2433bb36d5f9b14e4fea87c2d32d79abfe34e56808b891e471f4400fca2a336c"},
		{"long input", "0xf00b2433bb36d5f9b14e4fea87c2d32d79abfe34e56808b891e471f4400fca2a336c", "0x6b917c56ad7bea7d09132b9e1e29bb5d9aa7d32d067c638dfa886bbbf6874cdf"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input, err := hexutil.Decode(test.input)
			assert.NoError(t, err)
			result, err := utils.Keccak256(input)
			assert.NoError(t, err)

			assert.Equal(t, test.want, hexutil.Encode(result))
		})
	}
}

// From https://github.com/ethereum/EIPs/blob/master/EIPS/eip-55.md#test-cases
var testAddresses = []string{
	"0x52908400098527886E0F7030069857D2E4169EE7",
	"0x8617E340B3D01FA5F11F306F4090FD50E238070D",
	"0xde709f2102306220921060314715629080e2fb77",
	"0x27b1fdb04752bbc536007a920d24acb045561c26",
	"0x5aAeb6053F3E94C9b9A09f33669435E7Ef1BeAed",
	"0xfB6916095ca1df60bB79Ce92cE3Ea74c37c5d359",
	"0xdbF03B407c01E7cD3CBea99509d93f8DDDC8C6FB",
	"0xD1220A0cf47c7B9Be7A2E6BA89F429762e7b9aDb",
}

func TestClient_EIP55CapitalizedAddress(t *testing.T) {
	t.Parallel()

	valid := utils.EIP55CapitalizedAddress
	for _, address := range testAddresses {
		assert.True(t, valid(address))
		assert.False(t, valid(strings.ToLower(address)) &&
			valid(strings.ToUpper(address)))
	}
}

func TestClient_ParseEthereumAddress(t *testing.T) {
	t.Parallel()

	parse := utils.ParseEthereumAddress
	for _, address := range testAddresses {
		a1, err := parse(address)
		assert.NoError(t, err)
		no0xPrefix := address[2:]
		a2, err := parse(no0xPrefix)
		assert.NoError(t, err)
		assert.True(t, a1 == a2)
		_, lowerErr := parse(strings.ToLower(address))
		_, upperErr := parse(strings.ToUpper(address))
		shouldBeError := multierr.Combine(lowerErr, upperErr)
		assert.Error(t, shouldBeError)
		assert.True(t, strings.Contains(shouldBeError.Error(), no0xPrefix))
	}
	_, notHexErr := parse("0xCeci n'est pas une chaîne hexadécimale")
	assert.Error(t, notHexErr)
	_, tooLongErr := parse("0x0123456789abcdef0123456789abcdef0123456789abcdef")
	assert.Error(t, tooLongErr)
}

func TestWaitGroupChan(t *testing.T) {
	t.Parallel()

	wg := &sync.WaitGroup{}
	wg.Add(2)

	ch := utils.WaitGroupChan(wg)

	select {
	case <-ch:
		t.Fatal("should not fire immediately")
	default:
	}

	wg.Done()

	select {
	case <-ch:
		t.Fatal("should not fire until finished")
	default:
	}

	go func() {
		time.Sleep(2 * time.Second)
		wg.Done()
	}()

	cltest.CallbackOrTimeout(t, "WaitGroupChan fires", func() {
		<-ch
	}, 5*time.Second)
}

func TestDependentAwaiter(t *testing.T) {
	t.Parallel()

	da := utils.NewDependentAwaiter()
	da.AddDependents(2)

	select {
	case <-da.AwaitDependents():
		t.Fatal("should not fire immediately")
	default:
	}

	da.DependentReady()

	select {
	case <-da.AwaitDependents():
		t.Fatal("should not fire until finished")
	default:
	}

	go func() {
		time.Sleep(2 * time.Second)
		da.DependentReady()
	}()

	cltest.CallbackOrTimeout(t, "dependents are now ready", func() {
		<-da.AwaitDependents()
	}, 5*time.Second)
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

func TestEVMBytesToUint64(t *testing.T) {
	t.Parallel()

	require.Equal(t, uint64(257), utils.EVMBytesToUint64([]byte{0x01, 0x01}))
	require.Equal(t, uint64(257), utils.EVMBytesToUint64([]byte{0x00, 0x00, 0x01, 0x01}))
	require.Equal(t, uint64(299140445700113), utils.EVMBytesToUint64([]byte{0x00, 0x01, 0x10, 0x11, 0x10, 0x01, 0x00, 0x11}))

	// overflows without erroring
	require.Equal(t, uint64(17), utils.EVMBytesToUint64([]byte{0xFF, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x11}))
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

func Test_StartStopOnce_StopWaitsForStartToFinish(t *testing.T) {
	t.Parallel()

	once := utils.StartStopOnce{}

	ch := make(chan int, 3)

	ready := make(chan bool)

	go func() {
		assert.NoError(t, once.StartOnce("slow service", func() (err error) {
			ch <- 1
			ready <- true
			<-time.After(time.Millisecond * 500) // wait for StopOnce to happen
			ch <- 2

			return nil
		}))
	}()

	go func() {
		<-ready // try stopping halfway through startup
		assert.NoError(t, once.StopOnce("slow service", func() (err error) {
			ch <- 3

			return nil
		}))
	}()

	require.Equal(t, 1, <-ch)
	require.Equal(t, 2, <-ch)
	require.Equal(t, 3, <-ch)
}

func Test_StartStopOnce_MultipleStartNoBlock(t *testing.T) {
	t.Parallel()

	once := utils.StartStopOnce{}

	ch := make(chan int, 3)

	ready := make(chan bool)
	next := make(chan bool)

	go func() {
		ch <- 1
		assert.NoError(t, once.StartOnce("slow service", func() (err error) {
			ready <- true
			<-next // continue after the other StartOnce call fails

			return nil
		}))
		<-next
		ch <- 2

	}()

	go func() {
		<-ready // try starting halfway through startup
		assert.Error(t, once.StartOnce("slow service", func() (err error) {
			return nil
		}))
		next <- true
		ch <- 3
		next <- true

	}()

	require.Equal(t, 1, <-ch)
	require.Equal(t, 3, <-ch) // 3 arrives before 2 because it returns immediately
	require.Equal(t, 2, <-ch)
}

func TestAllEqual(t *testing.T) {
	t.Parallel()

	require.False(t, utils.AllEqual(1, 2, 3, 4, 5))
	require.True(t, utils.AllEqual(1, 1, 1, 1, 1))
	require.False(t, utils.AllEqual(1, 1, 1, 2, 1, 1, 1))
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

func TestUint256ToBytes(t *testing.T) {
	t.Parallel()

	v := big.NewInt(0).Sub(utils.MaxUint256, big.NewInt(1))
	uint256, err := utils.Uint256ToBytes(v)
	assert.NoError(t, err)

	b32 := utils.Uint256ToBytes32(v)
	assert.Equal(t, uint256, b32)

	large := big.NewInt(0).Add(utils.MaxUint256, big.NewInt(1))
	_, err = utils.Uint256ToBytes(large)
	assert.Error(t, err, "too large to convert to uint256")

	negative := big.NewInt(-1)
	assert.Panics(t, func() {
		_, _ = utils.Uint256ToBytes(negative)
	}, "failed to round-trip uint256 back to source big.Int")
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

func TestRetryWithBackoff(t *testing.T) {
	t.Parallel()

	var counter atomic.Int32
	ctx, cancel := context.WithCancel(testutils.Context(t))

	utils.RetryWithBackoff(ctx, func() bool {
		return false
	})

	retry := func() bool {
		return counter.Add(1) < 3
	}

	go utils.RetryWithBackoff(ctx, retry)

	assert.Eventually(t, func() bool {
		return counter.Load() == 3
	}, testutils.WaitTimeout(t), testutils.TestInterval)

	cancel()

	utils.RetryWithBackoff(ctx, retry)
	assert.Equal(t, int32(4), counter.Load())
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

func TestCheckUint256(t *testing.T) {
	t.Parallel()

	large := big.NewInt(0).Add(utils.MaxUint256, big.NewInt(1))
	err := utils.CheckUint256(large)
	assert.Error(t, err, "number out of range for uint256")

	negative := big.NewInt(-123)
	err = utils.CheckUint256(negative)
	assert.Error(t, err, "number out of range for uint256")

	err = utils.CheckUint256(big.NewInt(123))
	assert.NoError(t, err)
}

func TestRandUint256(t *testing.T) {
	t.Parallel()

	for i := 0; i < 1000; i++ {
		uint256 := utils.RandUint256()
		assert.NoError(t, utils.CheckUint256(uint256))
	}
}

func TestHexToUint256(t *testing.T) {
	t.Parallel()

	b, err := utils.HexToUint256("0x00")
	assert.NoError(t, err)
	assert.Zero(t, b.Cmp(big.NewInt(0)))

	b, err = utils.HexToUint256("0xFFFFFFFF")
	assert.NoError(t, err)
	assert.Zero(t, b.Cmp(big.NewInt(4294967295)))
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

func TestContextFromChanWithDeadline(t *testing.T) {
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
		ctx, cancel := utils.ContextFromChanWithDeadline(ch, testutils.TestInterval)
		defer cancel()

		assertCtxCancelled(ctx, t)
	})

	t.Run("stopped", func(t *testing.T) {
		t.Parallel()

		ch := make(chan struct{})
		ctx, cancel := utils.ContextFromChanWithDeadline(ch, testutils.WaitTimeout(t))
		defer cancel()

		ch <- struct{}{}

		assertCtxCancelled(ctx, t)
	})
}

func TestStartStopOnceState_String(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "Unstarted", utils.StartStopOnce_Unstarted.String())
	assert.Equal(t, "Started", utils.StartStopOnce_Started.String())
	assert.Equal(t, "Starting", utils.StartStopOnce_Starting.String())
	assert.Equal(t, "Stopping", utils.StartStopOnce_Stopping.String())
	assert.Equal(t, "Stopped", utils.StartStopOnce_Stopped.String())
	assert.Equal(t, "unrecognized state: 123", utils.StartStopOnceState(123).String())
}

func TestStartStopOnce(t *testing.T) {
	t.Parallel()

	var callsCount atomic.Int32
	incCount := func() {
		callsCount.Add(1)
	}

	var s utils.StartStopOnce
	ok := s.IfStarted(incCount)
	assert.False(t, ok)
	ok = s.IfNotStopped(incCount)
	assert.True(t, ok)
	assert.Equal(t, int32(1), callsCount.Load())

	err := s.StartOnce("foo", func() error { return nil })
	assert.NoError(t, err)

	assert.True(t, s.IfStarted(incCount))
	assert.Equal(t, int32(2), callsCount.Load())

	err = s.StopOnce("foo", func() error { return nil })
	assert.NoError(t, err)
	ok = s.IfNotStopped(incCount)
	assert.False(t, ok)
	assert.Equal(t, int32(2), callsCount.Load())
}

func TestStartStopOnce_StartErrors(t *testing.T) {
	var s utils.StartStopOnce

	err := s.StartOnce("foo", func() error { return errors.New("foo") })
	assert.Error(t, err)

	var callsCount atomic.Int32
	incCount := func() {
		callsCount.Add(1)
	}

	assert.False(t, s.IfStarted(incCount))
	assert.Equal(t, int32(0), callsCount.Load())

	err = s.StartOnce("foo", func() error { return nil })
	require.Error(t, err)
	assert.Contains(t, err.Error(), "foo has already been started once")
	err = s.StopOnce("foo", func() error { return nil })
	require.Error(t, err)
	assert.Contains(t, err.Error(), "foo cannot be stopped from this state; state=StartFailed")

	assert.Equal(t, utils.StartStopOnce_StartFailed, s.LoadState())
}

func TestStartStopOnce_StopErrors(t *testing.T) {
	var s utils.StartStopOnce

	err := s.StartOnce("foo", func() error { return nil })
	require.NoError(t, err)

	var callsCount atomic.Int32
	incCount := func() {
		callsCount.Add(1)
	}

	err = s.StopOnce("foo", func() error { return errors.New("explodey mcsplode") })
	assert.Error(t, err)

	assert.False(t, s.IfStarted(incCount))
	assert.Equal(t, int32(0), callsCount.Load())
	assert.True(t, s.IfNotStopped(incCount))
	assert.Equal(t, int32(1), callsCount.Load())

	err = s.StartOnce("foo", func() error { return nil })
	require.Error(t, err)
	assert.Contains(t, err.Error(), "foo has already been started once")
	err = s.StopOnce("foo", func() error { return nil })
	require.Error(t, err)
	assert.Contains(t, err.Error(), "foo cannot be stopped from this state; state=StopFailed")

	assert.Equal(t, utils.StartStopOnce_StopFailed, s.LoadState())
}

func TestStartStopOnce_Ready_Healthy(t *testing.T) {
	t.Parallel()

	var s utils.StartStopOnce
	assert.Error(t, s.Ready())
	assert.Error(t, s.Healthy())

	err := s.StartOnce("foo", func() error { return nil })
	assert.NoError(t, err)
	assert.NoError(t, s.Ready())
	assert.NoError(t, s.Healthy())
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

func TestTryParseHex(t *testing.T) {
	t.Parallel()

	t.Run("0x prefix missing", func(t *testing.T) {
		t.Parallel()

		_, err := utils.TryParseHex("abcd")
		assert.Error(t, err)
	})

	t.Run("wrong hex characters", func(t *testing.T) {
		t.Parallel()

		_, err := utils.TryParseHex("0xabcdzzz")
		assert.Error(t, err)
	})

	t.Run("valid hex string", func(t *testing.T) {
		t.Parallel()

		b, err := utils.TryParseHex("0x1234")
		assert.NoError(t, err)
		assert.Equal(t, []byte{0x12, 0x34}, b)
	})

	t.Run("prepend odd length with zero", func(t *testing.T) {
		t.Parallel()

		b, err := utils.TryParseHex("0x123")
		assert.NoError(t, err)
		assert.Equal(t, []byte{0x1, 0x23}, b)
	})
}
