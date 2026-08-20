package whitespace

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/smartcontractkit/chainlink/v2/tools/githooks/internal/filefilter"
)

// Config holds options for running whitespace fixer.
type Config struct {
	CheckOnly bool
	Stdout    io.Writer
	Stderr    io.Writer
}

// Result contains information about processed files.
type Result struct {
	ModifiedFiles []string
}

// FixContent routes content to the appropriate language-aware whitespace fixer based on file extension.
func FixContent(filePath string, content []byte) ([]byte, bool, error) {
	if len(content) == 0 {
		return content, false, nil
	}

	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".go":
		return FixGo(content)
	case ".py":
		return FixPython(content)
	case ".md", ".markdown", ".mdown", ".mkdn":
		return FixMarkdown(content)
	default:
		return FixGeneric(content)
	}
}

// FixFile checks and trims extraneous whitespace for the file at filePath.
func FixFile(filePath string, checkOnly bool) (bool, error) {
	cleanPath := filepath.Clean(filePath)
	eligible, err := filefilter.IsEligibleFile(cleanPath)
	if err != nil {
		return false, fmt.Errorf("failed to check eligibility for %s: %w", cleanPath, err)
	}
	if !eligible {
		return false, nil
	}

	info, err := os.Stat(cleanPath)
	if err != nil {
		return false, fmt.Errorf("failed to stat %s: %w", cleanPath, err)
	}

	content, err := os.ReadFile(cleanPath)
	if err != nil {
		return false, fmt.Errorf("failed to read %s: %w", cleanPath, err)
	}

	fixed, changed, err := FixContent(cleanPath, content)
	if err != nil {
		return false, fmt.Errorf("failed to fix whitespace in %s: %w", cleanPath, err)
	}
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

// Run executes the whitespace fixer on all specified files.
func Run(ctx context.Context, repoRoot string, files []string, cfg Config) (*Result, error) {
	if len(files) == 0 {
		return &Result{}, nil
	}

	var (
		wg       sync.WaitGroup
		mu       sync.Mutex
		modified []string
		errs     []error
	)

	for _, file := range files {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		filePath := file
		if !filepath.IsAbs(filePath) {
			filePath = filepath.Join(repoRoot, file)
		}

		wg.Add(1)
		go func(p string, rel string) {
			defer wg.Done()

			changed, err := FixFile(p, cfg.CheckOnly)
			if err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
				return
			}

			if changed {
				mu.Lock()
				modified = append(modified, rel)
				mu.Unlock()
			}
		}(filePath, file)
	}

	wg.Wait()

	if len(errs) > 0 {
		return nil, fmt.Errorf("encountered errors during whitespace fix: %w", errors.Join(errs...))
	}

	sort.Strings(modified)
	return &Result{ModifiedFiles: modified}, nil
}
