# Data Streams Deployments and Configurations

This module provides a set of workflows for deploying and configuring **Data Streams** contracts. These workflows streamline the process of setting up and maintaining the necessary on-chain components, ensuring consistent and reliable management across different environments.

## Overview

- **Contracts Location**: The primary contracts for Data Streams reside under [contracts/src/v0.8/llo-feeds](../../contracts/src/v0.8/llo-feeds).

## Interaction with Job Distributor

In addition to managing contract deployments, this module contains **changesets** that facilitate and standardize interaction with the **Job Distributor**. These changesets ensure:

- Consistent integration patterns between Data Streams contracts and the Job Distributor.
- Streamlined updates and rollouts across multiple contracts and environments.
- Improved maintainability and clarity in how Data Streams and the Job Distributor communicate and synchronize.
