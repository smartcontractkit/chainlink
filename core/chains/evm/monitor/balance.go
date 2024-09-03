package monitor

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"sync"
	"time"

	gethCommon "github.com/ethereum/go-ethereum/common"
	pkgerrors "github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	httypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/keystore"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

type (
	// BalanceMonitor checks the balance for each key on every new head
	BalanceMonitor interface {
		httypes.HeadTrackable
		GetEthBalance(gethCommon.Address) *assets.Eth
		services.Service
	}

	balanceMonitor struct {
		services.Service
		eng *services.Engine

		ethClient      evmclient.Client
		chainID        *big.Int
		chainIDStr     string
		ethKeyStore    keystore.Eth
		ethBalances    map[gethCommon.Address]*assets.Eth
		ethBalancesMtx sync.RWMutex
		sleeperTask    *utils.SleeperTask
	}

	NullBalanceMonitor struct{}
)

var _ BalanceMonitor = (*balanceMonitor)(nil)

// NewBalanceMonitor returns a new balanceMonitor
func NewBalanceMonitor(ethClient evmclient.Client, ethKeyStore keystore.Eth, lggr logger.Logger) *balanceMonitor {
	chainId := ethClient.ConfiguredChainID()
	bm := &balanceMonitor{
		ethClient:   ethClient,
		chainID:     chainId,
		chainIDStr:  chainId.String(),
		ethKeyStore: ethKeyStore,
		ethBalances: make(map[gethCommon.Address]*assets.Eth),
	}
	bm.Service, bm.eng = services.Config{
		Name:  "BalanceMonitor",
		Start: bm.start,
		Close: bm.close,
	}.NewServiceEngine(lggr)
	bm.sleeperTask = utils.NewSleeperTask(&worker{bm: bm})
	return bm
}

func (bm *balanceMonitor) start(ctx context.Context) error {
	// Always query latest balance on start
	(&worker{bm}).WorkCtx(ctx)
	return nil
}

// Close shuts down the BalanceMonitor, should not be used after this
func (bm *balanceMonitor) close() error {
	return bm.sleeperTask.Stop()
}

// OnNewLongestChain checks the balance for each key
func (bm *balanceMonitor) OnNewLongestChain(_ context.Context, _ *evmtypes.Head) {
	bm.eng.Debugw("BalanceMonitor: signalling balance worker")
	ok := bm.sleeperTask.WakeUpIfStarted()
	if !ok {
		bm.eng.Debugw("BalanceMonitor: ignoring OnNewLongestChain call, balance monitor is not started", "state", bm.sleeperTask.State())
	}
}

func (bm *balanceMonitor) updateBalance(ethBal assets.Eth, address gethCommon.Address) {
	bm.promUpdateEthBalance(&ethBal, address)

	bm.ethBalancesMtx.Lock()
	oldBal := bm.ethBalances[address]
	bm.ethBalances[address] = &ethBal
	bm.ethBalancesMtx.Unlock()

	lgr := logger.Named(bm.eng, "BalanceLog")
	lgr = logger.With(lgr,
		"address", address.Hex(),
		"ethBalance", ethBal.String(),
		"weiBalance", ethBal.ToInt())

	if oldBal == nil {
		lgr.Infof("ETH balance for %s: %s", address.Hex(), ethBal.String())
		return
	}

	if ethBal.Cmp(oldBal) != 0 {
		lgr.Infof("New ETH balance for %s: %s", address.Hex(), ethBal.String())
	}
}

func (bm *balanceMonitor) GetEthBalance(address gethCommon.Address) *assets.Eth {
	bm.ethBalancesMtx.RLock()
	defer bm.ethBalancesMtx.RUnlock()
	return bm.ethBalances[address]
}

var promETHBalance = promauto.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "eth_balance",
		Help: "Each Ethereum account's balance",
	},
	[]string{"account", "evmChainID"},
)

func (bm *balanceMonitor) promUpdateEthBalance(balance *assets.Eth, from gethCommon.Address) {
	balanceFloat, err := ApproximateFloat64(balance)

	if err != nil {
		bm.eng.Error(fmt.Errorf("updatePrometheusEthBalance: %v", err))
		return
	}

	promETHBalance.WithLabelValues(from.Hex(), bm.chainIDStr).Set(balanceFloat)
}

type worker struct {
	bm *balanceMonitor
}

func (*worker) Name() string {
	return "BalanceMonitorWorker"
}

func (w *worker) Work() {
	// Used with SleeperTask
	w.WorkCtx(context.Background())
}

func (w *worker) WorkCtx(ctx context.Context) {
	enabledAddresses, err := w.bm.ethKeyStore.EnabledAddressesForChain(ctx, w.bm.chainID)
	if err != nil {
		w.bm.eng.Error("BalanceMonitor: error getting keys", err)
	}

	var wg sync.WaitGroup

	wg.Add(len(enabledAddresses))
	for _, address := range enabledAddresses {
		go func(k gethCommon.Address) {
			defer wg.Done()
			w.checkAccountBalance(ctx, k)
		}(address)
	}
	wg.Wait()
}

// Approximately ETH block time
const ethFetchTimeout = 15 * time.Second

func (w *worker) checkAccountBalance(ctx context.Context, address gethCommon.Address) {
	ctx, cancel := context.WithTimeout(ctx, ethFetchTimeout)
	defer cancel()

	bal, err := w.bm.ethClient.BalanceAt(ctx, address, nil)
	if err != nil {
		w.bm.eng.Errorw(fmt.Sprintf("BalanceMonitor: error getting balance for key %s", address.Hex()),
			"err", err,
			"address", address,
		)
	} else if bal == nil {
		w.bm.eng.Errorw(fmt.Sprintf("BalanceMonitor: error getting balance for key %s: invariant violation, bal may not be nil", address.Hex()),
			"err", err,
			"address", address,
		)
	} else {
		ethBal := assets.Eth(*bal)
		w.bm.updateBalance(ethBal, address)
	}
}

func (*NullBalanceMonitor) GetEthBalance(gethCommon.Address) *assets.Eth {
	return nil
}

// Start does noop for NullBalanceMonitor.
func (*NullBalanceMonitor) Start(context.Context) error                                { return nil }
func (*NullBalanceMonitor) Close() error                                               { return nil }
func (*NullBalanceMonitor) Ready() error                                               { return nil }
func (*NullBalanceMonitor) OnNewLongestChain(ctx context.Context, head *evmtypes.Head) {}

func ApproximateFloat64(e *assets.Eth) (float64, error) {
	ef := new(big.Float).SetInt(e.ToInt())
	weif := new(big.Float).SetInt(evmtypes.WeiPerEth)
	bf := new(big.Float).Quo(ef, weif)
	f64, _ := bf.Float64()
	if f64 == math.Inf(1) || f64 == math.Inf(-1) {
		return math.Inf(1), pkgerrors.New("assets.Eth.Float64: Could not approximate Eth value into float")
	}
	return f64, nil
}
