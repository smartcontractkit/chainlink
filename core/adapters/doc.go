// Package adapters contain the core adapters used by the Chainlink node.
//
// HTTPGet
//
// The HTTPGet adapter is used to grab the JSON data from the given URL.
//  { "type": "HTTPGet", "url": "https://some-api-example.net/api" }
//
// HTTPPost
//
// Sends a POST request to the specified URL and will return the response.
//  { "type": "HTTPPost", "url": "https://weiwatchers.com/api" }
//
// JSONParse
//
// The JSONParse adapter will obtain the value(s) for the given field(s).
//  { "type": "JSONParse", "path": ["someField"] }
//
// EthBool
//
// The EthBool adapter will take the given values and format them for
// the Ethereum blockhain in boolean value.
//  { "type": "EthBool" }
//
// EthBytes32
//
// The EthBytes32 adapter will take the given values and format them for
// the Ethereum blockhain.
//  { "type": "EthBytes32" }
//
// EthInt256
//
// The EthInt256 adapter will take a given signed 256 bit integer and format
// it to hex for the Ethereum blockchain.
//   { "type": "EthInt256" }
//
// EthUint256
//
// The EthUint256 adapter will take a given 256 bit integer and format it
// in hex for the Ethereum blockchain.
//  { "type": "EthUint256" }
//
// EthTx
//
// The EthTx adapter will write the data to the given address and functionSelector.
//   {
//     "type": "EthTx",
//     "address": "0x0000000000000000000000000000000000000000",
//     "functionSelector": "0xffffffff"
//   }
//
// Multiplier
//
// The Multiplier adapter multiplies the given input value times another specified
// value.
//   { "type": "Multiply", "times": 100 }
//
// Bridge
//
// The Bridge adapter is used to send and receive data to and from external adapters.
// The adapter will POST to the target adapter URL with an "id" field for the TaskRunID
// and a "data" field.
// For example:
//  {"id":"b8004e2989e24e1d8e4449afad2eb480","data":{}}
//
// Random
//
// The Random adapter generates a cryptographically secure random number in the interval
// specified by the start and end parameters that default to 0 if not specified.
// For example:
//  {"start":"500","end":"1500"}
//  {end":"100"}
//  {start:"-100"}
//
package adapters
