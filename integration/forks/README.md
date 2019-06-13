# Integration test for ommers and re-orgs

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
   This means that container `geth1_1`'s block with hash a922c6…b0a076 was imported and accepted by container `geth2_1`.
