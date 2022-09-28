package keeper

import (
	"context"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/chains/evm/log"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/utils"
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
	SyncInterval             time.Duration
	MinIncomingConfirmations uint32
	Logger                   logger.Logger
	SyncUpkeepQueueSize      uint32
	EffectiveKeeperAddress   common.Address
	newTurnEnabled           bool
}

type RegistrySynchronizer struct {
	chStop                   chan struct{}
	newTurnEnabled           bool
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
	utils.StartStopOnce
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
		mbLogs:                   utils.NewMailbox[log.Broadcast](5000), // Arbitrary limit, better to have excess capacity
		minIncomingConfirmations: opts.MinIncomingConfirmations,
		orm:                      opts.ORM,
		effectiveKeeperAddress:   opts.EffectiveKeeperAddress,
		logger:                   logger.Sugared(opts.Logger.Named("RegistrySynchronizer")),
		syncUpkeepQueueSize:      opts.SyncUpkeepQueueSize,
		newTurnEnabled:           opts.newTurnEnabled,
	}
}

// Start starts RegistrySynchronizer.
func (rs *RegistrySynchronizer) Start(context.Context) error {
	return rs.StartOnce("RegistrySynchronizer", func() error {
		rs.wgDone.Add(2)
		go rs.run()

		var upkeepPerformedFilter [][]log.Topic
		upkeepPerformedFilter = nil
		if !rs.newTurnEnabled {
			upkeepPerformedFilter = [][]log.Topic{
				{},
				{},
				{
					log.Topic(rs.effectiveKeeperAddress.Hash()),
				},
			}
		}

		logListenerOpts, err := rs.registryWrapper.GetLogListenerOpts(rs.minIncomingConfirmations, upkeepPerformedFilter)
		if err != nil {
			return errors.Wrap(err, "Unable to fetch log listener opts from wrapper")
		}
		lbUnsubscribe := rs.logBroadcaster.Register(rs, *logListenerOpts)

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
