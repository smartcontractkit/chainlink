package soltxm

import (
	"context"

	solanaGo "github.com/gagliardetto/solana-go"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-solana/pkg/solana"
	solanaClient "github.com/smartcontractkit/chainlink-solana/pkg/solana/client"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/config"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/utils"
)

const MaxQueueLen = 1000

var (
	_ services.ServiceCtx = (*Txm)(nil)
	_ solana.TxManager    = (*Txm)(nil)
)

// Txm manages transactions for the solana blockchain.
// simple implementation with no persistently stored txs
type Txm struct {
	starter    utils.StartStopOnce
	lggr       logger.Logger
	tc         func() (solanaClient.ReaderWriter, error)
	queue      chan *solanaGo.Transaction
	stop, done chan struct{}
	cfg        config.Config
}

// NewTxm creates a txm. Uses simulation so should only be used to send txes to trusted contracts i.e. OCR.
func NewTxm(tc func() (solanaClient.ReaderWriter, error), cfg config.Config, lggr logger.Logger) *Txm {
	lggr = lggr.Named("Txm")
	return &Txm{
		starter: utils.StartStopOnce{},
		tc:      tc,
		lggr:    lggr,
		queue:   make(chan *solanaGo.Transaction, MaxQueueLen), // queue can support 1000 pending txs
		stop:    make(chan struct{}),
		done:    make(chan struct{}),
		cfg:     cfg,
	}
}

// Start subscribes to queuing channel and processes them.
func (txm *Txm) Start(context.Context) error {
	return txm.starter.StartOnce("solanatxm", func() error {
		go txm.run()
		return nil
	})
}

func (txm *Txm) run() {
	defer close(txm.done)
	ctx, cancel := utils.ContextFromChan(txm.stop)
	defer cancel()
	for {
		select {
		case tx := <-txm.queue:
			// TODO: this section could be better optimized for sending TXs quickly
			// fetch client
			client, err := txm.tc()
			if err != nil {
				txm.lggr.Errorw("failed to get client", "err", err)
				continue
			}
			// process tx
			sig, err := client.SendTx(ctx, tx)
			if err != nil {
				txm.lggr.Criticalw("failed to send transaction", "err", err)
				continue
			}
			txm.lggr.Debugw("successfully sent transaction", "signature", sig.String())
		case <-txm.stop:
			return
		}
	}
}

// TODO: transaction confirmation
// use ConfirmPollPeriod() in config

// Enqueue enqueue a msg destined for the solana chain.
func (txm *Txm) Enqueue(accountID string, msg *solanaGo.Transaction) error {
	select {
	case txm.queue <- msg:
	default:
		txm.lggr.Errorw("failed to enqeue tx", "queueLength", len(txm.queue), "tx", msg)
		return errors.Errorf("failed to enqueue transaction for %s", accountID)
	}
	return nil
}

// Close close service
func (txm *Txm) Close() error {
	return txm.starter.StopOnce("solanatxm", func() error {
		close(txm.stop)
		<-txm.done
		return nil
	})
}

// Healthy service is healthy
func (txm *Txm) Healthy() error {
	return nil
}

// Ready service is ready
func (txm *Txm) Ready() error {
	return nil
}
