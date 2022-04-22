package forwarders

import (
	"bytes"
	"context"
	"encoding/json"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	evmlogpoller "github.com/smartcontractkit/chainlink/core/chains/evm/logpoller"
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

type FwdMgr struct {
	logger    logger.Logger
	ORM       ORM
	evmClient evmclient.Client
	// fwdrsSenders is supposed to be long lived until triggered by db change
	fwdrsSenders map[common.Address][]common.Address
	// TODO(samhassan) destSenders should be an LRU capped cache
	destSenders map[common.Address][]common.Address
	config      Config
	logpoller   *evmlogpoller.LogPoller
	authRvr     authorized_receiver.AuthorizedReceiverInterface
	configSt    offchain_aggregator_wrapper.OffchainAggregatorInterface
	chStop      chan struct{}
	utils.StartStopOnce
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
	}
	return &fwdMgr
}

//Start starts the forwarder manager, init forwarder cache and listen to auth events for all forwarders
func (f *FwdMgr) Start(ctx context.Context) error {
	return f.StartOnce("EVMForwarderManager", func() error {
		fwdrs, cnt, err := f.ORM.FindForwardersByChain(utils.Big(*f.evmClient.ChainID()))
		if err != nil {
			return errors.Errorf("Error retrieving forwarders for chain %d: %s", f.evmClient.ChainID(), err)
		}

		if cnt != 0 {
			f.initForwardersCache(ctx, fwdrs)
			f.subscribeForwardersLogs(fwdrs)
			bs, _ := json.Marshal(f.fwdrsSenders)
			f.logger.Criticalf("state of the world cache %v", string(bs))
		}
		f.logpoller.Replay(ctx, 1)
		go f.runLoop()
		return nil
	})
}

func (f *FwdMgr) initForwardersCache(ctx context.Context, fwdrs []Forwarder) {
	for _, fwdr := range fwdrs {
		// Set of forwarders authorised to send to to address
		c, err := authorized_receiver.NewAuthorizedReceiverCaller(fwdr.Address, f.evmClient)
		if err != nil {
			f.logger.Criticalf("Failed to init forwarder caller: %s", err.Error())
			continue
		}
		opts := bind.CallOpts{Context: nil, Pending: false}
		senders, err := c.GetAuthorizedSenders(&opts) // This line should be cached
		if err != nil {
			f.logger.Criticalf("Error calling auth senders %s", err)
			continue
		}
		f.fwdrsSenders[fwdr.Address] = senders
	}
}
func (f *FwdMgr) subscribeForwardersLogs(fwdrs []Forwarder) {
	for _, fwdr := range fwdrs {
		f.logger.Criticalf("tracking forwarder %s in the poller", fwdr.Address.String())
		f.logpoller.MergeFilter(
			[]common.Hash{authorized_receiver.AuthorizedReceiverAuthorizedSendersChanged{}.Topic()},
			fwdr.Address)
	}
}

func (f *FwdMgr) runLoop() {
	tick := time.After(0)

	for {
		select {
		case <-tick:
			tick = time.After(utils.WithJitter(time.Duration(1 * time.Second)))
			// TODO(samhassan): We should optimise on this by tracking latest seen block within forwarder manager
			logs, err := f.logpoller.LatestLogEventSigsAddrsWithConfs(AuthTopics, f.collectAddresses(), 3)
			if err != nil {
				f.logger.Criticalf("Failed to retrieve latest log round %s", err)
				continue
			}
			if len(logs) == 0 {
				f.logger.Critical("Empty auth update round, skipping")
			}

			for _, log := range logs {
				f.handleAuthChange(log)
			}

		case <-f.chStop:
			return
		}
	}
}

func (f *FwdMgr) handleAuthChange(log evmlogpoller.Log) error {
	// there must be a better way to do this comparison that doesn't give me epilepsy looking at.
	ethLog := types.Log{
		Address: log.Address,
		Data:    log.Data,
		Topics:  topics(log),
	}

	switch {
	case bytes.Equal(log.EventSig, AuthTopics[0][:]):
		event, err := f.authRvr.ParseAuthorizedSendersChanged(ethLog)
		if err != nil {
			return errors.Errorf("Failed to parse auth change log")
		}

		if _, ok := f.destSenders[event.Raw.Address]; ok {
			f.logger.Criticalf("got auth event, updating %s senders list", event.Raw.Address)
			f.destSenders[event.Raw.Address] = event.Senders
		} else if _, ok := f.fwdrsSenders[event.Raw.Address]; ok {
			f.logger.Criticalf("got auth event, updating %s senders list", event.Raw.Address)
			f.fwdrsSenders[event.Raw.Address] = event.Senders
		}
	case bytes.Equal(log.EventSig, AuthTopics[0][:]):
		// ConfigSet event
		event, err := f.configSt.ParseConfigSet(ethLog)
		if err != nil {
			return errors.Errorf("Failed to parse auth change log")
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

// Close closes the EVMForwarderManager service.
func (f *FwdMgr) Close() error {
	return f.StopOnce("EVMForwarderManager", func() error {
		close(f.chStop)
		return nil
	})
}
