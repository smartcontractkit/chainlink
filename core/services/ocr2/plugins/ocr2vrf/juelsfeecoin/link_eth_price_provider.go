package juelsfeecoin

import (
	"context"
	"math/big"
	"runtime"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-vrf/types"

	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/aggregator_v3_interface"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// linkEthPriceProvider provides conversation rate between Link and native token using price feeds
type linkEthPriceProvider struct {
	aggregator             aggregator_v3_interface.AggregatorV3InterfaceInterface
	timeout                time.Duration
	interval               time.Duration
	lock                   sync.RWMutex
	stop                   chan struct{}
	currentJuelsPerFeeCoin *big.Int
	lggr                   logger.Logger
}

var _ types.JuelsPerFeeCoin = (*linkEthPriceProvider)(nil)

func NewLinkEthPriceProvider(
	linkEthFeedAddress common.Address,
	client evmclient.Client,
	timeout time.Duration,
	interval time.Duration,
	logger logger.Logger,
) (types.JuelsPerFeeCoin, error) {
	aggregator, err := aggregator_v3_interface.NewAggregatorV3Interface(linkEthFeedAddress, client)
	if err != nil {
		return nil, errors.Wrap(err, "new aggregator v3 interface")
	}

	if timeout >= interval {
		return nil, errors.New("timeout must be less than interval")
	}

	p := &linkEthPriceProvider{
		aggregator:             aggregator,
		timeout:                timeout,
		interval:               interval,
		currentJuelsPerFeeCoin: big.NewInt(0),
		stop:                   make(chan struct{}),
		lggr:                   logger,
	}

	// Begin updating JuelsPerFeeCoin.
	// Stop fetching price updates on garbage collection, as to avoid a leaked goroutine.
	go p.run()
	runtime.SetFinalizer(p, func(p *linkEthPriceProvider) { p.stop <- struct{}{} })

	return p, nil
}

// Run updates the JuelsPerFeeCoin value at a regular interval, until stopped.
// Do not block the main thread, such that updates are always timely.
func (p *linkEthPriceProvider) run() {
	ticker := time.NewTicker(p.interval)
	for {
		select {
		case <-ticker.C:
			go p.updateJuelsPerFeeCoin()
		case <-p.stop:
			ticker.Stop()
			return
		}
	}
}

// JuelsPerFeeCoin returns the current JuelsPerFeeCoin value, threadsafe.
func (p *linkEthPriceProvider) JuelsPerFeeCoin() (*big.Int, error) {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return p.currentJuelsPerFeeCoin, nil
}

// Get current JuelsPerFeeCoin value from aggregator contract.
// If the RPC call fails, log the error and return.
func (p *linkEthPriceProvider) updateJuelsPerFeeCoin() {
	// Ensure writes to currentJuelsPerFeeCoin are threadsafe.
	p.lock.Lock()
	defer p.lock.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()
	roundData, err := p.aggregator.LatestRoundData(&bind.CallOpts{Context: ctx})

	// For RPC issues, set the most recent price to 0. This way, stale prices will not be transmitted,
	// and zero-observations can be ignored in OCR and on-chain.
	if err != nil {
		p.currentJuelsPerFeeCoin = big.NewInt(0)
		return
	}

	// Update JuelsPerFeeCoin to the obtained value.
	p.currentJuelsPerFeeCoin = roundData.Answer
}
