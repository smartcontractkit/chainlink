package config

type ClientErrors interface {
	NonceTooLow() string
	NonceTooHigh() string
	ReplacementTransactionUnderpriced() string
	LimitReached() string
	TransactionAlreadyInMempool() string
	TerminallyUnderpriced() string
	InsufficientEth() string
	TxFeeExceedsCap() string
	L2FeeTooLow() string
	L2FeeTooHigh() string
	L2Full() string
	TransactionAlreadyMined() string
	Fatal() string
	ServiceUnavailable() string
}
