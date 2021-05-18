package keeper

import (
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/keeper_registry_wrapper"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/log"
	"github.com/smartcontractkit/chainlink/core/utils"
	"gorm.io/gorm"
)

// MailRoom holds the log mailboxes for all the log types that keeper cares about
type MailRoom struct {
	mbUpkeepCanceled   *utils.Mailbox
	mbSyncRegistry     *utils.Mailbox
	mbUpkeepPerformed  *utils.Mailbox
	mbUpkeepRegistered *utils.Mailbox
}

func NewRegistrySynchronizer(
	job job.Job,
	contract *keeper_registry_wrapper.KeeperRegistry,
	db *gorm.DB,
	jrm job.ORM,
	logBroadcaster log.Broadcaster,
	syncInterval time.Duration,
	minConfirmations uint64,
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
		orm:              NewORM(db),
		StartStopOnce:    utils.StartStopOnce{},
		wgDone:           sync.WaitGroup{},
	}
}

// RegistrySynchronizer conforms to the Service, Listener, and HeadRelayable interfaces
var _ job.Service = (*RegistrySynchronizer)(nil)
var _ log.Listener = (*RegistrySynchronizer)(nil)

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
	wgDone           sync.WaitGroup
	utils.StartStopOnce
}

func (rs *RegistrySynchronizer) Start() error {
	return rs.StartOnce("RegistrySynchronizer", func() error {
		rs.wgDone.Add(2)
		go rs.run()

		logListenerOpts := log.ListenerOpts{
			Contract: rs.contract,
			Logs: []generated.AbigenLog{
				keeper_registry_wrapper.KeeperRegistryKeepersUpdated{},
				keeper_registry_wrapper.KeeperRegistryConfigSet{},
				keeper_registry_wrapper.KeeperRegistryUpkeepCanceled{},
				keeper_registry_wrapper.KeeperRegistryUpkeepRegistered{},
				keeper_registry_wrapper.KeeperRegistryUpkeepPerformed{},
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
	if !rs.OkayToStop() {
		return errors.New("RegistrySynchronizer is already stopped")
	}
	close(rs.chStop)
	rs.wgDone.Wait()
	return nil
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
