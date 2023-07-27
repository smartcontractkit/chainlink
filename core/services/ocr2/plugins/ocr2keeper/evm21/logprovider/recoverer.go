package logprovider

import (
	"context"
	"errors"
	"sync"
	"time"

	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg"
	keepersflows "github.com/smartcontractkit/ocr2keepers/pkg/v3/flows"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

var (
	ErrNotFound = errors.New("not found")
)

// CoordinatedProposal contains all required values to construct a complete
// UpkeepPayload for use in a runner
// TODO:remove this struct and use the one from ocr2keeper
type CoordinatedProposal struct {
	UpkeepID ocr2keepers.UpkeepIdentifier
	Trigger  ocr2keepers.Trigger
	Block    ocr2keepers.BlockKey
}

type LogRecoverer interface {
	GetRecoverables() ([]ocr2keepers.UpkeepPayload, error)
}

type logRecoverer struct {
	lggr logger.Logger

	cancel context.CancelFunc

	interval time.Duration
	lock     *sync.RWMutex

	pending []ocr2keepers.UpkeepPayload
	visited map[string]bool

	poller logpoller.LogPoller
}

var _ keepersflows.RecoverableProvider = &logRecoverer{}

func NewLogRecoverer(lggr logger.Logger, poller logpoller.LogPoller, interval time.Duration) LogRecoverer {
	if interval == 0 {
		interval = 30 * time.Second
	}
	return &logRecoverer{
		lggr:     lggr.Named("LogRecoverer"),
		lock:     &sync.RWMutex{},
		interval: interval,
		poller:   poller,
	}
}

func (r *logRecoverer) Start(pctx context.Context) error {
	ctx, cancel := context.WithCancel(pctx)
	r.lock.Lock()
	r.cancel = cancel
	interval := r.interval
	r.lock.Unlock()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	r.lggr.Debug("Starting log recoverer")

	for {
		select {
		case <-ticker.C:
			r.recover(ctx)
		case <-ctx.Done():
			return nil
		}
	}
}

func (r *logRecoverer) Close() error {
	r.lock.Lock()
	defer r.lock.Unlock()

	if r.cancel != nil {
		r.cancel()
	}
	return nil
}

func (r *logRecoverer) GetRecoverables() ([]ocr2keepers.UpkeepPayload, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	if len(r.pending) == 0 {
		return nil, nil
	}

	pending := r.pending
	r.pending = make([]ocr2keepers.UpkeepPayload, 0)

	for _, p := range pending {
		r.visited[p.ID] = true
	}

	return pending, nil
}

func (r *logRecoverer) recover(ctx context.Context) {
	// TODO: implement
}
