package forwarders

import (
	"context"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	evmlogpoller "github.com/smartcontractkit/chainlink/core/chains/evm/logpoller"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/authorized_forwarder"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/authorized_receiver"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/offchain_aggregator_wrapper"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/sqlx"
)

var forwardABI = evmtypes.MustGetABI(authorized_forwarder.AuthorizedForwarderABI).Methods["forward"]

type FwdMgr struct {
	utils.StartStopOnce
	ORM       ORM
	evmClient evmclient.Client
	logger    logger.SugaredLogger
	logpoller *evmlogpoller.LogPoller

	// TODO(samhassan): sendersCache should be an LRU capped cache
	// https://app.shortcut.com/chainlinklabs/story/37884/forwarder-manager-uses-lru-for-caching-dest-addresses
	sendersCache map[common.Address][]common.Address
	latestBlock  int64

	authRcvr    authorized_receiver.AuthorizedReceiverInterface
	offchainAgg offchain_aggregator_wrapper.OffchainAggregatorInterface

	ctx    context.Context
	cancel context.CancelFunc

	cacheMu sync.RWMutex
	chStop  chan struct{}
	wg      sync.WaitGroup
}

func NewFwdMgr(db *sqlx.DB, client evmclient.Client, logpoller *evmlogpoller.LogPoller, l logger.Logger, cfg pg.LogConfig) *FwdMgr {
	lggr := logger.Sugared(l.Named("EVMForwarderManager"))
	fwdMgr := FwdMgr{
		logger:       lggr,
		evmClient:    client,
		ORM:          NewORM(db, lggr, cfg),
		logpoller:    logpoller,
		sendersCache: make(map[common.Address][]common.Address),
		chStop:       make(chan struct{}),
		cacheMu:      sync.RWMutex{},
		wg:           sync.WaitGroup{},
		latestBlock:  0,
	}
	return &fwdMgr
}

// Start starts Forwarder Manager.
func (f *FwdMgr) Start() error {
	return f.StartOnce("EVMForwarderManager", func() error {
		f.logger.Debug("Initializing EVM forwarder manager")

		fwdrs, err := f.ORM.FindForwardersByChain(utils.Big(*f.evmClient.ChainID()))
		if err != nil {
			return errors.Wrapf(err, "Failed to retrieve forwarders for chain %d", f.evmClient.ChainID())
		}
		f.ctx, f.cancel = context.WithCancel(context.Background())
		if len(fwdrs) != 0 {
			f.initForwardersCache(f.ctx, fwdrs)
			f.subscribeForwardersLogs(fwdrs)
		}

		f.authRcvr, err = authorized_receiver.NewAuthorizedReceiver(common.Address{}, f.evmClient)
		if err != nil {
			return errors.Wrap(err, "Failed to init AuthorizedReceiver")
		}

		f.offchainAgg, err = offchain_aggregator_wrapper.NewOffchainAggregator(common.Address{}, f.evmClient)
		if err != nil {
			return errors.Wrap(err, "Failed to init OffchainAggregator")
		}

		f.wg.Add(1)
		go f.runLoop()
		return nil
	})
}

// TODO(samhassan): this should be aware of job type to decide how to fetch senders list.
// 	This is necessary to support ocr1 jobs.
// 	https://app.shortcut.com/chainlinklabs/story/15448/ocr1-feeds-jobs-should-detect-if-they-are-configured-to-use-a-forwarder-contract
func (f *FwdMgr) MaybeForwardTransaction(from common.Address, to common.Address, encodedPayload []byte) (fwdAddr common.Address, fwdPayload []byte, err error) {

	senders, err := f.getContractSenders(to)
	if err != nil {
		return to, encodedPayload, errors.Wrap(err, "Skipping forwarding transaction")
	}

	// Gets current forwarders that are in `to` senders
	fwdrs, err := f.ORM.FindForwardersInListByChain(utils.Big(*f.evmClient.ChainID()), senders)
	if err != nil {
		return to, encodedPayload, errors.Wrap(err, "Skipping forwarding transaction")
	}

	for _, fwdr := range fwdrs {
		eoas, err := f.getContractSenders(fwdr.Address)
		if err != nil {
			f.logger.Errorw("Failed to get forwarder senders", "err", err)
			continue
		}
		for _, eoa := range eoas {
			if eoa != from {
				continue
			}
			forwardedPayload, err := f.getForwardedPayload(to, encodedPayload)
			if err != nil {
				f.logger.AssumptionViolationw("Forwarder encoding failed, this should never happen",
					"err", err, "to", to, "payload", encodedPayload)
				continue
			}
			return fwdr.Address, forwardedPayload, nil
		}
	}

	return to, encodedPayload, errors.New("Skipping forwarding transaction")
}

func (f *FwdMgr) getForwardedPayload(dest common.Address, origPayload []byte) ([]byte, error) {
	callArgs, err := forwardABI.Inputs.Pack(dest, origPayload)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to pack forwarder payload")
	}

	dataBytes := append(forwardABI.ID, callArgs...)
	return dataBytes, nil
}

func (f *FwdMgr) getContractSenders(addr common.Address) ([]common.Address, error) {
	if senders, ok := f.getCachedSenders(addr); ok {
		return senders, nil
	}
	senders, err := f.getAuthorizedSenders(addr)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to call getAuthorizedSenders on %s", addr)
	}
	f.setCachedSenders(addr, senders)
	f.subscribeSendersChangedLogs(addr)
	return senders, nil
}

func (f *FwdMgr) getAuthorizedSenders(addr common.Address) ([]common.Address, error) {
	c, err := authorized_receiver.NewAuthorizedReceiverCaller(addr, f.evmClient)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to init forwarder caller")
	}
	opts := bind.CallOpts{Context: f.ctx, Pending: false}
	senders, err := c.GetAuthorizedSenders(&opts)
	if err != nil {
		return nil, err
	}
	return senders, nil
}

func (f *FwdMgr) initForwardersCache(ctx context.Context, fwdrs []Forwarder) {
	for _, fwdr := range fwdrs {
		senders, err := f.getAuthorizedSenders(fwdr.Address)
		if err != nil {
			f.logger.Warnw("Failed to call getAuthorizedSenders on forwarder", fwdr, "err", err)
			continue
		}
		f.setCachedSenders(fwdr.Address, senders)

	}
}
func (f *FwdMgr) subscribeForwardersLogs(fwdrs []Forwarder) {
	for _, fwdr := range fwdrs {
		f.subscribeSendersChangedLogs(fwdr.Address)
	}
}

func (f *FwdMgr) subscribeSendersChangedLogs(addr common.Address) {
	f.logpoller.MergeFilter(
		[]common.Hash{authorized_receiver.AuthorizedReceiverAuthorizedSendersChanged{}.Topic()},
		addr)
}

func (f *FwdMgr) setCachedSenders(addr common.Address, senders []common.Address) {
	f.cacheMu.Lock()
	defer f.cacheMu.Unlock()
	f.sendersCache[addr] = senders
}

func (f *FwdMgr) getCachedSenders(addr common.Address) ([]common.Address, bool) {
	f.cacheMu.RLock()
	defer f.cacheMu.RUnlock()
	addrs, ok := f.sendersCache[addr]
	return addrs, ok
}

func (f *FwdMgr) runLoop() {
	defer f.wg.Done()
	tick := time.After(0)

	for ; ; tick = time.After(utils.WithJitter(time.Duration(time.Minute))) {
		select {
		case <-tick:
			addrs := f.collectAddresses()
			logs, err := f.logpoller.LatestLogEventSigsAddrs(
				f.latestBlock,
				[]common.Hash{
					authorized_receiver.AuthorizedReceiverAuthorizedSendersChanged{}.Topic(),
					offchain_aggregator_wrapper.OffchainAggregatorConfigSet{}.Topic(),
				},
				addrs,
			)
			if err != nil {
				f.logger.Errorw("Failed to retrieve latest log round", "err", err)
				continue
			}
			if len(logs) == 0 {
				f.logger.Debugf("Empty auth update round for addrs: %s, skipping", addrs)
				continue
			}
			f.logger.Debugf("Handling new %d auth updates", len(logs))
			for _, log := range logs {
				if err = f.handleAuthChange(log); err != nil {
					f.logger.Warnw("Error handling auth change", "TxHash", log.TxHash, "err", err)
				}
			}

		case <-f.chStop:
			return
		}
	}
}

func (f *FwdMgr) handleAuthChange(log evmlogpoller.Log) error {
	if f.latestBlock >= log.BlockNumber {
		return nil
	}

	f.latestBlock = log.BlockNumber

	ethLog := types.Log{
		Address:   log.Address,
		Data:      log.Data,
		Topics:    log.GetTopics(),
		TxHash:    log.TxHash,
		BlockHash: log.BlockHash,
	}

	switch {
	case ethLog.Topics[0] == authorized_receiver.AuthorizedReceiverAuthorizedSendersChanged{}.Topic():
		event, err := f.authRcvr.ParseAuthorizedSendersChanged(ethLog)
		if err != nil {
			return errors.New("Failed to parse senders change log")
		}
		f.setCachedSenders(event.Raw.Address, event.Senders)
	case ethLog.Topics[0] == offchain_aggregator_wrapper.OffchainAggregatorConfigSet{}.Topic():
		// ConfigSet event
		event, err := f.offchainAgg.ParseConfigSet(ethLog)
		if err != nil {
			return errors.New("Failed to parse config set log")
		}
		f.setCachedSenders(event.Raw.Address, event.Transmitters)
	}

	return nil
}

func (f *FwdMgr) collectAddresses() (addrs []common.Address) {
	f.cacheMu.RLock()
	defer f.cacheMu.RUnlock()
	for addr := range f.sendersCache {
		addrs = append(addrs, addr)
	}
	return
}

// Stop cancels all outgoings calls and stops internal ticker loop.
func (f *FwdMgr) Stop() error {
	return f.StopOnce("EVMForwarderManager", func() (err error) {
		f.cancel()
		close(f.chStop)
		f.wg.Wait()
		return nil
	})
}
