## Setup

Before starting, you will need:
1. A working [Go](https://go.dev/doc/install) installation
2. EVM chain endpoint URLs
   - The endpoint can be a local node, or an externally hosted node, e.g. [alchemy](alchemy.com) or [infura](infura.io)
   - Both the HTTPS and WSS URLs of your endpoint are needed
3. The chain ID corresponding to your chain, you can find the chain ID for your chosen chain [here](https://chainlist.org/)
4. The private key of an account funded with LINK, and the chain's native token (to pay transaction fees)
   - Steps for exporting your private key from Metamask can be found [here](https://metamask.zendesk.com/hc/en-us/articles/360015289632-How-to-Export-an-Account-Private-Key)
5. The LINK address, LINK-ETH feed address, fast gas feed address for your chain
6. Install [docker](https://docs.docker.com/get-docker/) for CLI and GUI (optional)
7. \[Optional\] get a [tenderly API key](https://docs.tenderly.co/other/platform-access/how-to-generate-api-access-tokens) and find your [username / project name](https://docs.tenderly.co/other/platform-access/how-to-find-the-project-slug-username-and-organization-name).

The example .env in this repo is for the Polygon Mumbai testnet. You can use [this faucet](https://faucets.chain.link/mumbai) to send testnet LINK
to your wallet ahead of executing the next steps

>Note: Be careful with your key. When using testnets, it's best to use a separate account that does not hold real funds.

## Run OCR2Keepers locally

Build a local copy of the chainlink docker image by running this command in the root directory of the chainlink repo:

```bash
docker build -t chainlink:local -f ./core/chainlink.Dockerfile .
```

Next, from the root directory again, `cd` into the chaincli directory:

```shell
cd core/scripts/chaincli
```

Build `chaincli` by running the following command:

```shell
go build
```

Create the `.env` file based on the example `.env.example`, adding the node endpoint URLs and the private key of your wallet

### Keeper Registry
Next, use chaincli to deploy the registry:

Example:
```shell
./chaincli keeper registry deploy
```

Other options include:
- `./chaincli keeper registry update`: update existing keeper registry
- `./chaincli keeper registry withdraw`: cancel upkeeps and withdraw funds from registry
- `./chaincli keeper registry verify <contract-addr> <constructor-args>`: verify keeper registry contract

As the `keeper registry deploy` command executes, _two_ address are written to the terminal:

- KeeperRegistry2.0 Logic _(can be ignored)_
- KeeperRegistry2.0

The second address, `KeeperRegistry2.0` is the address you need; in the `.env` file, set `KEEPER_REGISTRY_ADDRESS` variable to the `KeeperRegistry2.0` address.

Note that this command doesn't run contract verification by default. If you want to run verification (eth, op and arb supported), config your .env and add the `--verify=true` flag in command.

If you already have keeper registry contract deployed and want to run only contract verification, you can use the following command:

```shell
./chaincli keeper registry verify <contract-addr> <constructor-args>
```

### Bootstrap Nodes
Run the following `bootstrap` command to start bootstrap nodes:

Example:
```shell
./chaincli bootstrap
```

Other options include:
- `--ui-port`: default `5688`, the Chainlink node UI listen port
- `--p2pv2-port`: default `8000`, the Chainlink node P2P listen port
- `--force | -f`: default `false`, if existing containers should be forcefully removed

The output of this command will show the tcp address of the deployed bootstrap node in the following format: `<p2p-key>@bootstrap:8000`.
Copy this entire string, including the `@bootstrap:8000` suffix, and the set the `BOOTSTRAP_NODE_ADDR` variable to this address in the `.env` file.

### Keeper launch and test
Once the bootstrap node is running, run the following command to launch the ocr2keeper nodes:

Example:
```shell
./chaincli keeper launch-and-test
```

Other options include:
- `--withdraw | -w`: default `true`, if funds should be withdrawn and upkeeps should be canceled after the test
- `--export-logs | -l`: default `false`, if container logs should be exported to ./ directory
- `--force | -f`: default `false`, if existing containers should be forcefully removed

You can also combine the `bootstrap` and `launch-and-test` commands into a single command:

```shell
./chaincli keeper launch-and-test --bootstrap
```
In the output of this command, you will see the http address of the nodes, e.g. `http://localhost:6688`. This is the Chainlink Operator GUI. You can use the default username `notreal@fakeemail.ch` and password `fj293fbBnlQ!f9vNs` to log in.

### Logs
Now that the nodes are running, you can use the `logs` subcommand to stream the output of the containers to your local terminal:

Example:
```shell
./chaincli keeper logs
```

Other options include:
- `--container-pattern`: default `^/keeper-\d+$`, regex pattern of container names to listen to for logs
- `--grep [string terms]`: default `empty string`, comma separated list of terms logs must include
- `--grepv [string terms]`: default `empty string`, comma separated list of terms logs must not include


You can use the `grep` and `grepv` flags to filter log lines, e.g. to only show output of the ocr2keepers plugin across the nodes, run:

```shell
./chaincli keeper logs --grep keepers-plugin
```

---