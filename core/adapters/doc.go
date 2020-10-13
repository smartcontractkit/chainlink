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
// NOTE: For security, since the URL is untrusted, HTTPGet imposes some
// restrictions on which IPs may be fetched. Local network and multicast IPs
// are disallowed by default and attempting to connect will result in an error.
//
//
// HTTPPost
//
// Sends a POST request to the specified URL and will return the response.
//  { "type": "HTTPPost", "params": {"post": "https://weiwatchers.com/api" }}
//
// NOTE: For security, since the URL is untrusted, HTTPPost imposes some
// restrictions on which IPs may be fetched. Local network and multicast IPs
// are disallowed by default and attempting to connect will result in an error.
//
// HTTPGetWithUnrestrictedNetworkAccess
//
// Identical to HTTPGet except there are no IP restrictions. Use with caution.
//
// HTTPPostWithUnrestrictedNetworkAccess
//
// Identical to HTTPPost except there are no IP restrictions. Use with caution.
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
//   WARNING: The Random apdater's output is NOT the randomness you are looking
//   WARNING: for! The node must send the output onchain for verification by the
//   WARNING: method VRFCoordinator.sol#fulfillRandomnessRequest, which will
//   WARNING: pass the actual random output back to the consuming contract.
//   WARNING: Don't use the output of this adapter in any other way, unless you
//   WARNING: thoroughly understand the cryptography in use here, and the exact
//   WARNING: security guarantees it provides. See notes in VRFCoordinator.sol
//   WARNING: for more info.
//
//   WARNING: This system guarantees that the oracle cannot independently
//   WARNING: concoct a random output to suit itself, but it does not protect
//   WARNING: against collusion between the oracle and the provider of the seed
//   WARNING: the oracle uses to generate the randomness. It also does not
//   WARNING: protect against the oracle simply refusing to respond to a
//   WARNING: randomness request, if it doesn't like the output it would be
//   WARNING: required to provide. Solutions to these limitations are planned.
//
// Here is an example of a Random task specification. For an example of a full
// jobspec using this, see ../internal/testdata/randomness_job.json.
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
// A "random" task must be initiated by a "randomnesslog" initiator which
// explicitly specifies which ethereum address the logs will be emitted from,
// such as
//
// {"initiators": [{"type": "randomnesslog","address": "0xvrfCoordinatorAddr"}]}
//
// This prevents the node from responding to potentially hostile log requests
// from other contracts, which could be crafted to prematurely reveal the random
// output if someone learns a prospective input seed prior to its use in the VRF.
//
package adapters
