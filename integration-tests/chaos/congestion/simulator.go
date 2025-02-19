package congestion

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"sync"
	"sync/atomic"
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

const (
	backoffPeriodOnError = time.Second
	defaultTxPerBlock    = 10
	defaultTip           = 5_000_000_000
)

type Phase struct {
	Duration   uint64  `toml:"duration"` // in seconds
	Congestion float64 `toml:"congestion"`
}

// Input - defines configuration for chain congestion simulation.
// Simulation has three phases:
// RampUp - period when fees are gradually increased until reach FeesTargetPercent. Congestion is maintained at RampUp.Congestion level.
// Plateau - fees maintained at FeesTargetPercent, congestion - Plateau.Congestion
// CoolDown - fees gradually are decreased to the initial values.
// Delays between full cycles are defined by InitialDelay and Period. During this time simulation does not produce any transactions and
// base fee might fluctuate according EIP-1559 spec depending on Anvil configuration.
type Input struct {
	Enabled             bool    `toml:"enabled"`
	InitialDelay        uint64  `toml:"initial_delay"` // defines delay before the first surge in seconds
	Period              uint64  `toml:"period"`        // defines delay between surges in seconds
	RampUp              Phase   `toml:"ramp_up"`
	Plateau             Phase   `toml:"plateau"`
	CoolDown            Phase   `toml:"cool_down"`
	FeesIncreasePercent float64 `toml:"fees_increase_percent"` // e.g. 0 - fees remain constant; 100 - fees are increased by 100%

	// following values are optional and calculated automatically
	NumberOfTxsPerBlock  int    `toml:"number_of_txs_per_block"`  // defines number of transaction to be produced to fill a block (default: defaultTxPerBlock)
	PayloadSize          *int   `toml:"payload_size"`             // defines payload size that allows NumberOfTxsPerBlock transactions to fully fill a block. (Calculated automatically)
	GasLimit             uint64 `toml:"gas_limit"`                // calculated automatically if not specified
	InitialBaseFeePerGas uint64 `toml:"initial_base_fee_per_gas"` // defaults to base fee of a block fetched on the start, if not specified
	InitialTipPerGas     uint64 `toml:"initial_tip_per_gas"`      // defines initial tip (default: defaultTip)
}

func (i *Input) defaults() {
	if i.NumberOfTxsPerBlock == 0 {
		i.NumberOfTxsPerBlock = defaultTxPerBlock
	}

	if i.InitialTipPerGas == 0 {
		i.InitialTipPerGas = defaultTip
	}
}

type AnvilClient interface {
	AnvilSetNextBlockBaseFeePerGas(gas *big.Int) error
}

// Simulator - produces transactions to generate specified congestion and manipulates fees according to configuration.
// Refer to Input for more details.
type Simulator struct {
	services.Service
	eng *services.Engine

	t *testing.T
	Input
	client      *seth.Client
	anvilClient AnvilClient
	lggr        logger.SugaredLogger

	emitter contracts.LogEmitter
	payload []string

	expectedCurrentTip        atomic.Uint64
	expectedCurrentBaseFee    atomic.Uint64
	expectedCurrentCongestion atomic.Value
	phase                     atomic.Value

	observationsHandler func(chainState)
}

func NewSimulator(t *testing.T, input Input, client *seth.Client, anvilClient AnvilClient, lggr logger.SugaredLogger) (*Simulator, error) {
	input.defaults()

	if input.NumberOfTxsPerBlock > len(client.Addresses) {
		return nil, fmt.Errorf("number_of_txs_per_block %d cannot be greater than number of addresses available for seth.Client %d. Increase number of addresses", input.NumberOfTxsPerBlock, len(client.Addresses))
	}

	if input.RampUp.Duration == 0 && input.Plateau.Duration == 0 && input.CoolDown.Duration == 0 {
		return nil, fmt.Errorf("expected at least one phase to have a positive duration")
	}

	s := &Simulator{
		t:           t,
		Input:       input,
		client:      client,
		anvilClient: anvilClient,
	}
	s.phase.Store(phaseInactive)

	s.Service, s.eng = services.Config{
		Name:  "CongestionSimulator",
		Start: s.start,
	}.NewServiceEngine(lggr)
	s.lggr = s.eng.SugaredLogger
	t.Cleanup(func() {
		_ = s.Close()
	})

	return s, nil
}

type phaseName int

const (
	phaseInactive phaseName = iota
	phaseRampUp
	phasePlateau
	phaseCoolDown
)

type chainState struct {
	PhaseName           phaseName
	Congestion          float64
	TipDeltaPercent     float64
	BaseFeeDeltaPercent float64
}

func (s *Simulator) start(ctx context.Context) error {
	if !s.Input.Enabled {
		s.lggr.Infow("Congestion simulator is not enabled - won't run")
		return nil
	}
	s.expectedCurrentCongestion.Store(float64(0))
	s.expectedCurrentTip.Store(s.InitialTipPerGas)

	if s.InitialBaseFeePerGas == 0 {
		header, err := s.client.Client.HeaderByNumber(ctx, nil)
		if err != nil {
			return fmt.Errorf("failed to get header by number to capture inital base fee: %v", err)
		}

		s.InitialBaseFeePerGas = header.BaseFee.Uint64()
		s.lggr.Debugf("InitialBaseFeePerGas was not provided using: %d", s.InitialBaseFeePerGas)
	}

	s.expectedCurrentBaseFee.Store(s.InitialBaseFeePerGas)

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
	s.payload = []string{string(make([]byte, *s.Input.PayloadSize))}

	if s.Input.GasLimit == 0 {
		gasLimit, err := s.estimateTxGasLimit(ctx)
		if err != nil {
			return fmt.Errorf("failed to estimate gas limit: %v", err)
		}
		s.Input.GasLimit = gasLimit
	}

	ticker := services.TickerConfig{Initial: time.Second * time.Duration(s.InitialDelay)}.NewTicker(time.Second * time.Duration(s.Period))
	s.eng.GoTick(ticker, func(ctx context.Context) {
		s.runTick(ctx)
		ticker.Reset() // ensure that each iteration is delayed by s.Period
	})
	s.eng.Go(func(ctx context.Context) {
		s.listenHeadsUntilCanceled(ctx, s.reportObservationOnNewHead)
	})
	return nil
}

func (s *Simulator) startTransactionProducer(ctx context.Context, wg *sync.WaitGroup) {
	numberOfSenders := s.NumberOfTxsPerBlock
	work := make(chan struct{}, numberOfSenders)
	wg.Add(numberOfSenders)
	for i := range numberOfSenders {
		go s.runSender(ctx, wg, i, work)
	}

	txsToSend := numberOfSenders / 2
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.listenHeadsUntilCanceled(ctx, func(ctx context.Context, head *types.Header) error {
			congestion := float64(head.GasUsed) / float64(head.GasLimit)
			congestionTarget := s.expectedCurrentCongestion.Load().(float64)
			if txsToSend < numberOfSenders && congestion < congestionTarget*0.95 {
				txsToSend++
			} else if txsToSend > 0 && congestion > congestionTarget*1.05 {
				txsToSend--
			}

			s.lggr.Infof("Sending %d transactions; current congestion: %.2f target: %.2f", txsToSend, congestion, congestionTarget)
			for range txsToSend {
				select {
				case work <- struct{}{}:
				case <-ctx.Done():
					return nil
				}
			}
			return nil
		})
	}()
}

func (s *Simulator) runPhaseController(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	s.lggr.Debugf("Starting ramp up phase")
	s.phase.Store(phaseRampUp)
	s.updateCurrentFees(0)
	s.handleVolatilePhase(ctx, s.RampUp, true)

	s.lggr.Debugf("Starting plateue phase")
	s.phase.Store(phasePlateau)
	s.expectedCurrentCongestion.Store(s.Plateau.Congestion)
	timer := time.NewTimer(time.Second * time.Duration(s.Plateau.Duration))
	defer timer.Stop()
	select {
	case <-timer.C:
		break
	case <-ctx.Done():
		return
	}

	s.lggr.Debugf("Starting cooldown phase")
	s.phase.Store(phaseCoolDown)
	s.handleVolatilePhase(ctx, s.CoolDown, false)

	s.phase.Store(phaseInactive)
}

func (s *Simulator) handleVolatilePhase(ctx context.Context, phase Phase, isIncrease bool) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	phaseTimer := time.NewTimer(time.Second * time.Duration(phase.Duration))
	defer phaseTimer.Stop()
	phaseStart := time.Now()
	s.expectedCurrentCongestion.Store(phase.Congestion)
phaseLoop:
	for {
		select {
		case <-ticker.C:
		case <-phaseTimer.C:
			break phaseLoop
		case <-ctx.Done():
			return
		}

		elapsedRatio := min(float64(time.Since(phaseStart))/float64(time.Duration(phase.Duration)*time.Second), 1)
		if !isIncrease {
			// count backwards. If 0.1 of cooldown phase has passed, we should be at the same price point as if 0.9 of ramp up phase has passed
			elapsedRatio = 1 - elapsedRatio
		}
		s.updateCurrentFees(elapsedRatio)
	}

	if isIncrease {
		s.updateCurrentFees(1)
	} else {
		s.updateCurrentFees(0)
	}
}

// updateCurrentFees - sets fees according to expected distance from target. If targetDistanceRatio = 1, current fees s.FeesIncreasePercent compared to initial.
func (s *Simulator) updateCurrentFees(targetDistanceRatio float64) {
	targetPercent := s.FeesIncreasePercent * targetDistanceRatio
	baseFee := uint64(float64(s.InitialBaseFeePerGas) * (1 + targetPercent/100))
	s.lggr.Debugf("Updating fees to %d according to distance %.2f", baseFee, targetDistanceRatio)
	s.expectedCurrentBaseFee.Store(baseFee)
	tip := uint64(float64(s.InitialTipPerGas) * (1 + targetPercent/100))
	s.expectedCurrentTip.Store(tip)
}

func (s *Simulator) runTick(ctx context.Context) {
	var wg sync.WaitGroup
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(s.RampUp.Duration+s.Plateau.Duration+s.CoolDown.Duration))
	defer cancel()
	wg.Add(1)
	go s.runPhaseController(ctx, &wg)
	s.startTransactionProducer(ctx, &wg)
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.listenHeadsUntilCanceled(ctx, s.updateBaseFee)
	}()
	wg.Wait()
}

func (s *Simulator) updateBaseFee(_ context.Context, _ *types.Header) error {
	err := s.anvilClient.AnvilSetNextBlockBaseFeePerGas(big.NewInt(0).SetUint64(s.expectedCurrentBaseFee.Load()))
	if err != nil {
		return fmt.Errorf("failed to set next block base fee per gas: %w", err)
	}

	return nil
}

func (s *Simulator) runSender(ctx context.Context, wg *sync.WaitGroup, key int, send <-chan struct{}) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case <-send:
			_, err := s.emitter.EmitLogStringsFromKeyAsync(s.payload, key, func(o *bind.TransactOpts) {
				o.GasLimit = s.Input.GasLimit
				o.GasTipCap = big.NewInt(0).SetUint64(s.expectedCurrentTip.Load())
				o.GasFeeCap = big.NewInt(0).Add(big.NewInt(0).SetUint64(s.expectedCurrentBaseFee.Load()), o.GasTipCap)
			})
			if err != nil {
				s.lggr.Errorw("failed to emit log to fill block", "error", err)
			}
		}
	}
}

func (s *Simulator) reportObservationOnNewHead(ctx context.Context, header *types.Header) error {
	congestion := float64(header.GasUsed) / float64(header.GasLimit)
	baseFeePercent := float64(header.BaseFee.Uint64()) / float64(s.Input.InitialBaseFeePerGas) * 100
	block, err := s.client.Client.BlockByNumber(ctx, header.Number)
	if err != nil {
		return fmt.Errorf("failed to fetch block by number: %w", err)
	}
	tipPercent := float64(s.avgTip(block).Uint64()) / float64(s.InitialTipPerGas) * 100

	state := chainState{
		PhaseName:           s.phase.Load().(phaseName),
		Congestion:          congestion,
		TipDeltaPercent:     tipPercent - 100,
		BaseFeeDeltaPercent: baseFeePercent - 100,
	}
	if s.observationsHandler != nil {
		s.observationsHandler(state)
	}

	s.lggr.Debugw("New chain state report", "chainState", state)
	return nil
}

func (s *Simulator) avgTip(block *types.Block) *big.Int {
	if block.Transactions().Len() == 0 {
		return big.NewInt(0)
	}
	total := big.NewInt(0)
	for _, tx := range block.Transactions() {
		total = total.Add(tx.EffectiveGasTipValue(block.BaseFee()), total)
	}

	return big.NewInt(0).Div(total, big.NewInt(int64(block.Transactions().Len())))
}

func (s *Simulator) listenHeadsUntilCanceled(ctx context.Context, onNewBlock func(ctx context.Context, head *types.Header) error) {
	for {
		err := s.listenHeads(ctx, onNewBlock)
		if err != nil {
			if ctx.Err() != nil {
				return
			}

			s.lggr.Errorw("failed to run onNewBlock", "error", err)
		}
		time.Sleep(backoffPeriodOnError)
	}
}

func (s *Simulator) listenHeads(ctx context.Context, oneNewHead func(ctx context.Context, head *types.Header) error) error {
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
			err = oneNewHead(ctx, head)
			if err != nil {
				if ctx.Err() != nil {
					return ctx.Err()
				}
				s.lggr.Errorw("failed to run oneNewHead", "error", err)
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
	maxPayloadSize := int(gasLimit / 8 / uint64(s.NumberOfTxsPerBlock)) // assume that LogDataGas is 8
	const step = 100
	// Find the largest payload that won't overflow gas limit if we send NumberOfTxsPerBlock txs
	i := sort.Search(maxPayloadSize/step, func(i int) bool {
		payloadSize := maxPayloadSize - i*step
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

		estimatedCongestion := float64(s.NumberOfTxsPerBlock) * float64(receipt.GasUsed) / float64(gasLimit)
		s.lggr.Infof("emitted payload size: %d that results in %d gas used. If we send %d transactions, congestions will be %.2f", payloadSize, receipt.GasUsed, s.NumberOfTxsPerBlock, estimatedCongestion)
		return uint64(s.NumberOfTxsPerBlock)*receipt.GasUsed <= gasLimit
	})

	payloadSize = maxPayloadSize - i*step
	s.lggr.Infof("Using payload of size %d to fill gas limit: %d with %d txs", payloadSize, gasLimit, s.NumberOfTxsPerBlock)

	return payloadSize, nil
}

func (s *Simulator) estimateTxGasLimit(ctx context.Context) (uint64, error) {
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
