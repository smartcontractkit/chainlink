package monitor

import (
	"context"
	"time"

	"github.com/gagliardetto/solana-go"

	solanaClient "github.com/smartcontractkit/chainlink-solana/pkg/solana/client"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/solkey"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// Config defines the monitor configuration.
type Config interface {
	BalancePollPeriod() time.Duration
}

// Keystore provides the keys to be monitored.
type Keystore interface {
	GetAll() ([]solkey.Key, error)
}

// NewBalanceMonitor returns a balance monitoring services.Service which reports the luna balance of all ks keys to prometheus.
func NewBalanceMonitor(chainID string, cfg Config, lggr logger.Logger, ks Keystore, newReader func() (solanaClient.Reader, error)) services.ServiceCtx {
	return newBalanceMonitor(chainID, cfg, lggr, ks, newReader)
}

func newBalanceMonitor(chainID string, cfg Config, lggr logger.Logger, ks Keystore, newReader func() (solanaClient.Reader, error)) *balanceMonitor {
	b := balanceMonitor{
		chainID:   chainID,
		cfg:       cfg,
		lggr:      lggr.Named("BalanceMonitor"),
		ks:        ks,
		newReader: newReader,
		stop:      make(chan struct{}),
		done:      make(chan struct{}),
	}
	b.updateFn = b.updateProm
	return &b
}

type balanceMonitor struct {
	utils.StartStopOnce
	chainID   string
	cfg       Config
	lggr      logger.Logger
	ks        Keystore
	newReader func() (solanaClient.Reader, error)
	updateFn  func(acc solana.PublicKey, lamports uint64) // overridable for testing

	reader solanaClient.Reader

	stop, done chan struct{}
}

func (b *balanceMonitor) Name() string {
	return b.lggr.Name()
}

func (b *balanceMonitor) Start(context.Context) error {
	return b.StartOnce("SolanaBalanceMonitor", func() error {
		go b.monitor()
		return nil
	})
}

func (b *balanceMonitor) Close() error {
	return b.StopOnce("SolanaBalanceMonitor", func() error {
		close(b.stop)
		<-b.done
		return nil
	})
}

func (b *balanceMonitor) HealthReport() map[string]error {
	return map[string]error{b.Name(): b.Healthy()}
}

func (b *balanceMonitor) monitor() {
	defer close(b.done)

	tick := time.After(utils.WithJitter(b.cfg.BalancePollPeriod()))
	for {
		select {
		case <-b.stop:
			return
		case <-tick:
			b.updateBalances()
			tick = time.After(utils.WithJitter(b.cfg.BalancePollPeriod()))
		}
	}
}

// getReader returns the cached solanaClient.Reader, or creates a new one if nil.
func (b *balanceMonitor) getReader() (solanaClient.Reader, error) {
	if b.reader == nil {
		var err error
		b.reader, err = b.newReader()
		if err != nil {
			return nil, err
		}
	}
	return b.reader, nil
}

func (b *balanceMonitor) updateBalances() {
	keys, err := b.ks.GetAll()
	if err != nil {
		b.lggr.Errorw("Failed to get keys", "err", err)
		return
	}
	if len(keys) == 0 {
		return
	}
	reader, err := b.getReader()
	if err != nil {
		b.lggr.Errorw("Failed to get client", "err", err)
		return
	}
	var gotSomeBals bool
	for _, k := range keys {
		// Check for shutdown signal, since Balance blocks and may be slow.
		select {
		case <-b.stop:
			return
		default:
		}
		acc := k.PublicKey()
		lamports, err := reader.Balance(acc)
		if err != nil {
			b.lggr.Errorw("Failed to get balance", "account", acc.String(), "err", err)
			continue
		}
		gotSomeBals = true
		b.updateFn(acc, lamports)
	}
	if !gotSomeBals {
		// Try a new client next time.
		b.reader = nil
	}
}
