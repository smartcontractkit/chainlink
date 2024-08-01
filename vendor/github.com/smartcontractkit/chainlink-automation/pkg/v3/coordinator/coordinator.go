package coordinator

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	common "github.com/smartcontractkit/chainlink-common/pkg/types/automation"

	internalutil "github.com/smartcontractkit/chainlink-automation/internal/util"
	"github.com/smartcontractkit/chainlink-automation/pkg/util"
	"github.com/smartcontractkit/chainlink-automation/pkg/v3/config"
	"github.com/smartcontractkit/chainlink-automation/pkg/v3/types"
)

const (
	cadence           = time.Second
	defaultCacheClean = time.Duration(30) * time.Second
)

type coordinator struct {
	closer internalutil.Closer
	logger *log.Logger

	eventsProvider   types.TransmitEventProvider
	upkeepTypeGetter types.UpkeepTypeGetter

	cache   *util.Cache[record]
	visited *util.Cache[bool]

	minimumConfirmations int
	performLockoutWindow time.Duration
}

var _ types.Coordinator = (*coordinator)(nil)

type record struct {
	checkBlockNumber      common.BlockNumber
	isTransmissionPending bool // false = transmitted
	transmitType          common.TransmitEventType
	transmitBlockNumber   common.BlockNumber
}

func NewCoordinator(transmitEventProvider types.TransmitEventProvider, upkeepTypeGetter types.UpkeepTypeGetter, conf config.OffchainConfig, logger *log.Logger) *coordinator {
	performLockoutWindow := time.Duration(conf.PerformLockoutWindow) * time.Millisecond
	return &coordinator{
		logger:               logger,
		eventsProvider:       transmitEventProvider,
		upkeepTypeGetter:     upkeepTypeGetter,
		cache:                util.NewCache[record](performLockoutWindow),
		visited:              util.NewCache[bool](performLockoutWindow),
		minimumConfirmations: conf.MinConfirmations,
		performLockoutWindow: performLockoutWindow,
	}
}

func (c *coordinator) Accept(reportedUpkeep common.ReportedUpkeep) bool {
	if v, ok := c.cache.Get(reportedUpkeep.WorkID); !ok {
		c.cache.Set(reportedUpkeep.WorkID, record{
			checkBlockNumber:      reportedUpkeep.Trigger.BlockNumber,
			isTransmissionPending: true,
		}, util.DefaultCacheExpiration)
		return true
	} else if v.checkBlockNumber < reportedUpkeep.Trigger.BlockNumber {
		c.cache.Set(reportedUpkeep.WorkID, record{
			checkBlockNumber:      reportedUpkeep.Trigger.BlockNumber,
			isTransmissionPending: true,
		}, util.DefaultCacheExpiration)
		return true
	}
	// We are already waiting on a higher checkBlockNumber so no need to accept this report
	return false
}

func (c *coordinator) ShouldTransmit(reportedUpkeep common.ReportedUpkeep) bool {
	if v, ok := c.cache.Get(reportedUpkeep.WorkID); !ok {
		// We never saw this report, so don't try to transmit
		// Can happen in edge cases when plugin restarts after shouldAccept was called
		return false
	} else if reportedUpkeep.Trigger.BlockNumber < v.checkBlockNumber {
		// We already accepted a report for a higher checkBlockNumber, so don't try to transmit
		return false
	} else if reportedUpkeep.Trigger.BlockNumber == v.checkBlockNumber {
		return v.isTransmissionPending
	} else {
		// We never saw this report for such a high block number, so don't try to transmit
		// Can happen in edge cases when plugin restarts after shouldAccept was called
		return false
	}
}

func (c *coordinator) PreProcess(_ context.Context, payloads []common.UpkeepPayload) ([]common.UpkeepPayload, error) {
	res := make([]common.UpkeepPayload, 0)
	for _, payload := range payloads {
		if c.ShouldProcess(payload.WorkID, payload.UpkeepID, payload.Trigger) {
			res = append(res, payload)
		}
	}
	return res, nil
}

func (c *coordinator) FilterResults(results []common.CheckResult) ([]common.CheckResult, error) {
	res := make([]common.CheckResult, 0)
	for _, result := range results {
		if c.ShouldProcess(result.WorkID, result.UpkeepID, result.Trigger) {
			res = append(res, result)
		}
	}
	return res, nil
}

func (c *coordinator) FilterProposals(proposals []common.CoordinatedBlockProposal) ([]common.CoordinatedBlockProposal, error) {
	res := make([]common.CoordinatedBlockProposal, 0)
	for _, proposal := range proposals {
		if v, ok := c.cache.Get(proposal.WorkID); ok {
			if v.isTransmissionPending {
				// This workID has a pending transmit, should not process it
				continue
			} else if c.upkeepTypeGetter(proposal.UpkeepID) == types.LogTrigger && v.transmitType == common.PerformEvent {
				// For log triggers if workID was performed then skip
				// However for conditional triggers, allow proposals to be made for newer check block numbers
				continue
			}
		}
		res = append(res, proposal)
	}
	return res, nil
}

func (c *coordinator) ShouldProcess(workID string, upkeepID common.UpkeepIdentifier, trigger common.Trigger) bool {
	if v, ok := c.cache.Get(workID); ok {
		if v.isTransmissionPending {
			// This workID has a pending transmit, should not process it
			return false
		} else {
			switch c.upkeepTypeGetter(upkeepID) {
			case types.LogTrigger:
				switch v.transmitType {
				case common.PerformEvent:
					// For log triggers, a particular workID should only ever be performed once
					return false
				default:
					// There was an attempt to perform this workID, but it failed, so should be processed again
					return true
				}
			case types.ConditionTrigger:
				switch v.transmitType {
				case common.PerformEvent:
					// For conditionals, a particular workID should only be checked after or on its last perform block
					return trigger.BlockNumber >= v.transmitBlockNumber
				default:
					// There was an attempt to check this workID, but it failed, so should be processed again
					return true
				}
			}
		}
	}
	// If we have never seen this workID before, then we should process it
	return true
}

func (c *coordinator) checkEvents(ctx context.Context) error {
	events, err := c.eventsProvider.GetLatestEvents(ctx)
	if err != nil {
		return err
	}

	skipped := 0
	for _, event := range events {
		if event.Confirmations < int64(c.minimumConfirmations) {
			skipped++
			continue
		}

		// ensure we don't process the same event twice
		visitedID := c.visitedID(event)
		_, ok := c.visited.Get(visitedID)
		if ok {
			continue
		}

		v, ok := c.cache.Get(event.WorkID)
		if !ok {
			c.logger.Printf("Ignoring event in transaction %s of type %d for upkeepID %s, workID %s as it was not found in cache", hex.EncodeToString(event.TransactionHash[:]), event.Type, event.UpkeepID.String(), event.WorkID)
			continue
		}
		c.visited.Set(visitedID, true, c.performLockoutWindow)
		r := record{
			isTransmissionPending: false,
			transmitType:          event.Type,
			transmitBlockNumber:   event.TransmitBlock,
		}
		if event.CheckBlock == v.checkBlockNumber {
			c.logger.Printf("Got event in transaction %s of type %d for upkeepID %s, workID %s and check block %v", hex.EncodeToString(event.TransactionHash[:]), event.Type, event.UpkeepID.String(), event.WorkID, event.CheckBlock)
			r.checkBlockNumber = v.checkBlockNumber
			c.cache.Set(event.WorkID, r, util.DefaultCacheExpiration)
		} else if event.CheckBlock > v.checkBlockNumber {
			c.logger.Printf("Got event in transaction %s of type %d for upkeepID %s, workID %s from newer report (block %v) while waiting for (block %v)", hex.EncodeToString(event.TransactionHash[:]), event.Type, event.UpkeepID.String(), event.WorkID, event.CheckBlock, v.checkBlockNumber)
			r.checkBlockNumber = event.CheckBlock
			c.cache.Set(event.WorkID, r, util.DefaultCacheExpiration)
		}
		// otherwise this is an old event, ignore it
	}
	c.logger.Printf("Skipped %d events as confirmations are less than minimum confirmations (%d)", skipped, c.minimumConfirmations)

	return nil
}

func (c *coordinator) run(ctx context.Context) {
	timer := time.NewTimer(cadence)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			startTime := time.Now()

			if err := c.checkEvents(context.Background()); err != nil {
				c.logger.Printf("failed to check for transmit events: %s", err)
			}

			// attempt to adhere to a cadence of at least every second
			// a slow DB will cause the cadence to increase. these cases are logged
			diff := time.Since(startTime)
			if diff > cadence {
				c.logger.Printf("check transmit events took %dms to complete; expected cadence is %dms; check database indexes and other performance improvements", diff/time.Millisecond, cadence/time.Millisecond)
				// start again immediately
				timer.Reset(time.Microsecond)
			} else {
				// wait the difference between the cadence and the time taken
				timer.Reset(cadence - diff)
			}
		case <-ctx.Done():
			return
		}
	}
}

// Start starts all subprocesses
func (c *coordinator) Start(pctx context.Context) error {
	ctx, cancel := context.WithCancel(pctx)
	defer cancel()

	if !c.closer.Store(cancel) {
		return fmt.Errorf("process already running")
	}

	go c.cache.Start(defaultCacheClean)
	go c.visited.Start(defaultCacheClean)

	c.run(ctx)

	return nil
}

// Close terminates all subprocesses
func (c *coordinator) Close() error {
	if !c.closer.Close() {
		return fmt.Errorf("process not running")
	}

	c.cache.Stop()
	c.visited.Stop()

	return nil
}

func (c *coordinator) visitedID(e common.TransmitEvent) string {
	return fmt.Sprintf("%s_%x_%d", e.WorkID, e.TransactionHash, e.TransmitBlock)
}
