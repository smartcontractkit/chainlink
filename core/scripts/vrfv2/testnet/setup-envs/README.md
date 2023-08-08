## NOTE:
* Script will delete all existing jobs on the node!
* Currently works only with 0 or 1 VRF Keys on the node! Otherwise, will stop execution!
* Currently possible to fund all nodes with one amount of native tokens
## Commands:
1. If using Docker Compose
   1. create `.env` file in `core/scripts/vrfv2/testnet/docker` (can use `sample.env` file as an example)
   2. go to `core/scripts/vrfv2/testnet/docker` folder and start containers - `docker compose up`
2. Populate `./creds/` folder with relevant credentials for the nodes
3. Ensure that following env variables are set
```
export ETH_URL=
export ETH_CHAIN_ID=
export ACCOUNT_KEY=
```
3. execute from `core/scripts/vrfv2/testnet/setup-envs` folder
```
go run . \
--vrf-primary-node-url=http://localhost:6610 \
--vrf-primary-creds-file ./creds/vrf-primary-node.txt \
--vrf-backup-node-url=http://localhost:6611 \
--vrf-bk-creds-file ./creds/vrf-backup-node.txt \
--bhs-node-url=http://localhost:6612 \
--bhs-creds-file ./creds/bhs-node.txt \
--bhs-backup-node-url=http://localhost:6613 \
--bhs-bk-creds-file ./creds/bhs-backup-node.txt \
--bhf-node-url=http://localhost:6614 \
--bhf-creds-file ./creds/bhf-node.txt \
--num-eth-keys 5 \
--num-vrf-keys 1 \
--sending-key-funding-amount 100000000000000000

```

Optional parameters - will not be deployed if specified (NOT WORKING YET)
```
   --link-address 0x8606681e2295B2C4fD44D08E4E8a4D3180071559 \
   --link-eth-feed 0xD46fbB21875EBA71aF1b669c25EEe3692e5B9F13 \
   --bhs-address 0xbD84DbFaC527bd150384B75BA6cD75286F48da28 \
   --batch-bhs-address 0x126AeF1346A81003B8f1B6FBE8742536B2D22F71 \
   --coordinator-address 0x24d42DcD17C92d99100dce9D6133A7496e988A79 \
   --batch-coordinator-address 0xA276bCB4e67d787ce97F8787f2F454A4253Ce0Da 
```