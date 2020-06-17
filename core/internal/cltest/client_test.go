package cltest

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

// fail unless the struct fields on object match those in fields.
func assertFieldsMatch(t *testing.T, fields []string,
	object interface{}, targetMethod string) {
	structReflection := reflect.ValueOf(object).Elem()
	assert.Equal(t, len(fields), structReflection.NumField(),
		"number of fields in models.Log has changed; "+
			"please update %s", targetMethod)
	type_ := structReflection.Type()
	for i, field := range fields {
		assert.Equal(t, field, type_.Field(i).Name, "fields in %s have "+
			"changed; please update %s",
			structReflection.Type().Name(), targetMethod)
	}
}

func TestChainlinkEthLogFromGethLogGetsAllFields(t *testing.T) {
	fields := []string{"Address", "Topics", "Data", "BlockNumber", "TxHash",
		"TxIndex", "BlockHash", "Index", "Removed"}
	assertFieldsMatch(t, fields, &models.Log{}, "chainlinkEthLogFromGethLog")
}

func TestGetTxReceiptGetsAllFields(t *testing.T) {
	fields := []string{"BlockNumber", "BlockHash", "Hash", "Logs"}
	assertFieldsMatch(t, fields, &models.TxReceipt{}, "GetTxReceipt")
}

func TestGetBlockByNumberGetsAllFields(t *testing.T) {
	fields := []string{"GasPrice"}
	assertFieldsMatch(t, fields, &models.Transaction{}, "GetBlockByNumber")
	ofields := []string{"Number", "Transactions"}
	assertFieldsMatch(t, ofields, &models.Block{}, "GetBlockByNumber")
}

func TestSubscribeToNewHeadsGetsAllFields(t *testing.T) {
	fields := []string{"ParentHash", "UncleHash", "Coinbase", "Root", "TxHash",
		"ReceiptHash", "Bloom", "Difficulty", "Number", "GasLimit", "GasUsed",
		"Time", "Extra", "MixDigest", "Nonce"}
	assertFieldsMatch(t, fields, &gethTypes.Header{}, "SubscribeToNewHeads")
}
