// Package logpoller is a service for querying EVM log data.
//
// It can be thought of as a more performant and sophisticated version
// of eth_getLogs https://ethereum.org/en/developers/docs/apis/json-rpc/#eth_getlogs.
// Having a local table of relevant, continually canonical logs allows us to 2 main advantages:
//   - Have hundreds of jobs/clients querying for logs without overloading the underlying RPC provider.
//   - Do more sophisticated querying (filter by confirmations/time/log contents, efficiently join between the logs table
//     and other tables on the node, etc.)
//
// Guarantees provided by the poller:
//   - Queries always return the logs from the _current_ canonical chain (same as eth_getLogs). In particular
//     that means that querying unfinalized logs may change between queries but finalized logs remain stable.
//     The threshold between unfinalized and finalized logs is the finalityDepth parameter, chosen such that with
//     exceedingly high probability logs finalityDepth deep cannot be reorged.
//   - After calling RegisterFilter with a particular event, it will never miss logs for that event
//     despite node crashes and reorgs. The granularity of the filter is always at least one block (more when backfilling).
//   - Old logs stored in the db will only be deleted if all filters matching them have explicit retention periods set, and all
//     of them have expired.  Default retention of 0 on any matching filter guarantees permanent retention.
//   - After calling Replay(fromBlock), all blocks including that one to the latest chain tip will be polled
//     with the current filter. This can be used on first time job add to specify a start block from which you wish to capture
//     existing logs.
package logpoller
