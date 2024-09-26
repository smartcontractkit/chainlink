package config

import (
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
)

func derefOrDefault(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}

type clientErrorsConfig struct {
	c toml.ClientErrors
}

func (c *clientErrorsConfig) NonceTooLow() string  { return derefOrDefault(c.c.NonceTooLow) }
func (c *clientErrorsConfig) NonceTooHigh() string { return derefOrDefault(c.c.NonceTooHigh) }

func (c *clientErrorsConfig) ReplacementTransactionUnderpriced() string {
	return derefOrDefault(c.c.ReplacementTransactionUnderpriced)
}

func (c *clientErrorsConfig) LimitReached() string { return derefOrDefault(c.c.LimitReached) }

func (c *clientErrorsConfig) TransactionAlreadyInMempool() string {
	return derefOrDefault(c.c.TransactionAlreadyInMempool)
}

func (c *clientErrorsConfig) TerminallyUnderpriced() string {
	return derefOrDefault(c.c.TerminallyUnderpriced)
}

func (c *clientErrorsConfig) InsufficientEth() string { return derefOrDefault(c.c.InsufficientEth) }
func (c *clientErrorsConfig) TxFeeExceedsCap() string { return derefOrDefault(c.c.TxFeeExceedsCap) }
func (c *clientErrorsConfig) L2FeeTooLow() string     { return derefOrDefault(c.c.L2FeeTooLow) }
func (c *clientErrorsConfig) L2FeeTooHigh() string    { return derefOrDefault(c.c.L2FeeTooHigh) }
func (c *clientErrorsConfig) L2Full() string          { return derefOrDefault(c.c.L2Full) }

func (c *clientErrorsConfig) TransactionAlreadyMined() string {
	return derefOrDefault(c.c.TransactionAlreadyMined)
}

func (c *clientErrorsConfig) Fatal() string { return derefOrDefault(c.c.Fatal) }

func (c *clientErrorsConfig) ServiceUnavailable() string {
	return derefOrDefault(c.c.ServiceUnavailable)
}
func (c *clientErrorsConfig) TooManyResults() string { return derefOrDefault(c.c.TooManyResults) }
