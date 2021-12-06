# Using the Ownerless Consumer Example

The [ownerless consumer example contract](../../../contracts/src/v0.8/tests/VRFOwnerlessConsumerExample.sol)
allows anyone to request randomness from VRF V1 without needing to deploy their
own consuming contract. It does not hold any ETH or LINK; a caller must send it
LINK and spend that LINK on a randomness request within the same transaction.

This guide covers requesting randomness and optionally deploying the contract.

## Setup

Before starting, you will need:
 1. An EVM chain endpoint URL
 2. The chain ID corresponding to your chain
 3. The private key of an account funded with LINK, and the chain's native token
    (to pay transaction fees)
 4. [The LINK address, VRF coordinator address, and key hash](https://docs.chain.link/docs/vrf-contracts/) 
    for your chain
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
export COORDINATOR=<COORDINATOR ADDRESS>
export KEY_HASH=<KEY HASH>
```

Now "cd" into the VRF V1 scripts directory:

```shell
cd <YOUR LOCAL CHAINLINK REPO>/core/scripts/vrfv1
```

## Getting a Consumer

Since this contract is ownerless, you can use an existing instance instead of
deploying your own. To use an existing instance, copy the command corresponding
to the chain you want to use below, otherwise go to the 
[deployment](#deploying-a-new-consumer) section.

Once you have chosen or deployed a consumer, run:
```shell
export CONSUMER=<YOUR CONSUMER ADDRESS>
```

### Existing Consumers

#### Testnets

##### Ethereum Rinkeby Testnet

```0x1b7D5F1bD3054474cC043207aA1e7f8C152d263F```

#### BSC Testnet

```0x640F2D8fd734cb53a6938CeC4CfC0543BbcC0348```

#### Polygon Mumbai Testnet

```0x640F2D8fd734cb53a6938CeC4CfC0543BbcC0348```

### Deploying a New Consumer

To deploy the contract, run:
```shell
go run main.go ownerless-consumer-deploy --coordinator-address=$COORDINATOR --link-address=$LINK
```

You should see output like:
```
Ownerless Consumer: <YOUR CONSUMER ADDRESS> TX Hash: <YOUR TX HASH>
```

## Requesting Randomness

Since the ownerless consumer does not hold LINK funds, it can only request
randomness through a transferAndCall from the 
[LINK contract](../../../contracts/src/v0.4/LinkToken.sol). The transaction has
the following steps:
1. An externally owned account (controlled by your private key) initiates a
   transferAndCall on the LinkToken contract.
2. The LinkToken contract transfers funds to the ownerless consumer.
3. The ownerless consumer requests randomness from the
   [VRF Coordinator](../../../contracts/src/v0.6/VRFCoordinator.sol), using the
   LINK from step 2 to pay for it.

To request randomness for your chosen consumer, run:
```shell
go run main.go ownerless-consumer-request --link-address=$LINK --consumer-address=$CONSUMER --key-hash=$KEY_HASH
```

You should see the output:
```
TX Hash: <YOUR TX HASH>
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
