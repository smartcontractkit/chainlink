package vault

import (
	"context"
	"errors"
	"sync"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
)

var ErrDirectReadNotReady = errors.New("vault direct read path is not ready")

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
		return nil, ErrDirectReadNotReady
	}
	return reader.GetSecretsDirect(ctx, req)
}
