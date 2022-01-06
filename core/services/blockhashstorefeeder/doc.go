// Package blockhashstorefeeder provides the Blockhash Store Feeder job type,
// which listens to on-chain "request" events that don't have an associated "response"
// event, and after a configurable block delay, stores the blockhash of the block the
// request was included in into the BlockhashStore contract's storage.
package blockhashstorefeeder
