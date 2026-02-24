package concurrency_test

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink/devenv/products/automation/concurrency"
)

type result struct {
	integer int
}

func (r result) GetResult() int {
	return r.integer
}

type config struct{}

func TestExecute(t *testing.T) {
	type tc struct {
		name            string
		concurrency     int
		tasks           int
		expectedResults int
		expectedErrors  int
		errorFn         func(*big.Int) error
	}

	tcs := []tc{
		{
			name:            "concurrency 1, tasks 1",
			concurrency:     1,
			tasks:           1,
			expectedResults: 1,
			expectedErrors:  0,
		},
		{
			name:            "concurrency 10, tasks 1",
			concurrency:     10,
			tasks:           1,
			expectedResults: 1,
			expectedErrors:  0,
		},
		{
			name:            "concurrency 2, tasks 5",
			concurrency:     10,
			tasks:           1,
			expectedResults: 1,
			expectedErrors:  0,
		},
		{
			name:            "concurrency 0, tasks 1",
			concurrency:     0,
			tasks:           1,
			expectedResults: 1,
			expectedErrors:  0,
		},
		{
			name:            "concurrency -1, tasks 1",
			concurrency:     -1,
			tasks:           1,
			expectedResults: 1,
			expectedErrors:  0,
		},
		{
			name:            "concurrency 1, tasks 0",
			concurrency:     1,
			tasks:           0,
			expectedResults: 0,
			expectedErrors:  0,
		},
		{
			name:            "concurrency 1, tasks 100, sporadic errors",
			concurrency:     1,
			tasks:           100,
			expectedResults: 0,
			expectedErrors:  1,
			errorFn: func(integer *big.Int) error {
				if integer.Int64()%2 == 0 {
					return errors.New("even number error")
				}
				return nil
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			tc := tc
			t.Parallel()

			executor := concurrency.NewConcurrentExecutor[int, result, config](framework.L)

			processorFn := func(resultCh chan result, errCh chan error, keyNum int, payload config) {
				randomInt, err := rand.Int(rand.Reader, big.NewInt(100))
				if err != nil {
					errCh <- err
					return
				}

				if tc.errorFn != nil {
					err := tc.errorFn(randomInt)
					if err != nil {
						errCh <- err
						return
					}
				}

				resultCh <- result{integer: int(randomInt.Int64())}
			}

			configs := []config{}
			for i := 0; i < tc.tasks; i++ {
				configs = append(configs, struct{}{})
			}

			results, err := executor.Execute(tc.expectedResults, configs, processorFn)
			if tc.expectedErrors == 0 {
				require.NoError(t, err, "Error executing concurrently")
				require.Len(t, results, tc.expectedResults, "Wrong result number")
			} else {
				require.Error(t, err, "No error returned when executing concurrently")
				require.Empty(t, results, "Expected no results")
				require.GreaterOrEqual(t, len(executor.GetErrors()), tc.expectedErrors, "Wrong error number")
			}
		})
	}
}

func TestExecuteFailFast(t *testing.T) {
	type tc struct {
		name        string
		executor    *concurrency.ConcurrentExecutor[int, result, config]
		processorFn func(resultCh chan result, errCh chan error, keyNum int, payload config)
		failFast    bool
	}

	tcs := []tc{
		{
			name:     "fail fast enabled",
			executor: concurrency.NewConcurrentExecutor[int, result, config](framework.L),
			failFast: true,
			processorFn: func(resultCh chan result, errCh chan error, keyNum int, payload config) {
				time.Sleep(10 * time.Millisecond)
				errCh <- errors.New("always fail, fail fast enabled")
			},
		},
		{
			name:     "fail fast disabled",
			executor: concurrency.NewConcurrentExecutor(framework.L, concurrency.WithoutFailFast[int, result, config]()),
			failFast: false,
			processorFn: func(resultCh chan result, errCh chan error, keyNum int, payload config) {
				errCh <- errors.New("always fail, fail fast disabled")
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			tc := tc
			t.Parallel()

			expectedExecutions := 1000

			configs := []config{}
			for i := 0; i < expectedExecutions; i++ {
				configs = append(configs, struct{}{})
			}

			results, err := tc.executor.Execute(100, configs, tc.processorFn)
			require.Error(t, err, "No error returned when executing concurrently")
			require.Empty(t, results, "Expected no results")
			if tc.failFast {
				fmt.Println(len(tc.executor.GetErrors()))
				require.Less(t, len(tc.executor.GetErrors()), expectedExecutions, "With fail fast enabled not all tasks should be executed")
			} else {
				require.Len(t, tc.executor.GetErrors(), expectedExecutions, "With fail fast disabled all tasks should be executed")
			}
		})
	}
}

func TestExecuteSimple(t *testing.T) {
	type tc struct {
		name            string
		concurrency     int
		tasks           int
		expectedResults int
		expectedErrors  int
		errorFn         func(*big.Int) error
	}

	tcs := []tc{
		{
			name:            "concurrency 2, tasks 5",
			concurrency:     10,
			tasks:           1,
			expectedResults: 1,
			expectedErrors:  0,
		},
		{
			name:            "concurrency 1, tasks 100, sporadic errors",
			concurrency:     1,
			tasks:           100,
			expectedResults: 0,
			expectedErrors:  1,
			errorFn: func(integer *big.Int) error {
				if integer.Int64()%2 == 0 {
					return errors.New("even number error")
				}
				return nil
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			tc := tc
			t.Parallel()

			executor := concurrency.NewConcurrentExecutor[int, result, concurrency.NoTaskType](framework.L)

			processorFn := func(resultCh chan result, errCh chan error, keyNum int) {
				randomInt, err := rand.Int(rand.Reader, big.NewInt(100))
				if err != nil {
					errCh <- err
					return
				}

				if tc.errorFn != nil {
					err := tc.errorFn(randomInt)
					if err != nil {
						errCh <- err
						return
					}
				}

				resultCh <- result{integer: int(randomInt.Int64())}
			}

			results, err := executor.ExecuteSimple(tc.concurrency, tc.tasks, processorFn)
			if tc.expectedErrors == 0 {
				require.NoError(t, err, "Error executing concurrently")
				require.Len(t, results, tc.expectedResults, "Wrong result number")
			} else {
				require.Error(t, err, "No error returned when executing concurrently")
				require.Empty(t, results, "Expected no results")
				require.GreaterOrEqual(t, len(executor.GetErrors()), tc.expectedErrors, "Wrong error number")
			}
		})
	}
}

func TestParentContext(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	executor := concurrency.NewConcurrentExecutor(framework.L, concurrency.WithContext[int, result, config](ctx))

	processorFn := func(resultCh chan result, errCh chan error, keyNum int, _ config) {
		time.Sleep(10 * time.Millisecond)
		randomInt, err := rand.Int(rand.Reader, big.NewInt(100))
		if err != nil {
			errCh <- err
			return
		}

		resultCh <- result{integer: int(randomInt.Int64())}
	}

	taskCount := 1000

	configs := []config{}
	for i := 0; i < taskCount; i++ {
		configs = append(configs, struct{}{})
	}

	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	results, err := executor.Execute(10, configs, processorFn)

	require.NoError(t, err, "Error executing concurrently")
	require.GreaterOrEqual(t, len(results), 1, "Wrong result number")
	require.Empty(t, executor.GetErrors(), "No errors expected")
}
