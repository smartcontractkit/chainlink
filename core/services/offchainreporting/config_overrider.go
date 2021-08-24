package offchainreporting

import (
	"context"
	"math"
	"time"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/utils"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting/types"
	"github.com/tevino/abool"
)

type ConfigOverriderImpl struct {
	utils.StartStopOnce
	logger *logger.Logger

	flags           *ContractFlags
	contractAddress ethkey.EIP55Address

	flagsUpdateInterval time.Duration
	isHibernating       *abool.AtomicBool

	// Start/Stop lifecycle
	ctx       context.Context
	ctxCancel context.CancelFunc
	chDone    chan struct{}
}

// InitialHibernationStatus - hibernation state set until the first successful update from the chain
const InitialHibernationStatus = false

func NewConfigOverriderImpl(
	logger *logger.Logger,
	contractAddress ethkey.EIP55Address,
	flags *ContractFlags,
	pollInterval time.Duration,
) (*ConfigOverriderImpl, error) {

	if !flags.ContractExists() {
		return nil, errors.Errorf("OCR: Flags contract instance is missing, the contract does not exist: %s. "+
			"Please create the contract or remove the FLAGS_CONTRACT_ADDRESS configuration variable", contractAddress.Address())
	}

	ctx, cancel := context.WithCancel(context.Background())
	co := ConfigOverriderImpl{
		utils.StartStopOnce{},
		logger,
		flags,
		contractAddress,
		pollInterval,
		abool.NewBool(InitialHibernationStatus),
		ctx,
		cancel,
		make(chan struct{}),
	}

	return &co, nil
}

func (c *ConfigOverriderImpl) Start() error {
	return c.StartOnce("OCRConfigOverrider", func() (err error) {
		if err := c.updateFlagsStatus(); err != nil {
			c.logger.Errorw("Error updating hibernation status at OCR job start. Will default to not hibernating, until next successful update.", "err", err)
		}

		go c.eventLoop()
		return nil
	})
}

func (c *ConfigOverriderImpl) Close() error {
	return c.StopOnce("OCRContractTracker", func() error {
		c.ctxCancel()
		<-c.chDone
		return nil
	})
}

func (c *ConfigOverriderImpl) eventLoop() {
	defer close(c.chDone)
	ticker := time.NewTicker(utils.WithJitter(c.flagsUpdateInterval))
	defer ticker.Stop()
	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			if err := c.updateFlagsStatus(); err != nil {
				c.logger.Errorw("Error updating hibernation status", "err", err)
			}
		}
	}
}

func (c *ConfigOverriderImpl) updateFlagsStatus() error {
	isFlagLowered, err := c.flags.IsLowered(c.contractAddress.Address())
	if err != nil {
		return err
	}

	isHibernating := !isFlagLowered
	if isHibernating {
		if c.isHibernating.SetToIf(false, true) {
			c.logger.Infow("Setting hibernation state to 'hibernating'")
		}
	}
	if !isHibernating {
		if c.isHibernating.SetToIf(true, false) {
			c.logger.Infow("Setting hibernation state to 'not hibernating'")
		}
	}
	return nil
}

func (c *ConfigOverriderImpl) ConfigOverride() *ocrtypes.ConfigOverride {
	if c.isHibernating.IsSet() {
		c.logger.Debugw("Returning config override")
		return c.configOverrideInstance()
	}
	c.logger.Debugw("Config override returned as nil")
	return nil
}

func (c *ConfigOverriderImpl) configOverrideInstance() *ocrtypes.ConfigOverride {
	return &ocrtypes.ConfigOverride{AlphaPPB: math.MaxUint64, DeltaC: 23*time.Hour + time.Duration(c.contractAddress.Big().Int64()%3600)*time.Second}
}
