// Package adapters contain the core adapters used by the Chainlink node.
//
// Bridge
//
// The Bridge adapter is used to send and receive data to and from external adapters.
// The adapter will POST to the target adapter URL with an "id" field for the TaskRunID
// and a "data" field.
// For example:
//  {"id": "b8004e2989e24e1d8e4449afad2eb480", "data": {}}
//
// Compare
//
// The Compare adapter is used to compare the previous task's result
// against a specified value. Just like an `if` statement, the compare
// adapter will save `true` or `false` in the task run's result.
//  { "type": "Compare", "params": {"operator": "eq", "value": "Hello" }}
//
// HTTPGet
//
// The HTTPGet adapter is used to grab the JSON data from the given URL.
//  { "type": "HTTPGet", "params": {"get": "https://some-api-example.net/api" }}
//
// HTTPPost
//
// Sends a POST request to the specified URL and will return the response.
//  { "type": "HTTPPost", "params": {"post": "https://weiwatchers.com/api" }}
//
// JSONParse
//
// The JSONParse adapter will obtain the value(s) for the given field(s).
//  { "type": "JSONParse", "params": {"path": ["someField"] }}
//
// EthBool
//
// The EthBool adapter will take the given values and format them for
// the Ethereum blockchain in boolean value.
//  { "type": "EthBool" }
//
// EthBytes32
//
// The EthBytes32 adapter will take the given values and format them for
// the Ethereum blockchain.
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
//     "type": "EthTx", "params": {
//       "address": "0x0000000000000000000000000000000000000000",
//       "functionSelector": "0xffffffff"
//     }
//   }
//
// Multiplier
//
// The Multiplier adapter multiplies the given input value times another specified
// value.
//   { "type": "Multiply", "params": {"times": 100 }}
//
// Random
//
// Random adapter generates a number between 0 and 2**256-1
// WARNING: The random adapter as implemented is not verifiable.
// Outputs from this adapters are not verifiable onchain as a fairly-drawn random samples.
// As a result, the oracle potentially has complete discretion to instead deliberately choose
// values with favorable onchain outcomes. Don't use it for a lottery, for instance, unless 
// you fully trust the oracle not to pick its own tickets.
// We intend to either improve it in the future, or introduce a verifiable alternative.
// For now it is provided as an alternative to making web requests for random numbers,
// which is similarly unverifiable and has additional possible points of failure.
//  { "type": "Random" }
//
package adapters
