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
// EthTxABIEncode
//
// The EthTxABIEncode adapter serializes the contents of a json object as transaction data
// calling an arbitrary function of a smart contract. See
// https://solidity.readthedocs.io/en/v0.5.11/abi-spec.html#formal-specification-of-the-encoding
// for the serialization format. We currently support all types that solidity contracts
// as of solc v0.5.11 can decode, i.e. address, bool, bytes*, int*, uint*, arrays (e.g. address[2]),
// bytes (variable length), string (variable length), and slices (e.g. uint256[] or address[2][]).
//
// The ABI of the function to be called is specified in the functionABI field, using the ABI JSON
// format used by solc and vyper.
// For example,
//
//   {
//     "type": "EthTxABIEncode",
//     "functionABI": {
//       "name": "example"
//       "inputs": [
//         {"name": "x", "type": "uint256"},
//         {"name": "y", "type": "bool[2][]"}
//         {"name": "z", "type": "string"}
//       ]
//     },
//     "address": "0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef"
//   }
//
// will encode a transaction to a function example(uint256 x, bool[2][] y, string z) for a contract
// at address 0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef.
//
// Upon use, the json input to an EthTxABIEncode task is expected to have a
// corresponding map of argument names to compatible data in its
// `result` field such as
//
//   {
//     "result": {
//       "x": "680564733841876926926749227234810109236",
//       "y": [[true, false], [false, false]],
//       "z": "hello world! привет мир!"
//     }
//   }
//
// which will result in a call to
// `example(uint256,bool[2][],string)` at address
// 0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef, with arguments
//
//   example(
//     680564733841876926926749227234810109236,
//     [[true, false], [false, false]],
//     "hello world! привет мир!"
//   )
//
// The result from EthTxABIEncode is the hash of the resulting transaction, if it
// was successfully transmitted, or an error if not.
//
// ABI types must be represented in JSON as follows:
//   address:
//     - a hexstring containing exactly 20 bytes, e.g.
//       "0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef"
//   bool:
//     - true or false
//   bytes<n>:
//     - a hexstring containing exactly n bytes, e.g.
//       "0xdeadbeef" for a bytes8
//     - an array of numbers containing exactly n bytes, e.g.
//       [1,2,3,4,5,6,7,8,9,10] for a bytes10
//   int<n>:
//     - a decimal number, e.g. from "-128" to "127" for an int8
//     - a positive hex number, e.g. "0xdeadbeef" for an int64
//     - for n <= 48, a number, e.g. 4294967296 for an int40
//       (for larger n, json numbers aren't suitable due to them commonly
//       being interpreted as floats who lose precision about 2**53)
//   uint<n>:
//     - a decimal number, e.g. from "0" to "255" for an uint8
//     - a positive hex number, e.g. "0xdeadbeef" for an uint64
//     - for n <= 48, a number, e.g. 4294967296 for an uint40
//       (for larger n, json numbers aren't suitable due to them commonly
//       being interpreted as floats who lose precision about 2**53)
//   arrays:
//     - a json array of the appropriate length, e.g. ["0x1", "-1"] for int8[2]
//   ==============
//   bytes:
//     - a hexstring of variable length, e.g. "0xdeadbeefc0ffee"
//     - an array of numbers of variable length, e.g.
//       [1,2,1,2,1,2]
//   string:
//     - a utf8 string, e.g. "hello world! привет мир!"
//   slice:
//     - an array of variable length, e.g. ["0x1", "-2", 3] for
//       an int128[]
//
package adapters
