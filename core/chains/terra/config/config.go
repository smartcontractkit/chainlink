package config

import (
	"net/url"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/smartcontractkit/chainlink-terra/pkg/terra/db"

	"github.com/smartcontractkit/chainlink/core/store/models"
)

type Chain struct {
	BlockRate             *models.Duration
	BlocksUntilTxTimeout  *int64
	ConfirmPollPeriod     *models.Duration
	FallbackGasPriceULuna *decimal.Decimal
	FCDURL                *models.URL
	GasLimitMultiplier    *decimal.Decimal
	MaxMsgsPerBatch       *int64
	OCR2CachePollPeriod   *models.Duration
	OCR2CacheTTL          *models.Duration
	TxMsgTimeout          *models.Duration
}

func (c *Chain) SetFromDB(cfg *db.ChainCfg) error {
	if cfg == nil {
		return nil
	}
	if cfg.BlockRate != nil {
		c.BlockRate = models.MustNewDuration(cfg.BlockRate.Duration())
	}
	if cfg.BlocksUntilTxTimeout.Valid {
		c.BlocksUntilTxTimeout = &cfg.BlocksUntilTxTimeout.Int64
	}
	if cfg.ConfirmPollPeriod != nil {
		c.ConfirmPollPeriod = models.MustNewDuration(cfg.ConfirmPollPeriod.Duration())
	}
	if cfg.FallbackGasPriceULuna.Valid {
		s := cfg.FallbackGasPriceULuna.String
		d, err := decimal.NewFromString(s)
		if err != nil {
			return errors.Wrapf(err, "invalid decimal FallbackGasPriceULuna: %s", s)
		}
		c.FallbackGasPriceULuna = &d
	}
	if cfg.FCDURL.Valid {
		s := cfg.FCDURL.String
		d, err := url.Parse(s)
		if err != nil {
			return errors.Wrapf(err, "invalid FCDURL: %s", s)
		}
		c.FCDURL = (*models.URL)(d)
	}
	if cfg.GasLimitMultiplier.Valid {
		d := decimal.NewFromFloat(cfg.GasLimitMultiplier.Float64)
		c.GasLimitMultiplier = &d
	}
	if cfg.MaxMsgsPerBatch.Valid {
		c.MaxMsgsPerBatch = &cfg.MaxMsgsPerBatch.Int64
	}
	if cfg.OCR2CachePollPeriod != nil {
		c.OCR2CachePollPeriod = models.MustNewDuration(cfg.OCR2CachePollPeriod.Duration())
	}
	if cfg.OCR2CacheTTL != nil {
		c.OCR2CacheTTL = models.MustNewDuration(cfg.OCR2CacheTTL.Duration())
	}
	if cfg.TxMsgTimeout != nil {
		c.TxMsgTimeout = models.MustNewDuration(cfg.TxMsgTimeout.Duration())
	}
	return nil
}

type Node struct {
	Name          string
	TendermintURL *models.URL
}

func (n *Node) SetFromDB(db db.Node) error {
	n.Name = db.Name
	if db.TendermintURL != "" {
		u, err := url.Parse(db.TendermintURL)
		if err != nil {
			return err
		}
		n.TendermintURL = (*models.URL)(u)
	}
	return nil
}
