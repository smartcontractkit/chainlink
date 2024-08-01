package config

import (
	"errors"
	"time"

	"github.com/gagliardetto/solana-go/rpc"

	"github.com/smartcontractkit/chainlink-common/pkg/config"
)

// Global solana defaults.
var defaultConfigSet = configSet{
	BalancePollPeriod:   5 * time.Second,        // poll period for balance monitoring
	ConfirmPollPeriod:   500 * time.Millisecond, // polling for tx confirmation
	OCR2CachePollPeriod: time.Second,            // cache polling rate
	OCR2CacheTTL:        time.Minute,            // stale cache deadline
	TxTimeout:           time.Minute,            // timeout for send tx method in client
	TxRetryTimeout:      10 * time.Second,       // duration for tx rebroadcasting to RPC node
	TxConfirmTimeout:    30 * time.Second,       // duration before discarding tx as unconfirmed
	SkipPreflight:       true,                   // to enable or disable preflight checks
	Commitment:          rpc.CommitmentConfirmed,
	MaxRetries:          new(uint), // max number of retries (default = *new(uint) = 0). when config.MaxRetries < 0, interpreted as MaxRetries = nil and rpc node will do a reasonable number of retries

	// fee estimator
	FeeEstimatorMode:        "fixed",
	ComputeUnitPriceMax:     1_000,
	ComputeUnitPriceMin:     0,
	ComputeUnitPriceDefault: 0,
	FeeBumpPeriod:           3 * time.Second, // set to 0 to disable fee bumping
	BlockHistoryPollPeriod:  5 * time.Second,
}

//go:generate mockery --name Config --output ./mocks/ --case=underscore --filename config.go
type Config interface {
	BalancePollPeriod() time.Duration
	ConfirmPollPeriod() time.Duration
	OCR2CachePollPeriod() time.Duration
	OCR2CacheTTL() time.Duration
	TxTimeout() time.Duration
	TxRetryTimeout() time.Duration
	TxConfirmTimeout() time.Duration
	SkipPreflight() bool
	Commitment() rpc.CommitmentType
	MaxRetries() *uint

	// fee estimator
	FeeEstimatorMode() string
	ComputeUnitPriceMax() uint64
	ComputeUnitPriceMin() uint64
	ComputeUnitPriceDefault() uint64
	FeeBumpPeriod() time.Duration
	BlockHistoryPollPeriod() time.Duration
}

// opt: remove
type configSet struct {
	BalancePollPeriod   time.Duration
	ConfirmPollPeriod   time.Duration
	OCR2CachePollPeriod time.Duration
	OCR2CacheTTL        time.Duration
	TxTimeout           time.Duration
	TxRetryTimeout      time.Duration
	TxConfirmTimeout    time.Duration
	SkipPreflight       bool
	Commitment          rpc.CommitmentType
	MaxRetries          *uint

	FeeEstimatorMode        string
	ComputeUnitPriceMax     uint64
	ComputeUnitPriceMin     uint64
	ComputeUnitPriceDefault uint64
	FeeBumpPeriod           time.Duration
	BlockHistoryPollPeriod  time.Duration
}

type Chain struct {
	BalancePollPeriod       *config.Duration
	ConfirmPollPeriod       *config.Duration
	OCR2CachePollPeriod     *config.Duration
	OCR2CacheTTL            *config.Duration
	TxTimeout               *config.Duration
	TxRetryTimeout          *config.Duration
	TxConfirmTimeout        *config.Duration
	SkipPreflight           *bool
	Commitment              *string
	MaxRetries              *int64
	FeeEstimatorMode        *string
	ComputeUnitPriceMax     *uint64
	ComputeUnitPriceMin     *uint64
	ComputeUnitPriceDefault *uint64
	FeeBumpPeriod           *config.Duration
	BlockHistoryPollPeriod  *config.Duration
}

func (c *Chain) SetDefaults() {
	if c.BalancePollPeriod == nil {
		c.BalancePollPeriod = config.MustNewDuration(defaultConfigSet.BalancePollPeriod)
	}
	if c.ConfirmPollPeriod == nil {
		c.ConfirmPollPeriod = config.MustNewDuration(defaultConfigSet.ConfirmPollPeriod)
	}
	if c.OCR2CachePollPeriod == nil {
		c.OCR2CachePollPeriod = config.MustNewDuration(defaultConfigSet.OCR2CachePollPeriod)
	}
	if c.OCR2CacheTTL == nil {
		c.OCR2CacheTTL = config.MustNewDuration(defaultConfigSet.OCR2CacheTTL)
	}
	if c.TxTimeout == nil {
		c.TxTimeout = config.MustNewDuration(defaultConfigSet.TxTimeout)
	}
	if c.TxRetryTimeout == nil {
		c.TxRetryTimeout = config.MustNewDuration(defaultConfigSet.TxRetryTimeout)
	}
	if c.TxConfirmTimeout == nil {
		c.TxConfirmTimeout = config.MustNewDuration(defaultConfigSet.TxConfirmTimeout)
	}
	if c.SkipPreflight == nil {
		c.SkipPreflight = &defaultConfigSet.SkipPreflight
	}
	if c.Commitment == nil {
		c.Commitment = (*string)(&defaultConfigSet.Commitment)
	}
	if c.MaxRetries == nil && defaultConfigSet.MaxRetries != nil {
		i := int64(*defaultConfigSet.MaxRetries)
		c.MaxRetries = &i
	}
	if c.FeeEstimatorMode == nil {
		c.FeeEstimatorMode = &defaultConfigSet.FeeEstimatorMode
	}
	if c.ComputeUnitPriceMax == nil {
		c.ComputeUnitPriceMax = &defaultConfigSet.ComputeUnitPriceMax
	}
	if c.ComputeUnitPriceMin == nil {
		c.ComputeUnitPriceMin = &defaultConfigSet.ComputeUnitPriceMin
	}
	if c.ComputeUnitPriceDefault == nil {
		c.ComputeUnitPriceDefault = &defaultConfigSet.ComputeUnitPriceDefault
	}
	if c.FeeBumpPeriod == nil {
		c.FeeBumpPeriod = config.MustNewDuration(defaultConfigSet.FeeBumpPeriod)
	}
	if c.BlockHistoryPollPeriod == nil {
		c.BlockHistoryPollPeriod = config.MustNewDuration(defaultConfigSet.BlockHistoryPollPeriod)
	}
}

type Node struct {
	Name *string
	URL  *config.URL
}

func (n *Node) ValidateConfig() (err error) {
	if n.Name == nil {
		err = errors.Join(err, config.ErrMissing{Name: "Name", Msg: "required for all nodes"})
	} else if *n.Name == "" {
		err = errors.Join(err, config.ErrEmpty{Name: "Name", Msg: "required for all nodes"})
	}
	if n.URL == nil {
		err = errors.Join(err, config.ErrMissing{Name: "URL", Msg: "required for all nodes"})
	}
	return
}
