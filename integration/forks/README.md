# Integration test for ommers and re-orgs

## Docker images

From `integration/forks/bootnode`

```
docker build . -t bootnode:latest
```

From `integration/forks`
```
docker build . -t geth:latest
```

## Invoking this framework

Run `docker-compose --compatibility up`.

To trigger a block-reorg on `chainlink`/`geth1`, disconnect `geth2` from the network with

```
docker network disconnect forks_gethnet forks_geth2_1
```

Here `forks_gethnet` is the appropriate network name as reported by `docker network  ls | grep gethnet | awk '{ print $2 }'`, and `forks_geth2_1` is the appropriate container name as reported by `docker ps -a | grep geth2 | awk '{ print $(NF) }'`. These names may change depending on `docker-compose` and the directory in which you run this test.

Once `geth1`'s logs indicate that it's mined a few solo blocks, reconnect with `docker network connect <network name> <container name>`. You should see output like this in the logs:

```
geth1_1      | INFO [06-13|18:06:26.284] Commit new mining work                   number=156 sealhash=78d46c…e7a4bb uncles=1 txs=0 gas=0 fees=0 elapsed=113.312µs
chainlink_1  | 2019-06-13T18:06:26Z [DEBUG] Received new head #156 (0x9c)                      services/head_tracker.go:208     blockHash=0x8c7278a3efa544c74429cc0101dadbc39ff30282f03039745a0970d0668d196e blockHeight=156
chainlink_1  | 2019-06-13T18:06:26Z [DEBUG] Received new head #157 (0x9d)                      services/head_tracker.go:208     blockHash=0x4801e4572085b45014f0524390bae66b5b40136ab9b02bb4c1685f9fc73248ce blockHeight=157
geth2_1      | INFO [06-13|18:06:26.343] Successfully sealed new block            number=175 sealhash=cdfef5…f041a4 hash=04916d…c6a75c elapsed=613.169ms
geth2_1      | INFO [06-13|18:06:26.343] 🔗 block reached canonical chain          number=168 hash=12ae21…28054d
geth2_1      | INFO [06-13|18:06:26.343] 🔨 mined potential block                  number=175 hash=04916d…c6a75c
geth2_1      | INFO [06-13|18:06:26.344] Commit new mining work                   number=176 sealhash=4d3dfc…0f3c3f uncles=0 txs=0 gas=0 fees=0 elapsed=107.324µs
chainlink_1  | 2019-06-13T18:06:26Z [DEBUG] Received new head #158 (0x9e)                      services/head_tracker.go:208     blockHash=0x134690c8746d36c45743655bf3875b89a5354f8ae512835fb4aeb02ca829e480 blockHeight=158
geth1_1      | INFO [06-13|18:06:26.401] 😱 block lost                             number=153 hash=e0d317…74ff50
geth1_1      | INFO [06-13|18:06:26.401] 😱 block lost                             number=154 hash=bd5019…a4d07f
geth1_1      | INFO [06-13|18:06:26.401] 😱 block lost                             number=155 hash=c94194…cad7ab
geth1_1      | INFO [06-13|18:06:26.401] Commit new mining work                   number=174 sealhash=b731b1…29c274 uncles=0 txs=0 gas=0 fees=0 elapsed=221.367µs
chainlink_1  | 2019-06-13T18:06:26Z [DEBUG] Received new head #159 (0x9f)                      services/head_tracker.go:208     blockHash=0x95152d3d66bf9a37733009e73667bb87628824d1cd03f22c993b9c79055d3699 blockHeight=159
chainlink_1  | 2019-06-13T18:06:26Z [DEBUG] Received new head #160 (0xa0)                      services/head_tracker.go:208     blockHash=0x110b7aafcc70198741ef7fa726188a83d79666c2331c14772b4b9f99e475c283 blockHeight=160
chainlink_1  | 2019-06-13T18:06:26Z [DEBUG] Received new head #161 (0xa1)                      services/head_tracker.go:208     blockHash=0xc5703765b4d5d68453935628cb65fdf80967f88d5a87d55fb6e12cf61be58446 blockHeight=161
chainlink_1  | 2019-06-13T18:06:26Z [DEBUG] Received new head #162 (0xa2)                      services/head_tracker.go:208     blockHash=0x7246e2c8e3fc1c3f113bf95f9d01ce77c1ba747c5cf66b9e81294928d23edb9a blockHeight=162
chainlink_1  | 2019-06-13T18:06:26Z [DEBUG] Received new head #163 (0xa3)                      services/head_tracker.go:208     blockHash=0xf889f05c445d10a91b86a98dbf1759735f6e5fc4fae61c91aaaa06d2e1ddca1b blockHeight=163
chainlink_1  | 2019-06-13T18:06:26Z [DEBUG] Received new head #164 (0xa4)                      services/head_tracker.go:208     blockHash=0x94ff71d1b0b75af79879a7daca2acf1cffe781203d7ee355c0518f2e2b7fca32 blockHeight=164
chainlink_1  | 2019-06-13T18:06:26Z [DEBUG] Received new head #165 (0xa5)                      services/head_tracker.go:208     blockHash=0x7c493b8ed0bf880518c3e37a278480dfc26791f536f03498b1ec952a2c1b79f4 blockHeight=165
chainlink_1  | 2019-06-13T18:06:26Z [DEBUG] Received new head #166 (0xa6)                      services/head_tracker.go:208     blockHash=0xaf2351113c7242aef9400af21a29efc14133434d846395dc17c770bcc4b2b215 blockHeight=166
chainlink_1  | 2019-06-13T18:06:26Z [DEBUG] Received new head #167 (0xa7)                      services/head_tracker.go:208     blockHash=0x77dbce91f4fc5612a8b37c174faa1085e4db5c5c49382c990b3ea3699e8dccc1 blockHeight=167
chainlink_1  | 2019-06-13T18:06:26Z [DEBUG] Received new head #168 (0xa8)                      services/head_tracker.go:208     blockHash=0x24cb977e9246d18079aec6fab9f49a54a4ffcd7b2eeadb7b955883e50f749dd1 blockHeight=168
chainlink_1  | 2019-06-13T18:06:26Z [DEBUG] Received new head #169 (0xa9)                      services/head_tracker.go:208     blockHash=0xdfff32a55c9491930525e79db4b2db7d265de1b0d5172fda08d62b5f8a54a24a blockHeight=169
chainlink_1  | 2019-06-13T18:06:26Z [DEBUG] Received new head #170 (0xaa)                      services/head_tracker.go:208     blockHash=0x50b719df858d08451fc82cdf483e72472929fa0a4aa955558a8dc1c874635c3c blockHeight=170
chainlink_1  | 2019-06-13T18:06:26Z [DEBUG] Received new head #171 (0xab)                      services/head_tracker.go:208     blockHash=0x3c28ab4476da97167261ca4ab4ee45a0955d4c3be8469094dcf437e5a1b51dfc blockHeight=171
chainlink_1  | 2019-06-13T18:06:26Z [DEBUG] Received new head #172 (0xac)                      services/head_tracker.go:208     blockHash=0x25aa630d9893ff04f64a28aaef8c93df7f63502d5bea2026d572de34b221c73f blockHeight=172
chainlink_1  | 2019-06-13T18:06:26Z [DEBUG] Received new head #173 (0xad)                      services/head_tracker.go:208     blockHash=0xc39a4adbbf31c30802ef67a59e668f0baf2f092618b0c4fee7ee6a64f01071a5 blockHeight=173
geth1_1      | INFO [06-13|18:06:26.499] Imported new chain segment               blocks=1  txs=0 mgas=0.000 elapsed=112.136ms mgasps=0.000 number=174 hash=de9866…b68a3e dirty=15.55KiB
chainlink_1  | 2019-06-13T18:06:26Z [DEBUG] Received new head #174 (0xae)                      services/head_tracker.go:208     blockHash=0x78974a112416d3d6add0b80c2bdb47a96cf360ed6f578e2162fa9097e963d878 blockHeight=174
geth1_1      | INFO [06-13|18:06:26.596] Commit new mining work                   number=175 sealhash=1d436e…d06e4b uncles=0 txs=0 gas=0 fees=0 elapsed=69.755µs
geth1_1      | INFO [06-13|18:06:26.695] Imported new chain segment               blocks=1  txs=0 mgas=0.000 elapsed=192.261ms mgasps=0.000 number=175 hash=04916d…c6a75c dirty=15.79KiB
chainlink_1  | 2019-06-13T18:06:26Z [DEBUG] Received new head #175 (0xaf)                      services/head_tracker.go:208     blockHash=0x145cab063239c7d376ea75ecc90bddae0c4e585c8e0d925864b9bb25c06a93be blockHeight=175
geth1_1      | INFO [06-13|18:06:26.697] Commit new mining work                   number=176 sealhash=4d3dfc…0f3c3f uncles=0 txs=0 gas=0 fees=0 elapsed=90.669µs
geth2_1      | INFO [06-13|18:06:27.067] Generating DAG in progress               epoch=1 percentage=68 elapsed=2m1.490s
```

The `block lost` entries indicate that the blocks `geth1` solo-mined have been discarded in favor of `geth2`'s longer/heavier chain.


## Construction of this initial blockchain state

To reproduce a rough facsimile of the artifacts in this directory, follow these steps. Be careful: `geth` is not very user-friendly, and, given inconsistent inputs, will behave in unintuitive, cryptic, and disappointing ways.

1. Copy and modify `geth-config-1.toml`. The original came from the output of `geth dumpconfig`.
2. Run `initialize-node-data.sh <n>` for `<n>` ranging over the `<n>`s in the `geth-config-<n>.toml`s. That script has more details about what it does.
3. Put all the enode addresses output by `initialize-node-data.sh` in all the *other* `geth-config-<n>.toml`'s `StaticNodes` lists. E.g., if `initialize-node-data.sh 2` outputs `enode://addr@ipaddr:port`, the `StaticNodes` list in `geth-config-1.toml` should contain `"enode://addr@ipaddr:port"`. Note that `ipaddr` will need to change, based on the IP address you assign to `geth2` in the next step.
4. Make stanzas corresponding to each `geth-config-<n>.toml` in `docker-compose.yaml`, like
   ```
     geth1:
       image: ethereum/client-go
       restart: on-failure
       command:
         --mine
         --miner.threads 1
         --config /root/geth-config-1.toml
         --unlock "0x9ca9d2d5e04012c9ed24c0e513c9bfaa4a2dd77f"
         --password /run/secrets/node_password
         --nat extip:172.16.1.100
       volumes:
         - .:/root
       networks:
         gethnet:
           ipv4_address: 172.16.1.100
       ports: []
       secrets:
         - node_password
   ```
   (But this is just an example. Crib from the stanza in `docker-compose.yaml`, because that'll break noisily if it goes out of date.)

   Note that it's essential for `--miner.threads` to be positive. The default value is `0`, and if you leave it at that, `geth` will indicate that it's started mining with `Commit new mining work`, but will not actually mine.

   You don't need the account information (`--unlock`, `--password`) on instances which chainlink won't connect to.

5. Adjust the IP addresses in the `geth-config-<n>.toml`s to match those in the corresponding `docker-compose.yaml` stanzas. That includes the IP addresses in the enode addresses.
6. Run `docker-compose up geth1 ... geth<n>`. Check the containers' outputs to verify that they are sharing blocks. For instance, you should see matching hashes between containers, like this:
   ```
   geth2_1      | INFO [06-13|07:27:36.830] Imported new chain segment               blocks=1 txs=0 mgas=0.000 elapsed=4.949ms   mgasps=0.000 number=8 hash=a922c6…b0a076 dirty=1.91KiB
       ^                                                                                                                                                    ^^^^^^^^^^^^^
   geth1_1      | INFO [06-13|07:27:59.247] 🔗 block reached canonical chain          number=8  hash=a922c6…b0a076
       ^                                                                                             ^^^^^^^^^^^^^
   ```
   This means that container `geth1_1`'s block with hash a922c6…b0a076 was imported and accepted by container `geth2_1**.

## Constructing the trigger transaction

The goal of this framework is to test how chainlink responds to logs from the
ethereum blockchain in the context of a re-org. Therefore, we need to construct
a transaction which emits a log. Here is how I did this. You may want to
automate this, if you're adjusting your solidity contract a lot.

1. Get the binary for `Trigger.sol` with `solc --bin contracts/Trigger.sol`.
2. Place that binary output in the assignment to `contract_data` in
   `generate_tx.rb`.
3. `sudo apt install ruby ruby-dev ruby-bundler`
4. Build `digest-sha3-ruby` with `CFLAGS+-Wno-format-security`. (Straight `gem
   install eth` fails for me, [as described
   here](https://github.com/phusion/digest-sha3-ruby/issues/7).)

   `git clone https://github.com/izetex/digest-sha3-ruby`

   `cd digest-sha3-ruby; gem build digest-sha3.gemspec`

   `sudo gem install digest-sha3-1.1.0.gem`
5. `sudo gem install eth`
6. Test the script by running `ruby generate_tx.rb`. The output from this is
   assigned to `${TRANSACTION_HEX}` in `Makefile`.
