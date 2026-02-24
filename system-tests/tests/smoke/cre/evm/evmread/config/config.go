package config

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
)

type TestCase int

const (
	TestCaseEVMReadBalance TestCase = iota
	TestCaseEVMReadHeaders
	TestCaseEVMReadEvents
	TestCaseEVMReadCallContract
	TestCaseEVMReadTransactionReceipt
	TestCaseEVMReadBTx
	TestCaseEVMEstimateGas
	TestCaseEVMReadNotFoundTx
	TestCaseLen
)

func (tc TestCase) String() string {
	switch tc {
	case TestCaseEVMReadBalance:
		return "EVMReadBalance"
	case TestCaseEVMReadHeaders:
		return "EVMReadHeaders"
	case TestCaseEVMReadEvents:
		return "EVMReadEvents"
	case TestCaseEVMReadCallContract:
		return "EVMReadCallContract"
	case TestCaseEVMReadTransactionReceipt:
		return "EVMReadTransactionReceipt"
	case TestCaseEVMReadBTx:
		return "EVMReadBTx"
	case TestCaseEVMEstimateGas:
		return "EVMEstimateGas"
	case TestCaseEVMReadNotFoundTx:
		return "EVMReadNotFoundTx"
	default:
		return fmt.Sprintf("unknown TestCase: %d", tc)
	}
}

type Config struct {
	ChainSelector    uint64
	TestCase         TestCase
	WorkflowName     string
	ContractAddress  []byte
	AccountAddress   []byte
	ExpectedBalance  *big.Int
	TxHash           []byte
	ExpectedReceipt  *types.Receipt
	ExpectedBinaryTx []byte
}
