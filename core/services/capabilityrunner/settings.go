package capabilityrunner

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sync"

	"github.com/smartcontractkit/chainlink-common/pkg/types/core"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// NewFileBackedSettings wraps b so that every update is written to path before
// any subscriber sees it.
//
// Construct this once, alongside the SettingsBroadcaster it wraps, and share it:
// two instances over the same path would race each other over the file.
func NewFileBackedSettings(lggr logger.Logger, b core.SettingsBroadcaster, path string) *FileBackedSettings {
	return &FileBackedSettings{
		lggr:     lggr.Named("FileBackedSettings"),
		upstream: b,
		path:     path,
	}
}

// FileBackedSettings is a core.SettingsBroadcaster that dumps every update to
// disk before any subscriber observes it.
//
// Runner binaries have no RPC surface, so they pick settings up off the
// filesystem rather than over the LOOP. This holds a single subscription to the
// real broadcaster and fans updates out itself, so each update is written to
// path exactly once, by one writer, before any subscriber is handed it.
//
// An update whose write fails is not forwarded at all: forwarding it would make
// subscribers reload a file that does not match the update they were handed.
type FileBackedSettings struct {
	lggr     logger.Logger
	upstream core.SettingsBroadcaster
	path     string

	mu      sync.Mutex
	subs    []chan core.SettingsUpdate
	started bool
}

var _ core.SettingsBroadcaster = (*FileBackedSettings)(nil)
var _ core.SettingsBroadcasterLocal = (*FileBackedSettings)(nil)

// Subscribe returns a channel of updates that have already been written to disk.
//
// The returned channel holds only the newest update: settings are state, not
// events, so a subscriber busy reloading one update sees the latest value when
// it comes back rather than a backlog. New subscribers are not seeded with the
// current value - the reloader picks up whatever is already on disk when it
// starts, which covers the same ground without a redundant reload.
func (f *FileBackedSettings) Subscribe(ctx context.Context) (<-chan core.SettingsUpdate, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if !f.started {
		// Detached from ctx: the subscription belongs to the singleton, not to
		// whichever job happened to subscribe first.
		up, err := f.upstream.Subscribe(context.WithoutCancel(ctx))
		if err != nil {
			return nil, err
		}
		f.started = true
		go f.fanOut(up)
	}

	ch := make(chan core.SettingsUpdate, 1)
	f.subs = append(f.subs, ch)
	return ch, nil
}

func (f *FileBackedSettings) Unsubscribe(ch <-chan core.SettingsUpdate) {
	f.mu.Lock()
	defer f.mu.Unlock()
	i := slices.IndexFunc(f.subs, func(c chan core.SettingsUpdate) bool { return c == ch })
	if i < 0 {
		return
	}
	close(f.subs[i])
	f.subs = slices.Delete(f.subs, i, i+1)
}

// fanOut writes each upstream update to disk and then hands it to every
// subscriber. It runs for the lifetime of the process, ending only if the real
// broadcaster closes the subscription.
func (f *FileBackedSettings) fanOut(up <-chan core.SettingsUpdate) {
	for update := range up {
		if err := writeFileAtomic(f.path, update.Settings); err != nil {
			f.lggr.Errorw("failed to dump settings; not forwarding update", "err", err, "hash", update.Hash, "path", f.path)
			continue
		}
		f.lggr.Debugw("dumped settings", "hash", update.Hash, "path", f.path)

		// Sends happen under the lock so they cannot race Unsubscribe closing a
		// channel. They never block, so holding it is bounded.
		f.mu.Lock()
		for _, ch := range f.subs {
			send(ch, update)
		}
		f.mu.Unlock()
	}

	f.mu.Lock()
	defer f.mu.Unlock()
	for _, ch := range f.subs {
		close(ch)
	}
	f.subs, f.started = nil, false
}

// send delivers update to ch without blocking, displacing an older undelivered
// update if one is still sitting there.
func send(ch chan core.SettingsUpdate, update core.SettingsUpdate) {
	for {
		select {
		case ch <- update:
			return
		default:
		}
		select {
		case <-ch: // drop the stale value and retry
		default: // consumer took it in the meantime
		}
	}
}

// writeFileAtomic writes contents to path via a temp file in the same directory
// and a rename, so the runner never reads a torn file.
func writeFileAtomic(path, contents string) error {
	dir, name := filepath.Split(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create settings dir: %w", err)
	}
	tmp, err := os.CreateTemp(dir, name+".tmp-*")
	if err != nil {
		return fmt.Errorf("failed to create temp %s: %w", name, err)
	}
	if _, err = tmp.WriteString(contents); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmp.Name())
		return fmt.Errorf("failed to write %s: %w", name, err)
	}
	if err = tmp.Close(); err != nil {
		_ = os.Remove(tmp.Name())
		return fmt.Errorf("failed to close %s: %w", name, err)
	}
	if err = os.Rename(tmp.Name(), path); err != nil {
		_ = os.Remove(tmp.Name())
		return fmt.Errorf("failed to move %s into place: %w", name, err)
	}
	return nil
}
