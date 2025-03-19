# Foundry Guide

We lock Foundry to a specific version in the `GNUmakefile`.
To ensure you have the correct local version run `make foundry`.
When you see formatting or gas differences between local and CI, it often means there is a version mismatch.
We use a locked version to avoid formatting or gas changes that suddenly pop up in CI when breaking changes are pushed to the nightly Foundry feed.


## How to start a new Foundry project

There are several files to modify when starting a new Solidity project.
Everything starts with a foundry profile in `contracts/foundry.toml`,
this profile name will be used in most of the other files.
We will assume the profile is called `newproject`.

The foundry profile should look similar to this.

```toml
[profile.newproject]
solc_version = '0.8.24'
src = 'src/v0.8/newproject'
test = 'src/v0.8/newproject/test'
optimizer_runs = 1_000_000
evm_version = 'paris'
```

After that, we have to enable CI by editing the following files.

- `.github/CODEOWNERS`
  - Add `newproject` in three places.
    - `/contracts/**/*newproject* @smartcontractkit/newproject`
    - `/contracts/src/v0.8/*newproject* @smartcontractkit/newproject`
    - `/core/gethwrappers/*newproject* @smartcontractkit/newproject`
  - Look at the file layout for the correct positions for each of these lines. Please keep the ordering alphabetical. 
- `.github/workflows/solidity-foundry.yml`
  - Add `newproject` to the `Define test matrix` section.
    - Set the min coverage >=98%.
    - Enable run-gas-snapshot.
    - Enable run-forge-fmt.
  - Add `newproject` to the `Checkout the repo` step.
- `.github/workflows/solidity-hardhat.yml`
  - Add `newproject` to the ignored list to avoid hardhat CI running for `newproject` changes.
- `contracts/GNUmakefile`
  - Add `newproject` to the ALL_FOUNDRY_PRODUCTS list in alphabetical order.
- `contracts/.prettierignore`
  - Add `src/v0.8/newproject/**` .

To enable geth wrapper generation, you will also have to create the following files.

- `contracts/scripts`
  - Create a file called `native_solc_compile_all_newproject`.
    - See example below.
- `core/gethwrappers`
  - Create a folder `newproject`.
  - It's easiest to copy another projects folder and replace the contracts in `go_generate.go` with your own contracts.
    - `ccip` is a good version to copy.
    - Remove the contents of the `generated` folder.
    - Remove the contents of the `generated-wrapper-dependency-versions-do-not-edit.txt` file.
    - Remove the contents of the `mocks` folder.
- If you need mocks, define them in `.mockery.yaml`.

```bash
#!/usr/bin/env bash

set -e

echo " ┌──────────────────────────────────────────────┐"
echo " │       Compiling Newproject contracts...      │"
echo " └──────────────────────────────────────────────┘"

PROJECT="newproject"

CONTRACTS_DIR="$( cd "$(dirname "$0")" >/dev/null 2>&1 ; cd ../ && pwd -P )"
export FOUNDRY_PROFILE="$PROJECT"

compileContract () {
  local contract
  contract=$(basename "$1")
  echo "Compiling" "$contract"

  local command
  command="forge build $CONTRACTS_DIR/src/v0.8/$PROJECT/"$1.sol" \
       --root $CONTRACTS_DIR \
       --extra-output-files bin abi \
       -o $CONTRACTS_DIR/solc/$PROJECT/$contract"
  $command
}

compileContract newproject/SingleContract.sol
compileContract newproject/OtherContract.sol

```

You should now have a fully set-up project with CI enabled.
Create a PR that introduces this setup without adding all the project's Solidity complexity, ideally before you start.
This is important
because the people approving the PR for this CI are different people from the ones approving the Solidity code.

## Testing with Foundry

We aim for (near) 100% line coverage.
Line coverage can sometimes be misleading though, so you should also look at the branch coverage. 
The CI will only take line coverage into account, which means branch coverage spot checks are highly recommended.
Setting the line coverage requirement to ~100% means you will almost guarantee all branches are also taken.

We have a strict layout and naming convention for Foundry tests. 
This is to ensure consistency within the Chainlink codebases
and make it easier for developers to work on different projects.
If your Foundry project does not yet follow the structures described below, please consider refactoring it.
The test naming structure is especially important as CI depends on it for its snapshot generation.


### Test file layout

Each foundry project has its own folder in the appropriate Solidity version folder. Within this folder there is a `test`
folder that contains all tests. This test folder mirrors the non-test folder structure with the exception that for each 
contract to be tested, there is a folder with that contract's name. Within that folder, there is a test file for each
function that is tested and optionally a setup file which is shared between the function tests. Each file has a single
contract with the name `<Contract>_<function>` e.g. `contract OnRamp_getFee is OnRampSetup`. 

Consider the following folder structure.
```
├── Router.sol
├── FeeQuoter.sol
├── onRamp
│   ├── OnRamp.sol
│   └── AnotherOnRamp.sol
```

The folder including tests would look like this.

```
├── Router.sol
├── FeeQuoter.sol
├── onRamp
│   ├── OnRamp.sol
│   └── AnotherOnRamp.sol
├── test
│   ├── Router
│   │   ├── Router.ccipSend.t.sol
│   │   ├── Router.recoverTokens.t.sol
│   │   ├── RouterSetup.t.sol
│   │   └── ....
│   ├── FeeQuoter
│   │   ├── FeeQuoter.onReport.t.sol
│   │   ├── FeeQuoter.updatePrices.t.sol
│   │   └── ....
│   ├── onRamp
│   │   ├── OnRamp
│   │   │   ├── OnRamp.constructor.t.sol
│   │   │   ├── OnRamp.getFee.t.sol
│   │   │   └── ....
│   │   ├── AnotherOnRamp
│   │   │   ├── AnotherOnRamp.constructor.t.sol
│   │   │   ├── AnotherOnRamp.getFee.t.sol
│   │   │   └── ....
```

### Test naming

Tests are named according to the following format:

```
test_FunctionName_Description for standard tests.
test_FunctionName_RevertWhen_Condition for tests expecting a revert.
testFuzz_FunctionName_Description for fuzz tests.
testFork_FunctionName_Description for tests that fork from a network.
```

Including the function name first will group tests for the same function together in the gas snapshot. Using this format
will automatically exclude fuzz, fork and reverting tests from the gas snapshot. This leads to a less flaky snapshot
with fewer merge conflicts.

Examples of correct test names for a function `getFee` are:

```
test_getFee - the base success case
test_getFee_MultipleFeeTokens - another success case with a specific scenario
test_getFee_RevertWhen_CursedByRMN - getFee reverts when it's cursed by the RMN. The error name should be used as condition when there is a single tests that checks for it
testFuzz_getFee_OnlyFeeTokens - a fuzz test that asserts that only fee tokens are used
testFork_getFee_UniswapV3MainnetFee - a fork test that uses Uniswap V3 on mainnet to get the fee
```


### What to test

Foundry unit tests should cover at least the following

- The happy path
- All emitted events.
  - Use `vm.expectEmit()`.
  - Since all state updates should emit an event, this implicitly means we test all state updates.
- All revert reasons.
  - Use `vm.expectRevert(...)`.

Consider if a fuzz test makes sense.
It often doesn't, but when it does, it can be very powerful. 
Fork tests can be considered when the code relies on existing contracts or their state. 
Focus on unit tests before exploring more advanced testing.

## Best practices

Check out the official [Foundry best practices section](https://book.getfoundry.sh/tutorials/best-practices).

- There should be no code between `vm.expectEmit`/`vm.expectRevert` and the function call.
  - Test setup should be done before the `vm.expect` call.
- Set the block number and timestamp in `foundry.toml`.
  - It is preferred to set these values to some reasonable value close to reality.
  - There are already globally applied values in the `foundry.toml` file in this repo.
- Reference errors and events from the original contracts, do not duplicate them.
- Prefer `makeAddr("string describing the contract");` over `address(12345);`.
- Pin the fork test block number to cache the results of the RPC.
- If you see something being done in existing code, that doesn't mean it is automatically correct. 
  - This document will evolve over time, and often it won't make sense to go back and refactor an entire codebase when our preferences change.


## Tips and tricks

- Use `make snapshot` to generate the correct snapshot for the selected Foundry profile.
  - Use `make snapshot-diff` to see the diff between the local snapshot and your latest changes.
- use `make wrappers` to generate the gethwrappers for the selected Foundry profile.
- Use `vm.recordLogs();` to record all logs emitted
  - Use `vm.getRecordedLogs()` to get the logs emitted. 
  - This way you can assert that a log was *not* emitted.
- Run `forge coverage --report lcov` to output code coverage
  - This can be rendered as inline coverage using e.g. Coverage Gutters for VSCode
- You can provide inline config for fuzz/invariant tests
- You can find the function selectors for a given function or error using `cast sig <FUNC_SIG>`
  - Run `forge selectors list` to see the entire list of selectors split by the contract name.

