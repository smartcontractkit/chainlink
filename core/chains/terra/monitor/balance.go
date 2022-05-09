package monitor

import (
	"context"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/smartcontractkit/chainlink-terra/pkg/terra/client"

	"github.com/smartcontractkit/chainlink/core/chains/terra/denom"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/terrakey"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// Config defines the monitor configuration.
type Config interface {
	BlockRate() time.Duration
}

// Keystore provides the keys to be monitored.
type Keystore interface {
	GetAll() ([]terrakey.Key, error)
}

// NewBalanceMonitor returns a balance monitoring services.ServiceCtx which reports the luna balance of all ks keys to prometheus.
func NewBalanceMonitor(chainID string, cfg Config, lggr logger.Logger, ks Keystore, newReader func(string) (client.Reader, error)) services.ServiceCtx {
	return newBalanceMonitor(chainID, cfg, lggr, ks, newReader)
}

func newBalanceMonitor(chainID string, cfg Config, lggr logger.Logger, ks Keystore, newReader func(string) (client.Reader, error)) *balanceMonitor {
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
	newReader func(string) (client.Reader, error)
	updateFn  func(acc sdk.AccAddress, bal *sdk.DecCoin) // overridable for testing

	reader client.Reader

	stop, done chan struct{}
}

// Start starts balance monitor for terra.
func (b *balanceMonitor) Start(context.Context) error {
	return b.StartOnce("TerraBalanceMonitor", func() error {
		go b.monitor()
		return nil
	})
}

func (b *balanceMonitor) Close() error {
	return b.StopOnce("TerraBalanceMonitor", func() error {
		close(b.stop)
		<-b.done
		return nil
	})
}

func (b *balanceMonitor) monitor() {
	defer close(b.done)

	tick := time.After(utils.WithJitter(b.cfg.BlockRate()))
	for {
		select {
		case <-b.stop:
			return
		case <-tick:
			b.updateBalances()
			tick = time.After(utils.WithJitter(b.cfg.BlockRate()))
		}
	}
}

// getReader returns the cached client.Reader, or creates a new one if nil.
func (b *balanceMonitor) getReader() (client.Reader, error) {
	if b.reader == nil {
		var err error
		b.reader, err = b.newReader("")
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
		acc := sdk.AccAddress(k.PublicKey().Address())
		bal, err := reader.Balance(acc, "uluna")
		if err != nil {
			b.lggr.Errorw("Failed to get balance", "account", acc, "err", err)
			continue
		}
		gotSomeBals = true
		balLuna, err := denom.ConvertToLuna(*bal)
		if err != nil {
			b.lggr.Errorw("Failed to convert uluna to luna", "account", acc, "err", err)
			continue
		}
		b.updateFn(acc, &balLuna)
	}
	if !gotSomeBals {
		// Try a new client next time.
		b.reader = nil
	}
}
