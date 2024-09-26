package cache

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLazyFetchPass(t *testing.T) {
	counterFunction := createPassingCounter()

	counter, _ := counterFunction()
	require.Equal(t, 1, counter)

	lazyCounter := LazyFetch(counterFunction)
	counter, _ = lazyCounter()
	require.Equal(t, 2, counter)

	counter, _ = lazyCounter()
	require.Equal(t, 2, counter)
}

func TestLazyFetchFail(t *testing.T) {
	counterFunction := createFailingCounter()

	_, err := counterFunction()
	require.Equal(t, "counter 1 failed", err.Error())

	lazyCounter := LazyFetch(counterFunction)
	_, err = lazyCounter()
	require.Equal(t, "counter 2 failed", err.Error())

	_, err = lazyCounter()
	require.Equal(t, "counter 2 failed", err.Error())
}

func TestLazyFetchMultipleRoutines(t *testing.T) {
	routines := 100
	counterFunction := LazyFetch(createPassingCounter())

	var wg sync.WaitGroup
	wg.Add(routines)

	for i := 0; i < routines; i++ {
		go func() {
			counter, _ := counterFunction()
			require.Equal(t, 1, counter)
			wg.Done()
		}()
	}

	wg.Wait()
}

func createFailingCounter() func() (int, error) {
	counter := 0
	return func() (int, error) {
		counter++
		return 0, fmt.Errorf("counter %d failed", counter)
	}
}

func createPassingCounter() func() (int, error) {
	counter := 0
	return func() (int, error) {
		counter++
		return counter, nil
	}
}
