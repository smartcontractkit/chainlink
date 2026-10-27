package vault

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
)

var ErrDirectReadNotReady = errors.New("vault direct read path is not ready")

const directReadWaitPollInterval = 2 * time.Second

// TODO: match with errors.Is once libocr exports errKeyValueStateUnavailable.
const keyValueStateUnavailableSubstring = "key value state is unavailable"

func isKeyValueStateUnavailable(err error) bool {
	return err != nil && strings.Contains(err.Error(), keyValueStateUnavailableSubstring)
}

// DirectSecretsReader serves GetSecrets from local replicated state, returning only this node's shares.
type DirectSecretsReader interface {
	GetSecretsDirect(ctx context.Context, req *vaultcommon.GetSecretsRequest) (*vaultcommon.GetSecretsResponse, error)
}

// LazyDirectSecretsReader holds the reader of the running reporting plugin instance.
type LazyDirectSecretsReader struct {
	mu     sync.Mutex
	reader DirectSecretsReader
}

func NewLazyDirectSecretsReader() *LazyDirectSecretsReader {
	return &LazyDirectSecretsReader{}
}

func (l *LazyDirectSecretsReader) Get() DirectSecretsReader {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.reader
}

func (l *LazyDirectSecretsReader) Set(r DirectSecretsReader) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.reader = r
}

// Clear is a no-op if r has already been replaced, so a closing instance can't clear its successor.
func (l *LazyDirectSecretsReader) Clear(r DirectSecretsReader) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.reader == r {
		l.reader = nil
	}
}

func (s *Capability) handleGetSecretsDirect(ctx context.Context, requestID string, req *vaultcommon.GetSecretsRequest) (resp *vaultcommon.GetSecretsResponse, err error) {
	start := s.clock.Now()
	defer func() {
		d := s.clock.Since(start)
		s.lifecycle.RecordDirectGetSecrets(ctx, d, err)
		s.lggr.Debugw("served direct get secrets request", "requestID", requestID, "duration", d, "err", err)
	}()

	return s.getSecretsDirectWithRetry(ctx, req, start.Add(s.expiresAfter))
}

// getSecretsDirectWithRetry polls until a reader is registered and the
// replicated state is available, or the deadline passes.
func (s *Capability) getSecretsDirectWithRetry(ctx context.Context, req *vaultcommon.GetSecretsRequest, deadline time.Time) (*vaultcommon.GetSecretsResponse, error) {
	deadlineCh := s.clock.After(deadline.Sub(s.clock.Now()))
	var lastErr error
	for {
		if reader := s.directReader.Get(); reader != nil {
			resp, err := reader.GetSecretsDirect(ctx, req)
			switch {
			case err == nil:
				return resp, nil
			case !isKeyValueStateUnavailable(err):
				return nil, err
			default:
				lastErr = err
				s.lggr.Debugw("vault replicated state unavailable, retrying direct read", "err", err)
			}
		}
		select {
		case <-deadlineCh:
			if lastErr != nil {
				return nil, fmt.Errorf("%w: vault replicated state still unavailable, last error: %w", ErrDirectReadNotReady, lastErr)
			}
			return nil, ErrDirectReadNotReady
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-s.clock.After(directReadWaitPollInterval):
		}
	}
}
