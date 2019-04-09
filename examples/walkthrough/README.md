# Walkthrough

## Setup

### [Install Dependencies](https://github.com/smartcontractkit/chainlink#install)

For this walkthrough, you will need:

- [Go](https://golang.org/doc/install#install)
- [dep](https://github.com/golang/dep#installation)
- [Node](https://nodejs.org/en/download/package-manager/)
- [Yarn](https://yarnpkg.com/lang/en/docs/install/#mac-stable)
- [Docker](https://www.docker.com/get-started)
- [direnv](https://direnv.net/) (optional)

If you're on a mac with [Homebrew](https://brew.sh/), you can run:

`brew install go dep node yarn docker direnv`

### Install Chainlink


```bash
go get -d github.com/smartcontractkit/chainlink
cd $GOPATH/src/github.com/smartcontractkit/chainlink
make install
chainlink help
```

## Run

Spin up a local Ethereum node:

```bash
./tools/bin/devnet
```

In a separate window:

```bash
CHAINLINK_DEV=true chainlink local node
... login ...
> [INFO]  Link Balance for 0x79dBA5B14cBA2560360c2eF48e9329aC7Ab21573
```

In a final window, which you'll use for the rest of the walktrhough, save the Chainlink node's Ethereum address:
```bash
ORACLE_NODE=0x79dBA5B14cBA2560360c2eF48e9329aC7Ab21573
```

## Set up your LINK and Oracle contracts
In the same window where you set ORACLE_NODE, move to the solidity directory: `cd solidity/`

```
./bin/fund_address # fund local wallet
./bin/fund_address $ORACLE_NODE # fund the oracle node's wallet
```

Deploy the LINK token:

```
./bin/deploy LinkToken.sol
> LinkToken.sol successfully deployed: 0x9af9c91e1f5e22d1a7ea8e8af3cb3b3f858a619d
LINK_TOKEN=0x9af9c91e1f5e22d1a7ea8e8af3cb3b3f858a619d
```

Deploy an Oracle:

```
./bin/deploy Oracle.sol $LINK_TOKEN
> Oracle.sol successfully deployed: 0x22f9c91E1f5E22D1a7eA8E8AF3CB3b3f858a6122
ORACLE_CONTRACT=0x22f9c91E1f5E22D1a7eA8E8AF3CB3b3f858a6122
```
Transfer Oracle ownership to your node


```
./bin/transfer_ownership $ORACLE_CONTRACT $ORACLE_NODE
> ownership of 0x22f9c91E1f5E22D1a7eA8E8AF3CB3b3f858a6122 transferred to 0x79dBA5B14cBA2560360c2eF48e9329aC7Ab21573
```

### Create your first job spec:

Log into your node at [http://localhost:6688](http://localhost:6688), and then click the `Create Job` button in the top right of the dashboard.

Edit the following JSON so that it includes your Oracle address for the initiators address field:
```
{
  "initiators": [{
      "type": "runlog",
      "params": {"address": "0x22f9c91E1f5E22D1a7eA8E8AF3CB3b3f858a6122"}
  }],
  "tasks": [
    {"type": "httpget"},
    {"type": "jsonparse"},
    {"type": "ethbytes32"},
    {"type": "ethtx"}
  ]
}
```

Once it has successfully created, grab your Job Spec ID, it should look be hex again, but shorter: `8452ff74ebe745e0ab9f7edddb16ecb0 `


# Use your oracle:
Deploy a data consumer:

```
./bin/deploy BasicConsumer.sol $LINK_TOKEN $ORACLE_CONTRACT 32bac7e83e5e478084678fb9c0ac6704
> BasicConsumer.sol successfully deployed: 0x8e368fb378ff79efb8181d3203edd7ad4d70736e
CONSUMER_CONTRACT=0x8e368fb378ff79efb8181d3203edd7ad4d70736e
```

Request an Ethereum price:

```
./bin/update_eth_price $CONSUMER_CONTRACT
> FAILED!!!
```

Check price on contract:

```
./bin/view_eth_price $CONSUMER_CONTRACT
> No price listed
```

Your Consumer needs LINK. Transfer some LINK:

```
./bin/transfer_tokens $LINK_TOKEN $CONSUMER_CONTRACT
> 1000 LINK successfully sent to 0x55ad1706ca8cf0ac593b918c105944487d0737b2
```

Request an Ethereum price again:

```
./bin/update_eth_price $CONSUMER_CONTRACT
> price successfully requested
```

Check price on contract:

```
./bin/view_eth_price $CONSUMER_CONTRACT
> current ETH price: 230.04
```
