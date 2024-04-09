package config

import (
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
)

var defaultRegex = ""

func derefOrDefault(s *string) string {
	if s == nil {
		return defaultRegex
	}

	return *s
}

type clientErrorsRegexConfig struct {
	c toml.ClientErrorsRegex
}

func (c *clientErrorsRegexConfig) NonceTooLow() string {
	return derefOrDefault(c.c.NonceTooLow)
}

func (c *clientErrorsRegexConfig) NonceTooHigh() string {
	return derefOrDefault(c.c.NonceTooHigh)
}

func (c *clientErrorsRegexConfig) ReplacementTransactionUnderpriced() string {
	return derefOrDefault(c.c.ReplacementTransactionUnderpriced)
}

func (c *clientErrorsRegexConfig) LimitReached() string {
	return derefOrDefault(c.c.LimitReached)
}

func (c *clientErrorsRegexConfig) TransactionAlreadyInMempool() string {
	return derefOrDefault(c.c.TransactionAlreadyInMempool)
}

func (c *clientErrorsRegexConfig) TerminallyUnderpriced() string {
	return derefOrDefault(c.c.TerminallyUnderpriced)
}

func (c *clientErrorsRegexConfig) InsufficientEth() string {
	return derefOrDefault(c.c.InsufficientEth)
}

func (c *clientErrorsRegexConfig) TxFeeExceedsCap() string {
	return derefOrDefault(c.c.TxFeeExceedsCap)
}

func (c *clientErrorsRegexConfig) L2FeeTooLow() string {
	return derefOrDefault(c.c.L2FeeTooLow)
}

func (c *clientErrorsRegexConfig) L2FeeTooHigh() string {
	return derefOrDefault(c.c.L2FeeTooHigh)
}

func (c *clientErrorsRegexConfig) L2Full() string {
	return derefOrDefault(c.c.L2Full)
}

func (c *clientErrorsRegexConfig) TransactionAlreadyMined() string {
	return derefOrDefault(c.c.TransactionAlreadyMined)
}

func (c *clientErrorsRegexConfig) Fatal() string {
	return derefOrDefault(c.c.Fatal)
}

func (c *clientErrorsRegexConfig) ServiceUnavailable() string {
	return derefOrDefault(c.c.ServiceUnavailable)
}
