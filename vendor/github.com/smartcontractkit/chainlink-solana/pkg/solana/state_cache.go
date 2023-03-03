package solana

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"

	"github.com/smartcontractkit/chainlink-relay/pkg/utils"

	"github.com/smartcontractkit/chainlink-solana/pkg/solana/client"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/config"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/logger"
)

var (
	configVersion uint8 = 1
)

type StateCache struct {
	utils.StartStopOnce
	// on-chain program + 2x state accounts (state + transmissions)
	StateID solana.PublicKey

	stateLock sync.RWMutex
	state     State
	stateTime time.Time

	// dependencies
	reader client.Reader
	cfg    config.Config
	lggr   logger.Logger

	// polling
	done   chan struct{}
	ctx    context.Context
	cancel context.CancelFunc
}

func NewStateCache(stateID solana.PublicKey, cfg config.Config, reader client.Reader, lggr logger.Logger) *StateCache {
	return &StateCache{
		StateID: stateID,
		reader:  reader,
		lggr:    lggr,
		cfg:     cfg,
	}
}

// Start polling
func (c *StateCache) Start() error {
	return c.StartOnce("pollState", func() error {
		c.done = make(chan struct{})
		ctx, cancel := context.WithCancel(context.Background())
		c.ctx = ctx
		c.cancel = cancel
		// We synchronously update the config on start so that
		// when OCR starts there is config available (if possible).
		// Avoids confusing "contract has not been configured" OCR errors.
		err := c.fetchState(c.ctx)
		if err != nil {
			c.lggr.Warnf("error in initial PollState.fetchState %s", err)
		}
		go c.PollState()
		return nil
	})
}

// PollState contains the state and transmissions polling implementation
func (c *StateCache) PollState() {
	defer close(c.done)
	c.lggr.Debugf("Starting state polling for state: %s", c.StateID)
	tick := time.After(0)
	for {
		select {
		case <-c.ctx.Done():
			c.lggr.Debugf("Stopping state polling for state: %s", c.StateID)
			return
		case <-tick:
			// async poll both ocr2 states
			start := time.Now()
			err := c.fetchState(c.ctx)
			if err != nil {
				c.lggr.Errorf("error in PollState.fetchState %s", err)
			}
			// Note negative duration will be immediately ready
			tick = time.After(utils.WithJitter(c.cfg.OCR2CachePollPeriod()) - time.Since(start))
		}
	}
}

// Close stops the polling
func (c *StateCache) Close() error {
	return c.StopOnce("pollState", func() error {
		c.cancel()
		<-c.done
		return nil
	})
}

// ReadState reads the latest state from memory with mutex and errors if timeout is exceeded
func (c *StateCache) ReadState() (State, error) {
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()

	var err error
	if time.Since(c.stateTime) > c.cfg.OCR2CacheTTL() {
		err = errors.New("error in ReadState: stale state data, polling is likely experiencing errors")
	}
	return c.state, err
}

func (c *StateCache) fetchState(ctx context.Context) error {

	c.lggr.Debugf("fetch state for account: %s", c.StateID.String())
	state, _, err := GetState(ctx, c.reader, c.StateID, c.cfg.Commitment())
	if err != nil {
		return err
	}

	c.lggr.Debugf("state fetched for account: %s, result (config digest): %v", c.StateID, hex.EncodeToString(state.Config.LatestConfigDigest[:]))

	// acquire lock and write to state
	c.stateLock.Lock()
	defer c.stateLock.Unlock()
	c.state = state
	c.stateTime = time.Now()
	return nil
}

func GetState(ctx context.Context, reader client.AccountReader, account solana.PublicKey, commitment rpc.CommitmentType) (State, uint64, error) {
	res, err := reader.GetAccountInfoWithOpts(ctx, account, &rpc.GetAccountInfoOpts{
		Commitment: commitment,
		Encoding:   "base64",
	})
	if err != nil {
		return State{}, 0, fmt.Errorf("failed to fetch state account at address '%s': %w", account.String(), err)
	}

	// check for nil pointers
	if res == nil || res.Value == nil || res.Value.Data == nil {
		return State{}, 0, errors.New("nil pointer returned in GetState.GetAccountInfoWithOpts")
	}

	var state State
	if err := bin.NewBinDecoder(res.Value.Data.GetBinary()).Decode(&state); err != nil {
		return State{}, 0, fmt.Errorf("failed to decode state account data: %w", err)
	}

	// validation for config version
	if configVersion != state.Version {
		return State{}, 0, fmt.Errorf("decoded config version (%d) does not match expected config version (%d)", state.Version, configVersion)
	}

	blockNum := res.RPCContext.Context.Slot
	return state, blockNum, nil
}
