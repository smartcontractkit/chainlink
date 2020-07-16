package cltest

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/core/eth"
)

// fail unless the struct fields on object match those in fields.
func assertFieldsMatch(t *testing.T, fields []string,
	object interface{}, targetMethod string) {
	structReflection := reflect.ValueOf(object).Elem()
	assert.Equal(t, len(fields), structReflection.NumField(),
		"number of fields in eth.Log has changed; "+
			"please update chainlinkEthLogFromGethLog")
	type_ := structReflection.Type()
	for i, field := range fields {
		assert.Equal(t, field, type_.Field(i).Name, "fields in %s have "+
			"changed; please update chainlinkEthLogFromGethLog",
			structReflection.Type().Name, targetMethod)
	}
}

func TestChainlinkEthLogFromGethLogGetsAllFields(t *testing.T) {
	fields := []string{"Address", "Topics", "Data", "BlockNumber", "TxHash",
		"TxIndex", "BlockHash", "Index", "Removed"}
	assertFieldsMatch(t, fields, &eth.Log{}, "chainlinkEthLogFromGethLog")
}

func TestGetTxReceiptGetsAllFields(t *testing.T) {
	fields := []string{"BlockNumber", "BlockHash", "Hash", "Logs"}
	assertFieldsMatch(t, fields, &eth.TxReceipt{}, "GetTxReceipt")
}

func TestGetBlockByNumberGetsAllFields(t *testing.T) {
	fields := []string{"GasPrice"}
	assertFieldsMatch(t, fields, &eth.Transaction{}, "GetBlockByNumber")
	ofields := []string{"Number", "Transactions"}
	assertFieldsMatch(t, ofields, &eth.Block{}, "GetBlockByNumber")
}

func TestSubscribeToNewHeadsGetsAllFields(t *testing.T) {
	fields := []string{"ParentHash", "UncleHash", "Coinbase", "Root", "TxHash",
		"ReceiptHash", "Bloom", "Difficulty", "Number", "GasLimit", "GasUsed",
		"Time", "Extra", "Nonce", "GethHash", "ParityHash"}
	assertFieldsMatch(t, fields, &eth.BlockHeader{}, "SubscribeToNewHeads")
}
