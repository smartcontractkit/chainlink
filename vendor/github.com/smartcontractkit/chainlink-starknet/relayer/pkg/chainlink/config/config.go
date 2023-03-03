package config

import (
	"net/url"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-relay/pkg/logger"
	"github.com/smartcontractkit/chainlink-relay/pkg/utils"

	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/db"
	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/ocr2"
	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/txm"
)

var DefaultConfigSet = ConfigSet{
	OCR2CachePollPeriod: 5 * time.Second,
	OCR2CacheTTL:        time.Minute,
	RequestTimeout:      10 * time.Second,
	TxTimeout:           time.Minute,
	TxSendFrequency:     5 * time.Second,
	TxMaxBatchSize:      100,
}

type ConfigSet struct {
	OCR2CachePollPeriod time.Duration
	OCR2CacheTTL        time.Duration

	// client config
	RequestTimeout time.Duration

	// txm config
	TxTimeout       time.Duration
	TxSendFrequency time.Duration
	TxMaxBatchSize  int
}

type Config interface {
	txm.Config // txm config

	// ocr2 config
	ocr2.Config

	// client config
	RequestTimeout() time.Duration

	Update(db.ChainCfg)
}

var _ Config = (*config)(nil)

type config struct {
	defaults  ConfigSet
	dbCfg     db.ChainCfg
	dbCfgLock sync.RWMutex
	lggr      logger.Logger
}

func NewConfig(dbCfg db.ChainCfg, lggr logger.Logger) *config {
	return &config{
		defaults: DefaultConfigSet,
		dbCfg:    dbCfg,
		lggr:     lggr,
	}
}

func (c *config) Update(dbCfg db.ChainCfg) {
	c.dbCfgLock.Lock()
	c.dbCfg = dbCfg
	c.dbCfgLock.Unlock()
}

func (c *config) OCR2CachePollPeriod() time.Duration {
	c.dbCfgLock.RLock()
	ch := c.dbCfg.OCR2CachePollPeriod
	c.dbCfgLock.RUnlock()
	if ch != nil {
		return ch.Duration()
	}
	return c.defaults.OCR2CachePollPeriod
}

func (c *config) OCR2CacheTTL() time.Duration {
	c.dbCfgLock.RLock()
	ch := c.dbCfg.OCR2CacheTTL
	c.dbCfgLock.RUnlock()
	if ch != nil {
		return ch.Duration()
	}
	return c.defaults.OCR2CacheTTL
}

func (c *config) RequestTimeout() time.Duration {
	c.dbCfgLock.RLock()
	ch := c.dbCfg.RequestTimeout
	c.dbCfgLock.RUnlock()
	if ch != nil {
		return ch.Duration()
	}
	return c.defaults.RequestTimeout
}

func (c *config) TxTimeout() time.Duration {
	c.dbCfgLock.RLock()
	ch := c.dbCfg.TxTimeout
	c.dbCfgLock.RUnlock()
	if ch != nil {
		return ch.Duration()
	}
	return c.defaults.TxTimeout
}

func (c *config) TxSendFrequency() time.Duration {
	c.dbCfgLock.RLock()
	ch := c.dbCfg.TxSendFrequency
	c.dbCfgLock.RUnlock()
	if ch != nil {
		return ch.Duration()
	}
	return c.defaults.TxSendFrequency
}

func (c *config) TxMaxBatchSize() int {
	c.dbCfgLock.RLock()
	ch := c.dbCfg.TxMaxBatchSize
	c.dbCfgLock.RUnlock()
	if ch.Valid {
		return int(ch.Int64)
	}
	return c.defaults.TxMaxBatchSize
}

type Chain struct {
	OCR2CachePollPeriod *utils.Duration
	OCR2CacheTTL        *utils.Duration
	RequestTimeout      *utils.Duration
	TxTimeout           *utils.Duration
	TxSendFrequency     *utils.Duration
	TxMaxBatchSize      *int64
}

func (c *Chain) SetFromDB(cfg *db.ChainCfg) error {
	if cfg == nil {
		return nil
	}

	if cfg.OCR2CachePollPeriod != nil {
		c.OCR2CachePollPeriod = utils.MustNewDuration(cfg.OCR2CachePollPeriod.Duration())
	}
	if cfg.OCR2CacheTTL != nil {
		c.OCR2CacheTTL = utils.MustNewDuration(cfg.OCR2CacheTTL.Duration())
	}
	if cfg.RequestTimeout != nil {
		c.RequestTimeout = utils.MustNewDuration(cfg.RequestTimeout.Duration())
	}
	if cfg.TxTimeout != nil {
		c.TxTimeout = utils.MustNewDuration(cfg.TxTimeout.Duration())
	}
	if cfg.TxSendFrequency != nil {
		c.TxSendFrequency = utils.MustNewDuration(cfg.TxSendFrequency.Duration())
	}
	if cfg.TxMaxBatchSize.Valid {
		c.TxMaxBatchSize = &cfg.TxMaxBatchSize.Int64
	}

	return nil
}

func (c *Chain) SetDefaults() {
	if c.OCR2CachePollPeriod == nil {
		c.OCR2CachePollPeriod = utils.MustNewDuration(DefaultConfigSet.OCR2CachePollPeriod)
	}
	if c.OCR2CacheTTL == nil {
		c.OCR2CacheTTL = utils.MustNewDuration(DefaultConfigSet.OCR2CacheTTL)
	}
	if c.RequestTimeout == nil {
		c.RequestTimeout = utils.MustNewDuration(DefaultConfigSet.RequestTimeout)
	}
	if c.TxTimeout == nil {
		c.TxTimeout = utils.MustNewDuration(DefaultConfigSet.TxTimeout)
	}
	if c.TxSendFrequency == nil {
		c.TxSendFrequency = utils.MustNewDuration(DefaultConfigSet.TxSendFrequency)
	}
	if c.TxMaxBatchSize == nil {
		i := int64(DefaultConfigSet.TxMaxBatchSize)
		c.TxMaxBatchSize = &i
	}
}

type Node struct {
	Name *string
	URL  *utils.URL
}

func (n *Node) SetFromDB(db db.Node) error {
	if db.Name != "" {
		n.Name = &db.Name
	}
	if db.URL != "" {
		u, err := url.Parse(db.URL)
		if err != nil {
			return err
		}
		n.URL = (*utils.URL)(u)
	}
	return nil
}
