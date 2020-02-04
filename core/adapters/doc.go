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
// Quotient
//
// The Quotient adapter gives the result of x / y where x is a specified value (dividend)
// and y is the input value (result).
// This can be useful for inverting outputs, e.g. if you have a USD/ETH conversion
// rate and you want to flip it to ETH/USD you can use this adapter with a dividend of 1 to get
// 1 / result.
//
// value.
//   { "type": "Quotient", "params": {"dividend": 1 }}
//
// Random
//
// Random adapter generates proofs of randomness verifiable against a public key
//
// WARNING: The Random apdater's output is NOT the randomness you are looking
// for! The node should send it to VRFCoordinator.sol#fulfillRandomnessRequest,
// for verification and to pass the actual random output back to the consuming
// contract. Don't use the output of this adapter in any other way, unless you
// thoroughly understand the cryptography in use here, and the exact security
// guarantees it provides. See the notes in VRFCoordinator.sol for more info.
//
// WARNING: This system guarantees that the oracle cannot independently concoct
// a random output to suit itself, but it does not protect against collusion
// between the oracle and the provider of the seed the oracle uses to generate
// the randomness. It also does not protect against the oracle simply refusing
// to respond to a randomness request, if it doesn't like the output it would be
// required to provide. Solutions to these limitations are planned.
//
// Here is an example of a Random task specification. For an example of a full
// jobspec using this, see ../internal/fixtures/web/randomness_job.json.
//
//  {
//    "type": "Random",
//    "params": {
//    	"publicKey":
//        "0x79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f8179800"
//    }
//  }
//
// The publicKey must be the concatenation of its hex representation of its the
// secp256k1 point's x-ordinate as a uint256, followed by 00 if the y-ordinate
// is even, or 01 if it's odd. (Note that this is NOT an RFC 5480 section 2.2
// public-key representation. DO NOT prefix with 0x02, 0x03 or 0x04.)
//
// The chainlink node must know the corresponding secret key. Such a key pair
// can be created with the `chainlink local vrf create` command, and exported to
// a keystore with `vrf export <keystore-path>`.
//
// E.g. `chainlink local vrf create -p <password-file>` will log the public key
// under the field "public id".
//
// To see the public keys which have already been imported, use the command
// `chainlink local vrf list`. See `chainlink local vrf help` for more
// key-manipulation commands.
//
// The adapter output should be passed via EthTx to VRFCoordinator.sol's method
// fulfillRandomnessRequest.
//
// EthTxABIEncode
//
// The EthTxABIEncode adapter serializes the contents of a json object as
// transaction data calling an arbitrary function of a smart contract. See
// https://solidity.readthedocs.io/en/v0.5.11/abi-spec.html#formal-specification-of-the-encoding
// for the serialization format. We currently support all types that solidity
// contracts as of solc v0.5.11 can decode, i.e. address, bool, bytes1, ...,
// bytes32, int8, ..., int256, uint8, ..., uint256, arrays (e.g. address[2]),
// bytes (variable length), string (variable length), and slices (e.g. uint256[]
// or address[2][]).
//
// The ABI of the function to be called is specified in the functionABI field,
// using the ABI JSON format used by solc and vyper. For example,
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
//     - a hexstring containing at most 20 bytes, e.g.
//       "0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef"
//   bool:
//     - true or false
//   bytes<n> (where 1 <= n <= 32):
//     - a hexstring containing exactly n bytes, e.g.
//       "0xdeadbeef" for a bytes8
//     - an array of numbers containing exactly n bytes, e.g.
//       [1,2,3,4,5,6,7,8,9,10] for a bytes10
//   int<n> (where n ∈ {8, 16, 24, ..., 256}):
//     - a decimal number, e.g. from "-128" to "127" for an int8
//     - a positive hex number, e.g. "0xdeadbeef" for an int64
//     - for n <= 48, a number, e.g. 4294967296 for an int40
//       (for larger n, json numbers aren't suitable due to them commonly
//       being interpreted as doubles which cannot represent all integers above
//       2**53)
//   uint<n> (where n ∈ {8, 16, 24, ..., 256}):
//     - a decimal number, e.g. from "0" to "255" for an uint8
//     - a positive hex number, e.g. "0xdeadbeef" for an uint64
//     - for n <= 48, a number, e.g. 4294967296 for an uint40
//       (for larger n, json numbers aren't suitable due to them commonly
//       being interpreted as doubles which cannot represent all integers above
//       2**53)
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
