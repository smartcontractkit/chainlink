# Provisioning a CRIB keystone cluster 

Kudos to Functions team for inspiration.

This document outlines the steps to provision a CRIB keystone cluster for testing OCR3. 

## Usage

### Deploying a Chainlink node Cluster
First, we'll need to deploy a chainlink node cluster via CRIB. You'll want to follow the instructions in the [CRIB README](../../../crib/README.md) to set up your environment and deploy the cluster.

 **NOTE**: You'll want to deploy the keystone template, not the default devspace file. To do this, execute the following when deploying the cluster:
  
```bash
DEVSPACE_CONFIG=devspace.keystone.template.yaml devspace deploy
```

For convenience, setting the TTL to be a much longer value is helpful, otherwise the testnet native tokens that you send to nodes will be lost. Re-run this command every once in a while if you're still testing.

```bash
devspace run ttl ${namespace} 7d 
```

### Provision On-Chain Resources 
Next, we'll provision on-chain resources, namely:
1. Deploy the forwarder contract
2. Deploy OCR3 config contract
3. Setting the configuration for the OCR3 contract
4. Funding transmitters with native tokens

#### Pre-requisites
##### Blockchain Node
An HTTP URL to a blockchain node, such as a Geth node. This should be the same blockchain node that you used to deploy the chainlink node cluster.

##### Private Key
A private key to a testing wallet to use for deployment and funding. This wallet should have some native token on the chain you're deploying to. For Sepolia, around 2 ETH should be sufficient. 

The easiest way to set this up is to download [Metamask](https://metamask.io/) and create a new wallet. Once you have created a wallet, you can export the private key by clicking on the three dots next to the wallet name, selecting "Account Details", and then "Show Private Key".

#### Script Execution

Take a look at `01_deploy_contracts-sample.sh` for an example of how to do this. 

Make a copy of `01_deploy_contracts-sample.sh` and rename it to `01_deploy_contracts.sh`. Then, replace the placeholders with your own values. The default value for `ocr_config.json` is sufficient for testing.

When the on-chain resources are deployed, a json file within `artefacts` will be generated. This file will contain the addresses of the forwarder contract, the OCR3 config contract, and the block number at which the configuration was set. Be careful about deleting this file, as if you lose it, you will need to redeploy the contracts and run through all proceeding steps.

### Deploy Job Specs
The next step is to deploy the OCR3 job specs to the chainlink node cluster. This will create a bootstrapping job for the first node of the cluster (via alphabetical order) and an OCR job for each other node in the cluster.

Make a copy of `02_deploy_jobspecs-sample.sh` and rename it to `02_deploy_jobspecs.sh`. Then, replace the placeholders with your own values. The default p2p port is sufficient for testing.

### Generate CRIB Configuration

While we already have the chainlink node cluster deployed, we need to update the TOML configuration for each node to include fields for the forwarder contract along with each node's `from` address. 

Make a copy of `03_gen_crib-sample.sh` and rename it to `03_gen_crib.sh`. Then, replace the placeholders with your own values.

This will create a generated file within the `crib` directory. The script will output the command you'll need to run to deploy the configuration to the cluster. 

## Future Work
- Keystone workflow deployment
- Multi-DON support
 
