// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

contract IgnoreContractSize {
  // test contracts are excluded from forge build --sizes by default
  // --sizes exits with code 1 if any contract is over limit, which fails CI
  // for helper contracts that are not explicit test contracts
  // use this flag to exclude from --sizes
  bool public IS_SCRIPT = true;
}
