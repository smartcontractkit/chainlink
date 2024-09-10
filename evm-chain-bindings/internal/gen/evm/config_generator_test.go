package evm_test

import (
	"encoding/json"
	"fmt"
	"github.com/smartcontractkit/smart-contract-spec/internal/gen/evm"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSimpleChainReaderConfigGeneration(t *testing.T) {

	chainReaderConfig, err := evm.GenerateChainReaderChainWriterConfig([]string{"/Users/pablolagreca/Dev/evm-chain-bindings/contracts/ChainReaderTester.abi"})
	if err != nil {
		t.Fail()
	}

	chainReaderConfigJson, err := json.MarshalIndent(chainReaderConfig, "", "  ")

	if err != nil {
		t.Fail()
	}

	fmt.Println(string(chainReaderConfigJson))

	assert.NotNil(t, chainReaderConfig)
}
