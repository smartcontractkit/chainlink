## NOTE:
* Script will delete all existing jobs on the node!
* Currently works only with 0 or 1 VRF Keys on the node! Otherwise, will stop execution!
* Currently possible to fund all nodes with one amount of native tokens
## Commands:
1. If using Docker Compose
   1. create `.env` file in `core/scripts/common/vrf/docker` (can use `sample.env` file as an example)
   2. go to `core/scripts/common/vrf/docker` folder and start containers - `docker compose up`
2. Update [rpc-nodes.toml](..%2Fdocker%2Ftoml-config%2Frpc-nodes.toml) with relevant RPC nodes
3. Create files with credentials desirably outside `chainlink` repo (just not to push creds accidentally). Populate the files  with relevant credentials for the nodes
4. Ensure that following env variables are set
```
export ETH_URL=
export ETH_CHAIN_ID=
export ACCOUNT_KEY=
```
5. execute from `core/scripts/common/vrf/setup-envs` folder
```
go run . \
--vrf-version="v2plus" \
--vrf-primary-node-url=http://localhost:6610 \
--vrf-primary-creds-file <path_to_file_with_creds> \
--vrf-backup-node-url=http://localhost:6611 \
--vrf-bk-creds-file <path_to_file_with_creds> \
--bhs-node-url=http://localhost:6612 \
--bhs-creds-file <path_to_file_with_creds> \
--bhs-backup-node-url=http://localhost:6613 \
--bhs-bk-creds-file <path_to_file_with_creds> \
--bhf-node-url=http://localhost:6614 \
--bhf-creds-file <path_to_file_with_creds> \
--deploy-contracts-and-create-jobs="true" \
--subscription-balance="1e19" \
--subscription-balance-native="1e18" \
--batch-fulfillment-enabled="true" \
--min-confs=3 \
--num-eth-keys=1 \
--num-vrf-keys=1 \
--sending-key-funding-amount="1e17"
```

Optional parameters - will not be deployed if specified (NOT WORKING YET)
```
   --link-address <address> \
   --link-eth-feed <address> \
   --bhs-address <address> \
   --batch-bhs-address <address> \
   --coordinator-address <address> \
   --batch-coordinator-address <address> 
```