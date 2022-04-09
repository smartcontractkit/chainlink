package keeper

import (
	"context"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/core/chains/evm/log"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/keeper_registry_wrapper"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// RegistrySynchronizer conforms to the Service and Listener interfaces
var (
	_ job.ServiceCtx = (*RegistrySynchronizer)(nil)
	_ log.Listener   = (*RegistrySynchronizer)(nil)
)

// MailRoom holds the log mailboxes for all the log types that keeper cares about
type MailRoom struct {
	mbUpkeepCanceled   *utils.Mailbox[log.Broadcast]
	mbSyncRegistry     *utils.Mailbox[log.Broadcast]
	mbUpkeepPerformed  *utils.Mailbox[log.Broadcast]
	mbUpkeepRegistered *utils.Mailbox[log.Broadcast]
}

type RegistrySynchronizerOptions struct {
	Job                      job.Job
	Contract                 *keeper_registry_wrapper.KeeperRegistry
	ORM                      ORM
	JRM                      job.ORM
	LogBroadcaster           log.Broadcaster
	SyncInterval             time.Duration
	MinIncomingConfirmations uint32
	Logger                   logger.Logger
	SyncUpkeepQueueSize      uint32
}

type RegistrySynchronizer struct {
	chStop                   chan struct{}
	contract                 *keeper_registry_wrapper.KeeperRegistry
	interval                 time.Duration
	job                      job.Job
	jrm                      job.ORM
	logBroadcaster           log.Broadcaster
	mailRoom                 MailRoom
	minIncomingConfirmations uint32
	orm                      ORM
	logger                   logger.SugaredLogger
	wgDone                   sync.WaitGroup
	syncUpkeepQueueSize      uint32 //Represents the max number of upkeeps that can be synced in parallel
	utils.StartStopOnce
}

// NewRegistrySynchronizer is the constructor of RegistrySynchronizer
func NewRegistrySynchronizer(opts RegistrySynchronizerOptions) *RegistrySynchronizer {
	mailRoom := MailRoom{
		mbUpkeepCanceled:   utils.NewMailbox[log.Broadcast](50),
		mbSyncRegistry:     utils.NewMailbox[log.Broadcast](1),
		mbUpkeepPerformed:  utils.NewMailbox[log.Broadcast](300),
		mbUpkeepRegistered: utils.NewMailbox[log.Broadcast](50),
	}
	return &RegistrySynchronizer{
		chStop:                   make(chan struct{}),
		contract:                 opts.Contract,
		interval:                 opts.SyncInterval,
		job:                      opts.Job,
		jrm:                      opts.JRM,
		logBroadcaster:           opts.LogBroadcaster,
		mailRoom:                 mailRoom,
		minIncomingConfirmations: opts.MinIncomingConfirmations,
		orm:                      opts.ORM,
		logger:                   logger.Sugared(opts.Logger.Named("RegistrySynchronizer")),
		syncUpkeepQueueSize:      opts.SyncUpkeepQueueSize,
	}
}

// Start starts RegistrySynchronizer.
func (rs *RegistrySynchronizer) Start(context.Context) error {
	return rs.StartOnce("RegistrySynchronizer", func() error {
		rs.wgDone.Add(2)
		go rs.run()

		logListenerOpts := log.ListenerOpts{
			Contract: rs.contract.Address(),
			ParseLog: rs.contract.ParseLog,
			LogsWithTopics: map[common.Hash][][]log.Topic{
				keeper_registry_wrapper.KeeperRegistryKeepersUpdated{}.Topic():   nil,
				keeper_registry_wrapper.KeeperRegistryConfigSet{}.Topic():        nil,
				keeper_registry_wrapper.KeeperRegistryUpkeepCanceled{}.Topic():   nil,
				keeper_registry_wrapper.KeeperRegistryUpkeepRegistered{}.Topic(): nil,
				keeper_registry_wrapper.KeeperRegistryUpkeepPerformed{}.Topic():  nil,
			},
			MinIncomingConfirmations: rs.minIncomingConfirmations,
		}
		lbUnsubscribe := rs.logBroadcaster.Register(rs, logListenerOpts)

		go func() {
			defer lbUnsubscribe()
			defer rs.wgDone.Done()
			<-rs.chStop
		}()
		return nil
	})
}

func (rs *RegistrySynchronizer) Close() error {
	return rs.StopOnce("RegistrySynchronizer", func() error {
		close(rs.chStop)
		rs.wgDone.Wait()
		return nil
	})
}

func (rs *RegistrySynchronizer) run() {
	syncTicker := time.NewTicker(rs.interval)
	logTicker := time.NewTicker(time.Second)
	defer rs.wgDone.Done()
	defer syncTicker.Stop()
	defer logTicker.Stop()

	rs.fullSync()

	for {
		select {
		case <-rs.chStop:
			return
		case <-syncTicker.C:
			rs.fullSync()
		case <-logTicker.C:
			rs.processLogs()
		}
	}
}
