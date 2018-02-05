// Package adapters contains the core adapters used by the ChainLink node.
//
// HttpGet
//
// The HttpGet adapter is used to grab the JSON data from the given URL.
//  { "type": "HttpGet", "url": "https://some-api-example.net/api" }
//
// JsonParse
//
// The JsonParse adapter will obtain the value(s) for the given field(s).
//  { "type": "JsonParse", "path": ["someField"] }
//
// EthBytes32
//
// The EthBytes32 adapter will take the given values and format them for
// the Ethereum blockhain.
//  { "type": "EthBytes32" }
//
// EthTx
//
// The EthTx adapter will write the data to the given address and functionID.
//   {
//     "type": "EthTx",
//     "address": "0x0000000000000000000000000000000000000000",
//     "functionId": "0xffffffff"
//   }
//
package adapters
