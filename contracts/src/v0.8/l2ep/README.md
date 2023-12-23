# Overview

This folder contains the source code and tests for the Layer 2
Exchange Protocol (L2EP) contracts. It is organized as follows:

```text
.
├─/dev (stores the latest source code for L2EP)
├─/test (stores the Foundry tests for L2EP)
```

## The `/dev` Folder

The `/dev` folder contains subfolders for each chain that
has an L2EP solution implemented for it (e.g. `/scroll`, `/arbitrum`,
`/optimism`). It also contains a subfolder named `/interfaces`,
which stores shared interface types between all the supported
contracts. The top-level contracts (e.g. `CrossDomainOwnable.sol`)
serve as either abstract or parent contracts that are meant
to be reused for each indiviudal chain.

## The `/test` Folder

This folder is arranged as follows:

- `/mocks`: used for both Foundry test cases and Hardhat test cases (NOTE:
Hardhat test cases should be considered deprecated at this point)

- `/[version]`: test cases for a specific version of the L2EP contracts

### Testing Conventions and Methodology

By convention, each testing file should end in `.t.sol` (this has no effect on
Foundry, but it is a standard that other projects have also adopted). Each
testing file in this folder follows a similar structure.

```sh
L2EPTest Contract
|
|--- TestFile1.t.sol
      |
      |--- Base Contract
      |     |
      |     |--- Child Contract 1
      |     |--- Child Contract 2
      |     ...
      |
      ...
```

All test files contain a base contract defined at the top of the file. This
base contract inherits from a contract called `L2EPTest`. The `L2EPTest`
contract and base contracts have no test cases. Instead, the `L2EPTest`
contract is meant to store data that will be reused among all the base
contracts. Similarly, the base contract is meant to store data that will
be reused among any contracts that inherit it. As such, each test file
will define separate child contracts, and each will inherit from the base
contract + define its own set of tests.

The base contract defines a `setUp` function which is automatically called
exactly once before ***each*** of the tests are run in an inheriting contract.
The `setUp` function typically deploys a fresh set of test contracts so that
tests can run independently of each other. Alongside the `setUp` function,
the base contract can also define variables, constants, events, etc. that
are meant to be reused per test.

The name of the base contract follows the following convention:

```text
<NameOfContractBeingTested>Test
```

Each test in the inheriting contract has a name that follows the pattern:

```text
test_<NameOfTest>
```

The child contract names do not follow a naming convention, but they should
broadly describe the tests defined within it.

### Running Foundry Tests

#### Usage

First make sure you are in the contracts directory:

```sh
# Assuming you are currently in the /chainlink directory
cd ./contracts
```

If you already have foundry installed, you can use the following command
to run all L2EP tests:

```sh
FOUNDRY_PROFILE=l2ep forge test -vvv
```

To run a specific L2EP test, you can use a variation of the following command:

```sh
FOUNDRY_PROFILE=l2ep forge test -vvv --match-path ./src/v0.8/l2ep/test/v1_0_0/scroll/ScrollSequencerUptimeFeed.t.sol
```

A full list of flags for `forge test` can be found [here](https://book.getfoundry.sh/reference/forge/forge-test).

#### Coverage

First ensure that the correct files are being evaluated. For example, if only
v1 contracts are, being evaluated then temporarily change the L2EP profile in
`./foundry.toml`.

```sh
forge coverage --ir-minimum
```
