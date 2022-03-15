package soltxm

import (
	"fmt"

	"github.com/gagliardetto/solana-go"

	solanaClient "github.com/smartcontractkit/chainlink-solana/pkg/solana/client"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/config"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/utils"
)

var (
	_ services.Service = (*Txm)(nil)
	_ solana.TxManager = (*Txm)(nil)
)

// Txm manages transactions for the solana blockchain.
// simple implementation with no persistently stored txs
type Txm struct {
	starter    utils.StartStopOnce
	lggr       logger.Logger
	tc         func() (solanaClient.ReaderWriter, error)
	stop, done chan struct{}
	cfg        solana.Config
}

// NewTxm creates a txm. Uses simulation so should only be used to send txes to trusted contracts i.e. OCR.
func NewTxm(tc func() (solanaClient.ReaderWriter, error), cfg config.Config, lggr logger.Logger) *Txm {
	lggr = lggr.Named("Txm")
	return &Txm{
		starter: utils.StartStopOnce{},
		tc:      tc,
		lggr:    lggr,
		queue:   make(chan *solana.Transaction, 1000), // queue can support 1000 pending txs
		stop:    make(chan struct{}),
		done:    make(chan struct{}),
		cfg:     cfg,
	}
}

// Start subscribes to queuing channel and processes them.
func (txm *Txm) Start() error {
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
		case tx <- txm.queue:
			// fetch client
			client, err := tc()
			if err != nil {
				txm.lggr.Errorw("failed to get client", "err", err)
			}
			// process tx
			sig, err := client.SendTx(tx)
			if err != nil {
				txm.lggr.Errorw("failed to send transaction", "err", err)
			}
			txm.lggr.Debugw("successfully sent transaction", "signature", sig)
		case <-txm.stop:
			return
		}
	}
}

// TODO: transaction confirmation
// use ConfirmPollPeriod() in config

// Enqueue enqueue a msg destined for the solana chain.
func (txm *Txm) Enqueue(accountID string, msg *solana.Transaction) error {
	select {
	case txm.queue <- msg:
	default:
		txm.lggr.Errorw("failed to enqeue tx", "queueLength", len(txm.queue), "tx", msg)
		return fmt.Errorf("failed to enqueue transaction for %s", accountID)
	}
	return nil
}

// Close close service
func (txm *Txm) Close() error {
	txm.sub.Close()
	close(txm.stop)
	<-txm.done
	return nil
}

// Healthy service is healthy
func (txm *Txm) Healthy() error {
	return nil
}

// Ready service is ready
func (txm *Txm) Ready() error {
	return nil
}
