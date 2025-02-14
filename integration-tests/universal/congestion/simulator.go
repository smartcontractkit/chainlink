package congestion

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"

	"github.com/smartcontractkit/chainlink-common/pkg/services"

	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

const backoffPeriodOnError = time.Second
const defaultTxPerBlock = 10
const congestionMinTarget = 0.9
const maxBlocksToCalcCongestion = 3

// Input - defines configuration for chain congestion - state when block is filled with transactions
// with high gas price.
type Input struct {
	InitialDelay            time.Duration `toml:"initial_delay"` // defines delay before the first surge
	Period                  time.Duration `toml:"period"`        // defines delays between surges
	Duration                time.Duration `toml:"duration"`      // defines duration of a surge
	GasPriceIncreasePercent int64         `toml:"gas_price_increase_percent"`
	NumberOfTxsPerBlock     *int          `toml:"number_of_txs_per_block"` // defines number of transaction to be produced to fill a block (default: defaultTxPerBlock)
	PayloadSize             *int          `toml:"payload_size"`            // defines payload size of transactions to fill a block. Calculated automatically, if not specified
	GasLimit                *uint64       `toml:"gas_limit"`               // calculated automatically if not specified
	MinTxsInPool            *int          `toml:"min_txs_in_pool"`         // defines minimal number of transactions Simulator will try to maintain in the mem pool
}

func (i *Input) defaults() {
	if i.NumberOfTxsPerBlock == nil {
		i.NumberOfTxsPerBlock = new(int)
		*i.NumberOfTxsPerBlock = defaultTxPerBlock
	}

	if i.MinTxsInPool == nil {
		i.MinTxsInPool = new(int)
		*i.MinTxsInPool = *i.NumberOfTxsPerBlock * 2
	}
}

type Simulator struct {
	services.Service
	eng *services.Engine

	t *testing.T
	Input
	client *seth.Client
	lggr   logger.SugaredLogger

	emitter contracts.LogEmitter
	payload []string

	congestionStartedAt *uint64
}

func NewSimulator(t *testing.T, input Input, client *seth.Client) (*Simulator, error) {
	input.defaults()
	// required to calculate congestion
	if client.HeaderCache == nil {
		client.HeaderCache = seth.NewLFUBlockCache(maxBlocksToCalcCongestion)
	}

	if *input.MinTxsInPool > len(client.Addresses) {
		return nil, fmt.Errorf("min_txs_in_pool cannot be greater than number of addresses available for seth.Client %d. Increase number of addresses", len(client.Addresses))
	}

	s := &Simulator{
		t:      t,
		Input:  input,
		client: client,
	}

	s.Service, s.eng = services.Config{
		Name:  "CongestionSimulator",
		Start: s.start,
	}.NewServiceEngine(logger.Test(t))
	s.lggr = s.eng.SugaredLogger
	t.Cleanup(func() {
		_ = s.Close()
	})
	return s, nil
}

func (s *Simulator) start(ctx context.Context) error {
	// deploy logs emitter that we'll use to fill blocks
	var err error
	s.emitter, err = contracts.DeployLogEmitterContract(logging.GetTestLogger(s.t), s.client)
	if err != nil {
		return fmt.Errorf("failed to create emitter contract: %v", err)
	}

	if s.Input.PayloadSize == nil {
		payloadSize, err := s.estimatePayloadSizeToFillBlock(ctx)
		if err != nil {
			return fmt.Errorf("failed to estimate payload size to fill block: %v", err)
		}
		s.Input.PayloadSize = &payloadSize
	}

	if s.Input.GasLimit == nil {
		gasLimit, err := s.estimateGasLimit(ctx)
		if err != nil {
			return fmt.Errorf("failed to estimate gas limit: %v", err)
		}
		s.Input.GasLimit = &gasLimit
	}

	s.payload = []string{string(make([]byte, *s.Input.PayloadSize))}
	ticker := services.TickerConfig{Initial: s.InitialDelay}.NewTicker(s.Period)
	s.eng.GoTick(ticker, func(ctx context.Context) {
		s.runTick(ctx)
		ticker.Reset()
	})

	return nil
}

func (s *Simulator) runTick(ctx context.Context) {
	var wg sync.WaitGroup
	s.lggr.Info("starting congestion")
	ctx, cancel := context.WithTimeout(ctx, s.Duration)
	defer cancel()
	work := make(chan struct{}, *s.Input.MinTxsInPool)
	wg.Add(*s.Input.MinTxsInPool)
	for i := range *s.Input.MinTxsInPool {
		go s.runSender(ctx, &wg, i, work)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.runOnNewBlockUntilCanceled(ctx, func(ctx context.Context, head *types.Header) error {
			for range *s.Input.MinTxsInPool {
				select {
				case work <- struct{}{}:
				case <-ctx.Done():
					return nil
				}
			}
			return nil
		})
	}()
	s.congestionStartedAt = nil
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.runOnNewBlockUntilCanceled(ctx, s.reportCongestion)
	}()
	wg.Wait()
}

func (s *Simulator) runSender(ctx context.Context, wg *sync.WaitGroup, key int, send <-chan struct{}) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case <-send:
			_, err := s.emitter.EmitLogStringsFromKeyAsync(s.payload, key, func(o *bind.TransactOpts) {
				o.GasLimit = *s.Input.GasLimit
			})
			if err != nil {
				s.lggr.Errorw("failed to emit log to fill block", "error", err)
			}
		}
	}
}

func (s *Simulator) reportCongestion(_ context.Context, header *types.Header) error {
	if s.congestionStartedAt == nil {
		s.congestionStartedAt = new(uint64)
		*s.congestionStartedAt = header.Number.Uint64()
	}

	numberOfBlockToInclude := min(maxBlocksToCalcCongestion, header.Number.Uint64()-*s.congestionStartedAt+1)
	congestion, err := s.client.CalculateNetworkCongestionMetric(numberOfBlockToInclude, seth.CongestionStrategy_Simple)
	if err != nil {
		return fmt.Errorf("failed to calculate congestion metric: %v", err)
	}

	block, err := s.client.Client.BlockByNumber(context.Background(), header.Number)
	if err != nil {
		return fmt.Errorf("failed to get block by number: %w", err)
	}

	if congestion < congestionMinTarget {
		s.lggr.Errorf("at block %d congestion is %.2f, which is lower than taget %.2f. Number of txs %d", header.Number.Uint64(), congestion, congestionMinTarget, block.Transactions().Len())
		return nil
	}

	s.lggr.Infof("at block %d congestion is %.2f txs: %d", header.Number.Uint64(), congestion, block.Transactions().Len())
	return nil
}

func (s *Simulator) runOnNewBlockUntilCanceled(ctx context.Context, onNewBlock func(ctx context.Context, head *types.Header) error) {
	for {
		err := s.runOnNewBlock(ctx, onNewBlock)
		if err != nil {
			if ctx.Err() != nil {
				return
			}

			s.lggr.Errorw("failed to run onNewBlock", "error", err)
		}
		time.Sleep(backoffPeriodOnError)
	}
}

func (s *Simulator) runOnNewBlock(ctx context.Context, onNewBlock func(ctx context.Context, head *types.Header) error) error {
	heads := make(chan *types.Header)
	sub, err := s.client.Client.SubscribeNewHead(ctx, heads)
	if err != nil {
		return fmt.Errorf("failed to subscribe to new heads: %v", err)
	}
	defer sub.Unsubscribe()
	for {
		select {
		case head, ok := <-heads:
			if !ok {
				return fmt.Errorf("heads channel closed unexpectedly")
			}
			err = onNewBlock(ctx, head)
			if err != nil {
				if ctx.Err() != nil {
					return ctx.Err()
				}
				s.lggr.Errorw("failed to run onNewBlock", "error", err)
			}
		case <-ctx.Done():
			return ctx.Err()
		case err := <-sub.Err():
			return fmt.Errorf("subscription error: %v", err)
		}
	}
}

func (s *Simulator) estimatePayloadSizeToFillBlock(ctx context.Context) (payloadSize int, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("failed to estimate payload size to fill block: panic: %v", rec)
		}
	}()
	s.lggr.Infof("Estimating payload size to fill block with %d transactions", s.NumberOfTxsPerBlock)
	block, err := s.client.Client.BlockByNumber(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("error getting the latest block: %w", err)
	}

	gasLimit := block.GasLimit()
	maxPayloadSize := gasLimit / 8 / uint64(*s.NumberOfTxsPerBlock) // assume that LogDataGas is 8
	const step = 1000
	result := sort.Search(int(maxPayloadSize/step), func(payloadSize int) bool {
		payloadSize *= step
		var tx *types.Transaction
		tx, err = s.emitter.EmitLogString(string(make([]byte, payloadSize)))
		if err != nil {
			panic(fmt.Errorf("error emitting log emitter tx: %w", err))
		}
		var receipt *types.Receipt
		receipt, err = s.client.Client.TransactionReceipt(ctx, tx.Hash())
		if err != nil {
			panic(fmt.Errorf("failed to get the transaction by hash: %w", err))
		}

		txsToFill := int(gasLimit / receipt.GasUsed)
		s.lggr.Infof("emitted payload size: %d that results in %d transactions to fill a block: %v", payloadSize, txsToFill, err)
		if s.Input.GasLimit == nil {
			s.Input.GasLimit = &receipt.GasUsed
		}

		return txsToFill <= *s.NumberOfTxsPerBlock
	})

	return result * step, nil
}

func (s *Simulator) estimateGasLimit(ctx context.Context) (uint64, error) {
	tx, err := s.emitter.EmitLogString(string(make([]byte, *s.Input.PayloadSize)))
	if err != nil {
		return 0, fmt.Errorf("error emitting log emitter tx: %w", err)
	}
	receipt, err := s.client.Client.TransactionReceipt(ctx, tx.Hash())
	if err != nil {
		return 0, fmt.Errorf("failed to get the transaction by hash: %w", err)
	}

	return receipt.GasUsed, nil
}
