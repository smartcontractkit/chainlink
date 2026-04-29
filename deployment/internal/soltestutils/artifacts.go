package soltestutils

import (
	"context"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"

	"golang.org/x/sync/singleflight"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment/utils/solutils"
)

// Cross-process coordination (multiple `go test` invocations sharing the same
// programs_cache directory) is out of scope. singleflight deduplicates concurrent
// callers within a single test binary; separate processes should use isolated
// cache directories or an external file lock.
var (
	dlGroup singleflight.Group
	dlMu    sync.RWMutex
	dlDone  = map[string]bool{}
)

const (
	ccipKey   = "ccip"
	solanaKey = "solana"
)

// downloadFunc is a function type for downloading program artifacts.
type downloadFunc func(t *testing.T) string

// downloadChainlinkSolanaProgramArtifacts downloads the Chainlink Solana program artifacts.
//
// The artifacts that are downloaded contain both the CCIP and MCMS program artifacts (even though
// this is called "CCIP" program artifacts).
func downloadChainlinkSolanaProgramArtifacts(t *testing.T) string {
	t.Helper()
	cachePath := programsCachePath()
	doDownload(t, solanaKey, func(ctx context.Context) error {
		return solutils.DownloadChainlinkSolanaProgramArtifacts(ctx, cachePath, "", nil)
	})
	return cachePath
}

// downloadChainlinkCCIPProgramArtifacts downloads the Chainlink CCIP program artifacts for the
// test environment.
//
// The artifacts that are downloaded contain both the CCIP and MCMS program artifacts (even though
// this is called "CCIP" program artifacts).
func downloadChainlinkCCIPProgramArtifacts(t *testing.T) string {
	t.Helper()
	cachePath := programsCachePath()
	doDownload(t, ccipKey, func(ctx context.Context) error {
		return solutils.DownloadChainlinkCCIPProgramArtifacts(ctx, cachePath, "", nil)
	})
	return cachePath
}

// doDownload runs fn at most once per key per process via singleflight.
//
// The singleflight closure deliberately does NOT capture t or t.Context(): the
// winning goroutine may outlive the test that triggered it, so it uses a bounded
// background context. Each caller inspects the returned error against its own live
// t and calls t.Fatalf in its own goroutine — never inside the closure.
func doDownload(t *testing.T, key string, fn func(ctx context.Context) error) {
	t.Helper()

	dlMu.RLock()
	already := dlDone[key]
	dlMu.RUnlock()
	if already {
		return
	}

	_, err, _ := dlGroup.Do(key, func() (interface{}, error) {
		// Re-check under the write lock in case another caller completed while we waited.
		dlMu.RLock()
		if dlDone[key] {
			dlMu.RUnlock()
			return nil, nil
		}
		dlMu.RUnlock()

		// Use a bounded background context — never tied to *testing.T so the download
		// is not cancelled if the winning goroutine's test finishes first.
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()

		if err := fn(ctx); err != nil {
			return nil, err
		}

		// Mark success inside the closure, under the write lock, before returning nil
		// so all waiting callers observe the completed state atomically.
		dlMu.Lock()
		dlDone[key] = true
		dlMu.Unlock()
		return nil, nil
	})

	// Each caller handles the error against its own live *testing.T.
	require.NoError(t, err)
}

// programsCachePath returns the path to the cache directory for the program artifacts.
//
// This is used to cache the program artifacts so that they do not need to be downloaded every time
// the tests are run.
//
// The cache directory is located in the same directory as the current file.
func programsCachePath() string {
	// Get the directory of the current file
	_, currentFile, _, _ := runtime.Caller(0)

	dir := filepath.Dir(currentFile)

	return filepath.Join(dir, "programs_cache")
}
