# Running your own Chainlink node on Ropsten

## Syncing a Ropsten Ethereum Node

### [DevNet](https://github.com/smartcontractkit/infrastructure/tree/master/image)

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
make ropsten
```

Note that by default, Devnet will run the websocket protocl on port `18546` instead of the default `8546` as with Geth and Parity.

### [Geth](https://github.com/ethereum/go-ethereum)

```
geth --testnet --ws --wsaddr 127.0.0.1 --wsport 8546 --wsorigins "*"
```

### [Parity](https://github.com/paritytech/parity)

```
parity --chain=ropsten --ws-interface 127.0.0.1 --ws-port 8546 --ws-origins "all"
```

## Deploying your own Oracle contract

- In Remix, import the contracts at `chainlink/solidity/contracts`
- Click on the `Oracle.sol` contract in the left side-bar
- You may need to update the import paths since Remix doesn't support folders
  - You can also use a tool like [truffle-flattener](https://www.npmjs.com/package/truffle-flattener) to simplify this process
- On the Compile tab, click on the "Start to compile" button near the top-right

![compile](./images/12-41-31.png)

- Change to the Run tab
- Select Oracle from the dropdown in the right panel
- Copy and paste the line below and enter it into the text field next to the Deploy button <br>
    **0x20fE562d797A42Dcb3399062AE9546cd06f63280**
- Click Deploy

![create](./images/12-42-32.png)

- Metamask will prompt you to Confirm the Transaction
- You will need to choose a Gas Price (use 4 if you don't know what to pick)
- Select Submit

![deploy contracts](./images/11-03-14.png)

- A link to Etherscan will display at the bottom, you can open that in a new tab to keep track of the transaction

![confirm contract deploy](./images/12-43-32.png)

- Once successful, you should have a new address for the deployed contract

![contract deploy successful](./images/07-25-49.png)

- Keep note of the Oracle contract's address, you will need it as an environment variable and for adding a JobSpec to the node.

## Environment Variables

You can use [direnv](https://github.com/direnv/direnv/) to set environment variable for you when you enter a directory.

Use the following environment variables as an example to configure your node for Ropsten:

    ETH_CHAIN_ID=3
    MIN_INCOMING_CONFIRMATIONS=0
    LINK_CONTRACT_ADDRESS=0x20fe562d797a42dcb3399062ae9546cd06f63280
    ORACLE_CONTRACT_ADDRESS=***Your deployed oracle contract address***

_If you want to require HTTPS, you'll need to [create your own self-signed certificates](https://github.com/smartcontractkit/chainlink/wiki/Creating-Self-Signed-Certificates). Otherwise, add `CHAINLINK_DEV=true` to your environment variables._

## Running the Node

Once environment variables are set, run the node with:

```bash
cd $GOPATH/src/github.com/smartcontractkit/chainlink
make install
chainlink node
```

*You can also [run the node with Docker](https://github.com/smartcontractkit/chainlink/wiki/Running-the-Docker-Image).*

When running the node for the first time, it will ask for a password and confirmation password. It will use this password to create a keystore file for you at `$ROOT/keys`. You will then be prompted to enter an API Email and Password. This will be used to authenticate sessions with the node so you can manage it via the command-line or GUI. It will then display the following, showing you the address and its current balance:

```
2018-05-07T17:20:50Z [WARN]  0 Balance. Chainlink node not fully functional, please deposit ETH into your address: 0xC9EED6F5018E6aB95c03FcDfe661e38e97018235 cmd/client.go:70        
2018-05-07T17:20:50Z [INFO]  ETH Balance for 0xC9EED6F5018E6aB95c03FcDfe661e38e97018235: 0.000000000000000000 cmd/client.go:71        
2018-05-07T17:20:50Z [INFO]  Link Balance for 0xC9EED6F5018E6aB95c03FcDfe661e38e97018235: 0.000000000000000000 cmd/client.go:74
```

## Getting Ropsten ETH

Visit the faucet [here](http://faucet.ropsten.be:3001/) and paste your node's address to receive Ropsten ETH.

## Transfer Ownership of the Oracle contract

Back in Remix, call the `transferOwnership` function with the address of your node as the input parameter. This will allow your node the ability to fulfill requests.

You can get the address of your node by running:

```bash
chainlink a
```

## Adding and Running Jobs

With your own Oracle contract, you can use your own node to fulfill requests. With the `RopstenConsumer.sol` example contract, you would simply change the value for `address constant ROPSTEN_ORACLE_ADDRESS = ` to your deployed oracle contract address and the value for each of the `*_SPEC_ID` constants to the Spec ID reported from your node.

In order to add the specs to your node, we have example [JobSpecs](https://github.com/smartcontractkit/chainlink/wiki/Job-Pipeline) in the `jobs/` directory. Adding specs can be done by using the command `chainlink c` with the path to the JobSpec file.

**You will need to modify each JobSpec to update the `"address"` for the `RunLog` initiator to your deployed oracle contract address.**

EthUint256:

```
chainlink c jobs/EthUint256Job.json
```

EthInt256

```
chainlink c jobs/EthInt256Job.json
```

EthBytes32

```
chainlink c jobs/EthBytes32.json
```

Once the specs have been added, you'll update their value in the `RopstenConsumer.sol` contract, like so (replacing the values below with your actual IDs):

```js
bytes32 constant PRICE_SPEC_ID = bytes32("0bdf244a39234ea6a416c9d37d66c701");
bytes32 constant CHANGE_SPEC_ID = bytes32("cbe3b0be008747da92f66105998bdad4");
bytes32 constant MARKET_SPEC_ID = bytes32("d107362cef564d35b02792d7733e1481");
```
