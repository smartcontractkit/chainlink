package config

import (
	"fmt"
	"math/big"
)

type TestCase int

const (
	TestCaseSolanaReadAccountInfo      TestCase = iota
	TestCaseSolanaGetBalance           TestCase = iota
	TestCaseSolanaGetMultipleAccounts  TestCase = iota
	TestCaseSolanaGetProgramAccounts   TestCase = iota
	TestCaseSolanaGetBlock             TestCase = iota
	TestCaseSolanaGetSlotHeight        TestCase = iota
	TestCaseSolanaGetTransaction       TestCase = iota
	TestCaseSolanaGetSignatureStatuses TestCase = iota
	TestCaseSolanaGetFeeForMessage     TestCase = iota
	TestCaseLen
)

func (tc TestCase) String() string {
	switch tc {
	case TestCaseSolanaReadAccountInfo:
		return "SolanaReadAccountInfo"
	case TestCaseSolanaGetBalance:
		return "SolanaGetBalance"
	case TestCaseSolanaGetMultipleAccounts:
		return "SolanaGetMultipleAccounts"
	case TestCaseSolanaGetProgramAccounts:
		return "SolanaGetProgramAccounts"
	case TestCaseSolanaGetBlock:
		return "SolanaGetBlock"
	case TestCaseSolanaGetSlotHeight:
		return "SolanaGetSlotHeight"
	case TestCaseSolanaGetTransaction:
		return "SolanaGetTransaction"
	case TestCaseSolanaGetSignatureStatuses:
		return "SolanaGetSignatureStatuses"
	case TestCaseSolanaGetFeeForMessage:
		return "SolanaGetFeeForMessage"
	default:
		return fmt.Sprintf("unknown TestCase: %d", tc)
	}
}

type Config struct {
	ChainSelector   uint64
	TestCase        TestCase
	WorkflowName    string
	AccountAddress  []byte
	ExpectedBalance *big.Int
	ProgramAddress  []byte
	TxSignature     []byte
	// EncodedMessage is a base64-encoded serialised Solana message, used by GetFeeForMessage.
	EncodedMessage string
}
