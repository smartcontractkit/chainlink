package keeper

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/log"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// RegistrySynchronizer conforms to the Service and Listener interfaces
var (
	_ job.ServiceCtx = (*RegistrySynchronizer)(nil)
	_ log.Listener   = (*RegistrySynchronizer)(nil)
)

type RegistrySynchronizerOptions struct {
	Job                      job.Job
	RegistryWrapper          RegistryWrapper
	ORM                      ORM
	JRM                      job.ORM
	LogBroadcaster           log.Broadcaster
	MailMon                  *utils.MailboxMonitor
	SyncInterval             time.Duration
	MinIncomingConfirmations uint32
	Logger                   logger.Logger
	SyncUpkeepQueueSize      uint32
	EffectiveKeeperAddress   common.Address
}

type RegistrySynchronizer struct {
	utils.StartStopOnce
	chStop                   chan struct{}
	registryWrapper          RegistryWrapper
	interval                 time.Duration
	job                      job.Job
	jrm                      job.ORM
	logBroadcaster           log.Broadcaster
	mbLogs                   *utils.Mailbox[log.Broadcast]
	minIncomingConfirmations uint32
	effectiveKeeperAddress   common.Address
	orm                      ORM
	logger                   logger.SugaredLogger
	wgDone                   sync.WaitGroup
	syncUpkeepQueueSize      uint32 //Represents the max number of upkeeps that can be synced in parallel
	mailMon                  *utils.MailboxMonitor
}

// NewRegistrySynchronizer is the constructor of RegistrySynchronizer
func NewRegistrySynchronizer(opts RegistrySynchronizerOptions) *RegistrySynchronizer {
	return &RegistrySynchronizer{
		chStop:                   make(chan struct{}),
		registryWrapper:          opts.RegistryWrapper,
		interval:                 opts.SyncInterval,
		job:                      opts.Job,
		jrm:                      opts.JRM,
		logBroadcaster:           opts.LogBroadcaster,
		mbLogs:                   utils.NewMailbox[log.Broadcast](5_000), // Arbitrary limit, better to have excess capacity
		minIncomingConfirmations: opts.MinIncomingConfirmations,
		orm:                      opts.ORM,
		effectiveKeeperAddress:   opts.EffectiveKeeperAddress,
		logger:                   logger.Sugared(opts.Logger.Named("RegistrySynchronizer")),
		syncUpkeepQueueSize:      opts.SyncUpkeepQueueSize,
		mailMon:                  opts.MailMon,
	}
}

// Start starts RegistrySynchronizer.
func (rs *RegistrySynchronizer) Start(context.Context) error {
	return rs.StartOnce("RegistrySynchronizer", func() error {
		rs.wgDone.Add(2)
		go rs.run()

		var upkeepPerformedFilter [][]log.Topic
		upkeepPerformedFilter = nil

		logListenerOpts, err := rs.registryWrapper.GetLogListenerOpts(rs.minIncomingConfirmations, upkeepPerformedFilter)
		if err != nil {
			return errors.Wrap(err, "Unable to fetch log listener opts from wrapper")
		}
		lbUnsubscribe := rs.logBroadcaster.Register(rs, *logListenerOpts)

		go func() {
			defer rs.wgDone.Done()
			defer lbUnsubscribe()
			<-rs.chStop
		}()

		rs.mailMon.Monitor(rs.mbLogs, "RegistrySynchronizer", "Logs", fmt.Sprint(rs.job.ID))

		return nil
	})
}

func (rs *RegistrySynchronizer) Close() error {
	return rs.StopOnce("RegistrySynchronizer", func() error {
		close(rs.chStop)
		rs.wgDone.Wait()
		return rs.mbLogs.Close()
	})
}

func (rs *RegistrySynchronizer) run() {
	syncTicker := utils.NewResettableTimer()
	defer rs.wgDone.Done()
	defer syncTicker.Stop()

	rs.fullSync()

	for {
		select {
		case <-rs.chStop:
			return
		case <-syncTicker.Ticks():
			rs.fullSync()
			syncTicker.Reset(rs.interval)
		case <-rs.mbLogs.Notify():
			rs.processLogs()
		}
	}
}
