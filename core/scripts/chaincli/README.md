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
