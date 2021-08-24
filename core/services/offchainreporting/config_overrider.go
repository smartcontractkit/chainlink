package offchainreporting

import (
	"context"
	"math"
	"math/big"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/utils"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting/types"
)

type ConfigOverriderImpl struct {
	utils.StartStopOnce
	logger          *logger.Logger
	flags           *ContractFlags
	contractAddress ethkey.EIP55Address

	flagsUpdateInterval      time.Duration
	lastStateChangeTimestamp time.Time
	isHibernating            bool
	DeltaCFromAddress        time.Duration

	// Start/Stop lifecycle
	ctx       context.Context
	ctxCancel context.CancelFunc
	chDone    chan struct{}

	mu sync.Mutex
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

	addressBig := contractAddress.Big()
	addressInt64 := addressBig.Mod(addressBig, big.NewInt(math.MaxInt64)).Int64()
	deltaC := 23*time.Hour + time.Duration(addressInt64%3600)*time.Second

	ctx, cancel := context.WithCancel(context.Background())
	co := ConfigOverriderImpl{
		utils.StartStopOnce{},
		logger,
		flags,
		contractAddress,
		pollInterval,
		time.Now(),
		InitialHibernationStatus,
		deltaC,
		ctx,
		cancel,
		make(chan struct{}),
		sync.Mutex{},
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
	shouldHibernate := !isFlagLowered

	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.isHibernating && shouldHibernate {
		c.logger.Infow("Setting hibernation state to 'hibernating'", "elapsedSinceLastChange", time.Since(c.lastStateChangeTimestamp))
		c.lastStateChangeTimestamp = time.Now()
	}

	if c.isHibernating && !shouldHibernate {
		c.logger.Infow("Setting hibernation state to 'not hibernating'", "elapsedSinceLastChange", time.Since(c.lastStateChangeTimestamp))
		c.lastStateChangeTimestamp = time.Now()
	}

	c.isHibernating = shouldHibernate
	return nil
}

func (c *ConfigOverriderImpl) ConfigOverride() *ocrtypes.ConfigOverride {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.isHibernating {
		c.logger.Debugw("Returning a config override")
		return c.configOverrideInstance()
	}
	c.logger.Debugw("Config override returned as nil")
	return nil
}

func (c *ConfigOverriderImpl) configOverrideInstance() *ocrtypes.ConfigOverride {
	return &ocrtypes.ConfigOverride{AlphaPPB: math.MaxUint64, DeltaC: c.DeltaCFromAddress}
}
