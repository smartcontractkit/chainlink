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
   * `--vrf-version` - "v2" or "v2plus"

#### VRF V2
```
go run . \
--vrf-version="v2" \
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
--num-eth-keys=1 \
--num-vrf-keys=1 \
--num-bhs-sending-keys= 1 \
--num-bhf-sending-keys=1 \
--sending-key-funding-amount="1e17" \
--deploy-contracts-and-create-jobs="true" \
--subscription-balance="1e19" \
--subscription-balance-native="1e18" \
--batch-fulfillment-enabled="true" \
--batch-fulfillment-gas-multiplier=1.1 \
--estimate-gas-multiplier=1.1 \
--poll-period="5s" \
--request-timeout="30m0s" \
--reverts-pipeline-enabled="true" \
--min-confs=3 \
--simulation-block="latest" \
--bhs-job-wait-blocks=30 \
--bhs-job-look-back-blocks=200 \
--bhs-job-poll-period="1s" \
--bhs-job-run-timeout="1m" \
--register-vrf-key-against-address=<vrf key will be registered against this address 
in order to call oracleWithdraw from this address> \
--deploy-vrfv2-owner="true" \
--use-test-coordinator="true"
```
#### VRF V2 Plus
* does not need to register VRF key against address 
* does not need to deploy VRFV2Owner contract
* does not need to use test coordinator

VRF V2 Plus example:
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
--num-eth-keys=1 \
--num-vrf-keys=1 \
--num-bhs-sending-keys= 1 \
--num-bhf-sending-keys=1 \
--sending-key-funding-amount="1e17" \
--deploy-contracts-and-create-jobs="true" \
--subscription-balance="1e19" \
--subscription-balance-native="1e18" \
--batch-fulfillment-enabled="true" \
--batch-fulfillment-gas-multiplier=1.1 \
--estimate-gas-multiplier=1.1 \
--poll-period="5s" \
--request-timeout="30m0s" \
--min-confs=3 \
--simulation-block="latest" \
--proving-key-max-gas-price="1e12" \
--flat-fee-native-ppm=500 \
--flat-fee-link-discount-ppm=100 \
--native-premium-percentage=1 \
--link-premium-percentage=1 \
--bhs-job-wait-blocks=30 \
--bhs-job-look-back-blocks=200 \
--bhs-job-poll-period="1s" \
--bhs-job-run-timeout="1m" 
```

Optional parameters - will not be deployed if specified 
```
   --link-address <address> \
   --link-eth-feed <address> \
```

WIP - Not working yet:
```
   --bhs-address <address> \
   --batch-bhs-address <address> \
   --coordinator-address <address> \
   --batch-coordinator-address <address> 
```


## Process Example

1. If the CL nodes do not have needed amount of ETH and VRF keys, you need to create them first:
```
go run . \
--vrf-version="v2" \
--vrf-primary-node-url=<url> \
--vrf-primary-creds-file <path_to_file_with_creds> \
--bhs-node-url=<url> \
--bhs-creds-file <path_to_file_with_creds> \
--num-eth-keys=3 \
--num-vrf-keys=1 \
--sending-key-funding-amount="1e17" \
--deploy-contracts-and-create-jobs="false" 
```
Then update corresponding deployment scripts in infra-k8s repo with the new ETH addresses, specifying max gas price for each key

e.g.:
```
[[EVM.KeySpecific]]
Key = '<eth key address>'
GasEstimator.PriceMax = '30 gwei'
```

2. If the CL nodes already have needed amount of ETH and VRF keys, you can deploy contracts and create jobs with the following command:
NOTE - nodes will be funded at least to the amount specified in `--sending-key-funding-amount` parameter.
```
go run . \
--vrf-version="v2" \
--vrf-primary-node-url=<url> \
--vrf-primary-creds-file <path_to_file_with_creds> \
--bhs-node-url=<url> \
--bhs-creds-file <path_to_file_with_creds> \
--num-eth-keys=3 \
--num-vrf-keys=1 \
--sending-key-funding-amount="1e17" \
--deploy-contracts-and-create-jobs="true" \
--subscription-balance="1e19" \
--subscription-balance-native="1e18" \
--batch-fulfillment-enabled="true" \
--min-confs=3 \
--register-vrf-key-against-address="<eoa address>" \
--deploy-vrfv2-owner="true" \
--link-address "<link address>" \
--link-eth-feed "<link eth feed address>" 
``` 


3. We can run sample rand request to see if the setup works.
   After previous script was done, we should see the command to run in the console:

   e.g. to trigger rand request:
      1. navigate to `core/scripts/vrfv2plus/testnet` or `core/scripts/vrfv2/testnet` folder
      2. set needed env variables
         ```
         export ETH_URL=
         export ETH_CHAIN_ID=
         export ACCOUNT_KEY=
         ```
      3. Trigger rand request (get this command from the console after running `setup-envs` script )  
         ```bash
         go run . eoa-load-test-request-with-metrics --consumer-address=<contract address> --sub-id=1 --key-hash=<keyhash> --request-confirmations <> --requests 1 --runs 1 --cb-gas-limit 1_000_000 
         ```
      4. Then to check that rand request was fulfilled (get this command from the console after running `setup-envs` script )
         ```bash
         go run . eoa-load-test-read-metrics --consumer-address=<contract address> 
         ```