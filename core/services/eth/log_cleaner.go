package eth

import (
	"context"
	"strconv"
	"sync"
	"time"
)

// Public types and methods

// LogCleaner is the interface of a separate worker which deletes old records from log_consumptions.
type LogCleaner interface {
	// Start the cleaner goroutine. Only one goroutine will start per Cleaner instance.
	// Subsequent calls of Start will be no-ops.
	Start()
	// Stop will gracefully terminate the cleaner goroutine.
	// If a detele transaction is in progress it will be rolled back.
	Stop()
	// Performs a cleanup round. This method can be called without Start()ing first.
	Clean()
}

// LogCleanerConfig is a config
type LogCleanerConfig struct {
	// DeleteRecordsOlderThan is a Postgresql `interval` type string.
	// See https://www.postgresql.org/docs/9.2/datatype-datetime.html#DATATYPE-INTERVAL-INPUT
	DeleteRecordsOlderThan string
	// TimeBetweenExecutions is the approximate interval to attempt to do the cleanup.
	// The execution of the cleaner at this interval is not guaranteed!
	// When the cleaner is Start()ed, the cleaner will run
	TimeBetweenExecutions time.Duration
	// NumRecordsToRemove is the max number of records that can be removed by the cleaner in a single execution.
	// If there are more logs than this setting specifies, they will removed in the next execution.
	NumRecordsToRemove uint
}

// DefaultLogCleanerConfig wisott
var DefaultLogCleanerConfig = &LogCleanerConfig{
	DeleteRecordsOlderThan: "7 days",
	TimeBetweenExecutions:  time.Hour,
	NumRecordsToRemove:     1000,
}

// Interface implementation

type logBroadcasterCleaner struct {
	stopCtx context.Context
	stopFn  context.CancelFunc
	once    *sync.Once
	orm     ormSubset
	logger  loggerSubset
	cfg     *LogCleanerConfig
}

func NewLogCleaner(orm ormSubset, logger loggerSubset, cfg *LogCleanerConfig) LogCleaner {
	ctx, cancel := context.WithCancel(context.Background())
	return &logBroadcasterCleaner{
		stopCtx: ctx,
		stopFn:  cancel,
		once:    new(sync.Once),
		orm:     orm,
		logger:  logger,
		cfg:     cfg,
	}
}

func (lbc *logBroadcasterCleaner) Start() {
	lbc.once.Do(func() {
		go lbc.runner()
	})
}

func (lbc *logBroadcasterCleaner) Stop() {
	lbc.stopFn()
}

func (lbc *logBroadcasterCleaner) Clean() {
	lbc.logger.Infow("starting the cleanup for log_consumptions records")
	numRecords, err := lbc.orm.RemoveOldLogConsumedContext(lbc.stopCtx, lbc.cfg.DeleteRecordsOlderThan, lbc.cfg.NumRecordsToRemove)
	if err != nil {
		lbc.logger.Warnw("failed to remove a slice of old log_consumptions records", "error", err.Error())
	} else {
		lbc.logger.Infow("successfully removed old log_consumptions records", "num_records", strconv.FormatInt(numRecords, 10))
	}
}

// Helpers

func (lbc *logBroadcasterCleaner) runner() {
	lbc.Clean()
	for {
		timer := time.NewTimer(lbc.durationUntilNextRun())
		select {
		case <-timer.C:
			lbc.Clean()
		case <-lbc.stopCtx.Done():
			timer.Stop()
			return
		}
	}
}

func (lbc *logBroadcasterCleaner) durationUntilNextRun() time.Duration {
	return time.Hour
}

type ormSubset interface {
	RemoveOldLogConsumedContext(context.Context, string, uint) (int64, error)
}

type loggerSubset interface {
	Warnw(msg string, keysAndValues ...interface{})
	Infow(msg string, keysAndValues ...interface{})
}
