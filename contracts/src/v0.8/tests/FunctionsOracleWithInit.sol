// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "../dev/functions/FunctionsOracle.sol";

contract FunctionsOracleWithInit is FunctionsOracle_v0 {
  constructor() {
    initialize();
  }
}
