# Provisioning a CRIB keystone cluster

Kudos to Functions team for inspiration.

This document outlines the steps to provision a CRIB keystone cluster for testing OCR3.

## Pre-requisites

### Blockchain Node

An HTTP URL to a blockchain node, such as a Geth node. This should be the same blockchain node that you used to deploy the chainlink node cluster.

### Private Key

A private key to a testing wallet to use for deployment and funding. This wallet should have some native token on the chain you're deploying to. For Sepolia, around 2 ETH should be sufficient.

The easiest way to set this up is to download [Metamask](https://metamask.io/) and create a new wallet. Once you have created a wallet, you can export the private key by clicking on the three dots next to the wallet name, selecting "Account Details", and then "Show Private Key".

## Usage

### Your first deployment

Using devspace, we can deploy a cluster and provision it via the `keystone` devspace profile. You'll want to follow the instructions in the [CRIB README](../../../crib/README.md) to set up your environment and deploy the cluster.

**NOTE**: You'll want to deploy using the `keystone` profile, not the default profile file.

```bash
# From /crib
devspace deploy --profile keystone
```

For convenience, setting the TTL to be a much longer value is helpful, otherwise the testnet native tokens that you send to nodes will be lost. You can set this in your crib `.env` file, or interactively via:

```bash
# From /crib
devspace run ttl ${namespace} 7d
```

Everytime the interactive command is run, the TTL is reset.

### Iterate
Let's say you made some changes to the codebase, and you want to see that reflected within the cluster. Simply redeploy via:
```bash
devspace deploy --profile keystone
```

### Restarting from a fresh slate

If you want to redeploy all resources, then you'll want to do the following:

```bash
# From /crib
devspace purge --profile keystone # Remove all k8s resources
./cribbit.sh crib-<new-namespace> # Purge currently leaves some hanging resources, make a new namespace
devspace deploy --profile keysone --clean # Wipe any keystone related persisted data, like artefacts and caches.
```

## What does Provisioning a CRIB keystone cluster do?

### Provision On-Chain Resources

This will provision on-chain resources, namely:

1. Deploy the forwarder contract
2. Deploy OCR3 config contract
3. Setting the configuration for the OCR3 contract
4. Funding transmitters with native tokens

When the on-chain resources are deployed, a json file within `artefacts` will be generated. This file will contain the addresses of the forwarder contract, the OCR3 config contract, and the block number at which the configuration was set. Be careful about deleting this file, as if you lose it, you will need to redeploy the contracts and run through all proceeding steps.

### Job Spec Deployment

The next step is to deploy the OCR3 job specs to the chainlink node cluster. This will create a bootstrapping job for the first node of the cluster (determined via alphabetical order) and an OCR job for each other node in the cluster.

### Update Per-Node TOML Configuration

While we already have the chainlink node cluster deployed, we need to update the TOML configuration for each node to configure the `ChainWriter`.
After updated TOML configuration overrides are generated per node, the cluster is redeployed such that the updates that effect without wiping the databases.

## Future Work

### Keystone workflow deployment
Workflow style job spec deployments are not currently support, but it should be a minor modification to the existing OCR job spec deployment logic

### Multi-DON support
Multiple DONs are not currently supported
- the devspace profile will need to be expanded so that we have multiple deployments, one per DON.
- network policy / open ports will likely have to be adjusted in the chart

### Smarter jobspec deployment
Currently, job specs deployment logic is dumb. The scripts don't check if the jobspec to deploy already exists. If you need to redeploy a job spec that has the same name as a currently uploaded one, you'll want to delete the existing job specs via `./04_delete_ocr3_jobs.sh`.
