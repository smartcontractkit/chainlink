

# Apocalypse stress test environment

## Usage

1. Build the Docker images for Chainlink, Geth, and Parity:
    ```sh
    $ ./scripts/env-build-images
    ```
2. Install [Blockade](https://github.com/worstcase/blockade)
3. Install node modules:
    ```sh
    $ npm i
    ```
4. Spin up the environment:
    ```sh
    $ blockade up
    ```
    This will start 2 Geth nodes, 1 Parity node (non-mining), 2 Chainlink nodes, and a Blockscout instance for inspecting the chain/network.  The Geth nodes will take a while to generate their mining DAGs, so the scripts in `scenarios` wait until that process is complete before proceeding.
5. Run one of the scenarios:
    ```sh
    $ node scenarios/flux-monitor.js
    ```
6. Open Blockscout at http://localhost:4000, sit back, and watch the havoc.
7. Use Blockade to simulate various adverse network conditions:
    - `blockade flaky gethnet gethnet2 paritynet`
    - `blockade slow chainlink_neil`
    - `blockade duplicate chainlink_nelly`
    - `blockade partition gethnet,paritynet,chainlink_neil gethnet2,chainlink_nelly`
8. When you're finished with the environment, run `./scripts/env-destroy`.  You might also need to delete the containers, volumes, and the Docker network.

## Scenarios

**Chain reorgs**

To force a chain reorg:
- Run `blockade partition gethnet,paritynet,chainlink_neil gethnet2,chainlink_nelly` to partition the containers into 2 distinct networks
- Run `node scenarios/tx-tornado.js`, which will spam the Ethereum nodes with hundreds of transactions per second
- Wait until several blocks are mined (use Blockscout for this)
- Run `blockade join`
- Check Blockscout's "reorgs" view to see if a reorg occurred

**Flux Monitor**





## Modifying the environment

**IP addresses**

IP addresses are assigned deterministically based on the contents of `blockade.yaml`, but there's no explicit way to assign specific IPs.  Usually, this isn't an issue, as you can simply use a container's name in place of its IP.  However, Geth does not allow `enodes` to be specified with anything other than explicit IP addresses.  So if you change the containers in `blockade.yaml`, you'll need to run `blockade up`, take note of the new IP addresses, and modify them in `gethnet/Dockerfile`.

**Personas**

If you need to modify the "personas" in the environment, there's a helper script in the `scripts` directory which modifies a JSON accounts database located in the `config` subdirectory.

```sh
$ ./scripts/config add-persona sergey
$ ./scripts/config rm-persona sergey
```

**Geth/Parity config**

Geth and Parity configs are spread across TOML config files and CLI flags.

Geth:
- gethnet/config.toml
- gethnet/Dockerfile

Parity:
- paritynet/miner.toml
- paritynet/start.sh

