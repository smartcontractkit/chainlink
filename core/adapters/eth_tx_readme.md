# EthTx Adapter

I'm doing an audit of the EthTx Adapter for "transactional safety", i.e. does
it create records via the ORM in an atomic way, such that if the node were to
crash in the middle of the EthTx's operation, and were to be restarted we'd
maintain the following properties:

  1. We wouldn't send more than one successful ETH transaction
  2. The EthTx adapter would eventually be able to complete if all success
     conditions were met
  3. The database is not left in an invalid state
  4. No dangling records are created

## Sequence diagrams

Before the Perform method on the EthTx adapter, we see the following:

 JobRunner                  Store  Adapter

    |------------------------>|
    |<- []JobRun -------------|
    |                                 *
    |---Perform---------------------->|
    |                                 |

Initial run, are we connected to the chain? If not...

 Adapter         TxManager

    |- Connected() ->|
    |<------- false  |
    x

After this, job is marked as "pending_connection".

If the chain is connected:

 Adapter             TxManager

    |-Connected()->|
    |<--------true |
    |
    |- PendingConfirmations()
    |<- true
    |
    |- CreateTxWithGas ->|

If the chain is connected, and status is not "pending_confirmations":

 Adapter                TxManager                          Chain           Store

    |-------- Connected() ->|
    |<- true ---------------|
    |
    |- PendingConfirmations()
    |<- false
    |
    |-- CreateTxWithGas() ->|
    |                       |---------------------------------- CreateTx() ->|
    |                       |<- error ---------------------------------------|
    |                       |                                                |
    |                       |------------------------------ AddTxAttempt() ->|
    |                       |<- error ---------------------------------------|
    |                       |                                                |
    |                       |----- eth_sendRawTransaction  ->|               |
    |                       |<- (hash, error) ---------------|               |
    |<- (tx, error) --------|                                                |
    |
    x

" " but the status is "pending_confirmations":

 Adapter                TxManager                          Chain           Store

    |- BumpGasUntilSafe() ->|
    |                       |----------- eth_blockNumber() ->|
    |                       |<- blockId ---------------------|
    |                       |                                |
    |                       |--------------------------- First(TxAttempt) ->|
    |                       |<--TxAttempt ----------------------------------|
    |                       |                                |              |
    |                       |----------------------------------- FindTx() ->|
    |                       |<- Tx -----------------------------------------|
    |                       |                                |              |
    |                       |-----------------------------TxAttemptsFor() ->|
    |                       |<- TxAttempts ---------------------------------|
    |                       |                                |              |
    |                       |- eth_getTransactionReceipt() ->|              |
    |                       |<- receipt ---------------------|              |
    |                       |                                |              |
    |                       |------------------------------- markTxSafe() ->|
    |                       |<- error --------------------------------------|
    |<- (recept, error) ----|                                |              |
    |
    x

After any of these potential workflows, the JobRunner commits the job result,
status and any receipts to the database:

 JobRunner                  Store  Adapter

    |---Perform---------------------->|
    |<- (input, error) ---------------|
    |
    |--------- SaveJobrun() ->|
    |<- error ----------------|
    |
    x

## Potential bugs

### eth_sendRawTransaction is repeatedly invoked

The status of eth_sendRawTransaction is saved to the adapter's `input` record,
after a transaction is submitted. As you can see in the above sequence
diagrams, this data is not saved until after the adapter completes. The `input`
record is checked every time `CreateTxWithGas` is called, to ensure no existing
transaction has been attempted.

### Dangling TX + Attempts

When the TxManager is connected, and there are no pending confirmations, the
EthTx adapter can create a Transaction. This transaction is committed to the
database individually.

If after this point the node crashes, the task can be executed again creating a
dangling transaction. The same goes for the transaction that immediately
follows.

If the ETH node was down for a long time, this could result in many dangling
records.

#### Solutions

  * Check for existing transactions before creating.
  * Create TX and its attempts in a transaction
  * Have a foreign key relationship between a TX and its job

### createTxWithNonceReload is recursive

The nonceReloadLimit of 1 limits this from doing any damage, though it might be
clearer to write this as a loop.

## Domain Diagrams for Tx related records

  [ Tx ] ---- 1:M [ TxAttempts ]
