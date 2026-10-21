package vault

import (
	"context"
	"errors"
	"sync"
	"time"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
)

var ErrDirectReadNotReady = errors.New("vault direct read path is not ready")

const directReadWaitPollInterval = 2 * time.Second

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

	reader := s.directReader.Get()
	if reader == nil {
		reader, err = s.waitForDirectReader(ctx, start.Add(s.expiresAfter))
		if err != nil {
			return nil, err
		}
	}
	return reader.GetSecretsDirect(ctx, req)
}

// waitForDirectReader polls the lazy reader until the running reporting plugin
// instance's directReader is available or the deadline passes.
func (s *Capability) waitForDirectReader(ctx context.Context, deadline time.Time) (DirectSecretsReader, error) {
	deadlineCh := s.clock.After(deadline.Sub(s.clock.Now()))
	for {
		select {
		case <-deadlineCh:
			return nil, ErrDirectReadNotReady
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-s.clock.After(directReadWaitPollInterval):
			if reader := s.directReader.Get(); reader != nil {
				return reader, nil
			}
		}
	}
}
