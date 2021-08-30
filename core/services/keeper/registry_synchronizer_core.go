package keeper

import (
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/keeper_registry_wrapper"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/log"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// RegistrySynchronizer conforms to the Service and Listener interfaces
var (
	_ job.Service  = (*RegistrySynchronizer)(nil)
	_ log.Listener = (*RegistrySynchronizer)(nil)
)

type RegistrySynchronizer struct {
	chStop           chan struct{}
	contract         *keeper_registry_wrapper.KeeperRegistry
	interval         time.Duration
	job              job.Job
	jrm              job.ORM
	logBroadcaster   log.Broadcaster
	mailRoom         MailRoom
	minConfirmations uint64
	orm              ORM
	logger           *logger.Logger
	wgDone           sync.WaitGroup
	utils.StartStopOnce
}

// MailRoom holds the log mailboxes for all the log types that keeper cares about
type MailRoom struct {
	mbUpkeepCanceled   *utils.Mailbox
	mbSyncRegistry     *utils.Mailbox
	mbUpkeepPerformed  *utils.Mailbox
	mbUpkeepRegistered *utils.Mailbox
}

// NewRegistrySynchronizer is the constructor of RegistrySynchronizer
func NewRegistrySynchronizer(
	job job.Job,
	contract *keeper_registry_wrapper.KeeperRegistry,
	orm ORM,
	jrm job.ORM,
	logBroadcaster log.Broadcaster,
	syncInterval time.Duration,
	minConfirmations uint64,
	logger *logger.Logger,
) *RegistrySynchronizer {
	mailRoom := MailRoom{
		mbUpkeepCanceled:   utils.NewMailbox(50),
		mbSyncRegistry:     utils.NewMailbox(1),
		mbUpkeepPerformed:  utils.NewMailbox(300),
		mbUpkeepRegistered: utils.NewMailbox(50),
	}
	return &RegistrySynchronizer{
		chStop:           make(chan struct{}),
		contract:         contract,
		interval:         syncInterval,
		job:              job,
		jrm:              jrm,
		logBroadcaster:   logBroadcaster,
		mailRoom:         mailRoom,
		minConfirmations: minConfirmations,
		orm:              orm,
		logger:           logger,
	}
}

func (rs *RegistrySynchronizer) Start() error {
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
			NumConfirmations: rs.minConfirmations,
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
