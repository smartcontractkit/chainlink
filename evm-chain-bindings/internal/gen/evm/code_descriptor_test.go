package evm_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleParsing(t *testing.T) {
	abiJSON := `[
		{
			"constant": true,
			"inputs": [],
			"name": "name",
			"outputs": [
				{
					"name": "",
					"type": "string"
				}
			],
			"payable": false,
			"stateMutability": "view",
			"type": "function"
		},
		{
			"constant": false,
			"inputs": [
				{
					"name": "newValue",
					"type": "uint256"
				}
			],
			"name": "setValue",
			"outputs": [],
			"payable": false,
			"stateMutability": "nonpayable",
			"type": "function"
		}
	]`

	contractDetails, err := ConvertABIToCodeDetails("SimpleContract", abiJSON)

	if err != nil {
		t.Fail()
	}

	assert.NotNil(t, contractDetails)
}
