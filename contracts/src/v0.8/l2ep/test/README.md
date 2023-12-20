# Overview

This folder stores tests for the Layer 2 Exchange Protocol (L2EP) contracts.
It is organized as follows:

- `/mocks` (deprecated): used for Hardhat test cases

- `/[version]`: test cases for a specific version of the L2EP contracts

## Testing Methodology

TODO

## Running Foundry Tests

### Setup

Assuming you already have foundry installed, the only prerequisite to running
tests is setting the foundry profile to L2EP:

```sh
export FOUNDRY_PROFILE=l2ep
```

### Usage

Use the following command to run all tests:

```sh
forge test -vvv
```

Use the following command to run a specific test:

```sh
forge test -vvv --mp ./path/to/foundry/test/file.t.sol 
```

### Coverage

First ensure that the correct files are being evaluated. For example, if only
v1 contracts are, being evaluated then temporarily change the L2EP profile in
`./foundry.toml`.

```sh
forge coverage --ir-minimum
```
