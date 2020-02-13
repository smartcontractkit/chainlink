package adapters

// This file contains functions used to transform an EthTx input to the raw
// bytes of an ethereum transaction.

import (
	"encoding/hex"
	"fmt"

	"chainlink/core/utils"

	"github.com/tidwall/gjson"
)

// The following are options to the "format" argument of the JSON representation
// for adapters.EthTx. They control how EthTx will transform its input to an
// ethereum transaction ready to send to an on-chain contract. There are
// examples of the transformations they effect in TestEthTxAdapter_Perform (see
// the "format", "input" and "output" fields of the test table called "tests")
// and TestEVMTranscodeJSONWithFormat, which also uses "format", "input" and
// "output" fields, but does not include the function selector or initial
// argument offset, like most of the TestEthTxAdapter_Perform tests.
const (
	// FormatRawHexWithFuncSelectorAndDataPrefix is the default format behavior.
	// Its output is prefixed with the ethTx function selector, followed by the
	// data prefix
	FormatRawHexWithFuncSelectorAndDataPrefix = ""
	// FormatBytes encodes the output as bytes. I.e., the string input will be
	// cast to bytes, and passed to the on-chain contract method as a solidity
	// bytes array argument. (No conversion from hex is done; the string input
	// must be the raw bytes!)
	FormatBytes = "bytes"
	// FormatPreformatted encodes the output, assumed to be hex, as bytes and
	// passes them as arguments. Caller is responsible for all formatting for the
	// EVM. Input must be 0x-prefixed
	FormatPreformattedHexArguments = "preformattedHexArguments"
	// FormatRawHex does no formatting at all. Caller is responsible for
	// formatting the function selector and offset, in addition to any arguments
	// to be passed with the transaction. Input must be 0x-prefixed
	//
	// Note that this option isn't adressed in EVMTranscodeJSONWithFormat, because
	// eth_tx.go's getTxData short-circuits, when it encounters this.
	FormatRawHex = "rawHex"
	// FormatUint256 encodes the output as bytes containing a uint256
	FormatUint256 = "uint256"
	// FormatInt256 encodes the output as bytes containing an int256
	FormatInt256 = "int256"
	// FormatBool encodes the output as bytes containing a bool
	FormatBool = "bool"
)

// EVMTranscodeJSONWithFormat given a JSON input and a format specifier, encode the
// value for use by the EVM
func EVMTranscodeJSONWithFormat(value gjson.Result, format string) ([]byte, error) {
	switch format {
	case FormatBytes:
		return utils.EVMTranscodeBytes(value)
	case FormatPreformattedHexArguments:
		if !utils.HasHexPrefix(value.Str) {
			return nil, fmt.Errorf("%s input must be 0x-prefixed, got %s",
				FormatPreformattedHexArguments, value.Str)
		}
		return hex.DecodeString(utils.RemoveHexPrefix(value.Str))
	case FormatUint256:
		data, err := utils.EVMTranscodeUint256(value)
		if err != nil {
			return []byte{}, err
		}
		return utils.EVMEncodeBytes(data), nil

	case FormatInt256:
		data, err := utils.EVMTranscodeInt256(value)
		if err != nil {
			return []byte{}, err
		}
		return utils.EVMEncodeBytes(data), nil

	case FormatBool:
		data, err := utils.EVMTranscodeBool(value)
		if err != nil {
			return []byte{}, err
		}
		return utils.EVMEncodeBytes(data), nil

	default:
		return []byte{}, fmt.Errorf("unsupported format: %s", format)
	}
}
