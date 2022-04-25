package forwarders

import (
	"context"
	"encoding/json"
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
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/sqlx"
)

// Config encompasses config used by fwdmgr
//go:generate mockery --recursive --name Config --output ./mocks/ --case=underscore --structname Config --filename config.go
type Config interface {
	EvmUseForwarders() bool
	LogSQL() bool
}

var AuthTopics = []common.Hash{
	authorized_receiver.AuthorizedReceiverAuthorizedSendersChanged{}.Topic(),
	offchain_aggregator_wrapper.OffchainAggregatorConfigSet{}.Topic(),
}

var ForwardABI = evmtypes.MustGetABI(authorized_forwarder.AuthorizedForwarderABI).Methods["forward"]

type FwdMgr struct {
	ORM       ORM
	config    Config
	evmClient evmclient.Client
	logger    logger.Logger
	logpoller *evmlogpoller.LogPoller

	fwdrsSenders map[common.Address][]common.Address
	// TODO(samhassan): destSenders should be an LRU capped cache
	destSenders map[common.Address][]common.Address
	latestBlock int64

	authRvr  authorized_receiver.AuthorizedReceiverInterface
	configSt offchain_aggregator_wrapper.OffchainAggregatorInterface

	ctx    context.Context
	cancel context.CancelFunc
	chStop chan struct{}
}

func NewFwdMgr(db *sqlx.DB, cfg Config, client evmclient.Client, logpoller *evmlogpoller.LogPoller, lggr logger.Logger) *FwdMgr {
	lggr = lggr.Named("FwdMgr")
	lggr.Infow("Initializing EVM forwarder manager")
	fwdMgr := FwdMgr{
		logger:       lggr,
		evmClient:    client,
		ORM:          NewORM(db, lggr, cfg),
		logpoller:    logpoller,
		config:       cfg,
		fwdrsSenders: make(map[common.Address][]common.Address),
		destSenders:  make(map[common.Address][]common.Address),
		chStop:       make(chan struct{}),
		latestBlock:  0,
	}
	return &fwdMgr
}

// Start starts the forwarder manager, init forwarder cache and listen to auth events for all forwarders
func (f *FwdMgr) Start() error {
	fwdrs, cnt, err := f.ORM.FindForwardersByChain(utils.Big(*f.evmClient.ChainID()))
	if err != nil {
		return errors.Errorf("Error retrieving forwarders for chain %d: %s", f.evmClient.ChainID(), err)
	}
	f.ctx, f.cancel = context.WithCancel(context.Background())
	if cnt != 0 {
		f.initForwardersCache(f.ctx, fwdrs)
		f.subscribeForwardersLogs(fwdrs)
		bs, _ := json.Marshal(f.fwdrsSenders)
		f.logger.Criticalf("state of the world cache %v", string(bs))
	}

	f.authRvr, err = authorized_receiver.NewAuthorizedReceiver(common.Address{}, f.evmClient)
	if err != nil {
		f.logger.Criticalf("failed to init receiver")
	}

	f.configSt, err = offchain_aggregator_wrapper.NewOffchainAggregator(common.Address{}, f.evmClient)
	if err != nil {
		f.logger.Criticalf("Failed to init receiver")
	}

	go f.runLoop()
	return nil
}

// TODO(samhassan): this should be aware of job type to decide how to fetch senders list.
func (f *FwdMgr) MaybeForwardTransaction(from common.Address, to common.Address, EncodedPayload []byte) (fwdAddr common.Address, fwdPayload []byte, err error) {

	senders, err := f.getContractSenders(to)
	if err != nil {
		return to, EncodedPayload, errors.Wrap(err, "Skipping forwarding transaction")
	}

	// TODO(samhassan): This block should be optimised to prevent folks from getting epilepsy just looking at it.
	for _, sender := range senders {
		if EOAs, ok := f.fwdrsSenders[sender]; ok {
			for _, EOA := range EOAs {
				if EOA == from {
					f.logger.Debugf("Found forwarder %s", sender.String())
					forwardedPayload, err := f.getForwardedPayload(to, EncodedPayload)
					if err != nil {
						f.logger.Criticalf("Forwarder encoding failed, this should never happen")
						continue
					}
					return sender, forwardedPayload, nil
				}
			}
		}
	}

	return to, EncodedPayload, errors.Errorf("Skipping forwarding transaction")
}

func (f *FwdMgr) getForwardedPayload(dest common.Address, origPayload []byte) ([]byte, error) {
	callArgs, err := ForwardABI.Inputs.Pack(dest, origPayload)
	if err != nil {
		return nil, err
	}

	dataBytes := append(ForwardABI.ID, callArgs...)
	return dataBytes, nil
}

func (f *FwdMgr) getContractSenders(addr common.Address) ([]common.Address, error) {
	if senders, ok := f.destSenders[addr]; ok {
		return senders, nil
	}
	senders, err := f.getAuthorizedSenders(addr)
	if err != nil {
		f.logger.Warnf("Failed to call getAuthorizedSenders on %s", addr)
		return nil, err
	}
	f.destSenders[addr] = senders
	f.subscribeSendersChangedLogs(addr)
	return senders, nil
}

func (f *FwdMgr) getAuthorizedSenders(addr common.Address) ([]common.Address, error) {
	c, err := authorized_receiver.NewAuthorizedReceiverCaller(addr, f.evmClient)
	if err != nil {
		f.logger.Errorf("Failed to init forwarder caller: %s", err.Error())
		return nil, err
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
			f.logger.Criticalf("Failed to call getAuthorizedSenders on forwarder %s: %s", fwdr, err)
			continue
		}
		f.fwdrsSenders[fwdr.Address] = senders
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

func (f *FwdMgr) runLoop() {
	tick := time.After(0)

	for {
		select {
		case <-tick:
			tick = time.After(utils.WithJitter(time.Duration(1 * time.Second)))
			addrs := f.collectAddresses()
			logs, err := f.logpoller.LatestLogEventSigsAddrsWithConfs(f.latestBlock, AuthTopics, addrs, 0)
			if err != nil {
				f.logger.Errorf("Failed to retrieve latest log round %s", err)
				continue
			}
			if len(logs) == 0 {
				f.logger.Debugf("Empty auth update round for addrs: %s, skipping", addrs)
				continue
			}

			f.logger.Infof("handling new %d auth updates", len(logs))

			for _, log := range logs {
				f.handleAuthChange(log)
			}

		case <-f.chStop:
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
		Topics:    topics(log),
		TxHash:    log.TxHash,
		BlockHash: log.BlockHash,
	}

	switch {
	case ethLog.Topics[0] == AuthTopics[0]:
		event, err := f.authRvr.ParseAuthorizedSendersChanged(ethLog)
		if err != nil {
			return errors.Errorf("Failed to parse senders change log")
		}
		if _, ok := f.destSenders[event.Raw.Address]; ok {
			f.destSenders[event.Raw.Address] = event.Senders
		} else if _, ok := f.fwdrsSenders[event.Raw.Address]; ok {
			f.fwdrsSenders[event.Raw.Address] = event.Senders
		}
	case ethLog.Topics[0] == AuthTopics[1]:
		// ConfigSet event
		event, err := f.configSt.ParseConfigSet(ethLog)
		if err != nil {
			return errors.Errorf("Failed to parse config set log")
		}
		f.destSenders[event.Raw.Address] = event.Transmitters
	}

	return nil
}

func topics(l evmlogpoller.Log) []common.Hash {
	var tps []common.Hash
	for _, topic := range l.Topics {
		tps = append(tps, common.BytesToHash(topic))
	}
	return tps
}

func (f *FwdMgr) collectAddresses() (addrs []common.Address) {
	for addr := range f.destSenders {
		addrs = append(addrs, addr)
	}

	for addr := range f.fwdrsSenders {
		addrs = append(addrs, addr)
	}

	return
}

// Stop cancels all outgoings calls and stops internal ticker loop.
func (f *FwdMgr) Stop() error {
	f.cancel()
	close(f.chStop)
	return nil
}
