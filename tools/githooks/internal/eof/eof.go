package eof

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"sync"

	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink/v2/tools/githooks/internal/filefilter"
)

// Config holds options for running end-of-file fixer.
type Config struct {
	CheckOnly bool
}

// Result contains information about processed files.
type Result struct {
	ModifiedFiles []string
}

// FixContent ensures non-empty content ends with exactly one newline (\n)
// while leaving empty content untouched.
func FixContent(content []byte) ([]byte, bool) {
	if len(content) == 0 {
		return content, false
	}

	trimmed := bytes.TrimRight(content, "\r\n")
	fixed := make([]byte, len(trimmed)+1)
	copy(fixed, trimmed)
	fixed[len(trimmed)] = '\n'

	if bytes.Equal(content, fixed) {
		return content, false
	}

	return fixed, true
}

// FixFile checks and fixes trailing newlines for the file at filePath.
func FixFile(filePath string, checkOnly bool) (bool, error) {
	cleanPath := filepath.Clean(filePath)
	if !filefilter.IsEligiblePath(cleanPath) {
		return false, nil
	}

	info, err := os.Stat(cleanPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to stat %s: %w", cleanPath, err)
	}

	if info.IsDir() || !info.Mode().IsRegular() {
		return false, nil
	}

	content, err := os.ReadFile(cleanPath)
	if err != nil {
		return false, fmt.Errorf("failed to read %s: %w", cleanPath, err)
	}

	if filefilter.IsBinary(content) {
		return false, nil
	}

	fixed, changed := FixContent(content)
	if !changed {
		return false, nil
	}

	if checkOnly {
		return true, nil
	}

	if err := os.WriteFile(cleanPath, fixed, info.Mode().Perm()); err != nil { // #nosec G703,G304 -- cleanPath is local file in repo
		return false, fmt.Errorf("failed to write %s: %w", cleanPath, err)
	}

	return true, nil
}

// Run executes the end-of-file fixer on all specified files concurrently using a worker pool.
func Run(ctx context.Context, repoRoot string, files []string, cfg Config) (*Result, error) {
	if len(files) == 0 {
		return &Result{}, nil
	}

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(max(1, min(len(files), runtime.GOMAXPROCS(0)*4)))

	var (
		mu       sync.Mutex
		modified []string
	)

	for _, file := range files {
		abs := file
		if !filepath.IsAbs(abs) {
			abs = filepath.Join(repoRoot, file)
		}
		g.Go(func() error {
			if ctx.Err() != nil {
				return ctx.Err()
			}

			changed, err := FixFile(abs, cfg.CheckOnly)
			if err != nil {
				return err
			}

			if changed {
				mu.Lock()
				modified = append(modified, file)
				mu.Unlock()
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("encountered errors during end-of-file fix: %w", err)
	}

	sort.Strings(modified)
	return &Result{ModifiedFiles: modified}, nil
}
