package utils_test

import (
	"context"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
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
	for _, tt := range tests {
		test := tt
		secret := utils.NewSecret(test.numOfBytes)
		assert.Equal(t, test.wantStrLen, len(secret))
	}
}

func TestUtils_IsEmptyAddress(t *testing.T) {
	tests := []struct {
		name string
		addr common.Address
		want bool
	}{
		{"zero address", common.Address{}, true},
		{"non-zero address", cltest.NewAddress(), false},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			actual := utils.IsEmptyAddress(test.addr)
			assert.Equal(t, test.want, actual)
		})
	}
}

func TestUtils_StringToHex(t *testing.T) {
	tests := []struct {
		utf8 string
		hex  string
	}{
		{"abc", "0x616263"},
		{"Hi Mom!", "0x4869204d6f6d21"},
		{"", "0x"},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.utf8, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.hex, utils.StringToHex(test.utf8))
		})
	}
}

func TestUtils_BackoffSleeper(t *testing.T) {
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
	valid := utils.EIP55CapitalizedAddress
	for _, address := range testAddresses {
		assert.True(t, valid(address))
		assert.False(t, valid(strings.ToLower(address)) &&
			valid(strings.ToUpper(address)))
	}
}

func TestClient_ParseEthereumAddress(t *testing.T) {
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

func TestMaxUint32(t *testing.T) {
	tests := []struct {
		name        string
		expectation uint32
		vals        []uint32
	}{
		{"single", 9, []uint32{9}},
		{"positives", 5, []uint32{3, 4, 5}},
		{"equal", 3, []uint32{3, 3, 3}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := utils.MaxUint32(test.vals[0], test.vals[1:len(test.vals)]...)
			assert.Equal(t, test.expectation, actual)
		})
	}
}

func TestMaxInt(t *testing.T) {
	tests := []struct {
		name        string
		expectation int
		vals        []int
	}{
		{"negatives", -1, []int{-1, -2, -9}},
		{"positives", 5, []int{3, 4, 5}},
		{"both", 5, []int{-1, -2, -9, 3, 4, 5}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := utils.MaxInt(test.vals[0], test.vals[1:len(test.vals)]...)
			assert.Equal(t, test.expectation, actual)
		})
	}
}

func TestMinUint(t *testing.T) {
	tests := []struct {
		name        string
		expectation uint
		vals        []uint
	}{
		{"single", 9, []uint{9}},
		{"positives", 3, []uint{3, 4, 5}},
		{"equal", 3, []uint{3, 3, 3}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := utils.MinUint(test.vals[0], test.vals[1:len(test.vals)]...)
			assert.Equal(t, test.expectation, actual)
		})
	}
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

	q := utils.NewBoundedQueue(3)
	require.True(t, q.Empty())
	require.False(t, q.Full())

	q.Add(1)
	require.False(t, q.Empty())
	require.False(t, q.Full())

	x := q.Take().(int)
	require.Equal(t, 1, x)

	iface := q.Take()
	require.Nil(t, iface)
	require.True(t, q.Empty())
	require.False(t, q.Full())

	q.Add(1)
	q.Add(2)
	q.Add(3)
	q.Add(4)
	require.True(t, q.Full())

	x = q.Take().(int)
	require.Equal(t, 2, x)
	require.False(t, q.Empty())
	require.False(t, q.Full())

	x = q.Take().(int)
	require.Equal(t, 3, x)
	require.False(t, q.Empty())
	require.False(t, q.Full())

	x = q.Take().(int)
	require.Equal(t, 4, x)
	require.True(t, q.Empty())
	require.False(t, q.Full())
}

func TestBoundedPriorityQueue(t *testing.T) {
	t.Parallel()

	q := utils.NewBoundedPriorityQueue(map[uint]uint{
		1: 3,
		2: 1,
	})
	require.True(t, q.Empty())

	q.Add(1, 1)
	require.False(t, q.Empty())

	x := q.Take().(int)
	require.Equal(t, 1, x)
	require.True(t, q.Empty())

	iface := q.Take()
	require.Nil(t, iface)
	require.True(t, q.Empty())

	q.Add(2, 1)
	q.Add(1, 2)
	q.Add(1, 3)
	q.Add(1, 4)

	x = q.Take().(int)
	require.Equal(t, 2, x)
	require.False(t, q.Empty())

	x = q.Take().(int)
	require.Equal(t, 3, x)
	require.False(t, q.Empty())

	x = q.Take().(int)
	require.Equal(t, 4, x)
	require.False(t, q.Empty())

	x = q.Take().(int)
	require.Equal(t, 1, x)
	require.True(t, q.Empty())

	iface = q.Take()
	require.Nil(t, iface)

	q.Add(2, 1)
	q.Add(2, 2)

	x = q.Take().(int)
	require.Equal(t, 2, x)
	require.True(t, q.Empty())

	iface = q.Take()
	require.Nil(t, iface)
}

func TestEVMBytesToUint64(t *testing.T) {
	require.Equal(t, uint64(257), utils.EVMBytesToUint64([]byte{0x01, 0x01}))
	require.Equal(t, uint64(257), utils.EVMBytesToUint64([]byte{0x00, 0x00, 0x01, 0x01}))
	require.Equal(t, uint64(299140445700113), utils.EVMBytesToUint64([]byte{0x00, 0x01, 0x10, 0x11, 0x10, 0x01, 0x00, 0x11}))

	// overflows without erroring
	require.Equal(t, uint64(17), utils.EVMBytesToUint64([]byte{0xFF, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x11}))
}

func TestCombinedContext(t *testing.T) {
	t.Run("cancels when an inner context is canceled", func(t *testing.T) {
		innerCtx, innerCancel := context.WithCancel(context.Background())
		defer innerCancel()

		chStop := make(chan struct{})

		ctx, cancel := utils.CombinedContext(innerCtx, chStop, 1*time.Hour)
		defer cancel()

		innerCancel()

		select {
		case <-ctx.Done():
		case <-time.After(5 * time.Second):
			t.Fatal("context didn't cancel")
		}
	})

	t.Run("cancels when a channel is closed", func(t *testing.T) {
		innerCtx, innerCancel := context.WithCancel(context.Background())
		defer innerCancel()

		chStop := make(chan struct{})

		ctx, cancel := utils.CombinedContext(innerCtx, chStop, 1*time.Hour)
		defer cancel()

		close(chStop)

		select {
		case <-ctx.Done():
		case <-time.After(5 * time.Second):
			t.Fatal("context didn't cancel")
		}
	})

	t.Run("cancels when a duration elapses", func(t *testing.T) {
		innerCtx, innerCancel := context.WithCancel(context.Background())
		defer innerCancel()

		chStop := make(chan struct{})

		ctx, cancel := utils.CombinedContext(innerCtx, chStop, 1*time.Second)
		defer cancel()

		select {
		case <-ctx.Done():
		case <-time.After(5 * time.Second):
			t.Fatal("context didn't cancel")
		}
	})

	t.Run("doesn't cancel if none of its children cancel", func(t *testing.T) {
		innerCtx, innerCancel := context.WithCancel(context.Background())
		defer innerCancel()

		chStop := make(chan struct{})

		ctx, cancel := utils.CombinedContext(innerCtx, chStop, 1*time.Hour)
		defer cancel()

		select {
		case <-ctx.Done():
			t.Fatal("context canceled")
		case <-time.After(5 * time.Second):
		}
	})
}

func Test_WithJitter(t *testing.T) {
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
		once.StartOnce("slow service", func() (err error) {
			ch <- 1
			ready <- true
			<-time.After(time.Millisecond * 500) // wait for StopOnce to happen
			ch <- 2

			return nil
		})

	}()

	go func() {
		<-ready // try stopping halfway through startup
		once.StopOnce("slow service", func() (err error) {
			ch <- 3

			return nil
		})
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
		once.StartOnce("slow service", func() (err error) {
			ready <- true
			<-next // continue after the other StartOnce call fails

			return nil
		})
		<-next
		ch <- 2

	}()

	go func() {
		<-ready // try starting halfway through startup
		once.StartOnce("slow service", func() (err error) {
			return nil
		})
		next <- true
		ch <- 3
		next <- true

	}()

	require.Equal(t, 1, <-ch)
	require.Equal(t, 3, <-ch) // 3 arrives before 2 because it returns immediately
	require.Equal(t, 2, <-ch)
}
