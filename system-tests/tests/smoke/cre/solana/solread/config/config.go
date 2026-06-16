package config

import (
	"fmt"
	"math/big"
)

type TestCase int

const (
	TestCaseEVMReadAccountInfo TestCase = iota
	TestCaseLen
)

func (tc TestCase) String() string {
	switch tc {
	case TestCaseEVMReadAccountInfo:
		return "EVMReadAccountInfo"
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
}
