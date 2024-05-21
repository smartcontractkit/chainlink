# BulletproofTxManager Architecture Overview

# Diagrams

## Finite state machine

### `evm.txes.state`

`unstarted`
|
|
v
`in_progress` (only one per key)
| \
| \
v v
`fatal_error` `unconfirmed`
| ^
| |
v |
`confirmed`

### `eth_tx_attempts.state`

`in_progress`
| ^
| |
v |
`broadcast`

# Data structures

Key:

⚫️ - has never been broadcast to the network

🟠 - may or may not have been broadcast to the network

🔵 - has definitely been broadcast to the network

EB - EthBroadcaster

EC - EthConfirmer

`evm.txes` has five possible states:

- EB ⚫️ `unstarted`
- EB 🟠 `in_progress`
- EB/EC ⚫️ `fatal_error`
- EB/EC 🔵 `unconfirmed`
- EB/EC 🔵 `confirmed`

`eth_tx_attempts` has two possible states:

- EB/EC 🟠 `in_progress`
- EB/EC 🔵 `broadcast`

An attempt may have 0 or more `eth_receipts` indicating that the transaction has been mined into a block. This block may or may not exist as part of the canonical longest chain.

# Components

BulletproofTxManager is split into three components, each of which has a clearly delineated set of responsibilities.

## EthTx

Conceptually, **EthTx** defines the transaction.

**EthTx** is responsible for generating the transaction criteria and inserting the initial `unstarted` row into the `evm.txes` table.

**EthTx** guarantees that the transaction is defined with the following criteria:

- From address
- To address
- Encoded payload
- Value (eth)
- Gas limit

Only one transaction may be created per **EthTx** task.

EthTx should wait until it's transaction confirms before marking the task as completed.

## EthBroadcaster

Conceptually, **EthBroadcaster** assigns a nonce to a transaction and ensures that it is valid. It alone maintains the next usable sequence for a transaction.

**EthBroadcaster** monitors `evm.txes` for transactions that need to be broadcast, assigns nonces and ensures that at least one eth node somewhere has placed the transaction into its mempool.

It does not guarantee eventual confirmation!

A whole host of other things can subsequently go wrong such as transactions being evicted from the mempool, eth nodes crashing, netsplits between eth nodes, chain re-orgs etc. Responsibility for ensuring eventual inclusion into the longest chain falls on the shoulders of **EthConfirmer**.

**EthBroadcaster** makes the following guarantees:

- A gapless, monotonically increasing sequence of nonces for `evm.txes` (scoped to key).
- Transition of `evm.txes` from `unstarted` to either `fatal_error` or `unconfirmed`.
- If final state is `fatal_error` then the nonce is unassigned, and it is impossible that this transaction could ever be mined into a block.
- If final state is `unconfirmed` then a saved `eth_transaction_attempt` exists.
- If final state is `unconfirmed` then an eth node somewhere has accepted this transaction into its mempool at least once.

**EthConfirmer** must serialize access on a per-key basis since nonce assignment needs to be tightly controlled. Multiple keys can however be processed in parallel. Serialization is enforced with an advisory lock scoped to the key.

## EthConfirmer

Conceptually, **EthConfirmer** adjusts the gas price as necessary to get a transaction mined into a block on the longest chain.

**EthConfirmer** listens to new heads and performs four separate tasks in sequence every time we become aware of a longer chain.

### 1. Mark "broadcast before"

When we receive a block we can be sure that any currently `unconfirmed` transactions were broadcast before this block was received, so we set `broadcast_before_block_num` on all transaction attempts made since we saw the last block.

It is important to know how long a transaction has been waiting for inclusion, so we can know for how many blocks a transaction has been waiting for inclusion in order to decide if we need to bump gas.

### 2. Check for receipts

Find all `unconfirmed` transactions and ask the eth node for a receipt. If there is a receipt, we save it and move this transaction into `confirmed` state.

### 3. Bump gas if necessary

Find all `unconfirmed` transactions where all attempts have remained unconfirmed for more than `ETH_GAS_BUMP_THRESHOLD` blocks. Create a new `eth_transaction_attempt` for each, with a higher gas price.

### 4. Re-org protection

Find all transactions confirmed within the past `ETH_FINALITY_DEPTH` blocks and verify that they have at least one receipt in the current longest chain. If any do not, then rebroadcast those transactions.

**EthConfirmer** makes the following guarantees:

- All transactions will eventually be confirmed on the canonical longest chain, unless a reorg occurs that is deeper than `ETH_FINALITY_DEPTH` blocks.
- In the case that an external wallet used the nonce, we will ensure that _a_ transaction exists at this nonce up to a depth of `ETH_FINALITY_DEPTH` blocks but it most likely will not be the transaction in our database.

Note that since checking for inclusion in the longest chain can now be done cheaply, without any calls to the eth node, `ETH_FINALITY_DEPTH` can be set to something quite large without penalty (e.g. 50 or 100).

**EthBroadcaster** runs are designed to be serialized. Running it concurrently with itself probably can't get the data into an inconsistent state, but it might hit database conflicts or double-send transactions. Serialization is enforced with an advisory lock.

# Head Tracker limitations

The design of **EthConfirmer** relies on an unbroken chain of heads in our database. If there is a break in the chain of heads, our re-org protection is limited to this break.

For example if we have heads at heights:

1

2

4

Then a reorg that happened at block height 3 or above will not be detected and any transactions mined in those blocks may be left erroneously marked as confirmed.

Currently, the design of the head tracker opens us up to gaps in the head sequence. This can occur in several scenarios:

1. CL Node goes offline for more than one or two blocks
2. Eth node is behind a load balancer and gets switched out for one that has different block timing
3. Websocket connection is broken and resubscribe does not occur right away

For this reason, I propose that follow-up work should be undertaken to ensure that the head tracker has some facility for backfilling heads up to`ETH_FINALITY_DEPTH`.
