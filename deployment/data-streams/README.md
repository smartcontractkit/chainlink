# Data Streams Deployments and Configurations

This module provides a set of workflows for deploying and configuring **Data Streams** contracts. These workflows streamline the process of setting up and maintaining the necessary on-chain components, ensuring consistent and reliable management across different environments.

## Overview

- **Contracts Location**: The primary contracts for Data Streams reside under [contracts/src/v0.8/llo-feeds](../../contracts/src/v0.8/llo-feeds).

## Interaction with Job Distributor

In addition to managing contract deployments, this module contains **changesets** that facilitate and standardize interaction with the **Job Distributor**. These changesets ensure:

- Consistent integration patterns between Data Streams contracts and the Job Distributor.
- Streamlined updates and rollouts across multiple contracts and environments.
- Improved maintainability and clarity in how Data Streams and the Job Distributor communicate and synchronize.


## Development
Follow the core [README](../../README.md) for steps on setting up development environment to run tests locally.

### Testing
#### Generating Mocks
Contract mocks should be used sparingly but can be helpful if testing complexity is high. 

1. Add the contract that satisfies the interface under [test/mocks](https://github.com/smartcontractkit/chainlink/tree/develop/contracts/src/v0.8/llo-feeds/v0.5.0/test/mocks) directory.
2. Register the contract to be compiled in the [compile script](https://github.com/smartcontractkit/chainlink/blob/develop/contracts/scripts/native_solc_compile_all_llo-feeds)
3. Run `make wrappers` specifying llo foundry profile
```bash
make FOUNDRY_PROFILE=llo-feeds wrappers
```
 

