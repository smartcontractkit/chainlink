---
"chainlink": minor
---

Implemented new chain agnostic MultiNode design along with the corresponding EVM implementation. The chain agnostic components enable Multinode to be integrated with Solana and other non-EVM chains. Previously the Multinode was coupled with EVM specific actions, and was called to execute these actions direclty. With this change, the MultiNode's responsibility has been simplified to focus on RPC selection along with performing health checks. Chain specific actions will instead be executed on the RPC directly after being selected by MultiNode. The Chain Agnostic MultiNode provides improved reliability and metrics for all chain integrations using it.

These are following main components:
Node: Common component which wraps an RPC with state information, health checks, and an alive loop to handle state changes along with maintaining chain information.
RPCClient: Chain-specific RPC wrapper which implements required interface for MultiNode along with any chain-specific functionality needed.
MultiNode: Perform RPCClient selection and performs health checks on all RPCs.
TransactionSender: Chain agnostic component which broadcasts transactions to all healthy RPCs and aggregates results. A chain-specific error classifier must be implemented.

MultiNode picks the "best" RPC based on one of the configurable criteria:
- Priority defined in the config.
- Highest latest block.
- Round-robin within the same priority level (or using other configurable selection algorithms)

Benefits of Chain Agnostic MultiNode:
  Reliability: Improved RPC reliability scaleable to all chains
  Maintainability: Can apply changes across all chain integrations through the use of common code
  Extendability: Can add new health checks, RPC selection and ranking algorithms
  Integration Speed: Much faster to integrate MultiNode with new chains
  Reduced Generics: Significantly less bulky code!

#updated #changed #internal
