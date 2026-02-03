package automation

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/initial/log_emitter"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
	"github.com/smartcontractkit/chainlink/devenv/contracts"
	"github.com/smartcontractkit/chainlink/devenv/products/automation"
	ctf_concurrency "github.com/smartcontractkit/chainlink/devenv/products/automation/concurrency"
)

type DeploymentData struct {
	ConsumerContracts []contracts.KeeperConsumer
	TriggerContracts  []contracts.LogEmitter
	TriggerAddresses  []common.Address
	LoadConfigs       []Load
}

type deployedContractData struct {
	consumerContract contracts.KeeperConsumer
	triggerContract  contracts.LogEmitter
	triggerAddress   common.Address
	loadConfig       Load
}

func (d deployedContractData) GetResult() deployedContractData {
	return d
}

type task struct {
	deployTrigger bool
}

func deployConsumerAndTriggerContracts(l zerolog.Logger, tc loadtestcase, chainClient *seth.Client, multicallAddress common.Address, maxConcurrency int, automationDefaultLinkFunds *big.Int, linkToken contracts.LinkToken) (DeploymentData, error) {
	data := DeploymentData{}

	concurrency, err := automation.GetAndAssertCorrectConcurrency(chainClient, 1)
	if err != nil {
		return DeploymentData{}, err
	}

	if concurrency > maxConcurrency {
		concurrency = maxConcurrency
		l.Debug().
			Msgf("Concurrency is higher than max concurrency, setting concurrency to %d", concurrency)
	}

	l.Debug().
		Int("Number of Upkeeps", tc.UpkeepCount).
		Int("Concurrency", concurrency).
		Msg("Deployment parallelisation info")

	tasks := []task{}
	for i := 0; i < tc.UpkeepCount; i++ {
		if tc.SharedTrigger {
			if i == 0 {
				tasks = append(tasks, task{deployTrigger: true})
			} else {
				tasks = append(tasks, task{deployTrigger: false})
			}
			continue
		}
		tasks = append(tasks, task{deployTrigger: true})
	}

	var deployContractFn = func(deployedCh chan deployedContractData, errorCh chan error, keyNum int, task task) {
		data := deployedContractData{}
		consumerContract, err := contracts.DeployAutomationSimpleLogTriggerConsumerFromKey(chainClient, tc.IsStreamsLookup, keyNum)
		if err != nil {
			errorCh <- errors.Wrapf(err, "Error deploying simple log trigger contract")
			return
		}

		data.consumerContract = consumerContract

		loadCfg := Load{
			NumberOfEvents:                tc.NumberOfEvents,
			NumberOfSpamMatchingEvents:    tc.NumberOfSpamMatchingEvents,
			NumberOfSpamNonMatchingEvents: tc.NumberOfSpamNonMatchingEvents,
			CheckBurnAmount:               tc.CheckBurnAmount,
			PerformBurnAmount:             tc.PerformBurnAmount,
			UpkeepGasLimit:                tc.UpkeepGasLimit,
			SharedTrigger:                 tc.SharedTrigger,
			Feeds:                         []string{},
		}

		if tc.IsStreamsLookup {
			loadCfg.Feeds = tc.Feeds
		}

		data.loadConfig = loadCfg

		if !task.deployTrigger {
			deployedCh <- data
			return
		}

		triggerContract, err := contracts.DeployLogEmitterContractFromKey(l, chainClient, keyNum)
		if err != nil {
			errorCh <- errors.Wrapf(err, "Error deploying log emitter contract")
			return
		}

		data.triggerContract = triggerContract
		data.triggerAddress = triggerContract.Address()
		deployedCh <- data
	}

	executor := ctf_concurrency.NewConcurrentExecutor[deployedContractData, deployedContractData, task](l)
	results, err := executor.Execute(concurrency, tasks, deployContractFn)
	if err != nil {
		return DeploymentData{}, err
	}

	for _, result := range results {
		if result.GetResult().triggerContract != nil {
			data.TriggerContracts = append(data.TriggerContracts, result.GetResult().triggerContract)
			data.TriggerAddresses = append(data.TriggerAddresses, result.GetResult().triggerAddress)
		}
		data.ConsumerContracts = append(data.ConsumerContracts, result.GetResult().consumerContract)
		data.LoadConfigs = append(data.LoadConfigs, result.GetResult().loadConfig)
	}

	// if there's more than 1 upkeep and it's a shared trigger, then we should use only the first address in triggerAddresses
	// as triggerAddresses array
	if tc.SharedTrigger {
		if len(data.TriggerAddresses) == 0 {
			return DeploymentData{}, errors.New("No trigger addresses found")
		}
		triggerAddress := data.TriggerAddresses[0]
		data.TriggerAddresses = make([]common.Address, 0)
		for i := 0; i < tc.UpkeepCount; i++ {
			data.TriggerAddresses = append(data.TriggerAddresses, triggerAddress)
		}
	}

	sendErr := automation.SendLinkFundsToDeploymentAddresses(chainClient, concurrency, tc.UpkeepCount, tc.UpkeepCount/concurrency, multicallAddress, automationDefaultLinkFunds, linkToken)
	if sendErr != nil {
		return DeploymentData{}, sendErr
	}

	return data, nil
}

type LogTriggerConfig struct {
	Address                       string
	NumberOfEvents                int64
	NumberOfSpamMatchingEvents    int64
	NumberOfSpamNonMatchingEvents int64
}

type LogTriggerGun struct {
	data             [][]byte
	addresses        []string
	multiCallAddress string
	client           *seth.Client
	logger           zerolog.Logger

	keyPool *KeyPool
}

func generateCallData(int1 int64, int2 int64, count int64) []byte {
	abi, err := log_emitter.LogEmitterMetaData.GetAbi()
	if err != nil {
		panic(err)
	}
	data, err := abi.Pack("EmitLog4", big.NewInt(int1), big.NewInt(int2), big.NewInt(count))
	if err != nil {
		panic(err)
	}
	return data
}

func NewLogTriggerUser(
	logger zerolog.Logger,
	triggerConfigs []LogTriggerConfig,
	client *seth.Client,
	multicallAddress string,
	keyPool *KeyPool,
) (*LogTriggerGun, error) {
	var data [][]byte
	var addresses []string

	// we need to sync nodes manually, because we are not using ephemeral addresses
	if err := client.NonceManager.UpdateNonces(); err != nil {
		return nil, err
	}

	for _, c := range triggerConfigs {
		if c.NumberOfEvents > 0 {
			d := generateCallData(1, 1, c.NumberOfEvents)
			data = append(data, d)
			addresses = append(addresses, c.Address)
		}
		if c.NumberOfSpamMatchingEvents > 0 {
			d := generateCallData(1, 2, c.NumberOfSpamMatchingEvents)
			data = append(data, d)
			addresses = append(addresses, c.Address)
		}
		if c.NumberOfSpamNonMatchingEvents > 0 {
			d := generateCallData(2, 2, c.NumberOfSpamNonMatchingEvents)
			data = append(data, d)
			addresses = append(addresses, c.Address)
		}
	}

	return &LogTriggerGun{
		addresses:        addresses,
		data:             data,
		logger:           logger,
		multiCallAddress: multicallAddress,
		client:           client,
		keyPool:          keyPool,
	}, nil
}

func (m *LogTriggerGun) Call(_ *wasp.Generator) *wasp.Response {
	var wg sync.WaitGroup
	var dividedData [][][]byte
	d := m.data
	chunkSize := 100
	for i := 0; i < len(d); i += chunkSize {
		end := min(i+chunkSize, len(d))
		dividedData = append(dividedData, d[i:end])
	}

	resultCh := make(chan *wasp.Response, len(dividedData))

	for _, a := range dividedData {
		wg.Add(1)
		go func(a [][]byte, m *LogTriggerGun) {
			defer wg.Done()

			keyCtx, cancelFn := context.WithTimeout(context.Background(), m.keyPool.RecommendedCheckoutTimeout())
			defer cancelFn()
			keyIndex, nonce, err := m.keyPool.CheckoutKey(keyCtx)
			if err != nil {
				m.logger.Error().Err(err).Msg("Error checking out key from key pool")
				_ = m.keyPool.DiagnoseAndMonitor(60 * time.Second)
				if strings.Contains(err.Error(), "all keys have pending transactions") {
					dropCtx, dropCancelFn := context.WithTimeout(context.Background(), m.keyPool.RecommendedDropTimeout())
					defer dropCancelFn()
					dropped, dropErr := m.keyPool.DropPendingTxs(dropCtx)
					if dropErr != nil {
						m.logger.Error().Err(dropErr).Msg("Error dropping pending transactions")
					} else {
						m.logger.Info().Int("dropped", dropped).Msg("Dropped pending transactions")
					}
				}
				resultCh <- &wasp.Response{Error: err.Error(), Failed: true}
				return
			}

			tx, err := contracts.MultiCallLogTriggerLoadGen(m.client, keyIndex+1, big.NewInt(int64(nonce)), m.multiCallAddress, m.addresses, a) //nolint:gosec // we will never have that many keys to cause an overflow
			if err != nil {
				m.logger.Error().Err(err).Msg("Error calling MultiCallLogTriggerLoadGen")
				_ = m.keyPool.DiagnoseAndMonitor(60 * time.Second)
				resultCh <- &wasp.Response{Error: err.Error(), Failed: true}
				return
			}
			m.keyPool.RecordPendingTx(keyIndex, tx.Hash())
			resultCh <- &wasp.Response{}
		}(a, m)
	}
	wg.Wait()
	close(resultCh)

	r := &wasp.Response{}
	for result := range resultCh {
		if result.Failed {
			r.Failed = true
			if r.Error != "" {
				r.Error += "; " + result.Error
			} else {
				r.Error = result.Error
			}
		}
	}

	return r
}

// intListStats helper calculates some statistics on an int list: avg, median, 90pct, 99pct, max
//
//nolint:revive // we know what each int64 we return means
func IntListStats(in []int64) (float64, int64, int64, int64, int64) {
	length := len(in)
	if length == 0 {
		return 0, 0, 0, 0, 0
	}
	slices.Sort(in)
	var sum int64
	for _, num := range in {
		sum += num
	}
	return float64(sum) / float64(length), in[int(math.Floor(float64(length)*0.5))], in[int(math.Floor(float64(length)*0.9))], in[int(math.Floor(float64(length)*0.99))], in[length-1]
}

type KeyPool struct {
	mu        sync.Mutex
	dropMu    sync.Mutex // protects DropPendingTxs from concurrent calls
	client    *ethclient.Client
	rpcClient *rpc.Client
	logger    zerolog.Logger
	addresses []common.Address
	nextIndex int

	isAnvil bool
	// pendingTxs tracks the last pending tx hash for each key index
	pendingTxs   map[int]common.Hash
	checkTimeout time.Duration
}

func NewKeyPool(logger zerolog.Logger, rpcURL string, addrs []common.Address, isAnvil bool) (*KeyPool, error) {
	rpcClient, err := rpc.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to dial RPC: %w", err)
	}

	logger.Info().
		Int("keyCount", len(addrs)).
		Bool("isAnvil", isAnvil).
		Msg("Initialized KeyPool")

	return &KeyPool{
		client:       ethclient.NewClient(rpcClient),
		rpcClient:    rpcClient,
		logger:       logger,
		addresses:    addrs,
		pendingTxs:   make(map[int]common.Hash),
		checkTimeout: 5 * time.Second,
		isAnvil:      isAnvil,
	}, nil
}

// RecommendedCheckoutTimeout returns a recommended timeout for CheckoutKey.
// Based on: keyCount * perKeyTimeout * safetyMultiplier
// The safety multiplier accounts for potential RPC slowdowns.
func (p *KeyPool) RecommendedCheckoutTimeout() time.Duration {
	// Each key check has checkTimeout (5s default), but RPC should be <100ms normally
	// Use 1s per key as baseline with 30s minimum
	timeout := time.Duration(len(p.addresses)) * time.Second
	if timeout < 30*time.Second {
		timeout = 30 * time.Second
	}
	return timeout
}

// RecommendedDropTimeout returns a recommended timeout for DropPendingTxs.
// Based on: number of pending txs * perDropTimeout
func (p *KeyPool) RecommendedDropTimeout() time.Duration {
	p.mu.Lock()
	pendingCount := len(p.pendingTxs)
	p.mu.Unlock()

	// 1s per pending tx with 30s minimum
	timeout := time.Duration(pendingCount) * time.Second
	if timeout < 30*time.Second {
		timeout = 30 * time.Second
	}
	return timeout
}

// CheckoutKey finds the next key with no pending transactions.
// Returns key index and the nonce to use.
func (p *KeyPool) CheckoutKey(ctx context.Context) (int, uint64, error) {
	checkedCount := 0
	skippedPending := 0

	for i := 0; i < len(p.addresses); i++ {
		p.mu.Lock()
		idx := p.nextIndex
		p.nextIndex = (p.nextIndex + 1) % len(p.addresses)
		p.mu.Unlock()

		checkedCount++
		nonce, hasPending, err := p.checkKey(ctx, idx)
		if err != nil {
			p.logger.Trace().
				Err(err).
				Int("keyIndex", idx).
				Msg("Error checking key, trying next")
			continue // try next key
		}
		if hasPending {
			skippedPending++
			p.logger.Trace().
				Int("keyIndex", idx).
				Str("address", p.addresses[idx].Hex()).
				Msg("Key has pending tx, skipping")
			continue
		}

		// Clear tracked tx since it's confirmed
		p.mu.Lock()
		delete(p.pendingTxs, idx)
		p.mu.Unlock()

		p.logger.Debug().
			Int("keyIndex", idx).
			Uint64("nonce", nonce).
			Int("checkedKeys", checkedCount).
			Int("skippedPending", skippedPending).
			Msg("Checked out key")
		return idx, nonce, nil
	}

	p.logger.Warn().
		Int("totalKeys", len(p.addresses)).
		Int("checkedKeys", checkedCount).
		Int("skippedPending", skippedPending).
		Msg("All keys have pending transactions")
	return -1, 0, fmt.Errorf("all %d keys have pending transactions", len(p.addresses))
}

// checkKey returns current nonce and whether key has pending txs
func (p *KeyPool) checkKey(ctx context.Context, idx int) (nonce uint64, hasPending bool, err error) {
	start := time.Now()
	ctx, cancel := context.WithTimeout(ctx, p.checkTimeout)
	defer cancel()

	addr := p.addresses[idx]

	var pending, latest uint64
	var pendingErr, latestErr error
	var pendingDuration, latestDuration time.Duration
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		callStart := time.Now()
		pending, pendingErr = p.client.PendingNonceAt(ctx, addr)
		pendingDuration = time.Since(callStart)
	}()
	go func() {
		defer wg.Done()
		callStart := time.Now()
		latest, latestErr = p.client.NonceAt(ctx, addr, nil)
		latestDuration = time.Since(callStart)
	}()
	wg.Wait()

	totalDuration := time.Since(start)

	// Log timing at trace level (only visible with very verbose logging)
	p.logger.Trace().
		Int("keyIndex", idx).
		Dur("pendingNonceMs", pendingDuration).
		Dur("latestNonceMs", latestDuration).
		Dur("totalMs", totalDuration).
		Msg("checkKey RPC timing")

	// Warn if RPC calls are slow (>500ms)
	if pendingDuration > 500*time.Millisecond || latestDuration > 500*time.Millisecond {
		p.logger.Warn().
			Int("keyIndex", idx).
			Dur("pendingNonceMs", pendingDuration).
			Dur("latestNonceMs", latestDuration).
			Msg("Slow RPC detected in checkKey")
	}

	if pendingErr != nil {
		return 0, false, pendingErr
	}
	if latestErr != nil {
		return 0, false, latestErr
	}

	return latest, pending > latest, nil
}

// RecordPendingTx records a pending transaction hash for a key.
// Call this after successfully sending a transaction.
func (p *KeyPool) RecordPendingTx(idx int, txHash common.Hash) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.pendingTxs[idx] = txHash
	p.logger.Debug().
		Int("keyIndex", idx).
		Str("txHash", txHash.Hex()).
		Int("totalPending", len(p.pendingTxs)).
		Msg("Recorded pending tx")
}

// DropPendingTxs drops all tracked pending transactions from Anvil's mempool.
// Call this manually when CheckoutKey returns an error and you want to recover.
// Returns the number of transactions dropped and any errors encountered.
func (p *KeyPool) DropPendingTxs(ctx context.Context) (int, error) {
	if !p.isAnvil {
		p.logger.Debug().Msg("DropPendingTxs called but not running on Anvil, skipping")
		return 0, nil
	}

	// Prevent multiple goroutines from dropping simultaneously
	p.dropMu.Lock()
	defer p.dropMu.Unlock()

	p.mu.Lock()
	// Copy the map to avoid holding the lock during RPC calls
	toDropMap := make(map[int]common.Hash, len(p.pendingTxs))
	for k, v := range p.pendingTxs {
		toDropMap[k] = v
	}
	p.mu.Unlock()

	if len(toDropMap) == 0 {
		p.logger.Debug().Msg("No pending transactions to drop")
		return 0, nil
	}

	p.logger.Info().
		Int("pendingTxCount", len(toDropMap)).
		Msg("Dropping pending transactions from Anvil mempool")

	var dropped int
	var errs []string

	for idx, txHash := range toDropMap {
		if err := p.dropTransaction(ctx, txHash); err != nil {
			errs = append(errs, fmt.Sprintf("key %d (%s): %v", idx, txHash.Hex(), err))
			continue
		}

		p.mu.Lock()
		delete(p.pendingTxs, idx)
		p.mu.Unlock()

		dropped++
		p.logger.Debug().
			Int("keyIndex", idx).
			Str("txHash", txHash.Hex()).
			Msg("Dropped pending transaction")
	}

	if len(errs) > 0 {
		p.logger.Warn().
			Int("dropped", dropped).
			Int("failed", len(errs)).
			Msg("Finished dropping pending transactions with some failures")
		return dropped, fmt.Errorf("failed to drop some transactions: %s", strings.Join(errs, "; "))
	}

	p.logger.Info().
		Int("dropped", dropped).
		Msg("Successfully dropped all pending transactions")
	return dropped, nil
}

// dropTransaction drops a transaction from Anvil's mempool using anvil_dropTransaction
func (p *KeyPool) dropTransaction(ctx context.Context, txHash common.Hash) error {
	ctx, cancel := context.WithTimeout(ctx, p.checkTimeout)
	defer cancel()

	var result interface{}
	if err := p.rpcClient.CallContext(ctx, &result, "anvil_dropTransaction", txHash); err != nil {
		return fmt.Errorf("anvil_dropTransaction failed: %w", err)
	}
	return nil
}

// StartBlockMonitor starts a background goroutine that monitors block progress.
// It logs the current block number every `interval` for `duration`.
// Useful for debugging when the chain might be stuck.
// Returns a cancel function to stop monitoring early.
func (p *KeyPool) StartBlockMonitor(duration time.Duration, interval time.Duration) context.CancelFunc {
	ctx, cancel := context.WithTimeout(context.Background(), duration)

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		var lastBlock uint64
		iterations := 0

		p.logger.Info().
			Dur("duration", duration).
			Dur("interval", interval).
			Msg("Starting block progress monitor")

		for {
			select {
			case <-ctx.Done():
				p.logger.Info().
					Int("iterations", iterations).
					Msg("Block monitor stopped")
				return
			case <-ticker.C:
				iterations++
				callStart := time.Now()
				block, err := p.client.BlockNumber(ctx)
				callDuration := time.Since(callStart)

				if err != nil {
					p.logger.Error().
						Err(err).
						Dur("rpcLatencyMs", callDuration).
						Msg("Block monitor: failed to get block number")
					continue
				}

				blockDelta := int64(0)
				if lastBlock > 0 {
					blockDelta = int64(block) - int64(lastBlock)
				}

				logEvent := p.logger.Info().
					Uint64("blockNumber", block).
					Int64("blockDelta", blockDelta).
					Dur("rpcLatencyMs", callDuration)

				if blockDelta == 0 && lastBlock > 0 {
					logEvent.Msg("Block monitor: NO PROGRESS - chain may be stuck!")
				} else {
					logEvent.Msg("Block monitor: progress update")
				}

				lastBlock = block
			}
		}
	}()

	return cancel
}

// DiagnoseAndMonitor runs diagnostics when checkout fails.
// It starts block monitoring and logs current state.
// Call this when CheckoutKey returns an error.
func (p *KeyPool) DiagnoseAndMonitor(monitorDuration time.Duration) context.CancelFunc {
	p.mu.Lock()
	pendingCount := len(p.pendingTxs)
	pendingKeys := make([]int, 0, pendingCount)
	for k := range p.pendingTxs {
		pendingKeys = append(pendingKeys, k)
	}
	p.mu.Unlock()

	p.logger.Warn().
		Int("totalKeys", len(p.addresses)).
		Int("pendingTxCount", pendingCount).
		Ints("pendingKeyIndexes", pendingKeys).
		Msg("Diagnosis: KeyPool state at failure")

	// Start block monitor for the specified duration, checking every 2 seconds
	return p.StartBlockMonitor(monitorDuration, 2*time.Second)
}
