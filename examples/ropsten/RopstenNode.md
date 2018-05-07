# Running your own Chainlink node on Ropsten

## Syncing a Ropsten Ethereum Node

### [DevNet](https://github.com/smartcontractkit/devnet)

- Clone the ropository

```bash
git clone https://github.com/smartcontractkit/devnet.git
```

- Enter into the directory

```bash
cd devnet
```

- Run the command to invoke DevNet on Ropsten

```bash
$ make ropsten
```

### [Geth](https://github.com/ethereum/go-ethereum)

```
$ geth --testnet --ws --wsaddr 127.0.0.1 --wsport 8546 --wsorigins "*"
```

### [Parity](https://github.com/paritytech/parity)

```
$ parity --chain=ropsten --ws-interface 127.0.0.1 --ws-port 8546 --ws-origins "all"
```

## Environment Variables

Use the following environment variables as an example to configure your node for Ropsten:

    LOG_LEVEL=debug
    ROOT=~/.ropsten
    ETH_URL="ws://localhost:18546"
    ETH_CHAIN_ID=3
    TX_MIN_CONFIRMATIONS=2
    TASK_MIN_CONFIRMATIONS=2
    USERNAME=chainlink
    PASSWORD=twochains
    LINK_CONTRACT_ADDRESS=0x20fe562d797a42dcb3399062ae9546cd06f63280
    MINIMUM_CONTRACT_PAYMENT=1000000000000

## Running the Node

Once environment variables are set, run the node with:

```bash
$ chainlink node
```

When running the node for the first time, it will ask for a password and confirmation password. It will use this password to create a keystore file for you at `$ROOT/keys`. It will then display the following, showing you the address and its current balance:

```
2018-05-07T17:20:50Z [WARN]  0 Balance. Chainlink node not fully functional, please deposit ETH into your address: 0xC9EED6F5018E6aB95c03FcDfe661e38e97018235 cmd/client.go:70        
2018-05-07T17:20:50Z [INFO]  ETH Balance for 0xC9EED6F5018E6aB95c03FcDfe661e38e97018235: 0.000000000000000000 cmd/client.go:71        
2018-05-07T17:20:50Z [INFO]  Link Balance for 0xC9EED6F5018E6aB95c03FcDfe661e38e97018235: 0.000000000000000000 cmd/client.go:74
```

## Getting Ropsten ETH

Visit the faucet [here](http://faucet.ropsten.be:3001/) and paste your node's address to receive Ropsten ETH.

## Adding Jobs

We have example [JobSpecs](https://github.com/smartcontractkit/chainlink/wiki/Job-Pipeline) in the `jobs/` directory. They can be used on your node once you have [deployed your oracle contract](./OracleContract.md) by replacing the address with that of your deployed contract.

Adding jobs can be done by using the command `chainlink c` with the path to the JobSpec file.

EthUint256:

```
$ chainlink c jobs/EthUint256Job.json
```

EthInt256

```
$ chainlink c jobs/EthInt256Job.json
```

EthBytes32

```
$ chainlink c jobs/EthBytes32.json
```