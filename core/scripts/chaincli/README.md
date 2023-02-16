## Setup

Before starting, you will need:
1. An EVM chain endpoint URL
2. The chain ID corresponding to your chain
3. The private key of an account funded with LINK, and the chain's native token
   (to pay transaction fees)
4. [The LINK address, LINK-ETH feed address, fast gas feed address](https://docs.chain.link/docs/chainlink-keepers/introduction/#onboarding-steps)
   for your chain
5. [Go](https://go.dev/doc/install)
6. Running at least 2 nodes with the keeper job. Have some balance on each of them.

The endpoint URL can be a locally running node, or an externally hosted one like
[alchemy](https://www.alchemy.com/). Your chain ID will be a number
corresponding to the chain you pick. For example the Rinkeby testnet has chain
ID 4. Your private key can be exported from [MetaMask](https://metamask.zendesk.com/hc/en-us/articles/360015289632-How-to-Export-an-Account-Private-Key).

Note: Be careful with your key. When using testnets, it's best to use a separate
account that does not hold real funds.

1. "cd" into the keeper scripts  directory
```shell
cd <YOUR LOCAL CHAINLINK REPO>/core/scripts/chaincli
```
2. Create `.env` file based on the example `.env.example`

To see all available commands, run the following:
```bash
go run main.go --help
```

### Run OCR2Keepers on the local env

First, decide which CL node version to use, or build a new one using...

```bash
docker build -t chainlink:local -f ./core/chainlink.Dockerfile .
```

Before start, there should be `.env` file with all required environment variables. Example for Goerli network:
```.dotenv
CHAINLINK_DOCKER_IMAGE=chainlink:local
NODE_URL=<wss-rpc-node-addr>
CHAIN_ID=5
PRIVATE_KEY=<wallet-private-key>
LINK_TOKEN_ADDR=0x326C977E6efc84E512bB9C30f76E30c160eD06FB
LINK_ETH_FEED=0xb4c4a493AB6356497713A78FFA6c60FB53517c63
FAST_GAS_FEED=0x22134617ae0f6ca8d89451e5ae091c94f7d743dc
FUND_CHAINLINK_NODE=500000000000000000000 # 5 ETH

# Keepers config
PAYMENT_PREMIUM_PBB=200000000
FLAT_FEE_MICRO_LINK=1
CHECK_GAS_LIMIT=6500000
STALENESS_SECONDS=90000
GAS_CEILING_MULTIPLIER=1
MIN_UPKEEP_SPEND=0
MAX_PERFORM_GAS=5000000
MAX_CHECK_DATA_SIZE=2000
MAX_PERFORM_DATA_SIZE=2000
FALLBACK_GAS_PRICE=200000000
FALLBACK_LINK_PRICE=3684210526315790
TRANSCODER=0x97aFFbaE5d31965eAA427Dd4DD6Cd22271561853
REGISTRAR=0x0000000000000000000000000000000000000000
KEEPER_REGISTRY_ADDRESS=<registry-address-from-first-step>
BOOTSTRAP_NODE_ADDR=<bootstrap-node-addr-from-second-step>

KEEPER_OCR2=true
KEEPER_REGISTRY_VERSION=4

BLOCK_COUNT_PER_TURN=20
KEEPERS_COUNT=4
UPKEEP_TEST_RANGE=1000
UPKEEP_AVERAGE_ELIGIBILITY_CADENCE=20
UPKEEP_COUNT=1
UPKEEP_ADD_FUNDS_AMOUNT=5000000000000000000 # 5 LINK
```

1. First we need to deploy the registry if there is no one already deployed:
```shell
$ chaincli keeper registry deploy
```
The address should be in the output of this command. The address should be defined within `KEEPER_REGISTRY_ADDRESS` evar.

2. Then we should get the bootstrap node up and running using the registry contract:
```shell
$ chaincli bootstrap <registry-contract-address>
```
The output will show the tcp address of the deployed bootstrap node in the following format: `<p2p-key>@bootstrap:8000`.
This address should be defined within `BOOTSTRAP_NODE_ADDR` evar which gonna be used in the next step during OCR2Keeper nodes setup.

3. Once we have a bootstrap node up and running, ocr2keeper nodes are ready to be created.
```shell
$ chaincli keeper launch-and-test
```
