package forwarders

import (
	"context"
	"errors"
	"slices"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	pkgerrors "github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/offchain_aggregator_wrapper"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/operatorforwarder/generated/authorized_forwarder"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/operatorforwarder/generated/authorized_receiver"
	evmclient "github.com/smartcontractkit/chainlink-evm/pkg/client"
	evmlogpoller "github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils/big"
)

var forwardABI = evmtypes.MustGetABI(authorized_forwarder.AuthorizedForwarderABI).Methods["forward"]
var authChangedTopic = authorized_receiver.AuthorizedReceiverAuthorizedSendersChanged{}.Topic()

type Config interface {
	FinalityDepth() uint32
}

type FwdMgr struct {
	services.Service
	eng *services.Engine

	ORM       ORM
	evmClient evmclient.Client
	cfg       Config
	logger    logger.SugaredLogger
	logpoller evmlogpoller.LogPoller

	// TODO(samhassan): sendersCache should be an LRU capped cache
	// https://smartcontract-it.atlassian.net/browse/ARCHIVE-22505
	sendersCache map[common.Address][]common.Address
	latestBlock  int64

	authRcvr    authorized_receiver.AuthorizedReceiverInterface
	offchainAgg offchain_aggregator_wrapper.OffchainAggregatorInterface

	cacheMu sync.RWMutex
}

func NewFwdMgr(ds sqlutil.DataSource, client evmclient.Client, logpoller evmlogpoller.LogPoller, lggr logger.Logger, cfg Config) *FwdMgr {
	fm := FwdMgr{
		cfg:          cfg,
		evmClient:    client,
		ORM:          NewORM(ds),
		logpoller:    logpoller,
		sendersCache: make(map[common.Address][]common.Address),
	}
	fm.Service, fm.eng = services.Config{
		Name:  "ForwarderManager",
		Start: fm.start,
	}.NewServiceEngine(lggr)
	fm.logger = logger.Sugared(fm.eng)
	return &fm
}

func (f *FwdMgr) start(ctx context.Context) error {
	chainId := f.evmClient.ConfiguredChainID()

	fwdrs, err := f.ORM.FindForwardersByChain(ctx, big.Big(*chainId))
	if err != nil {
		return pkgerrors.Wrapf(err, "Failed to retrieve forwarders for chain %d", chainId)
	}
	if len(fwdrs) != 0 {
		f.initForwardersCache(ctx, fwdrs)
		if err = f.subscribeForwardersLogs(ctx, fwdrs); err != nil {
			return err
		}
	}

	f.authRcvr, err = authorized_receiver.NewAuthorizedReceiver(common.Address{}, f.evmClient)
	if err != nil {
		return pkgerrors.Wrap(err, "Failed to init AuthorizedReceiver")
	}

	f.offchainAgg, err = offchain_aggregator_wrapper.NewOffchainAggregator(common.Address{}, f.evmClient)
	if err != nil {
		return pkgerrors.Wrap(err, "Failed to init OffchainAggregator")
	}

	f.eng.Go(f.runLoop)
	return nil
}

func FilterName(addr common.Address) string {
	return evmlogpoller.FilterName("ForwarderManager AuthorizedSendersChanged", addr.String())
}

func (f *FwdMgr) ForwarderFor(ctx context.Context, addr common.Address) (forwarder common.Address, err error) {
	// Gets forwarders for current chain.
	fwdrs, err := f.ORM.FindForwardersByChain(ctx, big.Big(*f.evmClient.ConfiguredChainID()))
	if err != nil {
		return common.Address{}, err
	}

	for _, fwdr := range fwdrs {
		eoas, err := f.getContractSenders(ctx, fwdr.Address)
		if err != nil {
			f.logger.Errorw("Failed to get forwarder senders", "forwarder", fwdr.Address, "err", err)
			continue
		}
		for _, eoa := range eoas {
			if eoa == addr {
				return fwdr.Address, nil
			}
		}
	}
	return common.Address{}, ErrForwarderForEOANotFound
}

// ErrForwarderForEOANotFound defines the error triggered when no valid forwarders were found for EOA
var ErrForwarderForEOANotFound = errors.New("cannot find forwarder for given EOA")

func (f *FwdMgr) ForwarderForOCR2Feeds(ctx context.Context, eoa, ocr2Aggregator common.Address) (forwarder common.Address, err error) {
	fwdrs, err := f.ORM.FindForwardersByChain(ctx, big.Big(*f.evmClient.ConfiguredChainID()))
	if err != nil {
		return common.Address{}, err
	}

	offchainAggregator, err := ocr2aggregator.NewOCR2Aggregator(ocr2Aggregator, f.evmClient)
	if err != nil {
		return common.Address{}, err
	}

	transmitters, err := offchainAggregator.GetTransmitters(&bind.CallOpts{Context: ctx})
	if err != nil {
		return common.Address{}, pkgerrors.Errorf("failed to get ocr2 aggregator transmitters: %s", err.Error())
	}

	for _, fwdr := range fwdrs {
		if !slices.Contains(transmitters, fwdr.Address) {
			f.logger.Criticalw("Forwarder is not set as a transmitter", "forwarder", fwdr.Address, "ocr2Aggregator", ocr2Aggregator, "err", err)
			continue
		}

		eoas, err := f.getContractSenders(ctx, fwdr.Address)
		if err != nil {
			f.logger.Errorw("Failed to get forwarder senders", "forwarder", fwdr.Address, "err", err)
			continue
		}
		for _, addr := range eoas {
			if addr == eoa {
				return fwdr.Address, nil
			}
		}
	}
	return common.Address{}, ErrForwarderForEOANotFound
}

func (f *FwdMgr) ConvertPayload(dest common.Address, origPayload []byte) ([]byte, error) {
	databytes, err := f.getForwardedPayload(dest, origPayload)
	if err != nil {
		if err != nil {
			f.logger.AssumptionViolationw("Forwarder encoding failed, this should never happen",
				"err", err, "to", dest, "payload", origPayload)
			f.eng.EmitHealthErr(err)
		}
	}
	return databytes, nil
}

func (f *FwdMgr) getForwardedPayload(dest common.Address, origPayload []byte) ([]byte, error) {
	callArgs, err := forwardABI.Inputs.Pack(dest, origPayload)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "Failed to pack forwarder payload")
	}

	dataBytes := append(forwardABI.ID, callArgs...)
	return dataBytes, nil
}

func (f *FwdMgr) getContractSenders(ctx context.Context, addr common.Address) ([]common.Address, error) {
	if senders, ok := f.getCachedSenders(addr); ok {
		return senders, nil
	}
	senders, err := f.getAuthorizedSenders(ctx, addr)
	if err != nil {
		return nil, pkgerrors.Wrapf(err, "Failed to call getAuthorizedSenders on %s", addr)
	}
	f.setCachedSenders(addr, senders)
	if err = f.subscribeSendersChangedLogs(ctx, addr); err != nil {
		return nil, err
	}
	return senders, nil
}

func (f *FwdMgr) getAuthorizedSenders(ctx context.Context, addr common.Address) ([]common.Address, error) {
	c, err := authorized_receiver.NewAuthorizedReceiverCaller(addr, f.evmClient)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "Failed to init forwarder caller")
	}
	opts := bind.CallOpts{Context: ctx, Pending: false}
	senders, err := c.GetAuthorizedSenders(&opts)
	if err != nil {
		return nil, err
	}
	return senders, nil
}

func (f *FwdMgr) initForwardersCache(ctx context.Context, fwdrs []Forwarder) {
	for _, fwdr := range fwdrs {
		senders, err := f.getAuthorizedSenders(ctx, fwdr.Address)
		if err != nil {
			f.logger.Warnw("Failed to call getAuthorizedSenders on forwarder", "forwarder", fwdr.Address, "err", err)
			continue
		}
		f.setCachedSenders(fwdr.Address, senders)
	}
}

func (f *FwdMgr) subscribeForwardersLogs(ctx context.Context, fwdrs []Forwarder) error {
	for _, fwdr := range fwdrs {
		if err := f.subscribeSendersChangedLogs(ctx, fwdr.Address); err != nil {
			return err
		}
	}
	return nil
}

func (f *FwdMgr) subscribeSendersChangedLogs(ctx context.Context, addr common.Address) error {
	if err := f.logpoller.Ready(); err != nil {
		f.logger.Warnw("Unable to subscribe to AuthorizedSendersChanged logs", "forwarder", addr, "err", err)
		return nil
	}

	err := f.logpoller.RegisterFilter(
		ctx,
		evmlogpoller.Filter{
			Name:      FilterName(addr),
			EventSigs: []common.Hash{authChangedTopic},
			Addresses: []common.Address{addr},
		})
	return err
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

func (f *FwdMgr) runLoop(ctx context.Context) {
	ticker := services.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := f.logpoller.Ready(); err != nil {
				f.logger.Warnw("Skipping log syncing", "err", err)
				continue
			}

			addrs := f.collectAddresses()
			if len(addrs) == 0 {
				f.logger.Debug("Skipping log syncing, no forwarders tracked.")
				continue
			}

			logs, err := f.logpoller.LatestLogEventSigsAddrsWithConfs(
				ctx,
				f.latestBlock,
				[]common.Hash{authChangedTopic},
				addrs,
				evmtypes.Confirmations(f.cfg.FinalityDepth()),
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

		case <-ctx.Done():
			return
		}
	}
}

func (f *FwdMgr) handleAuthChange(log evmlogpoller.Log) error {
	if f.latestBlock > log.BlockNumber {
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

	if ethLog.Topics[0] == authChangedTopic {
		event, err := f.authRcvr.ParseAuthorizedSendersChanged(ethLog)
		if err != nil {
			return pkgerrors.New("Failed to parse senders change log")
		}
		f.setCachedSenders(event.Raw.Address, event.Senders)
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
