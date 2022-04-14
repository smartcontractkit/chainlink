package keeper

import (
	"context"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/core/chains/evm/log"
	registry1_1 "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/keeper_registry_wrapper1_1"
	registry1_2 "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/keeper_registry_wrapper1_2"
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
	Version                  RegistryVersion
	Contract1_1              *registry1_1.KeeperRegistry
	Contract1_2              *registry1_2.KeeperRegistry
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
	version                  RegistryVersion
	contract1_1              *registry1_1.KeeperRegistry
	contract1_2              *registry1_2.KeeperRegistry
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
		version:                  opts.Version,
		contract1_1:              opts.Contract1_1,
		contract1_2:              opts.Contract1_2,
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

		var logListenerOpts log.ListenerOpts
		switch rs.version {
		case RegistryVersion_1_0, RegistryVersion_1_1:
			logListenerOpts = log.ListenerOpts{
				Contract: rs.contract1_1.Address(),
				ParseLog: rs.contract1_1.ParseLog,
				LogsWithTopics: map[common.Hash][][]log.Topic{
					registry1_1.KeeperRegistryKeepersUpdated{}.Topic():   nil,
					registry1_1.KeeperRegistryConfigSet{}.Topic():        nil,
					registry1_1.KeeperRegistryUpkeepCanceled{}.Topic():   nil,
					registry1_1.KeeperRegistryUpkeepRegistered{}.Topic(): nil,
					registry1_1.KeeperRegistryUpkeepPerformed{}.Topic():  nil,
				},
				MinIncomingConfirmations: rs.minIncomingConfirmations,
			}
		case RegistryVersion_1_2:
			// TODO (sc-36399) support all v1.2 logs
			logListenerOpts = log.ListenerOpts{
				Contract: rs.contract1_2.Address(),
				ParseLog: rs.contract1_2.ParseLog,
				LogsWithTopics: map[common.Hash][][]log.Topic{
					registry1_2.KeeperRegistryKeepersUpdated{}.Topic():   nil,
					registry1_2.KeeperRegistryConfigSet{}.Topic():        nil,
					registry1_2.KeeperRegistryUpkeepCanceled{}.Topic():   nil,
					registry1_2.KeeperRegistryUpkeepRegistered{}.Topic(): nil,
					registry1_2.KeeperRegistryUpkeepPerformed{}.Topic():  nil,
				},
				MinIncomingConfirmations: rs.minIncomingConfirmations,
			}
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
