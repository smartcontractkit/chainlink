package ocr

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"sync"
	"time"

	"github.com/pkg/errors"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting/types"

	"github.com/smartcontractkit/chainlink-common/pkg/services"

	"github.com/smartcontractkit/chainlink-evm/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type ConfigOverriderImpl struct {
	services.StateMachine
	logger          logger.Logger
	flags           *ContractFlags
	contractAddress types.EIP55Address

	pollTicker               utils.TickerBase
	lastStateChangeTimestamp time.Time
	isHibernating            bool
	DeltaCFromAddress        time.Duration

	// Start/Stop lifecycle
	chStop services.StopChan
	chDone chan struct{}

	mu sync.RWMutex
}

// InitialHibernationStatus - hibernation state set until the first successful update from the chain
const InitialHibernationStatus = false

type DeltaCConfig interface {
	DeltaCOverride() time.Duration
	DeltaCJitterOverride() time.Duration
}

func NewConfigOverriderImpl(
	logger logger.Logger,
	cfg DeltaCConfig,
	contractAddress types.EIP55Address,
	flags *ContractFlags,
	pollTicker utils.TickerBase,
) (*ConfigOverriderImpl, error) {
	if !flags.ContractExists() {
		return nil, errors.Errorf("OCRConfigOverrider: Flags contract instance is missing, the contract does not exist: %s. "+
			"Please create the contract or remove the OCR.TransmitterAddress configuration variable", contractAddress.Address())
	}

	addressBig := contractAddress.Big()
	jitterSeconds := int64(cfg.DeltaCJitterOverride() / time.Second)
	addressSeconds := addressBig.Mod(addressBig, big.NewInt(jitterSeconds)).Uint64()
	deltaC := cfg.DeltaCOverride() + time.Duration(addressSeconds)*time.Second

	co := ConfigOverriderImpl{
		services.StateMachine{},
		logger,
		flags,
		contractAddress,
		pollTicker,
		time.Now(),
		InitialHibernationStatus,
		deltaC,
		make(chan struct{}),
		make(chan struct{}),
		sync.RWMutex{},
	}

	return &co, nil
}

// Start starts ConfigOverriderImpl.
func (c *ConfigOverriderImpl) Start(context.Context) error {
	return c.StartOnce("OCRConfigOverrider", func() (err error) {
		if err := c.updateFlagsStatus(); err != nil {
			c.logger.Errorw("OCRConfigOverrider: Error updating hibernation status at OCR job start. Will default to not hibernating, until next successful update.", "err", err)
		}

		go c.eventLoop()
		return nil
	})
}

func (c *ConfigOverriderImpl) Close() error {
	return c.StopOnce("OCRContractTracker", func() error {
		close(c.chStop)
		<-c.chDone
		return nil
	})
}

func (c *ConfigOverriderImpl) eventLoop() {
	defer close(c.chDone)
	ctx, cancel := c.chStop.NewCtx()
	defer cancel()
	c.pollTicker.Resume()
	defer c.pollTicker.Destroy()
	for {
		select {
		case <-ctx.Done():
			return
		case <-c.pollTicker.Ticks():
			if err := c.updateFlagsStatus(); err != nil {
				c.logger.Errorw("OCRConfigOverrider: Error updating hibernation status", "err", err)
			}
		}
	}
}

func (c *ConfigOverriderImpl) updateFlagsStatus() error {
	isFlagLowered, err := c.flags.IsLowered(c.contractAddress.Address())
	if err != nil {
		return errors.Wrap(err, "Failed to check if flag is lowered")
	}
	shouldHibernate := !isFlagLowered

	c.mu.Lock()
	defer c.mu.Unlock()

	wasUpdated := (!c.isHibernating && shouldHibernate) || (c.isHibernating && !shouldHibernate)
	if wasUpdated {
		c.logger.Infow(
			fmt.Sprintf("OCRConfigOverrider: Setting hibernation state to '%v'", shouldHibernate),
			"elapsedSinceLastChange", time.Since(c.lastStateChangeTimestamp),
		)
		c.lastStateChangeTimestamp = time.Now()
		c.isHibernating = shouldHibernate
	}
	return nil
}

func (c *ConfigOverriderImpl) ConfigOverride() *ocrtypes.ConfigOverride {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.isHibernating {
		c.logger.Debugw("OCRConfigOverrider: Returning a config override")
		return c.configOverrideInstance()
	}
	c.logger.Debugw("OCRConfigOverrider: Config override returned as nil")
	return nil
}

func (c *ConfigOverriderImpl) configOverrideInstance() *ocrtypes.ConfigOverride {
	return &ocrtypes.ConfigOverride{AlphaPPB: math.MaxUint64, DeltaC: c.DeltaCFromAddress}
}
