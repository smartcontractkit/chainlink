To run the solidity-contract tests run `yarn workspace chainlinkv0.5 setup` from the root of the repo, then run `yarn test` from the evm directory. You can also run `yarn test <test-file>` to run tests for a specific contract.

The actual test script run on CI is in ../tools/ci/truffle_test.

To update the Slither detection database, run `yarn build && slither --triage-mode .`.
