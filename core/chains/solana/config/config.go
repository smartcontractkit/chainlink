package config

import (
	"net/url"

	"github.com/smartcontractkit/chainlink-solana/pkg/solana/db"

	"github.com/smartcontractkit/chainlink/core/store/models"
)

type Chain struct {
	BalancePollPeriod   *models.Duration
	ConfirmPollPeriod   *models.Duration
	OCR2CachePollPeriod *models.Duration
	OCR2CacheTTL        *models.Duration
	TxTimeout           *models.Duration
	TxRetryTimeout      *models.Duration
	TxConfirmTimeout    *models.Duration
	SkipPreflight       *bool
	Commitment          *string
	MaxRetries          *int64
}

func (c *Chain) SetFromDB(cfg *db.ChainCfg) error {
	if cfg == nil {
		return nil
	}

	if cfg.BalancePollPeriod != nil {
		c.BalancePollPeriod = models.MustNewDuration(cfg.BalancePollPeriod.Duration())
	}
	if cfg.ConfirmPollPeriod != nil {
		c.ConfirmPollPeriod = models.MustNewDuration(cfg.ConfirmPollPeriod.Duration())
	}
	if cfg.OCR2CachePollPeriod != nil {
		c.OCR2CachePollPeriod = models.MustNewDuration(cfg.OCR2CachePollPeriod.Duration())
	}
	if cfg.OCR2CacheTTL != nil {
		c.OCR2CacheTTL = models.MustNewDuration(cfg.OCR2CacheTTL.Duration())
	}
	if cfg.TxTimeout != nil {
		c.TxTimeout = models.MustNewDuration(cfg.TxTimeout.Duration())
	}
	if cfg.TxRetryTimeout != nil {
		c.TxRetryTimeout = models.MustNewDuration(cfg.TxRetryTimeout.Duration())
	}
	if cfg.TxConfirmTimeout != nil {
		c.TxConfirmTimeout = models.MustNewDuration(cfg.TxConfirmTimeout.Duration())
	}
	if cfg.SkipPreflight.Valid {
		c.SkipPreflight = &cfg.SkipPreflight.Bool
	}
	if cfg.Commitment.Valid {
		c.Commitment = &cfg.Commitment.String
	}
	if cfg.MaxRetries.Valid {
		c.MaxRetries = &cfg.MaxRetries.Int64
	}
	return nil
}

type Node struct {
	Name string
	URL  *models.URL
}

func (n *Node) SetFromDB(db db.Node) error {
	n.Name = db.Name
	if db.SolanaURL != "" {
		u, err := url.Parse(db.SolanaURL)
		if err != nil {
			return err
		}
		n.URL = (*models.URL)(u)
	}
	return nil
}
