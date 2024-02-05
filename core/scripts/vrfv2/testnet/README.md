# Using the External Subscription Owner Example

The [external subscription owner example contract](../../../../contracts/src/v0.8/tests/VRFExternalSubOwnerExample.sol)
allows its owner to request random words from VRF V2 if it is added as a
consumer for a funded VRF subscription.

This guide covers:
 1. Deploying the contract
 2. Creating, funding, checking balance, and adding a consumer to a VRF V2 
    subscription
 3. Requesting randomness from the contract

## Setup

Before starting, you will need:
1. An EVM chain endpoint URL
2. The chain ID corresponding to your chain
3. The private key of an account funded with LINK, and the chain's native token
   (to pay transaction fees)
4. [The LINK address, VRF coordinator address, and key hash](https://docs.chain.link/docs/vrf/v2/supported-networks/)
   for your chain.
5. [Go](https://go.dev/doc/install)

The endpoint URL can be a locally running node, or an externally hosted one like
[alchemy](https://www.alchemy.com/). Your chain ID will be a number
corresponding to the chain you pick. For example the Rinkeby testnet has chain
ID 4. Your private key can be exported from [MetaMask](https://metamask.zendesk.com/hc/en-us/articles/360015289632-How-to-Export-an-Account-Private-Key).

Note: Be careful with your key. When using testnets, it's best to use a separate
account that does not hold real funds.

Run the following command to set up your environment:

```shell
export ETH_URL=<YOUR ETH URL>
export ETH_CHAIN_ID=<YOUR CHAIN ID>
export ACCOUNT_KEY=<YOUR PRIVATE KEY>
export LINK=<LINK ADDRESS>
export LINK_ETH_FEED=<ADDRESS OFF LINK/ETH FEED>
export COORDINATOR=<COORDINATOR ADDRESS>
export KEY_HASH=<KEY HASH>
export ORACLE_ADDRESS=<YOUR ORACLE NODE ADDRESS>
export PUB_KEY=<YOUR UNCOMPRESSED PUBLIC KEY>
```

By default, the script automatically estimates gas limits for operations. Optionally, `ETH_GAS_LIMIT_DEFAULT` environment variable can be set to override gas limit for operations. 

Now "cd" into the VRF V2 testnet scripts directory:

```shell
cd <YOUR LOCAL CHAINLINK REPO>/core/scripts/vrfv2/testnet
```

## Deploying a full VRF Universe (BHS, Registered + Funded Coordinator, Consumer)

To deploy a full VRF environment on-chain, run:

```shell
go run . deploy-universe \
--subscription-balance=5000000000000000000 \ #5 LINK
--uncompressed-pub-key=<VRF Uncompressed Public Key> \
--vrf-primary-node-sending-keys="<sending-key1-address,sending-key2-address>" \ #used to fund the keys and for sample VRF Job Spec generation
--sending-key-funding-amount 100000000000000000 \ #0.1 ETH, fund addresses specified in vrf-primary-node-sending-keys
--batch-fulfillment-enabled false \ #only used for sample VRF Job Spec generation
--register-vrf-key-against-address=<"from this address you can perform `coordinator.oracleWithdraw` to withdraw earned funds from rand request fulfilments>
```
```shell
go run . deploy-universe \
--subscription-balance=5000000000000000000 \
--uncompressed-pub-key="0xf3706e247a7b205c8a8bd25a6e8c4650474da496151371085d45beeead27e568c1a5e8330c7fa718f8a31226efbff6632ed6f8ed470b637aa9be2b948e9dcef6" \
--batch-fulfillment-enabled false \
--register-vrf-key-against-address="0x23b5613fc04949F4A53d1cc8d6BCCD21ffc38C11"
```

## Deploying the Consumer Contract

To deploy the VRFExternalSubOwnerExample contract, run:

```shell
go run . eoa-consumer-deploy --coordinator-address=$COORDINATOR --link-address=$LINK
```

You should get the output:
```
Consumer address <YOUR CONSUMER ADDRESS> hash <YOUR TX HASH>
```

Run the command:
```shell
export CONSUMER=<YOUR CONSUMER ADDRESS>
```

## Setting up a VRF V2 Subscription

In order for your newly deployed consumer to make VRF requests, it needs to be
authorized for a funded subscription.

### Creating a Subscription

```shell
go run . eoa-create-sub --coordinator-address=$COORDINATOR
```

You should get the output:
```
Create sub TX hash <YOUR TX HASH>
```

In order to get the subscription ID created by your transaction, you should use
an online block explorer and input your transaction hash. Once the transaction
is confirmed you should see a log (on Etherscan, this is in the "Logs" tab of
the transaction details screen) with the created subscription details including
the decimal representation of your subscription ID.

Once you have found the ID, run:
```shell
export SUB_ID=<YOUR SUBSCRIPTION ID>
```

### Funding a Subscription

In order to fund your subscription with 10 LINK, run:
```shell
go run . eoa-fund-sub --coordinator-address $COORDINATOR --link-address=$LINK  --sub-id=$SUB_ID --amount=10000000000000000000 # 10e18 or 10 LINK
```

You should get the output:
```
Initial account balance: <YOUR LINK BEFORE FUNDING> <YOUR ADDRESS> Funding amount: 10000000000000000000
Funding sub 61 hash <YOUR FUNDING TX HASH>
```

### (Optional) Checking Subscription Balance

To check the LINK balance of your subscription, run:
```shell
go run . sub-balance --coordinator-address $COORDINATOR --sub-id=$SUB_ID
```

You should get the output:
```
sub id <YOUR SUB ID> balance: <YOUR SUB BALANCE>
```

### Adding a Consumer to Your Subscription

In order to authorize the consumer contract to use the new subscription, run the
command:
```shell
go run . eoa-add-sub-consumer --coordinator-address $COORDINATOR --sub-id=$SUB_ID --consumer-address=$CONSUMER
```

### Requesting Randomness

At this point, the consumer is authorized as a consumer of a funded 
subscription, and is ready to request random words.

To make a request, run:
```shell
go run . eoa-request --consumer-address=$CONSUMER --sub-id=$SUB_ID --key-hash=$KEY_HASH --num-words 1 
```

You should get the output:
```
TX hash: 0x599022228ffca10b0192e0b13bea64ff74f6dab2f0a3002b0825cbe22bd98249
```

You can put this transaction hash into a block explorer to check its progress.
Shortly after it's confirmed, usually only a few minutes, you should see a
second incoming transaction to your consumer containing the randomness
result.

## Debugging Reverted Transactions

A reverted transaction could have number of root causes, for example
insufficient funds / LINK, or incorrect contract addresses.

[Tenderly](https://dashboard.tenderly.co/explorer) can be useful for debugging
why a transaction failed. For example [this Rinkeby transaction](https://dashboard.tenderly.co/tx/rinkeby/0x71a7279033b47472ca453f7a19ccb685d0f32cdb4854a45052f1aaccd80436e9)
failed because a non-owner tried to request random words from 
[VRFExternalSubOwnerExample](../../../../contracts/src/v0.8/tests/VRFExternalSubOwnerExample.sol).

## Using the `BatchBlockhashStore` Contract

The `BatchBlockhashStore` contract acts as a proxy to the `BlockhashStore` contract, allowing callers to store
and fetch many blockhashes in a single transaction.

### Deploy a `BatchBlockhashStore` instance

```
go run . batch-bhs-deploy -bhs-address $BHS_ADDRESS
```

where `$BHS_ADDRESS` is an environment variable that points to an existing `BlockhashStore` contract. If one is not available,
you can easily deploy one using this command:

```
go run . bhs-deploy
```

### Store many blockhashes

```
go run . batch-bhs-store -batch-bhs-address $BATCH_BHS_ADDRESS -block-numbers 10298742,10298741,10298740,10298739
```

where `$BATCH_BHS_ADDRESS` points to the `BatchBlockhashStore` contract deployed above, and `-block-numbers` is a comma-separated
list of block numbers you want to store in a single transaction.

Please note that these block numbers must not be further than 256 from the latest head, otherwise the store will fail.

### Fetch many blockhashes

```
go run . batch-bhs-get -batch-bhs-address $BATCH_BHS_ADDRESS -block-numbers 10298742,10298741,10298740,10298739
```

where `$BATCH_BHS_ADDRESS` points to the `BatchBlockhashStore` contract deployed above, and `-block-numbers` is a comma-separated
list of block numbers you want to get in a single transaction.

### Store many blockhashes, possibly farther back than 256 blocks

In order to store blockhashes farther back than 256 blocks we can make use of the `storeVerifyHeader` method on the `BatchBlockhashStore`.

Here's how to use it:

```
go run . batch-bhs-storeVerify -batch-bhs-address $BATCH_BHS_ADDRESS -num-blocks 25 -start-block 10298739
```

where `$BATCH_BHS_ADDRESS` points to the `BatchBlockhashStore` contract deployed above, `-num-blocks` is the amount of blocks to store, and
`-start-block` is the block to start storing from, backwards. The block number specified by `-start-block` MUST be
in the blockhash store already, or this will not work.
 
### Batch BHS "Backwards Mode"

There may be a situation where you want to backfill a lot of blockhashes, down to a certain block number.

This is where "Backwrads Mode" comes in - you're going to need the following:

* A block number that has already been stored in the BHS. The closer it is to the target block range you want to store,
the better. You can view the most oldest "Store" transactions on the BHS contract that is still ahead of the block range you
are interested in. For example, if you want to store blocks 100 to 200, and 210 and 220 are available, specify `-start-block`
as `210`.
* A destination block number, where you want to stop storing after this one has been stored in the BHS. This number doesn't have
to be in the BHS already but must be less than the block specified for `--start-block`
* A batch size to use. This is how many stores we will attempt to do in a single transaction. A good value for this is usually 50-75
for big block ranges.
* The address of the batch BHS to use.

Example:

```
go run . batch-bhs-backwards -batch-bhs-address $BATCH_BHS_ADDRESS -start-block 25814538 -end-block 25811350 -batch-size 50
```

This script is simplistic on purpose, where we wait for the transaction to mine before proceeding with the next one. This
is to avoid issues where a transaction gets sent and not included on-chain, and subsequent calls to `storeVerifyHeader` will
fail.
