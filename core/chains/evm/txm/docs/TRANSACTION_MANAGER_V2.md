
# Transaction Manager V2

## Configs
- `EIP1559`: enables EIP-1559 mode. This means the transaction manager will create and broadcast Dynamic attempts. Set this to false to broadcast Legacy transactions.
- `BlockTime`: controls the interval of the backfill loop. This dictates how frequently the transaction manager will check for confirmed transactions, rebroadcast stuck ones, and fill any nonce gaps. Transactions are getting confirmed only during new blocks so it's best if you set this to a value close to the block time. At least one RPC call is made during each BlockTime interval so the absolute minimum should be 2s. A small jitter is applied so the timeout won't be exactly the same each time.
- `RetryBlockThreshold`: is the number of blocks to wait for a transaction stuck in the mempool before automatically rebroadcasting it with a new attempt.
- `EmptyTxLimitDefault`: sets default gas limit for empty transactions. Empty transactions are created in case there is a nonce gap or another stuck transaction in the mempool to fill a given nonce. These are empty transactions and they don't have any data or value.