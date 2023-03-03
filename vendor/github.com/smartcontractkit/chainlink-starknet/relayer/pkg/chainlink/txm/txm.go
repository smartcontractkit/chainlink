package txm

import (
	"context"
	"math/big"
	"sync"
	"time"

	"github.com/dontpanicdao/caigo"
	caigotypes "github.com/dontpanicdao/caigo/types"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-relay/pkg/logger"
	relaytypes "github.com/smartcontractkit/chainlink-relay/pkg/types"
	"github.com/smartcontractkit/chainlink-relay/pkg/utils"

	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/starknet"

	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/keys"
)

const (
	MaxQueueLen = 1000
)

type TxManager interface {
	Enqueue(caigotypes.Hash, caigotypes.FunctionCall) error
}

type Tx struct {
	sender caigotypes.Hash
	call   caigotypes.FunctionCall
}

type StarkTXM interface {
	relaytypes.Service
	TxManager
}

type starktxm struct {
	starter utils.StartStopOnce
	lggr    logger.Logger
	done    sync.WaitGroup
	stop    chan struct{}
	queue   chan Tx
	ks      keys.Keystore
	cfg     Config

	// TODO: use lazy loaded client
	client    *starknet.Client
	getClient func() (*starknet.Client, error)
}

func New(lggr logger.Logger, keystore keys.Keystore, cfg Config, getClient func() (*starknet.Client, error)) (StarkTXM, error) {
	return &starktxm{
		lggr:      logger.Named(lggr, "StarkNetTxm"),
		queue:     make(chan Tx, MaxQueueLen),
		stop:      make(chan struct{}),
		getClient: getClient,
		ks:        keystore,
		cfg:       cfg,
	}, nil
}

func (txm *starktxm) Name() string {
	return txm.lggr.Name()
}

func (txm *starktxm) Start(ctx context.Context) error {
	return txm.starter.StartOnce("starktxm", func() error {
		txm.done.Add(1) // waitgroup: tx sender
		go txm.run()
		return nil
	})
}

func (txm *starktxm) run() {
	defer txm.done.Done()

	// TODO: func not available without importing core
	// ctx, cancel := utils.ContextFromChan(txm.stop)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tick := time.After(0)

	for {
		select {
		case <-tick:
			start := time.Now()

			// fetch client if needed (before processing txs to preserve queue)
			if txm.client == nil {
				var err error
				txm.client, err = txm.getClient()
				if err != nil {
					txm.lggr.Errorw("unable to fetch client", "error", err)
					tick = time.After(utils.WithJitter(txm.cfg.TxSendFrequency()) - time.Since(start)) // reset tick
					txm.client = nil                                                                   // reset
					continue
				}
			}

			// calculate total txs to process
			txLen := len(txm.queue)
			if txLen > txm.cfg.TxMaxBatchSize() {
				txLen = txm.cfg.TxMaxBatchSize()
			}

			// fetch batch and split by sender
			txsBySender := map[caigotypes.Hash][]caigotypes.FunctionCall{}
			for i := 0; i < txLen; i++ {
				tx := <-txm.queue
				txsBySender[tx.sender] = append(txsBySender[tx.sender], tx.call)
			}
			txm.lggr.Infow("creating batch", "totalTxCount", txLen, "batchCount", len(txsBySender))
			txm.lggr.Debugw("batch details", "batches", txsBySender)

			// async process of tx batches
			var wg sync.WaitGroup
			wg.Add(len(txsBySender))
			for sender, txs := range txsBySender {
				go func(sender caigotypes.Hash, txs []caigotypes.FunctionCall) {

					// fetch key matching sender address
					key, err := txm.ks.Get(sender.String())
					if err != nil {
						txm.lggr.Errorw("failed to retrieve key", "id", sender, "error", err)
					} else {
						// parse key to expected format
						privKeyBytes := key.Raw()
						privKey := caigotypes.BigToHex(caigotypes.BytesToBig(privKeyBytes))

						// broadcast batch based on sender
						hash, err := txm.broadcastBatch(ctx, privKey, sender, txs)
						if err != nil {
							txm.lggr.Errorw("transaction failed to broadcast", "error", err, "batchTx", txs)
						} else {
							txm.lggr.Infow("transaction broadcast", "txhash", hash)
						}
					}
					wg.Done()
				}(sender, txs)
			}
			wg.Wait()
			tick = time.After(utils.WithJitter(txm.cfg.TxSendFrequency()) - time.Since(start))
		case <-txm.stop:
			return
		}
	}
}

const FEE_MARGIN uint64 = 115

func (txm *starktxm) broadcastBatch(ctx context.Context, privKey string, sender caigotypes.Hash, txs []caigotypes.FunctionCall) (txhash string, err error) {
	// create new account
	account, err := caigo.NewGatewayAccount(privKey, sender.String(), txm.client.Gw, caigo.AccountVersion1)
	if err != nil {
		return txhash, errors.Errorf("failed to create new account: %s", err)
	}

	// get fee for txm
	feeEstimate, err := account.EstimateFee(ctx, txs, caigotypes.ExecuteDetails{})
	if err != nil {
		return txhash, errors.Errorf("failed to estimate fee: %s", err)
	}

	fee, _ := big.NewInt(0).SetString(string(feeEstimate.OverallFee), 0)
	expandedFee := big.NewInt(0).Mul(fee, big.NewInt(int64(FEE_MARGIN)))
	max := big.NewInt(0).Div(expandedFee, big.NewInt(100))
	details := caigotypes.ExecuteDetails{
		MaxFee: max,
	}

	// TODO: investigate if nonce management is needed (nonce is requested queried by the sdk for now)
	// transmit txs
	execCtx, execCancel := context.WithTimeout(ctx, txm.cfg.TxTimeout()*time.Second)
	defer execCancel()
	res, err := account.Execute(execCtx, txs, details)
	if err != nil {
		return txhash, errors.Errorf("failed to invoke tx: %s", err)
	}

	// handle nil pointer
	if res == nil {
		return txhash, errors.New("execute response and error are nil")
	}

	return res.TransactionHash, nil
}

func (txm *starktxm) Close() error {
	return txm.starter.StopOnce("starktxm", func() error {
		close(txm.stop)
		txm.done.Wait()
		return nil
	})
}

func (txm *starktxm) Healthy() error {
	return txm.starter.Healthy()
}

func (txm *starktxm) Ready() error {
	return txm.starter.Ready()
}

func (txm *starktxm) HealthReport() map[string]error {
	return map[string]error{txm.Name(): txm.Healthy()}
}

func (txm *starktxm) Enqueue(sender caigotypes.Hash, tx caigotypes.FunctionCall) error {
	select {
	case txm.queue <- Tx{sender: sender, call: tx}:
	default:
		return errors.Errorf("failed to enqueue transaction: %+v", tx)
	}

	return nil
}
