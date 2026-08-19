package tidy_test

import (
	"context"
	"errors"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/tools/githooks/internal/tidy"
)

func TestRun(t *testing.T) {
	t.Parallel()

	t.Run("runs tidy in parallel across multiple modules", func(t *testing.T) {
		t.Parallel()

		var (
			mu         sync.Mutex
			calledDirs []string
			concurrent int
			maxConc    int
		)

		runner := func(ctx context.Context, dir string, args ...string) error {
			mu.Lock()
			calledDirs = append(calledDirs, dir)
			concurrent++
			if concurrent > maxConc {
				maxConc = concurrent
			}
			mu.Unlock()

			// small sleep to test concurrency
			time.Sleep(20 * time.Millisecond)

			mu.Lock()
			concurrent--
			mu.Unlock()

			assert.Equal(t, []string{"mod", "tidy"}, args)
			return nil
		}

		repoRoot := "/test/repo"
		modules := []string{".", "deployment", "tools/githooks"}

		err := tidy.Run(t.Context(), repoRoot, modules, tidy.Config{Runner: runner})
		require.NoError(t, err)

		mu.Lock()
		defer mu.Unlock()

		assert.Len(t, calledDirs, 3)
		assert.Contains(t, calledDirs, filepath.Join(repoRoot, "."))
		assert.Contains(t, calledDirs, filepath.Join(repoRoot, "deployment"))
		assert.Contains(t, calledDirs, filepath.Join(repoRoot, "tools/githooks"))
		assert.Greater(t, maxConc, 1, "expected parallel execution")
	})

	t.Run("returns error when a module tidy fails", func(t *testing.T) {
		t.Parallel()

		runner := func(ctx context.Context, dir string, args ...string) error {
			if dir == filepath.Join("/test/repo", "failing-mod") {
				return errors.New("tidy failed on failing-mod")
			}
			return nil
		}

		repoRoot := "/test/repo"
		modules := []string{"good-mod", "failing-mod"}

		err := tidy.Run(t.Context(), repoRoot, modules, tidy.Config{Runner: runner})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failing-mod")
	})

	t.Run("no modules to tidy", func(t *testing.T) {
		t.Parallel()

		called := false
		runner := func(ctx context.Context, dir string, args ...string) error {
			called = true
			return nil
		}

		err := tidy.Run(t.Context(), "/test/repo", nil, tidy.Config{Runner: runner})
		require.NoError(t, err)
		assert.False(t, called)
	})
}
